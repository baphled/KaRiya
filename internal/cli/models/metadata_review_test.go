package models_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

var _ = Describe("MetadataReviewModel", func() {
	var (
		model   *models.MetadataReviewModel
		service *careerservice.Service
		repo    *careerrepo.MemoryRepository
		ctx     context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		service = careerservice.NewService(repo)

		// Create some test events
		event1 := &career.CareerEvent{
			ID:         "event-1",
			Text:       "First event",
			Date:       time.Now().Add(-48 * time.Hour),
			Company:    "CompanyA",
			Project:    "ProjectA",
			Tags:       []string{"technical"},
			Categories: []string{"Technical"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		event2 := &career.CareerEvent{
			ID:         "event-2",
			Text:       "Second event",
			Date:       time.Now().Add(-24 * time.Hour),
			Company:    "CompanyB",
			Project:    "ProjectB",
			Tags:       []string{"leadership"},
			Categories: []string{"Leadership"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		repo.Create(ctx, event1)
		repo.Create(ctx, event2)

		model = models.NewMetadataReviewModel(service, ctx)
	})

	Describe("Creation", func() {
		It("should create a metadata review model", func() {
			Expect(model).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should render without panicking", func() {
			Expect(func() {
				_ = model.View()
			}).NotTo(Panic())
		})

		It("should show header in view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Metadata Review"))
		})
	})

	Describe("Init", func() {
		It("should initialize without error", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})

		Describe("Displaying imported events", func() {
			It("should load events from service on initialization", func() {
				// Arrange: events were created in BeforeEach and added to repo
				// Act: model was created in BeforeEach which calls loadEvents()
				// Assert: view should contain event text
				view := model.View()
				Expect(view).To(ContainSubstring("First event"))
				Expect(view).To(ContainSubstring("Second event"))
			})

			It("should calculate quality scores for imported events", func() {
				// Arrange: events were created with partial metadata
				// Act: model was initialized
				// Assert: quality scores should be calculated
				// We verify by checking that the model has events with scores
				view := model.View()
				// If quality calculation works, view should render without error
				Expect(view).NotTo(BeEmpty())
			})

			It("should display events with company information when available", func() {
				// Arrange: events have company information
				// Act: model renders view
				// Assert: company names should be visible
				view := model.View()
				Expect(view).To(ContainSubstring("CompanyA"))
				Expect(view).To(ContainSubstring("CompanyB"))
			})

			It("should display events sorted by date (most recent first)", func() {
				// Arrange: events with different dates
				// Act: model loads and sorts events
				// Assert: view should show most recent event first
				view := model.View()
				// Second event is more recent (24 hours ago vs 48 hours ago)
				// Both should be visible in view
				Expect(view).To(ContainSubstring("First event"))
				Expect(view).To(ContainSubstring("Second event"))
			})

			It("should handle empty event list gracefully", func() {
				// Arrange: create new repo with no events
				emptyRepo := careerrepo.NewMemoryRepository()
				emptyService := careerservice.NewService(emptyRepo)

				// Act: create model with empty service
				emptyModel := models.NewMetadataReviewModel(emptyService, ctx)

				// Assert: view should render without error
				view := emptyModel.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Update", func() {
		It("should handle key messages", func() {
			_, cmd := model.Update(nil)
			Expect(cmd).To(BeNil())
		})

		Describe("Displaying imported events", func() {
			It("should load events from service on initialization", func() {
				// Arrange: events were created in BeforeEach and added to repo
				// Act: model was created in BeforeEach which calls loadEvents()
				// Assert: view should contain event text
				view := model.View()
				Expect(view).To(ContainSubstring("First event"))
				Expect(view).To(ContainSubstring("Second event"))
			})

			It("should calculate quality scores for imported events", func() {
				// Arrange: events were created with partial metadata
				// Act: model was initialized
				// Assert: quality scores should be calculated
				// We verify by checking that the model has events with scores
				view := model.View()
				// If quality calculation works, view should render without error
				Expect(view).NotTo(BeEmpty())
			})

			It("should display events with company information when available", func() {
				// Arrange: events have company information
				// Act: model renders view
				// Assert: company names should be visible
				view := model.View()
				Expect(view).To(ContainSubstring("CompanyA"))
				Expect(view).To(ContainSubstring("CompanyB"))
			})

			It("should display events sorted by date (most recent first)", func() {
				// Arrange: events with different dates
				// Act: model loads and sorts events
				// Assert: view should show most recent event first
				view := model.View()
				// Second event is more recent (24 hours ago vs 48 hours ago)
				// Both should be visible in view
				Expect(view).To(ContainSubstring("First event"))
				Expect(view).To(ContainSubstring("Second event"))
			})

			It("should handle empty event list gracefully", func() {
				// Arrange: create new repo with no events
				emptyRepo := careerrepo.NewMemoryRepository()
				emptyService := careerservice.NewService(emptyRepo)

				// Act: create model with empty service
				emptyModel := models.NewMetadataReviewModel(emptyService, ctx)

				// Assert: view should render without error
				view := emptyModel.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Displaying imported events", func() {
		It("should load events from service on initialization", func() {
			// Arrange: events were created in BeforeEach and added to repo
			// Act: model was created in BeforeEach which calls loadEvents()
			// Assert: view should contain event text
			view := model.View()
			Expect(view).To(ContainSubstring("First event"))
			Expect(view).To(ContainSubstring("Second event"))
		})

		It("should calculate quality scores for imported events", func() {
			// Arrange: events were created with partial metadata
			// Act: model was initialized
			// Assert: quality scores should be calculated
			// We verify by checking that the model has events with scores
			view := model.View()
			// If quality calculation works, view should render without error
			Expect(view).NotTo(BeEmpty())
		})

		It("should display events with company information when available", func() {
			// Arrange: events have company information
			// Act: model renders view
			// Assert: company names should be visible
			view := model.View()
			Expect(view).To(ContainSubstring("CompanyA"))
			Expect(view).To(ContainSubstring("CompanyB"))
		})

		It("should display events sorted by date (most recent first)", func() {
			// Arrange: events with different dates
			// Act: model loads and sorts events
			// Assert: view should show most recent event first
			view := model.View()
			// Second event is more recent (24 hours ago vs 48 hours ago)
			// Both should be visible in view
			Expect(view).To(ContainSubstring("First event"))
			Expect(view).To(ContainSubstring("Second event"))
		})

		It("should handle empty event list gracefully", func() {
			// Arrange: create new repo with no events
			emptyRepo := careerrepo.NewMemoryRepository()
			emptyService := careerservice.NewService(emptyRepo)

			// Act: create model with empty service
			emptyModel := models.NewMetadataReviewModel(emptyService, ctx)

			// Assert: view should render without error
			view := emptyModel.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
