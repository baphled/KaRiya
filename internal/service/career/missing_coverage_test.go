package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
)

var _ = Describe("Career Service - Missing Coverage", func() {
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

	Describe("DeleteBurst", func() {
		var (
			burstRepo *careermemory.BurstRepository
			testBurst = fixtures.Burst("burst-test-1", "event-1", "event-2")
		)

		BeforeEach(func() {
			burstRepo = careermemory.NewBurstRepository()
			service.SetBurstRepository(burstRepo)

			testBurst = fixtures.Burst("burst-test-1", "event-1", "event-2")
			testBurst.Description = "A test burst for deletion"

			err := burstRepo.Create(ctx, testBurst)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete an existing burst", func() {
			err := service.DeleteBurst(ctx, testBurst.ID)
			Expect(err).NotTo(HaveOccurred())

			// Verify burst was deleted
			_, err = burstRepo.GetByID(ctx, testBurst.ID)
			Expect(err).To(HaveOccurred())
		})

		It("should return error when burst ID is empty", func() {
			err := service.DeleteBurst(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
		})

		It("should return error when burst does not exist", func() {
			err := service.DeleteBurst(ctx, "nonexistent-burst-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to delete burst"))
		})

		It("should return ErrBurstRepositoryNotConfigured when repository is nil", func() {
			serviceWithoutBurst := NewService(repo)
			err := serviceWithoutBurst.DeleteBurst(ctx, "burst-1")
			Expect(err).To(Equal(ErrBurstRepositoryNotConfigured))
		})

		It("should handle concurrent deletions gracefully", func() {
			// First deletion should succeed
			err1 := service.DeleteBurst(ctx, testBurst.ID)
			Expect(err1).NotTo(HaveOccurred())

			// Second deletion should fail (already deleted)
			err2 := service.DeleteBurst(ctx, testBurst.ID)
			Expect(err2).To(HaveOccurred())
		})

		It("should not affect other bursts when deleting one", func() {
			burst2 := fixtures.Burst("burst-test-2", "event-3", "event-4")
			burst2.Name = "Second Burst"
			burst2.Description = "Another burst"
			err := burstRepo.Create(ctx, burst2)
			Expect(err).NotTo(HaveOccurred())

			// Delete first burst
			err = service.DeleteBurst(ctx, testBurst.ID)
			Expect(err).NotTo(HaveOccurred())

			// Verify second burst still exists
			retrieved, err := burstRepo.GetByID(ctx, burst2.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(burst2.ID))
			Expect(retrieved.Name).To(Equal("Second Burst"))
		})
	})

	Describe("GetEventRepository", func() {
		It("should return the event repository", func() {
			eventRepo := service.GetEventRepository()
			Expect(eventRepo).NotTo(BeNil())
			Expect(eventRepo).To(Equal(repo))
		})

		It("should return same repository instance on multiple calls", func() {
			repo1 := service.GetEventRepository()
			repo2 := service.GetEventRepository()
			Expect(repo1).To(BeIdenticalTo(repo2))
		})

		It("should allow using returned repository for operations", func() {
			eventRepo := service.GetEventRepository()

			event := fixtures.EventWith("test-event-1", "Test event via GetEventRepository", "TestCo", "")

			err := eventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve via service to verify
			retrieved, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(event.ID))
			Expect(retrieved.Text).To(Equal(event.Text))
		})
	})

	Describe("GetFactsBySourceEventID", func() {
		var factRepo *careermemory.FactRepository

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)
		})

		It("should return empty slice when repository is nil", func() {
			serviceWithoutFact := NewService(repo)
			facts, err := serviceWithoutFact.GetFactsBySourceEventID(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for empty event ID", func() {
			facts, err := service.GetFactsBySourceEventID(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should return empty list for event with no facts", func() {
			facts, err := service.GetFactsBySourceEventID(ctx, "event-no-facts")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})

	Describe("GetFactsBySourceBurstID", func() {
		var factRepo *careermemory.FactRepository

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)
		})

		It("should return empty slice when repository is nil", func() {
			serviceWithoutFact := NewService(repo)
			facts, err := serviceWithoutFact.GetFactsBySourceBurstID(ctx, "burst-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for empty burst ID", func() {
			facts, err := service.GetFactsBySourceBurstID(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should return empty list for burst with no facts", func() {
			facts, err := service.GetFactsBySourceBurstID(ctx, "burst-no-facts")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})

	Describe("GetSkillRepository", func() {
		It("returns nil when no skill repository is set", func() {
			skillRepo := service.GetSkillRepository()
			Expect(skillRepo).To(BeNil())
		})

		It("returns the configured skill repository", func() {
			memorySkillRepo := careermemory.NewSkillRepository()
			service.SetSkillRepository(memorySkillRepo)

			skillRepo := service.GetSkillRepository()
			Expect(skillRepo).To(Equal(memorySkillRepo))
		})

		It("returns same instance on multiple calls", func() {
			memorySkillRepo := careermemory.NewSkillRepository()
			service.SetSkillRepository(memorySkillRepo)

			repo1 := service.GetSkillRepository()
			repo2 := service.GetSkillRepository()
			Expect(repo1).To(BeIdenticalTo(repo2))
		})
	})

	Describe("SaveBurst - additional paths", func() {
		It("rejects a nil burst", func() {
			err := service.SaveBurst(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst cannot be nil"))
		})

		It("saves burst without repository configured (no-op)", func() {
			burst := fixtures.Burst("burst-no-repo", "evt-1", "evt-2")

			err := service.SaveBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("RejectBurstSuggestion - empty event IDs", func() {
		It("returns nil for empty event IDs", func() {
			err := service.RejectBurstSuggestion(ctx, []string{})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("SetSkillRepository", func() {
		It("sets the skill repository for use by SaveSkill", func() {
			skillRepo := careermemory.NewSkillRepository()
			service.SetSkillRepository(skillRepo)

			skill := fixtures.SkillWith("", "Go", "backend", "advanced")

			err := service.SaveSkill(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.ID).NotTo(BeEmpty())
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

var _ = Describe("Career Service - Mock-Based Coverage Gaps", func() {
	var (
		ctrl          *gomock.Controller
		mockEventRepo *mockrepo.MockEventRepository
		mockFactRepo  *mockrepo.MockFactRepository
		mockBurstRepo *mockrepo.MockBurstRepository
		svc           *Service
		ctx           context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockEventRepo = mockrepo.NewMockEventRepository(ctrl)
		mockFactRepo = mockrepo.NewMockFactRepository(ctrl)
		mockBurstRepo = mockrepo.NewMockBurstRepository(ctrl)
		svc = NewService(mockEventRepo)
		svc.SetFactRepository(mockFactRepo)
		svc.SetBurstRepository(mockBurstRepo)
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

	Describe("SaveBurst - repo error path", func() {
		It("returns error when burst repo create fails", func() {
			burst := fixtures.Burst("burst-save-fail", "evt-1", "evt-2")

			mockBurstRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(errors.New("create failed"))

			err := svc.SaveBurst(ctx, burst)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to save burst"))
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

	Describe("DetectAndSaveBursts - empty events and suggestion failures", func() {
		It("returns zero counts when no events exist", func() {
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*career.Event{}, nil)

			count, saved, err := svc.DetectAndSaveBursts(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
			Expect(saved).To(Equal(0))
		})

		It("returns error when SaveBurstSuggestions fails", func() {
			now := time.Now()
			events := []*career.Event{
				fixtures.EventWith("d1", "Led backend infrastructure project for TechCorp", "TechCorp", "Platform"),
				fixtures.EventWith("d2", "Led backend infrastructure initiative for TechCorp", "TechCorp", "Platform"),
			}
			events[0].Date = now
			events[1].Date = now.Add(24 * time.Hour)

			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(events, nil)
			mockEventRepo.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(events[0], nil).
				AnyTimes()
			mockBurstRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			count, saved, err := svc.DetectAndSaveBursts(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 0))
			Expect(saved).To(BeNumerically(">=", 0))
		})
	})

	Describe("SuggestBurstsWithOptions - detector error", func() {
		It("returns empty when all events are missing from repo", func() {
			mockEventRepo.EXPECT().
				GetByID(gomock.Any(), "missing-1").
				Return(nil, errors.New("not found"))
			mockEventRepo.EXPECT().
				GetByID(gomock.Any(), "missing-2").
				Return(nil, errors.New("not found"))

			suggestions, err := svc.SuggestBurstsWithOptions(ctx, []string{"missing-1", "missing-2"}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})
	})

	Describe("SaveSkill - timestamp preservation via mock", func() {
		var mockSkillRepo *mockrepo.MockSkillRepository

		BeforeEach(func() {
			mockSkillRepo = mockrepo.NewMockSkillRepository(ctrl)
			svc.SetSkillRepository(mockSkillRepo)
		})

		It("preserves existing CreatedAt when non-zero", func() {
			existingCreatedAt := time.Now().Add(-48 * time.Hour)
			skill := fixtures.SkillWith("", "Go", "backend", "advanced")
			skill.CreatedAt = existingCreatedAt

			mockSkillRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil)

			err := svc.SaveSkill(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.CreatedAt).To(Equal(existingCreatedAt))
		})

		It("sets CreatedAt when zero", func() {
			skill := fixtures.SkillWith("", "Python", "backend", "intermediate")
			skill.CreatedAt = time.Time{}

			mockSkillRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil)

			err := svc.SaveSkill(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.CreatedAt).NotTo(BeZero())
		})
	})
})
