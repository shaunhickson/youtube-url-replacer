# Design: Rate Limiter for LinkLens Backend

## Overview
Implement a multi-layered rate limiting strategy to protect the `/resolve` endpoint from abuse, resource exhaustion, and DDoS attacks. This design accounts for the new modular `ResolverManager` architecture capable of handling both legacy YouTube Video IDs and arbitrary URLs.

## Goals
- **Availability:** Ensure the backend remains responsive for legitimate users.
- **Fairness:** Prevent single IP addresses from monopolizing resources.
- **Cost Control:** Limit the number of external API calls (YouTube, GitHub, etc.) triggered by a single request.
- **Observability:** Provide clear feedback (HTTP 429) and metrics on blocked requests.

## Attack Vectors & Mitigation
| Vector | Mitigation Layer | Strategy |
| :--- | :--- | :--- |
| **Volumetric DDoS** | Layer 1: IP Rate Limit | Token Bucket (per IP) |
| **Massive Payloads** | Layer 2: Body Size Limit | `http.MaxBytesReader` |
| **Batch Abuse** | Layer 3: Item Count Limit | Max 50 items (`len(videoIds) + len(urls)`) |

## Proposed Implementation

### 1. IP-based Rate Limiting (Layer 1)
Use a token bucket algorithm to limit requests per IP address.
- **Library:** `golang.org/x/time/rate`
- **Storage:** In-memory `sync.Map` or LRU cache (expiring old IPs to prevent memory leaks).
- **Configuration:**
    - `RATE_LIMIT_RPM`: Requests per minute (Default: **60**).
    - `RATE_LIMIT_BURST`: Max burst size (Default: **20**).

**Behavior:**
- Middleware extracts IP from `X-Forwarded-For` (since we run on Cloud Run) or `RemoteAddr`.
- If limit exceeded: Immediate HTTP 429 response.

### 2. Payload Limitation (Layers 2 & 3)

**Body Size:**
- Enforce a strict limit on the request body size (e.g., 10KB) to prevent memory exhaustion before parsing.

**Item Count:**
- The `ResolveRequest` struct contains both `videoIds` and `urls`.
- The handler must validate that `len(req.VideoIDs) + len(req.URLs) <= MAX_ITEMS_PER_REQUEST`.
- **Default Limit:** **50** items total per request.

### 3. Response Headers
Conform to RFC 6585 and common standards:
- `X-RateLimit-Limit`: The request limit per minute.
- `X-RateLimit-Remaining`: Requests left in the current window.
- `X-RateLimit-Reset`: Seconds until the limit resets.
- `Retry-After`: Seconds to wait before retrying (only on 429).

### 4. Code Structure

We will introduce a new package or file `backend/middleware/ratelimit.go` (or keep simple in `backend/middleware.go`).

```go
package main

// Middleware wrapper
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. IP Check
        if !limiter.Allow(getIP(r)) {
            w.WriteHeader(http.StatusTooManyRequests)
            return
        }
        
        // 2. Size Check (handled by http.MaxBytesReader in main or here)
        
        next.ServeHTTP(w, r)
    })
}
```

**Validation Logic (in `ServeHTTP` or separate validation middleware):**
```go
if len(req.VideoIDs) + len(req.URLs) > MaxItems {
    http.Error(w, "Too many items", http.StatusRequestEntityTooLarge)
    return
}
```

### 5. Configuration (Env Vars)
| Variable | Default | Description |
| :--- | :--- | :--- |
| `RATE_LIMIT_RPM` | 60 | Requests per minute per IP |
| `RATE_LIMIT_BURST` | 20 | Token bucket burst capacity |
| `MAX_ITEMS_PER_REQUEST` | 50 | Max combined URLs/IDs per request |
| `MAX_BODY_BYTES` | 10240 | Max request body size in bytes (10KB) |

## Future Improvements
- **Distributed Limiting:** Use Redis (MemoryStore) to share limits across Cloud Run instances if we scale out.
- **API Keys:** If we introduce a public API, switch to Key-based limiting instead of IP-based.
