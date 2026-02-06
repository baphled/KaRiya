# VHS Demo Generation Prompt

## Purpose

This prompt is triggered **after a feature is marked as "done"** to ensure all UI-related changes have visual documentation. This creates a consistent, automated workflow for generating VHS demo recordings.

---

## When to Trigger

This prompt should be invoked when:

1. **New Feature Completed** - A new intent, screen, or workflow is implemented
2. **Bug Fix Completed** - A UI bug has been fixed (update existing tape)
3. **Enhancement Completed** - An existing feature has been enhanced (update existing tape)
4. **Task Status Changed to "Done"** - Final verification before closing

---

## Workflow Overview

```
Task Complete
     |
     v
+------------------+
| Read Task File   |
| Examine Changes  |
+------------------+
     |
     v
+------------------+
| Determine Type   |
| (New/Bug/Enhance)|
+------------------+
     |
     +---> New Feature --> Create new tape set
     |
     +---> Bug Fix ------> Find & update existing tape
     |
     +---> Enhancement --> Find & update existing tape
     |
     v
+------------------+
| Generate Demos   |
| Verify Output    |
+------------------+
     |
     v
+------------------+
| Include in PR    |
+------------------+
```

---

## Phase 1: Task Analysis

### 1.1 Read the Task

```bash
# View the task file or commit history
git log --oneline -10
git show HEAD

# Identify what was changed
git diff HEAD~1 --stat
```

**Questions to answer:**
- What feature/area was modified?
- Is this a new feature, bug fix, or enhancement?
- What intent/screen was affected?
- What user-visible changes were made?

### 1.2 Examine the Code

```bash
# Find affected intents
ls internal/cli/intents/

# Find affected screens
ls internal/cli/screens/

# Check for new screens or modals
git diff HEAD~1 --name-only | grep -E "(screen|modal|intent)"
```

**Identify:**
- New intents created
- New screens added
- Modified workflows
- Changed navigation patterns
- New modals or overlays

### 1.3 Determine Change Type

| Change Type | Criteria | Action |
|-------------|----------|--------|
| **New Feature** | New intent, new workflow, new capability | Create new tape set (3 tapes) |
| **Bug Fix** | Fixes existing UI behavior | Update existing tape to verify fix |
| **Enhancement** | Improves existing feature | Update existing tape to show improvement |

---

## Phase 2: New Feature - Create Tape Set

### 2.1 Create Feature Directory

```bash
# Create directory for the feature
mkdir -p demos/vhs/features/<feature-name>

# Copy templates
cp demos/vhs/features/template/happy-path.tape demos/vhs/features/<feature-name>/
cp demos/vhs/features/template/sad-path.tape demos/vhs/features/<feature-name>/
cp demos/vhs/features/template/edge-cases.tape demos/vhs/features/<feature-name>/
```

### 2.2 Analyze the Feature

Before writing tapes, understand:

1. **Entry Point** - How does the user access this feature?
2. **Happy Path** - What is the successful workflow?
3. **Error Cases** - What validations exist? What errors can occur?
4. **Edge Cases** - Cancel, back navigation, empty states?
5. **Keyboard Shortcuts** - What keys trigger actions?

```bash
# Read the intent to understand states
view internal/cli/intents/<feature>/intent.go

# Read screens to understand UI
view internal/cli/screens/<feature>/*.go

# Check for keyboard handling
grep -n "KeyMsg\|tea.Key" internal/cli/intents/<feature>/*.go
```

### 2.3 Write Happy Path Tape

**Template structure:**

```tape
# <Feature Name> - Happy Path Demo
# ================================
# Demonstrates successful workflow completion.
#
# What this demo shows:
#   - <Step 1 description>
#   - <Step 2 description>
#   - <Success confirmation>
#
# Duration: ~<X>s
# Keyboard shortcuts used: <list>

Source demos/vhs/config.tape

Output demos/vhs/features/<feature-name>/happy-path.gif
Output demos/vhs/features/<feature-name>/happy-path.mp4

# Setup (hidden)
Hide
Type `rm -rf /tmp/kariya-demo && mkdir -p /tmp/kariya-demo`
Enter
Sleep 300ms
Type `cp demos/vhs/minimal-config.yaml /tmp/kariya-demo/config.yaml`
Enter
Sleep 300ms
Type "cd /home/baphled/Projects/KoRiya"
Enter
Sleep 300ms
# Import sample data if needed
# Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
# Enter
# Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter
Sleep 2s

# === STEP 1: Navigate to feature ===
# <Navigation description>
Down
Down
Enter
Sleep 1s
Screenshot demos/vhs/features/<feature-name>/step1-entry.png

# === STEP 2: Perform main action ===
# <Action description>
Type@100ms "<input data>"
Tab
Sleep 500ms
Screenshot demos/vhs/features/<feature-name>/step2-input.png

# === STEP 3: Submit and confirm ===
# <Submission description>
Enter
Sleep 2s
Screenshot demos/vhs/features/<feature-name>/step3-success.png

# Exit cleanly
Escape
Sleep 500ms
Ctrl+C
Sleep 500ms
```

