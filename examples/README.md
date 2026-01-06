# Example Language: SimpleScript

This directory contains an example of compiling a simple programming language grammar using tm2hsl.

## Files

- `grammar.json`: TextMate grammar for SimpleScript
- `language.toml`: Configuration file
- `specs/`: Test specifications

## Compilation

```bash
tm2hsl compile language.toml -o simplescript.hsl
```

## Grammar Features

- Keywords: `function`, `if`, `else`
- Identifiers: `[a-zA-Z_][a-zA-Z0-9_]*`
- Numbers: `\d+`
- Strings: `"[^"]*"`
- Comments: `//.*`

## Test

```bash
tm2hsl test language.toml --spec-dir specs/
```