# CV Variant System - Implementation Complete

**Date**: 2026-01-09  
**Status**: ✅ IMPLEMENTED (Task 24)  
**Related**: CV Generation System, Export Service, Variant Service

---

## Executive Summary

KaRiya now supports a **CV Variant System** with 16 built-in variants combining two dimensions: **Role Emphasis** (what to highlight) and **Length Format** (how much detail). This replaces the simple structure selection with a more powerful two-step process that automatically determines the appropriate CV structure.

---

## Implementation Overview

### Key Concepts

| Concept | Description | When Selected |
|---------|-------------|---------------|
| **Role Emphasis** | What aspect of experience to highlight (4 options) | Step 1 of variant selection |
| **Length Format** | CV density and detail level (4 options) | Step 2 of variant selection |
| **CV Structure** | Output format (4 structures) | Automatically determined |
| **Export Format** | File output format (text, markdown, yaml) | After preview |

### Role Emphases

| Role Emphasis | Primary Categories | Best For |
|---------------|-------------------|----------|
| `senior_backend` | technical, architecture | Backend-focused engineering roles |
| `staff_principal` | leadership, strategy, architecture | Senior IC and technical leadership |
| `consulting` | strategy, delivery, consulting | Consulting and client-facing roles |
| `language_agnostic` | technical, architecture | Polyglot developers, adaptability focus |

### Length Formats

| Length | Max Years | Max Companies | Target Pages |
|--------|-----------|---------------|--------------|
| `full` | Unlimited | Unlimited | 3+ |
| `standard` | 10 | Unlimited | 2-3 |
| `short` | 5 | 5 | 1-2 |
| `ultra_short` | 3 | 3 | 1 |

### CV Structures (4 Total)

| Structure | Description | Used By Variants |
|-----------|-------------|------------------|
| `standard` | Traditional CV with Experience, Projects, Skills | senior_backend, staff_principal |
| `narrative` | Language-agnostic with Core Strengths, Technologies, What I Bring | language_agnostic |
| `consulting` | Client-focused with Client Engagements, Technical Capabilities | consulting |
| `highlights` | One-page executive summary with Key Capabilities | All ultra_short variants |

### State Flow

```
Select Profile → Select Audience → Select Role Emphasis → Select Length → Generate → Preview → Export
                                          ↑                    ↑                          ↑
                                   4 role options        4 length options          Text | MD | YAML
```

---

## Architecture

### Variant Types

**File**: `internal/service/career/cv/variants.go`

```go
type RoleEmphasis string
const (
    RoleEmphasisSeniorBackend    RoleEmphasis = "senior_backend"
    RoleEmphasisStaffPrincipal   RoleEmphasis = "staff_principal"
    RoleEmphasisConsulting       RoleEmphasis = "consulting"
    RoleEmphasisLanguageAgnostic RoleEmphasis = "language_agnostic"
)

type LengthFormat string
const (
    LengthFull       LengthFormat = "full"
    LengthStandard   LengthFormat = "standard"
    LengthShort      LengthFormat = "short"
    LengthUltraShort LengthFormat = "ultra_short"
)

type CVVariant struct {
    ID              string
    Name            string
    Description     string
    RoleEmphasis    RoleEmphasis
    LengthFormat    LengthFormat
    BaseStructure   CVStructure
    BulletConfig    BulletConfig
    ProfileOverride *ProfileOverride
    IsBuiltIn       bool
}
```

### Variant Service

**File**: `internal/service/career/cv/variants.go`

```go
type VariantService interface {
    ListVariants() []*CVVariant
    GetVariant(id string) (*CVVariant, error)
    GetVariantByDimensions(role RoleEmphasis, length LengthFormat) (*CVVariant, error)
    ListRoleEmphases() []RoleEmphasisInfo
    ListLengthFormats() []LengthFormatInfo
}
```

### Export Service

**File**: `internal/service/career/cv/export_service.go`

