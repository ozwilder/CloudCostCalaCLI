# Copilot Instructions for CloudCostCalaCLI

## Project Overview

CloudCostCalaCLI is a Go-based command-line tool for calculating cloud infrastructure costs across AWS, Azure, and GCP. The project is in early stages with minimal established structure.

## Development Setup

### Go Environment
- **Language**: Go 1.21+ (check `go.mod` for minimum version)
- **Module**: Uses Go modules for dependency management (`go.mod`, `go.sum`)
- **Build Output**: Binary executable (platform-specific or cross-compiled)

### Initial Setup
```bash
# Initialize Go module (if not already done)
go mod init github.com/ozwilder/CloudCostCalaCLI

# Download dependencies
go mod download

# Verify dependencies
go mod tidy
```

## Project Structure (Recommended)

```
CloudCostCalaCLI/
├── .github/
│   └── copilot-instructions.md
├── cmd/
│   └── cloudcostcala/              # Main CLI application
│       └── main.go
├── internal/
│   ├── providers/                  # Cloud provider implementations
│   │   ├── aws/
│   │   │   ├── client.go
│   │   │   └── costs.go
│   │   ├── azure/
│   │   │   ├── client.go
│   │   │   └── costs.go
│   │   └── gcp/
│   │       ├── client.go
│   │       └── costs.go
│   ├── config/                     # Configuration file handling
│   │   ├── config.go
│   │   └── loader.go
│   ├── calculator/                 # Cost calculation logic
│   │   └── calculator.go
│   └── models/                     # Data structures
│       └── types.go
├── pkg/                             # Public/exported packages (if any)
│   └── output/                     # Output formatting (CSV, JSON, etc.)
│       └── formatter.go
├── tests/                           # Integration tests
│   └── *.go
├── config.example.json             # Example config file
├── go.mod
├── go.sum
├── Makefile                        # Build/test automation
├── README.md
└── LICENSE
```

## Build, Test & Validation

### Building
```bash
# Build for current platform
go build -o bin/cloudcostcala ./cmd/cloudcostcala

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/cloudcostcala-linux ./cmd/cloudcostcala
GOOS=darwin GOARCH=amd64 go build -o bin/cloudcostcala-macos ./cmd/cloudcostcala
GOOS=windows GOARCH=amd64 go build -o bin/cloudcostcala.exe ./cmd/cloudcostcala

# Using Makefile (when created)
make build
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestGetAWSCosts ./internal/providers/aws
```

### Linting & Code Quality
```bash
# Run golangci-lint (if configured)
golangci-lint run ./...

# Format code
go fmt ./...

# Check for common mistakes
go vet ./...

# Run tests with race detection
go test -race ./...
```

### Dependencies
```bash
# Add a dependency
go get github.com/user/package@latest

# Update all dependencies
go get -u ./...

# Clean up unused dependencies
go mod tidy
```

## Key Conventions

### Naming & Code Style
- **Packages**: Use lowercase, single-word names when possible (`providers`, `config`, `calculator`)
- **Functions**: Use CamelCase, exported functions start with uppercase (e.g., `GetAWSCosts`, `LoadConfig`)
- **Interfaces**: Name with `-er` suffix convention (e.g., `CostProvider`, `ConfigLoader`)
- **Constants**: Use UPPER_SNAKE_CASE for package-level constants
- **Error handling**: Always check and handle errors; avoid `panic()` except in initialization
- **Comments**: Export package documentation with `//` comments above exported identifiers

### Provider Interface Pattern
All cloud providers should implement a common interface:
```go
// internal/providers/provider.go
package providers

type CostProvider interface {
    GetCosts(ctx context.Context, opts Options) (*CostResult, error)
    Validate(ctx context.Context) error
}

// Implement in internal/providers/aws/client.go, azure/client.go, gcp/client.go
type AWSProvider struct {
    client *aws.Client
    config *Config
}

func (p *AWSProvider) GetCosts(ctx context.Context, opts Options) (*CostResult, error) {
    // Implementation
}
```

### Configuration Management
- Load config from file path specified by environment variable or flag
- Support JSON format for config files
- Use `encoding/json` for parsing
- Validate required fields on load
- Don't commit actual config files with real credentials to repository

### Error Handling
- Define custom error types in `internal/errors/errors.go` if needed
- Use `fmt.Errorf()` with wrapped errors: `fmt.Errorf("failed to load config: %w", err)`
- Log errors with context (e.g., which provider failed)

### Testing
- Test files: `*_test.go` in the same package as code being tested
- Use table-driven tests for multiple scenarios
- Mock external dependencies (cloud provider SDKs) for unit tests
- Use `testify/assert` or similar for cleaner assertions if desired
```go
// Example: internal/calculator/calculator_test.go
func TestCalculateTotal(t *testing.T) {
    tests := []struct {
        name     string
        input    Input
        expected float64
    }{
        {"single provider", Input{...}, 100.0},
        {"multi provider", Input{...}, 250.0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Calculate(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Authentication & Configuration

### Configuration File Approach
- **File Format**: JSON (defined in schema at `config.example.json`)
- **Default Locations**: 
  - Command line flag: `--config /path/to/config.json`
  - Environment variable: `CLOUDCOSTCALA_CONFIG`
  - Home directory: `~/.config/cloudcostcala/config.json`
- **Security**: 
  - Do NOT commit config files with real credentials
  - Add `config.json` to `.gitignore`
  - Document in README how to create config file
  - Consider future support for encrypted sensitive fields
- **Example Config**:
  ```json
  {
    "providers": {
      "aws": {
        "access_key_id": "***",
        "secret_access_key": "***",
        "region": "us-east-1"
      },
      "azure": {
        "subscription_id": "***",
        "client_id": "***",
        "client_secret": "***",
        "tenant_id": "***"
      },
      "gcp": {
        "project_id": "***",
        "service_account_key": "***"
      }
    }
  }
  ```

## Dependencies to Consider

- **AWS SDK**: `github.com/aws/aws-sdk-go-v2/...` (or v1)
- **Azure SDK**: `github.com/Azure/azure-sdk-for-go`
- **GCP SDK**: `cloud.google.com/go`
- **CLI Framework**: `github.com/spf13/cobra` or `github.com/urfave/cli` (if complex CLI needed)
- **Config**: `github.com/spf13/viper` (optional, for enhanced config management)
- **Logging**: `go.uber.org/zap` or standard `log` package
- **Testing**: `github.com/stretchr/testify` (optional, for assertions)

## Documentation

- **README.md**: Quick start, installation, basic usage examples
- **Code Comments**: Document exported functions and types above their declarations
- **Architecture Docs**: `docs/ARCHITECTURE.md` for high-level design (when applicable)
- **Configuration**: `config.example.json` with comments explaining each field

## CI/CD (When Ready)

When adding GitHub Actions:
- Run tests on multiple Go versions (1.21+, latest)
- Test on Linux, macOS, Windows
- Cross-compile binaries for common platforms
- Run linting (golangci-lint) before merge
- Check for race conditions: `go test -race ./...`
- Consider releasing binaries to GitHub Releases

## Notes

- This is a new project being rewritten in Go for better performance and portability
- Support all three cloud providers (AWS, Azure, GCP) from the start
- Cross-platform binary distribution is a strength of Go
- Use Go modules exclusively; no vendor directory needed unless required
- Configuration file-based authentication (no environment variable parsing from shell scripts)
