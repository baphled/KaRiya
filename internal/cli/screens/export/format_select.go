package export

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// FormatSelect allows users to select the export format.
type FormatSelect struct {
	*base.BaseSelectScreen[types.ExportFormat]
}

// NewFormatSelect creates a new format selection screen.
func NewFormatSelect(formats []types.ExportFormat, artifactType types.ExportArtifactType, breadcrumbs []string) *FormatSelect {
	renderer := func(item types.ExportFormat) string {
		descriptions := map[types.ExportFormat]string{
			types.ExportFormatJSON: "JavaScript Object Notation",
			types.ExportFormatCSV:  "Comma Separated Values",
			types.ExportFormatYAML: "YAML Ain't Markup Language",
			types.ExportFormatTXT:  "Plain text format",
			types.ExportFormatMD:   "Markdown format",
			types.ExportFormatPDF:  "Portable Document Format",
		}
		desc := descriptions[item]
		if desc == "" {
			desc = string(item) + " format"
		}
		return fmt.Sprintf("%s\n  %s", string(item), desc)
	}

	baseScreen := base.NewBaseSelectScreen(
		formats,
		renderer,
		breadcrumbs,
		fmt.Sprintf("Select Export Format for %s", artifactType),
	)

	return &FormatSelect{
		BaseSelectScreen: baseScreen,
	}
}
