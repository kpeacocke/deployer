# gh-deployer

[![CI Status](https://github.com/kpeacocke/deployer/actions/workflows/ci.yml/badge.svg)](https://github.com/kpeacocke/deployer/actions/workflows/ci.yml)
[![CodeQL](https://github.com/kpeacocke/deployer/actions/workflows/codeql.yml/badge.svg)](https://github.com/kpeacocke/deployer/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kpeacocke/deployer)](https://goreportcard.com/report/github.com/kpeacocke/deployer)
[![GoDoc](https://godoc.org/github.com/kpeacocke/deployer?status.svg)](https://godoc.org/github.com/kpeacocke/deployer)
[![Release](https://img.shields.io/github/release/kpeacocke/deployer.svg)](https://github.com/kpeacocke/deployer/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Go Version](https://img.shields.io/github/go-mod/go-version/kpeacocke/deployer)](https://github.com/kpeacocke/deployer/blob/main/go.mod)
[![GitHub Package](https://img.shields.io/github/v/release/kpeacocke/deployer?label=package&logo=github)](https://github.com/kpeacocke/deployer/pkgs/npm/gh-deployer)
[![GitHub issues](https://img.shields.io/github/issues/kpeacocke/deployer)](https://github.com/kpeacocke/deployer/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/kpeacocke/deployer)](https://github.com/kpeacocke/deployer/pulls)
[![Downloads](https://img.shields.io/github/downloads/kpeacocke/deployer/total)](https://github.com/kpeacocke/deployer/releases)

A Go-based GitHub release deployer with blue/green deployment, designed to run on Raspberry Pi and launch Python apps using Poetry.

> 🚀 **Automatic releases**: Every push to main automatically creates a new release using semantic versioning!

## Features

- **Automated Deployment**: Polls GitHub for latest releases and deploys automatically
- **Blue/Green Deployment**: Uses separate directories for zero-downtime deployments
- **Archive Support**: Extracts .tar.gz, .tgz, and .zip archives automatically
- **Checksum Verification**: Optional SHA256 checksum verification for security (configurable)
- **Custom Install Commands**: Run Poetry or other install steps during deployment
- **Health Checks**: Validates deployments before switching traffic
- **Atomic Symlink Switching**: Zero-downtime switchover between versions
- **Rollback Support**: Easy rollback to previous version with validation
- **Post-Deploy Hooks**: Optional scripts to run after deployment
- **Systemd Integration**: Startup-safe with systemd service support
- **Structured Logging**: Detailed logging of all deployment steps
- **Dry-Run Mode**: Test deployments without making changes

## Installation

### Quick Install Script (Recommended)

```bash
# Install latest version to /usr/local/bin
curl -fsSL https://raw.githubusercontent.com/kpeacocke/deployer/main/install.sh | bash

# Install specific version
curl -fsSL https://raw.githubusercontent.com/kpeacocke/deployer/main/install.sh | bash -s v1.0.0

# Install to custom location
curl -fsSL https://raw.githubusercontent.com/kpeacocke/deployer/main/install.sh | bash -s latest /opt/bin
```

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/kpeacocke/deployer/releases):

- **Linux AMD64**: `gh-deployer-linux-amd64.tar.gz`
- **Linux ARM64**: `gh-deployer-linux-arm64.tar.gz`
- **Linux ARMv7** (Raspberry Pi): `gh-deployer-linux-armv7.tar.gz`
- **macOS Intel**: `gh-deployer-darwin-amd64.tar.gz`
- **macOS Apple Silicon**: `gh-deployer-darwin-arm64.tar.gz`
- **Windows**: `gh-deployer-windows-amd64.zip`

### Install via Go

```bash
# Install latest version
go install github.com/kpeacocke/deployer@latest

# Install specific version
go install github.com/kpeacocke/deployer@v1.0.0
```

### GitHub Packages

The project is also published as a GitHub Package with pre-built binaries:

```bash
# View package details
# https://github.com/kpeacocke/deployer/pkgs/npm/gh-deployer
```

> 📦 **Go Module**: Available on [pkg.go.dev](https://pkg.go.dev/github.com/kpeacocke/deployer) with full documentation
> 📦 **GitHub Package**: Published to GitHub Packages with every release

### Build from Source

```bash
git clone https://github.com/kpeacocke/deployer.git
cd deployer
make build
```

## Quick Start

1. **Get the binary** (see installation options above)

2. **Configure the deployer:**
   Edit `config.yaml` with your repository and deployment settings:

   ```yaml
   repo: "your-user/your-repo"
   asset_suffix: ".tar.gz"
   github_token: "your-github-token"  # or set GITHUB_TOKEN env var
   install_dir: "/opt/myapp/deployments"
   current_symlink: "/opt/myapp/current"
   ```

3. **Test with dry run:**

   ```bash
   ./gh-deployer --dry-run
   ```

4. **Install and run as systemd service:**

   ```bash
   make install
   make systemd-service
   sudo cp gh-deployer.service /etc/systemd/system/
   sudo systemctl enable gh-deployer
   sudo systemctl start gh-deployer
   ```

## Development

### VS Code Setup (Recommended)

For the best development experience, open this project as a workspace:

```bash
# Clone and open as workspace (optimal extension management)
git clone https://github.com/kpeacocke/deployer.git
code deployer/gh-deployer.code-workspace
```

This workspace configuration:

- ✅ **Enables only Go-relevant extensions** (golang.go, YAML, Markdown, Git tools)
- ✅ **Disables language features** for Python, Ansible, Docker, TypeScript, etc.
- ✅ **Optimizes performance** (disabled minimap, telemetry, file watching exclusions)
- ✅ **Provides consistent setup** for all contributors

> **Note**: If you have many extensions installed globally (Ansible, Godot, Python, etc.), you may need to manually disable them for this workspace. See `.vscode/README.md` for details.

### Build Commands

- **Run tests:** `make test`
- **Test with coverage:** `make test-coverage`
- **Format code:** `make fmt`
- **Lint code:** `make lint`
- **All checks:** `make check`

## Configuration

See `config.yaml` for all configuration options. Key settings:

### Required Settings

- `repo`: GitHub repository to monitor (format: "owner/repo")
- `asset_suffix`: Filter releases by asset name suffix (e.g., ".tar.gz")
- `install_dir`: Root directory for blue/green deployments
- `current_symlink`: Symlink pointing to active deployment
- `state_file`: Path to store deployment state

### Optional Settings

- `check_interval_seconds`: How often to check for new releases (default: 300)
- `github_token`: GitHub API token (or set `GITHUB_TOKEN` env var)
- `run_command`: Command to run after extraction (e.g., "poetry install --no-dev")
- `post_deploy_script`: Script to run after successful deployment
- `verify_checksums`: Enable SHA256 checksum verification (default: false)
- `health_check_url`: URL to check before activating deployment
- `health_check_timeout`: Timeout for health checks in seconds (default: 30)

### Example Configuration

```yaml
repo: "myorg/myapp"
asset_suffix: ".tar.gz"
check_interval_seconds: 300
install_dir: "/opt/myapp/deployments"
current_symlink: "/opt/myapp/current"
run_command: "poetry install --no-dev"
post_deploy_script: "/opt/myapp/scripts/notify-deployment.sh"
state_file: "/opt/myapp/gh-deployer/state.yaml"
verify_checksums: true  # Requires checksums.txt in release
health_check_url: "http://localhost:8000/health"
health_check_timeout: 30
```

## Architecture

The deployer implements a blue/green deployment strategy:

- Two deployment slots (`blue` and `green`) with separate Poetry virtual environments
- Atomic symlink switching for zero-downtime deployments
- State persistence in `state.yaml`
- Rollback capability to previous version
- Health checks before activation

For detailed architecture information, see `.github/copilot-instructions.md`.

## Automatic Releases 🚀

This project uses **automatic semantic versioning** - every push to `main` triggers a release if there are new features or fixes!

### Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/) for automatic version bumping:

- `feat:` - New feature → **Minor version** (e.g., 1.0.0 → 1.1.0)
- `fix:` - Bug fix → **Patch version** (e.g., 1.0.0 → 1.0.1)
- `perf:` - Performance improvement → **Patch version**
- `BREAKING CHANGE:` - Breaking change → **Major version** (e.g., 1.0.0 → 2.0.0)

### Examples

```bash
git commit -m "feat: add health check endpoint"     # → 1.1.0
git commit -m "fix: resolve memory leak"            # → 1.0.1
git commit -m "feat!: redesign configuration API"   # → 2.0.0
```

Every successful commit to main automatically:

- ✅ Runs full test suite and linting
- 🏷️ Creates a new semantic version tag
- 📦 Builds binaries for all platforms
- 🚀 Publishes GitHub release with assets
- 📖 Updates CHANGELOG.md automatically

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Use conventional commits: `git commit -m "feat: add amazing feature"`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

Your changes will be automatically released when merged to main!
