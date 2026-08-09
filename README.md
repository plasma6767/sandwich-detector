# Sandwich Detector

Detects MEV sandwich attacks on Ethereum by analyzing transaction ordering within blocks.

## Setup

1. `cp .env.example .env`
2. Add your Alchemy API key to `.env`
3. `make run`

## Status

`make run` connects via Alchemy, fetches the latest block's transactions, and scans them for sandwich attacks (`internal/sandwich`). The current heuristic: the same sender hitting the same recipient twice with a different sender's transaction to that recipient in between. It doesn't yet decode swap data or check price impact, so it can false-positive on unrelated same-sender/same-recipient traffic.

## Learning: raw JSON-RPC vs ethclient

The app talks to Ethereum using `ethclient`, a Go library that wraps Ethereum's network protocol (JSON-RPC) for us. To see what that protocol actually looks like underneath the library, you can send a request by hand:

```bash
export $(grep -v '^#' .env | xargs)
curl -s -X POST "$ETH_RPC_URL" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

This sends the same "what's the latest block?" request our Go code sends, and gets back something like:

```json
{"jsonrpc":"2.0","id":1,"result":"0x1885ee5"}
```

`"0x1885ee5"` is a hex-encoded number - Ethereum returns numbers as hex text rather than plain decimal, partly because some values (like account balances) are too large for normal number types. This exact string is what `ethclient`'s `BlockNumber()` parses into the plain integer our code works with.
