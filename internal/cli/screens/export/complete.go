package export

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// CompleteResult contains the result data for the complete screen.
type CompleteResult struct {
	ArtifactType types.ExportArtifactType
	Format       types.ExportFormat
	Destination  types.ExportDestination
	FilePath     string
	Size         int64
}

// Complete displays the export completion screen.
type Complete struct {
	*base.BaseDetailScreen[*CompleteResult]
}

// NewComplete creates a new completion screen.
func NewComplete(result *CompleteResult, breadcrumbs []string) *Complete {
	renderer := func(data *CompleteResult, _, _ int) string {
		var b strings.Builder
		b.WriteString("Export Complete!\n\n")
		b.WriteString(fmt.Sprintf("Artifact: %s\n", data.ArtifactType))
		b.WriteString(fmt.Sprintf("Format: %s\n", data.Format))
		b.WriteString(fmt.Sprintf("Destination: %s\n", data.Destination))
		b.WriteString(fmt.Sprintf("File: %s\n", data.FilePath))
		b.WriteString(fmt.Sprintf("Size: %s\n", formatBytes(data.Size)))
		return b.String()
	}

	baseScreen := base.NewBaseDetailScreen(
		breadcrumbs,
		renderer,
		result,
	)
	baseScreen.SetFooter("Enter: Done  Esc: Back")

	return &Complete{
		BaseDetailScreen: baseScreen,
	}
}

// formatBytes converts bytes to human-readable format.
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
