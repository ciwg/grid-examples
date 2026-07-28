package kernel

import (
	"slices"
	"strings"
)

type ProtocolRoute struct {
	PackageID      string
	PackageVersion string
	ProtocolPCID   string
	Role           string
	Summary        string
	Families       []string
}

type registeredRoute struct {
	owner        *activePackage
	protocolPCID string
	role         string
	summary      string
	families     []string
}

// Intent: Expose the kernel's current claim-derived route table so routing is
// visible as a first-class service role instead of staying implicit in package
// activation alone. Source: DI-rutom
func (runtime *Runtime) ProtocolRoutes() []ProtocolRoute {
	routes := make([]ProtocolRoute, 0, len(runtime.routes))
	for _, route := range runtime.routes {
		routes = append(routes, ProtocolRoute{
			PackageID:      route.owner.manifest.ID,
			PackageVersion: route.owner.manifest.Version,
			ProtocolPCID:   route.protocolPCID,
			Role:           route.role,
			Summary:        route.summary,
			Families:       append([]string{}, route.families...),
		})
	}
	slices.SortFunc(routes, func(left, right ProtocolRoute) int {
		if diff := strings.Compare(left.ProtocolPCID, right.ProtocolPCID); diff != 0 {
			return diff
		}
		if diff := strings.Compare(left.Role, right.Role); diff != 0 {
			return diff
		}
		return strings.Compare(left.PackageID, right.PackageID)
	})
	return routes
}
