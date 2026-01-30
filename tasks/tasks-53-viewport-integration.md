# TASK-53: Viewport integration for content stretching and scrolling

## Summary

Integrate bubbletea's viewport component to make content areas stretch to fill available vertical space and scroll when content exceeds the allocated height. This applies to both screen content and modal content.

## Context

Currently, content flows naturally without height constraints. With pinned layout (TASK-52), we have a fixed header (logo + breadcrumbs) and footer (help text), leaving a middle area that should stretch to fill available space and scroll if content exceeds that space.

## Acceptance Criteria

- [ ] Content areas calculate available height: `terminalHeight - headerHeight - footerHeight`
- [ ] TableBehavior integrates with viewport for scrolling
- [ ] Screens with lists/tables use viewport for content
- [ ] Modals use viewport for long content
- [ ] Viewport scrolling works with keyboard (↑/↓, PgUp/PgDn)
- [ ] Scroll position persists across updates
- [ ] Viewport height updates on terminal resize

## Technical Notes

### Files to Modify

- `internal/cli/behaviors/table.go` - Add viewport integration
- `internal/cli/uikit/layout/screen_layout.go` - Calculate content height, pass to content
- `internal/cli/uikit/feedback/modal.go` - Add viewport for modal content
- `internal/cli/screens/*/` - Update screens to use viewport

### New Files

- `internal/cli/uikit/viewport/` - Viewport wrapper/helper package (optional)

### Viewport Integration Pattern

```go
import "github.com/charmbracelet/bubbles/viewport"

type MyScreen struct {
    *base.BaseScreen
    viewport viewport.Model
    content  string  // Full content (may exceed viewport)
}

func NewMyScreen(termInfo *terminal.Info) *MyScreen {
    // Calculate available height
    contentHeight := termInfo.Height - headerHeight - footerHeight
    
    vp := viewport.New(termInfo.Width, contentHeight)
    vp.SetContent(content)
    
    return &MyScreen{
        viewport: vp,
    }
}

func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    var cmd tea.Cmd
    s.viewport, cmd = s.viewport.Update(msg)
    return cmd, nil
}

func (s *MyScreen) View() string {
    return s.viewport.View()
}
```

### Dependencies

- `github.com/charmbracelet/bubbles/viewport` - Already in go.mod

## Testing Requirements

### Unit Tests

- [ ] Viewport initializes with correct height
- [ ] Viewport scrolls with keyboard input
- [ ] Viewport content updates correctly
- [ ] Viewport height updates on terminal resize
- [ ] Scroll position maintained across updates

### Integration Tests

- [ ] TableBehavior with viewport scrolls correctly
- [ ] Modal with long content scrolls
- [ ] Screens with lists use viewport

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Tests written FIRST (TDD)
- [ ] Tests pass with >= 95% coverage
- [ ] No pattern violations (`make check-patterns`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Committed with `make ai-commit`

## Notes

- This builds on TASK-52 (pinned layout)
- Viewport should be optional - screens without viewport still work
- TableBehavior should handle viewport integration internally
- Modal viewport should auto-enable for content exceeding max height
