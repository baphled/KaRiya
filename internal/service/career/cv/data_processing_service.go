package cv

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/google/uuid"
)

// DataProcessingService handles extraction and organization of career data
type DataProcessingService interface {
	// GroupEventsByCompany organizes events into company-based structure
	GroupEventsByCompany(ctx context.Context, events []*career.CareerEvent) (map[string]*CompanyGroup, error)

	// ExtractAchievements identifies key achievements from events and facts
	ExtractAchievements(ctx context.Context, event *career.CareerEvent, facts []*career.Fact) ([]*Achievement, error)

	// ExtractSkills identifies skills from events and facts
	ExtractSkills(ctx context.Context, events []*career.CareerEvent, facts []*career.Fact) (map[string]*SkillCategory, error)

	// CalculateMetrics extracts quantifiable metrics from event descriptions
	CalculateMetrics(ctx context.Context, text string) ([]*Metric, error)

	// ExtractProjectsFromEvents identifies unique projects within events
	ExtractProjectsFromEvents(ctx context.Context, events []*career.CareerEvent) ([]*ProjectGroup, error)
}

// DefaultDataProcessingService is the default implementation
type DefaultDataProcessingService struct {
	logger *logger.Logger
}

// NewDataProcessingService creates a new data processing service
func NewDataProcessingService(log *logger.Logger) DataProcessingService {
	return &DefaultDataProcessingService{
		logger: log,
	}
}

// CompanyGroup represents work at a specific company
type CompanyGroup struct {
	ID           string
	Company      string
	Position     string
	StartDate    time.Time
	EndDate      time.Time
	Description  string
	Projects     []*ProjectGroup
	Achievements []*Achievement
	Skills       []*Skill
	EventIDs     []string // Reference to source events
}

// ProjectGroup represents a specific project
type ProjectGroup struct {
	ID           string
	Name         string
	Description  string
	StartDate    time.Time
	EndDate      time.Time
	Role         string
	Achievements []*Achievement
	Skills       []*Skill
	EventIDs     []string // Reference to source events
}

// Achievement represents a measurable accomplishment
type Achievement struct {
	ID          string
	Description string
	Metrics     []*Metric
	EventID     string
	FactIDs     []string
	Confidence  float64
	ActionVerb  string
	Category    constants.CompetencyCategory // Primary competency category (BUG-008)
}

// Metric represents a quantifiable measure
type Metric struct {
	Type    string // "percentage", "count", "currency", "time", "ratio"
	Value   string
	Unit    string
	Context string
}

// Skill represents a professional capability
type Skill struct {
	Name         string
	Level        string // "beginner", "intermediate", "advanced", "expert"
	Projects     int
	Endorsements int
	Categories   []string
}

// SkillCategory groups related skills
type SkillCategory struct {
	Name   string
	Skills []*Skill
}

// GroupEventsByCompany organizes events into company-based structure.
// It detects separate tenures when events at OTHER companies exist between
// two periods at the same company (BUG-009 fix).
func (svc *DefaultDataProcessingService) GroupEventsByCompany(
	ctx context.Context, events []*career.CareerEvent,
) (map[string]*CompanyGroup, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(events) == 0 {
		return make(map[string]*CompanyGroup), nil
	}

	groups := make(map[string]*CompanyGroup)

	// Sort ALL events chronologically (oldest first) for tenure detection.
	sortedEvents := make([]*career.CareerEvent, len(events))
	copy(sortedEvents, events)
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].Date.Before(sortedEvents[j].Date)
	})

	// Group events by company.
	eventsByCompany := make(map[string][]*career.CareerEvent)
	for _, event := range events {
		company := event.Company
		if company == "" {
			company = constants.DefaultCompanyName
		}
		eventsByCompany[company] = append(eventsByCompany[company], event)
	}

	// Process each company group with tenure detection.
	for company, companyEvents := range eventsByCompany {
		// Detect separate tenures for this company.
		tenures := svc.detectTenures(companyEvents, sortedEvents)

		// Create a group for each tenure.
		for tenureIndex, tenureEvents := range tenures {
			// Sort tenure events by date (newest first) for display.
			sort.Slice(tenureEvents, func(i, j int) bool {
				return tenureEvents[i].Date.After(tenureEvents[j].Date)
			})

			// Determine position (use most common or first).
			position := svc.extractPosition(tenureEvents)

			// Calculate date range for this tenure.
			startDate, endDate := svc.calculateDateRange(tenureEvents)

			// Extract projects.
			projects := svc.extractProjects(tenureEvents)

			// Create unique key for multiple tenures: "Company A" or "Company A#2".
			groupKey := company
			if len(tenures) > 1 {
				groupKey = fmt.Sprintf("%s%s%d", company, constants.TenureSeparator, tenureIndex+1)
			}

			// Create company group.
			group := &CompanyGroup{
				ID:        uuid.New().String(),
				Company:   company,
				Position:  position,
				StartDate: startDate,
				EndDate:   endDate,
				Projects:  projects,
				Skills:    []*Skill{},
				EventIDs:  extractEventIDs(tenureEvents),
			}

			groups[groupKey] = group
		}
	}

	svc.logger.Info("Grouped %d events into %d companies/tenures", len(events), len(groups))

	return groups, nil
}

