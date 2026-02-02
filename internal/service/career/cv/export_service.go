package cv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"gopkg.in/yaml.v3"
)

// ClipboardWriter defines the interface for clipboard operations.
type ClipboardWriter interface {
	WriteAll(text string) error
	IsUnsupported() bool
}

// SystemClipboard implements ClipboardWriter using the system clipboard.
type SystemClipboard struct{}

// WriteAll writes text to the system clipboard.
func (s *SystemClipboard) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// IsUnsupported returns true if clipboard is not available in this environment.
// This performs an actual write test because clipboard.Unsupported is unreliable
// (it checks if utilities exist, not if they actually work).
func (s *SystemClipboard) IsUnsupported() bool {
	// First check the library's flag
	if clipboard.Unsupported {
		return true
	}

	// Test if clipboard actually works by attempting a small write
	// This catches cases where utilities exist but no display is available
	testErr := clipboard.WriteAll("")
	return testErr != nil
}

// ExportService handles exporting CVs to various formats.
type ExportService struct {
	logger    *logger.Logger
	clipboard ClipboardWriter
}

// ExportFormat defines the export format type.
type ExportFormat string

const (
	// ExportFormatText exports CV as plain text.
	ExportFormatText ExportFormat = "text"
	// ExportFormatMarkdown exports CV as markdown.
	ExportFormatMarkdown ExportFormat = "markdown"
	// ExportFormatYAML exports CV as YAML.
	ExportFormatYAML ExportFormat = "yaml"
)

// ExportResult contains the result of an export operation.
type ExportResult struct {
	Format   ExportFormat
	Content  string
	FilePath string
	SavedAt  time.Time
}

// NewExportService creates a new export service.
func NewExportService(log *logger.Logger) *ExportService {
	return &ExportService{
		logger:    log,
		clipboard: &SystemClipboard{},
	}
}

// NewExportServiceWithClipboard creates a new export service with a custom clipboard implementation.
func NewExportServiceWithClipboard(log *logger.Logger, clipWriter ClipboardWriter) *ExportService {
	return &ExportService{
		logger:    log,
		clipboard: clipWriter,
	}
}

