package models

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BulkOperationsModel", func() {
	var (
		repo   *careerrepo.MemoryRepository
		svc    *careerservice.Service
		cliSvc *service.CLIEventService
		ctx    context.Context
		events []*career.CareerEvent
		model  *BulkOperationsModel
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(svc)
		ctx = context.Background()

		// Create test events
		events = []*career.CareerEvent{
			{
				ID:        "event1",
				Text:      "Led team on project A",
				Date:      time.Now().Add(-24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event2",
				Text:      "Implemented feature B",
				Date:      time.Now().Add(-48 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event3",
				Text:      "Mentored junior developer",
				Date:      time.Now().Add(-72 * time.Hour),
				Company:   "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
	})

	Context("when creating bulk operations model", func() {
		It("should initialize with provided events", func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			Expect(model).NotTo(BeNil())
		})

		It("should have no events selected initially", func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			Expect(model.GetSelectedCount()).To(Equal(0))
		})

		It("should track all provided events", func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			Expect(model.GetEventCount()).To(Equal(3))
		})
	})

	Context("when selecting events", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
		})

		It("should toggle selection on space key", func() {
			model.ToggleSelection(0)
			Expect(model.GetSelectedCount()).To(Equal(1))
		})

		It("should select all events with 'a' key", func() {
			model.SelectAll()
			Expect(model.GetSelectedCount()).To(Equal(3))
		})

		It("should deselect all events with 'd' key", func() {
			model.SelectAll()
			model.DeselectAll()
			Expect(model.GetSelectedCount()).To(Equal(0))
		})

		It("should maintain selection state across toggles", func() {
			model.ToggleSelection(0)
			model.ToggleSelection(1)
			Expect(model.GetSelectedCount()).To(Equal(2))
			model.ToggleSelection(0)
			Expect(model.GetSelectedCount()).To(Equal(1))
		})
	})

	Context("when navigating events", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
		})

		It("should move focus down with arrow key", func() {
			initialFocus := model.GetFocusIndex()
			model.MoveDown()
			Expect(model.GetFocusIndex()).To(Equal(initialFocus + 1))
		})

		It("should move focus up with arrow key", func() {
			model.MoveDown()
			model.MoveDown()
			model.MoveUp()
			Expect(model.GetFocusIndex()).To(Equal(1))
		})

		It("should not move focus below last event", func() {
			model.MoveDown()
			model.MoveDown()
			model.MoveDown()
			Expect(model.GetFocusIndex()).To(Equal(2))
		})

		It("should not move focus above first event", func() {
			model.MoveUp()
			Expect(model.GetFocusIndex()).To(Equal(0))
		})
	})

	Context("when editing bulk metadata", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SelectAll()
		})

		It("should set bulk company field", func() {
			model.SetBulkCompany("NewCorp")
			Expect(model.GetBulkCompany()).To(Equal("NewCorp"))
		})

		It("should set bulk project field", func() {
			model.SetBulkProject("ProjectX")
			Expect(model.GetBulkProject()).To(Equal("ProjectX"))
		})

		It("should set bulk tags", func() {
			model.SetBulkTags([]string{"leadership", "technical"})
			Expect(model.GetBulkTags()).To(ContainElements("leadership", "technical"))
		})

		It("should set bulk categories", func() {
			model.SetBulkCategories([]string{"Leadership", "Technical"})
			Expect(model.GetBulkCategories()).To(ContainElements("Leadership", "Technical"))
		})

		It("should set apply-if-empty flag for company", func() {
			model.SetApplyIfEmptyCompany(true)
			Expect(model.GetApplyIfEmptyCompany()).To(BeTrue())
		})

		It("should set apply-if-empty flag for project", func() {
			model.SetApplyIfEmptyProject(true)
			Expect(model.GetApplyIfEmptyProject()).To(BeTrue())
		})
	})

	Context("when previewing bulk changes", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.ToggleSelection(0)
			model.ToggleSelection(1)
			model.SetBulkCompany("UpdatedCorp")
		})

		It("should generate preview of changes", func() {
			preview := model.GetPreview()
			Expect(preview).NotTo(BeEmpty())
		})

		It("should show which events will be affected", func() {
			preview := model.GetPreview()
			Expect(preview).To(ContainSubstring("2"))
		})

		It("should show which fields will be updated", func() {
			preview := model.GetPreview()
			Expect(preview).To(ContainSubstring("company"))
		})
	})

	Context("when submitting bulk changes", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SelectAll()
			model.SetBulkCompany("BulkCorp")
		})

		It("should return true when submitted", func() {
			model.Submit()
			Expect(model.IsSubmitted()).To(BeTrue())
		})

		It("should not be cancelled after submit", func() {
			model.Submit()
			Expect(model.IsCancelled()).To(BeFalse())
		})

		It("should provide summary of changes", func() {
			model.Submit()
			summary := model.GetSummary()
			Expect(summary).NotTo(BeNil())
			Expect(summary.EventsAffected).To(Equal(3))
		})
	})

	Context("when cancelling bulk operation", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
			model.SelectAll()
			model.SetBulkCompany("BulkCorp")
		})

		It("should return true when cancelled", func() {
			model.Cancel()
			Expect(model.IsCancelled()).To(BeTrue())
		})

		It("should not be submitted after cancel", func() {
			model.Cancel()
			Expect(model.IsSubmitted()).To(BeFalse())
		})

		It("should discard bulk changes on cancel", func() {
			model.Cancel()
			Expect(model.GetBulkCompany()).To(BeEmpty())
		})
	})

	Context("when rendering bulk operations view", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
		})

		It("should render without error", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should display event list", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Led team"))
		})

		It("should display selected count", func() {
			model.SelectAll()
			view := model.View()
			Expect(view).To(ContainSubstring("3"))
		})

		It("should display keyboard shortcuts", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("elec"))
		})
	})

	Context("when handling window resize", func() {
		BeforeEach(func() {
			model = NewBulkOperationsModel(events, svc, cliSvc, ctx)
		})

		It("should update width on resize", func() {
			model.SetSize(100, 30)
			Expect(model.GetWidth()).To(Equal(100))
		})

		It("should update height on resize", func() {
			model.SetSize(100, 30)
			Expect(model.GetHeight()).To(Equal(30))
		})

		It("should render with new dimensions", func() {
			model.SetSize(100, 30)
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		Describe("Bulk operations with field origin awareness", func() {
			It("should store field origins for bulk operations", func() {
				// Arrange
				events := []*career.CareerEvent{
					{ID: "1", Text: "Event 1", Company: "CompanyA"},
					{ID: "2", Text: "Event 2", Company: "CompanyB"},
				}
				model := NewBulkOperationsModel(events, svc, cliSvc, ctx)
				eventID := "1"
				origins := map[string]bool{
					"company": false, // default
					"project": false, // default
				}

				// Act
				model.SetFieldOrigins(eventID, origins)
				retrieved := model.GetFieldOrigins(eventID)

				// Assert
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved["company"]).To(BeFalse())
			})

			It("should only allow updates to default fields in bulk operations", func() {
				// Arrange
				events := []*career.CareerEvent{
					{ID: "1", Text: "Event 1", Company: ""},
					{ID: "2", Text: "Event 2", Company: ""},
				}
				model := NewBulkOperationsModel(events, svc, cliSvc, ctx)

				// Mark company as default (not from CSV)
				model.SetFieldOrigins("1", map[string]bool{"company": false})
				model.SetFieldOrigins("2", map[string]bool{"company": false})

				// Act
				canUpdate := model.CanUpdateField("1", "company")

				// Assert
				Expect(canUpdate).To(BeTrue())
			})

			It("should prevent updates to CSV fields in bulk operations", func() {
				// Arrange
				events := []*career.CareerEvent{
					{ID: "1", Text: "Event 1", Company: "CompanyA"},
				}
				model := NewBulkOperationsModel(events, svc, cliSvc, ctx)

				// Mark company as from CSV
				model.SetFieldOrigins("1", map[string]bool{"company": true})

				// Act
				canUpdate := model.CanUpdateField("1", "company")

				// Assert
				Expect(canUpdate).To(BeFalse())
			})

			It("should get default fields for an event in bulk operations", func() {
				// Arrange
				events := []*career.CareerEvent{
					{ID: "1", Text: "Event 1"},
				}
				model := NewBulkOperationsModel(events, svc, cliSvc, ctx)
				origins := map[string]bool{
					"company": false,
					"project": false,
					"tags":    true,
				}
				model.SetFieldOrigins("1", origins)

				// Act
				defaults := model.GetDefaultFields("1")

				// Assert
				Expect(defaults).To(HaveLen(2))
				Expect(defaults).To(ContainElement("company"))
				Expect(defaults).To(ContainElement("project"))
			})
		})

	})
})
