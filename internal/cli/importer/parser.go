package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// ParsedRow represents a single CSV row with its parsed data and validation status
type ParsedRow struct {
	RowNumber         int // 1-based row number in CSV (excluding header)
	RawData           map[string]string
	Event             *career.CareerEvent
	ValidationErrors  []string
	IsValid           bool
	IsDuplicate       bool
	DuplicateOf       string // ID of duplicate event if found
	MappedCategories  []string // Original categories from CSV
	MappedTags        []string // Original tags from CSV
}

// CSVParser handles parsing and validation of CSV files
type CSVParser struct {
	existingEvents *[]*career.CareerEvent
	dateFormats    []string
	categoryMapper  *CategoryMapper
	tagMapper       *TagMapper
	mapData         bool // Whether to auto-map categories and tags
}

// NewCSVParser creates a new CSV parser
func NewCSVParser(existingEvents []*career.CareerEvent) *CSVParser {
	return &CSVParser{
		existingEvents: &existingEvents,
		dateFormats: []string{
			"2006-01",      // YYYY-MM
			"2006-01-02",   // YYYY-MM-DD
			"01/02/2006",   // MM/DD/YYYY
			"02-01-2006",   // DD-MM-YYYY
			"January 2006", // Month YYYY
			"Jan 2006",     // Mon YYYY
		},
		categoryMapper: NewCategoryMapper(),
		tagMapper:      NewTagMapper(),
		mapData:        false, // Disabled by default for compatibility
	}
}

// NewCSVParserWithMapping creates a CSV parser with auto-mapping enabled
func NewCSVParserWithMapping(existingEvents []*career.CareerEvent) *CSVParser {
	parser := NewCSVParser(existingEvents)
	parser.mapData = true
	return parser
}

// SetMapping enables or disables automatic mapping of categories and tags
func (p *CSVParser) SetMapping(enabled bool) {
	p.mapData = enabled
}

// Parse reads and parses a CSV file (auto-detects delimiter)
func (p *CSVParser) Parse(reader io.Reader) ([]*ParsedRow, error) {
	// Read entire content into buffer to detect delimiter
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detect delimiter from the first line
	delimiter := p.detectDelimiter(string(data))

	// Create CSV reader with detected delimiter
	csvReader := csv.NewReader(strings.NewReader(string(data)))
	csvReader.Comma = delimiter

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Map column names to indices
	columnMap := make(map[string]int)
	for i, col := range header {
		// Trim whitespace from column names
		colName := strings.TrimSpace(col)
		columnMap[colName] = i
	}

	// Validate required columns
	requiredColumns := []string{"Text", "Date"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return nil, fmt.Errorf("missing required column: %s (found columns: %v)", col, getColumnNames(columnMap))
		}
	}

	var parsedRows []*ParsedRow
	rowNumber := 1

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row %d: %w", rowNumber, err)
		}

		// Skip empty rows
		if len(record) == 0 || (len(record) == 1 && strings.TrimSpace(record[0]) == "") {
			continue
		}

		// Build raw data map
		rawData := make(map[string]string)
		for col, idx := range columnMap {
			if idx < len(record) {
				rawData[col] = strings.TrimSpace(record[idx])
			}
		}

		// Parse row
		parsedRow := p.parseRow(rowNumber, rawData, columnMap)
		parsedRows = append(parsedRows, parsedRow)

		rowNumber++
	}

	// Check for duplicates after all rows are parsed
	p.detectDuplicates(parsedRows)

	return parsedRows, nil
}

// detectDelimiter detects if the CSV uses pipes or commas
func (p *CSVParser) detectDelimiter(content string) rune {
	// Get first line
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ','
	}

	firstLine := lines[0]

	// Count pipes and commas
	pipeCount := strings.Count(firstLine, "|")
	commaCount := strings.Count(firstLine, ",")

	// If more pipes than commas and at least 3 pipes, use pipe
	if pipeCount > commaCount && pipeCount >= 3 {
		return '|'
	}

	return ','
}

