package generatecv

// State represents the current state of the GenerateCV intent.
type State string

// State constants for the GenerateCV intent wizard workflow.
const (
	// StateConfiguring - User configures CV via wizard modal.
	StateConfiguring State = "configuring"

	// StateExtracting - Extracting technologies from user skills.
	StateExtracting State = "extracting"

	// StateGenerating - CV is being generated.
	StateGenerating State = "generating"

	// StateReview - User reviews CV metadata and statistics via ReviewScreen.
	StateReview State = "review"

	// StatePreview - User previews full CV content via CVPreviewScreen.
	StatePreview State = "preview"
)

// ExportFormat defines the export format type.
type ExportFormat string

// Export format constants.
const (
	// ExportFormatText exports CV as plain text.
	ExportFormatText ExportFormat = "text"

	// ExportFormatMarkdown exports CV as markdown.
	ExportFormatMarkdown ExportFormat = "markdown"

	// ExportFormatYAML exports CV as YAML.
	ExportFormatYAML ExportFormat = "yaml"
)

// ExportOption defines where to save the CV.
type ExportOption string

// Export option constants.
const (
	// ExportOptionSaveToFile saves CV to file.
	ExportOptionSaveToFile ExportOption = "save_to_file"

	// ExportOptionClipboard copies CV to clipboard.
	ExportOptionClipboard ExportOption = "clipboard"

	// ExportOptionCancel cancels export.
	ExportOptionCancel ExportOption = "cancel"
)
