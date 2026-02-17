# Design: Extension Options Enhancements

## Overview
This design expands the recently created Options Page to include advanced configuration for power users: self-hosting capabilities and visual theme controls.

## Objectives
- **Self-Hosting:** Allow users to specify a custom backend API URL.
- **Visual Customization:** Provide a toggle for Light/Dark mode for rich tooltips.
- **Persistence:** Ensure all settings are stored safely and applied instantly.

## 1. Data Schema Updates (chrome.storage.local)
The `Settings` interface in `settings.ts` will be expanded:
```typescript
interface Settings {
    // ... existing fields (enabled, domainList, etc.)
    apiUrl: string; // Default: https://youtube-replacer-backend-542312799814.us-east1.run.app
    theme: 'light' | 'dark' | 'system'; // Default: 'system'
}
```

## 2. Implementation Strategy

### A. Self-Hosting (API URL)
- **UI:** A text input in the Options page with a "Reset to Default" button.
- **Validation:** Ensure the URL is valid and uses `https` (unless it's `localhost` for development).
- **Application:** The `LinkLensOptimizer` in `content.ts` will fetch this URL from storage for every batch request.

### B. Theme Controls (Tooltip Styling)
- **UI:** A radio group or dropdown in the Options page (Light, Dark, System).
- **Application:** 
    - The `UIManager` in `ui.ts` will receive the theme setting.
    - The Shadow DOM CSS will be updated to support a `.dark` class or use CSS variables that change based on the setting.
    - If 'System' is selected, it will use `@media (prefers-color-scheme: dark)`.

## 3. UI/UX Plan

### Options Page Additions:
- **Connection Section:**
    - Input for "Backend API URL".
    - "Test Connection" button (verifies the `/health` endpoint of the custom URL).
- **Appearance Section:**
    - Selection for Tooltip Theme.

## 4. Security Considerations
- **URL Validation:** Prevent injection of malicious URLs.
- **Content Security Policy (CSP):** Ensure the extension can still connect to the user-defined domain (may require updates to `host_permissions`).

## 5. Test Cases
- [ ] **Custom URL:** Point to `http://localhost:8080` and verify resolution works (with local backend running).
- [ ] **Theme Switch:** Change to "Dark" and verify tooltip colors update immediately.
- [ ] **Reset:** Verify "Reset to Default" restores the original production backend URL.
