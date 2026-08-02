package packages

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const maxDockerWorkerStreamBytes = 16 << 20

var errDockerWorkerStreamLimit = errors.New("docker worker stream exceeds 16 MiB")

type boundedBuffer struct {
	bytes.Buffer
	limit int
}

func (buffer *boundedBuffer) Write(body []byte) (int, error) {
	remaining := buffer.limit - buffer.Len()
	if remaining <= 0 || len(body) > remaining {
		if remaining > 0 {
			if _, err := buffer.Buffer.Write(body[:remaining]); err != nil {
				return 0, err
			}
		}
		return len(body), errDockerWorkerStreamLimit
	}
	return buffer.Buffer.Write(body)
}

type DockerWorker struct {
	Image     string
	Args      []string
	CPUs      string
	Memory    string
	PIDsLimit int
	Timeout   time.Duration
}

// Intent: Run independently authored agents behind Docker-enforced limits so a
// route selection never grants ambient host authority. Source: DI-dovek
func (worker DockerWorker) Command(ctx context.Context) (*exec.Cmd, error) {
	if strings.TrimSpace(worker.Image) == "" {
		return nil, errors.New("docker worker image is required")
	}
	if strings.TrimSpace(worker.CPUs) == "" {
		return nil, errors.New("docker worker CPU limit is required")
	}
	if strings.TrimSpace(worker.Memory) == "" {
		return nil, errors.New("docker worker memory limit is required")
	}
	if worker.PIDsLimit < 1 {
		return nil, errors.New("docker worker PID limit must be positive")
	}
	if worker.Timeout <= 0 {
		return nil, errors.New("docker worker timeout must be positive")
	}
	arguments := []string{
		"run", "--rm", "--interactive",
		"--read-only",
		"--network", "none",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"--pids-limit", strconv.Itoa(worker.PIDsLimit),
		"--cpus", worker.CPUs,
		"--memory", worker.Memory,
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=64m",
		"--env", "PATH=/usr/bin:/bin",
		worker.Image,
	}
	arguments = append(arguments, worker.Args...)
	return exec.CommandContext(ctx, "docker", arguments...), nil
}

// Run gives a confined worker exact CBOR bytes on stdin and returns only its
// stdout bytes. The runtime remains responsible for decoding, validating, and
// persisting the worker's proposal. Source: DI-fofuh
func (worker DockerWorker) Run(ctx context.Context, input []byte) ([]byte, error) {
	deadline, cancel := context.WithTimeout(ctx, worker.Timeout)
	defer cancel()
	command, err := worker.Command(deadline)
	if err != nil {
		return nil, err
	}
	command.Stdin = bytes.NewReader(input)
	stdout := &boundedBuffer{limit: maxDockerWorkerStreamBytes}
	stderr := &boundedBuffer{limit: maxDockerWorkerStreamBytes}
	command.Stdout = stdout
	command.Stderr = stderr
	err = command.Run()
	if err != nil {
		if errors.Is(err, errDockerWorkerStreamLimit) {
			return nil, errDockerWorkerStreamLimit
		}
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("docker workflow worker: %s", detail)
	}
	return stdout.Bytes(), nil
}
