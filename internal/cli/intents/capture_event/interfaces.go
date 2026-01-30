package capture_event

import (
	"github.com/baphled/kariya/internal/cli/service"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// EventService defines the service interface for event operations used by
// the capture event intent. It provides access to both the CLI-level event
// service (for form operations) and the domain career service (for enrichment).
type EventService interface {
	// GetCLIService returns the CLI event service used for form submission.
	GetCLIService() *service.CLIEventService

	// GetCareerService returns the domain career service used for enrichment
	// operations such as burst detection and fact extraction.
	GetCareerService() *careerservice.Service
}
