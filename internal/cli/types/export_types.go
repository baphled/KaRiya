// Package types defines export-related types for the CLI export wizard.
// These types identify which category of career data the user wants to export.
package types

// ExportArtifactType identifies which category of career data the user wants
// to export. The export wizard presents these as the first selection step, and
// downstream services use the value to determine which repository to query and
// how to structure the output.
type ExportArtifactType string

const (
	// ExportTypeCV selects a generated curriculum vitae document, assembled from the
	// user's profile, facts, and CV configuration.
	ExportTypeCV ExportArtifactType = "cv"
	// ExportTypeEvents selects raw career timeline events, including text, date, company,
	// project, tags, and categories.
	ExportTypeEvents ExportArtifactType = "events"
	// ExportTypeFacts selects curated career facts that have been extracted or confirmed
	// from timeline events.
	ExportTypeFacts ExportArtifactType = "facts"
	// ExportTypeBursts selects burst analysis results, which group related events into
	// activity clusters with confidence scores.
	ExportTypeBursts ExportArtifactType = "bursts"
	// ExportTypeProfile selects the user's professional profile data such as name, summary,
	// and contact information.
	ExportTypeProfile ExportArtifactType = "profile"
)

// ExportFormat specifies the serialisation format for an exported artifact.
// The export wizard presents these as the second selection step after the user
// chooses an artifact type. Not every format is available for every artifact;
// the wizard filters options accordingly.
type ExportFormat string

const (
	// ExportFormatPDF produces a typeset PDF document, primarily used for CV exports.
	ExportFormatPDF ExportFormat = "pdf"
	// ExportFormatJSON produces machine-readable JSON, suitable for programmatic consumption
	// or backup.
	ExportFormatJSON ExportFormat = "json"
	// ExportFormatYAML produces human-readable YAML, useful for configuration-style data
	// and version-controlled backups.
	ExportFormatYAML ExportFormat = "yaml"
	// ExportFormatCSV produces comma-separated values, suitable for spreadsheet import or
	// tabular analysis of events, facts, or skills.
	ExportFormatCSV ExportFormat = "csv"
	// ExportFormatTXT produces unformatted plain text for simple copy-paste use.
	ExportFormatTXT ExportFormat = "txt"
	// ExportFormatMD produces GitHub-flavored Markdown, useful for documentation or README
	// integration.
	ExportFormatMD ExportFormat = "markdown"
)

// ExportDestination controls the delivery target for an exported artifact.
// The export wizard presents these as the final selection step after artifact
// type and format have been chosen.
type ExportDestination string

const (
	// ExportDestinationFile writes the artifact to a local file at a user-specified path.
	ExportDestinationFile ExportDestination = "file"
	// ExportDestinationClipboard copies the artifact content to the system clipboard for
	// immediate pasting into another application.
	ExportDestinationClipboard ExportDestination = "clipboard"
	// ExportDestinationEmail sends the artifact as an email attachment or body to a
	// user-specified address.
	ExportDestinationEmail ExportDestination = "email"
)
