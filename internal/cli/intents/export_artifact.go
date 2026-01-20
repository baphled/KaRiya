package intents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/types"
	careerdomain "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"gopkg.in/yaml.v3"
)

// Type aliases for export types from internal/cli/types package.
// These allow the intents package to use shorter names while the canonical
// definitions live in the shared types package.
type (
	ExportArtifactType = types.ExportArtifactType
	ExportFormat       = types.ExportFormat
	ExportDestination  = types.ExportDestination
)

// Re-export constants from types package for convenience.
const (
	ExportTypeCV      = types.ExportTypeCV
	ExportTypeEvents  = types.ExportTypeEvents
	ExportTypeFacts   = types.ExportTypeFacts
	ExportTypeBursts  = types.ExportTypeBursts
	ExportTypeProfile = types.ExportTypeProfile

	ExportFormatPDF  = types.ExportFormatPDF
	ExportFormatJSON = types.ExportFormatJSON
	ExportFormatYAML = types.ExportFormatYAML
	ExportFormatCSV  = types.ExportFormatCSV
	ExportFormatTXT  = types.ExportFormatTXT
	ExportFormatMD   = types.ExportFormatMD

	ExportDestinationFile      = types.ExportDestinationFile
	ExportDestinationClipboard = types.ExportDestinationClipboard
	ExportDestinationEmail     = types.ExportDestinationEmail
)

// ExportArtifactContext contains context for the ExportArtifact intent
type ExportArtifactContext struct {
	// Configuration
	ArtifactTypes    []ExportArtifactType
	SupportedFormats map[ExportArtifactType][]ExportFormat
	DefaultFormat    map[ExportArtifactType]ExportFormat
	Destinations     []ExportDestination

	// Services
	ExportService *cv.ExportService
	CareerService *career.Service

	// Repositories
	EventRepository careerrepo.Repository
	FactRepository  careerrepo.FactRepository
	BurstRepository careerrepo.BurstRepository

	// Context
	AppContext context.Context
}

// ExportConfiguration holds the current export configuration
type ExportConfiguration struct {
	ArtifactType ExportArtifactType
	Format       ExportFormat
	Destination  ExportDestination
	FilePath     string
	Email        string
}

// ExportArtifactResult contains the result of artifact export
type ExportArtifactResult struct {
	Success      bool
	ArtifactType ExportArtifactType
	Format       ExportFormat
	Destination  ExportDestination
	FilePath     string
	Size         int64
	Error        *IntentError
}

// ExportState represents the current state of the export intent
type ExportState string

const (
	ExportStateSelectType   ExportState = "select_type"
	ExportStateSelectFormat ExportState = "select_format"
	ExportStateSelectDest   ExportState = "select_destination"
	ExportStateConfigure    ExportState = "configure"
	ExportStatePreview      ExportState = "preview"
	ExportStateConfirm      ExportState = "confirm"
	ExportStateInProgress   ExportState = "in_progress"
	ExportStateComplete     ExportState = "complete"
	ExportStateFailed       ExportState = "failed"
)

// ExportProgressMsg represents progress of an export operation
type ExportProgressMsg struct {
	Percentage int
	Message    string
}

// ExportCompleteMsg represents completion of an export operation
type ExportCompleteMsg struct {
	Result *ExportArtifactResult
	Error  *IntentError
}

// ExportCancelledMsg represents cancellation of an export operation
type ExportCancelledMsg struct{}

// ExportErrorMsg represents an error during export
type ExportErrorMsg struct {
	Error *IntentError
}

// NewExportConfiguration creates a new export configuration with defaults
func NewExportConfiguration(artifactType ExportArtifactType, ctx *ExportArtifactContext) *ExportConfiguration {
	format := ctx.DefaultFormat[artifactType]
	return &ExportConfiguration{
		ArtifactType: artifactType,
		Format:       format,
		Destination:  ExportDestinationFile,
	}
}

// NewExportArtifactResult creates a new export result
func NewExportArtifactResult(success bool, artifactType ExportArtifactType, format ExportFormat, destination ExportDestination, filePath string, size int64) *ExportArtifactResult {
	return &ExportArtifactResult{
		Success:      success,
		ArtifactType: artifactType,
		Format:       format,
		Destination:  destination,
		FilePath:     filePath,
		Size:         size,
	}
}

