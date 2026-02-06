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
# Generate all numbered demo GIFs
make vhs-demos

# Generate all feature-organized tapes
make vhs-features-all

# Generate all tapes for a specific feature
make vhs-feature-all FEATURE=skills
make vhs-feature-all FEATURE=capture
make vhs-feature-all FEATURE=cv

# Feature-specific shortcuts
make vhs-skills-all      # All skills tapes
make vhs-config-all      # All config tapes
make vhs-cv-all          # All CV tapes
make vhs-facts-all       # All facts tapes
make vhs-import-all      # All import tapes
make vhs-bursts-all      # All bursts tapes
make vhs-timeline-all    # All timeline tapes
make vhs-capture-all     # All capture tapes
make vhs-onboarding-all  # All onboarding tapes

# Individual feature tapes
make vhs-skills-happy    # Skills happy path
make vhs-skills-ai       # AI skill inference
make vhs-cv-happy        # CV wizard happy path
make vhs-import-enrich   # Import with enrichment

# Legacy numbered tapes
make vhs-onboarding      # 01-onboarding.tape
make vhs-capture         # 02-capture-event.tape
make vhs-browse          # 03-browse-timeline.tape
make vhs-skills          # 04-manage-skills.tape
make vhs-configure       # 05-configure-system.tape
make vhs-bursts          # 06-manage-bursts.tape
make vhs-cv              # 07-generate-cv.tape
make vhs-facts           # 08-manage-facts.tape
make vhs-import          # 09-import-workflow.tape
make vhs-journey         # 10-import-to-cv-journey.tape
make vhs-burst-detection # 11-burst-detection.tape
make vhs-skill-inference # 12-skill-inference.tape
make vhs-burst-flow      # 13-burst-creation-flow.tape
make vhs-fact-extraction # 15-fact-extraction-flow.tape
make vhs-enrichment      # 16-enrichment-workflow.tape

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
├── tapes/                   # Demo tapes
│   │
│   │ # Numbered tapes (legacy, quick reference)
│   ├── 01-onboarding.tape
│   ├── 02-capture-event.tape
│   ├── ...
│   ├── 16-enrichment-workflow.tape
│   │
│   │ # Feature directories (organized by feature)
│   ├── onboarding/
│   │   └── happy-path.tape
│   │
│   ├── capture/
│   │   ├── happy-path-manual.tape    # Manual capture with full form
│   │   ├── happy-path-quick.tape     # Quick capture mode
│   │   └── sad-path-cancel.tape      # Cancel at various stages
│   │
│   ├── timeline/
│   │   ├── happy-path.tape           # Browse, filter, view details
│   │   └── edge-case-empty.tape      # Empty timeline
│   │
│   ├── skills/
│   │   ├── happy-path.tape           # Browse, view, filter, sort
│   │   ├── add-skill.tape            # Create skill manually
│   │   ├── ai-inference.tape         # AI skill inference
│   │   └── edge-case-empty.tape      # Empty skills list
│   │
│   ├── config/
│   │   ├── happy-path.tape           # Edit and save settings
│   │   └── sad-path-cancel.tape      # Cancel without saving
│   │
│   ├── bursts/
│   │   ├── happy-path-detect-and-accept.tape  # AI detect + accept
│   │   └── skill-inference-from-burst.tape    # Infer skills from burst
│   │
│   ├── cv/
│   │   ├── happy-path-wizard.tape    # Complete CV wizard
│   │   └── edge-case-no-events.tape  # Warning with no events
│   │
│   ├── facts/
│   │   ├── happy-path.tape           # Browse and view facts
│   │   ├── create-delete.tape        # CRUD operations
│   │   └── edge-case-empty.tape      # Empty facts list
│   │
│   └── import/
│       ├── happy-path.tape           # Basic CSV import
│       └── with-enrichment.tape      # Import with AI enrichment
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

### Numbered Tapes (Legacy)

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

### Feature Directory Tapes

Each feature has its own directory with organized tapes for different scenarios.

#### Onboarding (`tapes/onboarding/`)

| Tape | Description |
|------|-------------|
| `happy-path.tape` | Complete onboarding wizard flow |

#### Capture (`tapes/capture/`)

| Tape | Description |
|------|-------------|
| `happy-path-manual.tape` | Full manual capture with all form fields |
| `happy-path-quick.tape` | Quick capture mode |
| `sad-path-cancel.tape` | Cancel capture at various stages |

#### Timeline (`tapes/timeline/`)

| Tape | Description |
|------|-------------|
| `happy-path.tape` | Browse, filter, view event details |
| `edge-case-empty.tape` | Empty timeline state |

#### Skills (`tapes/skills/`)

| Tape | Description |
|------|-------------|
| `happy-path.tape` | Browse, view, filter, sort skills |
| `add-skill.tape` | Create a skill manually |
| `ai-inference.tape` | AI skill inference from events |
| `edge-case-empty.tape` | Empty skills list |

#### Config (`tapes/config/`)

| Tape | Description |
|------|-------------|
| `happy-path.tape` | Edit and save configuration settings |
| `sad-path-cancel.tape` | Cancel editing without saving |

#### Bursts (`tapes/bursts/`)

| Tape | Description |
|------|-------------|
| `happy-path-detect-and-accept.tape` | AI detect bursts and accept suggestions |
| `skill-inference-from-burst.tape` | Infer skills from burst events |

#### CV (`tapes/cv/`)

| Tape | Description |
|------|-------------|
| `happy-path-wizard.tape` | Complete CV generation wizard |
| `edge-case-no-events.tape` | Warning when no events exist |

#### Facts (`tapes/facts/`)

| Tape | Description |
|------|-------------|
| `happy-path.tape` | Browse and view extracted facts |
| `create-delete.tape` | Create and delete facts |
| `edge-case-empty.tape` | Empty facts list |

#### Import (`tapes/import/`)

| Tape | Description |
|------|-------------|
| `happy-path.tape` | Basic CSV import |
| `with-enrichment.tape` | Import with full AI enrichment pipeline |

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
