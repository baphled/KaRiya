# Tasks 12-16: Standardized Centered Views with Logo

## Overview

This task series implements a comprehensive standardization of all TUI views in KaRiya, ensuring every screen has:
- **Logo at the top** with configurable spacing (default: 2 lines from top edge)
- **Centered content** using SmartContainer throughout
- **Context-aware help** at the bottom with visual separator
- **Modal overlays** for errors, loading, progress, and success messages
- **Full terminal width** utilization for content
- **Consistent user experience** across all intents

## Task Breakdown

### Task 12: Core Components (4-5 hours)
**File**: `tasks-12-standardized-view-core-components.md`

Creates the foundational components for the standardized view system:

**Deliverables**:
- `StandardView` component with logo, content, and help layout
- Modal system with 5 types (Error, Loading, Progress, Success, Warning)
- Modal fade-in animation (150ms)
- Loading message rotator with automatic rotation (2s interval)
- Spinner animation for indeterminate operations
- SmartContainer dimming support for modal overlays

**Key Features**:
- Logo with configurable spacing from top edge
- Visual separator between content and help footer
- Adaptive modal sizing based on content
- Terminal bell for error modals
- Auto-dismiss for success modals (3s)
- Escape key dismissal for error modals
- Cancellable loading operations

### Task 13: Intent Infrastructure (2-3 hours)
**File**: `tasks-13-intent-infrastructure-updates.md`

Updates the base infrastructure to support standardized views across all intents:

**Deliverables**:
- Enhanced `BaseIntent` with terminal awareness
- Logo instance management in BaseIntent
- View helper functions for easy StandardView creation
- Modal helper functions (ShowErrorModal, ShowLoadingModal, etc.)
- Footer helper functions (NavigationFooter, FormFooter, etc.)
- State management helpers (SetLoading, SetError, SetSuccess)
- Terminal info propagation from app → router → intent

**Key Features**:
- TerminalAwareIntent implementation in BaseIntent
- Lazy logo initialization
- Standardized state management across intents
- Predefined footer templates for common patterns
- Helper methods that reduce boilerplate

### Task 14: Migrate Primary Intents (5 hours)
**File**: `tasks-14-migrate-capture-browse-intents.md`

Migrates the two most complex and frequently used intents:

**Deliverables**:
- CaptureEventIntent using StandardView
- BrowseTimelineIntent using StandardView
- Loading states with message rotation
- Success modals with auto-dismiss
- Error modals with bell alerts
- Context-aware help for all states

**Pattern Established**:
```go
func (i *SomeIntent) View() string {
    view := i.CreateViewWithBreadcrumbs("Main Menu", "Intent", i.getStateName())
    
    // Modal states
    if i.isLoading {
        ShowLoadingModal(view, i.loadingMessage, true)
    }
    if i.errorState != nil {
        ShowErrorModal(view, i.errorState)
    }
    if i.ShouldShowSuccess() {
        ShowSuccessModal(view, i.successMessage)
    }
    
    view.WithContent(i.getStateContent())
    view.WithHelp(i.getContextHelp()).WithFooterSeparator(true)
    
    return view.Render()
}
```

### Task 15: Migrate Remaining Intents (4.5 hours)
**File**: `tasks-15-migrate-remaining-intents.md`

Completes the migration by updating the remaining primary intents:

**Deliverables**:
- GenerateCVIntent with progress tracking
- ExportArtifactIntent with export progress
- ConfigureSystemIntent with configuration states
- All using StandardView pattern
- Progress modals for multi-step operations
- Loading message rotation for all long operations

**Special Features**:
- **GenerateCV**: Multi-step progress (0% → 25% → 50% → 75% → 100%)
  - "🔍 Analyzing career events..."
  - "📊 Calculating impact metrics..."
  - "✨ Generating professional bullets..."
  - "📝 Formatting final document..."
  - "✅ CV ready!"
- **ExportArtifact**: Export progress with file path in success modal
- **ConfigureSystem**: Configuration validation and apply progress

### Task 16: Testing and Polish (3-4 hours)
**File**: `tasks-16-standardview-testing-polish.md`

Ensures quality, consistency, and comprehensive documentation:

