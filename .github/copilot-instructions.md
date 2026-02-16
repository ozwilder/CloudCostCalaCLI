# Copilot Instructions for CloudCostCalaCLI

## Project Overview

CloudCostCalaCLI is a Go-based CLI tool that scans and discovers all cloud assets across AWS, Azure, and GCP cloud environments and enriches the inventory with actual billing data for accurate budget planning.

**Dual Input Approach:**
1. **Current Asset Inventory**: Real-time scan of live cloud resources (VMs, databases, containers, etc.)
2. **Cloud Billing Files**: Historical consumption data (CSV format) capturing ephemeral resources and actual usage patterns

The tool aggregates both data sources across projects/subscriptions/resource groups and outputs an Excel file with:
- Each asset type as one row
- **Average instances/assets in use per hour** (normalized from billing data)
- **Synthetic Units**: Standardized credit-points based on average hourly consumption for accurate budget planning

**Unified Metric**: Everything is translated to "average number of instances in use per hour" - not storage volume or traffic, but actual resource count deployed on average.

## Synthetic Units Model

Synthetic units are customer-specific conversion rules that translate cloud assets into abstract credit points based on **average hourly instance consumption**.

**Default Model** (can be customized in config):

| Asset Type | Conversion Rule |
|------------|-----------------|
| **VM/Server** | 1 VM instance in use (avg/hr) = 5 units |
| **Containers** | 1 container instance in use (avg/hr) = 2 units |
| **Database** | 1 DB instance in use (avg/hr) = 5 units |
| **Storage** | 1 storage service instance in use (avg/hr) = 5 units |
| **Serverless** (Functions, Lambda, etc.) | 10 concurrent function instances (avg/hr) = 5 units |

**Key Metric**: All consumption is normalized to "average instances deployed per hour" regardless of volumetric attributes (GB, vCores, traffic). This provides a simple, consistent capacity model.

## Data Sources

### 1. Cloud Asset Inventory (Real-time Scan)
Discovers currently deployed resources:
- Running EC2 instances, Azure VMs, GCP Compute instances
- Databases (RDS, Azure SQL, Cloud SQL)
- Container clusters (ECS, AKS, GKE)
- Storage services
- Serverless functions

### 2. Cloud Billing Files
Historical consumption data (typically monthly CSV exports) that shows:
- **Actual instance hours**: How many resources were used and for how long
- **Ephemeral resources**: Temporary instances, auto-scaled workloads, spot instances not in current inventory
- **Usage patterns**: Peak vs. average consumption throughout the period

The tool calculates **average instances in use per hour** for the billing period:
- Example: If AWS billing shows 200 EC2 instance-hours for a month (30 days), average is 200 / (30 * 24) = 0.28 instances
- This single normalized metric applies to all asset types

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
│   │   ├── normalizer.go          # Normalize to instance-hours
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

# Run specific test (e.g., instance-hour normalization)
go test -run TestNormalizeInstanceHours ./internal/billing

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
- **Exported Functions**: CamelCase starting with uppercase (e.g., `ScanAWSAssets`, `ParseBillingFile`, `NormalizeInstanceHours`)
- **Interfaces**: `-er` suffix (e.g., `CloudProvider`, `BillingParser`, `InstanceNormalizer`)
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
    ID                    string                 // Unique asset ID
    Type                  string                 // VM, Database, Container, Storage, Function
    Name                  string
    Cloud                 string                 // AWS, Azure, GCP
    Project               string                 // Project ID, Subscription ID, etc.
    CurrentInstanceCount  int                    // Currently deployed instances
    Metadata              map[string]interface{} // Additional properties
    SourceType            string                 // "inventory" or "billing"
}
```

### Billing File Processing & Normalization

```go
// internal/billing/parser.go
package billing

// BillingRecord represents a single line in a cloud billing file
type BillingRecord struct {
    ServiceName       string                 // EC2, RDS, Storage, etc.
    ResourceType      string                 // VM, Database, Container, etc.
    ResourceID        string
    InstanceHours     float64                // Already in instance-hours (normalized)
    TimePeriod        string                 // YYYY-MM
    Region            string
    Project           string
    Metadata          map[string]string
}

