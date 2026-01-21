# Bug 007: E2E Onboarding Tests Pollute User's Config File

**Status**: Fixed  
**Severity**: 🟠 High  
**Created**: 2026-01-21  
**Updated**: 2026-01-21  
**Fixed**: 2026-01-21  

---

## Bug Summary

E2E tests using `SetupWithOnboarding()` and `CompleteOnboarding()` write test data to the user's real `~/.kariya/config.yaml`, corrupting their profile settings.

---

## Affected Components

- [x] `internal/testutil/e2e/helpers.go` - `SetupWithOnboarding()` doesn't isolate config path
- [x] `internal/cli/app/app.go:285` - Saves to real config on onboarding completion
- [x] `internal/config/config.go` - `SaveConfig()` always uses real path

**Related Intents/Workflows**:
- Onboarding workflow
- E2E test infrastructure

---

## Reproduction Steps

1. Have a valid `~/.kariya/config.yaml` with your real profile data
2. Run the onboarding E2E tests: `go test ./internal/testutil/e2e/... -run "onboarding"`
3. Check your config file: `cat ~/.kariya/config.yaml`
4. Observe test values have overwritten your real profile

**Consistency**: Always

**Environment**:
- OS: Linux
- Terminal: Any
- Go Version: 1.25.6
- KaRiya Version: next branch

---

## Expected Behavior

Tests should be fully isolated and never modify files outside of temporary test directories. The user's `~/.kariya/config.yaml` should remain unchanged after running tests.

---

## Actual Behavior

The user's config file is overwritten with test values:

**Evidence**:
```yaml
# ~/.kariya/config.yaml after running tests
profile:
    name: Test User
    email: test@example.com
```

---

## Investigation Log

### [2026-01-21] - Initial Investigation
- User reported profile name and email set to test values
- Searched for hardcoded test values in codebase
- Found `config.go:DefaultConfig()` correctly sets empty strings
- Traced issue to E2E test flow

### [2026-01-21] - Code Review
- Reviewed `internal/testutil/e2e/helpers.go:149` - `SetupWithOnboarding()`
- Found database isolation uses `t.TempDir()` correctly
- Found config path is NOT isolated
- Reviewed `internal/cli/app/app.go:285` - saves config on wizard completion
- Reviewed `internal/testutil/e2e/onboarding_e2e_test.go:124` - uses test values

### [2026-01-21] - Test Review
- Existing tests: Database isolation works correctly
- Gap identified: Config file path not isolated in test setup

---

## Root Cause

**Status**: Identified

**Cause**:
`SetupWithOnboarding()` creates a temporary database but does not override the config file path. When the onboarding wizard completes, it calls `config.SaveConfig()` which writes to the real `~/.kariya/config.yaml`.

**Technical Details**:
- File: `internal/testutil/e2e/helpers.go`
- Function: `SetupWithOnboarding()`
- Line: `149-206`
- Reason: Only database path is isolated; config path uses default `GetConfigPath()`

**Code Flow**:
1. `SetupWithOnboarding()` at `helpers.go:149` - creates temp DB, no config isolation
2. `CompleteOnboarding("Test User", "test@example.com")` at `onboarding_e2e_test.go:124`
3. Wizard completion detected at `app.go:277`
4. `config.SaveConfig(m.appConfig)` called at `app.go:285`
5. `SaveConfig()` uses `GetConfigPath()` returning `~/.kariya/config.yaml`

---

## Fix Strategy

**Approach**: Add config path override mechanism for tests

### Option A: Config Path Override (Recommended)

Add a function to override the config path for testing purposes.

**Pros**:
- Minimal change to existing code
- Clear API for test isolation
- Follows existing pattern (database uses temp path)

**Cons**:
- Requires cleanup in test teardown

**Files to Change**:
- [ ] `internal/config/config.go` - Add `SetConfigPathForTesting(path string)` and `ResetConfigPath()`
- [ ] `internal/testutil/e2e/helpers.go` - Use override in `SetupWithOnboarding()` cleanup

### Option B: Inject ConfigSaver Interface

