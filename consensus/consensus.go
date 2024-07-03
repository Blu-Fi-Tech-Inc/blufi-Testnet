package consensus

import (
	"github.com/blu-fi-tech-inc/blufi-network/types"
	"github.com/blu-fi-tech-inc/blufi-network/core"
	"github.com/blu-fi-tech-inc/blufi-network/crypto"
	"errors"
)

// Consensus represents the consensus mechanism
type Consensus struct {
	stakeManager *StakeManager
	pos          *PoS
}

// NewConsensus initializes a new consensus mechanism
func NewConsensus(sm *StakeManager, pos *PoS) *Consensus {
	return &Consensus{
		stakeManager: sm,
		pos:          pos,
	}
}

// ValidateBlock checks if a block is valid
func (c *Consensus) ValidateBlock(block *types.Block) bool {
	validators, err := c.pos.SelectValidators(10) // Assume 10 validators
	if err != nil {
		return false
	}

	for _, validator := range validators {
		if validator.Address == block.Proposer {
			return true
		}
	}
	return false
}

// AddBlock adds a block to the blockchain if valid
func (c *Consensus) AddBlock(block *types.Block) bool {
	if c.ValidateBlock(block) {
		// Add block to the blockchain
		if err := core.GetBlockchain().AddBlock(block); err != nil {
			return false
		}
		return true
	}
	return false
}

// PoS represents the Proof of Stake consensus mechanism
type PoS struct {
	validatorSet []crypto.PublicKey
}

// NewPoS initializes a new PoS instance
func NewPoS(validators []crypto.PublicKey) *PoS {
	return &PoS{
		validatorSet: validators,
	}
}

// SelectValidators selects a set of validators for block validation
func (pos *PoS) SelectValidators(numValidators int) ([]Validator, error) {
	if len(pos.validatorSet) < numValidators {
		return nil, errors.New("not enough validators")
	}
	var validators []Validator
	for i := 0; i < numValidators; i++ {
		validators = append(validators, Validator{
			Address: pos.validatorSet[i].Address(),
		})
	}
	return validators, nil
}

// SelectValidator selects a validator for the given block
func (pos *PoS) SelectValidator(block *core.Block, chain *core.Blockchain) error {
	validators, err := pos.SelectValidators(1)
	if err != nil {
		return err
	}

	block.Proposer = validators[0].Address
	return nil
}

// Validator represents a blockchain validator
type Validator struct {
	Address crypto.Address
}

// StakeManager manages the staking process
type StakeManager struct {
	// Implementation details omitted for brevity
}

// NewStakeManager initializes a new StakeManager instance
func NewStakeManager() *StakeManager {
	return &StakeManager{}
}

// AddStake adds stake for a given address
func (sm *StakeManager) AddStake(address crypto.Address, amount int64) {
	// Implementation details omitted for brevity
}

// RemoveStake removes stake for a given address
func (sm *StakeManager) RemoveStake(address crypto.Address, amount int64) {
	// Implementation details omitted for brevity
}
