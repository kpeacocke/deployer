#!/bin/bash
set -e

REPO="kpeacocke/deployer"
INSTALL_DIR="/opt/displayboard/gh-deployer"

echo "Installing gh-deployer from latest release..."

# Get latest release info
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")
VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep "browser_download_url.*tar.gz" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find release asset"
    exit 1
fi

echo "Downloading version ${VERSION}..."
echo "URL: ${DOWNLOAD_URL}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf $TMP_DIR' EXIT

# Download and extract
curl -L "$DOWNLOAD_URL" | tar -xzf - -C "$TMP_DIR"

# Stop service if running
if systemctl is-active --quiet gh-deployer; then
    echo "Stopping gh-deployer service..."
    sudo systemctl stop gh-deployer
fi

# Install binary
echo "Installing to ${INSTALL_DIR}..."
sudo mkdir -p "$INSTALL_DIR"
sudo cp "$TMP_DIR/gh-deployer" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/gh-deployer"

# Copy post-deploy script if it exists
if [ -f "$TMP_DIR/restart-displayboard.sh" ]; then
    echo "Installing post-deploy script..."
    sudo cp "$TMP_DIR/restart-displayboard.sh" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/restart-displayboard.sh"
fi

# Restart service if it was running
if systemctl is-enabled --quiet gh-deployer 2>/dev/null; then
    echo "Starting gh-deployer service..."
    sudo systemctl start gh-deployer
    echo "✅ gh-deployer updated to ${VERSION} and restarted"
else
    echo "✅ gh-deployer installed to ${INSTALL_DIR}/gh-deployer"
    echo "Note: Service not enabled. Run 'sudo systemctl enable --now gh-deployer' to start."
fi
