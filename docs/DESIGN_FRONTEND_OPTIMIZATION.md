# Design: Frontend DOM Observer Optimization

## Overview
As LinkLens expands from YouTube-only detection to universal link detection, the overhead of scanning every `<a>` tag on every DOM mutation increases. On complex pages with infinite scroll (e.g., Twitter, Reddit), this can lead to layout thrashing, high CPU usage, and a degraded user experience.

## Objectives
- **Minimize Layout Thrashing:** Avoid repetitive DOM measurements during scans.
- **Reduce CPU Usage:** Use non-blocking execution for scanning logic.
- **Precision Scanning:** Only scan new or modified portions of the DOM.
- **Stability:** Ensure replacements don't trigger infinite mutation loops.

## Proposed Architecture

### 1. Throttling and Debouncing
Currently, we use a simple `setTimeout` for debouncing. We will move to a more robust approach:
- **`requestIdleCallback`**: Schedule scans during the browser's idle periods to prevent frame drops.
- **Time-based Throttling**: Ensure scans don't happen more than once every 500ms-1s, even if many mutations occur.

### 2. Targeted Mutation Analysis
Instead of `Array.from(document.querySelectorAll('a'))` on every mutation:
- **Analyze Mutation Records**: Extract only the `addedNodes` from the `MutationRecord` list.
- **Sub-tree Scanning**: Only scan the `<a>` tags within the newly added elements rather than the entire document.

### 3. Link Caching (State Management)
- **`processedLinks`**: Maintain a `Set` of URLs (hrefs) that have already been handled or rejected by heuristics.
- **`pendingResolutions`**: Track URLs currently in flight to the backend to avoid duplicate requests for the same link appearing multiple times.

### 4. Batched Backend Requests
- Collect URLs found during a scan and wait a brief window (e.g., 200ms) before sending the backend request to batch multiple links from a single page load or scroll event into one network call.

## Implementation Plan

### Extension (content.ts)
- Implement a `MutationManager` class to encapsulate observer logic.
- Replace global `querySelectorAll` with targeted `element.querySelectorAll` on `addedNodes`.
- Integrate `requestIdleCallback` for the `scanAndReplace` execution.
- Implement a 100-item "LRU-style" cache for recently seen but non-matching links to quickly skip them.

## Verification & Metrics
- **Performance Testing:** Use Chrome DevTools Performance tab to measure "Long Tasks" before and after optimization.
- **Stress Test:** Load a page with 1,000+ links and verify scroll smoothness.

## Test Cases
- [ ] **Infinite Scroll:** Scroll through 10 pages of a dynamic site; verify backend requests are batched.
- [ ] **Dynamic Content:** Adding a single link via console (`document.body.appendChild(...)`) should trigger a targeted scan.
- [ ] **Large DOM:** Verify no noticeable lag on a site like Wikipedia.
