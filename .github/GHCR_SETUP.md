# GitHub Container Registry (GHCR) Setup Guide

This guide explains how to configure GitHub Container Registry for Sentinel Mesh to publish public Docker images.

## Prerequisites

- Repository admin access
- GitHub account with package publishing permissions

## Step 1: Enable Package Publishing

### Configure Workflow Permissions

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Actions** → **General**
3. Scroll to **Workflow permissions**
4. Select **Read and write permissions**
5. Check ✅ **Allow GitHub Actions to create and approve pull requests**
6. Click **Save**

This allows GitHub Actions to push Docker images to GHCR.

## Step 2: Configure Package Visibility

After the first release workflow runs, packages will be created. To make them public:

### For Each Package

1. Go to your GitHub profile or organization
2. Click on **Packages** tab
3. Find the package (e.g., `sentinel-mesh/api`)
4. Click on the package name
5. Click **Package settings** (gear icon)
6. Scroll to **Danger Zone**
7. Click **Change visibility**
8. Select **Public**
9. Type the package name to confirm
10. Click **I understand, change package visibility**

### Packages to Make Public

You'll need to make these packages public after the first release:

- `sentinel-mesh/api`
- `sentinel-mesh/collector`
- `sentinel-mesh/processor`
- `sentinel-mesh/analyzer`
- `sentinel-mesh/alerting`
- `sentinel-mesh/ml-service`
- `sentinel-mesh/frontend`

## Step 3: Link Package to Repository

For better organization:

1. Open the package settings (from step 2)
2. Under **Connect repository**
3. Select `sentinel-mesh` from the dropdown
4. This links the package to your repository

## Step 4: Configure Package Permissions (Optional)

If you want to allow others to push images:

1. In package settings
2. Scroll to **Manage Actions access**
3. Click **Add Repository**
4. Select repositories that can publish
5. Set permission level (Read, Write, Admin)

## Step 5: Test Image Pulling

After making packages public, test that anyone can pull them:

```bash docs-drift:skip
# No authentication required for public images
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:latest

# Should pull successfully without login
```

## Verifying GHCR Configuration

### Check Package Visibility

```bash docs-drift:skip
# Try pulling without authentication
docker logout ghcr.io

# This should work if public
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:v0.1.0

# Should see output like:
# v0.1.0: Pulling from georg-nikola/sentinel-mesh/api
# ...
# Status: Downloaded newer image for ghcr.io/georg-nikola/sentinel-mesh/api:v0.1.0
```

### View on GitHub

Public packages appear at:
- `https://github.com/georg-nikola/sentinel-mesh/pkgs/container/sentinel-mesh%2Fapi`
- `https://github.com/georg-nikola/sentinel-mesh/pkgs/container/sentinel-mesh%2Fcollector`
- etc.

## Creating a Release

### Via GitHub UI

1. Go to **Releases** → **Create a new release**
2. Click **Choose a tag**
3. Create new tag (e.g., `v0.2.0`)
4. Fill in release title: `Sentinel Mesh v0.2.0`
5. Add release notes
6. Click **Publish release**

The release workflow will automatically:
- Build all Docker images
- Tag with version number
- Tag with `latest`
- Push to GHCR
- Test images
- Update release notes with image URLs

### Via GitHub CLI

```bash docs-drift:skip
# Create release
gh release create v0.2.0 \
  --title "Sentinel Mesh v0.2.0" \
  --notes "See CHANGELOG.md for details"

# The workflow will trigger automatically
```

### Via Git Tag

```bash docs-drift:skip
# Create and push tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# Create release from tag
gh release create v0.2.0 --generate-notes
```

## Manual Workflow Trigger

You can also manually trigger the release workflow:

1. Go to **Actions** → **Release**
2. Click **Run workflow**
3. Enter release tag (e.g., `v0.2.0`)
4. Click **Run workflow**

Or via CLI:
```bash docs-drift:skip
gh workflow run release.yml -f tag=v0.2.0
```

## Image Naming Convention

Images follow this pattern:
```
ghcr.io/georg-nikola/sentinel-mesh/<service>:<tag>
```

