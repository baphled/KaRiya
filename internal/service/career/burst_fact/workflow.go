package burst_fact

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// RejectionRecord tracks a rejected burst suggestion to prevent re-suggesting
type RejectionRecord struct {
	EventIDs   []string  // Sorted event IDs that were rejected
	RejectedAt time.Time // When the rejection occurred
}

// RejectionTracker maintains a record of rejected burst suggestions
type RejectionTracker struct {
	mu         sync.RWMutex
	rejections map[string]RejectionRecord // Key is sorted event IDs joined
}

// NewRejectionTracker creates a new rejection tracker
func NewRejectionTracker() *RejectionTracker {
	return &RejectionTracker{
		rejections: make(map[string]RejectionRecord),
	}
}

// RecordRejection marks a burst suggestion as rejected
func (rt *RejectionTracker) RecordRejection(eventIDs []string) {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	sorted := make([]string, len(eventIDs))
	copy(sorted, eventIDs)
	sort.Strings(sorted)

	key := rt.generateKey(sorted)
	rt.rejections[key] = RejectionRecord{
		EventIDs:   sorted,
		RejectedAt: time.Now(),
	}
}

// IsRejected checks if a burst suggestion has been rejected
func (rt *RejectionTracker) IsRejected(eventIDs []string) bool {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	sorted := make([]string, len(eventIDs))
	copy(sorted, eventIDs)
	sort.Strings(sorted)

	key := rt.generateKey(sorted)
	_, exists := rt.rejections[key]
	return exists
}

// FilterRejected removes rejected suggestions from a list
func (rt *RejectionTracker) FilterRejected(suggestions []BurstSuggestion) []BurstSuggestion {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	var filtered []BurstSuggestion
	for _, suggestion := range suggestions {
		if !rt.isRejectedLocked(suggestion.EventIDs) {
			filtered = append(filtered, suggestion)
		}
	}
	return filtered
}

// isRejectedLocked checks if a burst is rejected (must be called with lock held)
func (rt *RejectionTracker) isRejectedLocked(eventIDs []string) bool {
	sorted := make([]string, len(eventIDs))
	copy(sorted, eventIDs)
	sort.Strings(sorted)

	key := rt.generateKey(sorted)
	_, exists := rt.rejections[key]
	return exists
}

// generateKey creates a consistent key from event IDs
func (rt *RejectionTracker) generateKey(eventIDs []string) string {
	key := ""
	for i, id := range eventIDs {
		if i > 0 {
			key += ","
		}
		key += id
	}
	return key
}

// ClearRejections removes all rejection records
func (rt *RejectionTracker) ClearRejections() {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.rejections = make(map[string]RejectionRecord)
}

// GetRejectionCount returns the number of tracked rejections
func (rt *RejectionTracker) GetRejectionCount() int {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	return len(rt.rejections)
}

// BurstConfirmationWorkflow orchestrates the burst suggestion and confirmation workflow
type BurstConfirmationWorkflow struct {
	detector           *BurstDetector
	rejectionTracker   *RejectionTracker
	confirmedBursts    []BurstSuggestion
	pendingSuggestions []BurstSuggestion
	mu                 sync.RWMutex
}

// NewBurstConfirmationWorkflow creates a new confirmation workflow
func NewBurstConfirmationWorkflow(detector *BurstDetector) *BurstConfirmationWorkflow {
	return &BurstConfirmationWorkflow{
		detector:         detector,
		rejectionTracker: NewRejectionTracker(),
		confirmedBursts:  []BurstSuggestion{},
	}
}

// GenerateSuggestions generates burst suggestions and filters out rejected ones
func (bcw *BurstConfirmationWorkflow) GenerateSuggestions(
	ctx context.Context,
	events []career.CareerEvent,
	opts *DetectionOptions,
) ([]BurstSuggestion, error) {
	bcw.mu.Lock()
	defer bcw.mu.Unlock()

	// Detect bursts
	suggestions, err := bcw.detector.DetectBursts(ctx, events, opts)
	if err != nil {
		return nil, err
	}

	// Filter out rejected suggestions
	filtered := bcw.rejectionTracker.FilterRejected(suggestions)

	// Store pending suggestions
	bcw.pendingSuggestions = filtered

	return filtered, nil
}

// ConfirmSuggestion confirms a burst suggestion
func (bcw *BurstConfirmationWorkflow) ConfirmSuggestion(suggestion BurstSuggestion) error {
	bcw.mu.Lock()
	defer bcw.mu.Unlock()

	// Validate suggestion
	if err := bcw.detector.ValidateSuggestion(suggestion); err != nil {
		return fmt.Errorf("invalid suggestion: %w", err)
	}

	// Remove from pending
	bcw.removePendingSuggestionLocked(suggestion.EventIDs)

	// Add to confirmed
	bcw.confirmedBursts = append(bcw.confirmedBursts, suggestion)

	return nil
}

// RejectSuggestion rejects a burst suggestion
func (bcw *BurstConfirmationWorkflow) RejectSuggestion(eventIDs []string) error {
	bcw.mu.Lock()
	defer bcw.mu.Unlock()

	if len(eventIDs) < 2 {
		return fmt.Errorf("cannot reject suggestion with fewer than 2 events")
	}

	// Record rejection
	bcw.rejectionTracker.RecordRejection(eventIDs)

	// Remove from pending
	bcw.removePendingSuggestionLocked(eventIDs)

	return nil
}

// removePendingSuggestionLocked removes a suggestion from pending list (must be called with lock held)
func (bcw *BurstConfirmationWorkflow) removePendingSuggestionLocked(eventIDs []string) {
	var filtered []BurstSuggestion
	for _, s := range bcw.pendingSuggestions {
		if !eventIDsEqual(s.EventIDs, eventIDs) {
			filtered = append(filtered, s)
		}
	}
	bcw.pendingSuggestions = filtered
}

// GetPendingSuggestions returns the current list of pending suggestions
func (bcw *BurstConfirmationWorkflow) GetPendingSuggestions() []BurstSuggestion {
	bcw.mu.RLock()
	defer bcw.mu.RUnlock()

	result := make([]BurstSuggestion, len(bcw.pendingSuggestions))
	copy(result, bcw.pendingSuggestions)
	return result
}

// GetConfirmedBursts returns all confirmed bursts in this workflow
func (bcw *BurstConfirmationWorkflow) GetConfirmedBursts() []BurstSuggestion {
	bcw.mu.RLock()
	defer bcw.mu.RUnlock()

	result := make([]BurstSuggestion, len(bcw.confirmedBursts))
	copy(result, bcw.confirmedBursts)
	return result
}

// Reset clears all workflow state
func (bcw *BurstConfirmationWorkflow) Reset() {
	bcw.mu.Lock()
	defer bcw.mu.Unlock()

	bcw.pendingSuggestions = []BurstSuggestion{}
	bcw.confirmedBursts = []BurstSuggestion{}
	bcw.rejectionTracker.ClearRejections()
}

// eventIDsEqual compares two event ID lists
func eventIDsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sortedA := make([]string, len(a))
	copy(sortedA, a)
	sort.Strings(sortedA)

	sortedB := make([]string, len(b))
	copy(sortedB, b)
	sort.Strings(sortedB)

	for i := range sortedA {
		if sortedA[i] != sortedB[i] {
			return false
		}
	}
	return true
}