// detectTenures splits a company's events into separate tenure groups.
// A new tenure is detected when events at OTHER companies fall between
// two consecutive events at this company (BUG-009).
func (svc *DefaultDataProcessingService) detectTenures(
	companyEvents []*career.CareerEvent,
	allEventsSorted []*career.CareerEvent,
) [][]*career.CareerEvent {
	if len(companyEvents) <= 1 {
		return [][]*career.CareerEvent{companyEvents}
	}

	// Get the company name from first event.
	company := companyEvents[0].Company
	if company == "" {
		company = constants.DefaultCompanyName
	}

	// Sort company events chronologically (oldest first).
	sortedCompanyEvents := make([]*career.CareerEvent, len(companyEvents))
	copy(sortedCompanyEvents, companyEvents)
	sort.Slice(sortedCompanyEvents, func(i, j int) bool {
		return sortedCompanyEvents[i].Date.Before(sortedCompanyEvents[j].Date)
	})

	// Detect tenure boundaries by checking for intervening work at other companies.
	var tenures [][]*career.CareerEvent
	currentTenure := []*career.CareerEvent{sortedCompanyEvents[0]}

	for i := 1; i < len(sortedCompanyEvents); i++ {
		prevEvent := sortedCompanyEvents[i-1]
		currEvent := sortedCompanyEvents[i]

		// Check if there are events at OTHER companies between these two dates.
		hasIntervening := hasInterveningCompanyEvents(
			prevEvent.Date,
			currEvent.Date,
			company,
			allEventsSorted,
		)

		if hasIntervening {
			// Start a new tenure.
			tenures = append(tenures, currentTenure)
			currentTenure = []*career.CareerEvent{currEvent}
		} else {
			// Continue current tenure.
			currentTenure = append(currentTenure, currEvent)
		}
	}

	// Add the last tenure.
	tenures = append(tenures, currentTenure)

	return tenures
}

// hasInterveningCompanyEvents checks if any events at OTHER companies
// exist between two dates (exclusive of both endpoints).
// This is a shared utility used by both data processing and section builder.
func hasInterveningCompanyEvents(
	startDate, endDate time.Time,
	currentCompany string,
	allEventsSorted []*career.CareerEvent,
) bool {
	for _, event := range allEventsSorted {
		// BUG-012: Skip project-only events (empty Company) — they run
		// concurrently with employment and are not intervening work.
		if event.Company == "" {
			continue
		}

		// Skip events at the same company.
		if event.Company == currentCompany {
			continue
		}

		// Check if this event falls between the two dates (exclusive).
		if event.Date.After(startDate) && event.Date.Before(endDate) {
			return true
		}
	}
	return false
}

// ExtractAchievements identifies key achievements from events and facts
func (svc *DefaultDataProcessingService) ExtractAchievements(
	ctx context.Context, event *career.CareerEvent, facts []*career.Fact,
) ([]*Achievement, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var achievements []*Achievement

	// Extract from event text
	baseAchievement := &Achievement{
		ID:          uuid.New().String(),
		Description: svc.enhanceBulletText(event.Text),
		EventID:     event.ID,
		Confidence:  0.8,
		ActionVerb:  svc.extractActionVerb(event.Text),
		Category:    svc.extractPrimaryCategory(event.Categories), // BUG-008
	}

	// Extract metrics from event text
	metrics, err := svc.CalculateMetrics(ctx, event.Text)
	if err != nil {
		svc.logger.Warn("Failed to extract metrics: %v", err)
	}
	baseAchievement.Metrics = metrics

	achievements = append(achievements, baseAchievement)

	// Extract from related facts
	for _, fact := range facts {
		if fact.SourceEventID == event.ID {
			achievement := &Achievement{
				ID:          uuid.New().String(),
				Description: fact.Text,
				EventID:     event.ID,
				FactIDs:     []string{fact.ID},
				Confidence:  0.9, // Facts are higher confidence
				ActionVerb:  svc.extractActionVerb(fact.Text),
				Category:    svc.extractPrimaryCategory(fact.CompetencyCategories), // BUG-008
			}
			achievements = append(achievements, achievement)
		}
	}

	return achievements, nil
}

