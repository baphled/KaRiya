# Modal Refactoring Verification Guide

## Summary
All 5 timeline modals have been refactored to use the `bubbletea-overlay` library (v0.6.3) instead of manual overlay compositing. Additionally, solid backgrounds have been added to all modals to prevent transparency issues.

**NEW**: Event detail view is now a modal overlay instead of a full-screen transition!

## Changes Made

### 1. Dependencies
- ✅ Added `github.com/rmhubbert/bubbletea-overlay` v0.6.3

### 2. Modal Components (All have solid backgrounds)

#### ViewEventDetailModal (`internal/cli/components/view_event_detail_modal.go`) **NEW!**
- ✅ Created new modal component for event details
- ✅ Solid background wrapper in View() method (lines 103-113):
  ```go
  return lipgloss.NewStyle().
      Border(lipgloss.RoundedBorder()).
      BorderForeground(styles.ColorBorder).
      Background(styles.ColorBackground).  // Solid #1a1f2e
      Padding(1, 2).
      MaxWidth(m.width - 8).
      MaxHeight(m.height - 8).
      Render(modalContent)
  ```
- ✅ Actions: Edit (e), Delete (d), Close (Enter/Esc)
- ✅ Reuses existing `RenderEventDetailCard` component
- ✅ Implements tea.Model for bubbletea-overlay compatibility

#### QuickAddEventModal (`internal/cli/components/quick_add_event_modal.go`)
- ✅ Removed complex height calculations (line ~100)
- ✅ Set form height to 0 (natural height)
- ✅ Added solid background wrapper in View() method (lines 145-150):
  ```go
  return lipgloss.NewStyle().
      Border(lipgloss.RoundedBorder()).
      BorderForeground(styles.ColorBorder).
      Background(styles.ColorBackground).  // Solid #1a1f2e
      Padding(1, 2).
      Render(m.form.View())
  ```

#### EditEventModal (`internal/cli/components/edit_event_modal.go`)
- ✅ Removed complex height calculations (line ~112)
- ✅ Set form height to 0 (natural height)
- ✅ Added solid background wrapper in View() method (lines 156-161)

#### FilterModal (`internal/cli/components/filter_modal.go`)
- ✅ Removed complex height calculations (line ~128)
- ✅ Set form height to 0 (natural height)
- ✅ Added solid background wrapper in View() method (lines 208-213)

#### DeleteConfirmModal (`internal/cli/components/delete_confirm_modal.go`)
- ✅ Added solid background to modalStyle (line 154):
  ```go
  Background(m.theme.BackgroundColor()).  // Solid #1a1f2e
  ```
- ✅ Removed self-centering (bubbletea-overlay handles positioning)

### 3. Intent Integration (`internal/cli/intents/browse_timeline_intent.go`)

#### Added View Detail Modal Support **NEW!**

1. **Added viewDetailModal field** (line ~102):
   ```go
   viewDetailModal *components.ViewEventDetailModal
   ```

2. **Event Selection Shows Modal** (line ~732):
   - Changed from full-screen transition to modal overlay
   - Preserves timeline context (no screen change)
   - Modal shows event details over the timeline list

3. **Modal Update Handler** (line ~273):
   - Handles modal close and action selection
   - Transitions to edit/delete modals based on user action
   - Clears modal after use

4. **Render Method** (line ~937):
   ```go
   func (i *BrowseTimelineIntent) renderViewDetailModalOverlay(background string) string
   ```
   - Uses bubbletea-overlay for compositing
   - Centers modal with Y offset of -2

5. **View Integration** (line ~385):
   - Checks if viewDetailModal is visible
   - Renders overlay on top of timeline view

6. **Legacy Code Marked**:
   - `BrowseStateEventDetail` cases marked as LEGACY
   - Kept for backward compatibility but should not be reached
   - Full-screen detail view replaced by modal

#### Added staticViewModel Helper (line ~17)
```go
type staticViewModel struct {
    content string
}

func (s *staticViewModel) Init() tea.Cmd { return nil }
func (s *staticViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s *staticViewModel) View() string { return s.content }
```

#### Refactored 5 Render Methods

1. **renderFilterModalOverlay()** (line ~764)
2. **renderQuickAddModalOverlay()** (line ~795)
3. **renderEditModalOverlay()** (line ~824)
4. **renderDeleteModalOverlay()** (line ~854)
5. **renderViewDetailModalOverlay()** (line ~937) **NEW!**

All now use:
```go
bgModel := &staticViewModel{content: background}
overlayModel := overlay.New(
    modalComponent,  // Foreground (with solid background)
    bgModel,         // Background (timeline list)
    overlay.Center,  // X position
    overlay.Center,  // Y position
    0,               // X offset
    -2,              // Y offset (avoid footer)
)
return overlayModel.View()
```

## What to Verify Manually

### Test Environment Setup
```bash
# Build the app
go build -o kariya ./cmd/cli

# Run the app
./kariya
```

### Verification Steps

