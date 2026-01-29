package career

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career/mocks"
)

var _ = Describe("Career Service", func() {
	var (
		mockRepo *mocks.TestMockRepository
		service  *Service
		ctx      context.Context
	)

	BeforeEach(func() {
		mockRepo = mocks.NewTestMockRepository()
		service = NewService(mockRepo)
		ctx = context.Background()
	})

	Describe("CaptureEvent", func() {
		var testEvent *career.Event

		BeforeEach(func() {
			testEvent = &career.Event{
				Text:    "Developed a high-performance backend service",
				Date:    time.Now().AddDate(0, 0, -10),
				Tags:    []string{"technical", "project"},
				Company: "Test Company",
			}
		})

		Context("with valid event and TimelineJournaling mode", func() {
			It("should capture event successfully", func() {
				mockRepo.SetCreateBehavior(nil)

				err := service.CaptureEvent(ctx, testEvent, TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.ID).NotTo(BeEmpty())
				Expect(testEvent.CreatedAt).NotTo(BeZero())
				Expect(testEvent.UpdatedAt).NotTo(BeZero())
				Expect(testEvent.CreatedAt).To(Equal(testEvent.UpdatedAt))
				Expect(mockRepo.CreateCalled()).To(BeTrue())
			})

			It("should generate a valid UUID for the event", func() {
				mockRepo.SetCreateBehavior(nil)

				err := service.CaptureEvent(ctx, testEvent, TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
				Expect(testEvent.ID).NotTo(BeEmpty())
				Expect(len(testEvent.ID)).To(Equal(36))
			})

			It("should preserve event data", func() {
				mockRepo.SetCreateBehavior(nil)
				originalText := testEvent.Text
				originalDate := testEvent.Date
				originalTags := testEvent.Tags
				originalCompany := testEvent.Company

				err := service.CaptureEvent(ctx, testEvent, TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.Text).To(Equal(originalText))
				Expect(testEvent.Date).To(Equal(originalDate))
				Expect(testEvent.Tags).To(Equal(originalTags))
				Expect(testEvent.Company).To(Equal(originalCompany))
			})
		})

		Context("with valid event and CVBackfill mode", func() {
			It("should capture older events successfully", func() {
				mockRepo.SetCreateBehavior(nil)
				testEvent.Date = time.Now().AddDate(-2, 0, 0)

				err := service.CaptureEvent(ctx, testEvent, CVBackfill)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.ID).NotTo(BeEmpty())
				Expect(mockRepo.CreateCalled()).To(BeTrue())
			})

			It("should allow very old dates", func() {
				mockRepo.SetCreateBehavior(nil)
				testEvent.Date = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

				err := service.CaptureEvent(ctx, testEvent, CVBackfill)
				Expect(err).NotTo(HaveOccurred())
				Expect(testEvent.ID).NotTo(BeEmpty())
			})
		})

		Context("with valid event and ManualEntry mode", func() {
			It("should capture event successfully", func() {
				mockRepo.SetCreateBehavior(nil)

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.ID).NotTo(BeEmpty())
				Expect(mockRepo.CreateCalled()).To(BeTrue())
			})

			It("should allow any valid date", func() {
				mockRepo.SetCreateBehavior(nil)
				testEvent.Date = time.Now().AddDate(-5, 0, 0)

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
				Expect(testEvent.ID).NotTo(BeEmpty())
			})
		})

		Context("with TimelineJournaling mode and old event", func() {
			It("should reject events older than 30 days", func() {
				testEvent.Date = time.Now().AddDate(0, 0, -31)

				err := service.CaptureEvent(ctx, testEvent, TimelineJournaling)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("within 30 days"))
				Expect(mockRepo.CreateCalled()).To(BeFalse())
			})

			It("should accept events within 30 days", func() {
				mockRepo.SetCreateBehavior(nil)
				testEvent.Date = time.Now().AddDate(0, 0, -29)

				err := service.CaptureEvent(ctx, testEvent, TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
				Expect(testEvent.ID).NotTo(BeEmpty())
			})
		})

		Context("with invalid event data", func() {
			It("should reject empty text", func() {
				testEvent.Text = ""

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).To(HaveOccurred())
				Expect(mockRepo.CreateCalled()).To(BeFalse())
			})

			It("should reject future dates", func() {
				testEvent.Date = time.Now().AddDate(0, 0, 1)

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).To(HaveOccurred())
				Expect(mockRepo.CreateCalled()).To(BeFalse())
			})

			It("should reject invalid tags", func() {
				testEvent.Tags = []string{"invalid_tag"}

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).To(HaveOccurred())
				Expect(mockRepo.CreateCalled()).To(BeFalse())
			})

			It("should reject overly long text", func() {
				testEvent.Text = string(make([]byte, 2001))

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).To(HaveOccurred())
				Expect(mockRepo.CreateCalled()).To(BeFalse())
			})
		})

		Context("with invalid capture mode", func() {
			It("should reject unknown capture mode", func() {
				invalidMode := EventCaptureMode("invalid_mode")

				err := service.CaptureEvent(ctx, testEvent, invalidMode)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid event capture mode"))
				Expect(mockRepo.CreateCalled()).To(BeFalse())
			})
		})

		Context("with repository errors", func() {
			It("should return error when repository fails to create", func() {
				mockRepo.SetCreateBehavior(errors.New("database error"))

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
			})

			It("should still generate ID even if repository fails", func() {
				mockRepo.SetCreateBehavior(errors.New("database error"))

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).To(HaveOccurred())
				Expect(testEvent.ID).NotTo(BeEmpty())
			})
		})

		Context("with pre-existing ID", func() {
			It("should preserve provided ID", func() {
				mockRepo.SetCreateBehavior(nil)
				testEvent.ID = "custom-id-12345"

				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
				Expect(testEvent.ID).To(Equal("custom-id-12345"))
			})
		})
	})

	Describe("UpdateEvent", func() {
		var testEvent *career.Event
		var existingEvent *career.Event

		BeforeEach(func() {
			existingEvent = &career.Event{
				ID:        "existing-id",
				Text:      "Original event text",
				Date:      time.Now().AddDate(0, 0, -10),
				Tags:      []string{"technical"},
				Company:   "Original Company",
				CreatedAt: time.Now().AddDate(0, 0, -5),
				UpdatedAt: time.Now().AddDate(0, 0, -5),
			}

			testEvent = &career.Event{
				ID:      "existing-id",
				Text:    "Updated event text",
				Date:    time.Now().AddDate(0, 0, -10),
				Tags:    []string{"leadership"},
				Company: "New Company",
			}
		})

		Context("with valid event", func() {
			It("should update event successfully", func() {
				mockRepo.SetGetByIDBehavior(existingEvent, nil)
				mockRepo.SetUpdateBehavior(nil)

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).NotTo(HaveOccurred())

				Expect(mockRepo.UpdateCalled()).To(BeTrue())
			})

			It("should preserve creation timestamp", func() {
				mockRepo.SetGetByIDBehavior(existingEvent, nil)
				mockRepo.SetUpdateBehavior(nil)

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.CreatedAt).To(Equal(existingEvent.CreatedAt))
			})

			It("should update modification timestamp", func() {
				mockRepo.SetGetByIDBehavior(existingEvent, nil)
				mockRepo.SetUpdateBehavior(nil)
				originalUpdatedAt := existingEvent.UpdatedAt

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.UpdatedAt).To(BeTemporally(">", originalUpdatedAt))
			})

			It("should update all event fields", func() {
				mockRepo.SetGetByIDBehavior(existingEvent, nil)
				mockRepo.SetUpdateBehavior(nil)

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).NotTo(HaveOccurred())

				Expect(testEvent.Text).To(Equal("Updated event text"))
				Expect(testEvent.Company).To(Equal("New Company"))
				Expect(testEvent.Tags).To(Equal([]string{"leadership"}))
			})
		})

		Context("with invalid event data", func() {
			It("should reject empty text", func() {
				testEvent.Text = ""

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).To(HaveOccurred())
				Expect(mockRepo.UpdateCalled()).To(BeFalse())
			})

			It("should reject invalid tags", func() {
				testEvent.Tags = []string{"invalid_tag"}

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).To(HaveOccurred())
				Expect(mockRepo.UpdateCalled()).To(BeFalse())
			})
		})

		Context("when event does not exist", func() {
			It("should return error if event not found", func() {
				mockRepo.SetGetByIDBehavior(nil, errors.New("event not found"))

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("event not found"))
				Expect(mockRepo.UpdateCalled()).To(BeFalse())
			})
		})

		Context("with repository errors", func() {
			It("should return error when update fails", func() {
				mockRepo.SetGetByIDBehavior(existingEvent, nil)
				mockRepo.SetUpdateBehavior(errors.New("database error"))

				err := service.UpdateEvent(ctx, testEvent)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
			})
		})
	})

	Describe("DeleteEvent", func() {
		Context("with valid event ID", func() {
			It("should delete event successfully", func() {
				mockRepo.SetDeleteBehavior(nil)

				err := service.DeleteEvent(ctx, "event-id-123")
				Expect(err).NotTo(HaveOccurred())

				Expect(mockRepo.DeleteCalled()).To(BeTrue())
			})
		})

		Context("when event does not exist", func() {
			It("should return error if event not found", func() {
				mockRepo.SetDeleteBehavior(errors.New("event not found"))

				err := service.DeleteEvent(ctx, "non-existent-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("event not found"))
			})
		})

		Context("with repository errors", func() {
			It("should return error when deletion fails", func() {
				mockRepo.SetDeleteBehavior(errors.New("database error"))

				err := service.DeleteEvent(ctx, "event-id-123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
			})
		})

		Context("with empty event ID", func() {
			It("should attempt deletion with empty ID", func() {
				mockRepo.SetDeleteBehavior(errors.New("invalid ID"))

				err := service.DeleteEvent(ctx, "")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetEventByID", func() {
		var testEvent *career.Event

		BeforeEach(func() {
			testEvent = &career.Event{
				ID:        "event-id-123",
				Text:      "Test event",
				Date:      time.Now().AddDate(0, 0, -10),
				Tags:      []string{"technical"},
				Company:   "Test Company",
				CreatedAt: time.Now().AddDate(0, 0, -5),
				UpdatedAt: time.Now().AddDate(0, 0, -5),
			}
		})

		Context("when event exists", func() {
			It("should retrieve event successfully", func() {
				mockRepo.SetGetByIDBehavior(testEvent, nil)

				event, err := service.GetEventByID(ctx, "event-id-123")
				Expect(err).NotTo(HaveOccurred())

				Expect(event).NotTo(BeNil())
				Expect(event.ID).To(Equal(testEvent.ID))
				Expect(event.Text).To(Equal(testEvent.Text))
			})

			It("should return complete event data", func() {
				mockRepo.SetGetByIDBehavior(testEvent, nil)

				event, err := service.GetEventByID(ctx, "event-id-123")
				Expect(err).NotTo(HaveOccurred())

				Expect(event.ID).To(Equal("event-id-123"))
				Expect(event.Text).To(Equal("Test event"))
				Expect(event.Company).To(Equal("Test Company"))
				Expect(event.Tags).To(Equal([]string{"technical"}))
			})
		})

		Context("when event does not exist", func() {
			It("should return error if event not found", func() {
				mockRepo.SetGetByIDBehavior(nil, errors.New("event not found"))

				event, err := service.GetEventByID(ctx, "non-existent-id")
				Expect(err).To(HaveOccurred())
				Expect(event).To(BeNil())
			})
		})

		Context("with repository errors", func() {
			It("should return error when retrieval fails", func() {
				mockRepo.SetGetByIDBehavior(nil, errors.New("database error"))

				event, err := service.GetEventByID(ctx, "event-id-123")
				Expect(err).To(HaveOccurred())
				Expect(event).To(BeNil())
			})
		})

		Context("with empty event ID", func() {
			It("should attempt retrieval with empty ID", func() {
				mockRepo.SetGetByIDBehavior(nil, errors.New("invalid ID"))

				event, err := service.GetEventByID(ctx, "")
				Expect(err).To(HaveOccurred())
				Expect(event).To(BeNil())
			})
		})
	})

	Describe("ListEvents", func() {
		var testEvents []*career.Event

		BeforeEach(func() {
			testEvents = []*career.Event{
				{
					ID:   "event-1",
					Text: "Technical event",
					Tags: []string{"technical"},
				},
				{
					ID:   "event-2",
					Text: "Leadership event",
					Tags: []string{"leadership"},
				},
				{
					ID:   "event-3",
					Text: "Project event",
					Tags: []string{"project"},
				},
			}
		})

		Context("with no filters", func() {
			It("should list all events", func() {
				mockRepo.SetListBehavior(testEvents, nil)

				events, err := service.ListEvents(ctx, mocks.EventListFilters{})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(events)).To(Equal(3))
			})
		})

		Context("with tag filters", func() {
			It("should list events matching tags", func() {
				filteredEvents := []*career.Event{testEvents[0]}
				filters := mocks.EventListFilters{
					Tags: []string{"technical"},
				}
				mockRepo.SetListBehavior(filteredEvents, nil)

				events, err := service.ListEvents(ctx, filters)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(events)).To(Equal(1))
				Expect(events[0].Tags).To(ContainElement("technical"))
			})
		})

		Context("with pagination", func() {
			It("should list events with limit and offset", func() {
				paginatedEvents := []*career.Event{testEvents[0], testEvents[1]}
				filters := mocks.EventListFilters{
					Limit:  2,
					Offset: 0,
				}
				mockRepo.SetListBehavior(paginatedEvents, nil)

				events, err := service.ListEvents(ctx, filters)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(events)).To(Equal(2))
			})
		})

		Context("when no events match filters", func() {
			It("should return empty list", func() {
				mockRepo.SetListBehavior([]*career.Event{}, nil)

				events, err := service.ListEvents(ctx, mocks.EventListFilters{})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(events)).To(Equal(0))
			})
		})

		Context("with repository errors", func() {
			It("should return error when list fails", func() {
				mockRepo.SetListBehavior(nil, errors.New("database error"))

				events, err := service.ListEvents(ctx, mocks.EventListFilters{})
				Expect(err).To(HaveOccurred())
				Expect(events).To(BeNil())
			})
		})

		Context("with date range filters", func() {
			It("should list events within date range", func() {
				startDate := time.Now().AddDate(0, 0, -20)
				endDate := time.Now()
				filters := mocks.EventListFilters{
					StartDate: &startDate,
					EndDate:   &endDate,
				}
				mockRepo.SetListBehavior(testEvents, nil)

				events, err := service.ListEvents(ctx, filters)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(events)).To(BeNumerically(">", 0))
			})
		})

		Context("with sorting options", func() {
			It("should list events with sort options", func() {
				filters := mocks.EventListFilters{
					SortBy:    "date",
					SortOrder: "desc",
				}
				mockRepo.SetListBehavior(testEvents, nil)

				events, err := service.ListEvents(ctx, filters)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(events)).To(Equal(3))
			})
		})
	})

	Describe("CountEvents", func() {
		Context("with no filters", func() {
			It("should return total event count", func() {
				mockRepo.SetCountBehavior(42, nil)

				count, err := service.CountEvents(ctx, mocks.EventListFilters{})
				Expect(err).NotTo(HaveOccurred())

				Expect(count).To(Equal(42))
			})
		})

		Context("with tag filters", func() {
			It("should return count of events matching tags", func() {
				filters := mocks.EventListFilters{
					Tags: []string{"technical"},
				}
				mockRepo.SetCountBehavior(10, nil)

				count, err := service.CountEvents(ctx, filters)
				Expect(err).NotTo(HaveOccurred())

				Expect(count).To(Equal(10))
			})
		})

		Context("when no events match filters", func() {
			It("should return zero count", func() {
				mockRepo.SetCountBehavior(0, nil)

				count, err := service.CountEvents(ctx, mocks.EventListFilters{})
				Expect(err).NotTo(HaveOccurred())

				Expect(count).To(Equal(0))
			})
		})

		Context("with repository errors", func() {
			It("should return error when count fails", func() {
				mockRepo.SetCountBehavior(0, errors.New("database error"))

				count, err := service.CountEvents(ctx, mocks.EventListFilters{})
				Expect(err).To(HaveOccurred())

				Expect(count).To(Equal(0))
			})
		})

		Context("with date range filters", func() {
			It("should count events within date range", func() {
				startDate := time.Now().AddDate(0, 0, -30)
				endDate := time.Now()
				filters := mocks.EventListFilters{
					StartDate: &startDate,
					EndDate:   &endDate,
				}
				mockRepo.SetCountBehavior(15, nil)

				count, err := service.CountEvents(ctx, filters)
				Expect(err).NotTo(HaveOccurred())

				Expect(count).To(Equal(15))
			})
		})
	})

	Describe("Service initialization", func() {
		Context("NewService", func() {
			It("should create service with provided repository", func() {
				newService := NewService(mockRepo)
				Expect(newService).NotTo(BeNil())
				Expect(newService.repo).To(Equal(mockRepo))
			})

			It("should initialize with default logger", func() {
				newService := NewService(mockRepo)
				Expect(newService.logger).NotTo(BeNil())
			})
		})
	})

	Describe("Timestamp handling", func() {
		Context("when capturing events", func() {
			It("should set CreatedAt and UpdatedAt to same time", func() {
				mockRepo.SetCreateBehavior(nil)
				testEvent := &career.Event{
					Text: "Test event",
					Date: time.Now().AddDate(0, 0, -5),
				}

				beforeCapture := time.Now()
				err := service.CaptureEvent(ctx, testEvent, ManualEntry)
				afterCapture := time.Now()

				Expect(err).NotTo(HaveOccurred())
				Expect(testEvent.CreatedAt).To(BeTemporally(">=", beforeCapture))
				Expect(testEvent.CreatedAt).To(BeTemporally("<=", afterCapture))
				Expect(testEvent.UpdatedAt).To(Equal(testEvent.CreatedAt))
			})
		})

		Context("when updating events", func() {
			It("should preserve CreatedAt but update UpdatedAt", func() {
				existingEvent := &career.Event{
					ID:        "event-id",
					Text:      "Original",
					Date:      time.Now().AddDate(0, 0, -10),
					CreatedAt: time.Now().AddDate(0, 0, -5),
					UpdatedAt: time.Now().AddDate(0, 0, -5),
				}
				mockRepo.SetGetByIDBehavior(existingEvent, nil)
				mockRepo.SetUpdateBehavior(nil)

				updateEvent := &career.Event{
					ID:   "event-id",
					Text: "Updated",
					Date: time.Now().AddDate(0, 0, -10),
				}

				originalCreatedAt := existingEvent.CreatedAt
				time.Sleep(10 * time.Millisecond)

				err := service.UpdateEvent(ctx, updateEvent)
				Expect(err).NotTo(HaveOccurred())

				Expect(updateEvent.CreatedAt).To(Equal(originalCreatedAt))
				Expect(updateEvent.UpdatedAt).To(BeTemporally(">", originalCreatedAt))
			})
		})
	})

	Describe("Mode-specific validation", func() {
		Context("TimelineJournaling mode", func() {
			It("should accept events within 30 days", func() {
				mockRepo.SetCreateBehavior(nil)
				event := &career.Event{
					Text: "Recent event",
					Date: time.Now().AddDate(0, 0, -15),
				}

				err := service.CaptureEvent(ctx, event, TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject events older than 30 days", func() {
				event := &career.Event{
					Text: "Old event",
					Date: time.Now().AddDate(0, 0, -40),
				}

				err := service.CaptureEvent(ctx, event, TimelineJournaling)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("30 days"))
			})
		})

		Context("CVBackfill mode", func() {
			It("should accept very old events", func() {
				mockRepo.SetCreateBehavior(nil)
				event := &career.Event{
					Text: "Very old event",
					Date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				}

				err := service.CaptureEvent(ctx, event, CVBackfill)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("ManualEntry mode", func() {
			It("should accept any valid date", func() {
				mockRepo.SetCreateBehavior(nil)
				event := &career.Event{
					Text: "Event from any time",
					Date: time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC),
				}

				err := service.CaptureEvent(ctx, event, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	// NOTE: Competency classification is handled by:
	// - internal/service/career/classification/classifier.go (with comprehensive tests)
	// - internal/service/career/burstfact/extractor.go (with comprehensive tests)
	// See those packages for competency-related tests.
})
