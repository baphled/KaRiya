# Custom CV Format Implementation Proposal

**Date**: 2026-01-08  
**Status**: Proposal  
**Related**: CV Generation System, Export Service

---

## Executive Summary

This document proposes adding a **language-agnostic custom CV format** to KaRiya's CV generation system. The format is based on a real-world CV showcasing a senior engineer's experience across multiple languages and technologies, emphasizing **pragmatic tool selection** over language identity.

---

## Problem Statement

Currently, KaRiya supports three export formats:
- **Text** (`.txt`) - Plain text with bullet points
- **Markdown** (`.md`) - Structured markdown
- **YAML** (`.yaml`) - Machine-readable data format

These formats are **data-centric** but don't provide a **narrative-focused, human-readable** CV template suitable for language-agnostic professionals who want to emphasize:

1. **Languages as tools, not identity**
2. **Cross-domain experience** (hardware, backend, frontend)
3. **Pragmatic problem-solving** over technology evangelism
4. **Long-term career narrative** over recent achievements only

---

## Proposed Solution

### Add "Custom" Export Format

Introduce a **fourth export format**: `ExportFormatCustom` that generates a markdown CV following this structure:

```markdown
# [Name]

**[Primary Role / Title]**  
[Location]  
Email: [email]  
GitHub: [github_url]  
Portfolio: [portfolio_url]

---

## Summary

[2-3 sentence professional summary emphasizing language-agnostic approach,
systems thinking, and years of experience]

---

## Core Strengths

- [Strength 1]
- [Strength 2]
- ...
- [Strength N]

---

## Languages & Technologies

**Languages:** [Ruby, Go, PHP, etc.]  
**Frontend:** [Vue.js, React, etc.]  
**Systems:** [Linux, SQL, APIs, etc.]  

---

## Selected Experience

### [Role Title] — [Company/Project]  
*[Time Period or "Independent Project"]*

- [Bullet point 1]
- [Bullet point 2]
- ...

---

### [Role Title] — [Company/Project]  
*[Time Period]*

- [Bullet point 1]
- ...

---

## What I Bring

- [Value proposition 1]
- [Value proposition 2]
- [Value proposition 3]
- [Value proposition 4]

---

**References available on request.**
```

---

## Technical Design

### 1. Update Export Format Enum

**File**: `internal/service/career/cv/export_service.go`

```go
const (
    ExportFormatText     ExportFormat = "text"
    ExportFormatMarkdown ExportFormat = "markdown"
    ExportFormatYAML     ExportFormat = "yaml"
    ExportFormatCustom   ExportFormat = "custom"  // NEW
)
```

### 2. Add Custom Export Method

**File**: `internal/service/career/cv/export_service.go`

```go
// ExportToCustom exports a CV to custom narrative format
func (es *ExportService) ExportToCustom(
    ctx context.Context, 
    cv *career.CVView, 
    sections []*career.CVSection, 
    bullets map[string][]*career.CVBullet,
    profile *CustomProfile,  // User profile data
) (string, error)
```

### 3. Add Custom Profile Data Structure

**File**: `internal/domain/career/profile.go` (new file)

```go
// CustomProfile contains user profile data for custom CV format
type CustomProfile struct {
    Name             string   // "Yomi Colledge"
    PrimaryRole      string   // "Senior Software Engineer / Technical Consultant"
    Location         string   // "Remote (UK)"
    Email            string   // "yomi@boodah.net"
    GitHubURL        string   // "https://github.com/baphled"
    PortfolioURL     string   // "http://boodah.net"
    
    // Derived from facts/events
    Languages        []string // ["Ruby", "Go", "PHP", ...]
    FrontendTech     []string // ["Vue.js"]
    SystemsTech      []string // ["Linux", "SQL", "APIs", ...]
    
    // Manual fields
    CoreStrengths    []string // ["Language-agnostic backend...", ...]
    WhatIBring       []string // ["Languages as tools...", ...]
}
```

### 4. Add Profile Management Intent

**File**: `internal/cli/intents/manage_profile.go` (new file)

New intent for creating and editing user profiles:
- `ManageProfileIntent` - CRUD operations for profiles
- Stores profiles in `$HOME/.kariya/profiles/` as YAML
- Referenced during CV generation for "Custom" format

### 5. Update GenerateCV Intent

**File**: `internal/cli/intents/generate_cv_intent.go`

When user selects "Custom" export format:
1. Check if profile exists
2. If not, prompt to create profile (transition to `ManageProfileIntent`)
3. If exists, load profile and pass to `ExportToCustom()`

### 6. Section Mapping Strategy

The custom format requires **narrative grouping** vs. current **chronological grouping**:

