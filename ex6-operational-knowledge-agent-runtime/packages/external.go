package packages

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type CASWrite struct {
	Alias string `json:"alias"`
	Body  string `json:"body"`
}

type CommandResult struct {
	Output  string     `json:"output,omitempty"`
	CAS     []CASWrite `json:"cas,omitempty"`
	Records [][]byte   `json:"records,omitempty"`
}

type Runner struct {
	Executable string
}

// Intent: Require installed packages to describe themselves at runtime so the
// runtime can compare the declared agent shape against the executable before it
// is trusted. Source: DI-moksu
func (runner Runner) Describe(ctx context.Context) (Manifest, error) {
	command := exec.CommandContext(ctx, runner.Executable, "describe")
	output, err := command.Output()
	if err != nil {
		return Manifest{}, fmt.Errorf("describe %s: %w", runner.Executable, err)
	}
	var manifest Manifest
	if err := json.Unmarshal(output, &manifest); err != nil {
		return Manifest{}, err
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	manifest.Executable = runner.Executable
	return manifest, nil
}

func (runner Runner) ValidateEnvelope(ctx context.Context, raw []byte) error {
	command := exec.CommandContext(ctx, runner.Executable, "validate")
	command.Stdin = bytes.NewReader(raw)
	output, err := command.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if trimmed == "" {
			trimmed = err.Error()
		}
		return fmt.Errorf("validate %s: %s", runner.Executable, trimmed)
	}
	return nil
}

func (runner Runner) RunCommand(ctx context.Context, commandKey string, args []string) (CommandResult, error) {
	commandArgs := append([]string{"run", commandKey}, args...)
	command := exec.CommandContext(ctx, runner.Executable, commandArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		return CommandResult{}, fmt.Errorf("run %s: %s", runner.Executable, strings.TrimSpace(string(output)))
	}
	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return CommandResult{}, nil
	}
	var result CommandResult
	if err := json.Unmarshal(output, &result); err == nil {
		return result, nil
	}
	return CommandResult{Output: trimmed}, nil
}
