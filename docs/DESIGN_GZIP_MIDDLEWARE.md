# Design: Gzip Compression Middleware

**Issue:** #46
**Status:** Approved

## Overview
Currently, the LinkLens Go backend returns uncompressed JSON responses. Adding Gzip compression middleware will significantly reduce the payload size for large OpenGraph metadata resolutions, lowering latency for users on slower connections and reducing egress bandwidth costs on Cloud Run.

## Implementation Details
- **Component:** `backend/middleware/gzip.go`
- **Logic:**
  - Check the incoming `Accept-Encoding` header for `gzip`.
  - If present, wrap the `http.ResponseWriter` with a custom writer that pipes output through `compress/gzip.Writer`.
  - Set the outgoing `Content-Encoding` header to `gzip`.
  - Ensure the gzip writer is flushed and closed when the response completes.
- **Integration:** The middleware will be added to the routing chain in `main.go`, alongside the existing `RequestLogger` and `RateLimiter`.
- **Testing:** Include unit tests to verify that compressed responses are valid and that clients without `Accept-Encoding: gzip` receive uncompressed responses safely.
