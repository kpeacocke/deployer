# Contributing to gh-deployer

Thank you for your interest in contributing to gh-deployer! This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct (see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)).

## How to Contribute

### Reporting Issues

1. **Search existing issues** first to avoid duplicates
2. **Use the issue templates** when available
3. **Provide detailed information** including:
   - Operating system and version
   - Go version
   - Steps to reproduce the issue
   - Expected vs actual behavior
   - Relevant logs or error messages

### Submitting Pull Requests

1. **Fork the repository** and create a feature branch
2. **Write clear commit messages** following conventional commits format
3. **Add tests** for new functionality
4. **Ensure all tests pass** with `make test`
5. **Follow the coding standards** (run `make fmt` and `make vet`)
6. **Update documentation** as needed
7. **Submit a pull request** with a clear description

### Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/kpeacocke/deployer.git
   cd deployer
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Run tests:**
   ```bash
   make test
   ```

4. **Build the application:**
   ```bash
   make build
   ```

### Coding Standards

- Follow standard Go formatting (`gofmt`)
- Write comprehensive tests for new features
- Use meaningful variable and function names
- Add comments for complex logic
- Follow the existing code style and patterns

### Testing

- Write unit tests for all new functions
- Include integration tests for complex workflows
- Ensure test coverage remains high
- Test edge cases and error conditions

### Documentation

- Update README.md for user-facing changes
- Add or update code comments for complex logic
- Update configuration examples if needed
- Consider updating the copilot instructions for architectural changes

## Release Process

Releases are managed by project maintainers:

1. Version bumping follows semantic versioning
2. Releases are tagged and published via GitHub Actions
3. Release notes are generated automatically from commit messages

## Getting Help

- **Check the documentation** in README.md and `.github/copilot-instructions.md`
- **Search existing issues** for similar questions
- **Create a new issue** if you need help

## Recognition

Contributors are recognized in our releases and documentation. Thank you for helping make gh-deployer better!
