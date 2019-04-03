package resp

import "time"

type Block struct {
	Version      uint32        `json:"version,int"` // Version of the block.  This is not the same as the protocol version.
	Height       uint32        `json:"height,int"`  //the height of block.
	Time         time.Time     `json:"time"`        // Time the block was created.  This is, unfortunately, encoded as a uint32 on the wire and therefore is limited to 2106.
	Nonce        uint32        `json:"nonce"`       // Nonce used to generate the block.
	Difficulty   uint32        `json:"difficulty"`  // Difficulty target for the block.
	Prev         string        `json:"prev"`        // Hash of the previous block header in the block chain.
	Next         string        `json:"next"`        // Hash of the next block header in the block chain.
	TxRoot       string        `json:"txroot"`      // Merkle tree reference to hash of all transactions for the block.
	Status       string        `json:"status"`      // The status of the whole system.
	Sign         string        `json:"sign"`        // The sign of miner.
	Hash         string        `json:"hash"`        //The Hash of this block.
	Hex          int64         `json:"hex"`         //The Hex of this block.
	TXs          int           `json:"txs"`         //The Transaction numbers of this block.
	Transactions []Transaction `json:"transactions"`
}

type BlockSummary struct {
	Height uint32    `json:"height,int"` //the height of block.
	Time   time.Time `json:"time"`       // Time the block was created.  This is, unfortunately, encoded as a uint32 on the wire and therefore is limited to 2106.
	Hash   string    `json:"hash"`       //The Hash of this block.
}
