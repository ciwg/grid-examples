package packages

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"
)

func TestDockerWorkerCommandEnforcesIsolation(t *testing.T) {
	command, err := (DockerWorker{Image: "example/agent:1", Args: []string{"worker", "--cbor"}, CPUs: "0.5", Memory: "128m", PIDsLimit: 64, Timeout: time.Second}).Command(context.Background())
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
	if !slices.Equal(command.Args[len(command.Args)-2:], []string{"worker", "--cbor"}) {
		t.Fatalf("worker command was not appended after image: %#v", command.Args)
	}
}

func TestBoundedBufferRejectsOversizedWorkerStream(t *testing.T) {
	buffer := &boundedBuffer{limit: 3}
	if _, err := buffer.Write([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if _, err := buffer.Write([]byte("d")); !errors.Is(err, errDockerWorkerStreamLimit) {
		t.Fatalf("stream limit error = %v", err)
	}
	if buffer.String() != "abc" {
		t.Fatalf("buffer retained %q", buffer.String())
	}
}

func TestDockerWorkerCommandRejectsIncompleteLimits(t *testing.T) {
	_, err := (DockerWorker{Image: "example/agent:1", CPUs: "0.5", Memory: "128m", Timeout: time.Second}).Command(context.Background())
	if err == nil {
		t.Fatal("expected PID-limit validation failure")
	}
	_, err = (DockerWorker{Image: "example/agent:1", CPUs: "0.5", Memory: "128m", PIDsLimit: 64}).Command(context.Background())
	if err == nil {
		t.Fatal("expected timeout validation failure")
	}
}
