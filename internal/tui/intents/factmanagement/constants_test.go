package factmanagement_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		Describe("State Values", func() {
			It("should define StateList", func() {
				Expect(factmanagement.StateList).To(Equal(factmanagement.State("list")))
			})

			It("should define StateView", func() {
				Expect(factmanagement.StateView).To(Equal(factmanagement.State("view")))
			})

			It("should define StateEditor", func() {
				Expect(factmanagement.StateEditor).To(Equal(factmanagement.State("editor")))
			})

			It("should define StateDeleteConfirm", func() {
				Expect(factmanagement.StateDeleteConfirm).To(Equal(factmanagement.State("delete_confirm")))
			})

			It("should define StateResults", func() {
				Expect(factmanagement.StateResults).To(Equal(factmanagement.State("results")))
			})

			It("should define StateCompleted", func() {
				Expect(factmanagement.StateCompleted).To(Equal(factmanagement.State("completed")))
			})
		})

		Describe("Type Safety", func() {
			It("should be a string type", func() {
				var state factmanagement.State = "test"
				Expect(string(state)).To(Equal("test"))
			})

			It("should allow comparison", func() {
				state := factmanagement.StateList
				Expect(state).To(Equal(factmanagement.StateList))
				Expect(state).NotTo(Equal(factmanagement.StateView))
			})
		})
	})
})
