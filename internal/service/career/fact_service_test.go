//nolint:errcheck // Test file - error handling for test setup is not relevant.
package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Career Service - Fact Methods", func() {
	var (
		repo    careerrepo.EventRepository
		service *Service
		ctx     context.Context
	)

	BeforeEach(func() {
		repo = careermemory.NewEventRepository()
		service = NewService(repo)
		ctx = context.Background()
	})

	Describe("ExtractFactsFromEvent", func() {
		It("should extract facts from a simple event", func() {
			event := fixtures.EventWith(uuid.New().String(), "Led development team through critical project", "TechCorp", "Platform")

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())
			Expect(facts[0].SourceEventID).To(Equal(event.ID))
		})

		It("should validate extracted facts", func() {
			event := fixtures.EventWith(uuid.New().String(), "Architected enterprise platform", "TechCorp", "Platform")

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			for _, fact := range facts {
				validationErr := fact.Validate()
				Expect(validationErr).NotTo(HaveOccurred())
			}
		})

		It("should infer competencies from event text", func() {
			event := fixtures.EventWith(uuid.New().String(), "Led engineering team and mentored junior developers", "TechCorp", "Platform")

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("mentoring"))
		})

		It("should return empty list for event with empty text", func() {
			event := fixtures.EventWith(uuid.New().String(), "", "", "") // Empty text

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
			event := fixtures.EventWith("", "Some event", "", "") // Empty ID

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should extract role fit from event text", func() {
			event := fixtures.EventWith(uuid.New().String(), "Managed engineering team and coordinated hiring", "TechCorp", "Platform")

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitEM))
		})

		It("should extract audience relevance from event text", func() {
			event := fixtures.EventWith(uuid.New().String(), "Led critical infrastructure initiative", "TechCorp", "Platform")

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].AudienceRelevance).To(ContainElement("peer"))
		})

		It("should extract strength signal from event text", func() {
			event := fixtures.EventWith(uuid.New().String(), "Delivered critical feature that improved performance", "TechCorp", "Platform")

			facts, err := service.ExtractFactsFromEvent(ctx, event)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].StrengthSignal).NotTo(BeEmpty())
		})
	})

	Describe("ExtractFactsFromBurst", func() {
		It("should extract facts from burst with multiple events", func() {
			event1 := fixtures.EventWith(uuid.New().String(), "Led development team", "TechCorp", "Platform")
			event1.Date = time.Now().Add(-60 * 24 * time.Hour)

			event2 := fixtures.EventWith(uuid.New().String(), "Architected microservices platform", "TechCorp", "Platform")
			event2.Date = time.Now().Add(-50 * 24 * time.Hour)

			// Store events in repository
			_ = service.CaptureEvent(ctx, event1, ManualEntry)
			_ = service.CaptureEvent(ctx, event2, ManualEntry)

			burst := fixtures.Burst(uuid.New().String(), event1.ID, event2.ID)
			burst.Name = "Platform Initiative"

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(facts)).To(BeNumerically(">=", 2))
			Expect(facts[0].SourceBurstID).To(Equal(burst.ID))
		})

		It("should validate all extracted burst facts", func() {
			event := fixtures.EventWith(uuid.New().String(), "Implemented critical feature", "TechCorp", "Platform")

			_ = service.CaptureEvent(ctx, event, ManualEntry)

			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.ID = uuid.New().String()
			burst.Name = "Feature Initiative"
			burst.EventIDs = []string{event.ID}

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
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.ID = "" // Empty ID to test validation

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should return empty list for burst with no events", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.ID = uuid.New().String()
			burst.EventIDs = []string{} // No events

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should handle missing events gracefully", func() {
			burst := fixtures.Burst(uuid.New().String(), uuid.New().String(), uuid.New().String())
			burst.Name = "Burst with Missing Events"

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			// Should return empty list, not error (no events found)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should combine competencies from multiple events", func() {
			event1 := fixtures.EventWith(uuid.New().String(), "Led development team", "TechCorp", "Platform")
			event1.Date = time.Now().Add(-60 * 24 * time.Hour)

			event2 := fixtures.EventWith(uuid.New().String(), "Mentored junior engineers", "TechCorp", "Platform")
			event2.Date = time.Now().Add(-50 * 24 * time.Hour)

			_ = service.CaptureEvent(ctx, event1, ManualEntry)
			_ = service.CaptureEvent(ctx, event2, ManualEntry)

			burst := fixtures.Burst(uuid.New().String(), event1.ID, event2.ID)
			burst.Name = "Team Growth"

			facts, err := service.ExtractFactsFromBurst(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			// Burst-level fact should combine competencies
			Expect(facts[0].CompetencyCategories).To(ContainElement("leadership"))
			Expect(facts[0].CompetencyCategories).To(ContainElement("mentoring"))
		})
	})

	Describe("ValidateFact", func() {
		It("should validate a correct fact", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "Led critical initiative", career.RoleFitPrincipal, []string{"leadership"}, []string{"peer", "hiring_manager"}, uuid.New().String())

			err := service.ValidateFact(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject nil fact", func() {
			err := service.ValidateFact(ctx, nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fact cannot be nil"))
		})

		It("should reject fact with empty text", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "", career.RoleFitSeniorIC, []string{"technical"}, []string{"peer"}, "")

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with no competency categories", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "Some achievement", career.RoleFitSeniorIC, []string{}, []string{"peer"}, "")

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with invalid role fit", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "Some achievement", "invalid_role", []string{"technical"}, []string{"peer"}, "")

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should accept fact with no source reference (manual entry)", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "Some achievement", career.RoleFitSeniorIC, []string{"technical"}, []string{"peer"}, "")

			err := service.ValidateFact(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject fact with aspirational language", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "Will lead critical initiatives", career.RoleFitPrincipal, []string{"leadership"}, []string{"peer"}, uuid.New().String())

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should reject fact with no audience relevance", func() {
			fact := fixtures.FactForValidation(uuid.New().String(), "Some achievement", career.RoleFitSeniorIC, []string{"technical"}, []string{}, uuid.New().String())

			err := service.ValidateFact(ctx, fact)

			Expect(err).To(HaveOccurred())
		})

		It("should accept fact with burst source reference", func() {
			fact := fixtures.FactFromBurst(uuid.New().String(), uuid.New().String())
			fact.Text = "Multi-phase initiative success"
			fact.CompetencyCategories = []string{"leadership"}
			fact.RoleFit = career.RoleFitPrincipal
			fact.AudienceRelevance = []string{"peer", "hiring_manager"}

			err := service.ValidateFact(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Integration: Event to Fact to Validation", func() {
		It("should successfully extract and validate facts from event", func() {
			event := fixtures.EventWith(uuid.New().String(), "Led platform architecture project and mentored team", "TechCorp", "")
			event.Date = time.Now().Add(-30 * 24 * time.Hour)

			facts, err := service.ExtractFactsFromEvent(ctx, event)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())

			for _, fact := range facts {
				err := service.ValidateFact(ctx, &fact)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should successfully extract and validate facts from burst", func() {
			event1 := fixtures.EventWith(uuid.New().String(), "Designed system architecture", "TechCorp", "")
			event1.Date = time.Now().Add(-60 * 24 * time.Hour)

			event2 := fixtures.EventWith(uuid.New().String(), "Led implementation and delivered solution", "TechCorp", "")
			event2.Date = time.Now().Add(-50 * 24 * time.Hour)

			_ = service.CaptureEvent(ctx, event1, ManualEntry)
			_ = service.CaptureEvent(ctx, event2, ManualEntry)

			burst := fixtures.Burst(uuid.New().String(), event1.ID, event2.ID)
			burst.Name = "System Architecture Initiative"

			facts, err := service.ExtractFactsFromBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).NotTo(BeEmpty())

			for _, fact := range facts {
				err := service.ValidateFact(ctx, &fact)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Describe("SetFactRepository", func() {
		It("should set the fact repository", func() {
			factRepo := careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)

			fact := fixtures.FactForSave("Test fact for repository", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, uuid.New().String())

			err := service.SaveFact(ctx, fact)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := factRepo.GetByID(ctx, fact.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(fact.ID))
		})
	})

	Describe("GetFactsBySourceEventID", func() {
		var factRepo *careermemory.FactRepository
		var event *career.Event

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)

			event = fixtures.EventWith(uuid.New().String(), "Led migration to microservices", "TechCorp", "")
			event.Date = time.Now().Add(-30 * 24 * time.Hour)

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
		var factRepo *careermemory.FactRepository
		var burst *career.Burst

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)

			event1 := fixtures.EventWith(uuid.New().String(), "Led platform migration", "", "")
			event1.Date = time.Now().Add(-60 * 24 * time.Hour)

			event2 := fixtures.EventWith(uuid.New().String(), "Architected microservices", "", "")
			event2.Date = time.Now().Add(-60 * 24 * time.Hour)

			err := service.CaptureEvent(ctx, event1, ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			err = service.CaptureEvent(ctx, event2, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			burst = fixtures.Burst(uuid.New().String(), event1.ID, event2.ID)
			burst.Name = "Platform Migration"

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
		var factRepo *careermemory.FactRepository

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)
		})

		It("should save a valid fact", func() {
			fact := fixtures.FactForSave("Led migration to microservices", []string{"technical", "leadership"}, career.RoleFitStaff, []string{"hiring_manager", "peer"}, uuid.New().String())

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
			fact := fixtures.FactForSave("", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, uuid.New().String())

			err := service.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
		})

		It("should return error for fact with aspirational language", func() {
			fact := fixtures.FactForSave("I will lead migration", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, uuid.New().String())

			err := service.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aspirational language"))
		})

		It("should update existing fact if ID is set", func() {
			fact := fixtures.FactForSave("Original text", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, uuid.New().String())

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
			fact := fixtures.FactForSave("Test fact", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, uuid.New().String())

			err := serviceWithoutRepo.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrFactRepositoryNotConfigured)).To(BeTrue())
		})
	})

	Describe("DeleteFact", func() {
		var factRepo *careermemory.FactRepository
		var fact *career.Fact

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)

			fact = fixtures.FactForSave("Fact to delete", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, uuid.New().String())
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
			Expect(errors.Is(err, ErrFactRepositoryNotConfigured)).To(BeTrue())
		})
	})

	Describe("GetFactRepository", func() {
		It("should return nil when no fact repository is set", func() {
			factRepo := service.GetFactRepository()
			Expect(factRepo).To(BeNil())
		})

		It("should return the configured fact repository", func() {
			memoryFactRepo := careermemory.NewFactRepository()
			service.SetFactRepository(memoryFactRepo)

			factRepo := service.GetFactRepository()
			Expect(factRepo).To(Equal(memoryFactRepo))
		})
	})

	Describe("GetFactsBySourceEventID - nil repository", func() {
		It("should return empty slice when repository is nil", func() {
			serviceWithoutFact := NewService(repo)
			facts, err := serviceWithoutFact.GetFactsBySourceEventID(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})

	Describe("GetFactsBySourceBurstID - nil repository", func() {
		It("should return empty slice when repository is nil", func() {
			serviceWithoutFact := NewService(repo)
			facts, err := serviceWithoutFact.GetFactsBySourceBurstID(ctx, "burst-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})

	Describe("ExtractFactsFromBurst - detector error path", func() {
		It("returns empty facts when all burst events are missing", func() {
			burst := fixtures.Burst(uuid.New().String(), "missing-1", "missing-2")
			burst.Name = "Ghost burst"

			facts, err := service.ExtractFactsFromBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})
})

var _ = Describe("Career Service - Fact Methods (Mock-Based)", func() {
	var (
		ctrl          *gomock.Controller
		mockEventRepo *mockrepo.MockEventRepository
		mockFactRepo  *mockrepo.MockFactRepository
		svc           *Service
		ctx           context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockEventRepo = mockrepo.NewMockEventRepository(ctrl)
		mockFactRepo = mockrepo.NewMockFactRepository(ctrl)
		svc = NewService(mockEventRepo)
		svc.SetFactRepository(mockFactRepo)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("SaveFact - update failure path", func() {
		It("returns error when fact update fails", func() {
			factID := uuid.New().String()
			existingFact := fixtures.Fact(factID, "evt-1")
			fact := fixtures.FactForSave("Updated achievement text", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, "evt-1")
			fact.ID = factID

			mockFactRepo.EXPECT().
				GetByID(gomock.Any(), factID).
				Return(existingFact, nil)
			mockFactRepo.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(errors.New("update failed"))

			err := svc.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to update fact"))
		})

		It("returns error when fact create fails", func() {
			fact := fixtures.FactForSave("New fact for create error", []string{"technical"}, career.RoleFitStaff, []string{"peer"}, "evt-1")

			mockFactRepo.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("not found"))
			mockFactRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(errors.New("create failed"))

			err := svc.SaveFact(ctx, fact)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to create fact"))
		})
	})

	Describe("GetFactsBySourceEventID - repo error path", func() {
		It("returns error when repository query fails", func() {
			mockFactRepo.EXPECT().
				GetBySourceEventID(gomock.Any(), "evt-fail").
				Return(nil, errors.New("database error"))

			facts, err := svc.GetFactsBySourceEventID(ctx, "evt-fail")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to retrieve facts"))
			Expect(facts).To(BeNil())
		})
	})

	Describe("GetFactsBySourceBurstID - repo error path", func() {
		It("returns error when repository query fails", func() {
			mockFactRepo.EXPECT().
				GetBySourceBurstID(gomock.Any(), "burst-fail").
				Return(nil, errors.New("database error"))

			facts, err := svc.GetFactsBySourceBurstID(ctx, "burst-fail")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to retrieve facts"))
			Expect(facts).To(BeNil())
		})
	})

	Describe("DeleteFact - repo error path", func() {
		It("returns error when repository delete fails", func() {
			mockFactRepo.EXPECT().
				Delete(gomock.Any(), "fact-fail").
				Return(errors.New("delete failed"))

			err := svc.DeleteFact(ctx, "fact-fail")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to delete fact"))
		})
	})

	Describe("ExtractFactsFromAllEvents - progress callback", func() {
		It("calls progress callback for each event", func() {
			events := []*career.Event{
				fixtures.EventWith("id-p1", "Implemented authentication system", "TechCorp", ""),
				fixtures.EventWith("id-p2", "Led API redesign project", "TechCorp", ""),
				fixtures.EventWith("id-p3", "Mentored junior developers", "TechCorp", ""),
			}

			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(events, nil)
			mockFactRepo.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("not found")).
				AnyTimes()
			mockFactRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			var progressCalls []int
			progress := func(current, total int) {
				progressCalls = append(progressCalls, current)
			}

			count, _, err := svc.ExtractFactsFromAllEvents(ctx, progress)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">", 0))
			Expect(progressCalls).To(Equal([]int{1, 2, 3}))
		})

		It("continues when individual fact save fails", func() {
			events := []*career.Event{
				fixtures.EventWith("id-sf1", "Implemented critical system", "TechCorp", ""),
			}

			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(events, nil)
			mockFactRepo.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("not found")).
				AnyTimes()
			mockFactRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(errors.New("save failed")).
				AnyTimes()

			count, _, err := svc.ExtractFactsFromAllEvents(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
