package packages

import (
	"context"
	"slices"
	"testing"
)

func TestDockerWorkerCommandEnforcesIsolation(t *testing.T) {
	command, err := (DockerWorker{Image: "example/agent:1", CPUs: "0.5", Memory: "128m", PIDsLimit: 64}).Command(context.Background())
	if err != nil {
		t.Fatalf("build command: %v", err)
	}
	for _, required := range []string{"--read-only", "--network", "none", "--cap-drop", "ALL", "--security-opt", "no-new-privileges", "--pids-limit", "64", "--cpus", "0.5", "--memory", "128m"} {
		if !slices.Contains(command.Args, required) {
			t.Fatalf("missing required sandbox argument %q in %#v", required, command.Args)
		}
	}
	for _, argument := range command.Args {
		if argument == "--volume" || argument == "--mount" {
			t.Fatalf("unexpected host-mount argument in %#v", command.Args)
		}
	}
}

func TestDockerWorkerCommandRejectsIncompleteLimits(t *testing.T) {
	_, err := (DockerWorker{Image: "example/agent:1", CPUs: "0.5", Memory: "128m"}).Command(context.Background())
	if err == nil {
		t.Fatal("expected PID-limit validation failure")
	}
}
