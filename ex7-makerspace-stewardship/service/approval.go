package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const terminalApprovalPCID = "bafkreidztcisyvexrlia4eos7wko27e4rqt7ivbmikjbkr5tzbbts3rcd4"

// ApprovalRequest is bounded unsigned terminal data, never a Grid record.
// Intent: make the terminal wait for participant-agent approval rather than
// confusing a request or token with author evidence. Source: DI-hibok.
type ApprovalRequest struct {
	RequestID          string `json:"request_id"`
	TargetPCID         string `json:"target_pcid"`
	PayloadBase64      string `json:"payload_base64"`
	CreatedAt          string `json:"created_at"`
	ExpiresAt          string `json:"expires_at"`
	ApprovalToken      string `json:"approval_token"`
	State              string `json:"state"`
	SignedRecordBase64 string `json:"signed_record_base64,omitempty"`
}

func (a *App) CreateApprovalRequest(target string, payload []byte, now time.Time) (ApprovalRequest, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.store == nil || a.signer == nil || !isMakerspaceProtocol(target) {
		return ApprovalRequest{}, errors.New("participant signer and known target protocol are required")
	}
	canonical, err := json.Marshal(json.RawMessage(payload))
	if err != nil || !bytes.Equal(canonical, payload) {
		return ApprovalRequest{}, errors.New("approval payload must be canonical JSON")
	}
	requestID, err := randomApprovalValue(18)
	if err != nil {
		return ApprovalRequest{}, err
	}
	token, err := randomApprovalValue(32)
	if err != nil {
		return ApprovalRequest{}, err
	}
	request := ApprovalRequest{RequestID: requestID, TargetPCID: target, PayloadBase64: base64.StdEncoding.EncodeToString(payload), CreatedAt: now.UTC().Format(time.RFC3339), ExpiresAt: now.UTC().Add(10 * time.Minute).Format(time.RFC3339), ApprovalToken: token, State: "pending"}
	return request, a.writeApprovalRequest(request)
}

func (a *App) ApprovalRequest(id, token string, now time.Time) (ApprovalRequest, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	request, err := a.readApprovalRequest(id)
	if err != nil {
		return ApprovalRequest{}, err
	}
	if token == "" || token != request.ApprovalToken {
		return ApprovalRequest{}, errors.New("approval token is invalid")
	}
	if request.State == "pending" && requestExpired(request, now) {
		request.State = "expired"
		if err := a.writeApprovalRequest(request); err != nil {
			return ApprovalRequest{}, err
		}
	}
	return request, nil
}

func (a *App) ApproveRequest(id string, now time.Time) (ApprovalRequest, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.signer == nil {
		return ApprovalRequest{}, errors.New("participant signer is unavailable")
	}
	request, err := a.readApprovalRequest(id)
	if err != nil {
		return ApprovalRequest{}, err
	}
	if request.State != "pending" || requestExpired(request, now) {
		return ApprovalRequest{}, errors.New("approval request is not pending")
	}
	payload, err := base64.StdEncoding.DecodeString(request.PayloadBase64)
	if err != nil {
		return ApprovalRequest{}, errors.New("approval payload is invalid")
	}
	_, raw, err := a.signer.Sign(Record{Protocol: request.TargetPCID, ID: "terminal-" + request.RequestID, CreatedAt: now.UTC().Format(time.RFC3339), Payload: payload})
	if err != nil {
		return ApprovalRequest{}, err
	}
	if err := a.ingestLocked([][]byte{raw}); err != nil {
		return ApprovalRequest{}, err
	}
	request.State, request.SignedRecordBase64 = "approved", base64.StdEncoding.EncodeToString(raw)
	return request, a.writeApprovalRequest(request)
}

func (a *App) PendingApprovalRequests(now time.Time) ([]ApprovalRequest, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	directory := filepath.Join(a.store.root, "requests")
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	requests := make([]ApprovalRequest, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		request, err := a.readApprovalRequest(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			return nil, err
		}
		if request.State == "pending" && requestExpired(request, now) {
			request.State = "expired"
			if err := a.writeApprovalRequest(request); err != nil {
				return nil, err
			}
		}
		if request.State == "pending" {
			requests = append(requests, request)
		}
	}
	sort.Slice(requests, func(i, j int) bool { return requests[i].CreatedAt < requests[j].CreatedAt })
	return requests, nil
}

func (a *App) ingestLocked(records [][]byte) error {
	a.mu.Unlock()
	err := a.IngestRecords(records)
	a.mu.Lock()
	return err
}
func (a *App) requestPath(id string) string {
	return filepath.Join(a.store.root, "requests", id+".json")
}
func (a *App) writeApprovalRequest(request ApprovalRequest) error {
	if err := os.MkdirAll(filepath.Join(a.store.root, "requests"), 0o700); err != nil {
		return err
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}
	return os.WriteFile(a.requestPath(request.RequestID), raw, 0o600)
}
func (a *App) readApprovalRequest(id string) (ApprovalRequest, error) {
	raw, err := os.ReadFile(a.requestPath(id))
	if err != nil {
		return ApprovalRequest{}, err
	}
	var request ApprovalRequest
	if err := json.Unmarshal(raw, &request); err != nil || request.RequestID != id {
		return ApprovalRequest{}, fmt.Errorf("invalid approval request")
	}
	return request, nil
}
func requestExpired(request ApprovalRequest, now time.Time) bool {
	expiry, err := time.Parse(time.RFC3339, request.ExpiresAt)
	return err != nil || !now.Before(expiry)
}
func randomApprovalValue(size int) (string, error) {
	value := make([]byte, size)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}
