<!-- cspell:ignore gofmt golangci  -->

# Contributing to gh-godo

Thank you for your interest in contributing to gh-godo! This document provides guidelines and information for developers.

## Developer Resources

### Required Reading

For comprehensive development guidelines, please read the following RAG (Retrieval-Augmented Generation) files:

- [`.github/copilot-instructions.md`](.github/copilot-instructions.md) - Project overview and code standards
- [`.github/instructions/go-standards.instructions.md`](.github/instructions/go-standards.instructions.md) - Detailed Go development standards

These files contain essential information about:

- Project structure and organization
- Coding standards and conventions
- Testing practices
- Build processes
- Security considerations

## Prerequisites

- **Go 1.23 or higher** - This project targets Go 1.23+
- **Make** - For build automation
- **golangci-lint** - Can be installed via `make install-lint`
- **Git** - Configured with `.githooks` for pre-commit checks

## Getting Started

1. **Fork and Clone**

   ```bash
   git clone https://github.com/YOUR_USERNAME/gh-godo.git
   cd gh-godo
   ```

2. **Configure Git Hooks**

   ```bash
   git config core.hooksPath .githooks
   ```

3. **Install Dependencies**

   ```bash
   make deps
   ```

4. **Verify Setup**

   ```bash
   make build
   make test
   make lint
   ```

## Development Workflow

1. **Create a Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Follow the coding standards (see RAG files)
   - Write or update tests for your changes
   - Update documentation as needed

3. **Test Your Changes**

   ```bash
   # Run tests
   make test

   # Run tests with coverage
   make coverage

   # Run linter
   make lint

   # Format code
   make fmt

   # Run static analysis
   make vet
   ```

4. **Build and Verify**

   ```bash
   # Build for current platform
   make build

   # Test the binary
   ./godo --help
   ```

5. **Commit Your Changes**
   - The pre-commit hook will automatically run checks
   - Ensure all checks pass before pushing

6. **Submit a Pull Request**
   - Push your branch to your fork
   - Create a pull request with a clear description
   - Link any related issues

## Project Structure

```txt
.
├── main.go                 # Application entry point
├── cmd/                    # Command implementations
│   ├── root.go            # Root command definition
│   ├── mkissue.go         # mkissue command definition
│   └── mkissue/           # mkissue implementation
│       ├── mkissue.go     # Core logic
│       └── mkissue_test.go # Tests (alongside implementation)
├── exercises/              # Example files and templates
│   └── template.issue.md  # Issue file format contract
├── Makefile               # Build automation
├── .golangci.yml          # Linter configuration
├── .githooks/             # Git hooks
│   └── pre-commit         # Pre-commit checks
└── .github/               # GitHub configuration
    ├── copilot-instructions.md        # Project guidelines
    └── instructions/                   # Additional documentation
        └── go-standards.instructions.md
```

## Testing

### Test Organization

Tests follow the **Go community standard**: test files are placed alongside the code they test.

- Unit tests: `*_test.go` files in the same directory as the implementation
- Example: `cmd/mkissue/mkissue.go` has tests in `cmd/mkissue/mkissue_test.go`

This is the idiomatic Go approach and makes tests easy to find and maintain.

### Writing Tests

- Use **table-driven tests** for multiple test cases
- Test both success and error paths
- Include edge cases
- Aim for high coverage on business logic

Example table-driven test:

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "result", false},
        {"error case", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Building

### Local Development Build

```bash
make build
```

This creates the `godo` binary in the repository root.

### Multi-Platform Builds

```bash
make build-all
```

Builds for:

- Linux: amd64, arm64
- Darwin (macOS): amd64, arm64
- Windows: amd64, arm64

Binaries are placed in the `dist/` directory.

### Clean Build Artifacts

```bash
make clean
```

## Pre-commit Hooks

The repository uses pre-commit hooks to maintain code quality. When you commit, the following checks run automatically:

1. **Spell checking** (cspell)
2. **Markdown linting**
3. **Code formatting** (prettier, gofmt)
4. **Go vet** - Static analysis
5. **golangci-lint** - Comprehensive linting
6. **Go tests** - All unit tests
7. **Build verification** - Ensures code compiles

If any check fails, the commit is rejected. Fix the issues and try again.

## Coding Standards

### Go Programming

- **Follow Go idioms**: Use standard Go patterns and conventions
- **Format with gofmt**: Code must be formatted with `gofmt`
- **Use golangci-lint**: Address all linter warnings
- **Document exports**: All exported functions, types, and constants need documentation
- **Error handling**: Always check and wrap errors with context

### CLI Design

- **Named flags**: Use `--flag value` instead of positional arguments
- **Short and long forms**: Provide both (e.g., `-f` and `--file`)
- **Clear help text**: Every command and flag should have descriptive help
- **Cobra framework**: Use Cobra for all CLI commands

### Package Organization

- **`cmd/`**: CLI command definitions and implementations
- **`internal/`**: Internal packages (if needed, currently not used)
- **One purpose per package**: Keep packages focused and cohesive

## CI/CD

GitHub Actions workflows automatically:

- Run tests with race detector
- Check code formatting
- Run all linters
- Build binaries
- Verify pre-commit hooks pass

See `.github/workflows/` for workflow definitions.

## Issue File Format

The `mkissue` command uses a specific format defined in [`exercises/template.issue.md`](exercises/template.issue.md). This is the **contract** for issue files:

- YAML frontmatter with metadata (title, assignees, labels, etc.)
- Markdown body for issue content

When making changes to issue parsing, ensure compatibility with this format.

## Common Tasks

### Add a New Command

1. Create command file in `cmd/` (e.g., `cmd/newcommand.go`)
2. Create implementation package if needed (e.g., `cmd/newcommand/newcommand.go`)
3. Add tests alongside implementation
4. Register command in `cmd/root.go` init function
5. Update README.md with usage documentation
6. Run tests and linters

### Update Dependencies

```bash
# Update go.mod and go.sum
go get -u ./...
go mod tidy

# Verify everything still works
make test
make build
```

### Debug Build Issues

```bash
# Clean everything
make clean

# Rebuild from scratch
make build

# Run with verbose output
go build -v -o ./godo .
```

## Getting Help

- Check the [RAG files](.github/copilot-instructions.md) first
- Look at existing code for examples
- Review closed pull requests for similar changes
- Open an issue if you have questions

## Code Review Process

When you submit a pull request:

1. Automated checks must pass (CI/CD)
2. Code will be reviewed by maintainers
3. Address any feedback
4. Once approved, your PR will be merged

## License

By contributing, you agree that your contributions will be licensed under the same terms as the project.
