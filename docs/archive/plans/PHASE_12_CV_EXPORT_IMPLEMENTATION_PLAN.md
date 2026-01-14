---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Phase 12: CV Export and Save Functionality - Implementation Plan

**Status**: 🚀 **IN PROGRESS**

**Objective**: Complete the CV generation workflow by implementing export to file and clipboard functionality.

---

## Overview

Phase 12 extends the GenerateCV intent to include a complete export workflow. After generating and previewing a CV, users will be able to:
1. Select export format (Text, Markdown, YAML)
2. Choose save location or copy to clipboard
3. Save CV to file with proper naming
4. Confirm successful export

---

## Architecture

### Current State Machine (Phase 11)
```
SelectProfile → SelectAudience → Generating → Preview → Review → Confirm
```

### New State Machine (Phase 12)
```
SelectProfile → SelectAudience → Generating → Preview → Review → Confirm → Export
                                                                              ↓
                                                                        SelectFormat
                                                                              ↓
                                                                        SelectLocation
                                                                              ↓
                                                                        Exporting
                                                                              ↓
                                                                        ExportComplete
```

### New States

```go
const (
	GenerateCVStateExportSelectFormat ExportState = "export_select_format"
	GenerateCVStateExportSelectLocation ExportState = "export_select_location"
	GenerateCVStateExporting ExportState = "exporting"
	GenerateCVStateExportComplete ExportState = "export_complete"
)
```

---

## Export Workflow

### 1. Export Format Selection
**User sees**:
```
Export CV

Select export format:
▶ Text (plain text)
  Markdown (formatted)
  YAML (configuration)

↑/k up, ↓/j down, Enter to select, Esc to go back
```

**Supported Formats**:
- **Text**: Plain text with headers and bullets
- **Markdown**: GitHub-flavored markdown for easy viewing
- **YAML**: Structured YAML for data interchange

### 2. Save Location Selection
**User sees**:
```
Save Location

Where would you like to save the CV?

▶ Save to file (~/kariya-cvs/)
  Copy to clipboard
  Cancel

↑/k up, ↓/j down, Enter to select, Esc to go back
```

**Options**:
- **Save to file**: Saves to `~/kariya-cvs/` directory with auto-generated filename
- **Copy to clipboard**: Copies CV content to system clipboard
- **Cancel**: Returns to review state

### 3. Export Progress
**User sees**:
```
⏳ Exporting CV...

Format: Text
Destination: ~/kariya-cvs/
CV Name: Staff Engineer

Processing...
```

### 4. Export Complete
**User sees**:
```
✅ Export Complete!

Format: Text
Location: ~/kariya-cvs/Staff-Engineer-2026-01-04.txt
Size: 2.5 KB

Press Enter to continue, Esc to go back
```

---

## Implementation Details

### 1. Update GenerateCV Data Structures

**File**: `internal/cli/intents/generate_cv.go`

Add new types:
```go
type ExportState string

const (
	GenerateCVStateExportSelectFormat ExportState = "export_select_format"
	GenerateCVStateExportSelectLocation ExportState = "export_select_location"
	GenerateCVStateExporting ExportState = "exporting"
	GenerateCVStateExportComplete ExportState = "export_complete"
)

type ExportOption string

const (
	ExportOptionSaveToFile ExportOption = "save_to_file"
	ExportOptionClipboard ExportOption = "clipboard"
	ExportOptionCancel ExportOption = "cancel"
)

type ExportFormat string

const (
	ExportFormatText ExportFormat = "text"
	ExportFormatMarkdown ExportFormat = "markdown"
	ExportFormatYAML ExportFormat = "yaml"
)

// Add to GenerateCVModel
type GenerateCVModel struct {
	// ... existing fields ...
	exportState ExportState
	selectedExportFormat ExportFormat
	selectedExportOption ExportOption
	exportResult *ExportResult
	exportError error
}

// Add to GenerateCVResult
type GenerateCVResult struct {
	// ... existing fields ...
	ExportedPath string
	ExportedFormat string
}
```

Add message types:
```go
type ExportFormatSelectedMsg struct {
	Format ExportFormat
}

type ExportOptionSelectedMsg struct {
	Option ExportOption
}

type ExportCompleteMsg struct {
	Path string
	Error error
}
```

### 2. Update GenerateCVIntent

**File**: `internal/cli/intents/generate_cv_intent.go`

