package assets

import (
	"fmt"
	"math"

	"github.com/ozwilder/CloudCostCalaCLI/internal/config"
)

// ConvertToSyntheticUnits calculates synthetic units from average instances per hour
func ConvertToSyntheticUnits(assetType string, avgInstancesPerHour float64, rules config.SyntheticUnitsConfig) int {
	rule, exists := rules.Rules[assetType]
	if !exists {
		return 0 // Unknown asset type
	}

	// Simple formula: instances per hour * units per instance
	unitsPerInstance := rule.UnitsPerInstance
	totalUnits := int(math.Round(avgInstancesPerHour * float64(unitsPerInstance)))

	return totalUnits
}

// ConvertMultiple converts multiple asset types to synthetic units
func ConvertMultiple(avgInstancesByType map[string]float64, rules config.SyntheticUnitsConfig) map[string]int {
	result := make(map[string]int)

	for assetType, avgInstances := range avgInstancesByType {
		result[assetType] = ConvertToSyntheticUnits(assetType, avgInstances, rules)
	}

	return result
}

// PrintConversionExample shows how synthetic unit conversion works
func PrintConversionExample() {
	fmt.Println("\n=== Synthetic Unit Conversion ===")
	fmt.Println("Formula: Units = Average Instances Per Hour × Units Per Instance")
	fmt.Println("\nExamples (assuming default multipliers):")
	fmt.Printf("  1.0 VM × 5 = %d units\n", int(math.Round(1.0*5)))
	fmt.Printf("  1.5 VMs × 5 = %d units\n", int(math.Round(1.5*5)))
	fmt.Printf("  0.5 Database × 5 = %d units\n", int(math.Round(0.5*5)))
	fmt.Printf("  2.0 Containers × 2 = %d units\n", int(math.Round(2.0*2)))
}
