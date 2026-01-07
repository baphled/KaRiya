# KaRiya Keyboard Reference

**Quick Reference Guide for All Keyboard Shortcuts**

---

## Global Shortcuts (Available Everywhere)

| Key | Action | Notes |
|-----|--------|-------|
| `q` | Quit application | Works from any screen |
| `Ctrl+C` | Quit application | Standard interrupt |
| `?` | Show help | Context-sensitive help |
| `Esc` | Go back / Cancel | Primary back navigation key |

---

## Home Screen (Main Menu)

| Key | Action | Notes |
|-----|--------|-------|
| `c` | Capture new event | Opens event capture form |
| `l` | List events | View all captured events |
| `m` | Metadata review | Review and enrich event metadata |
| `?` | Show help | Display help information |

---

## Form Navigation (Event Capture & Editing)

| Key | Action | Context | Notes |
|-----|--------|---------|-------|
| `Tab` | Next field | Forms | Move to next input field |
| `Shift+Tab` | Previous field | Forms | Move to previous field |
| `↑` or `k` | Up / Previous | Lists in forms | Navigate tag/category lists |
| `↓` or `j` | Down / Next | Lists in forms | Navigate tag/category lists |
| `←` or `h` | Left / Previous | Multi-select | Move left in selections |
| `→` or `l` | Right / Next | Multi-select | Move right in selections |
| `Enter` | Confirm | Forms | Submit form or confirm selection |
| `Space` | Toggle | Checkboxes | Select/deselect tags and categories |
| `Esc` | Cancel | Forms | Cancel form entry and go back |
| `Ctrl+O` | Toggle fields | Forms | Show/hide optional fields (manual mode) |

**Form Fields Supported:**
- Text input (event description)
- Date input (event date)
- Optional fields (company, project)
- Tag selection (multiple)
- Category selection (multiple)
- Capture mode selection (timeline, backfill, manual)

---

## List View (Event Browsing)

| Key | Action | Context | Notes |
|-----|--------|---------|-------|
| `↑` or `k` | Move up | Lists | Select previous event |
| `↓` or `j` | Move down | Lists | Select next event |
| `Enter` | Select | Lists | View event details |
| `e` | Edit | List items | Edit selected event metadata |
| `d` | Delete | List items | Delete selected event |
| `f` | Filter | Lists | Show filter options |
| `s` | Sort | Lists | Show sort options |
| `/` | Search | Lists | Activate search mode |
| `b` | Bulk select | Lists | Enter bulk operations mode |
| `Esc` | Go back | Lists | Return to home screen |

**List Navigation Features:**
- Scrolling through events
- Single event selection
- Quick filtering and sorting
- Bulk operations on multiple events

---

## Metadata Review & Editing

| Key | Action | Context | Notes |
|-----|--------|---------|-------|
| `↑` or `k` | Previous event | Metadata review | Move to previous event in list |
| `↓` or `j` | Next event | Metadata review | Move to next event in list |
| `Enter` | Edit | Metadata review | Open metadata editor for event |
| `Tab` | Next field | Metadata editor | Move to next metadata field |
| `Shift+Tab` | Previous field | Metadata editor | Move to previous field |
| `Space` | Toggle | Checkboxes | Select/deselect in editor |
| `Enter` | Save | Metadata editor | Save changes and return to list |
| `Esc` | Cancel | Metadata editor | Discard changes and go back |

**Data Quality Scoring:**
- Events show quality score (0-100%)
- Missing fields indicate improvement areas
- Metadata enrichment improves event quality

---

## Bulk Operations

| Key | Action | Context | Notes |
|-----|--------|---------|-------|
| `Space` | Select / Deselect | Bulk ops | Toggle event selection |
| `Enter` | Execute | Bulk ops | Apply bulk operation |
| `↑` or `k` | Previous | Bulk ops | Move to previous event |
| `↓` or `j` | Next | Bulk ops | Move to next event |
| `Esc` | Cancel | Bulk ops | Cancel bulk operation |
| `a` | Select all | Bulk ops | Select all events in list |
| `n` | Deselect all | Bulk ops | Clear all selections |

**Available Bulk Operations:**
- Add tags to multiple events
- Add categories to multiple events
- Update company/project for events
- Delete multiple events

---

## Help System

| Key | Action | Context | Notes |
|-----|--------|---------|-------|
| `?` | Open help | Any screen | Show context-sensitive help |
| `↑` or `k` | Previous topic | Help view | Navigate help topics |
| `↓` or `j` | Next topic | Help view | Navigate help topics |
| `Enter` | Select | Help menus | Select topic or action |
| `Esc` | Close help | Help view | Return to previous screen |

