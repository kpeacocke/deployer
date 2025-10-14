# Project Structure

This document describes the organization and structure of the gh-deployer project.

## Directory Structure

```text
gh-deployer/
├── .github/                    # GitHub-specific files
│   ├── ISSUE_TEMPLATE/         # Issue templates
│   │   ├── bug_report.yml      # Bug report template
│   │   ├── feature_request.yml # Feature request template
│   │   ├── question.yml        # Support question template
│   │   └── config.yml          # Issue template configuration
│   ├── workflows/              # GitHub Actions workflows
│   │   ├── ci.yml              # Continuous integration
│   │   └── release.yml         # Release automation
│   ├── copilot-instructions.md # AI agent instructions
│   └── pull_request_template.md # PR template
├── .vscode/                    # VS Code configuration
│   ├── extensions.json         # Recommended extensions
│   ├── launch.json             # Debug configurations
│   ├── settings.json           # Workspace settings
│   ├── tasks.json              # Build tasks
│   └── README.md               # Extension management guide
├── docs/                       # Documentation (future)
├── examples/                   # Usage examples (future)
├── scripts/                    # Utility scripts (future)
├── gh-deployer.code-workspace  # VS Code workspace file
├── config.go                   # Configuration management
├── config_test.go              # Configuration tests
├── config.example.yaml         # Example configuration
├── deployer.go                 # Main deployment logic
├── github.go                   # GitHub API client
├── github_test.go              # GitHub client tests
├── integration_test.go         # Integration tests
├── main.go                     # Application entry point
├── state.go                    # State management
├── state_test.go               # State management tests
├── deploy.sh                   # Post-deployment script
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies checksums
├── .gitignore                  # Git ignore rules
├── .golangci.yml               # Linter configuration
├── .editorconfig               # Editor configuration
├── CHANGELOG.md                # Change log
├── CODE_OF_CONDUCT.md          # Code of conduct
├── CONTRIBUTING.md             # Contributing guidelines
├── LICENSE                     # MIT license
├── PROJECT_STRUCTURE.md        # This file
├── README.md                   # Project documentation
└── SECURITY.md                 # Security policy
```

## File Descriptions

### Core Application Files

- **`main.go`** - Application entry point with CLI parsing and graceful shutdown
- **`config.go`** - Configuration loading, parsing, and validation
- **`state.go`** - Deployment state management and persistence
- **`deployer.go`** - Core deployment logic and orchestration
- **`github.go`** - GitHub API client with authentication and rate limiting

### Test Files

- **`*_test.go`** - Unit tests for corresponding modules
- **`integration_test.go`** - End-to-end integration tests

### Configuration Files

- **`config.example.yaml`** - Example configuration with all options
- **`deploy.sh`** - Post-deployment hook script template
- **`.golangci.yml`** - Go linter configuration
- **`.editorconfig`** - Cross-editor configuration

### Build and Development

- **`Makefile`** - Build, test, and development automation
- **`go.mod`** - Go module dependencies
- **`.vscode/`** - VS Code workspace configuration and extension management guide
- **`gh-deployer.code-workspace`** - Complete workspace settings with extension control

### Documentation

- **`README.md`** - Main project documentation
- **`CONTRIBUTING.md`** - Development and contribution guidelines
- **`SECURITY.md`** - Security policy and vulnerability reporting
- **`CHANGELOG.md`** - Version history and changes
- **`CODE_OF_CONDUCT.md`** - Community guidelines

### GitHub Integration

- **`.github/workflows/`** - CI/CD automation
- **`.github/ISSUE_TEMPLATE/`** - Standardized issue reporting
- **`.github/copilot-instructions.md`** - AI agent guidance

## Development Workflow

1. **Setup**: Clone repo, run `make deps`
2. **Development**: Edit code, run `make test`
3. **Testing**: Run `make check` before committing
4. **Building**: Run `make build` to create binary
5. **Release**: Tag version, GitHub Actions handles the rest

## Dependencies

### Runtime Dependencies

- **Go 1.21+** - Programming language
- **gopkg.in/yaml.v3** - YAML parsing

### Development Dependencies

- **golangci-lint** - Code linting
- **make** - Build automation
- **git** - Version control

### Optional Dependencies

- **systemd** - Service management (Linux)

## Architecture Patterns

- **Single Binary** - No external runtime dependencies
- **Configuration-Driven** - YAML-based configuration
- **State Management** - File-based state persistence
- **Blue/Green Deployment** - Zero-downtime deployments
- **Health Checks** - Pre-activation validation
- **Graceful Shutdown** - Signal-based cleanup

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and workflow.
