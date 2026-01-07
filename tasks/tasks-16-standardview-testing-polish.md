# Task 16: StandardView Testing, Documentation, and Polish

## Overview
- **Goal**: Comprehensive testing, documentation, and polish for the standardized view system across all intents
- **Time Estimate**: 3-4 hours
- **Prerequisites**: Tasks 12-15 completed (all components and intents migrated)
- **Status**: Not Started

## Motivation
Ensure the standardized view system is robust, well-tested, and properly documented. Verify consistency across all intents and handle edge cases gracefully.

## Files to Create
- [ ] `docs/STANDARDVIEW_GUIDE.md` (developer guide)
- [ ] `docs/MODAL_PATTERNS.md` (modal usage patterns)
- [ ] `cmd/test_all_views/main.go` (visual test program)

## Files to Modify
- [ ] `docs/TUI_DEVELOPER_GUIDE.md` (update with StandardView section)
- [ ] `docs/TUI_STANDARDS.md` (update with new standards)
- [ ] `AGENTS.md` (update with StandardView architecture info)

## Implementation Checklist

### Phase 1: Preparation (10 min)
- [ ] Create task file
- [ ] Run compliance check (baseline): `make check-compliance`
- [ ] Verify Tasks 12-15 completed successfully
- [ ] Review all migrated intents
- [ ] Run full test suite: `go test ./...`
- [ ] Document current test coverage baseline

### Phase 2: Comprehensive Component Testing (60 min)

#### 2.1 StandardView edge case tests (20 min)
- [ ] Open `internal/cli/components/standard_view_test.go`
- [ ] Add test: Empty content handling
- [ ] Add test: Very long content (scrolling needed)
- [ ] Add test: Terminal too small for logo + content
- [ ] Add test: Terminal resize during display
- [ ] Add test: Nil terminal info handling
- [ ] Add test: Empty breadcrumbs
- [ ] Add test: Very long breadcrumbs (truncation)
- [ ] Add test: Empty help text
- [ ] Add test: Footer separator with different widths
- [ ] Commit: `test(components): add StandardView edge case tests`

#### 2.2 Modal edge case tests (20 min)
- [ ] Open `internal/cli/components/modal_test.go`
- [ ] Add test: Modal larger than terminal
- [ ] Add test: Very long error messages
- [ ] Add test: Progress at edge values (0.0, 1.0)
- [ ] Add test: Fade-in timing accuracy
- [ ] Add test: Auto-dismiss timing accuracy
- [ ] Add test: Multiple rapid modal changes
- [ ] Add test: Modal with empty title/message
- [ ] Add test: Action buttons overflow
- [ ] Commit: `test(components): add modal edge case tests`

#### 2.3 LoadingMessages edge case tests (20 min)
- [ ] Open `internal/cli/components/loading_messages_test.go`
- [ ] Add test: Empty message list
- [ ] Add test: Single message (no rotation)
- [ ] Add test: Rapid rotation requests
- [ ] Add test: Very long messages (truncation)
- [ ] Add test: Reset during rotation
- [ ] Add test: Concurrent access (race conditions)
- [ ] Commit: `test(components): add loading messages edge case tests`

### Phase 3: Terminal Size Testing (45 min)

#### 3.1 Create terminal size test suite (30 min)
- [ ] Create `internal/cli/components/terminal_size_test.go`
- [ ] Add test fixtures for common sizes:
  - Tiny: 80x24
  - Compact: 100x30
  - Normal: 120x40
  - Large: 160x50
  - XLarge: 200x60
  - Ultra-wide: 240x40
  - Tall: 120x80
- [ ] Test StandardView rendering at each size:
  - Logo fits properly
  - Content scales appropriately
  - Footer visible
  - Modals sized correctly
- [ ] Test modal behavior at each size:
  - Error modal sizing
  - Progress modal sizing
  - Success modal sizing
- [ ] Add test: Minimum viable terminal size (80x24)
- [ ] Add test: Content overflow handling
- [ ] Commit: `test(components): add comprehensive terminal size tests`

