package burst_fact

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// BurstSuggestion represents a suggested burst with confidence score and related events
type BurstSuggestion struct {
	EventIDs        []string  // IDs of related events
	ConfidenceScore float64   // 0.0-1.0 confidence score
	Name            string    // Suggested burst name (optional)
	Description     string    // Suggested burst description (optional)
}

// DetectionOptions controls burst detection behavior
type DetectionOptions struct {
	MinConfidence       float64       // Minimum confidence score to suggest burst (default 0.6)
	TemporalWindow      time.Duration // Time window for grouping (default 6 months)
	MinEventCount       int           // Minimum events for burst (default 2)
	MaxSuggestionsCount int           // Maximum suggestions to return (default 10)
}

// BurstDetector detects related events and suggests bursts
type BurstDetector struct {
	similarityScorer *SimilarityScorer
	temporalGrouper  *TemporalGrouper
}

// NewBurstDetector creates a new burst detector
func NewBurstDetector() *BurstDetector {
	return &BurstDetector{
		similarityScorer: NewSimilarityScorer(),
		temporalGrouper:  NewTemporalGrouper(),
	}
}

// DetectBursts finds related event groups and returns burst suggestions
// Algorithm:
// 1. Group events by temporal proximity (6-month window)
// 2. Within each temporal group, score event pairs for similarity
// 3. Build clusters of similar events (connected components)
// 4. Calculate confidence score for each cluster
// 5. Filter by minimum confidence and return suggestions
func (bd *BurstDetector) DetectBursts(
	ctx context.Context,
	events []career.CareerEvent,
	opts *DetectionOptions,
) ([]BurstSuggestion, error) {
	if opts == nil {
		opts = &DetectionOptions{
			MinConfidence:       0.6,
			TemporalWindow:      6 * 30 * 24 * time.Hour, // ~6 months
			MinEventCount:       2,
			MaxSuggestionsCount: 10,
		}
	}

	if len(events) < opts.MinEventCount {
		return []BurstSuggestion{}, nil
	}

	// Step 1: Group events by temporal proximity using dates
	dates := make([]time.Time, len(events))
	eventsByDate := make(map[time.Time][]career.CareerEvent)
	for i, event := range events {
		dates[i] = event.Date
		eventsByDate[event.Date] = append(eventsByDate[event.Date], event)
	}

	temporalDateGroups := bd.temporalGrouper.GroupEventsByTemporal(dates)

	// Step 2: Build similarity matrix within each temporal group and find clusters
	suggestions := []BurstSuggestion{}
	for _, dateGroup := range temporalDateGroups {
		// Reconstruct event group from dates
		var eventGroup []career.CareerEvent
		for _, date := range dateGroup {
			eventGroup = append(eventGroup, eventsByDate[date]...)
		}

		if len(eventGroup) < opts.MinEventCount {
			continue
		}

		// Build similarity scores between all pairs
		similarityMatrix := bd.buildSimilarityMatrix(eventGroup)

		// Step 3: Find clusters using connected components
		clusters := bd.findClusters(eventGroup, similarityMatrix, opts.MinConfidence)

		// Step 4: Convert clusters to suggestions
		for _, cluster := range clusters {
			if len(cluster) >= opts.MinEventCount {
				suggestion := bd.clusterToSuggestion(cluster, similarityMatrix)
				suggestions = append(suggestions, suggestion)
			}
		}
	}

	// Step 5: Sort by confidence and limit results
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].ConfidenceScore > suggestions[j].ConfidenceScore
	})

	if len(suggestions) > opts.MaxSuggestionsCount {
		suggestions = suggestions[:opts.MaxSuggestionsCount]
	}

	return suggestions, nil
}

