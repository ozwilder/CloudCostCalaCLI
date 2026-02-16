package assets

import (
	"github.com/ozwilder/CloudCostCalaCLI/internal/config"
	"github.com/ozwilder/CloudCostCalaCLI/internal/models"
)

// EnrichAssets merges current inventory with billing data
func EnrichAssets(assets []models.Asset, avgInstancesByType map[string]float64,
	rules config.SyntheticUnitsConfig) []models.EnrichedAsset {

	// Group current assets by type
	assetsByType := make(map[string]int)
	for _, asset := range assets {
		assetsByType[asset.Type]++
	}

	// Merge and create enriched assets
	enriched := make([]models.EnrichedAsset, 0)
	allTypes := mergeKeysStr(assetsByType, avgInstancesByType)

	for _, assetType := range allTypes {
		currentCount := assetsByType[assetType]
		avgInstances := avgInstancesByType[assetType]
		hasEphemeral := avgInstances > 0 && currentCount == 0

		enriched = append(enriched, models.EnrichedAsset{
			AssetType:             assetType,
			CurrentlyDeployed:     currentCount,
			AverageInstancesPerHr: avgInstances,
			HasEphemeralUsage:     hasEphemeral,
			CalculatedUnits:       ConvertToSyntheticUnits(assetType, avgInstances, rules),
		})
	}

	return enriched
}

// AggregateForOutput converts enriched assets to output format
func AggregateForOutput(enriched []models.EnrichedAsset) []models.AggregatedOutput {
	output := make([]models.AggregatedOutput, len(enriched))

	for i, e := range enriched {
		ephemeralCount := 0
		if e.HasEphemeralUsage {
			ephemeralCount = 1 // Simplified: at least 1 ephemeral
		}

		output[i] = models.AggregatedOutput{
			AssetType:           e.AssetType,
			CurrentCount:        e.CurrentlyDeployed,
			EphemeralCount:      ephemeralCount,
			AvgInstancesPerHour: e.AverageInstancesPerHr,
			SyntheticUnits:      e.CalculatedUnits,
		}
	}

	return output
}

// mergeKeys returns unique keys from two maps
func mergeKeys(m1, m2 map[string]interface{}) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)

	for k := range m1 {
		if !keys[k] {
			keys[k] = true
			result = append(result, k)
		}
	}

	for k := range m2 {
		if !keys[k] {
			keys[k] = true
			result = append(result, k)
		}
	}

	return result
}

// Helper overload for []string keys
func mergeKeysStr(m1 map[string]int, m2 map[string]float64) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)

	for k := range m1 {
		if !keys[k] {
			keys[k] = true
			result = append(result, k)
		}
	}

	for k := range m2 {
		if !keys[k] {
			keys[k] = true
			result = append(result, k)
		}
	}

	return result
}
