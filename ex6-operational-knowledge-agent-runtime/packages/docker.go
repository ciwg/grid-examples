package packages

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type DockerWorker struct {
	Image     string
	CPUs      string
	Memory    string
	PIDsLimit int
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
	return exec.CommandContext(ctx, "docker", arguments...), nil
}
