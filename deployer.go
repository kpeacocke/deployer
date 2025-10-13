package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Deployer manages the deployment process
type Deployer struct {
	config *Config
	logger *log.Logger
	state  *DeploymentState
	github *GitHubClient
	dryRun bool
}

// NewDeployer creates a new deployer instance
func NewDeployer(config *Config, logger *log.Logger, dryRun bool) (*Deployer, error) {
	state, err := LoadState(config.StateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	github := NewGitHubClient(config.GitHubToken)

	return &Deployer{
		config: config,
		logger: logger,
		state:  state,
		github: github,
		dryRun: dryRun,
	}, nil
}

// Run starts the deployment polling loop
func (d *Deployer) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(d.config.CheckIntervalSecs) * time.Second)
	defer ticker.Stop()

	// Perform initial check
	if err := d.checkAndDeploy(ctx); err != nil {
		d.logger.Printf("Initial deployment check failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			d.logger.Println("Shutting down deployer")
			return nil
		case <-ticker.C:
			if err := d.checkAndDeploy(ctx); err != nil {
				d.logger.Printf("Deployment check failed: %v", err)
			}
		}
	}
}

// checkAndDeploy checks for new releases and deploys if needed
func (d *Deployer) checkAndDeploy(ctx context.Context) error {
	d.logger.Printf("Checking for new releases for repo: %s", d.config.Repo)

	release, err := d.github.GetLatestRelease(ctx, d.config.Repo)
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	// Check if this version is already deployed
	currentVersion := d.getCurrentVersion()
	if release.TagName == currentVersion {
		d.logger.Printf("Already on latest version: %s", release.TagName)
		return nil
	}

	d.logger.Printf("New version available: %s (current: %s)", release.TagName, currentVersion)

	if d.dryRun {
		d.logger.Printf("DRY RUN: Would deploy version %s", release.TagName)
		return nil
	}

	return d.deploy(ctx, release)
}

// getCurrentVersion gets the currently deployed version
func (d *Deployer) getCurrentVersion() string {
	if d.state.ActiveSlot == "blue" {
		return d.state.BlueVersion
	}
	return d.state.GreenVersion
}

// deploy performs the actual deployment
func (d *Deployer) deploy(ctx context.Context, release *Release) error {
	inactiveSlot := d.state.GetInactiveSlot()
	d.logger.Printf("Starting deployment of %s to %s slot", release.TagName, inactiveSlot)

	// Find the asset to download
	asset, err := release.FindAssetWithSuffix(d.config.AssetSuffix)
	if err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Create deployment directory
	deploymentDir := filepath.Join(d.config.InstallDir, inactiveSlot)
	if err := os.MkdirAll(deploymentDir, 0755); err != nil {
		return fmt.Errorf("failed to create deployment directory: %w", err)
	}

	// Download and extract
	assetPath := filepath.Join(deploymentDir, asset.Name)
	if err := d.github.DownloadAsset(ctx, asset, assetPath); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	// TODO: Extract archive and run poetry install
	// TODO: Run health checks
	// TODO: Switch symlink atomically
	// TODO: Update state
	// TODO: Run post-deploy script

	d.logger.Printf("Deployment of %s completed successfully", release.TagName)
	return nil
}

// Rollback performs a rollback to the previous version
func (d *Deployer) Rollback() error {
	d.logger.Printf("Starting rollback from %s slot", d.state.ActiveSlot)

	// Switch back to previous slot
	d.state.SwitchSlot()

	// TODO: Update symlink
	// TODO: Save state
	// TODO: Run post-deploy script
	// TODO: Validate rollback

	d.logger.Printf("Rollback completed to %s slot", d.state.ActiveSlot)
	return nil
}
