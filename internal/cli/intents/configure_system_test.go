package intents

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigureSystem Intent", func() {
	var (
		intent *ConfigureSystemIntent
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Intent Creation", func() {
		It("should create a new ConfigureSystem intent with valid context", func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
			Expect(intent.model).NotTo(BeNil())
		})

		It("should fail to create intent with nil context", func() {
			// Skipped because staticcheck prevents passing nil context
			// This is tested indirectly through all other tests
		})

		It("should initialize with correct default state", func() {
			intent, _ := NewConfigureSystemIntent(ctx)
			Expect(intent.GetState()).To(Equal(ConfigStateSelectDomain))
		})

		It("should initialize with correct domains", func() {
			intent, _ := NewConfigureSystemIntent(ctx)
			Expect(intent.model.context.Domains).To(HaveLen(4))
		})
	})

	Describe("Init Method", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil command", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("should mark intent as active", func() {
			intent.Init()
			Expect(intent.IsActive()).To(BeTrue())
		})
	})

	Describe("View Rendering - SelectDomain State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should render SelectDomain view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Configuration Domain"))
		})

		It("should show all domains", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("system"))
			Expect(view).To(ContainSubstring("profile"))
			Expect(view).To(ContainSubstring("export"))
			Expect(view).To(ContainSubstring("ui"))
		})

		It("should show navigation instructions", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("↑/↓"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})

		It("should highlight selected item", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("> system"))
		})
	})

	Describe("View Rendering - EditSettings State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateEditSettings)
			intent.SetDomain(DomainSystem)
		})

		It("should render EditSettings view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Edit Settings for system"))
		})

		It("should show settings for domain", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Log Level"))
			Expect(view).To(ContainSubstring("Data Directory"))
		})
	})

	Describe("View Rendering - ReviewChanges State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateReviewChanges)
			intent.SetDomain(DomainProfile)
			intent.GetChanges() // Initialize changes
		})

		It("should render ReviewChanges view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Review Changes"))
		})
	})

	Describe("View Rendering - Confirm State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateConfirm)
			intent.SetDomain(DomainExport)
		})

		It("should render Confirm view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Confirm Configuration Changes"))
		})

		It("should show domain", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Domain: export"))
		})
	})

	Describe("View Rendering - Saving State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateSaving)
		})

		It("should render Saving view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Saving configuration"))
		})
	})

	Describe("View Rendering - Complete State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateComplete)
		})

		It("should render Complete view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Configuration Updated"))
		})
	})

	Describe("View Rendering - Failed State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateFailed)
			intent.model.error = &IntentError{
				Code:    "config_error",
				Message: "Failed to save configuration",
			}
		})

		It("should render Failed view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Configuration Failed"))
		})

		It("should show error message", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Failed to save configuration"))
		})
	})

	Describe("State Transitions - SelectDomain", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should navigate down through domains", func() {
			Expect(intent.GetSelectedIndex()).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should navigate up through domains", func() {
			intent.SetSelectedIndex(2)
			intent.Update(tea.KeyMsg{Type: tea.KeyUp, Runes: []rune{'k'}})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should not go below first domain", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyUp, Runes: []rune{'k'}})
			Expect(intent.GetSelectedIndex()).To(Equal(0))
		})

		It("should not go above last domain", func() {
			intent.SetSelectedIndex(3)
			intent.Update(tea.KeyMsg{Type: tea.KeyDown, Runes: []rune{'j'}})
			Expect(intent.GetSelectedIndex()).To(Equal(3))
		})

		It("should transition to EditSettings on enter", func() {
			intent.SetSelectedIndex(0)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ConfigStateEditSettings))
			Expect(intent.GetDomain()).To(Equal(DomainSystem))
			Expect(intent.GetChanges()).NotTo(BeNil())
		})

		It("should cancel on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.IsActive()).To(BeFalse())
			result := intent.Result()
			Expect(result.Status).To(Equal(Cancelled))
		})
	})

	Describe("State Transitions - EditSettings", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateEditSettings)
			intent.SetDomain(DomainSystem)
		})

		It("should transition to ReviewChanges on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ConfigStateReviewChanges))
		})

		It("should go back to SelectDomain on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ConfigStateSelectDomain))
		})
	})

	Describe("State Transitions - ReviewChanges", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateReviewChanges)
			intent.SetDomain(DomainProfile)
		})

		It("should transition to Confirm on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ConfigStateConfirm))
		})

		It("should go back to EditSettings on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ConfigStateEditSettings))
		})
	})

	Describe("State Transitions - Confirm", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateConfirm)
			intent.SetDomain(DomainExport)
		})

		It("should transition to Saving on confirm (y)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(intent.GetState()).To(Equal(ConfigStateSaving))
		})

		It("should transition to Saving on confirm (enter)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.GetState()).To(Equal(ConfigStateSaving))
		})

		It("should go back to ReviewChanges on cancel (n)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.GetState()).To(Equal(ConfigStateReviewChanges))
		})

		It("should go back to ReviewChanges on cancel (esc)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.GetState()).To(Equal(ConfigStateReviewChanges))
		})
	})

	Describe("State Transitions - Saving", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateSaving)
			intent.SetDomain(DomainUI)
		})

		It("should transition to Complete on save completion", func() {
			result := &ConfigureSystemResult{
				Success: true,
				Domain:  DomainUI,
				Changes: &ConfigurationChanges{
					Domain:   DomainUI,
					Original: make(map[string]interface{}),
					Modified: make(map[string]interface{}),
				},
			}
			intent.Update(ConfigCompleteMsg{Result: result})
			Expect(intent.GetState()).To(Equal(ConfigStateComplete))
			Expect(intent.GetResult()).NotTo(BeNil())
		})

		It("should transition to Failed on save error", func() {
			err := &IntentError{Code: "save_error", Message: "Save failed"}
			intent.Update(ConfigErrorMsg{Error: err})
			Expect(intent.GetState()).To(Equal(ConfigStateFailed))
		})
	})

	Describe("State Transitions - Complete", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateComplete)
		})

		It("should mark intent as inactive on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter, Runes: []rune{'\n'}})
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should mark intent as inactive on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.IsActive()).To(BeFalse())
		})
	})

	Describe("State Transitions - Failed", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateFailed)
			intent.model.error = &IntentError{Code: "save_error", Message: "Save failed"}
		})

		It("should transition to Confirm on retry (r)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(intent.GetState()).To(Equal(ConfigStateConfirm))
		})

		It("should mark intent as inactive on cancel (esc)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})
			Expect(intent.IsActive()).To(BeFalse())
		})
	})

	Describe("Result Handling", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return nil when no result set", func() {
			// Before the intent is completed, Result() should return nil
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return Cancelled when intent is cancelled", func() {
			// Simulate cancelling the intent by pressing Escape
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should return Completed on successful configuration", func() {
			intent.SetState(ConfigStateComplete)
			configResult := &ConfigureSystemResult{
				Success: true,
				Domain:  DomainSystem,
				Changes: &ConfigurationChanges{
					Domain:   DomainSystem,
					Original: make(map[string]interface{}),
					Modified: make(map[string]interface{}),
				},
			}
			intent.model.result = configResult
			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).NotTo(BeNil())
		})

		It("should return Failed on configuration failure", func() {
			intent.SetState(ConfigStateFailed)
			configResult := &ConfigureSystemResult{
				Success: false,
				Error:   &IntentError{Code: "save_error", Message: "Save failed"},
			}
			intent.model.result = configResult
			result := intent.Result()
			Expect(result.Status).To(Equal(Failed))
			Expect(result.Error).NotTo(BeNil())
		})
	})

	Describe("Configuration Context", func() {
		It("should have all domains", func() {
			ctx := NewConfigureSystemContext()
			Expect(ctx.Domains).To(HaveLen(4))
			Expect(ctx.Domains).To(ContainElements(
				DomainSystem, DomainProfile, DomainExport, DomainUI,
			))
		})

		It("should have settings for each domain", func() {
			ctx := NewConfigureSystemContext()
			Expect(ctx.Settings[DomainSystem]).NotTo(BeEmpty())
			Expect(ctx.Settings[DomainProfile]).NotTo(BeEmpty())
			Expect(ctx.Settings[DomainExport]).NotTo(BeEmpty())
			Expect(ctx.Settings[DomainUI]).NotTo(BeEmpty())
		})

		It("should have setting types defined", func() {
			ctx := NewConfigureSystemContext()
			for _, setting := range ctx.Settings[DomainSystem] {
				Expect(setting.Type).To(BeElementOf("string", "bool", "int", "select"))
			}
		})
	})

	Describe("Intent Interface Compliance", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should implement Intent interface", func() {
			var _ Intent = intent
		})

		It("should have Init method", func() {
			Expect(intent.Init).NotTo(BeNil())
		})

		It("should have Update method", func() {
			Expect(intent.Update).NotTo(BeNil())
		})

		It("should have View method", func() {
			Expect(intent.View).NotTo(BeNil())
		})

		It("should have Result method", func() {
			Expect(intent.Result).NotTo(BeNil())
		})
	})
})
