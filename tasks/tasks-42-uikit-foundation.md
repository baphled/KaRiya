# Task 41: UIKit Foundation - Standardized Component Library

## Overview
- **Goal**: Create a standardized UI component library (`internal/cli/uikit/`) with reusable primitives (buttons, inputs, badges, text) that eliminate ad-hoc styling and provide consistent theming
- **Time Estimate**: 2-3 days
- **Prerequisites**: Stable main branch, all tests passing

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

---

## Design Decisions (Confirmed)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Modal return type | `ModalResult[T]` | Type-safe, consistent with `IntentResult[T]` |
| Button design | Interactive models | Focus navigation for modal footers |
| Location | `internal/cli/uikit/` | Clean separation, dedicated UI toolkit |
| Migration | Incremental | Safer, testable in stages |
| Priority | Buttons/Inputs first | Foundational primitives |
| Inputs | Standalone + huh theming | Both standalone and form integration |
| Button scope | Full groups with navigation | Tab/arrow navigation between buttons |
| Theme system | Themes exclusively | Components require theme (use default if nil) |
| Testing | Unit + visual snapshots | Golden file comparison in testdata/ |

---

## Package Structure

```
internal/cli/uikit/
├── theme/
│   ├── theme.go           # Theme interface re-export + Default()
│   ├── aware.go           # ThemeAware embeddable base
│   └── theme_test.go
├── primitives/
│   ├── text.go            # Text, Title, Subtitle, Body, Muted, ErrorText
│   ├── text_test.go
│   ├── button.go          # Button with Primary/Secondary/Danger variants
│   ├── button_test.go
│   ├── button_group.go    # ButtonGroup with keyboard navigation
│   ├── button_group_test.go
│   ├── input.go           # Standalone themed text input
│   ├── input_test.go
│   ├── badge.go           # KeyBadge, StatusBadge, TagBadge
│   ├── badge_test.go
│   └── testdata/          # Golden files for snapshot tests
└── doc.go                 # Package documentation
```

---

## Phase 1: Theme Infrastructure

### Task 1.1: Theme Package

**Files to Create**:
- [ ] `internal/cli/uikit/theme/theme.go`
- [ ] `internal/cli/uikit/theme/aware.go`
- [ ] `internal/cli/uikit/theme/theme_test.go`

**TDD Checklist**:

#### RED Phase
- [ ] Test file created: `internal/cli/uikit/theme/theme_test.go`
- [ ] Tests written for:
  - [ ] `Default()` returns a valid theme
  - [ ] `Aware.SetTheme()` stores theme
  - [ ] `Aware.Theme()` returns stored theme
  - [ ] All color getters return valid lipgloss.Color
  - [ ] Nil theme uses default (forgiving behavior)
- [ ] Tests **FAIL** with error:
  ```
  [Paste actual error here]
  ```

#### GREEN Phase
- [ ] `theme.go` implemented:
  - [ ] Re-export `themes.Theme` interface
  - [ ] `Default()` function returns `themes.NewDefaultTheme()`
- [ ] `aware.go` implemented:
  - [ ] `Aware` struct with `theme` field
  - [ ] `SetTheme(t Theme)` method
  - [ ] `Theme() Theme` method (returns default if nil)
  - [ ] Color getter methods: `PrimaryColor()`, `SecondaryColor()`, `AccentColor()`, `ErrorColor()`, `SuccessColor()`, `WarningColor()`, `BorderColor()`, `BackgroundColor()`, `MutedColor()`
- [ ] Tests now **PASS**

#### REFACTOR Phase
- [ ] Remove any duplication
- [ ] Ensure consistent naming
- [ ] Add doc comments

**Acceptance Criteria**:
- [ ] `theme.Default()` returns usable theme
- [ ] `Aware` can be embedded in other structs
- [ ] All color getters work with nil theme (use default)
- [ ] Tests pass with race detector

**Estimated LOC**: ~150

---

## Phase 2: Text Primitives

