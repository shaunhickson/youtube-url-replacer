# Design: Universal Link Detection

## Overview
Currently, LinkLens only detects and replaces YouTube URLs. To achieve our goal of making the web transparent, we need to detect "raw" URLs from any domain where the link text provides no context (e.g., `https://example.com/p/123`).

## Objectives
- **Broad Support:** Detect links from any domain.
- **Contextual Integrity:** Only replace links where the visible text is a "raw" URL, avoiding interference with intentionally labeled links (e.g., "Click Here").
- **Efficiency:** Scanning the DOM should be fast and not cause layout thrashing.
- **Privacy:** Avoid sending unnecessary data to the backend.

## Heuristics for "Raw" URLs
A link is considered a candidate for resolution if:
1.  **Protocol:** It uses `http://` or `https://`.
2.  **Text Match:** The `innerText` of the anchor tag matches a URL pattern.
3.  **HREF Match:** The `innerText` is substantially similar to the `href`.
    - Example: `<a href="https://example.com/foo">https://example.com/foo</a>` -> **MATCH**
    - Example: `<a href="https://example.com/foo">example.com/foo</a>` -> **MATCH**
    - Example: `<a href="https://example.com/foo">My Blog Post</a>` -> **NO MATCH**

## Proposed Implementation (Extension)

### 1. Generic URL Regex
We need a robust regex to identify URLs in text.
```typescript
const URL_REGEX = /^(https?:\/\/)?([\w.-]+)\.([a-z]{2,})(:\d+)?(\/\S*)?$/i;
```

### 2. Scanning Logic
The `scanAndReplace` function will be updated:
- Iterate through all `<a>` tags.
- For each tag:
    - Get `href` and `innerText` (trimmed).
    - If `innerText` matches `URL_REGEX`:
        - Normalize both `href` and `innerText` (remove protocol, trailing slashes).
        - If they match or are very similar:
            - Add to `urlsToFetch`.

### 3. API Communication
Switch from sending `videoIds` to sending `urls`.
```json
{
  "urls": ["https://example.com/raw-link", "https://youtube.com/watch?v=..."]
}
```

## Proposed Implementation (Backend)
The backend already supports a `urls` field and has an `OpenGraphResolver`. No major changes are needed to the backend logic, but we should ensure the `OpenGraphResolver` is the last one in the chain.

## Security & Privacy
- **SSRF Protection:** The backend already implements `SafeTransport` to prevent resolving internal IP addresses.
- **Data Minimization:** Only URLs that pass the "raw" heuristic are sent to the backend.
- **User Privacy:** No cookies or identifying headers are sent with the resolution request.

## Test Cases
- [ ] `https://google.com` as text -> Should resolve to "Google".
- [ ] `bit.ly/xxxx` as text -> Should resolve to the destination title.
- [ ] `Click here` with a hidden URL -> Should NOT resolve.
- [ ] YouTube link with "My Video" -> Should NOT resolve (consistent with current behavior).
- [ ] YouTube link with `https://youtu.be/...` -> Should resolve.
