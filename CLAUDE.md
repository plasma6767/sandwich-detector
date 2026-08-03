# Sandwich Detector

Detects MEV sandwich attacks on Ethereum. I'm learning web3 backend -- explain concepts as you go, dont assume I know chain-specific terminology.

## Commands
- Run: `make run`
- Test: `make test`
- Integration test (needs real `.env`, hits live RPC): `make test-integration`

## Structure
- `cmd/detector/` - entry point, wiring only
- `internal/` - all real logic

## Rules
- Config comes from environment variables. Never hardcode URLs or keys.
- Never write to `.env`
- Never log or print ETH_RPC_URL - the Alchemy key is in the URL path. Log the host only, never the full URL.

## Workflow
- Write tests before implementation. Define what "correct" and "broken" look like as a test before writing the code that satisfies it.
- Applies most strongly to logic-bearing code (parsers, validators, detection heuristics). Pure wiring (`cmd/detector/main.go`) doesn't need its own unit tests - it's proven by the integration/smoke test instead.
- Keep unit tests (fast, no network, run via `make test`) separate from integration tests (real external dependencies like the Alchemy RPC endpoint, gated behind a build tag, run via `make test-integration`). `make test` must always stay fast and safe to run anywhere without secrets.
