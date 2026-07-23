package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pgstore "github.com/computerscienceiscool/grid-examples/ex5-operational-knowledge-system/promisegrid/store"
)

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
	cas                             *pgstore.CASStore
	identity                        *RuntimeIdentity
}

type PersistedDraft struct {
	Body      string `json:"body"`
	BodyCID   string `json:"body_cid,omitempty"`
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
	casRoot := filepath.Join(root, "cas", "objects")
	if err := os.MkdirAll(casRoot, 0o755); err != nil {
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
	eventsFile, err := pgstore.OpenAppendOnlyLog(eventPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	events, err := readEvents(eventsFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close())
	}
	knowledgeItemPath := filepath.Join(root, "knowledge-item-messages.jsonl")
	knowledgeItemFile, err := pgstore.OpenAppendOnlyLog(knowledgeItemPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close())
	}
	knowledgeItemRecords, err := readSignedKnowledgeItemRecords(knowledgeItemFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	knowledgeApprovalPath := filepath.Join(root, "knowledge-approval-messages.jsonl")
	knowledgeApprovalFile, err := pgstore.OpenAppendOnlyLog(knowledgeApprovalPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close())
	}
	knowledgeApprovalRecords, err := readSignedKnowledgeApprovalRecords(knowledgeApprovalFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close())
	}
	knowledgeEvidencePath := filepath.Join(root, "knowledge-evidence-messages.jsonl")
	knowledgeEvidenceFile, err := pgstore.OpenAppendOnlyLog(knowledgeEvidencePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close())
	}
	knowledgeEvidenceRecords, err := readSignedKnowledgeEvidenceRecords(knowledgeEvidenceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close())
	}
	operationalRunPath := filepath.Join(root, "operational-run-messages.jsonl")
	operationalRunFile, err := pgstore.OpenAppendOnlyLog(operationalRunPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close())
	}
	operationalRunRecords, err := readSignedOperationalRunRecords(operationalRunFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close())
	}
	operationalPlacePath := filepath.Join(root, "operational-place-messages.jsonl")
	operationalPlaceFile, err := pgstore.OpenAppendOnlyLog(operationalPlacePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close())
	}
	operationalPlaceRecords, err := readSignedOperationalPlaceRecords(operationalPlaceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close())
	}
	operationalResourcePath := filepath.Join(root, "operational-resource-messages.jsonl")
	operationalResourceFile, err := pgstore.OpenAppendOnlyLog(operationalResourcePath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close())
	}
	operationalResourceRecords, err := readSignedOperationalResourceRecords(operationalResourceFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close())
	}
	knowledgeLinkPath := filepath.Join(root, "knowledge-link-messages.jsonl")
	knowledgeLinkFile, err := pgstore.OpenAppendOnlyLog(knowledgeLinkPath)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close())
	}
	knowledgeLinkRecords, err := readSignedKnowledgeLinkRecords(knowledgeLinkFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.Join(err, eventsFile.Close(), knowledgeItemFile.Close(), knowledgeApprovalFile.Close(), knowledgeEvidenceFile.Close(), operationalRunFile.Close(), operationalPlaceFile.Close(), operationalResourceFile.Close(), knowledgeLinkFile.Close())
	}
	knowledgeResponsibilityPath := filepath.Join(root, "knowledge-responsibility-messages.jsonl")
	knowledgeResponsibilityFile, err := pgstore.OpenAppendOnlyLog(knowledgeResponsibilityPath)
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
		cas:                             pgstore.NewCASStore(casRoot),
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
	events, err = pgstore.ReadJSONL[OperationalEvent](file)
	if err != nil {
		return nil, fmt.Errorf("decode event: %w", err)
	}
	return events, nil
}

func (store *Store) AppendEvent(event OperationalEvent) error {
	return pgstore.AppendJSONL(store.events, event)
}

func readSignedKnowledgeItemRecords(file *os.File) (records []SignedKnowledgeItemRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedKnowledgeItemRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode knowledge-item record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeItemRecord(record SignedKnowledgeItemRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.knowledgeItemMessages, record)
}