#### 3.2 Create visual test program (15 min)
- [ ] Create `cmd/test_all_views/main.go`
- [ ] Add menu to select different test scenarios:
  - Basic StandardView (logo, content, footer)
  - StandardView with breadcrumbs
  - StandardView with error modal
  - StandardView with loading modal
  - StandardView with progress modal (0%, 25%, 50%, 75%, 100%)
  - StandardView with success modal
  - StandardView with very long content
  - StandardView with minimal content
  - All intents' actual views
- [ ] Add terminal size indicator in footer
- [ ] Add instructions for manual testing
- [ ] Commit: `feat(test): add visual test program for StandardView`

### Phase 4: Cross-Intent Consistency Testing (30 min)

#### 4.1 Create consistency test suite (30 min)
- [ ] Create `internal/cli/intents/consistency_test.go`
- [ ] Add test: All intents use StandardView
- [ ] Add test: All intents show logo with 2-line spacing
- [ ] Add test: All intents have breadcrumbs
- [ ] Add test: All intents show footer separator
- [ ] Add test: All intents have context-aware help
- [ ] Add test: All intents handle terminal info
- [ ] Add test: All intents use full width content
- [ ] Add helper function to test each intent:
```go
func testIntentStandardView(t *testing.T, intent Intent, stateName string) {
    view := intent.View()
    
    // Assert logo present
    assert.Contains(t, view, "╦╔═╔═╗╦═╗╦╦ ╦╔═╗")
    
    // Assert breadcrumbs present
    assert.Contains(t, view, ">")
    
    // Assert footer separator
    assert.Contains(t, view, "─────")
    
    // Assert help text
    assert.Contains(t, view, "q Quit")
}
```
- [ ] Commit: `test(intents): add cross-intent consistency tests`

### Phase 5: Performance Testing (30 min)

#### 5.1 Add render performance tests (20 min)
- [ ] Create `internal/cli/components/performance_test.go`
- [ ] Add benchmark: StandardView render time
- [ ] Add benchmark: Modal render time
- [ ] Add benchmark: LoadingMessageRotator performance
- [ ] Add benchmark: Full view render (logo + content + modal)
- [ ] Set performance targets:
  - StandardView render < 50ms
  - Modal render < 20ms
  - Full view render < 100ms
- [ ] Run benchmarks: `go test -bench=. ./internal/cli/components/`
- [ ] Document baseline performance
- [ ] Commit: `test(components): add performance benchmarks`

#### 5.2 Memory usage testing (10 min)
- [ ] Add benchmark with memory allocation tracking
- [ ] Test for memory leaks in rotation
- [ ] Test for excessive allocations in rendering
- [ ] Run: `go test -bench=. -benchmem ./internal/cli/components/`
- [ ] Document memory baseline
- [ ] Commit: `test(components): add memory benchmarks`

### Phase 6: Documentation (90 min)

#### 6.1 Create StandardView developer guide (40 min)
- [ ] Create `docs/STANDARDVIEW_GUIDE.md`
- [ ] Add sections:
  - **Overview**: What is StandardView and why use it
  - **Architecture**: How StandardView works
  - **Basic Usage**: Creating and rendering a StandardView
  - **Breadcrumbs**: Adding navigation context
  - **Content**: Setting and styling content
  - **Modals**: Using error, loading, progress, success modals
  - **Help Footer**: Creating context-aware help
  - **Terminal Awareness**: Handling different terminal sizes
  - **Best Practices**: Common patterns and tips
  - **Troubleshooting**: Common issues and solutions
- [ ] Add code examples for each section
- [ ] Add screenshots/ASCII art showing layouts
- [ ] Commit: `docs: add StandardView developer guide`

#### 6.2 Create modal patterns guide (25 min)
- [ ] Create `docs/MODAL_PATTERNS.md`
- [ ] Add sections:
  - **Modal Types**: Error, Loading, Progress, Success, Warning
  - **When to Use Each Type**: Decision guide
  - **Error Modals**: Error handling patterns
  - **Loading Modals**: Long operation handling
  - **Progress Modals**: Multi-step operation tracking
  - **Success Modals**: User feedback patterns
  - **Modal Timing**: Fade-in, auto-dismiss guidelines
  - **Accessibility**: Bell alerts, keyboard navigation
  - **Common Patterns**: Typical modal workflows
- [ ] Add code examples
- [ ] Add visual examples
- [ ] Commit: `docs: add modal patterns guide`

