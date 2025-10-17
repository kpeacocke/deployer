package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Repo               string        `yaml:"repo"`
	AssetSuffix        string        `yaml:"asset_suffix"`
	CheckIntervalSecs  int           `yaml:"check_interval_seconds"`
	InstallDir         string        `yaml:"install_dir"`
	CurrentSymlink     string        `yaml:"current_symlink"`
	RunCommand         string        `yaml:"run_command"`
	PostDeployScript   string        `yaml:"post_deploy_script"`
	StateFile          string        `yaml:"state_file"`
	GitHubToken        string        `yaml:"github_token,omitempty"`
	HealthCheckURL     string        `yaml:"health_check_url,omitempty"`
	HealthCheckTimeout int           `yaml:"health_check_timeout"`
	VerifyChecksums    bool          `yaml:"verify_checksums"`
	Logging            LoggingConfig `yaml:"logging"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    string `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// LoadConfig loads configuration from the specified file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{
		// Set defaults
		CheckIntervalSecs:  300,
		HealthCheckTimeout: 30,
		Logging: LoggingConfig{
			Level: "info",
		},
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables if present
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		config.GitHubToken = token
	}
	if os.Getenv("VERIFY_CHECKSUMS") == "true" {
		config.VerifyChecksums = true
	}

	// Validate required fields and set intelligent defaults
	if config.StateFile == "" {
		// Default to state.yaml in the install directory
		if config.InstallDir != "" {
			config.StateFile = config.InstallDir + "/state.yaml"
		} else {
			config.StateFile = "./state.yaml"
		}
	}

	// Validate required configuration
	if config.Repo == "" {
		return nil, fmt.Errorf("repo is required in configuration")
	}
	if config.InstallDir == "" {
		return nil, fmt.Errorf("install_dir is required in configuration")
	}
	if config.CurrentSymlink == "" {
		return nil, fmt.Errorf("current_symlink is required in configuration")
	}

	return config, nil
}
