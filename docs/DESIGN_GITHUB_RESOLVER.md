# Design: GitHub Platform Resolver

## Overview
To provide deeper context for developer-centric workflows, LinkLens will implement a specialized resolver for GitHub repository URLs. Instead of a generic page title, this resolver will surface repository metadata such as star counts, primary programming languages, and descriptions.

## Objectives
- **Rich Metadata:** Surface stars, forks, language, and repo description.
- **Performance:** Efficiently parse GitHub URLs and utilize the GitHub API.
- **Resilience:** Gracefully fall back to generic OpenGraph resolution if API limits are hit or the repo is private.
- **Scalability:** Support optional personal access tokens (PATs) for higher rate limits.

## 1. URL Detection & Parsing
The `GitHubResolver` will identify URLs matching the repository pattern:
- **Pattern:** `https://github.com/{owner}/{repo}`
- **Exclusions:** It should ignore non-repo pages like `/settings`, `/notifications`, `/explore`, etc.

## 2. Implementation Strategy

### A. Backend: GitHubResolver (Go)
The resolver will implement the standard `Resolver` interface.
- **API Client:** Use a dedicated HTTP client (with SSRF protection).
- **Endpoint:** `https://api.github.com/repos/{owner}/{repo}`
- **Headers:** 
  - `Accept: application/vnd.github.v3+json`
  - `Authorization: token ${GITHUB_TOKEN}` (if configured).

### B. Result Mapping
The resolver will populate the `Result` struct with rich data:
- `Title`: `{owner}/{repo}`
- `Description`: Repository description + metadata string (e.g., "★ 1.2k | TypeScript").
- `Platform`: `github`

## 3. Configuration
- `GITHUB_TOKEN`: Optional environment variable for an API token.
- `ENABLED_RESOLVERS`: Ensure `github` is included.

## 4. UI Integration
The frontend already supports platform-specific icons and rich tooltips. The `GitHubResolver` will trigger the `github` icon (Octocat) and show the repository details in the Shadow DOM tooltip.

## 5. Sequence of Resolution
1. **GitHubResolver:** Tries to fetch from API.
2. **Fallback:** If API returns 404 (private repo) or 403 (rate limited), return `nil` so the `OpenGraphResolver` can attempt a generic resolution (which might work for public repos even without API).

## 6. Test Cases
- [ ] **Public Repo:** `https://github.com/google/go` -> Should show stars and description.
- [ ] **Invalid Repo:** `https://github.com/not/a/real/repo` -> Should fail gracefully.
- [ ] **Non-Repo URL:** `https://github.com/settings` -> `CanHandle` should return false.
- [ ] **Rate Limit:** Verify fallback behavior when API is unavailable.
