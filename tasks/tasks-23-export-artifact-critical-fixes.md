# Task 23: ExportArtifact Critical Fixes

**Created**: 2026-01-08
**Status**: IN PROGRESS (Phase 5/10 Complete - 50% Done)
**Priority**: CRITICAL
**Estimated Time**: 4-6 hours (3 hours spent)
**Related**: Codebase Audit (2026-01-08)

---

## Progress Tracking

### Completed Phases ✅

**Phase 1: Critical Bug Fixes** (Tasks 1-3) - COMPLETE ✅
- ✅ Fixed `formatBytes()` function (was showing garbage characters)
- ✅ Fixed scroll percentage display (was showing garbage characters)
- ✅ Added 14 comprehensive tests for both fixes
- **Commit**: `4c0309f`, `3ab85a1`

**Phase 2: Format Updates** (Tasks 4-7) - COMPLETE ✅
- ✅ Added `ExportFormatYAML` constant
- ✅ Updated CV formats (removed PDF/JSON, added TXT/MD/YAML)
- ✅ Updated Events/Facts/Bursts formats (added YAML)
- ✅ Removed `ExportDestinationEmail` (only File/Clipboard remain)
- **Commit**: `2a631b9`

**Phase 3: Vim Navigation** (Tasks 8-10) - COMPLETE ✅
- ✅ Added `j` key for down navigation in format/destination selection
- ✅ Added `k` key for up navigation in format/destination selection
- ✅ Added 4 explicit vim navigation tests
- **Commit**: `4fa6269`

**Phase 4: Service Integration** (Tasks 11-14) - COMPLETE ✅
- ✅ Added helper functions (DefaultArtifactTypes, DefaultSupportedFormats, DefaultFormats, DefaultDestinations)
- ✅ Updated ExportArtifactContext to include ExportService, CVGenerationService, CareerService, and repositories
- ✅ Changed NewExportArtifactIntent to accept *ExportArtifactContext instead of context.Context
- ✅ Created NewTestExportArtifactContext helper for test consistency
- ✅ Updated all 110+ test calls across 4 test files
- ✅ Updated app.go to properly initialize ExportArtifact with all services
- ✅ Removed obsolete NewExportArtifactContext function
- **Commit**: `3ffcd79`
- **Test Results**: 110/110 ExportArtifact tests passing ✓

**Phase 5: CV Selection State** (Tasks 15-18) - COMPLETE ✅
- ✅ Added `CVConfigManager` field to `ExportArtifactContext`
- ✅ Implemented `loadAvailableCVs()` to fetch CV configs from ConfigManager
- ✅ Updated test helpers to include ConfigManager (NewTestExportArtifactContext, NewTestExportArtifactContextWithServices)
- ✅ Updated `registerAllIntents()` signature to accept ConfigManager parameter
- ✅ Updated app.go to pass ConfigManager to ExportArtifact intent
- ✅ Added 4 comprehensive tests for loadAvailableCVs() function
- ✅ Added tests for empty CV list scenario and error handling
- **Test Results**: 526/526 tests passing (all intents) ✓
- **Build Status**: Successful ✅

### Current Status: Phase 6 (Next)

**Phase 6: Format Mapping** (Tasks 19-20) - READY TO START
- Map intent export formats to service export formats
- Update exportCV() to use proper format mapping

### Remaining Phases (Tasks 19-37, 19 tasks)

- **Phase 6**: Format Mapping (Tasks 19-20) - Ready
- **Phase 7**: Real Export (Tasks 21-24)
- **Phase 8**: Real Previews (Tasks 25-30)
- **Phase 9**: LoadingRotator (Tasks 31-33)
- **Phase 10**: Final Verification (Tasks 34-37)

**Estimated Remaining Time**: 1.5-2 hours

---

## Overview

The ExportArtifact intent is **completely non-functional** - exports don't actually export data, previews show static mock data, and there are critical bugs that cause display issues.

**Current State**:
- Export operation is stubbed (always returns `/tmp/export.*` with hardcoded 1024 bytes)
- All preview data is static/hardcoded (no real user data)
- `formatBytes()` has a critical bug (displays garbage characters)
- Scroll percentage display has the same bug
- ExportService exists and is fully functional but not integrated

**Impact**: Users cannot export their CVs, events, or other artifacts. The feature is advertised but non-functional.

---

## Files to Modify

