package tester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/ferchd/tm2hsl/internal/compiler"
)

type Tester struct {
	compiler *compiler.Compiler
}

func NewTester() *Tester {
	return &Tester{
		compiler: compiler.NewCompiler(),
	}
}

type TestSpec struct {
	Name        string     `toml:"name"`
	Description string     `toml:"description,omitempty"`
	Cases       []TestCase `toml:"cases"`
}

type TestCase struct {
	Input    string   `toml:"input"`
	Expected []string `toml:"expected"`
	Line     int      `toml:"line,omitempty"`
	Column   int      `toml:"column,omitempty"`
}

type TestReport struct {
	Total    int
	Passed   int
	Failed   int
	Failures []TestFailure
}

type TestFailure struct {
	TestName string
	Input    string
	Expected string
	Actual   string
	Error    string
}

func (t *Tester) Run(configPath, specDir string) (*TestReport, error) {
	// Compilar el lenguaje primero
	result, err := t.compiler.Compile(configPath)
	if err != nil {
		return nil, fmt.Errorf("compilación fallida: %w", err)
	}

	// Cargar especificaciones de prueba
	specs, err := t.loadSpecs(specDir)
	if err != nil {
		return nil, fmt.Errorf("error cargando specs: %w", err)
	}

	// Ejecutar pruebas
	report := &TestReport{}

	for _, spec := range specs {
		for _, testCase := range spec.Cases {
			report.Total++

			// Ejecutar prueba
			actual, err := t.executeTest(result.Bytecode, testCase.Input)
			if err != nil {
				report.Failed++
				report.Failures = append(report.Failures, TestFailure{
					TestName: spec.Name,
					Input:    testCase.Input,
					Expected: strings.Join(testCase.Expected, ", "),
					Error:    err.Error(),
				})
				continue
			}

			// Comparar resultados
			if t.compareResults(actual, testCase.Expected) {
				report.Passed++
			} else {
				report.Failed++
				report.Failures = append(report.Failures, TestFailure{
					TestName: spec.Name,
					Input:    testCase.Input,
					Expected: strings.Join(testCase.Expected, ", "),
					Actual:   strings.Join(actual, ", "),
				})
			}
		}
	}

	return report, nil
}

func (t *Tester) loadSpecs(specDir string) ([]TestSpec, error) {
	files, err := ioutil.ReadDir(specDir)
	if err != nil {
		return nil, err
	}

	var specs []TestSpec

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".toml") {
			path := filepath.Join(specDir, file.Name())
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}

			var spec TestSpec
			if err := toml.Unmarshal(data, &spec); err != nil {
				return nil, fmt.Errorf("error parseando %s: %w", file.Name(), err)
			}

			specs = append(specs, spec)
		}
	}

	return specs, nil
}

func (t *Tester) executeTest(bytecode interface{}, input string) ([]string, error) {
	// Ejecutar el bytecode HSL contra la entrada
	// Esto requeriría un runtime HSL mínimo para pruebas
	// Por ahora, devolvemos un mock
	return []string{"source.test"}, nil
}

func (t *Tester) compareResults(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i := range actual {
		if actual[i] != expected[i] {
			return false
		}
	}

	return true
}