func (store *Store) LoadSignedKnowledgeItemRecordsAuthoritative() ([]SignedKnowledgeItemRecord, error) {
	records, err := readSignedKnowledgeItemRecords(store.knowledgeItemMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeItemRecords(records)
}

func readSignedKnowledgeApprovalRecords(file *os.File) (records []SignedKnowledgeApprovalRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedKnowledgeApprovalRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode knowledge-approval record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeApprovalRecord(record SignedKnowledgeApprovalRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.knowledgeApprovalMessages, record)
}

func (store *Store) LoadSignedKnowledgeApprovalRecordsAuthoritative() ([]SignedKnowledgeApprovalRecord, error) {
	records, err := readSignedKnowledgeApprovalRecords(store.knowledgeApprovalMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeApprovalRecords(records)
}

func readSignedKnowledgeEvidenceRecords(file *os.File) (records []SignedKnowledgeEvidenceRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedKnowledgeEvidenceRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode knowledge-evidence record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeEvidenceRecord(record SignedKnowledgeEvidenceRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.knowledgeEvidenceMessages, record)
}

func (store *Store) LoadSignedKnowledgeEvidenceRecordsAuthoritative() ([]SignedKnowledgeEvidenceRecord, error) {
	records, err := readSignedKnowledgeEvidenceRecords(store.knowledgeEvidenceMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeEvidenceRecords(records)
}

func readSignedOperationalRunRecords(file *os.File) (records []SignedOperationalRunRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedOperationalRunRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode operational-run record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedOperationalRunRecord(record SignedOperationalRunRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.operationalRunMessages, record)
}

func (store *Store) LoadSignedOperationalRunRecordsAuthoritative() ([]SignedOperationalRunRecord, error) {
	records, err := readSignedOperationalRunRecords(store.operationalRunMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedOperationalRunRecords(records)
}

func readSignedOperationalPlaceRecords(file *os.File) (records []SignedOperationalPlaceRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedOperationalPlaceRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode operational-place record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedOperationalPlaceRecord(record SignedOperationalPlaceRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.operationalPlaceMessages, record)
}

func (store *Store) LoadSignedOperationalPlaceRecordsAuthoritative() ([]SignedOperationalPlaceRecord, error) {
	records, err := readSignedOperationalPlaceRecords(store.operationalPlaceMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedOperationalPlaceRecords(records)
}

func readSignedOperationalResourceRecords(file *os.File) (records []SignedOperationalResourceRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedOperationalResourceRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode operational-resource record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedOperationalResourceRecord(record SignedOperationalResourceRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.operationalResourceMessages, record)
}

func (store *Store) LoadSignedOperationalResourceRecordsAuthoritative() ([]SignedOperationalResourceRecord, error) {
	records, err := readSignedOperationalResourceRecords(store.operationalResourceMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedOperationalResourceRecords(records)
}

func readSignedKnowledgeLinkRecords(file *os.File) (records []SignedKnowledgeLinkRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedKnowledgeLinkRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode knowledge-link record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeLinkRecord(record SignedKnowledgeLinkRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.knowledgeLinkMessages, record)
}

func (store *Store) LoadSignedKnowledgeLinkRecordsAuthoritative() ([]SignedKnowledgeLinkRecord, error) {
	records, err := readSignedKnowledgeLinkRecords(store.knowledgeLinkMessages)
	if err != nil {
		return nil, err
	}
	return store.hydrateSignedKnowledgeLinkRecords(records)
}

func readSignedKnowledgeResponsibilityRecords(file *os.File) (records []SignedKnowledgeResponsibilityRecord, err error) {
	records, err = pgstore.ReadJSONL[SignedKnowledgeResponsibilityRecord](file)
	if err != nil {
		return nil, fmt.Errorf("decode knowledge-responsibility record: %w", err)
	}
	return records, nil
}

func (store *Store) AppendSignedKnowledgeResponsibilityRecord(record SignedKnowledgeResponsibilityRecord) error {
	if err := store.cas.WriteEnvelopeBase64(record.EnvelopeCID, record.EnvelopeBase64); err != nil {
		return err
	}
	return pgstore.AppendJSONL(store.knowledgeResponsibilityMessages, record)
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
	cid, err := store.cas.WriteObject(data)
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
	body, err := store.cas.LoadObject(attachmentCID)
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

// Intent: Preserve app-level test access to the current CAS object location
// while the durable object-path rule itself now lives in the reusable store
// substrate. Source: DI-lemor
func (store *Store) casObjectPath(cid string) string {
	return store.cas.ObjectPath(cid)
}

func (store *Store) hydrateSignedKnowledgeItemRecords(records []SignedKnowledgeItemRecord) ([]SignedKnowledgeItemRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedKnowledgeItemRecord) string { return record.EnvelopeCID },
		func(record SignedKnowledgeItemRecord) string { return record.EnvelopeBase64 },
		func(record SignedKnowledgeItemRecord) string {
			return fmt.Sprintf("knowledge-item envelope %d", record.Sequence)
		},
		func(record *SignedKnowledgeItemRecord, base64Envelope string) { record.EnvelopeBase64 = base64Envelope },
	)
}

func (store *Store) hydrateSignedKnowledgeApprovalRecords(records []SignedKnowledgeApprovalRecord) ([]SignedKnowledgeApprovalRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedKnowledgeApprovalRecord) string { return record.EnvelopeCID },
		func(record SignedKnowledgeApprovalRecord) string { return record.EnvelopeBase64 },
		func(record SignedKnowledgeApprovalRecord) string {
			return fmt.Sprintf("knowledge-approval envelope %d", record.Sequence)
		},
		func(record *SignedKnowledgeApprovalRecord, base64Envelope string) {
			record.EnvelopeBase64 = base64Envelope
		},
	)
}

func (store *Store) hydrateSignedKnowledgeEvidenceRecords(records []SignedKnowledgeEvidenceRecord) ([]SignedKnowledgeEvidenceRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedKnowledgeEvidenceRecord) string { return record.EnvelopeCID },
		func(record SignedKnowledgeEvidenceRecord) string { return record.EnvelopeBase64 },
		func(record SignedKnowledgeEvidenceRecord) string {
			return fmt.Sprintf("knowledge-evidence envelope %d", record.Sequence)
		},
		func(record *SignedKnowledgeEvidenceRecord, base64Envelope string) {
			record.EnvelopeBase64 = base64Envelope
		},
	)
}

func (store *Store) hydrateSignedOperationalRunRecords(records []SignedOperationalRunRecord) ([]SignedOperationalRunRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedOperationalRunRecord) string { return record.EnvelopeCID },
		func(record SignedOperationalRunRecord) string { return record.EnvelopeBase64 },
		func(record SignedOperationalRunRecord) string {
			return fmt.Sprintf("operational-run envelope %d", record.Sequence)
		},
		func(record *SignedOperationalRunRecord, base64Envelope string) {
			record.EnvelopeBase64 = base64Envelope
		},
	)
}