// ExportToText exports a CV to plain text format.
func (es *ExportService) ExportToText(_ context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet) (string, error) {
	if cv == nil {
		return "", errors.New("CV view is nil")
	}

	var buf bytes.Buffer

	// Write header
	buf.WriteString(strings.ToUpper(cv.Name) + "\n")
	buf.WriteString(strings.Repeat("=", len(cv.Name)) + "\n\n")

	// Write metadata
	buf.WriteString(fmt.Sprintf("Target Role: %s\n", cv.TargetRole))
	buf.WriteString(fmt.Sprintf("Target Audience: %s\n", cv.TargetAudience))
	buf.WriteString(fmt.Sprintf("Generated: %s\n", cv.GeneratedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("Source Events: %d\n", cv.SourceEventCount))
	buf.WriteString(fmt.Sprintf("Source Facts: %d\n\n", cv.SourceFactCount))

	// Write sections
	for _, section := range sections {
		buf.WriteString(strings.ToUpper(section.Title) + "\n")
		buf.WriteString(strings.Repeat("-", len(section.Title)) + "\n\n")

		// Handle summary section (prose)
		if section.SectionType == "summary" && section.Summary != "" {
			buf.WriteString(section.Summary + "\n\n")
			continue
		}

		// Handle content groups (experience, projects, skills)
		if len(section.Content) == 0 {
			buf.WriteString("(No content)\n\n")
			continue
		}

		for _, group := range section.Content {
			// Write header with date range (if present)
			if group.Header != "" {
				if group.StartDate != "" && group.EndDate != "" {
					if group.StartDate == group.EndDate {
						buf.WriteString(fmt.Sprintf("%s - %s\n\n", group.Header, group.StartDate))
					} else {
						buf.WriteString(fmt.Sprintf("%s - %s - %s\n\n", group.Header, group.StartDate, group.EndDate))
					}
				} else {
					buf.WriteString(group.Header + "\n\n")
				}
			}

			// Write bullets
			for _, bullet := range group.Bullets {
				buf.WriteString(fmt.Sprintf("• %s\n", bullet.Text))
			}
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}

// ExportToMarkdown exports a CV to markdown format.
func (es *ExportService) ExportToMarkdown(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet) (string, error) {
	if cv == nil {
		return "", errors.New("CV view is nil")
	}

	var buf bytes.Buffer

	// Write header
	buf.WriteString(fmt.Sprintf("# %s\n\n", cv.Name))

	// Write metadata as comment
	buf.WriteString(fmt.Sprintf("<!-- Target Role: %s -->\n", cv.TargetRole))
	buf.WriteString(fmt.Sprintf("<!-- Target Audience: %s -->\n", cv.TargetAudience))
	buf.WriteString(fmt.Sprintf("<!-- Generated: %s -->\n", cv.GeneratedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("<!-- Source Events: %d, Facts: %d -->\n\n", cv.SourceEventCount, cv.SourceFactCount))

	// Write sections
	for _, section := range sections {
		buf.WriteString(fmt.Sprintf("## %s\n\n", section.Title))

		// Handle summary section (prose)
		if section.SectionType == "summary" && section.Summary != "" {
			buf.WriteString(section.Summary + "\n\n")
			continue
		}

		// Handle content groups (experience, projects, skills)
		if len(section.Content) == 0 {
			buf.WriteString("*(No content)*\n\n")
			continue
		}

		for _, group := range section.Content {
			// Write header with date range (if present)
			if group.Header != "" {
				if group.StartDate != "" && group.EndDate != "" {
					if group.StartDate == group.EndDate {
						buf.WriteString(fmt.Sprintf("### %s - _%s_\n\n", group.Header, group.StartDate))
					} else {
						buf.WriteString(fmt.Sprintf("### %s - _%s - %s_\n\n", group.Header, group.StartDate, group.EndDate))
					}
				} else {
					buf.WriteString(fmt.Sprintf("### %s\n\n", group.Header))
				}
			}

			// Write bullets
			for _, bullet := range group.Bullets {
				buf.WriteString(fmt.Sprintf("- %s\n", bullet.Text))
			}
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}

// ExportToYAML exports a CV to YAML format.
func (es *ExportService) ExportToYAML(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet) (string, error) {
	if cv == nil {
		return "", errors.New("CV view is nil")
	}

	// Build a structured output
	output := map[string]interface{}{
		"name":               cv.Name,
		"target_role":        cv.TargetRole,
		"target_audience":    cv.TargetAudience,
		"generated_at":       cv.GeneratedAt,
		"source_event_count": cv.SourceEventCount,
		"source_fact_count":  cv.SourceFactCount,
		"sections":           []map[string]interface{}{},
	}

	// Add sections (sections already have the correct structure with Content groups)
	output["sections"] = sections

	// Marshal to YAML
	data, err := yaml.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CV to YAML: %w", err)
	}

	return string(data), nil
}

// SaveToFile saves exported CV content to a file.
func (es *ExportService) SaveToFile(ctx context.Context, cvName string, format ExportFormat, content string) (string, error) {
	// Determine export directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	exportDir := filepath.Join(homeDir, ".kariya", "cv_exports")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(exportDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	extension := getFileExtension(format)
	filename := fmt.Sprintf("%s_%s%s", sanitizeFilename(cvName), timestamp, extension)
	filePath := filepath.Join(exportDir, filename)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	es.logger.Info("CV exported to %s format at %s", format, filePath)
	return filePath, nil
}

// getFileExtension returns the file extension for a given format.
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

// sanitizeFilename removes invalid filename characters.
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

// GetExportPath returns the default export directory path.
func (es *ExportService) GetExportPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".kariya", "cv_exports"), nil
}

// ErrClipboardUnsupported is returned when clipboard operations are not available in the environment.
// On Linux, this typically means no display server is available (e.g., running over SSH).
// On macOS and Windows, clipboard support is built-in and this error should not occur.
var ErrClipboardUnsupported = errors.New("clipboard not available in headless environment (SSH/no display). Use 'Save to file' instead")

// CopyToClipboard copies the given content to the system clipboard.
func (es *ExportService) CopyToClipboard(ctx context.Context, content string) error {
	if content == "" {
		return errors.New("content is empty")
	}

	if es.clipboard.IsUnsupported() {
		return ErrClipboardUnsupported
	}

	if err := es.clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	es.logger.Info("CV content copied to clipboard")
	return nil
}

// Export exports a CV using the specified structure and format.
// For YAML format, always uses standard structure (it's a data format).
// Uses default profile for narrative structure.
func (es *ExportService) Export(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, structure Structure, format ExportFormat) (string, error) {
	return es.ExportWithProfile(ctx, cv, sections, bullets, structure, format, nil)
}

// ExportWithProfile exports a CV using the specified structure, format, and profile config.
// For YAML format, always uses standard structure (it's a data format).
// If profileCfg is nil, uses default profile.
func (es *ExportService) ExportWithProfile(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, structure Structure, format ExportFormat, profileCfg *config.ProfileConfig) (string, error) {
	if cv == nil {
		return "", errors.New("CV view is nil")
	}

	// YAML always uses standard structure (it's data, not presentation)
	if format == ExportFormatYAML {
		return es.ExportToYAML(ctx, cv, sections, bullets)
	}

	// Route to structure-specific renderer
	switch structure {
	case CVStructureNarrative:
		return es.exportNarrativeWithProfile(ctx, cv, sections, format, profileCfg)
	case CVStructureConsulting:
		return es.exportConsultingWithProfile(ctx, cv, sections, bullets, format, profileCfg)
	case CVStructureHighlights:
		return es.exportHighlightsWithProfile(ctx, cv, sections, bullets, format, profileCfg)
	case CVStructureStandard:
		return es.exportStandard(ctx, cv, sections, bullets, format)
	default:
		return es.exportStandard(ctx, cv, sections, bullets, format)
	}
}

// exportStandard exports using the standard CV structure.
func (es *ExportService) exportStandard(
	ctx context.Context, cv *career.CVView, sections []*career.CVSection,
	bullets map[string][]*career.CVBullet, format ExportFormat,
) (string, error) {
	switch format {
	case ExportFormatText:
		return es.ExportToText(ctx, cv, sections, bullets)
	case ExportFormatMarkdown:
		return es.ExportToMarkdown(ctx, cv, sections, bullets)
	default:
		return "", fmt.Errorf("unknown export format: %s", format)
	}
}

// exportConsultingWithProfile exports using the consulting CV structure with optional profile config.
func (es *ExportService) exportConsultingWithProfile(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, format ExportFormat, profileCfg *config.ProfileConfig) (string, error) {
	switch format {
	case ExportFormatText:
		return es.exportConsultingText(ctx, cv, sections, bullets, profileCfg)
	case ExportFormatMarkdown:
		return es.exportConsultingMarkdown(ctx, cv, sections, bullets, profileCfg)
	default:
		return "", fmt.Errorf("unknown export format: %s", format)
	}
}

// exportConsultingText exports consulting CV to plain text format.
func (es *ExportService) exportConsultingText(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, profileCfg *config.ProfileConfig) (string, error) {
	var buf bytes.Buffer
	profile := NarrativeProfileFromConfig(profileCfg)

	// Profile header
	buf.WriteString(strings.ToUpper(cv.Name) + "\n")
	buf.WriteString(strings.Repeat("=", len(cv.Name)) + "\n\n")
	buf.WriteString(profile.Role + "\n")
	buf.WriteString(profile.Location + "\n")
	buf.WriteString(fmt.Sprintf("Email: %s\n", profile.Email))
	buf.WriteString(fmt.Sprintf("GitHub: %s\n\n", forms.GitHubURL(profile.GitHub)))

	buf.WriteString(strings.Repeat("-", 80) + "\n\n")

	// Summary section
	buf.WriteString("SUMMARY\n")
	buf.WriteString(strings.Repeat("-", 7) + "\n\n")
	summary := getSummaryFromSections(sections)
	if summary != "" {
		buf.WriteString(summary + "\n\n")
	}

	// Client Engagements section (experience sections with "Client Engagements" title)
	buf.WriteString("CLIENT ENGAGEMENTS\n")
	buf.WriteString(strings.Repeat("-", 18) + "\n\n")

	for _, section := range sections {
		if section.SectionType == "experience" {
			for _, group := range section.Content {
				// Company header with dates
				if group.Header != "" {
					if group.StartDate != "" && group.EndDate != "" {
						buf.WriteString(fmt.Sprintf("%s | %s - %s\n", group.Header, group.StartDate, group.EndDate))
					} else {
						buf.WriteString(group.Header + "\n")
					}
				}

				// Bullets from the bullets map
				sectionBullets := bullets[section.ID]
				for _, bullet := range sectionBullets {
					buf.WriteString(fmt.Sprintf("  * %s\n", bullet.Text))
				}
				buf.WriteString("\n")
			}
		}
	}

	// What I Bring section
	if len(profile.ValuePropositions) > 0 {
		buf.WriteString("WHAT I BRING\n")
		buf.WriteString(strings.Repeat("-", 12) + "\n\n")
		for _, prop := range profile.ValuePropositions {
			buf.WriteString(fmt.Sprintf("  * %s\n", prop))
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// exportConsultingMarkdown exports consulting CV to markdown format.
func (es *ExportService) exportConsultingMarkdown(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, profileCfg *config.ProfileConfig) (string, error) {
	var buf bytes.Buffer
	profile := NarrativeProfileFromConfig(profileCfg)

	// Profile header
	buf.WriteString(fmt.Sprintf("# %s\n\n", cv.Name))
	buf.WriteString(fmt.Sprintf("**%s** | %s\n\n", profile.Role, profile.Location))
	buf.WriteString(fmt.Sprintf("Email: %s | GitHub: %s\n\n", profile.Email, forms.GitHubURL(profile.GitHub)))
	buf.WriteString("---\n\n")

	// Summary section
	buf.WriteString("## Summary\n\n")
	summary := getSummaryFromSections(sections)
	if summary != "" {
		buf.WriteString(summary + "\n\n")
	}

	// Client Engagements section
	buf.WriteString("## Client Engagements\n\n")

	for _, section := range sections {
		if section.SectionType == "experience" {
			for _, group := range section.Content {
				// Company header with dates
				if group.Header != "" {
					if group.StartDate != "" && group.EndDate != "" {
						buf.WriteString(fmt.Sprintf("### %s\n*%s - %s*\n\n", group.Header, group.StartDate, group.EndDate))
					} else {
						buf.WriteString(fmt.Sprintf("### %s\n\n", group.Header))
					}
				}

				// Bullets from the bullets map
				sectionBullets := bullets[section.ID]
				for _, bullet := range sectionBullets {
					buf.WriteString(fmt.Sprintf("- %s\n", bullet.Text))
				}
				buf.WriteString("\n")
			}
		}
	}

	// What I Bring section
	if len(profile.ValuePropositions) > 0 {
		buf.WriteString("## What I Bring\n\n")
		for _, prop := range profile.ValuePropositions {
			buf.WriteString(fmt.Sprintf("- %s\n", prop))
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// exportHighlightsWithProfile exports using the highlights CV structure with optional profile config.
func (es *ExportService) exportHighlightsWithProfile(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, format ExportFormat, profileCfg *config.ProfileConfig) (string, error) {
	switch format {
	case ExportFormatText:
		return es.exportHighlightsText(ctx, cv, sections, bullets, profileCfg)
	case ExportFormatMarkdown:
		return es.exportHighlightsMarkdown(ctx, cv, sections, bullets, profileCfg)
	default:
		return "", fmt.Errorf("unknown export format: %s", format)
	}
}

// exportHighlightsText exports highlights CV to plain text format.
func (es *ExportService) exportHighlightsText(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, profileCfg *config.ProfileConfig) (string, error) {
	var buf bytes.Buffer
	profile := NarrativeProfileFromConfig(profileCfg)

	// Condensed profile header
	buf.WriteString(strings.ToUpper(cv.Name) + "\n")
	buf.WriteString(strings.Repeat("=", len(cv.Name)) + "\n")
	buf.WriteString(fmt.Sprintf("%s | %s | %s\n\n", profile.Role, profile.Location, profile.Email))

	// Short summary
	summary := getSummaryFromSections(sections)
	if summary != "" {
		buf.WriteString(summary + "\n\n")
	}

	buf.WriteString(strings.Repeat("-", 60) + "\n\n")

	// Key Capabilities (from core strengths or derive from facts)
	buf.WriteString("KEY CAPABILITIES\n")
	buf.WriteString(strings.Repeat("-", 16) + "\n\n")
	if len(profile.CoreStrengths) > 0 {
		// Use first 4-6 core strengths
		maxStrengths := 6
		if len(profile.CoreStrengths) < maxStrengths {
			maxStrengths = len(profile.CoreStrengths)
		}
		for i := range maxStrengths {
			buf.WriteString(fmt.Sprintf("  * %s\n", profile.CoreStrengths[i]))
		}
	} else {
		buf.WriteString("  * Technical leadership and architecture\n")
		buf.WriteString("  * System design and optimization\n")
		buf.WriteString("  * Cross-functional collaboration\n")
	}
	buf.WriteString("\n")

	// Selected Highlights (top 5 bullets by confidence)
	buf.WriteString("SELECTED HIGHLIGHTS\n")
	buf.WriteString(strings.Repeat("-", 19) + "\n\n")

	topBullets := es.getTopBulletsByConfidence(bullets, 5)
	for _, bullet := range topBullets {
		buf.WriteString(fmt.Sprintf("  * %s\n", bullet.Text))
	}
	buf.WriteString("\n")

	// Technologies section
	if len(profile.Languages) > 0 || len(profile.Systems) > 0 {
		buf.WriteString("TECHNOLOGIES\n")
		buf.WriteString(strings.Repeat("-", 12) + "\n\n")
		if len(profile.Languages) > 0 {
			buf.WriteString(fmt.Sprintf("Languages: %s\n", strings.Join(profile.Languages, ", ")))
		}
		if len(profile.Systems) > 0 {
			buf.WriteString(fmt.Sprintf("Systems: %s\n", strings.Join(profile.Systems, ", ")))
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// exportHighlightsMarkdown exports highlights CV to markdown format.
func (es *ExportService) exportHighlightsMarkdown(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, profileCfg *config.ProfileConfig) (string, error) {
	var buf bytes.Buffer
	profile := NarrativeProfileFromConfig(profileCfg)

	// Condensed profile header
	buf.WriteString(fmt.Sprintf("# %s\n\n", cv.Name))
	buf.WriteString(fmt.Sprintf("**%s** | %s | %s\n\n", profile.Role, profile.Location, profile.Email))

	// Short summary
	summary := getSummaryFromSections(sections)
	if summary != "" {
		buf.WriteString(summary + "\n\n")
	}

	buf.WriteString("---\n\n")

	// Key Capabilities
	buf.WriteString("## Key Capabilities\n\n")
	if len(profile.CoreStrengths) > 0 {
		maxStrengths := 6
		if len(profile.CoreStrengths) < maxStrengths {
			maxStrengths = len(profile.CoreStrengths)
		}
		for i := range maxStrengths {
			buf.WriteString(fmt.Sprintf("- %s\n", profile.CoreStrengths[i]))
		}
	} else {
		buf.WriteString("- Technical leadership and architecture\n")
		buf.WriteString("- System design and optimization\n")
		buf.WriteString("- Cross-functional collaboration\n")
	}
	buf.WriteString("\n")

	// Selected Highlights
	buf.WriteString("## Selected Highlights\n\n")
	topBullets := es.getTopBulletsByConfidence(bullets, 5)
	for _, bullet := range topBullets {
		buf.WriteString(fmt.Sprintf("- %s\n", bullet.Text))
	}
	buf.WriteString("\n")

	// Technologies
	if len(profile.Languages) > 0 || len(profile.Systems) > 0 {
		buf.WriteString("## Technologies\n\n")
		if len(profile.Languages) > 0 {
			buf.WriteString(fmt.Sprintf("**Languages:** %s\n\n", strings.Join(profile.Languages, ", ")))
		}
		if len(profile.Systems) > 0 {
			buf.WriteString(fmt.Sprintf("**Systems:** %s\n", strings.Join(profile.Systems, ", ")))
		}
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// getTopBulletsByConfidence returns the top N bullets sorted by confidence descending.
func (es *ExportService) getTopBulletsByConfidence(bullets map[string][]*career.CVBullet, n int) []*career.CVBullet {
	// Collect all bullets
	var allBullets []*career.CVBullet
	for _, sectionBullets := range bullets {
		allBullets = append(allBullets, sectionBullets...)
	}

	for i := range len(allBullets) - 1 {
		for j := i + 1; j < len(allBullets); j++ {
			if allBullets[j].Confidence > allBullets[i].Confidence {
				allBullets[i], allBullets[j] = allBullets[j], allBullets[i]
			}
		}
	}

	// Return top N
	if len(allBullets) > n {
		return allBullets[:n]
	}
	return allBullets
}

// exportNarrativeWithProfile exports using the narrative CV structure with optional profile config.
func (es *ExportService) exportNarrativeWithProfile(
	ctx context.Context, cv *career.CVView, sections []*career.CVSection,
	format ExportFormat, profileCfg *config.ProfileConfig,
) (string, error) {
	switch format {
	case ExportFormatText:
		return es.exportNarrativeTextWithProfile(ctx, cv, sections, profileCfg)
	case ExportFormatMarkdown:
		return es.exportNarrativeMarkdownWithProfile(ctx, cv, sections, profileCfg)
	default:
		return "", fmt.Errorf("unknown export format: %s", format)
	}
}

// exportNarrativeTextWithProfile exports narrative CV to plain text format with optional profile config.
func (es *ExportService) exportNarrativeTextWithProfile(ctx context.Context, cv *career.CVView, sections []*career.CVSection, profileCfg *config.ProfileConfig) (string, error) {
	var buf bytes.Buffer
	profile := NarrativeProfileFromConfig(profileCfg)

	// Profile header
	buf.WriteString(strings.ToUpper(profile.Name) + "\n")
	buf.WriteString(strings.Repeat("=", len(profile.Name)) + "\n\n")
	buf.WriteString(profile.Role + "\n")
	buf.WriteString(profile.Location + "\n")
	buf.WriteString(fmt.Sprintf("Email: %s\n", profile.Email))
	buf.WriteString("GitHub: " + forms.GitHubURL(profile.GitHub) + "\n")
	buf.WriteString("Portfolio: " + profile.Portfolio + "\n\n")

	buf.WriteString(strings.Repeat("-", 80) + "\n\n")

	// Summary section
	buf.WriteString("SUMMARY\n")
	buf.WriteString(strings.Repeat("-", 7) + "\n\n")
	summary := getSummaryFromSections(sections)
	if summary != "" {
		buf.WriteString(summary + "\n\n")
	} else {
		buf.WriteString("Experienced software engineer with strong technical leadership skills.\n\n")
	}

	// Core Strengths section
	buf.WriteString("CORE STRENGTHS\n")
	buf.WriteString(strings.Repeat("-", 14) + "\n\n")
	for _, strength := range profile.CoreStrengths {
		buf.WriteString(fmt.Sprintf("• %s\n", strength))
	}
	buf.WriteString("\n")

	// Languages & Technologies section
	buf.WriteString("LANGUAGES & TECHNOLOGIES\n")
	buf.WriteString(strings.Repeat("-", 24) + "\n\n")
	buf.WriteString(fmt.Sprintf("Languages: %s\n", strings.Join(profile.Languages, ", ")))
	buf.WriteString(fmt.Sprintf("Frontend: %s\n", strings.Join(profile.Frontend, ", ")))
	buf.WriteString(fmt.Sprintf("Systems: %s\n\n", strings.Join(profile.Systems, ", ")))

	// Selected Experience section (filtered by confidence)
	buf.WriteString("SELECTED EXPERIENCE\n")
	buf.WriteString(strings.Repeat("-", 19) + "\n\n")

	experienceSections := getExperienceSections(sections)
	for _, section := range experienceSections {
		for _, group := range section.Content {
			// Filter bullets by confidence
			highConfidenceBullets := filterBulletsByConfidence(group.Bullets, MinConfidenceForNarrative)
			if len(highConfidenceBullets) == 0 {
				continue
			}

			// Group header with dates
			if group.Header != "" {
				if group.StartDate != "" && group.EndDate != "" {
					buf.WriteString(group.Header + "\n")
					buf.WriteString(group.StartDate + " - " + group.EndDate + "\n\n")
				} else {
					buf.WriteString(group.Header + "\n\n")
				}
			}

			// High-confidence bullets only
			for _, bullet := range highConfidenceBullets {
				buf.WriteString(fmt.Sprintf("• %s\n", bullet.Text))
			}
			buf.WriteString("\n")
		}
	}

	// What I Bring section
	buf.WriteString("WHAT I BRING\n")
	buf.WriteString(strings.Repeat("-", 12) + "\n\n")
	for _, prop := range profile.ValuePropositions {
		buf.WriteString(fmt.Sprintf("• %s\n", prop))
	}
	buf.WriteString("\n")

	buf.WriteString(strings.Repeat("-", 80) + "\n")
	buf.WriteString("References available on request.\n")

	return buf.String(), nil
}

// exportNarrativeMarkdownWithProfile exports narrative CV to markdown format with optional profile config.
func (es *ExportService) exportNarrativeMarkdownWithProfile(ctx context.Context, cv *career.CVView, sections []*career.CVSection, profileCfg *config.ProfileConfig) (string, error) {
	var buf bytes.Buffer
	profile := NarrativeProfileFromConfig(profileCfg)

	// Profile header
	buf.WriteString("# " + profile.Name + "\n\n")
	buf.WriteString("**" + profile.Role + "**\n")
	buf.WriteString(profile.Location + "\n")
	buf.WriteString("Email: [" + profile.Email + "](mailto:" + profile.Email + ")\n")
	buf.WriteString("GitHub: " + forms.GitHubURL(profile.GitHub) + "\n")
	buf.WriteString("Portfolio: " + profile.Portfolio + "\n\n")

	buf.WriteString("---\n\n")

	// Summary section
	buf.WriteString("## Summary\n\n")
	summary := getSummaryFromSections(sections)
	if summary != "" {
		buf.WriteString(summary + "\n\n")
	} else {
		buf.WriteString("Experienced software engineer with strong technical leadership skills.\n\n")
	}

	// Core Strengths section
	buf.WriteString("## Core Strengths\n\n")
	for _, strength := range profile.CoreStrengths {
		buf.WriteString(fmt.Sprintf("- %s\n", strength))
	}
	buf.WriteString("\n")

	// Languages & Technologies section
	buf.WriteString("## Languages & Technologies\n\n")
	buf.WriteString(fmt.Sprintf("**Languages:** %s\n", strings.Join(profile.Languages, ", ")))
	buf.WriteString(fmt.Sprintf("**Frontend:** %s\n", strings.Join(profile.Frontend, ", ")))
	buf.WriteString(fmt.Sprintf("**Systems:** %s\n\n", strings.Join(profile.Systems, ", ")))

	// Selected Experience section (filtered by confidence)
	buf.WriteString("## Selected Experience\n\n")

	experienceSections := getExperienceSections(sections)
	for _, section := range experienceSections {
		for _, group := range section.Content {
			// Filter bullets by confidence
			highConfidenceBullets := filterBulletsByConfidence(group.Bullets, MinConfidenceForNarrative)
			if len(highConfidenceBullets) == 0 {
				continue
			}

			// Group header with dates
			if group.Header != "" {
				if group.StartDate != "" && group.EndDate != "" {
					buf.WriteString(fmt.Sprintf("### %s\n", group.Header))
					buf.WriteString(fmt.Sprintf("*%s - %s*\n\n", group.StartDate, group.EndDate))
				} else {
					buf.WriteString(fmt.Sprintf("### %s\n\n", group.Header))
				}
			}

			// High-confidence bullets only
			for _, bullet := range highConfidenceBullets {
				buf.WriteString(fmt.Sprintf("- %s\n", bullet.Text))
			}
			buf.WriteString("\n")
		}
	}

	// What I Bring section
	buf.WriteString("## What I Bring\n\n")
	for _, prop := range profile.ValuePropositions {
		buf.WriteString(fmt.Sprintf("- %s\n", prop))
	}
	buf.WriteString("\n")

	buf.WriteString("---\n\n")
	buf.WriteString("**References available on request.**\n")

	return buf.String(), nil
}
