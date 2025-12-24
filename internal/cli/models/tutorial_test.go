package models_test

import (
	"github.com/baphled/kariya/internal/cli/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TutorialModel", func() {
	var model *models.TutorialModel

	BeforeEach(func() {
		model = models.NewTutorialModel()
	})

	Context("when creating a new TutorialModel", func() {
		It("should initialize", func() {
			Expect(model).NotTo(BeNil())
		})

		It("should start at step 0", func() {
			Expect(model.CurrentStep()).To(Equal(0))
		})

		It("should not be completed initially", func() {
			Expect(model.IsCompleted()).To(BeFalse())
			Expect(model.IsSkipped()).To(BeFalse())
		})
	})

	Context("when rendering tutorial views", func() {
		It("should render welcome step", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Welcome to KaRiya"))
		})

		It("should provide a non-empty view", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Context("when tutorial is skipped", func() {
		It("should have IsSkipped method", func() {
			Expect(model.IsSkipped()).To(BeFalse())
		})

		It("should have IsCompleted method", func() {
			Expect(model.IsCompleted()).To(BeFalse())
		})
	})

	Context("when handling model methods", func() {
		It("should have Init method", func() {
			Expect(model.Init).NotTo(BeNil())
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})

		It("should have Update method", func() {
			Expect(model.Update).NotTo(BeNil())
		})

		It("should have View method", func() {
			Expect(model.View).NotTo(BeNil())
		})

		It("should have CurrentStep method", func() {
			Expect(model.CurrentStep).NotTo(BeNil())
			Expect(model.CurrentStep()).To(Equal(0))
		})
	})
})
