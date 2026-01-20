package types

// ExportArtifactType represents the type of artifact to export
type ExportArtifactType string

const (
	ExportTypeCV      ExportArtifactType = "cv"
	ExportTypeEvents  ExportArtifactType = "events"
	ExportTypeFacts   ExportArtifactType = "facts"
	ExportTypeBursts  ExportArtifactType = "bursts"
	ExportTypeProfile ExportArtifactType = "profile"
)

// ExportFormat represents the export file format
type ExportFormat string

const (
	ExportFormatPDF  ExportFormat = "pdf"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatYAML ExportFormat = "yaml"
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatTXT  ExportFormat = "txt"
	ExportFormatMD   ExportFormat = "markdown"
)

// ExportDestination represents where to export the artifact
type ExportDestination string

const (
	ExportDestinationFile      ExportDestination = "file"
	ExportDestinationClipboard ExportDestination = "clipboard"
	ExportDestinationEmail     ExportDestination = "email"
)
