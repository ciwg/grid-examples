package service_test

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/awareness"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/crdt"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/identity"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocol"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/protocols"
	"github.com/computerscienceiscool/grid-examples/ex3-grid-editor-websocket/service"
)

func TestPostSyncAppearsInFeedAndReplay(t *testing.T) {
	t.Parallel()
	appA, err := service.NewApp(filepath.Join(t.TempDir(), "relay-a"))
	if err != nil {
		t.Fatalf("new app a: %v", err)
	}
	record, err := appA.PostSync("demo", "browser-a", "browser-b", base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}), "browser", "", "")
	if err != nil {
		t.Fatalf("post sync: %v", err)
	}
	feed := appA.SyncFeed("demo", 0, 8)
	if len(feed.Messages) != 1 {
		t.Fatalf("message count mismatch: got %d", len(feed.Messages))
	}
	if feed.Messages[0].EnvelopeCID != record.EnvelopeCID {
		t.Fatalf("envelope cid mismatch: got %s want %s", feed.Messages[0].EnvelopeCID, record.EnvelopeCID)
	}
	rawMessages, _ := appA.PeerMessagesSince(0, 8)
	if len(rawMessages) != 1 {
		t.Fatalf("peer message count mismatch: got %d", len(rawMessages))
	}

	appB, err := service.NewApp(filepath.Join(t.TempDir(), "relay-b"))
	if err != nil {
		t.Fatalf("new app b: %v", err)
	}
	if err := appB.IngestRawBase64(rawMessages[0]); err != nil {
		t.Fatalf("ingest peer message: %v", err)
	}
	feedB := appB.SyncFeed("demo", 0, 8)
	if len(feedB.Messages) != 1 {
		t.Fatalf("replayed message count mismatch: got %d", len(feedB.Messages))
	}
	if feedB.Messages[0].MessageBase64 != record.MessageBase64 {
		t.Fatalf("message bytes mismatch: got %s want %s", feedB.Messages[0].MessageBase64, record.MessageBase64)
	}
}

func TestSyncFeedPagesBoundedHistory(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	for index := 0; index < 3; index++ {
		if _, err := app.PostSync("demo", "browser-a", "", base64.StdEncoding.EncodeToString([]byte{byte(index + 1)}), "browser", "", ""); err != nil {
			t.Fatalf("post sync %d: %v", index, err)
		}
	}
	first := app.SyncFeed("demo", 0, 2)
	if len(first.Messages) != 2 {
		t.Fatalf("first page count mismatch: got %d", len(first.Messages))
	}
	if first.NextOffset <= first.Messages[1].Offset {
		t.Fatalf("expected next offset after first page, got %d", first.NextOffset)
	}
	second := app.SyncFeed("demo", first.NextOffset, 2)
	if len(second.Messages) != 1 {
		t.Fatalf("second page count mismatch: got %d", len(second.Messages))
	}
}

func TestPostSyncStoresLatestSnapshotInState(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	textBase64 := base64.StdEncoding.EncodeToString([]byte("hello"))
	replicaBase64 := base64.StdEncoding.EncodeToString([]byte{9, 8, 7, 6})
	record, err := app.PostSync("demo", "browser-a", "", base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}), "browser", textBase64, replicaBase64)
	if err != nil {
		t.Fatalf("post sync with snapshot: %v", err)
	}
	state := app.State("demo")
	if !state.SnapshotPresent {
		t.Fatalf("expected snapshot present")
	}
	if state.TextBase64 != textBase64 {
		t.Fatalf("text snapshot mismatch: got %q want %q", state.TextBase64, textBase64)
	}
	if state.ReplicaBase64 != replicaBase64 {
		t.Fatalf("replica snapshot mismatch: got %q want %q", state.ReplicaBase64, replicaBase64)
	}
	if state.SnapshotOffset != record.Offset {
		t.Fatalf("snapshot offset mismatch: got %d want %d", state.SnapshotOffset, record.Offset)
	}
}

