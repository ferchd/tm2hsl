package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/alecthomas/kong"

	"github.com/ferchd/tm2hsl/internal/compiler"
	"github.com/ferchd/tm2hsl/internal/tester"
)

type CLI struct {
	Compile struct {
		Config       string `arg:"" name:"config" help:"Ruta al archivo language.toml"`
		Output       string `short:"o" help:"Archivo de salida HSL" default:"output.hsl"`
		ValidateOnly bool   `short:"v" help:"Only validate without generating bytecode"`
	} `cmd:"" help:"Compilar una gramática"`

	Test struct {
		Config  string `arg:"" name:"config" help:"Ruta al archivo language.toml"`
		SpecDir string `short:"s" help:"Directorio de especificaciones" default:"specs/"`
	} `cmd:"" help:"Run tests"`

	Version struct{} `cmd:"" help:"Mostrar versión"`
}

var version = "0.0.1-alpha"

func Execute() error {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("tm2hsl"),
		kong.Description("Compilador de lenguajes léxicos TextMate a HSL"),
		kong.UsageOnError())

	return ctx.Run(context.Background())
}

func (c *CLI) RunCompile(ctx *kong.Context) error {
	configPath, _ := filepath.Abs(c.Compile.Config)

	compiler := compiler.NewCompiler()
	result, err := compiler.Compile(configPath)
	if err != nil {
		return fmt.Errorf("compilation error: %w", err)
	}

	if c.Compile.ValidateOnly {
		fmt.Println("Grammar validated successfully")
		return nil
	}

	outputPath := c.Compile.Output
	if err := result.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("error writing bytecode: %w", err)
	}

	fmt.Printf("HSL bytecode generated: %s\n", outputPath)

	return nil
}

func (c *CLI) RunTest(ctx *kong.Context) error {
	configPath, _ := filepath.Abs(c.Test.Config)
	specDir := c.Test.SpecDir

	tester := tester.NewTester()
	report, err := tester.Run(configPath, specDir)
	if err != nil {
		return err
	}

	fmt.Printf("Test results: %d passed, %d failed\n",
		report.Passed, report.Failed)

	if report.Failed > 0 {
		for _, failure := range report.Failures {
			fmt.Printf("FAILED %s: %s\n", failure.TestName, failure.Error)
		}
		return fmt.Errorf("tests failed")
	}

	return nil
}

func (c *CLI) RunVersion(ctx *kong.Context) error {
	fmt.Printf("tm2hsl v%s\n", version)
	return nil
}
