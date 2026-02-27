package generatecv

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
	cvmodals "github.com/baphled/kariya/internal/cli/screens/cv/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// Intent orchestrates the multi-step CV generation workflow, managing
// state transitions between configuration, technology extraction, generation,
// review, preview, and export phases.
//
// Side effects:
//   - None at construction; state mutations occur through Update and Init.
type Intent struct {
	*intents.BaseIntent

	context *IntentContext
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

	wizardModal   *cvmodals.ConfigWizardModal
	progressModal *cvmodals.ProgressModal
	exportModal   *cvmodals.ExportModal

	// --- Screen Orchestration ---

	activeScreen  screens.Screen
	reviewScreen  *cv.ReviewScreen
	previewScreen *cv.CVPreviewScreen
}
