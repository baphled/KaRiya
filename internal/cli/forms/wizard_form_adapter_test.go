package forms_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WizardFormAdapter", func() {
	var (
		adapter *forms.WizardFormAdapter
		data    *forms.OnboardingFormData
	)

	BeforeEach(func() {
		data = &forms.OnboardingFormData{}
		huhForm := forms.NewOnboardingWizardForm(data, 80, 40)
		adapter = forms.NewWizardFormAdapter(huhForm, 3)
	})

	Describe("NewWizardFormAdapter", func() {
		It("should create an adapter wrapping a huh form", func() {
			Expect(adapter).NotTo(BeNil())
		})

		It("should start at step 0", func() {
			Expect(adapter.CurrentStep()).To(Equal(0))
		})

		It("should report the correct total steps", func() {
			Expect(adapter.TotalSteps()).To(Equal(3))
		})

		It("should not be completed initially", func() {
			Expect(adapter.IsCompleted()).To(BeFalse())
		})

		It("should not be aborted initially", func() {
			Expect(adapter.IsAborted()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("should return a command from the form", func() {
			cmd := adapter.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		It("should process key messages without panicking", func() {
			adapter.Init()

			cmd := adapter.Update(tea.KeyMsg{Type: tea.KeyEnter})

			_ = cmd
		})

		It("should process window size messages", func() {
			adapter.Init()

			cmd := adapter.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

			_ = cmd
		})
	})

	Describe("View", func() {
		It("should return non-empty form view", func() {
			adapter.Init()

			view := adapter.View()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetDimensions", func() {
		It("should update the form dimensions", func() {
			adapter.SetDimensions(100, 50)

			view := adapter.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("IsCompleted", func() {
		It("should return false when form is in progress", func() {
			adapter.Init()

			Expect(adapter.IsCompleted()).To(BeFalse())
		})
	})

	Describe("IsAborted", func() {
		It("should return false when form is in progress", func() {
			adapter.Init()

			Expect(adapter.IsAborted()).To(BeFalse())
		})
	})

	Describe("Form accessor", func() {
		It("should return the underlying form", func() {
			f := adapter.Form()

			Expect(f).NotTo(BeNil())
		})
	})

	Describe("SetForm", func() {
		It("should replace the underlying form", func() {
			newData := &forms.OnboardingFormData{Name: "New"}
			newHuhForm := forms.NewOnboardingWizardForm(newData, 60, 30)

			adapter.SetForm(newHuhForm, 3)

			Expect(adapter.Form()).NotTo(BeNil())
			Expect(adapter.TotalSteps()).To(Equal(3))
			Expect(adapter.CurrentStep()).To(Equal(0))
		})
	})

	Describe("SetCurrentStep", func() {
		It("should update the current step", func() {
			adapter.SetCurrentStep(2)

			Expect(adapter.CurrentStep()).To(Equal(2))
		})

		It("should clamp step to valid range", func() {
			adapter.SetCurrentStep(10)

			Expect(adapter.CurrentStep()).To(Equal(2))
		})

		It("should clamp negative step to 0", func() {
			adapter.SetCurrentStep(-1)

			Expect(adapter.CurrentStep()).To(Equal(0))
		})
	})

	Describe("interface compliance", func() {
		It("should satisfy the WizardForm interface methods", func() {
			Expect(adapter.Init).NotTo(BeNil())
			Expect(adapter.Update).NotTo(BeNil())
			Expect(adapter.View).NotTo(BeNil())
			Expect(adapter.IsCompleted).NotTo(BeNil())
			Expect(adapter.IsAborted).NotTo(BeNil())
			Expect(adapter.CurrentStep).NotTo(BeNil())
			Expect(adapter.TotalSteps).NotTo(BeNil())
			Expect(adapter.SetDimensions).NotTo(BeNil())
		})
	})
})
