package service

import "transaction-parser/internal/entity"

type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type jsonRPCResponse struct {
	JsonRPC string        `json:"jsonrpc"`
	Result  *entity.Block `json:"result"`
}

type jsonRPCLatestBlockResponse struct {
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}
