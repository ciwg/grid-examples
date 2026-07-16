package seller

import (
	"bytes"
	"context"

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
	if err := client.RegisterProfiles(ctx, []cid.Cid{
		protocol.OrderProfile.CID,
		protocol.PickPackProfile.CID,
		protocol.AccountingProfile.CID,
		protocol.ShipmentProfile.CID,
	}); err != nil {
		return err
	}
	for {
		if err := handleOneOrder(ctx, client, cfg); err != nil {
			return err
		}
	}
}

func handleOneOrder(ctx context.Context, client *agent.Client, cfg agent.Config) error {
	var submit protocol.OrderMessage
	_, requestCID, err := agent.ReceiveTyped(ctx, client, "intake", protocol.OrderProfile, &submit)
	if err != nil {
		return err
	}
	if submit.Kind != "submit" {
		return nil
	}
	// Intent: Keep the seller service subscribed indefinitely, but bound the
	// downstream warehouse/accounting/carrier work for each concrete order so
	// timeout handling remains per-order rather than per-process. Source: DI-kozod; DI-rokol
	orderCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	if err := agent.VerifyMessageCapability(submit.CapabilityToken, "intake", "seller", protocol.OrderProfile, submit.Kind); err != nil {
		return sendFailure(ctx, client, submit, requestCID, "seller_validation", "failed", "invalid intake capability token")
	}
	// Intent: Fail malformed intake orders before contacting downstream services so
	// demo runs stay deterministic and an empty item list cannot crash warehouse.
	// Source: DI-seller-validation
	if len(submit.Items) == 0 {
		return sendFailure(ctx, client, submit, requestCID, "seller_validation", "failed", "order must contain at least one item")
	}
	requestCIDParsed, err := cid.Parse(requestCID)
	if err != nil {
		return err
	}
	pickPackCID, pickPackResult, pickPackResultCID, err := runPickPack(orderCtx, client, submit, requestCIDParsed.Bytes())
	if err != nil {
		return sendFailure(ctx, client, submit, requestCID, "warehouse", "failed", err.Error())
	}
	if pickPackResult.Status == "refused" {
		return sendFinal(ctx, client, submit, requestCID, protocol.OrderMessage{
			Kind:               "final",
			CustomerOrderRef:   submit.CustomerOrderRef,
			ParentOrderCID:     requestCIDParsed.Bytes(),
			OrderStatus:        "refused",
			FailureStage:       "warehouse",
			WarehouseResultCID: mustCIDBytes(pickPackResultCID),
			Summary:            "warehouse refused request",
		})
	}
	accountingCID, accountingResult, accountingResultCID, err := runAccounting(orderCtx, client, submit, requestCIDParsed.Bytes())
	if err != nil {
		return sendFailure(ctx, client, submit, requestCID, "accounting", "failed", err.Error())
	}
	if accountingResult.Status == "refused" {
		return sendFinal(ctx, client, submit, requestCID, protocol.OrderMessage{
			Kind:                "final",
			CustomerOrderRef:    submit.CustomerOrderRef,
			ParentOrderCID:      requestCIDParsed.Bytes(),
			OrderStatus:         "refused",
			FailureStage:        "accounting",
			WarehouseResultCID:  mustCIDBytes(pickPackResultCID),
			AccountingResultCID: mustCIDBytes(accountingResultCID),
			PackageID:           pickPackResult.PackageID,
			Summary:             "accounting refused request",
		})
	}
	shipmentResult, shipmentResultCID, err := runShipment(orderCtx, client, submit, requestCIDParsed.Bytes(), pickPackCID.Bytes(), accountingCID.Bytes(), pickPackResult)
	if err != nil {
		return sendFailure(ctx, client, submit, requestCID, "carrier", "failed", err.Error())
	}
	return sendFinal(ctx, client, submit, requestCID, protocol.OrderMessage{
		Kind:                "final",
		CustomerOrderRef:    submit.CustomerOrderRef,
		ParentOrderCID:      requestCIDParsed.Bytes(),
		OrderStatus:         "fulfilled",
		WarehouseResultCID:  mustCIDBytes(pickPackResultCID),
		AccountingResultCID: mustCIDBytes(accountingResultCID),
		ShipmentResultCID:   mustCIDBytes(shipmentResultCID),
		PackageID:           pickPackResult.PackageID,
		TrackingNumber:      shipmentResult.TrackingNumber,
		LedgerEntryID:       accountingResult.LedgerEntryID,
		Summary:             "order fulfilled",
	})
}

