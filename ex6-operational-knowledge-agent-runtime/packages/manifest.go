package packages

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
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

type ImplementationClaim struct {
	ProtocolPCID   string   `json:"protocol_pcid"`
	Role           string   `json:"role"`
	RouteType      string   `json:"route_type,omitempty"`
	EmitsProtocols []string `json:"emits_protocols,omitempty"`
	Summary        string   `json:"summary,omitempty"`
}

// WorkflowAdapter declares an executable workflow contract supplied by an
// installed package. The runtime, rather than the package, owns invocation and
// durable writes after validating this declaration during activation.
// Intent: Keep executable workflow authority inside the package manifest that
// is self-checked at installation, while retaining Docker as the only worker
// boundary. Source: DI-fofuh
type WorkflowAdapter struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Command    []string `json:"command,omitempty"`
	InputPCID  string   `json:"input_pcid"`
	OutputPCID string   `json:"output_pcid"`
	CPUs       string   `json:"cpus"`
	Memory     string   `json:"memory"`
	PIDsLimit  int      `json:"pids_limit"`
	Timeout    string   `json:"timeout"`
}

type Manifest struct {
	ID               string                `json:"id"`
	Version          string                `json:"version"`
	Description      string                `json:"description,omitempty"`
	Executable       string                `json:"executable,omitempty"`
	Commands         []Command             `json:"commands,omitempty"`
	Families         []Family              `json:"families,omitempty"`
	Claims           []ImplementationClaim `json:"claims,omitempty"`
	WorkflowAdapters []WorkflowAdapter     `json:"workflow_adapters,omitempty"`
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
	claimedProtocols := map[string]struct{}{}
	for _, claim := range manifest.Claims {
		if strings.TrimSpace(claim.ProtocolPCID) == "" {
			return fmt.Errorf("claim protocol_pcid is required for package %s", manifest.ID)
		}
		if strings.TrimSpace(claim.Role) == "" {
			return fmt.Errorf("claim role is required for package %s", manifest.ID)
		}
		// Intent: Keep route-shape meaning inside package claims so parser and
		// transform hops use the same declaration surface as direct handlers.
		// Source: DI-lafek
		switch claim.NormalizedRouteType() {
		case "direct":
			if len(claim.EmitsProtocols) > 0 {
				return fmt.Errorf("direct claim %s for package %s must not declare emits_protocols", claim.ProtocolPCID, manifest.ID)
			}
		case "parser", "transform":
			if len(claim.EmitsProtocols) == 0 {
				return fmt.Errorf("%s claim %s for package %s must declare emits_protocols", claim.NormalizedRouteType(), claim.ProtocolPCID, manifest.ID)
			}
			seenEmits := map[string]struct{}{}
			for _, emit := range claim.EmitsProtocols {
				if strings.TrimSpace(emit) == "" {
					return fmt.Errorf("%s claim %s for package %s has empty emits_protocols entry", claim.NormalizedRouteType(), claim.ProtocolPCID, manifest.ID)
				}
				if _, exists := seenEmits[emit]; exists {
					return fmt.Errorf("%s claim %s for package %s has duplicate emits_protocols entry %s", claim.NormalizedRouteType(), claim.ProtocolPCID, manifest.ID, emit)
				}
				seenEmits[emit] = struct{}{}
			}
		default:
			return fmt.Errorf("unsupported route_type %q for package %s", claim.RouteType, manifest.ID)
		}
		key := claim.ProtocolPCID + "\x00" + claim.Role
		if _, exists := claimedProtocols[key]; exists {
			return fmt.Errorf("duplicate claim %q for package %s", claim.ProtocolPCID, manifest.ID)
		}
		claimedProtocols[key] = struct{}{}
	}
	claimedProtocolSet := map[string]struct{}{}
	for _, claim := range manifest.Claims {
		claimedProtocolSet[claim.ProtocolPCID] = struct{}{}
	}
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
		if _, exists := claimedProtocolSet[family.ProtocolPCID]; !exists {
			return fmt.Errorf("family %s protocol_pcid %s is not covered by a claim", family.Name, family.ProtocolPCID)
		}
		familyNames[family.Name] = struct{}{}
	}
	adapterNames := map[string]struct{}{}
	for _, adapter := range manifest.WorkflowAdapters {
		if strings.TrimSpace(adapter.Name) == "" || strings.TrimSpace(adapter.Image) == "" || strings.TrimSpace(adapter.InputPCID) == "" || strings.TrimSpace(adapter.OutputPCID) == "" {
			return fmt.Errorf("workflow adapter name, image, input_pcid, and output_pcid are required for package %s", manifest.ID)
		}
		if strings.TrimSpace(adapter.CPUs) == "" || strings.TrimSpace(adapter.Memory) == "" || adapter.PIDsLimit < 1 {
			return fmt.Errorf("workflow adapter %s requires CPU, memory, and positive PID limits", adapter.Name)
		}
		timeout, err := time.ParseDuration(adapter.Timeout)
		if err != nil || timeout <= 0 {
			return fmt.Errorf("workflow adapter %s requires a positive timeout", adapter.Name)
		}
		for _, argument := range adapter.Command {
			if strings.TrimSpace(argument) == "" {
				return fmt.Errorf("workflow adapter %s command contains an empty argument", adapter.Name)
			}
		}
		if _, exists := adapterNames[adapter.Name]; exists {
			return fmt.Errorf("duplicate workflow adapter %q for package %s", adapter.Name, manifest.ID)
		}
		adapterNames[adapter.Name] = struct{}{}
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
		equalFamilies(left.Families, right.Families) &&
		equalClaims(left.Claims, right.Claims) &&
		equalWorkflowAdapters(left.WorkflowAdapters, right.WorkflowAdapters)
}

