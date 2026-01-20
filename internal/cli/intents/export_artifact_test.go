package intents

import (
	"context"

	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NewTestExportArtifactContext creates a test context with default values
func NewTestExportArtifactContext() *ExportArtifactContext {
	return &ExportArtifactContext{
		ArtifactTypes:    DefaultArtifactTypes(),
		SupportedFormats: DefaultSupportedFormats(),
		DefaultFormat:    DefaultFormats(),
		Destinations:     DefaultDestinations(),
		ExportService:    nil,
		CareerService:    nil,
		EventRepository:  nil,
		FactRepository:   nil,
		BurstRepository:  nil,
		AppContext:       context.Background(),
	}
}

// NewTestExportArtifactContextWithServices creates a test context with real services
func NewTestExportArtifactContextWithServices() *ExportArtifactContext {
	log := logger.New(nil, logger.ErrorLevel)
	eventRepo := careerrepo.NewMemoryRepository()
	factRepo := careerrepo.NewMemoryFactRepository()
	burstRepo := careerrepo.NewMemoryBurstRepository()
	exportService := cv.NewExportService(log)

	return &ExportArtifactContext{
		ArtifactTypes:    DefaultArtifactTypes(),
		SupportedFormats: DefaultSupportedFormats(),
		DefaultFormat:    DefaultFormats(),
		Destinations:     DefaultDestinations(),
		ExportService:    exportService,
		CareerService:    nil,
		EventRepository:  eventRepo,
		FactRepository:   factRepo,
		BurstRepository:  burstRepo,
		AppContext:       context.Background(),
	}
}

var _ = Describe("ExportArtifact Intent (Wizard-Based)", func() {
	var intent *ExportArtifactIntent

	Describe("Intent Creation", func() {
		It("should create a new intent with valid context", func() {
			var err error
			intent, err = NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should fail with nil context", func() {
			_, err := NewExportArtifactIntent(nil)
			Expect(err).To(HaveOccurred())
		})

		It("should initialize with SelectType state before Init", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())
			// Before Init, state is SelectType (default)
			Expect(intent.GetState()).To(Equal(ExportStateSelectType))
		})

		It("should be active after creation", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())
			Expect(intent.IsActive()).To(BeTrue())
		})
	})

	Describe("Init Method", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
		})

		It("should initialize with wizard modal", func() {
			intent.Init()
			// After Init, wizard should be created
			Expect(intent.wizardModal).NotTo(BeNil())
		})

		It("should set state to Configure", func() {
			intent.Init()
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
		})

		It("should return a command for wizard init", func() {
			cmd := intent.Init()
			// Wizard init may return a command
			_ = cmd // No assertion - just verify no panic
		})
	})

	Describe("Wizard Interaction", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
		})

		Context("Wizard Visibility", func() {
			It("should show wizard modal initially", func() {
				Expect(intent.wizardModal).NotTo(BeNil())
				Expect(intent.wizardModal.IsVisible()).To(BeTrue())
			})
		})

		Context("Wizard Cancellation", func() {
			It("should cancel intent when wizard is cancelled at step 1", func() {
				// Simulate Escape at wizard step 1
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Intent should be cancelled
				Expect(intent.IsActive()).To(BeFalse())
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(Cancelled))
			})
		})

		Context("Ctrl+S Skip", func() {
			It("should skip wizard and go to preview on Ctrl+S", func() {
				// Simulate Ctrl+S to skip wizard with defaults
				intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

				// Wizard should be complete (hidden but preserved for back-nav), intent should move to preview
				Expect(intent.wizardModal).NotTo(BeNil())
				Expect(intent.wizardModal.IsVisible()).To(BeFalse())
				Expect(intent.GetState()).To(Equal(ExportStatePreview))
				Expect(intent.GetConfig()).NotTo(BeNil())
				// Defaults should be applied
				Expect(intent.GetConfig().ArtifactType).To(Equal(ExportTypeEvents))
				Expect(intent.GetConfig().Format).To(Equal(ExportFormatJSON))
				Expect(intent.GetConfig().Destination).To(Equal(ExportDestinationFile))
			})
		})
	})

	Describe("Preview State", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
			// Skip wizard to get to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should have a preview screen", func() {
			Expect(intent.activeScreen).NotTo(BeNil())
		})

		It("should proceed to confirm on Enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))
		})

		It("should go back to wizard on Escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			// Should go back to Configure state with wizard re-created
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
			Expect(intent.wizardModal).NotTo(BeNil())
		})
	})

	Describe("Confirm State", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
			// Skip wizard to get to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			// Go to confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))
		})

		It("should show confirm modal", func() {
			Expect(intent.confirmModal).NotTo(BeNil())
			Expect(intent.confirmModal.IsVisible()).To(BeTrue())
		})

		It("should go back to preview on 'n' key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should go back to preview on Escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
		})

		It("should render wizard configuration view", func() {
			view := intent.View()
			// Wizard shows step title
			Expect(view).To(ContainSubstring("Step 1"))
		})

		It("should render with breadcrumbs", func() {
			view := intent.View()
			// Breadcrumbs show navigation path
			Expect(view).To(ContainSubstring("Configure"))
		})

		It("should show wizard title", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Export Configuration"))
		})
	})

	Describe("Result", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
		})

		It("should return nil before completion", func() {
			Expect(intent.Result()).To(BeNil())
		})

		It("should return cancelled result on escape from wizard", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Global Keys", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
		})

		It("should toggle help on '?' key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(intent.IsHelpVisible()).To(BeTrue())
		})

		It("should pass 'q' key to wizard (not handled globally)", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			_ = cmd // No assertion - just verifying it doesn't panic
		})
	})

	Describe("GetState Helper", func() {
		It("should return correct state name for each state", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())

			intent.SetState(ExportStateConfigure)
			Expect(intent.getStateName()).To(Equal("Configure"))

			intent.SetState(ExportStateSelectType)
			Expect(intent.getStateName()).To(Equal("Select Type"))

			intent.SetState(ExportStateSelectFormat)
			Expect(intent.getStateName()).To(Equal("Select Format"))

			intent.SetState(ExportStateSelectDest)
			Expect(intent.getStateName()).To(Equal("Select Destination"))

			intent.SetState(ExportStatePreview)
			Expect(intent.getStateName()).To(Equal("Preview"))

			intent.SetState(ExportStateConfirm)
			Expect(intent.getStateName()).To(Equal("Confirm"))

			intent.SetState(ExportStateInProgress)
			Expect(intent.getStateName()).To(Equal("Exporting"))

			intent.SetState(ExportStateComplete)
			Expect(intent.getStateName()).To(Equal("Complete"))

			intent.SetState(ExportStateFailed)
			Expect(intent.getStateName()).To(Equal("Failed"))
		})
	})
})
