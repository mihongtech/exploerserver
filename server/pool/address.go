package pool

import (
	"github.com/linkchain-explorer/db"
	"github.com/linkchain-explorer/server/resp"
)

func getAddressInfo(params interface{}) (interface{}, error) {
	a, ok := params.(*AddressParams)
	if !ok {
		return nil, resp.BadRequestErr
	}
	db := db.NewDB()
	defer db.Close()
	row, err := db.Query("SELECT tx_id, spend_tx_id FROM tickets WHERE account_id = ?", a.Hash)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	var TxIDList []string
	for row.Next() {
		ticket := &resp.Ticket{}
		err := row.Scan(&ticket.TxID, &ticket.SpendTxID)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		if !contain(TxIDList, ticket.TxID) {
			TxIDList = append(TxIDList, ticket.TxID)
		}
		if ticket.SpendTxID.Valid && !contain(TxIDList, ticket.SpendTxID.String) {
			TxIDList = append(TxIDList, ticket.SpendTxID.String)
		}
	}
	var transactionList []resp.Transaction
	for _, TxID := range TxIDList {
		transaction := &resp.Transaction{TxID: TxID}
		err = getTx(db, transaction)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		transactionList = append(transactionList, *transaction)
	}
	return transactionList, nil
}

func contain(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