// getColumnNames returns a slice of column names for error messages
func getColumnNames(columnMap map[string]int) []string {
	var names []string
	for name := range columnMap {
		names = append(names, name)
	}
	return names
}
// parseRow parses a single row and creates a CareerEvent
func (p *CSVParser) parseRow(rowNumber int, rawData map[string]string, columnMap map[string]int) *ParsedRow {
	parsedRow := &ParsedRow{
		RowNumber:        rowNumber,
		RawData:          rawData,
		ValidationErrors: []string{},
		IsValid:          true,
		IsDuplicate:      false,
		MappedCategories: []string{},
		MappedTags:       []string{},
	}

	event := &career.CareerEvent{
		Tags:       []string{},
		Categories: []string{},
	}

	// Parse Text (required)
	text := rawData["Text"]
	if strings.TrimSpace(text) == "" {
		parsedRow.ValidationErrors = append(parsedRow.ValidationErrors, "Text is required")
		parsedRow.IsValid = false
	} else {
		event.Text = text
	}

	// Parse Date (required)
	dateStr := rawData["Date"]
	if strings.TrimSpace(dateStr) == "" {
		parsedRow.ValidationErrors = append(parsedRow.ValidationErrors, "Date is required")
		parsedRow.IsValid = false
	} else {
		parsedDate, err := p.parseDate(dateStr)
		if err != nil {
			parsedRow.ValidationErrors = append(parsedRow.ValidationErrors, fmt.Sprintf("Invalid date format: %s", dateStr))
			parsedRow.IsValid = false
		} else {
			event.Date = parsedDate
		}
	}

	// Parse Company (optional)
	if company, ok := rawData["Company"]; ok && strings.TrimSpace(company) != "" {
		event.Company = company
	}

	// Parse Project (optional)
	if project, ok := rawData["Project"]; ok && strings.TrimSpace(project) != "" {
		event.Project = project
	}

	// Parse Tags (optional, semicolon-separated)
	if tagsStr, ok := rawData["Tags"]; ok && strings.TrimSpace(tagsStr) != "" {
		rawTags := strings.Split(tagsStr, ";")

		// Store original tags
		for _, tag := range rawTags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				parsedRow.MappedTags = append(parsedRow.MappedTags, tag)
			}
		}

		// Map or validate tags
		if p.mapData {
			// Use mapper to convert tags to allowed tags
			event.Tags = p.tagMapper.MapTags(rawTags)
		} else {
			// Original validation-only mode
			for _, tag := range rawTags {
				tag = strings.TrimSpace(strings.ToLower(tag))
				if tag != "" {
					if !career.AllowedTags[tag] {
						parsedRow.ValidationErrors = append(parsedRow.ValidationErrors,
							fmt.Sprintf("Invalid tag: %s (allowed: %v)", tag, getAllowedTagsList()))
						parsedRow.IsValid = false
					} else {
						event.Tags = append(event.Tags, tag)
					}
				}
			}
		}
	}

	// Parse Categories (optional, semicolon-separated)
	if categoriesStr, ok := rawData["Categories"]; ok && strings.TrimSpace(categoriesStr) != "" {
		rawCategories := strings.Split(categoriesStr, ";")

		// Store original categories
		for _, cat := range rawCategories {
			cat = strings.TrimSpace(cat)
			if cat != "" {
				parsedRow.MappedCategories = append(parsedRow.MappedCategories, cat)
			}
		}

		// Map or validate categories
		if p.mapData {
			// Use mapper to convert categories to allowed categories
			event.Categories = p.categoryMapper.MapCategories(rawCategories)
		} else {
			// Original validation-only mode
			for _, cat := range rawCategories {
				cat = strings.TrimSpace(strings.ToLower(cat))
				if cat != "" {
					if !career.AllowedCategories[cat] {
						parsedRow.ValidationErrors = append(parsedRow.ValidationErrors,
							fmt.Sprintf("Invalid category: %s (allowed: %v)", cat, getAllowedCategoriesList()))
						parsedRow.IsValid = false
					} else {
						event.Categories = append(event.Categories, cat)
					}
				}
			}
		}
	}

	// Validate the event if we have valid text and date
	if parsedRow.IsValid && event.Text != "" && !event.Date.IsZero() {
		if err := event.Validate(); err != nil {
			parsedRow.ValidationErrors = append(parsedRow.ValidationErrors, err.Error())
			parsedRow.IsValid = false
		}
	}

	if parsedRow.IsValid {
		parsedRow.Event = event
	}

	return parsedRow
}

