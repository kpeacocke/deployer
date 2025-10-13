# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0.0 | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a vulnerability, please follow these steps:

### Reporting Process

1. **Do NOT create a public GitHub issue** for security vulnerabilities
2. **Email the maintainers** at [security contact] with details
3. **Include the following information:**
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact assessment
   - Any suggested fixes or mitigations

### Response Timeline

- **Initial Response**: Within 48 hours of report
- **Assessment**: Within 7 days of report
- **Fix Development**: Timeline depends on severity and complexity
- **Public Disclosure**: After fix is available and deployed

### Security Considerations for gh-deployer

Since gh-deployer manages application deployments, security is critical:

- **GitHub Token Security**: Store tokens securely, use environment variables
- **File System Access**: Deployer requires write access to deployment directories
- **Network Access**: Deployer makes HTTPS requests to GitHub API
- **Process Execution**: Deployer executes post-deployment scripts
- **Symlink Management**: Atomic symlink operations prevent race conditions

### Security Best Practices

When deploying gh-deployer:

1. **Run with minimal privileges** - Use dedicated deployer user account
2. **Secure configuration files** - Protect config.yaml with appropriate permissions
3. **Monitor logs** - Watch for unusual deployment patterns or failures  
4. **Validate releases** - Verify GitHub release authenticity
5. **Network security** - Use HTTPS for all API communication
6. **Regular updates** - Keep gh-deployer updated with latest security patches

### Known Security Considerations

- **Post-deploy scripts** run with deployer user privileges
- **Deployment directories** must be writable by deployer user
- **GitHub API access** requires valid authentication token
- **Symlink race conditions** are prevented by atomic operations

## Acknowledgments

We appreciate responsible disclosure of security vulnerabilities and will acknowledge security researchers who help improve the security of gh-deployer.
