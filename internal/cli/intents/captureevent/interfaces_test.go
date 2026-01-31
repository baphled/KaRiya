package captureevent_test

import (
	ce "github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/cli/service"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// mockEventService is a minimal test implementation of EventService.
type mockEventService struct {
	cliService    *service.CLIEventService
	careerService *careerservice.Service
}

func (m *mockEventService) GetCLIService() *service.CLIEventService {
	return m.cliService
}

func (m *mockEventService) GetCareerService() *careerservice.Service {
	return m.careerService
}

var _ = Describe("EventService Interface", func() {
	It("should be satisfied by a mock implementation (compile-time check)", func() {
		var svc ce.EventService = &mockEventService{}
		Expect(svc).NotTo(BeNil())
	})

	It("should return nil services when not configured", func() {
		svc := &mockEventService{}
		Expect(svc.GetCLIService()).To(BeNil())
		Expect(svc.GetCareerService()).To(BeNil())
	})
})
