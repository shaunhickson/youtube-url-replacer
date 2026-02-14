# Design: Modular URL Resolver Architecture

## Overview
The `youtube_replacer` backend currently has a hardcoded dependency on the YouTube API and specifically handles YouTube Video IDs. To support additional platforms (Vimeo, GitHub, Bitly, etc.), we need a modular architecture that allows for independent, pluggable resolvers.

## Objectives
- **Extensibility:** Add new platform support by simply adding a new file that implements an interface.
- **Maintainability:** Isolate platform-specific logic and API interactions.
- **Backward Compatibility:** Continue to support the existing YouTube Video ID resolution while transitioning to URL-based resolution.
- **Robustness:** Ensure a generic fallback mechanism for unknown URLs.

## Proposed Architecture

### 1. The Resolver Interface
Each platform support will be implemented as a `Resolver`.

```go
package resolvers

import (
    "context"
    "net/url"
)

type Result struct {
    Title       string
    Description string // Optional, for tooltips
    Platform    string // e.g., "YouTube", "GitHub"
}

type Resolver interface {
    // Name returns the unique identifier for this resolver
    Name() string
    
    // CanHandle returns true if this resolver can process the given URL
    CanHandle(u *url.URL) bool
    
    // Resolve returns a human-friendly title/description for the URL
    Resolve(ctx context.Context, u *url.URL) (*Result, error)
}
```

### 2. Resolver Manager (Registry)
A central manager will coordinate the resolution process.

```go
type ResolverManager struct {
    resolvers []Resolver
    cache     Cache
    httpClient *http.Client
}

func (m *ResolverManager) ResolveMulti(ctx context.Context, rawURLs []string) map[string]*Result {
    // 1. Normalize URLs
    // 2. Check Cache (using normalized URL as key)
    // 3. For each missing URL:
    //    a. Find the first Resolver that returns CanHandle(u) == true.
    //    b. Call resolver.Resolve(ctx, u).
    //    c. Store result in Cache.
    // 4. Return combined results.
}
```

### 3. Proposed File Structure
We will move resolvers into a dedicated package:
```
backend/
├── resolvers/
│   ├── interface.go      # Interface and Result struct
│   ├── manager.go        # ResolverManager implementation
│   ├── youtube.go        # Existing YouTube logic refactored
│   ├── github.go         # Future GitHub resolver
│   ├── opengraph.go      # Future generic fallback resolver
│   └── util.go           # Shared helpers (SSRF protection, HTTP client)
├── main.go               # Orchestration
├── handler.go            # HTTP API handling
└── cache.go              # Shared caching logic
```

### 4. API Changes (Backward Compatibility)
To support both the old "Video ID" approach and the new "URL" approach:

**Request Body:**
```json
{
  "urls": ["https://www.youtube.com/watch?v=dQw4w9WgXcQ", "https://github.com/shaunhickson/youtube-url-replacer"],
  "videoIds": ["dQw4w9WgXcQ"] // Legacy support
}
```

**Response Body:**
```json
{
  "titles": {
    "https://www.youtube.com/watch?v=dQw4w9WgXcQ": "Rick Astley - Never Gonna Give You Up (Official Music Video)",
    "https://github.com/shaunhickson/youtube-url-replacer": "shaunhickson/youtube-url-replacer: A browser extension... (★ 10)",
    "dQw4w9WgXcQ": "Rick Astley - Never Gonna Give You Up (Official Music Video)"
  }
}
```

### 5. Shared Utilities
Resolvers will benefit from shared logic implemented in `resolvers/util.go`:
- **SafeFetch:** An HTTP client wrapper with:
    - SSRF protection (blocking internal/private IP ranges).
    - Strict timeouts (e.g., 2 seconds).
    - User-Agent header (e.g., `youtube-url-replacer/1.0 (+https://github.com/shaunhickson/youtube-url-replacer)`).
- **MetadataExtractor:** Helper to parse HTML for `<title>`, `og:title`, etc.

## Refactoring Plan
1. **Define Interface:** Create `backend/resolvers/interface.go`.
2. **Implement Manager:** Create `backend/resolvers/manager.go`.
3. **Refactor YouTube:** Move existing `youtube.go` logic into `backend/resolvers/youtube.go` and implement the `Resolver` interface.
4. **Update Cache:** Update `Cache` to handle full URLs as keys (or create a new namespace).
5. **Update Handler:** Modify `backend/handler.go` to use `ResolverManager` and support the `urls` field in requests.
6. **Generic Fallback:** Implement the `OpenGraphResolver` to provide immediate value for all other links.

## Configuration
The system will be configurable via environment variables to control which resolvers are active and their respective settings:
- `ENABLED_RESOLVERS`: A comma-separated list of resolver names to enable (e.g., `youtube,github,opengraph`). If empty, all are enabled.
- `GITHUB_TOKEN`: Optional API token for the GitHub resolver to increase rate limits.
- `RESOLVER_TIMEOUT_MS`: Global timeout for a single resolution (default: 2000).
- `MAX_URLS_PER_REQUEST`: Maximum number of URLs to process in one API call (default: 50).

## Security Considerations
- **SSRF:** Crucial for the `OpenGraphResolver` and any resolver that fetches external content. We must validate that the resolved IP address of any URL is not a private, loopback, or link-local address.
- **Resource Exhaustion:** Limit the number of URLs processed in a single request and set strict timeouts for each resolution.
