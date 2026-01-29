# BUG-016: Skill detail modal shows wrong skill when viewing from list

## Summary

When viewing a skill from the skills list, the detail modal displays the wrong
skill. Subsequent edit and delete actions from within the detail modal also
operate on the wrong skill, risking data corruption.

## Steps to Reproduce

1. Navigate to Manage Skills
2. Wait for the skills list to load (with 2+ skills)
3. Use arrow keys to navigate to any skill other than the first one
4. Press Enter to view the skill detail
5. Observe the detail modal shows a different skill (the first one)
6. From the detail modal, press `e` to edit or `d` to delete
7. Observe the edit/delete operates on the wrong skill

## Expected Behavior

The detail modal should display the skill the user selected. Edit and delete
actions from the modal should operate on that same skill.

## Actual Behavior

The detail modal always shows the skill at index 0 of `i.skills` (or whichever
stale index `i.selectedIndex` holds), regardless of which skill the user
selected in the list screen. Edit and delete from the modal cascade the error.

## Root Cause

**File**: `internal/cli/intents/skills_management/helpers.go:269-275`

```go
func (i *Intent) openViewDetailModal() tea.Cmd {
    if len(i.skills) == 0 || i.selectedIndex >= len(i.skills) {
        return nil
    }
    skill := i.skills[i.selectedIndex]   // BUG: stale index
    i.selectedSkill = skill              // Overrides the correct value
    ...
}
```

There are two problems:

1. **Stale `i.selectedIndex`**: The intent maintains its own `selectedIndex`
   field which is only synced once during `handleSkillsLoaded()` via
   `syncTableSelection()`. After that, user navigation updates the screen's
   `tableBehavior.selectedIndex`, but the intent's `selectedIndex` stays at 0.

2. **Two separate TableBehavior instances**: The intent creates its own
   `tableBehavior` in `NewIntent()`, and the `SkillsListScreen` creates a
   separate one in `NewSkillsListScreen()`. The user navigates the screen's
   table, but `syncTableSelection()` reads from the intent's table which is
   never navigated.

The call chain that triggers the bug:

1. User presses Enter on the list screen
2. `SkillsListScreen.handleViewAction()` correctly retrieves the selected
   `*career.Skill` via `tableBehavior.GetSelectedItem()` and passes it in
   a `NavigateResult`
3. `handleNavigateData("view")` correctly sets `i.selectedSkill = skill`
4. Then calls `openViewDetailModal()` which **overrides** `i.selectedSkill`
   with `i.skills[i.selectedIndex]` (the wrong skill)

The edit and delete paths from the list screen (pressing `e` or `d` directly)
are NOT affected because they pass the skill directly to `openAddEditModal()`
and `openDeleteModal()` without calling `openViewDetailModal()`.

However, edit and delete from within the detail modal ARE affected because
`handleViewDetailModalUpdate()` uses the now-incorrect `i.selectedSkill`.

## Fix Plan

Modify `openViewDetailModal()` to use `i.selectedSkill` (already correctly set
by `handleNavigateData()`) instead of re-indexing via `i.skills[i.selectedIndex]`:

```go
func (i *Intent) openViewDetailModal() tea.Cmd {
    if i.selectedSkill == nil {
        return nil
    }

    skill := i.selectedSkill

    eventCount := 0
    if i.eventCounts != nil {
        eventCount = i.eventCounts[skill.ID]
    }

    width, height := i.getTerminalDimensions()

    i.viewDetailModal = modals.NewDetailModal(skill, i.Theme(), eventCount, nil)
    i.viewDetailModal.SetDimensions(width, height)
    i.viewDetailModal.Show()

    return nil
}
```

## Regression Tests

- Navigate to non-first skill and view: detail modal shows the correct skill
- View a skill then edit from modal: edit modal receives the correct skill
- View a skill then delete from modal: delete confirmation uses the correct skill ID
- View with no selected skill (nil): returns nil without panic

## Related

- `internal/cli/intents/skills_management/handlers.go` - `handleNavigateData()` correctly passes skill
- `internal/cli/screens/skills/skill_list.go` - Screen correctly uses `GetSelectedItem()`
- `internal/cli/intents/skills_management/types.go` - `selectedIndex` field (stale, unused correctly)

## Severity

- [x] High - Wrong skill can be edited or deleted, risking data corruption
