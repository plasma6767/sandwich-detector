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

	"github.com/ethereum/go-ethereum/common"
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

// TestLatestBlockTransactions_Integration fetches the real, current latest
// block and checks that TransactionsIn's output looks sane against it. The
// real chain's contents change constantly, so this checks properties that
// must always hold (valid ordering, non-blank senders) rather than exact
// values - unlike client_test.go, which checks exact values against a
// block we made up ourselves.
func TestLatestBlockTransactions_Integration(t *testing.T) {
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

	summaries, err := client.LatestBlockTransactions(ctx)
	if err != nil {
		t.Fatalf("LatestBlockTransactions() failed: %v", err)
	}

	for i, s := range summaries {
		if s.Position != i {
			t.Fatalf("transaction at index %d has Position %d, want %d", i, s.Position, i)
		}
		if s.From == (common.Address{}) {
			t.Fatalf("transaction %d has a blank From address", i)
		}
	}
}