// parseDate attempts to parse a date string in multiple formats
func (p *CSVParser) parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Try each format
	for _, format := range p.dateFormats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// detectDuplicates checks for duplicate events in the parsed rows and existing events
func (p *CSVParser) detectDuplicates(parsedRows []*ParsedRow) {
	if p.existingEvents == nil {
		return
	}

	// Build a map of existing events by (text, company, date)
	existingMap := make(map[string]*career.CareerEvent)
	for _, event := range *p.existingEvents {
		if event == nil {
			continue
		}
		key := p.buildDuplicateKey(event.Text, event.Company, event.Date)
		existingMap[key] = event
	}

	// Check each parsed row against existing events and other parsed rows
	for i, parsedRow := range parsedRows {
		if !parsedRow.IsValid || parsedRow.Event == nil {
			continue
		}

		// Check against existing events
		key := p.buildDuplicateKey(parsedRow.Event.Text, parsedRow.Event.Company, parsedRow.Event.Date)
		if existing, exists := existingMap[key]; exists {
			parsedRow.IsDuplicate = true
			parsedRow.DuplicateOf = existing.ID
			parsedRow.ValidationErrors = append(parsedRow.ValidationErrors,
				fmt.Sprintf("Duplicate of existing event (ID: %s)", existing.ID))
			continue
		}

		// Check against other parsed rows
		for j := 0; j < i; j++ {
			otherRow := parsedRows[j]
			if !otherRow.IsValid || otherRow.Event == nil {
				continue
			}

			if p.isSameDuplicateKey(parsedRow.Event, otherRow.Event) {
				parsedRow.IsDuplicate = true
				parsedRow.DuplicateOf = fmt.Sprintf("Row %d", otherRow.RowNumber)
				parsedRow.ValidationErrors = append(parsedRow.ValidationErrors,
					fmt.Sprintf("Duplicate of row %d", otherRow.RowNumber))
				break
			}
		}
	}
}

// buildDuplicateKey creates a unique key for duplicate detection
func (p *CSVParser) buildDuplicateKey(text, company string, date time.Time) string {
	// Use text + company + year-month as duplicate key
	// This allows for some flexibility in exact date matching
	dateKey := date.Format("2006-01")
	return fmt.Sprintf("%s|%s|%s", text, company, dateKey)
}

// isSameDuplicateKey checks if two events have the same duplicate key
func (p *CSVParser) isSameDuplicateKey(event1, event2 *career.CareerEvent) bool {
	key1 := p.buildDuplicateKey(event1.Text, event1.Company, event1.Date)
	key2 := p.buildDuplicateKey(event2.Text, event2.Company, event2.Date)
	return key1 == key2
}

// getAllowedTagsList returns a slice of allowed tags
func getAllowedTagsList() []string {
	var tags []string
	for tag := range career.AllowedTags {
		tags = append(tags, tag)
	}
	return tags
}

// getAllowedCategoriesList returns a slice of allowed categories
func getAllowedCategoriesList() []string {
	var categories []string
	for cat := range career.AllowedCategories {
		categories = append(categories, cat)
	}
	return categories
}

