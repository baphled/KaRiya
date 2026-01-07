# Quick Start: Standardized View Implementation

## What We're Building

A standardized layout system for all KaRiya TUI screens that ensures:
- ✅ Logo always visible at the top (2-line spacing)
- ✅ Content centered and using full terminal width
- ✅ Context-aware help at the bottom with visual separator
- ✅ Modal overlays for errors, loading, progress, and success
- ✅ Consistent user experience across all screens

## Task Sequence

Execute tasks in this order:

### 1️⃣ Task 12: Core Components (4-5 hours)
**File**: `tasks-12-standardized-view-core-components.md`

**Start here**: Create the foundational components
```bash
# Open the task file
cat tasks/tasks-12-standardized-view-core-components.md

# Key deliverables:
# - StandardView component
# - Modal system (Error, Loading, Progress, Success, Warning)
# - LoadingMessageRotator with spinner
```

**Quick wins**:
- SmartContainer already exists - we're building on top of it
- ASCIILogo already exists - we're integrating it
- Most infrastructure is in place - we're standardizing usage

### 2️⃣ Task 13: Infrastructure (2-3 hours)
**File**: `tasks-13-intent-infrastructure-updates.md`

**After Task 12**: Wire up the infrastructure
```bash
# Key deliverables:
# - Enhanced BaseIntent with terminal awareness
# - View helper functions
# - Modal helper functions
# - State management helpers
```

### 3️⃣ Task 14: Primary Intents (5 hours)
**File**: `tasks-14-migrate-capture-browse-intents.md`

**After Task 13**: Migrate the most complex intents
```bash
# CaptureEventIntent - Most complex, establishes patterns
# BrowseTimelineIntent - Second most used

# Pattern established here will be followed by all other intents
```

### 4️⃣ Task 15: Remaining Intents (4.5 hours)
**File**: `tasks-15-migrate-remaining-intents.md`

**After Task 14**: Complete the migration
```bash
# GenerateCVIntent - With progress tracking
# ExportArtifactIntent - With export progress
# ConfigureSystemIntent - With configuration states
```

### 5️⃣ Task 16: Testing & Polish (3-4 hours)
**File**: `tasks-16-standardview-testing-polish.md`

**After Tasks 12-15**: Ensure quality and consistency
```bash
# Comprehensive testing
# Documentation creation
# Performance verification
# Cross-intent consistency checks
```

## Getting Started Right Now

### Step 1: Review Task 12
```bash
# Read the task file
cat tasks/tasks-12-standardized-view-core-components.md

# Review prerequisites
go test ./internal/cli/components/...
make check-compliance
```

### Step 2: Create First Component
```bash
# Create the StandardView component file
touch internal/cli/components/standard_view.go

# Follow Phase 2 of Task 12
# Start with the struct definition
```

### Step 3: Follow the Checklist
Each task file has a detailed **Implementation Checklist** with:
- ✅ Checkboxes for each step
- ⏱️ Time estimates per phase
- 📝 Code examples where helpful
- 💾 Commit messages for atomic commits
- ✅ Testing instructions

### Step 4: Test Frequently
```bash
# After each phase, run tests
go test ./internal/cli/components/...

# Use race detector
go test -race ./internal/cli/components/...

# Check coverage
go test -cover ./internal/cli/components/...
```

## Key Patterns to Follow

### Standard View Pattern (from Task 14)
```go
func (i *SomeIntent) View() string {
    // Create view with breadcrumbs
    view := i.CreateViewWithBreadcrumbs("Main Menu", "Intent Name", i.getStateName())
    
    // Handle modal states
    if i.isLoading {
        ShowLoadingModal(view, i.loadingMessage, true)
    }
    if i.errorState != nil {
        ShowErrorModal(view, i.errorState)
    }
    if i.ShouldShowSuccess() {
        ShowSuccessModal(view, i.successMessage)
    }
    
    // Set content and help
    view.WithContent(i.getStateContent())
    view.WithHelp(i.getContextHelp()).WithFooterSeparator(true)
    
    return view.Render()
}
```

### State Content Pattern
```go
func (i *SomeIntent) getStateContent() string {
    switch i.state {
    case StateA:
        return i.renderStateA()
    case StateB:
        return i.renderStateB()
    default:
        return "Unknown state"
    }
}
```

### Context Help Pattern
```go
func (i *SomeIntent) getContextHelp() string {
    base := "q Quit  m Main Menu"
    
    switch i.state {
    case StateA:
        return CombineFooters(NavigationFooter(), base)
    case StateB:
        return CombineFooters("Enter Confirm  Esc Back", base)
    default:
        return base
    }
}
```

## Commit Message Format

Follow conventional commits throughout:

```bash
# Features
feat(components): add StandardView component
feat(intents): add view helper functions

# Refactoring
refactor(capture): update View method to use StandardView
refactor(browse): extract state content rendering

# Tests
test(components): add StandardView tests
test(intents): add consistency tests

# Documentation
docs(components): add StandardView developer guide
docs: update TUI standards for StandardView

# Fixes
fix(components): resolve modal sizing on small terminals
fix(intents): correct breadcrumb navigation
```

