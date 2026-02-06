//go:generate mockgen -destination=../../../testutil/mocks/intent/browse_event_service_mock.go -package=mockintent -mock_names=EventCRUDService=MockBrowseEventCRUDService,EventSkillService=MockBrowseEventSkillService github.com/baphled/kariya/internal/cli/intents/browsetimeline EventCRUDService,EventSkillService

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

// EventCRUDService defines the interface for event create, read, update, delete operations.
// This follows the Interface Segregation Principle by separating event lifecycle
// operations from skill-event relationship operations.
type EventCRUDService interface {
	DeleteEvent(ctx context.Context, eventID string) error
	ListEvents(ctx context.Context, filters *careerrepo.EventListFilters) ([]*career.Event, error)
	CaptureEvent(ctx context.Context, text string, date time.Time, mode careerservice.EventCaptureMode, opts ...service.Option) error
	UpdateEventMetadata(ctx context.Context, event *career.Event) error
}

// EventSkillService defines the interface for skill-event relationship operations.
// This follows the Interface Segregation Principle by separating skill linking
// operations from core event CRUD operations.
type EventSkillService interface {
	GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)
	LinkSkillToEvent(ctx context.Context, eventID string, skillID string) error
	UnlinkSkillFromEvent(ctx context.Context, eventID string, skillID string) error
	ListAllSkills(ctx context.Context) ([]*career.Skill, error)
}

// EventService combines EventCRUDService and EventSkillService for backwards
// compatibility. New code should prefer the segregated interfaces.
type EventService interface {
	EventCRUDService
	EventSkillService
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
