# CSV Format Guide for KaRiya

## Supported CSV Formats

KaRiya's CSV import parser supports **both comma-delimited and pipe-delimited CSV files** with automatic delimiter detection.

### Format 1: Comma-Delimited CSV (Standard)

```
Text,Date,Categories,Tags,Project,Company
"Event description",2024-01,Technical;Backend,php;api,ProjectName,CompanyName
```

**Characteristics**:
- Delimiter: `,` (comma)
- Quoting: Use quotes for fields with commas
- Headers: Single row with column names

### Format 2: Pipe-Delimited CSV (Your Format)

```
Text|Date|Categories|Tags|Project|Company
Event description here|2024-01|Technical;Backend|php;api|ProjectName|CompanyName
```

**Characteristics**:
- Delimiter: `|` (pipe)
- Quoting: No quoting needed (pipes rarely appear in text)
- Headers: Single row with column names
- **AUTO-DETECTED**: No configuration needed!

---

## Your CSV File Format

Your CSV file uses **pipe-delimited format** with the following structure:

```
Text | Date | Categories | Tags | Project | Company
```

**Example from your file**:

```
Architected and maintained custom Linux infrastructure with golden disc images for RWDMag | 2000-01 | SysAdmin;Infrastructure | linux;golden-discs |  | RWDMag
```

### Column Details

| Column | Type | Required | Format | Example |
|--------|------|----------|--------|---------|
| **Text** | String | ✅ Yes | 1-2000 characters | "Led cross-functional team..." |
| **Date** | Date | ✅ Yes | YYYY-MM (or YYYY-MM-DD) | 2024-01 or 2024-01-15 |
| **Categories** | String | ❌ No | Semicolon-separated | "Technical;Backend" |
| **Tags** | String | ❌ No | Semicolon-separated | "php;api;deployment" |
| **Project** | String | ❌ No | Project name | "Platform Migration" |
| **Company** | String | ❌ No | Company name | "TechCorp Inc." |

---

## Column Format Details

### Text Field

- **Requirements**:
  - Non-empty (1+ characters)
  - Maximum 2000 characters
  - Can contain any text including pipes (will be handled correctly)

- **Examples**:
  - ✅ "Architected and maintained custom Linux infrastructure"
  - ✅ "Led cross-functional team | improved performance by 40%"
  - ❌ "" (empty)
  - ❌ "This text is longer than 2000 characters..." (too long)

### Date Field

- **Requirements**:
  - Required for every event
  - Cannot be in the future
  - Supported formats:

| Format | Example | Notes |
|--------|---------|-------|
| YYYY-MM | 2024-01 | Year-Month (most common in your CSV) |
| YYYY-MM-DD | 2024-01-15 | Year-Month-Day |
| MM/DD/YYYY | 01/15/2024 | US format |
| DD-MM-YYYY | 15-01-2024 | EU format |
| Month YYYY | January 2024 | Text format |
| Mon YYYY | Jan 2024 | Abbreviated format |

- **Validation**:
  - Date cannot be after today
  - Empty dates will cause import errors

### Categories Field

- **Format**: Semicolon-separated list
- **Examples**:
  - `Technical;Backend` (2 categories)
  - `Leadership;Mentoring` (2 categories)
  - `Technical` (1 category)
  - (empty - optional)

- **Your CSV Values**:
  Your CSV uses composite categories like:
  - `Architecture;Assessment`
  - `Backend;Delivery`
  - `SysAdmin;Support`
  - `Leadership;Delivery`

- **KaRiya Conversion** (with mapping enabled):
  - Automatically converted to one of: `technical`, `leadership`, `product`, `consulting`, `research`, `mentoring`
  - See CSV_IMPORT_MAPPING_GUIDE.md for details

### Tags Field

- **Format**: Semicolon-separated list (or comma-separated, parser handles both)
- **Maximum**: 8 tags per event
- **Examples**:
  - `linux;scripts;deployment` (3 tags)
  - `php;ruby;api` (3 tags)
  - `leadership;mentoring` (2 tags)
  - (empty - optional)

- **Your CSV Values**:
  Your CSV uses domain-specific tags like:
  - `linux;golden-discs;scripts`
  - `php;zend;api`
  - `mentoring;code-review`
  - `docker;kubernetes;devops`

- **KaRiya Conversion** (with mapping enabled):
  - Automatically converted to allowed tags: `project`, `achievement`, `leadership`, `technical`, `consulting`, `research`, `product`, `mentoring`
  - See CSV_IMPORT_MAPPING_GUIDE.md for complete mapping dictionary

### Project Field

- **Format**: Simple text string
- **Maximum**: No length limit (reasonable names recommended)
- **Examples**:
  - `Platform Migration`
  - `n-vyro.io`
  - `QuikCV`
  - (empty - optional)

- **Storage**: Stored as-is, no conversion applied

### Company Field

- **Format**: Simple text string
- **Maximum**: No length limit (reasonable names recommended)
- **Examples**:
  - `RWDMag`
  - `We Are Friday`
  - `Nature Publishing Group`
  - (empty - optional)

- **Storage**: Stored as-is, no conversion applied

---

## Automatic Delimiter Detection

The KaRiya CSV parser **automatically detects** whether your file uses commas or pipes.

### Detection Algorithm

```
1. Read first line (header)
2. Count pipes (|) and commas (,)
3. If pipes > commas AND pipes >= 3: use pipe delimiter
4. Otherwise: use comma delimiter
```

### Examples

