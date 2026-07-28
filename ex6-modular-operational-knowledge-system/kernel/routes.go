package kernel

import (
	"slices"
	"strconv"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-modular-operational-knowledge-system/grid"
)

type RoutePlan struct {
	ProtocolPCID string                `json:"protocol_pcid"`
	Preferred    *RoutePlanCandidate   `json:"preferred,omitempty"`
	Candidates   []RoutePlanCandidate  `json:"candidates,omitempty"`
	Explanation  *RoutePlanExplanation `json:"explanation,omitempty"`
}

type RoutePlanCandidate struct {
	Route       grid.RouteRegistration        `json:"route"`
	Executable  bool                          `json:"executable"`
	Next        []RoutePlan                   `json:"next,omitempty"`
	Explanation RoutePlanCandidateExplanation `json:"explanation,omitempty"`
}

type RoutePlanExplanation struct {
	Order                    []string                `json:"order,omitempty"`
	Winner                   []string                `json:"winner,omitempty"`
	Comparisons              []RoutePlanComparison   `json:"comparisons,omitempty"`
	TraceSummary             *RoutePlanTraceSummary  `json:"trace_summary,omitempty"`
	DownstreamTraceSummaries []RoutePlanTraceSummary `json:"downstream_trace_summaries,omitempty"`
	Trace                    []RoutePlanTraceStep    `json:"trace,omitempty"`
}

type RoutePlanComparison struct {
	Left         RoutePlanComparisonSide `json:"left"`
	Right        RoutePlanComparisonSide `json:"right"`
	Winner       string                  `json:"winner"`
	DecisionPath []string                `json:"decision_path,omitempty"`
}

type RoutePlanComparisonSide struct {
	PackageID      string `json:"package_id"`
	PackageVersion string `json:"package_version"`
	ProtocolPCID   string `json:"protocol_pcid"`
	Role           string `json:"role"`
	RouteType      string `json:"route_type"`
}

type RoutePlanTraceStep struct {
	Step     int      `json:"step"`
	Protocol string   `json:"protocol_pcid"`
	Event    string   `json:"event"`
	Details  []string `json:"details,omitempty"`
}

type RoutePlanTraceFilter struct {
	Kind   string `json:"kind,omitempty"`
	Target string `json:"target,omitempty"`
}

type RoutePlanTraceSummary struct {
	ProtocolPCID string                `json:"protocol_pcid"`
	Scope        string                `json:"scope"`
	HopPath      string                `json:"hop_path,omitempty"`
	HopSummary   string                `json:"hop_summary,omitempty"`
	HopDepth     int                   `json:"hop_depth"`
	HopIndex     int                   `json:"hop_index"`
	TotalSteps   int                   `json:"total_steps"`
	ShownSteps   int                   `json:"shown_steps"`
	HiddenSteps  int                   `json:"hidden_steps"`
	Filter       *RoutePlanTraceFilter `json:"filter,omitempty"`
}

type RoutePlanCandidateExplanation struct {
	GlobalPolicy      grid.RoutePlanPolicy             `json:"global_policy"`
	ProtocolPolicy    *grid.RoutePlanPolicy            `json:"protocol_policy,omitempty"`
	RolePolicy        *grid.RoutePlanPolicy            `json:"role_policy,omitempty"`
	EffectivePolicy   grid.RoutePlanPolicy             `json:"effective_policy"`
	PreferredByPolicy bool                             `json:"preferred_by_policy"`
	AvoidedByPolicy   bool                             `json:"avoided_by_policy"`
	Downstream        []RoutePlanDownstreamExplanation `json:"downstream,omitempty"`
	Notes             []string                         `json:"notes,omitempty"`
}

type RoutePlanDownstreamExplanation struct {
	ProtocolPCID   string                 `json:"protocol_pcid"`
	Executable     bool                   `json:"executable"`
	PreferredRoute string                 `json:"preferred_route,omitempty"`
	TraceSummary   *RoutePlanTraceSummary `json:"trace_summary,omitempty"`
	Winner         []string               `json:"winner,omitempty"`
	Comparisons    []RoutePlanComparison  `json:"comparisons,omitempty"`
	Notes          []string               `json:"notes,omitempty"`
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
	return runtime.protocolRoutePlan(protocolPCID, map[string]struct{}{}, nil)
}

