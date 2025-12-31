package burst_fact

import (
	"context"
	"fmt"
	"sort"
	"strings"
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
		Name:            bd.generateBurstName(cluster),
		Description:     bd.generateBurstDescription(cluster),
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

// generateBurstName creates a meaningful name for a burst based on the events
func (bd *BurstDetector) generateBurstName(cluster []career.CareerEvent) string {
	if len(cluster) == 0 {
		return "Unnamed Burst"
	}

	// Extract common themes from event texts
	commonWords := bd.extractCommonWords(cluster)
	if len(commonWords) > 0 {
		// Use the most common meaningful word
		return fmt.Sprintf("%s Initiative", commonWords[0])
	}

	// Fallback to project-based naming
	projects := bd.extractProjects(cluster)
	if len(projects) > 0 {
		return fmt.Sprintf("%s Project", projects[0])
	}

	// Final fallback
	return fmt.Sprintf("%d-Event Burst", len(cluster))
}

// generateBurstDescription creates a description for a burst based on the events
func (bd *BurstDetector) generateBurstDescription(cluster []career.CareerEvent) string {
	if len(cluster) == 0 {
		return "A collection of related career events"
	}

	projects := bd.extractProjects(cluster)
	companies := bd.extractCompanies(cluster)

	var parts []string

	if len(projects) > 0 {
		if len(projects) == 1 {
			parts = append(parts, fmt.Sprintf("Related to %s project", projects[0]))
		} else {
			parts = append(parts, fmt.Sprintf("Spanning %s and related projects", projects[0]))
		}
	}

	if len(companies) > 0 {
		if len(companies) == 1 {
			parts = append(parts, fmt.Sprintf("at %s", companies[0]))
		} else {
			parts = append(parts, fmt.Sprintf("across %s and other organizations", companies[0]))
		}
	}

	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}

	return fmt.Sprintf("A burst of %d related career events", len(cluster))
}

// extractCommonWords finds common meaningful words across event texts
func (bd *BurstDetector) extractCommonWords(cluster []career.CareerEvent) []string {
	wordCount := make(map[string]int)

	for _, event := range cluster {
		words := strings.Fields(strings.ToLower(event.Text))
		for _, word := range words {
			// Filter out common words and focus on meaningful terms
			if bd.isMeaningfulWord(word) {
				wordCount[word]++
			}
		}
	}

	// Find words that appear in multiple events
	var commonWords []string
	for word, count := range wordCount {
		if count > 1 || (len(cluster) <= 2 && count >= 1) {
			commonWords = append(commonWords, strings.Title(word))
		}
	}

	// Sort by frequency (most common first)
	sort.Slice(commonWords, func(i, j int) bool {
		return wordCount[strings.ToLower(commonWords[i])] > wordCount[strings.ToLower(commonWords[j])]
	})

	return commonWords
}

// extractProjects extracts unique project names from the cluster
func (bd *BurstDetector) extractProjects(cluster []career.CareerEvent) []string {
	projectSet := make(map[string]bool)
	var projects []string

	for _, event := range cluster {
		if event.Project != "" {
			if !projectSet[event.Project] {
				projectSet[event.Project] = true
				projects = append(projects, event.Project)
			}
		}
	}

	return projects
}

// extractCompanies extracts unique company names from the cluster
func (bd *BurstDetector) extractCompanies(cluster []career.CareerEvent) []string {
	companySet := make(map[string]bool)
	var companies []string

	for _, event := range cluster {
		if event.Company != "" {
			if !companySet[event.Company] {
				companySet[event.Company] = true
				companies = append(companies, event.Company)
			}
		}
	}

	return companies
}

// isMeaningfulWord filters out common words to focus on meaningful terms
func (bd *BurstDetector) isMeaningfulWord(word string) bool {
	// Remove punctuation
	word = strings.Trim(word, ".,!?;:")

	// Skip very short words
	if len(word) < 3 {
		return false
	}

	// Skip common English words
	commonWords := map[string]bool{
		"the": true, "and": true, "for": true, "are": true, "but": true,
		"not": true, "you": true, "all": true, "can": true, "her": true,
		"was": true, "one": true, "our": true, "had": true, "have": true,
		"has": true, "will": true, "been": true, "this": true, "that": true,
		"with": true, "from": true, "they": true, "know": true, "want": true,
		"good": true, "much": true, "some": true, "time": true,
		"very": true, "when": true, "come": true, "here": true, "just": true,
		"like": true, "long": true, "make": true, "many": true, "over": true,
		"such": true, "take": true, "than": true, "them": true, "well": true,
		"were": true,
	}

	return !commonWords[word]
}

