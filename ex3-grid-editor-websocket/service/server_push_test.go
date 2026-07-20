package service

import (
	"encoding/base64"
	"path/filepath"
	"testing"
	"time"
)

func TestServerNotifiesSyncSubscribersImmediately(t *testing.T) {
	t.Parallel()

	app, err := NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)
	updates, unsubscribe := server.subscribeSync("demo")
	defer unsubscribe()

	if _, err := app.PostSync("demo", "browser-a", "", base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}), "browser", "", ""); err != nil {
		t.Fatalf("post sync: %v", err)
	}

	select {
	case <-updates:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for sync subscriber update")
	}
}

func TestServerNotifiesAwarenessSubscribersImmediately(t *testing.T) {
	t.Parallel()

	app, err := NewApp(filepath.Join(t.TempDir(), "relay"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewServer(app)
	updates, unsubscribe := server.subscribeAwareness("demo")
	defer unsubscribe()

	if err := app.UpdateAwareness("demo", "browser-a", 7, 7, true, "Alice", "#1d6fd6", "browser"); err != nil {
		t.Fatalf("update awareness: %v", err)
	}

	select {
	case <-updates:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for awareness subscriber update")
	}
}
