// Package config loads application settings from environment variables.
// It's the only package that reads environment variables directly - every
// other package receives configuration as plain data.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config holds the settings this program needs to run.
type Config struct {
	// ETHRPCURL is the Ethereum node endpoint (provided by Alchemy) used to
	// read chain data. It contains a secret API key in its path, so it must
	// never be logged in full - only the host should be logged.
	ETHRPCURL string
}

// Load reads ETH_RPC_URL from the environment and validates it, returning an
// error if it's missing or malformed rather than letting bad config surface
// later as a confusing connection failure.
func Load() (Config, error) {
	raw := os.Getenv("ETH_RPC_URL")
	if strings.TrimSpace(raw) == "" {
		return Config{}, fmt.Errorf("ETH_RPC_URL is not set (check your .env file)")
	}

	// url.Parse is lenient about malformed input, so we also check for a
	// scheme (e.g. "https") and host rather than trusting a nil error alone.
	parsed, err := url.Parse(raw)
	if err != nil {
		return Config{}, fmt.Errorf("ETH_RPC_URL is not a valid URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return Config{}, fmt.Errorf("ETH_RPC_URL is not a valid URL: missing scheme or host")
	}

	return Config{ETHRPCURL: raw}, nil
}