func (runtime *Runtime) ProtocolRoutePlanTrace(protocolPCID string) RoutePlan {
	trace := &routePlanTraceRecorder{}
	plan := runtime.protocolRoutePlan(protocolPCID, map[string]struct{}{}, trace)
	if plan.Explanation == nil {
		plan.Explanation = &RoutePlanExplanation{}
	}
	plan.Explanation.Trace = trace.steps
	plan.Explanation.TraceSummary = traceSummary(protocolPCID, "root", "root", "root", 0, 0, trace.steps, trace.steps, RoutePlanTraceFilter{})
	plan.Explanation.DownstreamTraceSummaries = collectDownstreamTraceSummaries(plan)
	return plan
}

func (runtime *Runtime) ProtocolRoutePlanTraceFocused(protocolPCID string, filter RoutePlanTraceFilter) RoutePlan {
	plan := runtime.ProtocolRoutePlanTrace(protocolPCID)
	if plan.Explanation == nil {
		return plan
	}
	fullTrace := append([]RoutePlanTraceStep{}, plan.Explanation.Trace...)
	plan.Explanation.Trace = filterRoutePlanTrace(plan.Explanation.Trace, filter)
	plan.Explanation.TraceSummary = traceSummary(protocolPCID, "root", "root", "root", 0, 0, fullTrace, plan.Explanation.Trace, filter)
	plan.Explanation.DownstreamTraceSummaries = filterDownstreamTraceSummaries(plan.Explanation.DownstreamTraceSummaries, filter)
	return plan
}

func (runtime *Runtime) protocolRoutePlan(protocolPCID string, seen map[string]struct{}, trace *routePlanTraceRecorder) RoutePlan {
	plan := RoutePlan{
		ProtocolPCID: protocolPCID,
	}
	if strings.TrimSpace(protocolPCID) == "" {
		return plan
	}
	if _, exists := seen[protocolPCID]; exists {
		if trace != nil {
			trace.record(protocolPCID, "cycle-skip", "planner skipped recursive revisit for an already-seen protocol")
		}
		return plan
	}
	if trace != nil {
		trace.record(protocolPCID, "plan-start", "planner started building a route plan for this protocol")
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
				subplan := runtime.protocolRoutePlan(emitted, nextSeen, trace)
				candidate.Next = append(candidate.Next, subplan)
				if subplan.Preferred == nil {
					candidate.Executable = false
				}
			}
		default:
			candidate.Executable = false
		}
		plan.Candidates = append(plan.Candidates, candidate)
		if trace != nil {
			trace.record(
				protocolPCID,
				"candidate",
				comparisonSideID(candidate.Route),
				"executable="+boolString(candidate.Executable),
			)
		}
	}
	slices.SortFunc(plan.Candidates, func(left, right RoutePlanCandidate) int {
		return runtime.compareRoutePlanCandidates(protocolPCID, left, right, trace)
	})
	for index := range plan.Candidates {
		plan.Candidates[index].Explanation = runtime.explainRoutePlanCandidate(protocolPCID, plan.Candidates[index])
	}
	for index := range plan.Candidates {
		if plan.Candidates[index].Executable {
			plan.Preferred = &plan.Candidates[index]
			if trace != nil {
				trace.record(protocolPCID, "preferred", comparisonSideID(plan.Candidates[index].Route), "planner selected this executable route as preferred")
			}
			break
		}
	}
	if plan.Preferred == nil && trace != nil {
		trace.record(protocolPCID, "no-preferred", "no executable route was available for this protocol")
	}
	plan.Explanation = runtime.explainRoutePlan(protocolPCID, plan.Candidates, plan.Preferred)
	if trace != nil && plan.Explanation != nil {
		// Intent: Give each recursive protocol hop its own trace summary so
		// downstream explanations can say which protocol scope they describe
		// without flattening everything into the root summary. Source: DI-zafek
		scopedTrace := filterRoutePlanTrace(trace.steps, RoutePlanTraceFilter{
			Kind:   "downstream",
			Target: protocolPCID,
		})
		plan.Explanation.TraceSummary = traceSummary(protocolPCID, "root", "root", "root", 0, 0, scopedTrace, scopedTrace, RoutePlanTraceFilter{})
	}
	return plan
}

