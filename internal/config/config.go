package config

type SyntheticUnitRule struct {
	UnitsPerInstance int `json:"unitsPerInstance"`
}

type SyntheticUnitsConfig struct {
	Rules map[string]SyntheticUnitRule `json:"rules"`
}

type ProvidersConfig struct {
	AWS struct {
		Enabled bool   `json:"enabled"`
		Regions []string `json:"regions"`
	} `json:"aws"`
	Azure struct {
		Enabled bool `json:"enabled"`
	} `json:"azure"`
	GCP struct {
		Enabled bool `json:"enabled"`
	} `json:"gcp"`
}

type BillingConfig struct {
	AWS struct {
		FilePath string `json:"filePath"`
		Format   string `json:"format"`
		Period   string `json:"period"`
	} `json:"aws"`
	Azure struct {
		FilePath string `json:"filePath"`
		Format   string `json:"format"`
		Period   string `json:"period"`
	} `json:"azure"`
	GCP struct {
		FilePath string `json:"filePath"`
		Format   string `json:"format"`
		Period   string `json:"period"`
	} `json:"gcp"`
}

type OutputConfig struct {
	Format                   string `json:"format"`
	Filename                 string `json:"filename"`
	IncludeEphemeralResources bool  `json:"includeEphemeralResources"`
	IncludeBillingMetrics    bool  `json:"includeBillingMetrics"`
}

type Config struct {
	Providers      ProvidersConfig      `json:"providers"`
	Billing        BillingConfig        `json:"billing"`
	SyntheticUnits SyntheticUnitsConfig `json:"syntheticUnits"`
	Output         OutputConfig         `json:"output"`
}