Add state handlers:
```go
// In Update()
case GenerateCVStateExportSelectFormat:
	return i.updateExportSelectFormat(msg)
case GenerateCVStateExportSelectLocation:
	return i.updateExportSelectLocation(msg)
case GenerateCVStateExporting:
	return i.updateExporting(msg)
case GenerateCVStateExportComplete:
	return i.updateExportComplete(msg)

// In View()
case GenerateCVStateExportSelectFormat:
	return i.viewExportSelectFormat()
case GenerateCVStateExportSelectLocation:
	return i.viewExportSelectLocation()
case GenerateCVStateExporting:
	return i.viewExporting()
case GenerateCVStateExportComplete:
	return i.viewExportComplete()
```

Update confirm handler:
```go
// In updateConfirm()
case "x", "e":  // x for export
	i.state.exportState = GenerateCVStateExportSelectFormat
	i.state.selectedIndex = 0
	return nil
```

Implement export handlers:
```go
func (i *GenerateCVIntent) updateExportSelectFormat(msg tea.Msg) tea.Cmd {
	// Handle format selection
}

func (i *GenerateCVIntent) updateExportSelectLocation(msg tea.Msg) tea.Cmd {
	// Handle location selection
}

func (i *GenerateCVIntent) updateExporting(msg tea.Msg) tea.Cmd {
	// Handle export progress
}

func (i *GenerateCVIntent) updateExportComplete(msg tea.Msg) tea.Cmd {
	// Handle export completion
}

func (i *GenerateCVIntent) exportCVAsync() tea.Cmd {
	// Async export operation
}
```

Implement view methods:
```go
func (i *GenerateCVIntent) viewExportSelectFormat() string {
	// Render format selection UI
}

func (i *GenerateCVIntent) viewExportSelectLocation() string {
	// Render location selection UI
}

func (i *GenerateCVIntent) viewExporting() string {
	// Render export progress UI
}

func (i *GenerateCVIntent) viewExportComplete() string {
	// Render export success UI
}
```

### 3. Export Async Operation

```go
func (i *GenerateCVIntent) exportCVAsync() tea.Cmd {
	return func() tea.Msg {
		ctx := i.context.AppContext

		// Get export content based on format
		var content string
		var err error

		switch i.state.selectedExportFormat {
		case ExportFormatText:
			content, err = i.context.ExportService.ExportToText(ctx, i.state.generatedCV, nil, nil)
		case ExportFormatMarkdown:
			content, err = i.context.ExportService.ExportToMarkdown(ctx, i.state.generatedCV, nil, nil)
		case ExportFormatYAML:
			content, err = i.context.ExportService.ExportToYAML(ctx, i.state.generatedCV, nil, nil)
		}

		if err != nil {
			return ExportCompleteMsg{Path: "", Error: err}
		}

		// Save based on option
		switch i.state.selectedExportOption {
		case ExportOptionSaveToFile:
			path, err := i.context.ExportService.SaveToFile(ctx, i.state.generatedCV.Name,
				cv.ExportFormat(i.state.selectedExportFormat), content)
			if err != nil {
				return ExportCompleteMsg{Path: "", Error: err}
			}
			return ExportCompleteMsg{Path: path, Error: nil}

		case ExportOptionClipboard:
			err := i.context.ExportService.CopyToClipboard(ctx, content)
			if err != nil {
				return ExportCompleteMsg{Path: "", Error: err}
			}
			return ExportCompleteMsg{Path: "clipboard", Error: nil}
		}

		return ExportCompleteMsg{Path: "", Error: fmt.Errorf("unknown export option")}
	}
}
```

---

## UI Mockups

### Export Format Selection
```
╭─────────────────────────────────────────╮
│ 📤 Export CV                            │
│                                         │
│ Select export format:                   │
│                                         │
│ ▶ Text (plain text)                     │
│   • Best for: Email, simple sharing     │
│                                         │
│   Markdown (formatted)                  │
│   • Best for: GitHub, documentation    │
│                                         │
│   YAML (configuration)                  │
│   • Best for: Data interchange         │
│                                         │
├─────────────────────────────────────────┤
│ ↑/k up, ↓/j down, Enter to select      │
╰─────────────────────────────────────────╯
```

### Save Location Selection
```
╭─────────────────────────────────────────╮
│ 💾 Save Location                        │
│                                         │
│ Format: Text                            │
│ CV Name: Staff Engineer                 │
│                                         │
│ Where to save?                          │
│                                         │
│ ▶ Save to file                          │
│   Location: ~/kariya-cvs/              │
│                                         │
│   Copy to clipboard                     │
│   Paste anywhere                        │
│                                         │
│   Cancel                                │
│   Return to review                      │
│                                         │
├─────────────────────────────────────────┤
│ ↑/k up, ↓/j down, Enter to select      │
╰─────────────────────────────────────────╯
```

