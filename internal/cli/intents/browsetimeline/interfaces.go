//go:generate mockgen -destination=../../../testutil/mocks/intent/browse_event_service_mock.go -package=mockintent -mock_names=EventService=MockBrowseEventService github.com/baphled/kariya/internal/cli/intents/browsetimeline EventService

// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// EventService defines the interface for event CRUD operations.
// This allows for mocking in tests.
type EventService interface {
	DeleteEvent(ctx context.Context, eventID string) error
	ListEvents(ctx context.Context, filters *careerrepo.EventListFilters) ([]*career.Event, error)
	CaptureEvent(ctx context.Context, text string, date time.Time, mode careerservice.EventCaptureMode, opts ...service.Option) error
	UpdateEventMetadata(ctx context.Context, event *career.Event) error
	GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)
	LinkSkillToEvent(ctx context.Context, eventID string, skillID string) error
	UnlinkSkillFromEvent(ctx context.Context, eventID string, skillID string) error
	ListAllSkills(ctx context.Context) ([]*career.Skill, error)
}

// SkillService defines the interface for skill operations.
type SkillService interface {
	Create(ctx context.Context, skill *career.Skill) error
}

// SkillInferenceService defines the interface for skill inference operations.
type SkillInferenceService interface {
	InferSkillsFromEvents(ctx context.Context, events []*career.Event) (*skillinference.InferenceResult, error)
	CreateSkillsFromSuggestions(ctx context.Context, suggestions []skillinference.SkillSuggestion) ([]*career.Skill, error)
}
