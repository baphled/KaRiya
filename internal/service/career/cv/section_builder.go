package cv

import (
	"context"
	"sort"
	"strings"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
)

// SkillsFormatConfig configures how the skills section is formatted
type SkillsFormatConfig struct {
	Format               string   // "flat" or "grouped"
	Limit                int      // max skills per section/group (0 = no limit)
	SelectedTechnologies []string // skill IDs to prioritize
}

// skillInfo holds skill data for formatting
type skillInfo struct {
	ID       string
	Name     string
	Category string
}

// SectionBuilder organizes CV bullets into logical sections
type SectionBuilder interface {
	// BuildSections organizes bullets into CV sections
	// skillsConfig: optional configuration for skills section formatting (Phase 11 - Task 40)
	BuildSections(
		ctx context.Context,
		bullets []*career.CVBullet,
		events []*career.CareerEvent,
		facts []*career.Fact,
		targetRole string,
		skillsConfig *SkillsFormatConfig,
	) ([]*career.CVSection, error)
}

// DefaultSectionBuilder is the default implementation of SectionBuilder
type DefaultSectionBuilder struct {
	skillRepo careerrepo.SkillRepository
	logger    *logger.Logger
}

// NewSectionBuilder creates a new SectionBuilder instance
func NewSectionBuilder(skillRepo careerrepo.SkillRepository, log *logger.Logger) *DefaultSectionBuilder {
	return &DefaultSectionBuilder{
		skillRepo: skillRepo,
		logger:    log,
	}
}

// BuildSections organizes bullets into CV sections
func (sb *DefaultSectionBuilder) BuildSections(ctx context.Context, bullets []*career.CVBullet, events []*career.CareerEvent, facts []*career.Fact, targetRole string, skillsConfig *SkillsFormatConfig) ([]*career.CVSection, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var sections []*career.CVSection
	order := 0

	// 1. Summary FIRST (prose)
	if sb.shouldIncludeSummary(targetRole) && len(bullets) > 0 {
		summarySection := sb.buildSummarySection(bullets, order)
		if summarySection != nil {
			sections = append(sections, summarySection)
			order++
		}
	}

	// 2. Experience SECOND (company-based, with dates)
	experienceSection := sb.buildExperienceSection(bullets, events, order, targetRole)
	if experienceSection != nil {
		sections = append(sections, experienceSection)
		order++
	}

	// 3. Projects THIRD (project-based, with dates)
	projectsSection := sb.buildProjectsSection(bullets, events, order, targetRole)
	if projectsSection != nil {
		sections = append(sections, projectsSection)
		order++
	}

	// 4. Technical Skills LAST (from event skills, prioritizing selected technologies)
	if skillsConfig == nil {
		skillsConfig = &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: []string{}}
	}
	skillsSection := sb.buildSkillsSection(ctx, events, skillsConfig, order)
	if skillsSection != nil {
		sections = append(sections, skillsSection)
	}

	sb.logger.Info("Built %d CV sections from %d bullets, %d events, %d facts", len(sections), len(bullets), len(events), len(facts))
	return sections, nil
}

// buildExperienceSection creates the experience section
func (sb *DefaultSectionBuilder) buildExperienceSection(
	bullets []*career.CVBullet, events []*career.CareerEvent,
	order int, targetRole string,
) *career.CVSection {
	if len(bullets) == 0 {
		sb.logger.Info("buildExperienceSection: no bullets provided")
		return nil
	}

	// Group bullets by company (only events WITH company)
	groups := sb.groupBulletsByCompany(bullets, events)

	sb.logger.Info("buildExperienceSection: found %d company groups from %d bullets", len(groups), len(bullets))

	if len(groups) == 0 {
		return nil
	}

	// Sort groups by most recent endDate first
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].endDate.After(groups[j].endDate)
	})

	// Get role-specific bullet cap per company
	maxBulletsPerCompany := sb.getBulletsPerCompanyForRole(targetRole)

	// Convert to SectionContentGroup array and apply per-company bullet cap
	content := make([]*career.SectionContentGroup, 0, len(groups))
	for _, group := range groups {
		bullets := group.bullets
		// Apply per-company bullet cap (bullets are already ranked, so just take the first N)
		if len(bullets) > maxBulletsPerCompany {
			bullets = bullets[:maxBulletsPerCompany]
		}
		content = append(content, &career.SectionContentGroup{
			Header:    group.header,
			StartDate: formatMonthYear(group.startDate),
			EndDate:   formatMonthYear(group.endDate),
			Bullets:   bullets,
		})
	}

	return &career.CVSection{
		ID:          uuid.New().String(),
		SectionType: "experience",
		Title:       "Experience",
		Order:       order,
		Content:     content,
	}
}

