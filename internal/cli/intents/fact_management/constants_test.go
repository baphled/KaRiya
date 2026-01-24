package fact_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/fact_management"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		Describe("State Values", func() {
			It("should define StateList", func() {
				Expect(fact_management.StateList).To(Equal(fact_management.State("list")))
			})

			It("should define StateView", func() {
				Expect(fact_management.StateView).To(Equal(fact_management.State("view")))
			})

			It("should define StateEditor", func() {
				Expect(fact_management.StateEditor).To(Equal(fact_management.State("editor")))
			})

			It("should define StateDeleteConfirm", func() {
				Expect(fact_management.StateDeleteConfirm).To(Equal(fact_management.State("delete_confirm")))
			})

			It("should define StateResults", func() {
				Expect(fact_management.StateResults).To(Equal(fact_management.State("results")))
			})

			It("should define StateCompleted", func() {
				Expect(fact_management.StateCompleted).To(Equal(fact_management.State("completed")))
			})
		})

		Describe("Type Safety", func() {
			It("should be a string type", func() {
				var state fact_management.State = "test"
				Expect(string(state)).To(Equal("test"))
			})

			It("should allow comparison", func() {
				state := fact_management.StateList
				Expect(state).To(Equal(fact_management.StateList))
				Expect(state).NotTo(Equal(fact_management.StateView))
			})
		})
	})
})
