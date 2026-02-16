# Copilot Instructions for CloudCostCalaCLI

## Project Overview

CloudCostCalaCLI is a Go-based CLI tool that scans and discovers all cloud assets across AWS, Azure, and GCP cloud environments. It aggregates assets across projects/subscriptions/resource groups/tenants and outputs an Excel file with:
- Each asset type as one row
- Instance count for each asset type
- **Synthetic Units**: A standardized credit-point system for budget planning and capacity allocation

The tool translates real cloud infrastructure into abstract "credit units" to help customers plan budgets and allocate resources.

## Synthetic Units Conversion Rules

| Asset Type | Conversion Rule |
|------------|-----------------|
| **VM/Server** | 1 VM = 5 units (regardless of size/type) |
| **Containers** | 4 vCores = 2 units (minimum 2 units, cap at 16 vCores) |
| **Storage** | Per storage unit = 5 units |
| **Database** (SQL, NoSQL, etc.) | 1 DB = 5 units |
| **Serverless** (Functions, Lambda, etc.) | Per 10 functions = 5 units |

## Development Setup

### Go Environment
- **Language**: Go 1.21+ (check `go.mod` for minimum version)
- **Module**: Uses Go modules for dependency management (`go.mod`, `go.sum`)
- **Build Output**: Cross-platform binary executable

### Initial Setup
```bash
# Initialize Go module
go mod init github.com/ozwilder/CloudCostCalaCLI

# Download dependencies
go mod download

# Verify dependencies
go mod tidy
```

## Project Structure

```
CloudCostCalaCLI/
├── .github/
│   └── copilot-instructions.md
├── cmd/
│   └── cloudcostcala/
│       └── main.go                 # CLI entry point with flags/config handling
├── internal/
│   ├── providers/                  # Cloud provider implementations
│   │   ├── aws/
│   │   │   ├── client.go          # AWS API client setup
│   │   │   ├── scanner.go         # Discover EC2, RDS, S3, Lambda, etc.
│   │   │   └── mapper.go          # Map AWS assets to internal model
│   │   ├── azure/
│   │   │   ├── client.go
│   │   │   ├── scanner.go         # Discover VMs, SQL, containers, etc.
│   │   │   └── mapper.go
│   │   └── gcp/
│   │       ├── client.go
│   │       ├── scanner.go         # Discover Compute, Cloud SQL, etc.
│   │       └── mapper.go
│   ├── assets/                     # Asset model & aggregation
│   │   ├── types.go               # Asset types (VM, DB, Container, etc.)
│   │   ├── aggregator.go          # Combine assets across projects/subscriptions
│   │   └── converter.go           # Convert assets to synthetic units
│   ├── config/
│   │   ├── config.go
│   │   └── loader.go
│   └── models/
│       ├── asset.go               # Core asset struct
│       └── synthetic_unit.go      # Synthetic unit calculation
├── pkg/
│   ├── output/
│   │   ├── excel.go              # Excel file generation (xls/xlsx)
│   │   └── formatter.go          # Format assets for output
│   └── discovery/
│       └── scanner.go            # Discovery orchestration across providers
├── config.example.json            # Example config with cloud credentials
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

## Build, Test & Validation

### Building
```bash
# Build for current platform
go build -o bin/cloudcostcala ./cmd/cloudcostcala

# Cross-compile for all platforms
make build-all  # (defined in Makefile)

# Build with specific platform
GOOS=linux GOARCH=amd64 go build -o bin/cloudcostcala-linux ./cmd/cloudcostcala
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test (e.g., synthetic unit conversion)
go test -run TestSyntheticUnitConversion ./internal/assets

# Run race detection
go test -race ./...
```

### Linting & Code Quality
```bash
# Format code
go fmt ./...

# Run linter (if configured)
golangci-lint run ./...

# Check for issues
go vet ./...
```

## Key Conventions

### Naming & Code Style
- **Packages**: lowercase, single-word names (`providers`, `assets`, `output`, `discovery`)
- **Exported Functions**: CamelCase starting with uppercase (e.g., `ScanAWSAssets`, `ConvertToSyntheticUnits`)
- **Interfaces**: `-er` suffix (e.g., `CloudProvider`, `AssetScanner`, `OutputFormatter`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `UNITS_PER_VM`, `MIN_CONTAINER_UNITS`)
- **Error Handling**: Check all errors; wrap with context using `fmt.Errorf("operation: %w", err)`
- **Comments**: Document exported types and functions with `//` comment above declaration

### Asset Discovery Pattern

Each cloud provider implements the `CloudProvider` interface:

```go
// internal/providers/provider.go
package providers

type CloudProvider interface {
    // Discover all assets in all projects/subscriptions
    DiscoverAssets(ctx context.Context) ([]Asset, error)
    
    // Validate credentials
    Validate(ctx context.Context) error
}

// Asset is the common internal representation
type Asset struct {
    ID              string     // Unique asset ID
    Type            string     // VM, Database, Container, Storage, Function
    Name            string
    Cloud           string     // AWS, Azure, GCP
    Project         string     // Project ID, Subscription ID, etc.
    Count           int        // Number of instances/units
    Metadata        map[string]interface{} // vCores, storage size, etc.
}
```

