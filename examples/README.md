# Deployment Examples

This directory contains example configurations and scripts for deploying applications with gh-deployer.

## Displayboard Python/Poetry Project

Example setup for deploying a Python project using Poetry to a Raspberry Pi or Linux system.

### Configuration

**config.yaml:**
```yaml
repo: "kpeacocke/displayboard"
asset_suffix: ".tar.gz"
install_dir: "/opt/displayboard/deployments"
current_symlink: "/opt/displayboard/current"
run_command: "poetry install --no-dev"
post_deploy_script: "/opt/displayboard/scripts/restart-displayboard.sh"
state_file: "/opt/displayboard/gh-deployer/state.yaml"
```

### Setup Steps

1. **Create required directories:**
   ```bash
   sudo mkdir -p /opt/displayboard/{deployments,scripts}
   sudo mkdir -p /var/log/displayboard
   sudo chown -R $USER:$USER /opt/displayboard
   sudo chown -R $USER:$USER /var/log/displayboard
   ```

2. **Copy the post-deploy script:**
   ```bash
   sudo cp examples/restart-displayboard.sh /opt/displayboard/scripts/
   sudo chmod +x /opt/displayboard/scripts/restart-displayboard.sh
   ```

3. **Install as systemd service (recommended):**
   ```bash
   # Copy the service file
   sudo cp examples/displayboard.service /etc/systemd/system/
   
   # Reload systemd
   sudo systemctl daemon-reload
   
   # Enable the service to start on boot
   sudo systemctl enable displayboard
   
   # Start the service
   sudo systemctl start displayboard
   
   # Check status
   sudo systemctl status displayboard
   ```

4. **Configure gh-deployer:**
   ```bash
   # Copy your config
   cp config.yaml /opt/displayboard/gh-deployer/config.yaml
   
   # Edit with your actual repository
   nano /opt/displayboard/gh-deployer/config.yaml
   ```

5. **Run gh-deployer:**
   ```bash
   # Test with dry-run first
   gh-deployer --config /opt/displayboard/gh-deployer/config.yaml --dry-run
   
   # Run normally
   gh-deployer --config /opt/displayboard/gh-deployer/config.yaml
   ```

### How It Works

1. **gh-deployer monitors** your GitHub repository for new releases
2. **Downloads and extracts** the release asset to an inactive slot (blue/green)
3. **Runs `poetry install --no-dev`** in the deployment directory to install dependencies
4. **Switches the symlink** `/opt/displayboard/current` to point to the new deployment
5. **Executes the post-deploy script** which:
   - Verifies the Poetry environment exists
   - Restarts the systemd service (or manually restarts the app)
   - Validates the service is running

### Manual Application Start

If not using systemd, you can manually start displayboard:

```bash
cd /opt/displayboard/current
sudo $(poetry env info --path)/bin/python -m displayboard.main -d
```

### Monitoring and Logs

```bash
# View deployer logs
journalctl -u gh-deployer -f

# View displayboard logs
journalctl -u displayboard -f

# View deployment script logs
tail -f /var/log/displayboard/deployer.log
```

### Rollback

If a deployment fails, you can rollback to the previous version:

```bash
# The deployer keeps both blue and green slots
# Manually switch the symlink back
ln -sfn /opt/displayboard/deployments/blue /opt/displayboard/current

# Restart the service
sudo systemctl restart displayboard
```

### Troubleshooting

**Poetry environment not found:**
```bash
cd /opt/displayboard/current
poetry install --no-dev
poetry env info
```

**Permission issues:**
```bash
sudo chown -R $USER:$USER /opt/displayboard
sudo chmod +x /opt/displayboard/scripts/restart-displayboard.sh
```

**Service not starting:**
```bash
sudo systemctl status displayboard
journalctl -u displayboard -n 50
```
