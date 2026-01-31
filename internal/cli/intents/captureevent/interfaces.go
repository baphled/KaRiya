package captureevent

import (
	"github.com/baphled/kariya/internal/cli/service"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// EventService defines the contract for accessing event-related services.
//
// Expected:
//   - Implementations must return fully initialised service instances
//   - Both methods must be safe to call concurrently
//
// This interface decouples the intent from concrete service construction.
type EventService interface {
	// GetCLIService returns the CLI event service used for form submission
	// and event creation.
	GetCLIService() *service.CLIEventService

	// GetCareerService returns the domain career service used for enrichment
	// operations such as burst detection and fact extraction.
	GetCareerService() *careerservice.Service
}
