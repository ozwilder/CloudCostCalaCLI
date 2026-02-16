# Copilot Instructions for CloudCostCalaCLI

## Project Overview

CloudCostCalaCLI is a Go-based CLI tool that scans and discovers all cloud assets across AWS, Azure, and GCP cloud environments and enriches the inventory with actual billing data for accurate budget planning.

**Dual Input Approach:**
1. **Current Asset Inventory**: Real-time scan of live cloud resources (VMs, databases, containers, etc.)
2. **Cloud Billing Files**: Historical consumption data (CSV format) capturing ephemeral resources and actual usage patterns

The tool aggregates both data sources across projects/subscriptions/resource groups and outputs an Excel file with:
- Each asset type as one row
- Instance count from current inventory
- **Real-life consumption metrics** from billing data
- **Synthetic Units**: Standardized credit-points based on actual usage for accurate budget planning

This hybrid approach captures both persistent and ephemeral resources, providing realistic budget estimates based on real cloud spending patterns.

## Synthetic Units Model

Synthetic units are customer-specific conversion rules that translate cloud assets into abstract credit points based on actual consumption.

**Default Model** (can be customized in config):

| Asset Type | Conversion Rule |
|------------|-----------------|
| **VM/Server** | 1 VM = 5 units (actual usage hours from billing) |
| **Containers** | 4 vCores = 2 units (based on actual consumed vCore-hours) |
| **Storage** | Per GB-month = 5 units |
| **Database** | 1 DB = 5 units (usage-based) |
| **Serverless** (Functions, Lambda, etc.) | Per 10 function invocations = 5 units |

**Key Difference from Asset Scan**: Billing data provides actual consumption (ephemeral resources, scaling events, temporary workloads) that wouldn't appear in a static inventory scan.

## Data Sources

### 1. Cloud Asset Inventory (Real-time Scan)
Discovers currently deployed resources:
- Running EC2 instances, Azure VMs, GCP Compute instances
- Databases (RDS, Azure SQL, Cloud SQL)
- Container clusters (ECS, AKS, GKE)
- Storage buckets
- Serverless functions

### 2. Cloud Billing Files
Historical consumption data (typically monthly CSV exports):
- **AWS**: Cost and Usage Report (CSV or Parquet)
- **Azure**: Cost Management export (CSV)
- **GCP**: BigQuery export or CSV from Cloud Billing

Billing data includes:
- All consumed resources (including ephemeral)
- Actual usage metrics (vCore-hours, GB-months, invocations, etc.)
- Cost breakdowns by service/region
- Reserved vs. on-demand usage

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
│   │   │   ├── scanner.go         # Discover current EC2, RDS, etc.
│   │   │   ├── billing.go         # Parse AWS Cost & Usage Report
│   │   │   └── mapper.go          # Map assets to internal model
│   │   ├── azure/
│   │   │   ├── client.go
│   │   │   ├── scanner.go         # Discover current VMs, SQL, etc.
│   │   │   ├── billing.go         # Parse Azure Cost Management CSV
│   │   │   └── mapper.go
│   │   └── gcp/
│   │       ├── client.go
│   │       ├── scanner.go         # Discover current Compute, Cloud SQL
│   │       ├── billing.go         # Parse GCP BigQuery/CSV export
│   │       └── mapper.go
│   ├── billing/                    # Billing file processing
│   │   ├── parser.go              # Multi-format billing file parser
│   │   ├── validator.go           # Validate billing file structure
│   │   ├── aggregator.go          # Aggregate billing metrics by resource type
│   │   └── enricher.go            # Enrich inventory with billing data
│   ├── assets/                     # Asset model & aggregation
│   │   ├── types.go               # Asset types (VM, DB, Container, etc.)
│   │   ├── aggregator.go          # Combine assets across projects/subscriptions
│   │   ├── converter.go           # Data-driven converter using config rules
│   │   └── enrichment.go          # Merge inventory and billing data
│   ├── config/
│   │   ├── config.go              # Config struct definitions
│   │   └── loader.go              # Load and validate config file
│   └── models/
│       ├── asset.go               # Core asset struct
│       ├── billing.go             # Billing record struct
│       └── enriched.go            # Combined asset + billing data
├── pkg/
│   ├── output/
│   │   ├── excel.go              # Excel file generation (.xlsx)
│   │   └── formatter.go          # Format assets for output
│   └── discovery/
│       └── scanner.go            # Discovery orchestration across providers
├── config.example.json            # Example config with conversion rules
├── sample-data/
│   ├── aws-billing-sample.csv    # Sample AWS Cost & Usage Report
│   ├── azure-billing-sample.csv  # Sample Azure export
│   └── gcp-billing-sample.csv    # Sample GCP export
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
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test (e.g., billing parser)
go test -run TestBillingParser ./internal/billing

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
- **Packages**: lowercase, single-word names (`providers`, `assets`, `billing`, `output`)
- **Exported Functions**: CamelCase starting with uppercase (e.g., `ScanAWSAssets`, `ParseBillingFile`)
- **Interfaces**: `-er` suffix (e.g., `CloudProvider`, `BillingParser`, `ResourceEnricher`)
- **Constants**: UPPER_SNAKE_CASE
- **Error Handling**: Check all errors; wrap with context using `fmt.Errorf("operation: %w", err)`
- **Comments**: Document exported types and functions with `//` comment above declaration