**Deliverables**:
- Comprehensive component tests (StandardView, Modal, LoadingMessages)
- Terminal size testing (80x24 to 240x40)
- Cross-intent consistency tests
- Performance benchmarks (< 100ms render target)
- Edge case handling (small terminals, nil info, overflow)
- Developer documentation (`STANDARDVIEW_GUIDE.md`, `MODAL_PATTERNS.md`)
- Visual test program for manual verification
- Updated existing documentation

**Quality Targets**:
- Test coverage > 87%
- All intents using StandardView (100% consistency)
- Performance < 100ms for full view render
- 0 regressions in existing functionality

## Implementation Timeline

### Week 1
- **Days 1-2**: Task 12 (Core Components)
- **Day 3**: Task 13 (Infrastructure)
- **Days 4-5**: Task 14 (CaptureEvent, BrowseTimeline)

### Week 2
- **Days 1-2**: Task 15 (GenerateCV, ExportArtifact, ConfigureSystem)
- **Days 3-4**: Task 16 (Testing, Documentation)
- **Day 5**: Final polish and verification

### Week 3
- Integration testing and bug fixes
- User acceptance testing
- Final documentation updates

## Key Design Decisions

### 1. Logo Display
- **Decision**: Show logo on every screen (all intents)
- **Spacing**: 2 lines from top edge (configurable)
- **Mode**: Static (no animation in intents)
- **Rationale**: Consistent branding, professional appearance

### 2. Content Width
- **Decision**: Use full terminal width for content
- **Rationale**: Maximizes available space for tables, forms, CV previews

### 3. Modal Behavior
- **Error Modals**: 
  - Escape key dismisses
  - Terminal bell for accessibility
  - No auto-dismiss (requires user acknowledgment)
- **Loading Modals**:
  - Cancellable by user (shows confirmation)
  - Message rotation every 2s for long operations
  - Spinner animation at 80ms intervals
- **Progress Modals**:
  - Show percentage and progress bar
  - Update as operation progresses
  - Useful for multi-step operations (CV generation, exports)
- **Success Modals**:
  - Auto-dismiss after 3 seconds
  - Include useful info (file paths, counts, etc.)
  - Provide immediate positive feedback

### 4. Help Footer
- **Decision**: Context-aware help text per state
- **Separator**: Visual horizontal line above help
- **Format**: "key Action  key Action  |  q Quit  m Main Menu"
- **Rationale**: Always visible, always relevant to current context

### 5. Transitions
- **Decision**: Simple fade-in for modals (150ms)
- **Future**: No complex transitions yet (deferred to future enhancement)
- **Rationale**: Subtle polish without complexity

## Success Metrics

### Completion Criteria
- [ ] All 5 tasks completed (100%)
- [ ] All primary intents using StandardView (100%)
- [ ] Test coverage > 87%
- [ ] Performance benchmarks met (< 100ms renders)
- [ ] 0 test failures
- [ ] 0 regressions
- [ ] Documentation complete

### User Experience Metrics
- [ ] Logo visible on every screen
- [ ] Consistent layout across all intents
- [ ] Smooth modal animations
- [ ] Responsive to terminal resize
- [ ] Clear context at all times (breadcrumbs, help)
- [ ] Professional, polished appearance

### Developer Experience Metrics
- [ ] Easy to create new standardized views
- [ ] Clear patterns to follow
- [ ] Comprehensive documentation
- [ ] Helper functions reduce boilerplate
- [ ] Tests provide confidence

## Files Created

### Components (Task 12)
- `internal/cli/components/standard_view.go`
- `internal/cli/components/standard_view_test.go`
- `internal/cli/components/modal.go`
- `internal/cli/components/modal_test.go`
- `internal/cli/components/loading_messages.go`
- `internal/cli/components/loading_messages_test.go`

### Infrastructure (Task 13)
- `internal/cli/intents/view_helpers.go`
- `internal/cli/intents/view_helpers_test.go`

### Testing (Task 16)
- `internal/cli/components/terminal_size_test.go`
- `internal/cli/components/performance_test.go`
- `internal/cli/intents/consistency_test.go`
- `cmd/test_all_views/main.go`

### Documentation (Task 16)
- `docs/STANDARDVIEW_GUIDE.md`
- `docs/MODAL_PATTERNS.md`

## Files Modified

### Infrastructure
- `internal/cli/intents/contract.go` (BaseIntent enhancements)
- `internal/cli/intents/router.go` (terminal info propagation)
- `internal/cli/app/app.go` (terminal info flow)
- `internal/cli/components/smart_container.go` (dimming support)

