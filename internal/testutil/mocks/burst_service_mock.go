// Package mocks provides test mocks for KaRiya services.
package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
)

// BurstServiceMock provides a configurable mock for BurstService interface.
type BurstServiceMock struct {
	events           []*career.CareerEvent
	facts            map[string][]*career.Fact
	suggestions      []burstfact.BurstSuggestion
	suggestError     error
	extractedFacts   []career.Fact
	extractError     error
	savedFacts       []*career.Fact
	listEventsError  error
	extractCallCount int
	saveFactError    error
	confirmError     error
	confirmCallCount int
}

// NewBurstServiceMock creates a new configurable BurstService mock.
func NewBurstServiceMock() *BurstServiceMock {
	return &BurstServiceMock{
		events: []*career.CareerEvent{},
		facts:  make(map[string][]*career.Fact),
	}
}

// SetEvents configures the events returned by GetEventByID and ListEvents.
func (m *BurstServiceMock) SetEvents(events []*career.CareerEvent) *BurstServiceMock {
	m.events = events
	return m
}

// SetSuggestions configures the suggestions returned by SuggestBursts.
func (m *BurstServiceMock) SetSuggestions(suggestions []burstfact.BurstSuggestion) *BurstServiceMock {
	m.suggestions = suggestions
	return m
}

// SetSuggestError configures an error to be returned by SuggestBursts.
func (m *BurstServiceMock) SetSuggestError(err error) *BurstServiceMock {
	m.suggestError = err
	return m
}

// SetExtractedFacts configures the facts returned by ExtractFactsFromBurst.
func (m *BurstServiceMock) SetExtractedFacts(facts []career.Fact) *BurstServiceMock {
	m.extractedFacts = facts
	return m
}

// SetExtractError configures an error to be returned by ExtractFactsFromBurst.
func (m *BurstServiceMock) SetExtractError(err error) *BurstServiceMock {
	m.extractError = err
	return m
}

// SetListEventsError configures an error to be returned by ListEvents.
func (m *BurstServiceMock) SetListEventsError(err error) *BurstServiceMock {
	m.listEventsError = err
	return m
}

// SetSaveFactError configures an error to be returned by SaveFact.
func (m *BurstServiceMock) SetSaveFactError(err error) *BurstServiceMock {
	m.saveFactError = err
	return m
}

// SetConfirmError configures an error to be returned by ConfirmBurst.
func (m *BurstServiceMock) SetConfirmError(err error) *BurstServiceMock {
	m.confirmError = err
	return m
}

// GetConfirmCallCount returns how many times ConfirmBurst was called.
func (m *BurstServiceMock) GetConfirmCallCount() int {
	return m.confirmCallCount
}

// SetFactsForBurst configures facts to be returned for a specific burst ID.
func (m *BurstServiceMock) SetFactsForBurst(burstID string, facts []*career.Fact) *BurstServiceMock {
	m.facts[burstID] = facts
	return m
}

// GetExtractCallCount returns how many times ExtractFactsFromBurst was called.
func (m *BurstServiceMock) GetExtractCallCount() int {
	return m.extractCallCount
}

// GetSavedFacts returns all facts that were saved via SaveFact.
func (m *BurstServiceMock) GetSavedFacts() []*career.Fact {
	return m.savedFacts
}

// ConfirmBurst implements BurstService.
// Mirrors the real service: sets Confirmed, ConfirmedAt, and UpdatedAt.
func (m *BurstServiceMock) ConfirmBurst(_ context.Context, burst *career.Burst) error {
	m.confirmCallCount++
	if m.confirmError != nil {
		return m.confirmError
	}
	if burst != nil {
		now := time.Now()
		burst.Confirmed = true
		burst.ConfirmedAt = &now
		burst.UpdatedAt = now
	}
	return nil
}

// GetFactsBySourceBurstID implements BurstService.
func (m *BurstServiceMock) GetFactsBySourceBurstID(_ context.Context, burstID string) ([]*career.Fact, error) {
	if facts, ok := m.facts[burstID]; ok {
		return facts, nil
	}
	return []*career.Fact{}, nil
}

