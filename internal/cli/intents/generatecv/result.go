package generatecv

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Result is the output data returned when the GenerateCV intent completes.
type Result struct {
	GeneratedCV     *career.CVView
	SelectedProfile *CVProfile
	AcceptedFields  map[string]bool
	ExportPath      string
	CVExportFormat  string
	ExportedAt      *time.Time
}
