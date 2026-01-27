package skills_management_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skills_management"
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
				intentCtx := skills_management.NewIntentContext(ctx, mockRepo)
				intent := &skills_management.Intent{}
				intent.SetContext(intentCtx)
				Expect(intent.GetContext()).To(Equal(intentCtx))
			})

			It("should have state field for state machine", func() {
				intent := &skills_management.Intent{}
				intent.SetState(skills_management.StateList)
				Expect(intent.GetState()).To(Equal(skills_management.StateList))
			})

			It("should have active field to track if intent is active", func() {
				intent := &skills_management.Intent{}
				intent.SetActive(true)
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Describe("State Transitions", func() {
			It("should allow transition from List to Detail", func() {
				intent := &skills_management.Intent{}
				intent.SetState(skills_management.StateList)
				intent.SetState(skills_management.StateDetail)
				Expect(intent.GetState()).To(Equal(skills_management.StateDetail))
			})

			It("should allow transition from List to Add", func() {
				intent := &skills_management.Intent{}
				intent.SetState(skills_management.StateList)
				intent.SetState(skills_management.StateAdd)
				Expect(intent.GetState()).To(Equal(skills_management.StateAdd))
			})

			It("should allow transition from Detail to Edit", func() {
				intent := &skills_management.Intent{}
				intent.SetState(skills_management.StateDetail)
				intent.SetState(skills_management.StateEdit)
				Expect(intent.GetState()).To(Equal(skills_management.StateEdit))
			})

			It("should allow transition from Detail to Delete", func() {
				intent := &skills_management.Intent{}
				intent.SetState(skills_management.StateDetail)
				intent.SetState(skills_management.StateDelete)
				Expect(intent.GetState()).To(Equal(skills_management.StateDelete))
			})
		})

		Describe("Active State", func() {
			It("should be inactive by default", func() {
				intent := &skills_management.Intent{}
				Expect(intent.IsActive()).To(BeFalse())
			})

			It("should be activatable", func() {
				intent := &skills_management.Intent{}
				intent.SetActive(true)
				Expect(intent.IsActive()).To(BeTrue())
			})

			It("should be deactivatable", func() {
				intent := &skills_management.Intent{}
				intent.SetActive(true)
				intent.SetActive(false)
				Expect(intent.IsActive()).To(BeFalse())
			})
		})
	})
})