**Help Topics:**
- Keyboard shortcuts reference
- Feature explanations
- Data input guidelines
- Troubleshooting tips

---

## Vim-Style Navigation

KaRiya supports vim-style navigation for power users:

| Key | Alternative | Action |
|-----|-------------|--------|
| `k` | `↑` | Move up |
| `j` | `↓` | Move down |
| `h` | `←` | Move left |
| `l` | `→` | Move right |
| `Esc` | `Backspace` | Go back / Cancel |

**Vim Navigation Benefits:**
- Hands stay on keyboard home row
- No arrow key movement needed
- Compatible with vim muscle memory
- Optional (arrow keys still work)

---

## Screen-Specific Reference

### Home Screen
```
┌─────────────────────────────────────┐
│  KaRiya - Career Journal            │
├─────────────────────────────────────┤
│  [c] Capture event                  │
│  [l] List events                    │
│  [m] Metadata review                │
│  [?] Help                           │
│  [q] Quit                           │
├─────────────────────────────────────┤
│  c: Capture | l: List | m: Metadata │
└─────────────────────────────────────┘
```

### Form Screen
```
┌─────────────────────────────────────┐
│  Capture Event                      │
├─────────────────────────────────────┤
│  Event Text:        [____________]  │
│  Date:              [____________]  │
│  Company:           [____________]  │
│  Mode:              [Timeline    ▼] │
│  Tags:              [Technical    ] │
│                     [Leadership    ] │
│  [Submit]  [Cancel]                 │
├─────────────────────────────────────┤
│  Tab: Next | ↑↓: Navigate | Esc: Back│
└─────────────────────────────────────┘
```

### List Screen
```
┌─────────────────────────────────────┐
│  Events (42 total)                  │
├─────────────────────────────────────┤
│  ▶ Event 1: Led team to success    │
│    Event 2: Developed backend ...  │
│    Event 3: Mentored junior devs   │
│  [f] Filter | [s] Sort | [/] Search│
├─────────────────────────────────────┤
│  ↑↓: Navigate | e: Edit | d: Delete │
└─────────────────────────────────────┘
```

---

## Common Key Combinations

### Navigation Across Screens
1. **Home to Capture**: `c` → Capture form
2. **Home to List**: `l` → Event list
3. **List to Details**: `↓` `Enter` → Event details
4. **Details to List**: `Esc` → Back to list
5. **List to Home**: `Esc` → Back to home

### Form Completion
1. **Start form**: Press `c` from home
2. **Fill text**: Type event description
3. **Move to date**: `Tab` or click
4. **Navigate tags**: `T` then `↑↓` to select
5. **Submit**: `Tab` to submit button, `Enter`

### Bulk Operations
1. **Enter bulk mode**: `b` from list
2. **Select events**: `Space` to toggle
3. **Select all**: `a` shortcut
4. **Execute**: `Enter` to apply operation
5. **Exit**: `Esc` to cancel or return

---

## Accessibility Features

### Keyboard-Only Navigation
✅ All features accessible without mouse
✅ Clear keyboard flow on each screen
✅ Visual focus indicators
✅ Tab order matches screen layout

### Help and Discoverability
✅ Help text shown on every screen
✅ Keyboard shortcuts listed in help footer
✅ Context-sensitive shortcuts displayed
✅ Clear instructions for each field

### Input Methods
✅ Text input with character counter
✅ Date parsing supports multiple formats
✅ Dropdown selections with arrow keys
✅ Checkbox toggles with space bar

---

## Troubleshooting

### Key Not Working
- Verify terminal supports the key (some don't support Escape alternatives)
- Try the arrow key alternative if available
- Check if terminal has keyboard translation enabled

### Navigation Issues
- Press `Esc` to return to safe state
- Use `?` to see available shortcuts for current screen
- Quit with `q` and restart if needed

### Form Stuck
- Press `Esc` to cancel form entry
- All unsaved data will be discarded
- Return to home and try again

---

## Quick Key Reference Card

**Print this table and keep it handy:**

```
GLOBAL: q=quit, ?=help, Esc=back, Ctrl+C=quit
HOME:   c=capture, l=list, m=metadata
FORM:   Tab=next, ↑↓=select, Space=toggle, Esc=cancel
LIST:   ↑↓=move, e=edit, d=delete, f=filter, s=sort, /=search
METADATA: Enter=edit, Space=select, Esc=cancel
VIM:    k=up, j=down, h=left, l=right, Esc=back
```

---

**Last Updated**: 2025-12-30
**Version**: 1.0
**Status**: Complete and ready for reference