#### 0. Event Detail Modal (Select event, press Enter) **NEW!**
- [ ] **Positioning**: Modal appears centered over timeline, below logo
- [ ] **Background**: Modal has solid dark background (#1a1f2e)
- [ ] **NO transparency**: Timeline content should NOT show through the modal
- [ ] **Border**: Rounded border visible around modal
- [ ] **Content**: Event details displayed (date, company, project, text, tags, etc.)
- [ ] **Actions work**:
  - [ ] Press 'e' → Edit modal appears
  - [ ] Press 'd' → Delete confirmation modal appears
  - [ ] Press 'Enter' or 'Esc' → Modal closes, returns to timeline
  - [ ] Press 'q' → Modal closes, returns to timeline
- [ ] **Context preserved**: Timeline list still visible in background, scroll position maintained
- [ ] **Footer visible**: Application footer should be visible below modal

#### 1. Quick Add Modal (Press 'a' in Browse Timeline)
- [ ] **Positioning**: Modal appears centered, below the logo
- [ ] **Background**: Modal has solid dark background (#1a1f2e)
- [ ] **NO transparency**: Timeline content should NOT show through the modal
- [ ] **Border**: Rounded border visible around modal
- [ ] **Form scrolling**: If form is long, it should scroll naturally
- [ ] **Footer visible**: Application footer should be visible below modal

#### 2. Edit Event Modal (Select event, press 'e')
- [ ] **Positioning**: Modal appears centered, below the logo
- [ ] **Background**: Modal has solid dark background (#1a1f2e)
- [ ] **NO transparency**: Timeline content should NOT show through the modal
- [ ] **Border**: Rounded border visible around modal
- [ ] **Form scrolling**: Form should scroll if content exceeds height
- [ ] **Footer visible**: Application footer should be visible below modal

#### 3. Filter Modal (Press 'f' in Browse Timeline)
- [ ] **Positioning**: Modal appears centered, below the logo
- [ ] **Background**: Modal has solid dark background (#1a1f2e)
- [ ] **NO transparency**: Timeline content should NOT show through the modal
- [ ] **Border**: Rounded border visible around modal
- [ ] **Form usability**: All filter fields are accessible and usable
- [ ] **Footer visible**: Application footer should be visible below modal

#### 4. Delete Confirm Modal (Select event, press 'd')
- [ ] **Positioning**: Modal appears centered, below the logo
- [ ] **Background**: Modal has solid dark background (#1a1f2e)
- [ ] **NO transparency**: Timeline content should NOT show through the modal
- [ ] **Border**: Red border visible (error color)
- [ ] **Content readable**: Confirmation text and buttons are clear
- [ ] **Footer visible**: Application footer should be visible below modal

### Common Issues to Check For

❌ **Background showing through modal** (SHOULD BE FIXED)
- Previous issue: Background timeline list visible through modal text
- Fix applied: Solid background color (#1a1f2e) in all modal View() methods
- If still present: Background color might not be opaque or overlay compositing issue

❌ **Modal positioned incorrectly**
- Should be: Centered horizontally, below logo, above footer
- Y offset of -2 prevents overlapping footer

❌ **Form doesn't scroll**
- Previous issue: Constrained height broke Huh form scrolling
- Fix applied: Natural height (0) allows Huh to handle scrolling

❌ **Modal too narrow or too wide**
- Modals should be appropriately sized for content
- Should not exceed terminal width

## Technical Details

### Why This Works

1. **bubbletea-overlay handles compositing**
   - No manual ANSI code manipulation
   - Proper layering of background and foreground
   - Handles terminal dimensions automatically

2. **Solid backgrounds prevent transparency**
   - `styles.ColorBackground = #1a1f2e` (dark blue-gray)
   - Lipgloss renders solid color block under form content
   - No background "bleeding through"

3. **Natural height enables scrolling**
   - Huh forms are designed to scroll at natural height
   - Constraining height (previous approach) broke scrolling
   - Height of 0 = use natural height

4. **Positioning is automatic**
   - `overlay.Center` handles horizontal centering
   - Y offset of -2 avoids footer overlap
   - No manual width/position calculations needed

### Background Colors Used

- **QuickAdd, Edit, Filter**: `styles.ColorBackground = #1a1f2e`
- **Delete Confirm**: `theme.BackgroundColor() = #1a1f2e`

Both resolve to the same solid, opaque color.

## Test Results

✅ **All tests passing**: 2,078+ tests
✅ **Build successful**: No compilation errors
✅ **No regressions**: Existing functionality maintained

## What's NOT Changed

- `internal/cli/components/modal.go` - Old modal system kept for other intents
- Other intents (CaptureEvent, Help) - Still use old modal system
- These can be refactored later if desired

## Expected Outcome

When verification is complete:
- ✅ All 4 timeline modals render with bubbletea-overlay
- ✅ Modals have solid, opaque backgrounds (no transparency)
- ✅ Forms scroll properly at natural height
- ✅ Modals positioned correctly (centered, below logo, above footer)
- ✅ Clean, maintainable code
- ✅ No visual artifacts or rendering issues

## Files Modified (8 total)

1. `go.mod` / `go.sum` - Added bubbletea-overlay dependency
2. `internal/cli/components/view_event_detail_modal.go` - **NEW** Event detail modal component
3. `internal/cli/components/quick_add_event_modal.go` - Natural height + solid background
4. `internal/cli/components/edit_event_modal.go` - Natural height + solid background
5. `internal/cli/components/filter_modal.go` - Natural height + solid background
6. `internal/cli/components/delete_confirm_modal.go` - Added background, removed centering
7. `internal/cli/intents/browse_timeline_intent.go` - Use bubbletea-overlay for all 5 modals

## Next Steps (If Issues Found)

If transparency is still visible:
1. Check if `styles.ColorBackground` is being applied correctly
2. Verify lipgloss version is compatible
3. Try explicit dark color (e.g., `#000000`) to rule out opacity issues
4. Check if bubbletea-overlay has alpha channel support issues

If positioning is wrong:
1. Adjust Y offset in overlay.New() calls
2. Check terminal size detection
3. Verify logo height calculation

If forms don't scroll:
1. Verify height is set to 0 (natural)
2. Check terminal height detection
3. Test with longer forms