### Synthetic Unit Conversion

```go
// internal/assets/converter.go
package assets

func ConvertToSyntheticUnits(asset Asset) int {
    switch asset.Type {
    case "VM":
        return 5 * asset.Count
    case "Container":
        vCores := asset.Metadata["vCores"].(int)
        return int(math.Max(2, float64(vCores/4))) * asset.Count
    case "Database":
        return 5 * asset.Count
    case "Storage":
        return 5 * asset.Count
    case "Function":
        // Per 10 functions = 5 units
        return (asset.Count + 9) / 10 * 5
    default:
        return 0
    }
}
```

### Asset Aggregation

```go
// internal/assets/aggregator.go
package assets

// AggregateAssets combines assets across all projects/subscriptions/resource groups
// Groups by asset type and sums instances
func AggregateAssets(assets []Asset) map[string]AggregatedAsset {
    // Group by type
    // Sum instances per type
    // Calculate total synthetic units
}

type AggregatedAsset struct {
    Type           string
    TotalInstances int
    TotalUnits     int
    Breakdown      map[string]int // Optional: instances by cloud/project
}
```

### Excel Output

Use `github.com/xuri/excelize` for Excel generation:

```go
// pkg/output/excel.go
package output

func WriteExcel(filename string, assets []AggregatedAsset) error {
    // Create columns: Asset Type | Instance Count | Synthetic Units
    // Write aggregated assets
    // Format and save .xlsx file
}
```

## Configuration Management

### Configuration File (JSON)
```json
{
  "providers": {
    "aws": {
      "enabled": true,
      "access_key_id": "***",
      "secret_access_key": "***",
      "regions": ["us-east-1", "eu-west-1"]
    },
    "azure": {
      "enabled": true,
      "subscription_id": "***",
      "client_id": "***",
      "client_secret": "***",
      "tenant_id": "***"
    },
    "gcp": {
      "enabled": true,
      "project_id": "***",
      "service_account_key": "***"
    }
  },
  "output": {
    "format": "excel",
    "filename": "cloud-assets-inventory.xlsx"
  }
}
```

### CLI Flags
```bash
cloudcostcala --config config.json --output assets.xlsx
cloudcostcala --aws-key *** --azure-subscription *** --output assets.xlsx
```

## Dependencies to Consider

- **AWS SDK**: `github.com/aws/aws-sdk-go-v2`
- **Azure SDK**: `github.com/Azure/azure-sdk-for-go`
- **GCP SDK**: `cloud.google.com/go`
- **CLI Framework**: `github.com/spf13/cobra` (for complex CLI)
- **Excel**: `github.com/xuri/excelize` (for .xlsx generation)
- **Config**: `github.com/spf13/viper` (for config file management)
- **Logging**: `go.uber.org/zap` or standard `log`
- **Testing**: `github.com/stretchr/testify` (assertions)

## Testing Strategy

- **Unit Tests**: Test synthetic unit conversion with different asset types
- **Integration Tests**: Mock cloud provider APIs to test discovery flow
- **Mock Providers**: Create mock implementations of CloudProvider for testing

```go
// Example test for synthetic unit conversion
func TestSyntheticUnitConversion(t *testing.T) {
    tests := []struct {
        asset    Asset
        expected int
    }{
        {Asset{Type: "VM", Count: 1}, 5},
        {Asset{Type: "VM", Count: 3}, 15},
        {Asset{Type: "Container", Count: 1, Metadata: map[string]interface{}{"vCores": 4}}, 2},
        {Asset{Type: "Function", Count: 10}, 5},
        {Asset{Type: "Function", Count: 15}, 10},
    }
    
    for _, tt := range tests {
        result := ConvertToSyntheticUnits(tt.asset)
        assert.Equal(t, tt.expected, result)
    }
}
```

## Documentation

- **README.md**: Quick start, installation, basic usage examples with output samples
- **config.example.json**: Well-commented configuration template
- **docs/ARCHITECTURE.md**: Asset discovery flow, aggregation strategy
- **docs/SYNTHETIC_UNITS.md**: Detailed explanation of synthetic unit model
- **Code Comments**: Document asset types, conversion logic, provider implementations

## Workflow Overview

```
User provides config (cloud credentials)
  ↓
CLI discovers assets across all clouds/projects/subscriptions
  ↓
Assets aggregated by type (VM, Database, Container, etc.)
  ↓
Each asset type converted to synthetic units
  ↓
Grouped results written to Excel file
  ↓
Excel output: [Asset Type | Instance Count | Synthetic Units]
```

## CI/CD (When Ready)

- Test on Linux, macOS, Windows
- Test against mock cloud provider APIs
- Validate Excel file generation
- Cross-compile binaries for distribution

## Notes

- Asset discovery across multi-cloud (AWS, Azure, GCP) requires handling different APIs and structures
- Aggregation must be idempotent (same results with same input)
- Synthetic unit calculations must be consistent and clearly documented
- Excel file should be human-readable with proper formatting
- Handle large environments efficiently (pagination, concurrent API calls where possible)