### Asset Discovery Pattern

Each cloud provider implements the `CloudProvider` interface for current inventory:

```go
// internal/providers/provider.go
package providers

type CloudProvider interface {
    // Discover all assets in all projects/subscriptions (current state)
    DiscoverAssets(ctx context.Context) ([]Asset, error)
    
    // Validate credentials
    Validate(ctx context.Context) error
}

// Asset is the common internal representation
type Asset struct {
    ID              string                 // Unique asset ID
    Type            string                 // VM, Database, Container, Storage, Function
    Name            string
    Cloud           string                 // AWS, Azure, GCP
    Project         string                 // Project ID, Subscription ID, etc.
    Count           int                    // Number of currently deployed instances
    Metadata        map[string]interface{} // vCores, storage size, etc.
    SourceType      string                 // "inventory" or "billing"
}
```

### Billing File Processing

```go
// internal/billing/parser.go
package billing

// BillingRecord represents a single line in a cloud billing file
type BillingRecord struct {
    ServiceName       string
    ResourceType      string
    ResourceID        string
    UsageMetric       string  // e.g., "vCore-hours", "GB-month", "invocations"
    UsageAmount       float64
    CostCurrency      string
    CostAmount        float64
    TimePeriod        string  // YYYY-MM
    Region            string
    Project           string
    Metadata          map[string]string
}

// BillingParser interface for cloud provider-specific parsers
type BillingParser interface {
    ParseFile(filepath string) ([]BillingRecord, error)
    ValidateFormat(filepath string) error
}

// ParseBillingFile routes to correct parser based on cloud provider
func ParseBillingFile(filepath, cloudProvider string) ([]BillingRecord, error) {
    switch cloudProvider {
    case "aws":
        return parseAWSCostAndUsage(filepath)
    case "azure":
        return parseAzureCostManagement(filepath)
    case "gcp":
        return parseGCPBillingExport(filepath)
    default:
        return nil, fmt.Errorf("unknown cloud provider: %s", cloudProvider)
    }
}
```

### Enrichment Pattern

Merge current inventory with billing data:

```go
// internal/assets/enrichment.go
package assets

// EnrichedAsset combines current inventory with billing metrics
type EnrichedAsset struct {
    Asset              *Asset            // Current deployed asset
    BillingMetrics    *BillingMetrics   // Actual consumption from billing
    CalculatedUnits   int               // Based on billing data
    IsEphemeral       bool              // Only found in billing, not in current inventory
}

type BillingMetrics struct {
    TotalMonthlyUsage float64           // vCore-hours, GB-months, etc.
    AverageUnitsUsed  float64           // Estimated units from billing
    PeakUsage         float64           // Peak usage in period
    Invocations       int64             // For serverless
    CostAmount        float64           // Actual cost
    Period            string            // YYYY-MM
}

// EnrichAssets merges current inventory with billing data
func EnrichAssets(assets []Asset, billingRecords []billing.BillingRecord, rules config.SyntheticUnitRules) []EnrichedAsset {
    // Map assets by resource ID
    // Match with billing records
    // Calculate synthetic units based on actual usage
    // Identify ephemeral resources (only in billing)
}
```

### Synthetic Unit Conversion (Config-Driven)

The conversion rules are loaded from config and applied to both current inventory AND billing data:

```go
// internal/assets/converter.go
package assets

// ConvertToSyntheticUnits uses actual usage metrics from billing
func ConvertToSyntheticUnits(asset Asset, metrics *BillingMetrics, rules config.SyntheticUnitRules) int {
    rule, exists := rules[asset.Type]
    if !exists {
        return 0
    }
    
    // Use billing metrics (actual usage) instead of just count
    return rule.CalculateFromUsage(asset, metrics)
}
```

Updated config rules to support usage-based calculations:

```json
{
  "syntheticUnits": {
    "rules": {
      "VM": {
        "type": "usage",
        "unitsPerHour": 0.0208,
        "billingMetric": "vCore-hours"
      },
      "Container": {
        "type": "usage",
        "baseUnit": 4,
        "unitsPerBase": 2,
        "minimum": 2,
        "maximumvCores": 16,
        "billingMetric": "vCore-hours"
      },
      "Database": {
        "type": "usage",
        "unitsPerInstance": 5,
        "billingMetric": "instance-hours"
      },
      "Storage": {
        "type": "usage",
        "unitsPerGB": 0.00014,
        "billingMetric": "GB-months"
      },
      "Function": {
        "type": "batch",
        "batchSize": 1000000,
        "unitsPerBatch": 5,
        "billingMetric": "invocations"
      }
    }
  }
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
  "billing": {
    "aws": {
      "filePath": "/path/to/aws-cost-usage-report.csv",
      "format": "csv",
      "period": "2024-01"
    },
    "azure": {
      "filePath": "/path/to/azure-cost-export.csv",
      "format": "csv",
      "period": "2024-01"
    },
    "gcp": {
      "filePath": "/path/to/gcp-billing-export.csv",
      "format": "csv",
      "period": "2024-01"
    }
  },
  "syntheticUnits": {
    "rules": {
      "VM": {
        "type": "usage",
        "unitsPerHour": 0.0208,
        "billingMetric": "vCore-hours"
      },
      "Container": {
        "type": "usage",
        "baseUnit": 4,
        "unitsPerBase": 2,
        "minimum": 2,
        "billingMetric": "vCore-hours"
      },
      "Database": {
        "type": "usage",
        "unitsPerInstance": 5,
        "billingMetric": "instance-hours"
      },
      "Storage": {
        "type": "usage",
        "unitsPerGB": 0.00014,
        "billingMetric": "GB-months"
      },
      "Function": {
        "type": "batch",
        "batchSize": 1000000,
        "unitsPerBatch": 5,
        "billingMetric": "invocations"
      }
    }
  },
  "output": {
    "format": "excel",
    "filename": "cloud-assets-inventory-enriched.xlsx",
    "includeEphemeralResources": true,
    "includeBillingMetrics": true
  }
}
```

### CLI Usage
```bash
# Scan current inventory + process billing files
cloudcostcala --config config.json --output assets.xlsx

# Optional: scan only without billing data
cloudcostcala --config config.json --inventory-only --output assets.xlsx

# Optional: process billing files only (for cost analysis)
cloudcostcala --config config.json --billing-only --output billing-analysis.xlsx
```

