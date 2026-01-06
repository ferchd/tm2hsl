package tester

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// TestCase represents a single tokenization test
type TestCase struct {
	Name     string             `toml:"name"`
	Input    string             `toml:"input"`
	Expected []TokenExpectation `toml:"expected"`
}

// TokenExpectation defines expected token
type TokenExpectation struct {
	Scope string `toml:"scope"`
	Text  string `toml:"text"`
}

// Tester runs tokenization tests
type Tester struct{}

// NewTester creates a new tester
func NewTester() *Tester {
	return &Tester{}
}

// TestReport contains test results
type TestReport struct {
	Passed   int
	Failed   int
	Failures []TestFailure
}

// TestFailure describes a failed test
type TestFailure struct {
	TestName string
	Error    string
}

// Run executes tests from spec directory
func (t *Tester) Run(configPath, specDir string) (*TestReport, error) {
	files, err := filepath.Glob(filepath.Join(specDir, "*.toml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list spec files: %w", err)
	}

	report := &TestReport{}
	for _, file := range files {
		if err := t.runSpecFile(file, report); err != nil {
			return nil, fmt.Errorf("failed to run spec %s: %w", file, err)
		}
	}

	return report, nil
}

func (t *Tester) runSpecFile(path string, report *TestReport) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var specs struct {
		Cases []TestCase `toml:"cases"`
	}

	if err := toml.Unmarshal(data, &specs); err != nil {
		return fmt.Errorf("failed to parse TOML: %w", err)
	}

	for _, tc := range specs.Cases {
		if err := t.runTestCase(tc, report); err != nil {
			report.Failed++
			report.Failures = append(report.Failures, TestFailure{
				TestName: tc.Name,
				Error:    err.Error(),
			})
		} else {
			report.Passed++
		}
	}

	return nil
}

func (t *Tester) runTestCase(tc TestCase, report *TestReport) error {
	// Run tokenization
	actual, err := RunGoldenTest("dummy-config", tc.Input) // TODO: Pass real config
	if err != nil {
		return fmt.Errorf("tokenization failed: %w", err)
	}

	// Compare results
	return t.compareResults(actual, tc.Expected)
}
