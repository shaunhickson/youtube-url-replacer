# Design Document: Test Fixtures for E2E Testing (Issue #82)

## 1. Product Alignment
*   **Roadmap Check:** This design directly addresses Phase 4, Issue #82: "Create test fixtures for complex DOM structures like Shadow DOM, React SPAs, and infinite scroll" in `docs/PRODUCT_ROADMAP.md`.
*   **Scope:** The scope is limited to creating static and dynamic local HTML/JS files in the `extension/tests/e2e/fixtures/` directory, and adding Playwright tests to execute against them. We will *not* add complex backend logic here; all API calls will continue to be mocked.
*   **User Value:** Ensures that LinkLens works seamlessly across modern, complex websites without breaking layouts or failing to resolve links inside modern web frameworks and isolated DOM contexts.

## 2. Architecture & Patterns
*   **Modularity:** Fixtures will be separate HTML files, served via `file://` protocol during Playwright tests. This avoids the need for an external web server or relying on volatile public websites.
*   **Simplicity:** We will use vanilla JS (or simple React CDN links) within these HTML files to simulate the complex environments.
*   **Fixture Scenarios:**
    1.  **Shadow DOM:** A custom web component (`<shadow-host>`) containing a link inside its shadow root. LinkLens must pierce the shadow boundary to detect and modify it.
    2.  **React SPA (Client-Side Routing):** A simple React application that dynamically renders links *after* initial page load. LinkLens's `MutationObserver` must catch these.
    3.  **Infinite Scroll:** A script that continuously appends new links to the DOM as the user scrolls. LinkLens must efficiently process these without degrading performance.

## 3. Security & Safety
*   **Isolation:** Since tests are executed against local files with mocked APIs, there are no external security implications (e.g. no SSRF risk during test runs).

## 4. Scalability & Performance
*   **Test Latency:** Using local files ensures these complex DOM tests run in milliseconds.
*   **Mutation Observer Efficiency:** The "Infinite Scroll" fixture is specifically designed to stress-test our `MutationObserver`. We can assert that the observer processes new nodes efficiently without causing browser lockups.

## 5. Testability
*   **Unit Tests:** N/A (These are E2E fixtures).
*   **E2E Integration:** We will create `complex-dom.spec.ts` which will load each fixture page and assert:
    *   Shadow DOM: Link inside `#shadow-root` is successfully resolved and modified.
    *   React SPA: Link rendered 2 seconds after page load is successfully resolved and modified.
    *   Infinite Scroll: 100 links injected via scroll are all successfully resolved.

## 6. Operations
*   **CI Execution:** These tests will automatically be picked up by our existing `playwright test` setup added in Issue #81. No new CI configuration is required.

## Verification Checklist
- [x] Aligns with Product Roadmap.
- [x] Specifies the exact scenarios (Shadow DOM, SPA, Infinite Scroll).
- [x] Re-uses the Playwright runner from #81.
- [x] Avoids external dependencies (uses local fixtures).
