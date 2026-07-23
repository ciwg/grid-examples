package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/service"
)

func main() {
	socketPath := flag.String("socket", "", "local unix socket path, or 'off' to force HTTP compatibility transport")
	serverURL := flag.String("server", "http://127.0.0.1:7045", "server URL")
	flag.Parse()
	resolvedServerURL, resolvedSocketPath := resolveCLITransportConfig(*serverURL, *socketPath)
	cli := &CLI{
		ServerURL:  resolvedServerURL,
		SocketPath: resolvedSocketPath,
	}
	exitCode, err := cli.run(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// Intent: Resolve the CLI's embodiment transport choice in one small helper so
// the explicit `-socket=off` compatibility path and the runtime-first direct
// socket path are both testable without subprocess harnessing. Source:
// DI-lurav
func resolveCLITransportConfig(serverURL string, socketOption string) (string, string) {
	resolvedServerURL := strings.TrimRight(serverURL, "/")
	resolvedSocketPath := strings.TrimSpace(socketOption)
	if strings.EqualFold(resolvedSocketPath, "off") {
		return resolvedServerURL, ""
	}
	if resolvedSocketPath == "" {
		resolvedSocketPath = discoverSocketPath(resolvedServerURL)
	}
	return resolvedServerURL, resolvedSocketPath
}

// Intent: Ask the runtime for its canonical socket path before relying on
// filesystem guesses so a custom `-data-root` still yields the direct terminal
// embodiment contract unless the operator explicitly opts into HTTP
// compatibility mode. Source: DI-sorek; DI-zorav
func discoverSocketPath(serverURL string) string {
	socketPath, err := socketPathFromMeta(serverURL)
	if err == nil {
		return socketPath
	}
	return defaultSocketPath()
}

func socketPathFromMeta(serverURL string) (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Get(serverURL + "/api/meta")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("meta discovery failed: %s", response.Status)
	}
	var meta service.Meta
	if err := json.NewDecoder(response.Body).Decode(&meta); err != nil {
		return "", err
	}
	if !meta.LocalUnixSocketEnabled || strings.TrimSpace(meta.LocalUnixSocketPath) == "" {
		return "", errors.New("runtime does not advertise a direct local socket path")
	}
	return meta.LocalUnixSocketPath, nil
}

// Intent: Keep a local best-effort fallback when the runtime is not yet
// reachable so CLI startup still has a stable socket guess before HTTP
// discovery exists. Source: DI-vorag; DI-sorek
func defaultSocketPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.Join(".operational-knowledge-system", "embodiment.sock")
	}
	current := cwd
	for {
		runtimeRoot := filepath.Join(current, ".operational-knowledge-system")
		if info, statErr := os.Stat(runtimeRoot); statErr == nil && info.IsDir() {
			return filepath.Join(runtimeRoot, "embodiment.sock")
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return filepath.Join(cwd, ".operational-knowledge-system", "embodiment.sock")
}

// Intent: Keep the CLI embodiment argument handling aligned with the shared run-record model so browser and CLI can create the same durable run records. Source: DI-zuvob
func (cli *CLI) run(args []string) (int, error) {
	if len(args) == 0 {
		return 2, fmt.Errorf("%s", usageText)
	}
	var err error
	switch args[0] {
	case "dashboard":
		err = cli.Dashboard()
	case "problem-review":
		err = cli.ProblemReview()
	case "pending-review":
		err = cli.PendingReview()
	case "places":
		err = cli.Places()
	case "new-place":
		if len(args) < 5 {
			return 2, fmt.Errorf("%s", usageText)
		}
		parentID := ""
		if len(args) > 5 {
			parentID = args[5]
		}
		err = cli.NewPlace(args[1], args[2], args[3], args[4], parentID)
	case "show-place":
		if len(args) != 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.ShowPlace(args[1])
	case "resources":
		err = cli.Resources()
	case "new-resource":
		if len(args) < 5 {
			return 2, fmt.Errorf("%s", usageText)
		}
		placeID := ""
		if len(args) > 5 {
			placeID = args[5]
		}
		err = cli.NewResource(args[1], args[2], args[3], args[4], placeID)
	case "show-resource":
		if len(args) != 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.ShowResource(args[1])
	case "responsibilities":
		err = cli.Responsibilities()
	case "show-responsibility":
		if len(args) != 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.ShowResponsibility(args[1])
	case "new-responsibility":
		if len(args) < 4 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.NewResponsibility(args[1], args[2], strings.Join(args[3:], " "))
	case "items":
		kind := ""
		if len(args) > 1 {
			kind = args[1]
		}
		err = cli.Items(kind)
	case "new-item":
		if len(args) < 6 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.NewItem(args[1], args[2], args[3], args[4], strings.Join(args[5:], " "))
	case "show-item":
		if len(args) != 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.ShowItem(args[1])
	case "snapshot-item":
		if len(args) < 4 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.SnapshotItem(args[1], args[2], strings.Join(args[3:], " "))
	case "runs":
		kind := ""
		if len(args) > 1 {
			kind = args[1]
		}
		err = cli.Runs(kind)
	case "record-run":
		if len(args) < 7 {
			return 2, fmt.Errorf("%s", usageText)
		}
		revision, convErr := strconv.Atoi(args[4])
		if convErr != nil {
			return 1, convErr
		}
		notes := args[6]
		placeID := ""
		resourceIDs := []string{}
		if len(args) > 7 {
			placeID = args[7]
		}
		if len(args) > 8 {
			resourceIDs = splitCSV(args[8])
		}
		err = cli.RecordRun(args[1], args[2], args[3], revision, args[5], notes, placeID, resourceIDs)
	case "show-run":
		if len(args) != 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.ShowRun(args[1])
	case "approve-item":
		if len(args) < 7 {
			return 2, fmt.Errorf("%s", usageText)
		}
		revision, convErr := strconv.Atoi(args[2])
		if convErr != nil {
			return 1, convErr
		}
		err = cli.Approve("/api/items/"+args[1]+"/approvals", args[3], revision, args[4], args[5], strings.Join(args[6:], " "))
	case "supersede-item":
		if len(args) < 3 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.post("/api/items/"+args[1]+"/supersede", map[string]any{
			"actor": args[2],
			"notes": strings.Join(args[3:], " "),
		})
	case "approve-run":
		if len(args) < 6 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.Approve("/api/runs/"+args[1]+"/approvals", args[2], 0, args[3], args[4], strings.Join(args[5:], " "))
	case "add-link":
		if len(args) < 7 {
			return 2, fmt.Errorf("%s", usageText)
		}
		notes := ""
		if len(args) > 7 {
			notes = strings.Join(args[7:], " ")
		}
		err = cli.AddLink(args[1], args[2], args[3], args[4], args[5], args[6], notes)
	case "add-evidence":
		if len(args) < 4 {
			return 2, fmt.Errorf("%s", usageText)
		}
		factsJSON := ""
		attachmentPath := ""
		if len(args) > 4 {
			factsJSON = args[4]
		}
		if len(args) > 5 {
			attachmentPath = args[5]
		}
		err = cli.AddEvidence(args[1], args[2], args[3], factsJSON, attachmentPath)
	case "search":
		if len(args) < 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.Search(args[1], args[2:])
	default:
		return 2, fmt.Errorf("%s", usageText)
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}

const usageText = "usage: oks-cli dashboard|problem-review|pending-review|places|new-place ACTOR KIND NAME SUMMARY [PARENT_ID]|show-place ID|resources|new-resource ACTOR KIND NAME SUMMARY [PLACE_ID]|show-resource ID|responsibilities|show-responsibility ID|new-responsibility ACTOR TITLE SUMMARY|items [kind]|new-item ACTOR KIND TITLE SUMMARY BODY|show-item ID|snapshot-item ITEM_ID ACTOR BODY|runs [kind]|record-run ACTOR KIND ITEM_ID REVISION OUTCOME NOTES [PLACE_ID] [RESOURCE_IDS_CSV]|show-run ID|approve-item ITEM_ID REVISION ACTOR ROLE DECISION NOTES|supersede-item ITEM_ID ACTOR [NOTES]|approve-run RUN_ID ACTOR ROLE DECISION NOTES|add-link ACTOR FROM_TYPE FROM_ID TO_TYPE TO_ID RELATION [NOTES]|add-evidence RUN_ID ACTOR SUMMARY [FACTS_JSON] [FILE]|search QUERY [kind=VALUE] [status=VALUE] [outcome=VALUE] [place_id=VALUE] [resource_id=VALUE] [responsibility_id=VALUE] [problem=true]"

type CLI struct {
	ServerURL  string
	SocketPath string
}

type cliLink struct {
	Relation string `json:"relation"`
	FromType string `json:"from_type"`
	FromID   string `json:"from_id"`
	ToType   string `json:"to_type"`
	ToID     string `json:"to_id"`
	Notes    string `json:"notes"`
}

type cliRunSummary struct {
	ID          string   `json:"id"`
	AliasID     string   `json:"alias_id"`
	Kind        string   `json:"kind"`
	ItemID      string   `json:"item_id"`
	Outcome     string   `json:"outcome"`
	Notes       string   `json:"notes"`
	ResourceIDs []string `json:"resource_ids"`
}

type cliPlaceDetail struct {
	ID            string          `json:"id"`
	Kind          string          `json:"kind"`
	Name          string          `json:"name"`
	Summary       string          `json:"summary"`
	ParentID      string          `json:"parent_id"`
	ChildPlaceIDs []string        `json:"child_place_ids"`
	ResourceIDs   []string        `json:"resource_ids"`
	RelatedRuns   []cliRunSummary `json:"related_runs"`
	Links         []cliLink       `json:"links"`
}

type cliResourceDetail struct {
	ID          string          `json:"id"`
	Kind        string          `json:"kind"`
	Name        string          `json:"name"`
	Summary     string          `json:"summary"`
	PlaceID     string          `json:"place_id"`
	RelatedRuns []cliRunSummary `json:"related_runs"`
	Links       []cliLink       `json:"links"`
}

type cliApproval struct {
	Role     string `json:"role"`
	Decision string `json:"decision"`
	Actor    string `json:"actor"`
	Notes    string `json:"notes"`
}

type cliEvidence struct {
	Summary        string            `json:"summary"`
	Facts          map[string]string `json:"facts"`
	AttachmentName string            `json:"attachment_name"`
}

type cliRevision struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Author  string `json:"author"`
}

type cliResponsibilityDetail struct {
	ID             string          `json:"id"`
	AliasID        string          `json:"alias_id"`
	Title          string          `json:"title"`
	Summary        string          `json:"summary"`
	Team           string          `json:"team"`
	LinkedItemIDs  []string        `json:"linked_item_ids"`
	LinkedRunIDs   []string        `json:"linked_run_ids"`
	LinkedRoleKeys []string        `json:"linked_role_keys"`
	RelatedRuns    []cliRunSummary `json:"related_runs"`
	Links          []cliLink       `json:"links"`
}

type cliItemDetail struct {
	ID                string          `json:"id"`
	AliasID           string          `json:"alias_id"`
	Kind              string          `json:"kind"`
	Status            string          `json:"status"`
	Title             string          `json:"title"`
	Summary           string          `json:"summary"`
	Tags              []string        `json:"tags"`
	ResponsibilityIDs []string        `json:"responsibility_ids"`
	CurrentRevision   int             `json:"current_revision"`
	Revisions         []cliRevision   `json:"revisions"`
	RelatedRuns       []cliRunSummary `json:"related_runs"`
	Approvals         []cliApproval   `json:"approvals"`
	Links             []cliLink       `json:"links"`
}

type cliRunDetail struct {
	ID                string        `json:"id"`
	AliasID           string        `json:"alias_id"`
	Kind              string        `json:"kind"`
	ItemID            string        `json:"item_id"`
	ItemKind          string        `json:"item_kind"`
	Revision          int           `json:"revision"`
	Actor             string        `json:"actor"`
	Outcome           string        `json:"outcome"`
	Notes             string        `json:"notes"`
	PlaceID           string        `json:"place_id"`
	ResourceIDs       []string      `json:"resource_ids"`
	ResponsibilityIDs []string      `json:"responsibility_ids"`
	Evidence          []cliEvidence `json:"evidence"`
	Approvals         []cliApproval `json:"approvals"`
	Links             []cliLink     `json:"links"`
}

type cliSearchItem struct {
	ID      string `json:"id"`
	AliasID string `json:"alias_id"`
	Kind    string `json:"kind"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type cliSearchRun struct {
	ID          string        `json:"id"`
	AliasID     string        `json:"alias_id"`
	Kind        string        `json:"kind"`
	ItemID      string        `json:"item_id"`
	Outcome     string        `json:"outcome"`
	Notes       string        `json:"notes"`
	ResourceIDs []string      `json:"resource_ids"`
	Approvals   []cliApproval `json:"approvals"`
}

// Intent: Treat shared search runs without an explicit approvals array as
// contract drift instead of silently reclassifying them as genuine unreviewed
// work in terminal review queues. Source: DI-davur
func (run *cliSearchRun) UnmarshalJSON(body []byte) error {
	type rawSearchRun struct {
		ID          string           `json:"id"`
		AliasID     string           `json:"alias_id"`
		Kind        string           `json:"kind"`
		ItemID      string           `json:"item_id"`
		Outcome     string           `json:"outcome"`
		Notes       string           `json:"notes"`
		ResourceIDs []string         `json:"resource_ids"`
		Approvals   *json.RawMessage `json:"approvals"`
	}

	var raw rawSearchRun
	if err := json.Unmarshal(body, &raw); err != nil {
		return err
	}
	if raw.Approvals == nil {
		return errors.New(`search run missing "approvals" array`)
	}
	trimmed := bytes.TrimSpace(*raw.Approvals)
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return errors.New(`search run "approvals" field is not an array`)
	}
	var approvals []cliApproval
	if err := json.Unmarshal(trimmed, &approvals); err != nil {
		return fmt.Errorf("search run approvals decode: %w", err)
	}

	run.ID = raw.ID
	run.AliasID = raw.AliasID
	run.Kind = raw.Kind
	run.ItemID = raw.ItemID
	run.Outcome = raw.Outcome
	run.Notes = raw.Notes
	run.ResourceIDs = raw.ResourceIDs
	run.Approvals = approvals
	return nil
}

type cliSearchResponse struct {
	Items []cliSearchItem `json:"items"`
	Runs  []cliSearchRun  `json:"runs"`
}

type cliProblemReview struct {
	ProblemRuns    int                     `json:"problem_runs"`
	PlaceGroups    []cliProblemReviewGroup `json:"place_groups"`
	ResourceGroups []cliProblemReviewGroup `json:"resource_groups"`
}

type cliProblemReviewGroup struct {
	GroupID           string          `json:"group_id"`
	Kind              string          `json:"kind"`
	Name              string          `json:"name"`
	ProblemCount      int             `json:"problem_count"`
	ReceivingProblems int             `json:"receiving_problems"`
	InventoryProblems int             `json:"inventory_problems"`
	HighlightExamples []string        `json:"highlights"`
	Runs              []cliRunSummary `json:"runs"`
}

type cliPendingReviewProjection struct {
	DraftItems     []cliSearchItem `json:"draft_items"`
	UnreviewedRuns []cliSearchRun  `json:"unreviewed_runs"`
	ProblemRuns    []cliSearchRun  `json:"problem_runs"`
}

var allowedSearchFilters = map[string]struct{}{
	"kind":              {},
	"status":            {},
	"outcome":           {},
	"place_id":          {},
	"resource_id":       {},
	"responsibility_id": {},
	"problem":           {},
}

func (cli *CLI) Dashboard() error { return cli.Show("/api/dashboard") }

// Intent: Keep place drilldowns on the shared place-detail projection while
// rendering the hierarchy, links, and related runs in a terminal-first review
// layout that is easier to act on than a raw JSON dump. Source: DI-luzom
func (cli *CLI) ShowPlace(placeID string) error {
	body, err := cli.entityBody("place", placeID)
	if err != nil {
		return err
	}
	var place cliPlaceDetail
	if err := json.Unmarshal(body, &place); err != nil {
		return err
	}
	fmt.Println(renderPlaceDetail(place))
	return nil
}

// Intent: Keep resource drilldowns on the shared resource-detail projection
// while rendering context, links, and related runs in a terminal-first review
// layout that is easier to act on than a raw JSON dump. Source: DI-luzom
func (cli *CLI) ShowResource(resourceID string) error {
	body, err := cli.entityBody("resource", resourceID)
	if err != nil {
		return err
	}
	var resource cliResourceDetail
	if err := json.Unmarshal(body, &resource); err != nil {
		return err
	}
	fmt.Println(renderResourceDetail(resource))
	return nil
}

// Intent: Keep responsibility drilldowns on the shared responsibility-detail
// projection while rendering linked items, runs, roles, and graph links in the
// same terminal-first review layout used by the newer context drilldowns.
// Source: DI-salup
func (cli *CLI) ShowResponsibility(responsibilityID string) error {
	body, err := cli.entityBody("responsibility", responsibilityID)
	if err != nil {
		return err
	}
	var responsibility cliResponsibilityDetail
	if err := json.Unmarshal(body, &responsibility); err != nil {
		return err
	}
	fmt.Println(renderResponsibilityDetail(responsibility))
	return nil
}

// Intent: Keep item drilldowns on the shared item-detail projection while
// rendering revisions, approvals, related runs, and typed links in a
// terminal-first review layout instead of a raw JSON dump. Source: DI-salup
func (cli *CLI) ShowItem(itemID string) error {
	body, err := cli.itemBody(itemID)
	if err != nil {
		return err
	}
	var item cliItemDetail
	if err := json.Unmarshal(body, &item); err != nil {
		return err
	}
	fmt.Println(renderItemDetail(item))
	return nil
}

// Intent: Let shell-first operators cut a durable item revision without
// switching to the browser or Neovim, while still reusing the existing item
// revision route and the item's current title, summary, and tags. Source:
// DI-muvok
func (cli *CLI) SnapshotItem(itemID string, actor string, body string) error {
	itemBody, err := cli.get("/api/items/" + itemID)
	if err != nil {
		return err
	}
	var item cliItemDetail
	if err := json.Unmarshal(itemBody, &item); err != nil {
		return err
	}
	return cli.post("/api/items/"+itemID+"/revisions", map[string]any{
		"actor":   actor,
		"title":   item.Title,
		"summary": item.Summary,
		"body":    body,
		"tags":    item.Tags,
	})
}

// Intent: Keep run drilldowns on the shared run-detail projection while
// rendering evidence, approvals, linked context, and follow-on handoff hints
// in a terminal-first review layout instead of a raw JSON dump. Source:
// DI-salup
func (cli *CLI) ShowRun(runID string) error {
	body, err := cli.runBody(runID)
	if err != nil {
		return err
	}
	var run cliRunDetail
	if err := json.Unmarshal(body, &run); err != nil {
		return err
	}
	fmt.Println(renderRunDetail(run))
	return nil
}

// Intent: Keep the shell-first problem hotspot queue on the shared grouped
// problem-review projection while rendering it as a review-oriented terminal
// summary instead of a raw JSON blob. Source: DI-ravum
func (cli *CLI) ProblemReview() error {
	body, err := cli.problemReviewBody()
	if err != nil {
		return err
	}
	var review cliProblemReview
	if err := json.Unmarshal(body, &review); err != nil {
		return fmt.Errorf("/api/problem-review decode: %w", err)
	}
	fmt.Println(renderProblemReview(review))
	return nil
}

// Intent: Keep the shell-first pending-review queue on the same projection
// family Neovim already uses so terminal triage does not invent a separate
// review endpoint or drift away from the staged editor workflow. Source:
// DI-vabok
func (cli *CLI) PendingReview() error {
	if strings.TrimSpace(cli.SocketPath) != "" {
		body, err := cli.localSocketOperationBody(service.LocalEmbodimentRequest{
			Operation: "pending_review",
		})
		if err != nil {
			return err
		}
		var projection cliPendingReviewProjection
		if err := json.Unmarshal(body, &projection); err != nil {
			return fmt.Errorf("pending_review decode: %w", err)
		}
		fmt.Println(renderPendingReview(projection.DraftItems, projection.UnreviewedRuns, projection.ProblemRuns))
		return nil
	}

	draftSearch, err := cli.getSearch("/api/search?status=draft")
	if err != nil {
		return err
	}
	allRunsSearch, err := cli.getSearch("/api/search")
	if err != nil {
		return err
	}
	problemRunsSearch, err := cli.getSearch("/api/search?problem=true")
	if err != nil {
		return err
	}
	unreviewedRuns := make([]cliSearchRun, 0, len(allRunsSearch.Runs))
	for _, run := range allRunsSearch.Runs {
		if len(run.Approvals) == 0 {
			unreviewedRuns = append(unreviewedRuns, run)
		}
	}
	fmt.Println(renderPendingReview(draftSearch.Items, unreviewedRuns, problemRunsSearch.Runs))
	return nil
}
func (cli *CLI) Places() error           { return cli.Show("/api/places") }
func (cli *CLI) Resources() error        { return cli.Show("/api/resources") }
func (cli *CLI) Responsibilities() error { return cli.Show("/api/responsibilities") }
func (cli *CLI) Items(kind string) error {
	path := "/api/items"
	if kind != "" {
		path += "?kind=" + kind
	}
	return cli.Show(path)
}
func (cli *CLI) Runs(kind string) error {
	path := "/api/runs"
	if kind != "" {
		path += "?kind=" + kind
	}
	return cli.Show(path)
}

func (cli *CLI) NewResponsibility(actor string, title string, summary string) error {
	return cli.post("/api/responsibilities", map[string]any{"actor": actor, "title": title, "summary": summary})
}

func (cli *CLI) NewPlace(actor string, kind string, name string, summary string, parentID string) error {
	return cli.post("/api/places", map[string]any{
		"actor":     actor,
		"kind":      kind,
		"name":      name,
		"summary":   summary,
		"parent_id": parentID,
	})
}

func (cli *CLI) NewResource(actor string, kind string, name string, summary string, placeID string) error {
	return cli.post("/api/resources", map[string]any{
		"actor":    actor,
		"kind":     kind,
		"name":     name,
		"summary":  summary,
		"place_id": placeID,
	})
}

func (cli *CLI) NewItem(actor string, kind string, title string, summary string, body string) error {
	return cli.post("/api/items", map[string]any{"actor": actor, "kind": kind, "title": title, "summary": summary, "body": body})
}

func (cli *CLI) RecordRun(actor string, kind string, itemID string, revision int, outcome string, notes string, placeID string, resourceIDs []string) error {
	return cli.post("/api/runs", map[string]any{
		"actor":        actor,
		"kind":         kind,
		"item_id":      itemID,
		"revision":     revision,
		"outcome":      outcome,
		"notes":        notes,
		"place_id":     placeID,
		"resource_ids": resourceIDs,
	})
}

// Intent: Let shell-first operators build typed operational context over the
// same validated graph contract the browser already uses, without inventing a
// second terminal-only link schema. Source: DI-vuteg
func (cli *CLI) AddLink(actor string, fromType string, fromID string, toType string, toID string, relation string, notes string) error {
	return cli.post("/api/links", map[string]any{
		"actor":     actor,
		"from_type": fromType,
		"from_id":   fromID,
		"to_type":   toType,
		"to_id":     toID,
		"relation":  relation,
		"notes":     notes,
	})
}

// Intent: Close the terminal evidence gap by letting shell-first operators use
// the same run evidence multipart surface the browser already uses, including
// optional facts JSON and optional copied attachments. Source: DI-zanub
func (cli *CLI) AddEvidence(runID string, actor string, summary string, factsJSON string, attachmentPath string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("actor", actor); err != nil {
		return err
	}
	if err := writer.WriteField("summary", summary); err != nil {
		return err
	}
	if strings.TrimSpace(factsJSON) != "" {
		if err := writer.WriteField("facts_json", factsJSON); err != nil {
			return err
		}
	}
	if strings.TrimSpace(attachmentPath) != "" {
		attachmentBody, err := os.ReadFile(attachmentPath)
		if err != nil {
			return err
		}
		part, err := writer.CreateFormFile("attachment", filepath.Base(attachmentPath))
		if err != nil {
			return err
		}
		if _, err := part.Write(attachmentBody); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}

	status, message, err := cli.roundTrip("POST", "/api/runs/"+runID+"/evidence", writer.FormDataContentType(), body.Bytes())
	if err != nil {
		return err
	}
	if status >= 300 {
		return fmt.Errorf("%s", strings.TrimSpace(string(message)))
	}
	var indented bytes.Buffer
	if err := json.Indent(&indented, message, "", "  "); err == nil {
		fmt.Println(indented.String())
		return nil
	}
	fmt.Println(string(message))
	return nil
}

// Intent: Keep CLI search on the shared `/api/search` projection while letting
// shell-first operators use the same structured filters and problem-only view
// that already power browser and Neovim drilldowns. Source: DI-mifot
func (cli *CLI) Search(query string, filters []string) error {
	if strings.TrimSpace(cli.SocketPath) != "" {
		options, err := buildSearchOptions(query, filters)
		if err != nil {
			return err
		}
		body, err := cli.localSocketOperationBody(service.LocalEmbodimentRequest{
			Operation:     "search",
			SearchOptions: &options,
		})
		if err != nil {
			return err
		}
		return printJSONBytes(body)
	}
	path, err := buildSearchPath(query, filters)
	if err != nil {
		return err
	}
	return cli.Show(path)
}

// Intent: Preserve trustworthy approval history by making the CLI carry the
// real approver identity through to the shared runtime instead of inventing a
// placeholder actor. Source: DI-tarok
func (cli *CLI) Approve(path string, actor string, revision int, role string, decision string, notes string) error {
	return cli.post(path, map[string]any{
		"actor":    actor,
		"revision": revision,
		"role":     role,
		"decision": decision,
		"notes":    notes,
	})
}

func (cli *CLI) Show(path string) error {
	body, err := cli.get(path)
	if err != nil {
		return err
	}
	return printJSONBytes(body)
}

func (cli *CLI) post(path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	status, message, err := cli.roundTrip("POST", path, "application/json", body)
	if err != nil {
		return err
	}
	if status >= 300 {
		return fmt.Errorf("%s", strings.TrimSpace(string(message)))
	}
	return printJSONBytes(message)
}

func (cli *CLI) get(path string) ([]byte, error) {
	status, body, err := cli.roundTrip("GET", path, "", nil)
	if err != nil {
		return nil, err
	}
	if status >= 300 {
		return nil, fmt.Errorf("%s", strings.TrimSpace(string(body)))
	}
	return body, nil
}

func (cli *CLI) getJSON(path string) (map[string]any, error) {
	body, err := cli.get(path)
	if err != nil {
		return nil, err
	}
	var projection map[string]any
	if err := json.Unmarshal(body, &projection); err != nil {
		return nil, err
	}
	return projection, nil
}

func (cli *CLI) getSearchArray(path string, field string) ([]any, error) {
	projection, err := cli.getJSON(path)
	if err != nil {
		return nil, err
	}
	return requireJSONArray(projection, field, path)
}

func (cli *CLI) getSearch(path string) (cliSearchResponse, error) {
	body, err := cli.get(path)
	if err != nil {
		return cliSearchResponse{}, err
	}
	var response cliSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return cliSearchResponse{}, fmt.Errorf("%s decode: %w", path, err)
	}
	return response, nil
}

func printJSON(value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return printJSONBytes(body)
}

func printJSONBytes(body []byte) error {
	var indented bytes.Buffer
	if err := json.Indent(&indented, body, "", "  "); err == nil {
		fmt.Println(indented.String())
		return nil
	}
	fmt.Println(string(body))
	return nil
}

func requireJSONArray(object map[string]any, field string, context string) ([]any, error) {
	value, ok := object[field]
	if !ok {
		return nil, fmt.Errorf("%s missing %q array", context, field)
	}
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("%s %q field is not an array", context, field)
	}
	return items, nil
}

func requireJSONObject(value any, context string) (map[string]any, error) {
	if value == nil {
		return nil, fmt.Errorf("%s is null", context)
	}
	object, ok := value.(map[string]any)
	if !ok || object == nil {
		return nil, fmt.Errorf("%s is not an object", context)
	}
	return object, nil
}

func renderPlaceDetail(place cliPlaceDetail) string {
	lines := []string{
		fmt.Sprintf("# Place %s", safeText(place.ID, "-")),
		fmt.Sprintf("name=%s kind=%s parent=%s", safeText(place.Name, "-"), safeText(place.Kind, "-"), safeText(place.ParentID, "-")),
	}
	if strings.TrimSpace(place.Summary) != "" {
		lines = append(lines, place.Summary)
	}
	lines = append(lines, "")
	lines = append(lines, "child places: "+joinOrNone(place.ChildPlaceIDs))
	lines = append(lines, "resources: "+joinOrNone(place.ResourceIDs))
	lines = append(lines, "")
	lines = append(lines, "related runs:")
	lines = append(lines, renderRunLines(place.RelatedRuns)...)
	lines = append(lines, "")
	lines = append(lines, "links:")
	lines = append(lines, renderLinkLines(place.Links)...)
	return strings.Join(lines, "\n")
}

func renderResourceDetail(resource cliResourceDetail) string {
	lines := []string{
		fmt.Sprintf("# Resource %s", safeText(resource.ID, "-")),
		fmt.Sprintf("name=%s kind=%s place=%s", safeText(resource.Name, "-"), safeText(resource.Kind, "-"), safeText(resource.PlaceID, "-")),
	}
	if strings.TrimSpace(resource.Summary) != "" {
		lines = append(lines, resource.Summary)
	}
	lines = append(lines, "")
	lines = append(lines, "related runs:")
	lines = append(lines, renderRunLines(resource.RelatedRuns)...)
	lines = append(lines, "")
	lines = append(lines, "links:")
	lines = append(lines, renderLinkLines(resource.Links)...)
	return strings.Join(lines, "\n")
}

func renderResponsibilityDetail(responsibility cliResponsibilityDetail) string {
	lines := []string{
		fmt.Sprintf("# Responsibility %s", safeDisplayID(responsibility.ID, responsibility.AliasID)),
		fmt.Sprintf("title=%s team=%s", safeText(responsibility.Title, "-"), safeText(responsibility.Team, "-")),
	}
	if strings.TrimSpace(responsibility.Summary) != "" {
		lines = append(lines, responsibility.Summary)
	}
	lines = append(lines, "")
	lines = append(lines, "items: "+joinOrNone(responsibility.LinkedItemIDs))
	lines = append(lines, "role keys: "+joinOrNone(responsibility.LinkedRoleKeys))
	lines = append(lines, "linked runs: "+joinOrNone(responsibility.LinkedRunIDs))
	lines = append(lines, "")
	lines = append(lines, "related runs:")
	lines = append(lines, renderRunLines(responsibility.RelatedRuns)...)
	lines = append(lines, "")
	lines = append(lines, "links:")
	lines = append(lines, renderLinkLines(responsibility.Links)...)
	return strings.Join(lines, "\n")
}

func renderItemDetail(item cliItemDetail) string {
	lines := []string{
		fmt.Sprintf("# Item %s", safeDisplayID(item.ID, item.AliasID)),
		fmt.Sprintf("title=%s kind=%s status=%s current_revision=%d", safeText(item.Title, "-"), safeText(item.Kind, "-"), safeText(item.Status, "-"), item.CurrentRevision),
	}
	if strings.TrimSpace(item.Summary) != "" {
		lines = append(lines, item.Summary)
	}
	lines = append(lines, "")
	lines = append(lines, "responsibilities: "+joinOrNone(item.ResponsibilityIDs))
	lines = append(lines, "")
	lines = append(lines, "revisions:")
	lines = append(lines, renderRevisionLines(item.Revisions)...)
	lines = append(lines, "")
	lines = append(lines, "approvals:")
	lines = append(lines, renderApprovalLines(item.Approvals)...)
	lines = append(lines, "")
	lines = append(lines, "related runs:")
	lines = append(lines, renderRunLines(item.RelatedRuns)...)
	lines = append(lines, "")
	lines = append(lines, "links:")
	lines = append(lines, renderLinkLines(item.Links)...)
	return strings.Join(lines, "\n")
}

func renderRunDetail(run cliRunDetail) string {
	lines := []string{
		fmt.Sprintf("# Run %s", safeDisplayID(run.ID, run.AliasID)),
		fmt.Sprintf("kind=%s outcome=%s item=%s item_kind=%s revision=%d actor=%s", safeText(run.Kind, "-"), safeText(run.Outcome, "-"), safeText(run.ItemID, "-"), safeText(run.ItemKind, "-"), run.Revision, safeText(run.Actor, "-")),
	}
	if strings.TrimSpace(run.Notes) != "" {
		lines = append(lines, run.Notes)
	}
	lines = append(lines, "")
	lines = append(lines, "place: "+safeText(run.PlaceID, "none"))
	lines = append(lines, "resources: "+joinOrNone(run.ResourceIDs))
	lines = append(lines, "responsibilities: "+joinOrNone(run.ResponsibilityIDs))
	lines = append(lines, "show item: oks-cli show-item "+safeText(run.ItemID, "-"))
	// Intent: Keep terminal drilldowns navigable by turning run context fields
	// into direct follow-on CLI commands instead of leaving the operator to
	// reconstruct the next lookup by hand. Source: DI-josav
	if strings.TrimSpace(run.PlaceID) != "" {
		lines = append(lines, "show place: oks-cli show-place "+run.PlaceID)
	}
	for _, resourceID := range run.ResourceIDs {
		if strings.TrimSpace(resourceID) == "" {
			continue
		}
		lines = append(lines, "show resource: oks-cli show-resource "+resourceID)
	}
	for _, responsibilityID := range run.ResponsibilityIDs {
		if strings.TrimSpace(responsibilityID) == "" {
			continue
		}
		lines = append(lines, "show responsibility: oks-cli show-responsibility "+responsibilityID)
	}
	lines = append(lines, "")
	lines = append(lines, "evidence:")
	lines = append(lines, renderEvidenceLines(run.Evidence)...)
	lines = append(lines, "")
	lines = append(lines, "approvals:")
	lines = append(lines, renderApprovalLines(run.Approvals)...)
	lines = append(lines, "")
	lines = append(lines, "links:")
	lines = append(lines, renderLinkLines(run.Links)...)
	return strings.Join(lines, "\n")
}

func renderPendingReview(draftItems []cliSearchItem, unreviewedRuns []cliSearchRun, problemRuns []cliSearchRun) string {
	lines := []string{
		"# Pending review",
		fmt.Sprintf("draft_items=%d unreviewed_runs=%d problem_runs=%d", len(draftItems), len(unreviewedRuns), len(problemRuns)),
		"",
		"draft items:",
	}
	lines = append(lines, renderSearchItemLines(draftItems)...)
	lines = append(lines, "")
	lines = append(lines, "unreviewed runs:")
	lines = append(lines, renderSearchRunLines(unreviewedRuns)...)
	lines = append(lines, "")
	lines = append(lines, "problem runs:")
	lines = append(lines, renderSearchRunLines(problemRuns)...)
	return strings.Join(lines, "\n")
}

func renderProblemReview(review cliProblemReview) string {
	lines := []string{
		"# Problem review",
		fmt.Sprintf("problem_runs=%d", review.ProblemRuns),
		"",
		"place groups:",
	}
	lines = append(lines, renderProblemGroups(review.PlaceGroups, "place")...)
	lines = append(lines, "")
	lines = append(lines, "resource groups:")
	lines = append(lines, renderProblemGroups(review.ResourceGroups, "resource")...)
	return strings.Join(lines, "\n")
}

func renderRunLines(runs []cliRunSummary) []string {
	if len(runs) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(runs)*3)
	for _, run := range runs {
		lines = append(lines, fmt.Sprintf("- %s kind=%s outcome=%s item=%s", safeDisplayID(run.ID, run.AliasID), safeText(run.Kind, "-"), safeText(run.Outcome, "-"), safeText(run.ItemID, "-")))
		lines = append(lines, "  show: oks-cli show-run "+safeText(run.ID, "-"))
		if len(run.ResourceIDs) > 0 {
			lines = append(lines, "  resources: "+strings.Join(run.ResourceIDs, ", "))
		}
		if strings.TrimSpace(run.Notes) != "" {
			lines = append(lines, "  "+run.Notes)
		}
	}
	return lines
}

func renderSearchItemLines(items []cliSearchItem) []string {
	if len(items) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(items)*3)
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("- %s kind=%s status=%s title=%s", safeDisplayID(item.ID, item.AliasID), safeText(item.Kind, "-"), safeText(item.Status, "-"), safeText(item.Title, "-")))
		lines = append(lines, "  show: oks-cli show-item "+safeText(item.ID, "-"))
		if strings.TrimSpace(item.Summary) != "" {
			lines = append(lines, "  "+item.Summary)
		}
	}
	return lines
}

func renderSearchRunLines(runs []cliSearchRun) []string {
	if len(runs) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(runs)*3)
	for _, run := range runs {
		lines = append(lines, fmt.Sprintf("- %s kind=%s outcome=%s item=%s approvals=%d", safeDisplayID(run.ID, run.AliasID), safeText(run.Kind, "-"), safeText(run.Outcome, "-"), safeText(run.ItemID, "-"), len(run.Approvals)))
		lines = append(lines, "  show: oks-cli show-run "+safeText(run.ID, "-"))
		if len(run.ResourceIDs) > 0 {
			lines = append(lines, "  resources: "+strings.Join(run.ResourceIDs, ", "))
		}
		if strings.TrimSpace(run.Notes) != "" {
			lines = append(lines, "  "+run.Notes)
		}
	}
	return lines
}

func renderProblemGroups(groups []cliProblemReviewGroup, groupType string) []string {
	if len(groups) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(groups)*6)
	for _, group := range groups {
		lines = append(lines, fmt.Sprintf("- %s kind=%s name=%s problems=%d receiving=%d inventory=%d", safeText(group.GroupID, "-"), safeText(group.Kind, "-"), safeText(group.Name, "-"), group.ProblemCount, group.ReceivingProblems, group.InventoryProblems))
		lines = append(lines, "  show: oks-cli show-"+groupType+" "+safeText(group.GroupID, "-"))
		if len(group.HighlightExamples) > 0 {
			lines = append(lines, "  highlights: "+strings.Join(group.HighlightExamples, " | "))
		}
		lines = append(lines, "  runs:")
		for _, runLine := range renderRunLines(group.Runs) {
			lines = append(lines, "  "+runLine)
		}
	}
	return lines
}

func renderLinkLines(links []cliLink) []string {
	if len(links) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(links)*2)
	for _, link := range links {
		lines = append(lines, fmt.Sprintf("- %s %s %s -> %s %s", safeText(link.Relation, "-"), safeText(link.FromType, "-"), safeText(link.FromID, "-"), safeText(link.ToType, "-"), safeText(link.ToID, "-")))
		if strings.TrimSpace(link.Notes) != "" {
			lines = append(lines, "  "+link.Notes)
		}
	}
	return lines
}

func renderRevisionLines(revisions []cliRevision) []string {
	if len(revisions) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(revisions)*2)
	for _, revision := range revisions {
		lines = append(lines, fmt.Sprintf("- r%d title=%s author=%s", revision.Number, safeText(revision.Title, "-"), safeText(revision.Author, "-")))
		if strings.TrimSpace(revision.Summary) != "" {
			lines = append(lines, "  "+revision.Summary)
		}
	}
	return lines
}

func renderApprovalLines(approvals []cliApproval) []string {
	if len(approvals) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(approvals)*2)
	for _, approval := range approvals {
		lines = append(lines, fmt.Sprintf("- role=%s decision=%s actor=%s", safeText(approval.Role, "-"), safeText(approval.Decision, "-"), safeText(approval.Actor, "-")))
		if strings.TrimSpace(approval.Notes) != "" {
			lines = append(lines, "  "+approval.Notes)
		}
	}
	return lines
}

func renderEvidenceLines(evidence []cliEvidence) []string {
	if len(evidence) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(evidence)*3)
	for _, entry := range evidence {
		lines = append(lines, fmt.Sprintf("- %s", safeText(entry.Summary, "-")))
		if len(entry.Facts) > 0 {
			lines = append(lines, "  facts: "+renderFactPairs(entry.Facts))
		}
		if strings.TrimSpace(entry.AttachmentName) != "" {
			lines = append(lines, "  attachment: "+entry.AttachmentName)
		}
	}
	return lines
}

func renderFactPairs(facts map[string]string) string {
	if len(facts) == 0 {
		return "none"
	}
	pairs := make([]string, 0, len(facts))
	for key, value := range facts {
		pairs = append(pairs, fmt.Sprintf("%s=%s", safeText(key, "-"), safeText(value, "-")))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, ", ")
}

func joinOrNone(values []string) string {
	if len(values) == 0 {
		return "none"
	}
	return strings.Join(values, ", ")
}

func safeText(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

// Intent: Keep the CLI human-readable after peer-visible entity IDs switch to
// canonical envelope CIDs by showing the preserved short alias when one is
// available, while still using canonical IDs underneath for stable commands
// and lookups. Source: DI-loruk
func safeDisplayID(canonicalID string, aliasID string) string {
	if strings.TrimSpace(aliasID) != "" {
		return aliasID
	}
	return safeText(canonicalID, "-")
}

func readResponseBody(response *http.Response) ([]byte, error) {
	body, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if readErr != nil || closeErr != nil {
		return nil, errors.Join(readErr, closeErr)
	}
	return body, nil
}

func splitCSV(input string) []string {
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// Intent: Keep the CLI's typed local search operation using the same filter
// vocabulary as the HTTP adapter so terminal search semantics stay stable while
// the direct socket contract becomes more runtime-native. Source: DI-monuv
func buildSearchOptions(query string, filters []string) (service.SearchOptions, error) {
	options := service.SearchOptions{Query: query}
	for _, filter := range filters {
		key, value, ok := strings.Cut(filter, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return service.SearchOptions{}, fmt.Errorf("invalid search filter %q; expected key=value", filter)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if _, ok := allowedSearchFilters[key]; !ok {
			return service.SearchOptions{}, fmt.Errorf("unsupported search filter %q", key)
		}
		if value == "" {
			return service.SearchOptions{}, fmt.Errorf("search filter %q requires a value", key)
		}
		switch key {
		case "kind":
			options.Kind = value
		case "status":
			options.Status = value
		case "outcome":
			options.Outcome = value
		case "place_id":
			options.PlaceID = value
		case "resource_id":
			options.ResourceID = value
		case "responsibility_id":
			options.ResponsibilityID = value
		case "problem":
			options.Problem = strings.EqualFold(value, "true")
		}
	}
	return options, nil
}

func buildSearchPath(query string, filters []string) (string, error) {
	options, err := buildSearchOptions(query, filters)
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("q", options.Query)
	if options.Kind != "" {
		values.Set("kind", options.Kind)
	}
	if options.Status != "" {
		values.Set("status", options.Status)
	}
	if options.Outcome != "" {
		values.Set("outcome", options.Outcome)
	}
	if options.PlaceID != "" {
		values.Set("place_id", options.PlaceID)
	}
	if options.ResourceID != "" {
		values.Set("resource_id", options.ResourceID)
	}
	if options.ResponsibilityID != "" {
		values.Set("responsibility_id", options.ResponsibilityID)
	}
	if options.Problem {
		values.Set("problem", "true")
	}
	// Intent: Encode CLI search queries and structured filters before they hit
	// the HTTP adapter so spaces and reserved URL characters survive the
	// embodiment boundary without inventing a second search contract. Source:
	// DI-sifeg; DI-mifot
	return "/api/search?" + values.Encode(), nil
}

// Intent: Route direct CLI item inspection through the typed local runtime
// contract when the Unix socket is primary, while preserving the HTTP adapter
// as explicit compatibility transport. Source: DI-monuv
func (cli *CLI) itemBody(itemID string) ([]byte, error) {
	if strings.TrimSpace(cli.SocketPath) != "" {
		return cli.localSocketOperationBody(service.LocalEmbodimentRequest{
			Operation: "inspect_item",
			ItemID:    itemID,
		})
	}
	return cli.get("/api/items/" + itemID)
}

// Intent: Route direct CLI run inspection through the typed local runtime
// contract when the Unix socket is primary, while preserving the HTTP adapter
// as explicit compatibility transport. Source: DI-monuv
func (cli *CLI) runBody(runID string) ([]byte, error) {
	if strings.TrimSpace(cli.SocketPath) != "" {
		return cli.localSocketOperationBody(service.LocalEmbodimentRequest{
			Operation: "inspect_run",
			RunID:     runID,
		})
	}
	return cli.get("/api/runs/" + runID)
}

// Intent: Route direct CLI place/resource/responsibility inspection through
// the typed local runtime contract when the Unix socket is primary, while
// preserving the HTTP adapter as explicit compatibility transport. Source:
// DI-monuv
func (cli *CLI) entityBody(entityType string, entityID string) ([]byte, error) {
	if strings.TrimSpace(cli.SocketPath) != "" {
		return cli.localSocketOperationBody(service.LocalEmbodimentRequest{
			Operation:  "inspect_entity",
			EntityType: entityType,
			EntityID:   entityID,
		})
	}
	switch entityType {
	case "place":
		return cli.get("/api/places/" + entityID)
	case "resource":
		return cli.get("/api/resources/" + entityID)
	case "responsibility":
		return cli.get("/api/responsibilities/" + entityID)
	default:
		return nil, fmt.Errorf("unsupported entity type %q", entityType)
	}
}

// Intent: Route the CLI problem-review queue through the typed local runtime
// contract when the Unix socket is primary, while preserving the HTTP adapter
// as explicit compatibility transport. Source: DI-monuv
func (cli *CLI) problemReviewBody() ([]byte, error) {
	if strings.TrimSpace(cli.SocketPath) != "" {
		return cli.localSocketOperationBody(service.LocalEmbodimentRequest{
			Operation: "problem_review",
		})
	}
	return cli.get("/api/problem-review")
}
