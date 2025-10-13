# Creating Releases

This repository is set up for automated releases when you push tags.

## How to Create a Release

1. **Make sure all changes are committed and pushed to main**
2. **Create and push a tag:**
   ```bash
   # Create a new version tag (use semantic versioning)
   git tag v1.0.0
   
   # Push the tag to trigger the release
   git push origin v1.0.0
   ```

3. **The GitHub Actions workflow will automatically:**
   - Run all tests and linting
   - Build binaries for multiple platforms:
     - Linux (AMD64, ARM64, ARMv7 for Raspberry Pi)
     - macOS (Intel and Apple Silicon)
     - Windows (AMD64)
   - Create checksums for all binaries
   - Create a GitHub release with all assets
   - Generate release notes from commits

## Go Package Publishing

### Automatic Module Publishing
- **Go modules are automatically available** when you push tags
- Your package will be available at: `https://pkg.go.dev/github.com/kpeacocke/deployer`
- Users can import it with: `go get github.com/kpeacocke/deployer@v1.0.0`

### Using as a Library
Other Go projects can use your code as a library:

```go
import "github.com/kpeacocke/deployer"

// Use the deployer components
client := NewGitHubClient("your-token")
config, err := LoadConfig("config.yaml")
```

### Binary Distribution
Users can download pre-built binaries from the GitHub releases page, or install via Go:

```bash
# Install the latest version
go install github.com/kpeacocke/deployer@latest

# Install a specific version
go install github.com/kpeacocke/deployer@v1.0.0
```

## Version Numbering

Use [Semantic Versioning](https://semver.org/):
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features, backward compatible)
- `v1.0.1` - Patch release (bug fixes)

## Testing a Release Locally

Before creating a real release, you can test the build process:

```bash
# Test multi-platform builds
make build-all

# Test the version flag works
go run . --version

# Test help output
go run . --help
```