# Auto-Rebase Workflow Setup

This document explains how to configure the auto-rebase workflow to preserve PR approvals and trigger CI workflows.

## The Problem

GitHub's security design prevents `GITHUB_TOKEN` from triggering workflows when force-pushing (to avoid recursive workflow runs). This causes two issues:

1. **CI workflows don't run** on rebased commits (status checks stuck in "Expected" state)
2. **PR approvals are lost** if we close/reopen PRs to trigger workflows

## The Solution

Use a **Personal Access Token (PAT)** instead of `GITHUB_TOKEN` for the auto-rebase workflow. PATs:
- ✅ Trigger CI workflows on force-push
- ✅ Preserve PR approvals
- ✅ Don't require closing/reopening PRs

## Setup Instructions

### Step 1: Create a Personal Access Token

1. Go to GitHub: **Settings** → **Developer settings** → **Personal access tokens** → **Fine-grained tokens**
2. Click **Generate new token**
3. Configure the token:
   - **Name**: `Auto-Rebase Workflow Token`
   - **Expiration**: 90 days (or custom)
   - **Repository access**: Only select repositories → `sentinel-mesh`
   - **Permissions**:
     - **Contents**: Read and write
     - **Pull requests**: Read and write
     - **Workflows**: Read and write (to trigger CI)

4. Click **Generate token**
5. **Copy the token** (you won't see it again!)

### Step 2: Add Token as Repository Secret

1. Go to your repository: **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Configure:
   - **Name**: `AUTO_REBASE_PAT`
   - **Secret**: Paste your token
4. Click **Add secret**

### Step 3: Verify Setup

The auto-rebase workflow will automatically use `AUTO_REBASE_PAT` if available, falling back to `GITHUB_TOKEN` if not.

To verify it's working:
1. Create a test PR
2. Approve it
3. Push a commit to main (triggers auto-rebase)
4. Check that:
   - PR is rebased
   - CI workflows run automatically
   - **Approval is preserved** ✅

## Without PAT Configuration

If `AUTO_REBASE_PAT` is not configured, the workflow uses `GITHUB_TOKEN` which means:
- ⚠️ PRs are rebased but CI workflows won't automatically trigger
- ⚠️ Manual intervention needed to trigger CI (close/reopen PR or push empty commit)

## Token Security

- Fine-grained tokens are scoped to specific repositories and permissions
- Set reasonable expiration dates (90 days recommended)
- Rotate tokens regularly
- Revoke tokens if compromised: **Settings** → **Developer settings** → **Personal access tokens** → **Revoke**

## Troubleshooting

### CI workflows still not triggering

Check that the PAT has the **Workflows** permission (Read and write).

### "Resource not accessible by integration" error

The PAT may have expired or been revoked. Create a new token and update the secret.

### Approvals still being cleared

Ensure you're using a PAT, not `GITHUB_TOKEN`. Check workflow logs to confirm which token is being used.

## Alternative: GitHub App

For organization-level deployments, consider creating a GitHub App instead of using a PAT:
- Better security (automatic token rotation)
- More granular permissions
- Works across multiple repositories
- Requires more complex setup

See: https://docs.github.com/en/apps/creating-github-apps

## References

- [GitHub security design for GITHUB_TOKEN](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow)
- [Creating fine-grained PATs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)
- [Encrypted secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)

<!-- AUTO_REBASE_PAT configured successfully on 2025-10-27 12:45:56 UTC -->
