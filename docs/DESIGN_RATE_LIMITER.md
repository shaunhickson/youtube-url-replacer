# Design: Rate Limiter for youtube_replacer Backend

## Overview
Implement a middleware-based rate limiter to protect the `/resolve` endpoint from abuse, resource exhaustion, and DDoS attacks.

## Goals
- Prevent single IP addresses from flooding the backend with requests.
- Mitigate cache exhaustion/poisoning by limiting payload size.
- Ensure the backend remains responsive for legitimate users.
- Provide clear feedback (HTTP 429) when limits are exceeded.

## Attack Vectors
1. **Endpoint Flooding:** Rapid-fire requests to `/resolve`.
2. **Batch Abuse:** Large arrays of `videoIds` in a single request to force massive YouTube API calls.
3. **Distributed Requests:** Many IPs making a moderate number of requests (partially mitigated by per-IP limits).

## Proposed Strategy

### 1. IP-based Rate Limiting (Middleware)
Use a token bucket algorithm to limit requests per IP.
- **Library:** `golang.org/x/time/rate`
- **Default Limit:** 60 requests per minute.
- **Burst:** 20 requests.
- **Storage:** In-memory LRU cache to track visitor states.

### 2. Payload Limitation
Limit the number of Video IDs allowed in a single `/resolve` request.
- **Max Video IDs:** 50 per request.

### 3. Response Headers
Include standard rate-limiting headers:
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`
- `Retry-After` (on 429)

## Technical Implementation

### Middleware Components
A new file `backend/middleware.go` will be created:
- `RateLimiter`: A struct managing the visitor map and limits.
- `LimitMiddleware`: A function that wraps `http.Handler` and performs the check.

### Handling 429
When a limit is reached, return:
- **Status Code:** `429 Too Many Requests`
- **Body:** `{"error": "Too many requests. Please try again later."}`

### Configuration
Add the following environment variables:
- `RATE_LIMIT_RPM` (default: 60)
- `RATE_LIMIT_BURST` (default: 20)
- `MAX_VIDEO_IDS` (default: 50)

## Future Improvements
- **Redis Support:** For distributed rate limiting across multiple Cloud Run instances.
- **Cloud Armor:** Integration with GCP Cloud Armor for edge-level protection.
- **User-based Limiting:** If authentication is added in the future.
