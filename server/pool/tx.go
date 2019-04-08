package pool

import (
	"database/sql"
	"math/big"

	"github.com/linkchain-explorer/db"
	"github.com/linkchain-explorer/server/resp"
)

// get transaction by hash
func getTxByHash(params interface{}) (interface{}, error) {
	t, ok := params.(*TransactionHashParams)
	if !ok {
		return nil, resp.BadRequestErr
	}
	db := db.NewDB()
	defer db.Close()

	// query transaction info
	row := db.QueryRow("SELECT height, tx_id, version, type FROM transactions WHERE tx_id=?", t.Hash)
	transaction := &resp.Transaction{}
	err := row.Scan(&transaction.Block.Height, &transaction.TxID, &transaction.Version, &transaction.Type)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	err = getTx(db, transaction)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	err = getBlockSummary(db, transaction)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	return transaction, nil
}

func getTx(db *sql.DB, transaction *resp.Transaction) error {
	from, err := getTxFrom(db, transaction)
	transaction.From = from
	to, err := getTxTo(db, transaction)
	transaction.To = to
	return err
}

func getTxFrom(db *sql.DB, transaction *resp.Transaction) ([]resp.Ticket, error) {
	rows, err := db.Query("SELECT tx_id, account_id, amount, spend_tx_id FROM tickets WHERE spend_tx_id=?", transaction.TxID)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	defer rows.Close()
	var from []resp.Ticket
	for rows.Next() {
		f := resp.Ticket{}
		var amount int64
		err := rows.Scan(&f.TxID, &f.AccountID, &amount, &f.SpendTxID)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		f.Amount = big.NewInt(amount)
		from = append(from, f)
	}
	return from, nil
}

func getTxTo(db *sql.DB, transaction *resp.Transaction) ([]resp.Ticket, error) {
	rows, err := db.Query("SELECT tx_id, account_id, amount, spend_tx_id FROM tickets WHERE tx_id=?", transaction.TxID)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	defer rows.Close()
	var to []resp.Ticket
	for rows.Next() {
		t := resp.Ticket{}
		var amount int64
		err := rows.Scan(&t.TxID, &t.AccountID, &amount, &t.SpendTxID)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		t.Amount = big.NewInt(amount)
		to = append(to, t)
	}
	return to, nil
}

func getBlockSummary(db *sql.DB, transaction *resp.Transaction) error {
	row := db.QueryRow("SELECT hash, time FROM blocks WHERE height=?", transaction.Block.Height)
	err := row.Scan(&transaction.Block.Hash, &transaction.Block.Time)
	if err != nil {
		return resp.InternalServerErr
	}
	return nil
}
