package kernel

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/computerscienceiscool/grid-examples/ex6-operational-knowledge-agent-runtime/packages"
)

const workflowManifestName = "workflow.json"

type WorkflowManifest struct {
	ID                string   `json:"id"`
	Version           string   `json:"version"`
	Summary           string   `json:"summary"`
	RequiredPackages  []string `json:"required_packages"`
	RequiredProtocols []string `json:"required_protocols"`
	Adapter           string   `json:"adapter,omitempty"`
	InputPCID         string   `json:"input_pcid,omitempty"`
	OutputPCID        string   `json:"output_pcid,omitempty"`
	InputSchema       string   `json:"input_schema,omitempty"`
	OutputSchema      string   `json:"output_schema,omitempty"`
}

// legacyWorkflowAdapterPCIDs identifies the retired adapter contracts emitted
// before canonical embedded schemas. Intent: Preserve retained artifacts while
// preventing new captures from extending an obsolete contract generation.
// Source: DI-lumek
var legacyWorkflowAdapterPCIDs = map[string]struct{}{
	"bafkreihjwthblvvsaxlngupwghkshl2lnwgcj5txrr3qpelxtb76stg7ae": {},
	"bafkreihjhnfom2j2avcjjujcbvy22dbayjkdmkjj6ca3fbfjlm7vm23nxy": {},
	"bafkreidndb65kuarxuv3eue6ij3qupblgfuvjmm6v4s5vcfo2d7acbbcwq": {},
	"bafkreie7k5xcmmvygwh5fqsbvruh5iivxsyduost7bzrwphhfhudx7ga4q": {},
	"bafkreig4pegj6ckn6hg2yci7i5tt4vcknd4e25ul7a7ugyel5gejjz7iti": {},
	"bafkreigkvvjof6vhurueeod4mqtghfosiwunlskzkuhjuskiy66cnx5yua": {},
	"bafkreiboes7s6tcaebjnlibd7fkwj62typezjdsipskyafvtuzf74ypx3i": {},
	"bafkreicfalhnnj67rctw63c6j4w6x5l7ntmqtvyndmqi26vc46ukavmiha": {},
	"bafkreih2mechxf4slowhcag6xac5fqn7wy7tw6pw2lj5nvhc5nthmrx54e": {},
	"bafkreib6qcz4g3lsc4yzfulqihsbczc4wkpo3fwm5f7dvgrznv4qubwppe": {},
	"bafkreibyegjb3p52b3hzf4lw3jwqu3hktxockiamrt2a3bee2gwdl46ja4": {},
	"bafkreib3rq3zyljjn4v7tunm2xqpy26i7mpr24bko2sdof524p7xnqnjo4": {},
	"bafkreigkwagey45deeh6cc2hirr53avia3rbt2vgixxcbvwiccctnobi2y": {},
	"bafkreifrf4xznrekx4lohyueahudtz6ju7qhgodpnbyrd4rytjauo7qgsm": {},
}

// WorkflowStatus is one operator-facing view of a retained workflow artifact.
type WorkflowStatus struct {
	Workflow Workflow         `json:"workflow"`
	EventCID string           `json:"event_cid"`
	Manifest WorkflowManifest `json:"manifest"`
	Ready    bool             `json:"ready"`
	Reason   string           `json:"reason,omitempty"`
}

// WorkflowVerification is the operator-facing execution readiness view for one
// retained artifact. It separates successful structural verification from the
// local activation and adapter conditions required to start a run.
type WorkflowVerification struct {
	Manifest          WorkflowManifest `json:"manifest"`
	Contract          string           `json:"contract"`
	AdapterAvailable  bool             `json:"adapter_available"`
	SchemaCASReady    bool             `json:"schema_cas_ready"`
	RegistryAllowed   bool             `json:"registry_allowed"`
	ImageAvailable    bool             `json:"image_available"`
	EligibleToExecute bool             `json:"eligible_to_execute"`
	Reason            string           `json:"reason,omitempty"`
}

