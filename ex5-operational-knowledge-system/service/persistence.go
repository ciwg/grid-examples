package service

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/protocols"
)

const maxEventLineBytes = 1 << 20

type Store struct {
	root                            string
	events                          *os.File
	knowledgeItemMessages           *os.File
	knowledgeApprovalMessages       *os.File
	knowledgeEvidenceMessages       *os.File
	operationalRunMessages          *os.File
	operationalPlaceMessages        *os.File
	operationalResourceMessages     *os.File
	knowledgeLinkMessages           *os.File
	knowledgeResponsibilityMessages *os.File
	eventPath                       string
	knowledgeItemPath               string
	knowledgeApprovalPath           string
	knowledgeEvidencePath           string
	operationalRunPath              string
	operationalPlacePath            string
	operationalResourcePath         string
	knowledgeLinkPath               string
	knowledgeResponsibilityPath     string
	draftPath                       string
	casRoot                         string
	identity                        *RuntimeIdentity
}

type PersistedDraft struct {
	Body      string `json:"body"`
	Version   int    `json:"version"`
	UpdatedAt string `json:"updated_at"`
}

// Intent: Keep durable operational truth in an ex5-local append-only event log
// plus the staged signed-family logs and copied attachments so the example can
// preserve history independently of any browser or CLI session state. Source:
// DI-radok; DI-zuvob; DI-mibor; DI-vosul; DI-kavup; DI-votek; DI-sarib;
// DI-vamok; DI-pivul
func OpenStore(root string) (*Store, []OperationalEvent, []SignedKnowledgeItemRecord, []SignedKnowledgeApprovalRecord, []SignedKnowledgeEvidenceRecord, []SignedOperationalRunRecord, []SignedOperationalPlaceRecord, []SignedOperationalResourceRecord, []SignedKnowledgeLinkRecord, []SignedKnowledgeResponsibilityRecord, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "attachments"), 0o755); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "cas", "objects"), 0o755); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "drafts"), 0o755); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	identity, err := LoadOrCreateRuntimeIdentity(filepath.Join(root, "identity", "knowledge-item-ed25519.seed"))
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	eventPath := filepath.Join(root, "events.jsonl")
	eventsFile, err := os.OpenFile(eventPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	events, err := readEvents(eventsFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close())
	}
	knowledgeItemPath := filepath.Join(root, "knowledge-item-messages.jsonl")
	knowledgeItemFile, err := os.OpenFile(knowledgeItemPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close())
	}
	knowledgeItemRecords, err := readSignedKnowledgeItemRecords(knowledgeItemFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	knowledgeApprovalPath := filepath.Join(root, "knowledge-approval-messages.jsonl")
	knowledgeApprovalFile, err := os.OpenFile(knowledgeApprovalPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	knowledgeApprovalRecords, err := readSignedKnowledgeApprovalRecords(knowledgeApprovalFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close())
	}
	knowledgeEvidencePath := filepath.Join(root, "knowledge-evidence-messages.jsonl")
	knowledgeEvidenceFile, err := os.OpenFile(knowledgeEvidencePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close())
	}
	knowledgeEvidenceRecords, err := readSignedKnowledgeEvidenceRecords(knowledgeEvidenceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close())
	}
	operationalRunPath := filepath.Join(root, "operational-run-messages.jsonl")
	operationalRunFile, err := os.OpenFile(operationalRunPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close())
	}
	operationalRunRecords, err := readSignedOperationalRunRecords(operationalRunFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close())
	}
	operationalPlacePath := filepath.Join(root, "operational-place-messages.jsonl")
	operationalPlaceFile, err := os.OpenFile(operationalPlacePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close())
	}
	operationalPlaceRecords, err := readSignedOperationalPlaceRecords(operationalPlaceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close())
	}
	operationalResourcePath := filepath.Join(root, "operational-resource-messages.jsonl")
	operationalResourceFile, err := os.OpenFile(operationalResourcePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close())
	}
	operationalResourceRecords, err := readSignedOperationalResourceRecords(operationalResourceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close())
	}
	knowledgeLinkPath := filepath.Join(root, "knowledge-link-messages.jsonl")
	knowledgeLinkFile, err := os.OpenFile(knowledgeLinkPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close())
	}
	knowledgeLinkRecords, err := readSignedKnowledgeLinkRecords(knowledgeLinkFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close())
	}
	knowledgeResponsibilityPath := filepath.Join(root, "knowledge-responsibility-messages.jsonl")
	knowledgeResponsibilityFile, err := os.OpenFile(knowledgeResponsibilityPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close())
	}
	knowledgeResponsibilityRecords, err := readSignedKnowledgeResponsibilityRecords(knowledgeResponsibilityFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := eventsFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := knowledgeItemFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := knowledgeApprovalFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := knowledgeEvidenceFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := operationalRunFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := operationalPlaceFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := operationalResourceFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := knowledgeLinkFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	if _, err := knowledgeResponsibilityFile.Seek(0, os.SEEK_END); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close(), knowledgeResponsibilityFile.Close())
	}
	store := &Store{
		root:                            root,
		events:                          eventsFile,
		knowledgeItemMessages:           knowledgeItemFile,
		knowledgeApprovalMessages:       knowledgeApprovalFile,
		knowledgeEvidenceMessages:       knowledgeEvidenceFile,
		operationalRunMessages:          operationalRunFile,
		operationalPlaceMessages:        operationalPlaceFile,
		operationalResourceMessages:     operationalResourceFile,
		knowledgeLinkMessages:           knowledgeLinkFile,
		knowledgeResponsibilityMessages: knowledgeResponsibilityFile,
		eventPath:                       eventPath,
		knowledgeItemPath:               knowledgeItemPath,
		knowledgeApprovalPath:           knowledgeApprovalPath,
		knowledgeEvidencePath:           knowledgeEvidencePath,
		operationalRunPath:              operationalRunPath,
		operationalPlacePath:            operationalPlacePath,
		operationalResourcePath:         operationalResourcePath,
		knowledgeLinkPath:               knowledgeLinkPath,
		knowledgeResponsibilityPath:     knowledgeResponsibilityPath,
		draftPath:                       filepath.Join(root, "drafts"),
		casRoot:                         filepath.Join(root, "cas", "objects"),
		identity:                        identity,
	}
	knowledgeItemRecords, err = store.hydrateSignedKnowledgeItemRecords(knowledgeItemRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeApprovalRecords, err = store.hydrateSignedKnowledgeApprovalRecords(knowledgeApprovalRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeEvidenceRecords, err = store.hydrateSignedKnowledgeEvidenceRecords(knowledgeEvidenceRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalRunRecords, err = store.hydrateSignedOperationalRunRecords(operationalRunRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalPlaceRecords, err = store.hydrateSignedOperationalPlaceRecords(operationalPlaceRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	operationalResourceRecords, err = store.hydrateSignedOperationalResourceRecords(operationalResourceRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeLinkRecords, err = store.hydrateSignedKnowledgeLinkRecords(knowledgeLinkRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	knowledgeResponsibilityRecords, err = store.hydrateSignedKnowledgeResponsibilityRecords(knowledgeResponsibilityRecords)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	return store, events, knowledgeItemRecords, knowledgeApprovalRecords, knowledgeEvidenceRecords, operationalRunRecords, operationalPlaceRecords, operationalResourceRecords, knowledgeLinkRecords, knowledgeResponsibilityRecords, nil
}

func readEvents(file *os.File) (events []OperationalEvent, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	// Intent: Replay the full stored event log within the server's current
	// request-size envelope so durable large revisions do not become
	// unreadable after restart. Source: DI-busor
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	events = []OperationalEvent{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var event OperationalEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, fmt.Errorf("decode event: %w", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (store *Store) AppendEvent(event OperationalEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err := store.events.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.events.Sync()
}

func readSignedKnowledgeItemRecords(file *os.File) (records []SignedKnowledgeItemRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedKnowledgeItemRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedKnowledgeItemRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode knowledge-item record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeItemRecord(record SignedKnowledgeItemRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.knowledgeItemMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.knowledgeItemMessages.Sync()
}

func (store *Store) LoadSignedKnowledgeItemRecordsAuthoritative() ([]SignedKnowledgeItemRecord, error) {
	records, err := readSignedKnowledgeItemRecords(store.knowledgeItemMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeItemRecords(records)
}

func readSignedKnowledgeApprovalRecords(file *os.File) (records []SignedKnowledgeApprovalRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedKnowledgeApprovalRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedKnowledgeApprovalRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode knowledge-approval record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeApprovalRecord(record SignedKnowledgeApprovalRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.knowledgeApprovalMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.knowledgeApprovalMessages.Sync()
}

func (store *Store) LoadSignedKnowledgeApprovalRecordsAuthoritative() ([]SignedKnowledgeApprovalRecord, error) {
	records, err := readSignedKnowledgeApprovalRecords(store.knowledgeApprovalMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeApprovalRecords(records)
}

func readSignedKnowledgeEvidenceRecords(file *os.File) (records []SignedKnowledgeEvidenceRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedKnowledgeEvidenceRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedKnowledgeEvidenceRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode knowledge-evidence record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeEvidenceRecord(record SignedKnowledgeEvidenceRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.knowledgeEvidenceMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.knowledgeEvidenceMessages.Sync()
}

func (store *Store) LoadSignedKnowledgeEvidenceRecordsAuthoritative() ([]SignedKnowledgeEvidenceRecord, error) {
	records, err := readSignedKnowledgeEvidenceRecords(store.knowledgeEvidenceMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeEvidenceRecords(records)
}

func readSignedOperationalRunRecords(file *os.File) (records []SignedOperationalRunRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedOperationalRunRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedOperationalRunRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode operational-run record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedOperationalRunRecord(record SignedOperationalRunRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.operationalRunMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.operationalRunMessages.Sync()
}

func (store *Store) LoadSignedOperationalRunRecordsAuthoritative() ([]SignedOperationalRunRecord, error) {
	records, err := readSignedOperationalRunRecords(store.operationalRunMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedOperationalRunRecords(records)
}

func readSignedOperationalPlaceRecords(file *os.File) (records []SignedOperationalPlaceRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedOperationalPlaceRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedOperationalPlaceRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode operational-place record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedOperationalPlaceRecord(record SignedOperationalPlaceRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.operationalPlaceMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.operationalPlaceMessages.Sync()
}

func (store *Store) LoadSignedOperationalPlaceRecordsAuthoritative() ([]SignedOperationalPlaceRecord, error) {
	records, err := readSignedOperationalPlaceRecords(store.operationalPlaceMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedOperationalPlaceRecords(records)
}

func readSignedOperationalResourceRecords(file *os.File) (records []SignedOperationalResourceRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedOperationalResourceRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedOperationalResourceRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode operational-resource record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedOperationalResourceRecord(record SignedOperationalResourceRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.operationalResourceMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.operationalResourceMessages.Sync()
}

func (store *Store) LoadSignedOperationalResourceRecordsAuthoritative() ([]SignedOperationalResourceRecord, error) {
	records, err := readSignedOperationalResourceRecords(store.operationalResourceMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedOperationalResourceRecords(records)
}

func readSignedKnowledgeLinkRecords(file *os.File) (records []SignedKnowledgeLinkRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedKnowledgeLinkRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedKnowledgeLinkRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode knowledge-link record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeLinkRecord(record SignedKnowledgeLinkRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.knowledgeLinkMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.knowledgeLinkMessages.Sync()
}

func (store *Store) LoadSignedKnowledgeLinkRecordsAuthoritative() ([]SignedKnowledgeLinkRecord, error) {
	records, err := readSignedKnowledgeLinkRecords(store.knowledgeLinkMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeLinkRecords(records)
}

func readSignedKnowledgeResponsibilityRecords(file *os.File) (records []SignedKnowledgeResponsibilityRecord, err error) {
	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	defer func() {
		if _, seekErr := file.Seek(0, os.SEEK_END); seekErr != nil {
			err = errors.Join(err, seekErr)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxEventLineBytes)
	records = []SignedKnowledgeResponsibilityRecord{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record SignedKnowledgeResponsibilityRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode knowledge-responsibility record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeResponsibilityRecord(record SignedKnowledgeResponsibilityRecord) error {
	if err := store.writeEnvelopeCAS(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := store.knowledgeResponsibilityMessages.Write(append(body, '\n')); err != nil {
		return err
	}
	return store.knowledgeResponsibilityMessages.Sync()
}

func (store *Store) LoadSignedKnowledgeResponsibilityRecordsAuthoritative() ([]SignedKnowledgeResponsibilityRecord, error) {
	records, err := readSignedKnowledgeResponsibilityRecords(store.knowledgeResponsibilityMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeResponsibilityRecords(records)
}

// Intent: Preserve evidence attachment history by storing each uploaded file at
// a unique immutable path instead of overwriting earlier evidence bytes when a
// later upload reuses the same filename, while also dual-writing the same
// bytes into the staged CAS sidecar for later portable blob exchange. Source:
// DI-busor; DI-ribek
func (store *Store) SaveAttachment(entityID string, filename string, data []byte) (string, string, int64, error) {
	dir := filepath.Join(store.root, "attachments", entityID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", 0, err
	}
	base := filepath.Base(strings.TrimSpace(filename))
	if base == "" || base == "." {
		base = "attachment.bin"
	}
	tempFile, err := os.CreateTemp(dir, "evidence-*-"+base)
	if err != nil {
		return "", "", 0, err
	}
	target := tempFile.Name()
	if _, err := tempFile.Write(data); err != nil {
		return "", "", 0, errors.Join(err, tempFile.Close())
	}
	if err := tempFile.Close(); err != nil {
		return "", "", 0, err
	}
	cid, err := store.writeCASObject(data)
	if err != nil {
		return "", "", 0, err
	}
	return target, cid, int64(len(data)), nil
}

// Intent: Keep imported evidence attachments locally usable after peer
// exchange by materializing CID-addressed blob bytes into the compatibility
// attachment tree when the original source-host path is not valid here.
// Source: DI-faruv
func (store *Store) MaterializeAttachmentFromCID(entityID string, attachmentName string, attachmentCID string) (string, error) {
	if strings.TrimSpace(attachmentCID) == "" {
		return "", nil
	}
	body, err := store.loadCASObject(attachmentCID)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(store.root, "attachments", entityID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	base := filepath.Base(strings.TrimSpace(attachmentName))
	if base == "" || base == "." {
		base = "attachment.bin"
	}
	target := filepath.Join(dir, "cid-"+attachmentCID+"-"+base)
	existing, err := os.ReadFile(target)
	if err == nil {
		if !bytes.Equal(existing, body) {
			return "", fmt.Errorf("materialized attachment %q already exists with different bytes", target)
		}
		return target, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.WriteFile(target, body, 0o644); err != nil {
		return "", err
	}
	return target, nil
}

func (store *Store) writeEnvelopeCAS(expectedCID string, envelopeBase64 string) error {
	envelopeBytes, err := base64.StdEncoding.DecodeString(envelopeBase64)
	if err != nil {
		return fmt.Errorf("decode envelope base64: %w", err)
	}
	cid, err := store.writeCASObject(envelopeBytes)
	if err != nil {
		return err
	}
	if cid != expectedCID {
		return fmt.Errorf("envelope CAS cid mismatch: got %q want %q", cid, expectedCID)
	}
	return nil
}

// Intent: Make CAS authoritative for the eight frozen family envelope bytes
// while allowing one-time backfill from the manifest copy so older runtimes can
// migrate into the stronger replay path without a destructive cutover. Source:
// DI-rovud
func (store *Store) authoritativeEnvelopeBase64(envelopeCID string, manifestBase64 string) (string, error) {
	envelopeBytes, err := store.loadCASObject(envelopeCID)
	if err == nil {
		return base64.StdEncoding.EncodeToString(envelopeBytes), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if strings.TrimSpace(manifestBase64) == "" {
		return "", fmt.Errorf("cas envelope %q missing and no manifest fallback present", envelopeCID)
	}
	if err := store.writeEnvelopeCAS(envelopeCID, manifestBase64); err != nil {
		return "", err
	}
	envelopeBytes, err = store.loadCASObject(envelopeCID)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(envelopeBytes), nil
}

func (store *Store) writeCASObject(data []byte) (string, error) {
	cid, err := protocols.CIDForBytes(data)
	if err != nil {
		return "", fmt.Errorf("cid cas object: %w", err)
	}
	target := store.casObjectPath(cid.String())
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

func (store *Store) loadCASObject(cid string) ([]byte, error) {
	body, err := os.ReadFile(store.casObjectPath(cid))
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

func (store *Store) casObjectPath(cid string) string {
	prefix := cid
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}
	return filepath.Join(store.casRoot, prefix, cid)
}

func (store *Store) hydrateSignedKnowledgeItemRecords(records []SignedKnowledgeItemRecord) ([]SignedKnowledgeItemRecord, error) {
	out := append([]SignedKnowledgeItemRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative knowledge-item envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedKnowledgeApprovalRecords(records []SignedKnowledgeApprovalRecord) ([]SignedKnowledgeApprovalRecord, error) {
	out := append([]SignedKnowledgeApprovalRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative knowledge-approval envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedKnowledgeEvidenceRecords(records []SignedKnowledgeEvidenceRecord) ([]SignedKnowledgeEvidenceRecord, error) {
	out := append([]SignedKnowledgeEvidenceRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative knowledge-evidence envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedOperationalRunRecords(records []SignedOperationalRunRecord) ([]SignedOperationalRunRecord, error) {
	out := append([]SignedOperationalRunRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative operational-run envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedOperationalPlaceRecords(records []SignedOperationalPlaceRecord) ([]SignedOperationalPlaceRecord, error) {
	out := append([]SignedOperationalPlaceRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative operational-place envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedOperationalResourceRecords(records []SignedOperationalResourceRecord) ([]SignedOperationalResourceRecord, error) {
	out := append([]SignedOperationalResourceRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative operational-resource envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedKnowledgeLinkRecords(records []SignedKnowledgeLinkRecord) ([]SignedKnowledgeLinkRecord, error) {
	out := append([]SignedKnowledgeLinkRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative knowledge-link envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

func (store *Store) hydrateSignedKnowledgeResponsibilityRecords(records []SignedKnowledgeResponsibilityRecord) ([]SignedKnowledgeResponsibilityRecord, error) {
	out := append([]SignedKnowledgeResponsibilityRecord(nil), records...)
	for i := range out {
		base64Envelope, err := store.authoritativeEnvelopeBase64(out[i].EnvelopeCID, out[i].EnvelopeBase64)
		if err != nil {
			return nil, fmt.Errorf("load authoritative knowledge-responsibility envelope %d: %w", out[i].Sequence, err)
		}
		out[i].EnvelopeBase64 = base64Envelope
	}
	return out, nil
}

// Intent: Restore the shared browser working bodies on startup without mixing
// them into the append-only event log, so durable revision history stays
// explicit while collaborative drafting can resume after a restart. Source:
// DI-lusov; DI-zoruk
func (store *Store) LoadDrafts() (map[string]PersistedDraft, error) {
	entries, err := os.ReadDir(store.draftPath)
	if err != nil {
		return nil, err
	}
	out := map[string]PersistedDraft{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		body, err := os.ReadFile(filepath.Join(store.draftPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		var draft PersistedDraft
		if err := json.Unmarshal(body, &draft); err != nil {
			return nil, fmt.Errorf("decode draft %s: %w", entry.Name(), err)
		}
		out[strings.TrimSuffix(entry.Name(), ".json")] = draft
	}
	return out, nil
}

// Intent: Persist the current shared working body separately from durable
// revision snapshots so browser collaboration can converge on one draft
// without rewriting historical revision events. Source: DI-lusov; DI-zoruk
func (store *Store) SaveDraft(entityID string, draft PersistedDraft) error {
	body, err := json.Marshal(draft)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(store.draftPath, entityID+".json"), body, 0o644)
}

func (store *Store) Close() error {
	if store == nil {
		return nil
	}
	if store.events == nil && store.knowledgeItemMessages == nil && store.knowledgeApprovalMessages == nil && store.knowledgeEvidenceMessages == nil && store.operationalRunMessages == nil && store.operationalPlaceMessages == nil && store.operationalResourceMessages == nil && store.knowledgeLinkMessages == nil && store.knowledgeResponsibilityMessages == nil {
		return nil
	}
	var err error
	if store.events != nil {
		err = errors.Join(err, store.events.Close())
	}
	if store.knowledgeItemMessages != nil {
		err = errors.Join(err, store.knowledgeItemMessages.Close())
	}
	if store.knowledgeApprovalMessages != nil {
		err = errors.Join(err, store.knowledgeApprovalMessages.Close())
	}
	if store.knowledgeEvidenceMessages != nil {
		err = errors.Join(err, store.knowledgeEvidenceMessages.Close())
	}
	if store.operationalRunMessages != nil {
		err = errors.Join(err, store.operationalRunMessages.Close())
	}
	if store.operationalPlaceMessages != nil {
		err = errors.Join(err, store.operationalPlaceMessages.Close())
	}
	if store.operationalResourceMessages != nil {
		err = errors.Join(err, store.operationalResourceMessages.Close())
	}
	if store.knowledgeLinkMessages != nil {
		err = errors.Join(err, store.knowledgeLinkMessages.Close())
	}
	if store.knowledgeResponsibilityMessages != nil {
		err = errors.Join(err, store.knowledgeResponsibilityMessages.Close())
	}
	return err
}
