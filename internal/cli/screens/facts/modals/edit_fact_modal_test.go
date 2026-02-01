package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/facts/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EditFactModal", func() {
	var (
		modal *modals.EditFactModal
		fact  *career.Fact
	)

	BeforeEach(func() {
		fact = fixtures.FactWith("fact-1", "Test fact text")
		modal = modals.NewEditFactModal(fact)
	})

	Describe("GetTitle", func() {
		It("returns the correct title", func() {
			Expect(modal.GetTitle()).To(Equal("Edit Fact"))
		})
	})

	Describe("GetContent", func() {
		It("returns non-empty content when modal is active", func() {
			Expect(modal.GetContent()).NotTo(BeEmpty())
		})
	})

	Describe("GetFooter", func() {
		It("returns non-empty footer instructions", func() {
			Expect(modal.GetFooter()).NotTo(BeEmpty())
		})

		It("contains navigation instructions", func() {
			footer := modal.GetFooter()
			Expect(footer).To(ContainSubstring("Enter"))
			Expect(footer).To(ContainSubstring("Esc"))
			Expect(footer).To(ContainSubstring("Tab"))
		})
	})

	Describe("IsComplete", func() {
		It("returns false when modal is not complete", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})
	})

	Describe("Result", func() {
		It("returns nil when modal is not complete", func() {
			Expect(modal.Result()).To(BeNil())
		})
	})

	Describe("View", func() {
		It("returns non-empty view when modal is active", func() {
			Expect(modal.View()).NotTo(BeEmpty())
		})
	})
})
