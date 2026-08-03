// Tests for Load(). Run via `make test` (or `go test ./internal/config/...`).
package config

import "testing"

// TestLoad_MissingEnvVar checks that Load returns an error, rather than
// panicking or silently succeeding, when ETH_RPC_URL isn't set.
func TestLoad_MissingEnvVar(t *testing.T) {
	// t.Setenv sets an environment variable for the duration of this test
	// only, and Go restores its previous value automatically afterward.
	t.Setenv("ETH_RPC_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected an error when ETH_RPC_URL is not set, got nil")
	}
}

// TestLoad_ValidEnvVar checks that a well-formed URL is read back unchanged.
func TestLoad_ValidEnvVar(t *testing.T) {
	const url = "https://eth-mainnet.g.alchemy.com/v2/some-key"
	t.Setenv("ETH_RPC_URL", url)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ETHRPCURL != url {
		t.Fatalf("got ETHRPCURL %q, want %q", cfg.ETHRPCURL, url)
	}
}

// TestLoad_InvalidURL checks that a malformed value is rejected with an
// error instead of being accepted as-is.
func TestLoad_InvalidURL(t *testing.T) {
	t.Setenv("ETH_RPC_URL", "://not-a-valid-url")

	_, err := Load()
	if err == nil {
		t.Fatal("expected an error for a malformed URL, got nil")
	}
}
