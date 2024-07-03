package consensus

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/blu-fi-tech-inc/blufi-network/types"
)

type Validator struct {
	Address string
	Stake   uint64
}

type ValidatorSelection struct {
	stakeManager *StakeManager
	mu           sync.RWMutex
}

func NewValidatorSelection(stakeManager *StakeManager) *ValidatorSelection {
	return &ValidatorSelection{
		stakeManager: stakeManager,
	}
}

func (vs *ValidatorSelection) SelectValidators(numValidators int) ([]Validator, error) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	stakeholders := vs.stakeManager.getStakeholders()
	if len(stakeholders) < numValidators {
		return nil, errors.New("not enough validators")
	}

	var totalStake uint64
	for _, stakeholder := range stakeholders {
		totalStake += stakeholder.Stake.Amount
	}

	selectedValidators := make([]Validator, 0, numValidators)
	for i := 0; i < numValidators; i++ {
		rand.Seed(time.Now().UnixNano())
		target := rand.Int63n(int64(totalStake))
		var cumulativeStake uint64

		for _, stakeholder := range stakeholders {
			cumulativeStake += stakeholder.Stake.Amount
			if cumulativeStake > uint64(target) {
				selectedValidators = append(selectedValidators, Validator{
					Address: stakeholder.Address,
					Stake:   stakeholder.Stake.Amount,
				})
				totalStake -= stakeholder.Stake.Amount
				delete(stakeholders, stakeholder.Address)
				break
			}
		}
	}

	return selectedValidators, nil
}

func (sm *StakeManager) getStakeholders() map[string]Stakeholder {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	copy := make(map[string]Stakeholder)
	for k, v := range sm.stakeholders {
		copy[k] = v
	}
	return copy
}
