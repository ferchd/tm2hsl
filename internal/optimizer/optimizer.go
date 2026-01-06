// optimizer.go - Solo transformaciones estructurales
package optimizer

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/ir"
)

// Optimizer - No entiende scopes ni semántica
type Optimizer struct {
	// Reglas de optimización puramente estructurales
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

// Optimize - Aplica pasos sin cambiar semántica
func (o *Optimizer) Optimize(program *ir.Program) (*ir.Program, error) {
	for _, pass := range o.passes {
		changed, err := pass.Apply(program)
		if err != nil {
			return nil, fmt.Errorf("pass %s failed: %w", pass.Name(), err)
		}
		if !changed {
			// Si no cambió nada, podemos detenernos
			break
		}
	}
	return program, nil
}
