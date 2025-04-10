# Contributing to Agent Runtime

Thank you for your interest in contributing to Agent Runtime! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please read it before contributing.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue on GitHub with the following information:

- A clear, descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Any relevant logs or screenshots
- Your environment (OS, Go version, etc.)

### Suggesting Enhancements

We welcome suggestions for enhancements! Please create an issue with:

- A clear, descriptive title
- A detailed description of the proposed enhancement
- Any relevant examples or mockups
- Why this enhancement would be useful

### Pull Requests

1. Fork the repository
2. Create a new branch from `main`
3. Make your changes
4. Run tests and ensure they pass
5. Submit a pull request

#### Pull Request Guidelines

- Follow the Go style guide and code conventions
- Include tests for new features or bug fixes
- Update documentation as needed
- Keep pull requests focused on a single change
- Link to any relevant issues

## Development Setup

### Prerequisites

- Go 1.24 or higher
- Docker (for running tests)
- Git

### Setting Up the Development Environment

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/agent_runtime.git
cd agent_runtime

# Add the upstream repository
git remote add upstream https://github.com/spectrumwebco/agent_runtime.git

# Install dependencies
go mod download

# Run tests
go test ./...
```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Run `golint` and `go vet` before submitting changes

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests after the first line

### Documentation

- Update the README.md if necessary
- Add comments to your code
- Update any relevant documentation in the `docs` directory

## Testing

- Write tests for new features and bug fixes
- Ensure all tests pass before submitting a pull request
- Aim for high test coverage

## Review Process

- All submissions require review
- Changes may be requested before a pull request is merged
- Be responsive to feedback

## Community

- Join our community discussions
- Help answer questions from other contributors
- Share your experiences using Agent Runtime

Thank you for contributing to Agent Runtime!
