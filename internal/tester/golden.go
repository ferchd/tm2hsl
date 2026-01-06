package tester

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/compiler"
)

// RunGoldenTest runs a golden test using compiled grammar
func RunGoldenTest(configPath, input string) ([]TokenExpectation, error) {
	// Compile grammar
	cmp := compiler.NewCompiler()
	result, err := cmp.Compile(configPath)
	if err != nil {
		return nil, fmt.Errorf("compilation failed: %w", err)
	}

	// TODO: Use compiled bytecode to tokenize input
	// For now, return stub tokens
	tokens := []TokenExpectation{
		{Scope: "source.test", Text: input},
	}

	_ = result // Use result
	return tokens, nil
}

// CompareResults compares actual vs expected tokens
func (t *Tester) compareResults(actual, expected []TokenExpectation) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("token count mismatch: got %d, expected %d", len(actual), len(expected))
	}

	for i, exp := range expected {
		if i >= len(actual) {
			return fmt.Errorf("missing token at position %d", i)
		}
		act := actual[i]
		if act.Scope != exp.Scope || act.Text != exp.Text {
			return fmt.Errorf("token mismatch at %d: got %q/%q, expected %q/%q",
				i, act.Scope, act.Text, exp.Scope, exp.Text)
		}
	}

	return nil
}
