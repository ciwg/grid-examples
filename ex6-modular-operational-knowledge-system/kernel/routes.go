package kernel

import (
	"slices"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
)

type RoutePlan struct {
	ProtocolPCID string               `json:"protocol_pcid"`
	Preferred    *RoutePlanCandidate  `json:"preferred,omitempty"`
	Candidates   []RoutePlanCandidate `json:"candidates,omitempty"`
}

type RoutePlanCandidate struct {
	Route      grid.RouteRegistration `json:"route"`
	Executable bool                   `json:"executable"`
	Next       []RoutePlan            `json:"next,omitempty"`
}

type registeredRoute struct {
	owner        *activePackage
	protocolPCID string
	role         string
	routeType    string
	emits        []string
	summary      string
	families     []string
}

// Intent: Expose the kernel's current claim-derived route table so routing is
// visible as a first-class service role instead of staying implicit in package
// activation alone. Source: DI-rutom
func (runtime *Runtime) ProtocolRoutes() []grid.RouteRegistration {
	routes := make([]grid.RouteRegistration, 0, len(runtime.routes))
	for _, route := range runtime.routes {
		routes = append(routes, grid.RouteRegistration{
			PackageID:      route.owner.manifest.ID,
			PackageVersion: route.owner.manifest.Version,
			ProtocolPCID:   route.protocolPCID,
			Role:           route.role,
			RouteType:      route.routeType,
			EmitsProtocols: append([]string{}, route.emits...),
			Summary:        route.summary,
			Families:       append([]string{}, route.families...),
		})
	}
	slices.SortFunc(routes, func(left, right grid.RouteRegistration) int {
		if diff := strings.Compare(left.ProtocolPCID, right.ProtocolPCID); diff != 0 {
			return diff
		}
		if diff := strings.Compare(left.Role, right.Role); diff != 0 {
			return diff
		}
		if diff := strings.Compare(left.RouteType, right.RouteType); diff != 0 {
			return diff
		}
		return strings.Compare(left.PackageID, right.PackageID)
	})
	return routes
}

// Intent: Let route consumers query the current route table by input protocol
// so they can discover direct handlers and parser/transform hops without
// scraping every registered route. Source: DI-fotav
func (runtime *Runtime) ProtocolRoutesForProtocol(protocolPCID string) []grid.RouteRegistration {
	filtered := []grid.RouteRegistration{}
	for _, route := range runtime.ProtocolRoutes() {
		if route.ProtocolPCID == protocolPCID {
			filtered = append(filtered, route)
		}
	}
	return filtered
}

// Intent: Choose a deterministic preferred route plan for one input protocol
// so route consumers can ask what the kernel would actually try, not just what
// routes happen to exist. Source: DI-pabut
func (runtime *Runtime) ProtocolRoutePlan(protocolPCID string) RoutePlan {
	return runtime.protocolRoutePlan(protocolPCID, map[string]struct{}{})
}

func (runtime *Runtime) protocolRoutePlan(protocolPCID string, seen map[string]struct{}) RoutePlan {
	plan := RoutePlan{
		ProtocolPCID: protocolPCID,
	}
	if strings.TrimSpace(protocolPCID) == "" {
		return plan
	}
	if _, exists := seen[protocolPCID]; exists {
		return plan
	}
	nextSeen := cloneSeen(seen)
	nextSeen[protocolPCID] = struct{}{}
	for _, route := range runtime.ProtocolRoutesForProtocol(protocolPCID) {
		candidate := RoutePlanCandidate{
			Route: route,
		}
		switch route.RouteType {
		case "", "direct":
			candidate.Executable = true
		case "parser", "transform":
			candidate.Executable = len(route.EmitsProtocols) > 0
			for _, emitted := range route.EmitsProtocols {
				subplan := runtime.protocolRoutePlan(emitted, nextSeen)
				candidate.Next = append(candidate.Next, subplan)
				if subplan.Preferred == nil {
					candidate.Executable = false
				}
			}
		default:
			candidate.Executable = false
		}
		plan.Candidates = append(plan.Candidates, candidate)
	}
	slices.SortFunc(plan.Candidates, func(left, right RoutePlanCandidate) int {
		return runtime.compareRoutePlanCandidates(protocolPCID, left, right)
	})
	for index := range plan.Candidates {
		if plan.Candidates[index].Executable {
			plan.Preferred = &plan.Candidates[index]
			break
		}
	}
	return plan
}