func (manifest Manifest) HasClaim(protocolPCID string, role string) bool {
	for _, claim := range manifest.Claims {
		if claim.ProtocolPCID == protocolPCID && claim.Role == role {
			return true
		}
	}
	return false
}

func (manifest Manifest) FamiliesForProtocol(protocolPCID string) []string {
	families := []string{}
	for _, family := range manifest.Families {
		if family.ProtocolPCID == protocolPCID {
			families = append(families, family.Name)
		}
	}
	slices.Sort(families)
	return families
}

func (claim ImplementationClaim) NormalizedRouteType() string {
	routeType := strings.TrimSpace(claim.RouteType)
	if routeType == "" {
		return "direct"
	}
	return routeType
}

func (claim ImplementationClaim) SortedEmitsProtocols() []string {
	emits := append([]string{}, claim.EmitsProtocols...)
	slices.Sort(emits)
	return emits
}

func sortManifest(manifest *Manifest) {
	slices.SortFunc(manifest.Commands, func(left, right Command) int {
		return strings.Compare(left.Key(), right.Key())
	})
	slices.SortFunc(manifest.Families, func(left, right Family) int {
		return strings.Compare(left.Name, right.Name)
	})
	slices.SortFunc(manifest.Claims, func(left, right ImplementationClaim) int {
		leftKey := left.ProtocolPCID + "\x00" + left.Role
		rightKey := right.ProtocolPCID + "\x00" + right.Role
		if diff := strings.Compare(leftKey, rightKey); diff != 0 {
			return diff
		}
		return strings.Compare(left.NormalizedRouteType(), right.NormalizedRouteType())
	})
	slices.SortFunc(manifest.WorkflowAdapters, func(left, right WorkflowAdapter) int {
		return strings.Compare(left.Name, right.Name)
	})
	for index := range manifest.Claims {
		manifest.Claims[index].RouteType = manifest.Claims[index].NormalizedRouteType()
		manifest.Claims[index].EmitsProtocols = manifest.Claims[index].SortedEmitsProtocols()
	}
}

func equalWorkflowAdapters(left []WorkflowAdapter, right []WorkflowAdapter) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index].Name != right[index].Name || left[index].Image != right[index].Image || left[index].InputPCID != right[index].InputPCID || left[index].OutputPCID != right[index].OutputPCID || left[index].CPUs != right[index].CPUs || left[index].Memory != right[index].Memory || left[index].PIDsLimit != right[index].PIDsLimit || left[index].Timeout != right[index].Timeout || !slices.Equal(left[index].Command, right[index].Command) {
			return false
		}
	}
	return true
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

func equalClaims(left []ImplementationClaim, right []ImplementationClaim) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index].ProtocolPCID != right[index].ProtocolPCID ||
			left[index].Role != right[index].Role ||
			left[index].Summary != right[index].Summary ||
			left[index].NormalizedRouteType() != right[index].NormalizedRouteType() ||
			!slices.Equal(left[index].SortedEmitsProtocols(), right[index].SortedEmitsProtocols()) {
			return false
		}
	}
	return true
}
