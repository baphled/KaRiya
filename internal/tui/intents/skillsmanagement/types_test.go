package skillsmanagement_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/tui/intents/skillsmanagement"
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
				intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
				intent := &skillsmanagement.Intent{}
				intent.SetContext(intentCtx)
				Expect(intent.GetContext()).To(Equal(intentCtx))
			})

			It("should have state field for state machine", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetState(skillsmanagement.StateList)
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			})

			It("should have active field to track if intent is active", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetActive(true)
				Expect(intent.IsActive()).To(BeTrue())
			})

			It("should report no active modal by default", func() {
				intent := &skillsmanagement.Intent{}
				Expect(intent.HasActiveModal()).To(BeFalse())
			})

			It("should return zero selected index by default", func() {
				intent := &skillsmanagement.Intent{}
				Expect(intent.GetSelectedIndex()).To(Equal(0))
			})
		})

		Describe("State Transitions", func() {
			It("should allow transition from List to Detail", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetState(skillsmanagement.StateList)
				intent.SetState(skillsmanagement.StateDetail)
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateDetail))
			})

			It("should allow transition from List to Add", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetState(skillsmanagement.StateList)
				intent.SetState(skillsmanagement.StateAdd)
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateAdd))
			})

			It("should allow transition from Detail to Edit", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetState(skillsmanagement.StateDetail)
				intent.SetState(skillsmanagement.StateEdit)
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateEdit))
			})

			It("should allow transition from Detail to Delete", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetState(skillsmanagement.StateDetail)
				intent.SetState(skillsmanagement.StateDelete)
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateDelete))
			})
		})

		Describe("Active State", func() {
			It("should be inactive by default", func() {
				intent := &skillsmanagement.Intent{}
				Expect(intent.IsActive()).To(BeFalse())
			})

			It("should be activatable", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetActive(true)
				Expect(intent.IsActive()).To(BeTrue())
			})

			It("should be deactivatable", func() {
				intent := &skillsmanagement.Intent{}
				intent.SetActive(true)
				intent.SetActive(false)
				Expect(intent.IsActive()).To(BeFalse())
			})
		})
	})
})
