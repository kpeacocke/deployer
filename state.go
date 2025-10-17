package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DeploymentState represents the current deployment state
type DeploymentState struct {
	ActiveSlot   string `yaml:"active_slot"`
	BlueVersion  string `yaml:"blue_version"`
	GreenVersion string `yaml:"green_version"`
}

// LoadState loads the deployment state from file
func LoadState(path string) (*DeploymentState, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default state if file doesn't exist
		return &DeploymentState{
			ActiveSlot:   "blue",
			BlueVersion:  "",
			GreenVersion: "",
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state DeploymentState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// SaveState saves the deployment state to file
func (s *DeploymentState) SaveState(path string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create state file directory: %w", err)
		}
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// GetInactiveSlot returns the inactive deployment slot
func (s *DeploymentState) GetInactiveSlot() string {
	if s.ActiveSlot == "blue" {
		return "green"
	}
	return "blue"
}

// SwitchSlot switches the active slot
func (s *DeploymentState) SwitchSlot() {
	s.ActiveSlot = s.GetInactiveSlot()
}