// Intent: Keep route planning protocol-aware and role-aware by evaluating each
// candidate with the effective global, per-protocol, and per-protocol-role
// policy that applies to that route. Source: DI-rivuk
func (runtime *Runtime) compareRoutePlanCandidates(protocolPCID string, left RoutePlanCandidate, right RoutePlanCandidate) int {
	if left.Executable != right.Executable {
		if left.Executable {
			return -1
		}
		return 1
	}
	if diff := runtime.comparePolicyPreference(protocolPCID, left.Route, right.Route); diff != 0 {
		return diff
	}
	if diff := compareRouteRegistrations(left.Route, right.Route); diff != 0 {
		return diff
	}
	if len(left.Next) != len(right.Next) {
		if len(left.Next) < len(right.Next) {
			return -1
		}
		return 1
	}
	return 0
}

func (runtime *Runtime) comparePolicyPreference(protocolPCID string, left grid.RouteRegistration, right grid.RouteRegistration) int {
	leftPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, left.Role)
	rightPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, right.Role)
	if diff := compareAvoided(left, right, leftPolicy, rightPolicy); diff != 0 {
		return diff
	}
	if diff := comparePreferred(left, right, leftPolicy, rightPolicy); diff != 0 {
		return diff
	}
	return 0
}

func compareAvoided(left grid.RouteRegistration, right grid.RouteRegistration, leftPolicy grid.RoutePlanPolicy, rightPolicy grid.RoutePlanPolicy) int {
	leftAvoided := routeAvoided(left, leftPolicy)
	rightAvoided := routeAvoided(right, rightPolicy)
	if leftAvoided == rightAvoided {
		return 0
	}
	if leftAvoided {
		return 1
	}
	return -1
}

func comparePreferred(left grid.RouteRegistration, right grid.RouteRegistration, leftPolicy grid.RoutePlanPolicy, rightPolicy grid.RoutePlanPolicy) int {
	leftPreferred := routePreferred(left, leftPolicy)
	rightPreferred := routePreferred(right, rightPolicy)
	if leftPreferred == rightPreferred {
		return 0
	}
	if leftPreferred {
		return -1
	}
	return 1
}

func routePreferred(route grid.RouteRegistration, policy grid.RoutePlanPolicy) bool {
	return slices.Contains(policy.PreferRouteTypes, route.RouteType) || slices.Contains(policy.PreferRoles, route.Role)
}

func routeAvoided(route grid.RouteRegistration, policy grid.RoutePlanPolicy) bool {
	return slices.Contains(policy.AvoidRouteTypes, route.RouteType) || slices.Contains(policy.AvoidRoles, route.Role)
}

func compareRouteRegistrations(left grid.RouteRegistration, right grid.RouteRegistration) int {
	if diff := compareRouteType(left.RouteType, right.RouteType); diff != 0 {
		return diff
	}
	if diff := compareRouteRole(left.Role, right.Role); diff != 0 {
		return diff
	}
	if diff := strings.Compare(left.PackageID, right.PackageID); diff != 0 {
		return diff
	}
	return strings.Compare(left.PackageVersion, right.PackageVersion)
}

func compareRouteType(left string, right string) int {
	leftRank := routeTypeRank(left)
	rightRank := routeTypeRank(right)
	if leftRank != rightRank {
		if leftRank < rightRank {
			return -1
		}
		return 1
	}
	return strings.Compare(left, right)
}

func routeTypeRank(routeType string) int {
	switch routeType {
	case "", "direct":
		return 0
	case "parser":
		return 1
	case "transform":
		return 2
	default:
		return 3
	}
}

func compareRouteRole(left string, right string) int {
	leftRank := routeRoleRank(left)
	rightRank := routeRoleRank(right)
	if leftRank != rightRank {
		if leftRank < rightRank {
			return -1
		}
		return 1
	}
	return strings.Compare(left, right)
}

func routeRoleRank(role string) int {
	switch role {
	case "handler":
		return 0
	case "domain-behavior":
		return 1
	case "family-validator":
		return 2
	case "parser":
		return 3
	case "transform":
		return 4
	default:
		return 5
	}
}

func cloneSeen(source map[string]struct{}) map[string]struct{} {
	cloned := make(map[string]struct{}, len(source))
	for key := range source {
		cloned[key] = struct{}{}
	}
	return cloned
}
