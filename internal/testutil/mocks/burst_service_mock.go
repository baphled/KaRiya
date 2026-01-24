// Package mocks provides test mocks for KaRiya services.
package mocks

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
)

// BurstServiceMock provides a configurable mock for BurstService interface.
type BurstServiceMock struct {
	events           []*career.CareerEvent
	facts            map[string][]*career.Fact
	suggestions      []burst_fact.BurstSuggestion
	suggestError     error
	extractedFacts   []career.Fact
	extractError     error
	savedFacts       []*career.Fact
	listEventsError  error
	extractCallCount int
	saveFactError    error
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
func (m *BurstServiceMock) SetSuggestions(suggestions []burst_fact.BurstSuggestion) *BurstServiceMock {
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
func (m *BurstServiceMock) ListEvents(_ context.Context, _ careerrepo.ListFilters) ([]*career.CareerEvent, error) {
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
func (m *BurstServiceMock) SuggestBursts(_ context.Context, _ []string) ([]burst_fact.BurstSuggestion, error) {
	if m.suggestError != nil {
		return nil, m.suggestError
	}
	return m.suggestions, nil
}
