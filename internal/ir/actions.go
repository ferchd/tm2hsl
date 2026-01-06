package ir

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/parser"
)

// Action - Atomic action in the state machine
type Action interface {
	Type() ActionType
	String() string
	Validate() error
}

// ActionType - Supported action types
type ActionType int

const (
	ActionPushScope    ActionType = iota // Equivalent to begin
	ActionPopScope                       // Equivalent to end
	ActionSetScope                       // Change current scope
	ActionEmitToken                      // Emit token
	ActionPushState                      // Push state stack
	ActionPopState                       // Pop state stack
	ActionCaptureGroup                   // Capture group
)

// PushScopeAction - Start new scope
type PushScopeAction struct {
	Scope string
	Index int // For numbered captures
}

func (a *PushScopeAction) Type() ActionType { return ActionPushScope }
func (a *PushScopeAction) String() string   { return fmt.Sprintf("push-scope:%s", a.Scope) }
func (a *PushScopeAction) Validate() error  { return nil }

// PopScopeAction - End scope
type PopScopeAction struct {
	Count int // Number of scopes to pop (usually 1)
}

func (a *PopScopeAction) Type() ActionType { return ActionPopScope }
func (a *PopScopeAction) String() string   { return "pop-scope" }
func (a *PopScopeAction) Validate() error  { return nil }

// SetScopeAction - Change current scope
type SetScopeAction struct {
	Scope string
}

func (a *SetScopeAction) Type() ActionType { return ActionSetScope }
func (a *SetScopeAction) String() string   { return fmt.Sprintf("set-scope:%s", a.Scope) }
func (a *SetScopeAction) Validate() error  { return nil }

// EmitTokenAction - Emit token with accumulated scopes
type EmitTokenAction struct {
	Text  string
	Start Position
	End   Position
}

func (a *EmitTokenAction) Type() ActionType { return ActionEmitToken }
func (a *EmitTokenAction) String() string   { return fmt.Sprintf("emit-token:%s", a.Text) }
func (a *EmitTokenAction) Validate() error  { return nil }

type Position struct {
	Offset int
	Line   int
	Column int
}

// CaptureGroupAction - Capture text for group
type CaptureGroupAction struct {
	GroupID int
	Name    string
}

func (a *CaptureGroupAction) Type() ActionType { return ActionCaptureGroup }
func (a *CaptureGroupAction) String() string {
	return fmt.Sprintf("capture-group:%d:%s", a.GroupID, a.Name)
}
func (a *CaptureGroupAction) Validate() error { return nil }

// ActionSequence - Sequence of actions
type ActionSequence struct {
	Actions []Action
}

// ActionBuilder - Constructor de acciones
type ActionBuilder struct {
	scopeStack []string
	captureID  int
}

func NewActionBuilder() *ActionBuilder {
	return &ActionBuilder{
		scopeStack: []string{},
	}
}

// FromTextMateCaptures - Convierte capturas TextMate a acciones
func (b *ActionBuilder) FromTextMateCaptures(captures map[int]parser.Capture) ([]Action, error) {
	var actions []Action

	for idx, capture := range captures {
		if idx == 0 {
			// Captura 0 es el match completo
			actions = append(actions, &EmitTokenAction{
				Text: "", // Se llena durante ejecución
			})
		} else {
			actions = append(actions, &CaptureGroupAction{
				GroupID: idx,
				Name:    capture.Name,
			})
		}
	}

	return actions, nil
}