// Intent: Keep route planning protocol-aware and role-aware by evaluating each
// candidate with the effective global, per-protocol, and per-protocol-role
// policy that applies to that route. Source: DI-rivuk
func (runtime *Runtime) compareRoutePlanCandidates(protocolPCID string, left RoutePlanCandidate, right RoutePlanCandidate, trace *routePlanTraceRecorder) int {
	if trace != nil {
		trace.record(protocolPCID, "compare", comparisonSideID(left.Route), comparisonSideID(right.Route))
	}
	if left.Executable != right.Executable {
		if left.Executable {
			if trace != nil {
				trace.record(protocolPCID, "compare-result", comparisonSideID(left.Route)+" outranked "+comparisonSideID(right.Route), "executable routes rank before non-executable routes")
			}
			return -1
		}
		if trace != nil {
			trace.record(protocolPCID, "compare-result", comparisonSideID(right.Route)+" outranked "+comparisonSideID(left.Route), "executable routes rank before non-executable routes")
		}
		return 1
	}
	if diff := runtime.comparePolicyPreference(protocolPCID, left.Route, right.Route, trace); diff != 0 {
		return diff
	}
	if diff := compareRouteRegistrations(left.Route, right.Route); diff != 0 {
		if trace != nil {
			if diff < 0 {
				trace.record(protocolPCID, "compare-result", comparisonSideID(left.Route)+" outranked "+comparisonSideID(right.Route), "deterministic registration ordering broke the tie")
			} else {
				trace.record(protocolPCID, "compare-result", comparisonSideID(right.Route)+" outranked "+comparisonSideID(left.Route), "deterministic registration ordering broke the tie")
			}
		}
		return diff
	}
	if len(left.Next) != len(right.Next) {
		if len(left.Next) < len(right.Next) {
			if trace != nil {
				trace.record(protocolPCID, "compare-result", comparisonSideID(left.Route)+" outranked "+comparisonSideID(right.Route), "fewer downstream hops broke the tie")
			}
			return -1
		}
		if trace != nil {
			trace.record(protocolPCID, "compare-result", comparisonSideID(right.Route)+" outranked "+comparisonSideID(left.Route), "fewer downstream hops broke the tie")
		}
		return 1
	}
	if trace != nil {
		trace.record(protocolPCID, "compare-result", comparisonSideID(left.Route)+" tied with "+comparisonSideID(right.Route), "all deterministic comparison steps were equal")
	}
	return 0
}

func (runtime *Runtime) comparePolicyPreference(protocolPCID string, left grid.RouteRegistration, right grid.RouteRegistration, trace *routePlanTraceRecorder) int {
	leftPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, left.Role)
	rightPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, right.Role)
	if diff := compareAvoided(left, right, leftPolicy, rightPolicy); diff != 0 {
		if trace != nil {
			if diff < 0 {
				trace.record(protocolPCID, "policy", comparisonSideID(left)+" outranked "+comparisonSideID(right), "avoid rules favored the left route")
			} else {
				trace.record(protocolPCID, "policy", comparisonSideID(right)+" outranked "+comparisonSideID(left), "avoid rules favored the right route")
			}
		}
		return diff
	}
	if diff := comparePreferred(left, right, leftPolicy, rightPolicy); diff != 0 {
		if trace != nil {
			if diff < 0 {
				trace.record(protocolPCID, "policy", comparisonSideID(left)+" outranked "+comparisonSideID(right), "prefer rules favored the left route")
			} else {
				trace.record(protocolPCID, "policy", comparisonSideID(right)+" outranked "+comparisonSideID(left), "prefer rules favored the right route")
			}
		}
		return diff
	}
	return 0
}

type routePlanTraceRecorder struct {
	steps []RoutePlanTraceStep
}

