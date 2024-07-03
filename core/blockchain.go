package core

import (
	"fmt"
	"sync"

	"github.com/blu-fi-tech-inc/blufi-network/types"
	"github.com/go-kit/log"
)

type Blockchain struct {
	logger          log.Logger
	store           Store
	lock            sync.RWMutex
	headers         []*Header
	blocks          []*Block
	txStore         map[types.Hash]*Transaction
	blockStore      map[types.Hash]*Block

	accountState    *AccountState

	stateLock       sync.RWMutex
	collectionState map[types.Hash]*CollectionTx
	mintState       map[types.Hash]*MintTx
	validator       Validator
	contractState   *State
}

// NewBlockchain creates a new Blockchain instance.
func NewBlockchain(store Store, l log.Logger, accountState *AccountState, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		store:           store,
		logger:          l,
		accountState:    accountState,
		collectionState: make(map[types.Hash]*CollectionTx),
		mintState:       make(map[types.Hash]*MintTx),
		blockStore:      make(map[types.Hash]*Block),
		txStore:         make(map[types.Hash]*Transaction),
		contractState:   NewState(),
		headers:         []*Header{},
		blocks:          []*Block{},
	}
	bc.validator = NewBlockValidator(bc)

	if err := bc.addBlockWithoutValidation(genesis); err != nil {
		return nil, err
	}
	return bc, nil
}

// SetValidator sets the block validator.
func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

// AddBlock adds a block to the blockchain after validation.
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

// handleNativeTransfer processes native token transfers.
func (bc *Blockchain) handleNativeTransfer(tx *Transaction) error {
	bc.logger.Log(
		"msg", "handle native token transfer",
		"from", tx.From,
		"to", tx.To,
		"value", tx.Value,
	)

	fromAddr, err := tx.From.Address()
	if err != nil {
		return err
	}
	toAddr, err := tx.To.Address()
	if err != nil {
		return err
	}

	return bc.accountState.Transfer(fromAddr, toAddr, tx.Value)
}

// handleNativeNFT processes native NFT transactions.
func (bc *Blockchain) handleNativeNFT(tx *Transaction) error {
	hash := tx.Hash(TxHasher{})

	switch t := tx.TxInner.(type) {
	case CollectionTx:
		bc.collectionState[hash] = &t
		bc.logger.Log("msg", "created new NFT collection", "hash", hash)
	case MintTx:
		_, ok := bc.collectionState[t.Collection];
		if !ok {
			return fmt.Errorf("collection (%s) does not exist on the blockchain", t.Collection)
		}
		bc.mintState[hash] = &t

		bc.logger.Log("msg", "created new NFT mint", "NFT", t.NFT, "collection", t.Collection)
	default:
		return fmt.Errorf("unsupported tx type %T", t)
	}

	return nil
}

// GetBlockByHash retrieves a block by its hash.
func (bc *Blockchain) GetBlockByHash(hash types.Hash) (*Block, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	block, ok := bc.blockStore[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash (%s) not found", hash)
	}

	return block, nil
}

// GetBlock retrieves a block by its height.
func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return bc.blocks[height], nil
}

// GetHeader retrieves a block header by its height.
func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return bc.headers[height], nil
}

// GetTxByHash retrieves a transaction by its hash.
func (bc *Blockchain) GetTxByHash(hash types.Hash) (*Transaction, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	tx, ok := bc.txStore[hash]
	if !ok {
		return nil, fmt.Errorf("could not find tx with hash (%s)", hash)
	}

	return tx, nil
}

// HasBlock checks if a block exists at a given height.
func (bc *Blockchain) HasBlock(height uint32) bool {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return height <= uint32(len(bc.headers)-1)
}

// Height returns the current height of the blockchain.
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

// handleTransaction processes a transaction.
func (bc *Blockchain) handleTransaction(tx *Transaction) error {
	if len(tx.Data) > 0 {
		bc.logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.Hash(TxHasher{}))

		vm := NewVM(tx.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
	}

	if tx.TxInner != nil {
		switch tx.TxInner.(type) {
		case *CollectionTx, *MintTx:
			if err := bc.handleNativeNFT(tx); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported tx type %T", tx.TxInner)
		}
	}

	if tx.Value > 0 {
		if err := bc.handleNativeTransfer(tx); err != nil {
			return err
		}
	}

	return nil
}

// addBlockWithoutValidation adds a block to the blockchain without validation.
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.stateLock.Lock()
	defer bc.stateLock.Unlock()

	for _, tx := range b.Transactions {
		if err := bc.handleTransaction(tx); err != nil {
			bc.logger.Log("error", err.Error())
			continue
		}
	}

	defer bc.lock.Unlock()

	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, b)
	bc.blockStore[b.Hash(BlockHasher{})] = b

	for _, tx := range b.Transactions {
		bc.txStore[tx.Hash(TxHasher{})] = tx
	}
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
