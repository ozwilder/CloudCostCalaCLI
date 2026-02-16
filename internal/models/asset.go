package models

type Asset struct {
	ID                   string                 `json:"id"`
	Type                 string                 `json:"type"` // VM, Database, Container, Storage, Function
	Name                 string                 `json:"name"`
	Cloud                string                 `json:"cloud"` // AWS, Azure, GCP
	Project              string                 `json:"project"`
	CurrentInstanceCount int                    `json:"current_instance_count"`
	Metadata             map[string]interface{} `json:"metadata"`
	SourceType           string                 `json:"source_type"` // inventory or billing
}

type BillingRecord struct {
	ServiceName    string
	ResourceType   string // VM, Database, Container, etc.
	ResourceID     string
	InstanceHours  float64
	TimePeriod     string // YYYY-MM
	Region         string
	Project        string
	Metadata       map[string]string
}

type EnrichedAsset struct {
	AssetType             string
	CurrentlyDeployed     int
	AverageInstancesPerHr float64
	HasEphemeralUsage     bool
	CalculatedUnits       int
}

type AggregatedOutput struct {
	AssetType             string
	CurrentCount          int
	EphemeralCount        int
	AvgInstancesPerHour   float64
	SyntheticUnits        int
}
