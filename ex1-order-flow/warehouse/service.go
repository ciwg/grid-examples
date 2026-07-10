package warehouse

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/protocol"
)

func Run(ctx context.Context, cfg agent.Config) (runErr error) {
	client, err := agent.DialKernel(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil && runErr == nil {
			runErr = closeErr
		}
	}()
	if err := client.RegisterProfiles(ctx, []cid.Cid{protocol.PickPackProfile.CID}); err != nil {
		return err
	}
	for {
		var request protocol.PickPackMessage
		_, requestCID, err := agent.ReceiveTyped(ctx, client, "seller", protocol.PickPackProfile, &request)
		if err != nil {
			return err
		}
		if request.Kind != "request" {
			continue
		}
		if err := agent.VerifyMessageCapability(request.CapabilityToken, "seller", "warehouse", protocol.PickPackProfile, request.Kind); err != nil {
			return err
		}
		result := protocol.PickPackMessage{
			Kind:             "result",
			CustomerOrderRef: request.CustomerOrderRef,
			ParentOrderCID:   request.ParentOrderCID,
			Notes:            "warehouse processed request",
		}
		switch request.Items[0].SKU {
		case "widget-oos":
			result.Status = "refused"
			result.Notes = "warehouse refusal for out-of-stock widget"
		default:
			result.Status = "packed"
			result.PackageID = "PKG-" + requestCID[len(requestCID)-8:]
			result.WeightGrams = 1250
			result.PackageCount = 1
		}
		capabilityToken, err := agent.IssueMessageCapability("warehouse", "seller", protocol.PickPackProfile, result.Kind)
		if err != nil {
			return err
		}
		result.CapabilityToken = capabilityToken
		payloadBytes, err := protocol.Marshal(result)
		if err != nil {
			return fmt.Errorf("marshal warehouse result: %w", err)
		}
		envelope, err := agent.BuildSignedEnvelope("warehouse", protocol.PickPackProfile, payloadBytes)
		if err != nil {
			return err
		}
		if _, err := client.Send(ctx, envelope, []string{requestCID}); err != nil {
			return err
		}
	}
}
