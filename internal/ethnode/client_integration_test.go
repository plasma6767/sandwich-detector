//go:build integration

// Integration test for the ethnode package. Unlike client_test.go, this
// makes a real network connection to Ethereum using ETH_RPC_URL, so it only
// runs when explicitly requested via `make test-integration` - never as
// part of plain `make test`.
package ethnode

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/plasma6767/sandwich-detector/internal/config"
)

// TestLatestBlockNumber_Integration connects to the real Ethereum node
// configured via ETH_RPC_URL and checks that a real, sane block number
// comes back. It's skipped automatically if no key is configured.
func TestLatestBlockNumber_Integration(t *testing.T) {
	if os.Getenv("ETH_RPC_URL") == "" {
		t.Skip("ETH_RPC_URL not set, skipping integration test")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := Dial(ctx, cfg.ETHRPCURL)
	if err != nil {
		t.Fatalf("Dial() failed: %v", err)
	}
	defer client.Close()

	if client.Host() == "" {
		t.Fatal("expected a non-empty host")
	}

	blockNumber, err := client.LatestBlockNumber(ctx)
	if err != nil {
		t.Fatalf("LatestBlockNumber() failed: %v", err)
	}

	// Ethereum mainnet passed this block height in 2023, so a real response
	// today should be comfortably above it. This is a loose sanity check,
	// not an exact value, since the real chain height keeps increasing.
	const sanityFloor = 18_000_000
	if blockNumber < sanityFloor {
		t.Fatalf("got block number %d, expected something above %d", blockNumber, sanityFloor)
	}
}
