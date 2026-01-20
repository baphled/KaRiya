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

var _ = Describe("ExportArtifact Intent (Screen-Based)", func() {
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

		It("should initialize with SelectType state", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())
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

		It("should initialize with a screen", func() {
			intent.Init()
			Expect(intent.activeScreen).NotTo(BeNil())
		})

		It("should return nil command", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("State Transitions via Update", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
		})

		Context("Type Selection", func() {
			It("should navigate with down key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(intent.GetState()).To(Equal(ExportStateSelectType))
			})

			It("should select type on Enter", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.GetState()).To(Equal(ExportStateSelectFormat))
				Expect(intent.GetConfig()).NotTo(BeNil())
				Expect(intent.GetConfig().ArtifactType).To(Equal(ExportTypeEvents))
			})

			It("should cancel on Escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.IsActive()).To(BeFalse())
			})
		})

		Context("Format Selection", func() {
			BeforeEach(func() {
				// Navigate to format selection
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.GetState()).To(Equal(ExportStateSelectFormat))
			})

			It("should select format on Enter", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.GetState()).To(Equal(ExportStateSelectDest))
				Expect(intent.GetConfig().Format).To(Equal(ExportFormatJSON))
			})

			It("should go back on Escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.GetState()).To(Equal(ExportStateSelectType))
			})
		})

		Context("Destination Selection", func() {
			BeforeEach(func() {
				// Navigate to destination selection
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select type
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select format
				Expect(intent.GetState()).To(Equal(ExportStateSelectDest))
			})

			It("should select destination and show preview", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.GetState()).To(Equal(ExportStatePreview))
				Expect(intent.GetConfig().Destination).To(Equal(ExportDestinationFile))
			})

			It("should go back on Escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.GetState()).To(Equal(ExportStateSelectFormat))
			})
		})

		Context("Preview", func() {
			BeforeEach(func() {
				// Navigate to preview
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select type
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select format
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select dest
				Expect(intent.GetState()).To(Equal(ExportStatePreview))
			})

			It("should proceed to confirm on Enter", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.GetState()).To(Equal(ExportStateConfirm))
			})

			It("should go back on Escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.GetState()).To(Equal(ExportStateSelectDest))
			})
		})

		Context("Confirm", func() {
			BeforeEach(func() {
				// Navigate to confirm
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select type
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select format
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select dest
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Preview
				Expect(intent.GetState()).To(Equal(ExportStateConfirm))
			})

			It("should go back on 'n' key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(intent.GetState()).To(Equal(ExportStatePreview))
			})

			It("should go back on Escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.GetState()).To(Equal(ExportStatePreview))
			})
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
		})

		It("should render type selection view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("events"))
		})

		It("should render with navigation breadcrumbs", func() {
			view := intent.View()
			// Breadcrumbs show navigation path like "Main Menu ▸ ... ▸ Select Type"
			Expect(view).To(ContainSubstring("Select Type"))
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

		It("should return cancelled result on escape from first screen", func() {
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
			// Help is handled globally by HandleGlobalKeys
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			// Help toggle is a state change, not a command
			Expect(intent.IsHelpVisible()).To(BeTrue())
		})

		It("should pass 'q' key to screen (not handled globally)", func() {
			// Note: 'q' is intentionally NOT handled globally to prevent accidental exits
			// Users should only be able to quit from main menu
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			// 'q' is delegated to the screen, which may or may not produce a command
			// The key point is it doesn't return tea.Quit directly
			_ = cmd // No assertion - just verifying it doesn't panic
		})
	})

	Describe("GetState Helper", func() {
		It("should return correct state name for each state", func() {
			intent, _ := NewExportArtifactIntent(NewTestExportArtifactContext())

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