// buildProjectsSection creates the projects section
func (sb *DefaultSectionBuilder) buildProjectsSection(
	bullets []*career.CVBullet, events []*career.CareerEvent,
	order int, targetRole string,
) *career.CVSection {
	if len(bullets) == 0 {
		return nil
	}

	// Group bullets by project (only events WITHOUT company but WITH project)
	groups := sb.groupBulletsByProject(bullets, events)

	if len(groups) == 0 {
		return nil
	}

	// Sort groups by most recent endDate first
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].endDate.After(groups[j].endDate)
	})

	// Get role-specific bullet cap per project (same as companies)
	maxBulletsPerProject := sb.getBulletsPerCompanyForRole(targetRole)

	// Convert to SectionContentGroup array and apply per-project bullet cap
	content := make([]*career.SectionContentGroup, 0, len(groups))
	for _, group := range groups {
		bullets := group.bullets
		// Apply per-project bullet cap (bullets are already ranked, so just take the first N)
		if len(bullets) > maxBulletsPerProject {
			bullets = bullets[:maxBulletsPerProject]
		}
		content = append(content, &career.SectionContentGroup{
			Header:  group.header,
			Bullets: bullets,
			// Note: Projects don't have dates (StartDate/EndDate intentionally omitted)
		})
	}

	return &career.CVSection{
		ID:          uuid.New().String(),
		SectionType: "projects",
		Title:       "Projects",
		Order:       order,
		Content:     content,
	}
}

// buildSkillsSection creates the technical skills section from event skills (Phase 11 - Task 40)
// Looks up skill names from IDs, supports flat/grouped formatting with limits
func (sb *DefaultSectionBuilder) buildSkillsSection(ctx context.Context, events []*career.CareerEvent, config *SkillsFormatConfig, order int) *career.CVSection {
	// Collect unique skill IDs from all events
	skillIDSet := make(map[string]bool)
	for _, event := range events {
		for _, skillID := range event.Skills {
			if skillID != "" {
				skillIDSet[skillID] = true
			}
		}
	}

	if len(skillIDSet) == 0 {
		return nil
	}

	// Convert to slice for lookup
	skillIDs := make([]string, 0, len(skillIDSet))
	for id := range skillIDSet {
		skillIDs = append(skillIDs, id)
	}

	// Look up skills from repository (get full Skill objects for category info)
	skills := make([]skillInfo, 0, len(skillIDs))
	if sb.skillRepo != nil {
		for _, skillID := range skillIDs {
			skill, err := sb.skillRepo.GetByID(ctx, skillID)
			if err == nil && skill != nil {
				skills = append(skills, skillInfo{
					ID:       skill.ID,
					Name:     skill.Name,
					Category: skill.Category,
				})
			} else {
				// Fallback: use ID if lookup fails
				skills = append(skills, skillInfo{
					ID:       skillID,
					Name:     skillID,
					Category: "",
				})
			}
		}
	} else {
		// No repository: use IDs as names (shouldn't happen in production)
		for _, skillID := range skillIDs {
			skills = append(skills, skillInfo{
				ID:       skillID,
				Name:     skillID,
				Category: "",
			})
		}
	}

	// Sort skills: selected technologies first, then alphabetically
	selectedSet := make(map[string]bool)
	for _, techID := range config.SelectedTechnologies {
		selectedSet[techID] = true
	}

	sort.Slice(skills, func(i, j int) bool {
		iSelected := selectedSet[skills[i].ID]
		jSelected := selectedSet[skills[j].ID]

		// Selected technologies come first
		if iSelected != jSelected {
			return iSelected
		}

		// Alphabetical order by name
		return skills[i].Name < skills[j].Name
	})

	// Format based on config
	var content []*career.SectionContentGroup

	if config.Format == "grouped" {
		// Grouped format: limit is per-category
		content = sb.buildGroupedSkills(skills, config.Limit)
	} else {
		// Flat format: limit is total across all skills
		if config.Limit > 0 && len(skills) > config.Limit {
			skills = skills[:config.Limit]
		}
		content = sb.buildFlatSkills(skills)
	}

	return &career.CVSection{
		ID:          uuid.New().String(),
		SectionType: "skills",
		Title:       "Technical Skills",
		Order:       order,
		Content:     content,
	}
}

