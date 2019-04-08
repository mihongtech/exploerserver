package client

import (
	"encoding/json"
	"fmt"

	"github.com/linkchain/common/util/log"
	"github.com/linkchain/rpc/rpcobject"

	"github.com/linkchain-explorer/db"
	"github.com/linkchain-explorer/server/pool"
)

func Sync() {
	b, err := GetBestBlock()
	if err != nil {
		return
	}
	dbHeight := GetDBLastBlock()
	// DB height equal best height, not need sync
	if dbHeight == int(b.Height) {
		return
	}
	// when DB is empty, sync block from height 0
	if dbHeight < 0 {
		SyncBlockByHeight(0, b.Height)
	} else {
		forkHeight, err := getForkHeight(b.Height)
		if err != nil {
			return
		}
		// when DB height higher than fork height, remove db block which higher than fork height
		if int(forkHeight) < dbHeight {
			removeDBBlock(forkHeight)
		}
		SyncBlockByHeight(int(forkHeight+1), b.Height)
	}
}

func removeDBBlock(height uint32) {
	db := db.NewDB()
	defer db.Close()
	db.Exec("DELETE FROM blocks WHERE height > ?", height)
	db.Exec("DELETE FROM transactions WHERE height > ?", height)
	db.Exec("DELETE FROM tickets WHERE height > ?", height)
}

// get the latest block in DB which on the same fork with linkchain
func getForkHeight(height uint32) (uint32, error) {
	bs, err := pool.GetDBBlockSummaryByHeight(height)
	if err != nil {
		return 0, err
	}
	if bs == nil {
		if height > 0 {
			return getForkHeight(height - 1)
		} else {
			return 0, nil
		}
	}
	return bs.Height, nil
}

// get the highest block height
func GetDBLastBlock() int {
	db := db.NewDB()
	defer db.Close()

	row := db.QueryRow("SELECT height FROM blocks ORDER BY height DESC")
	height := -1
	row.Scan(&height)
	return height
}

// get best block
func GetBestBlock() (*rpcobject.BlockRSP, error) {
	s, err := callRpc("getBestBlock", nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var block = &rpcobject.BlockRSP{}
	err = json.Unmarshal(s, block)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return block, nil

}

// sync blockchain block info
func SyncBlockByHeight(height int, best uint32) {
	if height > int(best) {
		return
	}
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

	log.Info(fmt.Sprintf("Sync block success: block height: %d", height))
	tx.Commit()
	db.Close()

	SyncBlockByHeight(height+1, best)
}
