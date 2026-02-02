package mocks

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
)

// EventListFilters is re-exported for test convenience.
type EventListFilters = repo.EventListFilters

// TestMockRepository is a test helper for mocking repository behavior.
type TestMockRepository struct {
	createErr     error
	updateErr     error
	deleteErr     error
	getByIDErr    error
	listErr       error
	countErr      error
	getByIDEvent  *career.Event
	getByIDEvents map[string]*career.Event
	getByIDErrors map[string]error
	listEvents    []*career.Event
	countResult   int

	createCalled  bool
	updateCalled  bool
	deleteCalled  bool
	getByIDCalled bool
	listCalled    bool
	countCalled   bool
}

// NewTestMockRepository creates a new behavior-based mock repository.
//
// Returns:
//   - A fully initialized TestMockRepository ready for use.
//
// Side effects:
//   - None.
func NewTestMockRepository() *TestMockRepository {
	return &TestMockRepository{
		getByIDEvents: make(map[string]*career.Event),
		getByIDErrors: make(map[string]error),
	}
}

// SetCreateBehavior sets the error for Create calls.
//
// Expected:
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetCreateBehavior(err error) {
	m.createErr = err
}

// SetUpdateBehavior sets the error for Update calls.
//
// Expected:
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetUpdateBehavior(err error) {
	m.updateErr = err
}

// SetDeleteBehavior sets the error for Delete calls.
//
// Expected:
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetDeleteBehavior(err error) {
	m.deleteErr = err
}

// SetGetByIDBehavior sets the event and error for GetByID calls.
//
// Expected:
//   - event must be valid.
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetGetByIDBehavior(event *career.Event, err error) {
	m.getByIDEvent = event
	m.getByIDErr = err
}

// SetEventByID sets a specific event to be returned for a given ID.
//
// Expected:
//   - Must be a valid string.
//   - event must be valid.
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetEventByID(eventID string, event *career.Event, err error) {
	if event != nil {
		m.getByIDEvents[eventID] = event
	}
	if err != nil {
		m.getByIDErrors[eventID] = err
	}
}

// SetListBehavior sets the events and error for List calls.
//
// Expected:
//   - event must be valid.
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetListBehavior(events []*career.Event, err error) {
	m.listEvents = events
	m.listErr = err
}

// SetCountBehavior sets the count and error for Count calls.
//
// Expected:
//   - int must be valid.
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *TestMockRepository) SetCountBehavior(count int, err error) {
	m.countResult = count
	m.countErr = err
}

// Create implements EventRepository interface.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) Create(_ context.Context, _ *career.Event) error {
	m.createCalled = true
	return m.createErr
}

// Update implements EventRepository interface.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) Update(_ context.Context, _ *career.Event) error {
	m.updateCalled = true
	return m.updateErr
}

// Delete implements EventRepository interface.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) Delete(_ context.Context, _ string) error {
	m.deleteCalled = true
	return m.deleteErr
}

// GetByID implements EventRepository interface.
func (m *TestMockRepository) GetByID(_ context.Context, eventID string) (*career.Event, error) {
	m.getByIDCalled = true

	// Check for per-ID mocking first
	if event, exists := m.getByIDEvents[eventID]; exists {
		if err, hasErr := m.getByIDErrors[eventID]; hasErr {
			return nil, err
		}
		return event, nil
	}

	// Fall back to default behavior
	return m.getByIDEvent, m.getByIDErr
}

// List implements EventRepository interface.
func (m *TestMockRepository) List(_ context.Context, _ repo.EventListFilters) ([]*career.Event, error) {
	m.listCalled = true
	return m.listEvents, m.listErr
}

// Count implements EventRepository interface.
func (m *TestMockRepository) Count(_ context.Context, _ repo.EventListFilters) (int, error) {
	m.countCalled = true
	return m.countResult, m.countErr
}

// CreateCalled returns whether Create was called.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) CreateCalled() bool {
	return m.createCalled
}

// UpdateCalled returns whether Update was called.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) UpdateCalled() bool {
	return m.updateCalled
}

// DeleteCalled returns whether Delete was called.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) DeleteCalled() bool {
	return m.deleteCalled
}

// GetByIDCalled returns whether GetByID was called.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) GetByIDCalled() bool {
	return m.getByIDCalled
}

// ListCalled returns whether List was called.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) ListCalled() bool {
	return m.listCalled
}

// CountCalled returns whether Count was called.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) CountCalled() bool {
	return m.countCalled
}

// LinkSkill implements EventRepository interface.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *TestMockRepository) LinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}
