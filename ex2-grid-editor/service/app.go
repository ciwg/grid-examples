package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/awareness"
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/cas"
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/crdt"
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/identity"
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/protocol"
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/protocols"
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/store"
)

type Meta struct {
	LocalID       string `json:"local_id"`
	DocumentPCID  string `json:"document_pcid"`
	AwarenessPCID string `json:"awareness_pcid"`
	DataRoot      string `json:"data_root"`
}

type RelayState struct {
	DocumentID   string                `json:"document_id"`
	NextOffset   uint64                `json:"next_offset"`
	MessageCount int                   `json:"message_count"`
	Awareness    []awareness.PeerState `json:"awareness"`
}

type SyncFeed struct {
	DocumentID string            `json:"document_id"`
	Messages   []crdt.SyncRecord `json:"messages"`
	NextOffset uint64            `json:"next_offset"`
}

type peerResponse struct {
	Messages   []string `json:"messages"`
	NextOffset uint64   `json:"next_offset"`
}

type App struct {
	dataRoot      string
	identity      *identity.Identity
	log           *store.Log
	cas           *cas.Store
	documentPCID  cid.Cid
	awarenessPCID cid.Cid

	mu         sync.Mutex
	maxLamport uint64
	seen       map[string]struct{}
	presence   map[string]awareness.Index
	syncFeeds  map[string][]crdt.SyncRecord
}

func NewApp(dataRoot string) (*App, error) {
	documentPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.LiveDocumentSpec))
	if err != nil {
		return nil, fmt.Errorf("document pCID: %w", err)
	}
	awarenessPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.LiveAwarenessSpec))
	if err != nil {
		return nil, fmt.Errorf("awareness pCID: %w", err)
	}
	identityPath := filepath.Join(dataRoot, "identity_ed25519_seed")
	identityValue, err := identity.LoadOrCreate(identityPath)
	if err != nil {
		return nil, fmt.Errorf("load identity: %w", err)
	}
	casValue, err := cas.Open(filepath.Join(dataRoot, "cas"))
	if err != nil {
		return nil, fmt.Errorf("open cas: %w", err)
	}
	logValue, err := store.Open(filepath.Join(dataRoot, "message-log.jsonl"))
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	app := &App{
		dataRoot:      dataRoot,
		identity:      identityValue,
		log:           logValue,
		cas:           casValue,
		documentPCID:  documentPCID,
		awarenessPCID: awarenessPCID,
		seen:          map[string]struct{}{},
		presence:      map[string]awareness.Index{},
		syncFeeds:     map[string][]crdt.SyncRecord{},
	}
	for _, entry := range logValue.All() {
		if err := app.replayEntry(entry); err != nil {
			return nil, fmt.Errorf("replay entry %d: %w", entry.Offset, err)
		}
	}
	return app, nil
}

func (app *App) Meta() Meta {
	return Meta{
		LocalID:       app.identity.KeyID(),
		DocumentPCID:  app.documentPCID.String(),
		AwarenessPCID: app.awarenessPCID.String(),
		DataRoot:      app.dataRoot,
	}
}

// Intent: Keep the relay non-canonical by signing and relaying exact Automerge
// change bytes while leaving replica ownership in the browser or sidecar.
// Source: DI-ramuv; DI-lumek; DI-larok
func (app *App) PostSync(documentID string, participantID string, recipientID string, messageBase64 string, embodiment string) (crdt.SyncRecord, error) {
	messageBytes, err := base64.StdEncoding.DecodeString(messageBase64)
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("decode change bytes: %w", err)
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	app.maxLamport++
	message := crdt.Message{
		Kind:          "change",
		DocumentID:    documentID,
		Author:        app.identity.KeyID(),
		ParticipantID: participantID,
		RecipientID:   recipientID,
		ChangeBytes:   messageBytes,
		Lamport:       app.maxLamport,
		Embodiment:    embodiment,
	}
	envelopeBytes, err := app.makeSignedEnvelope(app.documentPCID, message)
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("sign document change: %w", err)
	}
	record, err := app.ingestEnvelopeLocked(envelopeBytes, nil)
	if err != nil {
		return crdt.SyncRecord{}, err
	}
	return record, nil
}

