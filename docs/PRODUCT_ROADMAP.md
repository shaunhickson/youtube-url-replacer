# Product Roadmap: LinkLens

**Current Version:** 0.2 (Beta Candidate)
**Last Updated:** 2026-02-16

## 1. Product Vision & Strategy

### The "Big Picture" Goal
**To make the web transparent by default.**
The modern web is cluttered with opaque identifiers: shortened links (`bit.ly/xyz`), tracking redirects (`t.co/...`), and raw IDs (`youtube.com/watch?v=...`). Users are forced to "click and hope." LinkLens acts as the "Alt Text" for the internet's connective tissue, restoring context, safety, and clarity to every link before you interact with it.

### Strategic Pillars
1.  **Public Utility (The "Trust" Engine):**
    *   **Goal:** Provide the fastest, most private link resolver for the public web.
    *   **Monetization:** Free / Open Source (Community Edition).
    *   **Value:** Builds the user base, refines the parsing engine, and establishes the brand as "Safe."

2.  **Enterprise Context (The "Value" Engine):**
    *   **Goal:** Bring context to internal tools where clarity equals productivity.
    *   **Use Case:** A developer sees `LINEAR-123` in Slack. LinkLens resolves it to *"Fix Critical Auth Bug (High Priority)"* securely.
    *   **Monetization:** SaaS subscriptions for teams (Private Resolvers, SSO, Audit Logs).

## 2. Execution Roadmap

### Phase 1: Robustness & Scale (Completed)
*Goal: Ensure the foundation is solid, secure, and ready for traffic.*

*   [x] **Rate Limiting:** Protect the backend from abuse (Issue #13).
*   [x] **Testing & CI/CD:** Establish a reliable deployment pipeline with comprehensive coverage (Issue #8).
*   [x] **Security Audit:** Harden the system against SSRF and XSS (Issue #10, #33).
*   [x] **Observability:** Structured JSON logging for Cloud Run (Issue #15).
*   [x] **DevEx:** Comprehensive Makefile targets (Issue #16).
*   [x] **Project README:** Comprehensive getting started guide and architecture docs (Issue #41).
*   [x] **Collaboration Strategy:** Define how agents and humans work together (Issue #12).

### Phase 2: The "Universal" Pivot (Current Focus)
*Goal: Move beyond YouTube to support the wider web.*

*   [x] **Modular Architecture:** Refactor backend to support pluggable resolvers (See `docs/DESIGN_MODULAR_RESOLVER.md`).
*   [x] **Generic OpenGraph Resolver:** "Catch-all" support for any website with `og:title` tags.
*   [x] **Universal Link Detection:** Heuristics to identify "raw" URLs in text without false positives (Issue #18).
*   [x] **Frontend Optimization:** Efficient DOM scanning (Issue #17).
*   [ ] **Privacy Controls:** Allow/Block lists for domains (Issue #19).
*   [x] **URL Unshortener:** Automatically unwrap `bit.ly`, `t.co`, etc., to show the final destination.
*   [ ] **Marketing Website:** Landing page with live demo, features, and download links (Issue #39).
*   [ ] **Visual Update:** Extension UI/UX improvements to differentiate link types (icons, tooltips).

### Phase 3: Deep Integrations & Enterprise Value
*Goal: Provide specific, high-value data for popular platforms and teams.*

*   [ ] **GitHub Resolver:** Show repo description, stars, and language.
*   [ ] **Social Media:** Better context for X (Twitter), LinkedIn, etc.
*   [ ] **Enterprise Pilot:** Integration for private tools (Jira, Linear, Notion) via self-hosted or authenticated resolvers.

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
*   **Privacy First:** We resolve links, we don't track users. Logs never contain full URLs.
*   **Low Latency:** Resolutions must happen in milliseconds.
*   **Fail Open:** If the backend is down, the user just sees the original link. The web must not break.
