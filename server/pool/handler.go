package pool

type RequestHandler func(params interface{}) (interface{}, error)

var Handler = map[string]RequestHandler{
	"/rpc/block/best":       getBestBlock,
	"/rpc/block/hash":       getBlockByHash,
	"/rpc/block/list":       getBlockList,
	"/rpc/transaction/hash": getTxByHash,
	"/rpc/address/hash":     getAddressInfo,
	"/rpc/search/global":    globalSearch,
}
