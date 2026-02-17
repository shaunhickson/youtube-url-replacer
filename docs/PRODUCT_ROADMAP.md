# Product Roadmap: LinkLens

**Current Version:** 0.1 (Alpha)
**Last Updated:** 2026-02-16

## 1. Product Identity

*   **Name:** **LinkLens**
*   **Tagline:** "Know before you click."
*   **Mission:** To bring context and clarity to the opaque web by instantly transforming raw URLs into human-readable insights.

## 2. Strategic Phases

### Phase 1: Robustness & Scale (Current Focus)
*Goal: Ensure the foundation is solid, secure, and ready for traffic.*

*   [x] **Rate Limiting:** Protect the backend from abuse (Issue #13).
*   [x] **Testing & CI/CD:** Establish a reliable deployment pipeline with comprehensive coverage (Issue #8).
*   [x] **Security Audit (Design):** Harden the system, specifically against SSRF as we prepare to fetch external URLs (Issue #10).
*   [ ] **Security Implementation:** Implement SSRF protection and hardening (Issue #33).
*   [x] **Collaboration Strategy:** Define how agents and humans work together (Issue #12).
*   [ ] **Observability:** Structured JSON logging for Cloud Run (Issue #15).
*   [ ] **DevEx:** Comprehensive Makefile targets (test, lint, docker) (Issue #16).

### Phase 2: The "Universal" Pivot
*Goal: Move beyond YouTube to support the wider web.*

*   [ ] **Modular Architecture:** Refactor backend to support pluggable resolvers (See `docs/DESIGN_MODULAR_RESOLVER.md`).
*   [ ] **Generic OpenGraph Resolver:** "Catch-all" support for any website with `og:title` tags.
*   [ ] **Universal Link Detection:** Heuristics to identify "raw" URLs in text without false positives (Issue #18).
*   [ ] **Frontend Optimization:** Efficient DOM scanning (Issue #17).
*   [ ] **Privacy Controls:** Allow/Block lists for domains (Issue #19).
*   [ ] **URL Unshortener:** Automatically unwrap `bit.ly`, `t.co`, etc., to show the final destination.
*   [ ] **Visual Update:** Extension UI/UX improvements to differentiate link types (icons, tooltips).

### Phase 3: Deep Integrations & specialized Value
*Goal: Provide specific, high-value data for popular platforms.*

*   [ ] **GitHub Resolver:** Show repo description, stars, and language.
*   [ ] **Social Media:** Better context for X (Twitter), LinkedIn, etc.
*   [ ] **Enterprise (Future):** Potential integrations for private tools like Jira, Notion, or Linear.

## 3. Marketing & Messaging

**Elevator Pitch:**
> The Web is full of mystery links. Solve them with LinkLens.
> You see a link: `https://youtu.be/dQw4w9WgXcQ`. Is it a tutorial? A rickroll?
> With LinkLens, you see: **Rick Astley - Never Gonna Give You Up (Official Music Video)**.
> Works everywhere you browse. Secure, fast, and privacy-focused.

**Key Value Props:**
*   **Security:** Know where a link goes before you click it.
*   **Context:** Get the title/summary instantly.
*   **Productivity:** Save time clicking back and forth.

## 4. Technical Constraints & Pillars
*   **Design First:** Every roadmap item must have a distinct design phase producing an artifact (doc, schema, or prototype) before implementation begins.
*   **Privacy First:** We resolve links, we don't track users.
*   **Low Latency:** Resolutions must happen in milliseconds.
*   **Fail Open:** If the backend is down, the user just sees the original link. The web must not break.
