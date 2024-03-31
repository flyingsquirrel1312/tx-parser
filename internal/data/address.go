package data

import (
	"errors"
	"transaction-parser/internal/entity"
)

var (
	ErrAddressNotSubscribed     = errors.New("address not subscribed")
	ErrCapacityExceeded         = errors.New("capacity exceeded")
	ErrAddressAlreadySubscribed = errors.New("address already subscribed")
)

type AddressStorage interface {
	// Subscribe adds an address to the list of subscribed addresses.
	// Returns true if the address was successfully added, false otherwise.
	Subscribe(address string) error
	// IsSubscribed checks if an address is in the list of subscribed addresses.
	//Returns true if the address is in the list, false otherwise.

	GetTransactions(address string, n uint32) ([]*entity.Transaction, error)
}

type InMemoryAddressStorage struct {
	db *inMemoryDB
}

func (s *InMemoryAddressStorage) Subscribe(address string) error {
	if ok := s.db.subscribedAddresses.Has(address); ok {
		return ErrAddressAlreadySubscribed
	}
	if ok := s.db.subscribedAddresses.Add(address); !ok {
		return ErrCapacityExceeded
	}
	return nil
}

func (s *InMemoryAddressStorage) GetTransactions(address string, n uint32) ([]*entity.Transaction, error) {
	if ok := s.db.subscribedAddresses.Has(address); !ok {
		return nil, ErrAddressNotSubscribed
	}
	var transactions []*entity.Transaction
	txHashes, ok := s.db.addressToTxMap.Get(address)
	if !ok {
		return transactions, nil
	}
	if n > txHashes.Len() || n < 1 {
		n = txHashes.Len()
	}
	for _, txHash := range txHashes.LastN(n) {
		tx, ok := s.db.tnxMap.Get(txHash)
		if !ok {
			continue
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func newInMemoryAddressStorage(db *inMemoryDB) *InMemoryAddressStorage {
	return &InMemoryAddressStorage{
		db: db,
	}
}
