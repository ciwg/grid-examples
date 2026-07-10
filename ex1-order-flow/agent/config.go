package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	Role          string
	DataDir       string
	KernelAddr    string
	CollectorAddr string
	Timeout       time.Duration
}

func ConfigFromEnv(role string) (Config, error) {
	dataDir := os.Getenv("PG_DATA_DIR")
	if dataDir == "" {
		return Config{}, fmt.Errorf("PG_DATA_DIR is required")
	}
	kernelAddr := os.Getenv("PG_KERNEL_ADDR")
	if kernelAddr == "" {
		return Config{}, fmt.Errorf("PG_KERNEL_ADDR is required")
	}
	collectorAddr := os.Getenv("PG_COLLECTOR_ADDR")
	if collectorAddr == "" {
		return Config{}, fmt.Errorf("PG_COLLECTOR_ADDR is required")
	}
	timeout := 5 * time.Second
	if raw := os.Getenv("PG_TIMEOUT_SECONDS"); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse PG_TIMEOUT_SECONDS: %w", err)
		}
		timeout = time.Duration(seconds) * time.Second
	}
	return Config{
		Role:          role,
		DataDir:       filepath.Clean(dataDir),
		KernelAddr:    kernelAddr,
		CollectorAddr: collectorAddr,
		Timeout:       timeout,
	}, nil
}
