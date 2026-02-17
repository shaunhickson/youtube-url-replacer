# Design: Testing Strategy & CI/CD Pipeline

## Overview
This document defines the quality assurance strategy for **LinkLens**. As we move toward a modular architecture and automated agent-based development, a reliable CI/CD pipeline is critical to prevent regressions.

## Goals
- **Automated Verification:** Every commit to a feature branch is verified.
- **Fast Feedback:** CI should run in < 5 minutes.
- **Full Coverage:** Testing at Unit, Integration, and E2E levels.
- **Safe Deployment:** Automated deployment to Cloud Run only on successful `main` builds.

## Testing Layers

### 1. Backend (Go)
*   **Unit Tests:** Use standard `testing` package.
    *   Target: `backend/**/*.go`.
    *   Mocking: Use interfaces (like `Cache`, `Resolver`) to mock dependencies.
*   **Linting:** `golangci-lint` (Standard, rigorous Go linting).

### 2. Frontend / Extension (TypeScript/React)
*   **Unit Tests:** Use **Vitest** + **React Testing Library**.
    *   Why Vitest? Faster than Jest, native Vite integration.
    *   Target: Components (`Popup.tsx`), Utilities (`content.ts` logic).
*   **Linting:** `eslint` with `typescript-eslint`.

### 3. End-to-End (E2E)
*   **Tool:** **Playwright**.
*   **Scope:**
    *   Load the extension in a headless browser.
    *   Navigate to a test page (`test.html`).
    *   Verify that links are replaced.
*   *Note:* E2E tests are slower and flaky. We will run these on PRs but maybe not every commit if they get too slow.

## CI Pipeline (GitHub Actions)

We will implement a single workflow `.github/workflows/ci.yml`.

### Triggers
*   `push` to `main`.
*   `pull_request` to `main`.

### Jobs

#### 1. `backend-test`
*   **Image:** `golang:1.23`
*   **Steps:**
    1.  Checkout.
    2.  `go mod download`.
    3.  `golangci-lint run`.
    4.  `go test -v -race -cover ./...`.

#### 2. `frontend-test`
*   **Image:** `node:20`
*   **Steps:**
    1.  Checkout.
    2.  `npm ci` (extension dir).
    3.  `npm run lint`.
    4.  `npm run test`.
    5.  `npm run build` (Verify build works).

#### 3. `e2e-test` (Future/Optional for now)
*   *Deferred to Phase 2 to keep velocity high initially, unless critical.*

#### 4. `deploy` (Only on `main`)
*   **Needs:** `[backend-test, frontend-test]`
*   **Steps:**
    1.  Authenticate to Google Cloud.
    2.  Build & Push Docker Image.
    3.  Deploy to Cloud Run.

## Implementation Roadmap
1.  **Frontend Setup:** Install `vitest`, `jsdom`, `@testing-library/react`.
2.  **Backend Setup:** Ensure `golangci-lint` config exists.
3.  **GitHub Action:** Create `.github/workflows/ci.yml`.
4.  **Branch Protection:** Require CI to pass before merging.
