# Task 5.0: Metadata Review View Container Pattern Audit

**Date**: 2026-01-01
**Task**: Make metadata review consistent with our view container pattern
**Status**: Audit Complete - Ready for Implementation

## Executive Summary

The `MetadataReviewModel` View() method does not follow the standardized view container pattern established in other models (FactListModel, MetadataEditorModel, BurstListModel, FactListModel). This audit identifies the specific inconsistencies and proposes remediation steps.

**Current State**: Manual content array building with inline rendering
**Target State**: Use ScreenContainer with ListContainer pattern (like FactListModel)
**Impact**: 3 files to modify, 1 pattern to establish

---

## Current Implementation Analysis

### View() Method Structure (Lines 369-432)

**Current Pattern:**
```go
func (m *MetadataReviewModel) View() string {
    var content []string

    // Manual header construction
    header := components.NewHeader(...)
    content = append(content, header.View())

    // Manual status bar
    content = append(content, styles.InputHint.Render(...))

    // Manual error handling
    if m.err != nil {
        content = append(content, styles.ErrorText.Render(...))
    }

    // Manual list rendering
    for i, event := range m.events {
        // Inline rendering
    }

    // Manual footer handling
    m.helpFooter.SetWidth(...)
    footer := m.helpFooter.View()
    content = append(content, footer)

    return strings.Join(content, "\n")
}
```

**Issues Identified:**

1. ❌ **No Container Component**: Does not use `ScreenContainer` or `ListContainer`
2. ❌ **Manual Content Assembly**: Uses `[]string` slice with manual append operations
3. ❌ **Inconsistent Error Display**: Uses `styles.ErrorText` instead of `styles.ErrorBox`
4. ❌ **Manual Footer Width Management**: Sets footer width manually instead of delegating to container
5. ❌ **No Separation of Concerns**: Header, content, and footer logic mixed in View()
6. ❌ **Helper Methods Not Used**: Has `renderEventItem()`, `renderSelectedEvent()`, `renderExpandedEvent()` but doesn't compose with container

---

## Standardized Pattern Reference

### FactListModel Pattern (COMPLIANT ✅)

**Container Structure:**
```go
func (m *FactListModel) View() string {
    screen := components.NewScreenContainer(m.renderContent())
    screen.SetHeader(m.renderHeader())
    screen.SetFooter(m.renderFooter())
    return screen.Render()
}
```

**Key Characteristics:**
- Uses `ScreenContainer` for consistent layout
- Delegates width management to container
- Header and footer are set via chainable methods
- Content is pre-rendered before passing to container
- Error handling uses `styles.ErrorBox`
- Returns container's Render() output

### MetadataEditorModel Pattern (COMPLIANT ✅)

**Container Structure:**
- Uses `FormFieldContainer` for form field rendering
- Implements field-level errors via `FormFieldContainer.SetError()`
- Uses `HelpFooter()` for consistent footer display
- Proper separation of header, form content, and footer

### BurstListModel Pattern (COMPLIANT ✅)

- Uses `ScreenContainer` pattern
- Implements error display with `styles.ErrorBox`
- Proper footer handling via container

---

## Identified Inconsistencies

### 1. Error Display (Line 410)

**Current:**
```go
if m.err != nil {
    content = append(content, styles.ErrorText.Render(fmt.Sprintf("Error loading events: %v", m.err)))
    return strings.Join(content, "\n")
}
```

**Issue**: Uses `ErrorText` (field-level) instead of `ErrorBox` (model-level)

**Standard**: Model-level errors should use `styles.ErrorBox`

**Impact**: Visual inconsistency with other models

---

### 2. Container Pattern Not Used (Lines 369-432)

**Current:**
```go
func (m *MetadataReviewModel) View() string {
    var content []string
    // ... manual assembly
    return strings.Join(content, "\n")
}
```

**Issue**: Does not use `ScreenContainer` like other list-based models

**Standard**: List-based views should use `ScreenContainer` with helper methods

**Impact**:
- No consistent layout management
- Footer width management is manual
- Harder to maintain consistency as styles change

---

### 3. Footer Width Management (Lines 426-427)

**Current:**
```go
m.helpFooter.SetWidth(styles.MaxWidth(80))
footer := m.helpFooter.View()
```

**Issue**: Manual width management instead of delegating to container

**Standard**: Container should manage footer width automatically

**Impact**: Inconsistent footer sizing across models

---

### 4. No Helper Method Separation

**Current State:**
- Has `renderEventItem()`, `renderSelectedEvent()`, `renderExpandedEvent()`
- But View() method directly builds content array instead of using container

**Standard**: View() should be minimal, delegating to container and helper methods

**Impact**: View() method is 64 lines long (hard to read)

---

## Proposed Changes

### Change 1: Refactor View() Method

**File**: `internal/cli/models/metadata_review.go`

**Current (Lines 369-432):**
- Manual content array building
- 64 lines of implementation

**Proposed:**
- Use `ScreenContainer` pattern
- Delegate to `renderContent()` helper
- ~15 lines of implementation

**Implementation:**
```go
func (m *MetadataReviewModel) View() string {
    screen := components.NewScreenContainer(m.renderContent())
    screen.SetHeader(m.renderHeader())
    screen.SetFooter(m.renderFooter())
    return screen.Render()
}
```

---

### Change 2: Create renderHeader() Helper

**File**: `internal/cli/models/metadata_review.go`

**New Method:**
```go
func (m *MetadataReviewModel) renderHeader() string {
    header := components.NewHeader("Metadata Review", m.width)
    header.SetBreadcrumbs(m.breadcrumbs)
    return header.View()
}
```

**Purpose:**
- Separate header rendering logic
- Consistent with FactListModel pattern
- Easier to test and maintain

---

