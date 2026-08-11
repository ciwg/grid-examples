package service

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
)

const gridTag uint64 = 1735551332

// Record is Ex7's exact durable promise carriage. Intent: bind every local
// projection to independently verifiable participant evidence. Source: DI-sinov.
type Record struct {
	Protocol, ID, Signer, CreatedAt, KeyID string
	Payload                                json.RawMessage
	PublicKey, Signature                   []byte
}

func (r Record) signingBytes() ([]byte, error) { return r.marshal(nil) }

func (r Record) marshal(signature []byte) ([]byte, error) {
	p, err := cid.Decode(r.Protocol)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Signer) == "" || !json.Valid(r.Payload) {
		return nil, errors.New("invalid record fields")
	}
	canonical, err := json.Marshal(json.RawMessage(r.Payload))
	if err != nil || !bytes.Equal(canonical, r.Payload) {
		return nil, errors.New("payload is not canonical JSON")
	}
	if len(r.PublicKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid author public key")
	}
	if r.KeyID != keyID(r.PublicKey) {
		return nil, errors.New("author key ID does not match public key")
	}
	encoder, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, err
	}
	return encoder.Marshal(cbor.Tag{Number: gridTag, Content: []any{cbor.Tag{Number: 42, Content: p.Bytes()}, r.ID, r.Signer, r.CreatedAt, []byte(r.Payload), r.KeyID, r.PublicKey, signature}})
}

func keyID(publicKey ed25519.PublicKey) string {
	digest := sha256.Sum256(publicKey)
	return "ed25519:" + hex.EncodeToString(digest[:])
}

func (r Record) Sign(private ed25519.PrivateKey) (Record, []byte, error) {
	raw, err := r.signingBytes()
	if err != nil {
		return Record{}, nil, err
	}
	r.Signature = ed25519.Sign(private, raw)
	out, err := r.marshal(r.Signature)
	return r, out, err
}

func ParseRecord(raw []byte) (Record, error) {
	var tag cbor.Tag
	if err := cbor.Unmarshal(raw, &tag); err != nil {
		return Record{}, err
	}
	if tag.Number != gridTag {
		return Record{}, errors.New("not a Grid record")
	}
	values, ok := tag.Content.([]any)
	if !ok || len(values) != 8 {
		return Record{}, errors.New("invalid record slots")
	}
	selector, ok := values[0].(cbor.Tag)
	if !ok {
		return Record{}, errors.New("missing protocol selector")
	}
	p, err := cid.Cast(selector.Content.([]byte))
	if err != nil {
		return Record{}, err
	}
	r := Record{Protocol: p.String()}
	fields := []*string{&r.ID, &r.Signer, &r.CreatedAt}
	for i, f := range fields {
		v, ok := values[i+1].(string)
		if !ok {
			return Record{}, fmt.Errorf("record slot %d must be text", i+1)
		}
		*f = v
	}
	payload, ok := values[4].([]byte)
	if !ok {
		return Record{}, errors.New("payload must be bytes")
	}
	r.Payload = payload
	r.KeyID, _ = values[5].(string)
	r.PublicKey, _ = values[6].([]byte)
	r.Signature, _ = values[7].([]byte)
	signing, err := r.signingBytes()
	if err != nil || !ed25519.Verify(r.PublicKey, signing, r.Signature) {
		return Record{}, errors.New("invalid author signature")
	}
	canonical, err := r.marshal(r.Signature)
	if err != nil || !bytes.Equal(raw, canonical) {
		return Record{}, errors.New("record is not canonical")
	}
	return r, nil
}
