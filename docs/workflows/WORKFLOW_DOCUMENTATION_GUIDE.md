# Workflow Documentation Guide

This guide explains how to create comprehensive workflow documentation for KaRiya features, including written docs and VHS demo recordings.

---

## Overview

Every significant workflow in KaRiya should be documented with:

1. **Written documentation** - Markdown file describing states, navigation, and usage
2. **VHS demo tapes** - Automated terminal recordings showing the workflow
3. **Generated media** - GIFs/screenshots for README, PRs, and marketing

This approach ensures:
- Documentation stays current (regenerated from actual code)
- PRs include visual evidence of functionality
- Visual regression testing catches unintended UI changes
- Professional demos for marketing/onboarding

---

## Quick Start

### For Existing Workflows

```bash
# Generate demo for an existing workflow
make vhs-capture    # Capture event workflow
make vhs-browse     # Browse timeline workflow
make vhs-cv         # CV generation workflow
# etc.
```

### For New Features

```bash
# 1. Create feature demo directory
mkdir -p demos/vhs/features/my-feature

# 2. Copy templates
cp demos/vhs/features/template/*.tape demos/vhs/features/my-feature/

# 3. Edit tapes for your feature
# - happy-path.tape: Successful workflow
# - sad-path.tape: Error handling
# - edge-cases.tape: Cancel, empty states, etc.

# 4. Generate demos
make vhs-feature FEATURE=my-feature
```

---

## Creating Written Documentation

### File Location

Place workflow documentation in `docs/workflows/`:

```
docs/workflows/
├── README.md                     # Index of all workflows
├── WORKFLOW_DOCUMENTATION_GUIDE.md  # This guide
├── your-workflow.md              # Your workflow doc
└── diagrams/
    └── your_workflow.mermaid     # Generated state diagram
```

### Template Structure

Use this structure for workflow documentation:

```markdown
# [Workflow Name] Workflow

**Complexity**: [1-5 stars]
**Purpose**: [One sentence description]
**Implementation**: `internal/cli/intents/[intent].go`

---

## Overview

[2-3 paragraphs describing the workflow, when to use it, and prerequisites]

---

## Workflow States

[State diagram - embed Mermaid or link to generated diagram]

| State | Description | Next States |
|-------|-------------|-------------|
| ... | ... | ... |

---

## Step-by-Step Guide

### State 1: [Name]

**What happens**: [Description]
**User actions**: [Available shortcuts]
**Transitions**: [Where user can go]

[Screenshot or GIF]

### State 2: [Name]
...

---

## Keyboard Shortcuts

| Context | Shortcut | Action |
|---------|----------|--------|
| ... | ... | ... |

---

## Demo Videos

| Scenario | Demo |
|----------|------|
| Happy path | ![Happy path](../../demos/vhs/output/gifs/workflow-happy.gif) |
| Error handling | ![Sad path](../../demos/vhs/output/gifs/workflow-sad.gif) |
| Edge cases | ![Edge cases](../../demos/vhs/output/gifs/workflow-edge.gif) |

---

## Troubleshooting

### [Common Issue 1]
**Symptom**: ...
**Cause**: ...
**Solution**: ...

---

## Technical Details

[Implementation notes for developers]
```

---

## Creating VHS Demo Tapes

### Directory Structure

```
demos/vhs/
├── config.tape              # Shared settings (sourced by all tapes)
├── tapes/                   # Main workflow demos (for README)
│   ├── 01-onboarding.tape
│   ├── 02-capture-event.tape
│   └── ...
├── features/                # Feature-specific demos (for PRs)
│   ├── template/            # Copy these for new features
│   │   ├── happy-path.tape
│   │   ├── sad-path.tape
│   │   └── edge-cases.tape
│   └── your-feature/
│       ├── happy-path.tape
│       ├── sad-path.tape
│       └── edge-cases.tape
├── golden/                  # Baseline screenshots for regression
└── output/                  # Generated files (gitignored)
```

### Required Demo Scenarios

Every feature should have three demo tapes:

| Tape | Purpose | Key Elements |
|------|---------|--------------|
| **happy-path.tape** | Successful completion | Normal input, successful submission, confirmation |
| **sad-path.tape** | Error handling | Validation errors, clear messages, recovery |
| **edge-cases.tape** | Boundary conditions | Cancel, back nav, empty states, long input |

### VHS Tape Syntax

