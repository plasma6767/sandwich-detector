// Command detector is the entry point for the sandwich detector. It wires
// together the config and ethnode packages: load settings, connect to
// Ethereum, and confirm the connection works by fetching the latest block
// number. All real logic lives in internal/ - this file is wiring only.
package main

import (
	"context"
	"log"
	"time"

	"github.com/plasma6767/sandwich-detector/internal/config"
	"github.com/plasma6767/sandwich-detector/internal/ethnode"
	"github.com/plasma6767/sandwich-detector/internal/sandwich"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethnode.Dial(ctx, cfg.ETHRPCURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	// cfg.ETHRPCURL is never referenced again below this point - only
	// client.Host() is safe to log.
	blockNumber, err := client.LatestBlockNumber(ctx)
	if err != nil {
		log.Fatalf("failed to fetch latest block number: %v", err)
	}

	log.Printf("connected to %s, latest block: %d", client.Host(), blockNumber)

	txs, err := client.LatestBlockTransactions(ctx)
	if err != nil {
		log.Fatalf("failed to fetch latest block transactions: %v", err)
	}

	sandwiches := sandwich.Detect(txs)
	if len(sandwiches) == 0 {
		log.Printf("scanned %d transactions, no sandwiches found", len(txs))
		return
	}

	for _, s := range sandwiches {
		log.Printf("sandwich detected: front-run at position %d, victim at %d, back-run at %d",
			s.FrontRun.Position, s.Victim.Position, s.BackRun.Position)
	}
}