### Change 3: Create renderContent() Helper

**File**: `internal/cli/models/metadata_review.go`

**New Method:**
```go
func (m *MetadataReviewModel) renderContent() string {
    var content []string

    // Status bar
    statusText := fmt.Sprintf("Showing %d events | Filter: %s | Sort: %s | Press 'f' to filter, 's' to sort",
        len(m.events), m.filterMode, m.sortBy)
    content = append(content, styles.InputHint.Render(statusText))
    content = append(content, "")

    // Error handling - use ErrorBox instead of ErrorText
    if m.err != nil {
        content = append(content, styles.ErrorBox.Render(fmt.Sprintf("Error loading events: %v", m.err)))
        return strings.Join(content, "\n")
    }

    // Empty state
    if len(m.events) == 0 {
        content = append(content, styles.InfoText.Render("No events to review. Start capturing events to improve their metadata."))
        return strings.Join(content, "\n")
    }

    // Events list
    for i, event := range m.events {
        var eventContent string

        if i == m.selectedIdx {
            eventContent = m.renderSelectedEvent(event, i)
        } else {
            eventContent = m.renderEventItem(event, i)
        }

        content = append(content, eventContent)

        // Show expanded view if selected
        if i == m.expandedIdx {
            content = append(content, m.renderExpandedEvent(event))
        }
    }

    return strings.Join(content, "\n")
}
```

**Purpose:**
- Separate content rendering logic
- Makes View() method clean and readable
- Follows ScreenContainer pattern

---

### Change 4: Create renderFooter() Helper

**File**: `internal/cli/models/metadata_review.go`

**New Method:**
```go
func (m *MetadataReviewModel) renderFooter() string {
    return m.helpFooter.View()
}
```

**Purpose:**
- Consistent with ScreenContainer pattern
- Allows future footer customization
- Removes manual width management

---

### Change 5: Update helpFooter Initialization

**File**: `internal/cli/models/metadata_review.go`

**Current (Line 46):**
```go
helpFooter: components.NewHelpFooter("metadata_review", 80),
```

**Proposed:**
```go
helpFooter: components.NewHelpFooter("metadata_review", 0), // Width managed by ScreenContainer
```

**Or remove and initialize in Init():**
```go
func (m *MetadataReviewModel) Init() tea.Cmd {
    m.helpFooter = components.NewHelpFooter("metadata_review", m.width)
    return nil
}
```

---

### Change 6: Update WindowSizeMsg Handler

**File**: `internal/cli/models/metadata_review.go`

**Current (Lines 379-382):**
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    m.helpFooter.SetWidth(msg.Width)
    return m, nil
```

**Proposed:**
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    m.helpFooter.SetWidth(msg.Width) // Keep for footer consistency
    return m, nil
```

**Note**: Keep SetWidth() call for footer, as ScreenContainer may not directly manage HelpFooter width

---

## Files to Modify

| File | Changes | Lines | Priority |
|------|---------|-------|----------|
| `internal/cli/models/metadata_review.go` | Refactor View(), add helpers, fix error display | 369-432 + new helpers | HIGH |
| `internal/cli/models/metadata_review_test.go` | Update View() tests if they exist | TBD | MEDIUM |
| `docs/guides/ERROR_HANDLING_GUIDE.md` | Add MetadataReviewModel as compliant example | TBD | LOW |

---

## Compliance Checklist

### Before Refactoring
- ❌ Uses ScreenContainer
- ❌ Has renderHeader() helper
- ❌ Has renderContent() helper
- ❌ Has renderFooter() helper
- ❌ Uses ErrorBox for model-level errors
- ❌ Minimal View() method (< 20 lines)
- ❌ Proper separation of concerns

### After Refactoring
- ✅ Uses ScreenContainer
- ✅ Has renderHeader() helper
- ✅ Has renderContent() helper
- ✅ Has renderFooter() helper
- ✅ Uses ErrorBox for model-level errors
- ✅ Minimal View() method (< 20 lines)
- ✅ Proper separation of concerns

---

## Testing Strategy

### Existing Tests to Verify
1. View() output format remains consistent
2. Error display shows correctly with ErrorBox
3. Header, content, and footer are all rendered
4. Layout respects window size changes

### New Tests to Add
1. Test renderHeader() output
2. Test renderContent() output
3. Test renderFooter() output
4. Test ScreenContainer integration

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-----------|
| View output changes | Low | Medium | Run visual tests after refactoring |
| Layout inconsistency | Low | Medium | Compare with FactListModel pattern |
| Error display issues | Low | Low | Verify ErrorBox styling matches other models |
| Footer sizing issues | Low | Low | Test with various terminal widths |

---

## Summary

The `MetadataReviewModel` View() method needs refactoring to align with the established view container pattern. The main changes are:

1. **Use ScreenContainer** instead of manual content array
2. **Create helper methods** (renderHeader, renderContent, renderFooter)
3. **Fix error display** to use ErrorBox instead of ErrorText
4. **Improve code organization** for better maintainability

**Estimated Implementation Time**: 30-45 minutes
**Estimated Testing Time**: 15-20 minutes
**Total Effort**: ~1 hour

**Readiness**: ✅ Ready for implementation after user confirmation

---

## Implementation Notes

### Key Points
1. Keep existing helper methods (renderEventItem, renderSelectedEvent, renderExpandedEvent)
2. Ensure ScreenContainer handles layout properly
3. Verify error messages are still clear and actionable
4. Test with different terminal widths

### Backward Compatibility
- No public API changes
- Internal refactoring only
- View() output should remain visually identical

---

## Next Steps

1. ✅ User reviews and confirms this audit
2. Implement proposed changes
3. Run full test suite
4. Verify visual output matches original
5. Update AGENTS.md with completion status

