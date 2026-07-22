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
)

func main() {
	serverURL := flag.String("server", "http://127.0.0.1:7045", "server URL")
	flag.Parse()
	cli := &CLI{ServerURL: strings.TrimRight(*serverURL, "/")}
	exitCode, err := cli.run(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
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

const usageText = "usage: oks-cli dashboard|problem-review|pending-review|places|new-place ACTOR KIND NAME SUMMARY [PARENT_ID]|show-place ID|resources|new-resource ACTOR KIND NAME SUMMARY [PLACE_ID]|show-resource ID|responsibilities|show-responsibility ID|new-responsibility ACTOR TITLE SUMMARY|items [kind]|new-item ACTOR KIND TITLE SUMMARY BODY|show-item ID|runs [kind]|record-run ACTOR KIND ITEM_ID REVISION OUTCOME NOTES [PLACE_ID] [RESOURCE_IDS_CSV]|show-run ID|approve-item ITEM_ID REVISION ACTOR ROLE DECISION NOTES|supersede-item ITEM_ID ACTOR [NOTES]|approve-run RUN_ID ACTOR ROLE DECISION NOTES|add-link ACTOR FROM_TYPE FROM_ID TO_TYPE TO_ID RELATION [NOTES]|add-evidence RUN_ID ACTOR SUMMARY [FACTS_JSON] [FILE]|search QUERY [kind=VALUE] [status=VALUE] [outcome=VALUE] [place_id=VALUE] [resource_id=VALUE] [responsibility_id=VALUE] [problem=true]"

type CLI struct {
	ServerURL string
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
	Kind              string          `json:"kind"`
	Status            string          `json:"status"`
	Title             string          `json:"title"`
	Summary           string          `json:"summary"`
	ResponsibilityIDs []string        `json:"responsibility_ids"`
	CurrentRevision   int             `json:"current_revision"`
	Revisions         []cliRevision   `json:"revisions"`
	RelatedRuns       []cliRunSummary `json:"related_runs"`
	Approvals         []cliApproval   `json:"approvals"`
	Links             []cliLink       `json:"links"`
}

type cliRunDetail struct {
	ID                string        `json:"id"`
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

var allowedSearchFilters = map[string]struct{}{
	"kind":              {},
	"status":            {},
	"outcome":           {},
	"place_id":          {},
	"resource_id":       {},
	"responsibility_id": {},
	"problem":           {},
}

func (cli *CLI) Dashboard() error     { return cli.Show("/api/dashboard") }
func (cli *CLI) ProblemReview() error { return cli.Show("/api/problem-review") }

// Intent: Keep place drilldowns on the shared place-detail projection while
// rendering the hierarchy, links, and related runs in a terminal-first review
// layout that is easier to act on than a raw JSON dump. Source: DI-luzom
func (cli *CLI) ShowPlace(placeID string) error {
	body, err := cli.get("/api/places/" + placeID)
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
	body, err := cli.get("/api/resources/" + resourceID)
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
	body, err := cli.get("/api/responsibilities/" + responsibilityID)
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
	body, err := cli.get("/api/items/" + itemID)
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

// Intent: Keep run drilldowns on the shared run-detail projection while
// rendering evidence, approvals, linked context, and follow-on handoff hints
// in a terminal-first review layout instead of a raw JSON dump. Source:
// DI-salup
func (cli *CLI) ShowRun(runID string) error {
	body, err := cli.get("/api/runs/" + runID)
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

// Intent: Keep the shell-first pending-review queue on the same projection
// family Neovim already uses so terminal triage does not invent a separate
// review endpoint or drift away from the staged editor workflow. Source:
// DI-vabok
func (cli *CLI) PendingReview() error {
	draftItems, err := cli.getSearchArray("/api/search?status=draft", "items")
	if err != nil {
		return err
	}
	allRuns, err := cli.getSearchArray("/api/search", "runs")
	if err != nil {
		return err
	}
	problemRuns, err := cli.getSearchArray("/api/search?problem=true", "runs")
	if err != nil {
		return err
	}
	unreviewedRuns := []map[string]any{}
	for _, raw := range allRuns {
		run, err := requireJSONObject(raw, "/api/search runs entry")
		if err != nil {
			return err
		}
		approvals, err := requireJSONArray(run, "approvals", "/api/search runs entry")
		if err != nil {
			return err
		}
		if len(approvals) == 0 {
			unreviewedRuns = append(unreviewedRuns, run)
		}
	}
	return printJSON(map[string]any{
		"draft_items":     draftItems,
		"unreviewed_runs": unreviewedRuns,
		"problem_runs":    problemRuns,
	})
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

	request, err := http.NewRequest(http.MethodPost, cli.ServerURL+"/api/runs/"+runID+"/evidence", &body)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	message, err := readResponseBody(response)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
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
	response, err := http.Post(cli.ServerURL+path, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	message, err := readResponseBody(response)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
		return fmt.Errorf("%s", strings.TrimSpace(string(message)))
	}
	return printJSONBytes(message)
}

func (cli *CLI) get(path string) ([]byte, error) {
	response, err := http.Get(cli.ServerURL + path)
	if err != nil {
		return nil, err
	}
	body, err := readResponseBody(response)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 300 {
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
		fmt.Sprintf("# Responsibility %s", safeText(responsibility.ID, "-")),
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
		fmt.Sprintf("# Item %s", safeText(item.ID, "-")),
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
		fmt.Sprintf("# Run %s", safeText(run.ID, "-")),
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

func renderRunLines(runs []cliRunSummary) []string {
	if len(runs) == 0 {
		return []string{"- none"}
	}
	lines := make([]string, 0, len(runs)*3)
	for _, run := range runs {
		lines = append(lines, fmt.Sprintf("- %s kind=%s outcome=%s item=%s", safeText(run.ID, "-"), safeText(run.Kind, "-"), safeText(run.Outcome, "-"), safeText(run.ItemID, "-")))
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

func buildSearchPath(query string, filters []string) (string, error) {
	values := url.Values{}
	values.Set("q", query)
	for _, filter := range filters {
		key, value, ok := strings.Cut(filter, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return "", fmt.Errorf("invalid search filter %q; expected key=value", filter)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if _, ok := allowedSearchFilters[key]; !ok {
			return "", fmt.Errorf("unsupported search filter %q", key)
		}
		if value == "" {
			return "", fmt.Errorf("search filter %q requires a value", key)
		}
		values.Set(key, value)
	}
	// Intent: Encode CLI search queries and structured filters before they hit
	// the HTTP adapter so spaces and reserved URL characters survive the
	// embodiment boundary without inventing a second search contract. Source:
	// DI-sifeg; DI-mifot
	return "/api/search?" + values.Encode(), nil
}
