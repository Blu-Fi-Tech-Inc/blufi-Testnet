package consensus

import (
	"errors"
	"sync"
)

type Stake struct {
	Amount uint64
}

type Stakeholder struct {
	Address string
	Stake   Stake
}

type StakeManager struct {
	stakeholders map[string]Stakeholder
	mu           sync.RWMutex
}

func NewStakeManager() *StakeManager {
	return &StakeManager{
		stakeholders: make(map[string]Stakeholder),
	}
}

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

func (sm *StakeManager) GetStake(address string) (Stake, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	stakeholder, exists := sm.stakeholders[address]
	if !exists {
		return Stake{}, errors.New("stakeholder not found")
	}
	return stakeholder.Stake, nil
}
