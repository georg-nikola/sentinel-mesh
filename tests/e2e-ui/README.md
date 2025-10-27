# Sentinel Mesh UI E2E Tests

Playwright-based end-to-end tests for the Sentinel Mesh dashboard frontend.

## Overview

These tests verify the Vue.js dashboard functionality, user interactions, and visual behavior. They complement the shell-based infrastructure tests in `tests/e2e/`.

**Test Coverage**:
- Dashboard loading and basic functionality
- Metrics and data display
- User interactions and navigation
- Performance and responsiveness
- Error handling

## Prerequisites

- Node.js 18+ and npm
- Kubernetes cluster running with Sentinel Mesh deployed
- Frontend accessible at http://localhost:30000 (or custom BASE_URL)

## Installation

```bash
# Install dependencies
npm install

# Install Playwright browsers
npx playwright install chromium
```

## Running Tests

### All Tests

```bash
# Run all tests (headless)
npm test

# Run tests in headed mode (see browser)
npm run test:headed

# Run tests with Playwright UI (interactive)
npm run test:ui
```

### Specific Tests

```bash
# Run only Chromium tests
npm run test:chromium

# Run specific test file
npx playwright test tests/dashboard.spec.ts

# Run tests matching a title
npx playwright test -g "should load dashboard"
```

### Debugging

```bash
# Debug mode (step through tests)
npm run test:debug

# Generate test code interactively
npm run test:codegen
```

### View Reports

```bash
# Show HTML test report
npm run test:report

# View test results JSON
cat test-results.json
```

## Test Structure

```
tests/e2e-ui/
├── playwright.config.ts    # Playwright configuration
├── package.json             # Dependencies and scripts
├── tests/
│   └── dashboard.spec.ts    # Dashboard UI tests
├── playwright-report/       # HTML reports (generated)
└── test-results/            # Test artifacts (generated)
```

## Test Files

### `tests/dashboard.spec.ts`

Tests for the main dashboard page:
- **Basic Functionality**: Page load, title, Vue app mount
- **Metrics and Data**: Charts, API connectivity
- **User Interactions**: Buttons, links, state management
- **Performance**: Load times, Lighthouse metrics

## Configuration

### Environment Variables

- `BASE_URL`: Dashboard URL (default: `http://localhost:30000`)
- `CI`: Enable CI mode (retries, single worker)

### Custom Configuration

Edit `playwright.config.ts` to customize:
- Test timeout
- Retries
- Parallelization
- Browsers to test
- Screenshot/video settings

## Writing New Tests

### Basic Test Template

```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('#app');
  });

  test('should do something', async ({ page }) => {
    // Arrange
    const element = page.locator('.some-element');

    // Act
    await element.click();

    // Assert
    await expect(element).toHaveText('Expected Text');
  });
});
```

### Best Practices

1. **Use Data Test IDs**: Add `data-testid` attributes to Vue components
2. **Wait for State**: Use `waitForSelector` or `waitForLoadState`
3. **Descriptive Names**: Test names should describe behavior clearly
4. **Independent Tests**: Each test should be self-contained
5. **Clean Up**: Use `beforeEach` and `afterEach` hooks

## Integration with CI/CD

### GitHub Actions

```yaml
- name: Run Playwright UI Tests
  run: |
    cd tests/e2e-ui
    npm ci
    npx playwright install chromium --with-deps
    npm test
  env:
    BASE_URL: http://localhost:30000

- name: Upload test report
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: playwright-report
    path: tests/e2e-ui/playwright-report/
```

## Troubleshooting

### Tests Fail to Connect

```bash
# Verify frontend is accessible
curl http://localhost:30000

# Check Kubernetes service
kubectl get svc web-dashboard -n sentinel-mesh

# Verify NodePort
kubectl get svc web-dashboard -n sentinel-mesh -o jsonpath='{.spec.ports[0].nodePort}'
```

### Browser Not Found

```bash
# Reinstall browsers
npx playwright install chromium --with-deps
```

### Test Timeouts

- Increase timeout in `playwright.config.ts`
- Check if services are running slowly
- Use `page.waitForLoadState('networkidle')` for async content

### Failed Screenshots

- Screenshots saved to `test-results/`
- Videos saved to `test-results/` (on failure only)
- View with: `ls -la test-results/`

## Resources

- [Playwright Documentation](https://playwright.dev)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Vue.js Testing Guide](https://vuejs.org/guide/scaling-up/testing.html)
- [Sentinel Mesh E2E Tests](../e2e/README.md) - Infrastructure tests

## Contributing

When adding new UI tests:
1. Follow the existing test structure
2. Add descriptive test names
3. Include comments for complex interactions
4. Update this README if adding new test files
5. Run tests locally before committing

---

**Note**: These tests run against a live Kubernetes deployment. Ensure services are running before executing tests.