### Intents
- `internal/cli/intents/capture_event.go`
- `internal/cli/intents/capture_event_intent.go`
- `internal/cli/intents/capture_event_test.go`
- `internal/cli/intents/browse_timeline.go`
- `internal/cli/intents/browse_timeline_intent.go`
- `internal/cli/intents/browse_timeline_test.go`
- `internal/cli/intents/generate_cv.go`
- `internal/cli/intents/generate_cv_intent.go`
- `internal/cli/intents/generate_cv_test.go`
- `internal/cli/intents/export_artifact.go`
- `internal/cli/intents/export_artifact_intent.go`
- `internal/cli/intents/export_artifact_test.go`
- `internal/cli/intents/configure_system.go`
- `internal/cli/intents/configure_system_intent.go`
- `internal/cli/intents/configure_system_test.go`

### Documentation
- `docs/TUI_DEVELOPER_GUIDE.md` (StandardView section)
- `docs/TUI_STANDARDS.md` (updated standards)
- `AGENTS.md` (architecture updates)

## Dependencies

### External
- BubbleTea (tea.Model, tea.Cmd, tea.Msg)
- Lipgloss (styling, layout)
- Existing SmartContainer component
- Existing ASCIILogo component
- Existing terminal.Info infrastructure

### Internal
- Task 12 → Task 13 (infrastructure needs components)
- Task 13 → Task 14 (intents need infrastructure)
- Task 14 → Task 15 (patterns established for remaining intents)
- Tasks 12-15 → Task 16 (testing and docs need complete system)

## Rollback Plan

### Per-Task Rollback
Each task file includes a rollback plan with specific steps.

### Full Rollback
If the entire feature needs to be rolled back:

1. Identify the commit before Task 12 started
2. Run: `git log --oneline` to find commit hash
3. Create a rollback branch: `git checkout -b rollback-standardview`
4. Reset to before Task 12: `git reset --hard <commit-hash>`
5. Verify all tests pass: `go test ./...`
6. Create new branch for fixes: `git checkout -b fix-standardview-issues`
7. Address issues and re-implement with learnings

### Partial Rollback
If only specific intents have issues:

1. Identify commits for that intent
2. Revert specific commits: `git revert <commit-hash>`
3. Keep other intents' changes
4. Fix and re-implement just the problematic intent

## Notes

- **Atomic Commits**: Each phase within each task should be committed atomically
- **Test-Driven**: Write/update tests before or alongside implementation
- **Documentation**: Update docs as you go, not at the end
- **Manual Testing**: Visual verification is critical for UI changes
- **Terminal Sizes**: Test on minimum (80x24), standard (120x40), and large (200x60)
- **Accessibility**: Verify keyboard navigation, bell alerts, clear visual hierarchy
- **Performance**: Profile if render times exceed 100ms
- **Consistency**: Use established patterns from Task 14 for all subsequent migrations

## Resources

### Related Documentation
- `docs/TUI_STANDARDS.md` - TUI design principles
- `docs/TUI_DEVELOPER_GUIDE.md` - TUI component development
- `docs/TUI_INTENT_DIAGRAM.md` - Intent architecture
- `docs/LIPGLOSS_BUBBLES_GUIDE.md` - Styling guide
- `docs/rules/master-task-prompt.md` - Development workflow

### Code References
- `internal/cli/components/smart_container.go` - Existing centering logic
- `internal/cli/components/ascii_logo.go` - Logo component
- `internal/cli/app/app.go` - Main menu implementation (good example)
- `internal/cli/intents/contract.go` - Intent interfaces

## Questions & Issues

### During Implementation
If you encounter issues or have questions:

1. Check the specific task file for detailed guidance
2. Review the established patterns in Task 14
3. Consult the TUI_DEVELOPER_GUIDE.md
4. Test on different terminal sizes
5. Run tests frequently to catch regressions early

### After Implementation
Document lessons learned:
- What worked well?
- What could be improved?
- What patterns should future developers follow?
- What edge cases were discovered?

## Completion Status

- [ ] Task 12: Core Components
- [ ] Task 13: Infrastructure Updates
- [ ] Task 14: Migrate Primary Intents
- [ ] Task 15: Migrate Remaining Intents
- [ ] Task 16: Testing and Polish

**Overall Progress**: 0% (Not Started)

---

*This summary document tracks the standardized view implementation across Tasks 12-16. Update completion status as tasks are finished.*
