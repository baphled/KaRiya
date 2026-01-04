# CV Generation Implementation Plan - KaRiya

**Date**: January 4, 2026
**Status**: Planning Phase
**Priority**: High

---

## Executive Summary

This document outlines the comprehensive plan to implement **actual CV generation** in KaRiya, moving from the current placeholder implementation to a production-ready system that:

1. **Intelligently groups career events** by company and project
2. **Generates professional CV sections** (Experience, Skills, Projects, Education, Summary)
3. **Filters and tailors content** based on target role and audience
4. **Exports to multiple formats** (PDF, Word, Markdown, Text, YAML)
5. **Provides customization options** for job application tailoring

---

## Current State Analysis

### What Exists

✅ **Domain Models** (Fully Defined)
- `CareerEvent` - Represents a career milestone with company, project, tags, categories
- `Fact` - Extracted competencies from events with role fit and audience relevance
- `CVView` - In-memory CV representation (not stored in DB)
- `CVSection` - Individual CV sections (experience, skills, summary)
- `CVBullet` - Individual bullet points with traceability to source events/facts
- `CVConfig` - Configuration for CV generation (name, role, audience, filters)

✅ **Service Layer** (Partially Implemented)
- `CVGenerationService` - Orchestrates CV generation workflow
- `BulletGenerator` - Generates bullet points from events and facts
- `SectionBuilder` - Builds CV sections from bullets
- `ExportService` - Exports CV to multiple formats (Text, Markdown, YAML)
- `ConfigManager` - Manages CV generation configurations

✅ **Data Access** (Fully Implemented)
- `Repository` - Event CRUD operations
- `FactRepository` - Fact management
- `BurstRepository` - Career burst detection

✅ **CLI Integration** (Partially Implemented)
- `GenerateCVIntent` - Intent for CV generation workflow
- CV configuration models and components
- CV preview and export dialogs

### What Needs Enhancement

❌ **Company Grouping** - Not implemented
- Events should be grouped by company
- Projects within each company should be highlighted
- Chronological ordering within companies

❌ **Intelligent Bullet Generation** - Placeholder implementation
- Current: Creates minimal CV structure
- Needed: Extract achievements, impact, and metrics from events and facts
- Needed: Filter bullets based on role and audience relevance

❌ **Professional Section Building** - Basic implementation
- Current: Simple text sections
- Needed: Structured sections with proper formatting
- Needed: Skills organized by category
- Needed: Education and certifications
- Needed: Professional summary tailored to role

❌ **Advanced Export Formats** - Limited
- Current: Text, Markdown, YAML only
- Needed: PDF with professional styling
- Needed: Word (.docx) with formatting
- Needed: ATS-friendly plain text

❌ **Customization Options** - Not implemented
- Needed: Job description keyword matching
- Needed: Section reordering
- Needed: Bullet point selection/editing
- Needed: Format/styling preferences

---

## Architecture Overview

### Data Flow for CV Generation

```
┌─────────────────────────────────────────────────────────────────┐
│                     GenerateCVIntent                             │
│  (User selects role, audience, and customization options)       │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                  CVGenerationService                             │
│  (Orchestrates the generation workflow)                          │
└────────────────────┬────────────────────────────────────────────┘
                     │
         ┌───────────┼───────────┐
         ▼           ▼           ▼
    ┌─────────┐ ┌──────────┐ ┌──────────┐
    │ Events  │ │  Facts   │ │ Bursts   │
    │ Repo    │ │ Repo     │ │ Repo     │
    └─────────┘ └──────────┘ └──────────┘
         │           │           │
         └───────────┼───────────┘
                     │
                     ▼
         ┌───────────────────────────┐
         │  DataProcessingService    │
         │  - Group by company       │
         │  - Extract achievements  │
         │  - Calculate metrics     │
         └────────┬──────────────────┘
                  │
                  ▼
         ┌───────────────────────────┐
         │  BulletGenerator          │
         │  - Filter by role/audience│
         │  - Rank by relevance      │
         │  - Generate descriptions  │
         └────────┬──────────────────┘
                  │
                  ▼
         ┌───────────────────────────┐
         │  SectionBuilder           │
         │  - Organize by section    │
         │  - Extract skills         │
         │  - Build summaries        │
         └────────┬──────────────────┘
                  │
                  ▼
         ┌───────────────────────────┐
         │  CVView with Sections     │
         │  - Experience             │
         │  - Skills                 │
         │  - Projects               │
         │  - Education              │
         │  - Summary                │
         └────────┬──────────────────┘
                  │
    ┌─────────────┼─────────────┐
    ▼             ▼             ▼
  ┌────────┐  ┌────────┐  ┌──────────┐
  │ Export │  │Preview │  │Customize │
  │Service │  │Dialog  │  │Dialog    │
  └────────┘  └────────┘  └──────────┘
    │
    ├─────────────────────────────────────┐
    ▼             ▼             ▼          ▼
  ┌─────┐  ┌──────┐  ┌─────┐  ┌────┐
  │PDF  │  │DOCX  │  │MD   │  │TXT │
  └─────┘  └──────┘  └─────┘  └────┘
```