func (trace *routePlanTraceRecorder) record(protocolPCID string, event string, details ...string) {
	trace.steps = append(trace.steps, RoutePlanTraceStep{
		Step:     len(trace.steps) + 1,
		Protocol: protocolPCID,
		Event:    event,
		Details:  append([]string{}, details...),
	})
}

func boolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// Intent: Keep route-plan trace output readable on larger route sets by letting
// operators focus on one candidate path or one downstream protocol without
// changing the underlying planner behavior. Source: DI-dovak
func filterRoutePlanTrace(steps []RoutePlanTraceStep, filter RoutePlanTraceFilter) []RoutePlanTraceStep {
	kind := strings.TrimSpace(filter.Kind)
	target := strings.TrimSpace(filter.Target)
	if kind == "" || target == "" {
		return steps
	}
	filtered := []RoutePlanTraceStep{}
	for _, step := range steps {
		switch kind {
		case "candidate":
			if stepMatchesCandidate(step, target) {
				filtered = append(filtered, step)
			}
		case "downstream":
			if step.Protocol == target {
				filtered = append(filtered, step)
			}
		case "depth":
			if stepMatchesDepth(step, target) {
				filtered = append(filtered, step)
			}
		default:
			return steps
		}
	}
	return renumberTraceSteps(filtered)
}

func stepMatchesCandidate(step RoutePlanTraceStep, target string) bool {
	for _, detail := range step.Details {
		if strings.Contains(detail, target) {
			return true
		}
	}
	return false
}

func stepMatchesDepth(step RoutePlanTraceStep, target string) bool {
	depth, atLeast, ok := parseDepthFilter(target)
	if !ok {
		return true
	}
	stepDepth, found := traceStepDepth(step)
	if !found {
		return false
	}
	if atLeast {
		return stepDepth >= depth
	}
	return stepDepth == depth
}

func traceStepDepth(step RoutePlanTraceStep) (int, bool) {
	for _, detail := range step.Details {
		if strings.HasPrefix(detail, "depth=") {
			value := strings.TrimPrefix(detail, "depth=")
			depth, err := strconv.Atoi(value)
			if err == nil {
				return depth, true
			}
		}
	}
	return 0, false
}

func parseDepthFilter(target string) (int, bool, bool) {
	value := strings.TrimSpace(target)
	if value == "" {
		return 0, false, false
	}
	atLeast := strings.HasSuffix(value, "+")
	if atLeast {
		value = strings.TrimSuffix(value, "+")
	}
	depth, err := strconv.Atoi(value)
	if err != nil || depth < 0 {
		return 0, false, false
	}
	return depth, atLeast, true
}

func renumberTraceSteps(steps []RoutePlanTraceStep) []RoutePlanTraceStep {
	out := make([]RoutePlanTraceStep, len(steps))
	copy(out, steps)
	for index := range out {
		out[index].Step = index + 1
	}
	return out
}

// Intent: Expose nested protocol-hop summaries alongside the root trace so
// operators can see each downstream scope directly in traced route-plan output.
// Source: DI-rukav
func collectDownstreamTraceSummaries(plan RoutePlan) []RoutePlanTraceSummary {
	out := []RoutePlanTraceSummary{}
	collectDownstreamTraceSummariesFromCandidates(plan.Candidates, "root", 0, &out)
	slices.SortFunc(out, func(left, right RoutePlanTraceSummary) int {
		if left.HopDepth != right.HopDepth {
			if left.HopDepth < right.HopDepth {
				return -1
			}
			return 1
		}
		if diff := strings.Compare(left.HopPath, right.HopPath); diff != 0 {
			return diff
		}
		return strings.Compare(left.ProtocolPCID, right.ProtocolPCID)
	})
	assignHopIndices(out)
	return out
}