### Export Progress
```
╭─────────────────────────────────────────╮
│ ⏳ Exporting CV...                      │
│                                         │
│ Format: Text                            │
│ Destination: ~/kariya-cvs/              │
│ CV Name: Staff Engineer                 │
│                                         │
│ Processing...                           │
│ • Formatting content                    │
│ • Generating filename                   │
│ • Writing to disk                       │
│                                         │
├─────────────────────────────────────────┤
│ Press q to cancel                       │
╰─────────────────────────────────────────╯
```

### Export Complete
```
╭─────────────────────────────────────────╮
│ ✅ Export Complete!                     │
│                                         │
│ Format: Text                            │
│ Location:                               │
│ ~/kariya-cvs/Staff-Engineer-2026-01-04 │
│                                         │
│ File size: 2.5 KB                       │
│ Saved at: 2026-01-04 10:30:45          │
│                                         │
│ You can now share this file!            │
│                                         │
├─────────────────────────────────────────┤
│ Enter to continue, Esc to go back       │
╰─────────────────────────────────────────╯
```

---

## Error Handling

### Export Errors
```go
// If export format is invalid
ExportCompleteMsg{Path: "", Error: fmt.Errorf("unsupported format")}

// If file write fails
ExportCompleteMsg{Path: "", Error: fmt.Errorf("failed to write file: %v", err)}

// If clipboard copy fails
ExportCompleteMsg{Path: "", Error: fmt.Errorf("clipboard unavailable")}
```

### Error Display
```
╭─────────────────────────────────────────╮
│ ❌ Export Failed                        │
│                                         │
│ Error: Failed to write file              │
│ Reason: Permission denied               │
│                                         │
│ Try saving to a different location or   │
│ copy to clipboard instead.              │
│                                         │
├─────────────────────────────────────────┤
│ Esc to go back, q to cancel             │
╰─────────────────────────────────────────╯
```

---

## File Naming Convention

Generated filenames follow this pattern:
```
{ProfileName}-{Date}.{extension}

Examples:
Staff-Engineer-2026-01-04.txt
Principal-Engineer-2026-01-04.md
Staff-Engineer-2026-01-04.yaml
```

### Directory Structure
```
~/.kariya/cvs/
├── Staff-Engineer-2026-01-04.txt
├── Staff-Engineer-2026-01-04.md
├── Principal-Engineer-2025-12-20.txt
└── Staff-Engineer-2025-12-15.yaml
```

---

## Testing Strategy

### Unit Tests
1. Export format selection state transitions
2. Export location selection state transitions
3. Export async operation success/failure
4. File naming generation
5. Error handling and recovery

### Integration Tests
1. Complete export workflow (format → location → export → confirm)
2. File save with various formats
3. Clipboard copy operation
4. Error scenarios (no disk space, permission denied, etc.)

### Manual Testing
1. Export to text file
2. Export to markdown file
3. Export to YAML file
4. Copy to clipboard
5. Verify file contents
6. Verify file location
7. Test error scenarios

---

## Success Criteria

✅ **Functional Requirements**
- [ ] Export format selection works
- [ ] Export location selection works
- [ ] Text export produces valid output
- [ ] Markdown export produces valid output
- [ ] YAML export produces valid output
- [ ] Files saved to correct location
- [ ] Clipboard copy works
- [ ] Error handling is robust

✅ **Quality Requirements**
- [ ] Build succeeds with no errors
- [ ] All tests pass (or updated appropriately)
- [ ] No race conditions detected
- [ ] Code follows project conventions
- [ ] User feedback is clear and helpful

✅ **User Experience**
- [ ] Clear navigation between states
- [ ] Helpful error messages
- [ ] Progress feedback during export
- [ ] Success confirmation
- [ ] Easy to use and understand

---

## Timeline

**Estimated Duration**: 2-3 hours

1. **Planning & Analysis** (15 min) ✅ In Progress
2. **Update Data Structures** (30 min)
3. **Implement State Handlers** (60 min)
4. **Implement View Methods** (30 min)
5. **Implement Export Logic** (30 min)
6. **Testing** (30 min)
7. **Documentation** (15 min)

---

## Dependencies

- ExportService (already implemented)
- CVView and CVSection models (already available)
- Bubble Tea for UI (already integrated)
- File system access (standard library)
- Clipboard support (atotto/clipboard)

---

## Next Steps

1. Update GenerateCVContext with export fields
2. Add new states and messages to generate_cv.go
3. Implement state handlers in generate_cv_intent.go
4. Implement view methods
5. Test complete workflow
6. Update documentation

---

## Notes

- Export states are separate from main CV generation states to keep code organized
- ExportService is already fully implemented - we just need to integrate it
- File naming should be user-friendly and avoid special characters
- Error handling should gracefully fallback to clipboard if file save fails
- All operations should be async to keep UI responsive

---

*This plan provides a complete roadmap for implementing CV export functionality, completing the CV generation workflow.*

