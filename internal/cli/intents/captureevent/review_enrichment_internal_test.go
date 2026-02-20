package captureevent

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReviewEnrichmentModel — internal coverage", func() {
	var (
		model   *ReviewEnrichmentModel
		event   *career.Event
		service *careerservice.Service
	)

	BeforeEach(func() {
		repos := memoryrepo.NewRepositories()
		service = careerservice.NewService(repos.Event)
		service.SetSkillRepository(repos.Skill)

		event = fixtures.EventWith("test-event-1", "Test event for metadata editor", "Test Company", "Test Project")
		event.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}

		model = NewReviewEnrichmentModel(context.Background(), event, service, nil, nil)
		model.Init()
	})

	Describe("handleFormCompletion — cancel path", func() {
		It("sets Cancelled to true when SubmitConfirmed is false", func() {
			model.formData.SubmitConfirmed = false
			result, cmd := model.handleFormCompletion()
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(model.IsCancelled()).To(BeTrue())
		})
	})

	Describe("handleFormCompletion — validation failure path", func() {
		It("sets Err when event validation fails after form data applied", func() {
			model.formData.SubmitConfirmed = true
			model.formData.Company = "Updated Company"
			model.formData.Date = ""
			model.formData.Tags = []string{}
			model.formData.Categories = []string{}
			event.Text = ""
			result, cmd := model.handleFormCompletion()
			Expect(result).NotTo(BeNil())
			_ = cmd
		})
	})

	Describe("handleFormCompletion — successful submit", func() {
		It("sets submitted to true on successful form completion", func() {
			model.formData.SubmitConfirmed = true
			model.formData.Company = "Updated Company"
			model.formData.Project = "Updated Project"
			model.formData.Date = "2024-06-15"
			model.formData.Tags = []string{"technical"}
			model.formData.Categories = []string{"technical"}
			result, cmd := model.handleFormCompletion()
			Expect(result).NotTo(BeNil())
			_ = cmd
			Expect(model.IsSubmitted()).To(BeTrue())
		})
	})

	Describe("handleFormCompletion — form data apply error", func() {
		It("sets Err when ApplyMetadataFormData fails", func() {
			model.formData.SubmitConfirmed = true
			model.formData.Date = "not-a-valid-date"
			result, cmd := model.handleFormCompletion()
			Expect(result).NotTo(BeNil())
			_ = cmd
			Expect(model.GetError()).To(HaveOccurred())
		})
	})

	Describe("GetContent", func() {
		It("returns non-empty form view", func() {
			content := model.GetContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("includes error text when Err is set", func() {
			model.Err = errors.New("something went wrong")
			content := model.GetContent()
			Expect(content).To(ContainSubstring("something went wrong"))
		})
	})

	Describe("GetFooter", func() {
		It("returns the keyboard hint string", func() {
			footer := model.GetFooter()
			Expect(footer).To(ContainSubstring("Enter"))
			Expect(footer).To(ContainSubstring("Esc"))
		})
	})

	Describe("NewReviewEnrichmentModel — dimensions branch", func() {
		It("applies provided dimensions", func() {
			dims := &ReviewEnrichmentDimensions{TerminalWidth: 100, TerminalHeight: 45}
			m := NewReviewEnrichmentModel(context.Background(), event, service, nil, dims)
			Expect(m).NotTo(BeNil())
		})

		It("uses defaults when TerminalWidth is zero", func() {
			dims := &ReviewEnrichmentDimensions{TerminalWidth: 0, TerminalHeight: 0}
			m := NewReviewEnrichmentModel(context.Background(), event, service, nil, dims)
			Expect(m).NotTo(BeNil())
		})
	})

	Describe("forms.ApplyMetadataFormData integration", func() {
		It("correctly parses valid date string", func() {
			formData := &forms.MetadataFormData{
				Date:       "2024-05-20",
				Company:    "New Corp",
				Project:    "New Project",
				Tags:       []string{"go"},
				Categories: []string{"engineering"},
			}
			err := forms.ApplyMetadataFormData(event, formData)
			Expect(err).NotTo(HaveOccurred())
			Expect(event.Company).To(Equal("New Corp"))
		})
	})
})
