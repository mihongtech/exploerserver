package client

import (
	"encoding/json"

	"github.com/linkchain/common/util/log"
	"github.com/linkchain/rpc/rpcobject"
	"hum.tan/server/db"
	"hum.tan/server/server/pool"
)

// get the highest block height
func GetDBLastBlock() int {
	db := db.NewDB()
	defer db.Close()

	row := db.QueryRow("SELECT height FROM blocks ORDER BY height DESC")
	height := -1
	row.Scan(&height)
	return height
}

// sync blockchain block info
func SyncBlockByHeight(height int) {
	s, err := callRpc("getBlockByHeight", pool.BlockHeightParams{Height: height})
	if err != nil {
		log.Error(err.Error())
		return
	}
	var block = &rpcobject.BlockRSP{}
	err = json.Unmarshal(s, block)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// when block hash is empty stop sync
	if block.Hash == "" {
		return
	}

	db := db.NewDB()
	tx, _ := db.Begin()

	// save block
	_, err = tx.Exec("INSERT INTO blocks(hash, height, version, time, nonce, difficulty, prev, tx_root, status, sign, hex) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		block.Hash, block.Height, block.Header.Version, block.Header.Time, block.Header.Nonce, block.Header.Difficulty,
		block.Header.Prev.String(), block.Header.TxRoot.String(), block.Header.Status.String(), block.Header.Sign.String(), len(block.Hex)/2)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		db.Close()
		return
	}
	err = SaveTransactionDetail(tx, block.Height, block.TXIDs)
	if err != nil {
		log.Error(err.Error())
		tx.Rollback()
		db.Close()
		return
	}

	log.Info("Sync block success: block height: ", height)
	tx.Commit()
	db.Close()

	SyncBlockByHeight(height + 1)
}
