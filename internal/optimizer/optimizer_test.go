package optimizer

import (
	"testing"

	"github.com/ferchd/tm2hsl/internal/ir"
)

func TestNewOptimizer(t *testing.T) {
	o := NewOptimizer()
	if o == nil {
		t.Error("NewOptimizer() returned nil")
	}
	if len(o.passes) == 0 {
		t.Error("passes not initialized")
	}
}

func TestOptimizer_Optimize(t *testing.T) {
	o := NewOptimizer()

	program := ir.NewProgram("test", "source.test")

	optimized, err := o.Optimize(program)
	if err != nil {
		t.Errorf("Optimizer.Optimize() error = %v", err)
	}
	if optimized == nil {
		t.Error("Optimizer.Optimize() returned nil")
	}
}
