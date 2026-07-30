package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/builtin"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/grid"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/kernel"
	contextpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/context"
	inventorypkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/inventory"
	knowledgepkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/knowledge"
	linkspkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/links"
	maintenancepkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/maintenance"
	procedurespkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/procedures"
	receivingpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/receiving"
	runspkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/runs"
	trainingpkg "github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages/training"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func workflowPrint(workflow kernel.Workflow) error {
	return workflowsPrint([]kernel.Workflow{workflow})
}

func workflowsPrint(workflows []kernel.Workflow) error {
	output, err := json.MarshalIndent(workflows, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func run(ctx context.Context, args []string) error {
	root, err := defaultRuntimeRoot()
	if err != nil {
		return err
	}
	runtime, err := kernel.Open(root)
	if err != nil {
		return err
	}
	defer func() {
		_ = runtime.Close()
	}()
	if err := registerBuiltins(runtime); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("command is required")
	}
	switch {
	case matchesPrefix(args, "route", "list"):
		if len(args) != 2 {
			return errors.New("usage: route list")
		}
		return routeList(runtime)
	case matchesPrefix(args, "route", "scope", "list"):
		if len(args) != 3 {
			return errors.New("usage: route scope list")
		}
		return routeScopeList(runtime)
	case matchesPrefix(args, "route", "scope", "inspect"):
		if len(args) != 4 && (len(args) < 6 || len(args[4:])%2 != 0) {
			return errors.New("usage: route scope inspect <name> [sort <depth|label|summary>] [depth <n|n+>] [label <text>] [summary <text>]")
		}
		query := kernel.RouteScopeGroupQuery{}
		for index := 4; index < len(args); index += 2 {
			switch args[index] {
			case "sort":
				query.SortBy = args[index+1]
			case "depth":
				query.DepthFilter = args[index+1]
			case "label":
				query.LabelFilter = args[index+1]
			case "summary":
				query.SummaryFilter = args[index+1]
			default:
				return errors.New("usage: route scope inspect <name> [sort <depth|label|summary>] [depth <n|n+>] [label <text>] [summary <text>]")
			}
		}
		return routeScopeInspect(runtime, args[3], query)
	case matchesPrefix(args, "route", "scope", "set"):
		if len(args) < 6 || len(args[4:])%2 != 0 {
			return errors.New("usage: route scope set <name> <kind> <target> [<kind> <target> ...]")
		}
		return routeScopeSet(runtime, args[3], args[4:])
	case matchesPrefix(args, "route", "scope", "remove"):
		if len(args) != 4 {
			return errors.New("usage: route scope remove <name>")
		}
		return routeScopeRemove(runtime, args[3])
	case matchesPrefix(args, "route", "policy", "show"):
		if len(args) != 3 && len(args) != 4 && len(args) != 5 {
			return errors.New("usage: route policy show [<protocol-pcid> [<role>]]")
		}
		protocolPCID := ""
		role := ""
		if len(args) == 4 {
			protocolPCID = args[3]
		}
		if len(args) == 5 {
			protocolPCID = args[3]
			role = args[4]
		}
		return routePolicyShow(runtime, protocolPCID, role)
	case matchesPrefix(args, "route", "policy", "set"):
		if len(args) != 7 {
			return errors.New("usage: route policy set <prefer-route-types|-> <avoid-route-types|-> <prefer-roles|-> <avoid-roles|->")
		}
		return routePolicySet(runtime, args[3:])
	case matchesPrefix(args, "route", "policy", "set-for"):
		if len(args) != 8 {
			return errors.New("usage: route policy set-for <protocol-pcid> <prefer-route-types|-> <avoid-route-types|-> <prefer-roles|-> <avoid-roles|->")
		}
		return routePolicySetForProtocol(runtime, args[3], args[4:])
	case matchesPrefix(args, "route", "policy", "remove"):
		if len(args) != 4 {
			return errors.New("usage: route policy remove <protocol-pcid>")
		}
		return routePolicyRemove(runtime, args[3])
	case matchesPrefix(args, "route", "policy", "set-for-role"):
		if len(args) != 9 {
			return errors.New("usage: route policy set-for-role <protocol-pcid> <role> <prefer-route-types|-> <avoid-route-types|-> <prefer-roles|-> <avoid-roles|->")
		}
		return routePolicySetForRole(runtime, args[3], args[4], args[5:])
	case matchesPrefix(args, "route", "policy", "remove-role"):
		if len(args) != 5 {
			return errors.New("usage: route policy remove-role <protocol-pcid> <role>")
		}
		return routePolicyRemoveRole(runtime, args[3], args[4])
	case matchesPrefix(args, "route", "plan"):
		if len(args) != 3 && len(args) != 4 && len(args) < 6 {
			return errors.New("usage: route plan <protocol-pcid> [trace [candidate <package-id:role:route-type>|downstream <protocol-pcid>|depth <n|n+>|scope <preset-or-alias>] ...]")
		}
		trace := false
		filter := kernel.RoutePlanTraceFilter{}
		if len(args) == 4 {
			if args[3] != "trace" {
				return errors.New("usage: route plan <protocol-pcid> [trace [candidate <package-id:role:route-type>|downstream <protocol-pcid>|depth <n|n+>|scope <preset-or-alias>] ...]")
			}
			trace = true
		}
		if len(args) >= 6 {
			if args[3] != "trace" {
				return errors.New("usage: route plan <protocol-pcid> [trace [candidate <package-id:role:route-type>|downstream <protocol-pcid>|depth <n|n+>|scope <preset-or-alias>] ...]")
			}
			if len(args[4:])%2 != 0 {
				return errors.New("usage: route plan <protocol-pcid> [trace [candidate <package-id:role:route-type>|downstream <protocol-pcid>|depth <n|n+>|scope <preset-or-alias>] ...]")
			}
			trace = true
			clauses := []kernel.RoutePlanTraceFilterClause{}
			for index := 4; index < len(args); index += 2 {
				if args[index] != "candidate" && args[index] != "downstream" && args[index] != "depth" && args[index] != "scope" {
					return errors.New("usage: route plan <protocol-pcid> [trace [candidate <package-id:role:route-type>|downstream <protocol-pcid>|depth <n|n+>|scope <preset-or-alias>] ...]")
				}
				clauses = append(clauses, kernel.RoutePlanTraceFilterClause{
					Kind:   args[index],
					Target: args[index+1],
				})
			}
			filter = kernel.RoutePlanTraceFilter{Clauses: clauses}
		}
		return routePlan(runtime, args[2], trace, filter)
	case matchesPrefix(args, "route", "inspect"):
		if len(args) != 3 {
			return errors.New("usage: route inspect <protocol-pcid>")
		}
		return routeInspect(runtime, args[2])
	case matchesPrefix(args, "package", "list"):
		return packageList(runtime)
	case matchesPrefix(args, "package", "inspect"):
		if len(args) != 3 {
			return errors.New("usage: package inspect <package-id>")
		}
		return packageInspect(runtime, args[2])
	case matchesPrefix(args, "package", "install"):
		if len(args) != 3 {
			return errors.New("usage: package install <dir>")
		}
		manifest, err := runtime.InstallPackageDir(ctx, args[2])
		if err != nil {
			return err
		}
		fmt.Printf("installed %s\n", manifest.ID)
		return nil
	case matchesPrefix(args, "workflow", "capture"):
		if len(args) != 4 {
			return errors.New("usage: workflow capture <directory> <alias>")
		}
		workflow, err := runtime.CaptureWorkflowDir(args[2], args[3])
		if err != nil {
			return err
		}
		return workflowPrint(workflow)
	case matchesPrefix(args, "workflow", "import"):
		if len(args) != 4 {
			return errors.New("usage: workflow import <alias> <artifact-cid>")
		}
		if err := runtime.ImportWorkflow(kernel.Workflow{ID: args[2], ArtifactCID: args[3]}); err != nil {
			return err
		}
		return workflowPrint(runtime.Workflows()[len(runtime.Workflows())-1])
	case matchesPrefix(args, "workflow", "list"):
		if len(args) != 2 {
			return errors.New("usage: workflow list")
		}
		return workflowsPrint(runtime.Workflows())
	case matchesPrefix(args, "workflow", "inspect"):
		if len(args) != 3 {
			return errors.New("usage: workflow inspect <alias-or-cid>")
		}
		for _, workflow := range runtime.Workflows() {
			if workflow.ID == args[2] || workflow.ArtifactCID == args[2] {
				return workflowPrint(workflow)
			}
		}
		return errors.New("workflow is not imported")
	case matchesPrefix(args, "workflow", "activate"):
		if len(args) != 3 {
			return errors.New("usage: workflow activate <alias>")
		}
		if err := runtime.ActivateWorkflow(args[2]); err != nil {
			return err
		}
		return workflowsPrint(runtime.Workflows())
	case matchesPrefix(args, "workflow", "deactivate"):
		if len(args) != 3 {
			return errors.New("usage: workflow deactivate <alias>")
		}
		if err := runtime.DeactivateWorkflow(args[2]); err != nil {
			return err
		}
		return workflowsPrint(runtime.Workflows())
	case matchesPrefix(args, "workflow", "revoke"):
		if len(args) != 3 {
			return errors.New("usage: workflow revoke <alias>")
		}
		if err := runtime.RevokeWorkflow(args[2]); err != nil {
			return err
		}
		return workflowsPrint(runtime.Workflows())
	case matchesPrefix(args, "relay", "export"):
		if len(args) != 3 {
			return errors.New("usage: relay export <path>")
		}
		return relayExport(runtime, args[2])
	case matchesPrefix(args, "relay", "import"):
		if len(args) != 3 {
			return errors.New("usage: relay import <path>")
		}
		return relayImport(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "serve"):
		if len(args) != 3 {
			return errors.New("usage: relay serve <addr>")
		}
		return relayServe(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "pull"):
		if len(args) != 3 {
			return errors.New("usage: relay pull <peer-id>")
		}
		return relayPull(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "push"):
		if len(args) != 3 {
			return errors.New("usage: relay push <peer-id>")
		}
		return relayPush(ctx, runtime, args[2])
	case matchesPrefix(args, "relay", "peer", "local"):
		if len(args) != 4 {
			return errors.New("usage: relay peer local show")
		}
		if args[3] != "show" {
			return errors.New("usage: relay peer local show")
		}
		fmt.Printf("%s\t%s\n", runtime.LocalPeerID(), runtime.LocalPeerPublicKey())
		return nil
	case matchesPrefix(args, "relay", "policy", "claim", "list"):
		if len(args) != 4 {
			return errors.New("usage: relay policy claim list")
		}
		return relayPolicyClaimList(runtime)
	case matchesPrefix(args, "relay", "policy", "claim", "set"):
		if len(args) != 8 {
			return errors.New("usage: relay policy claim set <protocol-pcid> <role|*> <min-attesters> <any|peer-id,peer-id>")
		}
		return relayPolicyClaimSet(runtime, args[4:])
	case matchesPrefix(args, "relay", "policy", "claim", "set-weighted"):
		if len(args) != 10 {
			return errors.New("usage: relay policy claim set-weighted <protocol-pcid> <role|*> <min-attesters> <min-weight> <any|peer-id,peer-id> <any|class,class>")
		}
		return relayPolicyClaimSetWeighted(runtime, args[4:])
	case matchesPrefix(args, "relay", "policy", "claim", "set-federated"):
		if len(args) != 12 {
			return errors.New("usage: relay policy claim set-federated <protocol-pcid> <role|*> <min-attesters> <min-weight> <min-federations> <any|peer-id,peer-id> <any|class,class> <any|federation,federation>")
		}
		return relayPolicyClaimSetFederated(runtime, args[4:])
	case matchesPrefix(args, "relay", "policy", "claim", "remove"):
		if len(args) != 6 {
			return errors.New("usage: relay policy claim remove <protocol-pcid> <role|*>")
		}
		return runtime.RemoveClaimPolicy(args[4], args[5])
	case matchesPrefix(args, "relay", "peer", "list"):
		return relayPeerList(runtime)
	case matchesPrefix(args, "relay", "peer", "discover"):
		if len(args) != 4 && len(args) != 5 {
			return errors.New("usage: relay peer discover <peer-card-url> [seed]")
		}
		seed := false
		if len(args) == 5 {
			if args[4] != "seed" {
				return errors.New("usage: relay peer discover <peer-card-url> [seed]")
			}
			seed = true
		}
		return relayPeerDiscover(ctx, runtime, args[3], seed)
	case matchesPrefix(args, "relay", "peer", "allow"):
		if len(args) != 8 {
			return errors.New("usage: relay peer allow <peer-id> <batch-url> <import-url> <public-key> <pull|no-pull> <push|no-push>")
		}
		return relayPeerAllow(runtime, args[3:])
	case matchesPrefix(args, "relay", "peer", "promote"):
		if len(args) != 5 {
			return errors.New("usage: relay peer promote <peer-id> <pull|push|both>")
		}
		return relayPeerPromote(runtime, args[3], args[4])
	case matchesPrefix(args, "relay", "peer", "classify"):
		if len(args) != 6 {
			return errors.New("usage: relay peer classify <peer-id> <class> <weight>")
		}
		return relayPeerClassify(runtime, args[3], args[4], args[5])
	case matchesPrefix(args, "relay", "peer", "federate"):
		if len(args) != 5 {
			return errors.New("usage: relay peer federate <peer-id> <federation>")
		}
		return relayPeerFederate(runtime, args[3], args[4])
	case matchesPrefix(args, "relay", "peer", "revoke"):
		if len(args) != 5 {
			return errors.New("usage: relay peer revoke <peer-id>")
		}
		return runtime.RevokePeer(args[4])
	default:
		output, err := runtime.RunCommand(ctx, args)
		if err != nil {
			return err
		}
		if strings.TrimSpace(output) != "" {
			fmt.Println(output)
		}
		return nil
	}
}

func defaultRuntimeRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".moks"), nil
}

