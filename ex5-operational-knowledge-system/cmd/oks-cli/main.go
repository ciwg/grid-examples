package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
	case "responsibilities":
		err = cli.Responsibilities()
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
		err = cli.Show("/api/items/" + args[1])
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
		err = cli.RecordRun(args[1], args[2], args[3], revision, args[5], strings.Join(args[6:], " "))
	case "show-run":
		if len(args) != 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.Show("/api/runs/" + args[1])
	case "approve-item":
		if len(args) < 6 {
			return 2, fmt.Errorf("%s", usageText)
		}
		revision, convErr := strconv.Atoi(args[2])
		if convErr != nil {
			return 1, convErr
		}
		err = cli.Approve("/api/items/"+args[1]+"/approvals", revision, args[3], args[4], strings.Join(args[5:], " "))
	case "approve-run":
		if len(args) < 5 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.Approve("/api/runs/"+args[1]+"/approvals", 0, args[2], args[3], strings.Join(args[4:], " "))
	case "search":
		if len(args) < 2 {
			return 2, fmt.Errorf("%s", usageText)
		}
		err = cli.Show("/api/search?q=" + args[1])
	default:
		return 2, fmt.Errorf("%s", usageText)
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}

const usageText = "usage: oks-cli dashboard|responsibilities|new-responsibility ACTOR TITLE SUMMARY|items [kind]|new-item ACTOR KIND TITLE SUMMARY BODY|show-item ID|runs [kind]|record-run ACTOR KIND ITEM_ID REVISION OUTCOME NOTES|show-run ID|approve-item ITEM_ID REVISION ROLE DECISION NOTES|approve-run RUN_ID ROLE DECISION NOTES|search QUERY"

type CLI struct {
	ServerURL string
}

func (cli *CLI) Dashboard() error        { return cli.Show("/api/dashboard") }
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

func (cli *CLI) NewItem(actor string, kind string, title string, summary string, body string) error {
	return cli.post("/api/items", map[string]any{"actor": actor, "kind": kind, "title": title, "summary": summary, "body": body})
}

func (cli *CLI) RecordRun(actor string, kind string, itemID string, revision int, outcome string, notes string) error {
	return cli.post("/api/runs", map[string]any{
		"actor":    actor,
		"kind":     kind,
		"item_id":  itemID,
		"revision": revision,
		"outcome":  outcome,
		"notes":    notes,
	})
}

func (cli *CLI) Approve(path string, revision int, role string, decision string, notes string) error {
	return cli.post(path, map[string]any{
		"actor":    "boss",
		"revision": revision,
		"role":     role,
		"decision": decision,
		"notes":    notes,
	})
}

func (cli *CLI) Show(path string) error {
	response, err := http.Get(cli.ServerURL + path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
		return fmt.Errorf("%s", strings.TrimSpace(string(body)))
	}
	var indented bytes.Buffer
	if err := json.Indent(&indented, body, "", "  "); err == nil {
		fmt.Println(indented.String())
		return nil
	}
	fmt.Println(string(body))
	return nil
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
	defer response.Body.Close()
	message, err := io.ReadAll(response.Body)
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
