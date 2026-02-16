# Copilot Instructions for CloudCostCalaCLI

## Project Overview

CloudCostCalaCLI is a PowerShell-based command-line tool for calculating cloud infrastructure costs. The project is in early stages with minimal established structure.

## Development Setup

### PowerShell Environment
- **Language**: PowerShell (Core/7.0+ recommended, but compatible with Windows PowerShell 5.1)
- **File Extension**: `.ps1` for scripts, `.psm1` for modules, `.psd1` for module manifests
- **Execution Policy**: May need to set execution policy for running scripts locally
  ```powershell
  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```

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

### Running Tests (Pester)
When tests are added, use Pester:
```powershell
# Run all tests
Invoke-Pester -Path ./tests

# Run specific test file
Invoke-Pester -Path ./tests/Get-CloudCosts.Tests.ps1

# Run with coverage
Invoke-Pester -Path ./tests -CodeCoverage ./src
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
- Structure code to support multiple providers: AWS, Azure, GCP
- Use parameter validation or separate functions per provider
- Consider using config files or environment variables for API credentials

## Documentation

- **README.md**: Keep focused on quick start and basic usage
- **Module Help**: Embed help in function definitions using comment-based help
- **docs/**: Create additional docs for architecture, API integration details, examples

## CI/CD (When Ready)

When adding GitHub Actions:
- Use `-ErrorAction Stop` in scripts to fail fast
- Validate PowerShell syntax before merging
- Run Pester tests on pull requests
- Consider publishing to PowerShell Gallery when mature

## Notes

- This is a new project with minimal initial structure—establish conventions early
- PowerShell 7+ (cross-platform) vs Windows PowerShell 5.1 (Windows-only) choice will affect compatibility
- Cloud API authentication handling requires careful credential management
