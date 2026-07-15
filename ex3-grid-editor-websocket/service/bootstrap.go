package service

import (
	"fmt"
	"time"
)

const mutationCapabilityTTL = 30 * time.Minute

type AppOptions struct {
	RemoteAccessToken string
}

type RemoteSession struct {
	DocumentID    string            `json:"document_id"`
	ParticipantID string            `json:"participant_id"`
	ExpiresAt     string            `json:"expires_at"`
	Capabilities  map[string]string `json:"capabilities"`
}

func (app *App) RemoteAccessEnabled() bool {
	return app.remoteAccessToken != ""
}

func (app *App) ValidateRemoteAccessToken(raw string) bool {
	return app.remoteAccessToken != "" && raw != "" && raw == app.remoteAccessToken
}

func (app *App) IssueRemoteSession(documentID string, participantID string) (RemoteSession, error) {
	if !app.RemoteAccessEnabled() {
		return RemoteSession{}, fmt.Errorf("remote mutation bootstrap is disabled")
	}
	if err := validateDocumentID(documentID); err != nil {
		return RemoteSession{}, err
	}
	if err := validateParticipantID(participantID); err != nil {
		return RemoteSession{}, err
	}
	// Intent: Keep remote bootstrap local and provisional while moving the live
	// mutation path itself onto short-lived per-protocol capabilities that match
	// ex3's existing protocol split. Source: DI-povip
	syncToken, claims, err := issueMutationCapability(app.identity, participantID, documentID, app.documentPCID.String(), "mutate", mutationCapabilityTTL)
	if err != nil {
		return RemoteSession{}, err
	}
	awarenessToken, _, err := issueMutationCapability(app.identity, participantID, documentID, app.awarenessPCID.String(), "mutate", mutationCapabilityTTL)
	if err != nil {
		return RemoteSession{}, err
	}
	metadataToken, _, err := issueMutationCapability(app.identity, participantID, documentID, app.metadataPCID.String(), "mutate", mutationCapabilityTTL)
	if err != nil {
		return RemoteSession{}, err
	}
	publishToken, _, err := issueMutationCapability(app.identity, participantID, documentID, app.publishPCID.String(), "mutate", mutationCapabilityTTL)
	if err != nil {
		return RemoteSession{}, err
	}
	return RemoteSession{
		DocumentID:    documentID,
		ParticipantID: participantID,
		ExpiresAt:     claims.ExpiresAt.UTC().Format(time.RFC3339),
		Capabilities: map[string]string{
			"sync":      syncToken,
			"awareness": awarenessToken,
			"metadata":  metadataToken,
			"publish":   publishToken,
		},
	}, nil
}
