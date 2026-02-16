package output

import (
	"fmt"

	"github.com/ozwilder/CloudCostCalaCLI/internal/models"
	"github.com/xuri/excelize/v2"
)

// WriteExcel generates an Excel file with aggregated asset data
func WriteExcel(filename string, assets []models.AggregatedOutput) error {
	f := excelize.NewFile()

	// Create header
	headers := []string{"Asset Type", "Current Count", "Ephemeral Count", "Avg Instances/Hr", "Synthetic Units"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+rune(i))
		f.SetCellValue("Sheet1", cell, header)
		
		// Bold header
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"D3D3D3"}, Pattern: 1},
		})
		f.SetCellStyle("Sheet1", cell, cell, style)
	}

	// Add data rows
	for i, asset := range assets {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), asset.AssetType)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), asset.CurrentCount)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), asset.EphemeralCount)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("%.2f", asset.AvgInstancesPerHour))
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), asset.SyntheticUnits)
	}

	// Adjust column widths
	f.SetColWidth("Sheet1", "A", "A", 15)
	f.SetColWidth("Sheet1", "B", "B", 15)
	f.SetColWidth("Sheet1", "C", "C", 16)
	f.SetColWidth("Sheet1", "D", "D", 18)
	f.SetColWidth("Sheet1", "E", "E", 15)

	// Add totals row
	if len(assets) > 0 {
		totalRow := len(assets) + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", totalRow), "TOTAL")
		
		// Sum formulas
		f.SetCellFormula("Sheet1", fmt.Sprintf("B%d", totalRow), fmt.Sprintf("SUM(B2:B%d)", totalRow-1))
		f.SetCellFormula("Sheet1", fmt.Sprintf("C%d", totalRow), fmt.Sprintf("SUM(C2:C%d)", totalRow-1))
		f.SetCellFormula("Sheet1", fmt.Sprintf("D%d", totalRow), fmt.Sprintf("SUM(D2:D%d)", totalRow-1))
		f.SetCellFormula("Sheet1", fmt.Sprintf("E%d", totalRow), fmt.Sprintf("SUM(E2:E%d)", totalRow-1))
		
		// Bold totals row
		boldStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFF00"}, Pattern: 1},
		})
		for col := 'A'; col <= 'E'; col++ {
			f.SetCellStyle("Sheet1", fmt.Sprintf("%c%d", col, totalRow), fmt.Sprintf("%c%d", col, totalRow), boldStyle)
		}
	}

	// Save file
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// PrintSummaryTable prints asset data to console
func PrintSummaryTable(assets []models.AggregatedOutput) {
	fmt.Println("\n╔════════════════╦════════════════╦════════════════╦════════════════╦════════════════╗")
	fmt.Println("║  Asset Type    ║ Current Count  ║ Ephemeral Cnt  ║ Avg Inst/Hr    ║ Synthetic Unts ║")
	fmt.Println("╠════════════════╬════════════════╬════════════════╬════════════════╬════════════════╣")

	totalCurrent := 0
	totalEphemeral := 0
	totalAvgInstances := 0.0
	totalUnits := 0

	for _, asset := range assets {
		fmt.Printf("║ %-14s ║ %14d ║ %14d ║ %14.2f ║ %14d ║\n",
			asset.AssetType,
			asset.CurrentCount,
			asset.EphemeralCount,
			asset.AvgInstancesPerHour,
			asset.SyntheticUnits)

		totalCurrent += asset.CurrentCount
		totalEphemeral += asset.EphemeralCount
		totalAvgInstances += asset.AvgInstancesPerHour
		totalUnits += asset.SyntheticUnits
	}

	fmt.Println("╠════════════════╬════════════════╬════════════════╬════════════════╬════════════════╣")
	fmt.Printf("║ %-14s ║ %14d ║ %14d ║ %14.2f ║ %14d ║\n",
		"TOTAL",
		totalCurrent,
		totalEphemeral,
		totalAvgInstances,
		totalUnits)
	fmt.Println("╚════════════════╩════════════════╩════════════════╩════════════════╩════════════════╝\n")
}
