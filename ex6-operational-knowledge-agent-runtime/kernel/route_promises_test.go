package kernel

import (
	"path/filepath"
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/store"
)

func TestRoutePromisesRequireExplicitBindingReceiveAndDelivery(t *testing.T) {
	cas, err := store.OpenCAS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	registry, err := OpenRoutePromiseRegistry(cas)
	if err != nil {
		t.Fatal(err)
	}
	if registry.routeExecutable("inventory", "pcid:inventory.v1") {
		t.Fatal("package route became executable without explicit local evidence")
	}
	if err := registry.BindAgent(AgentBinding{AgentID: "inventory-app", PackageID: "inventory", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if registry.routeExecutable("inventory", "pcid:inventory.v1") {
		t.Fatal("binding alone made a route executable")
	}
	if err := registry.PublishReceivePromise(ReceivePromise{AgentID: "inventory-app", ProtocolPCID: "pcid:inventory.v1", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if registry.routeExecutable("inventory", "pcid:inventory.v1") {
		t.Fatal("receive promise alone made a route executable")
	}
	if err := registry.PublishDeliveryPromise(DeliveryPromise{AgentID: "local-router", RecipientAgentID: "inventory-app", ProtocolPCID: "pcid:inventory.v1", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if !registry.routeExecutable("inventory", "pcid:inventory.v1") {
		t.Fatal("complete explicit route evidence was not executable")
	}
}

func TestRoutePromiseReplayRetainsDisabledState(t *testing.T) {
	root := t.TempDir()
	cas, err := store.OpenCAS(filepath.Join(root, "cas"))
	if err != nil {
		t.Fatal(err)
	}
	registry, err := OpenRoutePromiseRegistry(cas)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.BindAgent(AgentBinding{AgentID: "inventory-app", PackageID: "inventory", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := registry.PublishReceivePromise(ReceivePromise{AgentID: "inventory-app", ProtocolPCID: "pcid:inventory.v1", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := registry.PublishDeliveryPromise(DeliveryPromise{AgentID: "local-router", RecipientAgentID: "inventory-app", ProtocolPCID: "pcid:inventory.v1", Enabled: true}); err != nil {
		t.Fatal(err)
	}
	if err := registry.PublishReceivePromise(ReceivePromise{AgentID: "inventory-app", ProtocolPCID: "pcid:inventory.v1", Enabled: false}); err != nil {
		t.Fatal(err)
	}
	reopened, err := OpenRoutePromiseRegistry(cas)
	if err != nil {
		t.Fatal(err)
	}
	if reopened.routeExecutable("inventory", "pcid:inventory.v1") {
		t.Fatal("disabled receive promise became executable after CAS replay")
	}
}
