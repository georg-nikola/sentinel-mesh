# Repository Configuration

This document describes the repository settings configured to enforce fork-based contributions.

## Current Configuration

### Branch Protection (Main Branch)

The `main` branch is protected with the following rules:

✅ **Configured via API**:
- **Required Pull Request Reviews**: 1 approval required
- **Dismiss Stale Reviews**: Enabled - new commits dismiss previous approvals
- **Required Status Checks**: All CI checks must pass
  - Test Go Services
  - Test ML Services
  - Test Frontend
  - Test Helm Charts
  - Security Scan
- **Strict Status Checks**: Branch must be up to date before merging
- **Required Linear History**: No merge commits allowed (rebase only)
- **Required Conversation Resolution**: All review threads must be resolved
- **Disallow Force Pushes**: Force pushes are blocked
- **Disallow Deletions**: Branch cannot be deleted
- **Allow Fork Syncing**: Contributors can sync their forks

⚠️ **Manual Configuration Needed**:
- **Enforce Admins**: Should be enabled via GitHub web UI
  - Go to: Settings → Branches → Branch protection rules → main
  - Check "Include administrators"
  - This ensures even repository owner must follow PR workflow

### Labels

Created labels for Dependabot and issue organization:
- `automated` - Automated updates and PRs
- `ci-cd` - CI/CD and GitHub Actions
- `dependencies` - Dependency updates
- `github-actions` - GitHub Actions workflows
- `go` - Go/Golang related
- `javascript` - JavaScript/Node.js related
- `frontend` - Frontend/Vue.js related
- `python` - Python related
- `ml-service` - ML service related
- `docker` - Docker/Containers related

### Contribution Workflow

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for complete contribution guidelines.

**Fork-Based Workflow**:
1. Contributors fork the repository
2. Create feature branches in their fork
3. Submit pull requests from fork to main repository
4. All PRs require review and passing CI before merge

**Why Fork-Based?**:
- External contributors automatically must fork (no write access)
- Keeps main repository clean
- Clear separation between development and official branches
- Standard open-source workflow

## Repository Settings Recommendations

The following settings should be configured in GitHub's web UI:

### General Settings
- [ ] Disable "Allow merge commits" (enforce rebase/squash only)
- [ ] Enable "Automatically delete head branches" after PR merge
- [ ] Disable "Allow rebase merging" (optional - based on preference)
- [ ] Enable "Allow squash merging"

### Pull Requests
- [ ] Enable "Require conversation resolution before merging"
- [ ] Enable "Automatically delete head branches"

### Branch Protection (Manual Step)
1. Go to: `Settings` → `Branches` → `Branch protection rules` → `main`
2. Enable "Include administrators"
3. Verify all other settings match configuration above

### Collaborators
- For any future collaborators, grant "Triage" or "Read" access only
- This prevents direct branch creation while allowing PR contributions
- Repository owner maintains "Admin" access

## Verification

To verify branch protection is working:

```bash docs-drift:skip
# Check branch protection status
gh api repos/georg-nikola/sentinel-mesh/branches/main/protection

# Attempt to push directly to main (should fail)
git push origin main
# Expected: Remote rejected - protected branch

# Verify labels exist
gh label list

# Verify Dependabot can use labels
# Check any Dependabot PR for proper labeling
```

## Maintenance

- Review branch protection rules quarterly
- Update this document when configuration changes
- Keep CONTRIBUTING.md in sync with actual requirements
- Monitor Dependabot PR labels to ensure they're working

---

**Last Updated**: 2025-10-26
**Configuration Level**: Automated + Manual steps required
**Status**: Active (requires manual enforce_admins enablement)