func (store *Store) hydrateSignedOperationalPlaceRecords(records []SignedOperationalPlaceRecord) ([]SignedOperationalPlaceRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedOperationalPlaceRecord) string { return record.EnvelopeCID },
		func(record SignedOperationalPlaceRecord) string { return record.EnvelopeBase64 },
		func(record SignedOperationalPlaceRecord) string {
			return fmt.Sprintf("operational-place envelope %d", record.Sequence)
		},
		func(record *SignedOperationalPlaceRecord, base64Envelope string) {
			record.EnvelopeBase64 = base64Envelope
		},
	)
}

func (store *Store) hydrateSignedOperationalResourceRecords(records []SignedOperationalResourceRecord) ([]SignedOperationalResourceRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedOperationalResourceRecord) string { return record.EnvelopeCID },
		func(record SignedOperationalResourceRecord) string { return record.EnvelopeBase64 },
		func(record SignedOperationalResourceRecord) string {
			return fmt.Sprintf("operational-resource envelope %d", record.Sequence)
		},
		func(record *SignedOperationalResourceRecord, base64Envelope string) {
			record.EnvelopeBase64 = base64Envelope
		},
	)
}

func (store *Store) hydrateSignedKnowledgeLinkRecords(records []SignedKnowledgeLinkRecord) ([]SignedKnowledgeLinkRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedKnowledgeLinkRecord) string { return record.EnvelopeCID },
		func(record SignedKnowledgeLinkRecord) string { return record.EnvelopeBase64 },
		func(record SignedKnowledgeLinkRecord) string {
			return fmt.Sprintf("knowledge-link envelope %d", record.Sequence)
		},
		func(record *SignedKnowledgeLinkRecord, base64Envelope string) { record.EnvelopeBase64 = base64Envelope },
	)
}

