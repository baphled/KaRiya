package generatecv

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/tui/intents"
	cvviews "github.com/baphled/kariya/internal/tui/views/cv"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Intent orchestrates the multi-step CV generation workflow, managing
// state transitions between configuration, technology extraction, generation,
// review, preview, and export phases.
//
// Side effects:
//   - None at construction; state mutations occur through Update and Init.
type Intent struct {
	*intents.BaseIntent

	context *IntentValidator
	state   State
	active  bool
	result  *intents.IntentResult[*Result]
	logger  *logger.Logger

	// --- Flattened state fields (previously in model) ---

	selectedProfile  *CVProfile
	selectedAudience string
	generatedCV      *career.CVView

	extractedTechnologies []*ExtractedTechnology
	focusAreaSuggestion   *FocusAreaSuggestion

	selectedTechnologyFocus cvsvc.TechnologyFocus
	selectedTechnologies    []string
	selectedFocusArea       cvsvc.FocusArea

	selectedSkillsFormat string
	selectedSkillsLimit  int
	selectedCVLength     string

	selectedExportFormat ExportFormat
	selectedExportOption ExportOption
	exportedPath         string
	exportError          error
	isExporting          bool
	exportReturnState    State

	// --- Modals ---

	wizardModal   *cvviews.ConfigWizard
	progressModal *cvviews.Progress
	exportModal   *cvviews.Export

	// --- View Orchestration ---

	activeView  widgets.View
	reviewView  *cvviews.Review
	previewView *cvviews.Preview
}
