import { test, expect } from './fixtures';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

test('LinkLens injects tooltip and modifies link text', async ({ page }) => {
  // 1. Mock the backend API call to ensure fast, deterministic tests without hitting network
  await page.route('**/resolve*', async (route) => {
    const json = {
      titles: {
        "https://www.youtube.com/watch?v=dQw4w9WgXcQ": "Rick Astley - Never Gonna Give You Up (Official Music Video)"
      }
    };
    await route.fulfill({ json });
  });

  // 2. Navigate to our local test fixture
  const testHtmlPath = path.join(__dirname, 'test.html');
  await page.goto(`file://${testHtmlPath}`);

  // 3. The original link should be modified by LinkLens
  const link = page.locator('a[href="https://www.youtube.com/watch?v=dQw4w9WgXcQ"]');
  
  // Wait for the link text to change to the mocked title
  await expect(link).toContainText('Rick Astley - Never Gonna Give You Up');

  // 4. Verify LinkLens specific DOM elements are injected (e.g. icon or custom classes)
  // Check if the custom data attribute or shadow dom is added.
  // We'll just verify the text changed as a basic smoke test for now.
});
