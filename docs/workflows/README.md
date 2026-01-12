# KaRiya Workflow Documentation

**Comprehensive guides for KaRiya's complex user journeys**

This directory contains detailed workflow guides for the two most complex workflows in KaRiya's TUI. Each guide includes state machines, keyboard shortcuts, navigation patterns, and troubleshooting.

---

## Available Workflows

### 1. [CV Generation Workflow](CV_GENERATION_WORKFLOW.md)

**Complexity**: ⭐⭐⭐⭐⭐ High (10 states)  
**Purpose**: Generate role and audience-specific CVs from career events  
**Implementation**: `internal/cli/intents/generate_cv_intent.go`  
**Test Coverage**: 41 tests, 100% passing

**Key Features**:
- Profile and audience selection
- Real-time CV generation with progress indicator
- Interactive scrollable preview with metadata
- Export to multiple formats (Text, Markdown, YAML)
- Full back navigation with context preservation

**When to use**: Creating tailored CVs for job applications, generating role-specific resumes

**Typical Duration**: 30-60 seconds (depending on export)

**Quick Start**:
```
Main Menu → Generate CV → Select Profile → Select Audience →
  Wait for Generation → Preview → Confirm → [Optional: Export] → Complete
```

---

### 2. [Event Capture Workflow](EVENT_CAPTURE_WORKFLOW.md)

**Complexity**: ⭐⭐⭐⭐ High (4 states + 3 modal sub-flows)  
**Purpose**: Capture career events with optional burst/fact extraction  
**Implementation**: `internal/cli/intents/capture_event_intent.go`  
**Test Coverage**: 30+ tests, 100% passing

**Key Features**:
- Two capture modes: Quick (minimal) and Manual (all fields)
- Huh forms integration with real-time validation
- Natural language date parsing ("today", "7 days ago")
- Optional metadata, burst, and fact editing via modals
- Async submission with error recovery

**When to use**: Adding new career events to your journal, backfilling past experiences

**Typical Duration**: 20-60 seconds (depending on detail level)

**Quick Start**:
```
Main Menu → Capture Event → Choose Strategy (Quick/Manual) →
  Fill Form → [Optional: Review] → Submit → Success
```

---

## Workflow Diagrams

All workflow diagrams are **generated programmatically** from the actual implementation to ensure accuracy. Diagrams are created using Mermaid syntax and can be viewed on GitHub, VS Code, or online.

### Generating Diagrams

```bash
# Generate all workflow diagrams
make generate-diagrams

# Or use the script directly
./scripts/generate_workflow_diagrams.sh
```

### Diagram Files

| Workflow | Diagram File | States | Format |
|----------|--------------|--------|--------|
| CV Generation | [cv_generation_flow.mermaid](diagrams/cv_generation_flow.mermaid) | 10 | Mermaid |
| Event Capture | [event_capture_flow.mermaid](diagrams/event_capture_flow.mermaid) | 4 + 3 modals | Mermaid |

### Viewing Diagrams

**On GitHub**: Mermaid diagrams render automatically in GitHub markdown files

**In VS Code**: Install the "Markdown Preview Mermaid Support" extension

**Online**: Copy diagram content and paste at https://mermaid.live

**In Documentation**: Diagrams are embedded in workflow guides

---

## Navigation Patterns

All KaRiya workflows follow consistent navigation patterns for predictability and ease of use.

### Universal Shortcuts

These shortcuts work in **every workflow, every state**:

| Shortcut | Action | Description |
|----------|--------|-------------|
| **Esc** | Go back / Cancel | Navigate to previous state OR cancel if root state |
| **m** | Main menu | Return to main menu from anywhere |
| **q** | Quit | Exit application |
| **Ctrl+C** | Force quit | Interrupt and quit immediately |
| **?** | Help | Toggle help modal (future) |

### State Types

Every workflow state falls into one of these types:

**1. Root State** (First state in workflow)
- **Esc behavior**: Cancels entire workflow, returns to main menu
- **Examples**: CaptureEvent.ChooseStrategy, GenerateCV.SelectProfile

**2. Intermediate State** (Has previous state)
- **Esc behavior**: Navigate back to previous state
- **Examples**: CaptureEvent.Form, GenerateCV.Preview

**3. Async Operation State** (Background processing)
- **Esc behavior**: Let operation complete in background, navigate back
- **Examples**: GenerateCV.Generating, CaptureEvent.Submit

**4. Final State** (Workflow complete or error)
- **Esc behavior**: Retry or cancel
- **Examples**: GenerateCV.ExportComplete, CaptureEvent.Submit (error)

### Error Handling Philosophy

**Errors are preserved on back navigation**:
- User can see what went wrong
- Error context maintained for troubleshooting
- Retry options always available

**Example**:
```
Submit (error) → Esc → Review (error still visible) → Fix issue → Submit again
```

---

## Keyboard Shortcuts

For complete keyboard reference, see [Keyboard Shortcuts Guide](../KEYBOARD_SHORTCUTS_GUIDE.md).

### Common Patterns

