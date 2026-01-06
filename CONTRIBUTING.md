# Contributor Guide

Welcome! We're thrilled that you want to contribute to **tm2hsl**. This guide will help you understand how to contribute effectively to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Environment Setup](#environment-setup)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Pull Requests](#pull-requests)
- [Issues](#issues)

## Code of Conduct

This project follows a code of conduct to ensure the community is inclusive and respectful. By participating, you agree to:

- Be respectful to all people
- Use inclusive language
- Accept constructive criticism
- Focus on what's best for the community
- Show empathy towards other members

## How to Contribute

### Types of Contributions

- **Report bugs**: Use GitHub Issues with clear details
- **Propose features**: Discuss ideas in GitHub Discussions
- **Improve documentation**: Fix errors or add examples
- **Write code**: Implement features or fixes
- **Add tests**: Improve coverage and edge cases
- **Performance**: Optimize algorithms or memory

### First Steps

1. **Fork** the repository
2. **Clone** your fork locally
3. Create a **branch** with a descriptive name
4. **Implement** your changes
5. **Test** thoroughly
6. **Commit** following conventional commits
7. **Push** and create a **Pull Request**

## Environment Setup

### Requirements

- Go 1.21+
- Git
- (Optional) Docker for advanced testing

### Automatic Setup

```bash
# Clone and setup
git clone https://github.com/your-username/tm2hsl.git
cd tm2hsl

# Automatic setup
./scripts/setup-dev.sh
```

### Manual Setup

```bash
# Dependencies
go mod download

# Development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Verify installation
make ci
```

## Development Workflow

### Branches

We use a simple branching model:

```
main                    # Main branch (releases)
├── feature/new-functionality
├── fix/bug-fix
├── docs/documentation-improvement
└── refactor/restructuring
```

### Commits

We follow [Conventional Commits](https://conventionalcommits.org/):

```bash
# Format
<type>[optional scope]: <description>

# Examples
feat(parser): add support for recursive includes
fix(compiler): fix memory leak in optimizer
docs(api): update usage examples
refactor(ir): simplify state representation

# Allowed types
feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
```

### Atomic Commits

Each commit should be:
- **Independent**: Works by itself
- **Small**: Focused changes
- **Tested**: Includes tests if applicable
- **Documented**: Updates docs if necessary

## Code Standards

### Go Guidelines

- **Formatting**: `go fmt` automatically
- **Linting**: `golangci-lint` passes without warnings
- **Imports**: Grouped (stdlib → third-party → internal)
- **Naming**: PascalCase for exported, camelCase for private
- **Errors**: Wrapped with `fmt.Errorf("msg: %w", err)`
- **Comments**: In Spanish for exported functions

### File Structure

```
tm2hsl/
├── cmd/           # CLI applications
├── internal/      # Private code
│   └── feature/   # One feature per directory
├── pkg/           # Public packages
├── docs/          # Documentation
├── examples/      # Examples
├── scripts/       # Development scripts
└── testdata/      # Test data
```

### Interfaces

- **Small**: Maximum 1-3 methods
- **Clearly defined**: Descriptive names
- **Documented**: Comments for each method

## Testing

### Types of Tests

- **Unit tests**: `*_test.go` alongside code
- **Integration tests**: `integration/` directory
- **Golden tests**: For output comparison
- **Benchmark tests**: Critical performance

### Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Specific tests
go test ./internal/parser -v

# Benchmarks
go test -bench=. ./...

# Race detection
go test -race ./...
```

### Writing Tests

```go
func TestParser_ParseGrammar(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        want     *Grammar
        wantErr  bool
    }{
        {
            name:  "valid grammar",
            input: `{"scopeName": "test"}`,
            want:  &Grammar{ScopeName: "test"},
        },
        // more cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseGrammar(strings.NewReader(tt.input))
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### Coverage

Maintain >80% coverage for new code. Run:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Pull Requests

### Pre-Submit Checklist

- [ ] `make ci` passes (build, test, lint)
- [ ] New tests for added functionality
- [ ] Documentation updated
- [ ] Commits follow conventional commits
- [ ] No breaking changes without prior discussion
- [ ] Performance not degraded

### PR Template

```markdown
## Description
[Brief description of changes]

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation
- [ ] Code/style
- [ ] Performance
- [ ] Build/DevOps

## Testing
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual tests
- [ ] Coverage maintained

## Checklist
- [ ] Code reviewed
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No merge conflicts
```

### Review Process

1. **Automated checks**: CI must pass
2. **Code review**: At least 1 approval
3. **Testing**: Reviewer runs tests locally
4. **Merge**: Squash merge with descriptive message

## Issues

### Bug Reports

Use the bug template with:

- **Clear title**: "Parser fails with specific regex"
- **Description**: Steps to reproduce
- **Environment**: Go version, OS, tm2hsl version
- **Logs**: Relevant output
- **Code**: Minimal reproducible example

### Feature Requests

- **Research first**: Search for similar issues
- **Describe clearly**: What, why, how
- **Consider alternatives**: Other possible solutions
- **Accept feedback**: May not be implemented

## Contribution Areas

### High Impact
- **Performance**: Critical optimizations
- **Testing**: More coverage and edge cases
- **Documentation**: Examples and guides
- **Tooling**: Development scripts

### Beginner-Friendly
- **Bug fixes**: Issues labeled `good first issue`
- **Documentation**: Typo corrections
- **Tests**: Add missing cases
- **Examples**: Sample languages

## Communication

- **Issues**: For bugs and features
- **Discussions**: For general questions
- **Discord/Slack**: For real-time chat (if available)

## Recognition

Contributors are recognized in:
- Contributor list in README
- Release notes
- Possible mention in documentation

Thank you for contributing to tm2hsl!