# List Model Refactoring Guide

## Overview

This guide documents the common patterns extracted from the list models (ListModel, BurstListModel, FactListModel) and provides instructions for refactoring these models to use the extracted patterns.

## Extracted Patterns

### 1. ListDeletionState (in `internal/cli/models/list_patterns.go`)

The `ListDeletionState` encapsulates all deletion workflow logic that was duplicated across all three list models.

#### Usage Pattern

**Before (Old Code):**
```go
type ListModel struct {
    // ... other fields ...
    deletionConfirm   *ConfirmationDialog
    deletingEventID   string
    deleteSuccessMsg  string
    deleteErrorMsg    string
    showDeleteMessage bool
}

// In Update():
if m.deletionConfirm != nil && m.deletionConfirm.IsConfirmed() {
    // perform deletion
}

// In View():
if m.deletionConfirm != nil {
    return m.deletionConfirm.View()
}
```

**After (New Code):**
```go
type ListModel struct {
    // ... other fields ...
    deletionState     *ListDeletionState
    showDeleteMessage bool
}

// In NewListModel():
m.deletionState = NewListDeletionState()

// In Update():
if m.deletionState.IsConfirming() {
    cmd := m.deletionState.UpdateConfirmation(msg)
    if m.deletionState.IsConfirmed() {
        return m.performEventDeletion()
    }
    if m.deletionState.IsCancelled() {
        m.deletionState.Clear()
        return m, nil
    }
    return m, cmd
}

// In performEventDeletion():
m.deletionState.SetSuccessMsg("Event deleted successfully")
// or
m.deletionState.SetErrorMsg(fmt.Sprintf("Error: %v", err))
```

#### ListDeletionState API

```go
// Create confirmation dialog
ShowConfirmation(title, itemName string)

// Check state
IsConfirming() bool
IsConfirmed() bool
IsCancelled() bool

// Set messages
SetSuccessMsg(msg string)
SetErrorMsg(msg string)

// Handle messages
UpdateConfirmation(msg tea.Msg) tea.Cmd

// Reset state
Clear()

// Public fields for rendering
ConfirmationDialog *ConfirmationDialog
DeletingItemID string
SuccessMsg string
ErrorMsg string
ShowMessage bool
```

### 2. ListNavigationKeyHandler (in `internal/cli/models/list_patterns.go`)

The `ListNavigationKeyHandler` consolidates all keyboard navigation logic that was duplicated across all three list models.

#### Key Bindings Handled

- **Navigation**: `up`, `k`, `down`, `j`
- **Page Navigation**: `pgup`, `ctrl+b`, `pgdn`, `ctrl+f`
- **Jump to End**: `home`, `g`, `end`, `G`
- **Item Actions**: `enter`, `v` (view), `e` (edit), `x`, `d` (delete), `space` (toggle)

#### Usage Pattern

**Before (Old Code):**
```go
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        keyStr := msg.String()
        switch keyStr {
        case "up", "k":
            m.cursorIdx--
            if m.cursorIdx < 0 {
                m.cursorIdx = 0
            }
            m.updateTableRows()
        case "down", "j":
            m.cursorIdx++
            if m.cursorIdx >= len(m.filtered) {
                m.cursorIdx = len(m.filtered) - 1
            }
            m.updateTableRows()
        case "enter":
            if m.cursorIdx >= 0 && m.cursorIdx < len(m.filtered) {
                event := m.filtered[m.cursorIdx]
                return m, func() tea.Msg {
                    return EventActionSelectedMsg{Event: event, Action: EventActionView}
                }
            }
        // ... many more cases ...
        }
    }
    // ...
}
```

**After (New Code):**
```go
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        keyStr := msg.String()

        // Handle navigation keys through key handler
        cmd := m.keyHandler.HandleNavigationKey(keyStr)
        if cmd != nil {
            return m, cmd
        }
    }
    // ...
}

// Implement ListItemCallbacks interface for key handler
func (m *ListModel) OnView(item interface{}) tea.Cmd {
    if event, ok := item.(*career.CareerEvent); ok {
        return func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionView} }
    }
    return nil
}

func (m *ListModel) OnEdit(item interface{}) tea.Cmd {
    if event, ok := item.(*career.CareerEvent); ok {
        return func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionEdit} }
    }
    return nil
}

func (m *ListModel) OnDelete(item interface{}) {
    if event, ok := item.(*career.CareerEvent); ok {
        m.deletionState.DeletingItemID = event.ID
        m.deletionState.ShowConfirmation("event", event.Text)
    }
}

func (m *ListModel) OnToggleSelection(item interface{}) {
    if event, ok := item.(*career.CareerEvent); ok {
        m.selectedEvents[event.ID] = !m.selectedEvents[event.ID]
    }
}

func (m *ListModel) HasSelectedItem() interface{} {
    return m.GetSelectedEvent()
}

func (m *ListModel) MoveUp(count int) {
    m.listContainer.MoveUp(count)
}

func (m *ListModel) MoveDown(count int) {
    m.listContainer.MoveDown(count)
}

func (m *ListModel) MoveToFirst() {
    m.listContainer.MoveToFirst()
}

func (m *ListModel) MoveToLast() {
    m.listContainer.MoveToLast()
}

func (m *ListModel) UpdateDisplay() {
    m.updateTableRows()
}

func (m *ListModel) GetRowCount() int {
    return len(m.table.Rows())
}

func (m *ListModel) GetCurrentIndex() int {
    return m.listContainer.GetSelectedIdx()
}

func (m *ListModel) GetPageSize() int {
    return m.pageSize
}
```