// ExtractSkills identifies skills from events and facts
func (svc *DefaultDataProcessingService) ExtractSkills(
	ctx context.Context, events []*career.CareerEvent, facts []*career.Fact,
) (map[string]*SkillCategory, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	skillMap := make(map[string][]*Skill)

	// Extract skills from event tags and categories
	for _, event := range events {
		for _, tag := range event.Tags {
			skill := &Skill{
				Name:       tag,
				Level:      svc.determineSkillLevel(tag),
				Projects:   1,
				Categories: []string{tag},
			}
			skillMap[tag] = append(skillMap[tag], skill)
		}
	}

	// Extract skills from facts
	for _, fact := range facts {
		for _, category := range fact.CompetencyCategories {
			// Use fact text as skill name if not too long
			skillName := svc.extractSkillNameFromFact(fact.Text)
			if skillName != "" {
				skill := &Skill{
					Name:         skillName,
					Level:        svc.determineSkillLevelFromFact(fact),
					Endorsements: 1,
					Categories:   []string{category},
				}
				skillMap[skillName] = append(skillMap[skillName], skill)
			}
		}
	}

	// Aggregate and deduplicate skills
	categories := make(map[string]*SkillCategory)
	for _, skillList := range skillMap {
		// Merge duplicate skills
		merged := svc.mergeSkills(skillList)

		// Determine category
		category := svc.determineSkillCategory(merged)
		if _, exists := categories[category]; !exists {
			categories[category] = &SkillCategory{
				Name:   category,
				Skills: []*Skill{},
			}
		}
		categories[category].Skills = append(categories[category].Skills, merged)
	}

	svc.logger.Info("Extracted %d skill categories from events and facts", len(categories))

	return categories, nil
}

// CalculateMetrics extracts quantifiable metrics from text
func (svc *DefaultDataProcessingService) CalculateMetrics(ctx context.Context, text string) ([]*Metric, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var metrics []*Metric

	// Percentage pattern: "25%" or "increased by 25 percent"
	percentPattern := regexp.MustCompile(`(\d+\.?\d*)\s*%|(\d+\.?\d*)\s*percent`)
	percentMatches := percentPattern.FindAllStringSubmatch(text, -1)
	for _, match := range percentMatches {
		value := match[1]
		if value == "" {
			value = match[2]
		}
		metrics = append(metrics, &Metric{
			Type:    "percentage",
			Value:   value,
			Unit:    "%",
			Context: extractContext(text, match[0]),
		})
	}

	// Count pattern: "12 people", "500 customers", "3 projects"
	countPattern := regexp.MustCompile(`(\d+)\s+(people|customers|users|clients|projects|teams|departments)`)
	countMatches := countPattern.FindAllStringSubmatch(text, -1)
	for _, match := range countMatches {
		metrics = append(metrics, &Metric{
			Type:    "count",
			Value:   match[1],
			Unit:    match[2],
			Context: extractContext(text, match[0]),
		})
	}

	// Currency pattern: "$1M", "$50,000", "£100k"
	currencyPattern := regexp.MustCompile(`[\$£€][\d,]+(?:\.?\d{2})?[KMB]?`)
	currencyMatches := currencyPattern.FindAllString(text, -1)
	for _, match := range currencyMatches {
		metrics = append(metrics, &Metric{
			Type:    "currency",
			Value:   match,
			Unit:    "currency",
			Context: extractContext(text, match),
		})
	}

	// Time pattern: "6 months", "3 years", "2 weeks"
	timePattern := regexp.MustCompile(`(\d+)\s+(months?|years?|weeks?|days?)`)
	timeMatches := timePattern.FindAllStringSubmatch(text, -1)
	for _, match := range timeMatches {
		metrics = append(metrics, &Metric{
			Type:    "time",
			Value:   match[1],
			Unit:    match[2],
			Context: extractContext(text, match[0]),
		})
	}

	// Ratio pattern: "3x", "10x faster"
	ratioPattern := regexp.MustCompile(`(\d+\.?\d*)x(?:\s+faster)?`)
	ratioMatches := ratioPattern.FindAllStringSubmatch(text, -1)
	for _, match := range ratioMatches {
		metrics = append(metrics, &Metric{
			Type:    "ratio",
			Value:   match[1],
			Unit:    "x",
			Context: extractContext(text, match[0]),
		})
	}

	return metrics, nil
}

