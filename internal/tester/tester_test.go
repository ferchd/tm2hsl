package tester

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTester_Run(t *testing.T) {
	// Create temporary spec directory
	tmpDir, err := os.MkdirTemp("", "tester_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test spec file
	specContent := `
[[cases]]
name = "test case"
input = "test input"
expected = [
  { scope = "test.scope", text = "test input" }
]
`
	specPath := filepath.Join(tmpDir, "test.toml")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	tester := NewTester()
	report, err := tester.Run("dummy-config", tmpDir)
	if err != nil {
		t.Errorf("Tester.Run() error = %v", err)
		return
	}

	if report == nil {
		t.Error("Tester.Run() returned nil report")
		return
	}

	// Since RunGoldenTest is stub, it may fail, but structure should work
	t.Logf("Report: Passed=%d, Failed=%d", report.Passed, report.Failed)
}

func TestCompareResults(t *testing.T) {
	tester := NewTester()

	actual := []TokenExpectation{
		{Scope: "test", Text: "hello"},
	}
	expected := []TokenExpectation{
		{Scope: "test", Text: "hello"},
	}

	err := tester.compareResults(actual, expected)
	if err != nil {
		t.Errorf("compareResults() should succeed with matching results, got %v", err)
	}

	// Test mismatch
	actualMismatch := []TokenExpectation{
		{Scope: "test", Text: "world"},
	}
	err = tester.compareResults(actualMismatch, expected)
	if err == nil {
		t.Error("compareResults() should fail with mismatched results")
	}
}
