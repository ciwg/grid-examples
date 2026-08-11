package protocol

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const (
	gridTag = uint64(0x67726964)
	cidTag  = uint64(42)
)

var (
	canonicalEncMode = mustCanonicalEncMode()
	strictDecMode    = mustStrictDecMode()
)

// Proof is the signer-owned proof in envelope slot 2.
type Proof struct {
	Algorithm string `cbor:"alg"`
	AgentID   string `cbor:"agent_id"`
	PublicKey []byte `cbor:"pub"`
	Signature []byte `cbor:"sig"`
}

// Envelope is Ex4's local-draft grid([42(pCID), payload, proof]) artifact.
type Envelope struct {
	PCID         cid.Cid
	PayloadBytes []byte
	Proof        Proof
}

// PromiseDraft carries a binary tag-42 pCID selector and raw profile payload
// for the local signer bridge.
type PromiseDraft struct {
	PCID         cid.Cid
	PayloadBytes []byte
}

// PromiseProof carries a browser-owned proof over a prepared draft.
type PromiseProof struct {
	Draft PromiseDraft
	Proof Proof
}

func (draft PromiseDraft) SignableBytes() ([]byte, error) {
	return Marshal([]any{rawPCIDTag(draft.PCID), cbor.RawMessage(draft.PayloadBytes)})
}

// IssueReport is the pCID-owned payload for a new issue report.
type IssueReport struct {
	AgentID     string `cbor:"agent_id" json:"agent_id"`
	IssuedAt    string `cbor:"issued_at" json:"issued_at"`
	Team        string `cbor:"team" json:"team"`
	Title       string `cbor:"title" json:"title"`
	Description string `cbor:"description" json:"description"`
	Severity    string `cbor:"severity" json:"severity"`
}

// IssueLifecycleUpdate is the pCID-owned payload for an issue comment,
// assignment, or status update.
type IssueLifecycleUpdate struct {
	AgentID         string `cbor:"agent_id" json:"agent_id"`
	IssuedAt        string `cbor:"issued_at" json:"issued_at"`
	IssueID         string `cbor:"issue_id" json:"issue_id"`
	Kind            string `cbor:"kind" json:"kind"`
	Comment         string `cbor:"comment,omitempty" json:"comment,omitempty"`
	AssigneeAgentID string `cbor:"assignee_agent_id,omitempty" json:"assignee_agent_id,omitempty"`
	Status          string `cbor:"status,omitempty" json:"status,omitempty"`
}

// IssueAttachmentReference is the pCID-owned payload for an attachment object.
type IssueAttachmentReference struct {
	AgentID       string `cbor:"agent_id" json:"agent_id"`
	IssuedAt      string `cbor:"issued_at" json:"issued_at"`
	IssueID       string `cbor:"issue_id" json:"issue_id"`
	AttachmentCID string `cbor:"attachment_cid" json:"attachment_cid"`
	Name          string `cbor:"name" json:"name"`
	ContentType   string `cbor:"content_type" json:"content_type"`
	Size          int64  `cbor:"size" json:"size"`
}

func mustCanonicalEncMode() cbor.EncMode {
	mode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		panic(fmt.Errorf("canonical enc mode: %w", err))
	}
	return mode
}

func mustStrictDecMode() cbor.DecMode {
	mode, err := cbor.DecOptions{DupMapKey: cbor.DupMapKeyEnforcedAPF}.DecMode()
	if err != nil {
		panic(fmt.Errorf("strict dec mode: %w", err))
	}
	return mode
}

func Marshal(value any) ([]byte, error) {
	return canonicalEncMode.Marshal(value)
}

func Unmarshal(data []byte, value any) error {
	return strictDecMode.Unmarshal(data, value)
}

func CIDForBytes(data []byte) (cid.Cid, error) {
	sum, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, fmt.Errorf("multihash bytes: %w", err)
	}
	return cid.NewCidV1(cid.Raw, sum), nil
}

func NewEnvelope(pcid cid.Cid, payloadBytes []byte, proof Proof) Envelope {
	return Envelope{PCID: pcid, PayloadBytes: append([]byte(nil), payloadBytes...), Proof: proof}
}

// Intent: Keep the profile-selected payload as raw slot 1 bytes so the signed
// local-draft wire shape remains grid([42(pCID), payload, proof]). Source: DI-gonok
func (envelope Envelope) SignableBytes() ([]byte, error) {
	return Marshal([]any{rawPCIDTag(envelope.PCID), cbor.RawMessage(envelope.PayloadBytes)})
}

func (envelope Envelope) Bytes() ([]byte, error) {
	outer := cbor.RawTag{Number: gridTag, Content: mustMarshal([]any{
		rawPCIDTag(envelope.PCID), cbor.RawMessage(envelope.PayloadBytes), envelope.Proof,
	})}
	return Marshal(outer)
}

func ParseEnvelope(envelopeBytes []byte) (Envelope, error) {
	var outer cbor.RawTag
	if err := Unmarshal(envelopeBytes, &outer); err != nil {
		return Envelope{}, fmt.Errorf("decode outer envelope: %w", err)
	}
	if outer.Number != gridTag {
		return Envelope{}, fmt.Errorf("unexpected outer tag %d", outer.Number)
	}
	var slots []cbor.RawMessage
	if err := Unmarshal(outer.Content, &slots); err != nil {
		return Envelope{}, fmt.Errorf("decode envelope slots: %w", err)
	}
	if len(slots) != 3 {
		return Envelope{}, fmt.Errorf("unexpected envelope slot count %d", len(slots))
	}
	var pcidTag cbor.RawTag
	if err := Unmarshal(slots[0], &pcidTag); err != nil {
		return Envelope{}, fmt.Errorf("decode pCID tag: %w", err)
	}
	if pcidTag.Number != cidTag {
		return Envelope{}, fmt.Errorf("unexpected pCID tag %d", pcidTag.Number)
	}
	var pcidBytes []byte
	if err := Unmarshal(pcidTag.Content, &pcidBytes); err != nil {
		return Envelope{}, fmt.Errorf("decode pCID bytes: %w", err)
	}
	pcid, err := cid.Cast(pcidBytes)
	if err != nil {
		return Envelope{}, fmt.Errorf("cast pCID bytes: %w", err)
	}
	var proof Proof
	if err := Unmarshal(slots[2], &proof); err != nil {
		return Envelope{}, fmt.Errorf("decode proof: %w", err)
	}
	return NewEnvelope(pcid, slots[1], proof), nil
}

func mustMarshal(value any) []byte {
	bytes, err := Marshal(value)
	if err != nil {
		panic(err)
	}
	return bytes
}

func rawPCIDTag(pcid cid.Cid) cbor.RawTag {
	return cbor.RawTag{Number: cidTag, Content: cbor.RawMessage(mustMarshal(pcid.Bytes()))}
}
