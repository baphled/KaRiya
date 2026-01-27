package mocks

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
)

// ListFilters is re-exported for test convenience.
type ListFilters = repo.ListFilters

// TestMockRepository is a test helper for mocking repository behavior.
type TestMockRepository struct {
	createErr     error
	updateErr     error
	deleteErr     error
	getByIDErr    error
	listErr       error
	countErr      error
	getByIDEvent  *career.CareerEvent
	getByIDEvents map[string]*career.CareerEvent // Per-ID event mapping
	getByIDErrors map[string]error               // Per-ID error mapping
	listEvents    []*career.CareerEvent
	countResult   int

	createCalled  bool
	updateCalled  bool
	deleteCalled  bool
	getByIDCalled bool
	listCalled    bool
	countCalled   bool
}

// NewTestMockRepository creates a new behavior-based mock repository.
func NewTestMockRepository() *TestMockRepository {
	return &TestMockRepository{
		getByIDEvents: make(map[string]*career.CareerEvent),
		getByIDErrors: make(map[string]error),
	}
}

// SetCreateBehavior sets the error for Create calls.
func (m *TestMockRepository) SetCreateBehavior(err error) {
	m.createErr = err
}

// SetUpdateBehavior sets the error for Update calls.
func (m *TestMockRepository) SetUpdateBehavior(err error) {
	m.updateErr = err
}

// SetDeleteBehavior sets the error for Delete calls.
func (m *TestMockRepository) SetDeleteBehavior(err error) {
	m.deleteErr = err
}

// SetGetByIDBehavior sets the event and error for GetByID calls.
func (m *TestMockRepository) SetGetByIDBehavior(event *career.CareerEvent, err error) {
	m.getByIDEvent = event
	m.getByIDErr = err
}

// SetEventByID sets a specific event to be returned for a given ID.
func (m *TestMockRepository) SetEventByID(eventID string, event *career.CareerEvent, err error) {
	if event != nil {
		m.getByIDEvents[eventID] = event
	}
	if err != nil {
		m.getByIDErrors[eventID] = err
	}
}

// SetListBehavior sets the events and error for List calls.
func (m *TestMockRepository) SetListBehavior(events []*career.CareerEvent, err error) {
	m.listEvents = events
	m.listErr = err
}

// SetCountBehavior sets the count and error for Count calls.
func (m *TestMockRepository) SetCountBehavior(count int, err error) {
	m.countResult = count
	m.countErr = err
}

// Create implements Repository interface.
func (m *TestMockRepository) Create(_ context.Context, _ *career.CareerEvent) error {
	m.createCalled = true
	return m.createErr
}

// Update implements Repository interface.
func (m *TestMockRepository) Update(_ context.Context, _ *career.CareerEvent) error {
	m.updateCalled = true
	return m.updateErr
}

// Delete implements Repository interface.
func (m *TestMockRepository) Delete(_ context.Context, _ string) error {
	m.deleteCalled = true
	return m.deleteErr
}

// GetByID implements Repository interface.
func (m *TestMockRepository) GetByID(_ context.Context, eventID string) (*career.CareerEvent, error) {
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

// List implements Repository interface.
func (m *TestMockRepository) List(_ context.Context, _ repo.ListFilters) ([]*career.CareerEvent, error) {
	m.listCalled = true
	return m.listEvents, m.listErr
}

// Count implements Repository interface.
func (m *TestMockRepository) Count(_ context.Context, _ repo.ListFilters) (int, error) {
	m.countCalled = true
	return m.countResult, m.countErr
}

// CreateCalled returns whether Create was called.
func (m *TestMockRepository) CreateCalled() bool {
	return m.createCalled
}

// UpdateCalled returns whether Update was called.
func (m *TestMockRepository) UpdateCalled() bool {
	return m.updateCalled
}

// DeleteCalled returns whether Delete was called.
func (m *TestMockRepository) DeleteCalled() bool {
	return m.deleteCalled
}

// GetByIDCalled returns whether GetByID was called.
func (m *TestMockRepository) GetByIDCalled() bool {
	return m.getByIDCalled
}

// ListCalled returns whether List was called.
func (m *TestMockRepository) ListCalled() bool {
	return m.listCalled
}

// CountCalled returns whether Count was called.
func (m *TestMockRepository) CountCalled() bool {
	return m.countCalled
}
