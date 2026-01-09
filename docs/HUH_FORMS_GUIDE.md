# Huh Forms Guide

**Last Updated**: 2026-01-07  
**Status**: Active - All KaRiya forms use `huh` library

---

## Table of Contents

1. [Overview](#overview)
2. [Why Huh?](#why-huh)
3. [Quick Start](#quick-start)
4. [Form Structure](#form-structure)
5. [Field Types](#field-types)
6. [Validation](#validation)
7. [Dynamic Forms](#dynamic-forms)
8. [Integration with BubbleTea](#integration-with-bubbletea)
9. [Theming](#theming)
10. [Common Patterns](#common-patterns)
11. [Testing Forms](#testing-forms)
12. [Migration Notes](#migration-notes)

---

## Overview

KaRiya uses Charm's **`huh`** library for all form handling in the TUI. The `huh` library provides:

- **Type-safe forms** with compile-time checking
- **Native BubbleTea integration** - forms are `tea.Model` instances
- **Built-in validation** with custom validator support
- **Accessible design** with keyboard navigation
- **Catppuccin theming** for beautiful, consistent UI
- **Dynamic field visibility** based on user input

**Library**: [github.com/charmbracelet/huh](https://github.com/charmbracelet/huh)  
**Version**: v0.8.0+

---

## Why Huh?

### Before (Manual Implementation)

```go
// 956 lines in internal/cli/models/form.go
type FormModel struct {
    inputs      []textinput.Model
    focusIndex  int
    // ... 50+ lines of boilerplate
}

func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
            // ... 200+ lines of focus management
        }
    }
    // ... 300+ lines of update logic
}

func (m *FormModel) View() string {
    // ... 200+ lines of rendering
}
```

**Total**: ~3,800 lines across 8 form implementations

### After (Huh)

```go
// ~50 lines per form
func NewBurstEditorForm(burst *domain.Burst) *huh.Form {
    return forms.NewForm(
        huh.NewGroup(
            forms.NewInput(forms.FieldConfig{
                Key:      "title",
                Title:    "Title",
                Validate: forms.Title,
            }),
            forms.NewText(forms.FieldConfig{
                Key:       "description",
                Title:     "Description",
                CharLimit: 1000,
                Validate:  forms.Description,
            }),
        ),
    )
}
```

**Total**: ~960 lines for all 8 forms (75% reduction)

**Benefits**:
- ✅ No manual focus management
- ✅ No manual validation state tracking
- ✅ No manual rendering logic
- ✅ Built-in accessibility
- ✅ Consistent UX across all forms
- ✅ Type-safe field access

---

## Quick Start

### 1. Import the Forms Package

```go
import (
    "github.com/charmbracelet/huh"
    "github.com/baphled/kariya/internal/cli/forms"
)
```

### 2. Create a Form

```go
form := forms.NewForm(
    huh.NewGroup(
        huh.NewInput().
            Key("name").
            Title("What's your name?").
            Validate(forms.Required),
        
        huh.NewText().
            Key("description").
            Title("Tell us about yourself").
            CharLimit(500),
    ),
)
```

### 3. Embed in Your Model

```go
type MyModel struct {
    form *huh.Form
}

func (m *MyModel) Init() tea.Cmd {
    return m.form.Init()
}

func (m *MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    
    // Check form state
    if forms.IsCompleted(m.form) {
        name := forms.GetString(m.form, "name")
        desc := forms.GetString(m.form, "description")
        // Handle completion
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

---

## Form Structure

### Groups

Forms are organized into **groups**. Each group is displayed on its own "page", with automatic pagination.

```go
form := forms.NewForm(
    // Page 1: Basic info
    huh.NewGroup(
        huh.NewInput().Key("name").Title("Name"),
        huh.NewInput().Key("email").Title("Email"),
    ),
    
    // Page 2: Details
    huh.NewGroup(
        huh.NewText().Key("bio").Title("Bio"),
        huh.NewSelect[string]().Key("role").Title("Role"),
    ),
)
```

**Navigation**:
- `Tab` / `Shift+Tab`: Next/previous field
- `Enter`: Next group/complete form
- `Esc`: Cancel form

### Conditional Groups

Hide entire groups based on logic:

```go
showOptional := strategy == "manual"

form := forms.NewForm(
    huh.NewGroup(
        huh.NewInput().Key("required").Title("Required Field"),
    ),
    
    huh.NewGroup(
        huh.NewInput().Key("optional").Title("Optional Field"),
    ).WithHideFunc(func() bool {
        return !showOptional
    }),
)
```

---

## Field Types

### Input (Single Line Text)

```go
huh.NewInput().
    Key("company").
    Title("Company Name").
    Placeholder("Acme Corp").
    CharLimit(100).
    Validate(forms.CompanyName)
```

**Use for**: Short text (names, emails, dates, single-line values)

### Text (Multi-Line Text Area)

```go
huh.NewText().
    Key("description").
    Title("Event Description").
    Placeholder("Describe what you accomplished...").
    CharLimit(2000).
    Validate(forms.EventText)
```

**Use for**: Long text (descriptions, multi-line content)

### Select (Single Choice)

```go
huh.NewSelect[string]().
    Key("strategy").
    Title("Capture Strategy").
    Options(
        huh.NewOption("Quick Capture", "quick"),
        huh.NewOption("Manual Capture", "manual"),
        huh.NewOption("Import CSV", "csv"),
    ).
    Value(&selectedStrategy) // Bind to variable
```

**Use for**: Single selection from options

### MultiSelect (Multiple Choices)

```go
huh.NewMultiSelect[string]().
    Key("tags").
    Title("Select Tags").
    Options(
        huh.NewOption("Go", "go"),
        huh.NewOption("Python", "python"),
        huh.NewOption("JavaScript", "javascript"),
    ).
    Limit(5). // Max selections
    Value(&selectedTags) // []string
```

**Use for**: Multiple selections with optional limit

### Confirm (Yes/No)

```go
huh.NewConfirm().
    Key("confirm").
    Title("Are you sure?").
    Description("This action cannot be undone").
    Affirmative("Yes!").
    Negative("No")
```

**Use for**: Boolean confirmations

### FilePicker (File Selection)

```go
huh.NewFilePicker().
    Key("file").
    Title("Select CSV File").
    CurrentDirectory(".").
    AllowedTypes([]string{".csv"})
```

**Use for**: File selection dialogs

---

## Validation

### Built-in Validators

Located in `internal/cli/forms/validators.go`:

```go
// Generic validators
forms.Required              // Non-empty
forms.MinLength(10)         // Min characters
forms.MaxLength(100)        // Max characters
forms.LengthRange(10, 100)  // Range
forms.DateFormat            // YYYY-MM-DD
forms.Email                 // Email format
forms.URL                   // URL format
forms.AlphaNumeric          // Letters & numbers only
forms.NoSpecialChars        // No special chars

// Domain-specific validators
forms.EventText             // Career event text (10-2000 chars, required)
forms.EventTextOptional     // Same but optional
forms.CompanyName           // Company name (2-100 chars)
forms.TagName               // Tag name (1-50 chars, no special chars)
forms.Title                 // Generic title (3-200 chars, required)
forms.Description           // Generic description (10-1000 chars, optional)
forms.ProfileName           // CV profile name (2-100 chars, required)
forms.AudienceName          // CV audience name (2-100 chars, required)
```

### Composing Validators

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
// Define custom validation logic
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

### Using forms.Custom

```go
huh.NewInput().
    Key("code").
    Validate(forms.Custom(
        func(val string) bool {
            return len(val) == 6 && regexp.MustCompile(`^\d+$`).MatchString(val)
        },
        "must be exactly 6 digits",
    ))
```

---

## Dynamic Forms

### Dynamic Field Visibility

Show/hide fields based on form state:

```go
var strategy string

form := forms.NewForm(
    huh.NewGroup(
        huh.NewSelect[string]().
            Key("strategy").
            Title("Capture Strategy").
            Options(
                huh.NewOption("Quick", "quick"),
                huh.NewOption("Manual", "manual"),
            ).
            Value(&strategy),
    ),
    
    huh.NewGroup(
        huh.NewInput().Key("company").Title("Company"),
        huh.NewInput().Key("date").Title("Date"),
    ).WithHideFunc(func() bool {
        return strategy != "manual" // Only show in manual mode
    }),
)
```

### Dynamic Options

Update options based on external state:

```go
form := forms.NewForm(
    huh.NewGroup(
        huh.NewSelect[string]().
            Key("profile").
            Title("Select Profile").
            OptionsFunc(func() []huh.Option[string] {
                // Fetch profiles dynamically
                profiles := fetchProfiles()
                opts := make([]huh.Option[string], len(profiles))
                for i, p := range profiles {
                    opts[i] = huh.NewOption(p.Name, p.ID)
                }
                return opts
            }, &someWatchedValue), // Re-fetch when value changes
    ),
)
```

### Dynamic Titles

Change field titles dynamically:

```go
huh.NewInput().
    Key("name").
    TitleFunc(func() string {
        if strategy == "quick" {
            return "Quick Event Name"
        }
        return "Detailed Event Name"
    }, &strategy)
```

---

## Integration with BubbleTea

### Standalone Form (Modal)

```go
type EditModal struct {
    form *huh.Form
}

func NewEditModal(data *MyData) *EditModal {
    form := forms.NewForm(
        huh.NewGroup(
            huh.NewInput().Key("field").Value(&data.Field),
        ),
    )
    
    return &EditModal{form: form}
}

func (m *EditModal) Init() tea.Cmd {
    return m.form.Init()
}

func (m *EditModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    return m, cmd
}

func (m *EditModal) View() string {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        Render(m.form.View())
}
```

### Embedded in Intent

```go
type CaptureEventModel struct {
    form   *huh.Form
    state  CaptureState
}

func (m *CaptureEventModel) Update(msg tea.Msg) tea.Cmd {
    switch m.state {
    case StateForm:
        form, cmd := m.form.Update(msg)
        if f, ok := form.(*huh.Form); ok {
            m.form = f
        }
        
        if forms.IsCompleted(m.form) {
            m.state = StateReview
            m.extractFormData()
        }
        
        if forms.IsAborted(m.form) {
            m.state = StateCancelled
        }
        
        return cmd
    }
    
    return nil
}
```

---

## Theming

### Default Theme (Catppuccin)

All KaRiya forms use the **Catppuccin** theme by default:

```go
form := forms.NewForm(groups...) // Automatically uses Catppuccin
```

### Color Reference

Available in `forms.FormColors`:

```go
forms.FormColors.Title       // #89B4FA (Blue)
forms.FormColors.Description // #94E2D5 (Teal)
forms.FormColors.Error       // #F38BA8 (Red)
forms.FormColors.Success     // #A6E3A1 (Green)
forms.FormColors.Placeholder // #6C7086 (Overlay0)
```

### Custom Theme

Override the default theme if needed:

```go
customTheme := huh.ThemeBase16() // Or any other theme
form := huh.NewForm(groups...).WithTheme(customTheme)
```

### Accessibility

Enable accessibility mode for screen readers:

```go
form := forms.NewFormWithAccessible(groups...)
```

---

## Common Patterns

### Pattern 1: Simple Edit Form

```go
func NewEditBurstForm(burst *domain.Burst) *huh.Form {
    return forms.NewForm(
        huh.NewGroup(
            forms.NewInput(forms.FieldConfig{
                Key:      "title",
                Title:    "Title",
                Validate: forms.Title,
            }).Value(&burst.Title),
            
            forms.NewText(forms.FieldConfig{
                Key:       "content",
                Title:     "Content",
                CharLimit: 1000,
                Validate:  forms.Description,
            }).Value(&burst.Content),
        ),
    )
}
```

### Pattern 2: Multi-Page Form

```go
func NewCaptureEventForm(strategy string) *huh.Form {
    return forms.NewForm(
        // Page 1: Event text
        huh.NewGroup(
            forms.NewText(forms.FieldConfig{
                Key:       "text",
                Title:     "Event Text",
                CharLimit: 2000,
                Validate:  forms.EventText,
            }),
        ),
        
        // Page 2: Metadata (optional)
        huh.NewGroup(
            forms.NewInput(forms.FieldConfig{
                Key:      "company",
                Title:    "Company",
                Validate: forms.CompanyName,
            }),
            forms.NewInput(forms.FieldConfig{
                Key:      "date",
                Title:    "Date",
                Validate: forms.DateFormat,
            }),
        ).WithHideFunc(func() bool {
            return strategy != "manual"
        }),
    )
}
```

### Pattern 3: Confirmation Dialog

```go
func NewDeleteConfirmForm(itemName string) *huh.Form {
    return forms.NewForm(
        huh.NewGroup(
            forms.NewConfirm(
                "confirm",
                fmt.Sprintf("Delete '%s'?", itemName),
                "This action cannot be undone.",
                "Yes, delete",
                "Cancel",
            ),
        ),
    )
}
```

### Pattern 4: Dynamic Tag Selection

```go
func NewTagSelectorForm(availableTags []string) *huh.Form {
    options := make([]forms.SelectOption, len(availableTags))
    for i, tag := range availableTags {
        options[i] = forms.SelectOption{Key: tag, Value: tag}
    }
    
    return forms.NewForm(
        huh.NewGroup(
            forms.NewMultiSelect(
                "tags",
                "Select Tags",
                "Choose one or more tags",
                options,
                10, // Max 10 selections
            ),
        ),
    )
}
```

---

## Testing Forms

### Unit Testing Form Creation

```go
var _ = Describe("BurstEditorForm", func() {
    It("should create form with correct fields", func() {
        burst := &domain.Burst{Title: "Test", Content: "Content"}
        form := NewEditBurstForm(burst)
        
        Expect(form).NotTo(BeNil())
        Expect(form.State).To(Equal(huh.StateNormal))
    })
})
```

### Integration Testing with BubbleTea

```go
var _ = Describe("EditModal Integration", func() {
    It("should complete form on valid input", func() {
        modal := NewEditModal(data)
        
        // Initialize
        modal.Init()
        
        // Simulate user input
        modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")})
        modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
        
        Expect(forms.IsCompleted(modal.form)).To(BeTrue())
        Expect(forms.GetString(modal.form, "field")).To(Equal("test"))
    })
})
```

### Testing Validation

```go
var _ = Describe("Form Validation", func() {
    It("should reject empty required fields", func() {
        err := forms.Required("")
        Expect(err).To(HaveOccurred())
    })
    
    It("should accept valid dates", func() {
        err := forms.DateFormat("2024-01-07")
        Expect(err).NotTo(HaveOccurred())
    })
    
    It("should reject invalid dates", func() {
        err := forms.DateFormat("2024-13-99")
        Expect(err).To(HaveOccurred())
    })
})
```

---

## Migration Notes

### From Manual Forms to Huh

**Step 1**: Identify form fields
```go
// Old
inputs := []textinput.Model{
    textinput.New().SetPlaceholder("Name"),
    textinput.New().SetPlaceholder("Email"),
}
```

**Step 2**: Convert to huh fields
```go
// New
form := forms.NewForm(
    huh.NewGroup(
        huh.NewInput().Key("name").Placeholder("Name"),
        huh.NewInput().Key("email").Placeholder("Email"),
    ),
)
```

**Step 3**: Replace Update logic
```go
// Old: Manual focus management (50+ lines)
func (m *OldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
            // ... more focus logic
        }
    }
}

// New: Delegate to huh (5 lines)
func (m *NewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    form, cmd := m.form.Update(msg)
    if f, ok := form.(*huh.Form); ok {
        m.form = f
    }
    return m, cmd
}
```

**Step 4**: Replace View logic
```go
// Old: Manual rendering (100+ lines)
func (m *OldModel) View() string {
    var b strings.Builder
    for i, input := range m.inputs {
        if i == m.focusIndex {
            b.WriteString("> ")
        }
        b.WriteString(input.View())
        b.WriteString("\n")
    }
    return b.String()
}

// New: Delegate to huh (1 line)
func (m *NewModel) View() string {
    return m.form.View()
}
```

### Removed Components

After migration, these old components are **removed**:
- ❌ `internal/cli/models/form.go` (956 lines)
- ❌ `internal/cli/components/form_container.go`
- ❌ `internal/cli/components/form_field_container.go`

All functionality is now handled by **huh** + **`internal/cli/forms/`** helpers.

---

## Additional Resources

- **Huh GitHub**: https://github.com/charmbracelet/huh
- **Huh Examples**: https://github.com/charmbracelet/huh/tree/main/examples
- **BubbleTea Docs**: https://github.com/charmbracelet/bubbletea
- **Catppuccin Theme**: https://github.com/catppuccin/catppuccin

---

## Summary

**Key Takeaways**:
1. Use `forms.NewForm()` to create forms with Catppuccin theme
2. Use pre-built validators from `forms` package
3. Embed `*huh.Form` in your BubbleTea models
4. Check `forms.IsCompleted()` and `forms.IsAborted()` for state
5. Extract values with `forms.GetString()`, `forms.GetBool()`, etc.
6. Use dynamic features (`WithHideFunc`, `OptionsFunc`) for conditional UI
7. Test forms with Ginkgo/Gomega like any other component

**Benefits**:
- 75% less code (3,800 → 960 lines)
- Consistent UX across all forms
- Built-in accessibility and theming
- Type-safe field access
- No manual focus/validation management

---

*For questions or suggestions about form handling, see the `internal/cli/forms/` package or open a GitHub issue.*
