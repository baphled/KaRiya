# KaRiya Naming Conventions

**Version**: 1.1
**Last Updated**: 2026-01-20
**Status**: Official Standard

> **UIKit Components**: For UIKit component naming, see [UIKIT_GUIDE.md](../UIKIT_GUIDE.md).
> New components should go in `uikit/` not `components/` (legacy).

This document defines naming conventions for all code in the KaRiya TUI application. Consistent naming improves readability, discoverability, and maintainability.

---

## Table of Contents

1. [General Principles](#general-principles)
2. [Intent Layer](#intent-layer)
3. [Screen Layer](#screen-layer)
4. [Component Layer](#component-layer)
5. [Domain Layer](#domain-layer)
6. [Repository Layer](#repository-layer)
7. [Form Layer](#form-layer)
8. [Message Types](#message-types)
9. [State Constants](#state-constants)
10. [File Naming](#file-naming)
11. [Test Files](#test-files)
12. [Quick Reference](#quick-reference)

---

## General Principles

1. **Use Go idioms**: Follow standard Go naming (PascalCase for exports, camelCase for private)
2. **Be descriptive**: Names should reveal intent
3. **Prefer shorter names in narrow scope**: Local variables can be short; exported types should be clear
4. **Avoid redundancy**: Don't repeat package name in type names
5. **Consistency over creativity**: Match existing patterns

---

## Intent Layer

### Intent Type Names

**Pattern**: `{Workflow}Intent`

| Good | Avoid |
|------|-------|
| `CaptureEventIntent` | `CaptureEventIntentModel` |
| `GenerateCVIntent` | `CVGeneratorIntent` |
| `ManageSkillsIntent` | `SkillsManagementIntent` |
| `BrowseTimelineIntent` | `TimelineBrowserIntent` |

### Intent State Types

**Pattern**: `{ShortWorkflow}State`

```go
// Good: Short, clear prefix
type CVState string
type SkillsState string
type CaptureState string
type ExportState string

// Avoid: Long prefixes
type GenerateCVState string    // Too long
type ManageSkillsState string  // Too long
```

### Intent State Constants

**Pattern**: `{ShortWorkflow}{StateName}`

```go
// Good: Concise, descriptive
const (
    CVSelectProfile     CVState = "select_profile"
    CVSelectAudience    CVState = "select_audience"
    CVGenerating        CVState = "generating"
    CVPreview           CVState = "preview"
    CVConfirm           CVState = "confirm"
)

const (
    SkillsList          SkillsState = "list"
    SkillsDetail        SkillsState = "detail"
    SkillsForm          SkillsState = "form"
    SkillsDelete        SkillsState = "delete"
    SkillsFilter        SkillsState = "filter"
)

// Avoid: Redundant prefixes
const (
    GenerateCVStateSelectProfile GenerateCVState = "select_profile"  // Too long
    SkillsStateList              SkillsState     = "list"           // Redundant "State"
)
```

### Intent Context Types

**Pattern**: `{Workflow}Context`

```go
type CaptureEventContext struct {
    Service       *service.CLIEventService
    CareerService *careerservice.Service
    PreviousEvent *career.CareerEvent  // nil for new
}

type ManageSkillsContext struct {
    Service         *careerservice.Service
    SkillRepository career.SkillRepository
    Ctx             context.Context
}
```

### Intent Result Types

**Pattern**: `{Workflow}Result`

```go
type CaptureEventResult struct {
    Event          *career.CareerEvent
    Bursts         []*career.Burst
    Facts          []*career.Fact
}

type ManageSkillsResult struct {
    Action string         // "created", "updated", "deleted", "cancelled"
    Skill  *career.Skill
}
```

---

## Screen Layer

### Screen Type Names

**Pattern**: `{Domain}{Purpose}` (intent-specific screens)

Screens are named by their **domain** and **purpose**, not by generic patterns:

```go
// Good: Domain-specific, self-documenting
type SkillsList struct { ... }
type SkillDetail struct { ... }
type SkillForm struct { ... }
type SkillDeleteConfirm struct { ... }

type CVProfileSelect struct { ... }
type CVAudienceSelect struct { ... }
type TechnologyFocusSelect struct { ... }
type CVPreview struct { ... }

type TimelineEventList struct { ... }
type TimelineEventDetail struct { ... }

// Avoid: Generic names that don't indicate domain
type ListScreen struct { ... }      // What list?
type SelectScreen struct { ... }    // Select what?
type DetailScreen struct { ... }    // Detail of what?
```

### Base Screen Type Names

**Pattern**: `Base{Pattern}Screen[T]`

Base screens use the `Base` prefix to indicate they're abstract/reusable:

```go
// Good: Clear base types
type BaseScreen struct { ... }
type BaseSelectScreen[T any] struct { ... }
type BaseListScreen[T any] struct { ... }
type BaseFormScreen struct { ... }
type BaseDetailScreen struct { ... }
type BaseConfirmScreen struct { ... }
type BaseProgressScreen struct { ... }
```

### Screen Config Types

**Pattern**: `{ScreenName}Config`

```go
type SkillsListConfig struct {
    Skills      []*career.Skill
    Filters     *SkillFilters
    SortBy      string
    Breadcrumbs []string
}

type CVProfileSelectConfig struct {
    Profiles    []*cv.ProfileConfig
    Default     int
    Breadcrumbs []string
}
```

### Screen Result Types

**Pattern**: `{Action}Result` (shared across screens)

```go
// Good: Action-oriented, reusable
type NavigateResult struct {
    Data interface{}
}

type CancelResult struct {
    Reason string  // "back", "escape", "cancel"
}

type SubmitResult struct {
    Data   interface{}
    Errors []ValidationError
}

type ErrorResult struct {
    Err error
}
```

---

## Component Layer

### Component Type Names

**Pattern**: `{Name}` or `{Name}Model`

Components are simple UI primitives with clear names:

```go
// Good: Simple, clear names
type StandardView struct { ... }
type KeyBadge struct { ... }
type HelpModal struct { ... }
type ASCIILogo struct { ... }
type BreadcrumbBar struct { ... }
type TableListContainer struct { ... }

// Avoid: Overly specific or redundant
type StandardViewComponent struct { ... }  // Redundant "Component"
type KeyBadgeUI struct { ... }             // Redundant "UI"
```

### Modal Type Names

**Pattern**: `{Type}Modal`

```go
type ErrorModal struct { ... }
type LoadingModal struct { ... }
type ProgressModal struct { ... }
type SuccessModal struct { ... }
type WarningModal struct { ... }
type HelpModal struct { ... }
```

---

## Domain Layer

### Domain Model Names

**Pattern**: Singular nouns matching business domain

```go
// Good: Singular, business-aligned
type CareerEvent struct { ... }
type Skill struct { ... }
type Burst struct { ... }
type Fact struct { ... }
type CVView struct { ... }
type CVSection struct { ... }
type CVBullet struct { ... }

// Avoid: Plural or technical names
type CareerEvents struct { ... }  // Plural
type SkillEntity struct { ... }   // Technical suffix
```

### Enum/Constant Types

**Pattern**: `{Concept}` (string type with const values)

```go
// Good: Clear type with const values
type TechnologyFocus string
const (
    TechnologyFocusLanguageAgnostic TechnologyFocus = "language_agnostic"
    TechnologyFocusGeneralist       TechnologyFocus = "generalist"
    TechnologyFocusSpecialist       TechnologyFocus = "specialist"
)

type LengthFormat string
const (
    LengthFull       LengthFormat = "full"
    LengthStandard   LengthFormat = "standard"
    LengthShort      LengthFormat = "short"
    LengthUltraShort LengthFormat = "ultra_short"
)
```

---

## Repository Layer

### Repository Interface Names

**Pattern**: `{Entity}Repository`

```go
type CareerEventRepository interface { ... }
type SkillRepository interface { ... }
type BurstRepository interface { ... }
type FactRepository interface { ... }
```

### Repository Implementation Names

**Pattern**: `{Storage}{Entity}Repository`

```go
type SQLiteCareerEventRepository struct { ... }
type SQLiteSkillRepository struct { ... }
type MemorySkillRepository struct { ... }  // For testing
```

### Filter Types

**Pattern**: `{Entity}Filters`

```go
type SkillFilters struct {
    Category  string
    Level     string
    MinEvents int
    SortBy    string
    SortOrder string
}

type EventFilters struct {
    DateFrom  *time.Time
    DateTo    *time.Time
    Company   string
    Tags      []string
}
```

---

## Form Layer

### Form Data Types

**Pattern**: `{Entity}FormData`

```go
type SkillFormData struct {
    Name            string
    Category        string
    Level           string
    YearsUsed       string
    SubmitConfirmed bool
}

type EventFormData struct {
    Text       string
    Date       string
    Company    string
    Project    string
    Tags       []string
    Categories []string
    Skills     []string
}
```

### Form Factory Functions

**Pattern**: `New{Entity}Form...`

```go
// Basic factory
func NewSkillForm(skill *career.Skill) *huh.Form

// With data
func NewSkillFormWithData(data *SkillFormData) *huh.Form

// With dimensions
func NewSkillFormWithDataAndDimensions(data *SkillFormData, width, height int) *huh.Form
```

### Form Conversion Functions

**Pattern**: `Get{Entity}FormData` / `Apply{Entity}FormData`

```go
func GetSkillFormData(skill *career.Skill) *SkillFormData
func ApplySkillFormData(skill *career.Skill, data *SkillFormData)
```

---

## Message Types

### Async Result Messages

**Pattern**: `{Entity}{Action}Msg`

```go
// Good: Entity + past tense action
type SkillsLoadedMsg struct { Skills []*career.Skill; Err error }
type SkillCreatedMsg struct { Skill *career.Skill; Err error }
type SkillUpdatedMsg struct { Skill *career.Skill; Err error }
type SkillDeletedMsg struct { ID string; Err error }

type CVGeneratedMsg struct { CV *career.CVView; Err error }
type EventSavedMsg struct { Event *career.CareerEvent; Err error }
```

### Form Completion Messages

**Pattern**: `{Entity}FormCompleteMsg`

```go
type SkillFormCompleteMsg struct {
    Data      *forms.SkillFormData
    Cancelled bool
}

type EventFormCompleteMsg struct {
    Data      *forms.EventFormData
    Cancelled bool
}
```

### Request/Action Messages

**Pattern**: `Request{Action}Msg` or `{Action}RequestedMsg`

```go
type RequestBrowseEventMsg struct { EventID string }
type DeleteConfirmedMsg struct {}
type FilterChangedMsg struct { Filters *SkillFilters }
```

---

## State Constants

### State Naming Convention

**Pattern**: `{ShortIntent}{StateName}`

| Intent | State Type | State Constants |
|--------|-----------|-----------------|
| CaptureEvent | `CaptureState` | `CaptureChooseStrategy`, `CaptureForm`, `CaptureReview` |
| GenerateCV | `CVState` | `CVSelectProfile`, `CVSelectAudience`, `CVGenerating` |
| ManageSkills | `SkillsState` | `SkillsList`, `SkillsDetail`, `SkillsForm` |
| BrowseTimeline | `TimelineState` | `TimelineList`, `TimelineDetail` |
| ExportArtifact | `ExportState` | `ExportSelectType`, `ExportSelectFormat` |

### State String Values

Use snake_case for state string values to match database/API conventions:

```go
const (
    SkillsList   SkillsState = "list"
    SkillsDetail SkillsState = "detail"
    SkillsForm   SkillsState = "form"
)
```

---

## File Naming

### Intent Files

| Type | Pattern | Example |
|------|---------|---------|
| Types/Context/Result | `{workflow}.go` | `manage_skills.go` |
| Implementation | `{workflow}_intent.go` | `manage_skills_intent.go` |
| Tests | `{workflow}_test.go` | `manage_skills_test.go` |
| Escape tests | `{workflow}_escape_test.go` | `manage_skills_escape_test.go` |
| Navigation tests | `{workflow}_navigation_test.go` | `manage_skills_navigation_test.go` |

### Screen Files

| Type | Pattern | Example |
|------|---------|---------|
| Base screens | `{pattern}_screen.go` | `select_screen.go`, `list_screen.go` |
| Domain screens | `{purpose}.go` | `list.go`, `detail.go`, `form.go` |
| Tests | `{screen}_test.go` | `list_test.go` |

**Directory Organization**:
```
screens/
  base/
    select_screen.go
  skills/
    list.go      # SkillsList
    detail.go    # SkillDetail
    form.go      # SkillForm
```

### Component Files

| Type | Pattern | Example |
|------|---------|---------|
| Component | `{name}.go` | `standard_view.go`, `key_badge.go` |
| Tests | `{name}_test.go` | `standard_view_test.go` |

### Form Files

| Type | Pattern | Example |
|------|---------|---------|
| Form config | `{entity}_form.go` | `skill_form.go` |
| Tests | `{entity}_form_test.go` | `skill_form_test.go` |

---

## Test Files

### Test File Organization

```
{package}/
├── {file}.go
├── {file}_test.go           # Main unit tests
├── {file}_escape_test.go    # Escape key behavior tests
├── {file}_navigation_test.go # Navigation flow tests
├── {file}_global_keys_test.go # Global key handling tests
└── {file}_integration_test.go # Integration tests (if needed)
```

### Test Suite Names (Ginkgo)

**Pattern**: `Describe("{TypeName}", ...)`

```go
var _ = Describe("ManageSkillsIntent", func() {
    Describe("Init", func() { ... })
    Describe("Update", func() {
        Context("when in list state", func() { ... })
        Context("when in detail state", func() { ... })
    })
    Describe("View", func() { ... })
})
```

---

## Quick Reference

| Concept | Pattern | Example |
|---------|---------|---------|
| **Intent** | `{Workflow}Intent` | `ManageSkillsIntent` |
| **Intent State Type** | `{Short}State` | `SkillsState` |
| **Intent State Value** | `{Short}{Name}` | `SkillsList`, `CVSelectProfile` |
| **Intent Context** | `{Workflow}Context` | `ManageSkillsContext` |
| **Intent Result** | `{Workflow}Result` | `ManageSkillsResult` |
| **Domain Screen** | `{Domain}{Purpose}` | `SkillsList`, `CVProfileSelect` |
| **Base Screen** | `Base{Pattern}Screen[T]` | `BaseListScreen[T]` |
| **Screen Config** | `{Screen}Config` | `SkillsListConfig` |
| **Screen Result** | `{Action}Result` | `NavigateResult`, `CancelResult` |
| **Component** | `{Name}` | `StandardView`, `KeyBadge` |
| **Domain Model** | Singular noun | `Skill`, `CareerEvent` |
| **Repository** | `{Entity}Repository` | `SkillRepository` |
| **Form Data** | `{Entity}FormData` | `SkillFormData` |
| **Async Message** | `{Entity}{Action}Msg` | `SkillCreatedMsg` |
| **Intent File** | `{workflow}_intent.go` | `manage_skills_intent.go` |
| **Screen File** | `{purpose}.go` | `list.go`, `detail.go` |

---

## Examples

### Complete Intent Example

```go
// manage_skills.go - Data types
package intents

type SkillsState string

const (
    SkillsList   SkillsState = "list"
    SkillsDetail SkillsState = "detail"
    SkillsForm   SkillsState = "form"
)

type ManageSkillsContext struct {
    Service         *careerservice.Service
    SkillRepository career.SkillRepository
}

type ManageSkillsResult struct {
    Action string
    Skill  *career.Skill
}

// manage_skills_intent.go - Implementation
type ManageSkillsIntent struct {
    state        SkillsState
    context      *ManageSkillsContext
    result       *IntentResult[*ManageSkillsResult]
    activeScreen screens.Screen
}
```

### Complete Screen Example

```go
// screens/skills/list.go
package skills

type SkillsList struct {
    *base.BaseListScreen[*career.Skill]
    groupedByCategory map[string][]*career.Skill
}

type SkillsListConfig struct {
    Skills      []*career.Skill
    Breadcrumbs []string
}

func NewSkillsList(cfg SkillsListConfig) *SkillsList {
    return &SkillsList{
        BaseListScreen: base.NewListScreen(base.ListConfig[*career.Skill]{
            Items:        cfg.Skills,
            ItemRenderer: renderSkill,
            Breadcrumbs:  cfg.Breadcrumbs,
        }),
        groupedByCategory: groupByCategory(cfg.Skills),
    }
}
```

---

## Migration Notes

When refactoring existing code to follow these conventions:

1. **Update state constants first**: Change `SkillsStateList` → `SkillsList`
2. **Update screen names**: Generic → Domain-specific
3. **Update file names**: Follow new patterns
4. **Update tests**: Match new naming
5. **Update documentation**: References to old names

**Example Migration**:
```go
// Before
type ManageSkillsState string
const ManageSkillsStateList ManageSkillsState = "list"

// After
type SkillsState string
const SkillsList SkillsState = "list"
```
