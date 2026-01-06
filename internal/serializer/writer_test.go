package serializer

import (
	"bytes"
	"testing"

	"github.com/ferchd/tm2hsl/internal/ir"
)

func TestNewBytecodeWriter(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBytecodeWriter(&buf)
	if bw == nil {
		t.Error("NewBytecodeWriter() returned nil")
	}
	if bw.w == nil {
		t.Error("writer not set")
	}
}

func TestBytecodeWriter_WriteBytecode(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBytecodeWriter(&buf)

	sm := &ir.StateMachine{
		Name:    "test",
		Initial: 0,
		States:  make(map[ir.StateID]*ir.State),
		Tokens:  make(map[ir.TokenID]ir.TokenDef),
		Actions: make(map[ir.ActionID]ir.Action),
	}

	err := bw.WriteBytecode(sm)
	if err != nil {
		t.Errorf("BytecodeWriter.WriteBytecode() error = %v", err)
	}
}
