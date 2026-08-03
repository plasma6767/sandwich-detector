package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	ETHRPCURL string
}

func Load() (Config, error) {
	raw := os.Getenv("ETH_RPC_URL")
	if strings.TrimSpace(raw) == "" {
		return Config{}, fmt.Errorf("ETH_RPC_URL is not set (check your .env file)")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return Config{}, fmt.Errorf("ETH_RPC_URL is not a valid URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return Config{}, fmt.Errorf("ETH_RPC_URL is not a valid URL: missing scheme or host")
	}

	return Config{ETHRPCURL: raw}, nil
}
