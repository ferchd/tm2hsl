package compiler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompiler_Compile_Integration(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "tm2hsl_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test language.toml config
	configContent := `
name = "Test Language"
scope = "source.test"

[grammar]
path = "test.tmLanguage.json"
`
	configPath := filepath.Join(tmpDir, "language.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create test TextMate grammar
	grammarContent := `{
		"scopeName": "source.test",
		"name": "Test Language",
		"patterns": [
			{
				"name": "keyword.control",
				"match": "\\b(if|else)\\b"
			}
		]
	}`
	grammarPath := filepath.Join(tmpDir, "test.tmLanguage.json")
	if err := os.WriteFile(grammarPath, []byte(grammarContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Run compilation
	compiler := NewCompiler()
	result, err := compiler.Compile(configPath)
	if err != nil {
		t.Errorf("Compiler.Compile() error = %v", err)
		return
	}

	if result == nil {
		t.Error("Compiler.Compile() returned nil result")
		return
	}

	if result.Bytecode == nil {
		t.Error("Compilation result has nil bytecode")
	}

	// Check stats
	if result.Stats.StateCount < 0 {
		t.Error("Invalid state count in stats")
	}
}

func TestCompiler_Compile_InvalidConfig(t *testing.T) {
	compiler := NewCompiler()

	// Test with non-existent config
	_, err := compiler.Compile("nonexistent.toml")
	if err == nil {
		t.Error("Compiler.Compile() should fail with invalid config")
	}
}

func TestCompiler_Compile_InvalidGrammar(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "tm2hsl_invalid_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create config pointing to invalid grammar
	configContent := `
name = "Test Language"
scope = "source.test"

[grammar]
path = "invalid.json"
`
	configPath := filepath.Join(tmpDir, "language.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create invalid grammar
	grammarContent := `{invalid json}`
	grammarPath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(grammarPath, []byte(grammarContent), 0644); err != nil {
		t.Fatal(err)
	}

	compiler := NewCompiler()
	_, err = compiler.Compile(configPath)
	if err == nil {
		t.Error("Compiler.Compile() should fail with invalid grammar")
	}
}