// ExtractProjectsFromEvents identifies unique projects within events
func (svc *DefaultDataProcessingService) ExtractProjectsFromEvents(
	ctx context.Context, events []*career.CareerEvent,
) ([]*ProjectGroup, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	projectMap := make(map[string][]*career.CareerEvent)

	// Group events by project
	for _, event := range events {
		project := event.Project
		if project == "" {
			project = "General"
		}
		projectMap[project] = append(projectMap[project], event)
	}

	var projects []*ProjectGroup

	// Create project groups
	for projectName, projectEvents := range projectMap {
		sort.Slice(projectEvents, func(i, j int) bool {
			return projectEvents[i].Date.After(projectEvents[j].Date)
		})

		startDate, endDate := svc.calculateDateRange(projectEvents)

		project := &ProjectGroup{
			ID:        uuid.New().String(),
			Name:      projectName,
			StartDate: startDate,
			EndDate:   endDate,
			Skills:    []*Skill{},
			EventIDs:  extractEventIDs(projectEvents),
		}

		projects = append(projects, project)
	}

	// Sort by most recent
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].EndDate.After(projects[j].EndDate)
	})

	return projects, nil
}

// Helper functions

// extractPosition determines the position from events
func (svc *DefaultDataProcessingService) extractPosition(events []*career.CareerEvent) string {
	// Count position mentions in event text
	positionCount := make(map[string]int)

	for _, event := range events {
		// Look for common position words
		text := strings.ToLower(event.Text)
		if strings.Contains(text, "engineer") {
			positionCount["Engineer"]++
		}
		if strings.Contains(text, "manager") {
			positionCount["Manager"]++
		}
		if strings.Contains(text, "lead") {
			positionCount["Lead"]++
		}
		if strings.Contains(text, "architect") {
			positionCount["Architect"]++
		}
		if strings.Contains(text, "director") {
			positionCount["Director"]++
		}
	}

	// Return most common position
	maxCount := 0
	position := "Professional"
	for pos, count := range positionCount {
		if count > maxCount {
			maxCount = count
			position = pos
		}
	}

	return position
}

// calculateDateRange determines start and end dates
func (svc *DefaultDataProcessingService) calculateDateRange(events []*career.CareerEvent) (time.Time, time.Time) {
	if len(events) == 0 {
		return time.Now(), time.Now()
	}

	startDate := events[0].Date
	endDate := events[0].Date

	for _, event := range events {
		if event.Date.Before(startDate) {
			startDate = event.Date
		}
		if event.Date.After(endDate) {
			endDate = event.Date
		}
	}

	return startDate, endDate
}

// extractProjects creates project groups from events
func (svc *DefaultDataProcessingService) extractProjects(events []*career.CareerEvent) []*ProjectGroup {
	projectMap := make(map[string][]*career.CareerEvent)

	for _, event := range events {
		project := event.Project
		if project == "" {
			continue
		}
		projectMap[project] = append(projectMap[project], event)
	}

	var projects []*ProjectGroup
	for projectName, projectEvents := range projectMap {
		startDate, endDate := svc.calculateDateRange(projectEvents)
		projects = append(projects, &ProjectGroup{
			ID:        uuid.New().String(),
			Name:      projectName,
			StartDate: startDate,
			EndDate:   endDate,
			Skills:    []*Skill{},
			EventIDs:  extractEventIDs(projectEvents),
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].EndDate.After(projects[j].EndDate)
	})

	return projects
}

// extractEventIDs gets IDs from events
func extractEventIDs(events []*career.CareerEvent) []string {
	var ids []string
	for _, event := range events {
		ids = append(ids, event.ID)
	}
	return ids
}