---

## Implementation Plan

### Phase 1: Data Processing Service (High Priority)

**Goal**: Extract and organize raw event data into structured information

**Components to Create**:

1. **DataProcessingService**
   ```go
   type DataProcessingService interface {
       // GroupEventsByCompany organizes events into company-based structure
       GroupEventsByCompany(ctx context.Context, events []*CareerEvent) (map[string]*CompanyGroup, error)

       // ExtractAchievements identifies key achievements from events
       ExtractAchievements(ctx context.Context, event *CareerEvent, facts []*Fact) ([]*Achievement, error)

       // ExtractSkills identifies skills from events and facts
       ExtractSkills(ctx context.Context, events []*CareerEvent, facts []*Fact) (map[string]*SkillCategory, error)

       // CalculateMetrics extracts quantifiable metrics from event descriptions
       CalculateMetrics(ctx context.Context, event *CareerEvent) ([]*Metric, error)
   }
   ```

2. **Data Structures**
   ```go
   type CompanyGroup struct {
       Company      string
       Position     string
       StartDate    time.Time
       EndDate      time.Time
       Projects     []*ProjectGroup
       Achievements []*Achievement
       Skills       []string
   }

   type ProjectGroup struct {
       Name         string
       Description  string
       StartDate    time.Time
       EndDate      time.Time
       Achievements []*Achievement
       Skills       []string
   }

   type Achievement struct {
       Description string
       Metrics     []*Metric
       EventID     string
       FactIDs     []string
       Confidence  float64
   }

   type Metric struct {
       Type        string // "percentage", "count", "currency", etc.
       Value       string
       Description string
   }

   type SkillCategory struct {
       Category string
       Skills   []*Skill
       Metrics  map[string]int // endorsements, projects, etc.
   }

   type Skill struct {
       Name         string
       Level        string // "beginner", "intermediate", "advanced", "expert"
       Projects     int    // number of projects using this skill
       Endorsements int    // from facts
   }
   ```

**Key Algorithms**:

- **Company Grouping**: Aggregate events by company, determine position title and dates
- **Project Extraction**: Identify unique projects within each company group
- **Achievement Extraction**: Extract quantifiable outcomes from event descriptions
- **Metric Parsing**: Use regex/NLP to identify metrics (%, numbers, currency, etc.)
- **Skill Extraction**: Combine event tags with fact-based skills

### Phase 2: Enhanced Bullet Generator (High Priority)

**Goal**: Generate professional, targeted bullet points

**Enhancements to `BulletGenerator`**:

```go
type EnhancedBulletGenerator interface {
    // GenerateBullets creates targeted bullets for a specific role/audience
    GenerateBullets(ctx context.Context,
        events []*CareerEvent,
        facts []*Fact,
        targetRole string,
        audiences []string) ([]*CVBullet, error)

    // FilterByRelevance filters bullets based on role/audience fit
    FilterByRelevance(ctx context.Context,
        bullets []*CVBullet,
        targetRole string,
        audiences []string) ([]*CVBullet, error)

    // RankBullets ranks bullets by relevance and impact
    RankBullets(bullets []*CVBullet) []*CVBullet

    // EnhanceBulletText improves bullet point wording
    EnhanceBulletText(ctx context.Context,
        bullet *CVBullet,
        targetRole string) (string, error)
}
```

