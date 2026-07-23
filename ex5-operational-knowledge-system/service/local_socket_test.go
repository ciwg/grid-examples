package service

import (
	"bufio"
	"encoding/json"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLocalEmbodimentServerHandlesRequestResponse(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	server := NewLocalEmbodimentServer(app, EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket")))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, server.socketPath)

	response := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:   "request",
		Method: "GET",
		Path:   "/api/meta",
	})
	if response.Type != "response" || response.Status != 200 {
		t.Fatalf("unexpected local socket response: %+v", response)
	}
	if !strings.Contains(response.Body, `"local_unix_socket_enabled":true`) {
		t.Fatalf("missing local socket capability metadata: %s", response.Body)
	}
}

func TestLocalEmbodimentServerHandlesTypedReadOperations(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	place, err := app.CreatePlace("alice", "area", "Receiving", "Inbound inspection area", "", nil)
	if err != nil {
		t.Fatalf("create place: %v", err)
	}
	resource, err := app.CreateResource("alice", "container", "RJ45 Bin", "Connector bin", place.ID, nil)
	if err != nil {
		t.Fatalf("create resource: %v", err)
	}
	responsibility, err := app.CreateResponsibility("alice", "Receiving lead", "Owns intake review", []string{"reviewer"}, nil)
	if err != nil {
		t.Fatalf("create responsibility: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindReceiving, "Start line", "Startup checklist", "# Start", nil, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	run, err := app.RecordRun("bob", RunKindReceiving, item.ID, item.CurrentRevision, "accepted_with_notes", "Outer wrap torn", "", "", place.ID, []string{resource.ID}, []string{responsibility.ID})
	if err != nil {
		t.Fatalf("record run: %v", err)
	}

	server := NewLocalEmbodimentServer(app, EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket")))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, server.socketPath)

	itemResponse := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:      "operation",
		Operation: "inspect_item",
		ItemID:    item.ID,
	})
	if itemResponse.Status != 200 {
		t.Fatalf("unexpected inspect item status: %+v", itemResponse)
	}
	var inspectedItem KnowledgeItem
	if err := json.Unmarshal([]byte(itemResponse.Body), &inspectedItem); err != nil {
		t.Fatalf("decode inspect item: %v", err)
	}
	if inspectedItem.ID != item.ID {
		t.Fatalf("unexpected inspected item id: %q", inspectedItem.ID)
	}

	entityResponse := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:       "operation",
		Operation:  "inspect_entity",
		EntityType: "place",
		EntityID:   place.ID,
	})
	if entityResponse.Status != 200 {
		t.Fatalf("unexpected inspect entity status: %+v", entityResponse)
	}
	var inspectedPlace Place
	if err := json.Unmarshal([]byte(entityResponse.Body), &inspectedPlace); err != nil {
		t.Fatalf("decode inspect entity: %v", err)
	}
	if inspectedPlace.ID != place.ID {
		t.Fatalf("unexpected inspected place id: %q", inspectedPlace.ID)
	}

	searchResponse := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:          "operation",
		Operation:     "search",
		SearchOptions: &SearchOptions{Query: "start"},
	})
	if searchResponse.Status != 200 {
		t.Fatalf("unexpected search status: %+v", searchResponse)
	}
	var searchProjection struct {
		Items []KnowledgeItem `json:"items"`
		Runs  []RunRecord     `json:"runs"`
	}
	if err := json.Unmarshal([]byte(searchResponse.Body), &searchProjection); err != nil {
		t.Fatalf("decode search projection: %v", err)
	}
	if len(searchProjection.Items) != 1 || searchProjection.Items[0].ID != item.ID {
		t.Fatalf("unexpected search items: %+v", searchProjection.Items)
	}
	if len(searchProjection.Runs) != 0 {
		t.Fatalf("unexpected search runs: %+v", searchProjection.Runs)
	}

	pendingResponse := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:      "operation",
		Operation: "pending_review",
	})
	if pendingResponse.Status != 200 {
		t.Fatalf("unexpected pending review status: %+v", pendingResponse)
	}
	var pending PendingReviewProjection
	if err := json.Unmarshal([]byte(pendingResponse.Body), &pending); err != nil {
		t.Fatalf("decode pending review: %v", err)
	}
	if len(pending.DraftItems) != 1 || pending.DraftItems[0].ID != item.ID {
		t.Fatalf("unexpected pending draft items: %+v", pending.DraftItems)
	}
	if len(pending.UnreviewedRuns) != 1 || pending.UnreviewedRuns[0].ID != run.ID {
		t.Fatalf("unexpected pending unreviewed runs: %+v", pending.UnreviewedRuns)
	}
	if len(pending.ProblemRuns) != 1 || pending.ProblemRuns[0].ID != run.ID {
		t.Fatalf("unexpected pending problem runs: %+v", pending.ProblemRuns)
	}

	problemResponse := localSocketRequest(t, server.socketPath, LocalEmbodimentRequest{
		Type:      "operation",
		Operation: "problem_review",
	})
	if problemResponse.Status != 200 {
		t.Fatalf("unexpected problem review status: %+v", problemResponse)
	}
	var review ProblemReview
	if err := json.Unmarshal([]byte(problemResponse.Body), &review); err != nil {
		t.Fatalf("decode problem review: %v", err)
	}
	if review.ProblemRuns != 1 {
		t.Fatalf("unexpected problem run count: %d", review.ProblemRuns)
	}
	if len(review.PlaceGroups) != 1 || review.PlaceGroups[0].GroupID != place.ID {
		t.Fatalf("unexpected place groups: %+v", review.PlaceGroups)
	}
}

