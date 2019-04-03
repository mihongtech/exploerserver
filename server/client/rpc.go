package client

import (
	"encoding/json"
)

var httpConfig = &Config{
	RPCUser:     "lc",
	RPCPassword: "lc",
	RPCServer:   "localhost:8083",
}

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

func NewRequest(id interface{}, method string, param interface{}) (*Request, error) {
	var rawParams json.RawMessage
	if param != nil {
		marshalledParam, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}
		rawParams = json.RawMessage(marshalledParam)
	}

	return &Request{
		Jsonrpc: "1.0",
		ID:      id,
		Method:  method,
		Params:  rawParams,
	}, nil
}

////rpc call
func callRpc(method string, params interface{}) ([]byte, error) {
	//param
	rawParams, err := NewRequest(1, method, params)
	if err != nil {
		//log.Error(method, "resp", resp)
		return nil, err
	}
	marshaledParams, err := json.Marshal(rawParams)
	if err != nil {
		//log.Error(method, "resp", resp)
		return nil, err
	}
	//response
	rawRet, err := SendPostRequest(marshaledParams, httpConfig)
	if err != nil {
		//log.Error(method, "resp", resp)
		return nil, err
	}

	return rawRet, nil
}
