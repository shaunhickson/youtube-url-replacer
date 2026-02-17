# Design: Enhanced UI/UX (Contextual Icons & Rich Tooltips)

## Overview
Currently, LinkLens replaces link text with a simple string (e.g., `[YT] Video Title`). While functional, it is visually jarring and provides limited information. This design proposes a "Surgical UI" approach that provides rich context with zero layout shift and minimal visual noise.

## 1. Design Philosophy: "Context at a Glance"
- **Non-Intrusive:** We must not break the original site's layout or font hierarchy.
- **Differentiated:** Users should instantly distinguish between a video, a repository, and a generic article.
- **Performance First:** UI elements are injected "Just-in-Time" and use optimized CSS/SVG.

## 2. Visual Elements

### A. Surgical Icons (Inline)
Instead of a bracketed prefix like `[YT]`, we will use small, grayscale-by-default SVG icons that take on the color of the surrounding text.
- **YouTube:** Play button icon.
- **GitHub:** Octocat/Repo icon.
- **Generic/Article:** Document icon.
- **Unshortened:** "Link" icon with a small arrow indicating expansion.

**Implementation:**
- Icons are injected as an `::before` pseudo-element or a small `<span>` with a background SVG mask. 
- **Constraint:** Icon height is fixed to `0.9em` to ensure it never exceeds the line-height of the parent text.

### B. Rich Tooltips (The "Lens")
The `title` attribute is too limited. We will implement a custom tooltip that appears on hover.
- **Content:** 
  - **Header:** Final Destination Domain (e.g., `youtube.com`).
  - **Title:** The resolved page title (Bold).
  - **Summary:** The first 150 characters of the `og:description` (if available).
- **Trigger:** 500ms hover delay to prevent "flicker" while moving the mouse across the page.

## 3. Technical Strategy: Shadow DOM
To solve the problem of "CSS bleeding" (where the host site's styles break our tooltips or our styles break the site), we will use a **Single-Instance Shadow DOM**.

1. **The Container:** A single `div` at the bottom of `document.body`.
2. **The Shadow Root:** All tooltip HTML and CSS live inside this root, isolated from the page.
3. **Positioning:** Use `getBoundingClientRect()` on the target link to position the isolated tooltip.

## 4. Proposed Data Flow (Backend to UI)
The backend `Result` struct already supports `Description` and `Platform`. We will utilize these fully:
```json
{
  "https://youtu.be/123": {
    "title": "Go Concurrency Patterns",
    "description": "Rob Pike discusses Google's approach to...",
    "platform": "youtube"
  }
}
```

## 5. Skepticism Mitigation (The "Why it won't be annoying" section)
- **Constraint 1 (No Layout Shift):** We only replace the text *after* resolution. To prevent the page jumping, the new text should ideally be similar in length or use CSS `text-overflow: ellipsis` if it exceeds a certain width.
- **Constraint 2 (Styling Match):** Our icons will use `fill: currentColor`. If the link is blue, the icon is blue. if the link is red, the icon is red. It will look like a native part of the site.
- **Constraint 3 (Accessibility):** We will maintain `aria-label` attributes on the links so screen readers read the resolved title properly.

## 6. Implementation Plan

### Extension
1. **CSS Modules:** Create an isolated CSS file for the Shadow DOM tooltip.
2. **UI Utility:** Create a `UIManager` class to handle the creation/positioning of the Shadow DOM.
3. **Content Script:** Update the "Update DOM" phase of `LinkLensOptimizer` to:
   - Apply the `link-lens-resolved` class to the anchor.
   - Inject the icon span.
   - Attach the hover event listeners for the rich tooltip.

## 7. Test Cases
- [ ] **Contrast Test:** Ensure icons are visible on both Dark and Light mode sites.
- [ ] **Overflow Test:** Verify long titles are truncated gracefully on narrow containers (e.g., sidebars).
- [ ] **Isolation Test:** Verify that a site's global `div { background: red !important; }` does not turn our tooltip red.
