package models_test

import (
	"github.com/baphled/kariya/internal/cli/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HelpModel", func() {
	var model *models.HelpModel

	BeforeEach(func() {
		model = models.NewHelpModel()
	})

	Context("when creating a new HelpModel", func() {
		It("should initialize successfully", func() {
			Expect(model).NotTo(BeNil())
		})

		It("should start at first help section", func() {
			Expect(model.CurrentSection()).To(Equal(0))
		})

		It("should not be closed initially", func() {
			Expect(model.IsClosed()).To(BeFalse())
		})
	})

	Context("when displaying help content", func() {
		It("should render help view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Overview"))
		})

		It("should provide a non-empty view", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain help text", func() {
			view := model.View()
			Expect(len(view) > 50).To(BeTrue())
		})
	})

	Context("when navigating help sections", func() {
		It("should have multiple sections", func() {
			sectionCount := model.SectionCount()
			Expect(sectionCount > 0).To(BeTrue())
		})

		It("should have a current section getter", func() {
			section := model.CurrentSection()
			Expect(section >= 0).To(BeTrue())
		})
	})

	Context("when handling help queries", func() {
		It("should search help content", func() {
			results := model.SearchHelp("event")
			Expect(results).NotTo(BeNil())
		})

		It("should return results for matching queries", func() {
			results := model.SearchHelp("capture")
			Expect(results).NotTo(BeEmpty())
		})

		It("should handle non-matching queries gracefully", func() {
			results := model.SearchHelp("xyzabc123notaword")
			Expect(results).NotTo(BeNil())
		})
	})

	Context("when handling model methods", func() {
		It("should have Init method", func() {
			Expect(model.Init).NotTo(BeNil())
		})

		It("should have Update method", func() {
			Expect(model.Update).NotTo(BeNil())
		})

		It("should have View method", func() {
			Expect(model.View).NotTo(BeNil())
		})

		It("should have SectionCount method", func() {
			count := model.SectionCount()
			Expect(count > 0).To(BeTrue())
		})

		It("should have CurrentSection method", func() {
			section := model.CurrentSection()
			Expect(section >= 0).To(BeTrue())
		})

		It("should have SearchHelp method", func() {
			results := model.SearchHelp("help")
			Expect(results).NotTo(BeNil())
		})

		It("should have IsClosed method", func() {
			Expect(model.IsClosed()).To(BeFalse())
		})
	})
})
