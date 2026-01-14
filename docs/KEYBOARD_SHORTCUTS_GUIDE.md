# KaRiya Keyboard Shortcuts Guide

**Complete Keyboard Reference for KaRiya TUI**

**Last Updated**: 2026-01-12
**Version**: 2.0 (Consolidated from KEYBOARD_REFERENCE.md and shortcuts/README.md)
**Audience**: End users, new team members

---

## Table of Contents

1. [Quick Reference Card](#quick-reference-card)
2. [Global Shortcuts](#global-shortcuts)
3. [Workflow-Specific Shortcuts](#workflow-specific-shortcuts)
4. [Vim-Style Navigation](#vim-style-navigation)
5. [Common Key Combinations](#common-key-combinations)
6. [Screen-Specific Examples](#screen-specific-examples)
7. [Accessibility Features](#accessibility-features)
8. [Troubleshooting](#troubleshooting)

---

## Quick Reference Card

**Print this section and keep it handy!**

```
╔════════════════════════════════════════════════════════════════════╗
║                    KaRiya Keyboard Reference                      ║
├────────────────────────────────────────────────────────────────────┤
║ GLOBAL (Work Everywhere)                                          ║
║ ─────────────────────────────────────────────────────────────────  ║
║ Esc    Go back / Cancel      q       Quit application             ║
║ m      Main menu             Ctrl+C  Quit (force)                 ║
║ ?      Help                                                        ║
║                                                                    ║
║ NAVIGATION                                                         ║
║ ─────────────────────────────────────────────────────────────────  ║
║ ↑/k    Previous item         ↓/j     Next item                    ║
║ ←/h    Move left             →/l     Move right                   ║
║ PgUp   Page up               PgDn    Page down                    ║
║ Home   First item            End     Last item                    ║
║ g      Jump to top           G       Jump to bottom               ║
║                                                                    ║
║ ACTIONS                                                            ║
║ ─────────────────────────────────────────────────────────────────  ║
║ Enter  Select / Confirm      Space   Toggle selection             ║
║ e      Edit                  d       Delete                       ║
║ c      Capture event         l       List events                  ║
║ f      Filter                s       Sort                         ║
║ /      Search                b       Bulk operations              ║
║                                                                    ║
║ FORMS                                                              ║
║ ─────────────────────────────────────────────────────────────────  ║
║ Tab         Next field       Shift+Tab   Previous field           ║
║ Ctrl+O      Toggle optional fields (Manual mode)                  ║
║ Ctrl+S      Submit form                                           ║
╚════════════════════════════════════════════════════════════════════╝
```

---

## Global Shortcuts

These shortcuts work everywhere in KaRiya, regardless of which screen you're on:

| Key | Action | Description |
|-----|--------|-------------|
| **`q`** | Quit application | Exit KaRiya completely |
| **`Ctrl+C`** | Force quit | Interrupt and quit immediately |
| **`?`** | Show help | Display context-sensitive help |
| **`m`** | Main menu | Return to main menu from any screen |
| **`Esc`** | Go back / Cancel | Navigate to previous screen or cancel current operation |

### Global Navigation

| Key | Action | Description |
|-----|--------|-------------|
| **`↑`** or **`k`** | Move up | Navigate to previous item in lists or menus |
| **`↓`** or **`j`** | Move down | Navigate to next item in lists or menus |
| **`←`** or **`h`** | Move left | Navigate left in menus or selectors |
| **`→`** or **`l`** | Move right | Navigate right in menus or selectors |
| **`PgUp`** | Page up | Scroll up one page in lists |
| **`PgDn`** | Page down | Scroll down one page in lists |
| **`Home`** | First item | Jump to first item in list |
| **`End`** | Last item | Jump to last item in list |
| **`g`** | Jump to top | Navigate to top of list (vim-style) |
| **`G`** | Jump to bottom | Navigate to bottom of list (vim-style) |

---

## Workflow-Specific Shortcuts

### Event Capture Workflow

**See**: [Event Capture Workflow Guide](workflows/EVENT_CAPTURE_WORKFLOW.md) for complete workflow documentation

| Screen | Key | Action | Description |
|--------|-----|--------|-------------|
| **Choose Strategy** | `↑/↓` or `k/j` | Navigate | Select Quick or Manual mode |
| | `Enter` | Confirm | Proceed to event form |
| | `Esc` | Cancel | Return to main menu |
| **Event Form** | `Tab` | Next field | Move to next form field |
| | `Shift+Tab` | Previous field | Move to previous field |
| | `Ctrl+O` | Toggle fields | Show/hide optional fields (Manual mode only) |
| | `Ctrl+S` | Submit | Save event (skip Review) |
| | `Enter` | Continue | Proceed to Review screen |
| | `Esc` | Back | Return to strategy selection |
| **Review** | `Ctrl+S` or `Enter` | Confirm | Submit event |
| | `e` | Edit metadata | Open metadata editor modal |
| | `b` | Edit bursts | Open burst suggestion modal |
| | `f` | Edit facts | Open fact editor modal |
| | `a` | Accept | Accept selected burst/fact |
| | `r` | Reject | Reject selected burst/fact |
| | `↑/↓` or `k/j` | Navigate | Move through bursts/facts |
| | `Esc` | Back | Return to form |
| **Submit** | `Enter` | Continue | Proceed after successful submission |
| | `r` | Retry | Retry failed submission |
| | `Esc` | Back | Return to Review (with error visible) |

### Browse Timeline Workflow

**See**: [Browse Timeline Guide](guides/BROWSE_TIMELINE_GUIDE.md) for complete documentation

| Screen | Key | Action | Description |
|--------|-----|--------|-------------|
| **Timeline List** | `↑/↓` or `k/j` | Navigate | Move through event list |
| | `PgUp/PgDn` | Page | Scroll by page |
| | `Home/End` | Jump | First/last event |
| | `g/G` | Jump | Top/bottom of list |
| | `Enter` | View details | Open event detail screen |
| | `f` | Filter | Show filter options |
| | `s` | Sort | Show sort options |
| | `/` | Search | Activate search mode |
| | `Esc` | Cancel | Return to main menu |
| **Event Detail** | `Enter` | Select | Confirm selection |
| | `e` | Edit | Edit event metadata |
| | `d` | Delete | Delete event (with confirmation) |
| | `Esc` | Back | Return to timeline list |

### Generate CV Workflow

**See**: [CV Generation Workflow Guide](workflows/CV_GENERATION_WORKFLOW.md) for complete workflow documentation

| Screen | Key | Action | Description |
|--------|-----|--------|-------------|
| **Select Profile** | `↑/↓` or `k/j` | Navigate | Move through profiles |
| | `Enter` | Select | Choose profile and proceed |
| | `Esc` | Cancel | Return to main menu |
| **Select Audience** | `↑/↓` or `k/j` | Navigate | Move through audience options |
| | `Enter` | Generate | Start CV generation |
| | `Esc` | Back | Return to profile selection |
| **Generating** | `Esc` | Background | Let generation complete in background |
| **Preview** | `↑/↓` or `k/j` | Scroll | Scroll through CV preview |
| | `PgUp/PgDn` | Page scroll | Scroll by page |
| | `e` or `Enter` | Edit | Open CV review/edit |
| | `c` | Continue | Proceed to confirmation |
| | `Esc` | Back | Return to audience selection |
| **Review** | `Enter` | Continue | Proceed to confirmation |
| | `Esc` | Back | Return to preview |
| **Confirm** | `y` or `Enter` | Complete | Finish without exporting |
| | `e` or `x` | Export | Proceed to export workflow |
| | `n` or `Esc` | Back | Return to review |
| **Export Format** | `↑/↓` or `k/j` | Navigate | Select format (Text, Markdown, YAML) |
| | `Enter` | Select | Choose format and proceed |
| | `Esc` | Back | Return to confirmation |
| **Export Location** | `↑/↓` or `k/j` | Navigate | Select location (File, Clipboard, Cancel) |
| | `Enter` | Export | Start export process |
| | `Esc` | Back | Return to format selection |
| **Exporting** | `Esc` | Background | Let export complete in background |
| **Export Complete** | `Enter` | Done | Complete workflow |
| | `Esc` | Retry | Return to location selection |

### Export Artifact Workflow

| Screen | Key | Action | Description |
|--------|-----|--------|-------------|
| **Select Type** | `↑/↓` or `k/j` | Navigate | Choose artifact type (events, facts, bursts) |
| | `Enter` | Select | Proceed to format selection |
| | `Esc` | Cancel | Return to main menu |
| **Select Format** | `↑/↓` or `k/j` | Navigate | Choose format (JSON, CSV, YAML, TXT) |
| | `Enter` | Select | Proceed to destination |
| | `Esc` | Back | Return to type selection |
| **Select Destination** | `↑/↓` or `k/j` | Navigate | Choose destination (File, Clipboard) |
| | `Enter` | Select | Proceed to configuration |
| | `Esc` | Back | Return to format selection |
| **Configure** | `Enter` | Continue | Proceed to preview |
| | `Esc` | Back | Return to destination |
| **Preview** | `↑/↓` or `k/j` | Scroll | Scroll through export preview |
| | `PgUp/PgDn` | Page scroll | Scroll by page |
| | `Enter` | Continue | Proceed to confirmation |
| | `Esc` | Back | Return to configuration |
| **Confirm** | `y` or `Enter` | Confirm | Start export |
| | `n` or `Esc` | Cancel | Return to preview |
| **In Progress** | `Esc` | Background | Let export complete |
| **Complete/Failed** | `Enter` | Done | Finish workflow |
| | `r` | Retry | Retry export (if failed) |
| | `Esc` | Back | Return to confirmation |

### Configure System Workflow

| Screen | Key | Action | Description |
|--------|-----|--------|-------------|
| **Select Domain** | `↑/↓` or `k/j` | Navigate | Choose domain (System, Profile, Export, UI) |
| | `Enter` | Select | Open settings editor |
| | `Esc` | Cancel | Return to main menu |
| **Edit Settings** | `↑/↓` or `k/j` | Navigate | Move through settings (when not editing) |
| | `Enter` | Edit/Save | Start editing OR save edited value |
| | `Esc` | Cancel | Cancel editing OR return to domain selection |
| | `Ctrl+S` | Save all | Save all changes and proceed to review |
| | `Type` | Edit | Edit setting value (when in edit mode) |
| **Review Changes** | `Enter` | Continue | Proceed to confirmation |
| | `Esc` | Back | Return to settings editor |
| **Confirm** | `y` or `Enter` | Confirm | Save configuration changes |
| | `n` or `Esc` | Cancel | Return to review |
| **Saving** | `Esc` | Background | Let save complete |
| **Complete/Failed** | `Enter` | Done | Finish workflow |
| | `r` | Retry | Retry save (if failed) |
| | `Esc` | Back | Return to confirmation |

---

## Vim-Style Navigation

KaRiya supports vim-style navigation for power users who prefer keeping hands on the home row:

| Vim Key | Alternative | Action | Context |
|---------|-------------|--------|---------|
| **`k`** | `↑` | Move up | Lists, forms, menus |
| **`j`** | `↓` | Move down | Lists, forms, menus |
| **`h`** | `←` | Move left | Menus, selectors |
| **`l`** | `→` | Move right | Menus, selectors |
| **`g`** | `Home` | Jump to top | Lists |
| **`G`** | `End` | Jump to bottom | Lists |

**Benefits**:
- ✅ Hands stay on keyboard home row
- ✅ No arrow key movement needed
- ✅ Compatible with vim muscle memory
- ✅ Optional (arrow keys still work)

**Note**: Both vim keys and arrow keys work simultaneously - use whichever you prefer!

---

## Common Key Combinations

### Navigating Through Workflows

**Home to Event Capture**:
```
c → Choose Strategy → Enter → Fill Form → Tab (navigate) → Enter (submit)
```

**Home to Browse Timeline**:
```
l → Timeline List → ↑/↓ (navigate) → Enter (view details)
```

**Viewing Event Details**:
```
Timeline List → ↓ (select) → Enter (view) → e (edit) OR d (delete) OR Esc (back)
```

**Generating a CV**:
```
Generate CV → Select Profile → Enter → Select Audience → Enter → 
    Wait for Generation → Preview → c (confirm) OR e/x (export)
```

**Exporting a CV**:
```
From CV Confirm → e/x → Select Format → Enter → Select Location → Enter → 
    Wait for Export → Enter (done)
```

### Form Completion Workflow

**Quick Capture** (minimal fields):
```
1. c                  - Start capture from home
2. Enter              - Select Quick mode
3. Type event text    - Fill event description
4. Tab                - Move to date field
5. Type date          - Enter date (e.g., "today", "7 days ago")
6. Ctrl+S OR Enter    - Submit immediately
```

**Manual Capture** (all fields):
```
1. c                  - Start capture from home
2. ↓ → Enter          - Select Manual mode
3. Type event text    - Fill event description
4. Tab                - Move to date
5. Tab                - Move to company
6. Tab                - Move to project
7. Ctrl+O             - Toggle optional fields if needed
8. Tab                - Move to tags
9. Space              - Select tags
10. Enter             - Proceed to Review
11. Ctrl+S            - Submit
```

### Bulk Operations

**Selecting and Tagging Multiple Events**:
```
1. l                  - List events from home
2. b                  - Enter bulk operations mode
3. Space              - Select/deselect current event
4. ↓                  - Move to next event
5. Space              - Select/deselect
6. a                  - Select all (optional)
7. Enter              - Apply bulk operation
```

---

## Screen-Specific Examples

### Home Screen

```
┌──────────────────────────────────────────────────────────────┐
│ KaRiya - Career Journal                                      │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ►  Capture Event (c)        Add a new career event         │
│     Browse Timeline (l)      View your career history       │
│     Generate CV (g)          Create role-specific CV        │
│     Export Artifact (e)      Export data to files           │
│     Configure System (s)     System settings                │
│     Help (?)                 View help documentation        │
│     Quit (q)                 Exit the application           │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ ↑/k ↓/j Navigate  Enter Select  ? Help  q Quit              │
└──────────────────────────────────────────────────────────────┘
```

**Available shortcuts**: `↑/↓` or `k/j` to navigate, `Enter` to select, `c/l/g/e/s/?/q` for direct access

### Event Form Screen

```
┌──────────────────────────────────────────────────────────────┐
│ Capture Career Event                                         │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  Event Text *                                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Led team to migrate platform to Ruby 3.1...           │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  Date *                                                      │
│  ┌────────────────────┐                                      │
│  │ 2024-01-07         │  (e.g., "today", "7 days ago")     │
│  └────────────────────┘                                      │
│                                                              │
│  Company               Project                               │
│  ┌──────────────┐     ┌──────────────┐                      │
│  │ Acme Corp    │     │ Migration    │                      │
│  └──────────────┘     └──────────────┘                      │
│                                                              │
│  Tags                  Categories                            │
│  ☑ Technical          ☑ Achievement                         │
│  ☐ Leadership         ☐ Project                             │
│  ☐ Mentoring          ☐ Research                            │
│                                                              │
│  [ Submit ]  [ Cancel ]                                      │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ Tab Next  Shift+Tab Prev  Space Toggle  Ctrl+O Toggle       │
│ Ctrl+S Submit  Esc Cancel  ? Help                           │
└──────────────────────────────────────────────────────────────┘
```

**Available shortcuts**: `Tab/Shift+Tab` for field navigation, `Space` to toggle checkboxes, `Ctrl+O` to show/hide optional fields (Manual mode), `Ctrl+S` or `Enter` to submit

### Timeline List Screen

```
┌──────────────────────────────────────────────────────────────┐
│ Browse Career Timeline (42 events)                           │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ▶  2024-01-07  Led team to migrate platform to Ruby 3.1    │
│     2024-01-05  Developed new API for customer dashboard     │
│     2024-01-03  Mentored junior developers on best practices │
│     2024-01-01  Architected microservices infrastructure     │
│     2023-12-28  Optimized database queries (50% improvement) │
│     ...                                                      │
│                                                              │
│  [1-10 of 42]  Page 1/5                                      │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ ↑/k ↓/j Navigate  PgUp/PgDn Page  Enter View  f Filter      │
│ s Sort  / Search  Esc Back  ? Help                          │
└──────────────────────────────────────────────────────────────┘
```

**Available shortcuts**: `↑/↓` or `k/j` to navigate, `Enter` to view details, `f` for filters, `s` for sort, `/` for search, `Esc` to go back

---

## Accessibility Features

### Keyboard-Only Navigation

✅ **All features accessible without mouse**
- Every function can be performed using keyboard shortcuts
- Clear tab order through forms and menus
- Visual focus indicators show current selection

✅ **Clear keyboard flow on each screen**
- Consistent shortcuts across all screens
- Context-aware help displayed at bottom of every screen
- Logical tab order matches visual layout

✅ **Visual focus indicators**
- `▶` marker shows selected item in lists
- Thick teal border shows focused form field
- Bold text with border shows focused buttons

✅ **Progressive disclosure**
- Tab order follows natural reading flow (top to bottom, left to right)
- Optional fields can be hidden/shown with `Ctrl+O`
- Help text always available with `?` key

### Help and Discoverability

✅ **Help text shown on every screen**
- Footer displays relevant shortcuts for current screen
- `?` key opens detailed help at any time
- Help text adapts to current context

✅ **Keyboard shortcuts listed in help footer**
- Global shortcuts always visible
- Context-specific shortcuts shown when relevant
- Shortcut descriptions explain what each key does

✅ **Context-sensitive shortcuts displayed**
- Form screens show form navigation shortcuts
- List screens show list navigation shortcuts
- Different states within workflows show appropriate shortcuts

✅ **Clear instructions for each field**
- Form fields show expected format (e.g., "today", "7 days ago")
- Validation errors display clearly with specific guidance
- Character counters show limits for text fields

### Input Methods

✅ **Text input with character counter**
- Shows remaining characters for limited fields
- Real-time validation feedback
- Clear error messages

✅ **Date parsing supports multiple formats**
- Natural language: "today", "yesterday", "7 days ago"
- ISO format: "2024-01-07"
- Relative dates: "2 weeks ago", "1 month ago"

✅ **Dropdown selections with arrow keys**
- `↑/↓` or `k/j` to navigate options
- `Enter` to select
- Visual indicator shows current selection

✅ **Checkbox toggles with space bar**
- `Space` to toggle on/off
- Visual checkmark (☑/☐) shows state
- Works with `↑/↓` navigation to select multiple

---

## Troubleshooting

### Key Not Working

**Problem**: A shortcut key doesn't respond

**Solutions**:
1. **Verify the context** - The shortcut may only work on certain screens
   - Press `?` to see available shortcuts for current screen
   - Check if you're in the correct workflow state

2. **Check terminal support** - Some terminals don't support all keys
   - Try the alternative key (e.g., `↑` instead of `k`)
   - Test with arrow keys if vim keys don't work

3. **Ensure no conflicts** - Another application may intercept the key
   - Close other applications that use global shortcuts
   - Check if terminal emulator has conflicting keybindings

4. **Restart if needed** - Temporary glitch may be resolved by restart
   - Press `q` to quit, then restart KaRiya
   - Check terminal settings if problem persists

### Navigation Issues

**Problem**: Can't navigate to expected screen or option

**Solutions**:
1. **Press `Esc` to return to safe state**
   - `Esc` always goes back one step
   - `m` returns to main menu from anywhere

2. **Use `?` to see available shortcuts**
   - Shows what you can do from current screen
   - Displays context-specific navigation options

3. **Quit and restart if stuck**
   - Press `q` to quit cleanly
   - Restart KaRiya and try again

### Form Stuck or Unresponsive

**Problem**: Can't submit form or navigate away

**Solutions**:
1. **Press `Esc` to cancel form entry**
   - All unsaved data will be discarded
   - Returns to previous screen safely

2. **Ensure required fields are filled**
   - Forms won't submit if required fields are empty
   - Look for error messages in red
   - Required fields marked with `*`

3. **Try `Ctrl+S` to force submit**
   - Bypasses normal validation for Quick mode
   - May show validation errors if fields incomplete

4. **Return to home and try again**
   - Press `m` to return to main menu
   - Start workflow over from beginning

### Forgotten Shortcuts

**Problem**: Can't remember a specific shortcut

**Solutions**:
1. **Press `?` at any time to open help**
   - Shows shortcuts for current screen
   - Organized by category

2. **Check footer for common shortcuts**
   - Bottom of every screen shows relevant keys
   - Global shortcuts always displayed

3. **Refer to this guide**
   - Print the Quick Reference Card (top of this document)
   - Keep it handy for quick lookups

4. **Use alternative keys**
   - Most shortcuts have alternatives (arrows vs vim keys)
   - Try common keys like `Enter`, `Esc`, `Tab`

---

## Related Documentation

### User Documentation
- **[CV Generation Workflow Guide](workflows/CV_GENERATION_WORKFLOW.md)** - Complete CV generation workflow with state diagrams
- **[Event Capture Workflow Guide](workflows/EVENT_CAPTURE_WORKFLOW.md)** - Event capture workflow with burst/fact extraction
- **[TUI Standards](TUI_STANDARDS.md)** - Complete TUI design standards and guidelines
- **[CLI Guide](CLI_GUIDE.md)** - Command-line usage and options
- **[Troubleshooting Guide](TROUBLESHOOTING.md)** - General troubleshooting for KaRiya

### Developer Documentation
- **[State Matrix](STATE_MATRIX.md)** - Complete state matrix (10 intents, 64 states, escape behavior)
- **[Centralized Key Handling Guide](development/CENTRALIZED_KEY_HANDLING.md)** - Developer guide for key handling
- **[Keyboard System Guide](development/KEYBOARD_SYSTEM_GUIDE.md)** - Comprehensive keyboard system implementation

---

## Support

For keyboard shortcut issues or feature requests:

- **GitHub Issues**: https://github.com/baphled/kariya/issues
- **Documentation**: https://github.com/baphled/kariya/tree/main/docs
- **Include**: Operating system, terminal emulator, keyboard layout, steps to reproduce

---

**Document Status**: ✅ Complete and Production-Ready
**Last Updated**: 2026-01-12
**Version**: 2.0 (Consolidated)
**Next Review**: 2026-04-01
