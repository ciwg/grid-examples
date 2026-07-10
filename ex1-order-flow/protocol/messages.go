package protocol

type ShipTo struct {
	Name       string `cbor:"name" json:"name"`
	Address1   string `cbor:"address1" json:"address1"`
	Address2   string `cbor:"address2,omitempty" json:"address2,omitempty"`
	City       string `cbor:"city" json:"city"`
	Region     string `cbor:"region" json:"region"`
	PostalCode string `cbor:"postal_code" json:"postal_code"`
	Country    string `cbor:"country" json:"country"`
}

type OrderItem struct {
	SKU      string `cbor:"sku" json:"sku"`
	Quantity uint64 `cbor:"quantity" json:"quantity"`
}

type OrderMessage struct {
	Kind               string      `cbor:"kind" json:"kind"`
	CustomerOrderRef   string      `cbor:"customer_order_ref" json:"customer_order_ref"`
	RequestedBy        string      `cbor:"requested_by,omitempty" json:"requested_by,omitempty"`
	ShipTo             *ShipTo     `cbor:"ship_to,omitempty" json:"ship_to,omitempty"`
	Items              []OrderItem `cbor:"items,omitempty" json:"items,omitempty"`
	ServiceLevel       string      `cbor:"service_level,omitempty" json:"service_level,omitempty"`
	PaymentRef         string      `cbor:"payment_ref,omitempty" json:"payment_ref,omitempty"`
	CapabilityToken    []byte      `cbor:"capability_token,omitempty" json:"capability_token,omitempty"`
	Notes              string      `cbor:"notes,omitempty" json:"notes,omitempty"`
	ParentOrderCID     []byte      `cbor:"parent_order_cid,omitempty" json:"parent_order_cid,omitempty"`
	OrderStatus        string      `cbor:"order_status,omitempty" json:"order_status,omitempty"`
	FailureStage       string      `cbor:"failure_stage,omitempty" json:"failure_stage,omitempty"`
	WarehouseResultCID []byte      `cbor:"warehouse_result_cid,omitempty" json:"warehouse_result_cid,omitempty"`
	AccountingResultCID []byte     `cbor:"accounting_result_cid,omitempty" json:"accounting_result_cid,omitempty"`
	ShipmentResultCID  []byte      `cbor:"shipment_result_cid,omitempty" json:"shipment_result_cid,omitempty"`
	PackageID          string      `cbor:"package_id,omitempty" json:"package_id,omitempty"`
	TrackingNumber     string      `cbor:"tracking_number,omitempty" json:"tracking_number,omitempty"`
	LedgerEntryID      string      `cbor:"ledger_entry_id,omitempty" json:"ledger_entry_id,omitempty"`
	Summary            string      `cbor:"summary,omitempty" json:"summary,omitempty"`
}

type PickPackMessage struct {
	Kind            string      `cbor:"kind" json:"kind"`
	CustomerOrderRef string     `cbor:"customer_order_ref" json:"customer_order_ref"`
	ParentOrderCID  []byte      `cbor:"parent_order_cid" json:"parent_order_cid"`
	Items           []OrderItem `cbor:"items,omitempty" json:"items,omitempty"`
	CapabilityToken []byte      `cbor:"capability_token,omitempty" json:"capability_token,omitempty"`
	ServiceLevel    string      `cbor:"service_level,omitempty" json:"service_level,omitempty"`
	Status          string      `cbor:"status,omitempty" json:"status,omitempty"`
	PackageID       string      `cbor:"package_id,omitempty" json:"package_id,omitempty"`
	WeightGrams     uint64      `cbor:"weight_grams,omitempty" json:"weight_grams,omitempty"`
	PackageCount    uint64      `cbor:"package_count,omitempty" json:"package_count,omitempty"`
	Notes           string      `cbor:"notes,omitempty" json:"notes,omitempty"`
}

type AccountingMessage struct {
	Kind            string `cbor:"kind" json:"kind"`
	CustomerOrderRef string `cbor:"customer_order_ref" json:"customer_order_ref"`
	ParentOrderCID  []byte `cbor:"parent_order_cid" json:"parent_order_cid"`
	PaymentRef      string `cbor:"payment_ref,omitempty" json:"payment_ref,omitempty"`
	CapabilityToken []byte `cbor:"capability_token,omitempty" json:"capability_token,omitempty"`
	Currency        string `cbor:"currency,omitempty" json:"currency,omitempty"`
	AmountCents     uint64 `cbor:"amount_cents,omitempty" json:"amount_cents,omitempty"`
	Status          string `cbor:"status,omitempty" json:"status,omitempty"`
	LedgerEntryID   string `cbor:"ledger_entry_id,omitempty" json:"ledger_entry_id,omitempty"`
	InvoiceRef      string `cbor:"invoice_ref,omitempty" json:"invoice_ref,omitempty"`
	Notes           string `cbor:"notes,omitempty" json:"notes,omitempty"`
}

type ShipmentMessage struct {
	Kind               string  `cbor:"kind" json:"kind"`
	CustomerOrderRef   string  `cbor:"customer_order_ref" json:"customer_order_ref"`
	ParentOrderCID     []byte  `cbor:"parent_order_cid" json:"parent_order_cid"`
	ParentPickPackCID  []byte  `cbor:"parent_pick_pack_cid,omitempty" json:"parent_pick_pack_cid,omitempty"`
	ParentAccountingCID []byte `cbor:"parent_accounting_cid,omitempty" json:"parent_accounting_cid,omitempty"`
	CapabilityToken    []byte  `cbor:"capability_token,omitempty" json:"capability_token,omitempty"`
	PackageID          string  `cbor:"package_id,omitempty" json:"package_id,omitempty"`
	WeightGrams        uint64  `cbor:"weight_grams,omitempty" json:"weight_grams,omitempty"`
	ShipTo             *ShipTo `cbor:"ship_to,omitempty" json:"ship_to,omitempty"`
	ServiceLevel       string  `cbor:"service_level,omitempty" json:"service_level,omitempty"`
	Status             string  `cbor:"status,omitempty" json:"status,omitempty"`
	CarrierName        string  `cbor:"carrier_name,omitempty" json:"carrier_name,omitempty"`
	TrackingNumber     string  `cbor:"tracking_number,omitempty" json:"tracking_number,omitempty"`
	LabelArtifactCID   []byte  `cbor:"label_artifact_cid,omitempty" json:"label_artifact_cid,omitempty"`
	Notes              string  `cbor:"notes,omitempty" json:"notes,omitempty"`
}

type KernelRegisterMessage struct {
	Role          string   `cbor:"role" json:"role"`
	ReceivePCIDs  [][]byte `cbor:"receive_pcids" json:"receive_pcids"`
}
