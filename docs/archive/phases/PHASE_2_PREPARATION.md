---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Phase 2: Visual Consistency and Layout Standardization - Preparation Guide

**Status**: Ready for Implementation
**Date Prepared**: 2025-12-30
**Prerequisites**: Phase 1 Complete ✅

---

## Overview

Phase 2 will implement reusable visual components to standardize layout across all TUI screens. This phase builds directly on Phase 1's navigation system foundation.

### Phase 2 Goals

1. **Create NavigationMenu Component** (Tasks 3.1-3.5)
   - Reusable menu with configurable items
   - Support for horizontal and vertical layouts
   - Full keyboard navigation
   - Responsive rendering

2. **Create Header Component** (Tasks 4.1-4.3)
   - Consistent screen headers
   - Title, subtitle, breadcrumb support
   - Navigation context display

3. **Create Footer Component** (Tasks 4.4-4.6)
   - Status message display
   - Mode/context indicators
   - Responsive sizing

4. **Create ListItem Component** (Tasks 5.1-5.3)
   - Consistent item rendering
   - Selection indicators
   - Text truncation with ellipsis

### Phase 2 Estimated Effort
- **Duration**: 3-4 hours
- **Tasks**: 12 implementation tasks
- **Tests**: 150+ new tests expected
- **Commits**: 4-5 expected

---

## Phase 1 Recap

### What Was Built
- ✅ Navigation system (19 keys, centralized)
- ✅ Help footer component (100 tests)
- ✅ Keyboard standardization (Escape key)
- ✅ Comprehensive documentation (1,600+ lines)

### Current State
- ✅ 607+ tests passing
- ✅ 0 race conditions
- ✅ 80%+ code coverage
- ✅ All Phase 1 tasks complete

### Foundation for Phase 2
- Navigation constants ready to use
- Help footer component available
- Clear architectural patterns established
- Test infrastructure in place

---

## Phase 2 Task Breakdown

### Section 3: NavigationMenu Component (Tasks 3.1-3.5)

#### Task 3.1: Create navigation_menu.go
**Objective**: Create reusable menu component
**Key Features**:
- NavigationMenuModel struct
- Menu item structure (Label, Shortcut, Description)
- Configurable layouts (horizontal, vertical)
- Focus management

**Expected Structure**:
```go
type MenuItem struct {
    Label       string
    Shortcut    string
    Description string
    Action      func()
}

type NavigationMenuModel struct {
    items      []MenuItem
    focused    int
    layout     LayoutType
    width      int
    height     int
}
```

#### Task 3.2: Implement menu item structure
**Focus**: Menu configuration and item management
- NewNavigationMenu() constructor
- AddItem() method
- SetLayout() (horizontal/vertical)
- GetSelectedItem() method

#### Task 3.3: Implement keyboard navigation
**Focus**: User input handling
- ↑/↓ for vertical navigation
- ←/→ for horizontal navigation
- Enter to select
- Escape to cancel

#### Tasks 3.4-3.5: Testing (20+ tests)
**Focus**: Layout and edge cases
- Layout rendering verification
- Keyboard navigation tests
- Single/multiple items
- Responsive sizing (40-300+ chars)

---

### Section 4: Header Component (Tasks 4.1-4.3)

#### Task 4.1: Create header.go
**Objective**: Consistent screen headers
**Features**:
- Title and subtitle support
- Breadcrumb navigation
- Status indicators
- Responsive sizing

**Expected Structure**:
```go
type HeaderModel struct {
    title      string
    subtitle   string
    breadcrumb []string
    width      int
    style      lipgloss.Style
}
```

#### Task 4.2: Implement breadcrumb display
**Focus**: Navigation context
- Breadcrumb parsing (e.g., "Home > List > Details")
- Clickable navigation (for future)
- Truncation for narrow terminals

#### Task 4.3: Testing (10+ tests)
**Focus**: Title, subtitle, breadcrumb rendering
- Responsive at various widths
- Breadcrumb display correctness

---

### Section 5: Footer Component (Tasks 4.4-4.6)