#### ListItemCallbacks Interface

```go
type ListItemCallbacks interface {
    // Item action callbacks
    OnView(item interface{}) tea.Cmd
    OnEdit(item interface{}) tea.Cmd
    OnDelete(item interface{})
    OnToggleSelection(item interface{})
    HasSelectedItem() interface{}

    // Navigation callbacks
    MoveUp(count int)
    MoveDown(count int)
    MoveToFirst()
    MoveToLast()
    UpdateDisplay()

    // Query callbacks
    GetRowCount() int
    GetCurrentIndex() int
    GetPageSize() int
}
```

## Refactoring Checklist

For each list model (ListModel, BurstListModel, FactListModel):

### Step 1: Update Struct Definition

- [ ] Replace individual deletion fields with single `deletionState *ListDeletionState`
- [ ] Add `keyHandler *ListNavigationKeyHandler`
- [ ] Remove duplicate key handling code
- [ ] Keep pagination fields (`currentPage`, `pageSize`, `totalCount`)

### Step 2: Update Initialization

- [ ] Initialize `deletionState` in constructor:
  ```go
  deletionState: NewListDeletionState()
  ```
- [ ] Initialize `keyHandler` after constructor setup:
  ```go
  m.keyHandler = NewListNavigationKeyHandler(m)
  ```

### Step 3: Implement ListItemCallbacks Interface

- [ ] Implement `OnView()` returning appropriate ActionSelectedMsg
- [ ] Implement `OnEdit()` returning appropriate ActionSelectedMsg
- [ ] Implement `OnDelete()` calling `deletionState.ShowConfirmation()`
- [ ] Implement `OnToggleSelection()` updating the selection map
- [ ] Implement `HasSelectedItem()` returning current item or nil
- [ ] Implement movement methods (`MoveUp()`, `MoveDown()`, etc.)
- [ ] Implement `UpdateDisplay()` calling `updateTableRows()`
- [ ] Implement query methods (`GetRowCount()`, `GetCurrentIndex()`, `GetPageSize()`)

### Step 4: Simplify Update() Method

- [ ] Replace deletion confirmation logic with:
  ```go
  if m.deletionState.IsConfirming() {
      cmd := m.deletionState.UpdateConfirmation(msg)
      if m.deletionState.IsConfirmed() {
          return m.performXXXDeletion()
      }
      if m.deletionState.IsCancelled() {
          m.deletionState.Clear()
          return m, nil
      }
      return m, cmd
  }
  ```
- [ ] Replace all key handling with:
  ```go
  cmd := m.keyHandler.HandleNavigationKey(keyStr)
  if cmd != nil {
      return m, cmd
  }
  ```
- [ ] Remove duplicate navigation switch cases

### Step 5: Simplify Deletion Logic

- [ ] In `performXXXDeletion()`, update to use:
  ```go
  m.deletionState.SetSuccessMsg("XXX deleted successfully")
  // or
  m.deletionState.SetErrorMsg(fmt.Sprintf("Error: %v", err))
  ```

### Step 6: Update View() Method

- [ ] Check deletion state:
  ```go
  if m.deletionState.IsConfirming() {
      return m.deletionState.ConfirmationDialog.View()
  }
  if m.showDeleteMessage && m.deletionState.ShowMessage {
      return m.renderDeleteMessage()
  }
  ```

## Code Reduction Benefits

### Before (Typical Deletion Flow - ~80 lines per model)
```
- 4-5 deletion-related fields
- 15-20 lines of deletion confirmation in Update()
- 20-30 lines of performXXXDeletion() implementation
- 10-15 lines in View() for deletion rendering
```

### After (With Extracted Patterns - ~20 lines per model)
```
- 1 deletion state field
- 8-10 lines of unified deletion handling in Update()
- 8-10 lines in performXXXDeletion() (just service call and message)
- 4-5 lines in View() for deletion rendering
```

## Total Refactoring Impact

**Estimated Code Reduction:**
- ListModel: ~120 lines → ~80 lines (33% reduction)
- BurstListModel: ~120 lines → ~80 lines (33% reduction)
- FactListModel: ~140 lines → ~90 lines (36% reduction)
- **Total**: ~380 lines of duplicated code consolidated into ~150 lines of shared patterns

## Testing Strategy

After refactoring each model:

1. Run specific model tests:
   ```bash
   go test ./internal/cli/models -run ListModel -v
   go test ./internal/cli/models -run BurstListModel -v
   go test ./internal/cli/models -run FactListModel -v
   ```

2. Run full test suite:
   ```bash
   make test
   ```

3. Verify key bindings still work:
   - Navigation (up/down/pgup/pgdn/home/end)
   - Item actions (enter/e/d)
   - Deletion workflow
   - Selection (space)

## Future Improvements

After completing this refactoring, consider:

1. **Extract Pagination Pattern** - Create `ListPaginationHelper` for pagination logic
2. **Extract Filter/Sort Pattern** - Create `ListFilterSortHelper` for filtering and sorting
3. **Create GenericListModel Base** - Consolidate all common list logic into a reusable base
4. **Unified Test Suite** - Create generic list model tests that apply to all three models

## See Also

- `internal/cli/models/list_patterns.go` - Extracted pattern implementations
- `internal/cli/models/list.go` - Example refactoring target
- `internal/cli/models/burst_list.go` - Another refactoring target
- `internal/cli/models/fact_list.go` - Third refactoring target

