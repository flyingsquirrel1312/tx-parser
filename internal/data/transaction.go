package data

import (
	"fmt"
	"sync"
	"transaction-parser/internal/entity"
	"transaction-parser/pkg/lib"
)

type TransactionStorage interface {
	Create(transaction *entity.Transaction)
}

type InMemoryTransactionStorage struct {
	db       *inMemoryDB
	mutex    sync.RWMutex
	capacity uint32
}

func (s *InMemoryTransactionStorage) Create(transaction *entity.Transaction) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// only store transactions that involve subscribed addresses
	if db.subscribedAddresses.Has(transaction.From) || db.subscribedAddresses.Has(transaction.To) {
		s.db.tnxMap.Set(transaction.Hash, transaction)
		fmt.Println(transaction.From, transaction.To)
		s.writeToAddressMap(transaction.From, transaction.Hash)
		s.writeToAddressMap(transaction.To, transaction.Hash)
	}
}

func (s *InMemoryTransactionStorage) writeToAddressMap(address string, txHash string) {
	tnxs, ok := s.db.addressToTxMap.Get(address)
	if !ok {
		tnxs = lib.NewCircularBuffer[string](s.capacity)
	}
	tnxs.Enqueue(txHash)
	s.db.addressToTxMap.Set(address, tnxs)
}

func newInMemoryStorage(db *inMemoryDB, capacity uint32) *InMemoryTransactionStorage {
	return &InMemoryTransactionStorage{
		db:       db,
		mutex:    sync.RWMutex{},
		capacity: capacity,
	}
}
