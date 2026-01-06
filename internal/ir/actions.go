package ir

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/parser"
)

// Action - Acción atómica en la máquina de estados
type Action interface {
	Type() ActionType
	String() string
	Validate() error
}

// ActionType - Tipos de acciones soportadas
type ActionType int

const (
	ActionPushScope    ActionType = iota // Equivalente a begin
	ActionPopScope                       // Equivalente a end
	ActionSetScope                       // Cambiar scope actual
	ActionEmitToken                      // Emitir token
	ActionPushState                      // Push state stack
	ActionPopState                       // Pop state stack
	ActionCaptureGroup                   // Capturar grupo
)

// PushScopeAction - Comenzar nuevo scope
type PushScopeAction struct {
	Scope string
	Index int // Para capturas numeradas
}

func (a *PushScopeAction) Type() ActionType { return ActionPushScope }
func (a *PushScopeAction) String() string   { return fmt.Sprintf("push-scope:%s", a.Scope) }
func (a *PushScopeAction) Validate() error  { return nil }

// PopScopeAction - Terminar scope
type PopScopeAction struct {
	Count int // Número de scopes a pop (usualmente 1)
}

func (a *PopScopeAction) Type() ActionType { return ActionPopScope }
func (a *PopScopeAction) String() string   { return "pop-scope" }
func (a *PopScopeAction) Validate() error  { return nil }

// SetScopeAction - Cambiar scope actual
type SetScopeAction struct {
	Scope string
}

func (a *SetScopeAction) Type() ActionType { return ActionSetScope }
func (a *SetScopeAction) String() string   { return fmt.Sprintf("set-scope:%s", a.Scope) }
func (a *SetScopeAction) Validate() error  { return nil }

// EmitTokenAction - Emitir token con scopes acumulados
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

// CaptureGroupAction - Capturar texto para grupo
type CaptureGroupAction struct {
	GroupID int
	Name    string
}

func (a *CaptureGroupAction) Type() ActionType { return ActionCaptureGroup }
func (a *CaptureGroupAction) String() string {
	return fmt.Sprintf("capture-group:%d:%s", a.GroupID, a.Name)
}
func (a *CaptureGroupAction) Validate() error { return nil }

// ActionSequence - Secuencia de acciones
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
