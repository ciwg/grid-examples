package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/cas"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocols"
	promiseStore "github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/store"
	"github.com/ipfs/go-cid"
)

type App struct {
	dataRoot        string
	store           *Store
	mu              sync.Mutex
	issues          map[string]*Issue
	nextIssueNumber int
	nextSequence    uint64
	promiseCAS      *cas.Store
	artifactStore   *promiseStore.ArtifactStore
	enrollments     map[string]identity.Enrollment
	reportPCID      string
	updatePCID      string
	attachmentPCID  string
	prepared        map[string]preparedPromise
}

type preparedPromise struct {
	pcid         string
	payloadBytes []byte
	agentID      string
	expiresAt    time.Time
}

// PromiseDraft is local HTTPS adapter input, not a PromiseGrid artifact.
type PromiseDraft struct {
	Profile string          `json:"profile"`
	AgentID string          `json:"agent_id"`
	Payload json.RawMessage `json:"payload"`
}

// PromiseProof supplies the browser or CLI proof for a previously prepared draft.
type PromiseProof struct {
	DraftID   string `json:"draft_id"`
	PublicKey []byte `json:"public_key"`
	Signature []byte `json:"signature"`
}

type PreparedPromise struct {
	DraftID       string `json:"draft_id"`
	SignableBytes []byte `json:"signable_bytes"`
}

func NewApp(dataRoot string) (*App, error) {
	store, events, err := OpenStore(dataRoot)
	if err != nil {
		return nil, err
	}
	app := &App{
		dataRoot:    dataRoot,
		store:       store,
		issues:      map[string]*Issue{},
		enrollments: map[string]identity.Enrollment{},
		prepared:    map[string]preparedPromise{},
	}
	promiseCAS, err := cas.Open(filepath.Join(dataRoot, "cas"))
	if err != nil {
		return nil, err
	}
	artifactStore, err := promiseStore.Open(dataRoot)
	if err != nil {
		return nil, err
	}
	app.promiseCAS, app.artifactStore = promiseCAS, artifactStore
	records, err := artifactStore.AgentBindings.Load()
	if err != nil {
		return nil, err
	}
	for _, enrollment := range records {
		app.enrollments[string(enrollment.AgentID)] = enrollment
	}
	app.reportPCID, err = profileCID(protocols.IssueReportSpec)
	if err != nil {
		return nil, err
	}
	app.updatePCID, err = profileCID(protocols.IssueLifecycleUpdateSpec)
	if err != nil {
		return nil, err
	}
	app.attachmentPCID, err = profileCID(protocols.IssueAttachmentReferenceSpec)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := app.applyEventLocked(event); err != nil {
			return nil, fmt.Errorf("replay event %d: %w", event.Sequence, err)
		}
		if event.Sequence > app.nextSequence {
			app.nextSequence = event.Sequence
		}
	}
	return app, nil
}

func profileCID(name string) (string, error) {
	value, err := protocol.CIDForBytes(protocols.MustRead(name))
	if err != nil {
		return "", err
	}
	return value.String(), nil
}

