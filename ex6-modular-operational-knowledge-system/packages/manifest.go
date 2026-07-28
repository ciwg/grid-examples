package packages

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Command struct {
	Path    []string `json:"path"`
	Summary string   `json:"summary"`
}

func (command Command) Key() string {
	return strings.Join(command.Path, " ")
}

type Family struct {
	Name         string `json:"name"`
	ProtocolPCID string `json:"protocol_pcid"`
}

type Manifest struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	Description string    `json:"description,omitempty"`
	Executable  string    `json:"executable,omitempty"`
	Commands    []Command `json:"commands,omitempty"`
	Families    []Family  `json:"families,omitempty"`
}

func LoadManifest(path string) (Manifest, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return Manifest{}, err
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	if manifest.Executable != "" && !filepath.IsAbs(manifest.Executable) {
		manifest.Executable = filepath.Join(filepath.Dir(path), manifest.Executable)
	}
	return manifest, nil
}

func (manifest Manifest) Validate() error {
	if strings.TrimSpace(manifest.ID) == "" {
		return errors.New("package id is required")
	}
	if strings.TrimSpace(manifest.Version) == "" {
		return errors.New("package version is required")
	}
	commandKeys := map[string]struct{}{}
	for _, command := range manifest.Commands {
		if len(command.Path) == 0 {
			return errors.New("command path is required")
		}
		for _, segment := range command.Path {
			if strings.TrimSpace(segment) == "" {
				return fmt.Errorf("command path contains empty segment for package %s", manifest.ID)
			}
		}
		key := command.Key()
		if _, exists := commandKeys[key]; exists {
			return fmt.Errorf("duplicate command %q", key)
		}
		commandKeys[key] = struct{}{}
	}
	familyNames := map[string]struct{}{}
	for _, family := range manifest.Families {
		if strings.TrimSpace(family.Name) == "" {
			return fmt.Errorf("family name is required for package %s", manifest.ID)
		}
		if strings.TrimSpace(family.ProtocolPCID) == "" {
			return fmt.Errorf("protocol_pcid is required for family %s", family.Name)
		}
		if _, exists := familyNames[family.Name]; exists {
			return fmt.Errorf("duplicate family %q", family.Name)
		}
		familyNames[family.Name] = struct{}{}
	}
	return nil
}

func (manifest Manifest) Equal(other Manifest) bool {
	left := manifest
	right := other
	sortManifest(&left)
	sortManifest(&right)
	return left.ID == right.ID &&
		left.Version == right.Version &&
		left.Description == right.Description &&
		left.Executable == right.Executable &&
		equalCommands(left.Commands, right.Commands) &&
		equalFamilies(left.Families, right.Families)
}

func sortManifest(manifest *Manifest) {
	slices.SortFunc(manifest.Commands, func(left, right Command) int {
		return strings.Compare(left.Key(), right.Key())
	})
	slices.SortFunc(manifest.Families, func(left, right Family) int {
		return strings.Compare(left.Name, right.Name)
	})
}

func equalCommands(left []Command, right []Command) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index].Summary != right[index].Summary || !slices.Equal(left[index].Path, right[index].Path) {
			return false
		}
	}
	return true
}

func equalFamilies(left []Family, right []Family) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