### 2.4 Write Sad Path Tape

**Focus on:**
- Validation errors
- Error messages
- Recovery from errors
- User-friendly error handling

```tape
# <Feature Name> - Sad Path Demo
# ==============================
# Demonstrates error handling and validation.
#
# Scenarios covered:
#   - <Validation error 1>
#   - <Error message display>
#   - <Recovery demonstration>

Source demos/vhs/config.tape

Output demos/vhs/features/<feature-name>/sad-path.gif
Output demos/vhs/features/<feature-name>/sad-path.mp4

# ... setup (same as happy path) ...

# === SCENARIO 1: Validation Error ===
# Try to submit without required data
Enter
Sleep 1s
Screenshot demos/vhs/features/<feature-name>/validation-error.png

# === SCENARIO 2: Recovery ===
# Fix the error and retry
Type@100ms "valid-input"
Enter
Sleep 2s
Screenshot demos/vhs/features/<feature-name>/error-recovered.png

# Exit cleanly
Escape
Ctrl+C
```

### 2.5 Write Edge Cases Tape

**Focus on:**
- Cancel at various points
- Back navigation
- Empty states
- Boundary values (long input, special characters)
- Rapid navigation

```tape
# <Feature Name> - Edge Cases Demo
# ================================
# Demonstrates boundary conditions and special scenarios.
#
# Scenarios covered:
#   - Cancel mid-workflow
#   - Back navigation with state
#   - Empty state handling
#   - Long input handling

Source demos/vhs/config.tape

Output demos/vhs/features/<feature-name>/edge-cases.gif
Output demos/vhs/features/<feature-name>/edge-cases.mp4

# ... setup (same as happy path) ...

# === SCENARIO 1: Cancel Workflow ===
Down
Enter
Sleep 500ms
Escape
Sleep 500ms
Screenshot demos/vhs/features/<feature-name>/cancel-confirmed.png

# === SCENARIO 2: Empty State ===
# Navigate to list with no items (use empty-config.yaml)
# ...

# === SCENARIO 3: Long Input ===
Type@50ms "This is a very long input to test text wrapping and overflow handling"
Sleep 1s
Screenshot demos/vhs/features/<feature-name>/long-input.png

# Exit cleanly
Escape
Ctrl+C
```

### 2.6 Generate and Verify

```bash
# Generate all tapes for the feature
make vhs-feature FEATURE=<feature-name>

# Verify outputs exist
ls -la demos/vhs/features/<feature-name>/*.gif
ls -la demos/vhs/features/<feature-name>/*.png

# Watch the GIFs to verify they work
# (open in browser or image viewer)
```

---

## Phase 3: Bug Fix / Enhancement - Update Existing Tape

### 3.1 Find the Existing Tape

```bash
# Find tapes related to the feature
find demos/vhs -name "*.tape" | xargs grep -l "<feature-keyword>"

# Or list by directory
ls demos/vhs/features/<feature>/
```

**Common tape locations:**
- `demos/vhs/features/<feature>/` - Feature-specific demos

### 3.2 Determine What to Update

| Change | Update Required |
|--------|-----------------|
| Fixed validation error | Update sad-path.tape to show correct behavior |
| Changed navigation | Update happy-path.tape with new navigation |
| Added keyboard shortcut | Add shortcut demo to relevant tape |
| Fixed visual bug | Update screenshots in affected tape |
| Added new option | Extend happy-path.tape with new option |

### 3.3 Update the Tape

```bash
# Edit the relevant tape
view demos/vhs/features/<feature>/happy-path.tape

# Make changes to reflect the fix/enhancement
# Add new steps, update navigation, change timing
```

**Common updates:**
- Add/remove navigation steps
- Update timing (`Sleep` durations)
- Add new screenshots for new states
- Update comments to reflect changes

### 3.4 Regenerate and Verify

```bash
# Regenerate the specific tape
vhs demos/vhs/features/<feature>/happy-path.tape

# Or regenerate all tapes for the feature
make vhs-feature FEATURE=<feature>

# Compare with previous version (if golden exists)
make vhs-golden-compare
```

---

## Phase 4: Include in PR

### 4.1 Stage Demo Files

```bash
# Stage GIF files (for PR display)
git add demos/vhs/features/<feature-name>/*.gif

# Stage tape files (for reproducibility)
git add demos/vhs/features/<feature-name>/*.tape

# Stage key screenshots (optional)
git add demos/vhs/features/<feature-name>/*.png
```

### 4.2 Add to PR Description

Include in the PR description:

```markdown
## Demo Evidence

### Happy Path
![Happy path demo](./demos/vhs/features/<feature-name>/happy-path.gif)

### Error Handling
![Sad path demo](./demos/vhs/features/<feature-name>/sad-path.gif)

### Edge Cases
![Edge cases demo](./demos/vhs/features/<feature-name>/edge-cases.gif)

### Key Screenshots

| State | Screenshot |
|-------|------------|
| Initial | ![](./demos/vhs/features/<feature-name>/step1-entry.png) |
| Success | ![](./demos/vhs/features/<feature-name>/step3-success.png) |
```

