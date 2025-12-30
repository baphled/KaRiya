# KaRiya CSV Import Guide

## Overview

The KaRiya Career Journal CLI supports importing career events from CSV files. This feature allows you to bulk-import events from spreadsheets, with built-in validation, duplicate detection, interactive review, and automatic metadata enrichment before importing.

## CSV Format

### Required Columns
- **Text** (required): Event description (1-2000 characters)
- **Date** (required): Event date in flexible format

### Optional Columns
- **Categories**: Semicolon-separated competency categories
- **Tags**: Semicolon-separated tags
- **Project**: Project name
- **Company**: Company name

### Example CSV

```csv
Text,Date,Categories,Tags,Project,Company
Directed end-to-end technical strategy and delivery for an early-stage startup,2024-07,Delivery;Strategy,leadership;technical;consulting,Platform Architecture & MVP,FullSpektrum
Defined system architecture and technology stack to articulate product differentiation,2024-07,Architecture;Strategy,technical;architecture;consulting,Platform Architecture & MVP,FullSpektrum
Mentored engineering contributors and established accountability-focused delivery practices,2024-10,Leadership;Delivery,mentoring;leadership;process,Engineering Enablement,FullSpektrum
Led a major platform migration to improve system performance and reliability,2022-06,Architecture;Delivery,technical;migration;architecture,Core Platform,Mindful Chef
```

## Date Formats

The CSV parser supports multiple date formats:

| Format | Example |
|--------|---------|
| YYYY-MM | 2024-07 |
| YYYY-MM-DD | 2024-07-15 |
| MM/DD/YYYY | 07/15/2024 |
| DD-MM-YYYY | 15-07-2024 |
| Month YYYY | July 2024 |
| Mon YYYY | Jul 2024 |

## Allowed Tags

```
project, achievement, leadership, technical, consulting, research, product, mentoring
```

## Allowed Categories

```
technical, leadership, product, consulting, research, mentoring
```

## Using the Import Feature

### Basic Usage

```bash
./kariya-cli
# Then press 'i' from home screen to start import
```

### Interactive Review Process

1. **File Selection**: Choose CSV file to import
2. **File Parsing**: CSV file is parsed and validated
3. **Duplicate Detection**: System checks for duplicates against existing events
4. **Review Screen**: Interactive UI shows all parsed rows with status:
   - ✓ Valid (green) - Ready to import
   - ✗ Invalid (red) - Has validation errors
   - ⚠ Duplicate (yellow) - Matches existing event
5. **Selection**: Choose which rows to import
6. **Confirmation**: Review summary before finalizing
7. **Import**: Events are created in the database
8. **Metadata Review**: Automatically navigate to metadata enrichment screen

### Review Screen Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `space` | Toggle selection |
| `a` | Select all valid rows |
| `d` | Deselect all rows |
| `Enter` | Start import |
| `Esc` / `q` | Cancel |

## Validation Rules

### Text Field
- Required (cannot be empty)
- Maximum 2000 characters
- Whitespace trimmed

### Date Field
- Required (cannot be empty)
- Cannot be in the future
- Supports multiple formats (see Date Formats above)

### Tags
- Must be from allowed tags list
- Semicolon-separated (e.g., "technical;leadership")
- Maximum 8 tags per event
- Case-insensitive

### Categories
- Must be from allowed categories list
- Semicolon-separated (e.g., "technical;leadership")
- Maximum 2 categories per event

### Company (Optional)
- Maximum 200 characters
- Whitespace trimmed

### Project (Optional)
- Maximum 200 characters
- Whitespace trimmed

## Post-Import Metadata Review

After successful import, you'll automatically navigate to the **Metadata Review Screen** where you can:

### Individual Event Editing
1. Navigate to an imported event
2. Press `Enter` to edit
3. Update metadata fields:
   - Date
   - Company
   - Project
   - Tags
   - Categories
4. Press `Enter` to save or `Esc` to cancel

### Bulk Metadata Operations
1. At metadata review screen, select multiple events (press `Space`)
2. Press `a` to select all imported events
3. Press `e` to enter bulk operations mode
4. Choose fields to update:
   - Company (with "apply if empty" option)
   - Project (with "apply if empty" option)
   - Tags (with "apply if empty" option)
   - Categories (with "apply if empty" option)
5. Preview changes
6. Confirm to apply

### Data Quality Indicators
- Each event shows a quality score (0-100)
- Quality levels: Incomplete, Basic, Enriched, Complete
- See which fields are missing to improve quality
- Visual color coding helps identify events needing attention

### Filtering & Sorting
- Filter by quality level (view only incomplete events)
- Sort by date, company, or creation order
- Focus on events needing the most attention

## Example Workflow

### Step 1: Prepare CSV File

Create `career_events.csv`:

```csv
Text,Date,Categories,Tags,Project,Company
Led technical architecture redesign,2023-06,Architecture;Delivery,technical;architecture,Platform Redesign,TechCorp
Mentored junior developers on Go best practices,2023-07,Leadership;Mentoring,mentoring;leadership,Engineering,TechCorp
Consulted on cloud migration strategy,2023-08,Consulting;Architecture,consulting;technical,Cloud Migration,CloudServices Inc
```

### Step 2: Start Import from Home Screen

```bash
./kariya-cli
# Press 'i' from home screen
```

### Step 3: Review Events

The import screen shows:
```
Total: 3 | Valid: 3 | Invalid: 0 | Duplicates: 0

✓ SEL [1] Led technical architecture redesign
✓ SEL [2] Mentored junior developers on Go best practices
✓ SEL [3] Consulted on cloud migration strategy

↑/k - Move up | ↓/j - Move down | space - Toggle | a - Select all
Enter - Import | Esc/q - Cancel
```