// enhanceBulletText improves bullet point wording
func (svc *DefaultDataProcessingService) enhanceBulletText(text string) string {
	// Remove weak phrases
	weakPhrases := []string{"worked on", "was involved in", "helped with", "participated in"}
	result := text
	for _, phrase := range weakPhrases {
		result = strings.ReplaceAll(result, phrase, "")
	}

	// Capitalize first letter
	if result != "" {
		result = strings.ToUpper(result[:1]) + result[1:]
	}

	return strings.TrimSpace(result)
}

// extractActionVerb extracts the primary action verb
func (svc *DefaultDataProcessingService) extractActionVerb(text string) string {
	actionVerbs := []string{
		"led", "managed", "developed", "designed", "implemented", "built",
		"created", "established", "launched", "improved", "optimized",
		"increased", "decreased", "reduced", "accelerated", "streamlined",
		"architected", "engineered", "delivered", "deployed", "automated",
		"coordinated", "directed", "mentored", "coached", "trained",
		"analyzed", "evaluated", "assessed", "reviewed", "audited",
	}

	lowerText := strings.ToLower(text)
	for _, verb := range actionVerbs {
		if strings.Contains(lowerText, verb) {
			return verb
		}
	}

	return "accomplished"
}

// determineSkillLevel determines skill level from tag
func (svc *DefaultDataProcessingService) determineSkillLevel(tag string) string {
	switch strings.ToLower(tag) {
	case "technical":
		return "expert"
	case "leadership":
		return "advanced"
	case "product":
		return "advanced"
	default:
		return "intermediate"
	}
}

// determineSkillLevelFromFact determines skill level from fact
func (svc *DefaultDataProcessingService) determineSkillLevelFromFact(fact *career.Fact) string {
	switch fact.RoleFit {
	case career.RoleFitPrincipal:
		return "expert"
	case career.RoleFitStaff:
		return "advanced"
	case career.RoleFitSeniorIC:
		return "advanced"
	case career.RoleFitEM:
		return "advanced"
	default:
		return "intermediate"
	}
}

// extractSkillNameFromFact extracts a skill name from fact text
func (svc *DefaultDataProcessingService) extractSkillNameFromFact(text string) string {
	// Extract first few words as skill name
	words := strings.Fields(text)
	if len(words) > 5 {
		return strings.Join(words[:5], " ")
	}
	return text
}

// mergeSkills combines duplicate skills
func (svc *DefaultDataProcessingService) mergeSkills(skills []*Skill) *Skill {
	if len(skills) == 0 {
		return &Skill{}
	}

	merged := &Skill{
		Name:       skills[0].Name,
		Level:      skills[0].Level,
		Projects:   0,
		Categories: []string{},
	}

	for _, skill := range skills {
		merged.Projects += skill.Projects
		merged.Endorsements += skill.Endorsements
		merged.Categories = append(merged.Categories, skill.Categories...)
	}

	// Remove duplicate categories
	merged.Categories = removeDuplicates(merged.Categories)

	return merged
}

// determineSkillCategory determines the category for a skill
func (svc *DefaultDataProcessingService) determineSkillCategory(skill *Skill) string {
	if len(skill.Categories) > 0 {
		return skill.Categories[0]
	}

	// Default categories based on skill name
	lowerName := strings.ToLower(skill.Name)
	if strings.Contains(lowerName, "technical") || strings.Contains(lowerName, "engineer") {
		return "Technical"
	}
	if strings.Contains(lowerName, "leadership") || strings.Contains(lowerName, "manage") {
		return "Leadership"
	}
	if strings.Contains(lowerName, "product") {
		return "Product"
	}

	return "Other"
}

// extractContext extracts surrounding context from text
func extractContext(text string, match string) string {
	idx := strings.Index(text, match)
	if idx == -1 {
		return ""
	}

	start := idx - 20
	if start < 0 {
		start = 0
	}
	end := idx + len(match) + 20
	if end > len(text) {
		end = len(text)
	}

	return strings.TrimSpace(text[start:end])
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// extractPrimaryCategory extracts the primary (first) category from a list of categories
// and converts it to the type-safe CompetencyCategory constant (BUG-008)
func (svc *DefaultDataProcessingService) extractPrimaryCategory(categories []string) constants.CompetencyCategory {
	if len(categories) == 0 {
		return ""
	}
	// Use first category as primary
	primary := strings.ToLower(categories[0])
	if constants.IsValidCompetencyCategory(primary) {
		return constants.CompetencyCategory(primary)
	}
	return ""
}
