# Forms Development Workflow Guide

**Purpose**: Step-by-step workflows for creating and integrating forms in KaRiya  
**Audience**: AI agents and developers working with forms  
**Status**: Active - Use for ALL form-related development  

---

## Table of Contents

1. [Quick Decision Tree](#quick-decision-tree)
2. [Creating a New Form Configuration](#creating-a-new-form-configuration)
3. [Adding a Form to a Modal](#adding-a-form-to-a-modal)
4. [Adding a Form to a Model](#adding-a-form-to-a-model)
5. [Creating Custom Validators](#creating-custom-validators)
6. [Handling Empty Fields](#handling-empty-fields)
7. [Integration Checklists](#integration-checklists)
8. [Common Pitfalls](#common-pitfalls)
9. [Testing Strategy](#testing-strategy)

---

## Quick Decision Tree

### When to Create a New Form Configuration?

**Create new form configuration when:**
- ✅ You have a new domain object to edit (e.g., Profile, Audience)
- ✅ You need a reusable form across multiple modals/models
- ✅ The form has complex validation logic
- ✅ The form is used in multiple places

**Use existing form configuration when:**
- ❌ Form is single-use inline edit
- ❌ Form has 1-2 simple fields
- ❌ No validators needed

### Modal vs Model Form Integration?

**Use Modal Integration when:**
- ✅ Inline editing within a larger intent (e.g., metadata edit in capture flow)
- ✅ Need to preserve original on cancel
- ✅ Small, focused form (2-6 fields)
- ✅ Need overlay presentation

**Use Model Integration when:**
- ✅ Full-screen form experience
- ✅ Multi-step form flow
- ✅ Strategy-aware fields (quick vs manual)
- ✅ Direct repository persistence

### Form Builder Variants?

| Function | When to Use |
|----------|-------------|
| `NewForm()` | Basic form, no scrolling needed |
| `NewFormWithHeight()` | Form may exceed terminal height |
| `NewFormWithFixedConfirm()` | Always show submit button (recommended for modals) |
| `NewThemedForm()` | Need custom theme (rare, default is Catppuccin) |
| `NewFormWithDimensions()` | Both width and height constraints |

---

## Creating a New Form Configuration

**Time**: 30-45 minutes  
**Prerequisites**: Domain object exists, validators identified  
**Location**: `internal/cli/forms/your_form.go`

### Step 1: Define FormData Structure

```go
// internal/cli/forms/your_form.go
package forms

import (
    "github.com/baphled/kariya/internal/domain/career"
    "github.com/charmbracelet/huh"
)

// YourFormData holds the form data for your feature.
type YourFormData struct {
    Field1          string   // Simple fields
    Field2          []string // For MultiSelect
    Field3          string   // For Select
    SubmitConfirmed bool     // REQUIRED for confirm button
}
```

**✅ Checklist**:
- [ ] All domain fields represented
- [ ] Slices for MultiSelect fields
- [ ] SubmitConfirmed bool included
- [ ] Documentation comment added

### Step 2: Create Form Builder Functions

Create 3 variants for flexibility:

```go
// NewYourForm creates a form from a domain object.
func NewYourForm(obj *career.YourDomain) *huh.Form {
    data := GetYourFormData(obj)
    return NewYourFormWithData(data)
}

// NewYourFormWithData creates a form with initial data.
// This variant allows external data binding for more control.
func NewYourFormWithData(data *YourFormData) *huh.Form {
    return NewYourFormWithDataAndDimensions(data, 0, 0)
}

// NewYourFormWithDataAndDimensions creates a form with dimensions.
// When height > 0, the form becomes scrollable.
// The confirm button is fixed at the bottom, always visible.
func NewYourFormWithDataAndDimensions(data *YourFormData, width, height int) *huh.Form {
    // Initialize submit confirmation
    data.SubmitConfirmed = false
    
    // Create fields group (scrollable)
    fieldsGroup := huh.NewGroup(
        NewInput(FieldConfig{
            Key:         "field1",
            Title:       "Field 1",
            Description: "Description of field 1",
            Placeholder: "Enter value...",
            CharLimit:   100,
            Validate:    Required, // Use existing validator
        }).Value(&data.Field1),
        
        // For MultiSelect
        huh.NewMultiSelect[string]().
            Key("field2").
            Title("Field 2").
            Description("Select options").
            Options(
                huh.NewOption("Option 1", "option1"),
                huh.NewOption("Option 2", "option2"),
            ).
            Value(&data.Field2).
            Limit(10),
    )
    
    return NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
}
```

**✅ Checklist**:
- [ ] All 3 variants created (basic, with data, with dimensions)
- [ ] `SubmitConfirmed` initialized to `false`
- [ ] Used `NewFormWithFixedConfirm()` for confirm button
- [ ] Used `NewInput()` / `NewText()` helpers
- [ ] Validators applied to all fields
- [ ] `.Value(&data.Field)` binding for all fields

### Step 3: Create Domain Conversion Functions

```go
// GetYourFormData extracts form data from a domain object.
func GetYourFormData(obj *career.YourDomain) *YourFormData {
    // Ensure slices are not nil for MultiSelect binding
    field2 := obj.Field2
    if field2 == nil {
        field2 = []string{}
    }
    
    return &YourFormData{
        Field1: obj.Field1,
        Field2: field2,
        Field3: string(obj.Field3), // Convert enum to string if needed
    }
}

// ApplyYourFormData applies the form data to a domain object.
func ApplyYourFormData(obj *career.YourDomain, data *YourFormData) error {
    obj.Field1 = data.Field1
    obj.Field2 = data.Field2
    obj.Field3 = career.EnumType(data.Field3) // Convert back to enum
    
    return nil // Or return validation errors
}
```

**✅ Checklist**:
- [ ] `GetXFormData()` function created
- [ ] `ApplyXFormData()` function created
- [ ] Nil slices handled (convert to empty slice)
- [ ] Enum conversions handled
- [ ] Error handling for validation

### Step 4: Write Tests

```go
// internal/cli/forms/your_form_test.go
package forms_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("YourForm", func() {
    var testObj *career.YourDomain
    
    BeforeEach(func() {
        testObj = &career.YourDomain{
            Field1: "Test Value",
            Field2: []string{"option1"},
        }
    })
    
    Describe("NewYourForm", func() {
        It("should create form with domain data", func() {
            form := forms.NewYourForm(testObj)
            Expect(form).NotTo(BeNil())
        })
    })
    
    Describe("Data Conversion", func() {
        It("should extract data from domain", func() {
            data := forms.GetYourFormData(testObj)
            Expect(data.Field1).To(Equal("Test Value"))
        })
        
        It("should apply data to domain", func() {
            data := &forms.YourFormData{Field1: "New Value"}
            forms.ApplyYourFormData(testObj, data)
            Expect(testObj.Field1).To(Equal("New Value"))
        })
        
        It("should handle roundtrip conversion", func() {
            data := forms.GetYourFormData(testObj)
            newObj := &career.YourDomain{}
            forms.ApplyYourFormData(newObj, data)
            Expect(newObj.Field1).To(Equal(testObj.Field1))
        })
    })
})
```

**✅ Checklist**:
- [ ] Test file created
- [ ] Form creation tested
- [ ] Data extraction tested
- [ ] Data application tested
- [ ] Roundtrip conversion tested
- [ ] Nil slice handling tested

### Step 5: Document in FORMS_GUIDE.md

Add your form to the "Form Configurations" section:

```markdown
### X. YourForm

**Purpose**: Description of what this form edits  
**File**: [`internal/cli/forms/your_form.go`](../internal/cli/forms/your_form.go)  
**Fields**: N (list field names)

#### Data Structure
[Code example]

#### Usage
[Code examples]

#### Domain Conversion
[Code examples]

#### Validators
[List validators used]
```

**✅ Checklist**:
- [ ] Section added to FORMS_GUIDE.md
- [ ] File added to "Forms Package Structure" table
- [ ] Test file added to "Tests" section
- [ ] Code references updated

---

## Adding a Form to a Modal

**Time**: 20-30 minutes  
**Prerequisites**: Form configuration exists in `internal/cli/forms/`  
**Location**: `internal/cli/intents/modals.go`

### When to Use This Pattern

Use modals for inline editing within a larger intent flow:
- ✅ Edit metadata during event capture
- ✅ Edit burst name/description during review
- ✅ Edit fact details during fact extraction

### Step 1: Create Modal Structure

```go
// internal/cli/intents/modals.go

type EditYourModal struct {
    // Preserved original (NEVER mutated)
    original  *career.YourDomain
    
    // Working copy
    modified  *career.YourDomain
    
    // Result
    result    *ModalEditResult[*career.YourDomain]
    
    // Huh form
    form      *huh.Form
    formData  *forms.YourFormData
    
    // Dimensions (for responsive layout)
    width, height int
}
```

**✅ Checklist**:
- [ ] `original` field (preserve on cancel)
- [ ] `modified` field (working copy)
- [ ] `result` field (ModalEditResult[T])
- [ ] `form` field (*huh.Form)
- [ ] `formData` field (*forms.YourFormData)
- [ ] `width, height` fields (responsive)

### Step 2: Implement Constructor

```go
func NewEditYourModal(obj *career.YourDomain) *EditYourModal {
    // Get form data
    formData := forms.GetYourFormData(obj)
    
    // Create form with dimensions
    form := forms.NewYourFormWithDataAndDimensions(
        formData,
        defaultWidth-4,  // Leave margin for modal chrome
        forms.DefaultFormHeight(defaultHeight),
    )
    
    return &EditYourModal{
        original: obj,
        modified: copyYourDomain(obj), // Deep copy
        form:     form,
        formData: formData,
        width:    defaultWidth,
        height:   defaultHeight,
    }
}

// Helper for deep copy
func copyYourDomain(obj *career.YourDomain) *career.YourDomain {
    return &career.YourDomain{
        Field1: obj.Field1,
        Field2: append([]string{}, obj.Field2...), // Copy slices
    }
}
```

**✅ Checklist**:
- [ ] Form data extracted from domain
- [ ] Form created with dimensions (width-4 for chrome)
- [ ] Original preserved (not mutated)
- [ ] Modified is deep copy
- [ ] Default dimensions set

### Step 3: Implement Update Method

```go
func (m *EditYourModal) Update(msg tea.Msg) tea.Cmd {
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
        // CRITICAL: Check SubmitConfirmed
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
```

**✅ Checklist**:
- [ ] WindowSizeMsg handled
- [ ] Form updated with type assertion
- [ ] SubmitConfirmed checked on completion
- [ ] AcceptedResult created on confirm
- [ ] CancelledResult created on abort or cancel button

### Step 4: Implement Helper Methods

```go
func (m *EditYourModal) syncModified() {
    forms.ApplyYourFormData(m.modified, m.formData)
}

func (m *EditYourModal) createAcceptedResult() {
    m.result = &ModalEditResult[*career.YourDomain]{
        Original: m.original,
        Modified: m.modified,
        Accepted: true,
        Changes:  m.computeChanges(),
    }
}

func (m *EditYourModal) createCancelledResult() {
    m.result = &ModalEditResult[*career.YourDomain]{
        Original: m.original,
        Modified: m.original, // Restore original
        Accepted: false,
        Changes:  make(map[string]interface{}),
    }
}

func (m *EditYourModal) computeChanges() map[string]interface{} {
    changes := make(map[string]interface{})
    
    if m.original.Field1 != m.modified.Field1 {
        changes["field1"] = m.modified.Field1
    }
    // ... check other fields
    
    return changes
}
```

**✅ Checklist**:
- [ ] syncModified() applies form data
- [ ] createAcceptedResult() preserves original + modified
- [ ] createCancelledResult() restores original
- [ ] computeChanges() tracks what changed

### Step 5: Implement View and Interface Methods

```go
func (m *EditYourModal) View() string {
    if m.result != nil && m.result.Accepted {
        return ""
    }
    return m.form.View()
}

func (m *EditYourModal) Result() *ModalEditResult[*career.YourDomain] {
    return m.result
}

func (m *EditYourModal) IsComplete() bool {
    return m.result != nil
}

// For overlay rendering
func (m *EditYourModal) GetTitle() string {
    return "Edit Your Feature"
}

func (m *EditYourModal) GetContent() string {
    if m.result != nil && m.result.Accepted {
        return ""
    }
    return m.form.View()
}

func (m *EditYourModal) GetFooter() string {
    return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next  |  Shift+Tab: Prev"
}
```

**✅ Checklist**:
- [ ] View() implemented
- [ ] Result() returns ModalEditResult
- [ ] IsComplete() checks result != nil
- [ ] GetTitle() for overlay
- [ ] GetContent() for overlay
- [ ] GetFooter() for navigation

### Step 6: Write Modal Tests

```go
var _ = Describe("EditYourModal", func() {
    It("should preserve original on cancel", func() {
        original := &career.YourDomain{Field1: "Original"}
        modal := NewEditYourModal(original)
        
        // Simulate cancellation
        modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
        
        result := modal.Result()
        Expect(result.Accepted).To(BeFalse())
        Expect(result.Modified.Field1).To(Equal("Original"))
    })
    
    It("should apply changes on confirm", func() {
        original := &career.YourDomain{Field1: "Original"}
        modal := NewEditYourModal(original)
        
        // Modify form data
        modal.formData.Field1 = "Modified"
        modal.formData.SubmitConfirmed = true
        
        // Simulate completion
        modal.syncModified()
        modal.createAcceptedResult()
        
        result := modal.Result()
        Expect(result.Accepted).To(BeTrue())
        Expect(result.Modified.Field1).To(Equal("Modified"))
        Expect(result.Original.Field1).To(Equal("Original"))
    })
})
```

**✅ Checklist**:
- [ ] Cancel test (preserves original)
- [ ] Confirm test (applies changes)
- [ ] Overlay interface tests (GetTitle, GetContent, GetFooter)
- [ ] WindowSizeMsg test (responsive)

---

## Adding a Form to a Model

**Time**: 30-45 minutes  
**Prerequisites**: Form configuration exists  
**Location**: `internal/cli/models/your_model.go`

### When to Use This Pattern

Use models for full-screen form experiences:
- ✅ Strategy-aware forms (quick vs manual)
- ✅ Multi-step workflows
- ✅ Direct repository persistence
- ✅ Full keyboard shortcut integration

### Step 1: Create Model Structure

```go
type YourFormModel struct {
    form     *huh.Form
    data     *forms.YourFormData
    
    // Services
    service  *service.YourService
    
    // State
    saved    bool
    cancelled bool
    err      error
    
    // Dimensions
    width, height int
}
```

### Step 2: Implement Constructor

```go
func NewYourFormModel(service *service.YourService) *YourFormModel {
    data := forms.NewYourFormData() // Or from existing object
    
    form := forms.NewYourFormWithDataAndDimensions(
        data,
        0, // Will be set on first WindowSizeMsg
        0,
    )
    
    return &YourFormModel{
        form:    form,
        data:    data,
        service: service,
    }
}
```

### Step 3: Implement Init/Update/View

```go
func (m *YourFormModel) Init() tea.Cmd {
    return m.form.Init()
}

func (m *YourFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.form = m.form.
            WithWidth(msg.Width).
            WithHeight(forms.DefaultFormHeight(msg.Height))
        return m, nil
        
    case tea.KeyMsg:
        // Custom shortcuts (optional)
        if msg.Type == tea.KeyCtrlS {
            m.data.SubmitConfirmed = true
            return m.handleFormCompletion()
        }
    }
    
    // Update form
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    // Check completion
    if forms.IsCompleted(m.form) {
        return m.handleFormCompletion()
    }
    
    if forms.IsAborted(m.form) {
        m.cancelled = true
        return m, tea.Quit
    }
    
    return m, cmd
}

func (m *YourFormModel) handleFormCompletion() (tea.Model, tea.Cmd) {
    // Check SubmitConfirmed
    if !m.data.SubmitConfirmed {
        m.cancelled = true
        return m, tea.Quit
    }
    
    // Persist to repository
    err := m.service.SaveYourData(m.data)
    if err != nil {
        m.err = err
        return m, nil
    }
    
    m.saved = true
    return m, tea.Quit
}

func (m *YourFormModel) View() string {
    if m.err != nil {
        return errorView(m.err)
    }
    return m.form.View()
}
```

**✅ Checklist**:
- [ ] Init() calls form.Init()
- [ ] WindowSizeMsg handled
- [ ] Custom shortcuts implemented (optional)
- [ ] SubmitConfirmed checked
- [ ] Service integration for persistence
- [ ] Error handling implemented
- [ ] View() renders form or error

---

## Creating Custom Validators

**Time**: 5-10 minutes  
**Location**: `internal/cli/forms/validators.go`

### Step 1: Identify Validation Rules

Before creating a validator, check if you can compose existing ones:

```go
// Can this be composed?
validator := forms.Compose(
    forms.Required,
    forms.MinLength(5),
    forms.MaxLength(100),
)
```

### Step 2: Create Validator Function

```go
// Simple validator
func YourFieldName(value string) error {
    if strings.TrimSpace(value) == "" {
        return nil // Optional field
    }
    
    // Your validation logic
    if !isValid(value) {
        return fmt.Errorf("invalid format")
    }
    
    return nil
}

// Parameterized validator
func YourFieldWithParam(min, max int) func(string) error {
    return func(value string) error {
        // Validation logic using parameters
        return nil
    }
}
```

### Step 3: Write Validator Tests

```go
var _ = Describe("YourFieldName Validator", func() {
    It("should pass for valid values", func() {
        err := forms.YourFieldName("valid value")
        Expect(err).NotTo(HaveOccurred())
    })
    
    It("should fail for invalid values", func() {
        err := forms.YourFieldName("invalid")
        Expect(err).To(HaveOccurred())
    })
    
    It("should allow empty for optional fields", func() {
        err := forms.YourFieldName("")
        Expect(err).NotTo(HaveOccurred())
    })
})
```

### Step 4: Document in FORMS_GUIDE.md

Add to "Domain-Specific Validators" table:

```markdown
| `YourFieldName` | Your field name | Rules here |
```

**✅ Checklist**:
- [ ] Validator function created
- [ ] Handles empty strings (if optional)
- [ ] Tests written
- [ ] Documented in FORMS_GUIDE.md

---

## Handling Empty Fields

**CRITICAL**: Forms in KaRiya have specific requirements for handling empty fields, nil slices, and placeholders.

### Issue 1: Nil Slices Break MultiSelect

**Problem**: MultiSelect fields bind to `[]string` slices. If the domain object has `nil` for a slice field, huh will panic or behave unexpectedly.

**❌ Wrong**:
```go
func GetMetadataFormData(event *career.CareerEvent) *MetadataFormData {
    return &MetadataFormData{
        Tags:       event.Tags,       // May be nil!
        Categories: event.Categories, // May be nil!
    }
}
```

**✅ Correct**:
```go
func GetMetadataFormData(event *career.CareerEvent) *MetadataFormData {
    // Ensure slices are not nil for proper multi-select binding
    tags := event.Tags
    if tags == nil {
        tags = []string{}
    }
    categories := event.Categories
    if categories == nil {
        categories = []string{}
    }
    
    return &MetadataFormData{
        Tags:       tags,
        Categories: categories,
    }
}
```

**Pattern**: Always initialize nil slices to empty slices `[]string{}` in `GetXFormData()` functions.

### Issue 2: Empty Input Fields Show Only First Placeholder Character

**Problem**: The huh library has a display bug where empty input fields only show the first character of the placeholder text unless an explicit prompt is set.

**❌ Wrong**:
```go
huh.NewInput().
    Key("company").
    Title("Company").
    Placeholder("Enter company name...") // Only shows "E" when empty!
```

**✅ Correct**:
```go
forms.NewInput(forms.FieldConfig{
    Key:         "company",
    Title:       "Company",
    Placeholder: "Enter company name...",
}).Prompt("> ") // Explicit prompt fixes display issue
```

**Better**: Use the `forms.NewInput()` helper which automatically sets `Prompt("> ")`:

```go
forms.NewInput(forms.FieldConfig{
    Key:         "company",
    Title:       "Company",
    Placeholder: "Enter company name...",
}) // forms.NewInput() sets Prompt("> ") automatically
```

**Pattern**: Always use `forms.NewInput()` or `forms.NewText()` helpers instead of raw `huh.NewInput()`.

### Issue 3: Empty Strings vs Nil for Optional Fields

**Problem**: Some domain fields use `*string` or `*int` (pointer types) to represent optional values. Forms use plain `string` or `int`.

**Example**: Skill years of experience

**❌ Wrong**:
```go
type SkillFormData struct {
    YearsUsed *int // Don't use pointer in form data!
}
```

**✅ Correct**:
```go
type SkillFormData struct {
    YearsUsed string // Store as string, parse/validate when applying
}

// Extract from domain
func GetSkillFormData(skill *career.Skill) *SkillFormData {
    yearsStr := ""
    if skill.YearsUsed != nil {
        yearsStr = strconv.Itoa(*skill.YearsUsed)
    }
    return &SkillFormData{
        YearsUsed: yearsStr,
    }
}

// Apply to domain
func ApplySkillFormData(skill *career.Skill, data *SkillFormData) error {
    if strings.TrimSpace(data.YearsUsed) != "" {
        years, err := strconv.Atoi(data.YearsUsed)
        if err != nil {
            return fmt.Errorf("invalid years format")
        }
        skill.YearsUsed = &years
    } else {
        skill.YearsUsed = nil
    }
    return nil
}
```

**Pattern**: Use plain types (`string`, `int`) in FormData, convert to/from pointer types in Get/Apply functions.

### Issue 4: Empty String Validation

**Problem**: Optional validators must handle empty strings gracefully.

**❌ Wrong**:
```go
func CompanyName(value string) error {
    if len(value) < 2 {
        return fmt.Errorf("too short") // Fails on empty strings!
    }
    return nil
}
```

**✅ Correct**:
```go
func CompanyName(value string) error {
    if strings.TrimSpace(value) == "" {
        return nil // Allow empty for optional fields
    }
    
    if len(value) < 2 {
        return fmt.Errorf("must be at least 2 characters")
    }
    
    return nil
}
```

**Pattern**: Always check for empty strings first in optional validators, return `nil` to allow.

### Issue 5: MultiSelect Returns Different Types

**Problem**: The huh library can return MultiSelect values as either `[]string` or `[]interface{}` depending on context.

**❌ Wrong**:
```go
func GetTags(form *huh.Form) []string {
    return form.Get("tags").([]string) // May panic!
}
```

**✅ Correct**:
```go
func GetStrings(form *huh.Form, key string) []string {
    val := form.Get(key)
    if val == nil {
        return []string{}
    }
    
    switch v := val.(type) {
    case []string:
        return v
    case []interface{}:
        result := make([]string, 0, len(v))
        for _, item := range v {
            if s, ok := item.(string); ok {
                result = append(result, s)
            }
        }
        return result
    default:
        return []string{}
    }
}
```

**Pattern**: Always use `forms.GetStrings()` helper for MultiSelect fields, never type assert directly.

### Checklist: Empty Field Handling

When creating forms, verify:

- [ ] **Nil Slices**: All MultiSelect fields check for `nil` and convert to `[]string{}` in GetXFormData()
- [ ] **Prompt Fix**: All inputs use `forms.NewInput()` helper (not raw `huh.NewInput()`)
- [ ] **Optional Validators**: All optional field validators return `nil` for empty strings
- [ ] **Pointer Fields**: Pointer types (*string, *int) converted to/from plain types in Get/Apply functions
- [ ] **MultiSelect Extraction**: Use `forms.GetStrings()` for MultiSelect, never direct type assertion
- [ ] **Empty String Trim**: Use `strings.TrimSpace()` when checking for empty
- [ ] **Test Empty Cases**: Test with empty/nil values in test suite

### Code Examples

#### Complete GetXFormData() with All Empty Handling

```go
func GetMetadataFormData(event *career.CareerEvent) *MetadataFormData {
    // Handle nil slices (Issue 1)
    tags := event.Tags
    if tags == nil {
        tags = []string{}
    }
    categories := event.Categories
    if categories == nil {
        categories = []string{}
    }
    skills := event.Skills
    if skills == nil {
        skills = []string{}
    }
    
    // Handle pointer fields (Issue 3)
    company := event.Company
    project := event.Project
    
    // Format date (always has value in domain)
    dateStr := event.Date.Format("2006-01-02")
    
    return &MetadataFormData{
        Date:       dateStr,
        Company:    company,
        Project:    project,
        Tags:       tags,
        Categories: categories,
        Skills:     skills,
    }
}
```

#### Complete ApplyXFormData() with Validation

```go
func ApplyMetadataFormData(event *career.CareerEvent, data *MetadataFormData) error {
    // Parse and validate date
    parsedDate, err := forms.ParseDateString(data.Date)
    if err != nil {
        return err
    }
    
    // Trim whitespace from string fields (Issue 4)
    company := strings.TrimSpace(data.Company)
    project := strings.TrimSpace(data.Project)
    
    // Apply to domain
    event.Date = parsedDate
    event.Company = company
    event.Project = project
    event.Tags = data.Tags           // Already []string from form
    event.Categories = data.Categories // Already []string from form
    event.Skills = data.Skills       // Already []string from form
    
    return nil
}
```

#### Complete Form Creation with All Fixes

```go
func NewMetadataEditorFormWithDataAndDimensions(
    data *MetadataFormData,
    availableTags, availableCategories []string,
    availableSkills []*career.Skill,
    width, height int,
) *huh.Form {
    // Initialize submit confirmation
    data.SubmitConfirmed = false
    
    // Prepare MultiSelect options
    tagOptions := make([]huh.Option[string], len(availableTags))
    for i, tag := range availableTags {
        tagOptions[i] = huh.NewOption(tag, tag)
    }
    
    // Create fields group
    fieldsGroup := huh.NewGroup(
        // Use forms.NewInput() for prompt fix (Issue 2)
        forms.NewInput(forms.FieldConfig{
            Key:         "date",
            Title:       "Date",
            Placeholder: "YYYY-MM-DD or 'today'",
            Validate:    forms.DateFormat, // Handles empty (Issue 4)
        }).Value(&data.Date),
        
        forms.NewInput(forms.FieldConfig{
            Key:         "company",
            Title:       "Company",
            Placeholder: "Enter company name...",
            Validate:    forms.CompanyName, // Handles empty (Issue 4)
        }).Value(&data.Company),
        
        // MultiSelect with nil-safe binding (Issue 1)
        huh.NewMultiSelect[string]().
            Key("tags").
            Title("Tags").
            Options(tagOptions...).
            Value(&data.Tags). // Already initialized to []string{} if was nil
            Limit(10),
    )
    
    return forms.NewFormWithFixedConfirm(fieldsGroup, &data.SubmitConfirmed, width, height)
}
```

### Testing Empty Field Scenarios

```go
var _ = Describe("Empty Field Handling", func() {
    It("should handle nil slices", func() {
        event := &career.CareerEvent{
            Tags:       nil,
            Categories: nil,
        }
        
        data := forms.GetMetadataFormData(event)
        
        // Should be empty slice, not nil
        Expect(data.Tags).NotTo(BeNil())
        Expect(data.Tags).To(Equal([]string{}))
        Expect(data.Categories).NotTo(BeNil())
        Expect(data.Categories).To(Equal([]string{}))
    })
    
    It("should handle empty strings in optional fields", func() {
        data := &forms.MetadataFormData{
            Company: "",
            Project: "",
        }
        
        event := &career.CareerEvent{}
        err := forms.ApplyMetadataFormData(event, data)
        
        Expect(err).NotTo(HaveOccurred())
        Expect(event.Company).To(Equal(""))
        Expect(event.Project).To(Equal(""))
    })
    
    It("should validate optional fields with CompanyName validator", func() {
        // Empty is valid (optional)
        err := forms.CompanyName("")
        Expect(err).NotTo(HaveOccurred())
        
        // Too short is invalid
        err = forms.CompanyName("A")
        Expect(err).To(HaveOccurred())
        
        // Valid length
        err = forms.CompanyName("Acme Corp")
        Expect(err).NotTo(HaveOccurred())
    })
})
```

---

## Integration Checklists

### Form Configuration Checklist

- [ ] **FormData structure** defined with all fields
- [ ] **SubmitConfirmed bool** included
- [ ] **3 builder variants** created (basic, with data, with dimensions)
- [ ] **Validators** applied to all fields
- [ ] **`.Value(&field)`** binding on all fields
- [ ] **GetXFormData()** function created
- [ ] **ApplyXFormData()** function created
- [ ] **Nil slices** handled (convert to empty slice)
- [ ] **NewFormWithFixedConfirm()** used (not NewForm)
- [ ] **Tests** written (creation, extraction, application, roundtrip)
- [ ] **Documented** in FORMS_GUIDE.md

### Modal Integration Checklist

- [ ] **Original** field preserved (never mutated)
- [ ] **Modified** field is deep copy
- [ ] **FormData** field for binding
- [ ] **Dimensions** tracked (width, height)
- [ ] **WindowSizeMsg** handled
- [ ] **SubmitConfirmed** checked on completion
- [ ] **syncModified()** applies form data to modified
- [ ] **createAcceptedResult()** preserves original + modified
- [ ] **createCancelledResult()** restores original
- [ ] **Interface methods** implemented (GetTitle, GetContent, GetFooter)
- [ ] **Tests** written (cancel, confirm, overlay)

### Model Integration Checklist

- [ ] **Form** embedded in model
- [ ] **FormData** embedded in model
- [ ] **Services** injected for persistence
- [ ] **WindowSizeMsg** handled
- [ ] **SubmitConfirmed** checked
- [ ] **Error handling** implemented
- [ ] **Init()** calls form.Init()
- [ ] **Custom shortcuts** implemented (if needed)
- [ ] **Tests** written

### Theme Integration Checklist

- [ ] **Default theme** used (Catppuccin via NewForm)
- [ ] **OR custom theme** via NewThemedForm()
- [ ] **NOT** using huh.ThemeBase() or other non-KaRiya themes
- [ ] **Responsive** dimensions using DefaultFormHeight()

---

## Common Pitfalls

### Pitfall 1: Forgetting `.Value(&field)` Binding

**❌ Wrong**:
```go
huh.NewInput().Key("name").Title("Name") // No binding!
```

**✅ Correct**:
```go
huh.NewInput().
    Key("name").
    Title("Name").
    Value(&data.Name) // Binds to struct field
```

### Pitfall 2: Not Checking SubmitConfirmed on Cancel

**❌ Wrong**:
```go
if forms.IsCompleted(m.form) {
    m.syncModified()
    m.createAcceptedResult() // Wrong! User may have clicked "Cancel" on confirm button
}
```

**✅ Correct**:
```go
if forms.IsCompleted(m.form) {
    if m.formData.SubmitConfirmed {
        m.syncModified()
        m.createAcceptedResult()
    } else {
        m.createCancelledResult()
    }
}
```

### Pitfall 3: Hardcoding Dimensions

**❌ Wrong**:
```go
form := forms.NewFormWithHeight(30, groups...) // Hardcoded!
```

**✅ Correct**:
```go
height := forms.DefaultFormHeight(terminalHeight)
form := forms.NewFormWithHeight(height, groups...)
```

### Pitfall 4: Mutating Original in Modal

**❌ Wrong**:
```go
func NewEditModal(obj *Domain) *EditModal {
    return &EditModal{
        original: obj,
        modified: obj, // SAME REFERENCE!
    }
}
```

**✅ Correct**:
```go
func NewEditModal(obj *Domain) *EditModal {
    return &EditModal{
        original: obj,
        modified: copyDomain(obj), // Deep copy
    }
}
```

### Pitfall 5: Not Handling Nil Slices for MultiSelect

**CRITICAL**: This causes form panics! See [Handling Empty Fields](#handling-empty-fields) for complete documentation.

**❌ Wrong**:
```go
func GetFormData(obj *Domain) *FormData {
    return &FormData{
        Tags: obj.Tags, // May be nil! Will cause panic in MultiSelect
    }
}
```

**✅ Correct**:
```go
func GetFormData(obj *Domain) *FormData {
    // ALWAYS check for nil slices
    tags := obj.Tags
    if tags == nil {
        tags = []string{} // Initialize to empty slice
    }
    return &FormData{
        Tags: tags,
    }
}
```

**Additional Empty Field Issues**: See [Handling Empty Fields](#handling-empty-fields) section for:
- Empty placeholder display bug (use `forms.NewInput()`)
- Optional field validation (check for empty first)
- Pointer type conversion (*int, *string)
- MultiSelect type assertion safety

### Pitfall 6: Using Wrong Form Builder

**❌ Wrong**:
```go
// For modal with confirm button
form := forms.NewForm(groups...) // No fixed confirm!
```

**✅ Correct**:
```go
// For modal with fixed confirm button
form := forms.NewFormWithFixedConfirm(fieldsGroup, &confirmValue, width, height)
```

### Pitfall 7: Writing Custom Validators Instead of Composing

**❌ Wrong**:
```go
func MyCustomValidator(value string) error {
    if strings.TrimSpace(value) == "" {
        return fmt.Errorf("required")
    }
    if len(value) < 5 {
        return fmt.Errorf("too short")
    }
    if len(value) > 100 {
        return fmt.Errorf("too long")
    }
    return nil
}
```

**✅ Correct**:
```go
validator := forms.Compose(
    forms.Required,
    forms.MinLength(5),
    forms.MaxLength(100),
)
```

---

## Testing Strategy

### What to Test

**Form Configuration Tests**:
1. ✅ Form creation with domain object
2. ✅ Form creation with custom data
3. ✅ Data extraction (GetXFormData)
4. ✅ Data application (ApplyXFormData)
5. ✅ Roundtrip conversion (extract → apply)
6. ✅ Nil slice handling
7. ✅ Enum conversions (if applicable)

**Modal Tests**:
1. ✅ Cancel preserves original
2. ✅ Confirm applies changes
3. ✅ SubmitConfirmed on cancel returns cancelled result
4. ✅ Overlay interface methods return correct values
5. ✅ WindowSizeMsg updates dimensions

**Model Tests**:
1. ✅ Form initialization
2. ✅ WindowSizeMsg handling
3. ✅ Completion handling
4. ✅ Cancellation handling
5. ✅ Error handling
6. ✅ Repository integration (if applicable)

**Validator Tests**:
1. ✅ Valid values pass
2. ✅ Invalid values fail
3. ✅ Empty strings (for optional fields)
4. ✅ Boundary conditions
5. ✅ Edge cases

### Test Template

```go
var _ = Describe("YourForm", func() {
    Describe("Form Creation", func() {
        It("should create form with domain object", func() {
            // Test
        })
    })
    
    Describe("Data Conversion", func() {
        It("should extract data from domain", func() {
            // Test
        })
        
        It("should apply data to domain", func() {
            // Test
        })
        
        It("should handle roundtrip conversion", func() {
            // Test
        })
    })
    
    Describe("Edge Cases", func() {
        It("should handle nil slices", func() {
            // Test
        })
    })
})
```

---

## Quick Reference Commands

```bash
# Run form tests
go test -v ./internal/cli/forms/...

# Run modal tests
go test -v ./internal/cli/intents/modals_test.go

# Run specific form test
ginkgo -r --focus="YourForm" ./internal/cli/forms/

# Check form validator
go test -v -run TestValidators ./internal/cli/forms/
```

---

## Additional Resources

- **API Reference**: [`docs/FORMS_GUIDE.md`](../FORMS_GUIDE.md) - Complete technical documentation
- **TUI Patterns**: [`docs/TUI_DEVELOPER_GUIDE.md`](../TUI_DEVELOPER_GUIDE.md) - Generic TUI component patterns
- **Theme System**: [`docs/THEME_CUSTOMIZATION_GUIDE.md`](../THEME_CUSTOMIZATION_GUIDE.md) - Theme customization
- **Code Examples**: `internal/cli/forms/` - All form configurations

---

*For questions about forms workflow, see this guide first. For API details, see [`docs/FORMS_GUIDE.md`](../FORMS_GUIDE.md).*
