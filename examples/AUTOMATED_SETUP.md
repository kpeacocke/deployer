# Complete Setup Guide for Automated Deployments

This guide shows how to set up gh-deployer to run as a service on a Raspberry Pi, automatically deploying your Python application whenever you push a new release to GitHub.

## Overview

Once configured, this is what happens automatically:

1. You push code to GitHub and create a release
2. gh-deployer (running as a service) detects the new release
3. Downloads and extracts to a new deployment slot
4. Runs `poetry install --no-dev` to install dependencies
5. Copies your `.env` file to the new deployment
6. Switches the symlink to the new version
7. Restarts your displayboard application
8. You never SSH into the Pi! ✨

## Initial Setup (One-Time)

### 1. Install gh-deployer on the Pi

```bash
# Install gh-deployer
curl -fsSL https://raw.githubusercontent.com/kpeacocke/deployer/main/install.sh | bash

# Verify installation
gh-deployer --version
```

### 2. Create Directory Structure

```bash
# Create all required directories
sudo mkdir -p /opt/displayboard/{deployments,scripts,config,gh-deployer}
sudo mkdir -p /var/log/displayboard

# Set ownership to your user (pi)
sudo chown -R $USER:$USER /opt/displayboard
sudo chown -R $USER:$USER /var/log/displayboard
```

### 3. Create Your .env File

```bash
# Create your persistent .env file (this stays between deployments)
nano /opt/displayboard/config/.env
```

Add your environment variables:

```env
# Example .env for displayboard
DATABASE_URL=postgresql://user:pass@localhost/dbname
API_KEY=your-secret-api-key
DEBUG=false
LOG_LEVEL=info
```

```bash
# Secure the .env file
chmod 600 /opt/displayboard/config/.env
```

### 4. Set Up GitHub Token (Optional but Recommended)

```bash
# Create a GitHub personal access token at:
# https://github.com/settings/tokens
# Only needs "public_repo" access for public repos

# Save it to a file
echo "ghp_yourGitHubTokenHere" > /opt/displayboard/config/github-token
chmod 600 /opt/displayboard/config/github-token
```

### 5. Configure gh-deployer

```bash
# Create the config file
cat > /opt/displayboard/gh-deployer/config.yaml << 'EOF'
repo: "kpeacocke/displayboard"
asset_suffix: ".tar.gz"
check_interval_seconds: 300  # Check every 5 minutes
install_dir: "/opt/displayboard/deployments"
current_symlink: "/opt/displayboard/current"
run_command: "/home/$USER/.local/bin/poetry install --without dev"
post_deploy_script: "/opt/displayboard/scripts/restart-displayboard.sh"
state_file: "/opt/displayboard/gh-deployer/state.yaml"
github_token: ""  # Will be read from env var GITHUB_TOKEN
EOF

# Replace $USER with actual username
sed -i "s/\$USER/$USER/" /opt/displayboard/gh-deployer/config.yaml
```

### 6. Install the Scripts

After the next release of gh-deployer (v1.6.0+), get the example files:

```bash
# Download the latest release
cd /tmp
LATEST=$(curl -s https://api.github.com/repos/kpeacocke/deployer/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -L "https://github.com/kpeacocke/deployer/archive/$LATEST.tar.gz" | tar xz

# Copy example files
cd deployer-*/examples
cp restart-displayboard.sh /opt/displayboard/scripts/
chmod +x /opt/displayboard/scripts/restart-displayboard.sh

# Install systemd services
sudo cp gh-deployer.service /etc/systemd/system/
sudo cp displayboard.service /etc/systemd/system/
```

### 7. Set Up Hardware Permissions (For NeoPixel/GPIO Access)

If your displayboard uses NeoPixels or GPIO, set up permissions so it runs without sudo:

```bash
# Download permission setup files from your displayboard repo
cd /tmp
curl -L "https://github.com/kpeacocke/displayboard/archive/main.tar.gz" | tar xz
cd displayboard-main/permissions

# Run the setup script
sudo bash setup-permissions.sh

# Log out and back in for group changes to take effect
exit
# Then SSH back in

# Verify you're in the gpio and video groups
groups
# Should show: gpio video
```

This grants your user access to GPIO, /dev/mem, and DMA without requiring sudo.

### 8. Install Poetry (If Not Already Installed)

```bash
# Install Poetry
curl -sSL https://install.python-poetry.org | python3 -

# Add to PATH
export PATH="$HOME/.local/bin:$PATH"
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc

# Configure Poetry to create venvs in project directory
poetry config virtualenvs.in-project true

# Verify
poetry --version
```

### 9. Configure GitHub Token

```bash
# Edit the gh-deployer service to read the token (optional)
sudo nano /etc/systemd/system/gh-deployer.service
```

Make sure this line is present:

```ini
Environment="GITHUB_TOKEN_FILE=/opt/displayboard/config/github-token"
```

Or set it directly in the config:

```bash
# Add token to config
TOKEN=$(cat /opt/displayboard/config/github-token)
sed -i "s/github_token: \"\"/github_token: \"$TOKEN\"/" /opt/displayboard/gh-deployer/config.yaml
```

### 10. Fix Username in Services

```bash
# Replace 'pi' with your actual username in both service files
sudo sed -i "s/User=pi/User=$USER/" /etc/systemd/system/gh-deployer.service
sudo sed -i "s/Group=pi/Group=$USER/" /etc/systemd/system/gh-deployer.service
sudo sed -i "s/User=pi/User=$USER/" /etc/systemd/system/displayboard.service
sudo sed -i "s/Group=pi/Group=$USER/" /etc/systemd/system/displayboard.service
```

