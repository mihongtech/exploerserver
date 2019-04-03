package pool

type RequestHandler func(params interface{}) (interface{}, error)

var Handler = map[string]RequestHandler{
	"/block/best":       getBestBlock,
	"/block/hash":       getBlockByHash,
	"/block/list":       getBlockList,
	"/transaction/hash": getTxByHash,
	"/address/hash":     getAddressInfo,
}
