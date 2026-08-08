package kernel

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
	"github.com/ipfs/go-cid"
)

const RoutePromisesProtocolPCID = "bafkreieipf4oplqi2zwlw6igu5w3jx444gyq6yi5pqoa6ctulkyx5kxuue"

var routePromisesProtocolCID = cid.MustParse(RoutePromisesProtocolPCID)

type AgentBinding struct {
	AgentID   string `json:"agent_id"`
	PackageID string `json:"package_id"`
	Enabled   bool   `json:"enabled"`
}

type ReceivePromise struct {
	AgentID      string `json:"agent_id"`
	ProtocolPCID string `json:"protocol_pcid"`
	Enabled      bool   `json:"enabled"`
}

type DeliveryPromise struct {
	AgentID          string `json:"agent_id"`
	RecipientAgentID string `json:"recipient_agent_id"`
	ProtocolPCID     string `json:"protocol_pcid"`
	Enabled          bool   `json:"enabled"`
}

type routePromiseKind uint64

const (
	agentBindingKind routePromiseKind = iota
	receivePromiseKind
	deliveryPromiseKind
)

// routePromiseRecord is the retained, pCID-selected lifecycle evidence for one
// local binding or promise. Intent: Rebuild usable route evidence solely from
// immutable CAS records rather than from package installation or a mutable
// registry. Source: DI-kojab; DI-butam
type routePromiseRecord struct {
	Kind             routePromiseKind
	AgentID          string
	ProtocolPCID     string
	RecipientAgentID string
	PackageID        string
	Enabled          bool
	Parents          []cid.Cid
}

func (record routePromiseRecord) key() string {
	switch record.Kind {
	case agentBindingKind:
		return record.AgentID
	case receivePromiseKind:
		return record.AgentID + "\x00" + record.ProtocolPCID
	default:
		return record.AgentID + "\x00" + record.RecipientAgentID + "\x00" + record.ProtocolPCID
	}
}

func (record routePromiseRecord) validate() error {
	if strings.TrimSpace(record.AgentID) == "" {
		return errors.New("route promise agent ID is required")
	}
	if len(record.Parents) > 1 {
		return errors.New("route promise record has too many parents")
	}
	for _, parent := range record.Parents {
		if parent.Version() != 1 {
			return errors.New("route promise parent must be CIDv1")
		}
	}
	switch record.Kind {
	case agentBindingKind:
		if strings.TrimSpace(record.PackageID) == "" || record.ProtocolPCID != "" || record.RecipientAgentID != "" {
			return errors.New("agent binding fields are invalid")
		}
	case receivePromiseKind:
		if strings.TrimSpace(record.ProtocolPCID) == "" || record.RecipientAgentID != "" || record.PackageID != "" {
			return errors.New("receive promise fields are invalid")
		}
	case deliveryPromiseKind:
		if strings.TrimSpace(record.ProtocolPCID) == "" || strings.TrimSpace(record.RecipientAgentID) == "" || record.PackageID != "" {
			return errors.New("delivery promise fields are invalid")
		}
	default:
		return errors.New("route promise kind is invalid")
	}
	return nil
}

func encodeRoutePromiseRecord(record routePromiseRecord) ([]byte, error) {
	if err := record.validate(); err != nil {
		return nil, err
	}
	parents := make([]any, 0, len(record.Parents))
	for _, parent := range record.Parents {
		parents = append(parents, parent.Bytes())
	}
	return records.EncodeGrid(records.GridEnvelope{ProtocolPCID: routePromisesProtocolCID, Slots: []any{uint64(record.Kind), record.AgentID, record.ProtocolPCID, nullableText(record.RecipientAgentID), nullableText(record.PackageID), record.Enabled, parents}})
}

