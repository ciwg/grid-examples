package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

type CAS struct {
	root string
}

func OpenCAS(root string) (*CAS, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &CAS{root: root}, nil
}

// Intent: Keep large package payload bodies content-addressed and immutable so
// append-only records can point at durable bytes without rewriting history.
// Source: DI-moksu
func (cas *CAS) Put(body []byte) (string, error) {
	sum := sha256.Sum256(body)
	objectID := "sha256:" + hex.EncodeToString(sum[:])
	path := cas.pathFor(objectID)
	if existing, err := os.ReadFile(path); err == nil {
		if sha256.Sum256(existing) == sum {
			return objectID, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := cas.writeObject(path, body); err != nil {
		return "", err
	}
	return objectID, nil
}

func (cas *CAS) Get(objectID string) ([]byte, error) {
	return os.ReadFile(cas.pathFor(objectID))
}

// PutCID stores content under a CIDv1 raw-sha2-256 identifier.
// Intent: Allow new PromiseGrid-native records to retain binary-CID identity
// without invalidating existing sha256:<hex> references. Source: DI-bavuk
func (cas *CAS) PutCID(body []byte) (cid.Cid, error) {
	objectCID, err := cas.cidFor(body)
	if err != nil {
		return cid.Undef, err
	}
	path := cas.pathFor(objectCID.String())
	corruptCIDPath := false
	if existing, err := os.ReadFile(path); err == nil {
		actual, hashErr := cas.cidFor(existing)
		if hashErr != nil {
			return cid.Undef, hashErr
		}
		if actual == objectCID {
			return objectCID, nil
		}
		corruptCIDPath = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return cid.Undef, err
	}
	legacyPath := cas.pathFor(legacyObjectID(body))
	if existing, err := os.ReadFile(legacyPath); err == nil {
		actual, hashErr := cas.cidFor(existing)
		if hashErr != nil {
			return cid.Undef, hashErr
		}
		if actual == objectCID && !corruptCIDPath {
			return objectCID, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return cid.Undef, err
	}
	// Intent: Replace a known-bad object atomically so a retry repairs an
	// interrupted write instead of preserving a corrupt filename forever.
	// Source: DI-bavuk
	if err := cas.writeObject(path, body); err != nil {
		return cid.Undef, err
	}
	return objectCID, nil
}

// GetCID reads a CIDv1 object and falls back to its legacy filename when needed.
// Intent: Historical sha256:<hex> files remain readable while CID-addressed
// event replay becomes the authoritative path. Source: DI-bavuk
func (cas *CAS) GetCID(objectCID cid.Cid) ([]byte, error) {
	if err := validateCASCID(objectCID); err != nil {
		return nil, err
	}
	body, err := os.ReadFile(cas.pathFor(objectCID.String()))
	if errors.Is(err, os.ErrNotExist) {
		legacyID, legacyErr := legacyIDForCID(objectCID)
		if legacyErr != nil {
			return nil, legacyErr
		}
		body, err = os.ReadFile(cas.pathFor(legacyID))
	}
	if err != nil {
		return nil, err
	}
	actualCID, err := cas.cidFor(body)
	if err != nil {
		return nil, err
	}
	if actualCID != objectCID {
		return nil, errors.New("CAS object bytes do not match requested CID")
	}
	return body, nil
}

// ListCIDs returns the normalized CIDv1 view of all readable local CAS objects.
// Intent: A projection can rebuild from retained content-addressed event bytes
// instead of trusting a mutable index as its lifecycle source. Source: DI-bavuk
func (cas *CAS) ListCIDs() ([]cid.Cid, error) {
	entries, err := os.ReadDir(cas.root)
	if err != nil {
		return nil, err
	}
	objectCIDs := make([]cid.Cid, 0, len(entries))
	seen := map[string]struct{}{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		objectCID, err := cas.cidFromFilename(entry.Name())
		if err != nil {
			// Intent: CAS replay must retain readable objects even if an interrupted
			// write or unrelated file appears beside them. Source: DI-bavuk
			continue
		}
		if _, exists := seen[objectCID.String()]; exists {
			continue
		}
		seen[objectCID.String()] = struct{}{}
		objectCIDs = append(objectCIDs, objectCID)
	}
	slices.SortFunc(objectCIDs, func(left, right cid.Cid) int {
		return strings.Compare(left.String(), right.String())
	})
	return objectCIDs, nil
}

func (cas *CAS) writeObject(path string, body []byte) (result error) {
	temporary, err := os.CreateTemp(cas.root, ".cas-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		if temporaryPath == "" {
			return
		}
		if err := os.Remove(temporaryPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			result = errors.Join(result, err)
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if _, err := temporary.Write(body); err != nil {
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if err := temporary.Sync(); err != nil {
		if closeErr := temporary.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	temporaryPath = ""
	directory, err := os.Open(cas.root)
	if err != nil {
		return err
	}
	if err := directory.Sync(); err != nil {
		if closeErr := directory.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	return directory.Close()
}

func (cas *CAS) cidFor(body []byte) (cid.Cid, error) {
	digest, err := mh.Sum(body, mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, err
	}
	return cid.NewCidV1(cid.Raw, digest), nil
}

func (cas *CAS) cidFromFilename(name string) (cid.Cid, error) {
	if objectCID, err := cid.Decode(name); err == nil {
		if err := validateCASCID(objectCID); err != nil {
			return cid.Undef, err
		}
		return objectCID, nil
	}
	if !strings.HasPrefix(name, "sha256:") {
		return cid.Undef, fmt.Errorf("CAS object filename %q is not a CIDv1 or legacy sha256 identifier", name)
	}
	digest, err := hex.DecodeString(strings.TrimPrefix(name, "sha256:"))
	if err != nil {
		return cid.Undef, fmt.Errorf("legacy CAS identifier %q: %w", name, err)
	}
	if len(digest) != sha256.Size {
		return cid.Undef, fmt.Errorf("legacy CAS identifier %q has an invalid digest length", name)
	}
	multihash, err := mh.Encode(digest, mh.SHA2_256)
	if err != nil {
		return cid.Undef, err
	}
	return cid.NewCidV1(cid.Raw, multihash), nil
}

func validateCASCID(objectCID cid.Cid) error {
	if !objectCID.Defined() || objectCID.Version() != 1 || objectCID.Type() != cid.Raw {
		return errors.New("CAS object must use a CIDv1 raw codec")
	}
	decoded, err := mh.Decode(objectCID.Hash())
	if err != nil {
		return err
	}
	if decoded.Code != mh.SHA2_256 || len(decoded.Digest) != sha256.Size {
		return errors.New("CAS object must use sha2-256")
	}
	return nil
}

func legacyObjectID(body []byte) string {
	sum := sha256.Sum256(body)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// LegacyObjectID returns the deterministic compatibility identifier without
// persisting bytes. It lets callers validate a complete write proposal before
// beginning durable mutation. Source: DI-fofuh
func LegacyObjectID(body []byte) string {
	return legacyObjectID(body)
}

func legacyIDForCID(objectCID cid.Cid) (string, error) {
	if err := validateCASCID(objectCID); err != nil {
		return "", err
	}
	decoded, err := mh.Decode(objectCID.Hash())
	if err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(decoded.Digest), nil
}

func (cas *CAS) pathFor(objectID string) string {
	return filepath.Join(cas.root, objectID)
}
