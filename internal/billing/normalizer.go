package billing

import (
	"fmt"

	"github.com/ozwilder/CloudCostCalaCLI/internal/models"
)

// NormalizeToInstanceHours converts total instance-hours to average instances per hour
func NormalizeToInstanceHours(records []models.BillingRecord, billingPeriod string) map[string]float64 {
	daysInPeriod := getDaysInPeriod(billingPeriod)
	hoursInPeriod := float64(daysInPeriod * 24)

	normalized := make(map[string]float64)

	// Sum instance-hours by resource type
	for _, record := range records {
		normalized[record.ResourceType] += record.InstanceHours
	}

	// Convert total instance-hours to average instances per hour
	for resourceType := range normalized {
		normalized[resourceType] = normalized[resourceType] / hoursInPeriod
	}

	return normalized
}

// getDaysInPeriod returns number of days in a given month
// Expected format: YYYY-MM
func getDaysInPeriod(period string) int {
	if len(period) < 7 {
		return 30 // Default
	}

	month := period[5:7]
	switch month {
	case "01", "03", "05", "07", "08", "10", "12":
		return 31
	case "04", "06", "09", "11":
		return 30
	case "02":
		// Simplified: assume 28 (could check for leap year)
		return 28
	default:
		return 30
	}
}

// AggregateByType groups billing records by resource type and returns normalized instance-hours
func AggregateByType(records []models.BillingRecord, billingPeriod string) map[string]float64 {
	return NormalizeToInstanceHours(records, billingPeriod)
}

// GetBillingPeriod extracts period from records (assumes all same period)
func GetBillingPeriod(records []models.BillingRecord) string {
	if len(records) > 0 {
		return records[0].TimePeriod
	}
	return "2024-01"
}

// PrintNormalizationExample shows how normalization works
func PrintNormalizationExample(period string) {
	daysInPeriod := getDaysInPeriod(period)
	hoursInPeriod := daysInPeriod * 24

	fmt.Printf("\n=== Instance-Hour Normalization ===\n")
	fmt.Printf("Period: %s (%d days = %d hours)\n", period, daysInPeriod, hoursInPeriod)
	fmt.Printf("\nExample calculations:\n")
	fmt.Printf("  720 instance-hours / %d hours = %.2f avg instances/hr (1 VM all month)\n", hoursInPeriod, 720.0/float64(hoursInPeriod))
	fmt.Printf("  360 instance-hours / %d hours = %.2f avg instances/hr (0.5 VMs all month)\n", hoursInPeriod, 360.0/float64(hoursInPeriod))
}
