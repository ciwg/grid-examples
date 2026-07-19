package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/awareness"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/cas"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/crdt"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/identity"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/metadata"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocols"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/publish"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/store"
)

type Meta struct {
	LocalID       string `json:"local_id"`
	DocumentPCID  string `json:"document_pcid"`
	AwarenessPCID string `json:"awareness_pcid"`
	MetadataPCID  string `json:"metadata_pcid"`
	PublishPCID   string `json:"publish_pcid"`
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
	dataRoot          string
	identity          *identity.Identity
	remoteAccessToken string
	log               *store.Log
	cas               *cas.Store
	documentPCID      cid.Cid
	awarenessPCID     cid.Cid
	metadataPCID      cid.Cid
	publishPCID       cid.Cid

	mu                 sync.Mutex
	maxLamport         uint64
	seen               map[string]struct{}
	presence           map[string]awareness.Index
	syncFeeds          map[string][]crdt.SyncRecord
	metadata           map[string]metadata.Record
	published          map[string][]publish.Record
	publishByCID       map[string]publish.Record
	onSyncChanged      func(string)
	onAwarenessChanged func(string)
}

func NewApp(dataRoot string, options ...AppOptions) (*App, error) {
	var option AppOptions
	if len(options) > 0 {
		option = options[0]
	}
	documentPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.LiveDocumentSpec))
	if err != nil {
		return nil, fmt.Errorf("document pCID: %w", err)
	}
	awarenessPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.LiveAwarenessSpec))
	if err != nil {
		return nil, fmt.Errorf("awareness pCID: %w", err)
	}
	metadataPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.DocumentMetadataSpec))
	if err != nil {
		return nil, fmt.Errorf("metadata pCID: %w", err)
	}
	publishPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.PublishDocumentSpec))
	if err != nil {
		return nil, fmt.Errorf("publish pCID: %w", err)
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
		dataRoot:          dataRoot,
		identity:          identityValue,
		remoteAccessToken: option.RemoteAccessToken,
		log:               logValue,
		cas:               casValue,
		documentPCID:      documentPCID,
		awarenessPCID:     awarenessPCID,
		metadataPCID:      metadataPCID,
		publishPCID:       publishPCID,
		seen:              map[string]struct{}{},
		presence:          map[string]awareness.Index{},
		syncFeeds:         map[string][]crdt.SyncRecord{},
		metadata:          map[string]metadata.Record{},
		published:         map[string][]publish.Record{},
		publishByCID:      map[string]publish.Record{},
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
		MetadataPCID:  app.metadataPCID.String(),
		PublishPCID:   app.publishPCID.String(),
		DataRoot:      app.dataRoot,
	}
}

func (app *App) SetLiveChangeHooks(onSyncChanged func(string), onAwarenessChanged func(string)) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.onSyncChanged = onSyncChanged
	app.onAwarenessChanged = onAwarenessChanged
}

