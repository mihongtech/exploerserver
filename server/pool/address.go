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

	row, err := db.Query("SELECT tx_id, spend_tx_id, amount FROM tickets WHERE account_id = ?", a.Hash)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	var TxIDList []string
	var final int64
	for row.Next() {
		ticket := &resp.Ticket{}
		var amount int64
		err := row.Scan(&ticket.TxID, &ticket.SpendTxID, &amount)

		if err != nil {
			return nil, resp.InternalServerErr
		}
		if !ticket.SpendTxID.Valid {
			final += amount
		}
		if !contain(TxIDList, ticket.TxID) {
			TxIDList = append(TxIDList, ticket.TxID)
		}
		if ticket.SpendTxID.Valid && !contain(TxIDList, ticket.SpendTxID.String) {
			TxIDList = append(TxIDList, ticket.SpendTxID.String)
		}
	}

	if a.Page < 1 {
		a.Page = 1
	}

	if a.Size == 0 {
		a.Size = 10
	}

	count := len(TxIDList)
	if count/a.Size+1 < a.Page {
		a.Page = count / a.Size
	}

	addressResp := resp.Address{
		Page:  a.Page,
		Size:  a.Size,
		Total: count,
		Final: final,
	}
	var pagerTxIDList []string
	if a.Page*a.Size < count {
		pagerTxIDList = TxIDList[(a.Page-1)*a.Size : a.Page*a.Size]
	} else {
		pagerTxIDList = TxIDList[(a.Page-1)*a.Size:]
	}
	for _, TxID := range pagerTxIDList {
		transaction := &resp.Transaction{TxID: TxID}
		err = getTx(db, transaction)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		addressResp.List = append(addressResp.List, *transaction)
	}
	return addressResp, nil
}

func contain(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
