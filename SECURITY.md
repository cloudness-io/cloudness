# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| v*.*.*:  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please report it responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisories** (Preferred)
   - Go to the [Security tab](https://github.com/cloudness-io/cloudness/security/advisories) of this repository
   - Click "Report a vulnerability"
   - Fill out the form with details about the vulnerability

2. **Private Contact**
   - Contact the maintainers directly through GitHub

### What to Include

Please include the following information in your report:

- **Description** of the vulnerability
- **Steps to reproduce** the issue
- **Affected versions**
- **Potential impact** of the vulnerability
- **Suggested fix** (if you have one)

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 7 days
- **Resolution Timeline**: We aim to resolve critical vulnerabilities within 30 days
- **Credit**: We will credit you in the security advisory (unless you prefer to remain anonymous)

### Disclosure Policy

- We follow a coordinated disclosure process
- We request that you give us reasonable time to address the vulnerability before public disclosure
- We will work with you to understand and resolve the issue quickly
- Once a fix is available, we will publish a security advisory

## Security Best Practices

When deploying Cloudness, we recommend:

- **Keep updated**: Always run the latest version
- **Secure your Kubernetes cluster**: Follow [Kubernetes security best practices](https://kubernetes.io/docs/concepts/security/)
- **Use strong credentials**: Use strong, unique passwords for database and Redis connections
- **Enable TLS**: Use HTTPS for all external communications
- **Restrict network access**: Limit access to the Cloudness server and database
- **Regular backups**: Maintain regular backups of your database
- **Audit logs**: Monitor and review application logs regularly

## Security Features

Cloudness includes the following security features:

- **Authentication & Authorization**: Multi-tenant access control
- **Secrets Management**: Secure handling of sensitive configuration
- **Input Validation**: Protection against common injection attacks

## Past Security Advisories

No security advisories have been published yet.

---

Thank you for helping keep Cloudness and our users safe!
