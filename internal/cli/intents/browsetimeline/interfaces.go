// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// EventService defines the interface for event CRUD operations.
// This allows for mocking in tests.
type EventService interface {
	DeleteEvent(ctx context.Context, eventID string) error
	ListEvents(ctx context.Context, filters *careerrepo.EventListFilters) ([]*career.CareerEvent, error)
	CaptureEvent(ctx context.Context, text string, date time.Time, mode careerservice.EventCaptureMode, opts ...service.Option) error
	UpdateEventMetadata(ctx context.Context, event *career.CareerEvent) error
	GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)
}
