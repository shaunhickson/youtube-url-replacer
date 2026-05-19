# Design Document: Playwright E2E Testing Framework (Issue #81)

## 1. Product Alignment
*   **Roadmap Check:** This design directly addresses Phase 4, Issue #81: "Setup Playwright for Browser Extension E2E Testing" in `docs/PRODUCT_ROADMAP.md`.
*   **Scope:** Introduce Playwright as the End-to-End (E2E) testing framework for the LinkLens browser extension. This scope is strictly limited to adding the test harness, configuration, CI integration, and a single "happy path" smoke test to prove the framework functions. Test fixtures for complex DOM structures (Issue #82) will be handled separately.
*   **User Value:** E2E testing ensures that the extension successfully injects its content scripts, detects links, communicates with the backend, and modifies the DOM correctly without breaking the user's browsing experience. This provides confidence when shipping updates.

## 2. Architecture & Patterns
*   **Simplicity:** We will use `@playwright/test` within the `extension/` directory. Playwright has native support for loading Chrome extensions via `chromium.launchPersistentContext` with the `--disable-extensions-except` and `--load-extension` flags.
*   **Dependencies:**
    *   `@playwright/test`: The core test runner.
*   **Integration:**
    *   We will point Playwright to the compiled extension `dist/` directory.
    *   The test suite will load local HTML files (which will be expanded in Issue #82) to avoid relying on external websites that may change or block automated browsers.

## 3. Security & Safety
*   **Isolation:** E2E tests will run in an isolated Playwright browser context.
*   **Network:** Tests should eventually intercept backend requests using Playwright's `page.route()` to avoid hammering the production API, mocking successful and failed resolutions. This keeps tests fast, deterministic, and safe.

## 4. Scalability & Performance
*   **Execution Speed:** Mocking backend API calls (via `page.route()`) ensures tests run in milliseconds rather than waiting for actual network requests, minimizing CI duration.
*   **Parallelism:** Playwright runs tests in parallel by default, scaling well as the test suite grows.

## 5. Testability
*   **The Framework Itself:** This design *is* the testability layer for the frontend.
*   **Initial Test Coverage:** We will implement one "happy path" E2E test in `extension/tests/e2e/basic.spec.ts` that:
    1. Loads the extension in a Chromium context.
    2. Opens a simple local HTML file containing a `youtube.com/watch` link.
    3. Mocks the backend response.
    4. Asserts that the LinkLens tooltip/icon is injected and the link text is modified correctly.

## 6. Operations
*   **Configuration:** A `playwright.config.ts` will be created in the `extension/` directory.
*   **CI Integration:** We will update `.github/workflows/` (or create a new workflow) to run `npx playwright test` on every push to a feature branch.
*   **Running Locally:** Developers will run `npm run test:e2e` to execute the suite.

## Verification Checklist
- [x] Aligns with Product Roadmap.
- [x] Adheres to single-issue focus (fixtures deferred to #82).
- [x] Uses existing modular structure.
- [x] Does not introduce security risks.
- [x] Details CI operationalization.
