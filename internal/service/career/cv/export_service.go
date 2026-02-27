//go:generate mockgen -destination=../../../testutil/mocks/service/clipboard_writer_mock.go -package=mocksvc github.com/baphled/kariya/internal/service/career/cv ClipboardWriter

package cv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"gopkg.in/yaml.v3"
)

// audienceSummaryPrefixes maps audience identifiers to summary prefix text.
// When generating highlights for a specific audience, these prefixes replace the
// default profile summary to frame the candidate appropriately.
// Empty string or "master" uses the profile summary as-is (no prefix override).
var audienceSummaryPrefixes = map[string]string{
	"technical-peer":      "Technically deep engineer with",
	"hiring-manager":      "Results-driven engineer delivering",
	"recruiter":           "Versatile software engineer with",
	"engineering-manager": "Collaborative engineer who",
}

// ClipboardWriter defines the interface for clipboard operations.
type ClipboardWriter interface {
	WriteAll(text string) error
	IsUnsupported() bool
}

// SystemClipboard implements ClipboardWriter using the system clipboard.
type SystemClipboard struct{}

// WriteAll writes text to the system clipboard.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *SystemClipboard) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// IsUnsupported returns true if clipboard is not available in this environment.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
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
	logger        *logger.Logger
	clipboard     ClipboardWriter
	profileConfig *config.ProfileConfig
	skillRepo     careerrepo.SkillRepository
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

// yamlLink preserves field order: label, url.
type yamlLink struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}

// yamlJob preserves field order: company, position, start_date, end_date, description.
type yamlJob struct {
	Company     string `yaml:"company"`
	Position    string `yaml:"position"`
	StartDate   string `yaml:"start_date"`
	EndDate     string `yaml:"end_date"`
	Description string `yaml:"description"`
}

// yamlProject preserves field order: name, description, url, key_achievements.
type yamlProject struct {
	Name            string   `yaml:"name"`
	Description     string   `yaml:"description"`
	URL             string   `yaml:"url"`
	KeyAchievements []string `yaml:"key_achievements"`
}

// yamlCV is the top-level output struct preserving exact field order.
type yamlCV struct {
	FirstName  string              `yaml:"first_name"`
	LastName   string              `yaml:"last_name"`
	Email      string              `yaml:"email"`
	Location   string              `yaml:"location"`
	Phone      string              `yaml:"phone"`
	Links      []yamlLink          `yaml:"links"`
	Summary    string              `yaml:"summary"`
	Highlights string              `yaml:"highlights"`
	Jobs       []yamlJob           `yaml:"jobs"`
	Projects   []yamlProject       `yaml:"projects"`
	Skills     map[string][]string `yaml:"skills"`
}

// ExportResult contains the result of an export operation.
type ExportResult struct {
	Format   ExportFormat
	Content  string
	FilePath string
	SavedAt  time.Time
}

// NewExportService creates a new export service.
//
// Expected:
//   - logger must be valid.
//   - profileConfig may be nil (uses empty defaults).
//   - skillRepo may be nil (skills section will be empty).
//
// Returns:
//   - A fully initialized ExportService ready for use.
//
// Side effects:
//   - None.
func NewExportService(log *logger.Logger, profileConfig *config.ProfileConfig, skillRepo careerrepo.SkillRepository) *ExportService {
	return &ExportService{
		logger:        log,
		clipboard:     &SystemClipboard{},
		profileConfig: profileConfig,
		skillRepo:     skillRepo,
	}
}

// NewExportServiceWithClipboard creates a new export service with a custom clipboard implementation.
//
// Expected:
//   - log must be a valid logger.
//   - clipboardWriter must be a valid ClipboardWriter.
//   - profileConfig may be nil (uses empty defaults).
//   - skillRepo may be nil (skills section will be empty).
//
// Returns:
//   - A fully initialized ExportService ready for use.
//
// Side effects:
//   - None.
func NewExportServiceWithClipboard(log *logger.Logger, clipboardWriter ClipboardWriter, profileConfig *config.ProfileConfig, skillRepo careerrepo.SkillRepository) *ExportService {
	return &ExportService{
		logger:        log,
		clipboard:     clipboardWriter,
		profileConfig: profileConfig,
		skillRepo:     skillRepo,
	}
}

