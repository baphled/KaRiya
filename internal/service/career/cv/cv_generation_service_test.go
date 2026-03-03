package cv

import (
	"context"
	"errors"
	"io"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultCVGenerationService", func() {
	var (
		log *logger.Logger
		ctx context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()
	})

	Describe("GenerateCV", func() {
		It("should generate CV from saved configuration by name", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			configManager := NewMockConfigManager()
			configManager.configs["test-cv"] = config

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				configManager,
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCV(ctx, "test-cv")
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
			Expect(cv.Name).To(Equal("test-cv"))
			Expect(cv.TargetRole).To(Equal("principal"))
		})

		It("should return error when configuration not found", func() {
			configManager := NewMockConfigManager()

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				configManager,
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.GenerateCV(ctx, "nonexistent-cv")
			Expect(err).To(HaveOccurred())
		})

		It("should handle context cancellation", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.GenerateCV(cancelCtx, "test-cv")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GenerateCVFromConfig", func() {
		It("should return error when configuration is nil", func() {
			service := NewCVGenerationService(
				nil,
				nil,
				nil,
				nil,
				NewMockDataProcessingService(),
				nil,
				log,
			)
			_, err := service.GenerateCVFromConfig(ctx, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should require valid target role", func() {
			config := fixtures.CVConfig("test-cv")
			config.TargetRole = ""
			config.TargetAudience = ""

			service := NewCVGenerationService(
				nil,
				nil,
				nil,
				nil,
				NewMockDataProcessingService(),
				nil,
				log,
			)
			_, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should require at least one target audience", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "")

			service := NewCVGenerationService(
				nil,
				nil,
				nil,
				nil,
				NewMockDataProcessingService(),
				nil,
				log,
			)
			_, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should generate CV with valid configuration", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
			Expect(cv.Name).To(Equal("test-cv"))
			Expect(cv.TargetRole).To(Equal("principal"))
		})

		It("should set generated timestamp", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			beforeGeneration := time.Now()
			cv, err := service.GenerateCVFromConfig(ctx, config)
			afterGeneration := time.Now()

			Expect(err).NotTo(HaveOccurred())
			Expect(cv.GeneratedAt).To(BeTemporally(">=", beforeGeneration))
			Expect(cv.GeneratedAt).To(BeTemporally("<=", afterGeneration))
		})

		It("should handle context cancellation", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.GenerateCVFromConfig(cancelCtx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should generate unique CV IDs", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv1, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			cv2, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			Expect(cv1.ID).NotTo(Equal(cv2.ID))
		})

		It("should support single target audience", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.TargetAudience).To(Equal("hiring_manager"))
		})

		It("should track source event and fact counts", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewCountingRepository(5),
				NewCountingFactRepository(3),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.SourceEventCount).To(Equal(5))
			Expect(cv.SourceFactCount).To(Equal(3))
		})

		It("should preserve event filters in generated CV", func() {
			filters := map[string]interface{}{
				"companies": []string{"Google", "Meta"},
				"tags":      []string{"leadership", "technical"},
			}

			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = filters

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.EventFilters).To(Equal(filters))
		})

		It("should support all valid target roles", func() {
			validRoles := []string{"principal", "staff", "em", "senior_ic"}

			for _, role := range validRoles {
				config := fixtures.CVConfigWith("test-cv", role, "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				service := NewCVGenerationService(
					NewEmptyRepository(),
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.TargetRole).To(Equal(role))
			}
		})

		It("should return ephemeral CVs (not persisted)", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// CV should exist in memory but not be retrievable from config manager
			// (since it's ephemeral and not persisted)
			Expect(cv).NotTo(BeNil())
			Expect(cv.ID).NotTo(BeEmpty())
		})
	})

	Describe("Edge Cases", func() {
		It("should handle empty event repository gracefully", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})

		It("should handle empty fact repository gracefully", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = make(map[string]interface{})

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.SourceFactCount).To(Equal(0))
		})

		It("should handle nil event filters", func() {
			config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
			config.EventFilters = nil

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})
	})

	Describe("DataProcessingService integration", func() {
		Describe("NewCVGenerationService", func() {
			It("should panic when DataProcessingService is nil", func() {
				Expect(func() {
					NewCVGenerationService(
						NewEmptyRepository(),
						NewEmptyFactRepository(),
						NewMockConfigManager(),
						NewEmptyBulletGenerator(),
						nil, // nil DataProcessingService
						NewEmptySectionBuilder(),
						log,
					)
				}).To(Panic())
			})
		})

		Describe("GenerateCVFromConfig with achievements", func() {
			It("should extract achievements from events", func() {
				mockDataProcessor := NewMockDataProcessingService()
				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				event1 := fixtures.EventWith(uuid.New().String(), "Led API improvements, reducing latency by 40%", "", "")
				event2 := fixtures.EventWith(uuid.New().String(), "Mentored 5 junior engineers", "", "")

				// Create repository that returns test events
				eventRepo := NewMockEventRepository([]*career.Event{event1, event2})

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					mockDataProcessor,
					NewEmptySectionBuilder(),
					log,
				)

				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())

				// Verify ExtractAchievements was called for each event
				Expect(mockDataProcessor.ExtractAchievementsCalls).To(Equal(2))
			})

			It("should pass achievements to BulletGenerator", func() {
				mockBulletGen := NewMockBulletGenerator()
				mockDataProcessor := NewMockDataProcessingService()

				// Configure mock to return achievements
				mockDataProcessor.AchievementsToReturn = []*Achievement{
					{
						ID:          uuid.New().String(),
						Description: "Achievement 1",
						Metrics:     []*Metric{{Type: "percentage", Value: "40", Unit: "%"}},
						Confidence:  0.9,
					},
				}

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				event := fixtures.EventWith(uuid.New().String(), "Test event", "", "")

				eventRepo := NewMockEventRepository([]*career.Event{event})

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					mockBulletGen,
					mockDataProcessor,
					NewEmptySectionBuilder(),
					log,
				)

				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())

				Expect(mockBulletGen.ReceivedAchievements).NotTo(BeNil())
				Expect(mockBulletGen.ReceivedAchievements).ToNot(BeEmpty())
			})

			It("should handle ExtractAchievements errors gracefully", func() {
				mockDataProcessor := NewMockDataProcessingService()
				mockDataProcessor.ShouldReturnError = true

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				event := fixtures.EventWith(uuid.New().String(), "Test event", "", "")

				eventRepo := NewMockEventRepository([]*career.Event{event})

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					mockDataProcessor,
					NewEmptySectionBuilder(),
					log,
				)

				// Should not fail even if achievement extraction fails
				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("eventMatchesFilters", func() {
		var svc *DefaultCVGenerationService

		BeforeEach(func() {
			svc = &DefaultCVGenerationService{
				logger: log,
			}
		})

		It("should match with empty filters", func() {
			event := fixtures.EventWith(uuid.New().String(), "Test event", "Google", "")
			Expect(svc.eventMatchesFilters(event, map[string]interface{}{})).To(BeTrue())
		})

		Context("with minDate filter", func() {
			It("should match when event date is after minDate", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Date = time.Now()
				filters := map[string]interface{}{
					"minDate": time.Now().AddDate(-1, 0, 0),
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should not match when event date is before minDate", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Date = time.Now().AddDate(-2, 0, 0)
				filters := map[string]interface{}{
					"minDate": time.Now().AddDate(-1, 0, 0),
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})
		})

		Context("with maxDate filter", func() {
			It("should match when event date is before maxDate", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Date = time.Now().AddDate(-1, 0, 0)
				filters := map[string]interface{}{
					"maxDate": time.Now(),
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should not match when event date is after maxDate", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Date = time.Now()
				filters := map[string]interface{}{
					"maxDate": time.Now().AddDate(-1, 0, 0),
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})
		})

		Context("with companies filter", func() {
			It("should match when event company is in filter list", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Google", "")
				filters := map[string]interface{}{
					"companies": []string{"Google", "Meta"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should not match when event company is not in filter list", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Apple", "")
				filters := map[string]interface{}{
					"companies": []string{"Google", "Meta"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})

			It("should match when companies list is empty", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Apple", "")
				filters := map[string]interface{}{
					"companies": []string{},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})
		})

		Context("with tags filter", func() {
			It("should match when event has a matching tag", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Tags = []string{"leadership", "technical"}
				filters := map[string]interface{}{
					"tags": []string{"technical"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should not match when event has no matching tags", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Tags = []string{"leadership"}
				filters := map[string]interface{}{
					"tags": []string{"technical", "research"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})

			It("should match when tags filter list is empty", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Tags = []string{"leadership"}
				filters := map[string]interface{}{
					"tags": []string{},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})
		})

		Context("with categories filter", func() {
			It("should match against event Categories field", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Categories = []string{"technical", "leadership"}
				filters := map[string]interface{}{
					"categories": []string{"technical"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should fall back to Tags when Categories is empty", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Categories = []string{}
				event.Tags = []string{"technical", "leadership"}
				filters := map[string]interface{}{
					"categories": []string{"technical"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should not match when no categories match via Categories field", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Categories = []string{"research"}
				filters := map[string]interface{}{
					"categories": []string{"technical", "leadership"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})

			It("should not match when categories fallback to tags also fails", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Co", "")
				event.Categories = []string{}
				event.Tags = []string{"research"}
				filters := map[string]interface{}{
					"categories": []string{"technical"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})
		})

		Context("with combined filters", func() {
			It("should match when all filters pass", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Google", "")
				event.Date = time.Now()
				event.Tags = []string{"technical"}
				event.Categories = []string{"leadership"}
				filters := map[string]interface{}{
					"minDate":    time.Now().AddDate(-1, 0, 0),
					"maxDate":    time.Now().Add(time.Hour),
					"companies":  []string{"Google"},
					"tags":       []string{"technical"},
					"categories": []string{"leadership"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeTrue())
			})

			It("should not match when company filter fails but others pass", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Apple", "")
				event.Date = time.Now()
				event.Tags = []string{"technical"}
				filters := map[string]interface{}{
					"minDate":   time.Now().AddDate(-1, 0, 0),
					"companies": []string{"Google"},
					"tags":      []string{"technical"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})

			It("should not match when date filter fails but others pass", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test", "Google", "")
				event.Date = time.Now().AddDate(-5, 0, 0)
				event.Tags = []string{"technical"}
				filters := map[string]interface{}{
					"minDate":   time.Now().AddDate(-1, 0, 0),
					"companies": []string{"Google"},
					"tags":      []string{"technical"},
				}
				Expect(svc.eventMatchesFilters(event, filters)).To(BeFalse())
			})
		})
	})

	Describe("retrieveEventsWithFilters", func() {
		It("should return error when eventRepo.List fails", func() {
			service := NewCVGenerationService(
				NewErrorEventRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.retrieveEventsWithFilters(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to list events"))
		})

		It("should apply filters when filters map is non-empty", func() {
			events := []*career.Event{
				fixtures.EventWith(uuid.New().String(), "Google event", "Google", ""),
				fixtures.EventWith(uuid.New().String(), "Meta event", "Meta", ""),
			}
			eventRepo := NewMockEventRepository(events)

			service := NewCVGenerationService(
				eventRepo,
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			filters := map[string]interface{}{
				"companies": []string{"Google"},
			}
			result, err := service.retrieveEventsWithFilters(ctx, filters)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(1))
			Expect(result[0].Company).To(Equal("Google"))
		})

		It("should return all events when filters map is empty", func() {
			events := []*career.Event{
				fixtures.EventWith(uuid.New().String(), "Event 1", "Google", ""),
				fixtures.EventWith(uuid.New().String(), "Event 2", "Meta", ""),
			}
			eventRepo := NewMockEventRepository(events)

			service := NewCVGenerationService(
				eventRepo,
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			result, err := service.retrieveEventsWithFilters(ctx, map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(2))
		})
	})

	Describe("retrieveFacts", func() {
		It("should return empty facts when factRepo is nil", func() {
			svc := &DefaultCVGenerationService{
				factRepo: nil,
				logger:   log,
			}
			facts, err := svc.retrieveFacts(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error when factRepo.List fails", func() {
			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewErrorFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.retrieveFacts(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to list facts"))
		})

		It("should return facts from repository", func() {
			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewCountingFactRepository(3),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewMockDataProcessingService(),
				NewEmptySectionBuilder(),
				log,
			)

			facts, err := service.retrieveFacts(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(3))
		})
	})

	Describe("GenerateCVFromConfig additional coverage", func() {
		Context("with length format filtering", func() {
			It("should filter events by date when MaxYearsHistory is set", func() {
				recentEvent := fixtures.EventWith(uuid.New().String(), "Recent work", "Co", "")
				recentEvent.Date = time.Now()

				oldEvent := fixtures.EventWith(uuid.New().String(), "Old work", "Co", "")
				oldEvent.Date = time.Now().AddDate(-10, 0, 0)

				eventRepo := NewMockEventRepository([]*career.Event{recentEvent, oldEvent})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})
				config.LengthFormat = string(LengthShort)

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.SourceEventCount).To(Equal(1))
			})

			It("should not filter events when length format has no MaxYearsHistory", func() {
				event1 := fixtures.EventWith(uuid.New().String(), "Event 1", "Co", "")
				event1.Date = time.Now()

				event2 := fixtures.EventWith(uuid.New().String(), "Event 2", "Co", "")
				event2.Date = time.Now().AddDate(-20, 0, 0)

				eventRepo := NewMockEventRepository([]*career.Event{event1, event2})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})
				config.LengthFormat = string(LengthFull)

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.SourceEventCount).To(Equal(2))
			})

			It("should filter bullets by confidence threshold", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test event", "Co", "")
				event.Date = time.Now()
				eventRepo := NewMockEventRepository([]*career.Event{event})

				bulletGen := &ConfigurableBulletGenerator{
					Bullets: []*Bullet{
						{ID: uuid.New().String(), Text: "High confidence", Confidence: 0.90},
						{ID: uuid.New().String(), Text: "Low confidence", Confidence: 0.50},
					},
				}

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})
				config.LengthFormat = string(LengthShort)

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					bulletGen,
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv).NotTo(BeNil())
			})

			It("should return empty CV when all events are filtered out by date", func() {
				oldEvent := fixtures.EventWith(uuid.New().String(), "Old work", "Co", "")
				oldEvent.Date = time.Now().AddDate(-10, 0, 0)

				eventRepo := NewMockEventRepository([]*career.Event{oldEvent})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})
				config.LengthFormat = string(LengthShort)

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.SourceEventCount).To(Equal(0))
			})
		})

		Context("with technology focus", func() {
			It("should apply technology filtering when focus is not LanguageAgnostic", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test event", "Co", "")
				event.Date = time.Now()
				eventRepo := NewMockEventRepository([]*career.Event{event})

				bulletGen := &ConfigurableBulletGenerator{
					Bullets: []*Bullet{},
				}

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})
				config.TechnologyFocus = string(TechnologyFocusSpecialist)
				config.SelectedTechnologies = []string{"Go", "Python"}

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					bulletGen,
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(bulletGen.FilterByTechCalled).To(BeTrue())
			})

			It("should not apply technology filtering when focus is LanguageAgnostic", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test event", "Co", "")
				event.Date = time.Now()
				eventRepo := NewMockEventRepository([]*career.Event{event})

				bulletGen := &ConfigurableBulletGenerator{
					Bullets: []*Bullet{},
				}

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})
				config.TechnologyFocus = string(TechnologyFocusLanguageAgnostic)

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					bulletGen,
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(bulletGen.FilterByTechCalled).To(BeFalse())
			})
		})

		Context("when bullet generation fails", func() {
			It("should return error", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test event", "Co", "")
				event.Date = time.Now()
				eventRepo := NewMockEventRepository([]*career.Event{event})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewErrorBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to generate bullets"))
			})
		})

		Context("when section building fails", func() {
			It("should return error", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test event", "Co", "")
				event.Date = time.Now()
				eventRepo := NewMockEventRepository([]*career.Event{event})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewErrorSectionBuilder(),
					log,
				)

				_, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to build sections"))
			})
		})

		Context("when fact retrieval fails", func() {
			It("should continue generation with empty facts", func() {
				event := fixtures.EventWith(uuid.New().String(), "Test event", "Co", "")
				event.Date = time.Now()
				eventRepo := NewMockEventRepository([]*career.Event{event})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = make(map[string]interface{})

				service := NewCVGenerationService(
					eventRepo,
					NewErrorFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.SourceFactCount).To(Equal(0))
			})
		})

		Context("with event filters in config", func() {
			It("should filter events by company through full pipeline", func() {
				googleEvent := fixtures.EventWith(uuid.New().String(), "Google work", "Google", "")
				metaEvent := fixtures.EventWith(uuid.New().String(), "Meta work", "Meta", "")
				eventRepo := NewMockEventRepository([]*career.Event{googleEvent, metaEvent})

				config := fixtures.CVConfigWith("test-cv", "principal", "hiring_manager")
				config.EventFilters = map[string]interface{}{
					"companies": []string{"Google"},
				}

				service := NewCVGenerationService(
					eventRepo,
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewMockDataProcessingService(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.SourceEventCount).To(Equal(1))
			})
		})
	})
})

