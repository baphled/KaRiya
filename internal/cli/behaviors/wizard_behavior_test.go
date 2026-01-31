package behaviors_test

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testWizardData struct {
	Name  string
	Email string
}

type mockWizardForm struct {
	initCalled       bool
	updateCalled     bool
	viewCalled       bool
	completed        bool
	aborted          bool
	currentStep      int
	totalSteps       int
	lastMsg          tea.Msg
	dimensionsWidth  int
	dimensionsHeight int
	viewContent      string
}

func newMockWizardForm(totalSteps int) *mockWizardForm {
	return &mockWizardForm{
		totalSteps:  totalSteps,
		viewContent: "mock form view",
	}
}

func (m *mockWizardForm) Init() tea.Cmd {
	m.initCalled = true
	return nil
}

func (m *mockWizardForm) Update(msg tea.Msg) tea.Cmd {
	m.updateCalled = true
	m.lastMsg = msg
	return nil
}

func (m *mockWizardForm) View() string {
	m.viewCalled = true
	return m.viewContent
}

func (m *mockWizardForm) IsCompleted() bool {
	return m.completed
}

func (m *mockWizardForm) IsAborted() bool {
	return m.aborted
}

func (m *mockWizardForm) CurrentStep() int {
	return m.currentStep
}

func (m *mockWizardForm) TotalSteps() int {
	return m.totalSteps
}

func (m *mockWizardForm) SetDimensions(width, height int) {
	m.dimensionsWidth = width
	m.dimensionsHeight = height
}

var _ = Describe("WizardBehavior", func() {
	var (
		wizard *behaviors.WizardBehavior[testWizardData]
		form   *mockWizardForm
		data   *testWizardData
	)

	BeforeEach(func() {
		form = newMockWizardForm(3)
		data = &testWizardData{Name: "Test User", Email: "test@example.com"}
		wizard = behaviors.NewWizardBehavior[testWizardData](form, data)
	})

	Describe("NewWizardBehavior", func() {
		It("should create a wizard behavior with initial state", func() {
			Expect(wizard).NotTo(BeNil())
		})

		It("should be visible by default", func() {
			Expect(wizard.IsVisible()).To(BeTrue())
		})

		It("should not be completed initially", func() {
			Expect(wizard.IsCompleted()).To(BeFalse())
		})

		It("should not be cancelled initially", func() {
			Expect(wizard.IsCancelled()).To(BeFalse())
		})

		It("should not be skipped initially", func() {
			Expect(wizard.IsSkipped()).To(BeFalse())
		})

		It("should hold the provided data", func() {
			Expect(wizard.Data()).To(Equal(data))
		})

		It("should hold the provided form", func() {
			Expect(wizard.Form()).To(Equal(form))
		})
	})

	Describe("Init", func() {
		It("should delegate to the form's Init", func() {
			wizard.Init()

			Expect(form.initCalled).To(BeTrue())
		})

		It("should return nil when not visible", func() {
			wizard.Hide()

			cmd := wizard.Init()

			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should delegate messages to the form", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}

			wizard.Update(msg)

			Expect(form.updateCalled).To(BeTrue())
			Expect(form.lastMsg).To(Equal(msg))
		})

		It("should return nil when not visible", func() {
			wizard.Hide()

			cmd := wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(cmd).To(BeNil())
			Expect(form.updateCalled).To(BeFalse())
		})

		It("should mark as completed when form completes", func() {
			form.completed = true

			wizard.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(wizard.IsCompleted()).To(BeTrue())
			Expect(wizard.IsVisible()).To(BeFalse())
		})

		It("should mark as cancelled when form aborts", func() {
			form.aborted = true

			wizard.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(wizard.IsCancelled()).To(BeTrue())
			Expect(wizard.IsVisible()).To(BeFalse())
		})
	})

	Describe("View", func() {
		It("should delegate to the form's View", func() {
			view := wizard.View()

			Expect(form.viewCalled).To(BeTrue())
			Expect(view).To(Equal("mock form view"))
		})

		It("should return empty string when not visible", func() {
			wizard.Hide()

			view := wizard.View()

			Expect(view).To(BeEmpty())
		})
	})

	Describe("Visibility", func() {
		It("should hide when Hide is called", func() {
			wizard.Hide()

			Expect(wizard.IsVisible()).To(BeFalse())
		})

		It("should show when Show is called", func() {
			wizard.Hide()
			wizard.Show()

			Expect(wizard.IsVisible()).To(BeTrue())
		})
	})

	Describe("SkipWizard", func() {
		It("should mark as skipped and completed", func() {
			skipFn := wizard.Skip
			skipFn()

			Expect(wizard.IsSkipped()).To(BeTrue())
			Expect(wizard.IsCompleted()).To(BeTrue())
			Expect(wizard.IsVisible()).To(BeFalse())
		})
	})

	Describe("Complete", func() {
		It("should mark as completed and hide", func() {
			wizard.Complete()

			Expect(wizard.IsCompleted()).To(BeTrue())
			Expect(wizard.IsVisible()).To(BeFalse())
		})
	})

	Describe("Cancel", func() {
		It("should mark as cancelled and hide", func() {
			wizard.Cancel()

			Expect(wizard.IsCancelled()).To(BeTrue())
			Expect(wizard.IsVisible()).To(BeFalse())
		})
	})

	Describe("Reset", func() {
		It("should reset state flags", func() {
			wizard.Complete()
			Expect(wizard.IsCompleted()).To(BeTrue())

			wizard.Reset()

			Expect(wizard.IsCompleted()).To(BeFalse())
			Expect(wizard.IsCancelled()).To(BeFalse())
			Expect(wizard.IsSkipped()).To(BeFalse())
			Expect(wizard.IsVisible()).To(BeTrue())
		})
	})

	Describe("Step Information", func() {
		It("should return current step from form", func() {
			form.currentStep = 2

			Expect(wizard.CurrentStep()).To(Equal(2))
		})

		It("should return total steps from form", func() {
			Expect(wizard.TotalSteps()).To(Equal(3))
		})
	})

	Describe("SetForm", func() {
		It("should replace the form", func() {
			newForm := newMockWizardForm(5)

			wizard.SetForm(newForm)

			Expect(wizard.Form()).To(Equal(newForm))
			Expect(wizard.TotalSteps()).To(Equal(5))
		})
	})
})
