package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/linkchain/core/meta"
	"github.com/linkchain/rpc/rpcobject"

	"github.com/linkchain-explorer/server/pool"
)

func SaveTransactionDetail(tx *sql.Tx, height uint32, transactionHashList []string) error {
	for _, transactionHash := range transactionHashList {
		err := SaveTransaction(tx, height, transactionHash)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveTransaction(tx *sql.Tx, height uint32, transactionHash string) error {
	s, err := callRpc("getTxByHash", pool.TransactionHashParams{Hash: transactionHash})
	if err != nil {
		return err
	}
	transaction := &rpcobject.TransactionWithIDRSP{}
	err = json.Unmarshal(s, transaction)
	if err != nil {
		return err
	}
	if transaction.Tx != nil {
		err := SaveTransactionInfo(tx, height, transaction)
		if err != nil {
			return err
		}
		if transaction.Tx.From.Coins != nil {
			err := SaveTransactionFrom(tx, transaction.Tx.From.Coins, transactionHash)
			if err != nil {
				return err
			}
		}
		if transaction.Tx.To.Coins != nil {
			err := SaveTransactionTo(tx, height, transaction.Tx.To.Coins, transactionHash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SaveTransactionInfo(tx *sql.Tx, height uint32, transaction *rpcobject.TransactionWithIDRSP) error {
	_, err := tx.Exec("INSERT INTO transactions(height, tx_id, version, type) VALUES (?, ?, ?, ?)",
		height, transaction.ID, transaction.Tx.Version, transaction.Tx.Type)
	if err != nil {
		return err
	}
	return nil
}

func SaveTransactionFrom(tx *sql.Tx, coins []meta.FromCoin, transactionHash string) error {
	for _, coin := range coins {
		for _, ticket := range coin.Ticket {
			_, err := tx.Exec("UPDATE tickets SET spend_tx_id=? WHERE tx_id=? AND `index`=?", transactionHash, ticket.Txid.String(), ticket.Index)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SaveTransactionTo(tx *sql.Tx, height uint32, coins []meta.ToCoin, transactionHash string) error {
	var values []string
	for i, coin := range coins {
		values = append(values, fmt.Sprintf("(%d, '%s', '%s', %d, %d)",
			height, transactionHash, coin.Id, i, coin.Value.GetInt64()))
	}
	_, err := tx.Exec("INSERT INTO tickets(height, tx_id, account_id, `index`, amount) VALUES " + strings.Join(values, ","))
	if err != nil {
		return err
	}
	return nil
}
