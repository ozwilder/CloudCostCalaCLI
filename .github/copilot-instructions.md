# Copilot Instructions for CloudCostCalaCLI

## Project Overview

CloudCostCalaCLI is a PowerShell-based command-line tool for calculating cloud infrastructure costs. The project is in early stages with minimal established structure.

## Development Setup

### PowerShell Environment
- **Language**: Support both PowerShell 7+ (cross-platform) and Windows PowerShell 5.1
- **Compatibility Note**: Avoid PowerShell 7-only features (e.g., `using namespace` in certain contexts, native commands) to maintain Windows PowerShell compatibility
- **File Extension**: `.ps1` for scripts, `.psm1` for modules, `.psd1` for module manifests
- **Execution Policy**: May need to set execution policy for running scripts locally
  ```powershell
  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```
- **Testing**: Test on both PowerShell 7 and Windows PowerShell 5.1 when possible (especially before releases)

## Project Structure (Recommended)

While still being established, adopt this structure:

```
CloudCostCalaCLI/
├── .github/
│   └── copilot-instructions.md
├── src/
│   ├── CloudCostCalaCLI.psm1          # Main module file
│   ├── functions/                      # Individual function files
│   │   ├── Get-CloudCosts.ps1
│   │   ├── Calculate-Cost.ps1
│   │   └── ...
│   └── classes/                        # PowerShell classes if used
├── tests/
│   └── *.Tests.ps1                    # Pester test files
├── docs/
│   └── *.md                            # User documentation
├── CloudCostCalaCLI.psd1              # Module manifest (at root or in src/)
├── README.md
└── LICENSE
```

## Build, Test & Validation

### Testing (Pester)
Testing approach is undecided but consider adopting Pester if test coverage becomes important. Setup would look like:
```powershell
# Run all tests
Invoke-Pester -Path ./tests

# Run specific test file
Invoke-Pester -Path ./tests/Get-CloudCosts.Tests.ps1
```

### Linting (PSScriptAnalyzer)
When PSScriptAnalyzer is configured:
```powershell
# Analyze all scripts
Invoke-ScriptAnalyzer -Path ./src -Recurse

# Analyze with custom rules
Invoke-ScriptAnalyzer -Path ./src -IncludeRules @('PSUseConsistentWhitespace', 'PSAvoidUsingWildcardCharacters')
```

### Building Module
When a build process is established:
```powershell
# Typical approach: Copy files to output folder and validate module manifest
$OutPath = "./build/CloudCostCalaCLI"
Copy-Item -Path ./src/* -Destination $OutPath -Recurse
Test-ModuleManifest -Path ./CloudCostCalaCLI.psd1
```

## Key Conventions

### PowerShell Best Practices
- **Naming**: Use approved verbs (Get, Set, New, Remove, etc.) + Singular nouns in function names
  - Examples: `Get-CloudCosts`, `New-CostEstimate`, `Remove-Cache`
- **Parameters**: Use `[Parameter()]` attributes with clear descriptions
- **Documentation**: Use comment-based help with `.DESCRIPTION`, `.PARAMETER`, `.EXAMPLE`
- **Error Handling**: Use `$ErrorActionPreference = 'Stop'` at function start or error handling with try/catch

### Module Structure Pattern
```powershell
# File: src/CloudCostCalaCLI.psm1 (or .psd1)
# This would import all functions from the functions/ directory

# File: src/functions/Get-CloudCosts.ps1
function Get-CloudCosts {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory)]
        [string]$Provider  # AWS, Azure, GCP
    )
    
    <#
    .DESCRIPTION
    Retrieves current cloud costs from specified provider.
    
    .EXAMPLE
    Get-CloudCosts -Provider AWS
    #>
    
    # Implementation
}
```

### Cloud Provider Integration
- **Supported Providers**: AWS, Azure, GCP (all three)
- **Architecture**: Design for extensibility—consider using a provider interface/base class pattern or separate modules per provider
- **Example Pattern**:
  ```powershell
  # src/providers/AWS/Get-AWSCosts.ps1
  # src/providers/Azure/Get-AzureCosts.ps1
  # src/providers/GCP/Get-GCPCosts.ps1
  # src/CloudCostCalaCLI.psm1 (orchestrates all providers)
  ```
- **Authentication**: Use configuration files for credentials (see below)
- **Config File Format**: Define schema for storing API credentials securely (consider encrypted configs for sensitive data)

## Authentication & Configuration

### Configuration File Approach
The tool uses configuration files to manage cloud provider credentials. Key considerations:

- **Config File Location**: Define a standard location (e.g., `~/.CloudCostCalaCLI/config.json` or similar)
- **File Format**: Consider JSON or YAML for config files
- **Security**: Do NOT commit sensitive credentials to the repository
  - Use `.gitignore` to exclude config files with real credentials
  - Document how users should create their own config files
  - Consider encryption for sensitive fields (API keys, secrets)
- **Example Config Structure**:
  ```json
  {
    "providers": {
      "aws": {
        "accessKeyId": "***",
        "secretAccessKey": "***",
        "region": "us-east-1"
      },
      "azure": {
        "subscriptionId": "***",
        "clientId": "***",
        "clientSecret": "***",
        "tenantId": "***"
      },
      "gcp": {
        "projectId": "***",
        "serviceAccountKey": "***"
      }
    }
  }
  ```

## CI/CD (When Ready)

When adding GitHub Actions:
- Use `-ErrorAction Stop` in scripts to fail fast
- Validate PowerShell syntax before merging
- Test on both PowerShell 7 and Windows PowerShell 5.1 if possible
- Consider running against mock/test cloud APIs to validate integration logic

## Notes

- This is a new project with minimal initial structure—establish conventions early
- **Cross-platform support**: PowerShell 5.1 and 7+ compatibility requires care with syntax and cmdlets
- **Configuration security**: Protect sensitive credentials; use `.gitignore` for config files with real API keys
- Testing approach (Pester) can be decided later if needed
- Not planning PowerShell Gallery publication at this time, so can focus on internal quality
