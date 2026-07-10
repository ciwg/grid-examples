package carrier

import (
	"context"
	"fmt"
	"time"

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
	if err := client.RegisterProfiles(ctx, []cid.Cid{protocol.ShipmentProfile.CID}); err != nil {
		return err
	}
	for {
		var request protocol.ShipmentMessage
		_, requestCID, err := agent.ReceiveTyped(ctx, client, "seller", protocol.ShipmentProfile, &request)
		if err != nil {
			return err
		}
		if request.Kind != "request" {
			continue
		}
		if err := agent.VerifyMessageCapability(request.CapabilityToken, "seller", "carrier", protocol.ShipmentProfile, request.Kind); err != nil {
			return err
		}
		if request.CustomerOrderRef == "demo-carrier-timeout" {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(cfg.Timeout + 2*time.Second):
			}
		}
		result := protocol.ShipmentMessage{
			Kind:                "result",
			CustomerOrderRef:    request.CustomerOrderRef,
			ParentOrderCID:      request.ParentOrderCID,
			ParentPickPackCID:   request.ParentPickPackCID,
			ParentAccountingCID: request.ParentAccountingCID,
			PackageID:           request.PackageID,
			WeightGrams:         request.WeightGrams,
			Notes:               "carrier processed request",
		}
		result.Status = "booked"
		result.CarrierName = "demo-carrier"
		result.TrackingNumber = "TRACK-" + requestCID[len(requestCID)-10:]
		capabilityToken, err := agent.IssueMessageCapability("carrier", "seller", protocol.ShipmentProfile, result.Kind)
		if err != nil {
			return err
		}
		result.CapabilityToken = capabilityToken
		payloadBytes, err := protocol.Marshal(result)
		if err != nil {
			return fmt.Errorf("marshal shipment result: %w", err)
		}
		envelope, err := agent.BuildSignedEnvelope("carrier", protocol.ShipmentProfile, payloadBytes)
		if err != nil {
			return err
		}
		if _, err := client.Send(ctx, envelope, []string{requestCID}); err != nil {
			return err
		}
	}
}