// EnrollAgent records only a verified public local-admission binding.
func (app *App) EnrollAgent(enrollment identity.Enrollment, proof identity.EnrollmentProof) error {
	if err := identity.VerifyEnrollment(enrollment, proof); err != nil {
		return err
	}
	if _, err := validateIdentity(enrollment.Role); err != nil {
		return err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if existing, ok := app.enrollments[string(enrollment.AgentID)]; ok && !bytes.Equal(existing.PublicKey, enrollment.PublicKey) {
		return fmt.Errorf("agent %s already enrolled with another key", enrollment.AgentID)
	}
	if err := app.artifactStore.AgentBindings.Append(enrollment); err != nil {
		return err
	}
	app.enrollments[string(enrollment.AgentID)] = enrollment
	return nil
}

// PrepareEnrollment creates canonical proof bytes for a browser-owned key;
// enrollment remains an explicit service-scoped admission step. Source: DI-muzal
func (app *App) PrepareEnrollment(publicKey []byte, role string) (identity.Enrollment, []byte, error) {
	if _, err := validateIdentity(role); err != nil {
		return identity.Enrollment{}, nil, err
	}
	enrollment := identity.Enrollment{AgentID: identity.AgentIDForPublicKey(publicKey), PublicKey: append([]byte(nil), publicKey...), Role: role}
	claim, err := protocol.Marshal(enrollment)
	if err != nil {
		return identity.Enrollment{}, nil, fmt.Errorf("marshal enrollment claim: %w", err)
	}
	return enrollment, claim, nil
}

// Intent: Keep JSON limited to bounded adapter carriage while the service alone
// derives the pCID and canonical tag-42 signing bytes. Source: DI-kolaf
func (app *App) PreparePromise(draft PromiseDraft) (PreparedPromise, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	enrollment, enrolled := app.enrollments[draft.AgentID]
	if !enrolled {
		return PreparedPromise{}, fmt.Errorf("unaccepted signer")
	}
	pcid, payload, err := app.preparePayloadLocked(draft, enrollment.Role)
	if err != nil {
		return PreparedPromise{}, err
	}
	signable, err := protocol.NewEnvelope(pcid, payload, protocol.Proof{}).SignableBytes()
	if err != nil {
		return PreparedPromise{}, err
	}
	draftID, err := newDraftID()
	if err != nil {
		return PreparedPromise{}, err
	}
	app.prepared[draftID] = preparedPromise{pcid: pcid.String(), payloadBytes: payload, agentID: draft.AgentID, expiresAt: time.Now().Add(5 * time.Minute)}
	return PreparedPromise{DraftID: draftID, SignableBytes: signable}, nil
}

// Intent: Verify client-held proof before assembling exact envelope bytes; no
// service private key signs for an embodiment. Source: DI-pusip; DI-kolaf
func (app *App) FinalizePromise(proof PromiseProof) ([]byte, error) {
	app.mu.Lock()
	draft, found := app.prepared[proof.DraftID]
	if found {
		delete(app.prepared, proof.DraftID)
	}
	enrollment, enrolled := app.enrollments[draft.agentID]
	app.mu.Unlock()
	if !found || time.Now().After(draft.expiresAt) {
		return nil, fmt.Errorf("unknown or expired promise draft")
	}
	if !enrolled || !bytes.Equal(enrollment.PublicKey, proof.PublicKey) {
		return nil, fmt.Errorf("unaccepted signer")
	}
	pcid, err := cidForString(draft.pcid)
	if err != nil {
		return nil, err
	}
	envelope := protocol.NewEnvelope(pcid, draft.payloadBytes, protocol.Proof{Algorithm: "Ed25519", AgentID: draft.agentID, PublicKey: proof.PublicKey, Signature: proof.Signature})
	signable, err := envelope.SignableBytes()
	if err != nil {
		return nil, err
	}
	if !identity.Verify(enrollment.PublicKey, signable, proof.Signature) {
		return nil, fmt.Errorf("invalid promise proof")
	}
	return envelope.Bytes()
}

func newDraftID() (string, error) {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("random draft ID: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func cidForString(value string) (cid.Cid, error) {
	parsed, err := cid.Decode(value)
	if err != nil {
		return cid.Undef, fmt.Errorf("decode pCID: %w", err)
	}
	return parsed, nil
}

func (app *App) preparePayloadLocked(draft PromiseDraft, role string) (cid.Cid, []byte, error) {
	if draft.AgentID == "" {
		return cid.Undef, nil, fmt.Errorf("missing agent ID")
	}
	switch draft.Profile {
	case "issue-report":
		if role != RoleReporter {
			return cid.Undef, nil, fmt.Errorf("%s cannot create issues", draft.AgentID)
		}
		var payload protocol.IssueReport
		if err := json.Unmarshal(draft.Payload, &payload); err != nil {
			return cid.Undef, nil, fmt.Errorf("decode report draft: %w", err)
		}
		if err := validateTitle(payload.Title); err != nil {
			return cid.Undef, nil, err
		}
		if err := validateDescription(payload.Description); err != nil {
			return cid.Undef, nil, err
		}
		if err := validateSeverity(payload.Severity); err != nil {
			return cid.Undef, nil, err
		}
		payload.AgentID, payload.IssuedAt, payload.Team = draft.AgentID, time.Now().UTC().Format(time.RFC3339), defaultTeam
		bytes, err := protocol.Marshal(payload)
		if err != nil {
			return cid.Undef, nil, err
		}
		pcid, err := cidForString(app.reportPCID)
		return pcid, bytes, err
	case "issue-lifecycle-update":
		var payload protocol.IssueLifecycleUpdate
		if err := json.Unmarshal(draft.Payload, &payload); err != nil {
			return cid.Undef, nil, fmt.Errorf("decode lifecycle draft: %w", err)
		}
		if err := validateIssueID(payload.IssueID); err != nil {
			return cid.Undef, nil, err
		}
		switch payload.Kind {
		case "comment":
			if err := validateComment(payload.Comment); err != nil {
				return cid.Undef, nil, err
			}
		case "assignment":
			if role != RoleTriage {
				return cid.Undef, nil, fmt.Errorf("%s cannot assign issues", draft.AgentID)
			}
		case "status":
			if err := validateStatus(payload.Status); err != nil {
				return cid.Undef, nil, err
			}
		default:
			return cid.Undef, nil, fmt.Errorf("unsupported update kind %q", payload.Kind)
		}
		payload.AgentID, payload.IssuedAt = draft.AgentID, time.Now().UTC().Format(time.RFC3339)
		bytes, err := protocol.Marshal(payload)
		if err != nil {
			return cid.Undef, nil, err
		}
		pcid, err := cidForString(app.updatePCID)
		return pcid, bytes, err
	case "issue-attachment-reference":
		var payload protocol.IssueAttachmentReference
		if err := json.Unmarshal(draft.Payload, &payload); err != nil {
			return cid.Undef, nil, fmt.Errorf("decode attachment draft: %w", err)
		}
		if err := validateIssueID(payload.IssueID); err != nil {
			return cid.Undef, nil, err
		}
		if payload.AttachmentCID == "" || payload.Name == "" || payload.ContentType == "" || payload.Size <= 0 {
			return cid.Undef, nil, fmt.Errorf("incomplete attachment reference")
		}
		if _, err := app.promiseCAS.Get(payload.AttachmentCID); err != nil {
			return cid.Undef, nil, fmt.Errorf("attachment object unavailable: %w", err)
		}
		payload.AgentID, payload.IssuedAt = draft.AgentID, time.Now().UTC().Format(time.RFC3339)
		bytes, err := protocol.Marshal(payload)
		if err != nil {
			return cid.Undef, nil, err
		}
		pcid, err := cidForString(app.attachmentPCID)
		return pcid, bytes, err
	default:
		return cid.Undef, nil, fmt.Errorf("unsupported local profile %q", draft.Profile)
	}
}

// Intent: Accept only a signer-owned, locally enrolled, pCID-selected promise
// before updating Ex4's existing local workflow projection. Source: DI-gonok; DI-muzal
func (app *App) SubmitPromise(envelopeBytes []byte) (issue Issue, err error) {
	artifactCID := ""
	if observedCID, cidErr := protocol.CIDForBytes(envelopeBytes); cidErr == nil {
		artifactCID = observedCID.String()
	}
	defer func() {
		if err == nil {
			return
		}
		// Intent: Preserve bounded local evidence of rejected ingress without
		// treating a rejected artifact as an accepted promise. Source: DI-ninul
		if observationErr := app.artifactStore.Observations.Append(promiseStore.Observation{Reason: err.Error(), ArtifactCID: artifactCID, ObservedAt: time.Now().UTC().Format(time.RFC3339)}); observationErr != nil {
			err = fmt.Errorf("%w (record rejection observation: %v)", err, observationErr)
		}
	}()
	envelope, err := protocol.ParseEnvelope(envelopeBytes)
	if err != nil {
		return Issue{}, err
	}
	signable, err := envelope.SignableBytes()
	if err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	enrollment, ok := app.enrollments[envelope.Proof.AgentID]
	app.mu.Unlock()
	if !ok || !bytes.Equal(enrollment.PublicKey, envelope.Proof.PublicKey) || !identity.Verify(enrollment.PublicKey, signable, envelope.Proof.Signature) {
		return Issue{}, fmt.Errorf("unaccepted signer")
	}
	if envelope.Proof.AgentID == "" {
		return Issue{}, fmt.Errorf("missing agent ID")
	}
	artifactCID, err = app.promiseCAS.Put(envelopeBytes)
	if err != nil {
		return Issue{}, err
	}
	switch envelope.PCID.String() {
	case app.reportPCID:
		var payload protocol.IssueReport
		if err := protocol.Unmarshal(envelope.PayloadBytes, &payload); err != nil {
			return Issue{}, err
		}
		if payload.AgentID != envelope.Proof.AgentID {
			return Issue{}, fmt.Errorf("report agent mismatch")
		}
		issue, err = app.CreateIssue(enrollment.Role, payload.Title, payload.Description, payload.Severity)
	case app.updatePCID:
		var payload protocol.IssueLifecycleUpdate
		if err := protocol.Unmarshal(envelope.PayloadBytes, &payload); err != nil {
			return Issue{}, err
		}
		if payload.AgentID != envelope.Proof.AgentID {
			return Issue{}, fmt.Errorf("update agent mismatch")
		}
		switch payload.Kind {
		case "comment":
			issue, err = app.AddComment(enrollment.Role, payload.IssueID, payload.Comment)
		case "assignment":
			assignee := ""
			if payload.AssigneeAgentID != "" {
				var found bool
				assignee, found = app.enrollmentRole(payload.AssigneeAgentID)
				if !found {
					return Issue{}, fmt.Errorf("unknown assignee agent")
				}
			}
			issue, err = app.AssignIssue(enrollment.Role, payload.IssueID, assignee)
		case "status":
			issue, err = app.ChangeStatus(enrollment.Role, payload.IssueID, payload.Status)
		default:
			return Issue{}, fmt.Errorf("unsupported update kind %q", payload.Kind)
		}
	case app.attachmentPCID:
		var payload protocol.IssueAttachmentReference
		if err := protocol.Unmarshal(envelope.PayloadBytes, &payload); err != nil {
			return Issue{}, err
		}
		if payload.AgentID != envelope.Proof.AgentID {
			return Issue{}, fmt.Errorf("attachment agent mismatch")
		}
		issue, err = app.RecordAttachmentReference(enrollment.Role, payload)
	default:
		return Issue{}, fmt.Errorf("unsupported pCID %s", envelope.PCID)
	}
	if err != nil {
		return Issue{}, err
	}
	if err := app.artifactStore.AcceptedPromises.Append(promiseStore.AcceptedPromise{ArtifactCID: artifactCID, PCID: envelope.PCID.String(), AgentID: envelope.Proof.AgentID, AcceptedAt: time.Now().UTC().Format(time.RFC3339)}); err != nil {
		return Issue{}, err
	}
	return issue, nil
}

func (app *App) enrollmentRole(agentID string) (string, bool) {
	app.mu.Lock()
	defer app.mu.Unlock()
	enrollment, ok := app.enrollments[agentID]
	return enrollment.Role, ok
}

func (app *App) Meta() Meta {
	return Meta{
		DataRoot:   app.dataRoot,
		Statuses:   append([]string(nil), statuses...),
		Severities: append([]string(nil), severities...),
		Identities: append([]Identity(nil), identities...),
		Team:       defaultTeam,
		Profiles:   map[string]string{"issue-report": app.reportPCID, "issue-lifecycle-update": app.updatePCID, "issue-attachment-reference": app.attachmentPCID},
	}
}

func (app *App) ListIssues(status string, assignee string) ([]IssueSummary, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if status != "" {
		if err := validateStatus(status); err != nil {
			return nil, err
		}
	}
	if err := validateAssignee(assignee); err != nil {
		return nil, err
	}
	result := []IssueSummary{}
	for _, issue := range app.issues {
		if status != "" && issue.Status != status {
			continue
		}
		if assignee != "" && issue.Assignee != assignee {
			continue
		}
		result = append(result, issue.Summary())
	}
	sort.Slice(result, func(i int, j int) bool {
		if result[i].UpdatedAt == result[j].UpdatedAt {
			return result[i].ID > result[j].ID
		}
		return result[i].UpdatedAt > result[j].UpdatedAt
	})
	return result, nil
}

func (app *App) GetIssue(issueID string) (Issue, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	issue, ok := app.issues[issueID]
	if !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	return cloneIssue(issue), nil
}

// Intent: Keep the first bug-tracker slice durable-first by persisting issue
// creation as an append-only event and projecting queue/detail state from that
// history instead of mutating a canonical row in place. Source: DI-nunit
func (app *App) CreateIssue(actor string, title string, description string, severity string) (Issue, error) {
	identity, err := validateIdentity(actor)
	if err != nil {
		return Issue{}, err
	}
	if identity.Role != RoleReporter {
		return Issue{}, fmt.Errorf("%s cannot create issues", actor)
	}
	if err := validateTitle(title); err != nil {
		return Issue{}, err
	}
	if err := validateDescription(description); err != nil {
		return Issue{}, err
	}
	if err := validateSeverity(severity); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	issueID := fmt.Sprintf("BUG-%04d", app.nextIssueNumber+1)
	event := IssueEvent{
		IssueID:        issueID,
		Type:           "created",
		Actor:          actor,
		Title:          strings.TrimSpace(title),
		Description:    normalizeLineEndings(description),
		Severity:       severity,
		Status:         StatusNew,
		Team:           defaultTeam,
		PreviousStatus: "",
	}
	if _, err := app.appendEventLocked(event); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func (app *App) AddComment(actor string, issueID string, comment string) (Issue, error) {
	if _, err := validateIdentity(actor); err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if err := validateComment(comment); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.issues[issueID]; !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	if _, err := app.appendEventLocked(IssueEvent{
		IssueID: issueID,
		Type:    "commented",
		Actor:   actor,
		Comment: normalizeLineEndings(comment),
	}); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func (app *App) AssignIssue(actor string, issueID string, assignee string) (Issue, error) {
	identity, err := validateIdentity(actor)
	if err != nil {
		return Issue{}, err
	}
	if identity.Role != RoleTriage {
		return Issue{}, fmt.Errorf("%s cannot assign issues", actor)
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if err := validateAssignee(assignee); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	issue, ok := app.issues[issueID]
	if !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	// Intent: Keep assignment meaningful in the fixed bug-tracker workflow by
	// only allowing triage to assign issues that are ready for or already in
	// active work, instead of attaching owners to untouched or already resolved
	// issues. Source: DI-gitam
	if issue.Status == StatusNew || issue.Status == StatusResolved {
		return Issue{}, fmt.Errorf("issue %q cannot be assigned while %s", issueID, issue.Status)
	}
	if _, err := app.appendEventLocked(IssueEvent{
		IssueID:          issueID,
		Type:             "assigned",
		Actor:            actor,
		Assignee:         assignee,
		PreviousAssignee: issue.Assignee,
	}); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

// Intent: Keep the workflow legible by allowing only the locked v1 transitions
// and recording reopen as a normal history event that clears active ownership
// instead of mutating away the earlier resolved state. Source: DI-ninuf;
// DI-gofub
func (app *App) ChangeStatus(actor string, issueID string, status string) (Issue, error) {
	identity, err := validateIdentity(actor)
	if err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if err := validateStatus(status); err != nil {
		return Issue{}, err
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	issue, ok := app.issues[issueID]
	if !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	if issue.Status == status {
		return Issue{}, fmt.Errorf("issue %q is already %s", issueID, status)
	}
	if err := validateTransition(identity.Role, issue.Status, status); err != nil {
		return Issue{}, err
	}
	// Intent: Treat the assignee as the active workflow owner for engineer work
	// so only the assigned engineer can start or resolve implementation work.
	// Source: DI-gitam
	if identity.Role == RoleEngineer && (status == StatusInProgress || status == StatusResolved) {
		if issue.Assignee == "" {
			return Issue{}, fmt.Errorf("issue %q is unassigned", issueID)
		}
		if issue.Assignee != actor {
			return Issue{}, fmt.Errorf("issue %q is assigned to %s", issueID, issue.Assignee)
		}
	}
	event := IssueEvent{
		IssueID:        issueID,
		Type:           "status_changed",
		Actor:          actor,
		Status:         status,
		PreviousStatus: issue.Status,
		Assignee:       issue.Assignee,
	}
	if issue.Status == StatusResolved && status == StatusTriaged {
		event.PreviousAssignee = issue.Assignee
		event.Assignee = ""
	}
	if _, err := app.appendEventLocked(event); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

func validateTransition(role string, from string, to string) error {
	switch {
	case role == RoleTriage && from == StatusNew && to == StatusTriaged:
		return nil
	case role == RoleEngineer && from == StatusTriaged && to == StatusInProgress:
		return nil
	case role == RoleEngineer && from == StatusInProgress && to == StatusResolved:
		return nil
	case (role == RoleReporter || role == RoleTriage) && from == StatusResolved && to == StatusTriaged:
		return nil
	default:
		return fmt.Errorf("status change %q -> %q is not allowed for %s", from, to, role)
	}
}

func (app *App) AddAttachment(actor string, issueID string, filename string, contentType string, bytes []byte) (Issue, error) {
	if _, err := validateIdentity(actor); err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(issueID); err != nil {
		return Issue{}, err
	}
	if len(bytes) == 0 {
		return Issue{}, fmt.Errorf("attachment is empty")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.issues[issueID]; !ok {
		return Issue{}, fmt.Errorf("issue %q not found", issueID)
	}
	nextSequence := app.nextSequence + 1
	attachmentID := fmt.Sprintf("ATT-%06d", nextSequence)
	safeName := sanitizeAttachmentName(filename)
	relativePath := filepath.Join("attachments", issueID, fmt.Sprintf("%06d-%s", nextSequence, safeName))
	if err := app.store.WriteAttachment(relativePath, bytes); err != nil {
		return Issue{}, err
	}
	event := IssueEvent{
		IssueID:               issueID,
		Type:                  "attachment_added",
		Actor:                 actor,
		AttachmentID:          attachmentID,
		AttachmentName:        safeName,
		AttachmentPath:        relativePath,
		AttachmentContentType: contentType,
		AttachmentSize:        int64(len(bytes)),
	}
	if _, err := app.appendEventLocked(event); err != nil {
		if removeErr := app.store.RemoveAttachment(relativePath); removeErr != nil {
			return Issue{}, fmt.Errorf("%v (attachment cleanup failed: %v)", err, removeErr)
		}
		return Issue{}, err
	}
	return cloneIssue(app.issues[issueID]), nil
}

// StoreAttachmentObject accepts opaque attachment bytes into local CAS; it
// does not change an issue until a signed reference promise is accepted.
// Intent: Separate untrusted object carriage from the promise that gives an
// object issue-specific meaning. Source: DI-ninul; DI-kolaf
func (app *App) StoreAttachmentObject(bytes []byte) (string, error) {
	if len(bytes) == 0 || len(bytes) > maxAttachmentBytes {
		return "", fmt.Errorf("attachment size is out of bounds")
	}
	return app.promiseCAS.Put(bytes)
}

// RecordAttachmentReference projects an accepted attachment-reference promise.
func (app *App) RecordAttachmentReference(role string, payload protocol.IssueAttachmentReference) (Issue, error) {
	if _, err := validateIdentity(role); err != nil {
		return Issue{}, err
	}
	if err := validateIssueID(payload.IssueID); err != nil {
		return Issue{}, err
	}
	if payload.AttachmentCID == "" || payload.Name == "" || payload.ContentType == "" || payload.Size <= 0 {
		return Issue{}, fmt.Errorf("incomplete attachment reference")
	}
	bytes, err := app.promiseCAS.Get(payload.AttachmentCID)
	if err != nil {
		return Issue{}, fmt.Errorf("attachment object unavailable: %w", err)
	}
	if int64(len(bytes)) != payload.Size {
		return Issue{}, fmt.Errorf("attachment size does not match object")
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	if _, ok := app.issues[payload.IssueID]; !ok {
		return Issue{}, fmt.Errorf("issue %q not found", payload.IssueID)
	}
	attachmentID := fmt.Sprintf("ATT-%06d", app.nextSequence+1)
	event := IssueEvent{IssueID: payload.IssueID, Type: "attachment_added", Actor: role, AttachmentID: attachmentID, AttachmentName: sanitizeAttachmentName(payload.Name), AttachmentPath: "cas:" + payload.AttachmentCID, AttachmentContentType: payload.ContentType, AttachmentSize: payload.Size}
	if _, err := app.appendEventLocked(event); err != nil {
		return Issue{}, err
	}
	return cloneIssue(app.issues[payload.IssueID]), nil
}

func (app *App) DownloadAttachment(issueID string, attachmentID string) (AttachmentDownload, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	if err := validateIssueID(issueID); err != nil {
		return AttachmentDownload{}, err
	}
	issue, ok := app.issues[issueID]
	if !ok {
		return AttachmentDownload{}, fmt.Errorf("issue %q not found", issueID)
	}
	for _, event := range issue.Timeline {
		if event.Type == "attachment_added" && event.AttachmentID == attachmentID {
			var bytes []byte
			var err error
			if strings.HasPrefix(event.AttachmentPath, "cas:") {
				bytes, err = app.promiseCAS.Get(strings.TrimPrefix(event.AttachmentPath, "cas:"))
			} else {
				bytes, err = app.store.ReadAttachment(event.AttachmentPath)
			}
			if err != nil {
				return AttachmentDownload{}, err
			}
			return AttachmentDownload{
				Name:        event.AttachmentName,
				ContentType: event.AttachmentContentType,
				Bytes:       bytes,
			}, nil
		}
	}
	return AttachmentDownload{}, fmt.Errorf("attachment %q not found", attachmentID)
}

func (app *App) appendEventLocked(event IssueEvent) (IssueEvent, error) {
	app.nextSequence++
	event.Sequence = app.nextSequence
	if event.Timestamp == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if err := app.store.Append(event); err != nil {
		app.nextSequence--
		return IssueEvent{}, err
	}
	if err := app.applyEventLocked(event); err != nil {
		return IssueEvent{}, err
	}
	return event, nil
}

func (app *App) applyEventLocked(event IssueEvent) error {
	if event.Sequence > app.nextSequence {
		app.nextSequence = event.Sequence
	}
	if number := issueNumber(event.IssueID); number > app.nextIssueNumber {
		app.nextIssueNumber = number
	}
	switch event.Type {
	case "created":
		issue := &Issue{
			ID:          event.IssueID,
			Title:       event.Title,
			Description: event.Description,
			Severity:    event.Severity,
			Status:      event.Status,
			Reporter:    event.Actor,
			Assignee:    "",
			Team:        defaultTeam,
			CreatedAt:   event.Timestamp,
			UpdatedAt:   event.Timestamp,
			Timeline:    []IssueEvent{event},
		}
		if event.Team != "" {
			issue.Team = event.Team
		}
		app.issues[event.IssueID] = issue
		return nil
	default:
		issue, ok := app.issues[event.IssueID]
		if !ok {
			return fmt.Errorf("issue %q not found for event %q", event.IssueID, event.Type)
		}
		switch event.Type {
		case "commented":
		case "assigned":
			issue.Assignee = event.Assignee
		case "status_changed":
			issue.Status = event.Status
			issue.Assignee = event.Assignee
		case "attachment_added":
		default:
			return fmt.Errorf("unknown event type %q", event.Type)
		}
		issue.UpdatedAt = event.Timestamp
		issue.Timeline = append(issue.Timeline, event)
		return nil
	}
}

func issueNumber(issueID string) int {
	if !strings.HasPrefix(issueID, "BUG-") {
		return 0
	}
	number, err := strconv.Atoi(strings.TrimPrefix(issueID, "BUG-"))
	if err != nil {
		return 0
	}
	return number
}
