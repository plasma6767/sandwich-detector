// Tests for TransactionsIn(). Like client_test.go, these build all their
// input by hand in memory - no network call - so they run fast and safely
// as part of `make test`.
package ethnode

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
)

// TestTransactionsIn builds a fake block containing three transactions -
// two ordinary transfers and one contract-creation transaction (which has
// no recipient) - and checks that TransactionsIn reports each one's sender,
// recipient, and position correctly.
func TestTransactionsIn(t *testing.T) {
	// A real transaction's sender isn't stored directly - it's derived by
	// recovering the signing key from the transaction's signature. So to
	// build a believable fake transaction, we generate real throwaway
	// private keys and actually sign with them, the same way a wallet would.
	aliceKey := generateKey(t)
	bobKey := generateKey(t)
	carolKey := generateKey(t)
	poolAddr := crypto.PubkeyToAddress(generateKey(t).PublicKey)

	aliceAddr := crypto.PubkeyToAddress(aliceKey.PublicKey)
	bobAddr := crypto.PubkeyToAddress(bobKey.PublicKey)
	carolAddr := crypto.PubkeyToAddress(carolKey.PublicKey)

	// LatestSignerForChainID gives us the signing rules used to sign (and
	// later recover the sender from) transactions for a given chain ID.
	// Chain ID 1 is Ethereum mainnet's ID; the exact number doesn't matter
	// for this test as long as signing and recovery use the same one.
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// tx0 and tx1 are ordinary transfers with a recipient. tx2 has a nil
	// To address, which is how Ethereum represents "this transaction
	// deploys a new contract" rather than sending to an existing address.
	// Each uses a different gas price (the fee offered to get included) so
	// the test can confirm TransactionsIn reports the right one per
	// transaction, not just a shared default.
	tx0 := signTx(t, signer, aliceKey, &poolAddr, big.NewInt(10))
	tx1 := signTx(t, signer, bobKey, &poolAddr, big.NewInt(20))
	tx2 := signTx(t, signer, carolKey, nil, big.NewInt(30))

	block := types.NewBlock(
		&types.Header{Number: big.NewInt(1)},
		&types.Body{Transactions: []*types.Transaction{tx0, tx1, tx2}},
		nil,
		trie.NewStackTrie(nil),
	)

	got, err := TransactionsIn(block, signer)
	if err != nil {
		t.Fatalf("TransactionsIn() failed: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d transactions, want 3", len(got))
	}

	checkSummary(t, got[0], 0, aliceAddr, &poolAddr, big.NewInt(10))
	checkSummary(t, got[1], 1, bobAddr, &poolAddr, big.NewInt(20))
	checkSummary(t, got[2], 2, carolAddr, nil, big.NewInt(30))
}

// TestTransactionsIn_WrongSigner checks that TransactionsIn returns an
// error (rather than a wrong answer or a panic) when asked to recover
// senders using a signer for the wrong chain ID.
func TestTransactionsIn_WrongSigner(t *testing.T) {
	key := generateKey(t)
	to := crypto.PubkeyToAddress(generateKey(t).PublicKey)

	signingSigner := types.LatestSignerForChainID(big.NewInt(1))
	tx := signTx(t, signingSigner, key, &to, big.NewInt(1))

	block := types.NewBlock(
		&types.Header{Number: big.NewInt(1)},
		&types.Body{Transactions: []*types.Transaction{tx}},
		nil,
		trie.NewStackTrie(nil),
	)

	wrongSigner := types.LatestSignerForChainID(big.NewInt(5))
	if _, err := TransactionsIn(block, wrongSigner); err == nil {
		t.Fatal("expected an error recovering sender with the wrong signer, got nil")
	}
}

// generateKey creates a throwaway private key for use as a fake
// transaction sender in tests.
func generateKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return key
}

// signTx builds and signs a minimal transaction from the given key to the
// given recipient (nil meaning "contract creation, no recipient") offering
// the given gas price - the fee the sender pays to get included.
func signTx(t *testing.T, signer types.Signer, key *ecdsa.PrivateKey, to *common.Address, gasPrice *big.Int) *types.Transaction {
	t.Helper()
	tx, err := types.SignNewTx(key, signer, &types.LegacyTx{
		Nonce:    0,
		To:       to,
		Value:    big.NewInt(0),
		Gas:      21000,
		GasPrice: gasPrice,
	})
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}
	return tx
}

// checkSummary asserts that a TxSummary matches the expected position,
// sender, recipient (a nil wantTo means "expected no recipient"), and gas
// price.
func checkSummary(t *testing.T, got TxSummary, wantPosition int, wantFrom common.Address, wantTo *common.Address, wantGasPrice *big.Int) {
	t.Helper()
	if got.Position != wantPosition {
		t.Errorf("got position %d, want %d", got.Position, wantPosition)
	}
	if got.From != wantFrom {
		t.Errorf("got from %s, want %s", got.From, wantFrom)
	}
	switch {
	case wantTo == nil && got.To != nil:
		t.Errorf("got to %s, want nil (contract creation)", got.To)
	case wantTo != nil && got.To == nil:
		t.Errorf("got to nil, want %s", wantTo)
	case wantTo != nil && got.To != nil && *got.To != *wantTo:
		t.Errorf("got to %s, want %s", got.To, wantTo)
	}
	if got.GasPrice == nil || got.GasPrice.Cmp(wantGasPrice) != 0 {
		t.Errorf("got gas price %s, want %s", got.GasPrice, wantGasPrice)
	}
}
