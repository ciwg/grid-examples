package kernel

import (
	"archive/tar"
	"bytes"
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
)

// Intent: Workflow state enters from retained CAS bytes, so malformed input
// must be rejected without panicking or corrupting replay. Source: DI-lumek
func FuzzDecodeWorkflowLifecycleEvent(f *testing.F) {
	raw, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{
		State:         WorkflowImported,
		WorkflowAlias: "fuzz",
		ArtifactCID:   workflowHandoffProtocolCID,
	})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(raw)
	f.Fuzz(func(_ *testing.T, body []byte) {
		_, _ = DecodeWorkflowLifecycleEvent(body)
	})
}

// Intent: Run projections must reject arbitrary retained event bytes without
// trusting malformed transition state. Source: DI-lumek
func FuzzDecodeWorkflowRunEvent(f *testing.F) {
	raw, err := encodeWorkflowRunEvent(workflowRunEvent{
		RunCID:   workflowHandoffProtocolCID,
		Workflow: "fuzz",
		State:    WorkflowRunRunning,
		Input:    workflowHandoffProtocolCID,
	})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(raw)
	f.Fuzz(func(_ *testing.T, body []byte) {
		_, _ = decodeWorkflowRunEvent(body)
	})
}

// Intent: Artifact verification must safely reject arbitrary CAS bytes before
// they can become an executable workflow declaration. Source: DI-lumek
func FuzzWorkflowArtifactVerification(f *testing.F) {
	f.Add([]byte(`{"id":"fuzz"}`))
	canonicalArtifact, err := canonicalWorkflowFuzzArtifact()
	if err != nil {
		f.Fatal(err)
	}
	f.Add(canonicalArtifact)
	f.Fuzz(func(t *testing.T, artifact []byte) {
		cas, err := store.OpenCAS(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		artifactCID, err := cas.PutCID(artifact)
		if err != nil {
			t.Fatal(err)
		}
		runtime := &Runtime{cas: cas}
		_, _ = runtime.WorkflowManifest(artifactCID.String())
	})
}

// Intent: Full CAS scans must tolerate malformed mixtures of artifacts and
// lifecycle/run events without panicking during local projection rebuild.
// Source: DI-lumek
func FuzzWorkflowProjectionRebuild(f *testing.F) {
	f.Add([]byte("artifact"), []byte("lifecycle"), []byte("nonce"), []byte("handoff"), []byte("run"))
	artifact, lifecycle, nonce, handoff, run, err := workflowProjectionFuzzSeed(f)
	if err != nil {
		f.Fatal(err)
	}
	f.Add(artifact, lifecycle, nonce, handoff, run)
	f.Fuzz(func(t *testing.T, artifact, lifecycle, nonce, handoff, run []byte) {
		root := t.TempDir()
		cas, err := store.OpenCAS(root)
		if err != nil {
			t.Fatal(err)
		}
		for _, body := range [][]byte{artifact, lifecycle, nonce, handoff, run} {
			if _, err := cas.PutCID(body); err != nil {
				t.Fatal(err)
			}
		}
		stateRoot := filepath.Join(root, "state")
		_, _ = OpenWorkflowRegistry(stateRoot, cas)
		_, _ = OpenWorkflowRunRegistry(stateRoot, cas)
	})
}

func canonicalWorkflowFuzzArtifact() ([]byte, error) {
	files := []struct {
		name string
		body []byte
	}{
		{"schemas/input-v1.json", []byte(`{"fields":["inventory_id","run_id","place_id","counter","quantity","outcome","notes"],"kind":"moks.workflow.adapter.input","version":1,"workflow":"inventory-receipt"}`)},
		{"schemas/output-v1.json", []byte(`{"fields":["inventory_id","run_id","place_id","counter","quantity","outcome","notes","stage"],"kind":"moks.workflow.adapter.output","stage":"completed","version":1,"workflow":"inventory-receipt"}`)},
		{"workflow.json", []byte(`{"id":"inventory-receipt","version":"1.0.0","summary":"fuzz","required_packages":[],"required_protocols":[],"adapter":"inventory-receipt","input_pcid":"bafkreie3xn5cs7in24a5aenl7kpyaa22e346wr4tcqm4evxgcn2v55yvne","output_pcid":"bafkreibkoh3hdusvgscanho5rchq4esqjhd5kcbopnzzedhd7sgvime4ne","input_schema":"schemas/input-v1.json","output_schema":"schemas/output-v1.json"}`)},
	}
	var artifact bytes.Buffer
	writer := tar.NewWriter(&artifact)
	for _, file := range files {
		if err := writer.WriteHeader(&tar.Header{Name: file.name, Mode: 0o644, Size: int64(len(file.body)), Format: tar.FormatUSTAR}); err != nil {
			return nil, err
		}
		if _, err := writer.Write(file.body); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return artifact.Bytes(), nil
}

func workflowProjectionFuzzSeed(f *testing.F) ([]byte, []byte, []byte, []byte, []byte, error) {
	cas, err := store.OpenCAS(f.TempDir())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	artifact := []byte("accepted artifact")
	artifactCID, err := cas.PutCID(artifact)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	lifecycle, err := EncodeWorkflowLifecycleEvent(WorkflowLifecycleEvent{State: WorkflowImported, WorkflowAlias: "fuzz", ArtifactCID: artifactCID})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	nonce, err := records.EncodeGrid(records.GridEnvelope{ProtocolPCID: workflowRunProtocolCID, Slots: []any{"workflow-run-nonce", bytes.Repeat([]byte{1}, 32)}})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	runCID, err := cas.PutCID(nonce)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	handoff, err := EncodeWorkflowHandoff(WorkflowHandoff{PCID: WorkflowHandoffProtocolPCID, Values: map[string]string{"subject": "fuzz"}})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	inputCID, err := cas.PutCID(handoff)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	run, err := encodeWorkflowRunEvent(workflowRunEvent{RunCID: runCID, Workflow: "fuzz", State: WorkflowRunRunning, Input: inputCID})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return artifact, lifecycle, nonce, handoff, run, nil
}
