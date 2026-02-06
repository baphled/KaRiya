---
name: breaking-changes
description: Managing backwards compatibility, deprecation, and migration strategies
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide the handling of breaking changes including identification, documentation, deprecation strategies, and migration paths for users.

## When to use me

Use this skill when:
- About to change a public API
- Removing or renaming exported functions/types
- Changing function signatures
- Modifying configuration format
- Updating database schema

## What Constitutes a Breaking Change

### Definitely Breaking

| Change | Example | Impact |
|--------|---------|--------|
| Remove exported function | Delete `ProcessEvent()` | Compile error for users |
| Change function signature | Add required parameter | Compile error |
| Remove exported type | Delete `EventConfig` struct | Compile error |
| Change struct field type | `ID string` -> `ID int` | Compile error |
| Remove struct field | Delete `Config.Debug` | Compile error if accessed |
| Change interface | Add method to interface | Compile error for implementers |
| Change behaviour | Function now returns error on empty input | Runtime behaviour change |

### Not Breaking

| Change | Why |
|--------|-----|
| Add new exported function | Existing code unaffected |
| Add optional parameter | Use functional options pattern |
| Add struct field | Existing code compiles |
| Bug fix | Expected behaviour restored |
| Performance improvement | Same API, faster |
| Add new package | Doesn't affect existing imports |

## Decision Framework

```
Is this a breaking change?
    │
    ├─ Does it change exported API?
    │   ├─ No → Not breaking
    │   └─ Yes → Check compatibility
    │
    ├─ Will existing code compile?
    │   ├─ Yes → Check runtime
    │   └─ No → BREAKING
    │
    └─ Will existing code behave the same?
        ├─ Yes → Not breaking
        └─ No → BREAKING (behaviour change)
```

## Deprecation Process

### Step 1: Mark as Deprecated

```go
// ProcessEvent processes a single event.
//
// Deprecated: Use ProcessEvents instead, which handles batching.
// This function will be removed in v2.0.0.
func ProcessEvent(e *Event) error {
    return ProcessEvents([]*Event{e})
}
```

### Step 2: Add Compile-Time Warning (Go 1.21+)

```go
//go:deprecated Use ProcessEvents instead
func ProcessEvent(e *Event) error {
    return ProcessEvents([]*Event{e})
}
```

### Step 3: Log Runtime Warning

```go
func ProcessEvent(e *Event) error {
    log.Warn("ProcessEvent is deprecated, use ProcessEvents instead")
    return ProcessEvents([]*Event{e})
}
```

### Step 4: Document in CHANGELOG

```markdown
### Deprecated
- `ProcessEvent` - Use `ProcessEvents` instead (removal planned v2.0.0)
```

### Step 5: Remove in Next Major Version

```go
// v2.0.0 - Remove deprecated function
// Delete the function entirely
```

## Migration Guide Template

Create `docs/migrations/v1-to-v2.md`:

```markdown
# Migrating from v1.x to v2.0

## Overview

This guide covers breaking changes in v2.0 and how to update your code.

## Breaking Changes Summary

| Change | v1.x | v2.0 | Migration |
|--------|------|------|-----------|
| Event processing | `ProcessEvent(e)` | `ProcessEvents([]*Event{e})` | Wrap in slice |
| Config format | YAML | TOML | Convert config file |
| Error handling | Returns nil | Returns error | Add error handling |

## Detailed Migration Steps

### 1. Event Processing API

**Before (v1.x):**
```go
err := kariya.ProcessEvent(event)
```

**After (v2.0):**
```go
results := kariya.ProcessEvents([]*kariya.Event{event})
if len(results.Errors) > 0 {
    return results.Errors[0]
}
```

**Why changed:** Batching improves performance by 10x for bulk operations.

### 2. Configuration Format

**Before (v1.x) - config.yaml:**
```yaml
database:
  path: ~/.kariya/data.db