// buildSimilarityMatrix creates a matrix of similarity scores between events
func (bd *BurstDetector) buildSimilarityMatrix(
	events []career.CareerEvent,
) map[string]map[string]float64 {
	matrix := make(map[string]map[string]float64)

	for i := range events {
		matrix[events[i].ID] = make(map[string]float64)
		for j := range events {
			if i == j {
				matrix[events[i].ID][events[j].ID] = 1.0
			} else if j > i {
				// Only compute upper triangle to save computation
				input1 := EventSimilarityInput{
					Text:     events[i].Text,
					Keywords: events[i].Tags,
					Company:  events[i].Company,
					Project:  events[i].Project,
				}
				input2 := EventSimilarityInput{
					Text:     events[j].Text,
					Keywords: events[j].Tags,
					Company:  events[j].Company,
					Project:  events[j].Project,
				}
				score := bd.similarityScorer.CombinedSimilarityScore(input1, input2)
				matrix[events[i].ID][events[j].ID] = score
				// Mirror for symmetric access
				if _, ok := matrix[events[j].ID]; !ok {
					matrix[events[j].ID] = make(map[string]float64)
				}
				matrix[events[j].ID][events[i].ID] = score
			}
		}
	}

	return matrix
}

// findClusters identifies groups of similar events using connected components
// An edge exists between two events if their similarity score >= minConfidence
func (bd *BurstDetector) findClusters(
	events []career.CareerEvent,
	similarityMatrix map[string]map[string]float64,
	minConfidence float64,
) [][]career.CareerEvent {
	visited := make(map[string]bool)
	clusters := [][]career.CareerEvent{}

	for _, event := range events {
		if visited[event.ID] {
			continue
		}

		// BFS to find connected component
		cluster := []career.CareerEvent{}
		queue := []string{event.ID}
		visited[event.ID] = true

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			// Find original event
			var currentEvent *career.CareerEvent
			for i := range events {
				if events[i].ID == current {
					currentEvent = &events[i]
					break
				}
			}
			if currentEvent != nil {
				cluster = append(cluster, *currentEvent)
			}

			// Add neighbors with sufficient similarity
			for _, neighbor := range events {
				if !visited[neighbor.ID] && similarityMatrix[current][neighbor.ID] >= minConfidence {
					visited[neighbor.ID] = true
					queue = append(queue, neighbor.ID)
				}
			}
		}

		if len(cluster) >= 2 { // Only keep clusters with 2+ events
			clusters = append(clusters, cluster)
		}
	}

	return clusters
}

// clusterToSuggestion converts an event cluster to a burst suggestion
func (bd *BurstDetector) clusterToSuggestion(
	cluster []career.CareerEvent,
	similarityMatrix map[string]map[string]float64,
) BurstSuggestion {
	// Extract event IDs
	eventIDs := make([]string, len(cluster))
	for i, event := range cluster {
		eventIDs[i] = event.ID
	}

	// Calculate average pairwise similarity as confidence
	totalScore := 0.0
	pairCount := 0
	for i := range cluster {
		for j := i + 1; j < len(cluster); j++ {
			totalScore += similarityMatrix[cluster[i].ID][cluster[j].ID]
			pairCount++
		}
	}

	confidenceScore := 0.5 // default if no pairs
	if pairCount > 0 {
		confidenceScore = totalScore / float64(pairCount)
	}

	// Clamp confidence to [0, 1]
	if confidenceScore > 1.0 {
		confidenceScore = 1.0
	}
	if confidenceScore < 0.0 {
		confidenceScore = 0.0
	}

	return BurstSuggestion{
		EventIDs:        eventIDs,
		ConfidenceScore: confidenceScore,
	}
}

// ValidateSuggestion checks if a burst suggestion is valid
func (bd *BurstDetector) ValidateSuggestion(suggestion BurstSuggestion) error {
	if len(suggestion.EventIDs) < 2 {
		return fmt.Errorf("burst must have at least 2 events")
	}

	if suggestion.ConfidenceScore < 0.0 || suggestion.ConfidenceScore > 1.0 {
		return fmt.Errorf("confidence score must be between 0.0 and 1.0")
	}

	// Check for duplicate event IDs
	seen := make(map[string]bool)
	for _, id := range suggestion.EventIDs {
		if seen[id] {
			return fmt.Errorf("duplicate event ID in suggestion: %s", id)
		}
		seen[id] = true
	}

	return nil
}