**Detected as Pipe-Delimited**:
```
Text|Date|Categories|Tags|Project|Company
This is row 1|2024-01|...
```
✅ 5 pipes detected, 0 commas → **Use pipes**

**Detected as Comma-Delimited**:
```
Text,Date,Categories,Tags,Project,Company
"This is row 1",2024-01,...
```
✅ 5 commas detected, 0 pipes → **Use commas**

### You Don't Need To Do Anything!

The parser automatically handles your pipe-delimited format. Just provide the CSV file and it will work.

---

## Whitespace Handling

All fields are **automatically trimmed** of leading/trailing whitespace.

**Before Trimming**:
```
Text                    |Date    |Categories
```

**After Trimming**:
```
Text|Date|Categories
```

This works with both pipe and comma delimiters.

---

## Data Quality Requirements

### What Will Be Accepted

✅ Text field with content + Valid date → **Valid event**
✅ Text field + Date + Categories + Tags → **Valid event**
✅ Text field + Date + Optional fields → **Valid event**

### What Will Be Rejected

❌ Empty Text field → **Invalid** (required)
❌ Empty Date field → **Invalid** (required)
❌ Future date → **Invalid** (cannot be in future)
❌ Text > 2000 characters → **Invalid** (too long)
❌ Invalid date format → **Invalid** (unparseable)

### Validation Errors

Invalid rows are reported with specific error messages:

```
Row 15: Text is required
Row 42: Invalid date format: "invalid-date"
Row 73: Date cannot be in the future
Row 88: Event text exceeds 2000 character limit
```

---

## Import Process

### Step 1: Prepare CSV File

- ✅ Ensure Text and Date columns exist (required)
- ✅ Include Categories, Tags, Project, Company as needed
- ✅ Use your existing format (pipe or comma-delimited)
- ✅ No conversion needed - parser handles it!

### Step 2: Run Import

```bash
# Using the CLI
./kariya --import path/to/career_entries.csv

# Or programmatically
parser := importer.NewCSVParserWithMapping(events)
file, _ := os.Open("career_entries.csv")
parsedRows, err := parser.Parse(file)
```

### Step 3: Review Results

Parser output includes:

```
Successfully parsed 220 rows
Valid: 215
Invalid: 5

Row 1: Valid event (categories mapped, tags mapped)
Row 2: Valid event (categories mapped, tags mapped)
...
Row 150: Invalid - Text is required
...
```

### Step 4: Review and Adjust Metadata

Use **Metadata Review Screen** to:
- ✅ See original → mapped categories/tags
- ✅ Verify mappings are correct
- ✅ Manually adjust any mis-classifications
- ✅ Add missing information
- ✅ Confirm and save

---

## Common Issues and Solutions

### Issue: "missing required column: Text"

**Cause**: Parser couldn't find Text column in header

**Solution**:
1. Verify CSV has a column named exactly `Text` (case-sensitive)
2. Check if delimiter was detected correctly
3. Try specifying delimiter explicitly (if using programmatic API)

**Your CSV**: ✅ Has Text column, should work fine

### Issue: "Error parsing CSV: [error message]"

**Cause**: Generic parsing error

**Solution**:
1. Check CSV file is not corrupted
2. Verify first line is header (not data)
3. Check for invalid characters or encoding issues
4. Validate date format in Date column

**Your CSV**: ✅ Should work with automatic detection

### Issue: Some rows imported, some rejected

**Cause**: Only rows with valid Text + Date accepted

**Solution**:
1. Check rejected rows' error messages
2. Fix missing Text or Date fields
3. Fix invalid date formats
4. Re-import

**Your CSV**: ✅ Should have most valid rows (only fix specific errors)

### Issue: Categories/tags not what I expected

**Cause**: Automatic mapping applied

**Solution**:
1. Check mapping dictionary in docs/CSV_IMPORT_MAPPING_GUIDE.md
2. Review categories/tags in metadata review screen
3. Manually adjust in metadata review screen if needed
4. Update mapping dictionary if your semantics are different

**Your CSV**: ✅ Mapping will convert domain tags to allowed tags

---

## Integration Example

### CLI Usage

```bash
# Import CSV with auto-mapping
./kariya --import career_entries.csv

# Review results in metadata review screen
# Adjust any mis-classified events
# Confirm import
```

### Programmatic Usage

```go
// Create parser with mapping (for your CSV format)
parser := importer.NewCSVParserWithMapping(existingEvents)

// Open CSV file
file, err := os.Open("career_entries.csv")
if err != nil {
    panic(err)
}
defer file.Close()

// Parse (automatic delimiter detection)
parsedRows, err := parser.Parse(file)
if err != nil {
    panic(err)
}

// Import valid events
for _, row := range parsedRows {
    if row.IsValid && row.Event != nil {
        id, err := service.CaptureEvent(ctx, row.Event, mode)
        if err != nil {
            log.Printf("Failed to import row %d: %v", row.RowNumber, err)
            continue
        }
        log.Printf("Imported event %s from row %d", id, row.RowNumber)
    } else {
        log.Printf("Row %d validation errors: %v", row.RowNumber, row.ValidationErrors)
    }
}
```

---

## Summary

✅ **Your CSV Format Supported**: Pipe-delimited format with whitespace-padded headers
✅ **Automatic Detection**: No configuration needed, delimiter auto-detected
✅ **Data Mapping**: Composite categories and domain tags automatically converted
✅ **Quality Validation**: Invalid rows reported with specific error messages
✅ **Manual Review**: Metadata review screen for adjustments
✅ **Production Ready**: Fully tested and documented

Your CSV file should import successfully with these enhancements!