func matchesPrefix(args []string, prefix ...string) bool {
	if len(args) < len(prefix) {
		return false
	}
	for index := range prefix {
		if args[index] != prefix[index] {
			return false
		}
	}
	return true
}

func packageList(runtime *kernel.Runtime) error {
	for _, manifest := range runtime.PackageManifests() {
		fmt.Printf("%s\t%s\n", manifest.ID, manifest.Version)
	}
	return nil
}

func routeList(runtime *kernel.Runtime) error {
	// Intent: Let operators inspect the explicit route table that the kernel has
	// derived from package claims so the current routing role is visible and
	// debuggable from the CLI. Source: DI-rutom
	for _, route := range runtime.ProtocolRoutes() {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			route.ProtocolPCID,
			route.Role,
			route.RouteType,
			route.PackageID,
			route.PackageVersion,
			strings.Join(route.Families, ","),
			strings.Join(route.EmitsProtocols, ","),
		)
	}
	return nil
}

func routeScopeList(runtime *kernel.Runtime) error {
	body, err := json.MarshalIndent(struct {
		Builtin []string               `json:"builtin"`
		Local   []grid.TraceScopeAlias `json:"local,omitempty"`
	}{
		Builtin: []string{"direct-hops", "deep-hops", "one-branch:<candidate-id>"},
		Local:   runtime.TraceScopeAliases(),
	}, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func routeScopeInspect(runtime *kernel.Runtime, name string, query kernel.RouteScopeGroupQuery) error {
	inspection, ok := runtime.InspectTraceScopeWithQuery(name, query)
	if !ok {
		return fmt.Errorf("unknown route scope: %s", name)
	}
	body, err := json.MarshalIndent(inspection, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

// Intent: Keep reusable local trace views under an explicit CLI family so
// operators can define routing-inspection aliases without modifying code.
// Source: DI-bemok
func routeScopeSet(runtime *kernel.Runtime, name string, args []string) error {
	clauses := make([]grid.TraceScopeClause, 0, len(args)/2)
	for index := 0; index < len(args); index += 2 {
		clauses = append(clauses, grid.TraceScopeClause{
			Kind:   args[index],
			Target: args[index+1],
		})
	}
	if err := runtime.SetTraceScopeAlias(grid.TraceScopeAlias{
		Name:    name,
		Clauses: clauses,
	}); err != nil {
		return err
	}
	fmt.Printf("route scope set %s\n", name)
	return nil
}

func routeScopeRemove(runtime *kernel.Runtime, name string) error {
	if err := runtime.RemoveTraceScopeAlias(name); err != nil {
		return err
	}
	fmt.Printf("route scope removed %s\n", name)
	return nil
}

func routeInspect(runtime *kernel.Runtime, protocolPCID string) error {
	// Intent: Expose a machine-readable per-protocol route query so route
	// consumers can ask what direct handlers or hops exist for one input pCID.
	// Source: DI-fotav
	body, err := json.MarshalIndent(runtime.ProtocolRoutesForProtocol(protocolPCID), "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func routePlan(runtime *kernel.Runtime, protocolPCID string, trace bool, filter kernel.RoutePlanTraceFilter) error {
	// Intent: Expose the kernel's preferred executable route choice for one
	// input pCID so route consumers can ask what the kernel would actually use.
	// Source: DI-pabut
	plan := runtime.ProtocolRoutePlan(protocolPCID)
	if trace {
		plan = runtime.ProtocolRoutePlanTraceFocused(protocolPCID, filter)
	}
	body, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

// Intent: Make planner policy inspectable both globally and for one input
// protocol so operators can see the inherited effective route preferences that
// the kernel will actually use. Source: DI-posek
func routePolicyShow(runtime *kernel.Runtime, protocolPCID string, role string) error {
	if strings.TrimSpace(protocolPCID) != "" && strings.TrimSpace(role) != "" {
		body, err := json.MarshalIndent(struct {
			ProtocolPCID string               `json:"protocol_pcid"`
			Role         string               `json:"role"`
			Global       grid.RoutePlanPolicy `json:"global"`
			Protocol     grid.RoutePlanPolicy `json:"protocol"`
			Effective    grid.RoutePlanPolicy `json:"effective"`
		}{
			ProtocolPCID: protocolPCID,
			Role:         role,
			Global:       runtime.RoutePlanPolicy(),
			Protocol:     runtime.EffectiveRoutePlanPolicy(protocolPCID),
			Effective:    runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, role),
		}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	}
	if strings.TrimSpace(protocolPCID) != "" {
		body, err := json.MarshalIndent(struct {
			ProtocolPCID string               `json:"protocol_pcid"`
			Global       grid.RoutePlanPolicy `json:"global"`
			Effective    grid.RoutePlanPolicy `json:"effective"`
		}{
			ProtocolPCID: protocolPCID,
			Global:       runtime.RoutePlanPolicy(),
			Effective:    runtime.EffectiveRoutePlanPolicy(protocolPCID),
		}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	}
	body, err := json.MarshalIndent(struct {
		Global        grid.RoutePlanPolicy               `json:"global"`
		Protocols     []grid.ProtocolRoutePlanPolicy     `json:"protocols,omitempty"`
		ProtocolRoles []grid.ProtocolRoleRoutePlanPolicy `json:"protocol_roles,omitempty"`
	}{
		Global:        runtime.RoutePlanPolicy(),
		Protocols:     runtime.ProtocolRoutePlanPolicies(),
		ProtocolRoles: runtime.ProtocolRoleRoutePlanPolicies(),
	}, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func routePolicySet(runtime *kernel.Runtime, args []string) error {
	policy := parseRoutePlanPolicyArgs(args)
	if err := runtime.SetRoutePlanPolicy(policy); err != nil {
		return err
	}
	fmt.Println("route policy set")
	return nil
}

func routePolicySetForProtocol(runtime *kernel.Runtime, protocolPCID string, args []string) error {
	policy := parseRoutePlanPolicyArgs(args)
	if err := runtime.SetProtocolRoutePlanPolicy(protocolPCID, policy); err != nil {
		return err
	}
	fmt.Printf("route policy set for %s\n", protocolPCID)
	return nil
}

func routePolicyRemove(runtime *kernel.Runtime, protocolPCID string) error {
	if err := runtime.RemoveProtocolRoutePlanPolicy(protocolPCID); err != nil {
		return err
	}
	fmt.Printf("route policy removed for %s\n", protocolPCID)
	return nil
}

func routePolicySetForRole(runtime *kernel.Runtime, protocolPCID string, role string, args []string) error {
	policy := parseRoutePlanPolicyArgs(args)
	if err := runtime.SetProtocolRoleRoutePlanPolicy(protocolPCID, role, policy); err != nil {
		return err
	}
	fmt.Printf("route policy set for %s role %s\n", protocolPCID, role)
	return nil
}

func routePolicyRemoveRole(runtime *kernel.Runtime, protocolPCID string, role string) error {
	if err := runtime.RemoveProtocolRoleRoutePlanPolicy(protocolPCID, role); err != nil {
		return err
	}
	fmt.Printf("route policy removed for %s role %s\n", protocolPCID, role)
	return nil
}

func parseRoutePlanPolicyArgs(args []string) grid.RoutePlanPolicy {
	return grid.RoutePlanPolicy{
		PreferRouteTypes: parsePolicyListArg(args[0]),
		AvoidRouteTypes:  parsePolicyListArg(args[1]),
		PreferRoles:      parsePolicyListArg(args[2]),
		AvoidRoles:       parsePolicyListArg(args[3]),
	}
}

func parsePolicyListArg(raw string) []string {
	if raw == "-" {
		return nil
	}
	return strings.Split(raw, ",")
}

func packageInspect(runtime *kernel.Runtime, id string) error {
	manifest, ok := runtime.PackageManifest(id)
	if !ok {
		return fmt.Errorf("unknown package: %s", id)
	}
	body, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func relayExport(runtime *kernel.Runtime, path string) error {
	batch, err := runtime.SignedExportBatch()
	if err != nil {
		return err
	}
	body, err := json.MarshalIndent(batch, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}

func relayImport(ctx context.Context, runtime *kernel.Runtime, path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var batch grid.Batch
	if err := json.Unmarshal(body, &batch); err != nil {
		return err
	}
	return runtime.ImportBatch(ctx, batch)
}

func relayPolicyClaimList(runtime *kernel.Runtime) error {
	for _, policy := range runtime.ClaimPolicies() {
		attesters := "any-known-peer"
		if len(policy.AllowedAttesters) > 0 {
			attesters = strings.Join(policy.AllowedAttesters, ",")
		}
		classes := "any-class"
		if len(policy.AllowedClasses) > 0 {
			classes = strings.Join(policy.AllowedClasses, ",")
		}
		federations := "any-federation"
		if len(policy.AllowedFederations) > 0 {
			federations = strings.Join(policy.AllowedFederations, ",")
		}
		fmt.Printf("%s\t%s\t%d\t%d\t%d\t%s\t%s\t%s\n", policy.ProtocolPCID, policy.Role, policy.MinAttesters, policy.MinTrustWeight, policy.MinFederations, attesters, classes, federations)
	}
	return nil
}

func relayPolicyClaimSet(runtime *kernel.Runtime, args []string) error {
	minAttesters, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	policy := grid.ClaimTrustPolicy{
		ProtocolPCID: args[0],
		Role:         args[1],
		MinAttesters: minAttesters,
	}
	if args[3] != "any" {
		policy.AllowedAttesters = strings.Split(args[3], ",")
	}
	if err := runtime.SetClaimPolicy(policy); err != nil {
		return err
	}
	fmt.Printf("policy set %s %s quorum=%d\n", policy.ProtocolPCID, policy.Role, policy.MinAttesters)
	return nil
}

func relayPolicyClaimSetWeighted(runtime *kernel.Runtime, args []string) error {
	minAttesters, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	minWeight, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	policy := grid.ClaimTrustPolicy{
		ProtocolPCID:   args[0],
		Role:           args[1],
		MinAttesters:   minAttesters,
		MinTrustWeight: minWeight,
	}
	if args[4] != "any" {
		policy.AllowedAttesters = strings.Split(args[4], ",")
	}
	if args[5] != "any" {
		policy.AllowedClasses = strings.Split(args[5], ",")
	}
	if err := runtime.SetClaimPolicy(policy); err != nil {
		return err
	}
	fmt.Printf("policy set %s %s quorum=%d weight=%d\n", policy.ProtocolPCID, policy.Role, policy.MinAttesters, policy.MinTrustWeight)
	return nil
}

func relayPolicyClaimSetFederated(runtime *kernel.Runtime, args []string) error {
	minAttesters, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	minWeight, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	minFederations, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}
	policy := grid.ClaimTrustPolicy{
		ProtocolPCID:   args[0],
		Role:           args[1],
		MinAttesters:   minAttesters,
		MinTrustWeight: minWeight,
		MinFederations: minFederations,
	}
	if args[5] != "any" {
		policy.AllowedAttesters = strings.Split(args[5], ",")
	}
	if args[6] != "any" {
		policy.AllowedClasses = strings.Split(args[6], ",")
	}
	if args[7] != "any" {
		policy.AllowedFederations = strings.Split(args[7], ",")
	}
	if err := runtime.SetClaimPolicy(policy); err != nil {
		return err
	}
	fmt.Printf("policy set %s %s quorum=%d weight=%d federations=%d\n", policy.ProtocolPCID, policy.Role, policy.MinAttesters, policy.MinTrustWeight, policy.MinFederations)
	return nil
}

func relayServe(ctx context.Context, runtime *kernel.Runtime, addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: relayHandler(ctx, runtime),
	}
	fmt.Printf("relay serving on %s\n", addr)
	return server.ListenAndServe()
}

func relayHandler(ctx context.Context, runtime *kernel.Runtime) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /relay/peer-card", func(writer http.ResponseWriter, request *http.Request) {
		// Intent: Let peers discover runtime-owned relay identity and endpoints
		// without silently granting any exchange permissions.
		// Source: DI-vemut
		card := grid.PeerCard{
			PeerID:      runtime.LocalPeerID(),
			PublicKey:   runtime.LocalPeerPublicKey(),
			BatchURL:    absoluteRelayURL(request, "/relay/batch"),
			ImportURL:   absoluteRelayURL(request, "/relay/import"),
			DiscoverURL: absoluteRelayURL(request, "/relay/peer-card"),
		}
		body, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(body)
	})
	mux.HandleFunc("GET /relay/batch", func(writer http.ResponseWriter, _ *http.Request) {
		batch, err := runtime.SignedExportBatch()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := json.MarshalIndent(batch, "", "  ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("X-Moks-Peer-ID", runtime.LocalPeerID())
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(body)
	})
	mux.HandleFunc("POST /relay/import", func(writer http.ResponseWriter, request *http.Request) {
		peerID := strings.TrimSpace(request.Header.Get("X-Moks-Peer-ID"))
		if peerID == "" {
			http.Error(writer, "missing X-Moks-Peer-ID header", http.StatusForbidden)
			return
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var batch grid.Batch
		if err := json.Unmarshal(body, &batch); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := runtime.ImportBatchFromPeer(ctx, peerID, batch, "push"); err != nil {
			http.Error(writer, err.Error(), http.StatusForbidden)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}

func absoluteRelayURL(request *http.Request, path string) string {
	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + request.Host + path
}

func relayPull(ctx context.Context, runtime *kernel.Runtime, peerID string) error {
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, peer.BatchURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Moks-Peer-ID", runtime.LocalPeerID())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("relay pull failed: %s: %s", response.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var batch grid.Batch
	if err := json.Unmarshal(body, &batch); err != nil {
		return err
	}
	return runtime.ImportBatchFromPeer(ctx, peerID, batch, "pull")
}

func relayPush(ctx context.Context, runtime *kernel.Runtime, peerID string) error {
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("peer not allowed: %s", peerID)
	}
	batch, err := runtime.SignedExportBatch()
	if err != nil {
		return err
	}
	body, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, peer.ImportURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Moks-Peer-ID", runtime.LocalPeerID())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return fmt.Errorf("relay push failed: %s: %s", response.Status, strings.TrimSpace(string(responseBody)))
	}
	return nil
}

func relayPeerDiscover(ctx context.Context, runtime *kernel.Runtime, cardURL string, seed bool) error {
	// Intent: Keep discovery separate from trust by fetching peer metadata and
	// only seeding a no-pull/no-push local entry when the operator explicitly
	// asks for it.
	// Source: DI-kasud
	card, err := fetchPeerCard(ctx, cardURL)
	if err != nil {
		return err
	}
	seeded := false
	if seed {
		if err := runtime.AllowPeer(grid.AllowedPeer{
			PeerID:            card.PeerID,
			BatchURL:          card.BatchURL,
			ImportURL:         card.ImportURL,
			PublicKey:         card.PublicKey,
			AllowPull:         false,
			AllowPush:         false,
			AttesterClass:     "peer",
			AttestationWeight: 1,
			Federation:        "independent",
		}); err != nil {
			return err
		}
		seeded = true
	}
	fmt.Printf(
		"peer_id: %s\npublic_key: %s\nbatch_url: %s\nimport_url: %s\ndiscover_url: %s\nseeded_untrusted: %t\nallow_command: moks relay peer allow %s %s %s %s no-pull no-push\nenable_pull_command: moks relay peer allow %s %s %s %s pull no-push\nenable_push_command: moks relay peer allow %s %s %s %s no-pull push\nenable_both_command: moks relay peer allow %s %s %s %s pull push\n",
		card.PeerID,
		card.PublicKey,
		card.BatchURL,
		card.ImportURL,
		card.DiscoverURL,
		seeded,
		card.PeerID,
		card.BatchURL,
		card.ImportURL,
		card.PublicKey,
		card.PeerID,
		card.BatchURL,
		card.ImportURL,
		card.PublicKey,
		card.PeerID,
		card.BatchURL,
		card.ImportURL,
		card.PublicKey,
		card.PeerID,
		card.BatchURL,
		card.ImportURL,
		card.PublicKey,
	)
	return nil
}

func fetchPeerCard(ctx context.Context, cardURL string) (grid.PeerCard, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cardURL, nil)
	if err != nil {
		return grid.PeerCard{}, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return grid.PeerCard{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return grid.PeerCard{}, fmt.Errorf("peer discovery failed: %s: %s", response.Status, strings.TrimSpace(string(body)))
	}
	var card grid.PeerCard
	if err := json.NewDecoder(response.Body).Decode(&card); err != nil {
		return grid.PeerCard{}, err
	}
	if err := card.Validate(); err != nil {
		return grid.PeerCard{}, err
	}
	return card, nil
}

func registerBuiltins(runtime *kernel.Runtime) error {
	if err := runtime.RegisterBuiltin(contextpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(knowledgepkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(inventorypkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(runspkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(linkspkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(maintenancepkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(receivingpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(procedurespkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(trainingpkg.Package()); err != nil {
		return err
	}
	if err := runtime.RegisterBuiltin(builtin.OpsPackage()); err != nil {
		return err
	}
	return nil
}

func relayPeerList(runtime *kernel.Runtime) error {
	for _, peer := range runtime.AllowedPeers() {
		fmt.Printf("%s\tpull=%t\tpush=%t\tclass=%s\tweight=%d\tfederation=%s\tbatch=%s\timport=%s\tpub=%s\n", peer.PeerID, peer.AllowPull, peer.AllowPush, peer.AttesterClass, peer.AttestationWeight, peer.Federation, peer.BatchURL, peer.ImportURL, peer.PublicKey)
	}
	return nil
}

func relayPeerAllow(runtime *kernel.Runtime, args []string) error {
	allowPull := args[4] == "pull"
	allowPush := args[5] == "push"
	if args[4] != "pull" && args[4] != "no-pull" {
		return errors.New("usage: relay peer allow <peer-id> <batch-url> <import-url> <public-key> <pull|no-pull> <push|no-push>")
	}
	if args[5] != "push" && args[5] != "no-push" {
		return errors.New("usage: relay peer allow <peer-id> <batch-url> <import-url> <public-key> <pull|no-pull> <push|no-push>")
	}
	return runtime.AllowPeer(grid.AllowedPeer{
		PeerID:            args[0],
		BatchURL:          args[1],
		ImportURL:         args[2],
		PublicKey:         args[3],
		AllowPull:         allowPull,
		AllowPush:         allowPush,
		AttesterClass:     "peer",
		AttestationWeight: 1,
		Federation:        "independent",
	})
}

func relayPeerPromote(runtime *kernel.Runtime, peerID string, mode string) error {
	// Intent: Promote a discovered peer's exchange policy without forcing the
	// operator to retype stored metadata, while keeping the trust step explicit.
	// Source: DI-lutep
	peer, ok := runtime.LookupPeer(peerID)
	if !ok {
		return fmt.Errorf("unknown peer: %s", peerID)
	}
	switch mode {
	case "pull":
		peer.AllowPull = true
	case "push":
		peer.AllowPush = true
	case "both":
		peer.AllowPull = true
		peer.AllowPush = true
	default:
		return errors.New("usage: relay peer promote <peer-id> <pull|push|both>")
	}
	if err := runtime.AllowPeer(peer); err != nil {
		return err
	}
	fmt.Printf("promoted %s pull=%t push=%t\n", peer.PeerID, peer.AllowPull, peer.AllowPush)
	return nil
}

func relayPeerClassify(runtime *kernel.Runtime, peerID string, attesterClass string, weightArg string) error {
	weight, err := strconv.Atoi(weightArg)
	if err != nil {
		return err
	}
	if err := runtime.SetPeerTrust(peerID, attesterClass, weight); err != nil {
		return err
	}
	fmt.Printf("classified %s class=%s weight=%d\n", peerID, attesterClass, weight)
	return nil
}

func relayPeerFederate(runtime *kernel.Runtime, peerID string, federation string) error {
	if err := runtime.SetPeerFederation(peerID, federation); err != nil {
		return err
	}
	fmt.Printf("federated %s federation=%s\n", peerID, federation)
	return nil
}