- [x] `internal/cli/intents/export_artifact.go` (main implementation) - Phases 1-4 complete
- [x] `internal/cli/intents/export_artifact_intent.go` (add ExportService dependency) - Phase 4 complete
- [x] `internal/cli/intents/export_artifact_test.go` (update tests for real export) - Phases 1-4 complete
- [x] `internal/cli/app/app.go` (service integration) - Phase 4 complete
- [x] `internal/cli/intents/benchmarks_test.go` (update benchmarks) - Phase 4 complete
- [x] `internal/cli/intents/consistency_test.go` (update consistency tests) - Phase 4 complete
- [x] `internal/cli/intents/export_artifact_escape_test.go` (update escape tests) - Phase 4 complete

---

## Implementation Plan

### Phase 1: Critical Bug Fixes (1 hour)

#### 1.1 Fix formatBytes() Function
**Location**: `internal/cli/intents/export_artifact.go:930-942`

**Current Code** (BROKEN):
```go
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return string(rune(bytes)) + " B"  // BUG: converts to character
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    units := []string{"", "K", "M", "G", "T"}
    return string(rune(bytes/div)) + " " + units[exp] + "B"  // BUG
}
```

**Fix**:
```go
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)  // FIX: format as number
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    units := []string{"", "K", "M", "G", "T"}
    return fmt.Sprintf("%d %sB", bytes/div, units[exp])  // FIX
}
```

**Tasks**:
- [x] Replace `string(rune(bytes))` with `fmt.Sprintf("%d", bytes)`
- [x] Add test for formatBytes with various sizes (1, 1024, 1048576)
- [x] Verify output: "1 B", "1 KB", "1 MB"

#### 1.2 Fix Scroll Percentage Display
**Location**: `internal/cli/intents/export_artifact.go:601-604`

**Current Code** (BROKEN):
```go
if len(m.previewLines) > viewHeight {
    scrollPercent := (m.scrollOffset * 100) / len(m.previewLines)
    s += "\n[" + string(rune(scrollPercent/10)) + "% scrolled]"  // BUG
}
```

**Fix**:
```go
if len(m.previewLines) > viewHeight {
    scrollPercent := (m.scrollOffset * 100) / len(m.previewLines)
    s += fmt.Sprintf("\n[%d%% scrolled]", scrollPercent)  // FIX
}
```

**Tasks**:
- [x] Replace `string(rune(scrollPercent/10))` with `fmt.Sprintf("%d%%", scrollPercent)`
- [x] Test with different scroll positions
- [x] Verify display shows "0%", "50%", "100%"

**Verification**:
```bash
# Run tests
go test ./internal/cli/intents -v -run TestExportArtifact

# Check for rune conversion bugs
rg "string\(rune\(" internal/cli/intents/export_artifact.go
```

---

### Phase 2: Service Integration (2-3 hours)

#### 2.1 Add ExportService Dependency
**Location**: `internal/cli/intents/export_artifact_intent.go`

**Add to ExportArtifactContext**:
```go
type ExportArtifactContext struct {
    ArtifactTypes     []ExportArtifactType
    SupportedFormats  map[ExportArtifactType][]ExportFormat
    ExportService     *cv.ExportService  // ADD THIS
    CVService         *careerservice.Service  // For fetching data
    EventRepository   careerrepo.Repository   // For events
    FactRepository    careerrepo.FactRepository  // For facts
}
```

**Tasks**:
- [ ] Add ExportService field to context
- [ ] Add repository fields for data fetching
- [ ] Update context initialization in app.go
- [ ] Pass dependencies from GlobalContext

#### 2.2 Implement Real Export Logic
**Location**: `internal/cli/intents/export_artifact.go:506-518`

**Current Code** (STUB):
```go
func (m *ExportArtifactModel) startExport() tea.Cmd {
    return func() tea.Msg {
        result := NewExportArtifactResult(
            true,
            m.config.ArtifactType,
            m.config.Format,
            m.config.Destination,
            "/tmp/export."+string(m.config.Format),  // HARDCODED
            1024,  // HARDCODED
        )
        return ExportCompleteMsg{Result: result}
    }
}
```

