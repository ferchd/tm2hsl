package normalizer

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/ir"
	"github.com/ferchd/tm2hsl/internal/parser"
)

// PatternConverter - Convierte patrones TextMate a IR
type PatternConverter struct {
	machineBuilder   *ir.MachineBuilder
	predicateBuilder *ir.PredicateBuilder
	actionBuilder    *ir.ActionBuilder
	scopeStack       []string
}

func NewPatternConverter() *PatternConverter {
	return &PatternConverter{
		machineBuilder:   ir.NewMachineBuilder(),
		predicateBuilder: ir.NewPredicateBuilder(true),
		actionBuilder:    ir.NewActionBuilder(),
		scopeStack:       []string{},
	}
}

// ConvertPattern - Convierte un patrón TextMate a estados IR
func (c *PatternConverter) ConvertPattern(pattern parser.GrammarRule) error {
	switch {
	case pattern.Match != "":
		return c.convertMatchPattern(pattern)
	case pattern.Begin != "" && pattern.End != "":
		return c.convertBeginEndPattern(pattern)
	case pattern.Include != "":
		return c.convertIncludePattern(pattern)
	default:
		// Patrón contenedor - convertir patrones hijos
		for _, child := range pattern.Patterns {
			if err := c.ConvertPattern(child); err != nil {
				return err
			}
		}
		return nil
	}
}

// convertMatchPattern - Convierte patrón match
func (c *PatternConverter) convertMatchPattern(pattern parser.GrammarRule) error {
	// Crear predicado desde regex
	predicate, err := c.predicateBuilder.FromTextMateMatch(pattern.Match)
	if err != nil {
		return fmt.Errorf("match pattern %s: %w", pattern.Name, err)
	}

	// Crear acciones desde captures
	actions, err := c.actionBuilder.FromTextMateCaptures(pattern.Captures)
	if err != nil {
		return fmt.Errorf("captures for %s: %w", pattern.Name, err)
	}

	// Añadir acción de push/pop si hay contentName
	if pattern.ContentName != "" {
		actions = append([]ir.Action{
			&ir.PushScopeAction{Scope: pattern.ContentName},
		}, actions...)
		actions = append(actions, &ir.PopScopeAction{Count: 1})
	}

	// Crear transición
	transition := &ir.Transition{
		Predicate: predicate,
		Actions:   actions,
		Priority:  c.calculatePriority(pattern),
	}

	// Añadir al builder
	stateID := c.machineBuilder.AddState()
	c.machineBuilder.AddTransition(stateID, transition)

	return nil
}

// convertBeginEndPattern - Convierte patrón begin/end
func (c *PatternConverter) convertBeginEndPattern(pattern parser.GrammarRule) error {
	startState := c.machineBuilder.AddState()
	endState := c.machineBuilder.AddState()

	// Transición begin
	beginPred, err := c.predicateBuilder.FromTextMateMatch(pattern.Begin)
	if err != nil {
		return fmt.Errorf("begin pattern %s: %w", pattern.Name, err)
	}

	beginActions, err := c.actionBuilder.FromTextMateCaptures(pattern.BeginCaptures)
	if err != nil {
		return fmt.Errorf("begin captures: %w", err)
	}

	// Push del scope principal
	if pattern.Name != "" {
		beginActions = append(beginActions, &ir.PushScopeAction{
			Scope: pattern.Name,
		})
	}

	// Push del scope interno si hay contentName
	if pattern.ContentName != "" {
		beginActions = append(beginActions, &ir.PushScopeAction{
			Scope: pattern.ContentName,
		})
	}

	c.machineBuilder.AddTransition(startState, &ir.Transition{
		Predicate: beginPred,
		Target:    endState,
		Actions:   beginActions,
		Priority:  c.calculatePriority(pattern),
	})

	// Transición end
	endPred, err := c.predicateBuilder.FromTextMateMatch(pattern.End)
	if err != nil {
		return fmt.Errorf("end pattern %s: %w", pattern.Name, err)
	}

	endActions, err := c.actionBuilder.FromTextMateCaptures(pattern.EndCaptures)
	if err != nil {
		return fmt.Errorf("end captures: %w", err)
	}

	// Pop de scopes
	popCount := 0
	if pattern.ContentName != "" {
		popCount++
	}
	if pattern.Name != "" {
		popCount++
	}
	if popCount > 0 {
		endActions = append(endActions, &ir.PopScopeAction{
			Count: popCount,
		})
	}

	c.machineBuilder.AddTransition(endState, &ir.Transition{
		Predicate: endPred,
		Actions:   endActions,
		Priority:  c.calculatePriority(pattern),
	})

	// Convertir patrones internos en el estado intermedio
	savedStack := c.scopeStack
	if pattern.Name != "" {
		c.scopeStack = append(c.scopeStack, pattern.Name)
	}

	for _, child := range pattern.Patterns {
		if err := c.ConvertPattern(child); err != nil {
			return err
		}
	}

	c.scopeStack = savedStack
	return nil
}

// calculatePriority - Calcula prioridad basada en características del patrón
func (c *PatternConverter) calculatePriority(pattern parser.GrammarRule) uint8 {
	priority := uint8(50) // Prioridad media por defecto

	// Patrones con nombre tienen mayor prioridad
	if pattern.Name != "" {
		priority += 10
	}

	// Patrones con captures tienen mayor prioridad
	if len(pattern.Captures) > 0 {
		priority += 5
	}

	return priority
}
