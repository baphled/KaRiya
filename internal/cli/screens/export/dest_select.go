package export

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// DestSelect allows users to select the export destination.
type DestSelect struct {
	*base.BaseSelectScreen[types.ExportDestination]
}

// NewDestSelect creates a new destination selection screen.
func NewDestSelect(destinations []types.ExportDestination, breadcrumbs []string) *DestSelect {
	renderer := func(item types.ExportDestination) string {
		descriptions := map[types.ExportDestination]string{
			types.ExportDestinationFile:      "Save to a file on disk",
			types.ExportDestinationClipboard: "Copy to system clipboard",
			types.ExportDestinationEmail:     "Send via email",
		}
		desc := descriptions[item]
		if desc == "" {
			desc = string(item)
		}
		return fmt.Sprintf("%s\n  %s", string(item), desc)
	}

	baseScreen := base.NewBaseSelectScreen(
		destinations,
		renderer,
		breadcrumbs,
		"Select Export Destination",
	)

	return &DestSelect{
		BaseSelectScreen: baseScreen,
	}
}
