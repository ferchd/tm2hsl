package optimizer

import (
	"sort"

	"github.com/ferchd/tm2hsl/internal/ir"
)

// OptimizationPass - Optimization step
type OptimizationPass interface {
	Name() string
	Apply(program *ir.Program) (bool, error)
}

// RemoveUnreachableStates - Removes unreachable states
type RemoveUnreachableStates struct{}

func (p *RemoveUnreachableStates) Name() string { return "remove-unreachable-states" }

func (p *RemoveUnreachableStates) Apply(program *ir.Program) (bool, error) {
	if len(program.StateTable) == 0 {
		return false, nil
	}

	reachable := make(map[uint32]bool)
	var visit func(uint32)
	visit = func(stateID uint32) {
		if reachable[stateID] {
			return
		}
		reachable[stateID] = true
		if stateID >= uint32(len(program.StateTable)) {
			return
		}
		state := program.StateTable[stateID]
		rules := program.RuleTable[state.RuleOffset : state.RuleOffset+uint32(state.RuleCount)]
		for _, rule := range rules {
			if rule.NextState >= 0 {
				visit(uint32(rule.NextState))
			}
		}
	}
	visit(0) // assume initial state is 0

	// Create new tables
	newStateTable := []ir.StateEntry{}
	newRuleTable := []ir.RuleEntry{}
	stateMap := make(map[uint32]uint32)
	newID := uint32(0)
	for oldID, state := range program.StateTable {
		if reachable[uint32(oldID)] {
			stateMap[uint32(oldID)] = newID
			newStateTable = append(newStateTable, state)
			newStateTable[len(newStateTable)-1].ID = newID
			newID++
		}
	}

	ruleOffset := uint32(0)
	for i := range newStateTable {
		state := &newStateTable[i]
		oldOffset := state.RuleOffset
		oldCount := uint32(state.RuleCount)
		rules := program.RuleTable[oldOffset : oldOffset+oldCount]
		newRules := []ir.RuleEntry{}
		for _, rule := range rules {
			newRule := rule
			if rule.NextState >= 0 {
				if !reachable[uint32(rule.NextState)] {
					continue
				}
				newRule.NextState = int32(stateMap[uint32(rule.NextState)])
			}
			newRules = append(newRules, newRule)
		}
		state.RuleOffset = ruleOffset
		state.RuleCount = uint16(len(newRules))
		newRuleTable = append(newRuleTable, newRules...)
		ruleOffset += uint32(len(newRules))
	}

	program.StateTable = newStateTable
	program.RuleTable = newRuleTable
	return len(newStateTable) < len(program.StateTable), nil
}

// MergeEquivalentStates - Merges equivalent states
type MergeEquivalentStates struct{}

func (p *MergeEquivalentStates) Name() string { return "merge-equivalent-states" }

func (p *MergeEquivalentStates) Apply(program *ir.Program) (bool, error) {
	// TODO: Implement equivalent state merging
	return false, nil
}

// SimplifyTransitions - Simplifies redundant transitions
type SimplifyTransitions struct{}

func (p *SimplifyTransitions) Name() string { return "simplify-transitions" }

func (p *SimplifyTransitions) Apply(program *ir.Program) (bool, error) {
	// TODO: Implement transition simplification
	return false, nil
}

// ReorderByPriority - Reorders rules by priority
type ReorderByPriority struct{}

func (p *ReorderByPriority) Name() string { return "reorder-by-priority" }

func (p *ReorderByPriority) Apply(program *ir.Program) (bool, error) {
	// Sort rules by priority descending (higher priority first)
	sort.Slice(program.RuleTable, func(i, j int) bool {
		return program.RuleTable[i].Priority > program.RuleTable[j].Priority
	})
	return true, nil // Always reorder for now
}