func runPickPack(ctx context.Context, client *agent.Client, submit protocol.OrderMessage, parentOrderCID []byte) (cid.Cid, protocol.PickPackMessage, string, error) {
	capabilityToken, err := agent.IssueMessageCapability("seller", "warehouse", protocol.PickPackProfile, "request")
	if err != nil {
		return cid.Undef, protocol.PickPackMessage{}, "", err
	}
	request := protocol.PickPackMessage{
		Kind:             "request",
		CustomerOrderRef: submit.CustomerOrderRef,
		ParentOrderCID:   parentOrderCID,
		Items:            append([]protocol.OrderItem(nil), submit.Items...),
		CapabilityToken:  capabilityToken,
		ServiceLevel:     submit.ServiceLevel,
	}
	payloadBytes, err := protocol.Marshal(request)
	if err != nil {
		return cid.Undef, protocol.PickPackMessage{}, "", err
	}
	envelope, err := agent.BuildSignedEnvelope("seller", protocol.PickPackProfile, payloadBytes)
	if err != nil {
		return cid.Undef, protocol.PickPackMessage{}, "", err
	}
	requestCID, err := client.Send(ctx, envelope, []string{mustCIDText(parentOrderCID)})
	if err != nil {
		return cid.Undef, protocol.PickPackMessage{}, "", err
	}
	for {
		var result protocol.PickPackMessage
		_, resultCID, err := agent.ReceiveTyped(ctx, client, "warehouse", protocol.PickPackProfile, &result)
		if err != nil {
			return cid.Undef, protocol.PickPackMessage{}, "", err
		}
		if err := agent.VerifyMessageCapability(result.CapabilityToken, "warehouse", "seller", protocol.PickPackProfile, result.Kind); err != nil {
			return cid.Undef, protocol.PickPackMessage{}, "", err
		}
		// Intent: Accept only the warehouse result that matches the request just
		// sent so delayed or replayed results from another order cannot satisfy the
		// current order's wait path. Source: DI-seller-correlation
		if !pickPackMatchesRequest(request, result) {
			continue
		}
		requestParsed, err := cid.Parse(requestCID)
		if err != nil {
			return cid.Undef, protocol.PickPackMessage{}, "", err
		}
		return requestParsed, result, resultCID, nil
	}
}

func runAccounting(ctx context.Context, client *agent.Client, submit protocol.OrderMessage, parentOrderCID []byte) (cid.Cid, protocol.AccountingMessage, string, error) {
	capabilityToken, err := agent.IssueMessageCapability("seller", "accounting", protocol.AccountingProfile, "request")
	if err != nil {
		return cid.Undef, protocol.AccountingMessage{}, "", err
	}
	request := protocol.AccountingMessage{
		Kind:             "request",
		CustomerOrderRef: submit.CustomerOrderRef,
		ParentOrderCID:   parentOrderCID,
		PaymentRef:       submit.PaymentRef,
		CapabilityToken:  capabilityToken,
		Currency:         "USD",
		AmountCents:      1999,
	}
	payloadBytes, err := protocol.Marshal(request)
	if err != nil {
		return cid.Undef, protocol.AccountingMessage{}, "", err
	}
	envelope, err := agent.BuildSignedEnvelope("seller", protocol.AccountingProfile, payloadBytes)
	if err != nil {
		return cid.Undef, protocol.AccountingMessage{}, "", err
	}
	requestCID, err := client.Send(ctx, envelope, []string{mustCIDText(parentOrderCID)})
	if err != nil {
		return cid.Undef, protocol.AccountingMessage{}, "", err
	}
	for {
		var result protocol.AccountingMessage
		_, resultCID, err := agent.ReceiveTyped(ctx, client, "accounting", protocol.AccountingProfile, &result)
		if err != nil {
			return cid.Undef, protocol.AccountingMessage{}, "", err
		}
		if err := agent.VerifyMessageCapability(result.CapabilityToken, "accounting", "seller", protocol.AccountingProfile, result.Kind); err != nil {
			return cid.Undef, protocol.AccountingMessage{}, "", err
		}
		// Intent: Keep accounting correlation bound to the current order request so
		// a stale accounting reply cannot satisfy a later order. Source:
		// DI-seller-correlation
		if !accountingMatchesRequest(request, result) {
			continue
		}
		requestParsed, err := cid.Parse(requestCID)
		if err != nil {
			return cid.Undef, protocol.AccountingMessage{}, "", err
		}
		return requestParsed, result, resultCID, nil
	}
}