func TestLocalEmbodimentServerRejectsSecondActiveRuntimeOnSameSocket(t *testing.T) {
	socketRoot := filepath.Join(t.TempDir(), "runtime")
	appA, err := NewApp(filepath.Join(socketRoot, "a"))
	if err != nil {
		t.Fatalf("new app A: %v", err)
	}
	socketPath := EmbodimentSocketPath(socketRoot)
	serverA := NewLocalEmbodimentServer(appA, socketPath)
	firstErr := make(chan error, 1)
	go func() {
		firstErr <- serverA.ListenAndServe()
	}()
	defer func() {
		_ = serverA.Close()
	}()
	waitForUnixSocket(t, socketPath)

	appB, err := NewApp(filepath.Join(socketRoot, "b"))
	if err != nil {
		t.Fatalf("new app B: %v", err)
	}
	serverB := NewLocalEmbodimentServer(appB, socketPath)
	secondErr := make(chan error, 1)
	go func() {
		secondErr <- serverB.ListenAndServe()
	}()
	select {
	case err := <-secondErr:
		if err == nil || !strings.Contains(err.Error(), "already owned by an active runtime") {
			t.Fatalf("unexpected second server error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("second server did not fail fast")
	}

	response := localSocketRequest(t, socketPath, LocalEmbodimentRequest{
		Type:   "request",
		Method: "GET",
		Path:   "/api/meta",
	})
	if response.Type != "response" || response.Status != 200 {
		t.Fatalf("first server no longer reachable after collision attempt: %+v", response)
	}
	select {
	case err := <-firstErr:
		if err != nil {
			t.Fatalf("first server exited unexpectedly: %v", err)
		}
	default:
	}
}

func TestLocalEmbodimentServerStreamsLiveDraftUpdates(t *testing.T) {
	app, err := NewApp(filepath.Join(t.TempDir(), "runtime"))
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	item, err := app.CreateKnowledgeItem("alice", KnowledgeKindProcedure, "Start line", "startup", "# Start", nil, nil)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	server := NewLocalEmbodimentServer(app, EmbodimentSocketPath(filepath.Join(t.TempDir(), "runtime-socket")))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Errorf("listen and serve: %v", err)
		}
	}()
	defer func() {
		_ = server.Close()
	}()
	waitForUnixSocket(t, server.socketPath)

	conn, err := net.DialTimeout("unix", server.socketPath, 2*time.Second)
	if err != nil {
		t.Fatalf("dial local socket: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(LocalEmbodimentRequest{
		Type:          "live-open",
		ItemID:        item.ID,
		ParticipantID: "nvim-a",
		DisplayName:   "Alice",
		Color:         "#123456",
	}); err != nil {
		t.Fatalf("encode live open: %v", err)
	}
	reader := bufio.NewReader(conn)
	initial := readLocalSocketResponse(t, reader)
	if initial.Type != "live-state" || initial.State.Body != "# Start" {
		t.Fatalf("unexpected initial live state: %+v", initial)
	}
	if err := json.NewEncoder(conn).Encode(LocalEmbodimentRequest{
		Type:          "live-update",
		ItemID:        item.ID,
		ParticipantID: "nvim-a",
		DisplayName:   "Alice",
		Color:         "#123456",
		Cursor:        4,
		Head:          4,
		BaseVersion:   initial.State.Version,
		UpdateBody:    true,
		Body:          "# Start\n\nEdited",
	}); err != nil {
		t.Fatalf("encode live update: %v", err)
	}
	updated := readLocalSocketResponse(t, reader)
	if updated.Type != "live-state" || updated.State.Body != "# Start\n\nEdited" {
		t.Fatalf("unexpected updated live state: %+v", updated)
	}
}

func waitForUnixSocket(t *testing.T, socketPath string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := net.DialTimeout("unix", socketPath, 100*time.Millisecond); err == nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("socket did not become ready: %s", socketPath)
}

func localSocketRequest(t *testing.T, socketPath string, request LocalEmbodimentRequest) LocalEmbodimentResponse {
	t.Helper()
	conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
	if err != nil {
		t.Fatalf("dial local socket: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	if err := json.NewEncoder(conn).Encode(request); err != nil {
		t.Fatalf("encode request: %v", err)
	}
	return readLocalSocketResponse(t, bufio.NewReader(conn))
}

func readLocalSocketResponse(t *testing.T, reader *bufio.Reader) LocalEmbodimentResponse {
	t.Helper()
	var response LocalEmbodimentResponse
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		t.Fatalf("decode local socket response: %v", err)
	}
	return response
}
