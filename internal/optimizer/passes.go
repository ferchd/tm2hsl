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
	// TODO: Implement unreachable state removal
	return false, nil
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
