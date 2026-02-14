# User Research: Expansion beyond YouTube URLs

## Overview
The `youtube_replacer` currently specializes in transforming YouTube links into human-friendly titles. This research explores expanding this capability to other common URL types that are frequently shared but lack immediate context or readability.

## Goals
- Identify high-value URL types for expansion.
- Understand the technical requirements for "unfurling" these links.
- Propose a strategy for prioritizing and implementing these new resolvers.

## High-Value URL Categories

### 1. Shortened URLs (Bitly, TinyURL, t.co)
- **Problem:** These URLs are completely opaque. Users don't know where they lead until they click, which can be a security concern.
- **Proposed Solution:** "Unwrap" the redirect and display either the final destination's `<title>` or at least the destination domain (e.g., `bit.ly/xxxx` -> `[nytimes.com] Article Title`).
- **Technical Implementation:** Perform a `HEAD` request to follow redirects, then fetch the final page's metadata.

### 2. Other Video Platforms (Vimeo, Dailymotion, Twitch)
- **Problem:** Similar to YouTube, these links are just IDs (e.g., `vimeo.com/123456789`).
- **Proposed Solution:** Fetch the video title and possibly the uploader's name.
- **Technical Implementation:** Use platform-specific APIs or generic OpenGraph (`og:title`) tags.

### 3. GitHub Repositories
- **Problem:** `github.com/shaunhickson/youtube-url-replacer` is informative but doesn't convey the project's purpose or popularity.
- **Proposed Solution:** `Owner/Repo: Brief Description (Stars)`.
- **Technical Implementation:** Use the GitHub API (requires an API key/token for higher rate limits) or scrape the repo page for meta tags.

### 4. Social Media Posts (X/Twitter, LinkedIn, Instagram)
- **Problem:** Shared posts often appear as `x.com/user/status/123...`, providing zero context about the content.
- **Proposed Solution:** Extract a snippet of the post text or the author's name and the post type (e.g., `X: @User "Post content snippet..."`).
- **Technical Implementation:** These platforms often have strict anti-scraping measures. Using official APIs or specific "embed" endpoints might be necessary.

### 5. Cloud Documents (Google Docs, Dropbox, Figma)
- **Problem:** Links to collaborative tools are often just strings of random characters.
- **Proposed Solution:** Display the document title and the application type.
- **Technical Implementation:** Most of these provide OpenGraph metadata even if the content itself is private/restricted.

## Technical Strategy: Generic vs. Specific Resolvers

### Generic OpenGraph Resolver
As a fallback for any unknown URL, the backend could attempt to fetch the page and look for OpenGraph tags:
- `og:title`
- `og:site_name`
- `og:description`

This would provide a "catch-all" human-friendly name for almost any modern website.

### Specialized Resolvers
For high-traffic sites (GitHub, Vimeo, etc.), we should implement specialized logic to ensure the most relevant data is extracted (e.g., star counts for GitHub, view counts for Vimeo).

## Implementation Roadmap
1. **Phase 1: URL Shortener Unwrapping.** Add a middleware or helper to resolve redirects before processing.
2. **Phase 2: Generic Metadata Fallback.** Implement a crawler that extracts `<title>` or `og:title` for any non-YouTube link.
3. **Phase 3: Specialized Resolvers.** Roll out dedicated support for GitHub and Vimeo.

## Security Considerations
- **SSRF Protection:** Ensure the backend cannot be used to scan internal networks by resolving URLs.
- **Crawler Identity:** Use a clear User-Agent and respect `robots.txt` where possible.
- **Timeout Management:** External sites may be slow; aggressive timeouts are required to prevent backend hangs.
