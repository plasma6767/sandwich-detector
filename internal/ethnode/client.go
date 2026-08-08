// Package ethnode wraps the third-party "ethclient" library so the rest of
// the project can connect to Ethereum without touching that library
// directly. It also guarantees the full RPC URL (which contains a secret
// API key) never leaves this package - only the safe-to-log host is exposed.
package ethnode

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/core/types"
	// ethclient is part of go-ethereum, the reference Ethereum implementation,
	// and handles the Ethereum network protocol for us.
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client is an open connection to an Ethereum node.
type Client struct {
	inner *ethclient.Client // the underlying library connection
	host  string            // e.g. "eth-mainnet.g.alchemy.com" - safe to log
}

// Dial opens a connection to the Ethereum node at rawURL. ctx lets the
// caller enforce a time limit on the connection attempt.
func Dial(ctx context.Context, rawURL string) (*Client, error) {
	// Extract the safe-to-log host before doing anything else, so it's
	// available even if the connection attempt below fails.
	host, err := hostOnly(rawURL)
	if err != nil {
		return nil, fmt.Errorf("eth RPC URL is invalid: %w", err)
	}

	inner, err := ethclient.DialContext(ctx, rawURL)
	if err != nil {
		// Only host is referenced here, never rawURL, so the key can't leak
		// through an error message.
		return nil, fmt.Errorf("failed to connect to eth node at %s: %w", host, err)
	}

	return &Client{inner: inner, host: host}, nil
}

// Host returns the node's server name. It's the only way to learn which
// node this client is connected to - the original URL and its embedded key
// are never exposed.
func (c *Client) Host() string {
	return c.host
}

// LatestBlockNumber returns the height of the most recently produced block -
// i.e. how many blocks have been processed so far.
func (c *Client) LatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.inner.BlockNumber(ctx)
}

// LatestBlockTransactions fetches the most recently produced block and
// returns an ordered summary of every transaction in it. It first asks the
// node which network it's on (the chain ID), since that determines the
// signing rules needed to recover each transaction's sender.
func (c *Client) LatestBlockTransactions(ctx context.Context) ([]TxSummary, error) {
	chainID, err := c.inner.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chain ID: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	// A nil block number means "the latest block" to ethclient.
	block, err := c.inner.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest block: %w", err)
	}

	return TransactionsIn(block, signer)
}

// Close releases the underlying network connection.
func (c *Client) Close() {
	c.inner.Close()
}

// hostOnly extracts just the host (and port, if present) from a URL,
// discarding the path - which is where Alchemy's secret API key lives
// (e.g. "/v2/<key>").
func hostOnly(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("could not parse URL: %w", err)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("URL has no host")
	}
	return parsed.Host, nil
}