Examples:
- `ghcr.io/georg-nikola/sentinel-mesh/api:0.1.0`
- `ghcr.io/georg-nikola/sentinel-mesh/api:0.1`
- `ghcr.io/georg-nikola/sentinel-mesh/api:0`
- `ghcr.io/georg-nikola/sentinel-mesh/api:latest`
- `ghcr.io/georg-nikola/sentinel-mesh/api:main-abc123`

## Using Images in Helm

Update Helm values to use GHCR:

```yaml docs-drift:skip
# values.yaml
image:
  registry: ghcr.io
  repository: georg-nikola/sentinel-mesh
  tag: "0.1.0"
  pullPolicy: IfNotPresent
```

Install:
```bash docs-drift:skip
helm install sentinel-mesh ./deployments/helm/sentinel-mesh \
  --set image.registry=ghcr.io \
  --set image.repository=georg-nikola/sentinel-mesh \
  --set image.tag=0.1.0
```

## Security Best Practices

### 1. Use Specific Tags

```bash docs-drift:skip
# ❌ Avoid using 'latest' in production
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:latest

# ✅ Use specific versions
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:0.1.0
```

### 2. Verify Image Signatures (Future)

When image signing is enabled:
```bash docs-drift:skip
# Verify signature
cosign verify ghcr.io/georg-nikola/sentinel-mesh/api:0.1.0
```

### 3. Scan Images Before Use

```bash docs-drift:skip
# Scan with Trivy
trivy image ghcr.io/georg-nikola/sentinel-mesh/api:0.1.0
```

## Troubleshooting

### Error: "denied: installation not allowed to Create organization package"

**Solution**: Follow Step 1 to enable workflow permissions.

### Error: "denied: permission_denied: write_package"

**Solution**:
1. Check workflow permissions are set to "Read and write"
2. Verify GITHUB_TOKEN has package write permissions
3. Check if package exists and workflow has access

### Images Not Showing as Public

**Solution**: Manually change visibility (Step 2) after first push.

### Cannot Pull Public Image

**Symptoms**: `Error: unauthorized` when pulling public image

**Solution**:
1. Verify package is actually public
2. Try `docker logout ghcr.io` first
3. Check package name is correct (lowercase)

### Workflow Fails to Push

**Check**:
```bash docs-drift:skip
# View workflow logs
gh run view --log-failed

# Look for authentication errors
# Verify GITHUB_TOKEN permissions
```

## Automation Tips

### Auto-update Package Visibility

Unfortunately, GHCR doesn't support setting visibility via API during creation. You must:
1. Let first release run (creates private packages)
2. Manually make them public (one-time setup)
3. Future releases will use existing public packages

### Helm Repository Setup

Consider setting up GitHub Pages for Helm charts:

```bash docs-drift:skip
# Package Helm chart
helm package deployments/helm/sentinel-mesh

# Update index
helm repo index . --url https://georg-nikola.github.io/sentinel-mesh

# Commit to gh-pages branch
git add index.yaml sentinel-mesh-*.tgz
git commit -m "Release v0.1.0"
git push origin gh-pages
```

Then users can:
```bash docs-drift:skip
helm repo add sentinel-mesh https://georg-nikola.github.io/sentinel-mesh
helm install sentinel-mesh sentinel-mesh/sentinel-mesh
```

## Monitoring Releases

### Check Latest Release

```bash docs-drift:skip
gh release view --json tagName,publishedAt,url
```

### List All Releases

```bash docs-drift:skip
gh release list
```

### View Release Assets

```bash docs-drift:skip
gh release view v0.1.0
```

## Cleanup Old Images

GHCR has storage limits. To clean up old versions:

1. Go to package settings
2. Find **Package versions**
3. Delete old/unused versions
4. Keep latest and important milestones

Or automate with retention policy (future feature).

## Resources

- [GHCR Documentation](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Publishing Docker images](https://docs.github.com/en/actions/publishing-packages/publishing-docker-images)
- [Managing package access](https://docs.github.com/en/packages/learn-github-packages/configuring-a-packages-access-control-and-visibility)

---

**Note**: After initial setup, the release process is fully automated. Just create a new GitHub release, and images will be built and published automatically!