### Task 1.2: Text Component

**Files to Create**:
- [ ] `internal/cli/uikit/primitives/text.go`
- [ ] `internal/cli/uikit/primitives/text_test.go`
- [ ] `internal/cli/uikit/primitives/testdata/text_*.golden`

**TDD Checklist**:

#### RED Phase
- [ ] Test file created: `internal/cli/uikit/primitives/text_test.go`
- [ ] Tests written for:
  - [ ] `NewText()` creates text with content and theme
  - [ ] `Style()` fluent method works
  - [ ] `Bold()` fluent method works
  - [ ] `Width()` fluent method constrains output
  - [ ] `Render()` produces styled output
  - [ ] Each `TextStyle` renders with correct colors
  - [ ] Convenience constructors: `Title()`, `Subtitle()`, `Body()`, `Muted()`, `ErrorText()`
- [ ] Snapshot tests for visual verification
- [ ] Tests **FAIL**

#### GREEN Phase
- [ ] `text.go` implemented:
  - [ ] `TextStyle` enum: `TextBody`, `TextTitle`, `TextSubtitle`, `TextMuted`, `TextError`, `TextSuccess`, `TextWarning`
  - [ ] `Text` struct embedding `theme.Aware`
  - [ ] Fluent API: `Style()`, `Bold()`, `Width()`
  - [ ] `Render()` method with lipgloss styling
  - [ ] Convenience constructors
- [ ] Tests now **PASS**

#### REFACTOR Phase
- [ ] Extract style application to helper
- [ ] Add doc comments

**Acceptance Criteria**:
- [ ] All text styles render with theme colors
- [ ] Fluent API chains correctly
- [ ] Width constraint works
- [ ] Snapshot tests pass

**Estimated LOC**: ~200

---

## Phase 3: Button Components

### Task 1.3: Button Component

**Files to Create**:
- [ ] `internal/cli/uikit/primitives/button.go`
- [ ] `internal/cli/uikit/primitives/button_test.go`
- [ ] `internal/cli/uikit/primitives/testdata/button_*.golden`

**TDD Checklist**:

#### RED Phase
- [ ] Tests written for:
  - [ ] `NewButton()` creates button with label and theme
  - [ ] `Variant()` fluent method (Primary, Secondary, Danger)
  - [ ] `Focused()` fluent method changes styling
  - [ ] `Disabled()` fluent method changes styling
  - [ ] `Width()` fluent method constrains output
  - [ ] `Render()` produces styled button
  - [ ] Convenience constructors: `PrimaryButton()`, `SecondaryButton()`, `DangerButton()`
- [ ] Snapshot tests for all variants and states
- [ ] Tests **FAIL**

#### GREEN Phase
- [ ] `button.go` implemented:
  - [ ] `ButtonVariant` enum: `ButtonPrimary`, `ButtonSecondary`, `ButtonDanger`
  - [ ] `Button` struct embedding `theme.Aware`
  - [ ] Fluent API: `Variant()`, `Focused()`, `Disabled()`, `Width()`
  - [ ] `Render()` method with:
    - Rounded border
    - Theme-appropriate colors per variant
    - Focus state: thick border, accent color, bold
    - Disabled state: faint, muted colors
  - [ ] Convenience constructors
- [ ] Tests now **PASS**

#### REFACTOR Phase
- [ ] Extract border styling to helper
- [ ] Add doc comments

**Acceptance Criteria**:
- [ ] All variants render correctly
- [ ] Focus state visually distinct (thick border, accent color)
- [ ] Disabled state visually distinct (faint)
- [ ] Snapshot tests pass

**Estimated LOC**: ~250

---

### Task 1.4: Button Group Component

**Files to Create**:
- [ ] `internal/cli/uikit/primitives/button_group.go`
- [ ] `internal/cli/uikit/primitives/button_group_test.go`
- [ ] `internal/cli/uikit/primitives/testdata/button_group_*.golden`

**TDD Checklist**:

#### RED Phase
- [ ] Tests written for:
  - [ ] `NewButtonGroup()` creates empty group
  - [ ] `Add()` adds button to group
  - [ ] `AddPrimary()`, `AddSecondary()`, `AddDanger()` convenience methods
  - [ ] `Horizontal()` sets layout direction
  - [ ] `FocusIndex()` returns current focus
  - [ ] `FocusedLabel()` returns focused button's label
  - [ ] `FocusNext()` moves focus forward (wraps)
  - [ ] `FocusPrev()` moves focus backward (wraps)
  - [ ] `FocusFirst()` moves to first button
  - [ ] `FocusLast()` moves to last button
  - [ ] `Update()` handles keyboard: Tab, Shift+Tab, Left, Right, h, l
  - [ ] `Render()` produces layout with correct focus
- [ ] Snapshot tests for horizontal and vertical layouts
- [ ] Tests **FAIL**

#### GREEN Phase
- [ ] `button_group.go` implemented:
  - [ ] `ButtonGroup` struct embedding `theme.Aware`
  - [ ] `buttons []*Button` slice
  - [ ] `focusedIndex int` field
  - [ ] `horizontal bool` field (default true)
  - [ ] Navigation methods with wrap-around
  - [ ] `Update()` handles:
    - `tab`, `right`, `l` → FocusNext()
    - `shift+tab`, `left`, `h` → FocusPrev()
    - `home` → FocusFirst()
    - `end` → FocusLast()
  - [ ] `Render()` joins buttons horizontally or vertically
- [ ] Tests now **PASS**

#### REFACTOR Phase
- [ ] Extract key handling to switch
- [ ] Add doc comments

**Acceptance Criteria**:
- [ ] Navigation wraps at boundaries
- [ ] Only focused button shows focus state
- [ ] Horizontal/vertical layouts work
- [ ] Keyboard handling is comprehensive
- [ ] Snapshot tests pass

**Estimated LOC**: ~300

---

## Phase 4: Input Component

### Task 1.5: Input Component

**Files to Create**:
- [ ] `internal/cli/uikit/primitives/input.go`
- [ ] `internal/cli/uikit/primitives/input_test.go`
- [ ] `internal/cli/uikit/primitives/testdata/input_*.golden`

**TDD Checklist**:

#### RED Phase
- [ ] Tests written for:
  - [ ] `NewInput()` creates input with theme
  - [ ] `Label()` fluent method sets label
  - [ ] `Placeholder()` fluent method sets placeholder
  - [ ] `Value()` fluent method sets initial value
  - [ ] `Error()` fluent method sets error message
  - [ ] `Width()` fluent method constrains output
  - [ ] `Focus()` returns command to focus
  - [ ] `Blur()` removes focus
  - [ ] `Update()` handles keyboard input
  - [ ] `View()` renders label + input + error
  - [ ] `GetValue()` returns current value
- [ ] Snapshot tests for states: empty, with value, focused, with error
- [ ] Tests **FAIL**

#### GREEN Phase
- [ ] `input.go` implemented:
  - [ ] `Input` struct embedding `theme.Aware`
  - [ ] Wraps `textinput.Model` from bubbles
  - [ ] Fields: `label`, `placeholder`, `error`, `width`
  - [ ] Fluent API methods
  - [ ] `View()` renders:
    - Label above (styled with theme)
    - Input field with themed border
    - Error below (red styling, if set)
    - Focus state: accent border
  - [ ] `Update()` delegates to textinput.Model
- [ ] Tests now **PASS**

#### REFACTOR Phase
- [ ] Extract rendering sections
- [ ] Add doc comments

**Acceptance Criteria**:
- [ ] Input wraps bubbles textinput correctly
- [ ] Label renders above input
- [ ] Error renders below input (red)
- [ ] Focus state shows accent border
- [ ] Typing works correctly
- [ ] Snapshot tests pass

**Estimated LOC**: ~300

---

## Phase 5: Badge Component

### Task 1.6: Badge Component