func collectDownstreamTraceSummariesFromCandidates(candidates []RoutePlanCandidate, pathPrefix string, depth int, out *[]RoutePlanTraceSummary) {
	for candidateIndex, candidate := range candidates {
		for nextIndex, next := range candidate.Next {
			hopPath := downstreamHopPath(pathPrefix, candidate.Route, candidateIndex, next.ProtocolPCID, nextIndex)
			hopSummary := downstreamHopSummary(candidate.Route, candidateIndex, next.ProtocolPCID, nextIndex)
			hopDepth := depth + 1
			if next.Explanation != nil && next.Explanation.TraceSummary != nil {
				summary := *next.Explanation.TraceSummary
				summary.Scope = "downstream"
				summary.HopPath = hopPath
				summary.HopSummary = hopSummary
				summary.HopDepth = hopDepth
				*out = append(*out, summary)
			}
			collectDownstreamTraceSummariesFromCandidates(next.Candidates, hopPath, hopDepth, out)
		}
	}
}

func assignHopIndices(summaries []RoutePlanTraceSummary) {
	counts := map[int]int{}
	for index := range summaries {
		depth := summaries[index].HopDepth
		counts[depth]++
		summaries[index].HopIndex = counts[depth]
	}
}

func filterDownstreamTraceSummaries(summaries []RoutePlanTraceSummary, filter RoutePlanTraceFilter) []RoutePlanTraceSummary {
	kind := strings.TrimSpace(filter.Kind)
	target := strings.TrimSpace(filter.Target)
	if kind == "" || target == "" {
		return summaries
	}
	filtered := []RoutePlanTraceSummary{}
	for _, summary := range summaries {
		switch kind {
		case "downstream":
			if summary.ProtocolPCID == target {
				filtered = append(filtered, summary)
			}
		case "depth":
			if summaryMatchesDepth(summary, target) {
				filtered = append(filtered, summary)
			}
		default:
			filtered = append(filtered, summary)
		}
	}
	return filtered
}

func summaryMatchesDepth(summary RoutePlanTraceSummary, target string) bool {
	depth, atLeast, ok := parseDepthFilter(target)
	if !ok {
		return true
	}
	if atLeast {
		return summary.HopDepth >= depth
	}
	return summary.HopDepth == depth
}

// Intent: Keep filtered route traces honest and scoped by showing which
// protocol hop the summary describes, whether it is root or downstream, and
// how many planner steps remain after filtering. Source: DI-zafek
func traceSummary(protocolPCID string, scope string, hopPath string, hopSummary string, hopDepth int, hopIndex int, full []RoutePlanTraceStep, shown []RoutePlanTraceStep, filter RoutePlanTraceFilter) *RoutePlanTraceSummary {
	summary := &RoutePlanTraceSummary{
		ProtocolPCID: protocolPCID,
		Scope:        scope,
		HopPath:      hopPath,
		HopSummary:   hopSummary,
		HopDepth:     hopDepth,
		HopIndex:     hopIndex,
		TotalSteps:   len(full),
		ShownSteps:   len(shown),
		HiddenSteps:  len(full) - len(shown),
	}
	if strings.TrimSpace(filter.Kind) != "" && strings.TrimSpace(filter.Target) != "" {
		filterCopy := filter
		summary.Filter = &filterCopy
	}
	return summary
}

// Intent: Give each downstream summary a stable path label so repeated hops to
// the same protocol remain distinguishable in traced route-plan output.
// Source: DI-vatuk
func downstreamHopPath(pathPrefix string, route grid.RouteRegistration, candidateIndex int, protocolPCID string, nextIndex int) string {
	return pathPrefix +
		" > " + comparisonSideID(route) +
		"#" + strconv.Itoa(candidateIndex+1) +
		" > " + protocolPCID +
		"#" + strconv.Itoa(nextIndex+1)
}

