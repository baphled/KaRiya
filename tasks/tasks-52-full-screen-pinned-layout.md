# TASK-001: Full-screen pinned layout - logo at top, footer at bottom

## Summary

Modify the TUI layout so the screen maximises the terminal: logo pinned to line 0 at the top, footer pinned to the bottom, and content flows immediately below the header. Currently everything is vertically centered which wastes screen real estate.

## Acceptance Criteria

- [ ] Logo renders at line 0 (no blank lines before it)
- [ ] Footer/help text renders on the last lines of the terminal
- [ ] Content flows immediately below the header (logo + breadcrumbs)
- [ ] Dynamic spacer fills the gap between content and footer
- [ ] Graceful overflow when content exceeds available height (spacer becomes 0)
- [ ] Horizontal centering preserved
- [ ] Modal overlays still render correctly on top of pinned layout
- [ ] Applies globally to all screens (intents and main menu)

## Technical Notes

### Files to Modify

- `internal/cli/uikit/layout/screen_layout.go` - Restructure `Render()` into 3 sections (header/content/footer) with spacer, change vertical alignment from `Center` to `Top`
- `internal/cli/app/menu.go` - Update `viewMenu()` to use same spacer approach with `Top` alignment
- `internal/cli/uikit/primitives/text.go` - (Optional) Add `PlaceInTerminal()` helper

### New Files

- `internal/cli/uikit/layout/suite_test.go` - Ginkgo test suite for layout package
- `internal/cli/uikit/layout/screen_layout_test.go` - Tests for pinned layout behaviour

### Key Change

Current (`screen_layout.go:253`):
```go
rendered := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, combined)
```

Target:
```go
// Build header, content, footer separately
// Calculate spacer: termHeight - headerH - contentH - footerH
// Join: header + content + spacer + footer
rendered := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Top, combined)
```

### Dependencies

- No external dependencies
- `LogoSpacing` field semantics change: spacing now means "after logo" only (no pre-logo blank lines)

### Patterns to Use

| Need | Use |
|------|-----|
| Layout | `layout.ScreenLayout` (UIKit) |
| Height calc | `lipgloss.Height()` |
| Placement | `lipgloss.Place()` with `lipgloss.Top` |

## Testing Requirements

### Unit Tests (screen_layout_test.go)

- [ ] Logo starts at line 0 of rendered output
- [ ] Footer appears on last lines of rendered output
- [ ] Content appears immediately after header section
- [ ] Spacer fills remaining vertical space (total height == terminal height)
- [ ] Overflow: spacer is 0 when content exceeds terminal height
- [ ] No logo mode: breadcrumbs start at line 0
- [ ] No footer mode: no spacer, content flows naturally
- [ ] Modal overlay renders correctly on pinned layout

### Existing Tests to Verify

- [ ] `view_helpers_test.go` - CreateStandardView with logo spacing
- [ ] `contract_test.go` - GetLogoSpacing/SetLogoSpacing defaults

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Tests written FIRST (TDD)
- [ ] Tests pass with >= 95% coverage
- [ ] No pattern violations (`make check-patterns`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Committed with `make ai-commit`
