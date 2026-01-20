package export

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// TypeSelect allows users to select the artifact type to export.
type TypeSelect struct {
	*base.BaseSelectScreen[types.ExportArtifactType]
}

// NewTypeSelect creates a new artifact type selection screen.
func NewTypeSelect(artifactTypes []types.ExportArtifactType, breadcrumbs []string) *TypeSelect {
	renderer := func(item types.ExportArtifactType) string {
		descriptions := map[types.ExportArtifactType]string{
			types.ExportTypeEvents:  "Export career events",
			types.ExportTypeFacts:   "Export extracted facts",
			types.ExportTypeBursts:  "Export career bursts",
			types.ExportTypeProfile: "Export profile data",
			types.ExportTypeCV:      "Export generated CV",
		}
		desc := descriptions[item]
		if desc == "" {
			desc = "Export " + string(item)
		}
		return fmt.Sprintf("%s\n  %s", string(item), desc)
	}

	baseScreen := base.NewBaseSelectScreen(
		artifactTypes,
		renderer,
		breadcrumbs,
		"Select Artifact Type",
	)

	return &TypeSelect{
		BaseSelectScreen: baseScreen,
	}
}
