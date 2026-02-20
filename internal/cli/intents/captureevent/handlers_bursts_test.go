package captureevent

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("createBurstFromSuggestion", func() {
	var model *BurstSuggestionModelNew

	BeforeEach(func() {
		suggestions := []burstfact.BurstSuggestion{
			{Name: "API Development", Description: "Built REST APIs", EventIDs: []string{"evt-1", "evt-2"}, ConfidenceScore: 0.9},
		}
		model = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
	})

	It("assigns a non-empty UUID to the burst ID", func() {
		suggestion := burstfact.BurstSuggestion{
			Name:        "Backend Services",
			Description: "Built scalable services",
			EventIDs:    []string{"evt-1", "evt-2"},
		}
		burst := model.createBurstFromSuggestion(suggestion)
		Expect(burst.ID).NotTo(BeEmpty())
	})

	It("sets CreatedAt to a non-zero time", func() {
		suggestion := burstfact.BurstSuggestion{
			Name:     "Backend Services",
			EventIDs: []string{"evt-1", "evt-2"},
		}
		burst := model.createBurstFromSuggestion(suggestion)
		Expect(burst.CreatedAt.IsZero()).To(BeFalse())
	})

	It("sets UpdatedAt to a non-zero time", func() {
		suggestion := burstfact.BurstSuggestion{
			Name:     "Backend Services",
			EventIDs: []string{"evt-1", "evt-2"},
		}
		burst := model.createBurstFromSuggestion(suggestion)
		Expect(burst.UpdatedAt.IsZero()).To(BeFalse())
	})

	It("produces a burst that passes Validate()", func() {
		suggestion := burstfact.BurstSuggestion{
			Name:        "Backend Services",
			Description: "Built scalable services",
			EventIDs:    []string{"evt-1", "evt-2"},
		}
		burst := model.createBurstFromSuggestion(suggestion)
		Expect(burst.Validate()).To(Succeed())
	})
})

var _ = Describe("burst persistence in postSaveReview", func() {
	var (
		intent    *Intent
		svc       *careerservice.Service
		burstRepo *memoryrepo.BurstRepository
	)

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		repos := memoryrepo.NewRepositories()
		burstRepo = repos.Burst.(*memoryrepo.BurstRepository)

		svc = careerservice.NewService(repos.Event)
		svc.SetBurstRepository(burstRepo)
		svc.SetSkillRepository(repos.Skill)
		svc.SetFactRepository(repos.Fact)

		intent.context.CareerService = svc
	})

	Context("with accepted bursts that exist in the repository", func() {
		It("calls ConfirmBurst to mark them as confirmed", func() {
			event := fixtures.EventWith("evt-saved-1", "Reviewed event with bursts", "", "")
			burst := fixtures.Burst("burst-to-confirm", "evt-saved-1", "evt-other-1")

			ctx := context.Background()
			Expect(burstRepo.Create(ctx, burst)).To(Succeed())
			Expect(burst.Confirmed).To(BeFalse())

			intent.currentState = StateReview
			intent.postSaveReview = true
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			reviewData := map[string]interface{}{
				"event":  event,
				"bursts": []*career.Burst{burst},
				"facts":  []*career.Fact{},
			}
			cmd := intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))

			confirmed, err := burstRepo.GetByID(ctx, "burst-to-confirm")
			Expect(err).NotTo(HaveOccurred())
			Expect(confirmed.Confirmed).To(BeTrue())
			Expect(confirmed.ConfirmedAt).NotTo(BeNil())
		})

		It("preserves EventIDs through ConfirmBurst", func() {
			event := fixtures.EventWith("evt-saved-2", "Reviewed event", "", "")
			burst := fixtures.Burst("burst-preserve-ids", "evt-saved-2", "evt-other-2")

			ctx := context.Background()
			Expect(burstRepo.Create(ctx, burst)).To(Succeed())

			intent.currentState = StateReview
			intent.postSaveReview = true
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			reviewData := map[string]interface{}{
				"event":  event,
				"bursts": []*career.Burst{burst},
				"facts":  []*career.Fact{},
			}
			cmd := intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))

			confirmed, err := burstRepo.GetByID(ctx, "burst-preserve-ids")
			Expect(err).NotTo(HaveOccurred())
			Expect(confirmed.EventIDs).To(Equal([]string{"evt-saved-2", "evt-other-2"}))
		})
	})

	Context("when CareerService is nil", func() {
		It("completes without confirming bursts", func() {
			event := fixtures.EventWith("evt-saved-3", "Reviewed event", "", "")
			burst := fixtures.Burst("burst-no-svc", "evt-saved-3", "evt-other-3")

			intent.context.CareerService = nil
			intent.currentState = StateReview
			intent.postSaveReview = true
			intent.reviewState = &ReviewInferredEventState{
				Event: event,
			}

			reviewData := map[string]interface{}{
				"event":  event,
				"bursts": []*career.Burst{burst},
				"facts":  []*career.Fact{},
			}
			cmd := intent.HandleSubmit(&screens.SubmitResult{FormData: reviewData})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Completed))
			Expect(burst.Confirmed).To(BeFalse())
		})
	})
})
