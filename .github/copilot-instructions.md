# AI Agent Instructions for gh-deployer

## Project Overview
This is a Go-based GitHub release deployer implementing blue/green deployment for Python applications. The deployer is designed to run on Raspberry Pi and manage Poetry-based Python apps with zero-downtime deployments.

## Architecture & Core Components

### State Management Pattern
- **`state.yaml`**: Single source of truth for deployment state with `active_slot`, `blue_version`, `green_version`
- **Blue/Green Slots**: Two separate Poetry venvs (`blue` and `green`) for atomic switching
- **Symlink Strategy**: `current_symlink` points to active deployment directory for zero-downtime switches

### Key File Relationships
- `config.yaml` → Runtime configuration (repo, intervals, paths)
- `state.yaml` → Deployment state persistence 
- `deploy.sh` → Post-deployment hooks (systemd service restarts, health checks)
- `main.go` → Service entry point (currently placeholder - needs implementation)

## Development Workflows

### Configuration Management
```yaml
# config.yaml structure is production-ready
repo: "your-user/your-repo"           # GitHub repo to monitor
asset_suffix: ".tar.gz"              # Release asset filter
install_dir: "/opt/myapp/deployments" # Blue/green deployment root
current_symlink: "/opt/myapp/current" # Active deployment pointer
run_command: "poetry run python main.py" # App startup command
```

### GitHub API Integration
```go
// Recommended patterns for GitHub API usage
// Use personal access token from environment or config
// Implement exponential backoff: 1s, 2s, 4s, 8s, 16s max
// Cache release info to minimize API calls
// Handle 403 rate limit responses gracefully
type GitHubClient struct {
    token string
    lastCheck time.Time
    cachedRelease *Release
}
```

### Deployment Flow (to implement)
1. Poll GitHub API for releases matching `asset_suffix`
2. Download to inactive slot (`blue` or `green`)
3. Extract and run `poetry install` in slot directory
4. Test deployment via health check
5. Atomically switch `current_symlink` 
6. Update `state.yaml` with new active slot
7. Execute `post_deploy_script` for service management

## Project Conventions

### Go Implementation Patterns
- Single binary deployment model (no external dependencies)
- YAML-based configuration (using `gopkg.in/yaml.v3`)
- File-based state persistence (avoid databases for Raspberry Pi)
- Graceful shutdown with systemd integration

### Poetry Integration Specifics
- Each blue/green slot maintains independent Poetry venvs
- Use `poetry install --no-dev` for production deployments  
- Virtual env isolation prevents dependency conflicts during rollback

### Error Handling & Rollback
- Always validate download integrity before extraction
- Keep previous deployment slot as automatic rollback target
- Log all state changes to support debugging
- Never modify `current_symlink` until new deployment is verified

#### Rollback Procedures
```go
// Rollback strategy - switch back to previous slot
func (d *Deployer) Rollback() error {
    // 1. Update state.yaml to previous active_slot
    // 2. Atomically switch current_symlink
    // 3. Restart services via deploy.sh
    // 4. Validate rollback with health checks
    // 5. Log rollback event with reason
}
```

#### Health Check Integration
- Implement HTTP health endpoint checking before symlink switch
- Use configurable timeout (default 30s) for health validation
- Support custom health check commands in `config.yaml`
- Always test new deployment before making it active

## Critical Implementation Notes

### Raspberry Pi Constraints
- Optimize for ARM architecture and limited memory
- Use efficient polling intervals (`check_interval_seconds: 300`)
- Implement exponential backoff for GitHub API rate limiting
- Handle network interruptions gracefully

### systemd Integration

#### Service File Best Practices
```ini
# /etc/systemd/system/gh-deployer.service
[Unit]
Description=GitHub Release Deployer
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
ExecStart=/usr/local/bin/gh-deployer
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=10
User=deployer
Group=deployer
WorkingDirectory=/opt/myapp/gh-deployer

[Install]
WantedBy=multi-user.target
```

#### Graceful Shutdown Pattern
```go
// Handle SIGTERM/SIGINT for graceful shutdown
func (d *Deployer) handleSignals() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        d.shutdown = true
        // Complete current deployment before exiting
    }()
}
```

#### Deploy Script Integration
- Use `deploy.sh` for coordinating with systemd service lifecycle
- Support `systemctl reload` for config changes without deployment interruption
- Implement service dependency management (database, web server, etc.)

## Testing Strategy

### Unit Testing Patterns
```go
// Mock GitHub API responses for consistent testing
type MockGitHubClient struct { releases []Release }

// Test state transitions thoroughly
func TestBlueGreenSwitch(t *testing.T) {
    // Test all state combinations: blue->green, green->blue
    // Verify state.yaml persistence after each transition
}
```

### Integration Testing
- Use `httptest.Server` for GitHub API mocking
- Test with real Poetry projects in temporary directories
- Validate symlink atomicity on target filesystem type
- Mock systemd interactions for CI environments

### End-to-End Testing
```bash
# Create test repository with real releases
# Test complete deployment cycle
./gh-deployer --config test-config.yaml --dry-run
```

## Monitoring & Logging

### Structured Logging Pattern
```go
// Use structured logging for easier parsing
log.WithFields(log.Fields{
    "deployment_id": deploymentID,
    "slot": targetSlot,
    "version": newVersion,
    "duration_ms": deploymentTime,
}).Info("Deployment completed successfully")
```

### Key Metrics to Track
- Deployment frequency and success rate
- Time between GitHub release and deployment completion
- Rollback frequency and triggers
- API rate limit consumption
- Disk space usage in blue/green slots

### Log Rotation & Retention
```yaml
# Configure in config.yaml
logging:
  level: "info"
  file: "/var/log/gh-deployer/deployer.log"
  max_size: "100MB"
  max_backups: 5
  max_age: 30 # days
```

### Health Monitoring Integration
- Expose `/health` endpoint for external monitoring
- Include deployment status, last successful deployment time
- Monitor Poetry venv health and dependency conflicts
- Alert on consecutive deployment failures

When implementing features, prioritize deployment safety over speed - failed deployments should never leave the system in an inconsistent state.