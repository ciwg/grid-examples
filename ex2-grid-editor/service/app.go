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
	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/document"
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

type DocumentView struct {
	DocumentID string                `json:"document_id"`
	Content    string                `json:"content"`
	ContentCID string                `json:"content_cid"`
	MessageCID string                `json:"message_cid"`
	Lamport    uint64                `json:"lamport"`
	Author     string                `json:"author"`
	Awareness  []awareness.PeerState `json:"awareness"`
}

type peerResponse struct {
	Messages   []string `json:"messages"`
	NextOffset uint64   `json:"next_offset"`
}

type App struct {
	dataRoot      string
	identity      *identity.Identity
	log           *store.Log
	documentPCID  cid.Cid
	awarenessPCID cid.Cid

	mu         sync.Mutex
	maxLamport uint64
	seen       map[string]struct{}
	documents  map[string]document.State
	presence   map[string]awareness.Index
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
	logValue, err := store.Open(filepath.Join(dataRoot, "message-log.jsonl"))
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	app := &App{
		dataRoot:      dataRoot,
		identity:      identityValue,
		log:           logValue,
		documentPCID:  documentPCID,
		awarenessPCID: awarenessPCID,
		seen:          map[string]struct{}{},
		documents:     map[string]document.State{},
		presence:      map[string]awareness.Index{},
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

// Intent: The local service is the single place that translates embodiment
// actions into signed grid messages so browser and Neovim stay two surfaces of
// one app contract rather than diverging peer implementations. Source:
// DI-lodug; DI-jilin
func (app *App) ReplaceDocument(documentID string, content string, embodiment string) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	previousCID := app.documents[documentID].MessageCID
	app.maxLamport++
	message, err := document.NewMessage(documentID, content, app.maxLamport, app.identity.KeyID(), embodiment, previousCID)
	if err != nil {
		return fmt.Errorf("new document message: %w", err)
	}
	envelopeBytes, err := app.makeSignedEnvelope(app.documentPCID, message)
	if err != nil {
		return fmt.Errorf("sign document message: %w", err)
	}
	return app.ingestEnvelopeLocked(envelopeBytes, true)
}

func (app *App) UpdateAwareness(documentID string, cursor int, head int, typing bool, displayName string, color string, embodiment string) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.maxLamport++
	message := awareness.Message{
		Kind:        "state",
		DocumentID:  documentID,
		Author:      app.identity.KeyID(),
		DisplayName: displayName,
		Color:       color,
		Cursor:      cursor,
		Head:        head,
		Typing:      typing,
		Lamport:     app.maxLamport,
		Embodiment:  embodiment,
	}
	envelopeBytes, err := app.makeSignedEnvelope(app.awarenessPCID, message)
	if err != nil {
		return fmt.Errorf("sign awareness message: %w", err)
	}
	return app.ingestEnvelopeLocked(envelopeBytes, true)
}

func (app *App) DocumentState(documentID string) DocumentView {
	app.mu.Lock()
	defer app.mu.Unlock()
	documentState := app.documents[documentID]
	presenceIndex := app.presence[documentID]
	peers := make([]awareness.PeerState, 0, len(presenceIndex))
	for _, peer := range presenceIndex {
		peers = append(peers, peer)
	}
	slices.SortFunc(peers, func(left awareness.PeerState, right awareness.PeerState) int {
		if left.Author < right.Author {
			return -1
		}
		if left.Author > right.Author {
			return 1
		}
		return 0
	})
	return DocumentView{
		DocumentID: documentID,
		Content:    documentState.Content,
		ContentCID: documentState.ContentCID,
		MessageCID: documentState.MessageCID,
		Lamport:    documentState.Lamport,
		Author:     documentState.Author,
		Awareness:  peers,
	}
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
	return app.ingestEnvelopeLocked(envelopeBytes, true)
}

func (app *App) StartPeerPolling(ctx context.Context, peerURLs []string, interval time.Duration) {
	for _, peerURL := range peerURLs {
		peerURL := peerURL
		go app.pollPeer(ctx, peerURL, interval)
	}
}

// Intent: Keep the first inter-host sync path simple and inspectable by polling
// peer append-only message logs over an internal adapter boundary instead of
// hiding replication in UI code. Source: DI-lodug; DI-jilin
func (app *App) pollPeer(ctx context.Context, peerURL string, interval time.Duration) {
	client := &http.Client{Timeout: 4 * time.Second}
	offset := uint64(0)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		messages, next, err := app.fetchPeerMessages(client, peerURL, offset)
		if err == nil {
			for _, raw := range messages {
				_ = app.IngestRawBase64(raw)
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
	return app.ingestEnvelopeLocked(envelopeBytes, false)
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

// Intent: Rebuild the local document and awareness projections only from exact
// signed envelope bytes so replay, verification, and persistence all share one
// truth path. Source: DI-tofug; DI-jilin
func (app *App) ingestEnvelopeLocked(envelopeBytes []byte, persist bool) error {
	envelopeCID, err := protocol.CIDForBytes(envelopeBytes)
	if err != nil {
		return fmt.Errorf("envelope cid: %w", err)
	}
	envelopeCIDString := envelopeCID.String()
	if _, ok := app.seen[envelopeCIDString]; ok {
		return nil
	}
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return fmt.Errorf("parse envelope: %w", err)
	}
	signable, err := envelope.SignableBytes()
	if err != nil {
		return fmt.Errorf("build signable bytes: %w", err)
	}
	if err := identity.VerifyProof(signable, envelope.Proof); err != nil {
		return fmt.Errorf("verify proof: %w", err)
	}
	switch envelope.PCID.String() {
	case app.documentPCID.String():
		var message document.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return fmt.Errorf("decode document payload: %w", err)
		}
		nextState, applied, err := document.Apply(app.documents[message.DocumentID], message, envelopeCIDString)
		if err != nil {
			return err
		}
		if applied {
			app.documents[message.DocumentID] = nextState
		}
		if message.Lamport > app.maxLamport {
			app.maxLamport = message.Lamport
		}
	case app.awarenessPCID.String():
		var message awareness.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return fmt.Errorf("decode awareness payload: %w", err)
		}
		index, _ := awareness.Apply(app.presence[message.DocumentID], message, envelopeCIDString)
		app.presence[message.DocumentID] = index
		if message.Lamport > app.maxLamport {
			app.maxLamport = message.Lamport
		}
	default:
		return fmt.Errorf("unknown pCID %s", envelope.PCID)
	}
	app.seen[envelopeCIDString] = struct{}{}
	if persist {
		if _, err := app.log.Append(envelopeBytes, envelope.PCID.String()); err != nil {
			return fmt.Errorf("append log: %w", err)
		}
	}
	return nil
}
