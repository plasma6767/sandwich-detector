package ethnode

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxSummary is the minimal information about one transaction that
// sandwich detection needs: where it sat in the block, who sent it, and
// who it was sent to. To is nil for contract-creation transactions, which
// have no recipient.
type TxSummary struct {
	Position int
	From     common.Address
	To       *common.Address
}

// TransactionsIn returns an ordered summary of every transaction in block,
// in the same order they appear in the block itself - the ordering a
// sandwich attack depends on. signer supplies the chain-specific rules
// needed to recover each transaction's sender from its signature.
func TransactionsIn(block *types.Block, signer types.Signer) ([]TxSummary, error) {
	txs := block.Transactions()
	summaries := make([]TxSummary, 0, len(txs))

	for i, tx := range txs {
		// A transaction doesn't store its sender directly - only a
		// signature. types.Sender recovers the sender's address from that
		// signature using the chain's signing rules, the same check a
		// node performs to validate the transaction.
		from, err := types.Sender(signer, tx)
		if err != nil {
			return nil, fmt.Errorf("could not determine sender for transaction %d: %w", i, err)
		}

		// tx.To() is nil for contract-creation transactions, which have no
		// recipient - we pass that through as-is rather than substituting
		// a placeholder address.
		summaries = append(summaries, TxSummary{
			Position: i,
			From:     from,
			To:       tx.To(),
		})
	}

	return summaries, nil
}