| Custom Section | Source Data | Mapping Strategy |
|----------------|-------------|------------------|
| Summary | `CVSection[type=summary]` | Use existing summary prose, customize for language-agnostic angle |
| Core Strengths | `Fact[competency_categories]` + `CVBullet[inclusion_reason]` | Extract top 6 competencies by frequency |
| Languages & Technologies | `Fact[competency_categories]` where category matches "language", "frontend", "systems" | Group by category, list as comma-separated |
| Selected Experience | `CVSection[type=experience]` + `CVSection[type=projects]` | Merge experience and projects, sort by impact/recency |
| What I Bring | `CVBullet` with highest `confidence` + `rank` | Extract top 4 value propositions from bullet analysis |

---

## Implementation Plan

### Phase 1: Core Export Functionality (Week 1)

**Goal**: Implement basic custom format export without profile management

**Tasks**:
1. Add `ExportFormatCustom` constant
2. Implement `ExportToCustom()` method with hardcoded profile
3. Add custom format to export selection UI
4. Write unit tests for custom export
5. Update export service tests

**Acceptance Criteria**:
- Custom format appears in export selection
- Custom export generates valid markdown
- All sections render correctly
- Tests pass

### Phase 2: Profile Management (Week 2)

**Goal**: Add profile CRUD functionality

**Tasks**:
1. Create `internal/domain/career/profile.go` with `CustomProfile` struct
2. Create `ProfileManager` service for YAML persistence
3. Create `ManageProfileIntent` for profile CRUD
4. Add profile selection to CV generation workflow
5. Write integration tests for profile workflow

**Acceptance Criteria**:
- Users can create/edit/delete profiles
- Profiles persist across sessions
- CV generation uses profile data
- Tests pass

### Phase 3: Data Extraction & Enrichment (Week 3)

**Goal**: Automatically derive profile fields from facts/events

**Tasks**:
1. Implement `ProfileEnricher` service
2. Extract languages/technologies from facts
3. Extract core strengths from competencies
4. Extract "What I Bring" from top bullets
5. Write tests for extraction logic

**Acceptance Criteria**:
- Profile fields auto-populate from data
- User can override auto-generated fields
- Extraction accuracy >80% on sample data
- Tests pass

### Phase 4: Polish & Documentation (Week 4)

**Goal**: Production-ready feature with docs

**Tasks**:
1. Add error handling and validation
2. Add help text and keyboard shortcuts
3. Write user guide for custom format
4. Add examples to documentation
5. Performance testing
6. Security review (email/URL validation)

**Acceptance Criteria**:
- All error cases handled gracefully
- User guide complete with examples
- Performance benchmarks pass
- Security checks pass
- Zero regressions

---

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────┐
│ 1. User initiates CV generation                        │
│    (GenerateCVIntent)                                   │
└────────────────────┬────────────────────────────────────┘
                     │
                     v
┌─────────────────────────────────────────────────────────┐
│ 2. Select Profile (if format=custom)                   │
│    - Load existing profile OR                           │
│    - Create new profile (ManageProfileIntent)           │
└────────────────────┬────────────────────────────────────┘
                     │
                     v
┌─────────────────────────────────────────────────────────┐
│ 3. Generate CV (CVGenerationService)                    │
│    - Retrieve events/facts                              │
│    - Generate bullets (EnhancedBulletGenerator)         │
│    - Build sections (SectionBuilder)                    │
└────────────────────┬────────────────────────────────────┘
                     │
                     v
┌─────────────────────────────────────────────────────────┐
│ 4. User selects "Custom" export format                 │
└────────────────────┬────────────────────────────────────┘
                     │
                     v
┌─────────────────────────────────────────────────────────┐
│ 5. ExportService.ExportToCustom()                       │
│    - Load profile                                       │
│    - Map sections to custom format                      │
│    - Generate narrative markdown                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     v
┌─────────────────────────────────────────────────────────┐
│ 6. Save to file OR copy to clipboard                   │
└─────────────────────────────────────────────────────────┘
```

---

## Example Output

### Input Data

**Events**:
- 20 career events spanning 20+ years
- Companies: Grand Union, RWDMag, Interface Radio
- Projects: n-vyro.io (solo project)
- Technologies: Ruby, Go, PHP, C/C++, JavaScript, Vue.js

**Facts**:
- 15 competency categories: "Backend Engineering", "Systems Design", "Linux Administration", etc.
- 8 languages: Ruby, Go, PHP, C/C++, JavaScript, Shell, etc.

**Profile**:
- Name: Yomi Colledge
- Role: Senior Software Engineer / Technical Consultant
- Email: yomi@boodah.net
- GitHub: https://github.com/baphled

### Generated Custom CV (Markdown)

```markdown
# Yomi Colledge

**Senior Software Engineer / Technical Consultant**  
Remote (UK)  
Email: yomi@boodah.net  
GitHub: https://github.com/baphled  
Portfolio: http://boodah.net

