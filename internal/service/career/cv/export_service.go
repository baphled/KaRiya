package cv

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"gopkg.in/yaml.v3"
)

// ExportService handles exporting CVs to various formats
type ExportService struct {
	logger *logger.Logger
}

// ExportFormat defines the export format type
type ExportFormat string

const (
	// ExportFormatText exports CV as plain text
	ExportFormatText ExportFormat = "text"
	// ExportFormatMarkdown exports CV as markdown
	ExportFormatMarkdown ExportFormat = "markdown"
	// ExportFormatYAML exports CV as YAML
	ExportFormatYAML ExportFormat = "yaml"
)

// ExportResult contains the result of an export operation
type ExportResult struct {
	Format   ExportFormat
	Content  string
	FilePath string
	SavedAt  time.Time
}

// NewExportService creates a new export service
func NewExportService(logger *logger.Logger) *ExportService {
	return &ExportService{
		logger: logger,
	}
}

// ExportToText exports a CV to plain text format
func (es *ExportService) ExportToText(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet) (string, error) {
	if cv == nil {
		return "", fmt.Errorf("CV view is nil")
	}

	var buf bytes.Buffer

	// Write header
	buf.WriteString(strings.ToUpper(cv.Name) + "\n")
	buf.WriteString(strings.Repeat("=", len(cv.Name)) + "\n\n")

	// Write metadata
	buf.WriteString(fmt.Sprintf("Target Role: %s\n", cv.TargetRole))
	buf.WriteString(fmt.Sprintf("Target Audiences: %s\n", strings.Join(cv.TargetAudience, ", ")))
	buf.WriteString(fmt.Sprintf("Generated: %s\n", cv.GeneratedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("Source Events: %d\n", cv.SourceEventCount))
	buf.WriteString(fmt.Sprintf("Source Facts: %d\n\n", cv.SourceFactCount))

	// Write sections
	for _, section := range sections {
		buf.WriteString(strings.ToUpper(section.Title) + "\n")
		buf.WriteString(strings.Repeat("-", len(section.Title)) + "\n")

		// Get bullets for this section
		sectionBullets := bullets[section.ID]
		if len(sectionBullets) == 0 {
			buf.WriteString("(No content)\n\n")
			continue
		}

		for _, bullet := range sectionBullets {
			buf.WriteString(fmt.Sprintf("• %s\n", bullet.Text))
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// ExportToMarkdown exports a CV to markdown format
func (es *ExportService) ExportToMarkdown(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet) (string, error) {
	if cv == nil {
		return "", fmt.Errorf("CV view is nil")
	}

	var buf bytes.Buffer

	// Write header
	buf.WriteString(fmt.Sprintf("# %s\n\n", cv.Name))

	// Write metadata as comment
	buf.WriteString(fmt.Sprintf("<!-- Target Role: %s -->\n", cv.TargetRole))
	buf.WriteString(fmt.Sprintf("<!-- Target Audiences: %s -->\n", strings.Join(cv.TargetAudience, ", ")))
	buf.WriteString(fmt.Sprintf("<!-- Generated: %s -->\n", cv.GeneratedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("<!-- Source Events: %d, Facts: %d -->\n\n", cv.SourceEventCount, cv.SourceFactCount))

	// Write sections
	for _, section := range sections {
		buf.WriteString(fmt.Sprintf("## %s\n\n", section.Title))

		// Get bullets for this section
		sectionBullets := bullets[section.ID]
		if len(sectionBullets) == 0 {
			buf.WriteString("*(No content)*\n\n")
			continue
		}

		for _, bullet := range sectionBullets {
			buf.WriteString(fmt.Sprintf("- %s\n", bullet.Text))
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// ExportToYAML exports a CV to YAML format
func (es *ExportService) ExportToYAML(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet) (string, error) {
	if cv == nil {
		return "", fmt.Errorf("CV view is nil")
	}

	// Build a structured output
	output := map[string]interface{}{
		"name":               cv.Name,
		"target_role":        cv.TargetRole,
		"target_audiences":   cv.TargetAudience,
		"generated_at":       cv.GeneratedAt,
		"source_event_count": cv.SourceEventCount,
		"source_fact_count":  cv.SourceFactCount,
		"sections":           []map[string]interface{}{},
	}

	// Add sections
	sectionsList := output["sections"].([]map[string]interface{})
	for _, section := range sections {
		sectionData := map[string]interface{}{
			"title":   section.Title,
			"type":    section.SectionType,
			"bullets": []string{},
		}

		// Get bullets for this section
		sectionBullets := bullets[section.ID]
		bulletsList := make([]string, 0, len(sectionBullets))
		for _, bullet := range sectionBullets {
			bulletsList = append(bulletsList, bullet.Text)
		}
		sectionData["bullets"] = bulletsList

		sectionsList = append(sectionsList, sectionData)
	}
	output["sections"] = sectionsList

	// Marshal to YAML
	data, err := yaml.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CV to YAML: %w", err)
	}

	return string(data), nil
}

// SaveToFile saves exported CV content to a file
func (es *ExportService) SaveToFile(ctx context.Context, cvName string, format ExportFormat, content string) (string, error) {
	// Determine export directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	exportDir := filepath.Join(homeDir, ".kariya", "cv_exports")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	extension := getFileExtension(format)
	filename := fmt.Sprintf("%s_%s%s", sanitizeFilename(cvName), timestamp, extension)
	filePath := filepath.Join(exportDir, filename)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	es.logger.Info("CV exported to %s format at %s", format, filePath)
	return filePath, nil
}

// getFileExtension returns the file extension for a given format
func getFileExtension(format ExportFormat) string {
	switch format {
	case ExportFormatText:
		return ".txt"
	case ExportFormatMarkdown:
		return ".md"
	case ExportFormatYAML:
		return ".yaml"
	default:
		return ".txt"
	}
}

// sanitizeFilename removes invalid filename characters
func sanitizeFilename(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// Remove invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "")
	}

	// Truncate to reasonable length
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}

// GetExportPath returns the default export directory path
func (es *ExportService) GetExportPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".kariya", "cv_exports"), nil
}

// CopyToClipboard copies the given content to the system clipboard
func (es *ExportService) CopyToClipboard(ctx context.Context, content string) error {
	if content == "" {
		return fmt.Errorf("content is empty")
	}

	if err := clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	es.logger.Info("CV content copied to clipboard")
	return nil
}
