# Contributing to Cloudness

Thank you for your interest in contributing to Cloudness! We welcome contributions from the community and are grateful for your support.

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the [issue tracker](https://github.com/cloudness-io/cloudness/issues) to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** (code snippets, screenshots, etc.)
- **Describe the behavior you observed and what you expected**
- **Include your environment details** (OS, Go version, Kubernetes version, etc.)

**Bug Report Template:**

```markdown
**Description:**
A clear description of the bug.

**Steps to Reproduce:**
1. Go to '...'
2. Click on '...'
3. See error

**Expected Behavior:**
What you expected to happen.

**Actual Behavior:**
What actually happened.

**Environment:**
- OS: [e.g., Ubuntu 22.04]
- Go Version: [e.g., 1.21.5]
- Cloudness Version: [e.g., v1.0.0]
- Kubernetes Version: [e.g., 1.28.0]

**Additional Context:**
Add any other context, logs, or screenshots.
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **List examples of how it would be used**

### Pull Requests

We actively welcome your pull requests:

1. **Fork the repository** and create your branch from `main`
2. **Follow the development setup** outlined below
3. **Make your changes** following our coding standards
4. **Add tests** if applicable
5. **Ensure all tests pass**
6. **Update documentation** if needed
7. **Submit a pull request**

**Pull Request Guidelines:**

- Keep PRs focused on a single feature or bug fix
- Write clear, descriptive commit messages
- Reference related issues in your PR description
- Ensure CI checks pass
- Request review from maintainers
- Be responsive to feedback

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Node.js (latest stable version)
- Docker (for local testing)
- Kubernetes cluster (for integration testing)

### Setup Steps

- **Clone the repository:**

```bash
git clone https://github.com/cloudness-io/cloudness.git
cd cloudness
```

- **Install dependencies:**

```bash
make dep
make tools
```

- **Run the development server:**

```bash
make dev
```

- **Access the UI:**

Open your browser at `http://localhost:7331`

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `make format` to format your code
- Run `go vet` to catch common mistakes
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused

**Example:**

```go
// ProcessDeployment handles the deployment of an application to Kubernetes.
// It validates the deployment spec, creates necessary resources, and monitors
// the deployment status.
func ProcessDeployment(ctx context.Context, spec *DeploymentSpec) error {
    if err := spec.Validate(); err != nil {
        return fmt.Errorf("invalid deployment spec: %w", err)
    }
    
    // Implementation...
    return nil
}
```

### Frontend Code

- Follow [Templ](https://templ.guide/) best practices for Go templates
- Use `make format` to format template files
- Write semantic HTML in templ components
- Use Tailwind CSS utility classes consistently
- Keep components small and reusable
- Use Alpine.js for interactive elements when needed

**Example:**

```templ
package components

import "github.com/cloudness-io/cloudness/types"

// DeploymentCard renders a deployment card component
templ DeploymentCard(deployment *types.Deployment) {
    <div class="bg-white rounded-lg shadow p-4">
        <h3 class="text-lg font-semibold">{ deployment.Name }</h3>
        <p class="text-gray-600">{ deployment.Status }</p>
    </div>
}
```

### Commit Messages

Use clear and meaningful commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**

```
feat(api): add endpoint for deployment rollback

Add new API endpoint to support rolling back deployments to previous versions.
Includes validation and error handling.

Closes #123
```

```
fix(auth): resolve token expiration issue

Fixed bug where JWT tokens were not properly validated after expiration.
Updated token validation logic to check expiry timestamp.

Fixes #456
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./app/auth/...

# Run with coverage
make test-coverage
```

### Writing Tests

- Write unit tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

**Example:**

```go
func TestProcessDeployment(t *testing.T) {
    tests := []struct {
        name    string
        spec    *DeploymentSpec
        wantErr bool
    }{
        {
            name: "valid deployment",
            spec: &DeploymentSpec{
                Name:  "test-app",
                Image: "nginx:latest",
            },
            wantErr: false,
        },
        {
            name:    "invalid deployment - missing name",
            spec:    &DeploymentSpec{Image: "nginx:latest"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ProcessDeployment(context.Background(), tt.spec)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessDeployment() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

- Update README.md if you change functionality
- Add comments to exported functions and templ components
- Update API documentation if you add/modify endpoints
- Include examples in documentation
- Regenerate templ files after changes: `templ generate`

## Review Process

1. **Automated Checks:** All PRs must pass CI checks (build, tests, linting)
2. **Code Review:** At least one maintainer must approve the PR
3. **Testing:** Ensure changes are tested (unit tests, integration tests)
4. **Documentation:** Verify documentation is updated if needed
5. **Merge:** Once approved and passing all checks, a maintainer will merge

## Community

- Join our [GitHub Discussions](https://github.com/cloudness-io/cloudness/discussions)
- Ask questions, share ideas, and help others
- Be respectful and constructive in all interactions

## License

By contributing to Cloudness, you agree that your contributions will be licensed under the Apache License 2.0.

## Questions?

If you have questions about contributing, feel free to:

- Open a [discussion](https://github.com/cloudness-io/cloudness/discussions)
- Reach out to the maintainers
- Check existing issues and PRs

Thank you for contributing to Cloudness! 🎉
