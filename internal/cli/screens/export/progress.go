package export

import (
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// Progress displays progress during export.
type Progress struct {
	*base.BaseProgressScreen
}

// NewProgress creates a new progress screen.
func NewProgress(artifactType types.ExportArtifactType, breadcrumbs []string) *Progress {
	baseScreen := base.NewBaseProgressScreen(
		breadcrumbs,
		"Exporting",
		"Exporting "+string(artifactType)+"...",
	)
	baseScreen.SetAllowCancel(false)

	return &Progress{
		BaseProgressScreen: baseScreen,
	}
}
