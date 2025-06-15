# Contributing to n8n-client-go

Thank you for your interest in contributing to n8n-client-go! This document will guide you through the contribution process.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Submitting Pull Requests](#submitting-pull-requests)
- [Style Guide](#style-guide)
  - [Commit Messages](#commit-messages)
  - [Code Style](#code-style)
- [Review Process](#review-process)
- [License](#license)

## Code of Conduct

This project and all participants are governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

1. **Check existing issues**: Before creating a new issue, please check if a similar report already exists.
2. **Create a new issue**: If you don't find an existing issue, create a new one with a clear and detailed description of the problem.
3. **Include relevant information**:
   - Go version
   - n8n version
   - Steps to reproduce the issue
   - Error messages
   - Screenshots (if applicable)

### Suggesting Enhancements

1. **Check existing suggestions**: Look for similar enhancement requests.
2. **Create a new issue**: Describe the proposed enhancement and explain why it would be useful.
3. **Be specific**: Include use cases and test cases if possible.

### Submitting Pull Requests

1. **Create a descriptive branch** for your feature or fix:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```
3. **Commit** your changes with descriptive messages.
4. **Push** the branch to your fork:
   ```bash
   git push origin your-branch-name
   ```
5. **Open a Pull Request** against the `main` branch.

## Style Guide

### Commit Messages

Use the following format for commit messages:

```
type(scope): brief description

Detailed description if needed

Fixes #123
```

**Commit Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, missing semi-colons, etc.
- `refactor`: Code changes that neither fix a bug nor add a feature
- `test`: Adding or fixing tests
- `chore`: Changes to the build process or tools

### Code Style

1. **Go fmt**: Ensure your code is formatted with `go fmt`.
2. **Comments**: Document exported functions and types following Go conventions.
3. **Descriptive names**: Use descriptive names for variables and functions.
4. **Keep functions short**: Functions should do one thing and do it well.
5. **Error handling**: Always check and handle errors appropriately.

## Review Process

1. Pull Requests will be reviewed by project maintainers.
2. Changes may be requested before the PR is merged.
3. Ensure all tests pass before submitting a PR.
4. Update documentation if necessary.

## License

By contributing to this project, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

Thank you for contributing to n8n-client-go! Your help is greatly appreciated. 🚀