// NewExportArtifactResultWithError creates a new export result with an error
func NewExportArtifactResultWithError(err *IntentError) *ExportArtifactResult {
	return &ExportArtifactResult{
		Success: false,
		Error:   err,
	}
}

// DefaultArtifactTypes returns the default set of exportable artifact types
func DefaultArtifactTypes() []ExportArtifactType {
	return []ExportArtifactType{
		ExportTypeEvents,
		ExportTypeFacts,
		ExportTypeBursts,
	}
}

// DefaultSupportedFormats returns the default format support mapping
func DefaultSupportedFormats() map[ExportArtifactType][]ExportFormat {
	return map[ExportArtifactType][]ExportFormat{
		ExportTypeEvents: {
			ExportFormatJSON,
			ExportFormatCSV,
			ExportFormatTXT,
			ExportFormatYAML,
		},
		ExportTypeFacts: {
			ExportFormatJSON,
			ExportFormatCSV,
			ExportFormatTXT,
			ExportFormatYAML,
		},
		ExportTypeBursts: {
			ExportFormatJSON,
			ExportFormatCSV,
			ExportFormatTXT,
			ExportFormatYAML,
		},
	}
}

// DefaultFormats returns the default format for each artifact type
func DefaultFormats() map[ExportArtifactType]ExportFormat {
	return map[ExportArtifactType]ExportFormat{
		ExportTypeEvents: ExportFormatJSON,
		ExportTypeFacts:  ExportFormatJSON,
		ExportTypeBursts: ExportFormatJSON,
	}
}

// DefaultDestinations returns the default set of export destinations
func DefaultDestinations() []ExportDestination {
	return []ExportDestination{
		ExportDestinationFile,
		ExportDestinationClipboard,
	}
}

// getFileExtensionForFormat returns the file extension for a given format
func getFileExtensionForFormat(format ExportFormat) string {
	switch format {
	case ExportFormatTXT:
		return ".txt"
	case ExportFormatMD:
		return ".md"
	case ExportFormatJSON:
		return ".json"
	case ExportFormatCSV:
		return ".csv"
	case ExportFormatYAML:
		return ".yaml"
	default:
		return ".txt"
	}
}

// marshalToJSON marshals any data to JSON format
func marshalToJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}

// marshalToYAML marshals any data to YAML format
func marshalToYAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	return string(bytes), nil
}

// marshalEventsToCSV marshals career events to CSV format
func marshalEventsToCSV(events []*careerdomain.CareerEvent) (string, error) {
	var sb strings.Builder
	sb.WriteString("id,date,text,company,project,tags,categories\n")

	for _, event := range events {
		tags := strings.Join(event.Tags, ";")
		categories := strings.Join(event.Categories, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			escapeCSV(event.ID),
			event.Date.Format("2006-01-02"),
			escapeCSV(event.Text),
			escapeCSV(event.Company),
			escapeCSV(event.Project),
			escapeCSV(tags),
			escapeCSV(categories),
		))
	}

	return sb.String(), nil
}