// buildFlatSkills creates a flat list of skills (one bullet per skill)
func (sb *DefaultSectionBuilder) buildFlatSkills(skills []skillInfo) []*career.SectionContentGroup {
	bullets := make([]*career.CVBullet, 0, len(skills))
	for _, skill := range skills {
		bullets = append(bullets, &career.CVBullet{
			ID:   uuid.New().String(),
			Text: skill.Name,
		})
	}

	return []*career.SectionContentGroup{
		{
			Header:  "",
			Bullets: bullets,
		},
	}
}

// buildGroupedSkills creates grouped skills by category (comma-separated per group)
func (sb *DefaultSectionBuilder) buildGroupedSkills(skills []skillInfo, limitPerGroup int) []*career.SectionContentGroup {
	// Group skills by category
	categoryMap := make(map[string][]string)
	for _, skill := range skills {
		category := skill.Category
		if category == "" {
			category = "Other"
		}
		categoryMap[category] = append(categoryMap[category], skill.Name)
	}

	// Sort categories alphabetically
	categories := make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	// Create one group per category, with category as header
	groups := make([]*career.SectionContentGroup, 0, len(categories))
	for _, category := range categories {
		skillNames := categoryMap[category]

		// Apply per-group limit if specified
		if limitPerGroup > 0 && len(skillNames) > limitPerGroup {
			skillNames = skillNames[:limitPerGroup]
		}

		// Create bullets for each skill in this category
		bullets := make([]*career.CVBullet, 0, len(skillNames))
		for _, skillName := range skillNames {
			bullets = append(bullets, &career.CVBullet{
				ID:   uuid.New().String(),
				Text: skillName,
			})
		}

		// Add group with category as header
		groups = append(groups, &career.SectionContentGroup{
			Header:  category,
			Bullets: bullets,
		})
	}

	return groups
}

// buildSummarySection creates a brief professional summary
func (sb *DefaultSectionBuilder) buildSummarySection(bullets []*career.CVBullet, order int) *career.CVSection {
	if len(bullets) == 0 {
		return nil
	}

	// Use top 2 bullets to create summary prose
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
		Summary:     content.String(),
		Content:     nil, // No content groups for summary
	}
}