**New Implementation**:
```go
func (m *ExportArtifactModel) startExport() tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        
        // Fetch data based on artifact type
        var content string
        var err error
        
        switch m.config.ArtifactType {
        case ExportArtifactTypeCV:
            content, err = m.exportCV(ctx)
        case ExportArtifactTypeEvents:
            content, err = m.exportEvents(ctx)
        case ExportArtifactTypeFacts:
            content, err = m.exportFacts(ctx)
        // ... other types
        }
        
        if err != nil {
            return ExportCompleteMsg{
                Result: &ExportArtifactResult{
                    Success: false,
                    Error:   err.Error(),
                },
            }
        }
        
        // Handle destination
        var filePath string
        var fileSize int64
        
        switch m.config.Destination {
        case ExportDestinationFile:
            filePath, err = m.saveToFile(content)
            if err != nil {
                return ExportCompleteMsg{/* error result */}
            }
            fileInfo, _ := os.Stat(filePath)
            fileSize = fileInfo.Size()
            
        case ExportDestinationClipboard:
            err = m.context.ExportService.CopyToClipboard(ctx, content)
            if err != nil {
                return ExportCompleteMsg{/* error result */}
            }
            filePath = "(clipboard)"
            fileSize = int64(len(content))
        }
        
        return ExportCompleteMsg{
            Result: &ExportArtifactResult{
                Success: true,
                ArtifactType: m.config.ArtifactType,
                Format: m.config.Format,
                Destination: m.config.Destination,
                FilePath: filePath,
                FileSize: fileSize,
            },
        }
    }
}
```

**Helper Methods to Add**:
```go
func (m *ExportArtifactModel) exportCV(ctx context.Context) (string, error) {
    // Get CV from context or fetch latest
    cv := m.selectedCV  // Assuming we added CV selection
    
    switch m.config.Format {
    case ExportFormatText:
        return m.context.ExportService.ExportToText(ctx, cv, sections, bullets)
    case ExportFormatMarkdown:
        return m.context.ExportService.ExportToMarkdown(ctx, cv, sections, bullets)
    case ExportFormatYAML:
        return m.context.ExportService.ExportToYAML(ctx, cv)
    }
}

func (m *ExportArtifactModel) exportEvents(ctx context.Context) (string, error) {
    // Fetch events from repository
    events, err := m.context.EventRepository.List(ctx, filters)
    if err != nil {
        return "", err
    }
    
    // Format based on export format
    return m.formatEvents(events, m.config.Format), nil
}

func (m *ExportArtifactModel) saveToFile(content string) (string, error) {
    // Use ExportService.SaveToFile with appropriate name
    cvName := "export"  // Or get from selected artifact
    return m.context.ExportService.SaveToFile(content, cvName, 
        cv.ExportFormat(m.config.Format))
}
```

**Tasks**:
- [ ] Replace stub startExport with real implementation
- [ ] Add exportCV() helper method
- [ ] Add exportEvents() helper method
- [ ] Add exportFacts() helper method
- [ ] Add saveToFile() helper method
- [ ] Add error handling for each export type
- [ ] Test each export path (CV, Events, Facts)

**Verification**:
```bash
# Test CV export
go test ./internal/cli/intents -v -run "TestExportArtifact.*CV"

# Test file creation
ls ~/.kariya/cv_exports/

# Test clipboard
# (requires manual verification)
```

---

### Phase 3: Real Preview Data (1-2 hours)

#### 3.1 Replace Mock Preview Data
**Location**: `internal/cli/intents/export_artifact.go:663-928`

**Current**: All `generate*Preview()` functions return hardcoded strings

**Fix**: Fetch real data from repositories

**Example for CV Preview**:
```go
func (m *ExportArtifactModel) generateCVPreview() string {
    if m.selectedCV == nil {
        return "No CV selected"
    }
    
    ctx := context.Background()
    
    // Generate actual preview using ExportService
    var preview string
    var err error
    
    switch m.config.Format {
    case ExportFormatJSON:
        // Use YAML export then convert to JSON format display
        preview, err = m.context.ExportService.ExportToYAML(ctx, m.selectedCV)
        if err != nil {
            return fmt.Sprintf("Error generating preview: %s", err)
        }
        // Format as JSON display (pretty-print)
        
    case ExportFormatText:
        preview, err = m.context.ExportService.ExportToText(ctx, 
            m.selectedCV, sections, bullets)
        
    case ExportFormatMarkdown:
        preview, err = m.context.ExportService.ExportToMarkdown(ctx,
            m.selectedCV, sections, bullets)
    }
    
    if err != nil {
        return fmt.Sprintf("Error: %s", err)
    }
    
    return preview
}
```

**Tasks**:
- [ ] Update generateCVPreview() to use real CV data
- [ ] Update generateEventsPreview() to fetch from repository
- [ ] Update generateFactsPreview() to fetch from repository
- [ ] Update generateBurstsPreview() to fetch from repository
- [ ] Update generateProfilePreview() to fetch from context
- [ ] Add error handling for missing data
- [ ] Test previews with actual data

