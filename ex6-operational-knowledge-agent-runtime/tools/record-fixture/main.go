package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/records"
)

// record-fixture emits deterministic base64 canonical Grid bytes for test
// agents. Intent: Keep shell fixtures on the production record codec rather
// than copying opaque CBOR literals. Built-in families use the fixed registry;
// external families must supply an explicit pCID. Source: DI-sidoh, DI-solan.
func main() {
	if len(os.Args) == 3 && os.Args[1] == "pcid" {
		protocolPCID := records.PackageProtocolPCID(os.Args[2])
		if protocolPCID == "" {
			fmt.Fprintln(os.Stderr, "external family requires an explicit frozen pCID")
			os.Exit(2)
		}
		fmt.Println(protocolPCID)
		return
	}
	if len(os.Args) == 2 && os.Args[1] == "inspect" {
		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		envelope, err := records.Parse(raw)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(envelope.Family)
		return
	}
	protocolPCID, family, recordID, signer, timestamp, payloadJSON := "", "", "", "", "", ""
	if len(os.Args) == 8 && os.Args[1] == "--pcid" {
		protocolPCID, family, recordID, signer, timestamp, payloadJSON = os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7]
	} else if len(os.Args) == 6 {
		family, recordID, signer, timestamp, payloadJSON = os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]
		protocolPCID = records.PackageProtocolPCID(family)
		if protocolPCID == "" {
			fmt.Fprintln(os.Stderr, "external family requires record-fixture --pcid <frozen-pcid> <family> <record-id> <signer> <timestamp> <payload-json>")
			os.Exit(2)
		}
	} else {
		fmt.Fprintln(os.Stderr, "usage: record-fixture [--pcid <frozen-pcid>] <family> <record-id> <signer> <timestamp> <payload-json>")
		os.Exit(2)
	}
	payload, err := records.CanonicalJSON([]byte(payloadJSON))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	envelope := records.Envelope{Family: family, ProtocolPCID: protocolPCID, RecordID: recordID, Signer: signer, Timestamp: timestamp, Payload: payload}
	if err := envelope.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(records.MustMarshal(envelope)))
}