func (app *App) UpdateAwareness(documentID string, participantID string, cursor int, head int, typing bool, displayName string, color string, embodiment string) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.maxLamport++
	message := awareness.Message{
		Kind:          "state",
		DocumentID:    documentID,
		Author:        app.identity.KeyID(),
		ParticipantID: participantID,
		DisplayName:   displayName,
		Color:         color,
		Cursor:        cursor,
		Head:          head,
		Typing:        typing,
		Lamport:       app.maxLamport,
		Embodiment:    embodiment,
	}
	envelopeBytes, err := app.makeSignedEnvelope(app.awarenessPCID, message)
	if err != nil {
		return fmt.Errorf("sign awareness message: %w", err)
	}
	_, err = app.ingestEnvelopeLocked(envelopeBytes, nil)
	return err
}

func (app *App) SyncFeed(documentID string, since uint64) SyncFeed {
	app.mu.Lock()
	defer app.mu.Unlock()
	records := app.syncFeeds[documentID]
	filtered := make([]crdt.SyncRecord, 0, len(records))
	for _, record := range records {
		if record.Offset >= since {
			filtered = append(filtered, record)
		}
	}
	return SyncFeed{
		DocumentID: documentID,
		Messages:   filtered,
		NextOffset: app.log.NextOffset(),
	}
}

func (app *App) AwarenessState(documentID string) []awareness.PeerState {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.awarenessLocked(documentID)
}

func (app *App) State(documentID string) RelayState {
	app.mu.Lock()
	defer app.mu.Unlock()
	return RelayState{
		DocumentID:   documentID,
		NextOffset:   app.log.NextOffset(),
		MessageCount: len(app.syncFeeds[documentID]),
		Awareness:    app.awarenessLocked(documentID),
	}
}

func (app *App) awarenessLocked(documentID string) []awareness.PeerState {
	presenceIndex := app.presence[documentID]
	peers := make([]awareness.PeerState, 0, len(presenceIndex))
	for _, peer := range presenceIndex {
		peers = append(peers, peer)
	}
	slices.SortFunc(peers, func(left awareness.PeerState, right awareness.PeerState) int {
		if left.ParticipantID < right.ParticipantID {
			return -1
		}
		if left.ParticipantID > right.ParticipantID {
			return 1
		}
		return 0
	})
	return peers
}

func (app *App) PeerMessagesSince(offset uint64) ([]string, uint64) {
	entries := app.log.EntriesSince(offset)
	messages := make([]string, 0, len(entries))
	for _, entry := range entries {
		messages = append(messages, entry.RawBase64)
	}
	return messages, app.log.NextOffset()
}

func (app *App) IngestRawBase64(raw string) error {
	envelopeBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("decode base64 envelope: %w", err)
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	_, err = app.ingestEnvelopeLocked(envelopeBytes, nil)
	return err
}

func (app *App) StartPeerPolling(ctx context.Context, peerURLs []string, interval time.Duration) {
	for _, peerURL := range peerURLs {
		peerURL := peerURL
		go app.pollPeer(ctx, peerURL, interval)
	}
}