func TestAwarenessTracksParticipantsSeparately(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if err := app.UpdateAwareness("demo", "browser-a", 4, 4, true, "Alice", "#ff0000", "browser"); err != nil {
		t.Fatalf("update awareness a: %v", err)
	}
	if err := app.UpdateAwareness("demo", "browser-b", 9, 11, false, "Bob", "#00ff00", "browser"); err != nil {
		t.Fatalf("update awareness b: %v", err)
	}
	peers := app.AwarenessState("demo")
	if len(peers) != 2 {
		t.Fatalf("peer count mismatch: got %d", len(peers))
	}
	if peers[0].ParticipantID == peers[1].ParticipantID {
		t.Fatalf("participant ids collapsed: %q", peers[0].ParticipantID)
	}
	if peers[0].LastSeenAt == "" || peers[1].LastSeenAt == "" {
		t.Fatalf("expected last seen timestamps for awareness peers")
	}
}

func TestAwarenessUsesAuthorFallbackAndLatestRelayState(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	signer := loadIdentityForTest(t, filepath.Join(t.TempDir(), "signer"))
	if err := app.IngestRawBase64(signAwarenessMessage(t, signer, awareness.Message{
		Kind:        "state",
		DocumentID:  "demo",
		Author:      signer.KeyID(),
		DisplayName: "Alice",
		Color:       "#123456",
		Cursor:      12,
		Head:        14,
		Typing:      true,
		Lamport:     1,
		Embodiment:  "nvim",
	})); err != nil {
		t.Fatalf("ingest awareness initial: %v", err)
	}
	if err := app.IngestRawBase64(signAwarenessMessage(t, signer, awareness.Message{
		Kind:        "state",
		DocumentID:  "demo",
		Author:      signer.KeyID(),
		DisplayName: "Alice stale",
		Color:       "#654321",
		Cursor:      1,
		Head:        1,
		Typing:      false,
		Lamport:     2,
		Embodiment:  "nvim",
	})); err != nil {
		t.Fatalf("ingest awareness stale: %v", err)
	}
	peers := app.AwarenessState("demo")
	if len(peers) != 1 {
		t.Fatalf("peer count mismatch: got %d", len(peers))
	}
	if peers[0].ParticipantID == "" {
		t.Fatalf("expected author fallback participant id")
	}
	if peers[0].Cursor != 1 {
		t.Fatalf("expected latest relay update to win, got cursor %d", peers[0].Cursor)
	}
	if peers[0].LastSeenAt == "" {
		t.Fatalf("expected awareness replay to keep last seen timestamp")
	}
}

func TestIngestRejectsAuthorProofMismatch(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	signer := loadIdentityForTest(t, filepath.Join(t.TempDir(), "signer"))
	err = app.IngestRawBase64(signSyncMessage(t, signer, crdt.Message{
		Kind:          "change",
		DocumentID:    "demo",
		Author:        "not-the-signer",
		ParticipantID: "browser-a",
		ChangeBytes:   []byte{1, 2, 3},
		Lamport:       1,
		Embodiment:    "browser",
	}))
	if err == nil {
		t.Fatalf("expected author/proof mismatch error")
	}
}

