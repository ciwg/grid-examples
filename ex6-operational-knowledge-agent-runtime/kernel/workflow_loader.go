package kernel

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const workflowManifestName = "workflow.json"

type WorkflowManifest struct {
	ID                string   `json:"id"`
	Version           string   `json:"version"`
	Summary           string   `json:"summary"`
	RequiredPackages  []string `json:"required_packages"`
	RequiredProtocols []string `json:"required_protocols"`
}

func (m WorkflowManifest) Validate() error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.Version) == "" || strings.TrimSpace(m.Summary) == "" {
		return errors.New("workflow id, version, and summary are required")
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
	artifact, err := runtime.cas.PutCID(body)
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
			return m, m.Validate()
		}
	}
	return WorkflowManifest{}, errors.New("workflow.json is missing from artifact")
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
	for _, id := range manifest.RequiredPackages {
		if _, ok := runtime.PackageManifest(id); !ok {
			return WorkflowManifest{}, fmt.Errorf("required package is not active: %s", id)
		}
	}
	for _, protocol := range manifest.RequiredProtocols {
		if len(runtime.ProtocolRoutesForProtocol(protocol)) == 0 {
			return WorkflowManifest{}, fmt.Errorf("required protocol has no route: %s", protocol)
		}
	}
	return manifest, nil
}
