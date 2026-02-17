# Agent Design Review Checklist

## Overview
This checklist is a mandatory self-review tool for AI Agents drafting design documents (`docs/DESIGN_*.md`). Before submitting a Design PR for human review, the Agent must verify the design against these criteria.

## 1. Product Alignment
- [ ] **Roadmap Check:** Does this design directly address an item in `docs/PRODUCT_ROADMAP.md`?
- [ ] **Scope:** Is the scope clearly defined? Does it avoid "scope creep"?
- [ ] **User Value:** Does the document explain *why* this is valuable to the user?

## 2. Architecture & Patterns
- [ ] **Modularity:** Does the design respect the existing modular architecture (e.g., `resolvers/` interface)?
- [ ] **Simplicity:** Is this the simplest possible solution? (YAGNI principle).
- [ ] **Dependencies:** Are new libraries/dependencies justified? (Check if existing tools can solve it).

## 3. Security & Safety
- [ ] **SSRF:** If fetching URLs, is there protection against internal network scanning?
- [ ] **Input Validation:** Is all user input (payloads, params) validated?
- [ ] **Auth/AuthZ:** Does it respect existing authentication (if any)?
- [ ] **Secrets:** Does it require new secrets? Are they documented as Env Vars (never in code)?

## 4. Scalability & Performance
- [ ] **Rate Limiting:** Does the design consider abuse vectors?
- [ ] **Database:** Are Firestore reads/writes minimized?
- [ ] **Latency:** Will this add significant latency to the user experience?

## 5. Testability
- [ ] **Unit Tests:** Does the design specify what unit tests will be written?
- [ ] **Integration:** How will we verify this works with external APIs?
- [ ] **E2E:** Is there a critical user flow that needs E2E testing?

## 6. Operations
- [ ] **Configuration:** Are all tunables exposed as Environment Variables?
- [ ] **Observability:** Does the design include logging (structured JSON) and metrics?
- [ ] **Rollback:** Is there a safe way to disable this feature if it breaks? (Feature Flag).

## Self-Correction Prompt
*If you find gaps during this review, **update the design document** before opening the PR.*