func runShipment(ctx context.Context, client *agent.Client, submit protocol.OrderMessage, parentOrderCID []byte, pickPackCID []byte, accountingCID []byte, pickPackResult protocol.PickPackMessage) (protocol.ShipmentMessage, string, error) {
	capabilityToken, err := agent.IssueMessageCapability("seller", "carrier", protocol.ShipmentProfile, "request")
	if err != nil {
		return protocol.ShipmentMessage{}, "", err
	}
	request := protocol.ShipmentMessage{
		Kind:                "request",
		CustomerOrderRef:    submit.CustomerOrderRef,
		ParentOrderCID:      parentOrderCID,
		ParentPickPackCID:   pickPackCID,
		ParentAccountingCID: accountingCID,
		CapabilityToken:     capabilityToken,
		PackageID:           pickPackResult.PackageID,
		WeightGrams:         pickPackResult.WeightGrams,
		ShipTo:              submit.ShipTo,
		ServiceLevel:        submit.ServiceLevel,
	}
	payloadBytes, err := protocol.Marshal(request)
	if err != nil {
		return protocol.ShipmentMessage{}, "", err
	}
	envelope, err := agent.BuildSignedEnvelope("seller", protocol.ShipmentProfile, payloadBytes)
	if err != nil {
		return protocol.ShipmentMessage{}, "", err
	}
	parents := []string{
		mustCIDText(parentOrderCID),
		mustCIDText(pickPackCID),
		mustCIDText(accountingCID),
	}
	if _, err := client.Send(ctx, envelope, parents); err != nil {
		return protocol.ShipmentMessage{}, "", err
	}
	for {
		var result protocol.ShipmentMessage
		_, resultCID, err := agent.ReceiveTyped(ctx, client, "carrier", protocol.ShipmentProfile, &result)
		if err != nil {
			return protocol.ShipmentMessage{}, "", err
		}
		if err := agent.VerifyMessageCapability(result.CapabilityToken, "carrier", "seller", protocol.ShipmentProfile, result.Kind); err != nil {
			return protocol.ShipmentMessage{}, "", err
		}
		// Intent: Keep carrier replies correlated to the exact shipment request so
		// delayed booking results cannot be attached to the wrong order.
		// Source: DI-seller-correlation
		if !shipmentMatchesRequest(request, result) {
			continue
		}
		return result, resultCID, nil
	}
}

func sendFailure(ctx context.Context, client *agent.Client, submit protocol.OrderMessage, parentCID string, stage string, status string, summary string) error {
	return sendFinal(ctx, client, submit, parentCID, protocol.OrderMessage{
		Kind:             "final",
		CustomerOrderRef: submit.CustomerOrderRef,
		ParentOrderCID:   mustCIDBytes(parentCID),
		OrderStatus:      status,
		FailureStage:     stage,
		Summary:          summary,
	})
}

func sendFinal(ctx context.Context, client *agent.Client, submit protocol.OrderMessage, parentCID string, final protocol.OrderMessage) error {
	capabilityToken, err := agent.IssueMessageCapability("seller", "intake", protocol.OrderProfile, "final")
	if err != nil {
		return err
	}
	final.CapabilityToken = capabilityToken
	payloadBytes, err := protocol.Marshal(final)
	if err != nil {
		return err
	}
	envelope, err := agent.BuildSignedEnvelope("seller", protocol.OrderProfile, payloadBytes)
	if err != nil {
		return err
	}
	_, err = client.Send(ctx, envelope, []string{parentCID})
	return err
}

func mustCIDBytes(text string) []byte {
	parsed, err := cid.Parse(text)
	if err != nil {
		panic(err)
	}
	return parsed.Bytes()
}

func mustCIDText(raw []byte) string {
	parsed, err := cid.Cast(raw)
	if err != nil {
		panic(err)
	}
	return parsed.String()
}

func pickPackMatchesRequest(request protocol.PickPackMessage, result protocol.PickPackMessage) bool {
	return request.CustomerOrderRef == result.CustomerOrderRef &&
		bytes.Equal(request.ParentOrderCID, result.ParentOrderCID)
}

func accountingMatchesRequest(request protocol.AccountingMessage, result protocol.AccountingMessage) bool {
	return request.CustomerOrderRef == result.CustomerOrderRef &&
		bytes.Equal(request.ParentOrderCID, result.ParentOrderCID)
}

func shipmentMatchesRequest(request protocol.ShipmentMessage, result protocol.ShipmentMessage) bool {
	return request.CustomerOrderRef == result.CustomerOrderRef &&
		bytes.Equal(request.ParentOrderCID, result.ParentOrderCID) &&
		bytes.Equal(request.ParentPickPackCID, result.ParentPickPackCID) &&
		bytes.Equal(request.ParentAccountingCID, result.ParentAccountingCID) &&
		request.PackageID == result.PackageID
}
