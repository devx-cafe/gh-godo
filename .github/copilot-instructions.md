# GitHub Copilot Instructions for gh-godo

## Project Overview

This repository contains `gh-godo`, a GitHub CLI extension written in Go. The extension provides utility commands for GitHub workflows and automation.

## Project Structure

- `main.go` - Entry point in the repository root
- `cmd/` - Command implementations using Cobra CLI framework
- `internal/` - Internal packages not meant for external use
- `Makefile` - Build automation and common tasks
- `.githooks/` - Git hooks including pre-commit checks
- `.github/workflows/` - CI/CD workflows

## Code Standards

### Go Programming

1. **Use Go 1.23+** - This project targets Go 1.23 or higher
2. **Follow standard Go conventions**:
   - Use `gofmt` for formatting
   - Run `go vet` for static analysis
   - Use `golangci-lint` for comprehensive linting
3. **Package organization**:
   - `cmd/` for command-line interface commands
   - `internal/` for internal application code
   - Each package should have focused responsibility

### CLI Design Principles

1. **Named flags over anonymous arguments**: Always use named flags (`--flag` or `-f`) instead of positional arguments
2. **Long and short forms**: Provide both long (`--file`) and short (`-f`) versions for commonly used flags
3. **Cobra framework**: Use the Cobra library for all CLI command implementations
4. **Help text**: Provide clear, concise help text for all commands and flags

### Testing

1. **Write table-driven tests**: Use Go's standard table-driven testing pattern
2. **Test coverage**: Aim for high test coverage on business logic
3. **Run tests before commit**: The pre-commit hook runs all tests automatically

### Build and Development

1. **Use Makefile targets**:
   - `make build` - Build the binary
   - `make test` - Run tests
   - `make lint` - Run linters
   - `make coverage` - Generate coverage report
   - `make build-all` - Build for all platforms

2. **Cross-platform builds**: Support Linux, Darwin (macOS), and Windows on AMD64 and ARM64

### Git Workflow

1. **Pre-commit hooks**: All commits are validated by `.githooks/pre-commit`
2. **Checks run before commit**:
   - Spell checking (cspell)
   - Markdown linting
   - Code formatting (prettier, gofmt)
   - Go vet
   - golangci-lint
   - Go tests
   - Build verification

## Making Changes

When modifying the code:

1. **Minimal changes**: Make the smallest possible change to achieve the goal
2. **Test first**: Run tests locally before committing
3. **Lint your code**: Ensure `make lint` passes
4. **Update tests**: Add or update tests for new functionality
5. **Document**: Update relevant documentation

## GH CLI Extension

This project is designed to be installed as a GitHub CLI extension:

```bash
gh extension install lakruzz/gh-godo
gh godo mkissue --file issue.md
```

The binary must be named `godo` and placed in the repository root to be recognized by the `gh` CLI.

## RAG Instructions Location

Additional instructions can be found in:

- `.github/copilot-instructions.md` (this file)
- `.github/instructions/*.instructions.md` - Specific topic instructions
