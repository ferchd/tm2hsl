package ir

import (
	"fmt"
	"regexp"
)

// Predicate - Interfaz para condiciones de transición
type Predicate interface {
	Type() PredicateType
	String() string
	Equal(other Predicate) bool
}

// PredicateType - Tipos de predicados soportados
type PredicateType int

const (
	PredicateChar       PredicateType = iota // Carácter simple
	PredicateCharSet                         // Conjunto de caracteres [a-z]
	PredicateCharClass                       // Clase \w, \d, \s
	PredicateString                          // String literal
	PredicateRegex                           // Expresión regular
	PredicateAny                             // Cualquier carácter
	PredicateEOF                             // Fin de archivo
	PredicateLookahead                       // Lookahead positivo/negativo
	PredicateLookbehind                      // Lookbehind positivo/negativo
)

// CharPredicate - Predicado de carácter único
type CharPredicate struct {
	Char rune
}

func (p *CharPredicate) Type() PredicateType { return PredicateChar }
func (p *CharPredicate) String() string      { return fmt.Sprintf("char('%c')", p.Char) }

// CharSetPredicate - Conjunto de caracteres
type CharSetPredicate struct {
	Chars   string      // Caracteres literales
	Ranges  []CharRange // Rangos [a-z]
	Negated bool        // [^abc]
}

type CharRange struct {
	Start, End rune
}

// RegexPredicate - Expresión regular compilada
type RegexPredicate struct {
	Pattern  string
	Compiled *regexp.Regexp // Mantenido para referencia
	Simple   bool           // True si es una regex simple (sin grupos complejos)
}

func (p *RegexPredicate) Type() PredicateType { return PredicateRegex }
func (p *RegexPredicate) String() string      { return fmt.Sprintf("regex:%s", p.Pattern) }
func (p *RegexPredicate) Equal(other Predicate) bool {
	if o, ok := other.(*RegexPredicate); ok {
		return p.Pattern == o.Pattern
	}
	return false
}

// AnyPredicate - Cualquier carácter
type AnyPredicate struct{}

func (p *AnyPredicate) Type() PredicateType { return PredicateAny }
func (p *AnyPredicate) String() string      { return "any" }
func (p *AnyPredicate) Equal(other Predicate) bool {
	_, ok := other.(*AnyPredicate)
	return ok
}

// LookaheadPredicate - Lookahead
type LookaheadPredicate struct {
	Predicate Predicate
	Positive  bool // true=positivo, false=negativo
}

// CompoundPredicate - Predicado compuesto (AND/OR)
type CompoundPredicate struct {
	Op    CompoundOp
	Preds []Predicate
}

type CompoundOp int

const (
	OpAnd CompoundOp = iota
	OpOr
)

// PredicateBuilder - Constructor de predicados
type PredicateBuilder struct {
	strict bool
}

func NewPredicateBuilder(strict bool) *PredicateBuilder {
	return &PredicateBuilder{strict: strict}
}

// FromTextMateMatch - Convierte regex TextMate a predicados IR
func (b *PredicateBuilder) FromTextMateMatch(pattern string) (Predicate, error) {
	if b.strict && containsUnsupportedRegex(pattern) {
		return nil, fmt.Errorf("unsupported regex pattern: %s", pattern)
	}

	// Simplificar regex TextMate a predicados básicos
	return b.simplifyRegex(pattern)
}

func (b *PredicateBuilder) simplifyRegex(pattern string) (Predicate, error) {
	// Convertir a predicados más simples cuando sea posible
	switch {
	case pattern == ".":
		return &AnyPredicate{}, nil
	case isCharClass(pattern):
		return b.buildCharClassPredicate(pattern)
	case isLiteral(pattern):
		return b.buildLiteralPredicate(pattern)
	default:
		// Mantener como regex genérica
		return &RegexPredicate{
			Pattern: pattern,
			Simple:  false,
		}, nil
	}

}

// Helper functions
func containsUnsupportedRegex(pattern string) bool {
	// TODO: Implement unsupported regex detection
	return false
}

func isCharClass(pattern string) bool {
	// TODO: Implement character class detection
	return false
}

func (b *PredicateBuilder) buildCharClassPredicate(pattern string) (Predicate, error) {
	// TODO: Implement character class predicate building
	return &AnyPredicate{}, nil
}

func isLiteral(pattern string) bool {
	// TODO: Implement literal detection
	return false
}

func (b *PredicateBuilder) buildLiteralPredicate(pattern string) (Predicate, error) {
	// TODO: Implement literal predicate building
	return &AnyPredicate{}, nil
}
