# Helper Extraction Guide

**Purpose**: Where to place helper functions when extracting from intent.go  
**Status**: REQUIRED for Check #18 compliance (intent.go file size <600 lines)  
**Enforcement**: `check-intent-architecture.sh` (Check #18)

---

## Overview

### Why Extract Helpers?

**Problem**: Intent files bloated with helper methods

```go
// File: my_intent.go (800 lines - VIOLATION)

type MyIntent struct {
    *BaseIntent
    // ... fields
}

func (i *MyIntent) Init() tea.Cmd { /* ... */ }
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd { /* ... */ }
func (i *MyIntent) View() string { /* ... */ }

// 500+ lines of helper methods below
func (i *MyIntent) validateInput(s string) error { /* ... */ }
func (i *MyIntent) formatDate(t time.Time) string { /* ... */ }
func (i *MyIntent) filterItems(items []*Item) []*Item { /* ... */ }
func (i *MyIntent) sortItems(items []*Item) []*Item { /* ... */ }
func (i *MyIntent) groupByTag(items []*Item) map[string][]*Item { /* ... */ }
func (i *MyIntent) calculateStats(items []*Item) *Stats { /* ... */ }
// ... 50 more helpers
```

**Solution**: Extract helpers to appropriate locations

- Business logic → `context.go`
- Private orchestration → Keep in `intent.go` (if ≤5)
- Utilities → `internal/util/` package
- Formatting → `internal/formatting/` package

**Goal**: Intent.go ≤ 600 lines (target 200-400 lines)

---

## Decision Tree

```
Is this helper a...

1. Business logic method?
   (DB queries, validation, calculations)
   → Move to context.go
   
2. Rendering method?
   (View, format, style)
   → Move to screens/{feature}/
   
3. Private orchestration helper?
   (State transitions, result building)
   → Keep in intent.go (MAX 5 helpers)
   
4. General utility?
   (Date formatting, string manipulation)
   → Move to internal/util/
   
5. Domain-specific formatting?
   (CV formatting, timeline rendering)
   → Move to internal/formatting/
```

---

## Extraction Patterns

### 1. Business Logic → context.go

**Identify business logic** methods:
- Database queries
- Validation logic
- Calculations
- Data transformations
- Filtering/sorting

#### BEFORE (Intent)
```go
// File: my_intent.go
func (i *MyIntent) validateBurst(burst *domain.Burst) error {
    if burst.Title == "" {
        return errors.New("title required")
    }
    if burst.StartDate.After(burst.EndDate) {
        return errors.New("start date must be before end date")
    }
    return nil
}

func (i *MyIntent) filterBurstsByTag(tag string) []*domain.Burst {
    var filtered []*domain.Burst
    for _, burst := range i.bursts {
        for _, t := range burst.Tags {
            if t == tag {
                filtered = append(filtered, burst)
                break
            }
        }
    }
    return filtered
}

func (i *MyIntent) calculateDuration(burst *domain.Burst) time.Duration {
    return burst.EndDate.Sub(burst.StartDate)
}
```

#### AFTER (Context)
```go
// File: intents/myfeature/context.go
func (c *MyFeatureContext) ValidateBurst(burst *domain.Burst) error {
    if burst.Title == "" {
        return errors.New("title required")
    }
    if burst.StartDate.After(burst.EndDate) {
        return errors.New("start date must be before end date")
    }
    return nil
}

func (c *MyFeatureContext) FilterBurstsByTag(tag string) []*domain.Burst {
    var filtered []*domain.Burst
    for _, burst := range c.bursts {
        for _, t := range burst.Tags {
            if t == tag {
                filtered = append(filtered, burst)
                break
            }
        }
    }
    return filtered
}

func (c *MyFeatureContext) CalculateDuration(burst *domain.Burst) time.Duration {
    return burst.EndDate.Sub(burst.StartDate)
}
```

**Usage in intent**:
```go
// Intent delegates to context
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    burst := result.Data.(*domain.Burst)
    
    // Delegate validation to context
    if err := i.context.ValidateBurst(burst); err != nil {
        return i.showError(err)
    }
    
    // ... continue
}
```

---

### 2. Rendering → screens/{feature}/

**Identify rendering** methods:
- View formatting
- Styling logic
- Layout composition

#### BEFORE (Intent)
```go
// File: my_intent.go
func (i *MyIntent) formatBurstItem(burst *domain.Burst) string {
    titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
    dateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
    
    title := titleStyle.Render(burst.Title)
    date := dateStyle.Render(burst.StartDate.Format("2006-01-02"))
    
    return fmt.Sprintf("%s - %s", title, date)
}

func (i *MyIntent) renderBurstList() string {
    // 80 lines of rendering
}
```

#### AFTER (Screen)
```go
// File: screens/myfeature/list_screen.go
func (s *ListScreen) formatBurstItem(burst *domain.Burst) string {
    theme := s.GetTheme()
    
    title := primitives.Title(burst.Title, theme)
    date := primitives.Muted(burst.StartDate.Format("2006-01-02"), theme)
    
    return fmt.Sprintf("%s - %s", title, date)
}

func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    // Use formatBurstItem as private helper
    items := make([]string, len(s.bursts))
    for i, burst := range s.bursts {
        items[i] = s.formatBurstItem(burst)
    }
    
    // ...
}
```

**See**: [Screen Extraction Guide](SCREEN_EXTRACTION_GUIDE.md)

---

### 3. Private Orchestration → Keep in intent.go

**Keep these helpers** in intent.go (MAX 5):
- State transition helpers
- Result builders
- Screen switching
- Modal management

#### GOOD (Keep in intent.go)
```go
// File: intents/myfeature/intent.go

// setCancelled marks the intent as cancelled.
func (i *MyIntent) setCancelled() {
    i.result = &IntentResult[*MyResult]{
        Status: Cancelled,
    }
    i.active = false
}

// setCompleted marks the intent as completed.
func (i *MyIntent) setCompleted(data *MyResult) {
    i.result = &IntentResult[*MyResult]{
        Status: Completed,
        Data:   data,
    }
    i.active = false
}

// switchToListScreen transitions to list screen.
func (i *MyIntent) switchToListScreen() tea.Cmd {
    i.listScreen = myfeature.NewListScreen(i.context.GetItems())
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    i.activeScreen = i.listScreen
    i.state = StateList
    return nil
}

// showError displays an error modal.
func (i *MyIntent) showError(err error) tea.Cmd {
    i.errorModal = feedback.NewErrorModal(err.Error())
    i.errorModal.Show()
    return nil
}

// showSuccess displays a success modal.
func (i *MyIntent) showSuccess(msg string) tea.Cmd {
    i.successModal = feedback.NewSuccessModal(msg)
    i.successModal.Show()
    return nil
}
```

**Guideline**: ≤5 private helpers in intent.go (orchestration only)

---

### 4. General Utilities → internal/util/

**Identify general utilities**:
- Date formatting
- String manipulation
- Collection helpers
- Type conversions

#### BEFORE (Intent)
```go
// File: my_intent.go
func (i *MyIntent) formatDate(t time.Time) string {
    if t.IsZero() {
        return "N/A"
    }
    return t.Format("Jan 02, 2006")
}

func (i *MyIntent) truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}

func (i *MyIntent) uniqueStrings(items []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, item := range items {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}
```

#### AFTER (Utility Package)
```go
// File: internal/util/date.go
package util

import "time"

// FormatDate formats a date for display.
func FormatDate(t time.Time) string {
    if t.IsZero() {
        return "N/A"
    }
    return t.Format("Jan 02, 2006")
}

// File: internal/util/string.go
package util

// TruncateString truncates a string to maxLen with ellipsis.
func TruncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}

// UniqueStrings returns unique strings from a slice.
func UniqueStrings(items []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, item := range items {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}
```

**Usage in intent**:
```go
import "github.com/baphled/kariya/internal/util"

func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    burst := result.Data.(*domain.Burst)
    formattedDate := util.FormatDate(burst.StartDate)
    // ...
}
```

---

### 5. Domain Formatting → internal/formatting/

**Identify domain-specific formatting**:
- CV generation
- Timeline rendering
- Report formatting
- Export logic

#### BEFORE (Intent)
```go
// File: generate_cv_intent.go (2,746 lines!)
func (i *GenerateCVIntent) formatEducationSection(edu []*domain.Education) string {
    // 200 lines of CV formatting
}

func (i *GenerateCVIntent) formatExperienceSection(exp []*domain.Experience) string {
    // 250 lines of CV formatting
}

func (i *GenerateCVIntent) formatSkillsSection(skills []*domain.Skill) string {
    // 180 lines of CV formatting
}
```

#### AFTER (Formatting Package)
```go
// File: internal/formatting/cv/education.go
package cv

import "github.com/baphled/kariya/internal/domain"

// FormatEducationSection formats education for CV.
func FormatEducationSection(edu []*domain.Education) string {
    // 200 lines of CV formatting
}

// File: internal/formatting/cv/experience.go
package cv

// FormatExperienceSection formats experience for CV.
func FormatExperienceSection(exp []*domain.Experience) string {
    // 250 lines of CV formatting
}

// File: internal/formatting/cv/skills.go
package cv

// FormatSkillsSection formats skills for CV.
func FormatSkillsSection(skills []*domain.Skill) string {
    // 180 lines of CV formatting
}
```

**Usage in context**:
```go
// File: intents/generate_cv/context.go
import "github.com/baphled/kariya/internal/formatting/cv"

func (c *GenerateCVContext) GenerateCV() (string, error) {
    educationSection := cv.FormatEducationSection(c.Education)
    experienceSection := cv.FormatExperienceSection(c.Experience)
    skillsSection := cv.FormatSkillsSection(c.Skills)
    
    return educationSection + experienceSection + skillsSection, nil
}
```

---

## Helper Classification Examples

### ✅ Keep in intent.go (Orchestration)

```go
// State transitions
func (i *MyIntent) setCancelled() { /* ... */ }
func (i *MyIntent) setCompleted(data *MyResult) { /* ... */ }

// Screen management
func (i *MyIntent) switchToListScreen() tea.Cmd { /* ... */ }
func (i *MyIntent) switchToDetailScreen(item *Item) tea.Cmd { /* ... */ }

// Modal management
func (i *MyIntent) showError(err error) tea.Cmd { /* ... */ }
```

**Guideline**: ≤5 helpers, orchestration only

---

### ❌ Move to context.go (Business Logic)

```go
// Validation
func (i *MyIntent) validateBurst(burst *domain.Burst) error { /* ... */ }

// Data operations
func (i *MyIntent) filterBurstsByTag(tag string) []*domain.Burst { /* ... */ }
func (i *MyIntent) sortBursts(bursts []*domain.Burst) []*domain.Burst { /* ... */ }

// Calculations
func (i *MyIntent) calculateDuration(burst *domain.Burst) time.Duration { /* ... */ }
func (i *MyIntent) calculateStats(bursts []*domain.Burst) *Stats { /* ... */ }

// Database operations
func (i *MyIntent) loadBursts() error { /* ... */ }
func (i *MyIntent) saveBurst(burst *domain.Burst) error { /* ... */ }
```

---

### ❌ Move to screens/ (Rendering)

```go
// Formatting for display
func (i *MyIntent) formatBurstItem(burst *domain.Burst) string { /* ... */ }
func (i *MyIntent) formatDate(t time.Time) string { /* ... */ }

// Rendering
func (i *MyIntent) renderHeader() string { /* ... */ }
func (i *MyIntent) renderList() string { /* ... */ }
func (i *MyIntent) renderDetail() string { /* ... */ }
```

---

### ❌ Move to internal/util/ (General Utilities)

```go
// String utilities
func (i *MyIntent) truncateString(s string, maxLen int) string { /* ... */ }
func (i *MyIntent) capitalizeFirst(s string) string { /* ... */ }

// Collection utilities
func (i *MyIntent) uniqueStrings(items []string) []string { /* ... */ }
func (i *MyIntent) filterEmpty(items []string) []string { /* ... */ }

// Date utilities
func (i *MyIntent) formatDate(t time.Time) string { /* ... */ }
func (i *MyIntent) parseDate(s string) (time.Time, error) { /* ... */ }
```

---

## File Size Targets

| File | Target Lines | Max Lines | Enforcement |
|------|-------------|-----------|-------------|
| `intent.go` | 200-400 | 600 | Check #18 (BLOCKING) |
| `context.go` | 100-200 | 500 | Warning only |
| `result.go` | 20-50 | 100 | Warning only |
| `constants.go` | 50-80 | 150 | Warning only |
| `messages.go` | 50-100 | 200 | Warning only |

**If intent.go >600 lines**: BLOCKED by Check #18

---

## Extraction Checklist

### Discovery
- [ ] Count lines in `intent.go` (`wc -l intent.go`)
- [ ] Identify all helper methods
- [ ] Classify each helper (business logic, rendering, orchestration, utility)
- [ ] Calculate line reduction potential

### Implementation
- [ ] Extract business logic helpers to `context.go`
- [ ] Extract rendering helpers to `screens/{feature}/`
- [ ] Extract general utilities to `internal/util/`
- [ ] Extract domain formatting to `internal/formatting/`
- [ ] Keep ≤5 orchestration helpers in `intent.go`

### Validation
- [ ] Run `make check-intent-architecture` (Check #18)
- [ ] Verify `intent.go` ≤600 lines (target 200-400)
- [ ] Verify `context.go` ≤500 lines
- [ ] Run unit tests
- [ ] Run integration tests

---

## Before/After Example

### BEFORE (All helpers in intent.go)

```
File: my_intent.go (1,200 lines - VIOLATION)

Intent implementation: 300 lines
Business logic helpers: 400 lines
Rendering helpers: 300 lines
Utility helpers: 200 lines

TOTAL: 1,200 lines in ONE file
VIOLATION: Exceeds 600 line limit
```

### AFTER (Helpers extracted)

```
File: intents/myfeature/intent.go (290 lines)
- Intent implementation: 240 lines
- Orchestration helpers: 50 lines (5 helpers)

File: intents/myfeature/context.go (180 lines)
- Business logic: 180 lines (moved from intent)

File: screens/myfeature/list_screen.go (120 lines)
- Rendering: 120 lines (moved from intent)

File: internal/util/date.go (40 lines)
File: internal/util/string.go (35 lines)
File: internal/util/collection.go (45 lines)
- Utilities: 120 lines (moved from intent)

TOTAL: 785 lines across 6 files (35% reduction)
COMPLIANCE: intent.go = 290 lines (PASSES)
```

---

## Common Pitfalls

### 1. Too Many Helpers in intent.go

```go
// ❌ VIOLATION - 15+ helpers
func (i *MyIntent) setCancelled() { /* ... */ }
func (i *MyIntent) setCompleted() { /* ... */ }
func (i *MyIntent) validateInput() { /* ... */ }
func (i *MyIntent) formatDate() { /* ... */ }
func (i *MyIntent) filterItems() { /* ... */ }
// ... 10 more helpers

// ✅ CORRECT - ≤5 orchestration helpers
func (i *MyIntent) setCancelled() { /* ... */ }
func (i *MyIntent) setCompleted() { /* ... */ }
func (i *MyIntent) switchToListScreen() { /* ... */ }
func (i *MyIntent) showError() { /* ... */ }
func (i *MyIntent) showSuccess() { /* ... */ }
```

**Guideline**: ≤5 private helpers in intent.go

---

### 2. Business Logic Not in Context

```go
// ❌ WRONG - Business logic in intent
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    burst := result.Data.(*domain.Burst)
    
    // WRONG: Validation in intent
    if burst.Title == "" {
        return i.showError(errors.New("title required"))
    }
    
    // WRONG: Database call in intent
    err := i.service.Save(i.getContext(), burst)
    // ...
}

// ✅ CORRECT - Delegate to context
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    burst := result.Data.(*domain.Burst)
    
    // Delegate to context
    if err := i.context.ValidateBurst(burst); err != nil {
        return i.showError(err)
    }
    
    ctx := i.GetContext()
    if err := i.context.SaveBurst(ctx, burst); err != nil {
        return i.showError(err)
    }
    
    return i.showSuccess("Burst saved")
}
```

---

### 3. Utilities Duplicated Across Intents

```go
// ❌ WRONG - Same utility in multiple intents
// File: intent1.go
func (i *Intent1) formatDate(t time.Time) string { /* ... */ }

// File: intent2.go
func (i *Intent2) formatDate(t time.Time) string { /* ... */ }

// File: intent3.go
func (i *Intent3) formatDate(t time.Time) string { /* ... */ }

// ✅ CORRECT - Centralized utility
// File: internal/util/date.go
func FormatDate(t time.Time) string { /* ... */ }

// Usage in all intents
import "github.com/baphled/kariya/internal/util"
formattedDate := util.FormatDate(burst.StartDate)
```

---

## Resources

### Documentation
- [Intent Migration Guide](INTENT_MIGRATION_TO_SUBDIRECTORY.md) - Full migration process
- [Screen Extraction Guide](SCREEN_EXTRACTION_GUIDE.md) - Extracting rendering
- [Intent Development Checklist](../checklists/INTENT_DEVELOPMENT_CHECKLIST.md) - Compliance

### Tools
- `make check-intent-architecture` - Validate file sizes (Check #18)
- `wc -l intents/myfeature/intent.go` - Count lines

---

**Last Updated**: 2026-01-22  
**Status**: Required for Check #18 compliance  
**Enforcement**: Automated by `check-intent-architecture.sh` (Check #18 blocks >600 lines)