func (m WorkflowManifest) Validate() error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.Version) == "" || strings.TrimSpace(m.Summary) == "" {
		return errors.New("workflow id, version, and summary are required")
	}
	for _, pcid := range []string{m.InputPCID, m.OutputPCID} {
		if pcid != "" {
			if err := validateWorkflowPCID(pcid); err != nil {
				return err
			}
		}
	}
	// Intent: Reject a half-declared execution contract at capture and verify
	// time instead of deferring an avoidable failure until a run starts.
	// Source: DI-lumek
	if m.Adapter != "" || m.InputPCID != "" || m.OutputPCID != "" {
		if strings.TrimSpace(m.Adapter) == "" || strings.TrimSpace(m.InputPCID) == "" || strings.TrimSpace(m.OutputPCID) == "" {
			return errors.New("workflow executable adapter, input pCID, and output pCID must be declared together")
		}
	}
	if m.InputSchema != "" || m.OutputSchema != "" {
		if strings.TrimSpace(m.InputSchema) == "" || strings.TrimSpace(m.OutputSchema) == "" {
			return errors.New("workflow input and output schema paths must be declared together")
		}
		for _, schema := range []string{m.InputSchema, m.OutputSchema} {
			clean := filepath.Clean(schema)
			if filepath.IsAbs(clean) || clean == "." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
				return errors.New("workflow schema path must stay inside the artifact")
			}
		}
	}
	seen := map[string]bool{}
	for _, v := range append(append([]string{}, m.RequiredPackages...), m.RequiredProtocols...) {
		if strings.TrimSpace(v) == "" || seen[v] {
			return errors.New("workflow dependencies must be non-empty and unique")
		}
		seen[v] = true
	}
	return nil
}

// CaptureWorkflowDir archives a validated workflow deterministically, stores it in CAS, and imports it inactive.
// Intent: Make the workflow basket generic and content-addressed without granting execution authority. Source: DI-lovek
func (runtime *Runtime) CaptureWorkflowDir(directory, alias string) (Workflow, error) {
	body, manifest, err := archiveWorkflowDir(directory)
	if err != nil {
		return Workflow{}, err
	}
	if err = manifest.Validate(); err != nil {
		return Workflow{}, err
	}
	if _, legacy := legacyWorkflowAdapterPCIDs[manifest.InputPCID]; legacy {
		return Workflow{}, errors.New("new workflow capture cannot declare a retired adapter pCID")
	}
	if _, legacy := legacyWorkflowAdapterPCIDs[manifest.OutputPCID]; legacy {
		return Workflow{}, errors.New("new workflow capture cannot declare a retired adapter pCID")
	}
	artifact, err := runtime.cas.PutCID(body)
	if err != nil {
		return Workflow{}, err
	}
	manifest, err = runtime.WorkflowManifest(artifact.String())
	if err != nil {
		return Workflow{}, err
	}
	if err = runtime.ImportWorkflow(Workflow{ID: alias, ArtifactCID: artifact.String()}); err != nil {
		return Workflow{}, err
	}
	return runtime.workflow(alias)
}
func archiveWorkflowDir(directory string) ([]byte, WorkflowManifest, error) {
	files := []string{}
	var total int64
	err := filepath.WalkDir(directory, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == directory {
			return nil
		}
		if e.Type()&os.ModeSymlink != 0 || (!e.Type().IsRegular() && !e.IsDir()) {
			return fmt.Errorf("unsafe workflow path: %s", path)
		}
		if e.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(directory, path)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			return errors.New("workflow path escapes source directory")
		}
		info, err := e.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		if len(files) >= 1000 || total > 10<<20 {
			return errors.New("workflow exceeds loader size limit")
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, WorkflowManifest{}, err
	}
	slices.Sort(files)
	var out bytes.Buffer
	w := tar.NewWriter(&out)
	var manifest WorkflowManifest
	for _, rel := range files {
		body, err := os.ReadFile(filepath.Join(directory, rel))
		if err != nil {
			return nil, manifest, err
		}
		name := filepath.ToSlash(rel)
		if err = w.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(body)), Format: tar.FormatUSTAR}); err != nil {
			return nil, manifest, err
		}
		if _, err = w.Write(body); err != nil {
			return nil, manifest, err
		}
		if name == workflowManifestName {
			if err = json.Unmarshal(body, &manifest); err != nil {
				return nil, manifest, err
			}
		}
	}
	if err = w.Close(); err != nil {
		return nil, manifest, err
	}
	if manifest.ID == "" {
		return nil, manifest, errors.New("workflow.json is required")
	}
	return out.Bytes(), manifest, nil
}
func (runtime *Runtime) workflow(alias string) (Workflow, error) {
	for _, w := range runtime.Workflows() {
		if w.ID == alias {
			return w, nil
		}
	}
	return Workflow{}, errors.New("workflow is not imported")
}
func (runtime *Runtime) WorkflowManifest(artifactID string) (WorkflowManifest, error) {
	artifact, err := runtime.workflowArtifactCID(artifactID)
	if err != nil {
		return WorkflowManifest{}, err
	}
	body, err := runtime.cas.GetCID(artifact)
	if err != nil {
		return WorkflowManifest{}, err
	}
	r := tar.NewReader(bytes.NewReader(body))
	for {
		h, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return WorkflowManifest{}, err
		}
		if h.Name == workflowManifestName {
			raw, err := io.ReadAll(r)
			if err != nil {
				return WorkflowManifest{}, err
			}
			var m WorkflowManifest
			if err = json.Unmarshal(raw, &m); err != nil {
				return WorkflowManifest{}, err
			}
			if err := m.Validate(); err != nil {
				return WorkflowManifest{}, err
			}
			if err := runtime.verifyWorkflowSchemas(body, m); err != nil {
				return WorkflowManifest{}, err
			}
			return m, nil
		}
	}
	return WorkflowManifest{}, errors.New("workflow.json is missing from artifact")
}

