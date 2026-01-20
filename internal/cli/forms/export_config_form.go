package forms

import (
	"github.com/charmbracelet/huh"
)

// ExportConfigFormData holds the export wizard configuration data.
type ExportConfigFormData struct {
	// Step 1: WHAT
	ArtifactType string // "events", "facts", "bursts", "cv", "profile"

	// Step 2: HOW
	Format      string // "json", "yaml", "csv", "txt", "markdown"
	Destination string // "file", "clipboard"

	// Completion flag
	SubmitConfirmed bool
}

// NewExportConfigFormData creates a new ExportConfigFormData with defaults.
func NewExportConfigFormData() *ExportConfigFormData {
	return &ExportConfigFormData{
		ArtifactType: "events",
		Format:       "json",
		Destination:  "file",
	}
}

// NewExportConfigForm creates the export configuration wizard form.
// This form guides users through selecting what to export and how.
//
// Step 1 (WHAT): Select artifact type (events, facts, bursts, cv, profile)
// Step 2 (HOW): Select format and destination
func NewExportConfigForm(data *ExportConfigFormData, width, height int) *huh.Form {
	// Step 1: WHAT - Artifact type selection
	step1Fields := []huh.Field{
		huh.NewSelect[string]().
			Key("artifact_type").
			Title("What would you like to export?").
			Description("Use ↑/↓ to navigate, Enter to select").
			Options(
				huh.NewOption("📅 Career Events", "events"),
				huh.NewOption("💡 Facts", "facts"),
				huh.NewOption("⚡ Bursts", "bursts"),
				huh.NewOption("📄 CV/Resume", "cv"),
				huh.NewOption("👤 Profile", "profile"),
			).
			Value(&data.ArtifactType),
	}

	// Step 2: HOW - Format and destination
	step2Fields := []huh.Field{
		huh.NewSelect[string]().
			Key("format").
			Title("Export Format").
			Description("Choose the file format").
			Options(
				huh.NewOption("JSON (structured data)", "json"),
				huh.NewOption("YAML (human-readable)", "yaml"),
				huh.NewOption("CSV (spreadsheet)", "csv"),
				huh.NewOption("Markdown (documentation)", "markdown"),
				huh.NewOption("Plain Text", "txt"),
			).
			Value(&data.Format),

		huh.NewSelect[string]().
			Key("destination").
			Title("Export Destination").
			Description("Where should we save it?").
			Options(
				huh.NewOption("💾 Save to File", "file"),
				huh.NewOption("📋 Copy to Clipboard", "clipboard"),
			).
			Value(&data.Destination),
	}

	// Create groups for each step
	groups := []*huh.Group{
		huh.NewGroup(step1Fields...).Title("Step 1: What to Export"),
		huh.NewGroup(step2Fields...).Title("Step 2: How to Export"),
	}

	// Use standard form creation with theme
	form := huh.NewForm(groups...).
		WithTheme(Theme()).
		WithWidth(width)

	if height > 0 {
		form = form.WithHeight(height)
	}

	return form
}

// GetArtifactTypeLabel returns a human-readable label for the artifact type.
func GetArtifactTypeLabel(artifactType string) string {
	switch artifactType {
	case "events":
		return "Career Events"
	case "facts":
		return "Facts"
	case "bursts":
		return "Bursts"
	case "cv":
		return "CV/Resume"
	case "profile":
		return "Profile"
	default:
		return artifactType
	}
}

// GetFormatLabel returns a human-readable label for the format.
func GetFormatLabel(format string) string {
	switch format {
	case "json":
		return "JSON"
	case "yaml":
		return "YAML"
	case "csv":
		return "CSV"
	case "markdown":
		return "Markdown"
	case "txt":
		return "Plain Text"
	default:
		return format
	}
}

// GetDestinationLabel returns a human-readable label for the destination.
func GetDestinationLabel(destination string) string {
	switch destination {
	case "file":
		return "File"
	case "clipboard":
		return "Clipboard"
	default:
		return destination
	}
}