### Step 4: Confirm Import

All valid rows are pre-selected. Press `Enter` to import.

### Step 5: Verify Results

```
✓ Success: 3 | ⚠ Skipped: 0 | ✗ Failed: 0

Imported 3 events successfully
```

Automatically navigates to Metadata Review screen.

### Step 6: Enrich Metadata (New!)

You're now in the Metadata Review screen with your imported events:

```
Metadata Review - 3 Events

[1] ◆ Led technical architecture redesign
    Date: 2023-06-01 | Company: TechCorp | Quality: Enriched ██████░░

[2] ◆ Mentored junior developers...
    Date: 2023-07-15 | Company: TechCorp | Quality: Enriched ██████░░

[3] ◆ Consulted on cloud migration strategy
    Date: 2023-08-20 | Company: CloudServices Inc | Quality: Basic ████░░░░

↑/↓ - Navigate | Enter - Edit | Space - Select | e - Bulk Edit | Esc - Back
```

**Option A: Edit Individual Events**
- Press `Enter` on an event to edit
- Add missing metadata (tags, categories)
- Press `Enter` to save

**Option B: Bulk Edit Multiple Events**
- Press `Space` to select events needing similar updates
- Press `a` to select all
- Press `e` to bulk edit
- Update fields (e.g., add "achievement" tag to all)
- Confirm changes

### Step 7: Verify Enrichment

Once enriched, quality scores improve:
```
[1] ◆ Led technical architecture redesign
    Date: 2023-06-01 | Company: TechCorp | Quality: Complete ██████████
```

## Troubleshooting

### "File not found"
- Check file path is correct
- Ensure file exists in current directory
- Use absolute path if needed: `/path/to/events.csv`

### "wrong number of fields"
- CSV has inconsistent column count
- Some rows missing columns
- Fix: Ensure all rows have same number of columns

### "Invalid date format"
- Date doesn't match any recognized format
- Check date format in CSV
- Fix: Use one of the supported formats

### "Duplicate of existing event"
- Event already exists in database
- Duplicate detection matched text + company + month
- Fix: Review and decide whether to import anyway
- Duplicates are marked but can be imported if desired

### All rows marked as invalid
- Check CSV format
- Verify required columns exist
- Check column names match exactly (case-sensitive)
- Fix: Ensure "Text" and "Date" columns are present

### Import succeeded but metadata review screen didn't open
- Check if there were any import errors
- Try navigating to metadata review manually (press 'm' from home)
- Verify imported events appear in the list

## Advanced Usage

### Importing with Specific Mode

The import always uses "Manual Entry" mode, which:
- Accepts any date in the past
- No time restrictions
- Most flexible for bulk imports

### Handling Large Files

For files with many events:
1. Split into smaller CSV files (100-500 events each)
2. Import one file at a time
3. Review and enrich metadata for each batch
4. This makes the metadata enrichment process more manageable

### Bulk Enrichment Strategy

For large imports with missing metadata:

1. **Import all events** from CSV
2. **Filter by quality level** - View only "Incomplete" events
3. **Group by company** - Use sort to group similar events
4. **Bulk edit by company** - Select all TechCorp events, add company metadata
5. **Bulk add tags** - Select similar events, add relevant tags
6. **Final review** - Check quality scores improved

### Updating Existing Events

CSV import creates new events. To update existing events:
1. Edit directly in CLI metadata editor
2. Or delete and re-import with corrected data
3. Use event ID to track which events to update

## Integration with Existing Features

### Capture Modes

CSV imports use **Manual Entry** mode, which:
- Accepts any date in the past
- No time restrictions
- Most flexible for bulk imports

### Event Classification

After import, events can be classified into competency categories:
- Automatic classification based on text and tags
- Categories can be overridden in metadata editor
- Use tags to hint at categories
- Bulk edit to apply categories to multiple events

### Event Management

Imported events are fully managed:
- Can be viewed in event list
- Can be edited in metadata editor
- Can be deleted
- Can be filtered and searched
- Can be enriched with bulk operations

## FAQ

**Q: Can I import events from Excel?**
A: Yes! Export your Excel file as CSV (File > Save As > CSV format), then use the import feature.

**Q: What if I have a very old event?**
A: CSV import supports any past date, so old events (from years ago) are fully supported.

**Q: Can I modify events after importing?**
A: Yes! Use the metadata editor or bulk operations to enrich imported events. The metadata review screen opens automatically after import.

**Q: How do I know if an import succeeded?**
A: The import results screen shows success count. You'll automatically be taken to the metadata review screen with your imported events.

**Q: Can I import duplicate events intentionally?**
A: The duplicate detection can be overridden by deselecting duplicates before import, but duplicates are marked for your awareness.

**Q: What's the maximum file size?**
A: No specific limit, but very large files (1000+ rows) should be split for better performance and easier metadata enrichment.

**Q: Can I undo an import?**
A: Currently no undo feature. Delete events individually if needed. Consider testing with a small CSV first.

**Q: Can I bulk edit imported events?**
A: Yes! After import, you're in the metadata review screen. Press 'Space' to select events, then 'e' to bulk edit metadata.

**Q: How do I enrich metadata for imported events?**
A: After import, you'll see the metadata review screen. Use individual editing (Enter key) or bulk operations (Space + 'e') to add company, project, tags, and categories.

**Q: What's the difference between Tags and Categories?**
A: Tags are flexible labels (up to 8 per event), while Categories are structured competency areas (max 2 per event). Use categories to classify skills/competencies.

## See Also

- [README.md](../README.md) - Main project documentation
- [CLI_GUIDE.md](./CLI_GUIDE.md) - Interactive CLI guide
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Troubleshooting guide
