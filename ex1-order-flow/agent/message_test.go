package agent

import (
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

func TestVerifyMessageCapabilityRejectsInvalidToken(t *testing.T) {
	tokenBytes, err := IssueMessageCapability("seller", "warehouse", protocol.PickPackProfile, "request")
	if err != nil {
		t.Fatalf("issue capability: %v", err)
	}
	if err := VerifyMessageCapability(tokenBytes, "seller", "accounting", protocol.PickPackProfile, "request"); err == nil {
		t.Fatal("VerifyMessageCapability accepted a token for the wrong audience")
	}
}
