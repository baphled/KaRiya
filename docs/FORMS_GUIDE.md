# KaRiya Forms Guide

**Last Updated**: 2026-01-20  
**Status**: Complete - All 7 form configurations documented  
**Library**: [github.com/charmbracelet/huh](https://github.com/charmbracelet/huh) v0.8.0  
**Form Wrappers**: 2 (CaptureForm, SkillForm)  
**Form Configurations**: 7 (Burst, Metadata, Fact, CaptureEvent, BurstSuggestion, Skill, CVConfig)

---

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Empty Field Handling](#empty-field-handling)
4. [Architecture](#architecture)
5. [Form Configurations](#form-configurations)
6. [Validators](#validators)
7. [Helper Functions](#helper-functions)
8. [Integration Patterns](#integration-patterns)
   - [Modal Integration](#modal-integration)
   - [Model Integration](#model-integration)
   - [Form Alignment and Wrapper Pattern](#form-alignment-and-the-wrapper-pattern)
9. [Theming System](#theming-system)
10. [Testing Forms](#testing-forms)
11. [Code References](#code-references)

---

## Overview

KaRiya uses Charm's **`huh`** library for all form handling in the TUI. This provides:

- **Type-safe forms** with compile-time checking
- **Native BubbleTea integration** - forms are `tea.Model` instances
- **Built-in validation** with custom validator support
- **Accessible design** with keyboard navigation
- **Catppuccin theming** for consistent UI
- **Dynamic field visibility** based on user input

### Why Huh?

**Before (Manual Implementation)**:
- 956 lines of manual form handling per form
- Manual focus management (200+ lines)
- Manual validation state tracking
- Manual rendering logic (200+ lines)
- **Total**: ~3,800 lines across 8 form implementations

**After (Huh)**:
- ~50-250 lines per form configuration
- Automatic focus management
- Built-in validation
- Automatic rendering
- **Total**: ~1,600 lines for all 9 forms (**60% reduction**)

### Benefits

✅ **No manual focus management** - huh handles Tab/Shift+Tab automatically  
✅ **Built-in validation** with custom validator support  
✅ **Type-safe field access** via `.Value(&field)` binding  
✅ **Consistent theming** via `GenerateHuhTheme()`  
✅ **Responsive layouts** with height/width support  
✅ **Accessible design** with keyboard navigation  

---

## Quick Start

### 1. Basic Form Creation

```go
import (
    "github.com/charmbracelet/huh"
    "github.com/baphled/kariya/internal/cli/forms"
)

// Create form data structure
type MyFormData struct {
    Name        string
    Description string
}

// Create the form
data := &MyFormData{}
form := forms.NewForm(
    huh.NewGroup(
        forms.NewInput(forms.FieldConfig{
            Key:         "name",
            Title:       "Name",
            Description: "Enter your name",
            Placeholder: "John Doe",
            CharLimit:   100,
            Validate:    forms.Required,
        }).Value(&data.Name),
        
        forms.NewText(forms.FieldConfig{
            Key:         "description",
            Title:       "Description",
            CharLimit:   500,
            Validate:    forms.Description,
        }).Value(&data.Description),
    ),
)
```

### 2. Embed in BubbleTea Model

```go
type MyModel struct {
    form *huh.Form
    data *MyFormData
}

func (m *MyModel) Init() tea.Cmd {
    return m.form.Init()
}

func (m *MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Update form
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    // Check completion
    if forms.IsCompleted(m.form) {
        // Handle form submission
        // data.Name and data.Description are now populated
    }
    
    if forms.IsAborted(m.form) {
        // Handle cancellation
    }
    
    return m, cmd
}

func (m *MyModel) View() string {
    return m.form.View()
}
```

### 3. Common Patterns

```go
// Responsive form (adapts to terminal size)
height := forms.DefaultFormHeight(terminalHeight)
form := forms.NewFormWithHeight(height, groups...)

// Fixed confirm button (always visible)
form := forms.NewFormWithFixedConfirm(fieldsGroup, &confirmValue, width, height)

// Theme-aware form
form := forms.NewThemedForm(myTheme, groups...)

// Check form state
completed := forms.IsCompleted(form)
aborted := forms.IsAborted(form)

// Extract values
name := forms.GetString(form, "name")
tags := forms.GetStrings(form, "tags") // For MultiSelect
```

---

## Empty Field Handling

**CRITICAL**: Forms have specific requirements for handling empty fields. Failure to follow these patterns will cause form crashes or unexpected behavior.

### Quick Checklist

When creating forms, always verify:

- [ ] ✅ **Nil Slices**: Convert `nil` to `[]string{}` in `GetXFormData()` for MultiSelect fields
- [ ] ✅ **Prompt Fix**: Use `forms.NewInput()` helper (never raw `huh.NewInput()`)
- [ ] ✅ **Optional Validators**: Return `nil` for empty strings in optional field validators
- [ ] ✅ **Pointer Fields**: Convert `*string`/`*int` to/from plain types in Get/Apply functions
- [ ] ✅ **MultiSelect Safe**: Use `forms.GetStrings()` helper (never direct type assertion)

### Common Issues

| Issue | Problem | Solution |
|-------|---------|----------|
| **Nil Slices** | MultiSelect panics with nil | Initialize to `[]string{}` in GetXFormData() |
| **Empty Placeholder** | Only shows first character | Use `forms.NewInput()` (sets Prompt automatically) |
| **Optional Validation** | Validator fails on empty | Check `TrimSpace() == ""` first, return nil |
| **Pointer Types** | Form data uses *int/*string | Store as string, convert in Get/Apply |
| **MultiSelect Type** | Panic on type assertion | Use `forms.GetStrings()` helper |

### Code Examples

**Nil Slice Handling**:
```go
// ❌ Wrong - May pass nil to MultiSelect
func GetMetadataFormData(event *career.CareerEvent) *MetadataFormData {
    return &MetadataFormData{
        Tags: event.Tags, // Panic if nil!
    }
}

// ✅ Correct - Always initialize to empty slice
func GetMetadataFormData(event *career.CareerEvent) *MetadataFormData {
    tags := event.Tags
    if tags == nil {
        tags = []string{}
    }
    return &MetadataFormData{
        Tags: tags, // Safe for MultiSelect
    }
}
```

**Prompt Fix for Empty Placeholders**:
```go
// ❌ Wrong - Only shows "E" when empty
huh.NewInput().
    Key("company").
    Placeholder("Enter company...")

// ✅ Correct - Use helper with auto-prompt
forms.NewInput(forms.FieldConfig{
    Key:         "company",
    Placeholder: "Enter company...",
}) // Automatically sets Prompt("> ")
```

**Optional Field Validation**:
```go
// ❌ Wrong - Fails on empty strings
func CompanyName(value string) error {
    if len(value) < 2 {
        return fmt.Errorf("too short")
    }
    return nil
}

// ✅ Correct - Allow empty for optional fields
func CompanyName(value string) error {
    if strings.TrimSpace(value) == "" {
        return nil // Optional field
    }
    if len(value) < 2 {
        return fmt.Errorf("must be at least 2 characters")
    }
    return nil
}
```

**See Complete Guide**: [`docs/rules/FORMS_WORKFLOW_GUIDE.md#handling-empty-fields`](rules/FORMS_WORKFLOW_GUIDE.md#handling-empty-fields) for comprehensive documentation with all 5 issues and full code examples.

---

## Architecture

### Forms Package Structure

The forms package is located in `internal/cli/forms/` with the following files:

| File | Lines | Purpose |
|------|-------|---------|
| `forms.go` | 340 | Core utilities, theme, helpers, field builders |
| `validators.go` | 306 | 25+ validators + date parsing utilities |
| `burst_form.go` | 112 | Burst editor form configuration |
| `metadata_form.go` | 217 | Event metadata form (6 fields) |
| `fact_form.go` | 162 | Fact editor form (4 fields + dropdown) |
| `capture_event_form.go` | 272 | Career event capture (strategy-aware) |
| `burst_suggestion_form.go` | 77 | Burst suggestion editing form |
| `skill_form.go` | 245 | Skill management form (4 fields) |
| `cv_config_form.go` | 185 | CV configuration wizard (3 steps, 8 fields) |

**Total**: ~1,916 lines of source code + ~1,556 lines of tests = **~3,472 lines**

### Theme Integration

The theme system bridges KaRiya themes with huh forms:

- **`internal/cli/themes/huh.go`** (208 lines) - Generates huh-compatible themes
- **`GenerateHuhTheme(theme)`** - Converts KaRiya theme to huh.Theme
- **Complete style coverage**: TextInput, Select, MultiSelect, Confirm, Card, Help

### Standard Form Pattern

All forms follow this consistent pattern:

```go
// 1. Form Data Structure
type XFormData struct {
    Field1          string
    Field2          []string
    SubmitConfirmed bool  // For confirm button
}

// 2. Form Creator Functions
func NewXForm(domainObj *Domain) *huh.Form
func NewXFormWithData(data *XFormData) *huh.Form
func NewXFormWithDataAndDimensions(data *XFormData, width, height int) *huh.Form

// 3. Domain Conversion Functions
func GetXFormData(domainObj *Domain) *XFormData
func ApplyXFormData(domainObj *Domain, data *XFormData) error
```

This pattern provides:
- **Type safety**: Form fields are bound to struct pointers
- **Flexibility**: Multiple constructor variants for different use cases
- **Clean separation**: Domain objects separate from form data
- **Bidirectional conversion**: Easy to extract and apply form data

---

## Form Configurations

### 1. BurstForm

**Purpose**: Edit burst name and description  
**File**: [`internal/cli/forms/burst_form.go`](../internal/cli/forms/burst_form.go)  
**Fields**: 2 (Name, Description)

#### Data Structure

```go
type BurstFormData struct {
    Name            string
    Description     string
    SubmitConfirmed bool
}
```

#### Usage

```go
// From domain object
burst := &career.Burst{
    Name:        "Sprint Planning Improvements",
    Description: "Improved sprint planning process",
}
form := forms.NewBurstEditorForm(burst)

// With custom data
data := &forms.BurstFormData{
    Name:        "My Burst",
    Description: "Description here",
}
form := forms.NewBurstEditorFormWithData(data)

// With dimensions (scrollable)
form := forms.NewBurstEditorFormWithDataAndDimensions(data, width, height)
```

#### Domain Conversion

```go
// Extract form data from domain
data := forms.GetBurstFormData(burst)

// Apply form data to domain
forms.ApplyBurstFormData(burst, data)
```

#### Validators

- **Name**: `forms.Title` (required, 3-200 chars)
- **Description**: `forms.Description` (optional, 10-1000 chars)

---

### 2. MetadataForm

**Purpose**: Edit event metadata (Date, Company, Project, Tags, Categories, Skills)  
**File**: [`internal/cli/forms/metadata_form.go`](../internal/cli/forms/metadata_form.go)  
**Fields**: 6 (Date, Company, Project, Tags, Categories, Skills)

#### Data Structure

```go
type MetadataFormData struct {
    Date            string
    Company         string
    Project         string
    Tags            []string
    Categories      []string
    Skills          []string
    SubmitConfirmed bool
}
```

#### Usage

```go
// From domain object
event := &career.CareerEvent{
    Date:       time.Now(),
    Company:    "Acme Corp",
    Project:    "Project Alpha",
    Tags:       []string{"backend", "api"},
    Categories: []string{"development"},
    Skills:     []string{"go", "postgresql"},
}

// Load available options
availableTags := []string{"backend", "frontend", "devops"}
availableCategories := []string{"development", "leadership"}
availableSkills := []*career.Skill{
    {ID: "1", Name: "Go"},
    {ID: "2", Name: "PostgreSQL"},
}

form := forms.NewMetadataEditorForm(
    event,
    availableTags,
    availableCategories,
    availableSkills,
)
```

#### Domain Conversion

```go
// Extract form data
data := forms.GetMetadataFormData(event)

// Apply form data (handles date parsing)
err := forms.ApplyMetadataFormData(event, data)
if err != nil {
    // Handle date parsing error
}
```

#### Date Parsing

The metadata form supports flexible date input:

```go
forms.ParseDateString("2024-01-12")    // Standard YYYY-MM-DD
forms.ParseDateString("today")          // Current date
forms.ParseDateString("7 days ago")     // Relative
forms.ParseDateString("2 weeks ago")    // Relative
forms.ParseDateString("1 month ago")    // Relative
```

#### Validators

- **Date**: `forms.DateFormat` (optional, YYYY-MM-DD or relative)
- **Company**: `forms.CompanyName` (optional, 2-100 chars)
- **Project**: `forms.MaxLength(100)` (optional)
- **Tags**: MultiSelect with limit of 10
- **Categories**: MultiSelect with limit of 5
- **Skills**: MultiSelect with limit of 10

---

### 3. FactForm

**Purpose**: Edit fact text and metadata (CompetencyCategories, RoleFit, AudienceRelevance)  
**File**: [`internal/cli/forms/fact_form.go`](../internal/cli/forms/fact_form.go)  
**Fields**: 4 (Text, CompetencyCategories, RoleFit, AudienceRelevance)

#### Data Structure

```go
type FactFormData struct {
    Text                 string
    CompetencyCategories []string // Multi-select
    RoleFit              string   // Single-select dropdown
    AudienceRelevance    []string // Multi-select
    SubmitConfirmed      bool
}
```

#### Usage

```go
// From domain object
fact := &career.Fact{
    Text:                 "Led migration to microservices",
    CompetencyCategories: []string{"technical", "leadership"},
    RoleFit:              career.RoleFitStaff,
    AudienceRelevance:    []string{"hiring_manager", "peer"},
}

form := forms.NewFactEditorForm(fact)

// With dimensions
data := forms.GetFactFormData(fact)
form := forms.NewFactEditorFormWithDataAndDimensions(data, width, height)
```

#### Domain Conversion

```go
// Extract form data
data := forms.GetFactFormData(fact)

// Apply form data
err := forms.ApplyFactFormData(fact, data)
```

#### Field Options

```go
// Role Fit Options (dropdown)
forms.RoleFitOptions() // Returns:
// - Principal
// - Engineering Manager
// - Staff
// - Senior IC

// Competency Categories (multi-select)
forms.CompetencyCategoryOptions() // Returns:
// - Technical
// - Leadership
// - Product
// - Consulting
// - Research
// - Mentoring

// Audience Relevance (multi-select)
forms.AudienceRelevanceOptions() // Returns:
// - Hiring Manager
// - Recruiter
// - Peer
```

#### Validators

- **Text**: `Compose(Required, MinLength(10), MaxLength(2000))`
- **CompetencyCategories**: MultiSelect with limit of 6
- **RoleFit**: Single select (required)
- **AudienceRelevance**: MultiSelect with limit of 3

---

### 4. CaptureEventForm

**Purpose**: Capture career events with strategy-aware fields  
**File**: [`internal/cli/forms/capture_event_form.go`](../internal/cli/forms/capture_event_form.go)  
**Fields**: 2-6 (depends on strategy)

#### Data Structure

```go
type CaptureEventFormData struct {
    Text            string
    Date            string
    Company         string
    Project         string
    Tags            []string
    Categories      []string
    SubmitConfirmed bool
}
```

#### Usage

```go
// Initialize form data
data := forms.NewCaptureEventFormData()

// Quick capture (text + date only)
form := forms.NewCaptureEventForm(data, "quick", width, height)

// Manual capture (all fields)
form := forms.NewCaptureEventForm(data, "manual", width, height)
```

#### Strategy-Aware Fields

**Quick Mode** (2 fields):
- Text (required, 10-2000 chars)
- Date (optional, defaults to today)

**Manual Mode** (6 fields):
- Text (required, 10-2000 chars)
- Date (optional, defaults to today)
- Company (optional, max 200 chars)
- Project (optional, max 200 chars)
- Tags (multi-select, limit 6, from `career.AllowedTags`)
- Categories (multi-select, limit 6, from `career.AllowedCategories`)

#### Dynamic Form Building

The form rebuilds when strategy changes:

```go
// Rebuild form with new strategy
model.data = forms.NewCaptureEventFormData()
model.form = forms.NewCaptureEventForm(
    model.data,
    newStrategy, // "quick" or "manual"
    width,
    height,
)
```

#### Validators

- **Text**: `Compose(Required, MinLength(10), MaxLength(2000))`
- **Date**: `forms.DateFormat` (optional)
- **Company**: `forms.MaxLength(200)` (optional)
- **Project**: `forms.MaxLength(200)` (optional)

---

### 5. BurstSuggestionForm

**Purpose**: Edit suggested burst name and description before confirming  
**File**: [`internal/cli/forms/burst_suggestion_form.go`](../internal/cli/forms/burst_suggestion_form.go)  
**Fields**: 2 (Name, Description)

#### Data Structure

```go
type BurstSuggestionFormData struct {
    Name        string
    Description string
}
```

#### Usage

```go
// From burst suggestion (service type)
suggestion := burstfact.BurstSuggestion{
    Name:            "Sprint Planning",
    Description:     "Improved sprint planning process",
    EventIDs:        []string{"1", "2", "3"},
    ConfidenceScore: 0.85,
}

form := forms.NewBurstSuggestionEditForm(suggestion)

// With custom data
data := &forms.BurstSuggestionFormData{
    Name:        "Custom Name",
    Description: "Custom Description",
}
form := forms.NewBurstSuggestionEditFormWithData(data)
```

#### Domain Conversion

```go
// Extract form data
data := forms.GetBurstSuggestionFormData(suggestion)

// Apply form data (preserves EventIDs and ConfidenceScore)
forms.ApplyBurstSuggestionFormData(&suggestion, data)
```

#### Field Limits

- **Name**: 100 characters (optional, auto-generated if empty)
- **Description**: 500 characters (optional)

**Note**: EventIDs and ConfidenceScore are preserved during editing and not exposed in the form.

---

### 6. SkillForm

**Purpose**: Add or edit skill information  
**File**: [`internal/cli/forms/skill_form.go`](../internal/cli/forms/skill_form.go)  
**Fields**: 4 (Name, Category, Level, YearsUsed)

#### Data Structure

```go
type SkillFormData struct {
    Name            string
    Category        string
    Level           string
    YearsUsed       string
    SubmitConfirmed bool
}
```

#### Usage

```go
// From domain object
skill := &career.Skill{
    Name:     "Kubernetes",
    Category: "devops",
    Level:    "advanced",
}
form := forms.NewSkillForm(skill)

// With custom data
data := &forms.SkillFormData{
    Name:     "React",
    Category: "frontend",
    Level:    "intermediate",
}
form := forms.NewSkillFormWithData(data)

// With dimensions for modal
form := forms.NewSkillFormWithDataAndDimensions(data, 80, 40)
```

#### Domain Conversion

```go
// Extract form data
data := forms.GetSkillFormData(skill)

// Apply form data
forms.ApplySkillFormData(skill, data)
```

#### Constants

```go
// Common skill categories
var CommonSkillCategories = []string{
    "backend", "frontend", "devops", "database", "cloud", "tooling", "other",
}

// Valid skill levels
var ValidSkillLevels = []string{
    "beginner", "intermediate", "advanced", "expert",
}
```

#### Validators

- **Name**: `SkillName` - Required, 1-100 characters
- **Category**: `SkillCategory` - Required, 1-50 characters  
- **Level**: `SkillLevel` - Optional, must be valid level if provided
- **YearsUsed**: `SkillYearsUsed` - Optional, 0-50 if provided

---

### 7. CVConfigForm

**Purpose**: Configure CV generation with multi-step wizard  
**File**: [`internal/cli/forms/cv_config_form.go`](../internal/cli/forms/cv_config_form.go)  
**Fields**: 8 (ProfileID, Audience, TechFocus, Technologies, FocusArea, SkillsFormat, SkillsLimit, CVLength)  
**Steps**: 3 (WHO, TECH, FORMAT)

#### Data Structure

```go
type CVConfigFormData struct {
    ProfileID       string
    Audience        string
    TechFocus       string
    Technologies    []string
    FocusArea       string
    SkillsFormat    string
    SkillsLimit     int    // max skills per category/total (0 = no limit)
    CVLength        string
    SubmitConfirmed bool
}
```

#### Usage

```go
// Create with profile options and extracted technologies
profileOpts := []forms.ProfileOption{
    {ID: "1", Name: "Software Engineer"},
    {ID: "2", Name: "Tech Lead"},
}
extractedTechs := []forms.ExtractedTechnology{
    {Name: "Go"}, {Name: "Python"}, {Name: "React"},
}

data := &forms.CVConfigFormData{}
form := forms.NewCVConfigForm(data, profileOpts, extractedTechs, 80, 40)
```

#### Wizard Steps

**Step 1 - WHO**: Profile and Audience selection
- Select CV Profile (required)
- Target Audience (Hiring Manager, Recruiter, Peer)

**Step 2 - TECH**: Technology focus
- Technology Focus (Highlight specific tech, Balanced)
- Technologies to highlight (MultiSelect from extracted technologies)
- Focus Area (Backend, Frontend, Full Stack, Infrastructure, Leadership)

**Step 3 - FORMAT**: Output formatting
- Skills Format (Flat list, Grouped by category)
- Skills Limit (5/10/15/20 per category, or no limit)
- CV Length (Concise, Standard, Comprehensive)

#### Supporting Types

```go
// Skills limit presets
type SkillsLimitOption struct {
    Value int
    Label string
}

func SkillsLimitOptions() []SkillsLimitOption // Returns preset options

// Profile option for select field
type ProfileOption struct {
    ID   string
    Name string
}

// Technology extracted from events
type ExtractedTechnology struct {
    Name string
}
```

---

## Validators

The forms package provides 25 validators in [`internal/cli/forms/validators.go`](../internal/cli/forms/validators.go) and [`internal/cli/forms/skill_form.go`](../internal/cli/forms/skill_form.go).

### Generic Validators

| Validator | Purpose | Signature |
|-----------|---------|-----------|
| `Required` | Non-empty value | `func(string) error` |
| `MinLength(n)` | Minimum characters | `func(int) func(string) error` |
| `MaxLength(n)` | Maximum characters | `func(int) func(string) error` |
| `LengthRange(min, max)` | Character range | `func(int, int) func(string) error` |
| `DateFormat` | YYYY-MM-DD or relative | `func(string) error` |
| `DateFormatRequired` | Required date | `func(string) error` |
| `Email` | Email format | `func(string) error` |
| `URL` | HTTP/HTTPS URL | `func(string) error` |
| `AlphaNumeric` | Letters and numbers only | `func(string) error` |
| `NoSpecialChars` | No special characters | `func(string) error` |
| `OneOf(options)` | Value must be in list | `func([]string) func(string) error` |
| `Custom(fn, msg)` | Custom validation | `func(func(string) bool, string) func(string) error` |

### Domain-Specific Validators

| Validator | Purpose | Rules |
|-----------|---------|-------|
| `EventText` | Career event text | Required, 10-2000 chars |
| `EventTextOptional` | Optional event text | 10-2000 chars if provided |
| `CompanyName` | Company name | Optional, 2-100 chars |
| `TagName` | Tag name | 1-50 chars, no special chars |
| `Title` | Generic title | Required, 3-200 chars |
| `Description` | Description | Optional, 10-1000 chars |
| `ProfileName` | CV profile name | Required, 2-100 chars |
| `AudienceName` | CV audience name | Required, 2-100 chars |

### Validator Composition

```go
// Combine multiple validators
customValidator := forms.Compose(
    forms.Required,
    forms.MinLength(5),
    forms.MaxLength(50),
    forms.NoSpecialChars,
)

huh.NewInput().
    Key("username").
    Validate(customValidator)
```

### Custom Validators

```go
// Using forms.Custom
huh.NewInput().
    Key("code").
    Validate(forms.Custom(
        func(val string) bool {
            return len(val) == 6 && regexp.MustCompile(`^\d+$`).MatchString(val)
        },
        "must be exactly 6 digits",
    ))

// Standalone function
func validateEvenNumber(value string) error {
    num, err := strconv.Atoi(value)
    if err != nil {
        return fmt.Errorf("must be a number")
    }
    if num%2 != 0 {
        return fmt.Errorf("must be an even number")
    }
    return nil
}

huh.NewInput().
    Key("even").
    Validate(validateEvenNumber)
```

### Date Parsing

The `ParseDateString` function supports multiple formats:

```go
// Standard format
date, err := forms.ParseDateString("2024-01-12")

// Quick input
date, err := forms.ParseDateString("today")

// Relative dates
date, err := forms.ParseDateString("7 days ago")
date, err := forms.ParseDateString("2 weeks ago")
date, err := forms.ParseDateString("1 month ago")
```

---

## Helper Functions

The forms package provides helper functions in [`internal/cli/forms/forms.go`](../internal/cli/forms/forms.go).

### Form Creation

```go
// Basic form with Catppuccin theme
form := forms.NewForm(groups...)

// Form with height (scrollable)
form := forms.NewFormWithHeight(height, groups...)

// Form with dimensions
form := forms.NewFormWithDimensions(width, height, groups...)

// Theme-aware form
form := forms.NewThemedForm(theme, groups...)

// Theme-aware with height
form := forms.NewThemedFormWithHeight(theme, height, groups...)

// Accessible form
form := forms.NewFormWithAccessible(groups...)

// Theme-aware accessible form
form := forms.NewThemedFormWithAccessible(theme, groups...)

// Fixed confirm button pattern
form := forms.NewFormWithFixedConfirm(fieldsGroup, &confirmValue, width, height)
```

### State Checking

```go
// Check if form completed
completed := forms.IsCompleted(form)  // form.State == huh.StateCompleted

// Check if form aborted/cancelled
aborted := forms.IsAborted(form)      // form.State == huh.StateAborted
```

### Value Extraction

```go
// Get string value
name := forms.GetString(form, "name")

// Get boolean value
confirmed := forms.GetBool(form, "confirm")

// Get integer value
count := forms.GetInt(form, "count")

// Get string slice (for MultiSelect)
tags := forms.GetStrings(form, "tags")
```

**Note**: `GetStrings` handles both `[]string` and `[]interface{}` types from huh's MultiSelect.

### Field Builders

```go
// Input field (single line)
input := forms.NewInput(forms.FieldConfig{
    Key:         "name",
    Title:       "Name",
    Description: "Enter your name",
    Placeholder: "John Doe",
    CharLimit:   100,
    Validate:    forms.Required,
})

// Text area (multi-line)
text := forms.NewText(forms.FieldConfig{
    Key:         "description",
    Title:       "Description",
    CharLimit:   500,
    Validate:    forms.Description,
})

// Select (single choice)
options := []forms.SelectOption{
    {Key: "option1", Value: "Option 1"},
    {Key: "option2", Value: "Option 2"},
}
select := forms.NewSelect("choice", "Choose One", "Description", options)

// MultiSelect (multiple choices)
multiSelect := forms.NewMultiSelect("tags", "Select Tags", "Description", options, 10)

// Confirm (yes/no)
confirm := forms.NewConfirm("confirm", "Are you sure?", "Description", "Yes", "No")
```

**Note**: All input fields use `Prompt("> ")` explicitly to fix a huh library display issue with placeholders.

### Dimension Management

```go
// Calculate form height (reserves 20 lines for overhead)
height := forms.DefaultFormHeight(terminalHeight)
// Overhead: logo (~7), breadcrumbs (~2), footer (~3), modal chrome (~4), padding (~4)
// Minimum: 10 lines

// Calculate fields height (reserves space for confirm button)
fieldsHeight := forms.FieldsHeight(terminalHeight)
// fieldsHeight = formHeight - ConfirmButtonHeight (5 lines)
// Minimum: 5 lines

// Confirm button height constant
forms.ConfirmButtonHeight // = 5
```

### Theme Functions

```go
// Get default Catppuccin theme (deprecated, use ThemedForm)
theme := forms.Theme()

// Get theme-aware huh.Theme
huhTheme := forms.ThemedForm(kaRiyaTheme)
```

### Color Constants

```go
forms.FormColors.Title       // #89B4FA (Catppuccin Blue)
forms.FormColors.Description // #94E2D5 (Catppuccin Teal)
forms.FormColors.Error       // #F38BA8 (Catppuccin Red)
forms.FormColors.Success     // #A6E3A1 (Catppuccin Green)
forms.FormColors.Placeholder // #6C7086 (Catppuccin Overlay0)
```

---

## Integration Patterns

### Modal Integration

Modals use the `ModalEditResult[T]` pattern for inline editing. See [`internal/cli/intents/modals.go`](../internal/cli/intents/modals.go).

#### ModalEditResult Pattern

```go
type ModalEditResult[T any] struct {
    Original T
    Modified T
    Accepted bool
    Changes  map[string]interface{}
}
```

#### Modal Structure

```go
type EditBurstModal struct {
    // Preserved original (never mutated)
    original  *career.Burst
    
    // Working copy
    modified  *career.Burst
    
    // Result
    result    *ModalEditResult[*career.Burst]
    
    // Huh form
    form      *huh.Form
    formData  *forms.BurstFormData
    
    // Dimensions
    width, height int
}
```

#### Complete Modal Implementation

```go
func NewEditBurstModal(burst *career.Burst) *EditBurstModal {
    // Get form data
    formData := forms.GetBurstFormData(burst)
    
    // Create form with dimensions
    form := forms.NewBurstEditorFormWithDataAndDimensions(
        formData,
        defaultWidth-4,  // Leave margin for modal chrome
        forms.DefaultFormHeight(defaultHeight),
    )
    
    return &EditBurstModal{
        original: burst,
        modified: copyBurst(burst),
        form:     form,
        formData: formData,
        width:    defaultWidth,
        height:   defaultHeight,
    }
}

func (m *EditBurstModal) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Update dimensions without losing state
        m.width = msg.Width
        m.height = msg.Height
        m.form = m.form.
            WithHeight(forms.DefaultFormHeight(m.height)).
            WithWidth(m.width - 4)
        return nil
    }
    
    // Update form
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    // Check completion
    if forms.IsCompleted(m.form) {
        // Check if user confirmed (not cancelled on confirm button)
        if m.formData.SubmitConfirmed {
            m.syncModified()
            m.createAcceptedResult()
        } else {
            m.createCancelledResult()
        }
    }
    
    if forms.IsAborted(m.form) {
        m.createCancelledResult()
    }
    
    return cmd
}

func (m *EditBurstModal) View() string {
    return m.form.View()
}

func (m *EditBurstModal) Result() *ModalEditResult[*career.Burst] {
    return m.result
}

func (m *EditBurstModal) IsComplete() bool {
    return m.result != nil
}

// Private methods
func (m *EditBurstModal) syncModified() {
    forms.ApplyBurstFormData(m.modified, m.formData)
}

func (m *EditBurstModal) createAcceptedResult() {
    m.result = &ModalEditResult[*career.Burst]{
        Original: m.original,
        Modified: m.modified,
        Accepted: true,
        Changes:  m.computeChanges(),
    }
}

func (m *EditBurstModal) createCancelledResult() {
    m.result = &ModalEditResult[*career.Burst]{
        Original: m.original,
        Modified: m.original, // Restore original
        Accepted: false,
        Changes:  make(map[string]interface{}),
    }
}
```

#### Using Modal Results

```go
// In parent intent
modal := NewEditBurstModal(burst)

// After modal completes
if modal.IsComplete() {
    result := modal.Result()
    
    if result.Accepted {
        // User confirmed changes
        burst = result.Modified
        
        // Check what changed
        for field, value := range result.Changes {
            log.Printf("Changed %s to %v", field, value)
        }
    } else {
        // User cancelled, original preserved
        burst = result.Original
    }
}
```

#### Modal Interface Methods

All modals implement these interface methods for overlay rendering:

```go
// GetTitle returns the modal title
func (m *EditModal) GetTitle() string {
    return "Edit Burst"
}

// GetContent returns form content without wrapper
func (m *EditModal) GetContent() string {
    if m.result != nil && m.result.Accepted {
        return ""
    }
    return m.form.View()
}

// GetFooter returns navigation instructions
func (m *EditModal) GetFooter() string {
    return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next  |  Shift+Tab: Prev"
}
```

---

### Model Integration

Models embed huh forms for full-screen editing. See examples in [`internal/cli/models/`](../internal/cli/models/).

#### Strategy-Aware Forms (CaptureForm)

**File**: `internal/cli/models/capture_form.go`

```go
type CaptureForm struct {
    form     *huh.Form
    data     *forms.CaptureEventFormData
    strategy string
    width, height int
}

// Rebuild form when strategy changes
func (m *CaptureForm) SetStrategy(newStrategy string) {
    m.strategy = newStrategy
    m.rebuildForm()
}

func (m *CaptureForm) rebuildForm() {
    m.form = forms.NewCaptureEventForm(
        m.data,
        m.strategy,
        m.width,
        m.height,
    )
}

// Handle keyboard shortcuts
func (m *CaptureForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Ctrl+S to submit
        if msg.Type == tea.KeyCtrlS {
            // Mark as completed
            m.data.SubmitConfirmed = true
            return m, tea.Quit
        }
    }
    
    // Standard form update
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    return m, cmd
}
```

#### Inline Edit Mode (BurstSuggestionModelNew)

**File**: `internal/cli/models/burst_suggestion_new.go`

```go
type BurstSuggestionModelNew struct {
    suggestions []burstfact.BurstSuggestion
    
    // Edit mode
    editMode    bool
    editIndex   int
    editForm    *huh.Form
    editData    *forms.BurstSuggestionFormData
    
    // Cache edited values
    editedNames map[int]string
    editedDescs map[int]string
}

func (m *BurstSuggestionModelNew) startEdit(index int) {
    m.editMode = true
    m.editIndex = index
    
    suggestion := m.suggestions[index]
    m.editData = forms.GetBurstSuggestionFormData(suggestion)
    m.editForm = forms.NewBurstSuggestionEditFormWithData(m.editData)
}

func (m *BurstSuggestionModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.editMode {
        // Update form
        form, cmd := m.editForm.Update(msg)
        if f, ok := form.(*huh.Form); ok {
            m.editForm = f
        }
        
        // Check completion
        if forms.IsCompleted(m.editForm) {
            m.saveEdits()
            m.exitEditMode()
        }
        
        if forms.IsAborted(m.editForm) {
            m.exitEditMode()
        }
        
        return m, cmd
    }
    
    // Handle normal navigation
    // ...
}

func (m *BurstSuggestionModelNew) View() string {
    if m.editMode {
        return m.editForm.View()
    }
    
    // Normal review view
    // ...
}

func (m *BurstSuggestionModelNew) saveEdits() {
    // Cache edited values
    m.editedNames[m.editIndex] = m.editData.Name
    m.editedDescs[m.editIndex] = m.editData.Description
    
    // Apply to suggestion
    suggestion := &m.suggestions[m.editIndex]
    forms.ApplyBurstSuggestionFormData(suggestion, m.editData)
}
```

#### Repository Integration (MetadataEditorModelNew)

**File**: `internal/cli/models/metadata_editor_new.go`

```go
type MetadataEditorModelNew struct {
    event    *career.CareerEvent
    original *career.CareerEvent // For revert
    
    form     *huh.Form
    formData *forms.MetadataFormData
    
    // Services
    cliService *service.CLIEventService
    
    // Selectors
    tagSelector      *components.TagSelector
    categorySelector *components.CategorySelector
    skillSelector    *components.SkillSelector
}

func NewMetadataEditorModelNew(
    event *career.CareerEvent,
    cliService *service.CLIEventService,
) *MetadataEditorModelNew {
    // Load available options from domain
    availableTags := []string{"backend", "frontend", "devops"}
    availableCategories := []string{"development", "leadership"}
    availableSkills, _ := cliService.GetAllSkills()
    
    // Create form
    formData := forms.GetMetadataFormData(event)
    form := forms.NewMetadataEditorFormWithData(
        formData,
        availableTags,
        availableCategories,
        availableSkills,
    )
    
    return &MetadataEditorModelNew{
        event:      event,
        original:   copyEvent(event),
        form:       form,
        formData:   formData,
        cliService: cliService,
    }
}

func (m *MetadataEditorModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Update form
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    // Handle completion
    if forms.IsCompleted(m.form) {
        return m.handleFormCompletion()
    }
    
    if forms.IsAborted(m.form) {
        m.Revert()
    }
    
    return m, cmd
}

func (m *MetadataEditorModelNew) handleFormCompletion() (tea.Model, tea.Cmd) {
    // Apply form data to event
    err := forms.ApplyMetadataFormData(m.event, m.formData)
    if err != nil {
        m.err = err
        return m, nil
    }
    
    // Persist to repository
    err = m.cliService.UpdateEventMetadata(m.event)
    if err != nil {
        m.err = err
        return m, nil
    }
    
    m.saved = true
    return m, tea.Quit
}

func (m *MetadataEditorModelNew) Revert() {
    m.event.Date = m.original.Date
    m.event.Company = m.original.Company
    m.event.Project = m.original.Project
    m.event.Tags = m.original.Tags
    m.event.Categories = m.original.Categories
    m.event.Skills = m.original.Skills
}
```

---

### Form Alignment and the Wrapper Pattern

**CRITICAL**: Forms in intents must use wrapper models to ensure proper alignment and responsive sizing. Direct use of `*huh.Form` in intents causes alignment issues.

#### The Problem

When creating forms directly in an intent:

```go
// ❌ Wrong - Form is left-aligned, doesn't respond to window size changes
type MyIntent struct {
    form     *huh.Form
    formData *forms.MyFormData
}

func (i *MyIntent) handleAddNew() tea.Cmd {
    i.formData = &forms.MyFormData{}
    
    width, height := 80, 24
    if termInfo := i.GetTerminalInfo(); termInfo != nil {
        width = termInfo.Width
        height = termInfo.Height
    }
    
    i.form = forms.NewMyFormWithDataAndDimensions(
        i.formData, width-4, forms.DefaultFormHeight(height),
    )
    return i.form.Init()
}

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Manual dimension update - error-prone!
        if i.form != nil {
            i.form = i.form.WithWidth(msg.Width - 4).WithHeight(forms.DefaultFormHeight(msg.Height))
        }
    }
    // ...
}
```

**Issues**:
1. Form dimensions captured at creation time may be stale
2. `WindowSizeMsg` handling scattered across intent code
3. Form doesn't update properly on terminal resize
4. Result: **Form is left-aligned instead of centered**

#### The Solution: Wrapper Models

Create a wrapper model (like `CaptureForm`, `SkillForm`) that:

1. **Starts with sensible defaults** (80x24)
2. **Handles `tea.WindowSizeMsg` internally**
3. **Emits completion messages** for the intent to handle
4. **Updates form dimensions dynamically**

```go
// ✅ Correct - Use a wrapper model
type HuhMyForm struct {
    *BaseStandardModel
    formData *forms.MyFormData
    form     *huh.Form
    width    int
    height   int
}

func NewHuhMyForm() *HuhMyForm {
    m := &HuhMyForm{
        BaseStandardModel: NewBaseStandardModel(),
        formData:          &forms.MyFormData{},
        width:             80,  // Sensible default
        height:            24,
    }
    m.rebuildForm()
    return m
}

func (m *HuhMyForm) rebuildForm() {
    m.form = forms.NewMyFormWithDataAndDimensions(
        m.formData,
        m.width-4,
        forms.DefaultFormHeight(m.height),
    )
}

func (m *HuhMyForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Handle window size internally
        m.width = msg.Width
        m.height = msg.Height
        m.form = m.form.WithHeight(forms.DefaultFormHeight(m.height)).WithWidth(m.width - 4)
        return m, nil
    }
    
    // Forward to form
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    // Check completion
    if m.form.State == huh.StateCompleted && m.formData.SubmitConfirmed {
        return m, func() tea.Msg {
            return MyFormCompleteMsg{Data: m.formData, Cancelled: false}
        }
    }
    
    if m.form.State == huh.StateAborted {
        return m, func() tea.Msg {
            return MyFormCompleteMsg{Data: nil, Cancelled: true}
        }
    }
    
    return m, cmd
}

func (m *HuhMyForm) View() string {
    return m.form.View()
}

func (m *HuhMyForm) GetFormData() *forms.MyFormData {
    return m.formData
}
```

#### Using the Wrapper in Intents

```go
// ✅ Correct - Intent uses wrapper model
type MyIntent struct {
    myForm *models.HuhMyForm  // Wrapper, not raw *huh.Form
}

func (i *MyIntent) handleAddNew() tea.Cmd {
    i.myForm = models.NewHuhMyForm()
    return i.myForm.Init()
}

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case models.MyFormCompleteMsg:
        return i.handleFormComplete(msg)
        
    case tea.WindowSizeMsg:
        // Update intent's terminal info
        i.UpdateTerminalInfo(termInfo)
        
        // Forward to form wrapper - it handles its own dimensions
        if i.myForm != nil {
            _, cmd := i.myForm.Update(msg)
            return cmd
        }
        return nil
        
    case tea.KeyMsg:
        if i.myForm != nil {
            if msg.Type == tea.KeyEsc {
                return i.handleFormCancel()
            }
            _, cmd := i.myForm.Update(msg)
            return cmd
        }
        // ...
    }
    
    // Forward other messages
    if i.myForm != nil {
        _, cmd := i.myForm.Update(msg)
        return cmd
    }
    
    return nil
}

func (i *MyIntent) renderForm() string {
    if i.myForm == nil {
        return "Form not initialized"
    }
    return i.myForm.View()
}
```

#### Existing Wrapper Models

| Wrapper | Purpose | File |
|---------|---------|------|
| `CaptureForm` | Career event capture | `internal/cli/models/capture_form.go` |
| `SkillForm` | Skill management | `internal/cli/models/skill_form.go` |

#### When to Create a Wrapper

Create a form wrapper when:
- Form is used in an **intent** (not a modal)
- Form needs to **respond to terminal resize**
- Form should be **centered/aligned** with other UI elements
- You want **consistent completion message handling**

Modals typically don't need wrappers because:
- They have fixed dimensions relative to terminal
- The modal container handles positioning
- They use `ModalEditResult[T]` pattern

#### Key Differences: Modal vs Intent Forms

| Aspect | Modal | Intent |
|--------|-------|--------|
| Form location | Overlay/popup | Full screen area |
| Dimension handling | Modal container | Wrapper model |
| Completion | `ModalEditResult[T]` | Custom message type |
| Positioning | Centered by modal | Wrapper ensures alignment |
| Window resize | Modal updates | Wrapper handles |

#### Troubleshooting Form Alignment

If forms appear left-aligned:

1. **Check for direct `*huh.Form` usage** - Should use wrapper model
2. **Verify `WindowSizeMsg` forwarding** - Wrapper must receive it
3. **Check `Init()` is called** - `i.myForm.Init()` must be returned
4. **Verify wrapper width calculation** - Should use `width-4` for margins

---

## Theming System

The theme system bridges KaRiya themes with huh forms via [`internal/cli/themes/huh.go`](../internal/cli/themes/huh.go) (208 lines).

### GenerateHuhTheme

Converts a KaRiya theme to a complete `huh.Theme`:

```go
// Generate theme from KaRiya theme
huhTheme := themes.GenerateHuhTheme(kaRiyaTheme)

// Create form with theme
form := huh.NewForm(groups...).WithTheme(huhTheme)

// Or use convenience wrapper
form := themes.NewThemedForm(kaRiyaTheme, groups...)
```

### Theme Structure

The generated `huh.Theme` includes:

```go
type huh.Theme struct {
    Form           FormStyles        // Form container
    Group          GroupStyles       // Group container
    FieldSeparator lipgloss.Style   // Spacing between fields
    Focused        FieldStyles       // Active field styling
    Blurred        FieldStyles       // Inactive field styling
    Help           help.Styles       // Help text styling
}
```

### Field Styles (Focused)

```go
FieldStyles{
    Base:               // Border, padding
    Title:              // Bold, primary color
    Description:        // Dim color
    ErrorIndicator:     // Error color
    ErrorMessage:       // Error color
    
    // Select/MultiSelect
    SelectSelector:     // Primary, bold
    Option:             // Normal foreground
    SelectedOption:     // Primary, bold
    SelectedPrefix:     // Success color
    UnselectedOption:   // Normal foreground
    UnselectedPrefix:   // Muted
    NextIndicator:      // Muted
    PrevIndicator:      // Muted
    
    // TextInput
    TextInput: {
        Cursor:         // Primary color
        Placeholder:   // Muted
        Prompt:        // Primary
        Text:          // Normal foreground
    }
    
    // Confirm buttons
    FocusedButton:     // Primary bg, bold
    BlurredButton:     // Alt bg, dim
    
    // Card
    Card:              // Rounded border
    NoteTitle:         // Primary, bold
    Next:              // Primary
}
```

### Field Styles (Blurred)

Same structure as focused, but with dimmed colors for inactive fields.

### Default Theme (Catppuccin)

All forms use Catppuccin theme by default:

```go
form := forms.NewForm(groups...)  // Catppuccin theme applied
```

Color constants:

```go
forms.FormColors.Title       // #89B4FA (Blue)
forms.FormColors.Description // #94E2D5 (Teal)
forms.FormColors.Error       // #F38BA8 (Red)
forms.FormColors.Success     // #A6E3A1 (Green)
forms.FormColors.Placeholder // #6C7086 (Overlay0)
```

### Custom Themes

Override the default theme:

```go
// Use built-in huh theme
form := huh.NewForm(groups...).WithTheme(huh.ThemeBase16())

// Use KaRiya theme
customTheme := themes.NewCustomTheme(...)
form := forms.NewThemedForm(customTheme, groups...)
```

### Accessibility

Enable accessibility mode:

```go
form := forms.NewFormWithAccessible(groups...)

// Or with custom theme
form := forms.NewThemedFormWithAccessible(customTheme, groups...)
```

---

## Testing Forms

The forms package has comprehensive tests in [`internal/cli/forms/*_test.go`](../internal/cli/forms/).

### Test Coverage

- **76+ tests** across all form files
- **100% pass rate**
- Tests for validators, form creation, data conversion, and edge cases

### Unit Testing Form Creation

```go
var _ = Describe("BurstForm", func() {
    var testBurst *career.Burst
    
    BeforeEach(func() {
        testBurst = &career.Burst{
            ID:          "burst-123",
            Name:        "Test Burst",
            Description: "Test description",
        }
    })
    
    It("should create form with burst data", func() {
        form := forms.NewBurstEditorForm(testBurst)
        
        Expect(form).NotTo(BeNil())
        Expect(form.State).To(Equal(huh.StateNormal))
    })
    
    It("should bind data correctly", func() {
        data := &forms.BurstFormData{
            Name:        "Bound Name",
            Description: "Bound Description",
        }
        
        form := forms.NewBurstEditorFormWithData(data)
        
        // Simulate user updating the data
        data.Name = "Updated Name"
        
        // The form has access to updated data via binding
        Expect(data.Name).To(Equal("Updated Name"))
        Expect(form).NotTo(BeNil())
    })
})
```

### Testing Validators

```go
var _ = Describe("Validators", func() {
    Describe("Required", func() {
        It("should pass for non-empty strings", func() {
            err := forms.Required("hello")
            Expect(err).NotTo(HaveOccurred())
        })
        
        It("should fail for empty strings", func() {
            err := forms.Required("")
            Expect(err).To(HaveOccurred())
        })
    })
    
    Describe("DateFormat", func() {
        It("should pass for valid dates", func() {
            err := forms.DateFormat("2024-01-12")
            Expect(err).NotTo(HaveOccurred())
        })
        
        It("should fail for invalid dates", func() {
            err := forms.DateFormat("2024-13-99")
            Expect(err).To(HaveOccurred())
        })
    })
})
```

### Testing Data Conversion

```go
var _ = Describe("Data Conversion", func() {
    It("should extract data from domain", func() {
        burst := &career.Burst{
            Name:        "Test Burst",
            Description: "Test description",
        }
        
        data := forms.GetBurstFormData(burst)
        
        Expect(data.Name).To(Equal("Test Burst"))
        Expect(data.Description).To(Equal("Test description"))
    })
    
    It("should apply data to domain", func() {
        burst := &career.Burst{}
        data := &forms.BurstFormData{
            Name:        "New Name",
            Description: "New Description",
        }
        
        forms.ApplyBurstFormData(burst, data)
        
        Expect(burst.Name).To(Equal("New Name"))
        Expect(burst.Description).To(Equal("New Description"))
    })
    
    It("should handle roundtrip conversion", func() {
        original := &career.Burst{
            Name:        "Original Name",
            Description: "Original Description",
        }
        
        // Extract
        data := forms.GetBurstFormData(original)
        
        // Apply to new burst
        copied := &career.Burst{}
        forms.ApplyBurstFormData(copied, data)
        
        Expect(copied.Name).To(Equal(original.Name))
        Expect(copied.Description).To(Equal(original.Description))
    })
})
```

### Testing Modal Integration

```go
var _ = Describe("EditBurstModal", func() {
    It("should implement overlay interface", func() {
        modal := NewEditBurstModal(burst)
        
        // Interface methods
        Expect(modal.GetTitle()).To(Equal("Edit Burst"))
        Expect(modal.GetContent()).NotTo(BeEmpty())
        Expect(modal.GetFooter()).To(ContainSubstring("Enter: Confirm"))
    })
    
    It("should preserve original on cancel", func() {
        original := &career.Burst{Name: "Original"}
        modal := NewEditBurstModal(original)
        
        // Simulate cancellation
        modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
        
        result := modal.Result()
        Expect(result.Accepted).To(BeFalse())
        Expect(result.Modified.Name).To(Equal("Original"))
    })
})
```

### Testing Dimension Handling

```go
var _ = Describe("Dimension Management", func() {
    It("should calculate form height with overhead", func() {
        height := forms.DefaultFormHeight(40)
        Expect(height).To(Equal(20)) // 40 - 20 overhead
    })
    
    It("should return minimum height for small terminals", func() {
        height := forms.DefaultFormHeight(25)
        Expect(height).To(Equal(10)) // Minimum is 10
    })
    
    It("should calculate fields height", func() {
        fieldsHeight := forms.FieldsHeight(50)
        expected := 50 - 20 - 5 // total - overhead - confirm button
        Expect(fieldsHeight).To(Equal(expected))
    })
})
```

### Testing Edge Cases

```go
var _ = Describe("Edge Cases", func() {
    It("should handle empty form data", func() {
        data := &forms.BurstFormData{}
        form := forms.NewBurstEditorFormWithData(data)
        
        Expect(form).NotTo(BeNil())
        view := form.View()
        Expect(view).NotTo(BeEmpty())
    })
    
    It("should handle nil slices in MultiSelect", func() {
        event := &career.CareerEvent{
            Tags:       nil,
            Categories: nil,
        }
        
        data := forms.GetMetadataFormData(event)
        
        // Should initialize to empty slices
        Expect(data.Tags).NotTo(BeNil())
        Expect(data.Categories).NotTo(BeNil())
    })
    

})
```

### Test Helper Pattern

```go
// SetTestResult allows tests to bypass huh form interaction
func (m *EditBurstModal) SetTestResult(result *ModalEditResult[*career.Burst]) {
    m.result = result
}

// In tests
modal := NewEditBurstModal(burst)
modal.SetTestResult(&ModalEditResult[*career.Burst]{
    Original: burst,
    Modified: modifiedBurst,
    Accepted: true,
    Changes:  map[string]interface{}{"name": "New Name"},
})
```

---

## Code References

### Core Files

- **Forms Package**: [`internal/cli/forms/`](../internal/cli/forms/)
  - [`forms.go`](../internal/cli/forms/forms.go) - Core utilities (340 lines)
  - [`validators.go`](../internal/cli/forms/validators.go) - Validators (307 lines)
  - [`burst_form.go`](../internal/cli/forms/burst_form.go) - Burst editor (113 lines)
  - [`metadata_form.go`](../internal/cli/forms/metadata_form.go) - Metadata editor (218 lines)
  - [`fact_form.go`](../internal/cli/forms/fact_form.go) - Fact editor (163 lines)
  - [`capture_event_form.go`](../internal/cli/forms/capture_event_form.go) - Event capture (142 lines)
  - [`burst_suggestion_form.go`](../internal/cli/forms/burst_suggestion_form.go) - Burst suggestion (78 lines)

- **Theme Integration**: [`internal/cli/themes/huh.go`](../internal/cli/themes/huh.go) (208 lines)

- **Modal Integration**: [`internal/cli/intents/modals.go`](../internal/cli/intents/modals.go)
  - `EditMetadataModal` (183 lines)
  - `EditBurstModal` (161 lines)
  - `EditFactModal` (161 lines)

- **Model Integration**: [`internal/cli/models/`](../internal/cli/models/)
  - [`capture_form.go`](../internal/cli/models/capture_form.go) (165 lines) - Career event capture wrapper
  - [`skill_form.go`](../internal/cli/models/skill_form.go) (123 lines) - Skill management wrapper
  - [`burst_suggestion_new.go`](../internal/cli/models/burst_suggestion_new.go) (544 lines)
  - [`metadata_editor_new.go`](../internal/cli/models/metadata_editor_new.go) (276 lines)
  - [`fact_editor_new.go`](../internal/cli/models/fact_editor_new.go) (243 lines)

### Tests

- **Forms Tests**: [`internal/cli/forms/*_test.go`](../internal/cli/forms/)
  - [`validators_test.go`](../internal/cli/forms/validators_test.go) (50+ tests)
  - [`burst_form_test.go`](../internal/cli/forms/burst_form_test.go) (8 tests)
  - [`metadata_form_test.go`](../internal/cli/forms/metadata_form_test.go) (12 tests)
  - [`fact_form_test.go`](../internal/cli/forms/fact_form_test.go) (6 tests)
  - [`burst_suggestion_form_test.go`](../internal/cli/forms/burst_suggestion_form_test.go) (multiple tests)

- **Modal Tests**: [`internal/cli/intents/modals_test.go`](../internal/cli/intents/modals_test.go) (151 lines)

### External Resources

- **Huh Library**: https://github.com/charmbracelet/huh
- **Huh Examples**: https://github.com/charmbracelet/huh/tree/main/examples
- **BubbleTea Docs**: https://github.com/charmbracelet/bubbletea
- **Catppuccin Theme**: https://github.com/catppuccin/catppuccin

---

## Summary

### Key Takeaways

1. **Use `forms.NewForm()`** to create forms with Catppuccin theme
2. **Use pre-built validators** from `forms` package
3. **Embed `*huh.Form`** in your BubbleTea models
4. **Check state** with `forms.IsCompleted()` and `forms.IsAborted()`
5. **Extract values** with `forms.GetString()`, `forms.GetBool()`, etc.
6. **Use dynamic features** (`WithHideFunc`, `OptionsFunc`) for conditional UI
7. **Follow standard pattern** (FormData, Creator functions, Domain conversion)
8. **Test comprehensively** with Ginkgo/Gomega

### Benefits Achieved

- ✅ **75% less code** (3,800 → 960 lines)
- ✅ **Consistent UX** across all forms
- ✅ **Built-in accessibility** and theming
- ✅ **Type-safe field access**
- ✅ **No manual focus/validation management**
- ✅ **100% test coverage** maintained
- ✅ **Zero regressions** in functionality

### Form Statistics

| Metric | Value |
|--------|-------|
| **Total Forms** | 5 configurations |
| **Source Code** | 1,361 lines |
| **Tests** | 729 lines (76+ tests) |
| **Validators** | 20+ (generic + domain-specific) |
| **Theme Integration** | 208 lines |
| **Code Reduction** | 75% (vs manual implementation) |
| **Test Pass Rate** | 100% |

---

*For questions or suggestions about form handling, see the [`internal/cli/forms/`](../internal/cli/forms/) package or open a GitHub issue.*
