package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	if err := os.MkdirAll(deploymentDir, 0o755); err != nil {
		return fmt.Errorf("failed to create deployment directory: %w", err)
	}

	// Download and extract
	assetPath := filepath.Join(deploymentDir, asset.Name)
	if err := d.github.DownloadAsset(ctx, asset, assetPath); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	// Optional checksum verification
	if d.config.VerifyChecksums {
		d.logger.Printf("Checksum verification enabled, looking for checksums asset")
		// Try to find a checksums file in release assets
		var checksumsAsset *Asset
		for _, a := range release.Assets {
			if strings.HasPrefix(a.Name, strings.TrimSuffix(asset.Name, filepath.Ext(asset.Name))) && strings.Contains(a.Name, "checksums") {
				checksumsAsset = &a
				break
			}
		}
		if checksumsAsset != nil {
			checksumPath := filepath.Join(deploymentDir, checksumsAsset.Name)
			if err := d.github.DownloadAsset(ctx, checksumsAsset, checksumPath); err != nil {
				return fmt.Errorf("failed to download checksums asset: %w", err)
			}
			m, err := ParseChecksums(checksumPath)
			if err != nil {
				return fmt.Errorf("failed to parse checksums: %w", err)
			}
			// lookup by base name
			base := filepath.Base(asset.Name)
			if expected, ok := m[base]; ok {
				if err := VerifyFileSHA256(assetPath, expected); err != nil {
					return fmt.Errorf("checksum verification failed: %w", err)
				}
				d.logger.Printf("Checksum verification passed for %s", base)
			} else {
				return fmt.Errorf("no checksum entry found for %s", base)
			}
		} else {
			d.logger.Printf("VERIFY_CHECKSUMS set but no checksums asset found; aborting")
			return errors.New("checksums verification requested but no checksums asset found")
		}
	}

	// Extract archive based on extension
	d.logger.Printf("Extracting asset: %s", asset.Name)
	var extractErr error
	if strings.HasSuffix(asset.Name, ".tar.gz") || strings.HasSuffix(asset.Name, ".tgz") {
		extractErr = ExtractTarGz(assetPath, deploymentDir)
	} else if strings.HasSuffix(asset.Name, ".zip") {
		extractErr = ExtractZip(assetPath, deploymentDir)
	} else {
		// Not an archive; assume it's a binary. Nothing to extract.
		d.logger.Printf("Asset is not an archive, skipping extraction")
		extractErr = nil
	}
	if extractErr != nil {
		return fmt.Errorf("failed to extract archive: %w", extractErr)
	}

	// Run install command if configured (e.g., poetry install)
	if d.config.RunCommand != "" {
		d.logger.Printf("Running install command: %s", d.config.RunCommand)
		if err := runCommand(deploymentDir, d.config.RunCommand); err != nil {
			return fmt.Errorf("run command failed: %w", err)
		}
		d.logger.Printf("Install command completed successfully")
	}

	// Health check if configured
	if d.config.HealthCheckURL != "" {
		d.logger.Printf("Performing health check on %s", d.config.HealthCheckURL)
		if err := performHealthCheck(d.config.HealthCheckURL, time.Duration(d.config.HealthCheckTimeout)*time.Second); err != nil {
			return fmt.Errorf("health check failed: %w", err)
		}
		d.logger.Printf("Health check passed")
	}

	// Atomically switch symlink
	d.logger.Printf("Switching symlink %s to %s", d.config.CurrentSymlink, deploymentDir)
	if err := switchSymlink(d.config.CurrentSymlink, deploymentDir); err != nil {
		return fmt.Errorf("failed to switch symlink: %w", err)
	}

	// Update state and save
	if d.state.ActiveSlot == "blue" {
		d.state.GreenVersion = release.TagName
	} else {
		d.state.BlueVersion = release.TagName
	}
	d.state.SwitchSlot()
	if err := d.state.SaveState(d.config.StateFile); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Run post-deploy script if configured
	if d.config.PostDeployScript != "" {
		d.logger.Printf("Running post-deploy script: %s", d.config.PostDeployScript)
		if err := runCommand("/", d.config.PostDeployScript); err != nil {
			d.logger.Printf("Warning: post-deploy script failed: %v", err)
		} else {
			d.logger.Printf("Post-deploy script completed successfully")
		}
	}

	d.logger.Printf("Deployment of %s to %s slot completed successfully", release.TagName, inactiveSlot)
	return nil
}

// Rollback performs a rollback to the previous version
func (d *Deployer) Rollback() error {
	currentSlot := d.state.ActiveSlot
	previousSlot := d.state.GetInactiveSlot()
	previousVersion := d.getCurrentVersion()

	d.logger.Printf("Starting rollback from %s slot (version %s) to %s slot",
		currentSlot, previousVersion, previousSlot)

	if d.dryRun {
		d.logger.Printf("DRY RUN: Would rollback to %s slot", previousSlot)
		return nil
	}

	// Switch back to previous slot
	d.state.SwitchSlot()

	// Update symlink to point to previous slot
	previousDir := filepath.Join(d.config.InstallDir, previousSlot)
	if err := switchSymlink(d.config.CurrentSymlink, previousDir); err != nil {
		return fmt.Errorf("failed to switch symlink during rollback: %w", err)
	}

	// Save state
	if err := d.state.SaveState(d.config.StateFile); err != nil {
		return fmt.Errorf("failed to save state during rollback: %w", err)
	}

	// Run post-deploy script if configured
	if d.config.PostDeployScript != "" {
		d.logger.Printf("Running post-deploy script after rollback")
		if err := runCommand("/", d.config.PostDeployScript); err != nil {
			d.logger.Printf("Warning: post-deploy script failed during rollback: %v", err)
		}
	}

	// Validate rollback with health check if configured
	if d.config.HealthCheckURL != "" {
		d.logger.Printf("Validating rollback with health check")
		if err := performHealthCheck(d.config.HealthCheckURL, time.Duration(d.config.HealthCheckTimeout)*time.Second); err != nil {
			return fmt.Errorf("rollback validation failed: %w", err)
		}
	}

	d.logger.Printf("Rollback completed to %s slot", previousSlot)
	return nil
}

// runCommand runs a shell command in the specified working directory
func runCommand(workingDir string, command string) error {
	// Use 'sh -c' for portability
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// performHealthCheck polls the health endpoint until timeout
func performHealthCheck(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 5 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			if resp.Body != nil {
				resp.Body.Close()
			}
			return nil
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("health check failed for %s within %s", url, timeout)
}

// switchSymlink atomically updates the symlink to point to newDir
func switchSymlink(symlinkPath, newDir string) error {
	// Create parent dir for symlink if needed
	if err := os.MkdirAll(filepath.Dir(symlinkPath), 0o755); err != nil {
		return err
	}
	tmpLink := symlinkPath + ".tmp"
	// Remove tmp if exists
	_ = os.Remove(tmpLink)
	if err := os.Symlink(newDir, tmpLink); err != nil {
		return err
	}
	if err := os.Rename(tmpLink, symlinkPath); err != nil {
		// cleanup
		_ = os.Remove(tmpLink)
		return err
	}
	return nil
}
