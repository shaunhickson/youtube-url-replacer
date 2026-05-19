# Design Document: Backend Integration Test Suite (Issue #83)

## 1. Product Alignment
*   **Roadmap Check:** Addresses Phase 4, Issue #83: "Backend Integration Test Suite" in `docs/PRODUCT_ROADMAP.md`.
*   **Scope:** Implement Go integration tests for the `backend/` component. These tests will verify the full lifecycle of a URL resolution request from the API layer down to the external HTTP fetcher. We will use `httptest` to mock external API responses (like YouTube or generic websites) to keep tests fast, reliable, and decoupled from real-world network fluctuations.
*   **User Value:** Ensures that backend changes (e.g., adding new resolvers or refactoring middleware) do not break the core resolution logic, guaranteeing accurate link titles for the user.

## 2. Architecture & Patterns
*   **Modularity:** Integration tests will be housed alongside the backend code, likely in a `backend/integration_test.go` or `backend/api_test.go` file.
*   **Simplicity:**
    *   We will use the standard `testing` and `net/http/httptest` packages. No external testing frameworks (like Ginkgo or Testify) will be introduced unless strictly necessary, adhering to Go idioms.
    *   We will start an in-memory test server (`httptest.NewServer`) that intercepts outbound requests from our resolvers and returns deterministic HTML or JSON responses.
*   **Test Scenarios:**
    1.  **YouTube Resolution:** Simulate a request with a `youtu.be` link. The mocked external server returns a mocked YouTube Data API response. Assert the final returned title matches.
    2.  **Generic OpenGraph Resolution:** Simulate a generic link. The mocked external server returns HTML with `<meta property="og:title" content="...">`. Assert the final returned title matches.
    3.  **Invalid URL Handling:** Send malformed URLs and assert proper 400 Bad Request responses.

## 3. Security & Safety
*   **SSRF Validation:** The integration tests will specifically verify that our SSRF protection logic correctly rejects internal IP ranges (e.g., `127.0.0.1` or `169.254.169.254`) by simulating requests to these URLs and asserting a 400 error.

## 4. Scalability & Performance
*   **Fast Execution:** By using `httptest.NewServer` instead of real network calls, the test suite will run in milliseconds, keeping CI fast.
*   **Parallelism:** We will use `t.Parallel()` where appropriate to run tests concurrently.

## 5. Testability
*   **Target Coverage:** The integration test suite will cover the `main.go` HTTP handlers, the `Resolver` interface logic, and the HTTP Transport layer.
*   **Execution:** Executed via `make test-backend` or `go test -v ./...`.

## 6. Operations
*   **CI Execution:** Our existing `ci.yml` already runs `go test -v -race -cover ./...`. The new integration tests will be automatically picked up and run.

## Verification Checklist
- [x] Aligns with Product Roadmap.
- [x] Defines scope (HTTP handler to external mock).
- [x] Uses existing modular architecture (Go stdlib `httptest`).
- [x] Explicitly tests security features (SSRF).
