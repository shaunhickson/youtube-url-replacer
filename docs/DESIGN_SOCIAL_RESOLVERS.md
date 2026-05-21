# Design Document: Social Media & Custom Resolvers (Phase 3)

## 1. Product Alignment
*   **Roadmap Check:** Addresses Phase 3, "Deep Integrations & Enterprise Value" -> Social Media context and additional custom resolvers.
*   **Scope:** Implement five new backend resolvers:
    1.  **Twitter (X)**
    2.  **LinkedIn**
    3.  **Reddit**
    4.  **Wikipedia**
    5.  **Spotify**
*   **User Value:** Provides precise, high-quality link previews for the most commonly shared platforms on the web, bypassing the strict anti-bot measures that cause generic OpenGraph resolvers to fail.

## 2. Architecture & Patterns
Each resolver will implement the `Resolver` interface (`Name`, `CanHandle`, `Resolve`) and will be registered in `backend/main.go` through the `ResolverManager`.

### Resolver Specifics
*   **Reddit & Wikipedia:** Free, unauthenticated APIs. 
    *   Reddit: Appends `.json` to the URL.
    *   Wikipedia: Uses the `en.wikipedia.org/api/rest_v1/page/summary/{title}` endpoint.
*   **Twitter (X), LinkedIn, & Spotify:** Require authenticated API calls to function reliably.

## 3. Configuration & Authentication
Because several platforms heavily restrict unauthenticated scraping, LinkLens requires API tokens to be provided in the backend `.env` file. If a required token is missing, that specific resolver will not be registered, and LinkLens will gracefully fall back to the generic `OpenGraphResolver`.

### Setup Instructions for Server Admins

#### 1. Twitter (X)
**Env Var:** `TWITTER_BEARER_TOKEN`
**How to obtain:**
1. Go to the [Twitter Developer Portal](https://developer.twitter.com/).
2. Create a Free or Basic Project & App.
3. Generate a "Bearer Token" under the "Keys and tokens" tab.
4. Add to `.env`: `TWITTER_BEARER_TOKEN=your_bearer_token_here`

#### 2. LinkedIn
**Env Var:** `LINKEDIN_ACCESS_TOKEN`
**How to obtain:**
1. Go to the [LinkedIn Developer Portal](https://developer.linkedin.com/).
2. Create an App and request the "Sign In with LinkedIn" or "Share on LinkedIn" product to access the `r_liteprofile` and `r_organization_social` scopes.
3. Generate an OAuth 2.0 Access Token (Note: LinkedIn tokens expire and require a refresh strategy, or an Enterprise API arrangement).
4. Add to `.env`: `LINKEDIN_ACCESS_TOKEN=your_access_token_here`

#### 3. Spotify
**Env Vars:** `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET`
**How to obtain:**
1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/).
2. Create an App.
3. Copy the Client ID and Client Secret. LinkLens will handle generating the temporary access tokens using the Client Credentials Flow.
4. Add to `.env`: 
   `SPOTIFY_CLIENT_ID=your_client_id_here`
   `SPOTIFY_CLIENT_SECRET=your_client_secret_here`

## 4. Frontend Requirements
The frontend extension will be updated to display platform-specific icons.
*   `extension/src/utils/ui.ts` will receive new SVG mappings for: `twitter`, `linkedin`, `reddit`, `wikipedia`, and `spotify`.

## 5. Security & Safety
*   API keys will be loaded from the environment securely and never exposed to the frontend.
*   Resolvers using external APIs will utilize `backend/transport/SafeHttpClient` to prevent SSRF vulnerabilities when following redirects.
*   Timeouts will be strictly enforced on all external API requests (e.g., 2 seconds) to ensure LinkLens remains fast.

## Verification Checklist
- [x] Specifies all five targeted platforms.
- [x] Includes clear instructions for generating and injecting required API keys.
- [x] Details frontend SVG updates.
- [x] Ensures fallback behavior if API keys are not present.
