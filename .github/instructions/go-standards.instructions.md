# Go Development Standards

## Go Version

- **Current version**: Go 1.23
- **Update policy**: Keep reasonably current with stable Go releases

## Project Layout

Follow the standard Go project layout:

```txt
/
├── main.go                 # Application entry point
├── cmd/                    # Command implementations
├── internal/               # Private application code
│   └── <package>/         # Internal packages
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
├── Makefile               # Build automation
└── .golangci.yml          # Linter configuration
```

## Coding Standards

### Formatting

- Use `gofmt` for all Go files (enforced by pre-commit hook)
- Use `goimports` for import organization (included in golangci-lint)

### Naming Conventions

- **Packages**: Short, lowercase, single-word names (e.g., `issue`, `cmd`)
- **Interfaces**: End with `-er` suffix when appropriate (e.g., `Reader`, `Writer`)
- **Variables**: Use camelCase
- **Constants**: Use MixedCaps or UPPER_CASE for exported constants
- **Exported names**: Start with uppercase letter
- **Unexported names**: Start with lowercase letter

### Error Handling

- Always check errors explicitly
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Return errors up the call stack; handle at the appropriate level
- Don't use panic for normal error handling

Example:

```go
content, err := os.ReadFile(filePath)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

### Testing

- **Location**: Tests go in `*_test.go` files alongside the code
- **Pattern**: Use table-driven tests for multiple test cases
- **Coverage**: Run `make coverage` to check coverage
- **Race detector**: Tests run with `-race` flag in CI

Example table-driven test:

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"case 1", "input1", "output1", false},
        {"error case", "bad", "", true},
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

## Dependencies

- **Minimize dependencies**: Only add well-maintained, necessary dependencies
- **Use go.mod**: Manage dependencies with Go modules
- **Version pinning**: Use specific versions in go.mod
- **Update regularly**: Keep dependencies updated for security

### Current Dependencies

- `github.com/spf13/cobra` - CLI framework (standard for Go CLI apps)

## Linting

### golangci-lint Configuration

The project uses `golangci-lint` with multiple linters enabled:

- **errcheck**: Unchecked errors
- **gosimple**: Code simplification
- **govet**: Suspicious constructs
- **staticcheck**: Static analysis
- **gofmt**: Code formatting
- **goimports**: Import organization
- **revive**: Comprehensive linting (golint replacement)
- **gosec**: Security issues
- **gocyclo**: Cyclomatic complexity (threshold: 15)
- **gocognit**: Cognitive complexity (threshold: 20)

### Running Linters

```bash
make lint              # Run all linters
make fmt               # Format code
make vet               # Run go vet
```

## Build Process

### Local Development

```bash
make build             # Build for current platform
make test              # Run tests
make coverage          # Generate coverage report
```

### Multi-Platform Builds

```bash
make build-all         # Build for:
                       # - linux/amd64, linux/arm64
                       # - darwin/amd64, darwin/arm64  
                       # - windows/amd64, windows/arm64
```

## Performance Considerations

- Use profiling for performance-critical code (`go test -cpuprofile`, `-memprofile`)
- Pre-allocate slices when size is known
- Use string builders for string concatenation
- Be mindful of memory allocations in hot paths

## Security

- Run `gosec` linter (included in golangci-lint)
- Never hard-code credentials
- Validate all external inputs
- Use `crypto/rand` for random numbers (not `math/rand`)
- Keep dependencies updated for security patches

## Documentation

- **Package comments**: Every package should have a package-level comment
- **Exported functions**: Document all exported functions, types, and constants
- **Examples**: Provide examples for complex functionality
- **README**: Keep README.md updated with usage instructions

## Git Hooks

The pre-commit hook runs:

1. Go formatting check (`gofmt`)
2. Go vet
3. golangci-lint
4. Go tests with race detector
5. Build verification

All checks must pass before commit.
