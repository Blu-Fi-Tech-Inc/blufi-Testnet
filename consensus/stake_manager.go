package consensus

import (
	"errors"
	"sync"
)

// Stake represents the amount of stake held by a stakeholder.
type Stake struct {
	Amount uint64
}

// Stakeholder represents an individual stakeholder with an address and their stake.
type Stakeholder struct {
	Address string
	Stake   Stake
}

// StakeManager manages the stakes of all stakeholders in the system.
type StakeManager struct {
	stakeholders map[string]Stakeholder
	mu           sync.RWMutex
}

// NewStakeManager creates a new instance of StakeManager.
func NewStakeManager() *StakeManager {
	return &StakeManager{
		stakeholders: make(map[string]Stakeholder),
	}
}

// AddStake adds a given amount of stake to the specified address.
func (sm *StakeManager) AddStake(address string, amount uint64) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	stakeholder, exists := sm.stakeholders[address]
	if !exists {
		stakeholder = Stakeholder{Address: address, Stake: Stake{Amount: amount}}
	} else {
		stakeholder.Stake.Amount += amount
	}
	sm.stakeholders[address] = stakeholder
	return nil
}

// GetStake returns the stake of the specified address.
func (sm *StakeManager) GetStake(address string) (Stake, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	stakeholder, exists := sm.stakeholders[address]
	if !exists {
		return Stake{}, errors.New("stakeholder not found")
	}
	return stakeholder.Stake, nil
}

// GetStakeholders returns a copy of all stakeholders.
func (sm *StakeManager) GetStakeholders() map[string]Stakeholder {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	copy := make(map[string]Stakeholder)
	for k, v := range sm.stakeholders {
		copy[k] = v
	}
	return copy
}
