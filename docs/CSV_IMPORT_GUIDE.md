# KaRiya CSV Import Guide

## Overview

The KaRiya Career Journal CLI supports importing career events from CSV files. This feature allows you to bulk-import events from spreadsheets, with built-in validation, duplicate detection, and interactive review before importing.

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
./kariya-cli --import events.csv
```

### Interactive Review Process

1. **File Parsing**: CSV file is parsed and validated
2. **Duplicate Detection**: System checks for duplicates against existing events
3. **Review Screen**: Interactive UI shows all parsed rows with status:
   - ✓ Valid (green) - Ready to import
   - ✗ Invalid (red) - Has validation errors
   - ⚠ Duplicate (yellow) - Matches existing event
4. **Selection**: Choose which rows to import
5. **Confirmation**: Review summary before finalizing
6. **Import**: Events are created in the database

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
- Case-insensitive (normalized to lowercase)
- No duplicate tags
- Maximum 8 tags per event

### Categories
- Must be from allowed categories list
- Semicolon-separated (e.g., "technical;leadership")
- Case-insensitive (normalized to lowercase)
- Optional field

### Company & Project
- Optional fields
- No length restrictions

## Duplicate Detection

The import system uses intelligent duplicate detection:

**Duplicate Key**: Text + Company + Year-Month

Two events are considered duplicates if they have:
- The same event text
- The same company
- The same year and month (day is ignored for flexibility)

**Examples**:
- Same text, same company, same month → **DUPLICATE**
- Same text, different company → **NOT a duplicate**
- Same text, same company, different month → **NOT a duplicate**

## Error Messages

### Common Validation Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "Text is required" | Empty text field | Add event description |
| "Date is required" | Empty date field | Add event date |
| "text cannot exceed 2000 characters" | Text too long | Shorten description |
| "date cannot be in the future" | Future date | Use past date |
| "Invalid tag: xxx" | Unknown tag | Use allowed tags |
| "Invalid category: xxx" | Unknown category | Use allowed categories |
| "Duplicate of existing event" | Matches existing event | Review or skip |

## Import Results

After import completes, you'll see:

- **Success**: Number of events successfully imported
- **Skipped**: Number of duplicate or invalid events not imported
- **Failed**: Number of events that failed to import

## Best Practices

### Preparing Your CSV

1. **Clean Data**
   - Remove extra whitespace
   - Ensure dates are in recognized format
   - Verify tags are from allowed list

2. **Organize**
   - Sort by date (newest or oldest first)
   - Group by company if possible
   - Review for duplicates before importing

3. **Validate**
   - Check text length (max 2000 chars)
   - Verify no future dates
   - Confirm all tags are valid

### During Import

1. **Review Carefully**
   - Read error messages for invalid rows
   - Understand why rows are marked as duplicates
   - Decide whether to import duplicates or skip them

2. **Selective Import**
   - Don't import all rows blindly
   - Skip rows with validation errors
   - Skip obvious duplicates
   - Import only rows you want

3. **Verify After**
   - List events to confirm import succeeded
   - Check event details are correct
   - Review any duplicates that were skipped

## Example Workflow

### Step 1: Prepare CSV File

Create `career_events.csv`:

```csv
Text,Date,Categories,Tags,Project,Company
Led technical architecture redesign,2023-06,Architecture;Delivery,technical;architecture,Platform Redesign,TechCorp
Mentored junior developers on Go best practices,2023-07,Leadership;Mentoring,mentoring;leadership,Engineering,TechCorp
Consulted on cloud migration strategy,2023-08,Consulting;Architecture,consulting;technical,Cloud Migration,CloudServices Inc
```

### Step 2: Run Import

```bash
./kariya-cli --import career_events.csv
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

All valid rows are pre-selected. Press Enter to import.

### Step 5: Verify Results

```
✓ Success: 3 | ⚠ Skipped: 0 | ✗ Failed: 0

Imported 3 events successfully

Press any key to continue...
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

### All rows marked as invalid
- Check CSV format
- Verify required columns exist
- Check column names match exactly (case-sensitive)
- Fix: Ensure "Text" and "Date" columns are present

## Advanced Usage

### Importing with Specific Mode

The import always uses "Manual Entry" mode, which allows any past date. This is the safest mode for bulk imports.

### Handling Large Files

For files with many events:
1. Split into smaller CSV files
2. Import one file at a time
3. Review each import before moving to next

### Updating Existing Events

CSV import creates new events. To update existing events:
1. Edit directly in CLI (when available)
2. Delete and re-import with corrected data
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
- Categories can be overridden manually
- Use tags to hint at categories

### Event Management

Imported events are fully managed:
- Can be viewed in event list
- Can be edited (when feature available)
- Can be deleted
- Can be filtered and searched

## FAQ

**Q: Can I import events from Excel?**
A: Yes! Export your Excel file as CSV (File > Save As > CSV format), then use the import feature.

**Q: What if I have a very old event?**
A: CSV import supports any past date, so old events (from years ago) are fully supported.

**Q: Can I modify events after importing?**
A: Yes, imported events can be edited directly in the CLI (when editing feature is available).

**Q: How do I know if an import succeeded?**
A: The import results screen shows success count. You can also list events to verify they were created.

**Q: Can I import duplicate events intentionally?**
A: The duplicate detection can be overridden by deselecting duplicates before import, but duplicates are marked for your awareness.

**Q: What's the maximum file size?**
A: No specific limit, but very large files (1000+ rows) should be split for better performance.

**Q: Can I undo an import?**
A: Currently no undo feature. Delete events individually if needed. Consider testing with a small CSV first.

## See Also

- [README.md](../README.md) - Main project documentation
- [CLI_GUIDE.md](./CLI_GUIDE.md) - Interactive CLI guide
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Troubleshooting guide

