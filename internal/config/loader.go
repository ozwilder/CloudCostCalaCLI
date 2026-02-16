package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required rules
	if cfg.SyntheticUnits.Rules == nil {
		cfg.SyntheticUnits.Rules = make(map[string]SyntheticUnitRule)
	}

	return &cfg, nil
}