// groupBulletsByCompany groups bullets by company from source events.
// It detects separate tenures when events at OTHER companies exist between
// two periods at the same company (BUG-009 fix).
func (sb *DefaultSectionBuilder) groupBulletsByCompany(bullets []*career.CVBullet, events []*career.CareerEvent) []*bulletGroup {
	// Create map of event ID to event.
	eventMap := make(map[string]*career.CareerEvent)
	for _, event := range events {
		eventMap[event.ID] = event
	}

	// Sort ALL events chronologically (oldest first) for tenure detection.
	sortedEvents := make([]*career.CareerEvent, len(events))
	copy(sortedEvents, events)
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].Date.Before(sortedEvents[j].Date)
	})

	// First pass: Associate each bullet with its primary company and earliest event date.
	var bulletInfos []bulletInfo
	skippedCount := 0

	for _, bullet := range bullets {
		companyCounts := make(map[string]int)
		var earliestDate, latestDate time.Time
		var earliestEventID string

		for _, eventID := range bullet.SourceEventIDs {
			if event, exists := eventMap[eventID]; exists {
				if event.Company != "" {
					companyCounts[event.Company]++
					if earliestDate.IsZero() || event.Date.Before(earliestDate) {
						earliestDate = event.Date
						earliestEventID = event.ID
					}
					if latestDate.IsZero() || event.Date.After(latestDate) {
						latestDate = event.Date
					}
				}
			}
		}

		// Determine primary company (most frequent).
		primaryCompany := ""
		maxCount := 0
		for company, count := range companyCounts {
			if count > maxCount {
				maxCount = count
				primaryCompany = company
			}
		}

		if primaryCompany == "" {
			skippedCount++
			continue
		}

		bulletInfos = append(bulletInfos, bulletInfo{
			bullet:        bullet,
			company:       primaryCompany,
			earliestDate:  earliestDate,
			latestDate:    latestDate,
			sourceEventID: earliestEventID,
		})
	}

	// Sort bullet infos by date (oldest first) for tenure detection.
	sort.Slice(bulletInfos, func(i, j int) bool {
		return bulletInfos[i].earliestDate.Before(bulletInfos[j].earliestDate)
	})

	// Group bullets by company, maintaining order.
	bulletsByCompany := make(map[string][]bulletInfo)
	for _, bi := range bulletInfos {
		bulletsByCompany[bi.company] = append(bulletsByCompany[bi.company], bi)
	}

	// Detect tenures for each company and create groups.
	var groupSlice []*bulletGroup

	for company, companyBullets := range bulletsByCompany {
		if len(companyBullets) == 0 {
			continue
		}

		// Detect tenure boundaries.
		tenures := sb.detectBulletTenures(companyBullets, company, sortedEvents)

		// Create a group for each tenure.
		for _, tenure := range tenures {
			group := &bulletGroup{
				header:  company,
				bullets: make([]*career.CVBullet, 0, len(tenure)),
			}

			for _, bi := range tenure {
				group.bullets = append(group.bullets, bi.bullet)

				// Update date range.
				if group.startDate.IsZero() || bi.earliestDate.Before(group.startDate) {
					group.startDate = bi.earliestDate
				}
				if group.endDate.IsZero() || bi.latestDate.After(group.endDate) {
					group.endDate = bi.latestDate
				}
			}

			groupSlice = append(groupSlice, group)
		}
	}

	sb.logger.Info("groupBulletsByCompany: created %d groups, skipped %d bullets without company", len(groupSlice), skippedCount)
	return groupSlice
}

// detectBulletTenures splits a company's bullets into separate tenure groups.
// A new tenure is detected when events at OTHER companies fall between
// two consecutive bullets at this company (BUG-009).
func (sb *DefaultSectionBuilder) detectBulletTenures(
	companyBullets []bulletInfo,
	company string,
	allEventsSorted []*career.CareerEvent,
) [][]bulletInfo {
	if len(companyBullets) <= 1 {
		return [][]bulletInfo{companyBullets}
	}

	var tenures [][]bulletInfo
	currentTenure := []bulletInfo{companyBullets[0]}

	for i := 1; i < len(companyBullets); i++ {
		prevBullet := companyBullets[i-1]
		currBullet := companyBullets[i]

		// Check if there are events at OTHER companies between these two dates.
		hasIntervening := sb.hasInterveningCompanyEvents(
			prevBullet.latestDate,
			currBullet.earliestDate,
			company,
			allEventsSorted,
		)

		if hasIntervening {
			// Start a new tenure.
			tenures = append(tenures, currentTenure)
			currentTenure = []bulletInfo{currBullet}
		} else {
			// Continue current tenure.
			currentTenure = append(currentTenure, currBullet)
		}
	}

	// Add the last tenure.
	tenures = append(tenures, currentTenure)

	return tenures
}

// hasInterveningCompanyEvents checks if any events at OTHER companies
// exist between two dates (exclusive of both endpoints).
func (sb *DefaultSectionBuilder) hasInterveningCompanyEvents(
	startDate, endDate time.Time,
	currentCompany string,
	allEventsSorted []*career.CareerEvent,
) bool {
	for _, event := range allEventsSorted {
		eventCompany := event.Company
		if eventCompany == "" {
			continue
		}

		// Skip events at the same company.
		if eventCompany == currentCompany {
			continue
		}

		// Check if this event falls between the two dates (exclusive).
		if event.Date.After(startDate) && event.Date.Before(endDate) {
			return true
		}
	}
	return false
}

