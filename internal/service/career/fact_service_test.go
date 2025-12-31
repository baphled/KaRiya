package career

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Career Service - Fact Methods", func() {
	var (
		repo    careerrepo.Repository
		service *Service
		ctx     context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		service = NewService(repo)
		ctx = context.Background()
	})

	Describe("ExtractFactsFromEvent", func() {
		It("should extract facts from a simple event", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development team through critical project",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				Project:   "Platform",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())
			Expect(facts[0].SourceEventID).To(Equal(event.ID))
		})

		It("should validate extracted facts", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Architected enterprise platform",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			for _, fact := range facts {
				validationErr := fact.Validate()
				Expect(validationErr).NotTo(HaveOccurred())
			}
		})

		It("should infer competencies from event text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led engineering team and mentored junior developers",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("mentoring"))
		})

		It("should return empty list for event with empty text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for nil event", func() {
			facts, err := service.ExtractFactsFromEvent(ctx, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event cannot be nil"))
			Expect(facts).To(BeNil())
		})

		It("should return error for event with empty ID", func() {
			event := &career.CareerEvent{
				ID:        "",
				Text:      "Some event",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should extract role fit from event text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Managed engineering team and coordinated hiring",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitEM))
		})

		It("should extract audience relevance from event text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led critical infrastructure initiative",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].AudienceRelevance).To(ContainElement("peer"))
		})

		It("should extract strength signal from event text", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Delivered critical feature that improved performance",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].StrengthSignal).NotTo(BeEmpty())
		})
	})

	Describe("ExtractFactsFromBurst", func() {
		It("should extract facts from burst with multiple events", func() {
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development team",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}
			event2 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Architected microservices platform",
				Date:      time.Now().Add(-50 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-50 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

			// Store events in repository
			_ = service.CaptureEvent(ctx, event1, ManualEntry)
			_ = service.CaptureEvent(ctx, event2, ManualEntry)

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Platform Initiative",
				EventIDs:  []string{event1.ID, event2.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(facts)).To(BeNumerically(">=", 2))
			Expect(facts[0].SourceBurstID).To(Equal(burst.ID))
		})

		It("should validate all extracted burst facts", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Implemented critical feature",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			_ = service.CaptureEvent(ctx, event, ManualEntry)

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Feature Initiative",
				EventIDs:  []string{event.ID},
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			for _, fact := range facts {
				validationErr := fact.Validate()
				Expect(validationErr).NotTo(HaveOccurred())
			}
		})

		It("should return error for nil burst", func() {
			facts, err := service.ExtractFactsFromBurst(ctx, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst cannot be nil"))
			Expect(facts).To(BeNil())
		})

		It("should return error for burst with empty ID", func() {
			burst := &career.Burst{
				ID:        "",
				Name:      "Empty ID Burst",
				EventIDs:  []string{uuid.New().String()},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should return empty list for burst with no events", func() {
			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Empty Burst",
				EventIDs:  []string{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should handle missing events gracefully", func() {
			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Burst with Missing Events",
				EventIDs:  []string{uuid.New().String(), uuid.New().String()},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			// Should return empty list, not error (no events found)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should combine competencies from multiple events", func() {
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led development team",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}
			event2 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Mentored junior engineers",
				Date:      time.Now().Add(-50 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-50 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

			_ = service.CaptureEvent(ctx, event1, ManualEntry)
			_ = service.CaptureEvent(ctx, event2, ManualEntry)

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Team Growth",
				EventIDs:  []string{event1.ID, event2.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			// Burst-level fact should combine competencies
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("mentoring"))
		})
	})

	Describe("ValidateFact", func() {
		It("should validate a correct fact", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Led critical initiative",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              career.RoleFitPrincipal,
				AudienceRelevance:    []string{"peer", "hiring_manager"},
				StrengthSignal:       "leadership capability",
				SourceEventID:        uuid.New().String(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject nil fact", func() {
			err := service.ValidateFact(ctx, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact cannot be nil"))
		})

		It("should reject fact with empty text", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"peer"},
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with no competency categories", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Some achievement",
				CompetencyCategories: []string{},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"peer"},
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with invalid role fit", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Some achievement",
				CompetencyCategories: []string{"technical"},
				RoleFit:              "invalid_role",
				AudienceRelevance:    []string{"peer"},
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with no source reference", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Some achievement",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"peer"},
				SourceEventID:        "", // No source
				SourceBurstID:        "", // No source
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with aspirational language", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Will lead critical initiatives",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              career.RoleFitPrincipal,
				AudienceRelevance:    []string{"peer"},
				SourceEventID:        uuid.New().String(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with no audience relevance", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Some achievement",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{}, // No audience
				SourceEventID:        uuid.New().String(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should accept fact with burst source reference", func() {
			fact := &career.Fact{
				ID:                   uuid.New().String(),
				Text:                 "Multi-phase initiative success",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              career.RoleFitPrincipal,
				AudienceRelevance:    []string{"peer", "hiring_manager"},
				SourceBurstID:        uuid.New().String(),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}

			err := service.ValidateFact(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Integration: Event to Fact to Validation", func() {
		It("should successfully extract and validate facts from event", func() {
			event := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led platform architecture project and mentored team",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			}

			// Extract facts
			facts, err := service.ExtractFactsFromEvent(ctx, event)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())

			// Validate each fact
			for _, fact := range facts {
				err := service.ValidateFact(ctx, &fact)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should successfully extract and validate facts from burst", func() {
			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Designed system architecture",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-60 * 24 * time.Hour),
			}
			event2 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led implementation and delivered solution",
				Date:      time.Now().Add(-50 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-50 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-50 * 24 * time.Hour),
			}

			_ = service.CaptureEvent(ctx, event1, ManualEntry)
			_ = service.CaptureEvent(ctx, event2, ManualEntry)

			burst := &career.Burst{
				ID:        uuid.New().String(),
				Name:      "System Architecture Initiative",
				EventIDs:  []string{event1.ID, event2.ID},
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			// Extract facts
			facts, err := service.ExtractFactsFromBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())

			// Validate each fact
			for _, fact := range facts {
				err := service.ValidateFact(ctx, &fact)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Describe("SetFactRepository", func() {
		It("should set the fact repository", func() {
			factRepo := careerrepo.NewMemoryFactRepository()
			service.SetFactRepository(factRepo)

			// Verify by creating and retrieving a fact
			fact := &career.Fact{
				Text:                 "Test fact for repository",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"peer"},
				StrengthSignal:       "technical",
				SourceEventID:        uuid.New().String(),
			}

			err := service.SaveFact(ctx, fact)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := factRepo.GetByID(ctx, fact.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(fact.ID))
		})
	})

	Describe("GetFactsBySourceEventID", func() {
		var factRepo *careerrepo.MemoryFactRepository
		var event *career.CareerEvent

		BeforeEach(func() {
			factRepo = careerrepo.NewMemoryFactRepository()
			service.SetFactRepository(factRepo)

			event = &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led migration to microservices",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			// Extract and save facts
			facts, err := service.ExtractFactsFromEvent(ctx, event)
			Expect(err).NotTo(HaveOccurred())
			for _, fact := range facts {
				err = service.SaveFact(ctx, &fact)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should retrieve facts for a specific event", func() {
			facts, err := service.GetFactsBySourceEventID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())
			for _, fact := range facts {
				Expect(fact.SourceEventID).To(Equal(event.ID))
			}
		})

		It("should return empty list for event with no facts", func() {
			facts, err := service.GetFactsBySourceEventID(ctx, "non-existent-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for empty event ID", func() {
			_, err := service.GetFactsBySourceEventID(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
		})
	})

	Describe("GetFactsBySourceBurstID", func() {
		var factRepo *careerrepo.MemoryFactRepository
		var burst *career.Burst

		BeforeEach(func() {
			factRepo = careerrepo.NewMemoryFactRepository()
			service.SetFactRepository(factRepo)

			event1 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Led platform migration",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}
			event2 := &career.CareerEvent{
				ID:        uuid.New().String(),
				Text:      "Architected microservices",
				Date:      time.Now().Add(-60 * 24 * time.Hour),
				CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			}

			// Save events to repository first (required for burst fact extraction)
			err := service.CaptureEvent(ctx, event1, ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			err = service.CaptureEvent(ctx, event2, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			burst = &career.Burst{
				ID:              uuid.New().String(),
				Name:            "Platform Migration",
				EventIDs:        []string{event1.ID, event2.ID},
				CompetencyFocus: "technical",
				CreatedAt:       time.Now().Add(-60 * 24 * time.Hour),
				UpdatedAt:       time.Now(),
			}

			// Extract and save facts
			facts, err := service.ExtractFactsFromBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
			for _, fact := range facts {
				err = service.SaveFact(ctx, &fact)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should retrieve facts for a specific burst", func() {
			facts, err := service.GetFactsBySourceBurstID(ctx, burst.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())
			for _, fact := range facts {
				Expect(fact.SourceBurstID).To(Equal(burst.ID))
			}
		})

		It("should return empty list for burst with no facts", func() {
			facts, err := service.GetFactsBySourceBurstID(ctx, "non-existent-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for empty burst ID", func() {
			_, err := service.GetFactsBySourceBurstID(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
		})
	})

	Describe("SaveFact", func() {
		var factRepo *careerrepo.MemoryFactRepository

		BeforeEach(func() {
			factRepo = careerrepo.NewMemoryFactRepository()
			service.SetFactRepository(factRepo)
		})

		It("should save a valid fact", func() {
			fact := &career.Fact{
				Text:                 "Led migration to microservices",
				CompetencyCategories: []string{"technical", "leadership"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"hiring_manager", "peer"},
				StrengthSignal:       "leadership",
				SourceEventID:        uuid.New().String(),
			}

			err := service.SaveFact(ctx, fact)
			Expect(err).NotTo(HaveOccurred())
			Expect(fact.ID).NotTo(BeEmpty())
			Expect(fact.CreatedAt).NotTo(BeZero())
			Expect(fact.UpdatedAt).NotTo(BeZero())
		})

		It("should return error for nil fact", func() {
			err := service.SaveFact(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact cannot be nil"))
		})

		It("should return error for invalid fact", func() {
			fact := &career.Fact{
				Text:                 "", // Empty text
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"peer"},
				StrengthSignal:       "technical",
				SourceEventID:        uuid.New().String(),
			}

			err := service.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
		})

		It("should return error for fact with aspirational language", func() {
			fact := &career.Fact{
				Text:                 "I will lead migration",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"peer"},
				StrengthSignal:       "technical",
				SourceEventID:        uuid.New().String(),
			}

			err := service.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aspirational language"))
		})

		It("should update existing fact if ID is set", func() {
			fact := &career.Fact{
				Text:                 "Original text",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"peer"},
				StrengthSignal:       "technical",
				SourceEventID:        uuid.New().String(),
			}

			err := service.SaveFact(ctx, fact)
			Expect(err).NotTo(HaveOccurred())
			originalID := fact.ID

			fact.Text = "Updated text"
			err = service.SaveFact(ctx, fact)
			Expect(err).NotTo(HaveOccurred())

			Expect(fact.ID).To(Equal(originalID))
			retrieved, err := factRepo.GetByID(ctx, fact.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.Text).To(Equal("Updated text"))
		})

		It("should return error when repository is not set", func() {
			serviceWithoutRepo := NewService(repo)
			fact := &career.Fact{
				Text:                 "Test fact",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"peer"},
				StrengthSignal:       "technical",
				SourceEventID:        uuid.New().String(),
			}

			err := serviceWithoutRepo.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact repository not configured"))
		})
	})

	Describe("DeleteFact", func() {
		var factRepo *careerrepo.MemoryFactRepository
		var fact *career.Fact

		BeforeEach(func() {
			factRepo = careerrepo.NewMemoryFactRepository()
			service.SetFactRepository(factRepo)

			fact = &career.Fact{
				Text:                 "Fact to delete",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitStaff,
				AudienceRelevance:    []string{"peer"},
				StrengthSignal:       "technical",
				SourceEventID:        uuid.New().String(),
			}
			err := service.SaveFact(ctx, fact)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete an existing fact", func() {
			err := service.DeleteFact(ctx, fact.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = factRepo.GetByID(ctx, fact.ID)
			Expect(err).To(MatchError(careerrepo.ErrFactNotFound))
		})

		It("should return error for empty fact ID", func() {
			err := service.DeleteFact(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact ID cannot be empty"))
		})

		It("should return error for non-existent fact", func() {
			err := service.DeleteFact(ctx, "non-existent-id")
			Expect(err).To(MatchError(careerrepo.ErrFactNotFound))
		})

		It("should return error when repository is not set", func() {
			serviceWithoutRepo := NewService(repo)
			err := serviceWithoutRepo.DeleteFact(ctx, fact.ID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact repository not configured"))
		})
	})

	Describe("GetFactRepository", func() {
		It("should return nil when no fact repository is set", func() {
			factRepo := service.GetFactRepository()
			Expect(factRepo).To(BeNil())
		})

		It("should return the configured fact repository", func() {
			memoryFactRepo := careerrepo.NewMemoryFactRepository()
			service.SetFactRepository(memoryFactRepo)

			factRepo := service.GetFactRepository()
			Expect(factRepo).To(Equal(memoryFactRepo))
		})
	})
})
