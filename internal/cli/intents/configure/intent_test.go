package configure_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/configure"
	"github.com/baphled/kariya/internal/config"
)

var _ = Describe("ConfigureSystem Intent", func() {
	var (
		intent    *configure.Intent
		intentCtx *configure.IntentContext
		cfg       *config.Config
	)

	BeforeEach(func() {
		cfg = config.DefaultConfig()
		intentCtx = &configure.IntentContext{
			Cfg:      cfg,
			Settings: configure.SettingsFromConfig(cfg),
		}
	})

	Describe("Intent Creation", func() {
		It("should create a new intent with valid context", func() {
			var err error
			intent, err = configure.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should return error with nil config", func() {
			intentCtx.Cfg = nil
			_, err := configure.NewIntent(intentCtx)
			Expect(err).To(HaveOccurred())
		})

		It("should initialise with correct default state", func() {
			intent, _ = configure.NewIntent(intentCtx)
			Expect(intent.GetState()).To(Equal(configure.ConfigStateSelectDomain))
		})

		It("should initialise with settings for all domains", func() {
			intent, _ = configure.NewIntent(intentCtx)
			settings := intent.GetSettings()
			Expect(settings).NotTo(BeNil())
			Expect(settings).To(HaveKey(configure.DomainSystem))
			Expect(settings).To(HaveKey(configure.DomainProfile))
			Expect(settings).To(HaveKey(configure.DomainExport))
			Expect(settings).To(HaveKey(configure.DomainUI))
		})

		It("should embed BaseIntent", func() {
			intent, _ = configure.NewIntent(intentCtx)
			Expect(intent.BaseIntent).NotTo(BeNil())
		})

		It("should not be active initially", func() {
			intent, _ = configure.NewIntent(intentCtx)
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should store the context", func() {
			intent, _ = configure.NewIntent(intentCtx)
			Expect(intent.GetContext()).To(Equal(intentCtx))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			intent, _ = configure.NewIntent(intentCtx)
		})

		It("should mark intent as active", func() {
			intent.Init()
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should return a command or nil", func() {
			cmd := intent.Init()
			_ = cmd
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			intent, _ = configure.NewIntent(intentCtx)
			intent.Init()
		})

		Context("in SelectDomain state", func() {
			It("should show all domains", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("System"))
				Expect(view).To(ContainSubstring("Profile"))
				Expect(view).To(ContainSubstring("Export"))
			})
		})

		Context("in EditSettings state", func() {
			BeforeEach(func() {
				intent.SetDomain(configure.DomainSystem)
				intent.SetState(configure.ConfigStateEditSettings)
			})

			It("should be in EditSettings state", func() {
				Expect(intent.GetState()).To(Equal(configure.ConfigStateEditSettings))
			})
		})

		Context("in Saving state", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateSaving)
			})

			It("should show saving modal", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Saving"),
					ContainSubstring("Loading"),
				))
			})
		})

		Context("in Complete state", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateComplete)
			})

			It("should show success modal", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Complete"),
					ContainSubstring("Success"),
					ContainSubstring("saved"),
				))
			})
		})

		Context("in Failed state", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateFailed)
			})

			It("should show error modal", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Failed"),
					ContainSubstring("Error"),
				))
			})
		})
	})

	Describe("State Transitions", func() {
		BeforeEach(func() {
			intent, _ = configure.NewIntent(intentCtx)
			intent.Init()
		})

		Context("from SelectDomain", func() {
			It("should cancel on escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.IsActive()).To(BeFalse())
				result := intent.Result()
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should cancel on q key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(intent.IsActive()).To(BeFalse())
			})
		})

		Context("from Saving", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateSaving)
				intent.SetDomain(configure.DomainUI)
			})

			It("should transition to Complete on save completion", func() {
				result := &configure.SystemResult{
					Success: true,
					Domain:  configure.DomainUI,
					Changes: &configure.ConfigurationChanges{
						Domain:   configure.DomainUI,
						Original: make(map[string]interface{}),
						Modified: make(map[string]interface{}),
					},
				}
				intent.Update(configure.ConfigCompleteMsg{Result: result})
				Expect(intent.GetState()).To(Equal(configure.ConfigStateComplete))
				Expect(intent.GetResult()).NotTo(BeNil())
			})

			It("should transition to Failed on save error", func() {
				err := &intents.IntentError{Code: "save_error", Message: "Save failed"}
				intent.Update(configure.ConfigErrorMsg{Error: err})
				Expect(intent.GetState()).To(Equal(configure.ConfigStateFailed))
			})
		})

		Context("from Complete", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateComplete)
			})

			It("should mark intent as inactive on enter", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.IsActive()).To(BeFalse())
			})

			It("should mark intent as inactive on escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.IsActive()).To(BeFalse())
			})
		})

		Context("from Failed", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateFailed)
			})

			It("should go back to edit on enter (to retry)", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.GetState()).To(Equal(configure.ConfigStateEditSettings))
			})
		})
	})

	Describe("Result Handling", func() {
		BeforeEach(func() {
			intent, _ = configure.NewIntent(intentCtx)
			intent.Init()
		})

		It("should return nil when no result set", func() {
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return Cancelled when intent is cancelled", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("Intent Interface Compliance", func() {
		BeforeEach(func() {
			intent, _ = configure.NewIntent(intentCtx)
		})

		It("should implement Intent interface", func() {
			var _ intents.Intent = intent
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

	Describe("Config Data Loading", func() {
		It("should have config available after creation", func() {
			intent, err := configure.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent.GetConfig()).NotTo(BeNil())
		})

		It("should have settings populated for all domains", func() {
			intent, _ = configure.NewIntent(intentCtx)
			settings := intent.GetSettings()
			Expect(settings).NotTo(BeNil())
			Expect(settings).To(HaveKey(configure.DomainSystem))
			Expect(settings).To(HaveKey(configure.DomainProfile))
			Expect(settings).To(HaveKey(configure.DomainExport))
			Expect(settings).To(HaveKey(configure.DomainUI))
		})

		It("should return settings for a specific domain", func() {
			intent, _ = configure.NewIntent(intentCtx)
			settings := intent.GetSettingsForDomain(configure.DomainSystem)
			Expect(settings).NotTo(BeNil())
			Expect(settings).ToNot(BeEmpty())

			var hasLogLevel, hasDataDir bool
			for _, s := range settings {
				if s.Key == "log_level" {
					hasLogLevel = true
				}
				if s.Key == "data_dir" {
					hasDataDir = true
				}
			}
			Expect(hasLogLevel).To(BeTrue(), "Should have log_level setting from config")
			Expect(hasDataDir).To(BeTrue(), "Should have data_dir setting from config")
		})
	})

	Describe("Modal Rendering", func() {
		BeforeEach(func() {
			intent, _ = configure.NewIntent(intentCtx)
			intent.Init()
		})

		Context("saving modal", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateSaving)
			})

			It("should show saving modal content", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Saving"),
					ContainSubstring("Loading"),
					ContainSubstring("⏳"),
				))
			})

			It("should show saving modal with proper styling", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("╭"),
					ContainSubstring("╮"),
					ContainSubstring("│"),
				))
			})
		})

		Context("success modal", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateComplete)
			})

			It("should show success modal content", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Success"),
					ContainSubstring("Complete"),
					ContainSubstring("✅"),
					ContainSubstring("saved"),
				))
			})

			It("should show success modal with proper styling", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("╭"),
					ContainSubstring("╮"),
				))
			})
		})

		Context("error modal", func() {
			BeforeEach(func() {
				intent.SetState(configure.ConfigStateFailed)
			})

			It("should show error modal content", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Failed"),
					ContainSubstring("Error"),
					ContainSubstring("⚠"),
				))
			})

			It("should show error modal with proper styling", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("╭"),
					ContainSubstring("╮"),
				))
			})
		})
	})
})
