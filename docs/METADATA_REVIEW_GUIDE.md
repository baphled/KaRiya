---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya Metadata Review Guide

## Overview

The Metadata Review feature enables you to review, validate, and enrich event metadata (dates, companies, projects, tags, categories) before automated processing. This guide covers how to use the metadata review system to improve your career event data quality.

## What is Metadata?

Metadata refers to the structured information about your career events:

- **Date**: When the event occurred (YYYY-MM-DD format)
- **Company**: Organization where the event took place
- **Project**: Specific project or initiative
- **Tags**: Flexible labels to categorize events (up to 8)
- **Categories**: Structured competency areas (up to 2)

## Data Quality Scoring

Events are automatically scored for data completeness on a 0-100 scale:

### Quality Levels

| Level | Score | Description | Example |
|-------|-------|-------------|---------|
| **Incomplete** | 0-25 | Only event description | "Implemented critical feature" |
| **Basic** | 26-50 | Description + date | "Implemented critical feature (2024-06-15)" |
| **Enriched** | 51-75 | Description + date + company/project | Above + "TechCorp / Platform" |
| **Complete** | 76-100 | All fields + tags/categories | Above + Tags & Categories |

### Quality Calculation

Points are awarded for:
- **Text**: 20 points (event description provided)
- **Date**: 20 points (event date specified)
- **Company**: 20 points (company name provided)
- **Project**: 20 points (project name provided)
- **Tags**: 15 points (tags assigned, up to 8)
- **Categories**: 15 points (category assigned, up to 2)
- **Match Quality**: 10 bonus points (high quality match)

### Visual Indicators

Events display color-coded quality bars:
- 🔴 Red: 0-25 (Incomplete)
- 🟡 Yellow: 26-50 (Basic)
- 🟢 Green: 51-75 (Enriched)
- 🟢 Bright Green: 76-100 (Complete)

## Accessing Metadata Review

### From Home Screen

Press `m` to open the metadata review screen directly.

### From Success Screen

After capturing an event:
1. Press the right arrow (`→`) to navigate to "Review Metadata"
2. Press `Enter` to open metadata review for the captured event

### From Import

After importing events from CSV:
- Automatically navigates to metadata review
- All imported events are pre-loaded
- Ready for immediate enrichment

### From List View

1. Open list view (press `l` from home)
2. Press `m` to open metadata review

## Metadata Review Screen

### Layout

```
┌─────────────────────────────────────────────────────────────┐
│ Metadata Review - 5 Events (3 selected)                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ [1] ◆ Led team through critical migration                  │
│     Date: 2024-06-15 | Company: TechCorp | Quality: Basic  │
│     Quality: ████░░░░░░ 45% | Missing: tags, categories   │
│                                                             │
│ [2] ◆ Implemented critical feature                         │
│     Date: 2024-06-10 | Company: TechCorp | Quality: Basic  │
│     Quality: ████░░░░░░ 45% | Missing: tags, categories   │
│                                                             │
│ [3] ◆ Architected new microservices platform              │
│     Date: 2024-05-20 | Company: TechCorp | Quality: Basic  │
│     Quality: ████░░░░░░ 45% | Missing: tags, categories   │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ ↑/↓ - Navigate | Enter - Edit | Space - Select | e - Bulk  │
│ f - Filter | s - Sort | Esc - Back                         │
└─────────────────────────────────────────────────────────────┘
```

### Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move to previous event |
| `↓` / `j` | Move to next event |
| `Enter` | Edit selected event |
| `Space` | Select/deselect event for bulk operations |
| `a` | Select all events |
| `d` | Deselect all events |
| `e` | Enter bulk operations mode |
| `f` | Filter by quality level |
| `s` | Change sort order |
| `Esc` / `q` | Return to previous screen |

## Individual Event Editing

### Opening the Editor

1. Navigate to an event in the metadata review screen
2. Press `Enter` to open the editor
3. The current metadata is loaded

### Editor Fields

```
┌─────────────────────────────────────────────────────────────┐
│ Edit Event: Led team through critical migration            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Date:       2024-06-15 (YYYY-MM-DD or relative)           │
│ Company:    TechCorp                                        │
│ Project:    Critical Migration                             │
│ Tags:       ☑ technical  ☐ leadership  ☐ mentoring  ...   │
│ Categories: ☑ technical  ☐ leadership  ☐ product  ...     │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ [Save]  [Cancel]  [Undo]                                   │
└─────────────────────────────────────────────────────────────┘
```

### Editing Each Field

#### Date Field
- Format: `YYYY-MM-DD` (e.g., 2024-06-15)
- Or use relative dates: `today`, `1 week ago`, `2 months ago`
- Validation: Cannot be in the future, must be valid date
- Leave blank to reset to capture date

#### Company Field
- Text input (max 200 characters)
- Optional
- Autocomplete from previous entries
- Whitespace automatically trimmed