func nullableText(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func decodeRoutePromiseRecord(raw []byte) (routePromiseRecord, error) {
	envelope, err := records.DecodeGrid(raw)
	if err != nil {
		return routePromiseRecord{}, err
	}
	if envelope.ProtocolPCID != routePromisesProtocolCID || len(envelope.Slots) != 7 {
		return routePromiseRecord{}, errors.New("grid envelope does not select the route promises protocol")
	}
	kind, ok := envelope.Slots[0].(uint64)
	if !ok {
		return routePromiseRecord{}, errors.New("route promise kind must be an unsigned integer")
	}
	agentID, ok := envelope.Slots[1].(string)
	if !ok {
		return routePromiseRecord{}, errors.New("route promise agent ID must be text")
	}
	protocolPCID, ok := envelope.Slots[2].(string)
	if !ok {
		return routePromiseRecord{}, errors.New("route promise protocol pCID must be text")
	}
	recipient, err := nullableTextValue(envelope.Slots[3])
	if err != nil {
		return routePromiseRecord{}, fmt.Errorf("route promise recipient: %w", err)
	}
	packageID, err := nullableTextValue(envelope.Slots[4])
	if err != nil {
		return routePromiseRecord{}, fmt.Errorf("route promise package ID: %w", err)
	}
	enabled, ok := envelope.Slots[5].(bool)
	if !ok {
		return routePromiseRecord{}, errors.New("route promise enabled must be boolean")
	}
	values, ok := envelope.Slots[6].([]any)
	if !ok {
		return routePromiseRecord{}, errors.New("route promise parents must be an array")
	}
	parents := make([]cid.Cid, 0, len(values))
	for _, value := range values {
		bytes, ok := value.([]byte)
		if !ok {
			return routePromiseRecord{}, errors.New("route promise parent must be CID bytes")
		}
		parent, err := cid.Cast(bytes)
		if err != nil {
			return routePromiseRecord{}, fmt.Errorf("route promise parent CID: %w", err)
		}
		parents = append(parents, parent)
	}
	record := routePromiseRecord{Kind: routePromiseKind(kind), AgentID: agentID, ProtocolPCID: protocolPCID, RecipientAgentID: recipient, PackageID: packageID, Enabled: enabled, Parents: parents}
	return record, record.validate()
}

func nullableTextValue(value any) (string, error) {
	if value == nil {
		return "", nil
	}
	text, ok := value.(string)
	if !ok {
		return "", errors.New("must be text or null")
	}
	return text, nil
}

type RoutePromiseRegistry struct {
	cas        *store.CAS
	mu         sync.RWMutex
	bindings   map[string]AgentBinding
	receives   map[string]ReceivePromise
	deliveries map[string]DeliveryPromise
	heads      map[string]cid.Cid
	conflicted map[string]bool
}

func OpenRoutePromiseRegistry(cas *store.CAS) (*RoutePromiseRegistry, error) {
	registry := &RoutePromiseRegistry{cas: cas}
	if err := registry.rebuild(); err != nil {
		return nil, err
	}
	return registry, nil
}

func (registry *RoutePromiseRegistry) rebuild() error {
	ids, err := registry.cas.ListCIDs()
	if err != nil {
		return err
	}
	candidates := map[string]routePromiseRecord{}
	for _, id := range ids {
		raw, err := registry.cas.GetCID(id)
		if err != nil {
			continue
		}
		record, err := decodeRoutePromiseRecord(raw)
		if err == nil {
			candidates[id.String()] = record
		}
	}
	accepted := map[string]routePromiseRecord{}
	for progress := true; progress; {
		progress = false
		for _, id := range ids {
			key := id.String()
			record, ok := candidates[key]
			if !ok {
				continue
			}
			if len(record.Parents) == 0 {
				accepted[key] = record
				delete(candidates, key)
				progress = true
				continue
			}
			parent, ok := accepted[record.Parents[0].String()]
			if !ok {
				continue
			}
			delete(candidates, key)
			progress = true
			if parent.Kind == record.Kind && parent.key() == record.key() {
				accepted[key] = record
			}
		}
	}
	children := map[string]int{}
	for _, record := range accepted {
		for _, parent := range record.Parents {
			children[parent.String()]++
		}
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.bindings = map[string]AgentBinding{}
	registry.receives = map[string]ReceivePromise{}
	registry.deliveries = map[string]DeliveryPromise{}
	registry.heads = map[string]cid.Cid{}
	registry.conflicted = map[string]bool{}
	for _, id := range ids {
		key := id.String()
		record, ok := accepted[key]
		if !ok || children[key] != 0 {
			continue
		}
		logicalKey := record.key()
		if _, exists := registry.heads[logicalKey]; exists {
			registry.conflicted[logicalKey] = true
			continue
		}
		registry.heads[logicalKey] = id
		switch record.Kind {
		case agentBindingKind:
			registry.bindings[logicalKey] = AgentBinding{AgentID: record.AgentID, PackageID: record.PackageID, Enabled: record.Enabled}
		case receivePromiseKind:
			registry.receives[logicalKey] = ReceivePromise{AgentID: record.AgentID, ProtocolPCID: record.ProtocolPCID, Enabled: record.Enabled}
		case deliveryPromiseKind:
			registry.deliveries[logicalKey] = DeliveryPromise{AgentID: record.AgentID, RecipientAgentID: record.RecipientAgentID, ProtocolPCID: record.ProtocolPCID, Enabled: record.Enabled}
		}
	}
	for key := range registry.conflicted {
		delete(registry.heads, key)
		delete(registry.bindings, key)
		delete(registry.receives, key)
		delete(registry.deliveries, key)
	}
	return nil
}

func (registry *RoutePromiseRegistry) BindAgent(binding AgentBinding) error {
	return registry.append(routePromiseRecord{Kind: agentBindingKind, AgentID: binding.AgentID, PackageID: binding.PackageID, Enabled: binding.Enabled})
}

func (registry *RoutePromiseRegistry) PublishReceivePromise(promise ReceivePromise) error {
	return registry.append(routePromiseRecord{Kind: receivePromiseKind, AgentID: promise.AgentID, ProtocolPCID: promise.ProtocolPCID, Enabled: promise.Enabled})
}

func (registry *RoutePromiseRegistry) PublishDeliveryPromise(promise DeliveryPromise) error {
	return registry.append(routePromiseRecord{Kind: deliveryPromiseKind, AgentID: promise.AgentID, RecipientAgentID: promise.RecipientAgentID, ProtocolPCID: promise.ProtocolPCID, Enabled: promise.Enabled})
}

func (registry *RoutePromiseRegistry) append(record routePromiseRecord) error {
	if err := record.validate(); err != nil {
		return err
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	key := record.key()
	if registry.conflicted[key] {
		return errors.New("route promise has competing heads")
	}
	if head, ok := registry.heads[key]; ok {
		record.Parents = []cid.Cid{head}
	}
	raw, err := encodeRoutePromiseRecord(record)
	if err != nil {
		return err
	}
	id, err := registry.cas.PutCID(raw)
	if err != nil {
		return err
	}
	registry.heads[key] = id
	switch record.Kind {
	case agentBindingKind:
		registry.bindings[key] = AgentBinding{AgentID: record.AgentID, PackageID: record.PackageID, Enabled: record.Enabled}
	case receivePromiseKind:
		registry.receives[key] = ReceivePromise{AgentID: record.AgentID, ProtocolPCID: record.ProtocolPCID, Enabled: record.Enabled}
	case deliveryPromiseKind:
		registry.deliveries[key] = DeliveryPromise{AgentID: record.AgentID, RecipientAgentID: record.RecipientAgentID, ProtocolPCID: record.ProtocolPCID, Enabled: record.Enabled}
	}
	return nil
}

func (registry *RoutePromiseRegistry) routeExecutable(packageID string, protocolPCID string) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	for _, binding := range registry.bindings {
		if !binding.Enabled || binding.PackageID != packageID {
			continue
		}
		receive, ok := registry.receives[receivePromiseKindKey(binding.AgentID, protocolPCID)]
		if !ok || !receive.Enabled {
			continue
		}
		for _, delivery := range registry.deliveries {
			if delivery.Enabled && delivery.RecipientAgentID == binding.AgentID && delivery.ProtocolPCID == protocolPCID {
				return true
			}
		}
	}
	return false
}

func receivePromiseKindKey(agentID string, protocolPCID string) string {
	return agentID + "\x00" + protocolPCID
}
