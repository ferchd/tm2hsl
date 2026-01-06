package ir

import (
	"github.com/ferchd/tm2hsl/internal/parser"
)

type StateMachine struct {
	Name    string
	Initial StateID
	States  map[StateID]*State
	Tokens  map[TokenID]TokenDef
	Actions map[ActionID]*Action
}

type StateID uint32
type TokenID uint32
type ActionID uint32

type State struct {
	ID          StateID
	Transitions []Transition
	IsFinal     bool
	OnEntry     []ActionID
	OnExit      []ActionID
}

type Transition struct {
	Predicate Predicate
	Target    StateID
	Priority  uint8
	Consume   bool
	Actions   []ActionID
}

type TokenDef struct {
	ID   TokenID
	Name string
}

func BuildFromAST(ast *parser.TextMateAST) (*StateMachine, error) {
	// TODO: Implement builder pattern for AST conversion
	return &StateMachine{
		Name:    ast.ScopeName,
		Initial: 0,
		States:  make(map[StateID]*State),
		Tokens:  make(map[TokenID]TokenDef),
		Actions: make(map[ActionID]*Action),
	}, nil
}
