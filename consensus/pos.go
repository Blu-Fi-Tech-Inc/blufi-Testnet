package consensus

import (
    "errors"
    "math/rand"
    "sync"
    "time"

    "github.com/blu-fi-tech-inc/blufi-network/types"
)

type PoS struct {
    stakeManager *StakeManager
    mu           sync.RWMutex
}

func NewPoS(stakeManager *StakeManager) *PoS {
    return &PoS{
        stakeManager: stakeManager,
    }
}

func (pos *PoS) ValidateBlock(block *types.Block) bool {
    pos.mu.RLock()
    defer pos.mu.RUnlock()

    validators, err := SelectValidators(pos.stakeManager, 10) // Assume 10 validators
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

func (pos *PoS) AddBlock(block *types.Block) bool {
    if pos.ValidateBlock(block) {
        // Add block to the blockchain (omitted for brevity)
        return true
    }
    return false
}

type Validator struct {
    Address string
    Stake   int
}

func SelectValidators(stakeManager *StakeManager, numValidators int) ([]Validator, error) {
    stakeManager.mu.Lock()
    defer stakeManager.mu.Unlock()

    if len(stakeManager.stakes) < numValidators {
        return nil, errors.New("not enough validators")
    }

    var totalStake int
    for _, stake := range stakeManager.stakes {
        totalStake += stake
    }

    selectedValidators := make([]Validator, 0, numValidators)
    for i := 0; i < numValidators; i++ {
        rand.Seed(time.Now().UnixNano())
        target := rand.Intn(totalStake)
        cumulativeStake := 0

        for address, stake := range stakeManager.stakes {
            cumulativeStake += stake
            if cumulativeStake > target {
                selectedValidators = append(selectedValidators, Validator{Address: address, Stake: stake})
                break
            }
        }
    }

    return selectedValidators, nil
}

type StakeManager struct {
    stakes map[string]int
    mu     sync.Mutex
}

func NewStakeManager() *StakeManager {
    return &StakeManager{
        stakes: make(map[string]int),
    }
}

func (sm *StakeManager) AddStake(address string, amount int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.stakes[address] += amount
}

func (sm *StakeManager) GetStake(address string) int {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    return sm.stakes[address]
}