#### 6.3 Update existing documentation (25 min)
- [ ] Open `docs/TUI_DEVELOPER_GUIDE.md`
- [ ] Add section: "Creating Views with StandardView"
  - Link to STANDARDVIEW_GUIDE.md
  - Add quick start example
  - Update component architecture diagram
- [ ] Open `docs/TUI_STANDARDS.md`
- [ ] Update "View Patterns" section:
  - Add StandardView as required pattern
  - Update keyboard shortcuts consistency
  - Add modal standards
  - Update help footer standards
- [ ] Open `AGENTS.md`
- [ ] Update "Architecture Overview" section:
  - Add StandardView component
  - Add modal system
  - Update component diagram
- [ ] Commit: `docs: update existing documentation for StandardView`

### Phase 7: Edge Case Handling (30 min)

#### 7.1 Add graceful degradation (20 min)
- [ ] Open `internal/cli/components/standard_view.go`
- [ ] Add handling for terminal too small:
  - Detect if height < logo height + 10 lines
  - Option to hide logo on very small terminals
  - Ensure content still visible
- [ ] Add handling for nil/invalid terminal info:
  - Use sensible defaults (120x40)
  - Log warning
  - Continue rendering
- [ ] Add content overflow handling:
  - Detect if content exceeds available space
  - Add scroll indicators (↓ More content below)
  - Consider truncation with "..." for very long content
- [ ] Commit: `feat(components): add graceful degradation for edge cases`

#### 7.2 Test graceful degradation (10 min)
- [ ] Add tests for small terminal handling
- [ ] Add tests for nil terminal info
- [ ] Add tests for content overflow
- [ ] Run tests: `go test -v ./internal/cli/components/...`
- [ ] Commit: `test(components): add graceful degradation tests`

### Phase 8: Polish and Refinement (45 min)

#### 8.1 Visual refinement (20 min)
- [ ] Review logo spacing consistency across all intents
- [ ] Review footer separator styling consistency
- [ ] Review breadcrumb styling (color, separator)
- [ ] Review modal border styles
- [ ] Ensure color scheme consistency (use existing styles package)
- [ ] Test with different color schemes if applicable
- [ ] Commit: `style(components): refine visual consistency`

#### 8.2 Animation refinement (15 min)
- [ ] Review modal fade-in timing (should be 150ms)
- [ ] Review loading spinner animation speed (80ms)
- [ ] Review loading message rotation timing (2s)
- [ ] Ensure animations feel smooth and natural
- [ ] Test on slower terminals/systems
- [ ] Commit: `style(components): refine animation timing`

#### 8.3 Accessibility improvements (10 min)
- [ ] Verify terminal bell on errors works
- [ ] Verify keyboard shortcuts are consistent
- [ ] Verify Escape key dismisses modals
- [ ] Verify help text is always visible
- [ ] Test with screen readers if possible
- [ ] Add ARIA-like labels in comments
- [ ] Commit: `feat(components): improve accessibility`

### Phase 9: Integration Testing (30 min)

#### 9.1 End-to-end workflow testing (20 min)
- [ ] Build application: `go build -o kariya ./cmd/kariya`
- [ ] Test complete workflow: Capture → Browse → Generate CV → Export
- [ ] Verify logo visible throughout
- [ ] Verify breadcrumbs update correctly
- [ ] Verify modals appear and dismiss properly
- [ ] Verify loading states during operations
- [ ] Verify success confirmations
- [ ] Test error recovery (trigger errors intentionally)
- [ ] Test cancellation workflows
- [ ] Document any issues found

#### 9.2 Multi-terminal size testing (10 min)
- [ ] Test on 80x24 terminal (minimum)
- [ ] Test on 120x40 terminal (standard)
- [ ] Test on 200x60 terminal (large)
- [ ] Test on ultra-wide terminal (240x40)
- [ ] Verify responsive behavior
- [ ] Verify no layout breaking
- [ ] Document any size-specific issues

### Phase 10: Final Verification (30 min)

#### 10.1 Run complete test suite (15 min)
- [ ] Run all tests: `go test -v ./...`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Run with coverage: `go test -cover ./...`
- [ ] Generate coverage report: `go test -coverprofile=coverage.out ./...`
- [ ] View coverage: `go tool cover -html=coverage.out`
- [ ] Verify coverage > 87% overall
- [ ] Verify all tests pass (100%)
- [ ] Fix any failures
- [ ] Commit if fixes needed: `fix: resolve final test failures`

