package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ozwilder/CloudCostCalaCLI/internal/assets"
	"github.com/ozwilder/CloudCostCalaCLI/internal/billing"
	"github.com/ozwilder/CloudCostCalaCLI/internal/config"
	"github.com/ozwilder/CloudCostCalaCLI/internal/models"
	"github.com/ozwilder/CloudCostCalaCLI/pkg/output"
)

func main() {
	configPath := flag.String("config", "config.example.json", "Path to configuration file")
	outputFile := flag.String("output", "cloud-assets-inventory.xlsx", "Output Excel file path")
	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║         CloudCostCalaCLI - Cloud Asset Inventory            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("\nConfiguration: %s\n", *configPath)

	// Collect assets from billing files
	allAssets := make([]models.Asset, 0)
	allBillingRecords := make([]models.BillingRecord, 0)

	// Process AWS billing
	if cfg.Billing.AWS.FilePath != "" {
		fmt.Println("\n[AWS] Processing billing file...")
		awsRecords, err := billing.ParseBillingFile(cfg.Billing.AWS.FilePath, "aws")
		if err != nil {
			log.Printf("Warning: Failed to parse AWS billing: %v", err)
		} else {
			allBillingRecords = append(allBillingRecords, awsRecords...)
			fmt.Printf("  ✓ Loaded %d AWS billing records\n", len(awsRecords))
		}
	}

	// Process Azure billing
	if cfg.Billing.Azure.FilePath != "" {
		fmt.Println("\n[Azure] Processing billing file...")
		azureRecords, err := billing.ParseBillingFile(cfg.Billing.Azure.FilePath, "azure")
		if err != nil {
			log.Printf("Warning: Failed to parse Azure billing: %v", err)
		} else {
			allBillingRecords = append(allBillingRecords, azureRecords...)
			fmt.Printf("  ✓ Loaded %d Azure billing records\n", len(azureRecords))
		}
	}

	// Process GCP billing
	if cfg.Billing.GCP.FilePath != "" {
		fmt.Println("\n[GCP] Processing billing file...")
		gcpRecords, err := billing.ParseBillingFile(cfg.Billing.GCP.FilePath, "gcp")
		if err != nil {
			log.Printf("Warning: Failed to parse GCP billing: %v", err)
		} else {
			allBillingRecords = append(allBillingRecords, gcpRecords...)
			fmt.Printf("  ✓ Loaded %d GCP billing records\n", len(gcpRecords))
		}
	}

	if len(allBillingRecords) == 0 {
		log.Fatal("No billing records loaded. Check config file paths.")
	}

	// Normalize billing data to instance-hours
	fmt.Println("\n[Processing] Normalizing billing metrics...")
	billingPeriod := billing.GetBillingPeriod(allBillingRecords)
	avgInstancesByType := billing.AggregateByType(allBillingRecords, billingPeriod)
	fmt.Printf("  ✓ Billing period: %s\n", billingPeriod)
	fmt.Printf("  ✓ Asset types found: %v\n", getKeys(avgInstancesByType))

	// Enrich assets with billing data
	fmt.Println("\n[Processing] Enriching assets...")
	enrichedAssets := assets.EnrichAssets(allAssets, avgInstancesByType, cfg.SyntheticUnits)
	fmt.Printf("  ✓ Enriched %d asset types\n", len(enrichedAssets))

	// Aggregate for output
	fmt.Println("\n[Processing] Aggregating results...")
	aggregated := assets.AggregateForOutput(enrichedAssets)

	// Print summary table
	output.PrintSummaryTable(aggregated)

	// Generate Excel file
	fmt.Printf("\n[Output] Generating Excel file: %s\n", *outputFile)
	if err := output.WriteExcel(*outputFile, aggregated); err != nil {
		log.Fatalf("Error writing Excel: %v", err)
	}
	fmt.Println("  ✓ Excel file generated successfully!")

	// Print examples
	fmt.Println("\n[Examples]")
	billing.PrintNormalizationExample(billingPeriod)
	assets.PrintConversionExample()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  Processing Complete!                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func getKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
