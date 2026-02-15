# Release Setup Guide

## ✅ Implemented

All binary distribution infrastructure has been implemented:

1. **Version injection via ldflags** ✅
   - `internal/cli/config.go` - Variables for Version, Commit, Date
   - `internal/cli/root.go` - Uses FullVersion() for rich output
   - `Makefile` - Adds ldflags to build and install targets
   - Verified: `./bin/flashdoc --version` shows full version info

2. **GoReleaser configuration** ✅
   - `.goreleaser.yaml` - Complete release automation
   - Multi-platform builds (darwin/linux/windows, amd64/arm64)
   - Homebrew formula auto-generation
   - Changelog generation
   - SHA-256 checksums

3. **GitHub Actions CI/CD** ✅
   - `.github/workflows/ci.yml` - Tests on push/PR
   - `.github/workflows/release.yml` - Release on tag push
   - Pre-release test gate
   - GoReleaser integration

4. **License** ✅
   - `LICENSE` - Beerware license

## 📋 Manual Steps Required

### 1. Create Homebrew tap repository

```bash
# Create a new GitHub repository
# Name: homebrew-tap
# Owner: nicovandenhove
# Description: Homebrew formulae for flashdoc
# Public: Yes

# Initialize it with a Formula directory
mkdir -p Formula
git init
git add Formula/.gitkeep
git commit -m "Initial commit"
git remote add origin git@github.com:nicovandenhove/homebrew-tap.git
git push -u origin main
```

### 2. Create GitHub Personal Access Token

1. Go to: https://github.com/settings/tokens?type=beta
2. Click "Generate new token" (fine-grained)
3. Configure:
   - **Token name**: `homebrew-tap-flashdoc`
   - **Expiration**: 1 year (or custom)
   - **Repository access**: Only select repositories → `homebrew-tap`
   - **Permissions**:
     - Repository permissions:
       - Contents: Read and write
       - Metadata: Read-only (automatic)
4. Generate and copy the token

### 3. Add GitHub Secret

1. Go to: https://github.com/nicovandenhove/flashdoc/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_TOKEN`
4. Value: Paste the token from step 2
5. Click "Add secret"

## 🚀 Release Process

Once the manual steps are complete:

```bash
# 1. Ensure you're on main and it's clean
git status

# 2. Run tests locally
go test -v ./features/...

# 3. Create and push a version tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

GitHub Actions will automatically:
- ✅ Run all tests
- 🏗️ Build binaries for all platforms
- 📦 Create archives with LICENSE and README
- 🔐 Generate SHA-256 checksums
- 📝 Generate changelog from commits
- 🍺 Update Homebrew formula in tap repo
- 🎉 Create GitHub Release with all artifacts

## 📥 User Installation Methods

After the first release (v0.2.0):

```bash
# Option 1: go install
go install github.com/nicovandenhove/flashdoc/cmd/flashdoc@latest

# Option 2: Homebrew
brew install nicovandenhove/tap/flashdoc

# Option 3: Download binary
# Visit: https://github.com/nicovandenhove/flashdoc/releases
# Download for your platform and add to PATH
```

## 🧪 Testing Before First Release

Test GoReleaser locally (requires goreleaser to be installed):

```bash
# Install goreleaser
brew install goreleaser/tap/goreleaser
# or: go install github.com/goreleaser/goreleaser/v2@latest

# Validate configuration
goreleaser check

# Test release (doesn't publish)
goreleaser release --snapshot --clean

# Check dist/ directory for generated artifacts
ls -lh dist/
```

## 🔍 Verification Checklist

Before pushing the first release tag:

- [ ] Homebrew tap repository created
- [ ] `HOMEBREW_TAP_TOKEN` secret added to flashdoc repo
- [ ] `goreleaser check` passes
- [ ] All BDD tests pass: `go test -v ./features/...`
- [ ] Local build works: `make build && ./bin/flashdoc --version`
- [ ] Snapshot release works: `goreleaser release --snapshot --clean`

## 📚 References

- [GoReleaser Docs](https://goreleaser.com)
- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Homebrew Tap Guide](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
