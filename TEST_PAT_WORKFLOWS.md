# AUTO_REBASE_PAT Workflows Permission Test

This commit tests if AUTO_REBASE_PAT has the "Workflows" permission.

**Test time:** 2025-10-27 21:05:00 UTC

## Expected Behavior

1. ✅ This commit triggers auto-rebase workflow
2. ✅ Auto-rebase rebases all open PRs
3. ✅ **CI Pipeline workflows trigger automatically on rebased PRs**
4. ✅ No manual intervention needed

## How to Verify

Check workflow runs after this commit:
- Auto-rebase: https://github.com/georg-nikola/sentinel-mesh/actions/workflows/auto-rebase-prs.yml
- CI Pipeline on PRs: https://github.com/georg-nikola/sentinel-mesh/actions/workflows/ci.yml

If CI runs on PRs automatically, the PAT is working correctly! 🎉

## What This Tests

The AUTO_REBASE_PAT must have:
- ✓ Contents: Read and write
- ✓ Pull requests: Read and write
- ✓ **Workflows: Read and write** ← Critical for triggering CI

If the PAT has all three permissions, force-pushing to PR branches will trigger CI workflows automatically.