// Intent: Keep downstream trace summaries readable at a glance by exposing a
// short textual hop description beside the full structured hop path.
// Source: DI-lupav
func downstreamHopSummary(route grid.RouteRegistration, candidateIndex int, protocolPCID string, nextIndex int) string {
	return comparisonSideID(route) +
		" [" + strconv.Itoa(candidateIndex+1) + "] -> " +
		protocolPCID +
		" [" + strconv.Itoa(nextIndex+1) + "]"
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

// Intent: Make route-plan output explainable by exposing the policy layers and
// deterministic ranking reasons that caused one candidate to beat another.
// Source: DI-lavik
func (runtime *Runtime) explainRoutePlan(protocolPCID string, candidates []RoutePlanCandidate, preferred *RoutePlanCandidate) *RoutePlanExplanation {
	explanation := &RoutePlanExplanation{
		Order: []string{
			"executable routes rank before non-executable routes",
			"global, protocol, and exact-role planner policy layers apply prefer and avoid rules",
			"route type, route role, package identity, and downstream hop count break ties deterministically",
		},
	}
	if preferred == nil {
		explanation.Winner = []string{"no executable route was available for this protocol"}
		return explanation
	}
	if len(candidates) == 1 {
		explanation.Winner = []string{"only one candidate route exists for this protocol"}
		return explanation
	}
	for _, candidate := range candidates {
		if candidate.Route.PackageID == preferred.Route.PackageID &&
			candidate.Route.ProtocolPCID == preferred.Route.ProtocolPCID &&
			candidate.Route.Role == preferred.Route.Role &&
			candidate.Route.RouteType == preferred.Route.RouteType &&
			candidate.Route.PackageVersion == preferred.Route.PackageVersion {
			continue
		}
		explanation.Winner = runtime.explainWinningComparison(protocolPCID, *preferred, candidate)
		break
	}
	if len(explanation.Winner) == 0 {
		explanation.Winner = []string{"the preferred route remained first after deterministic candidate ordering"}
	}
	explanation.Comparisons = runtime.explainAllRouteComparisons(protocolPCID, candidates)
	return explanation
}

func (runtime *Runtime) explainRoutePlanCandidate(protocolPCID string, candidate RoutePlanCandidate) RoutePlanCandidateExplanation {
	globalPolicy := runtime.RoutePlanPolicy()
	protocolPolicy, hasProtocolPolicy := runtime.ProtocolRoutePlanPolicy(protocolPCID)
	rolePolicy, hasRolePolicy := runtime.ProtocolRoleRoutePlanPolicy(protocolPCID, candidate.Route.Role)
	effectivePolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, candidate.Route.Role)
	explanation := RoutePlanCandidateExplanation{
		GlobalPolicy:      globalPolicy,
		EffectivePolicy:   effectivePolicy,
		PreferredByPolicy: routePreferred(candidate.Route, effectivePolicy),
		AvoidedByPolicy:   routeAvoided(candidate.Route, effectivePolicy),
	}
	if hasProtocolPolicy {
		explanation.ProtocolPolicy = &protocolPolicy
	}
	if hasRolePolicy {
		explanation.RolePolicy = &rolePolicy
	}
	explanation.Downstream = explainDownstreamPlans(candidate.Route, candidate.Next)
	explanation.Notes = routeExecutabilityNotes(candidate)
	if explanation.PreferredByPolicy {
		explanation.Notes = append(explanation.Notes, "effective policy prefers this route's type or role")
	}
	if explanation.AvoidedByPolicy {
		explanation.Notes = append(explanation.Notes, "effective policy avoids this route's type or role")
	}
	if !explanation.PreferredByPolicy && !explanation.AvoidedByPolicy {
		explanation.Notes = append(explanation.Notes, "effective policy is neutral for this route")
	}
	return explanation
}

