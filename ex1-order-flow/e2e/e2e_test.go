package e2e_test

import (
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/accounting"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/agent"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/analyzer"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/carrier"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/collector"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/intake"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/kernel"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/seller"
	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/warehouse"
)

func TestHappyPath(t *testing.T) {
	result, summary := runScenario(t, "happy-path.json", 2*time.Second)
	if result.OrderStatus != "fulfilled" {
		t.Fatalf("order status = %q, want fulfilled", result.OrderStatus)
	}
	if result.PromiseStatus != "kept" {
		t.Fatalf("promise status = %q, want kept", result.PromiseStatus)
	}
	if result.TrackingNumber == "" {
		t.Fatalf("tracking number is empty")
	}
	if result.LedgerEntryID == "" {
		t.Fatalf("ledger entry id is empty")
	}
	if summary.ArtifactCount == 0 || summary.MessageCount == 0 {
		t.Fatalf("summary = %#v, want retained artifacts and messages", summary)
	}
}

func TestWarehouseRefusal(t *testing.T) {
	result, _ := runScenario(t, "warehouse-refusal.json", 2*time.Second)
	if result.OrderStatus != "refused" || result.FailureStage != "warehouse" {
		t.Fatalf("warehouse refusal result = %#v", result)
	}
}

func TestAccountingRefusal(t *testing.T) {
	result, _ := runScenario(t, "accounting-refusal.json", 2*time.Second)
	if result.OrderStatus != "refused" || result.FailureStage != "accounting" {
		t.Fatalf("accounting refusal result = %#v", result)
	}
}

func TestCarrierTimeout(t *testing.T) {
	result, _ := runScenario(t, "carrier-timeout.json", 1*time.Second)
	if result.OrderStatus != "failed" || result.FailureStage != "carrier" {
		t.Fatalf("carrier timeout result = %#v", result)
	}
}

func TestEmptyItemsFailsValidation(t *testing.T) {
	result, _ := runScenario(t, "empty-items.json", 2*time.Second)
	if result.OrderStatus != "failed" || result.FailureStage != "seller_validation" {
		t.Fatalf("empty-items result = %#v", result)
	}
}

func runScenario(t *testing.T, fixtureName string, timeout time.Duration) (intake.Result, analyzer.Summary) {
	t.Helper()
	root := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	collectorAddr := "127.0.0.1:" + freePort(t)
	kernelAddr := "127.0.0.1:" + freePort(t)
	run := func(name string, fn func(context.Context, agent.Config) error) {
		t.Helper()
		cfg := agent.Config{
			Role:          name,
			DataDir:       filepath.Join(root, name),
			KernelAddr:    kernelAddr,
			CollectorAddr: collectorAddr,
			Timeout:       timeout,
		}
		go func() {
			err := fn(ctx, cfg)
			if err != nil && ctx.Err() == nil {
				t.Errorf("%s run: %v", name, err)
			}
		}()
	}
	go func() {
		service := &collector.Service{Address: collectorAddr, DataDir: filepath.Join(root, "collector")}
		if err := service.Run(ctx); err != nil && ctx.Err() == nil {
			t.Errorf("collector run: %v", err)
		}
	}()
	go func() {
		server := &kernel.Server{Address: kernelAddr}
		if err := server.Run(ctx); err != nil && ctx.Err() == nil {
			t.Errorf("kernel run: %v", err)
		}
	}()
	run("warehouse", warehouse.Run)
	run("accounting", accounting.Run)
	run("carrier", carrier.Run)
	run("seller", seller.Run)
	time.Sleep(500 * time.Millisecond)
	result, err := intake.Run(ctx, agent.Config{
		Role:          "intake",
		DataDir:       filepath.Join(root, "intake"),
		KernelAddr:    kernelAddr,
		CollectorAddr: collectorAddr,
		Timeout:       timeout,
	}, filepath.Join("..", "fixtures", fixtureName))
	if err != nil {
		t.Fatalf("intake run: %v", err)
	}
	time.Sleep(500 * time.Millisecond)
	summary, err := analyzer.Analyze(root)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	cancel()
	return result, summary
}

func freePort(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen free port: %v", err)
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			t.Errorf("close free-port listener: %v", closeErr)
		}
	}()
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}
	return port
}
