# Design: Rename Primary Branch to 'main'

## Overview
As part of our modernization efforts, we will rename the primary development branch from `master` to `main`. This involves technical git operations, CI/CD configuration updates, and documentation synchronization.

## Objectives
- **Standardization:** Adopt the industry-standard `main` branch name.
- **Continuity:** Ensure zero downtime for CI/CD pipelines.
- **Consistency:** Update all references in documentation.

## 1. Execution Steps

### Phase A: Local & Remote Migration
1. Rename local branch: `git branch -m master main`
2. Push new branch to origin: `git push -u origin main`

### Phase B: GitHub Configuration (Manual)
1. Set `main` as the default branch in GitHub Repository Settings.
2. Update branch protection rules from `master` to `main`.
3. (Optional) Delete the `master` branch on origin after verification.

### Phase C: CI/CD & Code Updates
Update the following files to replace `master` with `main`:
- **`.github/workflows/ci.yml`**:
    - `on.push.branches`
    - `on.pull_request.branches`
    - `deploy-backend.if` condition
- **`README.md`**: Update any setup instructions.
- **`docs/GITHUB_STRATEGY.md`**: Update the "Branching Strategy" section.

## 2. Implementation Plan
This design will be merged into `master`. 
The implementation PR will then be created from a new branch (or we can use the `main` branch directly once renamed) to perform the file updates.

## 3. Test Cases
- [ ] CI/CD triggers correctly on a PR targeting `main`.
- [ ] Deployment to Cloud Run triggers correctly on a push to `main`.
- [ ] Documentation accurately reflects the new branch name.
