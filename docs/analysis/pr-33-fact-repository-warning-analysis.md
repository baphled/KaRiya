---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# PR #33 Analysis: Fact Repository Warning Suppression

**PR**: https://github.com/baphled/KaRiya/pull/33  
**Title**: fix(cli): suppress expected fact repository warnings  
**Date**: 2026-01-08

## Problem Statement

When the fact repository is not configured (e.g., initialization fails or in certain test scenarios), the application generates **duplicate warnings**:

1. Service layer logs: `s.logger.Warn("Fact repository not configured")`
2. Service layer returns error: `fmt.Errorf("fact repository not configured")`
3. Calling code prints: `fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)`

This happens in **TWO locations**:
- `cmd/cli/main.go:307` - in `handleExtractFacts()`
- `internal/cli/importer/service.go:165` - in `ImportRows()`

## Proposed Fix in PR #33

The PR suppresses the warning by checking the error message string:

```go
if err.Error() != "fact repository not configured" {
    fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)
}
```

Applied in both locations.

## Assessment: ❌ NOT RECOMMENDED

### Problems with PR #33 Approach

#### 1. **String Matching is Fragile**
- Breaks if error message changes
- No compile-time safety
- Hard to maintain

#### 2. **Violates Separation of Concerns**
- Callers shouldn't need to know internal error messages
- Tight coupling between service and callers
- Makes refactoring difficult

#### 3. **Treats Symptom, Not Root Cause**
- The real issue: treating a configuration state as a runtime error
- Fact repository being unconfigured is **valid** - it's optional
- Should handle gracefully, not suppress warnings

#### 4. **Code Duplication**
- Same check logic in TWO places
- Will need to be repeated anywhere else SaveFact is called
- Maintenance burden

#### 5. **Inconsistent Error Handling**
- Other optional features (burst repository) have same pattern
- Would need similar suppression everywhere
- Doesn't scale

## Root Cause Analysis

### Why Does This Happen?

**Initialization Flow:**
```go
// cmd/cli/main.go:151-156
factRepo, err := career.NewSQLiteFactRepository(db)
if err != nil {
    fmt.Fprintf(errOut, "Warning: Failed to initialize fact repository: %v\n", err)
    // factRepo NOT set - remains nil
} else {
    svc.SetFactRepository(factRepo)
}
```

**Later Usage:**
```go
// Tries to save fact
if err := svc.SaveFact(ctx, &fact); err != nil {
    // err = "fact repository not configured" (expected)
    fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)  // Double warning!
}
```

### The Core Issue

**Fact repository is OPTIONAL**, not required. The application should treat it as such.

## Recommended Solutions

### Option A: Sentinel Error Pattern (BEST PRACTICE) ⭐

**Benefits:**
- Type-safe error checking
- Explicit intent
- Idiomatic Go
- Easy to test
- Works everywhere

**Implementation:**

```go
// internal/service/career/errors.go (NEW FILE)
package career

import "errors"

// Sentinel errors for optional features
var (
    ErrFactRepositoryNotConfigured  = errors.New("fact repository not configured")
    ErrBurstRepositoryNotConfigured = errors.New("burst repository not configured")
)
```

```go
// internal/service/career/service.go
func (s *Service) SaveFact(ctx context.Context, fact *domain.Fact) error {
    if s.factRepo == nil {
        // Don't log - this is expected state, not an error
        return ErrFactRepositoryNotConfigured
    }
    // ... rest of logic
}
```

```go
// cmd/cli/main.go
for _, fact := range facts {
    if err := svc.SaveFact(ctx, &fact); err != nil {
        // Silently skip if fact repository not configured (expected)
        if !errors.Is(err, careerservice.ErrFactRepositoryNotConfigured) {
            fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)
        }
    } else {
        factCount++
        // ...
    }
}
```

**Changes Required:**
1. Create `internal/service/career/errors.go` with sentinel errors
2. Update `service.go` to use `ErrFactRepositoryNotConfigured` (remove log)
3. Update callers to check with `errors.Is()`
4. Update tests to use `errors.Is()`

**Why This is Best:**
- ✅ Type-safe, compile-time checked
- ✅ Idiomatic Go pattern
- ✅ Self-documenting code
- ✅ Easy to test
- ✅ Scales to other optional features (burst repository)
- ✅ No string matching
- ✅ Can wrap errors properly
- ✅ Backward compatible (error message stays same for logging)

---

### Option B: Silent No-Op (SIMPLE BUT CONTROVERSIAL)

**Benefits:**
- Simplest implementation
- No caller changes needed
- Clear intent