**Files to Create**:
- [ ] `internal/cli/uikit/primitives/badge.go`
- [ ] `internal/cli/uikit/primitives/badge_test.go`
- [ ] `internal/cli/uikit/primitives/testdata/badge_*.golden`

**TDD Checklist**:

#### RED Phase
- [ ] Tests written for:
  - [ ] `NewBadge()` creates badge with label and theme
  - [ ] `Value()` fluent method sets value (for key badges)
  - [ ] `Variant()` fluent method sets variant
  - [ ] `Render()` produces styled output
  - [ ] `KeyBadge()` convenience constructor formats "key → action"
  - [ ] `StatusBadge()` convenience constructor
  - [ ] `TagBadge()` convenience constructor
- [ ] Snapshot tests for all variants
- [ ] Tests **FAIL**

#### GREEN Phase
- [ ] `badge.go` implemented:
  - [ ] `BadgeVariant` enum: `BadgeDefault`, `BadgeKey`, `BadgeStatus`, `BadgeTag`
  - [ ] `Badge` struct embedding `theme.Aware`
  - [ ] Fields: `label`, `value`, `variant`
  - [ ] Fluent API methods
  - [ ] `Render()` with variant-specific styling:
    - Key badge: `[key] action` format
    - Status badge: colored based on status
    - Tag badge: pill-style
  - [ ] Convenience constructors
- [ ] Tests now **PASS**

#### REFACTOR Phase
- [ ] Unify styling patterns
- [ ] Add doc comments

**Acceptance Criteria**:
- [ ] All variants render correctly
- [ ] Key badges show key-action format
- [ ] Styling consistent with existing `components.KeyBadge`
- [ ] Snapshot tests pass

**Estimated LOC**: ~150

---

## Phase 6: Package Documentation

### Task 1.7: Documentation

**Files to Create**:
- [ ] `internal/cli/uikit/doc.go`
- [ ] `internal/cli/uikit/theme/doc.go`
- [ ] `internal/cli/uikit/primitives/doc.go`

**Content**:
- [ ] Package overview
- [ ] Usage examples
- [ ] Design principles
- [ ] Migration guide from old components

**Estimated LOC**: ~100

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit is atomic (ONE logical change)

---

## Acceptance Criteria

### Code Quality
- [ ] All tests pass
- [ ] Coverage >85% for new code
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions

### Components
- [ ] Theme infrastructure works (Aware embeddable)
- [ ] Text renders all semantic styles
- [ ] Buttons render all variants with focus/disabled states
- [ ] ButtonGroup handles keyboard navigation
- [ ] Input wraps textinput with theming
- [ ] Badges render all variants

### API Design
- [ ] All components use fluent builder pattern
- [ ] All components require theme (use default if nil)
- [ ] Consistent naming conventions
- [ ] Clear separation of concerns

### Testing
- [ ] Unit tests for all public methods
- [ ] Snapshot tests for visual verification
- [ ] Golden files in testdata/

---

## Rollback Plan

This is a new package with no dependencies on existing code:
1. Delete `internal/cli/uikit/` directory
2. No other changes required

---

## Future Phases (Not in This Task)

After Phase 1 completes, subsequent tasks will add:

### Phase 2: Containers
- Box (Card, Panel, Modal frame)
- Stack (VStack, HStack)
- Overlay utilities

### Phase 3: Modals
- ModalResult[T] type
- BaseModal
- ConfirmModal
- FormModal[T]
- StatusModal

### Phase 4: Integration
- Huh theme adapter
- Widget footer
- Migration of existing components

---

## Estimated Totals

| Component | Est. LOC | Tests |
|-----------|----------|-------|
| theme/ | 150 | ~20 |
| text.go | 200 | ~15 |
| button.go | 250 | ~20 |
| button_group.go | 300 | ~25 |
| input.go | 300 | ~20 |
| badge.go | 150 | ~15 |
| doc.go files | 100 | - |
| **Total** | **~1,450** | **~115** |
