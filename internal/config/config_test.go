package config

import "testing"

func TestLoad_MissingEnvVar(t *testing.T) {
	t.Setenv("ETH_RPC_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected an error when ETH_RPC_URL is not set, got nil")
	}
}

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

func TestLoad_InvalidURL(t *testing.T) {
	t.Setenv("ETH_RPC_URL", "://not-a-valid-url")

	_, err := Load()
	if err == nil {
		t.Fatal("expected an error for a malformed URL, got nil")
	}
}
