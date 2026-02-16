package billing

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ozwilder/CloudCostCalaCLI/internal/models"
)

// ParseBillingFile reads a billing CSV and converts to BillingRecords
func ParseBillingFile(filePath, cloudProvider string) ([]models.BillingRecord, error) {
	switch cloudProvider {
	case "aws":
		return parseAWSBilling(filePath)
	case "azure":
		return parseAzureBilling(filePath)
	case "gcp":
		return parseGCPBilling(filePath)
	default:
		return nil, fmt.Errorf("unknown cloud provider: %s", cloudProvider)
	}
}

// parseAWSBilling handles AWS Cost and Usage Report format
func parseAWSBilling(filePath string) ([]models.BillingRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open AWS billing file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read AWS billing CSV: %w", err)
	}

	var billingRecords []models.BillingRecord

	// Skip header (first row)
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 6 {
			continue
		}

		// Expected columns: service,resourceType,resourceId,instanceHours,period,region
		serviceType := records[i][0]
		resourceType := mapAWSServiceToType(serviceType)
		resourceID := records[i][2]
		instanceHours, _ := strconv.ParseFloat(records[i][3], 64)
		period := records[i][4]
		region := records[i][5]

		billingRecords = append(billingRecords, models.BillingRecord{
			ServiceName:   serviceType,
			ResourceType:  resourceType,
			ResourceID:    resourceID,
			InstanceHours: instanceHours,
			TimePeriod:    period,
			Region:        region,
			Project:       "aws-default",
			Metadata:      make(map[string]string),
		})
	}

	return billingRecords, nil
}

// parseAzureBilling handles Azure Cost Management format
func parseAzureBilling(filePath string) ([]models.BillingRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Azure billing file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read Azure billing CSV: %w", err)
	}

	var billingRecords []models.BillingRecord

	// Skip header
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 6 {
			continue
		}

		serviceType := records[i][0]
		resourceType := mapAzureServiceToType(serviceType)
		resourceID := records[i][2]
		instanceHours, _ := strconv.ParseFloat(records[i][3], 64)
		period := records[i][4]
		region := records[i][5]

		billingRecords = append(billingRecords, models.BillingRecord{
			ServiceName:   serviceType,
			ResourceType:  resourceType,
			ResourceID:    resourceID,
			InstanceHours: instanceHours,
			TimePeriod:    period,
			Region:        region,
			Project:       "azure-default",
			Metadata:      make(map[string]string),
		})
	}

	return billingRecords, nil
}

// parseGCPBilling handles GCP billing export format
func parseGCPBilling(filePath string) ([]models.BillingRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GCP billing file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read GCP billing CSV: %w", err)
	}

	var billingRecords []models.BillingRecord

	// Skip header
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 6 {
			continue
		}

		serviceType := records[i][0]
		resourceType := mapGCPServiceToType(serviceType)
		resourceID := records[i][2]
		instanceHours, _ := strconv.ParseFloat(records[i][3], 64)
		period := records[i][4]
		region := records[i][5]

		billingRecords = append(billingRecords, models.BillingRecord{
			ServiceName:   serviceType,
			ResourceType:  resourceType,
			ResourceID:    resourceID,
			InstanceHours: instanceHours,
			TimePeriod:    period,
			Region:        region,
			Project:       "gcp-default",
			Metadata:      make(map[string]string),
		})
	}

	return billingRecords, nil
}

// Service type mappers
func mapAWSServiceToType(service string) string {
	service = strings.ToLower(service)
	if strings.Contains(service, "ec2") {
		return "VM"
	}
	if strings.Contains(service, "rds") {
		return "Database"
	}
	if strings.Contains(service, "lambda") {
		return "Function"
	}
	if strings.Contains(service, "ecs") {
		return "Container"
	}
	if strings.Contains(service, "s3") {
		return "Storage"
	}
	return "Other"
}

func mapAzureServiceToType(service string) string {
	service = strings.ToLower(service)
	if strings.Contains(service, "virtual machine") || strings.Contains(service, "vm") {
		return "VM"
	}
	if strings.Contains(service, "sql") {
		return "Database"
	}
	if strings.Contains(service, "function") {
		return "Function"
	}
	if strings.Contains(service, "container") {
		return "Container"
	}
	if strings.Contains(service, "storage") {
		return "Storage"
	}
	return "Other"
}

func mapGCPServiceToType(service string) string {
	service = strings.ToLower(service)
	if strings.Contains(service, "compute engine") {
		return "VM"
	}
	if strings.Contains(service, "cloud sql") {
		return "Database"
	}
	if strings.Contains(service, "cloud functions") {
		return "Function"
	}
	if strings.Contains(service, "gke") {
		return "Container"
	}
	if strings.Contains(service, "cloud storage") {
		return "Storage"
	}
	return "Other"
}
