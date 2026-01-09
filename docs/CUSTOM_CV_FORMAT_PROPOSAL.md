# CV Structure Selection - Implementation Complete

**Date**: 2026-01-09  
**Status**: ✅ IMPLEMENTED  
**Related**: CV Generation System, Export Service

---

## Executive Summary

KaRiya now supports **CV structure selection** allowing users to choose between **Standard** and **Narrative** CV formats before generation. This feature enables language-agnostic professionals to create CVs that emphasize pragmatic tool selection over language identity.

---

## Implementation Overview

### Key Concepts

| Concept | Description | When Selected |
|---------|-------------|---------------|
| **CV Structure** | How content is organized (sections, emphasis) | Before generation |
| **Export Format** | File output format (text, markdown, yaml) | After preview |

### CV Structures

| Structure | Description | Use Case |
|-----------|-------------|----------|
| `standard` | Traditional CV with Experience, Projects, Skills sections | Most job applications |
| `narrative` | Language-agnostic format with Core Strengths, Technologies, What I Bring | Emphasizing cross-domain experience |

### State Flow

```
Select Profile → Select Audience → Select Structure → Generate → Preview → Export Format
                                        ↑                           ↑
                                  Standard | Narrative         Text | Markdown | YAML
```

---

## Architecture

### CV Structure Type

**File**: `internal/cli/intents/generate_cv.go`

```go
type CVStructure string

const (
    CVStructureStandard  CVStructure = "standard"
    CVStructureNarrative CVStructure = "narrative"
)
```

### Export Service

**File**: `internal/service/career/cv/export_service.go`

```go
// Export exports a CV using the specified structure and format
func (es *ExportService) Export(
    ctx context.Context,
    cv *career.CVView,
    sections []*career.CVSection,
    bullets map[string][]*career.CVBullet,
    structure CVStructure,
    format ExportFormat,
) (string, error)

// ExportWithProfile exports using a custom profile configuration
func (es *ExportService) ExportWithProfile(
    ctx context.Context,
    cv *career.CVView,
    sections []*career.CVSection,
    bullets map[string][]*career.CVBullet,
    structure CVStructure,
    format ExportFormat,
    profileCfg *config.ProfileConfig,
) (string, error)
```

### Profile Configuration

**File**: `internal/config/config.go`

```go
type ProfileConfig struct {
    Name            string   `yaml:"name"`
    Email           string   `yaml:"email"`
    Title           string   `yaml:"title"`
    Location        string   `yaml:"location"`
    GitHub          string   `yaml:"github"`
    Portfolio       string   `yaml:"portfolio"`
    Languages       string   `yaml:"languages"`
    Frontend        string   `yaml:"frontend"`
    Systems         string   `yaml:"systems"`
    CoreStrengths   []string `yaml:"core_strengths"`
    WhatIBring      []string `yaml:"what_i_bring"`
    DefaultRole     string   `yaml:"default_role"`
    DefaultAudience string   `yaml:"default_audience"`
}
```

---

## User Workflow

### 1. Generate CV with Structure Selection

1. Navigate to **Generate CV** from main menu
2. Select target profile (e.g., Senior IC, Staff)
3. Select target audience (e.g., Technical, Executive)
4. **NEW**: Select CV structure:
   - **Standard** - Traditional format with Experience, Skills, Summary
   - **Narrative** - Language-agnostic format emphasizing pragmatic approach
5. Review generated CV preview
6. Export to Text, Markdown, or YAML

### 2. Configure Profile for Narrative CVs

1. Navigate to **Configure System** from main menu
2. Select **Profile** domain
3. Edit fields:
   - Name, Email, Title, Location
   - GitHub URL, Portfolio URL
   - Languages (comma-separated)
   - Frontend technologies (comma-separated)
   - Systems/Infrastructure (comma-separated)
4. Save changes
5. Profile is used when exporting narrative CVs

---

## Narrative CV Format

The narrative structure produces content emphasizing language-agnostic expertise:

```markdown
# Yomi Colledge

**Senior Software Engineer / Technical Consultant**  
Remote (UK)  
Email: yomi@boodah.net  
GitHub: https://github.com/baphled  
Portfolio: http://boodah.net

---

## Summary

Senior, language-agnostic software engineer with 20+ years of experience...

---

## Core Strengths

- Language-agnostic backend and systems engineering  
- System design and architectural ownership  
- Pragmatic problem decomposition  
- Legacy stabilisation and modernisation  
- Product-focused delivery  
- Linux-first operational mindset  

---

## Languages & Technologies

**Languages:** Ruby, Go, PHP, C/C++, JavaScript, Shell  
**Frontend:** Vue.js  
**Systems:** Linux, SQL, APIs, CI/CD, automation  

---

## Selected Experience

### Company Name
*Jan 2020 - Present*

- Achievement with high confidence (>= 0.75)
- Another significant accomplishment

---

## What I Bring

- Languages as tools, not identity  
- Calm handling of complexity  
- Clear thinking under constraints  
- Long-term maintainability focus  

---

**References available on request.**
```

---

## Files Implemented

### New Files

| File | Purpose |
|------|---------|
| `internal/cli/intents/generate_cv_structure_test.go` | Structure selection tests |
| `internal/cli/intents/generate_cv_preview_test.go` | Preview rendering tests |
| `internal/cli/intents/generate_cv_helpers.go` | Narrative preview helpers |
| `internal/service/career/cv/cv_helpers.go` | Export helpers and profile conversion |

### Modified Files

| File | Changes |
|------|---------|
| `internal/cli/intents/generate_cv.go` | CVStructure type, state, model fields, ProfileConfig |
| `internal/cli/intents/generate_cv_intent.go` | State handler, structure-aware preview and export |
| `internal/cli/intents/configure_system.go` | Profile fields (Title, Location, GitHub, etc.) |
| `internal/cli/app/app.go` | Load and pass profile config to GenerateCV |
| `internal/config/config.go` | Narrative profile fields |
| `internal/service/career/cv/export_service.go` | Export() and ExportWithProfile() methods |

---

## Test Coverage

| Category | Tests |
|----------|-------|
| Structure selection state | 22 |
| Preview rendering | 42 |
| Export (structure-aware) | 7 |
| Profile configuration | 11 |
| **Total new tests** | 82 |

---

## YAML Export Note

YAML format always uses the **standard** structure regardless of selection. This is intentional because YAML is a data interchange format, not a presentation format. The narrative structure is designed for human-readable output (Text, Markdown).

---

## Related Documentation

- [Narrative CV Guide](docs/guides/NARRATIVE_CV_GUIDE.md) - Detailed guide for narrative CVs
- [CV Generation Guide](docs/guides/CV_GENERATION_GUIDE.md) - Complete CV generation workflow
- [CV Troubleshooting](docs/guides/CV_TROUBLESHOOTING.md) - Common issues and solutions

---

## Implementation History

| Phase | Description | Commits |
|-------|-------------|---------|
| Phase 1 | Cleanup - Reset branch | Branch reset |
| Phase 2 | CV Structure types and selection state | 3 commits |
| Phase 3 | Structure-aware preview | 2 commits |
| Phase 4 | Structure-aware export | 3 commits |
| Phase 5 | Profile configuration | 1 commit |

---

*Implemented: 2026-01-09*  
*Task Reference: tasks/tasks-23-custom-cv-format.md*
