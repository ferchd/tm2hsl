package compiler

import (
	"fmt"

	"github.com/ferchd/tm2hsl/internal/config"
	"github.com/ferchd/tm2hsl/internal/ir"
	"github.com/ferchd/tm2hsl/internal/normalizer"
	"github.com/ferchd/tm2hsl/internal/optimizer"
	"github.com/ferchd/tm2hsl/internal/parser"
	"github.com/ferchd/tm2hsl/internal/serializer"
	"github.com/ferchd/tm2hsl/pkg/hsl"
)

type Compiler struct {
	config       *config.LanguageConfig
	grammar      *parser.TextMateAST
	stateMachine *ir.StateMachine
	irProgram    *ir.Program
	bytecode     *hsl.Bytecode
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(configPath string) (*CompilationResult, error) {
	var err error

	// 1. Load configuration
	if err = c.loadConfig(configPath); err != nil {
		return nil, err
	}

	// 2. Parse TextMate grammar
	if err = c.parseGrammar(); err != nil {
		return nil, err
	}

	// 3. Normalize to state machine
	if err = c.normalize(); err != nil {
		return nil, err
	}

	// 4. Convert to IR
	if err = c.buildIR(); err != nil {
		return nil, err
	}

	// 5. Optimize
	if err = c.optimize(); err != nil {
		return nil, err
	}

	// 6. Generar bytecode
	if err = c.generateBytecode(); err != nil {
		return nil, err
	}

	return &CompilationResult{
		Bytecode: c.bytecode,
		Stats:    c.irProgram.Statistics(),
	}, nil
}

func (c *Compiler) loadConfig(path string) error {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}
	c.config = cfg
	return nil
}

func (c *Compiler) parseGrammar() error {
	// TODO: Implement grammar parsing
	grammarPath := c.config.GrammarPath()
	if grammarPath == "" {
		return fmt.Errorf("no grammar specified in configuration")
	}

	// Stub implementation - create a minimal TextMateAST
	c.grammar = &parser.TextMateAST{
		ScopeName: "source.test",
		Name:      "Test Language",
		Patterns: []parser.GrammarRule{
			{
				Match: "function",
				Name:  "keyword.function",
			},
		},
		Repository: make(map[string]parser.GrammarRule),
	}
	return nil
}

func (c *Compiler) normalize() error {
	norm := normalizer.NewNormalizer()
	machine, err := norm.Normalize(c.grammar)
	if err != nil {
		return fmt.Errorf("normalization failed: %w", err)
	}
	c.stateMachine = machine
	return nil
}

func (c *Compiler) buildIR() error {
	c.irProgram = ir.NewProgram(c.config.Name, c.config.Scope)

	// Convert state machine to IR
	// Simplified implementation
	return nil
}

func (c *Compiler) optimize() error {
	opt := optimizer.NewOptimizer()
	optimized, err := opt.Optimize(c.irProgram)
	if err != nil {
		return fmt.Errorf("optimization failed: %w", err)
	}
	c.irProgram = optimized
	return nil
}

func (c *Compiler) generateBytecode() error {
	ser := serializer.NewSerializer()
	c.bytecode = ser.ConvertToBytecode(c.stateMachine)
	return nil
}

type CompilationResult struct {
	Bytecode *hsl.Bytecode
	Stats    ir.ProgramStats
}

func (r *CompilationResult) WriteToFile(path string) error {
	ser := serializer.NewSerializer()
	return ser.WriteToFile(r.Bytecode, path)
}
