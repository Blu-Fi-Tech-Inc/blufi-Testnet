package network

import (
	"sync"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
	"github.com/blu-fi-tech-inc/boriqua_project/types"
)

// TxPool manages a pool of transactions.
type TxPool struct {
	all       *TxSortedMap // All transactions in the pool
	pending   *TxSortedMap // Transactions pending inclusion
	maxLength int          // Maximum length of the total pool
}

// NewTxPool creates a new transaction pool with a maximum length.
func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

// Add adds a transaction to the pool.
func (p *TxPool) Add(tx *core.Transaction) {
	// Prune the oldest transaction in the 'all' pool if it reaches max length.
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	// Add transaction to 'all' and 'pending' pools if it's not already present.
	if !p.all.Contains(tx.Hash(core.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

// Contains checks if a transaction hash exists in the pool.
func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

// Pending returns a slice of transactions in the pending pool.
func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.All()
}

// ClearPending clears all transactions from the pending pool.
func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

// PendingCount returns the count of transactions in the pending pool.
func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

// TxSortedMap is a sorted map implementation for transactions.
type TxSortedMap struct {
	lock   sync.RWMutex                // Mutex for concurrent access
	lookup map[types.Hash]*core.Transaction // Map for fast lookup by hash
	txx    *types.List                 // Sorted list of transactions
}

// NewTxSortedMap creates a new instance of TxSortedMap.
func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txx:    types.NewList(),
	}
}

// First returns the first transaction in the sorted map.
func (t *TxSortedMap) First() *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.txx.First()
	if first == nil {
		return nil
	}
	return t.lookup[first.Hash(core.TxHasher{})]
}

// Get retrieves a transaction by its hash.
func (t *TxSortedMap) Get(h types.Hash) *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

// Add adds a transaction to the sorted map.
func (t *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TxHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = tx
		t.txx.Insert(tx)
	}
}

// Remove removes a transaction from the sorted map by its hash.
func (t *TxSortedMap) Remove(h types.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	tx := t.lookup[h]
	if tx == nil {
		return
	}

	t.txx.Remove(tx)
	delete(t.lookup, h)
}

// Count returns the number of transactions in the sorted map.
func (t *TxSortedMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

// Contains checks if a transaction hash exists in the sorted map.
func (t *TxSortedMap) Contains(h types.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, ok := t.lookup[h]
	return ok
}

// All returns all transactions in the sorted map as a slice.
func (t *TxSortedMap) All() []*core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	var all []*core.Transaction
	for _, tx := range t.lookup {
		all = append(all, tx)
	}
	return all
}

// Clear removes all transactions from the sorted map.
func (t *TxSortedMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lookup = make(map[types.Hash]*core.Transaction)
	t.txx.Clear()
}
