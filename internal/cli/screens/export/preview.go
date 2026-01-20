package export

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// Preview displays a preview of the export content.
type Preview struct {
	*base.BaseDetailScreen[string]
}

// NewPreview creates a new preview screen.
func NewPreview(content string, artifactType types.ExportArtifactType, format types.ExportFormat, breadcrumbs []string) *Preview {
	renderer := func(data string, _, _ int) string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Preview Export (%s as %s):\n\n", artifactType, format))
		b.WriteString("════════════════════════════════════════════════════════\n\n")
		b.WriteString(data)
		b.WriteString("\n\n════════════════════════════════════════════════════════")
		return b.String()
	}

	baseScreen := base.NewBaseDetailScreen(
		breadcrumbs,
		renderer,
		content,
	)
	baseScreen.SetFooter("↑/↓/j/k: Scroll  Enter: Continue  Esc: Back")

	return &Preview{
		BaseDetailScreen: baseScreen,
	}
}
