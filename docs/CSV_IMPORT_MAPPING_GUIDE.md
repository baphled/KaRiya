# CSV Import Mapping Guide

## Overview

The KaRiya CSV importer now includes an intelligent **data mapping layer** that automatically converts your CSV data format to KaRiya's supported format.

### The Problem

Your CSV file uses:
- **Composite Categories** like `Architecture;Assessment`, `Backend;Delivery`
- **Domain-Specific Tags** like `linux`, `api`, `ruby`, `devops`, `cms`

However, KaRiya supports only:
- **6 Competency Categories**: `technical`, `leadership`, `product`, `consulting`, `research`, `mentoring`
- **8 Tags**: `project`, `achievement`, `leadership`, `technical`, `consulting`, `research`, `product`, `mentoring`

### The Solution

The mapping system automatically converts your data format to valid KaRiya format using semantic analysis.

---

## How Mapping Works

### 1. Category Mapping

The system analyzes your composite categories (e.g., `Architecture;Assessment`) by looking for keyword matches and returns the most appropriate competency category.

#### Keyword Mapping Examples

**Technical Category Keywords**:
- Direct: `technical`, `backend`, `frontend`, `devops`
- Technology: `php`, `ruby`, `python`, `javascript`, `golang`, `docker`, `kubernetes`
- Infrastructure: `linux`, `sysadmin`, `infrastructure`, `deployment`
- Quality: `testing`, `debugging`, `performance`, `quality`

**Leadership Category Keywords**:
- Direct: `leadership`, `lead`, `lead-engineer`, `lead-developer`
- Management: `manager`, `strategy`, `decision`, `team-lead`

**Mentoring Category Keywords**:
- Direct: `mentoring`, `training`, `enablement`
- Development: `onboarding`, `coaching`, `knowledge-sharing`

**Product Category Keywords**:
- Direct: `product`, `ux`, `ui`
- Customer-focused: `customer`, `innovation`, `roadmap`

**Consulting Category Keywords**:
- Direct: `consulting`, `advisory`, `clients`
- Agency: `agency`, `client-engagement`

**Research Category Keywords**:
- Direct: `research`, `analysis`, `data`
- Methods: `investigation`, `experiment`, `prototype`

#### Mapping Examples

| CSV Category | Detected Keywords | Mapped Category |
|---|---|---|
| `Architecture;Assessment` | architecture (technical) | technical |
| `Backend;Delivery` | backend (technical), delivery | technical |
| `Leadership;Delivery` | leadership, delivery | leadership |
| `Mentoring;Enablement` | mentoring, enablement | mentoring |
| `DevOps;Automation` | devops (technical), automation | technical |

### 2. Tag Mapping

The system maps domain-specific tags to allowed tags using a comprehensive mapping dictionary.

#### Tag Mapping Examples

| CSV Tag | Maps To | Reason |
|---|---|---|
| `linux` | `technical` | System administration is a technical skill |
| `api` | `technical` | API development is a technical skill |
| `ruby` | `technical` | Programming languages are technical |
| `leadership` | `leadership` | Direct match |
| `mentoring` | `mentoring` | Direct match |
| `customer` | `product` | Customer-focused work is product work |
| `consulting` | `consulting` | Direct match |
| `data` | `research` | Data work involves research/analysis |
| `unknown-tag` | (ignored) | Not in mapping dictionary |

#### Complete Mapping Dictionary

**→ technical**: backend, frontend, php, ruby, python, javascript, golang, go, c++, c#, java, rust, devops, infrastructure, sysadmin, linux, docker, kubernetes, database, sql, mongodb, api, rest, soap, grpc, ci-cd, deployment, automation, testing, quality, performance, debugging, architecture, design, firmware, embedded, iot, hardware

**→ leadership**: leadership, lead, management, manager, strategy, decision, team-lead, tech-lead

**→ product**: product, ux, ui, customer, innovation, roadmap

**→ consulting**: consulting, advisory, clients, agency

**→ research**: research, analysis, data, analytics

