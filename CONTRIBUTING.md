# Contributing to Sentinel Mesh

Thank you for your interest in contributing to Sentinel Mesh! This is a personal hobby project, but contributions are welcome.

## Fork-Based Contribution Workflow

This repository requires a **fork-based workflow** for all contributions. Direct branch creation is restricted to maintain a clean repository structure.

### Getting Started

1. **Fork the Repository**
   - Click the "Fork" button at the top right of the repository page
   - This creates a copy of the repository under your GitHub account

2. **Clone Your Fork**
   ```bash docs-drift:skip
   git clone https://github.com/YOUR-USERNAME/sentinel-mesh.git
   cd sentinel-mesh
   ```

3. **Add Upstream Remote**
   ```bash docs-drift:skip
   git remote add upstream https://github.com/georg-nikola/sentinel-mesh.git
   ```

4. **Create a Feature Branch**
   ```bash docs-drift:skip
   git checkout -b feature/your-feature-name
   ```

### Making Changes

1. **Keep Your Fork Updated**
   ```bash docs-drift:skip
   git fetch upstream
   git rebase upstream/main
   ```

2. **Make Your Changes**
   - Write clear, concise commit messages
   - Follow existing code style and conventions
   - Add tests for new features
   - Update documentation as needed

3. **Run Tests**
   ```bash docs-drift:skip
   # Run E2E tests
   ./tests/e2e/run-all.sh

   # Run Go tests
   cd services/api && go test ./...
   cd services/collector && go test ./...

   # Run Python tests
   cd ml-service && python -m pytest

   # Run frontend tests
   cd web && npm test
   ```

4. **Commit Your Changes**
   ```bash docs-drift:skip
   git add .
   git commit -m "Add feature: description of your changes"
   ```

### Submitting a Pull Request

1. **Push to Your Fork**
   ```bash docs-drift:skip
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request**
   - Go to your fork on GitHub
   - Click "Pull Request" button
   - Select `main` branch of `georg-nikola/sentinel-mesh` as the base
   - Provide a clear title and description
   - Reference any related issues

3. **PR Requirements**
   - All CI checks must pass (tests, linting, security scans)
   - At least 1 approving review required
   - All conversations must be resolved
   - Linear history required (use rebase, not merge commits)

### Code Review Process

- PRs are reviewed as time permits (hobby project)
- Be open to feedback and suggestions
- Address review comments promptly
- Keep PRs focused and reasonably sized

### Development Guidelines

#### Go Services (API, Collector)
- Follow Go best practices and idioms
- Use `gofmt` and `golint`
- Write table-driven tests
- Document exported functions and types

#### Python Services (ML Service)
- Follow PEP 8 style guide
- Use type hints where applicable
- Write docstrings for functions and classes
- Use `black` for formatting and `isort` for imports

#### Frontend (Vue.js)
- Follow Vue.js style guide
- Use TypeScript for type safety
- Write component tests with Vitest
- Ensure accessibility standards

#### Kubernetes Manifests
- Use proper resource limits and requests
- Include health and readiness probes
- Document resource requirements
- Follow security best practices

### Testing

All changes must include appropriate tests:
- Unit tests for business logic
- Integration tests for API endpoints
- E2E tests for critical user flows
- Load tests for performance-sensitive code

### Documentation

Update relevant documentation when:
- Adding new features
- Changing behavior
- Modifying APIs
- Updating configuration

### Branch Protection

The `main` branch is protected with:
- Required pull request reviews (1 approval)
- Required status checks (all CI must pass)
- No force pushes allowed
- No deletions allowed
- Conversation resolution required
- Linear history enforced

### Questions or Issues?

- Open an issue for bugs or feature requests
- Use Discussions for general questions
- Be respectful and constructive

## License

By contributing to Sentinel Mesh, you agree that your contributions will be licensed under the same license as the project (see LICENSE file).

---

Thank you for contributing to Sentinel Mesh! 🎉
