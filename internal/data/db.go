package data

import (
	"transaction-parser/internal/config"
	"transaction-parser/internal/entity"
	"transaction-parser/pkg/lib"
)

var (
	db *inMemoryDB
)

type inMemoryDB struct {
	subscribedAddresses *lib.ConcurrentSet[string]                          // list of subscribed addresses
	tnxMap              *lib.FIFOCache[string, *entity.Transaction]         // key: transaction hash, value: transaction
	addressToTxMap      *lib.FIFOCache[string, *lib.CircularBuffer[string]] // key: address, value: list of transaction hashes
}

type DB struct {
	Address     AddressStorage
	Transaction TransactionStorage
}

func NewDB(cfg *config.Config) *DB {
	db = &inMemoryDB{
		tnxMap:              lib.NewFIFOCache[string, *entity.Transaction](cfg.BufferSize),
		subscribedAddresses: lib.NewConcurrentSet[string](cfg.BufferSize),
		addressToTxMap:      lib.NewFIFOCache[string, *lib.CircularBuffer[string]](cfg.BufferSize),
	}
	return &DB{
		Address:     newInMemoryAddressStorage(db),
		Transaction: newInMemoryStorage(db, cfg.BufferSize),
	}
}