```tape
# Source shared config
Source demos/vhs/config.tape

# Output files (GIF and MP4)
Output demos/vhs/features/my-feature/happy-path.gif
Output demos/vhs/features/my-feature/happy-path.mp4

# Setup (hidden from recording)
Hide
Type `setup commands here`
Enter
Show

# Launch application
Type "./kariya --config /tmp/demo/config.yaml"
Enter
Sleep 2s

# Navigation
Down                    # Press down arrow
Up                      # Press up arrow
Enter                   # Press enter
Escape                  # Press escape
Tab                     # Press tab
Ctrl+C                  # Ctrl+C

# Typing
Type "fast typing"            # Types quickly
Type@100ms "slow typing"      # Types at 100ms per char

# Timing
Sleep 500ms             # Wait 500 milliseconds
Sleep 2s                # Wait 2 seconds

# Screenshots (for PR evidence)
Screenshot demos/vhs/features/my-feature/step1.png

# Special keys
Backspace
Space
Left / Right
```

### Best Practices

1. **Start with setup hidden**: Use `Hide`/`Show` to hide setup commands
2. **Use realistic timing**: Add `Sleep` after actions to show results
3. **Take screenshots at key moments**: Capture important states
4. **Add comments**: Document what each section does
5. **Test incrementally**: Run tape after each major addition

---

## Integrating Demos into PRs

### PR Description Template

When creating a PR that adds/modifies a workflow, include:

```markdown
## Demo Evidence

### Happy Path
![Happy path demo](./demos/vhs/features/my-feature/happy-path.gif)

### Error Handling
![Sad path demo](./demos/vhs/features/my-feature/sad-path.gif)

### Edge Cases
![Edge cases demo](./demos/vhs/features/my-feature/edge-cases.gif)

### Key Screenshots
| State | Screenshot |
|-------|------------|
| Initial | ![](./demos/vhs/features/my-feature/step1.png) |
| Success | ![](./demos/vhs/features/my-feature/success.png) |
```

### Generating PR Evidence

```bash
# Generate all demos for your feature
make vhs-feature FEATURE=my-feature

# Output will be in:
# - demos/vhs/features/my-feature/happy-path.gif
# - demos/vhs/features/my-feature/sad-path.gif
# - demos/vhs/features/my-feature/edge-cases.gif
# - demos/vhs/features/my-feature/*.png (screenshots)
```

---

## Visual Regression Testing

### Purpose

Detect unintended UI changes by comparing screenshots against known baselines.

### Workflow

```bash
# 1. Generate current screenshots
make vhs-golden-generate

# 2. Compare against baselines
make vhs-golden-compare

# 3. If changes are intentional, update baselines
make vhs-golden-update
```

### CI Integration

Add to CI workflow:

```yaml
visual-regression:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Install VHS
      run: |
        brew install vhs || \
        (curl -fsSL https://github.com/charmbracelet/vhs/releases/download/v0.7.1/vhs_0.7.1_linux_amd64.tar.gz | tar xz && sudo mv vhs /usr/local/bin/)
    - name: Build
      run: make build
    - name: Generate screenshots
      run: make vhs-golden-generate
    - name: Compare
      run: make vhs-golden-compare
    - name: Upload diff on failure
      if: failure()
      uses: actions/upload-artifact@v4
      with:
        name: visual-diff
        path: demos/vhs/output/diff/
```

---

## Checklist for New Workflows

- [ ] Written documentation created (`docs/workflows/your-workflow.md`)
- [ ] State diagram generated (`make generate-diagrams`)
- [ ] Happy path tape created and tested
- [ ] Sad path tape created (error scenarios)
- [ ] Edge cases tape created (cancel, back, empty states)
- [ ] Demos generated (`make vhs-feature FEATURE=your-feature`)
- [ ] README.md updated with workflow entry
- [ ] PR includes demo evidence

---

## Troubleshooting

### "VHS not installed"

```bash
# macOS
brew install vhs

# Linux (download from releases)
curl -fsSL https://github.com/charmbracelet/vhs/releases/download/v0.7.1/vhs_0.7.1_linux_amd64.tar.gz | tar xz
sudo mv vhs /usr/local/bin/
```

### "ttyd not found"

VHS requires ttyd for terminal emulation:

```bash
# macOS
brew install ttyd

# Linux
# Download from https://github.com/tsl0922/ttyd/releases
```

### Screenshots don't match baseline

1. Check `demos/vhs/output/diff/` for visual diff
2. If changes are intentional: `make vhs-golden-update`
3. If changes are bugs: fix the code and re-run

### Demo timing issues

Adjust `Sleep` durations:
- Increase if actions complete before visible
- Decrease if demos are too slow

---

## Related Documentation

- [VHS README](../../demos/vhs/README.md) - VHS setup and quick start
- [Development Workflow](../development/DEVELOPMENT_WORKFLOW.md) - General development process
- [TUI Developer Guide](../TUI_DEVELOPER_GUIDE.md) - Building TUI components
