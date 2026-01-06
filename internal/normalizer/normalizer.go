// normalizer.go - Transformaciones semánticas, evolución independiente
package normalizer

import (
    "github.com/ferchd/tm2hsl/ir"
    "github.com/ferchd/tm2hsl/parser"
)

// Normalizer - Convierte AST TextMate a IR formal
type Normalizer struct {
    supportedFeatures map[string]bool
    strictMode        bool
}

func NewNormalizer() *Normalizer {
    return &Normalizer{
        supportedFeatures: map[string]bool{
            "match":          true,
            "begin-end":      true,
            "captures":       true,
            "contentName":    true,
            "include-self":   true,  // $self
            "include-base":   true,  // $base
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
    
    machine := ir.NewStateMachine(ast.ScopeName)
    
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