**Bullet Generation Algorithm**:

1. For each event:
   - Extract base achievement from event text
   - Identify associated facts
   - Extract metrics if available
   - Build bullet structure: "Action + Context + Result"

2. Filter by role/audience:
   - Check fact.RoleFit matches targetRole
   - Check fact.AudienceRelevance includes audience
   - Remove duplicates and weak bullets

3. Rank by relevance:
   - Score based on role fit (0-1)
   - Score based on audience relevance (0-1)
   - Score based on metrics presence (0-1)
   - Score based on confidence (0-1)

4. Professional wording:
   - Convert weak verbs to strong action verbs
   - Ensure metric-driven language
   - Add context for clarity
   - Remove aspirational language

**Example Bullet Transformations**:

```
Raw Event: "Worked on team standup meeting"
↓
With Metrics: "Led daily standup for 12-person engineering team"
↓
With Impact: "Led daily standup for 12-person engineering team, improving communication efficiency by 25%"
↓
With Role Context (Principal): "Established and led daily standups for 12-person engineering team, driving 25% improvement in cross-team communication efficiency and reducing blocker resolution time from 48 to 24 hours"
```

### Phase 3: Enhanced Section Builder (High Priority)

**Goal**: Create professional, well-structured CV sections

**Enhancements to `SectionBuilder`**:

```go
type EnhancedSectionBuilder interface {
    // BuildSections creates all CV sections from bullets and events
    BuildSections(ctx context.Context,
        bullets []*CVBullet,
        events []*CareerEvent,
        facts []*Fact,
        targetRole string) ([]*CVSection, error)

    // BuildExperienceSection creates the experience section
    BuildExperienceSection(ctx context.Context,
        groupedEvents map[string]*CompanyGroup) (*CVSection, error)

    // BuildSkillsSection creates the skills section
    BuildSkillsSection(ctx context.Context,
        skills map[string]*SkillCategory) (*CVSection, error)

    // BuildProjectsSection creates a projects section
    BuildProjectsSection(ctx context.Context,
        projects []*ProjectGroup) (*CVSection, error)

    // BuildSummarySection creates a professional summary
    BuildSummarySection(ctx context.Context,
        achievements []*Achievement,
        targetRole string) (*CVSection, error)
}
```

**Section Structures**:

1. **Professional Summary** (Top of CV)
   - 2-3 lines tailored to target role
   - Highlight key achievements
   - Include years of experience

2. **Experience Section** (Main content)
   ```
   Company Name | Position Title | Start Date - End Date

   • Achievement bullet 1 with metrics
   • Achievement bullet 2 with metrics
   • Achievement bullet 3 with metrics

   [Projects section within company if relevant]
   ```

3. **Skills Section** (Organized by category)
   ```
   Technical: Java, Go, Python, Kubernetes
   Leadership: Team Management, Strategic Planning, Mentoring
   Product: Product Strategy, Go-to-Market, User Research
   ```

4. **Projects Section** (Optional, if significant projects exist)
   ```
   Project Name | Role | Start Date - End Date

   Description of project and your contribution
   • Key achievement from project
   ```

5. **Education Section** (If available)
   ```
   Degree Name | Institution | Graduation Date
   ```

### Phase 4: Export Service Enhancements (Medium Priority)

**Goal**: Export CVs to professional formats

**New Export Formats**:

1. **PDF Export**
   - Library: `github.com/go-pdf/fpdf` or `github.com/unidoc/unipdf`
   - Features:
     - Professional styling with margins and fonts
     - Company grouping with visual separation
     - Metrics highlighted
     - Consistent formatting

2. **Word (.docx) Export**
   - Library: `github.com/unidoc/unioffice` or `github.com/go-echarts/go-echarts`
   - Features:
     - Editable document
     - Professional formatting
     - Tables for company/project grouping
     - Bullet points with proper indentation

3. **ATS-Friendly Text Export**
   - Flat text format optimized for Applicant Tracking Systems
   - No special characters or formatting
   - Keywords in plain language
   - Easy parsing

**Implementation**:

```go
type AdvancedExportService interface {
    // ExportToPDF generates a professional PDF
    ExportToPDF(ctx context.Context,
        cv *CVView,
        sections []*CVSection,
        bullets map[string][]*CVBullet,
        options *ExportOptions) ([]byte, error)

    // ExportToDocx generates a Word document
    ExportToDocx(ctx context.Context,
        cv *CVView,
        sections []*CVSection,
        bullets map[string][]*CVBullet,
        options *ExportOptions) ([]byte, error)

    // ExportToATSText generates ATS-optimized text
    ExportToATSText(ctx context.Context,
        cv *CVView,
        sections []*CVSection,
        bullets map[string][]*CVBullet) (string, error)
}

type ExportOptions struct {
    Theme          string // "classic", "modern", "minimal"
    FontSize       int    // 10-14
    IncludeMetrics bool
    IncludeLinks   bool
    ColorScheme    string // "black_white", "blue", "green"
}
```

### Phase 5: Customization Features (Medium Priority)

**Goal**: Allow users to tailor CVs for specific job applications

**Customization Options**:

1. **Job Description Matching**
   ```go
   type JobMatcher interface {
       // MatchJobDescription analyzes a job posting
       MatchJobDescription(ctx context.Context,
           jobDescription string) (*JobAnalysis, error)

       // RecommendBullets suggests bullets that match the job
       RecommendBullets(ctx context.Context,
           jobAnalysis *JobAnalysis,
           availableBullets []*CVBullet) ([]*CVBullet, error)
   }

   type JobAnalysis struct {
       RequiredSkills    []string
       PreferredSkills   []string
       KeyResponsibilities []string
       MatchedBullets    []*CVBullet
       MissingKeywords   []string
   }
   ```

2. **Section Customization**
   - Reorder sections (e.g., Skills before Experience)
   - Show/hide sections
   - Customize section titles
   - Adjust bullet point selection per section

3. **Bullet Point Editor**
   ```go
   type BulletCustomizer interface {
       // SelectBullets allows user to select specific bullets
       SelectBullets(bullets []*CVBullet, selected []string) []*CVBullet

       // EditBullet allows user to edit bullet text
       EditBullet(bullet *CVBullet, newText string) *CVBullet

       // RankBullets allows user to reorder bullets
       RankBullets(bullets []*CVBullet, newOrder []int) []*CVBullet
   }
   ```

4. **Format Preferences**
   - Font selection
   - Color scheme
   - Layout style (chronological, functional, hybrid)
   - Spacing and margins

---

## Data Structures

### New Domain Models

```go
// internal/domain/career/cv_generation.go

// CompanyGroup represents work at a specific company
type CompanyGroup struct {
    Company      string
    Position     string
    StartDate    time.Time
    EndDate      time.Time
    Description  string
    Projects     []*ProjectGroup
    Achievements []*Achievement
    Skills       []*Skill
}

// ProjectGroup represents a specific project
type ProjectGroup struct {
    Name         string
    Description  string
    StartDate    time.Time
    EndDate      time.Time
    Role         string
    Achievements []*Achievement
    Skills       []*Skill
}

// Achievement represents a measurable accomplishment
type Achievement struct {
    Description string
    Metrics     []*Metric
    EventID     string
    FactIDs     []string
    Confidence  float64
    ActionVerb  string
}

// Metric represents a quantifiable measure
type Metric struct {
    Type        string // "percentage", "count", "currency", "time", "ratio"
    Value       string
    Unit        string
    Context     string
}

// SkillCategory groups related skills
type SkillCategory struct {
    Name   string
    Skills []*Skill
}

// Skill represents a professional capability
type Skill struct {
    Name         string
    Level        string // "beginner", "intermediate", "advanced", "expert"
    Projects     int
    Endorsements int
}
```

---

## Implementation Timeline

### Week 1: Data Processing Service
- [ ] Create DataProcessingService
- [ ] Implement company grouping algorithm
- [ ] Implement achievement extraction
- [ ] Implement skill extraction
- [ ] Write comprehensive tests
- [ ] Estimated: 40 hours