// NewExportServiceWithDeps creates a new export service with profile config and skill repository.
//
// Expected:
//   - log must be a valid logger.
//   - profileConfig may be nil (uses empty defaults).
//   - skillRepo may be nil (skills section will be empty).
//
// Returns:
//   - A fully initialized ExportService ready for use.
//
// Side effects:
//   - None.
func NewExportServiceWithDeps(log *logger.Logger, profileConfig *config.ProfileConfig, skillRepo careerrepo.SkillRepository) *ExportService {
	return &ExportService{
		logger:        log,
		clipboard:     &SystemClipboard{},
		profileConfig: profileConfig,
		skillRepo:     skillRepo,
	}
}

// ExportToText exports a CV to plain text format.
//
// Expected:
//   - ctx: A valid context (not cancelled).
//   - cv: A non-nil CVView containing CV metadata.
//   - sections: A slice of CV sections to export.
//   - bullets: A map of bullets indexed by section ID.
//
// Returns:
//   - string: The formatted plain text CV.
//   - error: Non-nil if cv is nil.
//
// Side effects:
//   - Writes formatted text to a buffer.
//   - Formats headers, metadata, and bullet points.
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
//
// Expected:
//   - cv must be a valid CVView.
//   - sections contains the CV sections.
//   - bullets maps section names to bullet points.
//
// Returns:
//   - A markdown string and nil on success.
//   - Empty string and error if cv is nil.
//
// Side effects:
//   - None.
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

// ExportToYAML exports a CV to YAML format (flat structure for external tools).
//
// Expected:
//   - cv must be a valid CVView.
//   - sections contains the CV sections.
//   - profileCfg is optional; if non-nil, uses the provided profile; otherwise uses es.profileConfig.
//
// Returns:
//   - A YAML string in flat format and nil on success.
//   - Empty string and error if cv is nil or marshalling fails.
//
// Side effects:
//   - Queries skillRepo if available to populate skills section.
func (es *ExportService) ExportToYAML(ctx context.Context, cv *career.CVView, sections []*career.CVSection, profileCfg *config.ProfileConfig) (string, error) {
	if cv == nil {
		return "", errors.New("CV view is nil")
	}

	var firstName, lastName, email, location, phone, linkedIn, gitHub, portfolio, summary string

	effectiveProfile := profileCfg
	if effectiveProfile == nil {
		effectiveProfile = es.profileConfig
	}

	if effectiveProfile != nil {
		firstName = effectiveProfile.FirstName
		lastName = effectiveProfile.LastName
		email = effectiveProfile.Email
		// Combine Location and Country
		if effectiveProfile.Location != "" && effectiveProfile.Country != "" {
			location = effectiveProfile.Location + " (" + effectiveProfile.Country + ")"
		} else if effectiveProfile.Location != "" {
			location = effectiveProfile.Location
		} else if effectiveProfile.Country != "" {
			location = effectiveProfile.Country
		}
		phone = effectiveProfile.Phone
		linkedIn = effectiveProfile.LinkedIn
		gitHub = effectiveProfile.GitHub
		portfolio = effectiveProfile.Portfolio
	}

	summary = es.getSummaryFromSectionsYAML(sections)

	links := es.buildYAMLLinks(linkedIn, gitHub, portfolio)

	jobs := es.buildYAMLJobs(sections, effectiveProfile)

	projects := es.buildYAMLProjects(sections)

	skillsLimit := 0
	if effectiveProfile != nil {
		skillsLimit = effectiveProfile.SkillsLimit
	}
	skills := es.buildYAMLSkills(ctx, skillsLimit)

	output := yamlCV{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		Location:   location,
		Phone:      phone,
		Links:      links,
		Summary:    summary,
		Highlights: generateHighlights(sections, effectiveProfile),
		Jobs:       jobs,
		Projects:   projects,
		Skills:     skills,
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(output); err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return buf.String(), nil
}

func (es *ExportService) getSummaryFromSectionsYAML(sections []*career.CVSection) string {
	for _, section := range sections {
		if section.SectionType == "summary" && section.Summary != "" {
			return section.Summary
		}
	}
	return ""
}

const (
	githubURLPrefix   = "https://github.com/"
	linkedInURLPrefix = "https://www.linkedin.com/in/"
)

func ensureURL(value, prefix string) string {
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}
	return prefix + value
}

