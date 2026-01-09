// Package e2e provides test fixtures for E2E testing.
package e2e

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// CreateSampleEvents generates test career events with varied data.
// Events are created with realistic data spanning multiple dates, companies, and categories.
func CreateSampleEvents(count int) []*career.CareerEvent {
	events := make([]*career.CareerEvent, count)

	companies := []string{"Acme Corp", "TechStart Inc", "BigCorp", "StartupXYZ", "MegaTech"}
	projects := []string{"Project Alpha", "Backend Refactor", "Customer Portal", "API Gateway", "Data Pipeline"}
	categories := [][]string{
		{"technical", "leadership"},
		{"leadership", "mentoring"},
		{"technical"},
		{"leadership", "product"},
		{"technical", "research"},
	}
	tags := [][]string{
		{"technical", "project"},
		{"leadership", "mentoring"},
		{"technical", "product"},
		{"leadership", "achievement"},
		{"technical", "research"},
	}

	eventTexts := []string{
		"Led the redesign of the authentication system, improving security and reducing login failures by 40%%",
		"Mentored three junior developers on best practices for code review and testing",
		"Implemented a new caching layer that reduced API response times by 60%%",
		"Facilitated weekly team retrospectives and introduced new process improvements",
		"Designed and implemented a microservices migration strategy for the legacy monolith",
		"Created comprehensive documentation for the deployment pipeline",
		"Led cross-team collaboration to deliver the Q4 release on schedule",
		"Optimized database queries resulting in 50%% reduction in page load times",
		"Introduced test-driven development practices across the engineering team",
		"Built a real-time notification system handling 10K messages per second",
	}

	now := time.Now()
	for i := 0; i < count; i++ {
		event := &career.CareerEvent{
			ID:         fmt.Sprintf("event_%d", i+1),
			Text:       eventTexts[i%len(eventTexts)],
			Date:       now.AddDate(0, 0, -i*7), // Each event 1 week apart
			Company:    companies[i%len(companies)],
			Project:    projects[i%len(projects)],
			Categories: categories[i%len(categories)],
			Tags:       tags[i%len(tags)],
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		events[i] = event
	}

	return events
}

// CreateSampleBursts generates test bursts linked to the provided events.
// Each burst groups 2-3 related events together.
func CreateSampleBursts(count int, events []*career.CareerEvent) []*career.Burst {
	bursts := make([]*career.Burst, count)

	burstNames := []string{
		"Authentication System Overhaul",
		"Team Mentoring Initiative",
		"Performance Optimization Sprint",
		"Process Improvement Quarter",
		"Microservices Migration",
	}

	burstDescriptions := []string{
		"A comprehensive effort to improve the security and reliability of our authentication system",
		"Focused mentoring and development of junior team members",
		"Performance optimization across multiple system components",
		"Improving team processes and development workflows",
		"Strategic migration from monolith to microservices architecture",
	}

	now := time.Now()
	for i := 0; i < count; i++ {
		// Assign 2-3 events to each burst
		eventIDs := make([]string, 0)
		startIdx := (i * 2) % len(events)
		eventCount := 2 + (i % 2) // 2 or 3 events per burst

		for j := 0; j < eventCount && startIdx+j < len(events); j++ {
			eventIDs = append(eventIDs, events[startIdx+j].ID)
		}

		burst := &career.Burst{
			ID:          fmt.Sprintf("burst_%d", i+1),
			Name:        burstNames[i%len(burstNames)],
			Description: burstDescriptions[i%len(burstDescriptions)],
			EventIDs:    eventIDs,
			Confirmed:   i%2 == 0, // Half confirmed, half not
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		bursts[i] = burst
	}

	return bursts
}

// CreateSampleFacts generates test facts linked to the provided events.
// Facts have varied competency categories, role fits, and audience relevance.
func CreateSampleFacts(count int, events []*career.CareerEvent) []*career.Fact {
	facts := make([]*career.Fact, count)

	factTexts := []string{
		"Reduced authentication failures by 40%% through system redesign",
		"Mentored 3 junior developers on testing and code review practices",
		"Improved API response times by 60%% via caching implementation",
		"Facilitated team retrospectives driving process improvements",
		"Designed microservices migration strategy for legacy systems",
		"Created comprehensive deployment documentation",
		"Led cross-team collaboration for on-time Q4 delivery",
		"Optimized database queries reducing page load by 50%%",
		"Introduced TDD practices across engineering team",
		"Built notification system handling 10K messages/second",
	}

	competencies := [][]string{
		{"technical", "product"},
		{"leadership", "mentoring"},
		{"technical", "research"},
		{"leadership", "consulting"},
		{"technical", "leadership"},
		{"product", "research"},
		{"leadership", "product"},
		{"technical", "mentoring"},
		{"leadership", "research"},
		{"technical", "consulting"},
	}

	roleFits := []career.RoleFit{career.RoleFitStaff, career.RoleFitSeniorIC, career.RoleFitPrincipal, career.RoleFitEM, career.RoleFitStaff}
	audiences := [][]string{
		{"hiring_manager", "peer"},
		{"hiring_manager", "recruiter"},
		{"peer"},
		{"hiring_manager"},
		{"hiring_manager", "peer", "recruiter"},
	}
	signals := []string{"high", "high", "medium", "medium", "high"}

	now := time.Now()
	for i := 0; i < count; i++ {
		// Link to an event
		eventIdx := i % len(events)
		sourceEventID := ""
		if len(events) > 0 {
			sourceEventID = events[eventIdx].ID
		}

		fact := &career.Fact{
			ID:                   fmt.Sprintf("fact_%d", i+1),
			Text:                 factTexts[i%len(factTexts)],
			CompetencyCategories: competencies[i%len(competencies)],
			RoleFit:              roleFits[i%len(roleFits)],
			AudienceRelevance:    audiences[i%len(audiences)],
			StrengthSignal:       signals[i%len(signals)],
			SourceEventID:        sourceEventID,
			CreatedAt:            now,
			UpdatedAt:            now,
		}
		facts[i] = fact
	}

	return facts
}

// CreateSampleProfiles returns a set of CV profiles for testing.
func CreateSampleProfiles() []CVProfile {
	return []CVProfile{
		{
			ID:             "profile_senior_ic",
			Name:           "Senior IC - Tech Lead",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			Description:    "Profile for senior individual contributor positions",
		},
		{
			ID:             "profile_em",
			Name:           "Engineering Manager",
			TargetRole:     "em",
			TargetAudience: "hiring_manager",
			Description:    "Profile for engineering manager positions",
		},
		{
			ID:             "profile_staff",
			Name:           "Staff Engineer",
			TargetRole:     "staff",
			TargetAudience: "hiring_manager",
			Description:    "Profile for staff engineer positions",
		},
		{
			ID:             "profile_principal",
			Name:           "Principal Engineer",
			TargetRole:     "principal",
			TargetAudience: "hiring_manager",
			Description:    "Profile for principal engineer positions",
		},
	}
}

// CVProfile represents a CV profile for testing (mirrors intents.CVProfile)
type CVProfile struct {
	ID             string
	Name           string
	TargetRole     string
	TargetAudience string
	Description    string
}

// PopulateTestData adds events, bursts, and facts to the test environment.
// Returns the environment for method chaining.
func (e *TestEnv) PopulateTestData(eventCount, burstCount, factCount int) *TestEnv {
	e.T.Helper()

	events := CreateSampleEvents(eventCount)
	for _, event := range events {
		e.AddEvent(event)
	}

	// Only create bursts if we have events to link them to
	if len(events) > 0 && burstCount > 0 {
		bursts := CreateSampleBursts(burstCount, events)
		for _, burst := range bursts {
			e.AddBurst(burst)
		}
	}

	// Only create facts if we have events to link them to
	if len(events) > 0 && factCount > 0 {
		facts := CreateSampleFacts(factCount, events)
		for _, fact := range facts {
			e.AddFact(fact)
		}
	}

	return e
}

// CreateMinimalEvent creates a single minimal valid event for testing.
func CreateMinimalEvent(id string) *career.CareerEvent {
	now := time.Now()
	return &career.CareerEvent{
		ID:        id,
		Text:      "Test event " + id,
		Date:      now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateMinimalBurst creates a single minimal valid burst for testing.
func CreateMinimalBurst(id string, eventIDs []string) *career.Burst {
	now := time.Now()
	return &career.Burst{
		ID:        id,
		Name:      "Test burst " + id,
		EventIDs:  eventIDs,
		Confirmed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateMinimalFact creates a single minimal valid fact for testing.
func CreateMinimalFact(id string, sourceEventID string) *career.Fact {
	now := time.Now()
	return &career.Fact{
		ID:                   id,
		Text:                 "Test fact " + id,
		CompetencyCategories: []string{"technical"},
		RoleFit:              "staff",
		AudienceRelevance:    []string{"hiring_manager"},
		StrengthSignal:       "high",
		SourceEventID:        sourceEventID,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}
