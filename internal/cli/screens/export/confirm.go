package export

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// Confirm displays a confirmation screen before export.
type Confirm struct {
	*base.BaseConfirmScreen
}

// NewConfirm creates a new confirmation screen.
func NewConfirm(artifactType types.ExportArtifactType, format types.ExportFormat, destination types.ExportDestination, breadcrumbs []string) *Confirm {
	message := fmt.Sprintf(
		"Confirm export?\n\n"+
			"Artifact: %s\n"+
			"Format: %s\n"+
			"Destination: %s",
		artifactType,
		format,
		destination,
	)

	baseScreen := base.NewBaseConfirmScreen(
		breadcrumbs,
		"Confirm Export",
		message,
	)
	baseScreen.SetYesText("Export")
	baseScreen.SetNoText("Cancel")

	return &Confirm{
		BaseConfirmScreen: baseScreen,
	}
}