#### Project Field
- Text input (max 200 characters)
- Optional
- Autocomplete from previous entries
- Whitespace automatically trimmed

#### Tags
- Multi-select from allowed tags
- Maximum 8 tags per event
- Case-insensitive matching
- Allowed tags: `project`, `achievement`, `leadership`, `technical`, `consulting`, `research`, `product`, `mentoring`
- Use `Space` to toggle, `Tab` to navigate

#### Categories
- Multi-select from allowed categories
- Maximum 2 categories per event
- Represents core competency areas
- Allowed categories: `technical`, `leadership`, `product`, `consulting`, `research`, `mentoring`
- Use `Space` to toggle, `Tab` to navigate

### Editor Controls

| Key | Action |
|-----|--------|
| `Tab` | Move to next field |
| `Shift+Tab` | Move to previous field |
| `Space` | Toggle checkbox (for tags/categories) |
| `Enter` | Save changes |
| `Esc` | Cancel without saving |
| `Ctrl+Z` | Undo changes (revert to original) |

### Saving Changes

1. Update desired fields
2. Press `Enter` to validate and save
3. If validation errors exist, they're shown in red
4. Fix errors and try again
5. On success, returns to metadata review screen

### Validation Rules

**Date**:
- Format: YYYY-MM-DD or recognized relative date
- Cannot be in the future
- Must be a valid calendar date

**Company**:
- Maximum 200 characters
- Optional field

**Project**:
- Maximum 200 characters
- Optional field

**Tags**:
- Must be from allowed list
- Maximum 8 total
- No duplicates
- Case-insensitive

**Categories**:
- Must be from allowed list
- Maximum 2 total
- No duplicates

## Bulk Operations

### When to Use Bulk Operations

Use bulk operations when:
- Multiple events need the same metadata update
- Adding company name to all events from a period
- Applying tags to a group of related events
- Setting categories for similar events

### Entering Bulk Operations Mode

1. At metadata review screen, select events with `Space`
2. Press `a` to select all events, or select specific ones
3. Press `e` to enter bulk operations mode
4. Or select events and press `e` directly

### Bulk Edit Interface

```
┌─────────────────────────────────────────────────────────────┐
│ Bulk Edit - 3 Events Selected                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Select fields to update:                                    │
│                                                             │
│ ☑ Company:     [TechCorp           ]                       │
│    ☑ Apply if empty                                         │
│                                                             │
│ ☑ Project:     [Critical Migration ]                       │
│    ☐ Apply if empty                                         │
│                                                             │
│ ☑ Tags:        ☑ technical  ☑ leadership  ☐ mentoring     │
│    ☑ Add if empty                                           │
│                                                             │
│ ☐ Categories:  ☐ technical  ☑ leadership  ☐ product       │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ [Preview]  [Cancel]                                         │
└─────────────────────────────────────────────────────────────┘
```

### Bulk Operation Modes

#### Apply to All
- Updates all selected events with the value
- Overwrites existing values

#### Apply if Empty
- Only updates events where the field is empty
- Preserves existing values
- Useful for adding missing metadata without overwriting

### Workflow

1. **Select Events**: Use `Space` to select multiple events
2. **Enter Bulk Edit**: Press `e`
3. **Choose Fields**: Toggle checkboxes for fields to update
4. **Set Values**: Enter values for selected fields
5. **Configure Options**: Choose "Apply if empty" for conditional updates
6. **Preview**: Press `Enter` to see changes
7. **Confirm**: Press `Enter` again to apply changes

### Preview Screen

```
┌─────────────────────────────────────────────────────────────┐
│ Preview Changes - 3 Events                                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Event 1: Led team through critical migration               │
│   Company: (empty) → TechCorp                              │
│   Tags: (none) → technical, leadership                     │
│                                                             │
│ Event 2: Implemented critical feature                      │
│   Company: (empty) → TechCorp                              │
│   Tags: (none) → technical, leadership                     │
│                                                             │
│ Event 3: Architected new microservices platform            │
│   Company: (empty) → TechCorp                              │
│   Tags: (none) → technical, leadership                     │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ [Confirm]  [Cancel]  [Undo]                                │
└─────────────────────────────────────────────────────────────┘
```

### Confirming Changes

1. Review the preview carefully
2. Press `Enter` to confirm and apply
3. Changes are applied to all selected events
4. Returns to metadata review screen
5. Quality scores are updated

## Filtering & Sorting

### Filtering by Quality Level

Press `f` at metadata review screen to filter:

- **All**: Show all events
- **Incomplete**: Quality 0-25 (missing most metadata)
- **Basic**: Quality 26-50 (has date but missing company/project)
- **Enriched**: Quality 51-75 (has basic metadata)
- **Complete**: Quality 76-100 (fully enriched)

**Use case**: Focus on "Incomplete" events to improve overall data quality.