## Dependencies to Consider

- **AWS SDK**: `github.com/aws/aws-sdk-go-v2`
- **Azure SDK**: `github.com/Azure/azure-sdk-for-go`
- **GCP SDK**: `cloud.google.com/go`
- **CSV Parsing**: `encoding/csv` (standard library) or `github.com/gocarina/gocsv`
- **Excel**: `github.com/xuri/excelize` (for .xlsx generation)
- **Config**: Standard `encoding/json` or `github.com/spf13/viper`
- **Logging**: `go.uber.org/zap` or standard `log`
- **Testing**: `github.com/stretchr/testify` (assertions)

## Testing Strategy

- **Unit Tests**: 
  - Billing file parsing for each cloud provider
  - Synthetic unit calculation with real usage metrics
  - Asset enrichment logic
- **Integration Tests**: 
  - End-to-end workflow: discover assets + parse billing + enrich + output
  - Mock billing files in CSV format
- **Sample Data**: Include sample billing files in `sample-data/` for testing

```go
// Example test for billing-driven synthetic unit conversion
func TestSyntheticUnitConversionFromBilling(t *testing.T) {
    rules := config.SyntheticUnitsRules{
        "VM": {Type: "usage", UnitsPerHour: 0.0208, BillingMetric: "vCore-hours"},
    }
    
    metrics := &BillingMetrics{
        TotalMonthlyUsage: 720, // 30 days * 24 hours
    }
    
    asset := Asset{Type: "VM", Count: 1}
    
    // 720 hours * 0.0208 = ~15 units
    result := ConvertToSyntheticUnits(asset, metrics, rules)
    assert.Equal(t, 15, result)
}
```

## Documentation

- **README.md**: Quick start, installation, usage examples with sample outputs
- **config.example.json**: Well-commented configuration template for all rule types and billing sources
- **docs/ARCHITECTURE.md**: Dual-source data model, enrichment flow
- **docs/BILLING_FILE_FORMAT.md**: Expected format for AWS, Azure, GCP billing files
- **docs/SYNTHETIC_UNITS.md**: Synthetic unit model and usage-based calculations
- **sample-data/**: Example billing files for testing
- **Code Comments**: Asset types, billing parsing, enrichment logic

## Workflow Overview

```
User provides config (credentials + billing file paths + conversion rules)
  ↓
├─→ Scan current assets across all clouds/projects/subscriptions
└─→ Parse billing files (AWS Cost & Usage, Azure Cost Mgmt, GCP export)
  ↓
Match and enrich: merge inventory with billing metrics
  ↓
Identify ephemeral resources (only in billing, not in current inventory)
  ↓
Calculate synthetic units based on actual usage from billing data
  ↓
Group results and write to Excel
  ↓
Excel output:
  [Asset Type | Current Count | Ephemeral Count | Monthly Usage | Synthetic Units]
```

## Key Benefits of Dual-Source Approach

1. **Ephemeral Resources**: Capture temporary instances, auto-scaled resources, spot instances
2. **Accurate Usage**: Based on actual consumption (vCore-hours, invocations, GB-months)
3. **Realistic Budget**: Account for scaling patterns and temporary workloads
4. **Cost Correlation**: Can correlate actual cloud costs with synthetic units
5. **Trend Analysis**: Historical billing data shows usage patterns

## CI/CD (When Ready)

- Test on Linux, macOS, Windows
- Test billing file parsing with sample CSV files
- Validate Excel file generation with enriched data
- Test config loading and validation
- Cross-compile binaries for distribution

## Notes

- Asset discovery + billing processing requires handling multiple data sources
- Billing file formats vary significantly by cloud provider (AWS, Azure, GCP)
- Match assets between inventory and billing by resource ID/name (may require fuzzy matching)
- Ephemeral resources only in billing are important for realistic capacity planning
- Synthetic unit calculations now use actual consumption metrics, not just instance counts
- Config file is the single source of truth for both credentials and conversion rules