**→ mentoring**: mentoring, training, enablement, onboarding, coaching, knowledge-base

**→ achievement**: achievement, delivered, completed, successful, milestone

**→ project**: project, contract, freelance

---

## Using Mapped CSV Import

### Option 1: Automatic Mapping (Recommended)

Use the mapping parser to automatically convert your CSV to valid format:

```go
// Create parser with mapping enabled
parser := importer.NewCSVParserWithMapping(existingEvents)

// Parse CSV
file, _ := os.Open("career_entries.csv")
parsedRows, err := parser.Parse(file)

// Valid events are automatically converted
for _, row := range parsedRows {
    if row.IsValid && row.Event != nil {
        // Event has been converted and validated
        service.CaptureEvent(ctx, row.Event, mode)
    }
}
```

### Option 2: Manual Mapping Toggle

```go
// Create standard parser
parser := importer.NewCSVParser(existingEvents)

// Enable mapping after creation
parser.SetMapping(true)

// Now parsing will use the mapping system
parsedRows, err := parser.Parse(file)
```

### Option 3: Review and Adjust

The parser preserves original data for review:

```go
for _, row := range parsedRows {
    if row.IsValid {
        // Original CSV categories
        fmt.Println("Original categories:", row.MappedCategories)
        // Mapped categories
        fmt.Println("Mapped categories:", row.Event.Categories)

        // Original CSV tags
        fmt.Println("Original tags:", row.MappedTags)
        // Mapped tags
        fmt.Println("Mapped tags:", row.Event.Tags)
    } else {
        // Review errors
        fmt.Println("Validation errors:", row.ValidationErrors)
    }
}
```

---

## What Happens During Mapping

### Process Flow

1. **Read CSV Row**: Extract categories, tags, text, date
2. **Store Original**: Save original categories/tags for reference
3. **Analyze**: Apply keyword matching to categories
4. **Map Categories**: Use semantic analysis to find best-fit competency category
5. **Map Tags**: Convert domain tags to allowed tags
6. **Validate**: Ensure all data meets KaRiya constraints
7. **Create Event**: Build CareerEvent with mapped data
8. **Store Originals**: Keep original data in ParsedRow for review

### Examples

#### Example 1: System Administration Event

```
Original CSV:
- Categories: SysAdmin;Infrastructure
- Tags: linux;golden-discs

After Mapping:
- Categories: ["technical"]
- Tags: ["technical"]

Reason:
- "SysAdmin" and "Infrastructure" map to "technical"
- "linux" is technical tag, "golden-discs" has no mapping
```

#### Example 2: Leadership Event

```
Original CSV:
- Categories: Leadership;Delivery
- Tags: mentoring;code-review

After Mapping:
- Categories: ["leadership"]
- Tags: ["mentoring"]

Reason:
- "Leadership" and "Delivery" → "leadership" (priority-based)
- "mentoring" is direct match, "code-review" has no mapping
```

#### Example 3: Complex Event

```
Original CSV:
- Categories: Architecture;Backend;Database
- Tags: design;api;performance;optimization

After Mapping:
- Categories: ["technical"]
- Tags: ["technical"]

Reason:
- Multiple technical keywords → "technical" category
- "design", "api", "performance" all map to "technical"
- "optimization" has no mapping and is filtered
```

---

## Quality & Review

### Validation

Even with mapping enabled, events are still validated for:
- Text: 1-2000 characters (required)
- Date: Not in future (required)
- Company/Project: Optional, stored as-is
- Categories/Tags: Validated after mapping

### Review After Import

After importing with mapping, use the **Metadata Review Screen** to:

1. **Verify Mappings**: See original → mapped values
2. **Adjust Categorization**: If mapping wasn't perfect
3. **Add Tags**: Fill in missing tags if needed
4. **Fix Issues**: Correct any misclassified events

---

## When to Use Mapping

### ✅ Use Mapping When

- Importing career data with domain-specific tags/categories
- Your tags are IT-related (languages, tools, technologies)
- Your categories represent skill areas rather than pure competencies
- You want automated conversion with manual review later

