package accounting

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
	if err := client.RegisterProfiles(ctx, []cid.Cid{protocol.AccountingProfile.CID}); err != nil {
		return err
	}
	for {
		var request protocol.AccountingMessage
		_, requestCID, err := agent.ReceiveTyped(ctx, client, "seller", protocol.AccountingProfile, &request)
		if err != nil {
			return err
		}
		if request.Kind != "request" {
			continue
		}
		if err := agent.VerifyMessageCapability(request.CapabilityToken, "seller", "accounting", protocol.AccountingProfile, request.Kind); err != nil {
			return err
		}
		result := protocol.AccountingMessage{
			Kind:             "result",
			CustomerOrderRef: request.CustomerOrderRef,
			ParentOrderCID:   request.ParentOrderCID,
			Notes:            "accounting processed request",
		}
		if request.PaymentRef == "pay-demo-dup" {
			result.Status = "refused"
			result.Notes = "duplicate payment reference refused"
		} else {
			result.Status = "recorded"
			result.LedgerEntryID = "LEDGER-" + requestCID[len(requestCID)-8:]
			result.InvoiceRef = "INV-" + requestCID[len(requestCID)-8:]
		}
		capabilityToken, err := agent.IssueMessageCapability("accounting", "seller", protocol.AccountingProfile, result.Kind)
		if err != nil {
			return err
		}
		result.CapabilityToken = capabilityToken
		payloadBytes, err := protocol.Marshal(result)
		if err != nil {
			return fmt.Errorf("marshal accounting result: %w", err)
		}
		envelope, err := agent.BuildSignedEnvelope("accounting", protocol.AccountingProfile, payloadBytes)
		if err != nil {
			return err
		}
		if _, err := client.Send(ctx, envelope, []string{requestCID}); err != nil {
			return err
		}
	}
}
