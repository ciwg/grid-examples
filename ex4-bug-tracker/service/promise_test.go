package service_test

import (
	"encoding/json"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocol"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/protocols"
	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/service"
)

func TestSubmitPromiseAcceptsEnrolledSignedReport(t *testing.T) {
	app, err := service.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	enrollment := identity.NewEnrollment(key, "reporter")
	proof, err := key.SignEnrollment(enrollment)
	if err != nil {
		t.Fatalf("sign enrollment: %v", err)
	}
	if err := app.EnrollAgent(enrollment, proof); err != nil {
		t.Fatalf("enroll: %v", err)
	}
	pcid, err := protocol.CIDForBytes(protocols.MustRead(protocols.IssueReportSpec))
	if err != nil {
		t.Fatalf("pCID: %v", err)
	}
	payload, err := protocol.Marshal(protocol.IssueReport{AgentID: string(key.AgentID()), IssuedAt: "2026-08-10T00:00:00Z", Team: "CORE", Title: "Signed issue", Description: "Signed report body", Severity: "High"})
	if err != nil {
		t.Fatalf("payload: %v", err)
	}
	envelope := protocol.NewEnvelope(pcid, payload, protocol.Proof{Algorithm: "Ed25519", AgentID: string(key.AgentID()), PublicKey: key.PublicKey()})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("signable: %v", err)
	}
	envelope.Proof.Signature = key.Sign(signable)
	wire, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("wire: %v", err)
	}
	issue, err := app.SubmitPromise(wire)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if issue.ID != "BUG-0001" {
		t.Fatalf("issue ID: got %s", issue.ID)
	}
}

func TestPrepareFinalizeAndSubmitPromise(t *testing.T) {
	app, err := service.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	enrollment := identity.NewEnrollment(key, service.RoleReporter)
	enrollmentProof, err := key.SignEnrollment(enrollment)
	if err != nil {
		t.Fatalf("sign enrollment: %v", err)
	}
	if err := app.EnrollAgent(enrollment, enrollmentProof); err != nil {
		t.Fatalf("enroll: %v", err)
	}
	payload, err := json.Marshal(map[string]string{"title": "Bridge", "description": "Prepared by service", "severity": service.SeverityHigh})
	if err != nil {
		t.Fatalf("marshal draft: %v", err)
	}
	prepared, err := app.PreparePromise(service.PromiseDraft{Profile: "issue-report", AgentID: string(key.AgentID()), Payload: payload})
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	wire, err := app.FinalizePromise(service.PromiseProof{DraftID: prepared.DraftID, PublicKey: key.PublicKey(), Signature: key.Sign(prepared.SignableBytes)})
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	envelope, err := protocol.ParseEnvelope(wire)
	if err != nil {
		t.Fatalf("parse finalized envelope: %v", err)
	}
	if envelope.PCID.String() == "" {
		t.Fatal("finalized envelope omitted pCID")
	}
	issue, err := app.SubmitPromise(wire)
	if err != nil {
		t.Fatalf("submit finalized envelope: %v", err)
	}
	if issue.Title != "Bridge" {
		t.Fatalf("title: got %q", issue.Title)
	}
}

func TestSubmitPromiseRejectsUnenrolledSigner(t *testing.T) {
	app, err := service.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	pcid, err := protocol.CIDForBytes(protocols.MustRead(protocols.IssueReportSpec))
	if err != nil {
		t.Fatalf("pCID: %v", err)
	}
	payload, err := protocol.Marshal(protocol.IssueReport{AgentID: string(key.AgentID()), IssuedAt: "2026-08-10T00:00:00Z", Team: "CORE", Title: "Rejected", Description: "Unenrolled signer", Severity: "Low"})
	if err != nil {
		t.Fatalf("payload: %v", err)
	}
	envelope := protocol.NewEnvelope(pcid, payload, protocol.Proof{AgentID: string(key.AgentID()), PublicKey: key.PublicKey()})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("signable: %v", err)
	}
	envelope.Proof.Signature = key.Sign(signable)
	wire, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("wire: %v", err)
	}
	if _, err := app.SubmitPromise(wire); err == nil {
		t.Fatal("expected unenrolled signer rejection")
	}
}

func TestSubmitPromiseAcceptsSignedAttachmentReference(t *testing.T) {
	app, err := service.NewApp(t.TempDir())
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	seed, err := app.CreateIssue(service.RoleReporter, "Attachment", "Attachment promise test", service.SeverityLow)
	if err != nil {
		t.Fatalf("seed issue: %v", err)
	}
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	enrollment := identity.NewEnrollment(key, service.RoleReporter)
	enrollmentProof, err := key.SignEnrollment(enrollment)
	if err != nil {
		t.Fatalf("sign enrollment: %v", err)
	}
	if err := app.EnrollAgent(enrollment, enrollmentProof); err != nil {
		t.Fatalf("enroll: %v", err)
	}
	attachmentCID, err := app.StoreAttachmentObject([]byte("stack trace"))
	if err != nil {
		t.Fatalf("store object: %v", err)
	}
	pcid, err := protocol.CIDForBytes(protocols.MustRead(protocols.IssueAttachmentReferenceSpec))
	if err != nil {
		t.Fatalf("pCID: %v", err)
	}
	payload, err := protocol.Marshal(protocol.IssueAttachmentReference{AgentID: string(key.AgentID()), IssuedAt: "2026-08-10T00:00:00Z", IssueID: seed.ID, AttachmentCID: attachmentCID, Name: "trace.log", ContentType: "text/plain", Size: int64(len("stack trace"))})
	if err != nil {
		t.Fatalf("payload: %v", err)
	}
	envelope := protocol.NewEnvelope(pcid, payload, protocol.Proof{Algorithm: "Ed25519", AgentID: string(key.AgentID()), PublicKey: key.PublicKey()})
	signable, err := envelope.SignableBytes()
	if err != nil {
		t.Fatalf("signable: %v", err)
	}
	envelope.Proof.Signature = key.Sign(signable)
	wire, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("wire: %v", err)
	}
	if _, err := app.SubmitPromise(wire); err != nil {
		t.Fatalf("submit: %v", err)
	}
	attachment, err := app.DownloadAttachment(seed.ID, "ATT-000002")
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	if string(attachment.Bytes) != "stack trace" {
		t.Fatalf("attachment bytes = %q", attachment.Bytes)
	}
}
