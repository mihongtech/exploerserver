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
	var TxIDMap = map[string]int{}
	for row.Next() {
		ticket := &resp.Ticket{}
		err := row.Scan(&ticket.TxID, &ticket.SpendTxID)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		TxIDMap[ticket.TxID] = 0
		if ticket.SpendTxID.Valid {
			TxIDMap[ticket.SpendTxID.String] = 0
		}
	}
	var transactionList []resp.Transaction
	for TxID := range TxIDMap {
		transaction := &resp.Transaction{TxID: TxID}
		err = getTx(db, transaction)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		transactionList = append(transactionList, *transaction)
	}
	return transactionList, nil
}
