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
		Config       string `arg:"" name:"config" help:"Path to language.toml"`
		Output       string `short:"o" help:"Output HSL file" default:"output.hsl"`
		ValidateOnly bool   `short:"v" help:"Only validate without generating bytecode"`
		Verbose      bool   `short:"V" help:"Enable verbose output"`
	} `cmd:"" help:"Compile a TextMate grammar to HSL bytecode"`

	Test struct {
		Config  string `arg:"" name:"config" help:"Path to language.toml"`
		SpecDir string `short:"s" help:"Directory with TOML test specs" default:"specs/"`
	} `cmd:"" help:"Run tokenization tests"`

	Version struct{} `cmd:"" help:"Show version"`
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

	cmp := compiler.NewCompiler()
	result, err := cmp.Compile(configPath)
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

	if c.Compile.Verbose {
		fmt.Printf("Compilation stats: %d regex, %d states, %d rules\n",
			result.Stats.RegexCount, result.Stats.StateCount, result.Stats.RuleCount)
	}
	fmt.Printf("HSL bytecode generated: %s\n", outputPath)

	return nil
}

func (c *CLI) RunTest(ctx *kong.Context) error {
	configPath, _ := filepath.Abs(c.Test.Config)
	specDir := c.Test.SpecDir

	tstr := tester.NewTester()
	report, err := tstr.Run(configPath, specDir)
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
