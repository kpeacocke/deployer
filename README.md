# gh-deployer

A Go-based GitHub release deployer with blue/green deployment, designed to run on Raspberry Pi and launch Python apps using Poetry.

## Features

- Polls GitHub for latest release
- Uses separate Poetry venvs for blue and green slots
- Atomic symlink switching
- Optional post-deploy hook
- Startup-safe with systemd
- Structured logging and monitoring
- Health checks and rollback support

## Quick Start

1. **Build the application:**
   ```bash
   make build
   ```

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

- **Run tests:** `make test`
- **Test with coverage:** `make test-coverage`
- **Format code:** `make fmt`
- **Lint code:** `make lint`
- **All checks:** `make check`

## Configuration

See `config.yaml` for all configuration options. Key settings:

- `repo`: GitHub repository to monitor (format: "owner/repo")
- `asset_suffix`: Filter releases by asset name suffix (e.g., ".tar.gz")
- `check_interval_seconds`: How often to check for new releases (default: 300)
- `install_dir`: Root directory for blue/green deployments
- `current_symlink`: Symlink pointing to active deployment
- `post_deploy_script`: Script to run after successful deployment

## Architecture

The deployer implements a blue/green deployment strategy:
- Two deployment slots (`blue` and `green`) with separate Poetry virtual environments
- Atomic symlink switching for zero-downtime deployments
- State persistence in `state.yaml`
- Rollback capability to previous version
- Health checks before activation

For detailed architecture information, see `.github/copilot-instructions.md`.
