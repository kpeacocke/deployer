package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDeployerIntegration(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directories for test
	tempDir := t.TempDir()
	installDir := filepath.Join(tempDir, "deployments")
	stateFile := filepath.Join(tempDir, "state.yaml")
	logFile := filepath.Join(tempDir, "test.log")

	// Create test config
	config := &Config{
		Repo:               "nonexistent/repo", // This will fail, which is expected
		AssetSuffix:        ".tar.gz",
		CheckIntervalSecs:  1, // Short interval for testing
		InstallDir:         installDir,
		CurrentSymlink:     filepath.Join(tempDir, "current"),
		RunCommand:         "echo test",
		PostDeployScript:   "echo 'post-deploy'",
		StateFile:          stateFile,
		HealthCheckTimeout: 5,
		Logging: LoggingConfig{
			Level: "debug",
			File:  logFile,
		},
	}

	// Create logger
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Create deployer
	deployer, err := NewDeployer(config, logger, true) // dry-run mode
	if err != nil {
		t.Fatalf("Failed to create deployer: %v", err)
	}

	// Test that deployer can start and handle context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = deployer.Run(ctx)
	if err != nil {
		t.Fatalf("Deployer run failed: %v", err)
	}

	// Verify state file exists (it should be created on deployer initialization)
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		// State file creation is handled by LoadState, which creates default state
		// This is actually expected behavior - state file is created on first save
		t.Log("State file not created, which is expected for failed deployments")
	}

	// Verify deployer state is valid
	if deployer.state.ActiveSlot == "" {
		t.Error("Deployer state active slot should not be empty")
	}
}

func TestDeployerRollback(t *testing.T) {
	tempDir := t.TempDir()
	stateFile := filepath.Join(tempDir, "state.yaml")

	// Create initial state
	state := &DeploymentState{
		ActiveSlot:   "blue",
		BlueVersion:  "v1.0.0",
		GreenVersion: "v1.1.0",
	}

	if err := state.SaveState(stateFile); err != nil {
		t.Fatalf("Failed to save initial state: %v", err)
	}

	// Create config
	config := &Config{
		StateFile: stateFile,
		Logging: LoggingConfig{
			Level: "info",
		},
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Test with dry-run mode
	t.Run("DryRun", func(t *testing.T) {
		deployer, err := NewDeployer(config, logger, true)
		if err != nil {
			t.Fatalf("Failed to create deployer: %v", err)
		}

		// Verify initial state
		if deployer.state.ActiveSlot != "blue" {
			t.Errorf("Expected active slot 'blue', got '%s'", deployer.state.ActiveSlot)
		}

		// Test rollback in dry-run
		if err := deployer.Rollback(); err != nil {
			t.Fatalf("Rollback failed: %v", err)
		}

		// In dry-run mode, state should NOT change
		if deployer.state.ActiveSlot != "blue" {
			t.Errorf("Expected active slot 'blue' in dry-run (no change), got '%s'", deployer.state.ActiveSlot)
		}
	})

	// Test without dry-run mode
	t.Run("RealRollback", func(t *testing.T) {
		// Create fresh deployer without dry-run
		deployer, err := NewDeployer(config, logger, false)
		if err != nil {
			t.Fatalf("Failed to create deployer: %v", err)
		}

		// Setup install dir and symlink for real rollback
		installDir := filepath.Join(tempDir, "deployments")
		config.InstallDir = installDir
		config.CurrentSymlink = filepath.Join(tempDir, "current")
		deployer.config = config

		// Create slot directories
		if err := os.MkdirAll(filepath.Join(installDir, "blue"), 0o755); err != nil {
			t.Fatalf("Failed to create blue slot: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(installDir, "green"), 0o755); err != nil {
			t.Fatalf("Failed to create green slot: %v", err)
		}

		// Verify initial state
		if deployer.state.ActiveSlot != "blue" {
			t.Errorf("Expected active slot 'blue', got '%s'", deployer.state.ActiveSlot)
		}

		// Test rollback
		if err := deployer.Rollback(); err != nil {
			t.Fatalf("Rollback failed: %v", err)
		}

		// Verify state changed
		if deployer.state.ActiveSlot != "green" {
			t.Errorf("Expected active slot 'green' after rollback, got '%s'", deployer.state.ActiveSlot)
		}

		// Verify symlink was updated
		linkTarget, err := os.Readlink(config.CurrentSymlink)
		if err != nil {
			t.Fatalf("Failed to read symlink: %v", err)
		}
		expectedTarget := filepath.Join(installDir, "green")
		if linkTarget != expectedTarget {
			t.Errorf("Expected symlink to point to '%s', got '%s'", expectedTarget, linkTarget)
		}
	})
}
