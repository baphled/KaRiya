package burstfact

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/google/uuid"
)

// Extractor handles fact extraction from events and bursts.
type Extractor struct {
	classifier *Classifier
}

// NewExtractor creates a new fact extractor.
//
// Expected:
//   - classifier must be valid.
//
// Returns:
//   - A fully initialized Extractor ready for use.
//
// Side effects:
//   - None.
func NewExtractor(classifier *Classifier) *Extractor {
	return &Extractor{
		classifier: classifier,
	}
}

// ExtractFromEvent extracts facts from a single career event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A []career.Fact value.
//
// Side effects:
//   - None.
func (e *Extractor) ExtractFromEvent(_ context.Context, event *career.Event) []career.Fact {
	if event == nil || event.Text == "" {
		return []career.Fact{}
	}

	// Extract primary fact from event
	facts := []career.Fact{}

	// Generate fact text (use event text as-is or summarize)
	factText := e.generateFactTextFromEvent(event)

	// Infer competencies
	competencies := e.inferCompetencies(event)
	if len(competencies) == 0 {
		competencies = []string{"technical"}
	}

	// Infer role fit from text AND event categories
	// Categories provide a strong signal for role classification
	roleFit := e.classifier.ClassifyRoleFitWithCategories(event.Text, event.Categories)

	// Infer audience relevance
	audiences := e.classifier.ClassifyAudienceRelevance(event.Text, roleFit)
	if len(audiences) == 0 {
		audiences = []string{"peer"}
	}

	// Extract strength signal
	strengthSignal := e.classifier.ExtractStrengthSignal(event.Text)

	// Create fact
	fact := career.Fact{
		ID:                   uuid.New().String(),
		Text:                 factText,
		CompetencyCategories: competencies,
		RoleFit:              roleFit,
		AudienceRelevance:    audiences,
		StrengthSignal:       strengthSignal,
		SourceEventID:        event.ID,
		CreatedAt:            event.CreatedAt,
		UpdatedAt:            event.UpdatedAt,
	}

	facts = append(facts, fact)

	// Extract tag-specific facts if tags provide additional context
	if len(event.Tags) > 0 {
		tagFacts := e.extractTagSpecificFacts(event)
		facts = append(facts, tagFacts...)
	}

	return facts
}

// ExtractFromBurst extracts facts from a burst (multiple related events).
//
// Expected:
//   - burst must be valid.
//   - event must be valid.
//
// Returns:
//   - A []career.Fact value.
//
// Side effects:
//   - None.
func (e *Extractor) ExtractFromBurst(ctx context.Context, burst *career.Burst, events []*career.Event) []career.Fact {
	if burst == nil || len(events) == 0 {
		return []career.Fact{}
	}

	facts := []career.Fact{}

	// Generate burst-level fact
	burstFactText := e.generateFactTextFromBurst(burst, events)
	burstCompetencies := e.inferBurstCompetencies(events)
	if len(burstCompetencies) == 0 {
		burstCompetencies = []string{"technical"}
	}

	// Aggregate categories from all burst events for role classification
	burstCategories := e.aggregateBurstCategories(events)
	burstRoleFit := e.classifier.ClassifyRoleFitWithCategories(burstFactText, burstCategories)
	burstAudiences := e.classifier.ClassifyAudienceRelevance(burstFactText, burstRoleFit)
	if len(burstAudiences) == 0 {
		burstAudiences = []string{"peer"}
	}

	burstStrengthSignal := e.classifier.ExtractStrengthSignal(burstFactText)

	burstFact := career.Fact{
		ID:                   uuid.New().String(),
		Text:                 burstFactText,
		CompetencyCategories: burstCompetencies,
		RoleFit:              burstRoleFit,
		AudienceRelevance:    burstAudiences,
		StrengthSignal:       burstStrengthSignal,
		SourceBurstID:        burst.ID,
		CreatedAt:            burst.CreatedAt,
		UpdatedAt:            burst.UpdatedAt,
	}

	facts = append(facts, burstFact)

	// Extract individual event facts for comprehensive coverage
	for _, event := range events {
		eventFacts := e.ExtractFromEvent(ctx, event)
		facts = append(facts, eventFacts...)
	}

	return facts
}