```

**After (v2.0) - config.toml:**
```toml
[database]
path = "~/.kariya/data.db"
```

**Migration script:**
```bash
# Convert YAML to TOML
kariya migrate-config --from yaml --to toml
```

## Compatibility Matrix

| Feature | v1.x | v2.0 | Notes |
|---------|------|------|-------|
| Go version | 1.21+ | 1.22+ | Minimum version increased |
| SQLite | 3.35+ | 3.40+ | New features required |
| Config format | YAML | TOML | Migration tool provided |

## Getting Help

If you encounter issues migrating:
1. Check the [FAQ](./migration-faq.md)
2. Open an issue with the `migration` label
3. Ask in discussions
```

## Avoiding Breaking Changes

### Use Functional Options

```go
// Instead of adding required parameters:

// BAD - Breaking
func NewClient(url string, timeout time.Duration) *Client

// GOOD - Not breaking
func NewClient(url string, opts ...ClientOption) *Client

type ClientOption func(*clientConfig)

func WithTimeout(d time.Duration) ClientOption {
    return func(c *clientConfig) {
        c.timeout = d
    }
}
```

### Use Interface Evolution

```go
// Instead of adding methods to interface:

// BAD - Breaking for implementers
type Handler interface {
    Handle(ctx context.Context) error
    HandleBatch(ctx context.Context, items []Item) error // Breaking!
}

// GOOD - Optional interface
type Handler interface {
    Handle(ctx context.Context) error
}

type BatchHandler interface {
    Handler
    HandleBatch(ctx context.Context, items []Item) error
}

// Check at runtime
if bh, ok := handler.(BatchHandler); ok {
    return bh.HandleBatch(ctx, items)
}
```

### Use Struct Embedding

```go
// Instead of changing struct:

// BAD - Might break if users rely on struct layout
type Config struct {
    Path    string
    Debug   bool
    Timeout time.Duration // New field - could break reflection
}

// GOOD - Embedded options
type Config struct {
    Path  string
    Debug bool
    opts  configOpts // Private, can change freely
}

func (c *Config) SetTimeout(d time.Duration) {
    c.opts.timeout = d
}
```

## Breaking Change Announcement

When a breaking change is unavoidable:

### In PR Description

```markdown
## Breaking Change

**What:** Removed `ProcessEvent` function
**Why:** Replaced with more efficient `ProcessEvents` batch API
**Migration:** Wrap single events in slice: `ProcessEvents([]*Event{e})`
**Timeline:** Deprecated in v1.5, removed in v2.0

### Migration Guide
See [v1-to-v2.md](docs/migrations/v1-to-v2.md)
```

### In Commit Message

```
feat!: replace ProcessEvent with ProcessEvents batch API

BREAKING CHANGE: ProcessEvent is removed.
Use ProcessEvents([]*Event{e}) for single events.

Migration guide: docs/migrations/v1-to-v2.md
```

### In Release Notes

```markdown
## Breaking Changes

### ProcessEvent Removed

The `ProcessEvent` function has been removed in favour of `ProcessEvents`.

**Before:**
```go
err := ProcessEvent(e)
```

**After:**
```go
results := ProcessEvents([]*Event{e})
```

See [migration guide](docs/migrations/v1-to-v2.md) for details.
```

## Checklist for Breaking Changes

Before introducing a breaking change:

- [ ] Is this change absolutely necessary?
- [ ] Can it be done in a backwards-compatible way?
- [ ] Is there a deprecation path?
- [ ] Is the migration guide written?
- [ ] Is the changelog updated?
- [ ] Is the version bump correct (MAJOR)?
- [ ] Are users warned in advance?
- [ ] Is there a migration tool if needed?

## Related Skills

- `release-management` - Version bumping for breaking changes
- `migration-strategies` - Database and data migrations
- `api-design` - Designing APIs that minimise breaking changes
- `documentation-writing` - Writing migration guides
