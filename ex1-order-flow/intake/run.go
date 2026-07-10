package intake

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

type Result struct {
	CustomerOrderRef string `json:"customer_order_ref"`
	OrderStatus      string `json:"order_status"`
	PromiseStatus    string `json:"promise_status"`
	FailureStage     string `json:"failure_stage"`
	PackageID        string `json:"package_id"`
	TrackingNumber   string `json:"tracking_number"`
	LedgerEntryID    string `json:"ledger_entry_id"`
	FinalOrderCID    string `json:"final_order_cid"`
	Notes            string `json:"notes"`
}

func Run(ctx context.Context, cfg agent.Config, fixturePath string) (result Result, err error) {
	client, err := agent.DialKernel(ctx, cfg)
	if err != nil {
		return Result{}, err
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if err := client.RegisterProfiles(ctx, []cid.Cid{protocol.OrderProfile.CID}); err != nil {
		return Result{}, err
	}
	fixtureBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		return Result{}, fmt.Errorf("read fixture: %w", err)
	}
	var submit protocol.OrderMessage
	if err := json.Unmarshal(fixtureBytes, &submit); err != nil {
		return Result{}, fmt.Errorf("decode fixture: %w", err)
	}
	submit.Kind = "submit"
	capabilityToken, err := agent.IssueMessageCapability("intake", "seller", protocol.OrderProfile, submit.Kind)
	if err != nil {
		return Result{}, err
	}
	submit.CapabilityToken = capabilityToken
	payloadBytes, err := protocol.Marshal(submit)
	if err != nil {
		return Result{}, err
	}
	envelope, err := agent.BuildSignedEnvelope("intake", protocol.OrderProfile, payloadBytes)
	if err != nil {
		return Result{}, err
	}
	if _, err = client.Send(ctx, envelope, nil); err != nil {
		return Result{}, err
	}
	var final protocol.OrderMessage
	_, finalCID, err := agent.ReceiveTyped(ctx, client, "seller", protocol.OrderProfile, &final)
	if err != nil {
		return Result{}, err
	}
	if err := agent.VerifyMessageCapability(final.CapabilityToken, "seller", "intake", protocol.OrderProfile, final.Kind); err != nil {
		return Result{}, err
	}
	result = Result{
		CustomerOrderRef: final.CustomerOrderRef,
		OrderStatus:      final.OrderStatus,
		FailureStage:     final.FailureStage,
		PackageID:        final.PackageID,
		TrackingNumber:   final.TrackingNumber,
		LedgerEntryID:    final.LedgerEntryID,
		FinalOrderCID:    finalCID,
		Notes:            "conforming signed final order message received for this request",
	}
	if final.Kind != "final" {
		result.PromiseStatus = "broken"
		result.Notes = "non-final order message received"
	} else {
		result.PromiseStatus = "kept"
	}
	resultPath := filepath.Join(cfg.DataDir, "result.json")
	resultBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(resultPath, append(resultBytes, '\n'), 0o644); err != nil {
		return Result{}, err
	}
	return result, nil
}
