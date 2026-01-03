# app.go Migration Quick Start Guide

**Version**: 1.0
**Date**: 2026-01-03
**Purpose**: Quick reference for implementing the app.go migration

---

## TL;DR - The 3 Critical Changes

### Change 1: Update() Method (5 minutes)

**Location**: `internal/cli/app/app.go` line 368

**Add at the very beginning of Update():**

```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // NEW: Check intent mode FIRST
    if m.inIntentMode {
        return m.handleIntentMessage(msg)
    }

    // EXISTING: Rest of the code continues as-is
    switch msg.(type) {
    case models.BackMsg:
        // ... existing logic
```

### Change 2: View() Method (5 minutes)

**Location**: `internal/cli/app/app.go` line 1273

**Add at the very beginning of View():**

```go
func (m *Model) View() string {
    // NEW: Check intent mode FIRST
    if m.inIntentMode && m.intentRouter != nil {
        return m.intentRouter.View()
    }

    // EXISTING: Rest of the code continues as-is
    switch m.currentScreen {
    case MainMenuScreen:
        // ... existing logic
```

### Change 3: Global Shortcuts (10 minutes)

**Location**: `internal/cli/app/app.go` line 400+ (BackMsg handler)

**Replace the BackMsg case:**

```go
case models.BackMsg:
    // NEW: Check if in intent mode
    if m.inIntentMode {
        cmd, err := m.intentRouter.Back()
        if err != nil {
            m.deactivateIntent()
        }
        return m, cmd
    }

    // EXISTING: Screen-based back navigation
    if m.currentScreen == ListScreen {
        m.previousScreen = m.currentScreen
        m.currentScreen = HomeScreen
        return m, nil
    }
    // ... rest of existing back handling
```

---

## Phase 1 Checklist

### Step 1: Implement Changes (20 minutes)
- [ ] Add intent mode check to Update()
- [ ] Add intent mode check to View()
- [ ] Update BackMsg handler
- [ ] Test locally

### Step 2: Update Menu Items (30 minutes)
- [ ] Find `handleMenuItemSelection()` method (line 1611)
- [ ] Replace screen-based activation with intent activation:

```go
case "c":
    return m, m.activateIntent("capture_event", make(map[string]interface{}))
case "l":
    return m, m.activateIntent("browse_timeline", make(map[string]interface{}))
case "g":
    return m, m.activateIntent("generate_cv", make(map[string]interface{}))
```

### Step 3: Testing (1 hour)
- [ ] Run `go test ./...`
- [ ] Manually test menu activation
- [ ] Manually test back navigation
- [ ] Manually test switching between intents

### Step 4: Verification
- [ ] All tests pass ✅
- [ ] No regressions ✅
- [ ] Intents activate from menu ✅
- [ ] Back navigation works ✅

---

## Code Locations Reference

| Component | File | Line | Purpose |
|-----------|------|------|---------|
| **Update()** | app.go | 368 | Main message handler |
| **View()** | app.go | 1273 | Main rendering |
| **BackMsg Handler** | app.go | 400+ | Back navigation |
| **Menu Selection** | app.go | 1611 | Menu item handling |
| **activateIntent()** | app.go | 1708 | Intent activation |
| **deactivateIntent()** | app.go | 1724 | Intent deactivation |
| **handleIntentMessage()** | app.go | 1733 | Intent message routing |

---

## Testing Commands

```bash
# Run all tests
go test -v ./...

# Run with race detector
go test -race ./...

# Run specific test file
go test -v ./internal/cli/app/...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

---

## Common Issues & Solutions

### Issue: "intentRouter is nil"
**Solution**: Check that NewModel() properly initializes the router (line 146-237)

### Issue: "inIntentMode is always false"
**Solution**: Check that activateIntent() is being called from menu

### Issue: "Intent doesn't receive messages"
**Solution**: Verify Update() check is at the very beginning before other logic

### Issue: "Intent view doesn't render"
**Solution**: Verify View() check is at the very beginning before other logic

### Issue: "Back from intent doesn't work"
**Solution**: Verify BackMsg handler checks inIntentMode before screen-based logic

---

## File Modification Summary

### app.go Changes

| Method | Current Lines | New Lines | Change |
|--------|---------------|-----------|--------|
| Update() | 906 | 920 | +14 lines (intent check) |
| View() | 180 | 190 | +10 lines (intent check) |
| handleMenuItemSelection() | 50+ | 60+ | +10 lines (intent activation) |
| **Total** | 1766 | 1786 | +20 lines |

**No lines removed** - All existing code remains, just with intent mode checks added

---

## Success Indicators

✅ **Phase 1 Complete When**:
1. Update() has intent mode check at top
2. View() has intent mode check at top
3. BackMsg handler checks inIntentMode
4. Menu items activate intents
5. All tests pass
6. No regressions detected

✅ **Can Proceed to Phase 2 When**:
- All 5 intents activate from menu
- Back navigation works in intents
- Result handlers work correctly
- No breaking changes to existing workflows

---

## Next Steps After Phase 1

1. **Phase 2**: Verify all 5 core intents work end-to-end
2. **Phase 3**: Create new intents for remaining screens
3. **Phase 4**: Remove legacy screen code
4. **Phase 5**: Optimize and cleanup

See `APP_GO_MIGRATION_PLAN.md` for complete details.

---

## Quick Links

- **Full Migration Plan**: `docs/APP_GO_MIGRATION_PLAN.md`
- **Integration Review**: `docs/APP_GO_INTENT_INTEGRATION_REVIEW.md`
- **Intent Architecture**: `docs/TUI_INTENT_DIAGRAM.md`
- **Intent Reference**: `internal/cli/intents/contract.go`
- **Router Implementation**: `internal/cli/intents/router.go`

---

**Time to Complete Phase 1**: ~2 hours (including testing)
**Complexity**: LOW (just 3 conditional checks)
**Risk**: MINIMAL (no existing code modified, only additions)


