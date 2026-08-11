package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocols"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/service"
)

func main() {
	var (
		serverURL = flag.String("server", "http://127.0.0.1:7035", "bug tracker server URL")
		user      = flag.String("user", "engineer", "identity to use")
		agentKey  = flag.String("agent-key", "", "explicit path for this CLI's Ed25519 private key")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
	}
	client := &CLI{ServerURL: strings.TrimRight(*serverURL, "/"), User: *user, AgentKeyPath: *agentKey}
	var err error
	switch args[0] {
	case "assigned":
		err = client.Assigned()
	case "show":
		if len(args) != 2 {
			usage()
		}
		err = client.Show(args[1])
	case "comment":
		if len(args) < 3 {
			usage()
		}
		err = client.Comment(args[1], strings.Join(args[2:], " "))
	case "start":
		if len(args) != 2 {
			usage()
		}
		err = client.ChangeStatus(args[1], service.StatusInProgress)
	case "resolve":
		if len(args) != 2 {
			usage()
		}
		err = client.ChangeStatus(args[1], service.StatusResolved)
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: bug-tracker-cli [--server URL] [--user USER] assigned|show|comment|start|resolve ...")
	os.Exit(2)
}

type CLI struct {
	ServerURL    string
	User         string
	AgentKeyPath string
	HTTPClient   *http.Client
}

func (cli *CLI) Assigned() error {
	var payload struct {
		Issues []service.IssueSummary `json:"issues"`
	}
	if err := cli.get("/api/issues?assignee="+cli.User, &payload); err != nil {
		return err
	}
	for _, issue := range payload.Issues {
		fmt.Printf("%s  %-12s  %s\n", issue.ID, issue.Status, issue.Title)
	}
	return nil
}

func (cli *CLI) Show(issueID string) error {
	var issue service.Issue
	if err := cli.get("/api/issues/"+issueID, &issue); err != nil {
		return err
	}
	fmt.Printf("%s\n", issue.ID)
	fmt.Printf("Title: %s\n", issue.Title)
	fmt.Printf("Status: %s\n", issue.Status)
	fmt.Printf("Severity: %s\n", issue.Severity)
	fmt.Printf("Reporter: %s\n", issue.Reporter)
	fmt.Printf("Assignee: %s\n", blankFallback(issue.Assignee, "unassigned"))
	fmt.Printf("Updated: %s\n", issue.UpdatedAt)
	fmt.Printf("\n%s\n\n", issue.Description)
	for _, event := range issue.Timeline {
		fmt.Printf("[%s] %s %s\n", event.Timestamp, event.Actor, event.Type)
		if event.Comment != "" {
			fmt.Printf("  %s\n", event.Comment)
		}
		if event.Status != "" && event.Type == "status_changed" {
			fmt.Printf("  %s -> %s\n", event.PreviousStatus, event.Status)
		}
		if event.Type == "attachment_added" {
			fmt.Printf("  %s (%d bytes)\n", event.AttachmentName, event.AttachmentSize)
		}
	}
	return nil
}

func (cli *CLI) Comment(issueID string, comment string) error {
	return cli.submitLifecycle(issueID, "comment", comment, "")
}

func (cli *CLI) ChangeStatus(issueID string, status string) error {
	return cli.submitLifecycle(issueID, "status", "", status)
}

// Intent: Make CLI mutations client-signed pCID-selected promises instead of
// legacy unsigned HTTP writes. Source: DI-gonok; DI-muzal
func (cli *CLI) submitLifecycle(issueID, kind, comment, status string) error {
	key, err := identity.LoadOrCreateAgentKey(cli.AgentKeyPath)
	if err != nil {
		return err
	}
	enrollment := identity.NewEnrollment(key, cli.User)
	enrollmentProof, err := key.SignEnrollment(enrollment)
	if err != nil {
		return err
	}
	requestBytes, err := protocol.Marshal(struct {
		Enrollment identity.Enrollment      `cbor:"enrollment"`
		Proof      identity.EnrollmentProof `cbor:"proof"`
	}{enrollment, enrollmentProof})
	if err != nil {
		return err
	}
	if err := cli.postCBOR("/api/agents/enroll", requestBytes, nil); err != nil {
		return err
	}
	pcid, err := protocol.CIDForBytes(protocols.MustRead(protocols.IssueLifecycleUpdateSpec))
	if err != nil {
		return err
	}
	payload, err := protocol.Marshal(protocol.IssueLifecycleUpdate{AgentID: string(key.AgentID()), IssuedAt: "", IssueID: issueID, Kind: kind, Comment: comment, Status: status})
	if err != nil {
		return err
	}
	envelope := protocol.NewEnvelope(pcid, payload, protocol.Proof{Algorithm: "Ed25519", AgentID: string(key.AgentID()), PublicKey: key.PublicKey()})
	signable, err := envelope.SignableBytes()
	if err != nil {
		return err
	}
	envelope.Proof.Signature = key.Sign(signable)
	wire, err := envelope.Bytes()
	if err != nil {
		return err
	}
	return cli.postCBOR("/api/promises", wire, nil)
}

func (cli *CLI) postCBOR(path string, payloadBytes []byte, target any) error {
	request, err := http.NewRequest(http.MethodPost, cli.ServerURL+path, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/cbor")
	response, err := cli.httpClient().Do(request)
	if err != nil {
		return err
	}
	body, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if readErr != nil {
		return readErr
	}
	if closeErr != nil {
		return fmt.Errorf("close response body: %w", closeErr)
	}
	if response.StatusCode >= 300 {
		return fmt.Errorf("%s", strings.TrimSpace(string(body)))
	}
	if target != nil {
		return json.Unmarshal(body, target)
	}
	return nil
}

func (cli *CLI) get(path string, target any) error {
	request, err := http.NewRequest(http.MethodGet, cli.ServerURL+path, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Bug-User", cli.User)
	response, err := cli.httpClient().Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		if closeErr := response.Body.Close(); closeErr != nil {
			return fmt.Errorf("%s (close response body: %v)", strings.TrimSpace(string(body)), closeErr)
		}
		return fmt.Errorf("%s", strings.TrimSpace(string(body)))
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			return fmt.Errorf("decode response: %w (close response body: %v)", err, closeErr)
		}
		return err
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("close response body: %w", err)
	}
	return nil
}

func (cli *CLI) post(path string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, cli.ServerURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Bug-User", cli.User)
	response, err := cli.httpClient().Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 {
		message, _ := io.ReadAll(response.Body)
		if closeErr := response.Body.Close(); closeErr != nil {
			return fmt.Errorf("%s (close response body: %v)", strings.TrimSpace(string(message)), closeErr)
		}
		return fmt.Errorf("%s", strings.TrimSpace(string(message)))
	}
	if target == nil {
		if _, err := io.Copy(io.Discard, response.Body); err != nil {
			if closeErr := response.Body.Close(); closeErr != nil {
				return fmt.Errorf("discard response body: %w (close response body: %v)", err, closeErr)
			}
			return err
		}
		if err := response.Body.Close(); err != nil {
			return fmt.Errorf("close response body: %w", err)
		}
		return nil
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			return fmt.Errorf("decode response: %w (close response body: %v)", err, closeErr)
		}
		return err
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("close response body: %w", err)
	}
	return nil
}

func blankFallback(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func (cli *CLI) httpClient() *http.Client {
	if cli.HTTPClient != nil {
		return cli.HTTPClient
	}
	return http.DefaultClient
}