func TestPublishDocumentResolvesLocallyAndStaysOutOfPeerFeed(t *testing.T) {
	t.Parallel()
	appA, err := service.NewApp(filepath.Join(t.TempDir(), "relay-a"))
	if err != nil {
		t.Fatalf("new app a: %v", err)
	}
	record, err := appA.PublishDocument(
		"demo",
		"browser-a",
		"current",
		"",
		"",
		"Demo publish",
		"Shared from relay a",
		[]byte("# demo\n\nhello"),
		[]byte{9, 8, 7, 6},
		"browser",
	)
	if err != nil {
		t.Fatalf("publish document: %v", err)
	}
	published := appA.Published("demo")
	if len(published) != 1 {
		t.Fatalf("published count mismatch: got %d", len(published))
	}
	if published[0].EnvelopeCID != record.EnvelopeCID {
		t.Fatalf("publish envelope cid mismatch: got %s want %s", published[0].EnvelopeCID, record.EnvelopeCID)
	}
	resolved, err := appA.ResolvePublished(record.EnvelopeCID)
	if err != nil {
		t.Fatalf("resolve published: %v", err)
	}
	if got := decodedBase64(t, resolved.TextBase64); got != "# demo\n\nhello" {
		t.Fatalf("resolved text mismatch: %q", got)
	}

	rawMessages, _ := appA.PeerMessagesSince(0, 16)
	for _, raw := range rawMessages {
		if raw == "" {
			t.Fatalf("unexpected empty peer message")
		}
	}
	if len(rawMessages) != 0 {
		t.Fatalf("expected publish manifest to stay out of peer feed, got %d messages", len(rawMessages))
	}
}

func TestResolvePublishedMissingManifestFails(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if _, err := app.ResolvePublished("bafymissing"); err == nil {
		t.Fatalf("expected missing manifest error")
	}
}

func TestPublishDocumentRejectsInvalidInputs(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if _, err := app.PublishDocument("demo", "browser-a", "bad-kind", "", "", "Title", "", []byte("text"), []byte{1}, "browser"); err == nil {
		t.Fatalf("expected invalid source kind error")
	}
	if _, err := app.PublishDocument("demo", "browser-a", "current", "", "", "", "", []byte("text"), []byte{1}, "browser"); err == nil {
		t.Fatalf("expected missing title error")
	}
}

func TestPublishDocumentRejectsOversizedReplica(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	oversized := make([]byte, 4<<20+1)
	if _, err := app.PublishDocument("demo", "browser-a", "current", "", "", "Title", "", []byte("text"), oversized, "browser"); err == nil {
		t.Fatalf("expected oversized replica error")
	}
}

func TestMetadataRoundTripsLocally(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	record, err := app.UpdateMetadata(
		"demo",
		"browser-a",
		"Demo title",
		"Longer description",
		"Short summary",
		[]string{"docs", "grid"},
		[]string{"favorites"},
		true,
		false,
		"browser",
	)
	if err != nil {
		t.Fatalf("update metadata: %v", err)
	}
	got := app.Metadata("demo")
	if got.EnvelopeCID != record.EnvelopeCID {
		t.Fatalf("metadata envelope cid mismatch: got %s want %s", got.EnvelopeCID, record.EnvelopeCID)
	}
	if got.Title != "Demo title" || got.Description != "Longer description" || got.Summary != "Short summary" {
		t.Fatalf("metadata fields mismatch: %#v", got)
	}
	if len(got.Tags) != 2 || len(got.Collections) != 1 {
		t.Fatalf("metadata labels mismatch: %#v", got)
	}
}

func TestMetadataSearchPrefersFavoritesAndFiltersArchived(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if _, err := app.UpdateMetadata("alpha", "browser-a", "Alpha", "", "", []string{"grid"}, []string{"team"}, false, false, "browser"); err != nil {
		t.Fatalf("update alpha metadata: %v", err)
	}
	if _, err := app.UpdateMetadata("beta", "browser-b", "Beta", "", "", []string{"grid"}, []string{"team"}, true, false, "browser"); err != nil {
		t.Fatalf("update beta metadata: %v", err)
	}
	if _, err := app.UpdateMetadata("archive", "browser-c", "Archive", "", "", []string{"grid"}, []string{"team"}, false, true, "browser"); err != nil {
		t.Fatalf("update archive metadata: %v", err)
	}

	results := app.SearchMetadata("grid", false)
	if len(results) != 2 {
		t.Fatalf("metadata search count mismatch: got %d", len(results))
	}
	if results[0].DocumentID != "beta" {
		t.Fatalf("expected favorite document first, got %s", results[0].DocumentID)
	}
	for _, result := range results {
		if result.DocumentID == "archive" {
			t.Fatalf("archived document unexpectedly returned when includeArchived=false")
		}
	}

	withArchived := app.SearchMetadata("grid", true)
	if len(withArchived) != 3 {
		t.Fatalf("metadata search with archived count mismatch: got %d", len(withArchived))
	}
}

