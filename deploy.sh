#!/bin/bash
set -e

# Post-deployment script for gh-deployer
# This script is executed after a successful deployment

echo "Post-deploy script running"
echo "Deployment completed at: $(date)"

# Example: Restart related services
# systemctl restart myapp
# systemctl reload nginx

# Example: Health check
# curl -f http://localhost:8080/health || exit 1

# Example: Send notification
# curl -X POST -H 'Content-type: application/json' \
#   --data '{"text":"Deployment completed successfully"}' \
#   "$SLACK_WEBHOOK_URL"

echo "Post-deploy script completed successfully"
