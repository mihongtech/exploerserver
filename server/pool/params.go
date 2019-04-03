package pool

import (
	"reflect"
)

var Params = map[string]reflect.Type{
	"/block/list":               reflect.TypeOf((*BlockListParams)(nil)),
	"/block/hash":               reflect.TypeOf((*BlockHashParams)(nil)),
	"/transaction/hash":         reflect.TypeOf((*TransactionHashParams)(nil)),
	"/address/hash":             reflect.TypeOf((*AddressParams)(nil)),
	"/wallet/import/privateKey": reflect.TypeOf((*ImportAccountParams)(nil)),
}

// block
type BlockHeightParams struct {
	Height int `json:"height"`
}

type BlockHashParams struct {
	Hash string `json:"hash"`
}

type BlockListParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// transaction
type TransactionHashParams struct {
	Hash string `json:"hash"`
}

// wallet
type ImportAccountParams struct {
	Signer string `json:"accountPrivateKey"`
}

// address
type AddressParams struct {
	Hash string `json:"hash"`
}