**Verification**:
```bash
# Test preview generation
go test ./internal/cli/intents -v -run "TestExportArtifact.*Preview"
```

---

### Phase 4: Consistency Fixes (30 min)

#### 4.1 Add Vim Navigation
**Location**: `internal/cli/intents/export_artifact.go`

**Add to updateSelectFormat()** (line 320-327):
```go
case "up", "k":  // ADD "k"
    if m.selectedIndex > 0 {
        m.selectedIndex--
    }
case "down", "j":  // ADD "j"
    formats := m.context.SupportedFormats[m.config.ArtifactType]
    if m.selectedIndex < len(formats)-1 {
        m.selectedIndex++
    }
```

**Add to updateSelectDest()** (line 350-357):
```go
case "up", "k":  // ADD "k"
    if m.selectedIndex > 0 {
        m.selectedIndex--
    }
case "down", "j":  // ADD "j"
    destinations := []ExportDestination{
        ExportDestinationFile,
        ExportDestinationClipboard,
        ExportDestinationEmail,
    }
    if m.selectedIndex < len(destinations)-1 {
        m.selectedIndex++
    }
```

**Tasks**:
- [ ] Add "k" to up navigation in updateSelectFormat
- [ ] Add "j" to down navigation in updateSelectFormat
- [ ] Add "k" to up navigation in updateSelectDest
- [ ] Add "j" to down navigation in updateSelectDest
- [ ] Test vim navigation works in both states

#### 4.2 Use LoadingRotator
**Location**: `internal/cli/intents/export_artifact.go`

**Current viewInProgress()** (shows static text):
```go
func (m *ExportArtifactModel) viewInProgress() string {
    s := "Exporting...\n\n"
    s += "Please wait while your artifact is being exported.\n"
    return s
}
```

**Fix**:
```go
func (m *ExportArtifactModel) viewInProgress() string {
    message := m.loadingRotator.GetMessage()
    s := message + "\n\n"
    s += "Please wait while your artifact is being exported.\n"
    return s
}

// Add tick command to rotate messages
func (m *ExportArtifactModel) tickLoadingRotator() tea.Cmd {
    return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
        return LoadingTickMsg{}
    })
}

// Handle in Update
case LoadingTickMsg:
    if m.state == ExportStateInProgress {
        m.loadingRotator.Rotate()
        return m, m.tickLoadingRotator()
    }
```

**Tasks**:
- [ ] Use loadingRotator.GetMessage() in viewInProgress
- [ ] Add tick command for rotation
- [ ] Handle LoadingTickMsg in Update
- [ ] Start rotation when entering InProgress state

**Verification**:
```bash
# Run all tests
go test ./internal/cli/intents -v -run TestExportArtifact

# Full test suite
go test ./... -v

# Race detector
go test -race ./internal/cli/intents
```

---

## Acceptance Criteria

### Must Have
- [ ] `formatBytes()` displays correct values ("1 KB", "1 MB", etc.)
- [ ] Scroll percentage displays correctly ("50%", not garbage)
- [ ] Export to file actually creates files in `~/.kariya/cv_exports/`
- [ ] Export to clipboard works (content in clipboard)
- [ ] CV previews show actual user CV data
- [ ] Event/Fact/Burst previews show actual user data
- [ ] All export formats work (Text, Markdown, YAML)
- [ ] All tests passing (maintain 2,078/2,078)
- [ ] Zero race conditions
- [ ] Build successful

### Should Have
- [ ] Vim j/k navigation works in format and destination selection
- [ ] LoadingRotator shows rotating messages during export
- [ ] Error messages are clear and actionable
- [ ] File paths are displayed to user after export

### Nice to Have
- [ ] Progress indicator for large exports
- [ ] Email destination implemented (currently stub)
- [ ] PDF export implemented (currently not supported)

---

## Testing Strategy

### TDD Approach

**Red**: Write failing test for each fix
```go
// Test formatBytes
It("should format bytes correctly", func() {
    Expect(formatBytes(0)).To(Equal("0 B"))
    Expect(formatBytes(1023)).To(Equal("1023 B"))
    Expect(formatBytes(1024)).To(Equal("1 KB"))
    Expect(formatBytes(1048576)).To(Equal("1 MB"))
})

// Test real export
It("should export CV to file", func() {
    result := performExport(cv, ExportFormatText, ExportDestinationFile)
    Expect(result.Success).To(BeTrue())
    Expect(result.FilePath).To(ContainSubstring(".kariya/cv_exports"))
    Expect(fileExists(result.FilePath)).To(BeTrue())
})
```

