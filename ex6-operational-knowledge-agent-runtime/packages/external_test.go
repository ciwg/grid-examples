package packages

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

func TestRunnerDescribeValidateAndRun(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "helper-agent.sh")
	fixtureGenerator, err := filepath.Abs(filepath.Join("..", "tools", "record-fixture"))
	if err != nil {
		t.Fatalf("resolve fixture generator: %v", err)
	}
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    cat <<'EOF'
{"id":"helper-agent","version":"0.1.0","description":"Test helper package","commands":[{"path":["helper","echo"],"summary":"Echo a string"}],"families":[{"name":"helper.echo.v1","protocol_pcid":"pcid:helper.echo.v1"}],"claims":[{"protocol_pcid":"pcid:helper.echo.v1","role":"family-validator","summary":"Validates helper echo envelopes."}]}
EOF
    ;;
  validate)
    family="$(go run "$FIXTURE_GENERATOR" inspect)"
    if [ "$family" != "helper.echo.v1" ]; then
      echo "wrong family" >&2
      exit 1
    fi
    ;;
  run)
    if [ "$2" != "helper echo" ]; then
      echo "unknown helper command" >&2
      exit 1
    fi
    printf '%s\n' "$3"
    ;;
  *)
    echo "unknown helper verb" >&2
    exit 1
    ;;
esac
`
	script = strings.Replace(script, "set -eu\n", "set -eu\nFIXTURE_GENERATOR='"+fixtureGenerator+"'\n", 1)
	script = strings.ReplaceAll(script, "pcid:helper.echo.v1", "bafkreihnjma5aoxdsxf2bj45q2hgtnsufyfb2gzguma2jaof6infm3ma6a")
	if err := os.WriteFile(executable, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper script: %v", err)
	}
	runner := Runner{Executable: executable}
	ctx := context.Background()
	manifest, err := runner.Describe(ctx)
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	if manifest.ID != "helper-agent" {
		t.Fatalf("unexpected id: %s", manifest.ID)
	}
	payload, err := records.CanonicalJSON([]byte(`{"message":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	raw := records.MustMarshal(records.Envelope{Family: "helper.echo.v1", ProtocolPCID: "bafkreihnjma5aoxdsxf2bj45q2hgtnsufyfb2gzguma2jaof6infm3ma6a", RecordID: "one", Signer: "helper", Timestamp: "2026-07-28T00:00:00Z", Payload: payload})
	if err := runner.ValidateEnvelope(ctx, raw); err != nil {
		t.Fatalf("validate: %v", err)
	}
	output, err := runner.RunCommand(ctx, "helper echo", []string{"hello"})
	if err != nil {
		t.Fatalf("run command: %v", err)
	}
	if output.Output != "hello" {
		t.Fatalf("unexpected output: %s", output.Output)
	}
}

func TestRunnerRunCommandStructuredResult(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "helper-agent.sh")
	fixtureGenerator, err := filepath.Abs(filepath.Join("..", "tools", "record-fixture"))
	if err != nil {
		t.Fatalf("resolve fixture generator: %v", err)
	}
	script := `#!/bin/sh
set -eu
case "$1" in
  describe)
    printf '{}\n'
    ;;
  validate)
    exit 0
    ;;
  run)
    record="$(go run "$FIXTURE_GENERATOR" --pcid bafkreihnjma5aoxdsxf2bj45q2hgtnsufyfb2gzguma2jaof6infm3ma6a helper.echo.v1 one helper 2026-07-28T00:00:00Z '{"body_ref":"$cas:body1"}')"
    printf '{"output":"created","cas":[{"alias":"body1","body":"hello body"}],"records":["%s"]}\n' "$record"
    ;;
  *)
    exit 1
    ;;
esac
`
	script = strings.Replace(script, "set -eu\n", "set -eu\nFIXTURE_GENERATOR='"+fixtureGenerator+"'\n", 1)
	if err := os.WriteFile(executable, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper script: %v", err)
	}
	runner := Runner{Executable: executable}
	result, err := runner.RunCommand(context.Background(), "helper create", nil)
	if err != nil {
		t.Fatalf("run structured command: %v", err)
	}
	if result.Output != "created" || len(result.CAS) != 1 || len(result.Records) != 1 {
		t.Fatalf("unexpected structured result: %#v", result)
	}
	record, err := records.Parse(result.Records[0])
	if err != nil {
		t.Fatalf("parse canonical record: %v", err)
	}
	if record.Family != "helper.echo.v1" {
		t.Fatalf("record family = %s", record.Family)
	}
}

