// normalizer.go - Transformaciones semánticas, evolución independiente
package normalizer

import (
	"fmt"
	"regexp"

	"github.com/ferchd/tm2hsl/internal/ir"
	"github.com/ferchd/tm2hsl/internal/parser"
)

// Normalizer - Convierte AST TextMate a IR formal
type Normalizer struct {
	supportedFeatures map[string]bool
	strictMode        bool
}

func NewNormalizer() *Normalizer {
	return &Normalizer{
		supportedFeatures: map[string]bool{
			"match":        true,
			"begin-end":    true,
			"captures":     true,
			"contentName":  true,
			"include-self": true, // $self
			"include-base": true, // $base
			// Features NO soportados en v0:
			// "include-repository": false,
			// "begin-captures":     false,
			// "end-captures":       false,
			// "while":              false,
		},
		strictMode: true,
	}
}

// Normalize - Transformación semántica principal
func (n *Normalizer) Normalize(ast *parser.TextMateAST) (*ir.StateMachine, error) {
	if err := n.validateAST(ast); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	machine, err := ir.BuildFromAST(ast)
	if err != nil {
		return nil, fmt.Errorf("failed to build state machine: %w", err)
	}

	// 1. Expandir includes y flatten patterns
	expanded := n.expandPatterns(ast.Patterns, ast.Repository)

	// 2. Convertir a estados y transiciones
	for _, pattern := range expanded {
		if err := n.convertPattern(pattern, machine); err != nil {
			return nil, err
		}
	}

	// 3. Resolver referencias y optimizar estructura
	n.resolveReferences(machine)

	// 4. Aplicar transformaciones semánticas específicas
	n.applySemanticTransforms(machine)

	return machine, nil
}

// validateAST - Rechaza features no soportados
func (n *Normalizer) validateAST(ast *parser.TextMateAST) error {
	var unsupported []string

	// Verificar repository (no soportado en v0)
	if len(ast.Repository) > 0 {
		unsupported = append(unsupported, "repository")
	}

	// Verificar includes complejos
	if hasComplexIncludes(ast) {
		unsupported = append(unsupported, "complex-includes")
	}

	if len(unsupported) > 0 && n.strictMode {
		return fmt.Errorf("unsupported features: %v", unsupported)
	}

	return nil
}

// hasComplexIncludes - Stub implementation
func hasComplexIncludes(ast *parser.TextMateAST) bool {
	return false
}

// convertPattern - Mapeo explícito de conceptos TextMate a IR
func (n *Normalizer) convertPattern(pattern parser.GrammarRule, machine *ir.StateMachine) error {
	switch {
	case pattern.Match != "":
		return n.convertMatchPattern(pattern, machine)
	case pattern.Begin != "" && pattern.End != "":
		return n.convertBeginEndPattern(pattern, machine)
	case pattern.Include != "":
		return n.convertInclude(pattern, machine)
	default:
		return fmt.Errorf("unsupported pattern type")
	}
}

// applySemanticTransforms - Lógica inteligente aquí, no en parser
func (n *Normalizer) applySemanticTransforms(machine *ir.StateMachine) {
	// Ejemplo: transformar capturas con nombre a tokens
	n.normalizeCaptureNames(machine)

	// Ejemplo: eliminar estados redundantes
	n.removeRedundantStates(machine)

	// Ejemplo: optimizar transiciones con misma prioridad
	n.optimizeTransitionOrder(machine)
}

// expandPatterns - Expands includes and flattens patterns
func (n *Normalizer) expandPatterns(patterns []parser.GrammarRule, repository map[string]parser.GrammarRule) []parser.GrammarRule {
	var expanded []parser.GrammarRule
	for _, pattern := range patterns {
		if pattern.Include != "" {
			if expandedPatterns := n.expandInclude(pattern.Include, repository); expandedPatterns != nil {
				expanded = append(expanded, expandedPatterns...)
			}
		} else {
			expanded = append(expanded, pattern)
		}
	}
	return expanded
}

// expandInclude - Expands a single include reference
func (n *Normalizer) expandInclude(include string, repository map[string]parser.GrammarRule) []parser.GrammarRule {
	if rule, exists := repository[include]; exists {
		return []parser.GrammarRule{rule}
	}
	// Handle $self, $base, etc.
	switch include {
	case "$self":
		// Return self patterns, but for simplicity, ignore
		return nil
	case "$base":
		return nil
	default:
		return nil
	}
}

// resolveReferences - Stub implementation
func (n *Normalizer) resolveReferences(machine *ir.StateMachine) {
	// TODO: implement
}

// convertMatchPattern - Converts a match pattern to IR states
func (n *Normalizer) convertMatchPattern(pattern parser.GrammarRule, machine *ir.StateMachine) error {
	if machine.States == nil {
		machine.States = make(map[ir.StateID]*ir.State)
	}
	startStateID := ir.StateID(len(machine.States))
	machine.States[startStateID] = &ir.State{
		ID:          startStateID,
		Transitions: []ir.Transition{},
		IsFinal:     true,
	}

	// Create predicate
	predicate := &ir.RegexPredicate{
		Pattern:  pattern.Match,
		Compiled: regexp.MustCompile(pattern.Match),
	}

	// Create actions (stub)
	var actions []ir.ActionID

	transition := ir.Transition{
		Predicate: predicate,
		Target:    startStateID,
		Actions:   actions,
		Priority:  0,
	}

	machine.States[startStateID].Transitions = append(machine.States[startStateID].Transitions, transition)

	return nil
}

// convertBeginEndPattern - Stub implementation
func (n *Normalizer) convertBeginEndPattern(pattern parser.GrammarRule, machine *ir.StateMachine) error {
	return nil
}

// convertInclude - Stub implementation
func (n *Normalizer) convertInclude(pattern parser.GrammarRule, machine *ir.StateMachine) error {
	return nil
}

// normalizeCaptureNames - Stub implementation
func (n *Normalizer) normalizeCaptureNames(machine *ir.StateMachine) {
	// TODO: implement
}

// removeRedundantStates - Stub implementation
func (n *Normalizer) removeRedundantStates(machine *ir.StateMachine) {
	// TODO: implement
}

// optimizeTransitionOrder - Stub implementation
func (n *Normalizer) optimizeTransitionOrder(machine *ir.StateMachine) {
	// TODO: implement
}
