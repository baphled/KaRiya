# Bug Debugging Guide

**Purpose**: Quick reference for systematic bug investigation and resolution

---

## Quick Start

```bash
# 1. Create new bug report
cp bugs/BUG_TEMPLATE.md bugs/bug-XXX-short-description.md

# 2. Edit bug report with details
vim bugs/bug-XXX-short-description.md

# 3. Update bugs/README.md active bugs table

# 4. Follow 6-phase debugging workflow below
```

---

## 6-Phase Debugging Workflow

### Phase 1: Reproduction (30-60 min)

**Goal**: Confirm issue exists and is reproducible

```bash
# Start application
go run ./cmd/cli

# OR build and run
go build -o kariya ./cmd/cli
./kariya
```

**Checklist**:
- [ ] Document exact steps to reproduce
- [ ] Verify issue occurs 100% of time (or X%)
- [ ] Test in clean environment
- [ ] Capture screenshots/error messages
- [ ] Note environment details (OS, terminal, Go version)

**Update Bug Report**:
- Fill in "Reproduction Steps"
- Update "Consistency" field
- Add "Environment" details
- Document "Actual Behavior" with evidence

---

### Phase 2: Investigation (1-3 hours)

**Goal**: Gather context and identify potential causes

#### A. Review Related Tests

```bash
# Run existing tests
go test -v ./internal/cli/intents -run TestNamePattern

# Run with Ginkgo
ginkgo -r --focus="FeatureName" ./internal/cli/intents

# Check if tests cover this scenario
grep -r "test description" internal/cli/intents/
```

**Questions**:
- Do tests exist for this scenario?
- Do they pass or fail?
- What's the gap between tests and reality?

#### B. Review Code

```bash
# Find relevant files
find . -name "*feature_name*" -type f

# Search for key functions/patterns
grep -r "functionName" internal/

# Check recent changes
git log --oneline --follow path/to/file.go
```

**Look For**:
- Message handling order
- State transitions
- Error handling
- Edge cases
- Similar patterns elsewhere

#### C. Check Documentation

```bash
# Search docs for feature
grep -r "feature" docs/

# Review standards
cat docs/TUI_STANDARDS.md
cat docs/development/NAVIGATION_TESTING_GUIDE.md
```

**Questions**:
- What does documentation claim?
- Does implementation match docs?
- Are standards followed?

**Update Bug Report**:
- Add "Investigation Log" entries with timestamps
- Note findings in "Code Review" section
- Document test review results

---

### Phase 3: Root Cause Analysis (1-2 hours)

**Goal**: Identify exact cause of failure

#### Techniques

**1. Trace Execution**:
```go
// Add debug prints temporarily
func Update(msg tea.Msg) tea.Cmd {
    fmt.Printf("[DEBUG] Received msg: %#v\n", msg)
    // ... rest of function
}
```

**2. Check Message Flow**:
```
User Action
  → Terminal Event
  → BubbleTea tea.Msg
  → Intent.Update(msg)
  → Sub-component.Update(msg)?
  → State Transition?
```

**3. Compare Working vs Broken**:
- Find similar working code
- Identify differences
- Determine critical difference

**4. Use Debugger**:
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug application
dlv debug ./cmd/cli
(dlv) break internal/cli/intents/capture_event_intent.go:298
(dlv) continue
```

**Update Bug Report**:
- Fill in "Root Cause" section with technical details
- Include file, function, line number
- Explain WHY it fails (not just WHAT fails)

---

### Phase 4: Fix Implementation (2-6 hours)

**Goal**: Implement minimal, surgical fix

#### Fix Strategy

**1. Design Approach**:
- List multiple options
- Evaluate pros/cons
- Select best approach
- Document rationale

**2. Implement Fix**:
```bash
# Create branch (optional)
git checkout -b fix/bug-XXX-description

# Edit files
vim path/to/file.go

# Test fix locally
go run ./cmd/cli
```

**3. Add/Update Tests**:
```bash
# Add unit tests
vim internal/cli/intents/feature_test.go

# Add E2E tests
vim internal/testutil/e2e/feature_e2e_test.go

# Run tests
go test -v ./internal/cli/intents/feature_test.go
ginkgo -r ./internal/testutil/e2e
```

**4. Verify Fix**:
```bash
# Run full test suite
go test ./...

# Run with race detector
go test -race ./...

# Check coverage
go test -cover ./...

# Run linting
staticcheck ./...

# Manual test
go run ./cmd/cli
```

**Update Bug Report**:
- Document "Fix Strategy" with chosen approach
- List all files modified
- Add code snippets showing changes

---

### Phase 5: Testing & Verification (1-2 hours)

**Goal**: Ensure fix works and no regressions

#### Test Levels

**1. Unit Tests**:
```bash
# Run specific tests
go test -v ./internal/cli/intents -run TestFeature

# Run all intent tests
go test -v ./internal/cli/intents/...
```

**2. Integration Tests**:
```bash
# Run E2E tests
ginkgo -r ./internal/testutil/e2e
```

**3. Manual Testing**:
- Follow reproduction steps
- Verify issue no longer occurs
- Test edge cases
- Test related features

**4. Regression Testing**:
```bash
# Full test suite
make test

# Compliance check
make check-compliance