// GetEventByID implements BurstService.
func (m *BurstServiceMock) GetEventByID(_ context.Context, id string) (*career.CareerEvent, error) {
	for _, e := range m.events {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, fmt.Errorf("event not found: %s", id)
}

// ListEvents implements BurstService.
func (m *BurstServiceMock) ListEvents(_ context.Context, _ careerrepo.EventListFilters) ([]*career.CareerEvent, error) {
	if m.listEventsError != nil {
		return nil, m.listEventsError
	}
	return m.events, nil
}

// ExtractFactsFromBurst implements BurstService.
func (m *BurstServiceMock) ExtractFactsFromBurst(_ context.Context, _ *career.Burst) ([]career.Fact, error) {
	m.extractCallCount++
	if m.extractError != nil {
		return nil, m.extractError
	}
	return m.extractedFacts, nil
}

// SaveFact implements BurstService.
func (m *BurstServiceMock) SaveFact(_ context.Context, fact *career.Fact) error {
	if m.saveFactError != nil {
		return m.saveFactError
	}
	m.savedFacts = append(m.savedFacts, fact)
	return nil
}

// SuggestBursts implements BurstService.
func (m *BurstServiceMock) SuggestBursts(_ context.Context, _ []string) ([]burstfact.BurstSuggestion, error) {
	if m.suggestError != nil {
		return nil, m.suggestError
	}
	return m.suggestions, nil
}

// BurstRepositoryMock provides a configurable mock for BurstRepository interface.
// Use this to simulate database failures in E2E tests.
type BurstRepositoryMock struct {
	bursts      map[string]*career.Burst
	createError error
	updateError error
	deleteError error
	getError    error
	listError   error
	createCalls int
	updateCalls int
	deleteCalls int
}

// NewBurstRepositoryMock creates a new configurable BurstRepository mock.
func NewBurstRepositoryMock() *BurstRepositoryMock {
	return &BurstRepositoryMock{
		bursts: make(map[string]*career.Burst),
	}
}

// SetCreateError configures an error to be returned by Create.
func (m *BurstRepositoryMock) SetCreateError(err error) *BurstRepositoryMock {
	m.createError = err
	return m
}

// SetUpdateError configures an error to be returned by Update.
func (m *BurstRepositoryMock) SetUpdateError(err error) *BurstRepositoryMock {
	m.updateError = err
	return m
}

// SetDeleteError configures an error to be returned by Delete.
func (m *BurstRepositoryMock) SetDeleteError(err error) *BurstRepositoryMock {
	m.deleteError = err
	return m
}

// SetGetError configures an error to be returned by GetByID.
func (m *BurstRepositoryMock) SetGetError(err error) *BurstRepositoryMock {
	m.getError = err
	return m
}

// SetListError configures an error to be returned by List.
func (m *BurstRepositoryMock) SetListError(err error) *BurstRepositoryMock {
	m.listError = err
	return m
}

// AddBurst adds a burst to the mock repository (for test setup).
func (m *BurstRepositoryMock) AddBurst(burst *career.Burst) *BurstRepositoryMock {
	if burst.ID == "" {
		burst.ID = fmt.Sprintf("burst-%d", len(m.bursts)+1)
	}
	m.bursts[burst.ID] = burst
	return m
}

// GetCreateCalls returns how many times Create was called.
func (m *BurstRepositoryMock) GetCreateCalls() int {
	return m.createCalls
}

// GetUpdateCalls returns how many times Update was called.
func (m *BurstRepositoryMock) GetUpdateCalls() int {
	return m.updateCalls
}

// GetDeleteCalls returns how many times Delete was called.
func (m *BurstRepositoryMock) GetDeleteCalls() int {
	return m.deleteCalls
}

// Create implements BurstRepository.
func (m *BurstRepositoryMock) Create(_ context.Context, burst *career.Burst) error {
	m.createCalls++
	if m.createError != nil {
		return m.createError
	}
	if burst.ID == "" {
		burst.ID = fmt.Sprintf("burst-%d", len(m.bursts)+1)
	}
	m.bursts[burst.ID] = burst
	return nil
}

// GetByID implements BurstRepository.
func (m *BurstRepositoryMock) GetByID(_ context.Context, id string) (*career.Burst, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	if burst, ok := m.bursts[id]; ok {
		return burst, nil
	}
	return nil, careerrepo.ErrBurstNotFound
}

// Update implements BurstRepository.
func (m *BurstRepositoryMock) Update(_ context.Context, burst *career.Burst) error {
	m.updateCalls++
	if m.updateError != nil {
		return m.updateError
	}
	if _, ok := m.bursts[burst.ID]; !ok {
		return careerrepo.ErrBurstNotFound
	}
	m.bursts[burst.ID] = burst
	return nil
}

// Delete implements BurstRepository.
func (m *BurstRepositoryMock) Delete(_ context.Context, id string) error {
	m.deleteCalls++
	if m.deleteError != nil {
		return m.deleteError
	}
	if _, ok := m.bursts[id]; !ok {
		return careerrepo.ErrBurstNotFound
	}
	delete(m.bursts, id)
	return nil
}

// List implements BurstRepository.
func (m *BurstRepositoryMock) List(_ context.Context, _ careerrepo.BurstListFilters) ([]*career.Burst, error) {
	if m.listError != nil {
		return nil, m.listError
	}
	result := make([]*career.Burst, 0, len(m.bursts))
	for _, b := range m.bursts {
		result = append(result, b)
	}
	return result, nil
}

// Count implements BurstRepository.
func (m *BurstRepositoryMock) Count(_ context.Context, _ careerrepo.BurstListFilters) (int, error) {
	if m.listError != nil {
		return 0, m.listError
	}
	return len(m.bursts), nil
}
