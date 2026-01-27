package manage_skills_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/manage_skills"
)

var _ = Describe("Types", func() {
	Describe("Intent struct", func() {
		var (
			ctx      context.Context
			mockRepo *MockSkillRepository
		)

		BeforeEach(func() {
			ctx = context.Background()
			mockRepo = NewMockSkillRepository()
		})

		Describe("Fields", func() {
			It("should have context field for input parameters", func() {
				intentCtx := manage_skills.NewIntentContext(ctx, mockRepo)
				intent := &manage_skills.Intent{}
				intent.SetContext(intentCtx)
				Expect(intent.GetContext()).To(Equal(intentCtx))
			})

			It("should have state field for state machine", func() {
				intent := &manage_skills.Intent{}
				intent.SetState(manage_skills.StateList)
				Expect(intent.GetState()).To(Equal(manage_skills.StateList))
			})

			It("should have active field to track if intent is active", func() {
				intent := &manage_skills.Intent{}
				intent.SetActive(true)
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Describe("State Transitions", func() {
			It("should allow transition from List to Detail", func() {
				intent := &manage_skills.Intent{}
				intent.SetState(manage_skills.StateList)
				intent.SetState(manage_skills.StateDetail)
				Expect(intent.GetState()).To(Equal(manage_skills.StateDetail))
			})

			It("should allow transition from List to Add", func() {
				intent := &manage_skills.Intent{}
				intent.SetState(manage_skills.StateList)
				intent.SetState(manage_skills.StateAdd)
				Expect(intent.GetState()).To(Equal(manage_skills.StateAdd))
			})

			It("should allow transition from Detail to Edit", func() {
				intent := &manage_skills.Intent{}
				intent.SetState(manage_skills.StateDetail)
				intent.SetState(manage_skills.StateEdit)
				Expect(intent.GetState()).To(Equal(manage_skills.StateEdit))
			})

			It("should allow transition from Detail to Delete", func() {
				intent := &manage_skills.Intent{}
				intent.SetState(manage_skills.StateDetail)
				intent.SetState(manage_skills.StateDelete)
				Expect(intent.GetState()).To(Equal(manage_skills.StateDelete))
			})
		})

		Describe("Active State", func() {
			It("should be inactive by default", func() {
				intent := &manage_skills.Intent{}
				Expect(intent.IsActive()).To(BeFalse())
			})

			It("should be activatable", func() {
				intent := &manage_skills.Intent{}
				intent.SetActive(true)
				Expect(intent.IsActive()).To(BeTrue())
			})

			It("should be deactivatable", func() {
				intent := &manage_skills.Intent{}
				intent.SetActive(true)
				intent.SetActive(false)
				Expect(intent.IsActive()).To(BeFalse())
			})
		})
	})
})
