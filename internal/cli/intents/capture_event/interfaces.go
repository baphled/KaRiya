package capture_event

import (
	"github.com/baphled/kariya/internal/cli/service"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// EventService defines the service interface for event operations used by
// the capture event intent.
type EventService interface {
	GetCLIService() *service.CLIEventService
	GetCareerService() *careerservice.Service
}
