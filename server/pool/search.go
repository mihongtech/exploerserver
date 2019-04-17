package pool

import (
	"database/sql"
	"strconv"

	"github.com/mihongtech/linkchain-explorer/db"
	"github.com/mihongtech/linkchain-explorer/server/resp"
)

func globalSearch(params interface{}) (interface{}, error) {
	p, ok := params.(*GlobalSearchParams)
	if !ok {
		return nil, resp.BadRequestErr
	}
	height, err := strconv.Atoi(p.Keyword)
	if err != nil {
		return searchByHash(p.Keyword)
	} else {
		return searchBlock(height)
	}
}

func searchBlock(height int) (interface{}, error) {
	db := db.NewDB()
	defer db.Close()

	row := db.QueryRow("SELECT hash FROM blocks WHERE height=?", height)
	var hash string
	err := row.Scan(&hash)
	if err == sql.ErrNoRows {
		return map[string]string{"path": "null"}, nil
	}
	if err != nil {
		return nil, resp.InternalServerErr
	}
	return map[string]string{
		"path":  "block",
		"param": hash,
	}, nil
}

func searchByHash(keyword string) (interface{}, error) {
	db := db.NewDB()
	defer db.Close()
	return searchBlockByHash(db, keyword)
}

// search block
func searchBlockByHash(db *sql.DB, keyword string) (interface{}, error) {
	row := db.QueryRow("SELECT hash FROM blocks WHERE hash=?", keyword)
	var hash string
	err := row.Scan(&hash)
	if err == sql.ErrNoRows {
		return searchTransactionByHash(db, keyword)
	}
	if err != nil {
		return nil, resp.InternalServerErr
	}
	return map[string]string{
		"path":  "block",
		"param": hash,
	}, nil
}

func searchTransactionByHash(db *sql.DB, keyword string) (interface{}, error) {
	row := db.QueryRow("SELECT tx_id FROM transactions WHERE tx_id=?", keyword)
	var transaction string
	err := row.Scan(&transaction)
	if err == sql.ErrNoRows {
		return searchTransactionByAddress(db, keyword)
	}
	if err != nil {
		return nil, resp.InternalServerErr
	}
	return map[string]string{
		"path":  "transaction",
		"param": transaction,
	}, nil
}

func searchTransactionByAddress(db *sql.DB, keyword string) (interface{}, error) {
	row := db.QueryRow("SELECT account_id FROM tickets WHERE account_id=?", keyword)
	var address string
	err := row.Scan(&address)
	if err == sql.ErrNoRows {
		return map[string]string{"path": "null"}, nil
	}
	if err != nil {
		return nil, resp.InternalServerErr
	}
	return map[string]string{
		"path":  "address",
		"param": address,
	}, nil
}
