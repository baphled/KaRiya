package cv

import (
	"context"
	"sort"
	"strings"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/google/uuid"
)

// SectionBuilder organizes CV bullets into logical sections
type SectionBuilder interface {
	// BuildSections organizes bullets into CV sections
	BuildSections(ctx context.Context, bullets []*career.CVBullet, events []*career.CareerEvent, targetRole string) ([]*career.CVSection, error)
}

// DefaultSectionBuilder is the default implementation of SectionBuilder
type DefaultSectionBuilder struct {
	logger *logger.Logger
}

// NewSectionBuilder creates a new SectionBuilder instance
func NewSectionBuilder(log *logger.Logger) *DefaultSectionBuilder {
	return &DefaultSectionBuilder{
		logger: log,
	}
}

// BuildSections organizes bullets into CV sections
func (sb *DefaultSectionBuilder) BuildSections(ctx context.Context, bullets []*career.CVBullet, events []*career.CareerEvent, targetRole string) ([]*career.CVSection, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var sections []*career.CVSection
	order := 0

	// Build experience section (always first if we have bullets)
	if len(bullets) > 0 {
		experienceSection := sb.buildExperienceSection(bullets, events, order)
		if experienceSection != nil {
			sections = append(sections, experienceSection)
			order++
		}
	}

	// Build skills section (from fact-based bullets)
	skillsSection := sb.buildSkillsSection(bullets, order)
	if skillsSection != nil {
		sections = append(sections, skillsSection)
		order++
	}

	// Build summary section (optional, based on role fit)
	if sb.shouldIncludeSummary(targetRole) && len(bullets) > 0 {
		summarySection := sb.buildSummarySection(bullets, order)
		if summarySection != nil {
			sections = append(sections, summarySection)
		}
	}

	sb.logger.Info("Built %d CV sections from %d bullets", len(sections), len(bullets))
	return sections, nil
}

// buildExperienceSection creates the experience section
func (sb *DefaultSectionBuilder) buildExperienceSection(bullets []*career.CVBullet, events []*career.CareerEvent, order int) *career.CVSection {
	if len(bullets) == 0 {
		return nil
	}

	// Group bullets by company/project from source events
	groupedBullets := sb.groupBulletsByCompany(bullets, events)

	// Build content
	var content strings.Builder
	for _, group := range groupedBullets {
		if group.company != "" {
			content.WriteString("**" + group.company + "**\n")
		}
		for _, bullet := range group.bullets {
			content.WriteString("• " + bullet.Text + "\n")
		}
		content.WriteString("\n")
	}

	return &career.CVSection{
		ID:          uuid.New().String(),
		SectionType: "experience",
		Title:       "Experience",
		Order:       order,
		Content:     strings.TrimSpace(content.String()),
	}
}

// buildSkillsSection creates the skills/competencies section
func (sb *DefaultSectionBuilder) buildSkillsSection(bullets []*career.CVBullet, order int) *career.CVSection {
	// Filter for fact-based bullets (higher quality for skills section)
	factBullets := make([]*career.CVBullet, 0)
	for _, bullet := range bullets {
		if len(bullet.SourceFactIDs) > 0 {
			factBullets = append(factBullets, bullet)
		}
	}

	if len(factBullets) == 0 {
		return nil // Skip skills section if no facts
	}

	// Build content as bullet list
	var content strings.Builder
	for _, bullet := range factBullets {
		content.WriteString("• " + bullet.Text + "\n")
	}

	return &career.CVSection{
		ID:          uuid.New().String(),
		SectionType: "skills",
		Title:       "Core Competencies",
		Order:       order,
		Content:     strings.TrimSpace(content.String()),
	}
}

