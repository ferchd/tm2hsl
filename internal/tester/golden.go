// golden.go - Pruebas de referencia
package tester

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestCase - Entrada completa para golden test
type TestCase struct {
	Name     string
	Grammar  string // Archivo .tmLanguage.json
	Input    string // Texto a tokenizar
	Expected []TokenExpectation
	Skip     bool `json:"skip,omitempty"`
}

type TokenExpectation struct {
	Text   string   // Substring del input
	Scopes []string // Jerarquía de scopes esperada
	Line   int      // Línea (1-indexed)
	Column int      // Columna (0-indexed)
}

// RunGoldenTest - Ejecuta prueba de referencia
func RunGoldenTest(t *testing.T, testDir string, tc TestCase) {
	t.Helper()

	// 1. Cargar gramática
	grammarPath := filepath.Join(testDir, tc.Grammar)
	grammar, err := parser.LoadGrammarFile(grammarPath)
	if err != nil {
		t.Fatalf("failed to load grammar: %v", err)
	}

	// 2. Normalizar a IR
	normalizer := normalizer.NewNormalizer()
	machine, err := normalizer.Normalize(grammar)
	if err != nil {
		t.Fatalf("failed to normalize: %v", err)
	}

	// 3. Optimizar
	optimizer := optimizer.NewOptimizer()
	optimized, err := optimizer.Optimize(machine)
	if err != nil {
		t.Fatalf("failed to optimize: %v", err)
	}

	// 4. Serializar a bytecode
	bytecode, err := serializer.Serialize(optimized)
	if err != nil {
		t.Fatalf("failed to serialize: %v", err)
	}

	// 5. Ejecutar en máquina virtual de prueba
	vm := virtualmachine.NewVM(bytecode)
	tokens, err := vm.Tokenize(tc.Input)
	if err != nil {
		t.Fatalf("failed to tokenize: %v", err)
	}

	// 6. Comparar con resultado esperado
	compareTokens(t, tokens, tc.Expected)
}

// compareTokens - Comparación estricta
func compareTokens(t *testing.T, actual []Token, expected []TokenExpectation) {
	if len(actual) != len(expected) {
		t.Errorf("token count mismatch: got %d, want %d", len(actual), len(expected))
	}

	for i := 0; i < min(len(actual), len(expected)); i++ {
		act := actual[i]
		exp := expected[i]

		if act.Text != exp.Text {
			t.Errorf("token[%d].text: got %q, want %q", i, act.Text, exp.Text)
		}

		if !equalScopes(act.Scopes, exp.Scopes) {
			t.Errorf("token[%d].scopes: got %v, want %v", i, act.Scopes, exp.Scopes)
		}
	}
}