### 11. Enable and Start Services

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable services to start on boot
sudo systemctl enable gh-deployer
sudo systemctl enable displayboard

# Start gh-deployer (it will handle displayboard deployments)
sudo systemctl start gh-deployer

# Check status
sudo systemctl status gh-deployer
sudo systemctl status displayboard
```

## Monitoring

### View Logs

```bash
# Watch gh-deployer logs (shows deployment activity)
journalctl -u gh-deployer -f

# Watch displayboard logs (shows your application output)
journalctl -u displayboard -f

# View deployment script logs
tail -f /var/log/displayboard/deployer.log
```

### Check Deployment Status

```bash
# See current deployment
ls -la /opt/displayboard/current

# View state file
cat /opt/displayboard/gh-deployer/state.yaml

# Check both slots
ls -la /opt/displayboard/deployments/
```

## Creating Releases (Your Workflow)

### Using Conventional Commits (Automatic)

```bash
# Your normal workflow - gh-deployer automates the rest!
git add .
git commit -m "feat: add new dashboard widget"
git push

# GitHub Actions will automatically:
# 1. Run tests
# 2. Create a new release with semantic versioning
# 3. Build and upload the deployment archive
# 4. gh-deployer on your Pi will detect it and deploy!
```

### Manual Release

```bash
# Create a release archive
cd your-displayboard-project
tar -czf displayboard-deployment-v1.0.0.tar.gz .

# Create a GitHub release and upload the archive
gh release create v1.0.0 displayboard-deployment-v1.0.0.tar.gz
```

## How .env File is Handled

The `.env` file is stored in `/opt/displayboard/config/.env` (outside the deployment directories).

On every deployment:

1. New code is extracted to `/opt/displayboard/deployments/{blue|green}`
2. The post-deploy script copies `.env` from the config directory
3. Your application loads the `.env` from its current directory
4. When you need to update environment variables, edit `/opt/displayboard/config/.env`
5. Restart the service: `sudo systemctl restart displayboard`

## Updating Environment Variables

```bash
# Edit the persistent .env file
nano /opt/displayboard/config/.env

# Restart to apply changes (no redeployment needed)
sudo systemctl restart displayboard
```

## Troubleshooting

### gh-deployer not detecting releases

```bash
# Check logs
journalctl -u gh-deployer -n 100

# Check GitHub token
cat /opt/displayboard/config/github-token

# Test manually
gh-deployer --config /opt/displayboard/gh-deployer/config.yaml --dry-run
```

### Deployment failed

```bash
# Check deployment logs
tail -n 100 /var/log/displayboard/deployer.log

# Check if .env was copied
ls -la /opt/displayboard/current/.env

# Check Poetry environment
cd /opt/displayboard/current
poetry env info
```

### Application not starting

```bash
# Check service status
sudo systemctl status displayboard

# View recent logs
journalctl -u displayboard -n 50

# Check if .env exists
ls -la /opt/displayboard/current/.env

# Try starting manually
cd /opt/displayboard/current
$(poetry env info --path)/bin/python -m displayboard.main -d
```

### Rollback to previous version

```bash
# Stop the service
sudo systemctl stop displayboard

# Check current state
cat /opt/displayboard/gh-deployer/state.yaml

# Manually switch to the other slot
# If currently on green, switch to blue
sudo ln -sfn /opt/displayboard/deployments/blue /opt/displayboard/current

# Copy .env to old deployment
cp /opt/displayboard/config/.env /opt/displayboard/current/.env

# Start the service
sudo systemctl start displayboard
```

## Security Best Practices

1. **Never commit .env to git** - keep secrets in `/opt/displayboard/config/.env`
2. **Restrict .env permissions**: `chmod 600 /opt/displayboard/config/.env`
3. **Secure GitHub token**: `chmod 600 /opt/displayboard/config/github-token`
4. **Use minimal GitHub token scope** - only `public_repo` for public repos
5. **Regular updates**: Keep gh-deployer updated for security patches

## Complete File Structure

```text
/opt/displayboard/
├── config/
│   ├── .env                          # Your persistent environment variables
│   └── github-token                  # GitHub API token
├── deployments/
│   ├── blue/                         # Blue deployment slot
│   │   ├── .env -> copied from config
│   │   ├── pyproject.toml
│   │   └── ... your app files
│   └── green/                        # Green deployment slot
│       ├── .env -> copied from config
│       ├── pyproject.toml
│       └── ... your app files
├── current -> deployments/green      # Symlink to active deployment
├── scripts/
│   └── restart-displayboard.sh       # Post-deploy script
└── gh-deployer/
    ├── config.yaml                   # Deployer configuration
    └── state.yaml                    # Deployment state (auto-generated)

/var/log/displayboard/
└── deployer.log                      # Deployment logs

/etc/systemd/system/
├── gh-deployer.service               # Deployer service
└── displayboard.service              # Your app service
```

## What Happens When You're Away

1. You push code and create a release from anywhere (laptop, phone via GitHub web)
2. Within 5 minutes, gh-deployer on your Pi detects it
3. Downloads to inactive slot (blue or green)
4. Runs `poetry install`
5. Copies your `.env` file
6. Switches symlink
7. Restarts your app
8. Your displayboard is running the new version!

**You never need to SSH into the Pi!** 🎉