### ⚠️ Be Cautious When

- Your tags have custom business meanings
- You need 100% accuracy without review
- Tag/category semantics are non-standard

### ❌ Don't Use Mapping When

- Your CSV already uses KaRiya's format (technical, leadership, product, etc.)
- You prefer strict validation with no automatic conversion
- You want to review and manually adjust every tag

---

## Configuration Examples

### Example 1: Import with Default Strict Validation

```go
parser := importer.NewCSVParser(events)
// Mapping disabled by default
parsedRows, err := parser.Parse(file)

// Only events that exactly match allowed values will be valid
for _, row := range parsedRows {
    if row.IsValid {
        service.CaptureEvent(ctx, row.Event, mode)
    } else {
        // Handle validation errors
    }
}
```

### Example 2: Import with Auto-Mapping

```go
parser := importer.NewCSVParserWithMapping(events)
// Mapping enabled by default
parsedRows, err := parser.Parse(file)

// All events are attempted to be mapped to valid format
for _, row := range parsedRows {
    if row.IsValid {
        service.CaptureEvent(ctx, row.Event, mode)
    }
}
```

### Example 3: Import, Map, Then Review

```go
// Step 1: Parse with mapping
parser := importer.NewCSVParserWithMapping(events)
parsedRows, err := parser.Parse(file)

// Step 2: Import valid events
validCount := 0
for _, row := range parsedRows {
    if row.IsValid && row.Event != nil {
        id, err := service.CaptureEvent(ctx, row.Event, mode)
        if err == nil {
            validCount++
        }
    }
}

// Step 3: Open metadata review screen to verify/adjust
// User can see original → mapped values
// User can make manual corrections if needed
fmt.Printf("Imported %d events. Review in metadata review screen.\n", validCount)
```

---

## Troubleshooting

### Issue: Event marked invalid after mapping

**Cause**: Event failed domain validation (e.g., text too long, date in future)

**Solution**:
1. Check `row.ValidationErrors` for specific issues
2. Fix the issue in the source CSV
3. Re-import

### Issue: Important tag was filtered out

**Cause**: Tag not in mapping dictionary

**Solution**:
1. Check `row.MappedTags` to see what was preserved
2. After import, use metadata review to manually add the tag
3. Or update the tag mapping dictionary

### Issue: Category mapped incorrectly

**Cause**: Keyword matching chose wrong category

**Solution**:
1. Use metadata review screen to correct categorization
2. Or provide more specific category keywords in CSV

---

## Integration with Metadata Review

After CSV import with mapping, events appear in the **Metadata Review Screen** where you can:

1. **See Mappings**: View original CSV categories/tags vs. mapped values
2. **Edit Individually**: Correct any misclassified events
3. **Bulk Correct**: If multiple events have the same issue
4. **Validate**: Confirm all changes are valid

This provides a safe way to import data with confidence that you can review and correct any mapping issues before final confirmation.

---

## FAQ

**Q: Will my original data be lost?**
A: No. Original CSV categories/tags are preserved in `ParsedRow` for review and reference.

**Q: Can I improve the mapping?**
A: Yes. The mapping dictionaries in `mapper.go` can be updated to add more keywords or tag mappings.

**Q: What if a tag doesn't map?**
A: Tags with no mapping are filtered out. You can add them manually in the metadata review screen.

**Q: Is mapping reversible?**
A: You can see original→mapped values in metadata review and adjust manually if needed.

**Q: What's the performance impact?**
A: Minimal. Mapping uses efficient keyword matching. Typical CSV with 100s of events imports in <1 second.

---

## Summary

The mapping system provides:
- ✅ Automatic conversion from domain-specific to KaRiya format
- ✅ Intelligent keyword-based categorization
- ✅ Preservation of original data for review
- ✅ Optional manual review and correction
- ✅ Full validation after mapping

This makes it easy to import career data from external sources while maintaining data quality and allowing fine-tuning through the metadata review interface.

