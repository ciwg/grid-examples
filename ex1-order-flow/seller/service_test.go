package seller

import (
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

func TestPickPackMatchesRequest(t *testing.T) {
	request := protocol.PickPackMessage{
		CustomerOrderRef: "demo-001",
		ParentOrderCID:   []byte("order-a"),
	}
	matching := protocol.PickPackMessage{
		CustomerOrderRef: "demo-001",
		ParentOrderCID:   []byte("order-a"),
	}
	mismatched := protocol.PickPackMessage{
		CustomerOrderRef: "demo-002",
		ParentOrderCID:   []byte("order-a"),
	}
	if !pickPackMatchesRequest(request, matching) {
		t.Fatal("expected matching pick-pack response to be accepted")
	}
	if pickPackMatchesRequest(request, mismatched) {
		t.Fatal("expected mismatched pick-pack response to be rejected")
	}
}

func TestAccountingMatchesRequest(t *testing.T) {
	request := protocol.AccountingMessage{
		CustomerOrderRef: "demo-001",
		ParentOrderCID:   []byte("order-a"),
	}
	matching := protocol.AccountingMessage{
		CustomerOrderRef: "demo-001",
		ParentOrderCID:   []byte("order-a"),
	}
	mismatched := protocol.AccountingMessage{
		CustomerOrderRef: "demo-001",
		ParentOrderCID:   []byte("order-b"),
	}
	if !accountingMatchesRequest(request, matching) {
		t.Fatal("expected matching accounting response to be accepted")
	}
	if accountingMatchesRequest(request, mismatched) {
		t.Fatal("expected mismatched accounting response to be rejected")
	}
}

func TestShipmentMatchesRequest(t *testing.T) {
	request := protocol.ShipmentMessage{
		CustomerOrderRef:    "demo-001",
		ParentOrderCID:      []byte("order-a"),
		ParentPickPackCID:   []byte("pick-a"),
		ParentAccountingCID: []byte("acct-a"),
		PackageID:           "PKG-001",
	}
	matching := protocol.ShipmentMessage{
		CustomerOrderRef:    "demo-001",
		ParentOrderCID:      []byte("order-a"),
		ParentPickPackCID:   []byte("pick-a"),
		ParentAccountingCID: []byte("acct-a"),
		PackageID:           "PKG-001",
	}
	mismatched := protocol.ShipmentMessage{
		CustomerOrderRef:    "demo-001",
		ParentOrderCID:      []byte("order-a"),
		ParentPickPackCID:   []byte("pick-a"),
		ParentAccountingCID: []byte("acct-a"),
		PackageID:           "PKG-002",
	}
	if !shipmentMatchesRequest(request, matching) {
		t.Fatal("expected matching shipment response to be accepted")
	}
	if shipmentMatchesRequest(request, mismatched) {
		t.Fatal("expected mismatched shipment response to be rejected")
	}
}
