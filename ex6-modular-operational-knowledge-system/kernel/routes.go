package kernel

import (
	"slices"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
)

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
