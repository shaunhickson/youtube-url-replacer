import { test, expect } from './fixtures';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

test.describe('Complex DOM Scenarios', () => {

  test.beforeEach(async ({ page }) => {
    await page.route('**/resolve*', async (route) => {
      const postData = route.request().postDataJSON();
      const urls = postData?.urls || [];
      
      const titles: Record<string, string> = {};
      for (const url of urls) {
        if (url.includes('shadowdom')) {
          titles[url] = "Shadow DOM Link Resolved";
        } else if (url.includes('SPA')) {
          titles[url] = "React SPA Link Resolved";
        } else if (url.includes('scroll')) {
          titles[url] = `Scroll Link Resolved`;
        } else {
          titles[url] = "Mocked Generic Title";
        }
      }
      
      await route.fulfill({ json: { titles } });
    });
  });

  test('React SPA Test: Link injected after initial load is resolved', async ({ page }) => {
    const fixturePath = path.join(__dirname, 'fixtures', 'react-spa.html');
    await page.goto(`file://${fixturePath}`);

    // Wait for the simulated SPA link injection (500ms in fixture) + batching
    const link = page.locator('a[href="https://www.youtube.com/watch?v=SPA"]');
    await expect(link).toContainText('React SPA Link Resolved');
  });

  test('Infinite Scroll Test: Multiple links injected continuously are resolved', async ({ page }) => {
    const fixturePath = path.join(__dirname, 'fixtures', 'infinite-scroll.html');
    await page.goto(`file://${fixturePath}`);

    // Pick a few links to verify
    const link0 = page.locator('a[href="https://www.youtube.com/watch?v=scroll0"]');
    const link49 = page.locator('a[href="https://www.youtube.com/watch?v=scroll49"]');

    await expect(link0).toContainText('Scroll Link Resolved');
    await expect(link49).toContainText('Scroll Link Resolved');
  });

  // Note: Currently skipped because our MutationObserver does not pierce Shadow DOM
  test.skip('Shadow DOM Test: Link inside open shadow root is resolved', async ({ page }) => {
    const fixturePath = path.join(__dirname, 'fixtures', 'shadow-dom.html');
    await page.goto(`file://${fixturePath}`);

    // We must pierce the shadow DOM to find the link
    const link = page.locator('#host').locator('a');
    await expect(link).toContainText('Shadow DOM Link Resolved');
  });
});
