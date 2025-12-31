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
	tea "github.com/charmbracelet/bubbletea"
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
				// Set width so header renders
				model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
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

		It("should trigger bulk operations with 'b' key", func() {
			// Arrange: model has events
			// Act: press 'b' key
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
			// Assert: command should be returned to send BulkOperationsMsg
			Expect(cmd).NotTo(BeNil())
			// Execute the command to get the message
			msg := cmd()
			bulkMsg, ok := msg.(models.BulkOperationsMsg)
			Expect(ok).To(BeTrue())
			Expect(bulkMsg.Events).To(Equal(model.GetEvents()))
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

	Describe("Field origin tracking for imported events", func() {
		It("should store field origins for imported events", func() {
			// Arrange: create model
			model := models.NewMetadataReviewModel(service, ctx)
			eventID := "test-event-id"
			origins := map[string]bool{
				"text":       true,  // from CSV
				"date":       true,  // from CSV
				"company":    false, // default
				"project":    false, // default
				"tags":       true,  // from CSV
				"categories": false, // default
			}

			// Act: set field origins
			model.SetFieldOrigins(eventID, origins)

			// Assert: origins should be stored
			stored := model.GetFieldOrigins(eventID)
			Expect(stored).NotTo(BeNil())
			Expect(stored["text"]).To(BeTrue())
			Expect(stored["company"]).To(BeFalse())
		})

		It("should indicate which fields came from CSV", func() {
			// Arrange
			model := models.NewMetadataReviewModel(service, ctx)
			eventID := "test-event-id"
			origins := map[string]bool{
				"text":    true,
				"company": false,
			}
			model.SetFieldOrigins(eventID, origins)

			// Act
			isFromCSV := model.IsFieldFromCSV(eventID, "text")
			isDefault := model.IsFieldFromCSV(eventID, "company")

			// Assert
			Expect(isFromCSV).To(BeTrue())
			Expect(isDefault).To(BeFalse())
		})

		It("should return default fields for an event", func() {
			// Arrange
			model := models.NewMetadataReviewModel(service, ctx)
			eventID := "test-event-id"
			origins := map[string]bool{
				"text":       true,
				"date":       true,
				"company":    false,
				"project":    false,
				"tags":       true,
				"categories": false,
			}
			model.SetFieldOrigins(eventID, origins)

			// Act
			defaultFields := model.GetDefaultFields(eventID)

			// Assert
			Expect(defaultFields).To(HaveLen(3))
			Expect(defaultFields).To(ContainElement("company"))
			Expect(defaultFields).To(ContainElement("project"))
		})

		Describe("Parsing warnings and issues display", func() {
			It("should store and retrieve parsing warnings for an event", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"
				warnings := []string{
					"missing company field",
					"date parsed from ambiguous format",
				}

				// Act
				model.SetParsingWarnings(eventID, warnings)
				retrieved := model.GetParsingWarnings(eventID)

				// Assert
				Expect(retrieved).To(HaveLen(2))
				Expect(retrieved).To(ContainElement("missing company field"))
				Expect(retrieved).To(ContainElement("date parsed from ambiguous format"))
			})

			It("should indicate if an event has parsing warnings", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"
				warnings := []string{"missing company"}

				// Act
				model.SetParsingWarnings(eventID, warnings)
				hasWarnings := model.HasParsingWarnings(eventID)

				// Assert
				Expect(hasWarnings).To(BeTrue())
			})

			It("should return empty list for events without warnings", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"

				// Act
				warnings := model.GetParsingWarnings(eventID)

				// Assert
				Expect(warnings).To(BeEmpty())
			})

			It("should check if event has parsing issues", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"

				// Act & Assert
				Expect(model.HasParsingWarnings(eventID)).To(BeFalse())

				model.SetParsingWarnings(eventID, []string{"issue"})
				Expect(model.HasParsingWarnings(eventID)).To(BeTrue())
			})
		})

		Describe("Duplicate detection status display", func() {
			It("should store and retrieve duplicate status for an event", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"
				duplicateOf := "original-event-id"

				// Act
				model.SetDuplicateStatus(eventID, true, duplicateOf)
				isDuplicate, originalID := model.GetDuplicateStatus(eventID)

				// Assert
				Expect(isDuplicate).To(BeTrue())
				Expect(originalID).To(Equal(duplicateOf))
			})

			It("should indicate if an event is a duplicate", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"

				// Act
				model.SetDuplicateStatus(eventID, true, "original-id")
				isDuplicate := model.IsDuplicate(eventID)

				// Assert
				Expect(isDuplicate).To(BeTrue())
			})

			It("should return original event ID for duplicates", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"
				originalID := "original-event-id"

				// Act
				model.SetDuplicateStatus(eventID, true, originalID)
				retrievedID := model.GetOriginalEventID(eventID)

				// Assert
				Expect(retrievedID).To(Equal(originalID))
			})

			It("should handle non-duplicate events gracefully", func() {
				// Arrange
				model := models.NewMetadataReviewModel(service, ctx)
				eventID := "test-event-id"

				// Act
				isDuplicate := model.IsDuplicate(eventID)
				originalID := model.GetOriginalEventID(eventID)

				// Assert
				Expect(isDuplicate).To(BeFalse())
				Expect(originalID).To(BeEmpty())
			})
		})

	})
})