func (es *ExportService) buildYAMLLinks(linkedIn, gitHub, portfolio string) []yamlLink {
	var links []yamlLink

	if linkedIn != "" {
		links = append(links, yamlLink{Label: "LinkedIn", URL: ensureURL(linkedIn, linkedInURLPrefix)})
	}

	if gitHub != "" {
		links = append(links, yamlLink{Label: "GitHub", URL: ensureURL(gitHub, githubURLPrefix)})
	}

	if portfolio != "" {
		links = append(links, yamlLink{Label: "Portfolio", URL: portfolio})
	}

	return links
}

func (es *ExportService) buildYAMLJobs(sections []*career.CVSection, effectiveProfile *config.ProfileConfig) []yamlJob {
	var jobs []yamlJob

	for _, section := range sections {
		if section.SectionType != "experience" {
			continue
		}

		for _, group := range section.Content {
			position := ""
			if effectiveProfile != nil {
				position = effectiveProfile.Title
			}

			jobs = append(jobs, yamlJob{
				Company:     group.Header,
				Position:    position,
				StartDate:   group.StartDate,
				EndDate:     group.EndDate,
				Description: formatBulletsAsDescription(group.Bullets),
			})
		}
	}

	return jobs
}

func (es *ExportService) buildYAMLProjects(sections []*career.CVSection) []yamlProject {
	var projects []yamlProject

	for _, section := range sections {
		if section.SectionType != "projects" {
			continue
		}

		for _, group := range section.Content {
			var keyAchievements []string
			for _, bullet := range group.Bullets {
				keyAchievements = append(keyAchievements, bullet.Text)
			}

			projects = append(projects, yamlProject{
				Name:            group.Header,
				Description:     formatBulletsAsDescription(group.Bullets),
				URL:             "",
				KeyAchievements: keyAchievements,
			})
		}
	}

	return projects
}

func formatBulletsAsDescription(bullets []*career.CVBullet) string {
	var sb strings.Builder
	for _, bullet := range bullets {
		sb.WriteString("- " + bullet.Text + "\n")
	}
	return sb.String()
}

func maxHighlightsFromProfile(profile *config.ProfileConfig) int {
	if profile != nil && profile.MaxHighlights > 0 {
		return profile.MaxHighlights
	}
	return 5
}

func generateHighlights(sections []*career.CVSection, profile *config.ProfileConfig) string {
	if profile != nil && len(profile.WhatIBring) > 0 {
		var sb strings.Builder
		for _, item := range profile.WhatIBring {
			sb.WriteString("- " + item + "\n")
		}
		return sb.String()
	}

	bullets := extractExperienceBullets(sections)
	if len(bullets) > 0 {
		sortBulletsByConfidenceAndRoleScore(bullets)
		maxBullets := maxHighlightsFromProfile(profile)
		if len(bullets) > maxBullets {
			bullets = bullets[:maxBullets]
		}
		var sb strings.Builder
		for _, bullet := range bullets {
			sb.WriteString("- " + bullet.Text + "\n")
		}
		return sb.String()
	}

	if profile != nil && len(profile.CoreStrengths) > 0 {
		var sb strings.Builder
		for _, strength := range profile.CoreStrengths {
			sb.WriteString("- " + strength + "\n")
		}
		return sb.String()
	}

	return ""
}

func extractExperienceBullets(sections []*career.CVSection) []*career.CVBullet {
	var bullets []*career.CVBullet
	for _, section := range sections {
		if section.SectionType != "experience" {
			continue
		}
		for _, group := range section.Content {
			bullets = append(bullets, group.Bullets...)
		}
	}
	return bullets
}

func sortBulletsByConfidenceAndRoleScore(bullets []*career.CVBullet) {
	slices.SortFunc(bullets, func(a, b *career.CVBullet) int {
		if a.Confidence > b.Confidence {
			return -1
		}
		if a.Confidence < b.Confidence {
			return 1
		}
		if a.RoleScore > b.RoleScore {
			return -1
		}
		if a.RoleScore < b.RoleScore {
			return 1
		}
		return 0
	})
}

func (es *ExportService) buildYAMLSkills(ctx context.Context, limit int) map[string][]string {
	skills := make(map[string][]string)

	if limit <= 0 {
		limit = 5
	}

	if es.skillRepo == nil {
		return skills
	}

	allSkills, err := es.skillRepo.List(ctx, nil)
	if err != nil {
		return skills
	}

	for _, skill := range allSkills {
		category := skill.Category
		if len(skills[category]) < limit {
			skills[category] = append(skills[category], skill.Name)
		}
	}

	// Sort skill names within each category
	for category := range skills {
		slices.Sort(skills[category])
	}

	return skills
}