// marshalEventsToText marshals career events to plain text format
func marshalEventsToText(events []*careerdomain.CareerEvent) (string, error) {
	var sb strings.Builder
	sb.WriteString("Career Events\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	for i, event := range events {
		sb.WriteString(fmt.Sprintf("Event %d\n", i+1))
		sb.WriteString(fmt.Sprintf("Date: %s\n", event.Date.Format("2006-01-02")))
		if event.Company != "" {
			sb.WriteString(fmt.Sprintf("Company: %s\n", event.Company))
		}
		if event.Project != "" {
			sb.WriteString(fmt.Sprintf("Project: %s\n", event.Project))
		}
		sb.WriteString(fmt.Sprintf("Text: %s\n", event.Text))
		if len(event.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(event.Tags, ", ")))
		}
		if len(event.Categories) > 0 {
			sb.WriteString(fmt.Sprintf("Categories: %s\n", strings.Join(event.Categories, ", ")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal Events: %d\n", len(events)))
	return sb.String(), nil
}

// marshalFactsToCSV marshals facts to CSV format
func marshalFactsToCSV(facts []*careerdomain.Fact) (string, error) {
	var sb strings.Builder
	sb.WriteString("id,text,competency_categories,role_fit,audience_relevance,strength_signal,source_event_id,source_burst_id\n")

	for _, fact := range facts {
		categories := strings.Join(fact.CompetencyCategories, ";")
		audiences := strings.Join(fact.AudienceRelevance, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n",
			escapeCSV(fact.ID),
			escapeCSV(fact.Text),
			escapeCSV(categories),
			escapeCSV(string(fact.RoleFit)),
			escapeCSV(audiences),
			escapeCSV(fact.StrengthSignal),
			escapeCSV(fact.SourceEventID),
			escapeCSV(fact.SourceBurstID),
		))
	}

	return sb.String(), nil
}

// marshalFactsToText marshals facts to plain text format
func marshalFactsToText(facts []*careerdomain.Fact) (string, error) {
	var sb strings.Builder
	sb.WriteString("Facts\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	for i, fact := range facts {
		sb.WriteString(fmt.Sprintf("Fact %d: %s\n", i+1, fact.Text))
		sb.WriteString(fmt.Sprintf("Competency Categories: %s\n", strings.Join(fact.CompetencyCategories, ", ")))
		sb.WriteString(fmt.Sprintf("Role Fit: %s\n", fact.RoleFit))
		sb.WriteString(fmt.Sprintf("Audience Relevance: %s\n", strings.Join(fact.AudienceRelevance, ", ")))
		sb.WriteString(fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal))
		if fact.SourceEventID != "" {
			sb.WriteString(fmt.Sprintf("Source Event: %s\n", fact.SourceEventID))
		}
		if fact.SourceBurstID != "" {
			sb.WriteString(fmt.Sprintf("Source Burst: %s\n", fact.SourceBurstID))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal Facts: %d\n", len(facts)))
	return sb.String(), nil
}

// marshalBurstsToCSV marshals bursts to CSV format
func marshalBurstsToCSV(bursts []*careerdomain.Burst) (string, error) {
	var sb strings.Builder
	sb.WriteString("id,name,description,event_ids,confirmed\n")

	for _, burst := range bursts {
		eventIDs := strings.Join(burst.EventIDs, ";")
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%t\n",
			escapeCSV(burst.ID),
			escapeCSV(burst.Name),
			escapeCSV(burst.Description),
			escapeCSV(eventIDs),
			burst.Confirmed,
		))
	}

	return sb.String(), nil
}

// marshalBurstsToText marshals bursts to plain text format
func marshalBurstsToText(bursts []*careerdomain.Burst) (string, error) {
	var sb strings.Builder
	sb.WriteString("Bursts\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	for i, burst := range bursts {
		sb.WriteString(fmt.Sprintf("Burst %d: %s\n", i+1, burst.Name))
		if burst.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", burst.Description))
		}
		sb.WriteString(fmt.Sprintf("Event IDs: %s\n", strings.Join(burst.EventIDs, ", ")))
		sb.WriteString(fmt.Sprintf("Confirmed: %t\n", burst.Confirmed))
		if burst.ConfirmedAt != nil {
			sb.WriteString(fmt.Sprintf("Confirmed At: %s\n", burst.ConfirmedAt.Format("2006-01-02 15:04:05")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal Bursts: %d\n", len(bursts)))
	return sb.String(), nil
}

// escapeCSV escapes a string for CSV format
func escapeCSV(s string) string {
	// If string contains comma, quote, or newline, wrap in quotes and escape quotes
	if strings.ContainsAny(s, ",\"\n") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

// SaveToDestination saves content to file or clipboard based on config.
// This is used by the ExportArtifactIntent for actual export operations.
func SaveToDestination(ctx context.Context, exportCtx *ExportArtifactContext, config *ExportConfiguration, name string, content string) (string, error) {
	var filePath string

	switch config.Destination {
	case ExportDestinationFile:
		// Save to file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}

		exportDir := filepath.Join(homeDir, ".kariya", "exports")
		if err := os.MkdirAll(exportDir, 0750); err != nil {
			return "", fmt.Errorf("failed to create export directory: %w", err)
		}

		timestamp := time.Now().Format("20060102_150405")
		extension := getFileExtensionForFormat(config.Format)
		filename := fmt.Sprintf("%s_%s%s", name, timestamp, extension)
		filePath = filepath.Join(exportDir, filename)

		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}

	case ExportDestinationClipboard:
		if exportCtx.ExportService == nil {
			return "", fmt.Errorf("export service not available")
		}
		if err := exportCtx.ExportService.CopyToClipboard(ctx, content); err != nil {
			return "", err
		}
		filePath = "clipboard"

	default:
		return "", fmt.Errorf("unsupported destination: %s", config.Destination)
	}

	return filePath, nil
}