Pass a config saver interface to the Model that can be mocked in tests.

**Pros**:
- More testable design
- No global state

**Cons**:
- Larger refactor
- Changes Model constructor signature

**Files to Change**:
- [ ] `internal/config/config.go` - Add `ConfigSaver` interface
- [ ] `internal/cli/app/app.go` - Accept `ConfigSaver` in constructor
- [ ] `internal/testutil/e2e/helpers.go` - Provide mock saver

### Option C: Test Mode Flag on Model

Add a flag to skip config saving in test mode.

**Pros**:
- Simple implementation

**Cons**:
- Test-specific code in production code
- Easy to forget to set flag

**Files to Change**:
- [ ] `internal/cli/app/app.go` - Add `SetTestMode(bool)`, skip save when true
- [ ] `internal/testutil/e2e/helpers.go` - Call `SetTestMode(true)`

**Selected Approach**: Option A  
**Rationale**: Minimal changes, follows existing isolation pattern, no production code changes

---

## Testing Plan

### Phase 1: Unit Tests
- [x] Test `SetConfigPathForTesting()` overrides `GetConfigPath()`
- [x] Test `ResetConfigPath()` restores default behavior
- [x] Test `SaveConfig()` uses override path when set

**Files**:
- `internal/config/config_test.go`

### Phase 2: Integration Tests
- [x] Test `SetupWithOnboarding()` uses temp config path
- [x] Test `CompleteOnboarding()` saves to temp path, not real path
- [x] Test cleanup resets config path

**Files**:
- `internal/testutil/e2e/onboarding_e2e_test.go` (BUG-007 regression tests)

### Phase 3: Regression Test
- [x] BUG-007 regression: Verify real config unchanged after onboarding tests

```go
It("BUG-007: should NOT write to user's real config file", func() {
    realConfigPath, _ := config.GetConfigPath()
    originalContent, originalExists := readFileIfExists(realConfigPath)
    
    env := e2e.SetupWithOnboarding(GinkgoT())
    defer env.Cleanup()
    env.CompleteOnboarding("Test User", "test@example.com")
    
    if originalExists {
        currentContent, _ := os.ReadFile(realConfigPath)
        Expect(currentContent).To(Equal(originalContent))
    } else {
        Expect(realConfigPath).NotTo(BeAnExistingFile())
    }
})
```

**Files**:
- `internal/testutil/e2e/onboarding_e2e_test.go`

### Phase 4: Manual Testing
- [ ] Run onboarding E2E tests
- [ ] Verify `~/.kariya/config.yaml` unchanged
- [ ] Verify temp config created in test temp directory

---

## Verification Checklist

### Code Quality
- [x] Fix implemented and tested
- [x] All tests passing (go test ./...)
- [ ] No race conditions (go test -race)
- [ ] Code coverage maintained (>80%)
- [x] Linting passing (staticcheck)

### Functionality
- [x] Issue no longer reproduces
- [x] Real config file remains unchanged after tests
- [x] Test config saved to temp directory
- [x] Cleanup properly resets config path

### Documentation
- [x] Code comments added/updated
- [x] Bug report updated with resolution

### Compliance
- [x] Follows project coding standards
- [ ] Atomic commits with clear messages
- [ ] AI attribution

---

## Related Files

### Implementation Files
- `internal/config/config.go:127` - `GetConfigPath()`
- `internal/config/config.go:235` - `SaveConfig()`
- `internal/cli/app/app.go:285` - Saves config on onboarding completion
- `internal/testutil/e2e/helpers.go:149` - `SetupWithOnboarding()`

### Test Files
- `internal/config/config_test.go`
- `internal/testutil/e2e/onboarding_e2e_test.go:124` - Uses `CompleteOnboarding()`

### Documentation Files
- `docs/development/SESSION_PROTOCOL.md`

---

## References

- Related bugs: None
- Documentation: E2E test helpers
- Pattern reference: Database isolation in `Setup()` at `helpers.go:86-87`

---

**Last Updated**: 2026-01-21  
**Updated By**: Claude Code
