// Tests for Detect(). These build fake transaction lists by hand - no
// network, no real block - so they run fast and safely as part of
// `make test`.
package sandwich

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/plasma6767/sandwich-detector/internal/ethnode"
)

func addr(hex string) common.Address {
	return common.HexToAddress(hex)
}

func addrPtr(hex string) *common.Address {
	a := addr(hex)
	return &a
}

// TestDetect_FindsSandwich checks the classic pattern: same sender trades
// the same pool twice, with a different sender's trade to that pool
// sandwiched in between.
func TestDetect_FindsSandwich(t *testing.T) {
	alice := addr("0x1")
	bob := addr("0x2")
	pool := addrPtr("0xaaaa")

	txs := []ethnode.TxSummary{
		{Position: 0, From: alice, To: pool, GasPrice: big.NewInt(30)},
		{Position: 1, From: bob, To: pool, GasPrice: big.NewInt(20)},
		{Position: 2, From: alice, To: pool, GasPrice: big.NewInt(10)},
	}

	got := Detect(txs)
	if len(got) != 1 {
		t.Fatalf("got %d sandwiches, want 1", len(got))
	}
	s := got[0]
	if s.FrontRun.Position != 0 || s.Victim.Position != 1 || s.BackRun.Position != 2 {
		t.Fatalf("got sandwich %+v, want front=0 victim=1 back=2", s)
	}
}

// TestDetect_NoFalsePositive_GasPriceNotBracketed checks that the same
// sender/recipient shape as a real sandwich is NOT flagged when the gas
// prices don't follow the attacker pattern (high on the first trade, lower
// on the victim's, lower still on the last).
func TestDetect_NoFalsePositive_GasPriceNotBracketed(t *testing.T) {
	alice := addr("0x1")
	bob := addr("0x2")
	pool := addrPtr("0xaaaa")

	txs := []ethnode.TxSummary{
		{Position: 0, From: alice, To: pool, GasPrice: big.NewInt(10)},
		{Position: 1, From: bob, To: pool, GasPrice: big.NewInt(20)},
		{Position: 2, From: alice, To: pool, GasPrice: big.NewInt(30)},
	}

	if got := Detect(txs); len(got) != 0 {
		t.Fatalf("got %d sandwiches, want 0 (gas prices don't bracket the victim)", len(got))
	}
}

// TestDetect_NoFalsePositive_DifferentSenders checks that transactions from
// three different senders to the same pool are never flagged - no sender
// repeats, so there's no attacker to have front-run and back-run anything.
func TestDetect_NoFalsePositive_DifferentSenders(t *testing.T) {
	pool := addrPtr("0xaaaa")

	txs := []ethnode.TxSummary{
		{Position: 0, From: addr("0x1"), To: pool},
		{Position: 1, From: addr("0x2"), To: pool},
		{Position: 2, From: addr("0x3"), To: pool},
	}

	if got := Detect(txs); len(got) != 0 {
		t.Fatalf("got %d sandwiches, want 0 (no repeated sender)", len(got))
	}
}

// TestDetect_NoFalsePositive_NoVictimBetween checks that the same sender
// trading the same pool twice in a row, with nothing in between, is not a
// sandwich - there's no victim to have profited off of.
func TestDetect_NoFalsePositive_NoVictimBetween(t *testing.T) {
	alice := addr("0x1")
	pool := addrPtr("0xaaaa")

	txs := []ethnode.TxSummary{
		{Position: 0, From: alice, To: pool},
		{Position: 1, From: alice, To: pool},
	}

	if got := Detect(txs); len(got) != 0 {
		t.Fatalf("got %d sandwiches, want 0 (no victim between same-sender txs)", len(got))
	}
}

// TestDetect_IgnoresContractCreation checks that a contract-creation
// transaction (no recipient) is never treated as a front-run, since there's
// no pool address for a victim or back-run to match against.
func TestDetect_IgnoresContractCreation(t *testing.T) {
	pool := addrPtr("0xaaaa")

	txs := []ethnode.TxSummary{
		{Position: 0, From: addr("0x1"), To: nil},
		{Position: 1, From: addr("0x2"), To: pool},
		{Position: 2, From: addr("0x1"), To: pool},
	}

	if got := Detect(txs); len(got) != 0 {
		t.Fatalf("got %d sandwiches, want 0 (front-run has no recipient)", len(got))
	}
}