// Intent: Make multi-hop route plans readable end to end by summarizing how
// each downstream emitted protocol resolved, including its nested winner and
// pairwise comparison detail. Source: DI-povak
func explainDownstreamPlans(route grid.RouteRegistration, next []RoutePlan) []RoutePlanDownstreamExplanation {
	out := []RoutePlanDownstreamExplanation{}
	for index, plan := range next {
		explanation := RoutePlanDownstreamExplanation{
			ProtocolPCID: plan.ProtocolPCID,
			Executable:   plan.Preferred != nil,
		}
		if plan.Preferred != nil {
			explanation.PreferredRoute = comparisonSideID(plan.Preferred.Route)
		}
		if plan.Explanation != nil {
			if plan.Explanation.TraceSummary != nil {
				traceSummary := *plan.Explanation.TraceSummary
				traceSummary.Scope = "downstream"
				traceSummary.HopPath = downstreamHopPath("root", route, 0, plan.ProtocolPCID, index)
				traceSummary.HopSummary = downstreamHopSummary(route, 0, plan.ProtocolPCID, index)
				traceSummary.HopDepth = 1
				traceSummary.HopIndex = index + 1
				explanation.TraceSummary = &traceSummary
			}
			explanation.Winner = append([]string{}, plan.Explanation.Winner...)
			explanation.Comparisons = append([]RoutePlanComparison{}, plan.Explanation.Comparisons...)
		}
		if plan.Preferred == nil {
			explanation.Notes = append(explanation.Notes, "no executable downstream route was available")
		} else {
			explanation.Notes = append(explanation.Notes, "downstream protocol resolved to a preferred nested route")
		}
		out = append(out, explanation)
	}
	return out
}

func routeExecutabilityNotes(candidate RoutePlanCandidate) []string {
	switch candidate.Route.RouteType {
	case "", "direct":
		if candidate.Executable {
			return []string{"direct routes are executable immediately"}
		}
		return []string{"direct route was marked non-executable unexpectedly"}
	case "parser", "transform":
		if len(candidate.Route.EmitsProtocols) == 0 {
			return []string{"route is not executable because it emits no downstream protocols"}
		}
		if candidate.Executable {
			return []string{"all emitted protocols resolved to executable downstream plans"}
		}
		return []string{"route is not executable because at least one emitted protocol has no preferred downstream plan"}
	default:
		return []string{"route type is unknown to the planner and is treated as non-executable"}
	}
}

func (runtime *Runtime) explainWinningComparison(protocolPCID string, winner RoutePlanCandidate, loser RoutePlanCandidate) []string {
	reasons := []string{}
	if winner.Executable != loser.Executable {
		if winner.Executable {
			reasons = append(reasons, "winner is executable while the next candidate is not")
		} else {
			reasons = append(reasons, "winner remained ahead despite executability parity")
		}
	}
	winnerPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, winner.Route.Role)
	loserPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, loser.Route.Role)
	winnerAvoided := routeAvoided(winner.Route, winnerPolicy)
	loserAvoided := routeAvoided(loser.Route, loserPolicy)
	if winnerAvoided != loserAvoided {
		if !winnerAvoided && loserAvoided {
			reasons = append(reasons, "winner is not avoided by its effective policy while the next candidate is avoided")
		} else {
			reasons = append(reasons, "winner remained ahead despite being avoided")
		}
	}
	winnerPreferred := routePreferred(winner.Route, winnerPolicy)
	loserPreferred := routePreferred(loser.Route, loserPolicy)
	if winnerPreferred != loserPreferred {
		if winnerPreferred && !loserPreferred {
			reasons = append(reasons, "winner is preferred by its effective policy while the next candidate is neutral")
		} else {
			reasons = append(reasons, "winner remained ahead without a preference advantage")
		}
	}
	if diff := compareRouteType(winner.Route.RouteType, loser.Route.RouteType); diff != 0 {
		reasons = append(reasons, "route type ordering broke the tie")
	}
	if diff := compareRouteRole(winner.Route.Role, loser.Route.Role); diff != 0 {
		reasons = append(reasons, "route role ordering broke the tie")
	}
	if winner.Route.PackageID != loser.Route.PackageID || winner.Route.PackageVersion != loser.Route.PackageVersion {
		reasons = append(reasons, "package identity broke the remaining tie")
	}
	if len(winner.Next) != len(loser.Next) {
		reasons = append(reasons, "downstream hop count broke the final tie")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "winner remained first after deterministic planner comparison")
	}
	return reasons
}

