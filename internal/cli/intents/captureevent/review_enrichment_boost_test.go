//nolint:errcheck // Test file - error handling for test setup is not relevant.
package captureevent_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/domain/capture"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReviewEnrichmentModel — ExtractInput", func() {
	var (
		model   *captureevent.ReviewEnrichmentModel
		service *careerservice.Service
		ctx     context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo := careermemory.NewEventRepository()
		skillRepo := careermemory.NewSkillRepository()
		service = careerservice.NewService(repo)
		service.SetSkillRepository(skillRepo)
	})

	Describe("ExtractInput", func() {
		Context("when the model is initialised from an event with metadata", func() {
			It("should return a MetadataInput matching the event fields", func() {
				event := fixtures.EventWith("extract-input-1", "Extract input test", "Acme Corp", "Project X")
				event.Date = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
				event.Tags = []string{"golang", "tui"}
				event.Categories = []string{"development", "backend"}
				event.Skills = []string{"go", "bubbletea"}

				model = captureevent.NewReviewEnrichmentModel(ctx, event, service, nil, nil)
				model.Init()

				input := model.ExtractInput()

				Expect(input).To(BeAssignableToTypeOf(capture.MetadataInput{}))
				Expect(input.Company).To(Equal("Acme Corp"))
				Expect(input.Project).To(Equal("Project X"))
				Expect(input.Date).To(Equal("2024-06-15"))
				Expect(input.Tags).To(ConsistOf("golang", "tui"))
				Expect(input.Categories).To(ConsistOf("development", "backend"))
				Expect(input.Skills).To(ConsistOf("go", "bubbletea"))
			})
		})

		Context("when the model is initialised from an event with empty metadata", func() {
			It("should return a MetadataInput with empty fields", func() {
				event := fixtures.EventWith("extract-input-2", "Extract input empty", "", "")
				event.Date = time.Time{}
				event.Tags = nil
				event.Categories = nil
				event.Skills = nil

				model = captureevent.NewReviewEnrichmentModel(ctx, event, service, nil, nil)
				model.Init()

				input := model.ExtractInput()

				Expect(input).To(BeAssignableToTypeOf(capture.MetadataInput{}))
				Expect(input.Company).To(BeEmpty())
				Expect(input.Project).To(BeEmpty())
				Expect(input.Tags).To(BeEmpty())
				Expect(input.Categories).To(BeEmpty())
				Expect(input.Skills).To(BeEmpty())
			})
		})
	})
})