// BillingParser interface for cloud provider-specific parsers
type BillingParser interface {
    // Parse file and normalize all metrics to instance-hours
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

### Instance-Hour Normalization

```go
// internal/billing/normalizer.go
package billing

// NormalizeToInstanceHours converts any resource usage to average instances per hour
func NormalizeToInstanceHours(records []BillingRecord, billingPeriod string) map[string]float64 {
    // billingPeriod format: "2024-01" (January 2024)
    daysInPeriod := calculateDaysInPeriod(billingPeriod)
    hoursInPeriod := float64(daysInPeriod * 24)
    
    normalized := make(map[string]float64)
    
    for _, record := range records {
        // Sum instance-hours by resource type
        normalized[record.ResourceType] += record.InstanceHours
    }
    
    // Convert total instance-hours to average instances per hour
    for resourceType := range normalized {
        normalized[resourceType] = normalized[resourceType] / hoursInPeriod
    }
    
    return normalized
}

// Example calculation:
// AWS reports: EC2 instance-hours = 720 (30 days * 24 hours * 1 instance)
// Result: 720 / (30 * 24) = 1 average instance per hour
// 
// If 3 instances ran for 10 days each: (3 * 10 * 24) / (30 * 24) = 1 average
```

### Enrichment Pattern

Merge current inventory with billing metrics:

```go
// internal/assets/enrichment.go
package assets

// EnrichedAsset combines current inventory with billing metrics
type EnrichedAsset struct {
    AssetType              string
    CurrentlyDeployed      int                 // From current inventory
    AverageInstancesPerHr  float64             // From billing (normalized)
    HasEphemeralUsage      bool                // True if billing shows usage not in current inventory
    CalculatedUnits        int                 // Based on average instances
}

// EnrichAssets merges current inventory with billing data
// Normalizes all metrics to "average instances per hour"
func EnrichAssets(assets []Asset, billingRecords []BillingRecord, 
    billingPeriod string, rules config.SyntheticUnitRules) []EnrichedAsset {
    
    // Group current assets by type and count
    assetsByType := groupAssets(assets)
    
    // Normalize billing records to average instances per hour
    avgInstancesByType := NormalizeToInstanceHours(billingRecords, billingPeriod)
    
    // Merge and create enriched assets
    enriched := make([]EnrichedAsset, 0)
    
    // Include types from both inventory and billing
    allTypes := mergeKeys(assetsByType, avgInstancesByType)
    
    for _, assetType := range allTypes {
        enriched = append(enriched, EnrichedAsset{
            AssetType:              assetType,
            CurrentlyDeployed:      assetsByType[assetType],
            AverageInstancesPerHr:  avgInstancesByType[assetType],
            HasEphemeralUsage:      avgInstancesByType[assetType] > 0 && assetsByType[assetType] == 0,
            CalculatedUnits:        calculateUnits(assetType, avgInstancesByType[assetType], rules),
        })
    }
    
    return enriched
}
```

### Synthetic Unit Conversion (Config-Driven with Instance Metric)

The conversion rules now use normalized instance-hour metrics:

```go
// internal/assets/converter.go
package assets

// ConvertToSyntheticUnits uses average hourly instance count from billing
func ConvertToSyntheticUnits(assetType string, avgInstancesPerHour float64, 
    rules config.SyntheticUnitRules) int {
    
    rule, exists := rules[assetType]
    if !exists {
        return 0
    }
    
    // Simple multiplication: instances per hour * units per instance
    unitsPerInstance := rule.UnitsPerInstance
    totalUnits := int(math.Round(avgInstancesPerHour * float64(unitsPerInstance)))
    
    return totalUnits
}
```

Updated config rules:

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
        "unitsPerInstance": 5
      },
      "Container": {
        "unitsPerInstance": 2
      },
      "Database": {
        "unitsPerInstance": 5
      },
      "Storage": {
        "unitsPerInstance": 5
      },
      "Function": {
        "unitsPerInstance": 0.5
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

## Configuration Management

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
  - Instance-hour normalization with different period lengths
  - Synthetic unit calculation with average instance metrics
  - Asset enrichment logic
- **Integration Tests**: 
  - End-to-end workflow: discover assets + parse billing + normalize + enrich + output
  - Mock billing files in CSV format with known instance-hour counts
- **Sample Data**: Include sample billing files in `sample-data/` for testing

```go
// Example test for instance-hour normalization
func TestNormalizeToInstanceHours(t *testing.T) {
    // January 2024 = 31 days = 744 hours
    billingPeriod := "2024-01"
    
    records := []BillingRecord{
        {ResourceType: "VM", InstanceHours: 744},     // 1 instance all month
        {ResourceType: "VM", InstanceHours: 372},     // 0.5 instances all month
        {ResourceType: "Database", InstanceHours: 744}, // 1 DB all month
    }
    
    result := NormalizeToInstanceHours(records, billingPeriod)
    
    // VM: (744 + 372) / 744 = 1.5 average instances per hour
    assert.InDelta(t, 1.5, result["VM"], 0.01)
    
    // Database: 744 / 744 = 1 average instance per hour
    assert.InDelta(t, 1.0, result["Database"], 0.01)
}

// Example test for synthetic unit conversion
func TestSyntheticUnitConversion(t *testing.T) {
    rules := config.SyntheticUnitsRules{
        "VM": {UnitsPerInstance: 5},
        "Database": {UnitsPerInstance: 5},
    }
    
    // 1.5 average VMs * 5 units per VM = 7.5 → 8 units
    units := ConvertToSyntheticUnits("VM", 1.5, rules)
    assert.Equal(t, 8, units)
    
    // 1.0 average DB * 5 units per DB = 5 units
    units = ConvertToSyntheticUnits("Database", 1.0, rules)
    assert.Equal(t, 5, units)
}
```

## Documentation

- **README.md**: Quick start, installation, usage examples with sample outputs
- **config.example.json**: Well-commented configuration template
- **docs/ARCHITECTURE.md**: Dual-source data model, enrichment flow
- **docs/BILLING_FILE_FORMAT.md**: Expected format for AWS, Azure, GCP billing files
- **docs/INSTANCE_HOUR_NORMALIZATION.md**: How volumetric metrics are converted to instance-hour average
- **docs/SYNTHETIC_UNITS.md**: Synthetic unit model based on average hourly instances
- **sample-data/**: Example billing files for testing
- **Code Comments**: Asset types, billing parsing, enrichment logic

## Workflow Overview

```
User provides config (credentials + billing file paths + conversion rules)
  ↓
├─→ Scan current assets across all clouds/projects/subscriptions
└─→ Parse billing files (AWS Cost & Usage, Azure Cost Mgmt, GCP export)
  ↓
Normalize all billing metrics to "average instances in use per hour"
  ↓
Match and enrich: merge inventory with normalized billing metrics
  ↓
Identify ephemeral resources (in billing but not in current inventory)
  ↓
Calculate synthetic units (average instances per hour × units per instance)
  ↓
Group results and write to Excel
  ↓
Excel output:
  [Asset Type | Current Count | Avg Instances/Hr | Ephemeral? | Synthetic Units]
```

## Key Benefits of Instance-Hour Normalization

1. **Unified Metric**: All consumption reduced to "average instances per hour"
2. **No Volumetric Complexity**: Ignores GB, vCores, traffic - focuses on resource count
3. **Simple Calculation**: Units = Avg Instances × Units-Per-Instance (linear, predictable)
4. **Accurate Scaling**: Captures auto-scaled and ephemeral workloads
5. **Easy to Understand**: "1 VM for a month" = 1 instance-hour metric
6. **Customer-Agnostic**: Works identically across AWS, Azure, GCP

## Billing File Format Requirements

Billing files must include columns for:
- Service name / Resource type
- Resource identifier
- Instance-hours (or raw metric that can be converted to instance-hours)
- Billing period

Tool normalizes all metrics to instance-hours during parsing.

## CI/CD (When Ready)

- Test on Linux, macOS, Windows
- Test billing file parsing with sample CSV files
- Test instance-hour normalization with various period lengths
- Validate Excel file generation with enriched data
- Test config loading and validation
- Cross-compile binaries for distribution

## Notes

- All consumption metrics converted to "average instances per hour" - volumetric aspects (GB, vCores, traffic) ignored
- Synthetic unit calculation is simple: avg instances × units per instance
- Ephemeral resources are critical for realistic capacity planning
- Billing file parsing is cloud-provider-specific but normalization is universal
- Config file defines conversion rates and billing file locations