### Sorting Options

Press `s` at metadata review screen to sort:

| Option | Order |
|--------|-------|
| **By Date** | Chronological (oldest first) |
| **By Company** | Alphabetical by company name |
| **By Creation** | Newest first (capture order) |

**Use cases**:
- Sort by company to bulk-edit events from same organization
- Sort by date to review career timeline
- Sort by creation to review recent captures

## Common Workflows

### Workflow 1: Quick Enrichment After Capture

1. Capture event (press `c`)
2. On success screen, press `→` to "Review Metadata"
3. Press `Enter` to edit
4. Add company, project, tags
5. Press `Enter` to save
6. Done!

### Workflow 2: Bulk Enrich Imported Events

1. Import CSV (press `i` from home)
2. Auto-navigates to metadata review with imported events
3. Press `a` to select all
4. Press `e` for bulk operations
5. Set Company field with "Apply if empty"
6. Add relevant tags
7. Press `Enter` twice to confirm
8. Done!

### Workflow 3: Focus on Low-Quality Events

1. Open metadata review (press `m` from home)
2. Press `f` to filter
3. Select "Incomplete" or "Basic"
4. View only events needing work
5. Use bulk operations to enrich similar events
6. Press `f` again to see updated quality scores

### Workflow 4: Organize by Company

1. Open metadata review
2. Press `s` to sort by company
3. All TechCorp events grouped together
4. Select all TechCorp events (press `Space` multiple times)
5. Press `e` for bulk operations
6. Add relevant tags/categories for that company
7. Confirm changes

### Workflow 5: Complete Event Details

1. Open metadata review
2. Find event with low quality score
3. Press `Enter` to edit
4. Add all missing fields one by one
5. Watch quality score improve as you fill fields
6. Save when complete

## Tips & Best Practices

### Data Quality Tips

1. **Add Company First**: Most important for grouping and analysis
2. **Use Consistent Names**: Always use "TechCorp", not "Tech Corp" or "TechCorp Inc"
3. **Project Names**: Be specific (e.g., "Platform Migration" not "Project")
4. **Tag Strategy**: Use overlapping tags (e.g., "technical" + "achievement")
5. **Regular Enrichment**: Enrich events soon after capture while memory is fresh

### Efficiency Tips

1. **Bulk Operations**: Group similar events and bulk-edit metadata
2. **Filtering**: Focus on low-quality events first
3. **Sorting**: Group by company to process similar events together
4. **Autocomplete**: Reuse company/project names from history
5. **Keyboard**: Master keyboard shortcuts for faster navigation

### Quality Tips

1. **Aim for "Enriched"**: 51-75 quality is good starting point
2. **Complete is Best**: 76-100 quality provides most value
3. **Tags Matter**: Tags help with future analysis and grouping
4. **Categories Help**: Competency categories useful for career planning
5. **Review Regularly**: Periodic review keeps data quality high

## Troubleshooting

### Issue: Quality score isn't updating

**Solution**: 
- Make sure you pressed `Enter` to save changes
- Check that fields were actually updated (not just focused)
- Try editing again and verify changes are saved

### Issue: Can't select multiple events

**Solution**:
- Make sure you're using `Space` (not `Enter`)
- `Enter` opens the editor instead of selecting
- Try pressing `a` to select all events

### Issue: Bulk operations not applying changes

**Solution**:
- Verify events are selected (press `a` to select all)
- Check that fields are enabled (checkbox should be checked)
- Verify values are entered for selected fields
- Make sure you confirmed the preview

### Issue: Validation errors when saving

**Solution**:
- Check error message (shown in red)
- Fix the specific field mentioned
- Date format: Use YYYY-MM-DD or relative dates
- Company/Project: Max 200 characters
- Tags: Must be from allowed list, max 8
- Categories: Must be from allowed list, max 2

### Issue: Can't find an event

**Solution**:
- Use filter to show all events (press `f`)
- Use search in list view (press `/` then type)
- Sort by date to find events from specific period
- Try navigating with arrow keys

## Advanced Features

### Undo Changes

Press `Ctrl+Z` while editing to undo changes and revert to original values.

### Apply if Empty

In bulk operations, use "Apply if empty" to conditionally update only events missing a field. This preserves existing values.

### Autocomplete

Company and Project fields support autocomplete from previous entries. Start typing and suggestions appear.

### Keyboard Navigation

Master these shortcuts for efficient editing:
- `Tab`: Move between fields
- `Shift+Tab`: Move to previous field
- `Space`: Toggle selections
- `Ctrl+Z`: Undo changes

## See Also

- [CLI_GUIDE.md](./CLI_GUIDE.md) - Main CLI guide with all features
- [CSV_IMPORT_GUIDE.md](./CSV_IMPORT_GUIDE.md) - Import events from CSV
- [README.md](./README.md) - Project overview
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - General troubleshooting