**Green**: Implement minimal fix
**Refactor**: Clean up implementation

### Verification Commands
```bash
# Unit tests
go test ./internal/cli/intents -v -run TestExportArtifact

# Integration test (manual)
go build -o kariya ./cmd/cli
./kariya
# Navigate to: Export > CV > Text > File
# Verify file created in ~/.kariya/cv_exports/

# Clipboard test (manual)
# Export to clipboard, then paste - should show CV content

# Full test suite
go test ./... -v
go test -race ./...

# Compliance check
make check-compliance
```

---

## Risk Assessment

### High Risk
1. **ExportService integration breaks existing code**
   - Mitigation: Add dependencies incrementally, test after each
   
2. **Preview generation is slow with large datasets**
   - Mitigation: Implement pagination, limit preview size

### Medium Risk
3. **File system permissions issues**
   - Mitigation: Use os.UserHomeDir(), handle permission errors gracefully

4. **Clipboard doesn't work on all platforms**
   - Mitigation: Check platform, show clear error if unsupported

### Low Risk
5. **Test failures due to mock data changes**
   - Mitigation: Update tests to use test fixtures instead of mocks

---

## References

### Related Documents
- `docs/rules/master-task-prompt.md` - Task execution workflow
- `docs/TUI_STANDARDS.md` - TUI design standards
- `AGENTS.md` - Project overview

### Related Files
- `internal/service/career/cv/export_service.go` - Export service (284 lines, fully functional)
- `internal/cli/intents/generate_cv_intent.go` - Example of ExportService usage (lines 940-987)
- `internal/cli/intents/export_artifact.go` - Main implementation file (943 lines)
- `internal/cli/intents/export_artifact_intent.go` - Context and dependencies (175 lines)

### Service Methods Available
```go
// From export_service.go
ExportToText(ctx, cv, sections, bullets) (string, error)
ExportToMarkdown(ctx, cv, sections, bullets) (string, error)
ExportToYAML(ctx, cv) (string, error)
SaveToFile(content, cvName, format) (string, error)
CopyToClipboard(ctx, content) error
GetExportPath() (string, error)
```

---

## Notes

### Key Decisions
1. **Use existing ExportService** - Don't duplicate export logic
2. **Fix bugs first** - Small, safe changes before big refactoring
3. **TDD approach** - Write tests for each phase
4. **Incremental integration** - Add dependencies one at a time

### Things to Watch
- File system permissions on different platforms
- Clipboard library compatibility
- Large file export performance
- Preview generation with missing data

---

## Session Progress Summary

### Session 1: 2026-01-08 (2 hours)

**Completed Work**:
- ✅ Phase 1: Critical Bug Fixes (formatBytes, scroll percentage) - 14 tests added
- ✅ Phase 2: Format Updates (YAML support, removed PDF/Email) - Updated all format mappings
- ✅ Phase 3: Vim Navigation (j/k keys) - 4 navigation tests added
- ✅ Phase 4: Service Integration (ExportService, repositories) - 110/110 tests passing

**Commits**:
- `4c0309f` - test: add tests for formatBytes() utility function
- `3ab85a1` - test: add tests for scroll percentage display in preview
- `2a631b9` - refactor(export): update export formats and remove email destination
- `4fa6269` - feat(export): add vim j/k navigation to format and destination selection
- `3ffcd79` - feat(export): integrate services and repositories into ExportArtifact intent

**Files Modified**: 8 files (+432/-84 lines)
**Test Status**: 110/110 passing ✅
**Build Status**: Successful ✅

**Next Session**: Continue with Phase 6 (Format Mapping)

### Session 2: 2026-01-08 (1 hour)

**Completed Work**:
- ✅ Phase 5: CV Selection State (ConfigManager integration) - 4 new tests added

**Commits**:
- `f5f43b6` - feat(export): implement CV selection with ConfigManager integration

**Files Modified**: 4 files (+157/-19 lines)
**Test Status**: 526/526 passing ✅
**Build Status**: Successful ✅

**Implementation Details**:
- Added CVConfigManager to ExportArtifactContext
- Implemented loadAvailableCVs() to fetch from ConfigManager and convert to CVView
- Updated test helpers and app.go initialization
- All existing CV selection tests continue to pass (navigation, view rendering, etc.)

**Next Session**: Continue with Phase 6 (Format Mapping)

---

**Last Updated**: 2026-01-08
**Author**: AI Assistant (via OpenCode)
**Status**: In Progress (50% complete - 5/10 phases done)