**Drawbacks:**
- Hides fact repository unavailability
- Callers can't distinguish "not configured" from "success"
- May surprise users

**Implementation:**

```go
// internal/service/career/service.go
func (s *Service) SaveFact(ctx context.Context, fact *domain.Fact) error {
    if s.factRepo == nil {
        // Fact storage is optional - silently skip if not configured
        return nil  // <- Change: return nil instead of error
    }
    // ... rest of logic
}
```

**When to Use:**
- If fact storage is truly optional and users don't need to know
- If initialization failure is already logged clearly

**Why NOT Recommended:**
- Users might want to know if facts aren't being saved
- Debugging becomes harder (silent failures)
- Less explicit than Option A

---

### Option C: Check Before Calling (LEAST INVASIVE)

**Benefits:**
- Doesn't change service behavior
- Explicit at call site
- Easy to understand

**Drawbacks:**
- Exposes internal state
- Must remember to check everywhere
- Still duplicates logic

**Implementation:**

```go
// internal/service/career/service.go
// Add method to expose state
func (s *Service) HasFactRepository() bool {
    return s.factRepo != nil
}
```

```go
// cmd/cli/main.go & importer/service.go
if svc.HasFactRepository() {
    if err := svc.SaveFact(ctx, &fact); err != nil {
        fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)
    } else {
        factCount++
    }
}
// Else: silently skip
```

**Why NOT Recommended:**
- Leaks abstraction
- Callers need domain knowledge
- More verbose
- Still some duplication

---

## Recommendation: Option A (Sentinel Error)

Implement **Option A** for these reasons:

1. **Go Best Practice**: Sentinel errors are the idiomatic way to handle known error conditions
2. **Type Safety**: Compile-time checking via `errors.Is()`
3. **Scalability**: Works for burst repository and any future optional features
4. **Clarity**: Makes intent explicit in code
5. **Maintainability**: Easy to refactor, no string matching
6. **Testing**: Easy to test with `errors.Is()` in assertions

### Implementation Plan

**Phase 1: Create Sentinel Errors**
```bash
# Create new file
touch internal/service/career/errors.go
```

**Phase 2: Update Service Layer**
- Replace `fmt.Errorf("fact repository not configured")` with `ErrFactRepositoryNotConfigured`
- Remove `s.logger.Warn()` call (not an error condition)
- Same for burst repository

**Phase 3: Update Callers**
- `cmd/cli/main.go`: Use `errors.Is()` in `handleExtractFacts()`
- `internal/cli/importer/service.go`: Use `errors.Is()` in `ImportRows()`

**Phase 4: Update Tests**
- Replace `Expect(err.Error()).To(ContainSubstring("fact repository not configured"))` 
- With `Expect(errors.Is(err, careerservice.ErrFactRepositoryNotConfigured)).To(BeTrue())`

**Phase 5: Update Documentation**
- Document that fact/burst repositories are optional
- Add to error handling guide

### Breaking Changes

**None** - This is fully backward compatible:
- Error message stays the same for logging/display
- Existing tests can be updated incrementally
- Old string-based checks still work (but should migrate)

## Alternative: Accept PR #33 with Improvements

If we want a **quick fix** now and refactor later:

### Improvements to PR #33:

1. **Add constant for error message:**
```go
const errFactRepoNotConfigured = "fact repository not configured"

if err.Error() != errFactRepoNotConfigured {
    // ...
}
```

2. **Extract to helper function:**
```go
func isFactRepoNotConfigured(err error) bool {
    return err != nil && err.Error() == "fact repository not configured"
}

if !isFactRepoNotConfigured(err) {
    fmt.Fprintf(errOut, "Warning: Failed to save fact: %v\n", err)
}
```

3. **Add TODO comments:**
```go
// TODO: Replace with sentinel error pattern (errors.Is)
// See: docs/analysis/pr-33-fact-repository-warning-analysis.md
```

**Verdict:** ⚠️ **Accept with improvements as interim solution, but plan migration to Option A**

---

## Conclusion

**Short-term**: PR #33 can be accepted with improvements (constants, helper function, TODOs)

**Long-term**: Migrate to sentinel error pattern (Option A) for proper architecture

**Immediate Action**: 
1. Improve PR #33 with constants and helper
2. Create issue to track sentinel error migration
3. Document in AGENTS.md that this is technical debt

---

**References:**
- Go Error Handling Best Practices: https://go.dev/blog/error-handling-and-go
- Effective Go - Errors: https://go.dev/doc/effective_go#errors
- Go Sentinel Errors: https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
