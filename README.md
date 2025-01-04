# gh-deployer

A Go-based GitHub release deployer with blue/green deployment, designed to run on Raspberry Pi and launch Python apps using Poetry.

## Features

- Polls GitHub for latest release
- Uses separate Poetry venvs for blue and green slots
- Atomic symlink switching
- Optional post-deploy hook
- Startup-safe with systemd

See `config.yaml` for configuration.
