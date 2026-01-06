# tm2hsl

[![CI](https://github.com/ferchd/tm2hsl/workflows/CI/badge.svg)](https://github.com/ferchd/tm2hsl/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ferchd/tm2hsl)](https://goreportcard.com/report/github.com/ferchd/tm2hsl)
[![GoDoc](https://godoc.org/github.com/ferchd/tm2hsl?status.svg)](https://godoc.org/github.com/ferchd/tm2hsl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**tm2hsl** is a compiler that transforms TextMate syntax grammars into optimized HSL bytecode, revolutionizing language support in code editors.

## Vision

Modern editors (VSCode, Sublime Text, etc.) interpret complex TextMate grammars at runtime, causing latency, high memory usage, and limited scalability. **tm2hsl** changes this paradigm by compiling grammars once into deterministic bytecode that editors execute directly.

## Features

- **Ahead-of-time compilation**: Transforms TextMate grammars into optimized HSL bytecode
- **Deterministic execution**: Same input produces same bytecode
- **Massive scalability**: Efficient support for hundreds of languages
- **Layer separation**: Compiled languages vs. execution engine
- **Binary format**: Memory-mappable, versioned, and compact
- **Complete CLI**: Tools for development and testing

## Installation

### Pre-built Binaries

Download binaries from the [releases page](https://github.com/ferchd/tm2hsl/releases).

### From Source

```bash
# Clone the repository
git clone https://github.com/ferchd/tm2hsl.git
cd tm2hsl

# Setup development environment
make dev-setup

# Build the project
make build

# Install globally
make install
```

### Requirements

- Go 1.21 or higher
- Git

## Usage

### Basic Compilation

```bash
# Compile a TextMate grammar
tm2hsl compile --config language.toml --output output.hsl

# Validate without generating bytecode
tm2hsl compile --config language.toml --validate-only
```

### Configuration File

Create a `language.toml`:

```toml
name = "MyLanguage"
scope = "source.mylanguage"
grammar = "grammars/mylanguage.json"

[metadata]
version = "1.0.0"
description = "Support for MyLanguage"
```

## Architecture

```
tm2hsl
├── cmd/tm2hsl/          # Main CLI
├── internal/             # Private code
│   ├── cli/             # Command interface
│   ├── compiler/        # Compilation logic
│   ├── parser/          # TextMate parsing
│   ├── ir/              # Intermediate representation
│   ├── normalizer/      # Grammar normalization
│   ├── optimizer/       # Optimizations
│   ├── codegen/         # Bytecode generation
│   ├── serializer/      # HSL serialization
│   └── config/          # Configuration handling
├── pkg/                 # Public packages
│   ├── hsl/            # HSL bytecode format
│   └── textmate/       # TextMate types
└── docs/               # Documentation
```

### Compilation Flow

1. **Parsing**: Load and validate TextMate grammar (JSON/plist)
2. **Normalization**: Convert to deterministic state machine
3. **IR**: Generate optimized intermediate representation
4. **Optimization**: Apply structural transformations
5. **Bytecode**: Generate HSL binary bytecode
6. **Serialization**: Write final `.hsl` file

## HSL Format

HSL bytecode is a binary format designed for:

- **Sequential execution**: Efficient disk reading
- **Memory-mapping**: Zero-copy loading
- **Versioning**: Forward compatibility
- **Compression**: Optimized and deduplicated tables

### Bytecode Structure

```
HSL Header (64 bytes)
├── Magic: "HSL1"
├── Version: uint16
├── Checksum: uint32
└── Offset table...

String Table
Regex Table
State Table
Rule Table
Scope Table
```



## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guides.

### Quick Development

```bash
# Setup environment
./scripts/setup-dev.sh

# Iterative development
make          # build + test + lint
make build    # build only
make test     # tests only
```

### Conventional Commits

We use [Conventional Commits](https://conventionalcommits.org/) for messages:

```bash
feat: add support for recursive includes
fix: fix regex lookbehind parsing
docs: update HSL specification
```

## Documentation

- [HSL Specification](docs/HSL_SPEC.md) - Detailed bytecode format
- [Migration Guide](docs/MIGRATION.md) - From TextMate to HSL
- [Internal API](docs/API.md) - Developer reference
- [Examples](examples/) - Sample languages

## Project Status

**Current version**: 0.x (active development)

### Supported (v0)
- `match` rules with basic regex
- `begin`/`end` rules with content
- `contentName` for internal scopes
- `captures` with simple names
- Includes: `$self`, `$base`
- Line and block comments

### Not Supported (future)
- Repository with `#name` references
- Captures in `begin`/`end`
- `while` rules
- Complex back-references
- Advanced lookahead/lookbehind

## License

This project is licensed under Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [TextMate](https://macromates.com/) for the grammar format
- [VSCode](https://code.visualstudio.com/) for popularizing TextMate
- Open source community for inspiration and tools

## Contact

- **Issues**: [GitHub Issues](https://github.com/ferchd/tm2hsl/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ferchd/tm2hsl/discussions)
- **Email**: [fernando@example.com](mailto:fernando@example.com)

---

**tm2hsl**: Compiling languages, accelerating editors.