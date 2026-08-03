package ethnode

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	inner *ethclient.Client
	host  string
}

func Dial(ctx context.Context, rawURL string) (*Client, error) {
	host, err := hostOnly(rawURL)
	if err != nil {
		return nil, fmt.Errorf("eth RPC URL is invalid: %w", err)
	}

	inner, err := ethclient.DialContext(ctx, rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth node at %s: %w", host, err)
	}

	return &Client{inner: inner, host: host}, nil
}

func (c *Client) Host() string {
	return c.host
}

func (c *Client) LatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.inner.BlockNumber(ctx)
}

func (c *Client) Close() {
	c.inner.Close()
}

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