func (store *Store) hydrateSignedKnowledgeResponsibilityRecords(records []SignedKnowledgeResponsibilityRecord) ([]SignedKnowledgeResponsibilityRecord, error) {
	return pgstore.HydrateAuthoritativeEnvelopes(
		store.cas,
		records,
		func(record SignedKnowledgeResponsibilityRecord) string { return record.EnvelopeCID },
		func(record SignedKnowledgeResponsibilityRecord) string { return record.EnvelopeBase64 },
		func(record SignedKnowledgeResponsibilityRecord) string {
			return fmt.Sprintf("knowledge-responsibility envelope %d", record.Sequence)
		},
		func(record *SignedKnowledgeResponsibilityRecord, base64Envelope string) {
			record.EnvelopeBase64 = base64Envelope
		},
	)
}

// Intent: Restore the shared browser working bodies on startup without mixing
// them into the append-only event log, so durable revision history stays
// explicit while collaborative drafting can resume after a restart. Source:
// DI-lusov; DI-zoruk; DI-zunep
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
		entityID := strings.TrimSuffix(entry.Name(), ".json")
		draft, changed, err := store.hydrateDraftManifest(entityID, draft)
		if err != nil {
			return nil, fmt.Errorf("hydrate draft %s: %w", entry.Name(), err)
		}
		if changed {
			if err := store.saveDraftManifest(entityID, draft); err != nil {
				return nil, fmt.Errorf("rewrite draft %s: %w", entry.Name(), err)
			}
		}
		out[entityID] = draft
	}
	return out, nil
}

// Intent: Persist the current shared working body separately from durable
// revision snapshots so browser collaboration can converge on one draft
// without rewriting historical revision events. Source: DI-lusov; DI-zoruk;
// DI-zunep
func (store *Store) SaveDraft(entityID string, draft PersistedDraft) error {
	bodyCID, err := store.cas.WriteObject([]byte(draft.Body))
	if err != nil {
		return err
	}
	draft.BodyCID = bodyCID
	return store.saveDraftManifest(entityID, draft)
}

// Intent: Treat shared draft bodies as content-addressed local runtime state so
// reload trusts CAS over any stale inline manifest body while still backfilling
// older draft files during migration. Source: DI-zunep
func (store *Store) hydrateDraftManifest(entityID string, draft PersistedDraft) (PersistedDraft, bool, error) {
	changed := false
	if strings.TrimSpace(draft.BodyCID) == "" {
		bodyCID, err := store.cas.WriteObject([]byte(draft.Body))
		if err != nil {
			return PersistedDraft{}, false, err
		}
		draft.BodyCID = bodyCID
		changed = true
	}
	body, err := store.cas.LoadObject(draft.BodyCID)
	if err != nil {
		return PersistedDraft{}, false, err
	}
	if draft.Body != string(body) {
		draft.Body = string(body)
		changed = true
	}
	return draft, changed, nil
}

// Intent: Keep one small per-item local manifest pointing at the authoritative
// CAS-backed draft body so embodiments can resume shared drafting without
// promoting drafts into a new peer-visible family. Source: DI-zunep
func (store *Store) saveDraftManifest(entityID string, draft PersistedDraft) error {
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
