package compiler

import (
	"fmt"
	"path/filepath"

	"github.com/ferchd/tm2hsl/internal/codegen"
	"github.com/ferchd/tm2hsl/internal/config"
	"github.com/ferchd/tm2hsl/internal/ir"
	"github.com/ferchd/tm2hsl/internal/optimizer"
	"github.com/ferchd/tm2hsl/internal/parser"
	"github.com/ferchd/tm2hsl/internal/serializer"
	"github.com/ferchd/tm2hsl/pkg/hsl"
)

type Compiler struct {
	config    *config.LanguageConfig
	grammar   interface{}
	irProgram *ir.Program
	bytecode  *hsl.Bytecode
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(configPath string) (*CompilationResult, error) {
	var err error

	// 1. Cargar configuración
	if err = c.loadConfig(configPath); err != nil {
		return nil, err
	}

	// 2. Parsear gramática TextMate
	if err = c.parseGrammar(); err != nil {
		return nil, err
	}

	// 3. Normalizar a máquina de estados
	if err = c.normalize(); err != nil {
		return nil, err
	}

	// 4. Convertir a IR
	if err = c.buildIR(); err != nil {
		return nil, err
	}

	// 5. Optimizar
	c.optimize()

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
		return fmt.Errorf("error cargando configuración: %w", err)
	}
	c.config = cfg
	return nil
}

func (c *Compiler) parseGrammar() error {
	p := parser.NewParser()

	grammarPath := c.config.GrammarPath()
	if grammarPath == "" {
		return fmt.Errorf("no se especificó gramática en configuración")
	}

	ext := filepath.Ext(grammarPath)
	switch ext {
	case ".json":
		grammar, err := p.ParseFile(grammarPath)
		if err != nil {
			return fmt.Errorf("error parseando gramática: %w", err)
		}
		c.grammar = grammar
	case ".tmLanguage", ".plist":
		return fmt.Errorf("formatos plist no soportados aún")
	default:
		return fmt.Errorf("formato de gramática no soportado: %s", ext)
	}

	return nil
}

func (c *Compiler) normalize() error {
	// Convertir gramática TextMate a máquina de estados
	// Esto es una simplificación - en realidad necesitaríamos
	// convertir la gramática parseada al tipo correcto
	return nil
}

func (c *Compiler) buildIR() error {
	c.irProgram = ir.NewProgram(c.config.Name, c.config.Scope)

	// Convertir máquina de estados a IR
	// Implementación simplificada
	return nil
}

func (c *Compiler) optimize() {
	opt := optimizer.NewOptimizer(c.irProgram)
	c.irProgram = opt.Optimize()
}

func (c *Compiler) generateBytecode() error {
	gen := codegen.NewGenerator(c.irProgram)
	bytecode, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("error generando bytecode: %w", err)
	}
	c.bytecode = bytecode
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