// Intent: Keep the first inter-host relay path inspectable by polling peer
// signed-message feeds instead of hiding replication in the browser code.
// Source: DI-ramuv; DI-lumek
func (app *App) pollPeer(ctx context.Context, peerURL string, interval time.Duration) {
	client := &http.Client{Timeout: 4 * time.Second}
	offset := uint64(0)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		messages, next, err := app.fetchPeerMessages(client, peerURL, offset)
		if err == nil {
			for _, raw := range messages {
				if err := app.IngestRawBase64(raw); err != nil {
					continue
				}
			}
			offset = next
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (app *App) fetchPeerMessages(client *http.Client, peerURL string, offset uint64) ([]string, uint64, error) {
	response, err := client.Get(fmt.Sprintf("%s/api/peer/messages?since=%d", peerURL, offset))
	if err != nil {
		return nil, offset, err
	}
	defer response.Body.Close()
	var payload peerResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, offset, err
	}
	return payload.Messages, payload.NextOffset, nil
}

func (app *App) replayEntry(entry store.Entry) error {
	envelopeBytes, err := base64.StdEncoding.DecodeString(entry.RawBase64)
	if err != nil {
		return fmt.Errorf("decode replay envelope: %w", err)
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	_, err = app.ingestEnvelopeLocked(envelopeBytes, &entry)
	return err
}

func (app *App) makeSignedEnvelope(pcid cid.Cid, payload any) ([]byte, error) {
	payloadBytes, err := protocol.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	envelope := protocol.NewEnvelope(pcid, payloadBytes, protocol.Proof{})
	signable, err := envelope.SignableBytes()
	if err != nil {
		return nil, fmt.Errorf("signable bytes: %w", err)
	}
	proof, err := app.identity.SignProof(signable)
	if err != nil {
		return nil, fmt.Errorf("sign proof: %w", err)
	}
	envelope.Proof = proof
	return envelope.Bytes()
}

// Intent: Rebuild relay awareness and CRDT feed indices only from exact signed
// envelope bytes so CAS, verification, replay, and peer polling all share the
// same object identity path. Source: DI-ramuv; DI-zegov; DI-larok
func (app *App) ingestEnvelopeLocked(envelopeBytes []byte, existing *store.Entry) (crdt.SyncRecord, error) {
	envelopeCID, err := protocol.CIDForBytes(envelopeBytes)
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("envelope cid: %w", err)
	}
	envelopeCIDString := envelopeCID.String()
	if _, ok := app.seen[envelopeCIDString]; ok {
		return crdt.SyncRecord{}, nil
	}
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("parse envelope: %w", err)
	}
	signable, err := envelope.SignableBytes()
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("build signable bytes: %w", err)
	}
	if err := identity.VerifyProof(signable, envelope.Proof); err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("verify proof: %w", err)
	}
	address, err := app.cas.Put(envelopeBytes)
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("persist cas object: %w", err)
	}
	if address != envelopeCIDString {
		return crdt.SyncRecord{}, fmt.Errorf("cas address mismatch: got %s want %s", address, envelopeCIDString)
	}
	entry := store.Entry{}
	if existing != nil {
		entry = *existing
	} else {
		entry, err = app.log.Append(envelopeBytes, envelope.PCID.String())
		if err != nil {
			return crdt.SyncRecord{}, fmt.Errorf("append log: %w", err)
		}
	}
	var record crdt.SyncRecord
	switch envelope.PCID.String() {
	case app.documentPCID.String():
		var message crdt.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return crdt.SyncRecord{}, fmt.Errorf("decode document payload: %w", err)
		}
		record = crdt.SyncRecord{
			Offset:        entry.Offset,
			EnvelopeCID:   envelopeCIDString,
			ParticipantID: message.ParticipantID,
			RecipientID:   message.RecipientID,
			Author:        message.Author,
			MessageBase64: base64.StdEncoding.EncodeToString(message.ChangeBytes),
			Embodiment:    message.Embodiment,
			ReceivedAt:    entry.ReceivedAt,
		}
		app.syncFeeds[message.DocumentID] = append(app.syncFeeds[message.DocumentID], record)
		if message.Lamport > app.maxLamport {
			app.maxLamport = message.Lamport
		}
	case app.awarenessPCID.String():
		var message awareness.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return crdt.SyncRecord{}, fmt.Errorf("decode awareness payload: %w", err)
		}
		index, _ := awareness.Apply(app.presence[message.DocumentID], message, envelopeCIDString)
		app.presence[message.DocumentID] = index
		if message.Lamport > app.maxLamport {
			app.maxLamport = message.Lamport
		}
	default:
		return crdt.SyncRecord{}, fmt.Errorf("unknown pCID %s", envelope.PCID)
	}
	app.seen[envelopeCIDString] = struct{}{}
	return record, nil
}