---

## Summary

Senior, language-agnostic software engineer with 20+ years of experience delivering production systems across diverse stacks and domains. Strong systems thinker with a proven ability to select and adopt the right language or tooling to solve complex problems pragmatically. Comfortable operating across backend services, infrastructure-adjacent components, and product-facing systems.

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

### Founder / Solo Engineer — n-vyro.io  
*Independent Product Project*

- Designed and built a modular environmental control and monitoring product end-to-end.  
- Architected systems spanning hardware-facing components, backend services, and a web frontend.  
- Selected languages per subsystem responsibility across Ruby, Go, C/C++, and Node/Vue.  
- Demonstrates strong system design, autonomy, and language-agnostic delivery.

---

### Senior Software Engineer / Consultant  
*Multiple Clients & Agencies*

- Delivered backend and full-stack systems for media, broadcast, and consumer platforms.  
- Trusted with senior ownership across unfamiliar or legacy codebases.  
- Regularly engaged to stabilise, extend, or modernise production systems.

---

### Senior PHP Engineer — Grand Union  
*2008 – 2010*

- Built backend services for broadcaster-facing platforms used by national TV networks.  
- Delivered reliable, production-grade systems in fast-paced agency environments.

---

### Early Career — Systems & Web Engineering  
*RWDMag, Interface Radio, KeyOne, others*

- Linux systems administration and automation.  
- Web and backend development across PHP and early dynamic stacks.  
- Foundation of strong operational and systems thinking.

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

## Testing Strategy

### Unit Tests

**File**: `internal/service/career/cv/export_service_test.go`

```go
Describe("ExportToCustom", func() {
    Context("with complete profile and data", func() {
        It("generates valid markdown", func() { ... })
        It("includes all required sections", func() { ... })
        It("formats experience chronologically", func() { ... })
        It("groups technologies correctly", func() { ... })
    })
    
    Context("with minimal data", func() {
        It("handles missing profile gracefully", func() { ... })
        It("handles missing sections gracefully", func() { ... })
    })
    
    Context("edge cases", func() {
        It("sanitizes email addresses", func() { ... })
        It("validates URLs", func() { ... })
        It("handles special characters in name", func() { ... })
    })
})
```

### Integration Tests

**File**: `internal/cli/intents/generate_cv_custom_test.go`

```go
Describe("GenerateCVIntent - Custom Format", func() {
    It("prompts for profile if missing", func() { ... })
    It("loads existing profile", func() { ... })
    It("exports custom format successfully", func() { ... })
})
```

### Manual Testing Scenarios

1. **Happy Path**: Generate CV with complete profile → Export custom format → Verify output
2. **Profile Creation**: Generate CV without profile → Create profile → Export custom format
3. **Profile Update**: Edit existing profile → Regenerate CV → Verify changes
4. **Edge Cases**: Test with minimal data, special characters, long names, etc.

---

## Extensibility Considerations

This implementation is designed to support **future custom formats**:

1. **Template Engine**: Consider using `text/template` for custom format rendering
2. **User-Defined Templates**: Allow users to define their own markdown templates
3. **Format Registry**: Register formats dynamically (plugin architecture)
4. **Profile Presets**: Provide pre-built profiles for common roles (e.g., "Language-Agnostic Engineer", "Frontend Specialist", etc.)

**Future Enhancement**: `ExportFormatCustomTemplate` with user-provided Go templates

---

## Migration Path

This change is **additive** and **non-breaking**:

- Existing export formats (Text, Markdown, YAML) remain unchanged
- Users opt-in to custom format
- No database migrations required (profiles are file-based)
- Backward compatible with existing workflows

---

## Success Metrics

1. **Adoption**: >30% of users try custom format within first month
2. **Usage**: >10% of users create a profile
3. **Quality**: Generated CVs require <20% manual editing
4. **Performance**: Export completes in <500ms
5. **Satisfaction**: User feedback >4/5 stars

---

## Open Questions

1. **Profile Storage**: File-based (YAML) vs. database? → **Decision: File-based for Phase 1**
2. **Multiple Profiles**: Support multiple profiles per user? → **Decision: Yes, allow multiple**
3. **Template Engine**: Use Go templates or custom renderer? → **Decision: Custom renderer for Phase 1**
4. **Profile Enrichment**: Auto-generate vs. manual entry? → **Decision: Auto-generate with manual override**

---

## References

- [CV Generation Service](../internal/service/career/cv/cv_generation_service.go)
- [Export Service](../internal/service/career/cv/export_service.go)
- [CV Domain Models](../internal/domain/career/cv.go)
- [GenerateCV Intent](../internal/cli/intents/generate_cv_intent.go)

---

*Last Updated: 2026-01-08*