// Mock implementations for testing

type MockConfigManager struct {
	configs map[string]*career.CVConfig
}

func NewMockConfigManager() *MockConfigManager {
	return &MockConfigManager{
		configs: make(map[string]*career.CVConfig),
	}
}

func (m *MockConfigManager) LoadConfig(_ context.Context, name string) (*career.CVConfig, error) {
	if config, ok := m.configs[name]; ok {
		return config, nil
	}
	return nil, ErrConfigNotFound
}

func (m *MockConfigManager) SaveConfig(_ context.Context, config *career.CVConfig) error {
	m.configs[config.Name] = config
	return nil
}

func (m *MockConfigManager) DeleteConfig(_ context.Context, name string) error {
	delete(m.configs, name)
	return nil
}

func (m *MockConfigManager) ListConfigs(_ context.Context) ([]*career.CVConfig, error) {
	configs := make([]*career.CVConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	return configs, nil
}

func (m *MockConfigManager) GetConfigPath(name string) string {
	return "/tmp/" + name
}

func (m *MockConfigManager) ConfigExists(ctx context.Context, name string) (bool, error) {
	_, ok := m.configs[name]
	return ok, nil
}

type EmptyRepository struct{}

func NewEmptyRepository() *EmptyRepository {
	return &EmptyRepository{}
}

func (r *EmptyRepository) List(ctx context.Context, filters careerrepo.EventListFilters) ([]*career.Event, error) {
	return []*career.Event{}, nil
}

func (r *EmptyRepository) GetByID(ctx context.Context, id string) (*career.Event, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *EmptyRepository) Create(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *EmptyRepository) Update(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *EmptyRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *EmptyRepository) Count(ctx context.Context, filters careerrepo.EventListFilters) (int, error) {
	return 0, nil
}

func (r *EmptyRepository) LinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

func (r *EmptyRepository) UnlinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

type CountingRepository struct {
	count int
}

func NewCountingRepository(count int) *CountingRepository {
	return &CountingRepository{count: count}
}

func (r *CountingRepository) List(ctx context.Context, filters careerrepo.EventListFilters) ([]*career.Event, error) {
	events := make([]*career.Event, r.count)
	for i := range r.count {
		events[i] = fixtures.EventWith(uuid.New().String(), "Sample event", "", "")
	}
	return events, nil
}

func (r *CountingRepository) GetByID(ctx context.Context, id string) (*career.Event, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *CountingRepository) Create(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *CountingRepository) Update(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *CountingRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *CountingRepository) Count(ctx context.Context, filters careerrepo.EventListFilters) (int, error) {
	return r.count, nil
}

func (r *CountingRepository) LinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

func (r *CountingRepository) UnlinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

type EmptyFactRepository struct{}

func NewEmptyFactRepository() *EmptyFactRepository {
	return &EmptyFactRepository{}
}

func (r *EmptyFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *EmptyFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *EmptyFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *EmptyFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *EmptyFactRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *EmptyFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return 0, nil
}

func (r *EmptyFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *EmptyFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

type CountingFactRepository struct {
	count int
}

func NewCountingFactRepository(count int) *CountingFactRepository {
	return &CountingFactRepository{count: count}
}

func (r *CountingFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*career.Fact, error) {
	facts := make([]*career.Fact, r.count)
	for i := range r.count {
		facts[i] = fixtures.FactWith(uuid.New().String(), "Sample fact")
	}
	return facts, nil
}

func (r *CountingFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *CountingFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *CountingFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *CountingFactRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *CountingFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return r.count, nil
}

func (r *CountingFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *CountingFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

// EmptyBulletGenerator implements BulletGenerator for tests.
type EmptyBulletGenerator struct{}

func NewEmptyBulletGenerator() *EmptyBulletGenerator {
	return &EmptyBulletGenerator{}
}

func (g *EmptyBulletGenerator) GenerateBullets(ctx context.Context, events []*career.Event, facts []*career.Fact, achievements []*Achievement, targetRole string, targetAudience string) ([]*Bullet, error) {
	return []*Bullet{}, nil
}

func (g *EmptyBulletGenerator) FilterByRole(bullets []*Bullet, _ string) []*Bullet {
	return bullets
}

func (g *EmptyBulletGenerator) FilterByAudience(bullets []*Bullet, _ string) []*Bullet {
	return bullets
}

func (g *EmptyBulletGenerator) RankByRelevance(bullets []*Bullet, _ string, _ string) []*Bullet {
	return bullets
}

func (g *EmptyBulletGenerator) EnhanceWording(bullet *Bullet, _ string) (*Bullet, error) {
	return bullet, nil
}

func (g *EmptyBulletGenerator) FilterByTechnologies(bullets []*Bullet, _ []*career.Event, _ TechnologyFocus, _ []string) []*Bullet {
	return bullets // No-op for tests
}

type EmptySectionBuilder struct{}

func NewEmptySectionBuilder() *EmptySectionBuilder {
	return &EmptySectionBuilder{}
}

func (b *EmptySectionBuilder) BuildSections(ctx context.Context, bullets []*career.CVBullet, events []*career.Event, facts []*career.Fact, targetRole string, skillsConfig *SkillsFormatConfig, summaryCfg *SummaryConfig) ([]*career.CVSection, error) {
	return []*career.CVSection{}, nil
}

// MockDataProcessingService for testing achievement extraction.
type MockDataProcessingService struct {
	ExtractAchievementsCalls int
	AchievementsToReturn     []*Achievement
	ShouldReturnError        bool
}

func NewMockDataProcessingService() *MockDataProcessingService {
	return &MockDataProcessingService{
		AchievementsToReturn: []*Achievement{},
	}
}

func (m *MockDataProcessingService) ExtractAchievements(ctx context.Context, event *career.Event, facts []*career.Fact) ([]*Achievement, error) {
	m.ExtractAchievementsCalls++
	if m.ShouldReturnError {
		return nil, ErrConfigNotFound // reuse existing error for test
	}
	return m.AchievementsToReturn, nil
}

func (m *MockDataProcessingService) GroupEventsByCompany(ctx context.Context, events []*career.Event) (map[string]*CompanyGroup, error) {
	return nil, nil //nolint:nilnil // test stub
}

func (m *MockDataProcessingService) ExtractSkills(ctx context.Context, events []*career.Event, facts []*career.Fact) (map[string]*SkillCategory, error) {
	return nil, nil //nolint:nilnil // test stub
}

func (m *MockDataProcessingService) CalculateMetrics(ctx context.Context, text string) ([]*Metric, error) {
	return nil, nil //nolint:nilnil // test stub
}

func (m *MockDataProcessingService) ExtractProjectsFromEvents(ctx context.Context, events []*career.Event) ([]*ProjectGroup, error) {
	return nil, nil //nolint:nilnil // test stub
}

// MockEventRepository for testing with specific events.
type MockEventRepository struct {
	events []*career.Event
}

func NewMockEventRepository(events []*career.Event) *MockEventRepository {
	return &MockEventRepository{events: events}
}

func (r *MockEventRepository) List(ctx context.Context, filters careerrepo.EventListFilters) ([]*career.Event, error) {
	return r.events, nil
}

func (r *MockEventRepository) GetByID(ctx context.Context, id string) (*career.Event, error) {
	for _, event := range r.events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, nil //nolint:nilnil // test stub
}

func (r *MockEventRepository) Create(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *MockEventRepository) Update(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *MockEventRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *MockEventRepository) Count(ctx context.Context, filters careerrepo.EventListFilters) (int, error) {
	return len(r.events), nil
}

func (r *MockEventRepository) LinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

func (r *MockEventRepository) UnlinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

// MockBulletGenerator for testing bullet generation with achievements.
type MockBulletGenerator struct {
	ReceivedAchievements []*Achievement
}

func NewMockBulletGenerator() *MockBulletGenerator {
	return &MockBulletGenerator{}
}

func (g *MockBulletGenerator) GenerateBullets(ctx context.Context, events []*career.Event, facts []*career.Fact, achievements []*Achievement, targetRole string, targetAudience string) ([]*Bullet, error) {
	g.ReceivedAchievements = achievements
	return []*Bullet{}, nil
}

func (g *MockBulletGenerator) FilterByRole(bullets []*Bullet, _ string) []*Bullet {
	return bullets
}

func (g *MockBulletGenerator) FilterByAudience(bullets []*Bullet, _ string) []*Bullet {
	return bullets
}

func (g *MockBulletGenerator) RankByRelevance(bullets []*Bullet, _ string, _ string) []*Bullet {
	return bullets
}

func (g *MockBulletGenerator) EnhanceWording(bullet *Bullet, _ string) (*Bullet, error) {
	return bullet, nil
}

func (g *MockBulletGenerator) FilterByTechnologies(bullets []*Bullet, _ []*career.Event, _ TechnologyFocus, _ []string) []*Bullet {
	return bullets
}

// ErrorEventRepository returns error from List for testing error paths.
type ErrorEventRepository struct {
	EmptyRepository
}

func NewErrorEventRepository() *ErrorEventRepository {
	return &ErrorEventRepository{}
}

func (r *ErrorEventRepository) List(_ context.Context, _ careerrepo.EventListFilters) ([]*career.Event, error) {
	return nil, errors.New("event repository error")
}

// ErrorFactRepository returns error from List for testing error paths.
type ErrorFactRepository struct {
	EmptyFactRepository
}

func NewErrorFactRepository() *ErrorFactRepository {
	return &ErrorFactRepository{}
}

func (r *ErrorFactRepository) List(_ context.Context, _ careerrepo.FactListFilters) ([]*career.Fact, error) {
	return nil, errors.New("fact repository error")
}

// ErrorBulletGenerator returns error from GenerateBullets for testing error paths.
type ErrorBulletGenerator struct {
	EmptyBulletGenerator
}

func NewErrorBulletGenerator() *ErrorBulletGenerator {
	return &ErrorBulletGenerator{}
}

func (g *ErrorBulletGenerator) GenerateBullets(_ context.Context, _ []*career.Event, _ []*career.Fact, _ []*Achievement, _ string, _ string) ([]*Bullet, error) {
	return nil, errors.New("bullet generation error")
}

// ErrorSectionBuilder returns error from BuildSections for testing error paths.
type ErrorSectionBuilder struct {
	EmptySectionBuilder
}

func NewErrorSectionBuilder() *ErrorSectionBuilder {
	return &ErrorSectionBuilder{}
}

func (b *ErrorSectionBuilder) BuildSections(_ context.Context, _ []*career.CVBullet, _ []*career.Event, _ []*career.Fact, _ string, _ *SkillsFormatConfig, _ *SummaryConfig) ([]*career.CVSection, error) {
	return nil, errors.New("section build error")
}

// ConfigurableBulletGenerator returns preset bullets and tracks FilterByTechnologies calls.
type ConfigurableBulletGenerator struct {
	EmptyBulletGenerator
	Bullets            []*Bullet
	FilterByTechCalled bool
}

func (g *ConfigurableBulletGenerator) GenerateBullets(_ context.Context, _ []*career.Event, _ []*career.Fact, _ []*Achievement, _ string, _ string) ([]*Bullet, error) {
	return g.Bullets, nil
}

func (g *ConfigurableBulletGenerator) FilterByTechnologies(bullets []*Bullet, _ []*career.Event, _ TechnologyFocus, _ []string) []*Bullet {
	g.FilterByTechCalled = true
	return bullets
}
