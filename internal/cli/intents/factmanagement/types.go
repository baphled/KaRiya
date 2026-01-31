// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	factmodals "github.com/baphled/kariya/internal/cli/screens/facts/modals"
	domain "github.com/baphled/kariya/internal/domain/career"
)

// Intent orchestrates the fact management workflow.
type Intent struct {
	*intents.BaseIntent

	// Context holds input parameters and business logic.
	context *IntentContext

	// State machine fields.
	state  State
	active bool
	result *intents.IntentResult[*Result]

	// Table behavior for list navigation.
	tableBehavior *behaviors.TableBehavior[*domain.Fact]

	// Edit modal for creating/editing facts.
	editModal *factmodals.EditFactModal
}
