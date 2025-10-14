#!/bin/bash
set -e

# gh-deployer installation script
VERSION="${1:-latest}"
INSTALL_DIR="${2:-/usr/local/bin}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  armv7l) ARCH="armv7" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

case $OS in
  linux) PLATFORM="linux" ;;
  darwin) PLATFORM="darwin" ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Get latest version if not specified
if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -s https://api.github.com/repos/kpeacocke/deployer/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
fi

# Download and install
BINARY_NAME="gh-deployer-${PLATFORM}-${ARCH}"
DOWNLOAD_URL="https://github.com/kpeacocke/deployer/releases/download/${VERSION}/${BINARY_NAME}.tar.gz"

echo "Downloading gh-deployer ${VERSION} for ${PLATFORM}/${ARCH}..."
curl -sL "$DOWNLOAD_URL" | tar -xz

echo "Installing to ${INSTALL_DIR}/gh-deployer..."

# Create install directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
  if [ -w "$(dirname "$INSTALL_DIR")" ]; then
    mkdir -p "$INSTALL_DIR"
  else
    sudo mkdir -p "$INSTALL_DIR"
  fi
fi

# Install the binary
if [ -w "$INSTALL_DIR" ]; then
  mv "$BINARY_NAME" "${INSTALL_DIR}/gh-deployer"
  chmod +x "${INSTALL_DIR}/gh-deployer"
else
  sudo mv "$BINARY_NAME" "${INSTALL_DIR}/gh-deployer"
  sudo chmod +x "${INSTALL_DIR}/gh-deployer"
fi

echo "gh-deployer installed successfully!"
echo "Run 'gh-deployer --help' to get started."