// buildSummarySection creates a brief professional summary
func (sb *DefaultSectionBuilder) buildSummarySection(bullets []*career.CVBullet, order int) *career.CVSection {
	if len(bullets) == 0 {
		return nil
	}

	// Use top 2 bullets to create summary
	summaryBullets := bullets
	if len(summaryBullets) > 2 {
		summaryBullets = summaryBullets[:2]
	}

	var content strings.Builder
	content.WriteString("Experienced professional specializing in ")

	// Extract key terms from top bullets
	keyTerms := sb.extractKeyTerms(summaryBullets)
	if len(keyTerms) > 0 {
		content.WriteString(strings.Join(keyTerms, ", "))
		content.WriteString(". ")
	}

	content.WriteString("Proven track record of delivering high-quality results.")

	return &career.CVSection{
		ID:          uuid.New().String(),
		SectionType: "summary",
		Title:       "Professional Summary",
		Order:       order,
		Content:     content.String(),
	}
}

// groupBulletsByCompany groups bullets by company from source events
func (sb *DefaultSectionBuilder) groupBulletsByCompany(bullets []*career.CVBullet, events []*career.CareerEvent) []*bulletGroup {
	// Create map of event ID to company
	eventCompanies := make(map[string]string)
	for _, event := range events {
		if event.Company != "" {
			eventCompanies[event.ID] = event.Company
		}
	}

	// Group bullets
	groups := make(map[string]*bulletGroup)
	for _, bullet := range bullets {
		company := ""
		if len(bullet.SourceEventIDs) > 0 {
			company = eventCompanies[bullet.SourceEventIDs[0]]
		}

		if _, exists := groups[company]; !exists {
			groups[company] = &bulletGroup{
				company: company,
				bullets: make([]*career.CVBullet, 0),
			}
		}
		groups[company].bullets = append(groups[company].bullets, bullet)
	}

	// Convert to slice and sort
	var groupSlice []*bulletGroup
	for _, group := range groups {
		groupSlice = append(groupSlice, group)
	}

	// Sort by company name, with empty company first
	sort.Slice(groupSlice, func(i, j int) bool {
		if groupSlice[i].company == "" {
			return true
		}
		if groupSlice[j].company == "" {
			return false
		}
		return groupSlice[i].company < groupSlice[j].company
	})

	return groupSlice
}

// extractKeyTerms extracts key terms from bullets for summary
func (sb *DefaultSectionBuilder) extractKeyTerms(bullets []*career.CVBullet) []string {
	terms := make(map[string]bool)

	for _, bullet := range bullets {
		// Extract nouns and important terms
		words := strings.Fields(bullet.Text)
		for i, word := range words {
			word = strings.ToLower(word)
			// Skip common words and short words
			if len(word) > 4 && !isCommonWord(word) {
				terms[word] = true
			}
			// Also look for two-word combinations
			if i < len(words)-1 {
				phrase := word + " " + strings.ToLower(words[i+1])
				if len(phrase) > 6 {
					terms[phrase] = true
				}
			}
		}
	}

	// Convert to slice and limit to 3 terms
	var result []string
	count := 0
	for term := range terms {
		if count >= 3 {
			break
		}
		result = append(result, term)
		count++
	}

	sort.Strings(result)
	return result
}

// shouldIncludeSummary determines if a summary section should be included
func (sb *DefaultSectionBuilder) shouldIncludeSummary(targetRole string) bool {
	// Include summary for higher-level roles
	switch strings.ToLower(targetRole) {
	case "principal", "staff":
		return true
	default:
		return false
	}
}

// bulletGroup represents a group of bullets under a company/project
type bulletGroup struct {
	company string
	bullets []*career.CVBullet
}

// isCommonWord checks if a word is a common word to skip
func isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"and":   true,
		"the":   true,
		"with":  true,
		"from":  true,
		"to":    true,
		"for":   true,
		"of":    true,
		"in":    true,
		"on":    true,
		"at":    true,
		"by":    true,
		"or":    true,
		"was":   true,
		"were":  true,
		"been":  true,
		"being": true,
		"have":  true,
		"has":   true,
		"had":   true,
		"do":    true,
		"does":  true,
		"did":   true,
		"will":  true,
		"would": true,
		"could": true,
		"should": true,
		"may":   true,
		"might": true,
		"must":  true,
		"can":   true,
		"about": true,
		"as":    true,
		"be":    true,
		"but":   true,
		"you":   true,
		"me":    true,
		"him":   true,
		"her":   true,
		"it":    true,
		"us":    true,
		"them":  true,
	}
	return commonWords[word]
}

