package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	configContent := `repo: "test/repo"
asset_suffix: ".tar.gz"
check_interval_seconds: 60
install_dir: "/tmp/test"
current_symlink: "/tmp/current"
run_command: "poetry run python main.py"
post_deploy_script: "deploy.sh"
state_file: "/tmp/state.yaml"
health_check_timeout: 30
logging:
  level: "debug"
  file: "/tmp/test.log"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test values
	if config.Repo != "test/repo" {
		t.Errorf("Expected repo 'test/repo', got '%s'", config.Repo)
	}

	if config.CheckIntervalSecs != 60 {
		t.Errorf("Expected check interval 60, got %d", config.CheckIntervalSecs)
	}

	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.Logging.Level)
	}
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Create minimal config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "minimal-config.yaml")

	configContent := `repo: "test/repo"
asset_suffix: ".tar.gz"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test defaults
	if config.CheckIntervalSecs != 300 {
		t.Errorf("Expected default check interval 300, got %d", config.CheckIntervalSecs)
	}

	if config.HealthCheckTimeout != 30 {
		t.Errorf("Expected default health check timeout 30, got %d", config.HealthCheckTimeout)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", config.Logging.Level)
	}
}

func TestLoadConfigEnvironmentOverride(t *testing.T) {
	// Set environment variable
	os.Setenv("GITHUB_TOKEN", "test-token-123")
	defer os.Unsetenv("GITHUB_TOKEN")

	// Create config file without token
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `repo: "test/repo"
asset_suffix: ".tar.gz"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test environment override
	if config.GitHubToken != "test-token-123" {
		t.Errorf("Expected GitHub token from environment, got '%s'", config.GitHubToken)
	}
}
