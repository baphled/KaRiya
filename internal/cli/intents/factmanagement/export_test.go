package factmanagement

import (
	factmodals "github.com/baphled/kariya/internal/cli/screens/facts/modals"
	"github.com/charmbracelet/huh"
)

// TestGetEditModal exposes the edit modal for testing.
func (i *Intent) TestGetEditModal() *factmodals.EditFactModal {
	return i.editModal
}

// TestSetState exposes the state for testing.
func (i *Intent) TestSetState(state State) {
	i.state = state
}

// TestSetEditModalNil clears the edit modal for testing.
func (i *Intent) TestSetEditModalNil() {
	i.editModal = nil
}

// TestCompleteEditModal marks the edit modal form as completed for testing.
func (i *Intent) TestCompleteEditModal() {
	if i.editModal != nil {
		i.editModal.GetForm().State = huh.StateCompleted
	}
}

// TestAbortEditModal marks the edit modal form as aborted for testing.
func (i *Intent) TestAbortEditModal() {
	if i.editModal != nil {
		i.editModal.GetForm().State = huh.StateAborted
	}
}

// TestPrepareNewFactForSave sets required fields on the editing fact
// so that validation passes for new fact creation.
func (i *Intent) TestPrepareNewFactForSave(id, sourceEventID string) {
	if i.context.EditingFact != nil {
		i.context.EditingFact.ID = id
		i.context.EditingFact.SourceEventID = sourceEventID
	}
}