### 4.3 Commit Demo Files

```bash
# Commit separately from code changes
git add demos/vhs/features/<feature-name>/
make ai-commit FILE=/tmp/commit.txt

# Commit message example:
# docs(demos): add VHS demo for <feature-name>
```

---

## Decision Tree

```
Is this a UI-affecting change?
     |
     +---> No --> Skip VHS demo generation
     |
     +---> Yes
            |
            v
     Is this a new feature?
            |
            +---> Yes --> Create new tape set (3 tapes)
            |             Location: demos/vhs/features/<feature-name>/
            |
            +---> No
                   |
                   v
            Is this a bug fix or enhancement?
                   |
                   +---> Yes --> Find and update existing tape
                   |             Regenerate and verify
                   |
                   +---> No --> Skip (non-UI change)
```

---

## Checklist

### For New Features

- [ ] Created `demos/vhs/features/<feature-name>/` directory
- [ ] Created `happy-path.tape` with successful workflow
- [ ] Created `sad-path.tape` with error scenarios
- [ ] Created `edge-cases.tape` with boundary conditions
- [ ] All tapes have proper header comments
- [ ] All tapes use `Source demos/vhs/config.tape`
- [ ] Screenshots captured at key moments
- [ ] Demos generated successfully (`make vhs-feature FEATURE=<name>`)
- [ ] GIFs reviewed and look correct
- [ ] Demo files staged for commit
- [ ] PR description includes demo evidence

### For Bug Fixes / Enhancements

- [ ] Identified affected tape(s)
- [ ] Updated tape to reflect fix/enhancement
- [ ] Regenerated demo
- [ ] Compared with previous version (if applicable)
- [ ] Demo shows the fix/enhancement clearly
- [ ] Updated files staged for commit
- [ ] PR description mentions demo update

---

## Quick Reference

### Common VHS Commands

```tape
# Navigation
Down / Up / Left / Right    # Arrow keys
Enter                       # Enter key
Escape                      # Escape key
Tab                         # Tab key
Space                       # Space bar
Ctrl+C                      # Ctrl+C

# Typing
Type "fast text"            # Fast typing
Type@100ms "slow text"      # Slow typing (100ms per char)
Type@50ms "medium text"     # Medium typing

# Timing
Sleep 500ms                 # Half second
Sleep 1s                    # One second
Sleep 2s                    # Two seconds

# Visibility
Hide                        # Hide following commands
Show                        # Show following commands

# Capture
Screenshot path/to/file.png # Take screenshot
Output path/to/file.gif     # Set GIF output
Output path/to/file.mp4     # Set MP4 output
```

### Make Targets

```bash
# Generate all demos
make vhs-demos

# Generate specific feature
make vhs-feature FEATURE=<name>

# Generate all tapes for a feature
make vhs-feature-all FEATURE=<name>

# Feature shortcuts (for common features)
make vhs-skills-all
make vhs-capture-all
make vhs-timeline-all

# Golden file comparison
make vhs-golden-compare
make vhs-golden-update
```

### Timing Guidelines

| Action | Recommended Sleep |
|--------|-------------------|
| After launch | `Sleep 2s` |
| After navigation | `Sleep 500ms` - `1s` |
| After typing | `Sleep 500ms` |
| After submission | `Sleep 2s` |
| For screenshots | `Sleep 1s` before |
| After error display | `Sleep 1s` - `2s` |

---

## Examples

### Example 1: New Feature (Burst Management)

```bash
# 1. Create directory
mkdir -p demos/vhs/features/burst-management

# 2. Copy templates
cp demos/vhs/features/template/*.tape demos/vhs/features/burst-management/

# 3. Edit tapes for burst management workflow
# - happy-path.tape: Detect bursts, review, accept
# - sad-path.tape: No bursts found, validation errors
# - edge-cases.tape: Cancel detection, empty timeline

# 4. Generate
make vhs-feature FEATURE=burst-management

# 5. Verify and commit
git add demos/vhs/features/burst-management/
make ai-commit FILE=/tmp/commit.txt
```

### Example 2: Bug Fix (Skills Display)

```bash
# 1. Find existing tape
find demos/vhs -name "*.tape" | xargs grep -l "skills"
# Found: demos/vhs/features/skills/happy-path.tape

# 2. Edit to show fixed behavior
view demos/vhs/features/skills/happy-path.tape
# Update navigation or add verification step

# 3. Regenerate
vhs demos/vhs/features/skills/happy-path.tape

# 4. Commit update
git add demos/vhs/features/skills/happy-path.tape
git add demos/vhs/features/skills/happy-path.gif
make ai-commit FILE=/tmp/commit.txt
```

---

## Related Documentation

- [VHS README](../../demos/vhs/README.md) - VHS setup and configuration
- [Workflow Documentation Guide](../workflows/WORKFLOW_DOCUMENTATION_GUIDE.md) - Complete workflow docs
- [Master Task Prompt](./master-task-prompt.md) - Overall task workflow
- [Development Workflow](../development/DEVELOPMENT_WORKFLOW.md) - Development process

---

*Last Updated: 2025-02-05*
*Version: 1.0*
