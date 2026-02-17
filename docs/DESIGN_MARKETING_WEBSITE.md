# Design: LinkLens Marketing Website

## Overview
A simple, high-impact landing page to showcase the value of LinkLens, provide a live demo, and offer installation instructions. The site will emphasize our core mission: **Transparency for the Web.**

## Objectives
- **Demonstrate Value:** Show users exactly how their browsing experience improves.
- **Provide Trust:** Highlight privacy-first architecture and open-source nature.
- **Drive Adoption:** Clear links to browser stores and setup guides.

## 1. Page Structure

### A. Hero Section
- **Headline:** Transparency for the Web.
- **Sub-headline:** Stop "clicking and hoping." LinkLens reveals what's behind opaque URLs instantly.
- **CTA:** "Install Extension" (Primary) and "View on GitHub" (Secondary).

### B. "The Lens in Action" (Before & After Examples)
A visual comparison section showing "Raw" vs. "LinkLens Enhanced" links.
- **Example 1 (YouTube):** 
  - *Before:* `https://youtu.be/dQw4w9WgXcQ`
  - *After:* `[Icon] Rick Astley - Never Gonna Give You Up (Official Music Video)`
- **Example 2 (GitHub):**
  - *Before:* `https://github.com/google/go`
  - *After:* `[Icon] google/go (The Go programming language ★ 120k)`
- **Example 3 (Short Link):**
  - *Before:* `https://bit.ly/3x86n7r`
  - *After:* `[Icon] Google`

### C. Live Interactive Demo
An input field where users can paste any URL.
- **Action:** Calls the LinkLens backend API.
- **Output:** Renders the link exactly as it would appear in the browser (with icon and rich tooltip).

### D. Key Pillars
- **Privacy First:** We resolve links, we don't track you. Settings stay in your browser.
- **Universal:** From YouTube to GitHub to bit.ly, we've got you covered.
- **Open Source:** Auditable code for your peace of mind.

### E. Footer
- Links to GitHub, Privacy Policy, and License.

## 2. Tech Stack
- **Framework:** Next.js (React) for a fast, SEO-friendly site.
- **Styling:** Tailwind CSS for a modern, responsive design.
- **Icons:** Lucide-react (matching the extension's clean look).
- **Hosting:** GitHub Pages (using GitHub Actions for automated deployment).

## 3. Implementation Plan
1.  **Scaffold:** Initialize Next.js project in a new `/website` directory.
2.  **Components:** Build reusable "LinkCard" and "DemoBox" components.
3.  **API Integration:** Connect the demo box to the production backend.
4.  **Content:** Write copy and design "Before & After" visuals.
5.  **Deployment:** Configure CI/CD to push to GitHub Pages.

## 4. Visual Style
- **Aesthetic:** Clean, minimalist, "Developer-friendly."
- **Color Palette:** Professional blues, grays, and high-contrast text.
- **Typography:** Inter or System Sans-serif.

## 5. Test Cases
- [ ] **Mobile Responsive:** Verify the demo and examples look great on phones.
- [ ] **Demo Functionality:** Ensure the live resolution works for diverse links.
- [ ] **Accessibility:** Screen reader tests for the demo output.