// generateFactTextFromEvent creates a fact statement from an event.
func (e *Extractor) generateFactTextFromEvent(event *career.Event) string {
	// If event text is already concise and statement-like, use it as-is
	if len(event.Text) <= 200 {
		return event.Text
	}

	// For longer event descriptions, create a summary
	words := strings.Fields(event.Text)
	if len(words) > 30 {
		// Take first 20-30 words to create a concise fact
		summary := strings.Join(words[:30], " ")
		return summary + "..."
	}

	return event.Text
}

// generateFactTextFromBurst creates a fact statement from a burst.
func (e *Extractor) generateFactTextFromBurst(burst *career.Burst, events []*career.Event) string {
	// Start with burst name/description if available
	if burst.Name != "" {
		return burst.Name
	}

	// Otherwise, synthesize from related events
	if len(events) > 0 {
		// Take themes from first event
		firstEvent := events[0]
		if firstEvent.Company != "" && firstEvent.Project != "" {
			return fmt.Sprintf("Led initiative on %s project at %s", firstEvent.Project, firstEvent.Company)
		}
		if firstEvent.Company != "" {
			return "Multiple achievements at " + firstEvent.Company
		}
	}

	return fmt.Sprintf("Related achievements (%d events)", len(events))
}

// inferCompetencies infers competency categories from event text and tags.
func (e *Extractor) inferCompetencies(event *career.Event) []string {
	competencies := make(map[string]bool)

	// Use existing event categories if available
	if len(event.Categories) > 0 {
		for _, cat := range event.Categories {
			competencies[cat] = true
		}
	}

	// Infer from tags
	tagCompetencies := e.classifier.InferCompetencies(event.Text, event.Tags)
	for _, comp := range tagCompetencies {
		competencies[comp] = true
	}

	// Infer from text content
	lowerText := strings.ToLower(event.Text)

	if strings.Contains(lowerText, "lead") || strings.Contains(lowerText, "manage") || strings.Contains(lowerText, "direct") {
		competencies["leadership"] = true
	}
	if strings.Contains(lowerText, "mentor") || strings.Contains(lowerText, "coach") || strings.Contains(lowerText, "train") {
		competencies["mentoring"] = true
	}
	if strings.Contains(lowerText, "product") || strings.Contains(lowerText, "feature") || strings.Contains(lowerText, "roadmap") {
		competencies["product"] = true
	}
	if strings.Contains(lowerText, "consult") || strings.Contains(lowerText, "advise") || strings.Contains(lowerText, "optimize") {
		competencies["consulting"] = true
	}
	if strings.Contains(lowerText, "research") || strings.Contains(lowerText, "analyze") || strings.Contains(lowerText, "investigate") {
		competencies["research"] = true
	}

	// Default to technical if no other competencies detected
	if len(competencies) == 0 {
		competencies["technical"] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(competencies))
	for comp := range competencies {
		result = append(result, comp)
	}

	return result
}

// inferBurstCompetencies infers competencies from multiple burst events.
func (e *Extractor) inferBurstCompetencies(events []*career.Event) []string {
	competencies := make(map[string]bool)

	// Collect competencies from all events
	for _, event := range events {
		eventComps := e.inferCompetencies(event)
		for _, comp := range eventComps {
			competencies[comp] = true
		}
	}

	// Determine burst competency focus (CompetencyFocus in burst)
	// For now, return most common competencies
	if len(competencies) == 0 {
		return []string{"technical"}
	}

	result := make([]string, 0, len(competencies))
	for comp := range competencies {
		result = append(result, comp)
	}

	return result
}

// extractTagSpecificFacts creates additional facts from event tags.
func (e *Extractor) extractTagSpecificFacts(_ *career.Event) []career.Fact {
	// Don't create duplicate facts for every tag, only meaningful ones.
	// For now, return empty - tag context already captured in main fact.
	return []career.Fact{}
}

// aggregateBurstCategories collects unique categories from all burst events.
func (e *Extractor) aggregateBurstCategories(events []*career.Event) []string {
	categories := make(map[string]bool)
	for _, event := range events {
		for _, cat := range event.Categories {
			categories[cat] = true
		}
	}

	result := make([]string, 0, len(categories))
	for cat := range categories {
		result = append(result, cat)
	}
	return result
}