// groupBulletsByProject groups bullets by project from source events (where Company is empty)
func (sb *DefaultSectionBuilder) groupBulletsByProject(bullets []*career.CVBullet, events []*career.CareerEvent) []*bulletGroup {
	// Create map of event ID to event
	eventMap := make(map[string]*career.CareerEvent)
	for _, event := range events {
		eventMap[event.ID] = event
	}

	// Group bullets by project
	groups := make(map[string]*bulletGroup)
	for _, bullet := range bullets {
		// Count projects from source events (where Company is empty but Project is not)
		projectCounts := make(map[string]int)
		var eventDates []time.Time

		for _, eventID := range bullet.SourceEventIDs {
			if event, exists := eventMap[eventID]; exists {
				// Only count events WITHOUT company but WITH project
				if event.Company == "" && event.Project != "" {
					projectCounts[event.Project]++
					eventDates = append(eventDates, event.Date)
				}
			}
		}

		// Determine primary project (most frequent)
		primaryProject := ""
		maxCount := 0
		for project, count := range projectCounts {
			if count > maxCount {
				maxCount = count
				primaryProject = project
			}
		}

		// Skip bullets with no project
		if primaryProject == "" {
			continue
		}

		// Initialize group if needed
		if _, exists := groups[primaryProject]; !exists {
			groups[primaryProject] = &bulletGroup{
				header:  primaryProject,
				bullets: make([]*career.CVBullet, 0),
			}
		}

		// Add bullet to group and update date range
		group := groups[primaryProject]
		group.bullets = append(group.bullets, bullet)

		// Update date range
		for _, date := range eventDates {
			if group.startDate.IsZero() || date.Before(group.startDate) {
				group.startDate = date
			}
			if group.endDate.IsZero() || date.After(group.endDate) {
				group.endDate = date
			}
		}
	}

	// Convert to slice
	var groupSlice []*bulletGroup
	for _, group := range groups {
		groupSlice = append(groupSlice, group)
	}

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
	header    string    // Company or Project name
	startDate time.Time // Earliest event date in group
	endDate   time.Time // Latest event date in group
	bullets   []*career.CVBullet
}

// bulletInfo holds information about a bullet for tenure detection (BUG-009).
type bulletInfo struct {
	bullet        *career.CVBullet
	company       string
	earliestDate  time.Time
	latestDate    time.Time
	sourceEventID string
}

// formatMonthYear formats a time.Time as "Jan 2006"
func formatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2006")
}

// isCommonWord checks if a word is a common word to skip
func isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"and":     true,
		"the":     true,
		"for":     true,
		"with":    true,
		"from":    true,
		"that":    true,
		"this":    true,
		"have":    true,
		"been":    true,
		"were":    true,
		"will":    true,
		"your":    true,
		"their":   true,
		"would":   true,
		"about":   true,
		"which":   true,
		"when":    true,
		"make":    true,
		"like":    true,
		"time":    true,
		"just":    true,
		"know":    true,
		"take":    true,
		"people":  true,
		"into":    true,
		"year":    true,
		"could":   true,
		"them":    true,
		"some":    true,
		"than":    true,
		"then":    true,
		"now":     true,
		"look":    true,
		"only":    true,
		"come":    true,
		"its":     true,
		"over":    true,
		"think":   true,
		"also":    true,
		"back":    true,
		"after":   true,
		"use":     true,
		"two":     true,
		"how":     true,
		"our":     true,
		"work":    true,
		"first":   true,
		"well":    true,
		"way":     true,
		"even":    true,
		"new":     true,
		"want":    true,
		"because": true,
		"any":     true,
		"these":   true,
		"give":    true,
		"day":     true,
		"most":    true,
		"does":    true,
		"very":    true,
		"through": true,
		"being":   true,
		"each":    true,
		"much":    true,
		"made":    true,
		"many":    true,
		"must":    true,
		"before":  true,
		"such":    true,
		"where":   true,
		"those":   true,
		"both":    true,
		"during":  true,
		"same":    true,
		"until":   true,
		"while":   true,
		"too":     true,
		"try":     true,
	}
	return commonWords[word]
}

// getBulletsPerCompanyForRole returns the maximum bullets per company for a role
// Based on documented role-specific caps:
// - Principal: 3-4 bullets max
// - Staff: 4-5 bullets max
// - EM: 3-4 bullets max
// - Senior IC: 4-5 bullets max
func (sb *DefaultSectionBuilder) getBulletsPerCompanyForRole(targetRole string) int {
	switch strings.ToLower(targetRole) {
	case "principal":
		return 4 // 3-4 per docs, using upper bound
	case "staff":
		return 5 // 4-5 per docs, using upper bound
	case "em":
		return 4 // 3-4 per docs, using upper bound
	case "senior_ic":
		return 5 // 4-5 per docs, using upper bound
	default:
		return 4 // Conservative default
	}
}