func TestMetadataAppearsInPeerFeedAndReplays(t *testing.T) {
	t.Parallel()
	appA, err := service.NewApp(filepath.Join(t.TempDir(), "relay-a"))
	if err != nil {
		t.Fatalf("new app a: %v", err)
	}
	record, err := appA.UpdateMetadata("demo", "browser-a", "Shared title", "Description", "Summary", []string{"grid"}, []string{"docs"}, true, false, "browser")
	if err != nil {
		t.Fatalf("update metadata: %v", err)
	}

	rawMessages, _ := appA.PeerMessagesSince(0, 16)
	if len(rawMessages) != 1 {
		t.Fatalf("peer message count mismatch: got %d", len(rawMessages))
	}

	appB, err := service.NewApp(filepath.Join(t.TempDir(), "relay-b"))
	if err != nil {
		t.Fatalf("new app b: %v", err)
	}
	if err := appB.IngestRawBase64(rawMessages[0]); err != nil {
		t.Fatalf("ingest metadata message: %v", err)
	}
	got := appB.Metadata("demo")
	if got.EnvelopeCID != record.EnvelopeCID {
		t.Fatalf("replayed metadata envelope cid mismatch: got %s want %s", got.EnvelopeCID, record.EnvelopeCID)
	}
	if got.Title != "Shared title" || got.Favorite != true {
		t.Fatalf("replayed metadata mismatch: %#v", got)
	}
}

func loadIdentityForTest(t *testing.T, path string) *identity.Identity {
	t.Helper()
	value, err := identity.LoadOrCreate(path)
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}
	return value
}

func signSyncMessage(t *testing.T, signer *identity.Identity, message crdt.Message) string {
	t.Helper()
	documentPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.LiveDocumentSpec))
	if err != nil {
		t.Fatalf("document pCID: %v", err)
	}
	payloadBytes, err := protocol.Marshal(message)
	if err != nil {
		t.Fatalf("marshal sync payload: %v", err)
	}
	envelope := protocol.NewEnvelope(documentPCID, payloadBytes, protocol.Proof{})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("sync signable bytes: %v", err)
	}
	proof, err := signer.SignProof(signable)
	if err != nil {
		t.Fatalf("sign sync proof: %v", err)
	}
	envelope.Proof = proof
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("sync envelope bytes: %v", err)
	}
	return base64.StdEncoding.EncodeToString(envelopeBytes)
}

func signAwarenessMessage(t *testing.T, signer *identity.Identity, message awareness.Message) string {
	t.Helper()
	awarenessPCID, err := protocol.CIDForBytes(protocols.MustRead(protocols.LiveAwarenessSpec))
	if err != nil {
		t.Fatalf("awareness pCID: %v", err)
	}
	payloadBytes, err := protocol.Marshal(message)
	if err != nil {
		t.Fatalf("marshal awareness payload: %v", err)
	}
	envelope := protocol.NewEnvelope(awarenessPCID, payloadBytes, protocol.Proof{})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("awareness signable bytes: %v", err)
	}
	proof, err := signer.SignProof(signable)
	if err != nil {
		t.Fatalf("sign awareness proof: %v", err)
	}
	envelope.Proof = proof
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("awareness envelope bytes: %v", err)
	}
	return base64.StdEncoding.EncodeToString(envelopeBytes)
}

func decodedBase64(t *testing.T, value string) string {
	t.Helper()
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		t.Fatalf("decode base64: %v", err)
	}
	return string(bytes)
}
