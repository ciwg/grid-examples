package packages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// RegistryHostFromImage parses the registry hostname from a portable OCI
// digest reference. Intent: Never authorize image acquisition by raw-prefix
// comparison. Source: DI-hapak
func RegistryHostFromImage(image string) (string, error) {
	image = strings.TrimSpace(image)
	parts := strings.Split(image, "@sha256:")
	if len(parts) != 2 || parts[0] == "" || len(parts[1]) != 64 {
		return "", errors.New("image must be a registry-qualified sha256 digest reference")
	}
	for _, character := range parts[1] {
		if !(character >= '0' && character <= '9') && !(character >= 'a' && character <= 'f') {
			return "", errors.New("image digest must be lowercase hexadecimal")
		}
	}
	segments := strings.Split(parts[0], "/")
	if len(segments) < 2 || strings.TrimSpace(segments[0]) == "" {
		return "", errors.New("image must include registry host and repository")
	}
	host := strings.ToLower(segments[0])
	if strings.ContainsAny(host, "\\*?") {
		return "", errors.New("image registry host is invalid")
	}
	for _, segment := range segments[1:] {
		if segment == "" || strings.ContainsAny(segment, "@:\\*? ") {
			return "", errors.New("image repository path is invalid")
		}
	}
	return host, nil
}

// PullImage acquires an exact OCI digest and confirms Docker reports that same
// repository digest locally. Source: DI-zivut
func PullImage(ctx context.Context, image string) error {
	if _, err := RegistryHostFromImage(image); err != nil {
		return err
	}
	if output, err := exec.CommandContext(ctx, "docker", "pull", image).CombinedOutput(); err != nil {
		return fmt.Errorf("docker pull: %w: %s", err, strings.TrimSpace(string(output)))
	}
	available, err := ImageAvailable(ctx, image)
	if err != nil {
		return err
	}
	if available {
		return nil
	}
	return errors.New("docker did not retain the requested image digest")
}

// ImageAvailable reports whether Docker has the exact requested repository
// digest locally without contacting a registry. Source: DI-zivut
func ImageAvailable(ctx context.Context, image string) (bool, error) {
	if _, err := RegistryHostFromImage(image); err != nil {
		return false, err
	}
	output, err := exec.CommandContext(ctx, "docker", "image", "inspect", "--format", "{{json .RepoDigests}}", image).Output()
	if err != nil {
		return false, nil
	}
	var digests []string
	if err := json.Unmarshal(output, &digests); err != nil {
		return false, fmt.Errorf("decode local image digests: %w", err)
	}
	for _, digest := range digests {
		if digest == image {
			return true, nil
		}
	}
	return false, nil
}
