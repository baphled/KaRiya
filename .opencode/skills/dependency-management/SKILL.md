---
name: dependency-management
description: Manage Go modules safely - version constraints, updating dependencies, security patches, minimal dependencies
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide safe management of Go module dependencies - adding, updating, auditing, and minimising dependencies.

## When to use me

- Adding new dependencies
- Updating existing dependencies
- Auditing for vulnerabilities
- Cleaning up unused dependencies
- Resolving version conflicts

## Core Principles

1. **Minimal dependencies** - Every dependency is a liability
2. **Explicit versions** - Know exactly what you're using
3. **Regular updates** - Stay current, especially security patches
4. **Audit before adding** - Evaluate quality and maintenance

## Go Modules Basics

### Key Files

```
go.mod   - Direct dependencies and Go version
go.sum   - Cryptographic checksums for all dependencies
```

### Essential Commands

```bash
# Initialise module
go mod init github.com/baphled/kariya

# Add dependency (automatically on import)
go get github.com/example/package

# Add specific version
go get github.com/example/package@v1.2.3

# Update to latest
go get -u github.com/example/package

# Update all dependencies
go get -u ./...

# Update patch versions only (safer)
go get -u=patch ./...

# Remove unused dependencies
go mod tidy

# Verify checksums
go mod verify

# Show dependency graph
go mod graph

# Show why dependency is needed
go mod why github.com/example/package
```

## Adding Dependencies

### Before Adding - Evaluate

```markdown
## Dependency Evaluation: [package name]

### Need Assessment
- What problem does this solve?
- Can we solve it with stdlib?
- Can we write it ourselves reasonably?

### Quality Assessment
- [ ] Actively maintained (commits in last 6 months)
- [ ] Adequate test coverage
- [ ] Clear documentation
- [ ] Stable API (v1+)
- [ ] Reasonable dependency tree (not pulling in the world)
- [ ] No known security issues
- [ ] Compatible licence

### Risk Assessment
- What happens if abandoned?
- How hard to replace?
- How much of it do we use?
```

### Adding Safely

```bash
# 1. Add the dependency
go get github.com/example/package@v1.2.3

# 2. Check what it pulled in
go mod graph | grep example/package

# 3. Verify it builds
go build ./...

# 4. Run tests
go test ./...

# 5. Tidy up
go mod tidy

# 6. Commit go.mod and go.sum together
git add go.mod go.sum
git commit -m "deps: add github.com/example/package v1.2.3"
```

## Updating Dependencies

### Routine Updates

```bash
# 1. Check for updates
go list -m -u all

# 2. Update patch versions (safest)
go get -u=patch ./...

# 3. Test thoroughly
go test ./...

# 4. Tidy
go mod tidy

# 5. Commit
git add go.mod go.sum
git commit -m "deps: update patch versions"
```

### Major Version Updates

```bash
# Major versions are different module paths
# v1: github.com/example/package
# v2: github.com/example/package/v2

# Update import paths in code
sed -i 's|example/package"|example/package/v2"|g' **/*.go

# Get new version
go get github.com/example/package/v2

# Remove old version
go mod tidy

# Test thoroughly - major versions may have breaking changes
go test ./...
```

### Security Updates

```bash
# Check for vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Update vulnerable dependency
go get github.com/vulnerable/package@v1.2.4  # Patched version

# Verify fix
govulncheck ./...
```

## Version Constraints

### Version Syntax

```go
// go.mod
require (
    // Specific version
    github.com/example/a v1.2.3
    
    // Latest matching
    github.com/example/b v1.2.0
    
    // Pseudo-version (commit)
    github.com/example/c v0.0.0-20230101120000-abcdef123456
)
```

### Replace Directive

```go
// go.mod

// Use fork instead of original
replace github.com/original/package => github.com/fork/package v1.0.0

// Use local version (development)
replace github.com/example/package => ../local/package

// Pin to specific version despite transitive requirements
replace github.com/example/package => github.com/example/package v1.2.3
```

### Exclude Directive

```go
// go.mod

// Exclude known bad version
exclude github.com/example/package v1.2.0
```

## Dependency Hygiene

### Regular Maintenance

```bash
# Weekly/monthly routine

# 1. Remove unused
go mod tidy

# 2. Check for vulnerabilities
govulncheck ./...

# 3. Check for updates
go list -m -u all | grep -v '^\s*$'

# 4. Update patch versions
go get -u=patch ./...

# 5. Test
go test ./...

# 6. Commit if changes
git add go.mod go.sum
git commit -m "deps: routine maintenance"
```

### Audit Dependencies

```bash
# What do we depend on?
go list -m all

# Why do we need this?
go mod why github.com/example/package

# What does this depend on?
go mod graph | grep "^github.com/example/package"

# Licence check
go-licenses check ./...
```

## Vendoring (Optional)

```bash
# Create vendor directory
go mod vendor

# Build using vendor
go build -mod=vendor ./...

# Update vendor after go.mod changes
go mod vendor
```

### When to Vendor

**Vendor when:**
- Reproducible builds critical
- Building in isolated environment
- Dependencies might disappear

**Don't vendor when:**
- Regular application development
- Dependencies are stable
- CI/CD has internet access

## Troubleshooting

### "module declares its path as X but was required as Y"

```bash
# Module path mismatch - use replace
replace github.com/wrong/path => github.com/correct/path v1.0.0
```

### Version conflicts

```bash
# Find why version is required
go mod why -m github.com/example/package

# Force specific version
go get github.com/example/package@v1.2.3
```

### Checksum mismatch

```bash
# Clear cache and re-download
go clean -modcache
go mod download
```

### "cannot find module providing package"

```bash
# Module not in cache
go mod download

# Or specific module
go get github.com/example/package
```

## Best Practices

### DO

- Lock to specific versions in production
- Update regularly (at least security patches)
- Audit dependencies before adding
- Use `go mod tidy` before committing
- Commit `go.mod` and `go.sum` together
- Test after any dependency change

### DON'T

- Use `latest` in production
- Ignore security advisories
- Add dependencies without evaluation
- Forget to test after updates
- Commit `go.mod` without `go.sum`

## Related Skills

- `security` - Vulnerability management
- `go-expert` - Go language patterns
- `devops` - CI/CD integration
