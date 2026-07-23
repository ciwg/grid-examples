package store

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

type CASStore struct {
	root string
}

func NewCASStore(root string) *CASStore {
	return &CASStore{root: root}
}

func (cas *CASStore) ObjectPath(cid string) string {
	prefix := cid
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}
	return filepath.Join(cas.root, prefix, cid)
}

// Intent: Make portable blob and envelope bytes substrate-owned by CID so
// runtimes and relays can share one content-addressed write contract without
// inheriting ex5-local draft or attachment policy. Source: DI-lemor
func (cas *CASStore) WriteObject(data []byte) (string, error) {
	cid, err := protocols.CIDForBytes(data)
	if err != nil {
		return "", fmt.Errorf("cid cas object: %w", err)
	}
	target := cas.ObjectPath(cid.String())
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", err
	}
	existing, err := os.ReadFile(target)
	if err == nil {
		if !bytes.Equal(existing, data) {
			return "", fmt.Errorf("cas object %q already exists with different bytes", cid.String())
		}
		return cid.String(), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	tempFile, err := os.CreateTemp(filepath.Dir(target), "cas-*")
	if err != nil {
		return "", err
	}
	tempPath := tempFile.Name()
	if _, err := tempFile.Write(data); err != nil {
		return "", errors.Join(err, tempFile.Close(), os.Remove(tempPath))
	}
	if err := tempFile.Close(); err != nil {
		return "", errors.Join(err, os.Remove(tempPath))
	}
	if err := os.Rename(tempPath, target); err != nil {
		if errors.Is(err, os.ErrExist) {
			existing, readErr := os.ReadFile(target)
			if readErr != nil {
				return "", errors.Join(err, readErr)
			}
			if !bytes.Equal(existing, data) {
				return "", fmt.Errorf("cas object %q already exists with different bytes", cid.String())
			}
			if removeErr := os.Remove(tempPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				return "", removeErr
			}
			return cid.String(), nil
		}
		return "", errors.Join(err, os.Remove(tempPath))
	}
	return cid.String(), nil
}

// Intent: Keep CID verification substrate-owned so any runtime that reuses the
// store layer trusts the content-addressed bytes before higher ex5 policy runs.
// Source: DI-lemor
func (cas *CASStore) LoadObject(cid string) ([]byte, error) {
	body, err := os.ReadFile(cas.ObjectPath(cid))
	if err != nil {
		return nil, err
	}
	readCID, err := protocols.CIDForBytes(body)
	if err != nil {
		return nil, fmt.Errorf("cid cas object %q: %w", cid, err)
	}
	if readCID.String() != cid {
		return nil, fmt.Errorf("cas object %q bytes hash to %q", cid, readCID.String())
	}
	return body, nil
}

func (cas *CASStore) WriteEnvelopeBase64(expectedCID string, envelopeBase64 string) error {
	envelopeBytes, err := base64.StdEncoding.DecodeString(envelopeBase64)
	if err != nil {
		return fmt.Errorf("decode envelope base64: %w", err)
	}
	cid, err := cas.WriteObject(envelopeBytes)
	if err != nil {
		return err
	}
	if cid != expectedCID {
		return fmt.Errorf("envelope CAS cid mismatch: got %q want %q", cid, expectedCID)
	}
	return nil
}

// Intent: Make CAS authoritative for durable frozen-envelope bytes while
// keeping one-time manifest backfill inside the reusable persistence layer
// instead of scattering that migration rule across apps and relays. Source:
// DI-lemor
func (cas *CASStore) AuthoritativeEnvelopeBase64(envelopeCID string, manifestBase64 string) (string, error) {
	envelopeBytes, err := cas.LoadObject(envelopeCID)
	if err == nil {
		return base64.StdEncoding.EncodeToString(envelopeBytes), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if strings.TrimSpace(manifestBase64) == "" {
		return "", fmt.Errorf("cas envelope %q missing and no manifest fallback present", envelopeCID)
	}
	if err := cas.WriteEnvelopeBase64(envelopeCID, manifestBase64); err != nil {
		return "", err
	}
	envelopeBytes, err = cas.LoadObject(envelopeCID)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(envelopeBytes), nil
}