#### Task 4.4: Create footer.go
**Objective**: Consistent screen footers
**Features**:
- Status message display
- Mode/context indicators
- Responsive sizing

**Expected Structure**:
```go
type FooterModel struct {
    status    string
    mode      string
    context   string
    width     int
    style     lipgloss.Style
}
```

#### Task 4.5: Implement status display
**Focus**: Information presentation
- Status messages (e.g., "3/10 events")
- Mode display (e.g., "Capture Mode: Timeline")
- Context information

#### Task 4.6: Testing (10+ tests)
**Focus**: Status and mode display
- Various status message lengths
- Responsive sizing

---

### Section 6: ListItem Component (Tasks 5.1-5.3)

#### Task 5.1: Create list_item.go
**Objective**: Consistent list item rendering
**Features**:
- Title and subtitle
- Metadata fields
- Status indicators
- Selection highlighting

**Expected Structure**:
```go
type ListItemModel struct {
    title       string
    subtitle    string
    metadata    map[string]string
    selected    bool
    focused     bool
    width       int
    style       lipgloss.Style
}
```

#### Task 5.2: Implement styling and truncation
**Focus**: Professional rendering
- Text truncation with ellipsis
- Selected/focused indicators
- Metadata display
- Responsive layout

#### Task 5.3: Testing (15+ tests)
**Focus**: Rendering at various sizes
- Long title truncation
- Selected/focused states
- Metadata display
- Responsive behavior

---

### Section 7: Integration Tasks (Tasks 6.1-6.4)

After components are created:
- Integrate HelpFooter into form.go
- Integrate HelpFooter into list.go
- Integrate HelpFooter into metadata_review.go
- Integrate HelpFooter into remaining models

**Note**: These tasks depend on successful creation of components and tests

---

### Section 8: Integration Testing (Tasks 7.1-7.3)

After all components integrated:
- Test Escape key across all models
- Test vim navigation (hjkl) in forms/lists
- Test help footer context-awareness

---

## Implementation Strategy

### Step 1: Create Components (Parallel)
- Create header.go (Task 4.1-4.3)
- Create footer.go (Task 4.4-4.6)
- Create list_item.go (Task 5.1-5.3)
- Create navigation_menu.go (Task 3.1-3.5)

### Step 2: Write Tests (Parallel)
- Add tests for each component
- Test rendering at various widths
- Test edge cases

### Step 3: Integrate into Models (Sequential)
- Update models to use new components
- Verify existing tests still pass
- Add integration tests

### Step 4: Final Testing
- Run full test suite
- Check race conditions
- Verify coverage

---

## Code Patterns to Follow

### From Phase 1, use these patterns:

**Model Interface** (from help_footer.go):
```go
type ComponentModel struct {
    width   int
    height  int
}

func (m ComponentModel) Init() tea.Cmd { return nil }
func (m ComponentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m ComponentModel) View() string
```

**Responsive Rendering**:
```go
func (m ComponentModel) View() string {
    if m.width < 40 {
        return m.renderCompact()
    }
    return m.render()
}
```

**Testing Pattern**:
```go
It("should render on various terminal widths", func() {
    for width := 40; width <= 200; width += 20 {
        m.SetWidth(width)
        view := m.View()
        Expect(view).NotTo(BeEmpty())
    }
})
```

---

## Testing Requirements

### For Each Component
- ✅ Construction and initialization
- ✅ BubbleTea interface compliance
- ✅ View rendering
- ✅ Width management/responsiveness
- ✅ State changes
- ✅ Edge cases

### Overall
- ✅ 150+ new tests expected
- ✅ 100% pass rate required
- ✅ 80%+ coverage maintenance
- ✅ 0 race conditions allowed

---

## Success Criteria for Phase 2

| Criterion | Target | How to Verify |
|-----------|--------|---------------|
| Tasks Completed | 12/12 | All tasks marked done |
| Tests Passing | 100% | `go test ./...` ✅ |
| Race Conditions | 0 | `go test -race ./...` ✅ |
| Coverage | 80%+ | `go test -cover ./...` ✅ |
| Components | 4 | ls internal/cli/components/ |
| Integration | Complete | All models updated |
| Documentation | Updated | TUI_STANDARDS.md revised |