func TestManifestValidateRequiresEmitsForParserClaims(t *testing.T) {
	manifest := Manifest{
		ID:      "parser-agent",
		Version: "0.1.0",
		Claims: []ImplementationClaim{
			{ProtocolPCID: "pcid:raw.example.v1", Role: "parser", RouteType: "parser"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected parser claim validation failure")
	}
}

func TestManifestValidateWorkflowAdapters(t *testing.T) {
	manifest := Manifest{
		ID:      "procedure-execution-adapter",
		Version: "0.1.0",
		WorkflowAdapters: []WorkflowAdapter{{
			Name:       "procedure-execution",
			Image:      "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Command:    []string{"worker"},
			InputPCID:  "bafkreiahdp34nto2rnnqde26jw3xnkd6xnlalnr72sug3w7tjb3bhhoj4q",
			OutputPCID: "bafkreifmttp5fwt3yvxvkb7ni6kwg3j3arl7mbjsyzszf7s7crxrncch24",
			CPUs:       "0.5",
			Memory:     "128m",
			PIDsLimit:  64,
			Timeout:    "30s",
		}},
	}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("validate workflow adapter: %v", err)
	}
	changed := manifest
	changed.WorkflowAdapters = append([]WorkflowAdapter{}, manifest.WorkflowAdapters...)
	changed.WorkflowAdapters[0].OutputPCID = "bafkreih6yllp2v7e5bmerznebzmohniezsv64hpqe2m33h6jclq6rfzqdu"
	if manifest.Equal(changed) {
		t.Fatal("manifest equality ignored workflow adapter output contract")
	}
	changed.WorkflowAdapters[0].PIDsLimit = 0
	if err := changed.Validate(); err == nil {
		t.Fatal("accepted workflow adapter without PID limit")
	}
	changed = manifest
	changed.WorkflowAdapters = append([]WorkflowAdapter{}, manifest.WorkflowAdapters...)
	changed.WorkflowAdapters[0].Image = "example/procedure-execution:1"
	if err := changed.Validate(); err == nil {
		t.Fatal("accepted workflow adapter with mutable Docker tag")
	}
	changed = manifest
	changed.WorkflowAdapters = append([]WorkflowAdapter{}, manifest.WorkflowAdapters...)
	changed.WorkflowAdapters[0].Image = "sha256:short"
	if err := changed.Validate(); err == nil {
		t.Fatal("accepted workflow adapter with incomplete Docker digest")
	}
}

func TestRegistryHostFromImage(t *testing.T) {
	digest := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	host, err := RegistryHostFromImage("REGISTRY.example:5000/moks/worker@sha256:" + digest)
	if err != nil || host != "registry.example:5000" {
		t.Fatalf("registry host = %q, %v", host, err)
	}
	for _, image := range []string{"sha256:" + digest, "registry.example/moks/worker:latest", "registry.example/moks/worker@sha256:short", "registry.example/@sha256:" + digest, "registry.example/moks/worker:tag@sha256:" + digest} {
		if _, err := RegistryHostFromImage(image); err == nil {
			t.Fatalf("accepted invalid portable image %q", image)
		}
	}
}

func TestPullImageRequiresExactRetainedDigest(t *testing.T) {
	digest := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	image := "registry.example/moks/worker@sha256:" + digest
	docker := filepath.Join(t.TempDir(), "docker")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = pull ]; then exit 0; fi\n" +
		"if [ \"$1\" = image ]; then printf '%s\\n' '[\"registry.example/moks/worker@sha256:" + digest + "\"]'; exit 0; fi\n" +
		"exit 64\n"
	if err := os.WriteFile(docker, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", filepath.Dir(docker)+":"+os.Getenv("PATH"))
	if err := PullImage(context.Background(), image); err != nil {
		t.Fatalf("pull exact image: %v", err)
	}
	script = "#!/bin/sh\nif [ \"$1\" = pull ]; then exit 0; fi\nif [ \"$1\" = image ]; then printf '%s\\n' '[]'; exit 0; fi\nexit 64\n"
	if err := os.WriteFile(docker, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := PullImage(context.Background(), image); err == nil {
		t.Fatal("accepted image without requested retained digest")
	}
}