```go
// ExportWithProfile exports using a custom profile configuration
// Supports all 4 structures: standard, narrative, consulting, highlights
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

### 1. Generate CV with Variant Selection

1. Navigate to **Generate CV** from main menu
2. Select target profile (e.g., Senior IC, Staff)
3. Select target audience (e.g., Hiring Manager, Recruiter, Peer)
4. **Select Role Emphasis**:
   - **Senior Backend** - Technical depth, architecture focus
   - **Staff/Principal** - Leadership, strategy, cross-team impact
   - **Consulting** - Client engagements, delivery focus
   - **Language-Agnostic** - Adaptability, multi-language expertise
5. **Select Length Format**:
   - **Full (3+ pages)** - Complete history
   - **Standard (2-3 pages)** - Last 10 years
   - **Short (1-2 pages)** - Last 5 years
   - **Ultra-Short (1 page)** - Highlights only
6. Review generated CV preview
7. Export to Text, Markdown, or YAML

### 2. Configure Profile for Custom CVs

1. Navigate to **Configure System** from main menu
2. Select **Profile** domain
3. Edit fields:
   - Name, Email, Title, Location
   - GitHub URL, Portfolio URL
   - Languages (comma-separated)
   - Frontend technologies (comma-separated)
   - Systems/Infrastructure (comma-separated)
   - Core Strengths (for narrative/consulting)
   - What I Bring (value propositions)
4. Save changes
5. Profile is used when exporting CVs (with variant ProfileOverride applied)

---

## CV Structure Formats

### Narrative CV Format

The narrative structure (used by language_agnostic variants) produces content emphasizing language-agnostic expertise:

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

### Consulting CV Format

The consulting structure (used by consulting variants) emphasizes client engagements:

```markdown
# Jane Smith

**Senior Consulting Engineer**  
Remote (UK)  
Email: jane@example.com

---

## Summary

Consulting engineer with 10+ years delivering technical solutions across diverse client environments...

---

## Client Engagements

### Acme Corp
*Jan 2022 - Present*

- Led technical assessment and modernization roadmap
- Delivered microservices architecture reducing deployment time by 80%

### TechStartup Inc
*Jun 2021 - Dec 2021*

- Rapid assessment of legacy codebase
- Implemented CI/CD pipeline and automated testing

---

## What I Bring

- Rapid technical assessment and roadmapping
- Clear communication with technical and non-technical stakeholders
- Pragmatic solutions within budget constraints
- Knowledge transfer and team enablement
```

### Highlights CV Format

The highlights structure (used by all ultra_short variants) is a one-page executive summary:

```markdown
# Jane Smith | Senior Software Engineer | Remote (UK) | jane@example.com

---

## Summary

Senior engineer with 15+ years delivering production systems across diverse stacks.

---

## Key Capabilities

- Distributed systems architecture
- Technical leadership and mentorship
- Legacy modernization
- Cross-functional collaboration

---

## Selected Highlights

- Designed distributed caching layer reducing response times by 60%
- Led microservices migration for 50-engineer organization
- Mentored 12 engineers to senior level over 3 years

---

## Technologies

