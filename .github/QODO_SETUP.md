# Qodo Merge Setup Instructions

This repository is configured to use [Qodo Merge](https://www.qodo.ai/products/qodo-merge/) (formerly Codium PR-Agent) for automated AI-powered code reviews on pull requests.

## Installation Steps

### 1. Install Qodo Merge GitHub App

1. Visit the [Qodo Merge GitHub App](https://github.com/apps/qodo-merge) page
2. Click "Install" or "Configure"
3. Select the repositories you want to enable:
   - Choose "Only select repositories"
   - Select `georg-nikola/sentinel-mesh`
4. Click "Install" to grant permissions

### 2. Configure API Key (Optional - for advanced features)

If you want to use advanced features or custom AI models:

1. Sign up at [Qodo.ai](https://www.qodo.ai/)
2. Get your API key from the dashboard
3. Add it as a repository secret:
   ```bash
   gh secret set QODO_API_KEY --body "your-api-key-here"
   ```

### 3. Verify Installation

After installation, Qodo Merge will automatically:

1. **Review Pull Requests**: When a PR is opened or updated, Qodo will:
   - Analyze code changes
   - Provide inline code suggestions
   - Check for security issues
   - Suggest performance improvements
   - Review test coverage
   - Generate PR descriptions

2. **Add PR Labels**: Automatically label PRs based on content:
   - `security-review-required`
   - `performance-impact`
   - `needs-tests`
   - etc.

3. **Update CHANGELOG**: Suggest CHANGELOG entries for significant changes

## Configuration

The Qodo configuration is stored in `.github/qodo.toml`. Current settings:

- ✅ Automatic PR review on creation
- ✅ Security vulnerability scanning
- ✅ Performance review
- ✅ Test coverage analysis
- ✅ Code suggestions (up to 4 per PR)
- ✅ PR description generation
- ✅ CHANGELOG updates
- ✅ Similar issue detection

## Usage

### Automatic Reviews

Qodo will automatically review PRs when:
- A new PR is opened
- Commits are pushed to an existing PR
- PR is marked ready for review (from draft)

### Manual Commands

You can trigger specific actions by commenting on a PR:

```
/review              - Trigger a full code review
/describe            - Generate/update PR description
/improve             - Get code improvement suggestions
/ask <question>      - Ask questions about the PR
/test                - Generate test suggestions
/update_changelog    - Update CHANGELOG.md
/similar_issues      - Find similar issues/PRs
```

### Example PR Review Output

When Qodo reviews a PR, you'll see:

1. **Overall Assessment**
   - Security score
   - Code quality score
   - Test coverage assessment

2. **Inline Comments**
   - Specific code suggestions
   - Security warnings
   - Performance tips
   - Best practice recommendations

3. **PR Description Enhancement**
   - Summary of changes
   - Technical walkthrough
   - Breaking changes notice
   - Related issues

## Customization

Edit `.github/qodo.toml` to customize:

- Review strictness
- Number of suggestions
- Auto-commit behavior
- Language-specific rules
- Custom review prompts

## Troubleshooting

### Qodo not commenting on PRs

1. Verify the app is installed: https://github.com/apps/qodo-merge
2. Check repository permissions include PR write access
3. Ensure PRs are not from forks (forks have limited access)

### Too many/few suggestions

Adjust in `.github/qodo.toml`:
```toml
[pr_code_suggestions]
num_code_suggestions = 4  # Increase or decrease
```

### Disable specific features

Set `enabled = false` for any unwanted feature in `.github/qodo.toml`

## Resources

- [Qodo Merge Documentation](https://qodo-merge-docs.qodo.ai/)
- [GitHub App Permissions](https://github.com/apps/qodo-merge/installations)
- [Configuration Options](https://qodo-merge-docs.qodo.ai/usage-guide/configuration_options/)

---

**Note**: The free tier of Qodo Merge provides limited reviews per month. For unlimited reviews, consider upgrading to a paid plan at [qodo.ai](https://www.qodo.ai/pricing/).
