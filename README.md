# LinkLens

**Make the web transparent by default.**

LinkLens is a browser extension and backend service that automatically replaces opaque URLs (like raw YouTube links or shortened `bit.ly` links) with their human-readable titles. Stop "clicking and hoping" and start browsing with context.

## Features

- **Universal Link Resolution:** Automatically detects raw URLs in web pages and fetches their titles using OpenGraph metadata.
- **YouTube Specialization:** Dedicated support for YouTube links (watch, shorts, live, mobile) using the official YouTube API.
- **Privacy Focused:** Only sends "raw" links to the backend. No user tracking, no cookies, and SSRF-protected resolution.
- **High Performance:** Backend caching with Firestore/Redis support and low-latency Go implementation.
- **Modular Architecture:** Easily add support for new platforms (GitHub, Jira, etc.) via a pluggable resolver system.

## Project Structure

- `/backend`: Go-based resolution engine.
  - `/resolvers`: Pluggable logic for different platforms.
  - `/transport`: Security-hardened HTTP client (SSRF protection).
  - `/middleware`: Rate limiting and logging.
- `/extension`: React & TypeScript browser extension (Vite-powered).
  - `/src/content.ts`: DOM scanning and link replacement logic.
- `/docs`: Detailed design documents and roadmaps.

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- A YouTube Data API v3 Key (optional, for YouTube resolution)

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Set up environment variables:
   ```bash
   export YOUTUBE_API_KEY=your_key_here
   ```
3. Run the server:
   ```bash
   go run .
   ```

### Extension Setup

1. Navigate to the extension directory:
   ```bash
   cd extension
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Build the extension:
   ```bash
   npm run build
   ```
4. Load the `dist` folder into your browser (Chrome/Edge/Brave) via "Load unpacked" in the Extensions settings.

## Development

The project includes a `Makefile` with common tasks:

- `make build-backend`: Build the Go binary.
- `make build-extension`: Build the React extension.
- `make test-backend`: Run all Go tests.
- `make test-extension`: Run vitest for the extension.

## Contributing

We follow a **Design-First** workflow. Every major feature must have a design document in the `docs/` folder before implementation begins. See `docs/GITHUB_STRATEGY.md` for our collaboration guidelines.

## License

MIT
