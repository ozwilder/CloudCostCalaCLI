# CloudCostCalaCLI - Working Demo

## Overview

CloudCostCalaCLI successfully processes billing data from AWS, Azure, and GCP, normalizes metrics to instance-hours, and calculates synthetic units for budget planning.

## Demo Execution

### Step 1: Parse Billing Files

The tool reads CSV billing exports from all three cloud providers:

```
[AWS] Processing billing file...
  ✓ Loaded 7 AWS billing records

[Azure] Processing billing file...
  ✓ Loaded 6 Azure billing records

[GCP] Processing billing file...
  ✓ Loaded 6 GCP billing records
```

### Step 2: Normalize Metrics

All metrics converted to "average instances per hour":

```
Billing period: 2024-01 (31 days = 744 hours)

Example calculations:
  720 instance-hours / 744 hours = 0.97 avg instances/hr
  360 instance-hours / 744 hours = 0.48 avg instances/hr
```

### Step 3: Process Results

Asset types aggregated with enriched data:

```
Asset types found: [VM Database Container Function Storage]
Enriched 5 asset types
```

### Step 4: Generate Output

Console table display:

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

Excel file generated: `cloud-assets-inventory.xlsx`

## How It Works

### 1. Billing File Parsing

**Input CSV Format** (sample AWS):
```csv
service,resourceType,resourceId,instanceHours,period,region
EC2,VM,i-1234567890abcdef0,720,2024-01,us-east-1
RDS,Database,db-primary-1,744,2024-01,us-east-1
Lambda,Function,function-1,1000,2024-01,us-east-1
```

**Supported Cloud Providers:**
- AWS Cost and Usage Report
- Azure Cost Management export
- GCP BigQuery/CSV export

### 2. Instance-Hour Normalization

**Calculation:**
```
Average Instances Per Hour = Total Instance-Hours / (Days in Period × 24)

Examples:
- 720 instance-hours / 744 hours (31 days) = 0.97 avg instances
- 360 instance-hours / 744 hours = 0.48 avg instances
- 1000 invocations × 0.5 (10:1 ratio) = 50 instance-hours = 0.067 avg instances
```

All metrics unified to a single dimension: **average instances per hour**

### 3. Synthetic Unit Conversion

**Config-Driven Rules:**
```json
{
  "VM": { "unitsPerInstance": 5 },
  "Database": { "unitsPerInstance": 5 },
  "Container": { "unitsPerInstance": 2 },
  "Storage": { "unitsPerInstance": 5 },
  "Function": { "unitsPerInstance": 1 }
}
```

**Calculation:**
```
Synthetic Units = Average Instances Per Hour × Units Per Instance

Examples:
- 4.76 VMs × 5 units/VM = 23.8 → 24 units
- 3.00 Databases × 5 units/DB = 15 units
- 2.10 Containers × 2 units/Container = 4.2 → 4 units
```

### 4. Output Formats

**Console Table:**
- Real-time display during processing
- Color-coded for easy reading
- Running totals

**Excel Report:**
- Professional formatting
- Header styling
- Automatic totals row
- Column auto-sizing
- Ready for budget planning

## Sample Data

### Test Data Used

**AWS Billing:**
- 3 VMs: 720 + 360 + 240 = 1,320 instance-hours
- 1 Database: 744 instance-hours
- 1 Container: 480 instance-hours
- 1 Function: 1,000 invocations
- 1 Storage: 744 instance-hours

**Azure Billing:**
- 2 VMs: 744 + 372 = 1,116 instance-hours
- 1 Database: 744 instance-hours
- 1 Container: 480 instance-hours
- 1 Function: 500 invocations
- 1 Storage: 744 instance-hours

**GCP Billing:**
- 2 Compute instances: 744 + 360 = 1,104 instance-hours
- 1 Cloud SQL: 744 instance-hours
- 1 GKE node: 600 instance-hours
- 1 Cloud Function: 200 invocations
- 1 Cloud Storage: 744 instance-hours

### Aggregated Results

```
VM:        4.76 avg instances → 24 synthetic units
Database:  3.00 avg instances → 15 synthetic units
Container: 2.10 avg instances → 4 synthetic units
Function:  2.28 avg instances → 2 synthetic units
Storage:   3.00 avg instances → 15 synthetic units

TOTAL: 15.14 avg instances → 60 synthetic units
```

## Usage

### Build

```bash
cd /tmp/CloudCostCalaCLI
go build -o bin/cloudcostcala ./cmd/cloudcostcala
```

Or use Makefile:

```bash
make build
```

### Run

```bash
./bin/cloudcostcala --config config.example.json --output cloud-assets-inventory.xlsx
```

### Custom Configuration

Edit `config.example.json`:
- Set billing file paths to your CSV exports
- Adjust synthetic unit rules per customer
- Customize output filename

### Makefile Commands

```bash
make build      # Compile binary
make run        # Build and run with defaults
make clean      # Remove artifacts
make test       # Run tests (when added)
make fmt        # Format code
make lint       # Run linter
make all        # Clean, build, and run
```

## Key Features Demonstrated

✅ **Multi-cloud parsing** - AWS, Azure, GCP all in one run
✅ **Instance-hour normalization** - All metrics unified to single dimension
✅ **Config-driven conversion** - Rules in JSON, no code changes needed
✅ **Ephemeral resource detection** - Tracks resources only in billing, not in inventory
✅ **Professional output** - Console tables and Excel reports
✅ **Error handling** - Graceful degradation if some files missing
✅ **Formatted billing metrics** - Clear examples showing normalization math

## Architecture

See `.github/copilot-instructions.md` for detailed architecture and future development plans.

## Files Created

```
✓ cmd/cloudcostcala/main.go         - CLI entry point
✓ internal/config/config.go         - Configuration structures
✓ internal/config/loader.go         - Config file loading
✓ internal/models/asset.go          - Data models
✓ internal/billing/parser.go        - CSV parsing for 3 clouds
✓ internal/billing/normalizer.go    - Instance-hour conversion
✓ internal/assets/converter.go      - Synthetic unit calculation
✓ internal/assets/enrichment.go     - Asset enrichment logic
✓ pkg/output/excel.go               - Excel file generation
✓ config.example.json               - Configuration template
✓ sample-data/*.csv                 - Test billing files
✓ Makefile                          - Build automation
✓ README.md                         - Usage documentation
```

## Next Steps

The prototype demonstrates core functionality. Future enhancements:

1. **Cloud API Integration** - Direct asset discovery from AWS/Azure/GCP APIs
2. **Database Support** - Persist results in database
3. **Time-series Analysis** - Track trends over multiple months
4. **Advanced Filtering** - By region, project, service type
5. **Alerting** - Notify when units exceed budget
6. **Multi-tenant Support** - Handle multiple customers
7. **API Server** - REST endpoint for programmatic access