---

## Files to Create

### Components (4 new files)
1. `internal/cli/components/navigation_menu.go` (~250 lines)
2. `internal/cli/components/header.go` (~200 lines)
3. `internal/cli/components/footer.go` (~200 lines)
4. `internal/cli/components/list_item.go` (~250 lines)

### Tests (4 new files)
1. `internal/cli/components/navigation_menu_test.go` (~400 lines)
2. `internal/cli/components/header_test.go` (~300 lines)
3. `internal/cli/components/footer_test.go` (~300 lines)
4. `internal/cli/components/list_item_test.go` (~350 lines)

### Integration (multiple files)
- Update model files to use new components
- Update model test files
- Add integration tests

---

## Development Checklist

### Before Starting Phase 2
- [ ] Phase 1 all tests passing ✅
- [ ] Phase 1 documentation complete ✅
- [ ] Phase 1 code committed ✅
- [ ] Working directory clean ✅

### During Phase 2
- [ ] Create navigation_menu component
- [ ] Create header component
- [ ] Create footer component
- [ ] Create list_item component
- [ ] Write comprehensive tests (150+)
- [ ] Integrate into models
- [ ] Run full test suite
- [ ] Verify race detection
- [ ] Check coverage

### After Phase 2
- [ ] All 607+ existing tests still passing
- [ ] 150+ new tests passing
- [ ] 0 race conditions
- [ ] 80%+ coverage maintained
- [ ] Code committed and clean
- [ ] Documentation updated

---

## Estimated Timeline

| Task | Duration | Start | End |
|------|----------|-------|-----|
| NavigationMenu | 45 min | 0:00 | 0:45 |
| Header | 30 min | 0:45 | 1:15 |
| Footer | 30 min | 1:15 | 1:45 |
| ListItem | 45 min | 1:45 | 2:30 |
| Integration | 60 min | 2:30 | 3:30 |
| Testing | 30 min | 3:30 | 4:00 |
| **Total** | **~4 hours** | | |

---

## Resources Available

### From Phase 1
- `internal/cli/navigation/` - Use constants and help generation
- `internal/cli/styles/` - Use predefined styles
- `internal/cli/components/help_footer.go` - Reference implementation
- Documentation: `docs/TUI_STANDARDS.md`, `docs/TUI_DEVELOPER_GUIDE.md`

### Test Framework
- Ginkgo v2 for BDD testing
- Gomega for assertions
- Existing test patterns in help_footer_test.go

### Reference Implementations
- `help_footer.go` - BubbleTea component pattern
- Form models - Navigation and keyboard handling
- List models - Item rendering patterns

---

## Next Session Preparation

To maximize efficiency when starting Phase 2:

1. Review this preparation guide
2. Read `docs/TUI_DEVELOPER_GUIDE.md` (especially "Creating Components")
3. Review `help_footer.go` as reference implementation
4. Review `help_footer_test.go` for test patterns
5. Have Phase 1 completion report available

---

## Questions to Answer Before Starting

1. Should NavigationMenu support both horizontal and vertical layouts in phase 2?
   → Yes, specified in tasks 3.1-3.2

2. Should Header have clickable breadcrumbs in phase 2?
   → No, that's future work; just display for now

3. Should ListItem support icons/visual indicators?
   → Yes, specified in task 5.1 (status indicators)

4. When should integration with models happen?
   → After all 4 components are complete and tested (Tasks 6.1-6.4)

---

## Phase 2 Success Definition

Phase 2 will be considered complete when:

✅ All 4 components created and tested
✅ 150+ new tests passing (100% pass rate)
✅ 0 race conditions detected
✅ 80%+ code coverage maintained
✅ HelpFooter integrated into key models
✅ Integration tests passing
✅ All code committed with clear messages
✅ Documentation updated with new components

---

**This preparation guide ensures Phase 2 can start immediately with clear goals and implementation strategy.**

**Next Session**: Begin Phase 2 with NavigationMenu component implementation.