// Intent: Keep the relay non-canonical by signing and relaying exact Automerge
// change bytes while leaving replica ownership in the browser or sidecar.
// Source: DI-ramuv; DI-lumek; DI-larok
func (app *App) PostSync(documentID string, participantID string, recipientID string, messageBase64 string, embodiment string) (crdt.SyncRecord, error) {
	if err := validateDocumentID(documentID); err != nil {
		return crdt.SyncRecord{}, err
	}
	if err := validateParticipantID(participantID); err != nil {
		return crdt.SyncRecord{}, err
	}
	if err := validateRecipientID(recipientID); err != nil {
		return crdt.SyncRecord{}, err
	}
	if err := validateEmbodiment(embodiment); err != nil {
		return crdt.SyncRecord{}, err
	}
	messageBytes, err := base64.StdEncoding.DecodeString(messageBase64)
	if err != nil {
		return crdt.SyncRecord{}, fmt.Errorf("decode change bytes: %w", err)
	}
	if err := validateChangeBytes(messageBytes); err != nil {
		return crdt.SyncRecord{}, err
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
	if err := validateDocumentID(documentID); err != nil {
		return err
	}
	if err := validateParticipantID(participantID); err != nil {
		return err
	}
	if err := validateCursorValue("cursor", cursor); err != nil {
		return err
	}
	if err := validateCursorValue("head", head); err != nil {
		return err
	}
	if err := validateDisplayName(displayName); err != nil {
		return err
	}
	if err := validateColor(color); err != nil {
		return err
	}
	if err := validateEmbodiment(embodiment); err != nil {
		return err
	}
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

// Intent: Keep document-management metadata as a relay-signed current-time
// state instead of leaving description, tags, archive, and favorites purely in
// browser storage. Source: DI-loruk; DI-sukip
func (app *App) UpdateMetadata(documentID string, participantID string, title string, description string, summary string, tags []string, collections []string, favorite bool, archived bool, embodiment string) (metadata.Record, error) {
	if err := validateDocumentID(documentID); err != nil {
		return metadata.Record{}, err
	}
	if err := validateParticipantID(participantID); err != nil {
		return metadata.Record{}, err
	}
	if err := validateMetadataTitle(title); err != nil {
		return metadata.Record{}, err
	}
	if err := validateMetadataDescription(description); err != nil {
		return metadata.Record{}, err
	}
	if err := validateMetadataSummary(summary); err != nil {
		return metadata.Record{}, err
	}
	if err := validateMetadataLabels("tags", tags); err != nil {
		return metadata.Record{}, err
	}
	if err := validateMetadataLabels("collections", collections); err != nil {
		return metadata.Record{}, err
	}
	if err := validateEmbodiment(embodiment); err != nil {
		return metadata.Record{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	app.maxLamport++
	message := metadata.Message{
		Kind:          "metadata",
		DocumentID:    documentID,
		Author:        app.identity.KeyID(),
		ParticipantID: participantID,
		Title:         title,
		Description:   description,
		Summary:       summary,
		Tags:          append([]string(nil), tags...),
		Collections:   append([]string(nil), collections...),
		Favorite:      favorite,
		Archived:      archived,
		UpdatedAt:     time.Now().Format(time.RFC3339Nano),
		Lamport:       app.maxLamport,
		Embodiment:    embodiment,
	}
	envelopeBytes, err := app.makeSignedEnvelope(app.metadataPCID, message)
	if err != nil {
		return metadata.Record{}, fmt.Errorf("sign metadata message: %w", err)
	}
	return app.ingestMetadataEnvelopeLocked(envelopeBytes, nil)
}

// Intent: Make publish a separate current-time relay action that signs a
// durable manifest referencing CAS-backed bytes, rather than overloading live
// editing state or browser-local snapshot storage. Source: DI-tavul; DI-gosaf
func (app *App) PublishDocument(documentID string, participantID string, sourceKind string, sourceVersionID string, sourceVersionName string, title string, summary string, textBytes []byte, replicaBytes []byte, embodiment string) (publish.Record, error) {
	if err := validateDocumentID(documentID); err != nil {
		return publish.Record{}, err
	}
	if err := validateParticipantID(participantID); err != nil {
		return publish.Record{}, err
	}
	if err := validatePublishSourceKind(sourceKind); err != nil {
		return publish.Record{}, err
	}
	if err := validatePublishTitle(title); err != nil {
		return publish.Record{}, err
	}
	if err := validatePublishSummary(summary); err != nil {
		return publish.Record{}, err
	}
	if err := validatePublishBytes("text", textBytes); err != nil {
		return publish.Record{}, err
	}
	if err := validatePublishBytes("replica", replicaBytes); err != nil {
		return publish.Record{}, err
	}
	if err := validateEmbodiment(embodiment); err != nil {
		return publish.Record{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	textCID, err := app.cas.Put(textBytes)
	if err != nil {
		return publish.Record{}, fmt.Errorf("persist publish text: %w", err)
	}
	replicaCID, err := app.cas.Put(replicaBytes)
	if err != nil {
		return publish.Record{}, fmt.Errorf("persist publish replica: %w", err)
	}
	app.maxLamport++
	message := publish.Message{
		Kind:              "publish",
		DocumentID:        documentID,
		Author:            app.identity.KeyID(),
		ParticipantID:     participantID,
		SourceKind:        sourceKind,
		SourceVersionID:   sourceVersionID,
		SourceVersionName: sourceVersionName,
		Title:             title,
		Summary:           summary,
		TextCID:           textCID,
		ReplicaCID:        replicaCID,
		PublishedAt:       time.Now().Format(time.RFC3339Nano),
		Lamport:           app.maxLamport,
		Embodiment:        embodiment,
	}
	envelopeBytes, err := app.makeSignedEnvelope(app.publishPCID, message)
	if err != nil {
		return publish.Record{}, fmt.Errorf("sign publish manifest: %w", err)
	}
	record, err := app.ingestPublishEnvelopeLocked(envelopeBytes, nil)
	if err != nil {
		return publish.Record{}, err
	}
	return record, nil
}

func (app *App) SyncFeed(documentID string, since uint64, limit int) SyncFeed {
	app.mu.Lock()
	defer app.mu.Unlock()
	records := app.syncFeeds[documentID]
	filtered := make([]crdt.SyncRecord, 0, clampFeedLimit(limit))
	nextOffset := app.log.NextOffset()
	for _, record := range records {
		if record.Offset >= since {
			if len(filtered) >= clampFeedLimit(limit) {
				nextOffset = record.Offset
				break
			}
			filtered = append(filtered, record)
		}
	}
	if len(filtered) > 0 && nextOffset == app.log.NextOffset() {
		nextOffset = filtered[len(filtered)-1].Offset + 1
	}
	if len(filtered) == 0 {
		nextOffset = app.log.NextOffset()
	}
	return SyncFeed{
		DocumentID: documentID,
		Messages:   filtered,
		NextOffset: nextOffset,
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

func (app *App) Published(documentID string) []publish.Record {
	app.mu.Lock()
	defer app.mu.Unlock()
	records := app.published[documentID]
	if len(records) == 0 {
		return []publish.Record{}
	}
	return append([]publish.Record(nil), records...)
}

func (app *App) Metadata(documentID string) metadata.Record {
	app.mu.Lock()
	defer app.mu.Unlock()
	return cloneMetadataRecord(app.metadata[documentID], documentID)
}

func (app *App) SearchMetadata(query string, includeArchived bool) []metadata.Record {
	app.mu.Lock()
	defer app.mu.Unlock()
	normalized := strings.ToLower(strings.TrimSpace(query))
	results := make([]metadata.Record, 0, len(app.metadata))
	for documentID, record := range app.metadata {
		if record.Archived && !includeArchived {
			continue
		}
		candidate := cloneMetadataRecord(record, documentID)
		if normalized != "" && !metadataMatches(candidate, normalized) {
			continue
		}
		results = append(results, candidate)
	}
	slices.SortFunc(results, func(left metadata.Record, right metadata.Record) int {
		if left.Favorite != right.Favorite {
			if left.Favorite {
				return -1
			}
			return 1
		}
		if left.DocumentID < right.DocumentID {
			return -1
		}
		if left.DocumentID > right.DocumentID {
			return 1
		}
		return 0
	})
	return results
}

// Intent: Resolve published exchange manifests through relay-local CAS reads so
// importers can fetch a durable handoff object and its referenced bytes
// without pretending publish/import is the same thing as live sync.
// Source: DI-tavul; DI-gosaf
func (app *App) ResolvePublished(envelopeCID string) (publish.Resolved, error) {
	app.mu.Lock()
	record, ok := app.publishByCID[envelopeCID]
	app.mu.Unlock()
	if !ok {
		return publish.Resolved{}, fmt.Errorf("published manifest not found")
	}
	textBytes, err := app.cas.Get(record.TextCID)
	if err != nil {
		return publish.Resolved{}, fmt.Errorf("resolve publish text: %w", err)
	}
	replicaBytes, err := app.cas.Get(record.ReplicaCID)
	if err != nil {
		return publish.Resolved{}, fmt.Errorf("resolve publish replica: %w", err)
	}
	return publish.Resolved{
		Record:        record,
		TextBase64:    base64.StdEncoding.EncodeToString(textBytes),
		ReplicaBase64: base64.StdEncoding.EncodeToString(replicaBytes),
	}, nil
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

func (app *App) PeerMessagesSince(offset uint64, limit int) ([]string, uint64) {
	entries := app.log.EntriesSince(offset)
	maxEntries := clampFeedLimit(limit)
	messages := make([]string, 0, maxEntries)
	next := app.log.NextOffset()
	for _, entry := range entries {
		// Intent: Keep the peer relay feed limited to live document and awareness
		// traffic so publish/import exchange stays an explicit URL-resolved path
		// instead of silently piggybacking on the live sync channel. Source:
		// DI-tavul; DI-gosaf
		if entry.PCID != app.documentPCID.String() && entry.PCID != app.awarenessPCID.String() && entry.PCID != app.metadataPCID.String() {
			continue
		}
		if len(messages) >= maxEntries {
			next = entry.Offset
			break
		}
		messages = append(messages, entry.RawBase64)
	}
	if len(messages) == 0 {
		return messages, app.log.NextOffset()
	}
	if next == app.log.NextOffset() {
		lastOffset := offset
		for _, entry := range entries {
			if entry.PCID == app.documentPCID.String() || entry.PCID == app.awarenessPCID.String() || entry.PCID == app.metadataPCID.String() {
				lastOffset = entry.Offset
			}
		}
		next = lastOffset + 1
	}
	return messages, next
}

func (app *App) IngestRawBase64(raw string) error {
	envelopeBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("decode base64 envelope: %w", err)
	}
	if len(envelopeBytes) > maxChangeBytesLen*4 {
		return fmt.Errorf("envelope bytes too large")
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
					log.Printf("grid-relay peer ingest error from %s: %v", peerURL, err)
					continue
				}
			}
			offset = next
		} else {
			log.Printf("grid-relay peer poll error from %s: %v", peerURL, err)
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
	var payload peerResponse
	// Intent: Keep the copied example's HTTP teardown explicit so verification
	// can prove close-path errors are not silently ignored. Source: DI-rokod
	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	closeErr := response.Body.Close()
	if decodeErr != nil {
		if closeErr != nil {
			return nil, offset, fmt.Errorf("decode peer response: %w (close body: %v)", decodeErr, closeErr)
		}
		return nil, offset, decodeErr
	}
	if closeErr != nil {
		return nil, offset, fmt.Errorf("close peer response body: %w", closeErr)
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
	var syncDocumentID string
	var awarenessDocumentID string
	switch envelope.PCID.String() {
	case app.documentPCID.String():
		var message crdt.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return crdt.SyncRecord{}, fmt.Errorf("decode document payload: %w", err)
		}
		// Intent: Bind relay attribution to the actual signing key so peers cannot
		// sign with one key while claiming another payload author identity.
		// Source: DI-rabod
		if message.Author != envelope.Proof.KeyID {
			return crdt.SyncRecord{}, fmt.Errorf("document author %q does not match proof key", message.Author)
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
		syncDocumentID = message.DocumentID
		if message.Lamport > app.maxLamport {
			app.maxLamport = message.Lamport
		}
	case app.publishPCID.String():
		_, err = app.ingestPublishEnvelopeLocked(envelopeBytes, &entry)
		if err != nil {
			return crdt.SyncRecord{}, err
		}
		return crdt.SyncRecord{}, nil
	case app.metadataPCID.String():
		_, err = app.ingestMetadataEnvelopeLocked(envelopeBytes, &entry)
		if err != nil {
			return crdt.SyncRecord{}, err
		}
		return crdt.SyncRecord{}, nil
	case app.awarenessPCID.String():
		var message awareness.Message
		if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
			return crdt.SyncRecord{}, fmt.Errorf("decode awareness payload: %w", err)
		}
		// Intent: Apply the same signer-to-author binding to awareness as document
		// messages so cursor/presence attribution cannot drift from the proof key.
		// Source: DI-rabod
		if message.Author != envelope.Proof.KeyID {
			return crdt.SyncRecord{}, fmt.Errorf("awareness author %q does not match proof key", message.Author)
		}
		observedAt, err := time.Parse(time.RFC3339Nano, entry.ReceivedAt)
		if err != nil {
			return crdt.SyncRecord{}, fmt.Errorf("parse awareness received_at: %w", err)
		}
		index, _ := awareness.Apply(app.presence[message.DocumentID], message, envelopeCIDString, observedAt)
		app.presence[message.DocumentID] = index
		awarenessDocumentID = message.DocumentID
		if message.Lamport > app.maxLamport {
			app.maxLamport = message.Lamport
		}
	default:
		return crdt.SyncRecord{}, fmt.Errorf("unknown pCID %s", envelope.PCID)
	}
	app.seen[envelopeCIDString] = struct{}{}
	if syncDocumentID != "" && app.onSyncChanged != nil {
		app.onSyncChanged(syncDocumentID)
	}
	if awarenessDocumentID != "" && app.onAwarenessChanged != nil {
		app.onAwarenessChanged(awarenessDocumentID)
	}
	return record, nil
}

func (app *App) ingestPublishEnvelopeLocked(envelopeBytes []byte, existing *store.Entry) (publish.Record, error) {
	envelopeCID, err := protocol.CIDForBytes(envelopeBytes)
	if err != nil {
		return publish.Record{}, fmt.Errorf("publish envelope cid: %w", err)
	}
	envelopeCIDString := envelopeCID.String()
	if _, ok := app.seen[envelopeCIDString]; ok {
		return app.publishByCID[envelopeCIDString], nil
	}
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return publish.Record{}, fmt.Errorf("parse publish envelope: %w", err)
	}
	signable, err := envelope.SignableBytes()
	if err != nil {
		return publish.Record{}, fmt.Errorf("build publish signable bytes: %w", err)
	}
	if err := identity.VerifyProof(signable, envelope.Proof); err != nil {
		return publish.Record{}, fmt.Errorf("verify publish proof: %w", err)
	}
	address, err := app.cas.Put(envelopeBytes)
	if err != nil {
		return publish.Record{}, fmt.Errorf("persist publish envelope: %w", err)
	}
	if address != envelopeCIDString {
		return publish.Record{}, fmt.Errorf("publish cas address mismatch: got %s want %s", address, envelopeCIDString)
	}
	entry := store.Entry{}
	if existing != nil {
		entry = *existing
	} else {
		entry, err = app.log.Append(envelopeBytes, envelope.PCID.String())
		if err != nil {
			return publish.Record{}, fmt.Errorf("append publish log: %w", err)
		}
	}
	var message publish.Message
	if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
		return publish.Record{}, fmt.Errorf("decode publish payload: %w", err)
	}
	if message.Author != envelope.Proof.KeyID {
		return publish.Record{}, fmt.Errorf("publish author %q does not match proof key", message.Author)
	}
	// Intent: Bind durable publish attribution to the verified signing key so a
	// relay-signed exchange manifest cannot claim a different author than the
	// key that actually signed it. Source: DI-tavul; DI-gosaf
	record := publish.Record{
		Offset:            entry.Offset,
		EnvelopeCID:       envelopeCIDString,
		DocumentID:        message.DocumentID,
		Author:            message.Author,
		ParticipantID:     message.ParticipantID,
		SourceKind:        message.SourceKind,
		SourceVersionID:   message.SourceVersionID,
		SourceVersionName: message.SourceVersionName,
		Title:             message.Title,
		Summary:           message.Summary,
		TextCID:           message.TextCID,
		ReplicaCID:        message.ReplicaCID,
		PublishedAt:       message.PublishedAt,
		Embodiment:        message.Embodiment,
		ReceivedAt:        entry.ReceivedAt,
	}
	app.published[message.DocumentID] = append(app.published[message.DocumentID], record)
	app.publishByCID[record.EnvelopeCID] = record
	if message.Lamport > app.maxLamport {
		app.maxLamport = message.Lamport
	}
	app.seen[envelopeCIDString] = struct{}{}
	return record, nil
}

func (app *App) ingestMetadataEnvelopeLocked(envelopeBytes []byte, existing *store.Entry) (metadata.Record, error) {
	envelopeCID, err := protocol.CIDForBytes(envelopeBytes)
	if err != nil {
		return metadata.Record{}, fmt.Errorf("metadata envelope cid: %w", err)
	}
	envelopeCIDString := envelopeCID.String()
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return metadata.Record{}, fmt.Errorf("parse metadata envelope: %w", err)
	}
	var message metadata.Message
	if err := protocol.Unmarshal(envelope.PayloadBytes, &message); err != nil {
		return metadata.Record{}, fmt.Errorf("decode metadata payload: %w", err)
	}
	if _, ok := app.seen[envelopeCIDString]; ok {
		return cloneMetadataRecord(app.metadata[message.DocumentID], message.DocumentID), nil
	}
	signable, err := envelope.SignableBytes()
	if err != nil {
		return metadata.Record{}, fmt.Errorf("build metadata signable bytes: %w", err)
	}
	if err := identity.VerifyProof(signable, envelope.Proof); err != nil {
		return metadata.Record{}, fmt.Errorf("verify metadata proof: %w", err)
	}
	address, err := app.cas.Put(envelopeBytes)
	if err != nil {
		return metadata.Record{}, fmt.Errorf("persist metadata envelope: %w", err)
	}
	if address != envelopeCIDString {
		return metadata.Record{}, fmt.Errorf("metadata cas address mismatch: got %s want %s", address, envelopeCIDString)
	}
	entry := store.Entry{}
	if existing != nil {
		entry = *existing
	} else {
		entry, err = app.log.Append(envelopeBytes, envelope.PCID.String())
		if err != nil {
			return metadata.Record{}, fmt.Errorf("append metadata log: %w", err)
		}
	}
	if message.Author != envelope.Proof.KeyID {
		return metadata.Record{}, fmt.Errorf("metadata author %q does not match proof key", message.Author)
	}
	// Intent: Keep cross-document labels and descriptive fields relay-backed as
	// latest-state metadata so search, favorites, archive, and collections stay
	// shareable without pretending they are part of the live CRDT text stream.
	// Source: DI-loruk; DI-sukip
	record := metadata.Record{
		Offset:        entry.Offset,
		EnvelopeCID:   envelopeCIDString,
		DocumentID:    message.DocumentID,
		Author:        message.Author,
		ParticipantID: message.ParticipantID,
		Title:         message.Title,
		Description:   message.Description,
		Summary:       message.Summary,
		Tags:          append([]string(nil), message.Tags...),
		Collections:   append([]string(nil), message.Collections...),
		Favorite:      message.Favorite,
		Archived:      message.Archived,
		UpdatedAt:     message.UpdatedAt,
		Embodiment:    message.Embodiment,
		ReceivedAt:    entry.ReceivedAt,
		Lamport:       message.Lamport,
	}
	current := app.metadata[message.DocumentID]
	if metadataWins(record, current) {
		app.metadata[message.DocumentID] = record
	}
	if message.Lamport > app.maxLamport {
		app.maxLamport = message.Lamport
	}
	app.seen[envelopeCIDString] = struct{}{}
	return record, nil
}

func metadataWins(next metadata.Record, current metadata.Record) bool {
	if current.DocumentID == "" {
		return true
	}
	if next.Lamport != current.Lamport {
		return next.Lamport > current.Lamport
	}
	if next.Author != current.Author {
		return next.Author > current.Author
	}
	return next.EnvelopeCID > current.EnvelopeCID
}

func metadataMatches(record metadata.Record, query string) bool {
	values := []string{
		record.DocumentID,
		record.Title,
		record.Description,
		record.Summary,
		strings.Join(record.Tags, " "),
		strings.Join(record.Collections, " "),
	}
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}
	return false
}

func cloneMetadataRecord(record metadata.Record, documentID string) metadata.Record {
	record.DocumentID = firstNonEmpty(record.DocumentID, documentID)
	record.Tags = append([]string(nil), record.Tags...)
	record.Collections = append([]string(nil), record.Collections...)
	return record
}

func firstNonEmpty(left string, right string) string {
	if left != "" {
		return left
	}
	return right
}
