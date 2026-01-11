package intents

import (
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Edit Modals Overlay Interface", func() {
	Describe("EditMetadataModal", func() {
		var modal *EditMetadataModal

		BeforeEach(func() {
			modal = NewEditMetadataModal("Test event text", "2024-01-15", "Company", "Project", []string{"tag1"}, []string{"cat1"})
		})

		Describe("GetTitle", func() {
			It("returns the correct title", func() {
				Expect(modal.GetTitle()).To(Equal("Edit Event"))
			})
		})

		Describe("GetContent", func() {
			It("returns non-empty content when modal is active", func() {
				Expect(modal.GetContent()).NotTo(BeEmpty())
			})

			It("returns empty content when modal is complete", func() {
				modal.result = &ModalEditResult[*MetadataSnapshot]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}
				Expect(modal.GetContent()).To(BeEmpty())
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
	})

	Describe("EditBurstModal", func() {
		var (
			modal *EditBurstModal
			burst *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test Description",
			}
			modal = NewEditBurstModal(burst)
		})

		Describe("GetTitle", func() {
			It("returns the correct title", func() {
				Expect(modal.GetTitle()).To(Equal("Edit Burst"))
			})
		})

		Describe("GetContent", func() {
			It("returns non-empty content when modal is active", func() {
				Expect(modal.GetContent()).NotTo(BeEmpty())
			})

			It("returns empty content when modal is complete", func() {
				modal.result = &ModalEditResult[*career.Burst]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}
				Expect(modal.GetContent()).To(BeEmpty())
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
	})

	Describe("EditFactModal", func() {
		var (
			modal *EditFactModal
			fact  *career.Fact
		)

		BeforeEach(func() {
			fact = &career.Fact{
				ID:   "fact-1",
				Text: "Test fact text",
			}
			modal = NewEditFactModal(fact)
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

			It("returns empty content when modal is complete", func() {
				modal.result = &ModalEditResult[*career.Fact]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}
				Expect(modal.GetContent()).To(BeEmpty())
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
	})
})