### Week 2: Enhanced Bullet Generator
- [ ] Enhance BulletGenerator
- [ ] Implement role/audience filtering
- [ ] Implement bullet ranking
- [ ] Implement professional wording
- [ ] Write comprehensive tests
- [ ] Estimated: 35 hours

### Week 3: Enhanced Section Builder
- [ ] Enhance SectionBuilder
- [ ] Implement all section types
- [ ] Implement professional formatting
- [ ] Integrate with bullet generator
- [ ] Write comprehensive tests
- [ ] Estimated: 40 hours

### Week 4: Export Enhancements
- [ ] Implement PDF export
- [ ] Implement Word export
- [ ] Implement ATS text export
- [ ] Add export options
- [ ] Write comprehensive tests
- [ ] Estimated: 45 hours

### Week 5: Customization Features
- [ ] Implement job description matching
- [ ] Implement section customization
- [ ] Implement bullet editor
- [ ] Add format preferences
- [ ] Write comprehensive tests
- [ ] Estimated: 40 hours

### Week 6: Integration & Polish
- [ ] Integrate with GenerateCVIntent
- [ ] Add UI for customization
- [ ] Performance optimization
- [ ] Documentation
- [ ] Final testing
- [ ] Estimated: 30 hours

**Total Estimated Time**: 230 hours (6 weeks at 40 hours/week)

---

## Testing Strategy

### Unit Tests
- Test each component in isolation
- Mock dependencies
- Test edge cases and error handling
- Target: 90%+ coverage per component

### Integration Tests
- Test data flow through pipeline
- Test with real event/fact data
- Test export formats
- Test customization workflows

### Performance Tests
- Benchmark bullet generation
- Benchmark section building
- Benchmark export operations
- Target: < 5 seconds for complete CV generation

### Quality Tests
- Test with various event types
- Test with different roles/audiences
- Manual review of generated CVs
- User acceptance testing

---

## Success Criteria

### Functional Requirements
✅ CV generation based on role and audience
✅ Professional grouping by company and project
✅ Achievement extraction with metrics
✅ Skill extraction and categorization
✅ Professional summary generation
✅ Multiple export formats (PDF, Word, Text, Markdown)
✅ Job description keyword matching
✅ Bullet point customization
✅ Section reordering and customization

### Quality Requirements
✅ 90%+ code coverage
✅ All tests passing
✅ 0 race conditions
✅ < 5 second generation time
✅ Professional output quality
✅ No breaking changes to existing code

### User Experience Requirements
✅ Clear customization workflow
✅ Real-time preview of changes
✅ One-click export
✅ Multiple format options
✅ Intuitive UI

---

## Dependencies

### Required Libraries
- `github.com/go-pdf/fpdf` - PDF generation
- `github.com/unidoc/unioffice` - Word document generation
- Existing: `github.com/charmbracelet/lipgloss` - UI styling
- Existing: `github.com/charmbracelet/bubbles` - UI components

### Optional Libraries
- `github.com/nlopes/slack` - Export to Slack
- `github.com/go-echarts/go-echarts` - Charts for metrics

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Complex NLP for achievement extraction | Medium | Start with regex/heuristics, add ML later |
| Performance with large event sets | Medium | Implement pagination and lazy loading |
| Export format complexity | Low | Use proven libraries (fpdf, unioffice) |
| User customization scope creep | Medium | Implement MVP first, add features based on feedback |
| Maintaining backward compatibility | Low | Extend existing services, don't modify |

---

## Next Steps

1. **Create DataProcessingService** - Start with company grouping
2. **Enhance BulletGenerator** - Implement role/audience filtering
3. **Enhance SectionBuilder** - Create professional sections
4. **Add PDF/Word export** - Integrate export libraries
5. **Add customization UI** - Implement in GenerateCVIntent
6. **Comprehensive testing** - Ensure quality and performance

---

## References

- Existing: `internal/service/career/cv/cv_generation_service.go`
- Existing: `internal/service/career/cv/bullet_generator.go`
- Existing: `internal/service/career/cv/section_builder.go`
- Existing: `internal/service/career/cv/export_service.go`
- Existing: `internal/domain/career/cv.go`
- Existing: `internal/cli/intents/generate_cv_intent.go`


