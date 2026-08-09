// Package sandwich contains the detection heuristic itself: given an
// ordered list of transactions from one block, find the sandwich-attack
// pattern among them.
package sandwich

import "github.com/plasma6767/sandwich-detector/internal/ethnode"

// Sandwich is one detected attack: a victim transaction sitting between two
// transactions from the same attacker, all targeting the same recipient
// (typically a trading pool contract).
type Sandwich struct {
	FrontRun ethnode.TxSummary
	Victim   ethnode.TxSummary
	BackRun  ethnode.TxSummary
}

// Detect scans an ordered list of transactions for the sandwich pattern:
// the same sender appearing twice against the same recipient, with a
// different sender's transaction to that same recipient sandwiched in
// between. It returns every match found, in block order.
func Detect(txs []ethnode.TxSummary) []Sandwich {
	var found []Sandwich

	for i, front := range txs {
		// A contract-creation transaction has no recipient, so it can't be
		// the front-run half of a sandwich targeting a pool.
		if front.To == nil {
			continue
		}

		// Only check the nearest transaction with a matching sender and
		// recipient - if that pair has no victim between them, treating it
		// as a possible sandwich isn't useful, so we move on to the next
		// front-run candidate rather than searching further ahead.
		for k := i + 1; k < len(txs); k++ {
			back := txs[k]
			if back.To == nil || back.From != front.From || *back.To != *front.To {
				continue
			}

			for j := i + 1; j < k; j++ {
				victim := txs[j]
				if victim.To != nil && *victim.To == *front.To && victim.From != front.From {
					found = append(found, Sandwich{FrontRun: front, Victim: victim, BackRun: back})
					break
				}
			}
			break
		}
	}

	return found
}
