# Sandwich Detector

Detects MEV sandwich attacks on Ethereum. I'm learning web3 backend -- explain concepts as you go, dont assume I know chain-specific terminology.

## Commands
- Run: `make run`
- Test: `make test`

## Structure
- `cmd/detector/` - entry point, wiring only
- `internal/` - all real logic

## Rules
- Config comes from environment variables. Never hardcode URLs or keys.
- Never write to `.env`
- Never log or print ETH_RPC_URL - the Alchemy key is in the URL path. Log the host only, never the full URL.
