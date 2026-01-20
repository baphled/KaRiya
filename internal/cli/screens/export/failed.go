package export

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens/base"
)

// FailedError contains the error data for the failed screen.
type FailedError struct {
	Code    string
	Message string
}

// Failed displays the export failure screen.
type Failed struct {
	*base.BaseDetailScreen[*FailedError]
}

// NewFailed creates a new failure screen.
func NewFailed(err *FailedError, breadcrumbs []string) *Failed {
	renderer := func(data *FailedError, _, _ int) string {
		var b strings.Builder
		b.WriteString("Export Failed\n\n")
		b.WriteString(fmt.Sprintf("Error: %s\n", data.Message))
		if data.Code != "" {
			b.WriteString(fmt.Sprintf("Code: %s\n", data.Code))
		}
		return b.String()
	}

	baseScreen := base.NewBaseDetailScreen(
		breadcrumbs,
		renderer,
		err,
	)
	baseScreen.SetFooter("r: Retry  Esc: Back")
	baseScreen.AddAction("r", "retry")

	return &Failed{
		BaseDetailScreen: baseScreen,
	}
}
