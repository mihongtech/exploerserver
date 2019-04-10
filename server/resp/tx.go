package resp

import (
	"database/sql"
	"math/big"
)

type Transaction struct {
	TxID    string       `json:"tx_id"`   // The hash of the Transaction.
	Type    uint32       `json:"type"`    // The type of the Transaction.
	Version uint32       `json:"version"` // The version of the Transaction.  This is not the same as the Blocks version.
	From    []Ticket     `json:"from"`    // The accounts of the Transaction related to inputs.
	To      []Ticket     `json:"to"`      // The accounts of the Transaction related to outputs.
	Block   BlockSummary `json:"block"`   // The summary of the block
}

type Ticket struct {
	AccountID string         `json:"account_id"`  // AccountID of the Ticket.
	Amount    *big.Int       `json:"amount"`      // Amount of the Ticket.
	TxID      string         `json:"tx_id"`       // Transaction Hash of the Ticket generate.
	SpendTxID sql.NullString `json:"spend_tx_id"` // Transaction Hash of the Ticket spend.
	Index     uint32         `json:"index"`
}
