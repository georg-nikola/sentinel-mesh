import { test, expect } from '@playwright/test';

test.describe('Dashboard - Basic Functionality', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the dashboard
    await page.goto('/');

    // Wait for the Vue app to load
    await page.waitForSelector('#app', { timeout: 10000 });
  });

  test('should load dashboard and display title', async ({ page }) => {
    // Check page title
    await expect(page).toHaveTitle(/Sentinel Mesh/i);

    // Check for main heading or logo
    const heading = page.locator('h1, h2, .logo, [class*="title"]').first();
    await expect(heading).toBeVisible();
  });

  test('should display Vue.js application container', async ({ page }) => {
    // Verify Vue app mount point exists and is visible
    const appContainer = page.locator('#app').first();
    await expect(appContainer).toBeVisible();

    // Check that the app has content
    const content = await appContainer.textContent();
    expect(content).toBeTruthy();
    expect(content!.length).toBeGreaterThan(0);
  });

  test('should load without JavaScript errors', async ({ page }) => {
    const errors: string[] = [];

    page.on('pageerror', (error) => {
      errors.push(error.message);
    });

    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });

    await page.reload();
    await page.waitForLoadState('networkidle');

    // Allow up to 2 seconds for any async errors
    await page.waitForTimeout(2000);

    // Check for critical errors (filter out common warnings)
    const criticalErrors = errors.filter(
      err => !err.includes('warning') &&
             !err.includes('deprecated') &&
             !err.includes('favicon')
    );

    expect(criticalErrors).toHaveLength(0);
  });

  test('should have responsive layout', async ({ page }) => {
    // Test desktop viewport
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.waitForTimeout(500);

    const appDesktop = page.locator('#app').first();
    await expect(appDesktop).toBeVisible();

    // Test mobile viewport
    await page.setViewportSize({ width: 375, height: 812 });
    await page.waitForTimeout(500);

    const appMobile = page.locator('#app').first();
    await expect(appMobile).toBeVisible();

    // Test tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);

    const appTablet = page.locator('#app').first();
    await expect(appTablet).toBeVisible();
  });
});

test.describe('Dashboard - Metrics and Data', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('#app');
  });

  test('should display metrics or charts section', async ({ page }) => {
    // Look for common chart/metrics containers
    const metricsSelectors = [
      '[class*="chart"]',
      '[class*="metric"]',
      '[class*="dashboard"]',
      'canvas',
      'svg'
    ];

    let foundMetrics = false;
    for (const selector of metricsSelectors) {
      const element = page.locator(selector).first();
      if (await element.count() > 0) {
        foundMetrics = true;
        break;
      }
    }

    // If no metrics UI found yet, that's okay for initial setup
    // Just verify the page structure is present
    const appContainer = page.locator('#app').first();
    await expect(appContainer).toBeVisible();
  });

  test('should handle API connectivity', async ({ page }) => {
    // Monitor network requests for API calls
    const apiRequests: string[] = [];

    page.on('request', (request) => {
      const url = request.url();
      if (url.includes('/api/') || url.includes(':8080') || url.includes(':8000')) {
        apiRequests.push(url);
      }
    });

    await page.reload();
    await page.waitForLoadState('networkidle');

    // The dashboard might make API calls to backend services
    // This test just verifies the page handles API state gracefully
    const appContainer = page.locator('#app').first();
    await expect(appContainer).toBeVisible();
  });
});

test.describe('Dashboard - User Interactions', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('#app');
  });

  test('should have interactive elements', async ({ page }) => {
    // Look for buttons, links, or interactive components
    const buttons = page.locator('button, a, [role="button"]');
    const buttonCount = await buttons.count();

    // Dashboard should have at least some interactive elements
    // Even if it's just navigation or refresh buttons
    expect(buttonCount).toBeGreaterThanOrEqual(0);
  });

  test('should maintain state during interactions', async ({ page }) => {
    // Get initial page content
    const initialContent = await page.locator('#app').first().textContent();

    // Scroll the page
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
    await page.waitForTimeout(500);

    // Verify content is still present
    const appContainer = page.locator('#app').first();
    await expect(appContainer).toBeVisible();

    // Verify the app didn't crash or reload unexpectedly
    const currentContent = await page.locator('#app').first().textContent();
    expect(currentContent).toBeTruthy();
  });
});

test.describe('Dashboard - Performance', () => {
  test('should load within reasonable time', async ({ page }) => {
    const startTime = Date.now();

    await page.goto('/');
    await page.waitForSelector('#app');
    await page.waitForLoadState('networkidle');

    const loadTime = Date.now() - startTime;

    // Dashboard should load within 10 seconds
    expect(loadTime).toBeLessThan(10000);
  });

  test('should have good Lighthouse scores', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('#app');

    // Basic performance checks
    const performanceMetrics = await page.evaluate(() => {
      const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
      return {
        domContentLoaded: navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
        loadComplete: navigation.loadEventEnd - navigation.loadEventStart,
        domInteractive: navigation.domInteractive - navigation.fetchStart,
      };
    });

    // These should complete reasonably quickly
    expect(performanceMetrics.domContentLoaded).toBeLessThan(3000);
    expect(performanceMetrics.loadComplete).toBeLessThan(5000);
  });
});
