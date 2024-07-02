package core

import (
	"errors"
	"fmt"
)

// ErrBlockKnown is returned when a block is already known to the blockchain.
var ErrBlockKnown = errors.New("block already known")

// Validator is an interface that defines the ValidateBlock method.
type Validator interface {
	ValidateBlock(*Block) error
}

// BlockValidator is responsible for validating blocks against the blockchain.
type BlockValidator struct {
	bc *Blockchain
}

// NewBlockValidator creates a new BlockValidator for a given blockchain.
func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

// ValidateBlock validates a block to ensure it can be added to the blockchain.
func (v *BlockValidator) ValidateBlock(b *Block) error {
	// Check if the block is already known to the blockchain.
	if v.bc.HasBlock(b.Height) {
		return ErrBlockKnown
	}

	// Ensure the block height is the next expected height in the blockchain.
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf(
			"block (%s) with height (%d) is too high => current height (%d)",
			b.Hash(BlockHasher{}),
			b.Height,
			v.bc.Height(),
		)
	}

	// Retrieve the previous block header and verify it.
	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	// Verify the previous block hash matches the current block's PrevBlockHash.
	if hash := BlockHasher{}.Hash(prevHeader); hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	// Verify the integrity of the current block.
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