func (runtime *Runtime) explainAllRouteComparisons(protocolPCID string, candidates []RoutePlanCandidate) []RoutePlanComparison {
	comparisons := []RoutePlanComparison{}
	for leftIndex := 0; leftIndex < len(candidates); leftIndex++ {
		for rightIndex := leftIndex + 1; rightIndex < len(candidates); rightIndex++ {
			left := candidates[leftIndex]
			right := candidates[rightIndex]
			comparison := RoutePlanComparison{
				Left:         comparisonSide(left.Route),
				Right:        comparisonSide(right.Route),
				DecisionPath: runtime.explainPairwiseComparison(protocolPCID, left, right),
			}
			if runtime.compareRoutePlanCandidates(protocolPCID, left, right, nil) <= 0 {
				comparison.Winner = comparisonSideID(left.Route)
			} else {
				comparison.Winner = comparisonSideID(right.Route)
			}
			comparisons = append(comparisons, comparison)
		}
	}
	return comparisons
}

func (runtime *Runtime) explainPairwiseComparison(protocolPCID string, left RoutePlanCandidate, right RoutePlanCandidate) []string {
	reasons := []string{}
	if left.Executable != right.Executable {
		if left.Executable {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because it is executable and "+comparisonSideID(right.Route)+" is not")
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because it is executable and "+comparisonSideID(left.Route)+" is not")
		}
		return reasons
	}
	leftPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, left.Route.Role)
	rightPolicy := runtime.EffectiveRoutePlanPolicyForRole(protocolPCID, right.Route.Role)
	leftAvoided := routeAvoided(left.Route, leftPolicy)
	rightAvoided := routeAvoided(right.Route, rightPolicy)
	if leftAvoided != rightAvoided {
		if !leftAvoided && rightAvoided {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because its effective policy does not avoid it while "+comparisonSideID(right.Route)+" is avoided")
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because its effective policy does not avoid it while "+comparisonSideID(left.Route)+" is avoided")
		}
		return reasons
	}
	leftPreferred := routePreferred(left.Route, leftPolicy)
	rightPreferred := routePreferred(right.Route, rightPolicy)
	if leftPreferred != rightPreferred {
		if leftPreferred {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because its effective policy prefers it and "+comparisonSideID(right.Route)+" is neutral")
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because its effective policy prefers it and "+comparisonSideID(left.Route)+" is neutral")
		}
		return reasons
	}
	if diff := compareRouteType(left.Route.RouteType, right.Route.RouteType); diff != 0 {
		if diff < 0 {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because route type ordering prefers "+left.Route.RouteType+" over "+right.Route.RouteType)
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because route type ordering prefers "+right.Route.RouteType+" over "+left.Route.RouteType)
		}
		return reasons
	}
	if diff := compareRouteRole(left.Route.Role, right.Route.Role); diff != 0 {
		if diff < 0 {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because route role ordering prefers "+left.Route.Role+" over "+right.Route.Role)
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because route role ordering prefers "+right.Route.Role+" over "+left.Route.Role)
		}
		return reasons
	}
	if diff := strings.Compare(left.Route.PackageID, right.Route.PackageID); diff != 0 {
		if diff < 0 {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because package ID ordering broke the tie")
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because package ID ordering broke the tie")
		}
		return reasons
	}
	if diff := strings.Compare(left.Route.PackageVersion, right.Route.PackageVersion); diff != 0 {
		if diff < 0 {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because package version ordering broke the tie")
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because package version ordering broke the tie")
		}
		return reasons
	}
	if len(left.Next) != len(right.Next) {
		if len(left.Next) < len(right.Next) {
			reasons = append(reasons, comparisonSideID(left.Route)+" ranked ahead because it has fewer downstream hops")
		} else {
			reasons = append(reasons, comparisonSideID(right.Route)+" ranked ahead because it has fewer downstream hops")
		}
		return reasons
	}
	reasons = append(reasons, "the pair remained tied after all deterministic comparison steps")
	return reasons
}

func comparisonSide(route grid.RouteRegistration) RoutePlanComparisonSide {
	return RoutePlanComparisonSide{
		PackageID:      route.PackageID,
		PackageVersion: route.PackageVersion,
		ProtocolPCID:   route.ProtocolPCID,
		Role:           route.Role,
		RouteType:      route.RouteType,
	}
}

func comparisonSideID(route grid.RouteRegistration) string {
	return route.PackageID + ":" + route.Role + ":" + route.RouteType
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
