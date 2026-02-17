# Design: URL Unshortener Resolver

## Overview
Shortened URLs (e.g., `bit.ly/xyz`, `t.co/abc`) obscure the final destination of a link, creating a security risk and a lack of transparency. This resolver will expand these links to their final destination and provide the title of the resulting page.

## Objectives
- **Transparency:** Reveal the final destination URL and title.
- **Safety:** Protect users from clicking "mystery" links by providing context beforehand.
- **Robustness:** Handle multiple levels of redirects and detect redirect loops.

## Proposed Architecture

### 1. Implementation Strategy
The `UnshortenerResolver` will be a special type of resolver in the backend. Unlike the `YouTubeResolver`, which uses an API, the `UnshortenerResolver` will perform HTTP `HEAD` or `GET` requests to follow the redirect chain.

### 2. Detection (Heuristics)
The resolver will trigger for known shortener domains.
- **Domains:** `bit.ly`, `t.co`, `tinyurl.com`, `is.gd`, `buff.ly`, `goo.gl`, etc.
- **Generic Trigger:** Any URL that returns a 3xx status code (though we should start with a known list to avoid unnecessary overhead).

### 3. Resolution Logic
1.  **Normalization:** Ensure the URL is valid.
2.  **Redirect Following:**
    - Use an `http.Client` with `CheckRedirect` configured to limit the maximum number of hops (e.g., 5).
    - Detect infinite loops.
    - Maintain SSRF protection at every hop using `SafeTransport`.
3.  **Title Extraction:** Once the final destination is reached:
    - If the final URL has a specialized resolver (e.g., YouTube), delegate to it.
    - Otherwise, use the `OpenGraphResolver` logic to extract the title.

### 4. Integration with Resolver Manager
The `UnshortenerResolver` should ideally run *before* the generic `OpenGraphResolver` but potentially *after* specialized resolvers like `YouTubeResolver` (if the short link is known to be a specific platform). However, most shorteners are generic, so it might act as a pre-processor.

**Refined Flow in Manager:**
1.  Check Cache for original URL.
2.  If missing, try specialized resolvers (YouTube, etc.).
3.  If no specialized resolver matches, try `UnshortenerResolver`.
4.  If `UnshortenerResolver` expands the link:
    - Recursively resolve the *new* URL (this allows a `bit.ly` link to resolve via the `YouTubeResolver` if it points to a video).
5.  Fallback to `OpenGraphResolver`.

## Safety & Performance
- **Max Redirects:** Limit to 5 hops to prevent resource exhaustion.
- **SSRF Protection:** Re-validate the IP address of every hop in the redirect chain.
- **Timeout:** Total time for the entire chain must not exceed the global resolver timeout (e.g., 2s).
- **User-Agent:** Use the standard LinkLens User-Agent.

## Test Cases
- [ ] `https://bit.ly/3x86n7r` -> Should resolve to the final page title (e.g., "Google").
- [ ] Redirect Loop: `https://site-a.com` -> `https://site-b.com` -> `https://site-a.com` -> Should fail gracefully.
- [ ] Deep Redirect: 10 hops -> Should stop at 5 and return the intermediate title or an error.
- [ ] Short link pointing to a YouTube video -> Should resolve with "[YT] Video Title".
