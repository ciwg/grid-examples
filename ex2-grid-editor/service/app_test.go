package service_test

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex2-grid-editor/service"
)

func TestPostSyncAppearsInFeedAndReplay(t *testing.T) {
	t.Parallel()
	appA, err := service.NewApp(filepath.Join(t.TempDir(), "relay-a"))
	if err != nil {
		t.Fatalf("new app a: %v", err)
	}
	record, err := appA.PostSync("demo", "browser-a", "browser-b", base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}), "browser")
	if err != nil {
		t.Fatalf("post sync: %v", err)
	}
	feed := appA.SyncFeed("demo", 0)
	if len(feed.Messages) != 1 {
		t.Fatalf("message count mismatch: got %d", len(feed.Messages))
	}
	if feed.Messages[0].EnvelopeCID != record.EnvelopeCID {
		t.Fatalf("envelope cid mismatch: got %s want %s", feed.Messages[0].EnvelopeCID, record.EnvelopeCID)
	}
	rawMessages, _ := appA.PeerMessagesSince(0)
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
	feedB := appB.SyncFeed("demo", 0)
	if len(feedB.Messages) != 1 {
		t.Fatalf("replayed message count mismatch: got %d", len(feedB.Messages))
	}
	if feedB.Messages[0].MessageBase64 != record.MessageBase64 {
		t.Fatalf("message bytes mismatch: got %s want %s", feedB.Messages[0].MessageBase64, record.MessageBase64)
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
}

func TestAwarenessUsesAuthorFallbackAndLatestRelayState(t *testing.T) {
	t.Parallel()
	app, err := service.NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	if err := app.UpdateAwareness("demo", "", 12, 14, true, "Alice", "#123456", "nvim"); err != nil {
		t.Fatalf("update awareness initial: %v", err)
	}
	if err := app.UpdateAwareness("demo", "", 1, 1, false, "Alice stale", "#654321", "nvim"); err != nil {
		t.Fatalf("update awareness stale: %v", err)
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
}
