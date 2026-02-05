# KaRiya VHS Demo Suite

This directory contains VHS tape files for generating demo GIFs, documentation visuals, and golden file tests for KaRiya.

## Prerequisites

Install VHS:

```bash
# macOS
brew install vhs

# Linux (requires ttyd and ffmpeg)
# See: https://github.com/charmbracelet/vhs#installation
```

## Quick Start

```bash
# Generate all demo GIFs
make vhs-demos

# Generate specific demo
make vhs-onboarding
make vhs-capture
make vhs-capture-cancel    # Cancel workflow
make vhs-browse
make vhs-skills
make vhs-configure
make vhs-bursts
make vhs-cv
make vhs-facts
make vhs-import
make vhs-journey
make vhs-burst-detection   # AI feature
make vhs-skill-inference   # AI feature
make vhs-burst-flow        # Complete workflow: events -> burst
make vhs-fact-extraction   # Complete workflow: burst -> facts
make vhs-enrichment        # Full enrichment: events -> bursts -> facts -> skills

# Run visual regression tests
make vhs-golden-compare

# Update golden baselines
make vhs-golden-update
```

## Directory Structure

```
demos/vhs/
├── config.tape              # Shared VHS settings (sourced by all tapes)
├── empty-config.yaml        # Empty profile config (triggers onboarding)
├── minimal-config.yaml      # Pre-populated config (skips onboarding)
├── sample-events.csv        # 40 anonymized career events for demos
├── golden-test.tape         # Generates golden file screenshots
│
├── tapes/                   # Main demo tapes
│   ├── 01-onboarding.tape   # First-time user experience (~30s)
│   ├── 02-capture-event.tape # Capture new event (~25s)
│   ├── 03-browse-timeline.tape # Browse & filter events (~20s)
│   ├── 04-manage-skills.tape # Manage skills (~25s)
│   ├── 05-configure-system.tape # Configure system settings (~30s)
│   ├── 06-manage-bursts.tape # Manage career bursts (~25s)
│   ├── 07-generate-cv.tape  # CV wizard workflow (~40s)
│   ├── 08-manage-facts.tape # Manage extracted facts (~25s)
│   ├── 09-import-workflow.tape # CSV import (~30s)
│   ├── 10-import-to-cv-journey.tape # Complete journey (~90s)
│   ├── 11-burst-detection.tape # AI burst detection (~45s)
│   ├── 12-skill-inference.tape # AI skill inference (~50s)
│   ├── 13-burst-creation-flow.tape # Events -> burst workflow (~120s)
│   ├── 14-capture-cancel.tape # Cancel event capture (~20s)
│   ├── 15-fact-extraction-flow.tape # Fact extraction workflow (~60s)
│   └── 16-enrichment-workflow.tape # Full enrichment workflow (~150s)
│
├── golden/                  # Golden file baselines for visual regression
│   ├── main-menu.png
│   ├── capture-strategy.png
│   └── ...
│
├── scripts/
│   ├── compare-golden.sh    # Compare screenshots against baselines
│   └── update-golden.sh     # Update golden baselines
│
└── output/                  # Generated outputs (gitignored)
    ├── gifs/
    ├── mp4s/
    └── screenshots/
```

## Available Demos

| Demo | Duration | Description |
|------|----------|-------------|
| `01-onboarding` | ~30s | First-time user onboarding wizard |
| `02-capture-event` | ~25s | Capturing a new career event |
| `03-browse-timeline` | ~20s | Browsing and filtering events |
| `04-manage-skills` | ~25s | Managing skills (browsing, filtering, details) |
| `05-configure-system` | ~30s | System configuration (domains, settings) |
| `06-manage-bursts` | ~25s | Career burst management (browsing, details) |
| `07-generate-cv` | ~40s | CV generation wizard |
| `08-manage-facts` | ~25s | Extracted facts management (browsing, editing) |
| `09-import-workflow` | ~30s | CSV import workflow |
| `10-complete-journey` | ~90s | Full journey: import -> browse -> CV |
| `11-burst-detection` | ~45s | AI-powered career burst detection and review |
| `12-skill-inference` | ~50s | AI-powered skill inference from events |
| `13-burst-creation-flow` | ~120s | Complete workflow: create events -> detect burst -> facts |
| `14-capture-cancel` | ~20s | Cancel event capture at various stages |
| `15-fact-extraction-flow` | ~60s | Complete workflow: burst detection -> fact extraction |
| `16-enrichment-workflow` | ~150s | Full enrichment: events -> bursts -> facts -> skills |

