package ir

import (
	"fmt"
	"regexp"
)

// Predicate - Interface for transition conditions
type Predicate interface {
	Type() PredicateType
	String() string
	Equal(other Predicate) bool
}

// PredicateType - Supported predicate types
type PredicateType int

const (
	PredicateChar       PredicateType = iota // Simple character
	PredicateCharSet                         // Character set [a-z]
	PredicateCharClass                       // Class \w, \d, \s
	PredicateString                          // String literal
	PredicateRegex                           // Regular expression
	PredicateAny                             // Any character
	PredicateEOF                             // End of file
	PredicateLookahead                       // Positive/negative lookahead
	PredicateLookbehind                      // Positive/negative lookbehind
)

// CharPredicate - Single character predicate
type CharPredicate struct {
	Char rune
}

func (p *CharPredicate) Type() PredicateType { return PredicateChar }
func (p *CharPredicate) String() string      { return fmt.Sprintf("char('%c')", p.Char) }

// CharSetPredicate - Character set
type CharSetPredicate struct {
	Chars   string      // Literal characters
	Ranges  []CharRange // Ranges [a-z]
	Negated bool        // [^abc]
}

type CharRange struct {
	Start, End rune
}

// RegexPredicate - Compiled regular expression
type RegexPredicate struct {
	Pattern  string
	Compiled *regexp.Regexp // Kept for reference
	Simple   bool           // True if simple regex (no complex groups)
}

func (p *RegexPredicate) Type() PredicateType { return PredicateRegex }
func (p *RegexPredicate) String() string      { return fmt.Sprintf("regex:%s", p.Pattern) }
func (p *RegexPredicate) Equal(other Predicate) bool {
	if o, ok := other.(*RegexPredicate); ok {
		return p.Pattern == o.Pattern
	}
	return false
}

// AnyPredicate - Any character
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
	Positive  bool // true=positive, false=negative
}

// CompoundPredicate - Compound predicate (AND/OR)
type CompoundPredicate struct {
	Op    CompoundOp
	Preds []Predicate
}

type CompoundOp int

const (
	OpAnd CompoundOp = iota
	OpOr
)

// PredicateBuilder - Predicate constructor
type PredicateBuilder struct {
	strict bool
}

func NewPredicateBuilder(strict bool) *PredicateBuilder {
	return &PredicateBuilder{strict: strict}
}

// FromTextMateMatch - Converts TextMate regex to IR predicates
func (b *PredicateBuilder) FromTextMateMatch(pattern string) (Predicate, error) {
	if b.strict && containsUnsupportedRegex(pattern) {
		return nil, fmt.Errorf("unsupported regex pattern: %s", pattern)
	}

	// Simplify TextMate regex to basic predicates
	return b.simplifyRegex(pattern)
}

func (b *PredicateBuilder) simplifyRegex(pattern string) (Predicate, error) {
	// Convert to simpler predicates when possible
	switch {
	case pattern == ".":
		return &AnyPredicate{}, nil
	case isCharClass(pattern):
		return b.buildCharClassPredicate(pattern)
	case isLiteral(pattern):
		return b.buildLiteralPredicate(pattern)
	default:
		// Keep as generic regex
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
