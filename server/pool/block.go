package pool

import (
	"database/sql"

	"github.com/linkchain/common/util/log"

	"github.com/linkchain-explorer/db"
	"github.com/linkchain-explorer/server/resp"
)

func GetDBBlockSummaryByHeight(height uint32) (*resp.BlockSummary, error) {
	db := db.NewDB()
	defer db.Close()
	row := db.QueryRow("SELECT height, hash, time FROM blocks WHERE height=?", height)
	block := &resp.BlockSummary{}
	err := row.Scan(&block.Height, &block.Hash, &block.Time)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, resp.NotFoundErr
	}
	return block, nil
}

func getBestBlock(params interface{}) (interface{}, error) {
	db := db.NewDB()
	defer db.Close()
	row := db.QueryRow("SELECT height, hash, version, time, nonce, difficulty, prev, tx_root, status, sign, hex FROM blocks ORDER BY height DESC")
	block := &resp.Block{}
	err := row.Scan(&block.Height, &block.Hash, &block.Version, &block.Time, &block.Nonce, &block.Difficulty, &block.Prev, &block.TxRoot, &block.Status, &block.Sign, &block.Hex)
	if err != nil {
		return nil, resp.NotFoundErr
	}
	return block, nil
}

// get block list
func getBlockList(params interface{}) (interface{}, error) {
	p, ok := params.(*BlockListParams)
	if !ok {
		log.Error("Params error")
		return nil, resp.BadRequestErr
	}
	db := db.NewDB()
	defer db.Close()

	query := db.QueryRow("SELECT COUNT(*) FROM blocks")
	var count int
	err := query.Scan(&count)
	if err != nil {
		log.Error(err.Error())
		return nil, resp.InternalServerErr
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Size == 0 {
		p.Size = 10
	}

	if count/p.Size+1 < p.Page {
		p.Page = count / p.Size
	}

	rows, err := db.Query("SELECT hash, height, version, time, nonce, difficulty, prev, tx_root, status, sign, hex FROM blocks order by height DESC LIMIT ?, ?",
		(p.Page-1)*p.Size, p.Size)
	if err != nil {
		log.Error(err.Error())
		return nil, resp.InternalServerErr
	}
	defer rows.Close()

	var res resp.ListResp
	for rows.Next() {
		block := resp.Block{}
		err := rows.Scan(&block.Hash, &block.Height, &block.Version, &block.Time, &block.Nonce, &block.Difficulty, &block.Prev, &block.TxRoot, &block.Status, &block.Sign, &block.Hex)
		if err != nil {
			log.Error(err.Error())
			return nil, resp.InternalServerErr
		}

		// query transaction num
		query := db.QueryRow("SELECT COUNT(*) FROM transactions WHERE height=?", block.Height)
		var txs int
		err = query.Scan(&txs)
		if err != nil {
			log.Error(err.Error())
			return nil, resp.InternalServerErr
		}
		block.TXs = txs

		res.List = append(res.List, block)
	}
	res.Page = p.Page
	res.Size = p.Size
	res.Total = count
	return res, nil
}

// get block info by block hash
func getBlockByHash(params interface{}) (interface{}, error) {
	p, ok := params.(*BlockHashParams)
	if !ok {
		return nil, resp.BadRequestErr
	}
	db := db.NewDB()
	defer db.Close()

	// query block info
	row := db.QueryRow("SELECT height, hash, version, time, nonce, difficulty, prev, tx_root, status, sign, hex FROM blocks WHERE hash=?", p.Hash)
	block := &resp.Block{}
	err := row.Scan(&block.Height, &block.Hash, &block.Version, &block.Time, &block.Nonce, &block.Difficulty, &block.Prev, &block.TxRoot, &block.Status, &block.Sign, &block.Hex)
	if err != nil {
		return nil, resp.NotFoundErr
	}

	// query next block hash
	row = db.QueryRow("SELECT hash FROM blocks WHERE height=?", block.Height+1)
	var next string
	row.Scan(&next)
	block.Next = next

	// query block transactions
	rows, err := db.Query("SELECT tx_id, version, type FROM transactions WHERE height=?", block.Height)
	if err != nil {
		return nil, resp.InternalServerErr
	}
	defer rows.Close()

	for rows.Next() {
		var transaction = &resp.Transaction{}
		err := rows.Scan(&transaction.TxID, &transaction.Version, &transaction.Type)
		if err != nil {
			return nil, resp.InternalServerErr
		}
		err = getTx(db, transaction)
		if err != nil {
			return nil, err
		}
		block.Transactions = append(block.Transactions, *transaction)
	}
	return block, nil
}