# Build verification
make build
```

**Update Bug Report**:
- Check off all items in "Testing Plan"
- Document test results
- Note any unexpected findings

---

### Phase 6: Documentation (30 min)

**Goal**: Update docs and close out bug

#### Tasks

**1. Code Documentation**:
```go
// Add comments explaining fix
// BUGFIX: Handle global keys before delegating to form
// See bug-001 for context
if keyMsg, ok := msg.(tea.KeyMsg); ok {
    // ...
}
```

**2. Update Bug Report**:
- Fill in "Resolution Summary"
- Add commit hashes
- Mark status as "Closed"
- Add follow-up actions if needed

**3. Update Project Docs** (if needed):
```bash
# Update relevant guides
vim docs/TUI_STANDARDS.md
vim docs/development/NAVIGATION_TESTING_GUIDE.md
```

**4. Commit Changes**:
```bash
# Stage changes
git add path/to/fixed/file.go
git add path/to/test/file.go

# Commit with clear message
git commit -m "fix(intents): handle escape key before form delegation

Fixes escape key navigation in CaptureEvent forms by checking
global keys before delegating to Huh form component.

Root cause: Huh forms were consuming escape key events before
intent's HandleGlobalKeys could process them.

Solution: Reorder message handling to check global keys first.

Fixes: bug-001
"

# If AI-generated
git commit -m "fix(intents): handle escape key before form delegation

[commit message body]

AI-Generated-By: OpenCode (Claude)
Reviewed-By: [Your Name]
"
```

**5. Update README**:
```bash
# Move bug from Active to Closed
vim bugs/README.md
```

---

## Common Debugging Patterns

### Pattern 1: Message Not Reaching Handler

**Symptom**: Key press does nothing

**Check**:
1. Is message being sent? (Add debug print)
2. Is it reaching Update method? (Add debug print)
3. Is it being handled? (Check switch cases)
4. Is sub-component consuming it? (Check delegation order)

**Common Causes**:
- Sub-component consumes message first
- Wrong message type (KeyMsg vs KeyRunes)
- Missing case in switch statement
- Wrong state for handler

---

### Pattern 2: Tests Pass But App Fails

**Symptom**: Unit tests green, but manual test fails

**Check**:
1. Do tests reflect actual app flow?
2. Are tests bypassing real components?
3. Is test setup different from app setup?
4. Missing E2E tests?

**Common Causes**:
- Tests call methods directly (bypass message routing)
- Test mocks too much
- Missing integration tests
- Unit tests too isolated

---

### Pattern 3: Inconsistent Behavior

**Symptom**: Works sometimes, fails other times

**Check**:
1. Race conditions? (Run with -race flag)
2. State-dependent? (Check state transitions)
3. Timing-dependent? (Check async operations)
4. Environment-dependent? (Test on different systems)

**Common Causes**:
- Race conditions
- Uninitialized state
- Async timing issues
- Terminal/environment differences

---

## Tools & Commands

### Testing

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -cover ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific test
go test -v ./internal/cli/intents -run TestEscapeKey

# Run Ginkgo tests
ginkgo -r ./internal/cli/intents
ginkgo -r --focus="Escape" ./internal/cli/intents
```

### Code Analysis

```bash
# Static analysis
staticcheck ./...

# Code formatting
go fmt ./...
gofmt -d .

# Vet
go vet ./...

# Find TODOs
grep -r "TODO" internal/
grep -r "FIXME" internal/
grep -r "BUG" internal/
```

### Debugging

```bash
# Debug with Delve
dlv debug ./cmd/cli
(dlv) break file.go:123
(dlv) continue
(dlv) print variableName
(dlv) next
(dlv) step

# Trace execution
go run -trace trace.out ./cmd/cli
go tool trace trace.out

# CPU profiling
go test -cpuprofile cpu.prof -bench .
go tool pprof cpu.prof
```

### Git

```bash
# Find when bug introduced
git bisect start
git bisect bad HEAD
git bisect good v1.0.0
# ... git bisect will check out commits to test
git bisect reset

# View changes to file
git log --follow -p path/to/file.go

# Find commits touching function
git log -S "functionName" --source --all

# Blame specific lines
git blame path/to/file.go
```

---

## Best Practices

### DO

- ✅ Reproduce issue before investigating
- ✅ Document findings in bug report
- ✅ Write tests that fail before fix
- ✅ Implement minimal fix
- ✅ Verify no regressions
- ✅ Update documentation
- ✅ Use descriptive commit messages
- ✅ Add code comments explaining fix

### DON'T

- ❌ Skip reproduction step
- ❌ Make large refactors for small bugs
- ❌ Fix without tests
- ❌ Commit debugging code (fmt.Printf, etc.)
- ❌ Leave TODOs in code
- ❌ Close bug without verification
- ❌ Skip documentation updates

---

## Getting Help

### Resources

- **Documentation**: `/docs` directory
- **Development Rules**: `/docs/rules`
- **Testing Guides**: `/docs/development`
- **Bug Reports**: `/bugs` directory
- **Task Files**: `/tasks` directory

### Escalation

If stuck:
1. Review similar bugs in `/bugs`
2. Check related tasks in `/tasks`
3. Search documentation in `/docs`
4. Review test patterns in codebase
5. Ask for help (provide bug report link)

---

**Last Updated**: 2026-01-12
