# CloudCostCalaCLI
# CloudCostCalaCLI

A Go-based CLI tool for scanning cloud assets and enriching inventory with billing data to calculate synthetic units for accurate budget planning.

## Features

- **Multi-cloud support**: AWS, Azure, GCP
- **Dual-source data**: Current asset inventory + historical billing data
- **Instance-hour normalization**: All metrics unified to "average instances per hour"
- **Config-driven conversions**: Flexible synthetic unit rules
- **Excel output**: Professional formatted reports with totals and summaries

## Quick Start

### Build

```bash
make build
```

### Run

```bash
./bin/cloudcostcala --config config.example.json --output cloud-assets-inventory.xlsx
```

Or use make:

```bash
make run
```

### Configure

Edit `config.example.json` with your billing file paths:

```json
{
  "billing": {
    "aws": {
      "filePath": "path/to/aws-billing.csv",
      "format": "csv",
      "period": "2024-01"
    },
    ...
  },
  "syntheticUnits": {
    "rules": {
      "VM": { "unitsPerInstance": 5 },
      "Database": { "unitsPerInstance": 5 },
      ...
    }
  }
}
```

## Billing File Format

Billing files should be CSV with columns:
- `service`: Cloud service name
- `resourceType`: Mapped to asset type (VM, Database, Container, Storage, Function)
- `resourceId`: Unique resource identifier
- `instanceHours`: Total instance-hours for the period
- `period`: YYYY-MM format
- `region`: Cloud region

### Example

```csv
service,resourceType,resourceId,instanceHours,period,region
EC2,VM,i-1234567890abcdef0,720,2024-01,us-east-1
RDS,Database,db-primary-1,744,2024-01,us-east-1
```

## How It Works

1. **Parse Billing Files**: Reads CSV exports from AWS, Azure, GCP
2. **Normalize to Instance-Hours**: Converts all metrics to average instances per hour
   - Example: 720 instance-hours for 31-day month = 0.97 avg instances/hr
3. **Enrich Assets**: Merges current inventory with billing metrics
4. **Calculate Synthetic Units**: Applies config-defined conversion rules
   - Formula: Units = Average Instances/Hr × Units-Per-Instance
5. **Generate Report**: Creates Excel file with summary table

## Synthetic Units Model

Default conversion rules:

| Asset Type | Units Per Instance |
|------------|--------------------|
| VM | 5 |
| Database | 5 |
| Container | 2 |
| Storage | 5 |
| Function | 1 |

Customizable in config file - no code changes needed.

## Example Output

```
╔════════════════╦════════════════╦════════════════╦════════════════╦════════════════╗
║  Asset Type    ║ Current Count  ║ Ephemeral Cnt  ║ Avg Inst/Hr    ║ Synthetic Unts ║
╠════════════════╬════════════════╬════════════════╬════════════════╬════════════════╣
║ VM             ║              0 ║              1 ║           4.76 ║             24 ║
║ Database       ║              0 ║              1 ║           3.00 ║             15 ║
║ Container      ║              0 ║              1 ║           2.10 ║              4 ║
║ Function       ║              0 ║              1 ║           2.28 ║              2 ║
║ Storage        ║              0 ║              1 ║           3.00 ║             15 ║
╠════════════════╬════════════════╬════════════════╬════════════════╬════════════════╣
║ TOTAL          ║              0 ║              5 ║          15.14 ║             60 ║
╚════════════════╩════════════════╩════════════════╩════════════════╩════════════════╝
```

## Project Structure

```
CloudCostCalaCLI/
├── cmd/cloudcostcala/          # CLI entry point
├── internal/
│   ├── config/                 # Configuration loading
│   ├── models/                 # Data structures
│   ├── billing/                # Billing file parsing & normalization
│   ├── assets/                 # Asset enrichment & conversion
│   └── providers/              # Cloud provider implementations (future)
├── pkg/
│   └── output/                 # Excel generation
├── sample-data/                # Example billing files
├── config.example.json         # Configuration template
└── Makefile                    # Build commands
```

## Commands

```bash
# Build
make build

# Run with default config
make run

# Run with custom config
./bin/cloudcostcala --config my-config.json --output my-report.xlsx

# Test
make test

# Format code
make fmt

# Clean build artifacts
make clean
```

## Architecture

See `.github/copilot-instructions.md` for detailed architecture documentation.

## License

MIT 