## Development Workflow (Per Task)

### 1. Preparation
```bash
# Check baseline
make check-compliance
go test ./...

# Review task file
cat tasks/tasks-XX-name.md

# Start task
# Mark preparation phase items complete
```

### 2. Implementation
```bash
# For each phase:
# - Read phase description
# - Implement changes
# - Run tests
# - Commit atomically

# Example:
git add internal/cli/components/standard_view.go
git commit -m "feat(components): add StandardView struct and basic structure"
```

### 3. Testing
```bash
# After implementation:
go test -v ./...
go test -race ./...
go test -cover ./...

# Manual testing if applicable
go build -o kariya ./cmd/kariya
./kariya
```

### 4. Verification
```bash
# Before marking task complete:
make check-compliance
golangci-lint run ./...
go fmt ./...

# Review all commits
git log --oneline

# Update task checklist
# Mark all phases complete
```

## Quick Reference: Modal Types

```go
// Error Modal - Red border, bell, Esc to dismiss
ShowErrorModal(view, err)

// Loading Modal - Spinner, cancellable
ShowLoadingModal(view, "Loading data...", true)

// Progress Modal - Percentage bar
ShowProgressModal(view, "Exporting", "Writing file...", 0.65)

// Success Modal - Green border, auto-dismiss 3s
ShowSuccessModal(view, "✅ Operation completed!")
```

## Quick Reference: Footer Helpers

```go
// Predefined footers
NavigationFooter()      // ↑/k Up  ↓/j Down  Enter Select  Esc Back
FormFooter()           // Tab Next  Shift+Tab Previous  Enter Submit
ListFooter()           // Navigation + / Search
DetailViewFooter()     // ↑/↓ Scroll  Esc Back

// Combine footers
CombineFooters(NavigationFooter(), "q Quit", "m Main Menu")
// Output: ↑/k Up  ↓/j Down  Enter Select  Esc Back  |  q Quit  m Main Menu
```

## Common Issues & Solutions

### Issue: Tests failing after changes
**Solution**: 
- Review test expectations - they may check for old view structure
- Update assertions to look for logo, breadcrumbs, separator
- Add new tests for StandardView features

### Issue: View not centered
**Solution**:
- Verify terminal info is valid and propagated
- Check SmartContainer is using CenterBoth mode
- Ensure StandardView.Render() is calling SmartContainer

### Issue: Modal not appearing
**Solution**:
- Check `ShowModal` flag is true
- Verify modal content is not nil
- Ensure view.Render() is being called after modal is shown

### Issue: Loading messages not rotating
**Solution**:
- Check rotateInterval is set (default: 2s)
- Verify lastRotation is being updated
- Call Rotate() method, not just GetCurrent()

## Time Tracking

Expected total time: **19-21 hours**

| Task | Estimate | Actual | Status |
|------|----------|--------|--------|
| Task 12 | 4-5h | - | ⬜ Not Started |
| Task 13 | 2-3h | - | ⬜ Not Started |
| Task 14 | 5h | - | ⬜ Not Started |
| Task 15 | 4.5h | - | ⬜ Not Started |
| Task 16 | 3-4h | - | ⬜ Not Started |

## Success Checklist

When all tasks are complete, you should have:

- [ ] All intents use StandardView
- [ ] Logo visible on every screen
- [ ] Breadcrumbs show navigation context
- [ ] Footer separator visible
- [ ] Context-aware help on every screen
- [ ] Modal system working (errors, loading, progress, success)
- [ ] Loading messages rotate during long operations
- [ ] Success messages auto-dismiss after 3s
- [ ] Error messages dismissible with Esc
- [ ] Terminal resize handled gracefully
- [ ] All tests passing (100%)
- [ ] Code coverage > 87%
- [ ] Documentation complete
- [ ] No regressions

## Getting Help

If you get stuck:

1. **Check the task file** - Detailed guidance for each phase
2. **Review the summary** - `tasks-12-16-SUMMARY.md`
3. **Look at existing code** - `internal/cli/app/app.go` has good examples
4. **Run tests** - They often reveal what's missing
5. **Check documentation** - `docs/TUI_DEVELOPER_GUIDE.md`, `docs/TUI_STANDARDS.md`

## Ready to Start?

```bash
# Step 1: Open Task 12
cat tasks/tasks-12-standardized-view-core-components.md

# Step 2: Mark preparation phase as in progress
# Edit the task file and check off items as you complete them

# Step 3: Start with Phase 2.1 - Create basic structure
touch internal/cli/components/standard_view.go

# Step 4: Follow the checklist, commit atomically, test frequently

# Good luck! 🚀
```

---

**Remember**: 
- ✅ Follow the master-task-prompt workflow (5 phases)
- ✅ Make atomic commits with conventional messages
- ✅ Test frequently with race detector
- ✅ Update task checklists as you go
- ✅ Ask for help if stuck

**You've got this!** The task files are comprehensive and well-structured. Take it one phase at a time, and you'll have a beautiful, consistent TUI in no time.
