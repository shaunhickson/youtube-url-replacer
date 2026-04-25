# Design: Terraform Infrastructure Configuration

**Issue:** #44
**Status:** Approved

## Overview
Currently, the GCP infrastructure for the LinkLens backend (Cloud Run, Firestore, Secret Manager) is managed manually or through gcloud scripts. This creates a risk of configuration drift and makes deployments less reproducible. We will adopt Terraform as our Infrastructure as Code (IaC) tool.

## Implementation Details
We will create a foundational `terraform/` directory containing the definitions for our entire stack:

1.  **Google Cloud Run (`cloud_run.tf`):** Defines the `youtube-replacer-backend` service, deploying the Go container. It exposes the service to the public internet (`allUsers`).
2.  **Secret Manager (`secrets.tf`):** Provisions the `YOUTUBE_API_KEY` secret used by the backend resolver.
3.  **Firestore (`firestore.tf`):** Provisions the native Firestore database used for OpenGraph caching.
4.  **IAM (`cloud_run.tf`):** Creates a dedicated Service Account (`youtube-replacer-sa`) for the Cloud Run instance. Binds roles allowing it to read from Secret Manager and write to Firestore, adhering to the principle of least privilege.

## Note on Concurrency
This PR establishes the `terraform/` base directory. It is designed to be cleanly merged alongside the `feat/45-monitoring` branch (which adds `monitoring.tf` to this same directory).