// verifyWorkflowSchemas resolves each declared schema from the immutable
// artifact and retains it under its pCID locally for independent inspection.
// Intent: A received artifact must carry the exact schema bytes that define
// its input and output contracts. Source: DI-lumek
func (runtime *Runtime) verifyWorkflowSchemas(artifact []byte, manifest WorkflowManifest) error {
	if manifest.InputSchema == "" && manifest.OutputSchema == "" {
		return nil
	}
	for _, schema := range []struct{ path, pcid string }{{manifest.InputSchema, manifest.InputPCID}, {manifest.OutputSchema, manifest.OutputPCID}} {
		r := tar.NewReader(bytes.NewReader(artifact))
		for {
			header, err := r.Next()
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("workflow schema %q is missing from artifact", schema.path)
			}
			if err != nil {
				return err
			}
			if header.Name != schema.path {
				continue
			}
			body, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			actual, err := runtime.cas.PutCID(body)
			if err != nil {
				return err
			}
			if actual.String() != schema.pcid {
				return fmt.Errorf("workflow schema %q CID %s does not match declared pCID %s", schema.path, actual, schema.pcid)
			}
			break
		}
	}
	return nil
}

// ExtractWorkflow writes a retained workflow artifact into a new safe directory.
// Intent: Let operators inspect exact CAS workflow bytes without allowing an
// archive to escape the chosen destination. Source: DI-lovek
func (runtime *Runtime) ExtractWorkflow(aliasOrCID string, destination string) error {
	artifactID := aliasOrCID
	if workflow, err := runtime.workflow(aliasOrCID); err == nil {
		artifactID = workflow.ArtifactCID
	}
	artifact, err := runtime.workflowArtifactCID(artifactID)
	if err != nil {
		return err
	}
	body, err := runtime.cas.GetCID(artifact)
	if err != nil {
		return err
	}
	if _, err := os.Stat(destination); err == nil || !errors.Is(err, os.ErrNotExist) {
		return errors.New("workflow extraction destination must not exist")
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}
	reader := tar.NewReader(bytes.NewReader(body))
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		name := filepath.Clean(filepath.FromSlash(header.Name))
		if filepath.IsAbs(name) || name == "." || strings.HasPrefix(name, ".."+string(filepath.Separator)) || header.Typeflag != tar.TypeReg {
			return errors.New("workflow archive contains an unsafe entry")
		}
		path := filepath.Join(destination, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, reader); err != nil {
			if closeErr := file.Close(); closeErr != nil {
				return errors.Join(err, closeErr)
			}
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
}
func (runtime *Runtime) validateWorkflowDependencies(w Workflow) error {
	m, err := runtime.WorkflowManifest(w.ArtifactCID)
	if err != nil && (errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || err.Error() == "workflow.json is missing from artifact") {
		// Historical direct-CAS artifacts predate the loader manifest and retain
		// their existing lifecycle behavior. Source: DI-lovek
		return nil
	}
	if err != nil {
		return err
	}
	for _, id := range m.RequiredPackages {
		if _, ok := runtime.PackageManifest(id); !ok {
			return fmt.Errorf("required package is not active: %s", id)
		}
	}
	for _, p := range m.RequiredProtocols {
		if len(runtime.ProtocolRoutesForProtocol(p)) == 0 {
			return fmt.Errorf("required protocol has no route: %s", p)
		}
	}
	return nil
}

// VerifyWorkflow validates one retained artifact without changing its lifecycle state.
// Intent: Let operators check local readiness before choosing activation. Source: DI-lovek
func (runtime *Runtime) VerifyWorkflow(aliasOrCID string) (WorkflowManifest, error) {
	artifactID := aliasOrCID
	if workflow, err := runtime.workflow(aliasOrCID); err == nil {
		artifactID = workflow.ArtifactCID
	}
	manifest, err := runtime.WorkflowManifest(artifactID)
	if err != nil {
		return WorkflowManifest{}, err
	}
	if err := runtime.workflowDependencyError(manifest); err != nil {
		return WorkflowManifest{}, err
	}
	return manifest, nil
}

func (runtime *Runtime) workflowDependencyError(manifest WorkflowManifest) error {
	for _, id := range manifest.RequiredPackages {
		if _, ok := runtime.PackageManifest(id); !ok {
			return fmt.Errorf("required package is not active: %s", id)
		}
	}
	for _, protocol := range manifest.RequiredProtocols {
		if len(runtime.ProtocolRoutesForProtocol(protocol)) == 0 {
			return fmt.Errorf("required protocol has no route: %s", protocol)
		}
	}
	return nil
}

// VerifyWorkflowReadiness reports whether a verified artifact can execute on
// this runtime without changing its lifecycle state. Intent: Make the retained
// v1 compatibility boundary visible before an operator starts work. Source:
// DI-lumek
func (runtime *Runtime) VerifyWorkflowReadiness(aliasOrCID string) (WorkflowVerification, error) {
	artifactID := aliasOrCID
	if workflow, err := runtime.workflow(aliasOrCID); err == nil {
		artifactID = workflow.ArtifactCID
	}
	manifest, err := runtime.WorkflowManifest(artifactID)
	if err != nil {
		return WorkflowVerification{}, err
	}
	verification := WorkflowVerification{
		Manifest:         manifest,
		Contract:         "canonical",
		AdapterAvailable: runtime.workflowAdapterAvailable(manifest),
		SchemaCASReady:   manifest.InputSchema != "" && manifest.OutputSchema != "",
	}
	registryRequired := false
	if adapter, installed := runtime.workflowAdapters[manifest.Adapter]; installed {
		if host, hostErr := packages.RegistryHostFromImage(adapter.Image); hostErr == nil {
			registryRequired = true
			for _, allowed := range runtime.RegistryAllowList() {
				if allowed == host {
					verification.RegistryAllowed = true
					break
				}
			}
			available, availableErr := packages.ImageAvailable(context.Background(), adapter.Image)
			if availableErr == nil {
				verification.ImageAvailable = available
			}
		}
	}
	if _, legacyInput := legacyWorkflowAdapterPCIDs[manifest.InputPCID]; legacyInput {
		verification.Contract = "retained-v1"
	}
	if _, legacyOutput := legacyWorkflowAdapterPCIDs[manifest.OutputPCID]; legacyOutput {
		verification.Contract = "retained-v1"
	}
	if err := runtime.workflowDependencyError(manifest); err != nil {
		verification.Reason = err.Error()
		return verification, nil
	}
	workflow, err := runtime.workflow(aliasOrCID)
	if err != nil {
		for _, candidate := range runtime.Workflows() {
			if candidate.ArtifactCID == aliasOrCID {
				workflow = candidate
				err = nil
				break
			}
		}
	}
	if err != nil {
		verification.Reason = "workflow is not imported"
		return verification, nil
	}
	if workflow.State != WorkflowActive {
		verification.Reason = "workflow is not active"
		return verification, nil
	}
	if manifest.Adapter == "" || manifest.InputPCID == "" || manifest.OutputPCID == "" {
		verification.Reason = "workflow does not declare an executable adapter"
		return verification, nil
	}
	if !verification.AdapterAvailable {
		verification.Reason = "workflow adapter is unavailable"
		return verification, nil
	}
	if registryRequired && !verification.RegistryAllowed {
		verification.Reason = "adapter registry is not allowed"
		return verification, nil
	}
	if registryRequired && !verification.ImageAvailable {
		verification.Reason = "adapter image is not available locally"
		return verification, nil
	}
	if verification.Contract == "canonical" && !verification.SchemaCASReady {
		verification.Reason = "canonical workflow schemas are not ready in CAS"
		return verification, nil
	}
	verification.EligibleToExecute = true
	return verification, nil
}

func (runtime *Runtime) workflowAdapterAvailable(manifest WorkflowManifest) bool {
	if adapter, installed := runtime.workflowAdapters[manifest.Adapter]; installed {
		return adapter.InputPCID == manifest.InputPCID && adapter.OutputPCID == manifest.OutputPCID
	}
	return runtime.workflowOps[manifest.Adapter] != nil
}

// InspectWorkflowStatus summarizes local lifecycle and dependency readiness without mutation.
// Intent: Give operators one concise basket view before changing availability. Source: DI-lovek
func (runtime *Runtime) InspectWorkflowStatus(aliasOrCID string) (WorkflowStatus, error) {
	workflow, err := runtime.workflow(aliasOrCID)
	if err != nil {
		for _, candidate := range runtime.Workflows() {
			if candidate.ArtifactCID == aliasOrCID {
				workflow = candidate
				err = nil
				break
			}
		}
	}
	if err != nil {
		return WorkflowStatus{}, err
	}
	head, ok := runtime.workflows.headCID(workflow.ID)
	if !ok {
		return WorkflowStatus{}, errors.New("workflow lifecycle head is missing")
	}
	status := WorkflowStatus{Workflow: workflow, EventCID: head.String()}
	manifest, err := runtime.VerifyWorkflow(workflow.ID)
	if err != nil {
		status.Reason = err.Error()
		return status, nil
	}
	status.Manifest, status.Ready = manifest, true
	return status, nil
}