**Languages:** Go, Python, Ruby | **Systems:** Kubernetes, AWS, Terraform
```

---

## Files Implemented

### Core Variant Files

| File | Purpose |
|------|---------|
| `internal/service/career/cv/variants.go` | Variant types, 16 built-in variants, VariantService |
| `internal/service/career/cv/variants_test.go` | Variant tests (31 specs) |
| `internal/service/career/cv/role_emphasis.go` | Role emphasis configuration |
| `internal/service/career/cv/role_emphasis_test.go` | Role emphasis tests (11 specs) |
| `internal/service/career/cv/length_format.go` | Length format configuration |
| `internal/service/career/cv/length_format_test.go` | Length format tests (31 specs) |
| `internal/service/career/cv/structure_test.go` | Consulting/Highlights structure tests (24 specs) |
| `internal/service/career/cv/cv_helpers.go` | Export helpers, ProfileOverride functions |
| `internal/service/career/cv/cv_helpers_test.go` | Helper tests (13 specs) |

### Modified Files

| File | Changes |
|------|---------|
| `internal/cli/intents/generate_cv.go` | Variant types, states, model fields |
| `internal/cli/intents/generate_cv_intent.go` | Variant selection handlers and views |
| `internal/service/career/cv/export_service.go` | 4 structures: standard, narrative, consulting, highlights |
| `internal/service/career/cv/bullet_generator.go` | Audience filtering |
| `internal/service/career/cv/enhanced_bullet_generator.go` | Audience filtering, role emphasis scoring |

---

## Test Coverage

| Category | Tests |
|----------|-------|
| Variant system | 31 |
| Role emphasis | 11 |
| Length format | 31 |
| Consulting/Highlights structures | 24 |
| ProfileOverride helpers | 13 |
| Audience filtering | 15 |
| UI variant selection | ~40 |
| **Total new tests (Task 24)** | ~125 |

---

## YAML Export Note

YAML format always uses the **standard** structure regardless of variant selection. This is intentional because YAML is a data interchange format, not a presentation format. The other structures are designed for human-readable output (Text, Markdown).

---

## The 16 Built-In Variants

| Variant ID | Role Emphasis | Length | Structure |
|------------|---------------|--------|-----------|
| `senior_backend_full` | senior_backend | full | standard |
| `senior_backend_standard` | senior_backend | standard | standard |
| `senior_backend_short` | senior_backend | short | standard |
| `senior_backend_ultra_short` | senior_backend | ultra_short | highlights |
| `staff_principal_full` | staff_principal | full | standard |
| `staff_principal_standard` | staff_principal | standard | standard |
| `staff_principal_short` | staff_principal | short | standard |
| `staff_principal_ultra_short` | staff_principal | ultra_short | highlights |
| `consulting_full` | consulting | full | consulting |
| `consulting_standard` | consulting | standard | consulting |
| `consulting_short` | consulting | short | consulting |
| `consulting_ultra_short` | consulting | ultra_short | highlights |
| `language_agnostic_full` | language_agnostic | full | narrative |
| `language_agnostic_standard` | language_agnostic | standard | narrative |
| `language_agnostic_short` | language_agnostic | short | narrative |
| `language_agnostic_ultra_short` | language_agnostic | ultra_short | highlights |

---

## Related Documentation

- [CV Variants Guide](docs/guides/CV_VARIANTS_GUIDE.md) - Complete guide to variant selection
- [Narrative CV Guide](docs/guides/NARRATIVE_CV_GUIDE.md) - Detailed guide for narrative CVs
- [Consulting CV Guide](docs/guides/CONSULTING_CV_GUIDE.md) - Guide for consulting structure
- [CV Generation Guide](docs/guides/CV_GENERATION_GUIDE.md) - Complete CV generation workflow
- [CV Troubleshooting](docs/guides/CV_TROUBLESHOOTING.md) - Common issues and solutions

---

## Implementation History

| Task | Description | Key Commits |
|------|-------------|-------------|
| Task 23 | CV Structure Selection (Standard, Narrative) | Phase 1-5 |
| Task 24 | CV Variant System (16 variants, 4 structures) | 8 commits |

### Task 24 Key Commits

- `c2b21e7` - Audience filtering in bullet generator
- `eed48a8` - Audience filtering in EnhancedBulletGenerator
- `f2c0ffa` - CV variant types and 16 built-in variants
- `b8074e4` - Consulting and Highlights export structures
- `c5dea6b` - Role emphasis configuration and scoring
- `729d265` - Length format configuration and filtering
- `d43b7ca` - TUI variant-based CV generation workflow
- `b518e33` - ProfileOverride wiring to export

---

*Implemented: 2026-01-09*  
*Task References: tasks/tasks-23-custom-cv-format.md, tasks/tasks-24-flexible-cv-variants.md*
