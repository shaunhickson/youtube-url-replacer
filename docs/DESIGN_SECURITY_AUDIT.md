# Design: Comprehensive Security Audit & Hardening

## Overview
As LinkLens evolves from a YouTube-specific tool to a general-purpose URL resolver, the attack surface expands significantly. This document outlines the security architecture required to safely fetch and parse arbitrary content from the internet.

## Threat Model

### 1. Server-Side Request Forgery (SSRF)
**Risk:** High
**Vector:** An attacker submits a URL like `http://169.254.169.254/computeMetadata/v1/` or `http://localhost:8080/health`.
**Impact:** Access to internal cloud metadata (GCP secrets), local services, or loopback interfaces.
**Mitigation:**
- **Strict IP Validation:** Resolve the hostname *before* making the request.
- **Blocklist:** Deny private IP ranges (RFC 1918, RFC 4193, RFC 4291).
- **DNS Rebinding Protection:** Ensure the IP doesn't change between validation and fetch (use a custom `DialContext`).
- **Redirect Following:** Validate the IP of every redirect target.

### 2. Denial of Service (DoS) & Resource Exhaustion
**Risk:** Medium
**Vector:**
- **Slowloris:** Holding connections open.
- **Large Payloads:** Serving a 10GB file or infinite stream (gzip bomb).
- **Regex ReDoS:** Crafted HTML causing catastrophic backtracking in title extraction.
**Mitigation:**
- **Timeouts:** Aggressive timeouts (2s) for all external fetches.
- **Max Body Size:** Limit response reading to the first 1MB (sufficient for `<title>`).
- **Rate Limiting:** (Implemented in PR #27).

### 3. Cross-Site Scripting (XSS) in Extension
**Risk:** High
**Vector:** A resolved title contains `<script>alert(1)</script>`.
**Impact:** Code execution in the context of the victim's browser.
**Mitigation:**
- **Sanitization:** The Extension MUST treat all backend responses as untrusted text.
- **`innerText` vs `innerHTML`:** Always use `innerText` when updating the DOM.
- **CSP:** Enforce strict Content Security Policy in `manifest.json`.

### 4. Data Privacy
**Risk:** Medium
**Vector:** Logging user's browsing history via URLs sent to `/resolve`.
**Mitigation:**
- **No-Log Policy:** Do not log the specific URLs resolved in production logs.
- **Anonymization:** If metrics are needed, log only the domain, not the full path.

## Implementation Plan

### Phase 1: SSRF Hardening (Backend)
Create a `SafeTransport` for `http.Client`.

```go
// Proposed Logic
func isPrivateIP(ip net.IP) bool {
    // Check against 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, 127.0.0.0/8...
}
```

### Phase 2: Response Handling (Backend)
- Limit `io.LimitReader` to 1MB.
- Use `net/html` tokenizer instead of Regex for title extraction (safer & faster).

### Phase 3: Extension Hardening (Frontend)
- Audit `content.ts` to ensure no `innerHTML` usage.
- Review `manifest.json` permissions.

## Verification
- **SSRF Test Suite:** Attempt to resolve local/metadata IPs.
- **XSS Test Suite:** Mock backend returning malicious payloads.
