#!/bin/bash
set -e

# Post-deployment script for displayboard Python application
# This script restarts the displayboard service after a successful deployment

APP_DIR="/opt/displayboard/current"
SERVICE_NAME="displayboard"
LOG_FILE="/var/log/displayboard/deployer.log"
ENV_SOURCE="/opt/displayboard/config/.env"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

log "Starting post-deployment tasks for displayboard..."

# Navigate to the deployment directory
cd "$APP_DIR" || {
    log "ERROR: Failed to change to directory $APP_DIR"
    exit 1
}

# Copy .env file from persistent location if it exists
if [ -f "$ENV_SOURCE" ]; then
    log "Copying .env file from $ENV_SOURCE"
    cp "$ENV_SOURCE" "$APP_DIR/.env"
    chmod 600 "$APP_DIR/.env"
    log "✓ .env file copied successfully"
else
    log "WARNING: No .env file found at $ENV_SOURCE"
    log "Application may fail if it requires environment variables"
fi

# Verify poetry environment exists
if ! poetry env info --path &>/dev/null; then
    log "ERROR: Poetry environment not found"
    exit 1
fi

POETRY_ENV=$(poetry env info --path)
log "Poetry environment: $POETRY_ENV"

# Check if running as systemd service
if systemctl is-active --quiet "$SERVICE_NAME"; then
    log "Restarting systemd service: $SERVICE_NAME"
    sudo systemctl restart "$SERVICE_NAME"
    
    # Wait a moment and verify service started
    sleep 2
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log "✓ Service $SERVICE_NAME restarted successfully"
    else
        log "ERROR: Service $SERVICE_NAME failed to start"
        sudo systemctl status "$SERVICE_NAME" --no-pager
        exit 1
    fi
else
    # If not running as systemd service, start manually
    log "Starting displayboard manually (not running as systemd service)"
    
    # Kill any existing displayboard processes
    if pgrep -f "displayboard.main" > /dev/null; then
        log "Stopping existing displayboard processes..."
        sudo pkill -f "displayboard.main" || true
        sleep 2
    fi
    
    # Start the application in the background
    log "Starting displayboard application..."
    nohup sudo "$POETRY_ENV/bin/python" -m displayboard.main -d >> "$LOG_FILE" 2>&1 &
    
    # Store the PID
    echo $! > /var/run/displayboard.pid
    log "✓ Displayboard started with PID $(cat /var/run/displayboard.pid)"
fi

log "Post-deployment tasks completed successfully!"
