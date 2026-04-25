# Design: Legal Policies (Privacy & Terms)

**Issue:** #42
**Status:** Approved

## Overview
This document outlines the privacy and terms of service implementation for the LinkLens platform, ensuring full transparency about our data practices. As a privacy-first extension, we must be explicit about our "no-tracking" guarantees.

## Goals
- Establish a clear, plain-English Privacy Policy and Terms of Service.
- Assure users that no Personally Identifiable Information (PII) or user tracking is logged or sold.
- Outline the usage terms of our public APIs to prevent automated abuse (rate limiting).

## Implementation Detail
We will add two static routes to the Next.js marketing website:
1.  `/privacy`: Detail data collection limits. We only process URLs passed to the resolver, and our backend caches the public `URL -> Title` mapping, never associating it with an IP or user ID.
2.  `/terms`: Standard terms of use, specifying that the service is "as-is" and prohibiting abusive scraping/botting.

Both pages will be linked in the footer of `website/src/app/page.tsx`.