**Navigation** (Lists and Menus):
- `↑/k` and `↓/j`: Navigate up/down (arrow or vim keys)
- `PgUp/PgDn`: Page navigation
- `Home/End` or `g/G`: Jump to first/last item
- `Enter`: Select/Confirm

**Forms** (Huh-based):
- `Tab` / `Shift+Tab`: Navigate fields
- `Enter`: Submit (via button focus)
- `Ctrl+S`: Quick submit
- `Ctrl+O`: Toggle optional fields (Manual mode)
- `Space`: Toggle checkboxes

**Scrolling** (Preview screens):
- `↑/k` and `↓/j`: Scroll line by line
- `PgUp/PgDn`: Scroll page by page
- `Home/End`: Jump to top/bottom

---

## Common Workflows

### End-to-End: Capture Event → Generate CV

**Goal**: Capture a new achievement and immediately use it in a CV

**Steps**:
1. Capture event (20-40s)
   - Main menu → `c`
   - Quick capture → fill text → submit

2. Generate CV (30-45s)
   - Main menu → `g`
   - Select profile → select audience → generate → preview → confirm

**Total time**: ~60-90 seconds

---

### Iterative CV Refinement

**Goal**: Generate CV, review, adjust audience, regenerate

**Steps**:
1. Generate with Hiring Manager audience
2. Review preview → not satisfied
3. Press `Esc` to go back to audience selection
4. Change to Peer audience
5. Review new preview → better
6. Confirm or export

**Benefit**: Compare how different audiences filter your events

---

## Troubleshooting

### General Workflow Issues

**Stuck in a state**:
- Press `m` to return to main menu (always works)
- Press `q` to quit application (last resort)
- Check footer for available shortcuts

**Can't go back**:
- In async states (Generating, Submitting), wait for operation to complete
- Press `m` instead of `Esc` for immediate main menu return

**Keyboard shortcuts not working**:
- Ensure terminal window has focus
- Check if another application is intercepting keys
- Restart application: `q` then relaunch

### Workflow-Specific Troubleshooting

For detailed troubleshooting:
- **CV Generation**: See [CV Generation Workflow - Troubleshooting](CV_GENERATION_WORKFLOW.md#troubleshooting)
- **Event Capture**: See [Event Capture Workflow - Troubleshooting](EVENT_CAPTURE_WORKFLOW.md#troubleshooting)

---

## Related Documentation

### User Guides
- **[Keyboard Shortcuts Guide](../KEYBOARD_SHORTCUTS_GUIDE.md)** - Complete keyboard reference for all workflows
- **[CLI Guide](../CLI_GUIDE.md)** - Command-line usage and options
- **[Burst & Fact Extraction Guide](../BURST_FACT_EXTRACTION_GUIDE.md)** - Understanding bursts and facts
- **[CV Generation Guide](../guides/CV_GENERATION_GUIDE.md)** - User guide for creating CVs

### Developer Guides
- **[TUI Standards](../TUI_STANDARDS.md)** - TUI design standards and guidelines
- **[TUI Developer Guide](../TUI_DEVELOPER_GUIDE.md)** - Creating new TUI components
- **[Keyboard System Guide](../development/KEYBOARD_SYSTEM_GUIDE.md)** - Implementing keyboard shortcuts
- **[Intent Development Checklist](../INTENT_DEVELOPMENT_CHECKLIST.md)** - Creating new workflows

### Technical Documentation
- **[TUI Intent Diagram](../TUI_INTENT_DIAGRAM.md)** - Complete intent architecture specification
- **[StandardView Guide](../STANDARDVIEW_GUIDE.md)** - StandardView system documentation
- **[Performance Benchmarks](../PERFORMANCE_BENCHMARKS.md)** - Performance targets and metrics

---

## Contributing

When adding new workflow guides:

1. **Create the guide**: Use `CV_GENERATION_WORKFLOW.md` as a template
2. **Update diagram script**: Add new diagram to `scripts/generate_workflow_diagrams.sh`
3. **Run diagram generation**: `make generate-diagrams`
4. **Add to this README**: Include in "Available Workflows" section
5. **Cross-reference**: Link from AGENTS.md and relevant user guides
6. **Test navigation**: Verify all keyboard shortcuts work as documented

**Template structure**:
- Overview (purpose, when to use, prerequisites)
- Workflow states & navigation (diagram + state table)
- Step-by-step guide (each state detailed)
- Complete keyboard reference (comprehensive table)
- Navigation patterns (forward, back, error recovery)
- Common workflows (timed examples)
- Troubleshooting (specific issues and solutions)
- Technical details (implementation notes)

---

## Workflow Statistics

| Workflow | States | Modals | Shortcuts | Tests | Avg Duration |
|----------|--------|--------|-----------|-------|--------------|
| CV Generation | 10 | 0 | 25+ | 41 | 30-60s |
| Event Capture | 4 | 3 | 20+ | 30+ | 20-60s |
| **Total** | **14** | **3** | **45+** | **71+** | - |

---

**Document Status**: ✅ Complete and Production-Ready  
**Last Updated**: 2026-01-12  
**Workflows Documented**: 2/5 core intents (40%)  
**Next Workflows**: BrowseTimeline, ExportArtifact, ConfigureSystem
