// normalizer.go - Semantic transformations, independent evolution
package normalizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ferchd/tm2hsl/internal/ir"
	"github.com/ferchd/tm2hsl/internal/parser"
)

// Normalizer - Converts TextMate AST to formal IR
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
			// Features not supported in v0:
			// "include-repository": false,
			// "begin-captures":     false,
			// "end-captures":       false,
			// "while":              false,
		},
		strictMode: true,
	}
}

// Normalize - Main semantic transformation
func (n *Normalizer) Normalize(ast *parser.TextMateAST) (*ir.StateMachine, error) {
	if err := n.validateAST(ast); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	machine, err := ir.BuildFromAST(ast)
	if err != nil {
		return nil, fmt.Errorf("failed to build state machine: %w", err)
	}

	// 1. Expand includes and flatten patterns
	expanded := n.expandPatterns(ast.Patterns, ast.Repository)

	// 2. Convert to states and transitions
	for _, pattern := range expanded {
		if err := n.convertPattern(pattern, machine); err != nil {
			return nil, err
		}
	}

	// 3. Resolve references and optimize structure
	n.resolveReferences(machine)

	// 4. Apply specific semantic transformations
	n.applySemanticTransforms(machine)

	return machine, nil
}

// validateAST - Rejects unsupported features
func (n *Normalizer) validateAST(ast *parser.TextMateAST) error {
	var unsupported []string

	// Check repository (not supported in v0)
	if len(ast.Repository) > 0 {
		unsupported = append(unsupported, "repository")
	}

	// Check complex includes
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

// convertPattern - Explicit mapping of TextMate concepts to IR
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

// applySemanticTransforms - Smart logic here, not in parser
func (n *Normalizer) applySemanticTransforms(machine *ir.StateMachine) {
	// Example: transform named captures to tokens
	n.normalizeCaptureNames(machine)

	// Example: remove redundant states
	n.removeRedundantStates(machine)

	// Example: optimize transitions with same priority
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
	// Handle special includes
	switch {
	case include == "$self":
		// $self refers to the current grammar's patterns - for now, return empty to avoid recursion
		return nil
	case include == "$base":
		// $base would refer to base grammar patterns - not implemented yet
		return nil
	case strings.HasPrefix(include, "#"):
		// Named reference within repository, e.g., "#comment"
		name := strings.TrimPrefix(include, "#")
		if rule, exists := repository[name]; exists {
			return []parser.GrammarRule{rule}
		}
		return nil
	default:
		// Direct repository reference
		if rule, exists := repository[include]; exists {
			return []parser.GrammarRule{rule}
		}
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

	// Create actions from captures
	actions := n.createActionsFromCaptures(pattern.Captures, machine)

	transition := ir.Transition{
		Predicate: predicate,
		Target:    startStateID,
		Actions:   actions,
		Priority:  0,
	}

	machine.States[startStateID].Transitions = append(machine.States[startStateID].Transitions, transition)

	return nil
}

// convertBeginEndPattern - Converts begin/end patterns to IR states
func (n *Normalizer) convertBeginEndPattern(pattern parser.GrammarRule, machine *ir.StateMachine) error {
	if machine.States == nil {
		machine.States = make(map[ir.StateID]*ir.State)
	}

	startStateID := ir.StateID(len(machine.States))
	intermediateStateID := startStateID + 1
	endStateID := startStateID + 2

	// Start state
	machine.States[startStateID] = &ir.State{
		ID:          startStateID,
		Transitions: []ir.Transition{},
		IsFinal:     false,
	}

	// Intermediate state (inside the block)
	machine.States[intermediateStateID] = &ir.State{
		ID:          intermediateStateID,
		Transitions: []ir.Transition{},
		IsFinal:     false,
	}

	// End state
	machine.States[endStateID] = &ir.State{
		ID:          endStateID,
		Transitions: []ir.Transition{},
		IsFinal:     true,
	}

	// Begin transition: from start to intermediate
	beginPredicate := &ir.RegexPredicate{
		Pattern:  pattern.Begin,
		Compiled: regexp.MustCompile(pattern.Begin),
	}
	var beginActions []ir.ActionID // TODO: push scopes
	beginTransition := ir.Transition{
		Predicate: beginPredicate,
		Target:    intermediateStateID,
		Actions:   beginActions,
		Priority:  0,
	}
	machine.States[startStateID].Transitions = append(machine.States[startStateID].Transitions, beginTransition)

	// Process child patterns in intermediate state
	for _, child := range pattern.Patterns {
		if err := n.convertPattern(child, machine); err != nil {
			return err
		}
	}

	// End transition: from intermediate to end
	endPredicate := &ir.RegexPredicate{
		Pattern:  pattern.End,
		Compiled: regexp.MustCompile(pattern.End),
	}
	var endActions []ir.ActionID // TODO: pop scopes
	endTransition := ir.Transition{
		Predicate: endPredicate,
		Target:    endStateID,
		Actions:   endActions,
		Priority:  0,
	}
	machine.States[intermediateStateID].Transitions = append(machine.States[intermediateStateID].Transitions, endTransition)

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

// createActionsFromCaptures - Creates IR actions from capture definitions
func (n *Normalizer) createActionsFromCaptures(captures map[int]parser.Capture, machine *ir.StateMachine) []ir.ActionID {
	var actions []ir.ActionID
	for groupIndex, capture := range captures {
		actionID := ir.ActionID(len(machine.Actions))
		machine.Actions[actionID] = &ir.PushScopeAction{
			Scope: capture.Name,
			Index: groupIndex,
		}
		actions = append(actions, actionID)
	}
	return actions
}
