package config

type Config struct {
	RpcURL       string `json:"rpc_url"`
	Port         int    `json:"port"`
	BufferSize   uint32 `json:"buffer_size"`
	TickerPeriod uint32 `json:"ticker_period"`
	StartBlock   uint64 `json:"start_block"`
}