// SaveToFile saves exported CV content to a file.
//
// Expected:
//   - cvName is a non-empty string.
//   - format is a valid ExportFormat.
//   - content is the CV content to save.
//
// Returns:
//   - The file path where the CV was saved and nil on success.
//   - Empty string and error if directory creation or file write fails.
//
// Side effects:
//   - Creates ~/.kariya/cv_exports directory if it does not exist.
//   - Writes a file to disk with timestamp-based filename.
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
//
// Returns:
//   - The export directory path and nil on success.
//   - Empty string and error if home directory cannot be determined.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - cv must be a valid CVView.
//   - sections contains the CV sections.
//   - bullets maps section names to bullet points.
//   - structure is a valid Structure value.
//   - format is a valid ExportFormat.
//
// Returns:
//   - The exported CV content as a string and nil on success.
//   - Empty string and error if cv is nil or export fails.
//
// Side effects:
//   - None.
func (es *ExportService) Export(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, structure Structure, format ExportFormat) (string, error) {
	return es.ExportWithProfile(ctx, cv, sections, bullets, structure, format, nil)
}

// ExportWithProfile exports a CV using the specified structure, format, and profile config.
// For YAML format, always uses standard structure (it's a data format).
// If profileCfg is nil, uses default profile.
//
// Expected:
//   - cv: non-nil CVView with career data
//   - sections: slice of CVSection to include in export
//   - bullets: map of section IDs to CVBullet slices
//   - structure: valid CVStructure (Narrative, Consulting, Highlights, Standard)
//   - format: valid ExportFormat (Text, Markdown, YAML)
//   - profileCfg: optional ProfileConfig (uses default if nil)
//
// Returns:
//   - string: exported CV content
//   - error: nil on success, error if cv is nil or export fails
//
// Side effects:
//   - None.
func (es *ExportService) ExportWithProfile(ctx context.Context, cv *career.CVView, sections []*career.CVSection, bullets map[string][]*career.CVBullet, structure Structure, format ExportFormat, profileCfg *config.ProfileConfig) (string, error) {
	if cv == nil {
		return "", errors.New("CV view is nil")
	}

	// YAML uses flat structure for external tools
	if format == ExportFormatYAML {
		return es.ExportToYAML(ctx, cv, sections, profileCfg)
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

// ExportHighlightsForAudience exports a highlights-format CV tailored to a specific audience.
// It uses GetTopBulletsForAudience instead of getTopBulletsByConfidence and generates
// an audience-specific summary prefix. The existing exportHighlightsWithProfile remains
// unchanged for backward compatibility.
//
// Expected:
//   - view must be a valid CVView with populated sections.
//   - profile must be a valid ProfileConfig (nil uses defaults).
//   - audience identifies the target audience for bullet scoring and summary prefix.
//
// Returns:
//   - A string containing the formatted highlights CV in plain text.
//
// Side effects:
//   - None.
func ExportHighlightsForAudience(view career.CVView, profile config.ProfileConfig, audience string) string {
	var buf bytes.Buffer
	narrative := NarrativeProfileFromConfig(&profile)

	buf.WriteString(strings.ToUpper(view.Name) + "\n")
	buf.WriteString(strings.Repeat("=", len(view.Name)) + "\n")
	buf.WriteString(fmt.Sprintf("%s | %s | %s\n\n", narrative.Role, narrative.Location, narrative.Email))

	summaryPrefix, hasPrefix := audienceSummaryPrefixes[audience]
	if hasPrefix {
		buf.WriteString(summaryPrefix + "\n\n")
	} else {
		summary := getSummaryFromSections(view.Sections)
		if summary != "" {
			buf.WriteString(summary + "\n\n")
		}
	}

	buf.WriteString(strings.Repeat("-", 60) + "\n\n")

	buf.WriteString("KEY CAPABILITIES\n")
	buf.WriteString(strings.Repeat("-", 16) + "\n\n")
	if len(narrative.CoreStrengths) > 0 {
		maxStrengths := 6
		if len(narrative.CoreStrengths) < maxStrengths {
			maxStrengths = len(narrative.CoreStrengths)
		}
		for i := range maxStrengths {
			buf.WriteString(fmt.Sprintf("  * %s\n", narrative.CoreStrengths[i]))
		}
	} else {
		buf.WriteString("  * Technical leadership and architecture\n")
		buf.WriteString("  * System design and optimization\n")
		buf.WriteString("  * Cross-functional collaboration\n")
	}
	buf.WriteString("\n")

	buf.WriteString("SELECTED HIGHLIGHTS\n")
	buf.WriteString(strings.Repeat("-", 19) + "\n\n")

	topBullets := GetTopBulletsForAudience(view, audience, 5)
	for _, bullet := range topBullets {
		buf.WriteString(fmt.Sprintf("  * %s\n", bullet.Text))
	}
	buf.WriteString("\n")

	if len(narrative.Languages) > 0 || len(narrative.Systems) > 0 {
		buf.WriteString("TECHNOLOGIES\n")
		buf.WriteString(strings.Repeat("-", 12) + "\n\n")
		if len(narrative.Languages) > 0 {
			buf.WriteString(fmt.Sprintf("Languages: %s\n", strings.Join(narrative.Languages, ", ")))
		}
		if len(narrative.Systems) > 0 {
			buf.WriteString(fmt.Sprintf("Systems: %s\n", strings.Join(narrative.Systems, ", ")))
		}
		buf.WriteString("\n")
	}

	return buf.String()
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

	// Sort by confidence descending using standard library
	slices.SortFunc(allBullets, func(a, b *career.CVBullet) int {
		// Sort descending: higher confidence comes first
		if a.Confidence > b.Confidence {
			return -1
		}
		if a.Confidence < b.Confidence {
			return 1
		}
		return 0
	})

	// Return top N
	if len(allBullets) > n {
		return allBullets[:n]
	}
	return allBullets
}

// GetTopBulletsForAudience returns the top N bullets scored by audience relevance.
// Each bullet is scored as confidence * audienceRelevance[audience] when audience
// data exists. Falls back to confidence-only when AudienceRelevance is nil, empty,
// or the audience key is missing.
//
// Expected:
//   - view must contain sections with content groups containing bullets.
//   - audience identifies which audience relevance score to use.
//   - limit must be a positive integer.
//
// Returns:
//   - A slice of CVBullet pointers sorted by combined score descending, limited to limit.
//
// Side effects:
//   - None.
func GetTopBulletsForAudience(view career.CVView, audience string, limit int) []*career.CVBullet {
	var allBullets []*career.CVBullet
	for _, section := range view.Sections {
		for _, group := range section.Content {
			allBullets = append(allBullets, group.Bullets...)
		}
	}

	slices.SortFunc(allBullets, func(a, b *career.CVBullet) int {
		scoreA := audienceScore(a, audience)
		scoreB := audienceScore(b, audience)
		if scoreA > scoreB {
			return -1
		}
		if scoreA < scoreB {
			return 1
		}
		return 0
	})

	if len(allBullets) > limit {
		return allBullets[:limit]
	}
	return allBullets
}

// audienceScore computes the combined score for a bullet given an audience.
// Returns confidence * audienceRelevance[audience] when audience data exists,
// or confidence alone as fallback.
func audienceScore(bullet *career.CVBullet, audience string) float64 {
	if len(bullet.AudienceRelevance) == 0 {
		return bullet.Confidence
	}
	relevance, ok := bullet.AudienceRelevance[audience]
	if !ok {
		return bullet.Confidence
	}
	return bullet.Confidence * relevance
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
	buf.WriteString(fmt.Sprintf("GitHub: %s\n", forms.GitHubURL(profile.GitHub)))
	buf.WriteString(fmt.Sprintf("Portfolio: %s\n\n", profile.Portfolio))

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
					buf.WriteString(fmt.Sprintf("%s - %s\n\n", group.StartDate, group.EndDate))
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
	buf.WriteString(fmt.Sprintf("# %s\n\n", profile.Name))
	buf.WriteString(fmt.Sprintf("**%s**\n", profile.Role))
	buf.WriteString(profile.Location + "\n")
	buf.WriteString(fmt.Sprintf("Email: [%s](mailto:%s)\n", profile.Email, profile.Email))
	buf.WriteString(fmt.Sprintf("GitHub: %s\n", forms.GitHubURL(profile.GitHub)))
	buf.WriteString(fmt.Sprintf("Portfolio: %s\n\n", profile.Portfolio))

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
