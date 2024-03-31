package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"transaction-parser/internal/config"

	"transaction-parser/internal/data"
	"transaction-parser/internal/entity"
)

const (
	fetchTransactionTemplate = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":1}`
	fetchLatestBlock         = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}`
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) error
	GetTransactions(address string) ([]*entity.Transaction, error)
}

type EthereumParser struct {
	latestBlock  uint64
	currentBlock uint64
	rpc          string
	ctx          context.Context
	db           *data.DB
	mutex        sync.Mutex
}

func NewEthereumParser(ctx context.Context, cfg *config.Config) (*EthereumParser, error) {
	p := &EthereumParser{
		rpc:          cfg.RpcURL,
		ctx:          ctx,
		db:           data.NewDB(cfg),
		currentBlock: cfg.StartBlock,
	}
	err := p.Start(cfg.TickerPeriod)
	if err != nil {
		return nil, err

	}
	return p, nil
}

func (e *EthereumParser) GetCurrentBlock() int {
	return int(e.currentBlock)
}

func (e *EthereumParser) Subscribe(address string) error {
	return e.db.Address.Subscribe(address)
}

func (e *EthereumParser) GetTransactions(address string) ([]*entity.Transaction, error) {
	return e.db.Address.GetTransactions(address, 0)
}

func (e *EthereumParser) Start(period uint32) error {
	log.Printf("Starting the parser at block %d\n", e.currentBlock)
	latestBlock, err := e.getLatestBlock()
	if err != nil {
		return err
	}
	e.latestBlock = latestBlock
	if e.currentBlock == 0 {
		e.currentBlock = latestBlock
	}
	timer := time.NewTicker(time.Second * time.Duration(period))
	go func() {
		for {
			select {
			case <-timer.C:
				if e.currentBlock >= e.latestBlock {
					latestBlock, err := e.getLatestBlock()
					if err != nil {
						log.Println(err)
						continue
					}
					e.latestBlock = latestBlock
				}
				for e.currentBlock < e.latestBlock {
					err := e.parseCurrentBlock()
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			case <-e.ctx.Done():
				timer.Stop()
				return
			}
		}
	}()
	return nil
}

func (e *EthereumParser) parseCurrentBlock() error {
	block, err := e.getBlockInfo(e.currentBlock)
	if err != nil {
		return err
	}
	if block == nil {
		return errors.New("block is nil")
	}
	log.Printf("parsing block %d\n", e.currentBlock)
	for _, tx := range block.Transactions {
		e.db.Transaction.Create(tx)
	}
	e.currentBlock++
	return nil
}

func (e *EthereumParser) getBlockInfo(blockNumber uint64) (*entity.Block, error) {
	if blockNumber == 0 {
		return nil, errors.New("block number must be greater than 0")
	}
	body := &jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		ID:      1,
		Params:  []any{fmt.Sprintf("0x%x", blockNumber), true},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err

	}
	resp, err := http.Post(e.rpc, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var baseResponse jsonRPCResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	if err != nil {
		return nil, err
	}
	return baseResponse.Result, nil
}

func (e *EthereumParser) getLatestBlock() (uint64, error) {
	body := &jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		ID:      1,
		Params:  []any{"latest", false},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return 0, err

	}
	resp, err := http.Post(e.rpc, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return 0, err
	}
	var baseResponse jsonRPCLatestBlockResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	if err != nil {
		return 0, err
	}
	return hexStringToInt64(baseResponse.Result)
}

func hexStringToInt64(hexString string) (uint64, error) {
	hexString = strings.ToLower(hexString)
	if !strings.HasPrefix(hexString, "0x") {
		hexString = "0x" + hexString
	}
	intValue, err := strconv.ParseUint(hexString[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return intValue, nil
}
