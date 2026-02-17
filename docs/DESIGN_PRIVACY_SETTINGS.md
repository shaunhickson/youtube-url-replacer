# Design: Privacy Settings (Allow/Block Lists)

## Overview
As LinkLens becomes a universal link resolver, it is critical that users have control over what data is sent to the backend. This design introduces domain-level filtering to allow users to exclude sensitive sites or only enable the extension on trusted domains.

## Objectives
- **User Control:** Provide an easy-to-use interface for managing domain filters.
- **Privacy by Design:** Filtering happens locally in the browser before any network request is made.
- **Flexibility:** Support both "Blocklist" (exclude specific sites) and "Allowlist" (only include specific sites) modes.

## Proposed Architecture

### 1. Storage Schema (chrome.storage.local)
We will store settings in the following structure:
```json
{
  "enabled": true,
  "filterMode": "blocklist", // "blocklist" or "allowlist"
  "domainList": ["internal.corp.com", "bank.com"],
  "matchSubdomains": true
}
```

### 2. Filtering Logic (content.ts)
The `LinkLensOptimizer` will be updated to check the current domain and the target link domain against the user's filters.
- **Current Page Check:** If the site the user is currently browsing is blocked, the extension does nothing.
- **Link-by-Link Check:** For every candidate link:
    - If mode is `blocklist`: Ignore if domain is in `domainList`.
    - If mode is `allowlist`: Ignore if domain is NOT in `domainList`.

### 3. User Interface

#### A. Enhanced Popup
- **Quick Toggle:** "Disable on this site" (adds current domain to blocklist).
- **Status Indicator:** Shows if the current site is being filtered.
- **Link to Options:** "Manage Privacy Settings".

#### B. Options Page
A dedicated extension options page (`options.html`) to:
- Toggle between Blocklist and Allowlist modes.
- Manage the list of domains (add/remove).
- Toggle "Match Subdomains".

## Implementation Plan

### Extension
1.  **Storage Utilities:** Create helpers to read/write settings and perform domain matching.
2.  **Options Page:** Implement a simple React-based UI for managing the domain list.
3.  **Content Script:**
    - Listen for storage changes.
    - Perform filtering check before `enqueueLinks`.
4.  **Popup:** Add the "Disable on this site" quick action.

## Security Considerations
- **No Syncing:** Privacy settings remain local to the browser instance by default (using `chrome.storage.local`).
- **Input Validation:** Sanitize domain inputs to prevent regex-based performance issues or invalid patterns.

## Test Cases
- [ ] **Blocklist Mode:** Add `example.com` to blocklist; verify no links on `example.com` are resolved.
- [ ] **Allowlist Mode:** Add `youtube.com` to allowlist; verify only YouTube links are resolved on other pages.
- [ ] **Subdomain Matching:** Verify `corp.com` blocks `sub.corp.com` if enabled.
- [ ] **Global Toggle:** Turn off extension; verify all DOM scanning stops.
