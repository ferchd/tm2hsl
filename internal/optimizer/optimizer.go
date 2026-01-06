// optimizer.go - Only structural transformations
package optimizer

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/ir"
)

// Optimizer - Does not understand scopes or semantics
type Optimizer struct {
	// Purely structural optimization rules
	passes []OptimizationPass
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		passes: []OptimizationPass{
			&RemoveUnreachableStates{},
			&MergeEquivalentStates{},
			&SimplifyTransitions{},
			&ReorderByPriority{},
		},
	}
}

// Optimize - Applies steps without changing semantics
func (o *Optimizer) Optimize(program *ir.Program) (*ir.Program, error) {
	for _, pass := range o.passes {
		changed, err := pass.Apply(program)
		if err != nil {
			return nil, fmt.Errorf("pass %s failed: %w", pass.Name(), err)
		}
		if !changed {
			// If nothing changed, we can stop
			break
		}
	}
	return program, nil
}