#### 10.2 Final compliance and quality checks (15 min)
- [ ] Run compliance check: `make check-compliance`
- [ ] Run linter: `golangci-lint run ./...`
- [ ] Format all code: `go fmt ./...`
- [ ] Review all commits across Tasks 12-16
- [ ] Verify atomic commits
- [ ] Verify conventional commit messages
- [ ] Review code for TODOs or FIXMEs
- [ ] Update AGENTS.md with final status
- [ ] Commit: `docs: update project status for StandardView completion`

## Testing Instructions

### Automated Tests
```bash
# Run all tests
go test -v ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem ./internal/cli/components/

# Run specific test suites
go test -v ./internal/cli/components/ -run TestStandardView
go test -v ./internal/cli/components/ -run TestModal
go test -v ./internal/cli/intents/ -run TestConsistency
```

### Manual Testing with Visual Test Program
```bash
# Build and run visual test program
go build -o test_views ./cmd/test_all_views
./test_views

# Test different terminal sizes
# Resize terminal while program is running
# Verify responsive behavior
```

### Cross-Intent Manual Testing Checklist
For each intent:
- [ ] CaptureEvent
  - [ ] Logo visible with 2-line spacing
  - [ ] Breadcrumbs update per state
  - [ ] Footer separator visible
  - [ ] Context help state-aware
  - [ ] Modals work (error, loading, success)
- [ ] BrowseTimeline
  - [ ] All above items verified
- [ ] GenerateCV
  - [ ] All above items verified
  - [ ] Progress modal shows during generation
- [ ] ExportArtifact
  - [ ] All above items verified
  - [ ] Progress modal shows during export
- [ ] ConfigureSystem
  - [ ] All above items verified

## Acceptance Criteria
- [ ] All component tests pass (100%)
- [ ] All intent tests pass (100%)
- [ ] Consistency tests pass for all intents
- [ ] Edge case handling verified
- [ ] Terminal size tests pass for all sizes
- [ ] Performance benchmarks meet targets (< 100ms render)
- [ ] Memory usage acceptable (no leaks)
- [ ] STANDARDVIEW_GUIDE.md created and complete
- [ ] MODAL_PATTERNS.md created and complete
- [ ] Existing docs updated with StandardView info
- [ ] Visual test program works
- [ ] Graceful degradation implemented
- [ ] Animations smooth and natural
- [ ] Accessibility verified
- [ ] End-to-end workflows tested
- [ ] Multi-terminal size testing complete
- [ ] Code coverage > 87%
- [ ] All code passes linting
- [ ] All code formatted
- [ ] Compliance check passes
- [ ] AGENTS.md updated with completion status

## Rollback Plan
If critical issues are discovered:
1. Identify problematic task/commits
2. Run: `git revert <commit-hash>` for each commit in reverse order
3. Alternative: `git reset --hard <commit-before-task-12>`
4. Re-run full test suite to verify stability
5. Review issues and create new task files if needed

## Notes
- This task ensures the standardized view system is production-ready
- Comprehensive testing is critical for user-facing changes
- Documentation helps future developers maintain consistency
- Visual test program is valuable for ongoing development
- Performance benchmarks establish baseline for future optimizations
- Edge case handling prevents crashes and poor UX
- Cross-intent consistency is key to professional appearance

## Success Metrics
- **Test Coverage**: > 87% overall, > 90% for new components
- **Performance**: All renders < 100ms
- **Consistency**: 100% of intents use StandardView
- **Documentation**: Complete guides for developers
- **User Experience**: Smooth, responsive, professional
- **Quality**: 0 regressions, all tests passing

## Dependencies
- Requires Tasks 12-15 completed
- Uses all StandardView components
- Tests all migrated intents
- Comprehensive documentation of system

## Post-Completion
After this task:
1. All intents use standardized centered views
2. Logo visible on every screen
3. Consistent layout throughout application
4. Comprehensive testing ensures quality
5. Documentation supports future development
6. System is production-ready

Update AGENTS.md to mark StandardView implementation as complete.