## Usage

### Generating Demos

```bash
# All demos
make vhs-demos

# Individual demos
vhs demos/vhs/tapes/01-onboarding.tape
vhs demos/vhs/tapes/07-generate-cv.tape
```

### Visual Regression Testing

Golden file testing compares generated screenshots against known baselines to detect visual regressions.

```bash
# Generate screenshots
vhs demos/vhs/golden-test.tape

# Compare against baselines
./demos/vhs/scripts/compare-golden.sh

# Update baselines (after visual review)
./demos/vhs/scripts/update-golden.sh
```

### PR Documentation

When making UI changes, generate evidence for your PR:

```bash
# Generate demo for your feature
vhs demos/vhs/features/your-feature/demo.tape

# Include in PR description
```

## Config Files

Two config files are provided for different demo scenarios:

| File | Profile Data | Use Case |
|------|--------------|----------|
| `empty-config.yaml` | Empty (no name/email) | Triggers onboarding wizard |
| `minimal-config.yaml` | Pre-populated | Skips onboarding, goes to main menu |

**Example usage in tape files:**

```tape
# For onboarding demos (triggers wizard)
Type `cp demos/vhs/empty-config.yaml /tmp/kariya-demo/config.yaml`

# For feature demos (skips onboarding)
Type `cp demos/vhs/minimal-config.yaml /tmp/kariya-demo/config.yaml`
```

## Sample Data

The `sample-events.csv` file contains 40 anonymized career events spanning:

- Early career (2006-2010): Infrastructure, PHP, gaming
- Mid career (2010-2015): Ruby/Rails, DevOps, agency work
- Senior career (2015-2020): Architecture, consulting
- Recent (2020-2026): Operations, Go, TUI development

This data enables realistic CV generation demos without exposing actual career details.

## Customization

### VHS Settings

Edit `config.tape` to change global settings:

```tape
Set FontSize 16        # Font size
Set Width 1200         # Terminal width
Set Height 700         # Terminal height
Set Theme "Catppuccin Mocha"  # Color theme
Set TypingSpeed 50ms   # Typing animation speed
```

### Creating New Demos

1. Create a new tape file in `tapes/` or `features/`
2. Source the config: `Source demos/vhs/config.tape`
3. Set output paths: `Output demos/vhs/output/gifs/my-demo.gif`
4. Add VHS commands for your workflow
5. Test with: `vhs your-tape.tape`

## CI Integration

Visual regression tests can run in CI:

```yaml
visual-regression:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Install VHS
      run: brew install vhs
    - name: Generate screenshots
      run: vhs demos/vhs/golden-test.tape
    - name: Compare
      run: ./demos/vhs/scripts/compare-golden.sh
    - name: Upload diff on failure
      if: failure()
      uses: actions/upload-artifact@v4
      with:
        name: visual-diff
        path: demos/vhs/output/diff/
```

## Troubleshooting

### "ttyd not found"

VHS requires ttyd. Install it:

```bash
# macOS
brew install ttyd

# Linux
# Download from https://github.com/tsl0922/ttyd/releases
```

### Screenshots not matching

If golden file comparison fails:

1. Review the diff images in `demos/vhs/output/diff/`
2. If changes are intentional, update baselines: `./demos/vhs/scripts/update-golden.sh`
3. If changes are bugs, fix the code and re-run
