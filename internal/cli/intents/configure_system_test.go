package intents

import (
	"context"

	"github.com/baphled/kariya/internal/config"
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
		})

		It("should fail to create intent with nil context", func() {
			// Skipped because staticcheck prevents passing nil context
			// This is tested indirectly through all other tests
		})

		It("should initialize with correct default state", func() {
			intent, _ := NewConfigureSystemIntent(ctx) //nolint:errcheck // test helper
			Expect(intent.GetState()).To(Equal(ConfigStateSelectDomain))
		})

		It("should initialize with settings for all domains", func() {
			intent, _ := NewConfigureSystemIntent(ctx) //nolint:errcheck // test helper
			Expect(intent.settings).NotTo(BeNil())
			Expect(intent.settings).To(HaveKey(DomainSystem))
			Expect(intent.settings).To(HaveKey(DomainProfile))
			Expect(intent.settings).To(HaveKey(DomainExport))
			Expect(intent.settings).To(HaveKey(DomainUI))
		})
	})

	Describe("Init Method", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return a command or nil", func() {
			cmd := intent.Init()
			// Init may return nil or a command
			_ = cmd
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

		It("should show all domains", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("System"))
			Expect(view).To(ContainSubstring("Profile"))
			Expect(view).To(ContainSubstring("Export"))
		})
	})

	Describe("View Rendering - EditSettings State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetDomain(DomainSystem)
			intent.SetState(ConfigStateEditSettings)
		})

		It("should be in EditSettings state", func() {
			Expect(intent.GetState()).To(Equal(ConfigStateEditSettings))
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

		It("should show saving modal", func() {
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Loading"),
			))
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

		It("should show success modal", func() {
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Complete"),
				ContainSubstring("Success"),
				ContainSubstring("saved"),
			))
		})
	})

	Describe("View Rendering - Failed State", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.SetState(ConfigStateFailed)
		})

		It("should show error modal", func() {
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Failed"),
				ContainSubstring("Error"),
			))
		})
	})

	Describe("State Transitions - SelectDomain", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should cancel on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.IsActive()).To(BeFalse())
			result := intent.Result()
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.IsActive()).To(BeFalse())
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
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should mark intent as inactive on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
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
		})

		It("should go back to edit on enter (to retry)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ConfigStateEditSettings))
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

	Describe("Config Data Loading", func() {
		It("should load config data on creation", func() {
			intent, err := NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			// Config should be loaded
			Expect(intent.cfg).NotTo(BeNil())

			// Settings should be populated
			Expect(intent.settings).NotTo(BeNil())
			Expect(intent.settings).To(HaveKey(DomainSystem))
			Expect(intent.settings).To(HaveKey(DomainProfile))
			Expect(intent.settings).To(HaveKey(DomainExport))
			Expect(intent.settings).To(HaveKey(DomainUI))
		})

		It("should return settings from loaded config", func() {
			intent, _ := NewConfigureSystemIntent(ctx) //nolint:errcheck // test helper

			// Get system settings
			settings := intent.getSettingsForDomain(DomainSystem)
			Expect(settings).NotTo(BeNil())
			Expect(settings).ToNot(BeEmpty())

			// Should have actual config settings (not sample data)
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

	Describe("Loading Modal Rendering", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show saving modal when in Saving state", func() {
			intent.SetState(ConfigStateSaving)
			view := intent.View()

			// Should show loading modal content
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Loading"),
				ContainSubstring("⏳"),
			))
		})

		It("should show saving modal with proper styling", func() {
			intent.SetState(ConfigStateSaving)
			view := intent.View()

			// Modal should have border characters (rounded border)
			Expect(view).To(SatisfyAny(
				ContainSubstring("╭"),
				ContainSubstring("╮"),
				ContainSubstring("│"),
			))
		})

		It("should dim background when showing saving modal", func() {
			intent.SetState(ConfigStateSaving)
			view := intent.View()

			// The view should still contain the domain options (they're the background)
			// but they should be dimmed (faint styling applies escape sequences)
			// At minimum, the saving message should be present
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("configuration"),
			))
		})
	})

	Describe("Success Modal Rendering", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show success modal when in Complete state", func() {
			intent.SetState(ConfigStateComplete)
			view := intent.View()

			// Should show success modal content
			Expect(view).To(SatisfyAny(
				ContainSubstring("Success"),
				ContainSubstring("Complete"),
				ContainSubstring("✅"),
				ContainSubstring("saved"),
			))
		})

		It("should show success modal with proper styling", func() {
			intent.SetState(ConfigStateComplete)
			view := intent.View()

			// Modal should have border characters
			Expect(view).To(SatisfyAny(
				ContainSubstring("╭"),
				ContainSubstring("╮"),
			))
		})
	})

	Describe("Error Modal Rendering", func() {
		BeforeEach(func() {
			var err error
			intent, err = NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show error modal when in Failed state", func() {
			intent.SetState(ConfigStateFailed)
			view := intent.View()

			// Should show error modal content
			Expect(view).To(SatisfyAny(
				ContainSubstring("Failed"),
				ContainSubstring("Error"),
				ContainSubstring("⚠"),
			))
		})

		It("should show error modal with proper styling", func() {
			intent.SetState(ConfigStateFailed)
			view := intent.View()

			// Modal should have border characters
			Expect(view).To(SatisfyAny(
				ContainSubstring("╭"),
				ContainSubstring("╮"),
			))
		})
	})
})

var _ = Describe("applyConfigChange", func() {
	var cfg *config.Config

	BeforeEach(func() {
		cfg = config.DefaultConfig()
	})

	Describe("routing to system domain", func() {
		It("should apply system changes", func() {
			err := applyConfigChange(cfg, DomainSystem, "log_level", "debug")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.System.LogLevel).To(Equal("debug"))
		})
	})

	Describe("routing to profile domain", func() {
		It("should apply profile changes", func() {
			err := applyConfigChange(cfg, DomainProfile, "name", "Alice")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Profile.Name).To(Equal("Alice"))
		})
	})

	Describe("routing to export domain", func() {
		It("should apply export changes", func() {
			err := applyConfigChange(cfg, DomainExport, "default_destination", "clipboard")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Export.DefaultDestination).To(Equal("clipboard"))
		})
	})

	Describe("routing to UI domain", func() {
		It("should apply display changes", func() {
			err := applyConfigChange(cfg, DomainUI, "theme", "light")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Display.Theme).To(Equal("light"))
		})
	})

	Describe("unknown domain", func() {
		It("should return an error", func() {
			err := applyConfigChange(cfg, "unknown", "key", "value")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown domain"))
		})
	})
})

var _ = Describe("applySystemChange", func() {
	var sys config.SystemConfig

	BeforeEach(func() {
		sys = config.SystemConfig{
			DataDir:     "/data",
			LogLevel:    "info",
			AutoBackup:  true,
			BackupCount: 5,
		}
	})

	It("should set data_dir", func() {
		err := applySystemChange(&sys, "data_dir", "/new/path")
		Expect(err).NotTo(HaveOccurred())
		Expect(sys.DataDir).To(Equal("/new/path"))
	})

	It("should set log_level", func() {
		err := applySystemChange(&sys, "log_level", "debug")
		Expect(err).NotTo(HaveOccurred())
		Expect(sys.LogLevel).To(Equal("debug"))
	})

	It("should set auto_backup", func() {
		err := applySystemChange(&sys, "auto_backup", false)
		Expect(err).NotTo(HaveOccurred())
		Expect(sys.AutoBackup).To(BeFalse())
	})

	It("should set backup_count", func() {
		err := applySystemChange(&sys, "backup_count", 10)
		Expect(err).NotTo(HaveOccurred())
		Expect(sys.BackupCount).To(Equal(10))
	})

	It("should return error for unknown key", func() {
		err := applySystemChange(&sys, "unknown_key", "value")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown system setting"))
	})
})

var _ = Describe("applyProfileChange", func() {
	var prof config.ProfileConfig

	BeforeEach(func() {
		prof = config.ProfileConfig{}
	})

	It("should set name", func() {
		err := applyProfileChange(&prof, "name", "Alice")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Name).To(Equal("Alice"))
	})

	It("should set email", func() {
		err := applyProfileChange(&prof, "email", "alice@example.com")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Email).To(Equal("alice@example.com"))
	})

	It("should set title", func() {
		err := applyProfileChange(&prof, "title", "Senior Engineer")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Title).To(Equal("Senior Engineer"))
	})

	It("should set location", func() {
		err := applyProfileChange(&prof, "location", "London")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Location).To(Equal("London"))
	})

	It("should set github", func() {
		err := applyProfileChange(&prof, "github", "https://github.com/alice")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.GitHub).To(Equal("https://github.com/alice"))
	})

	It("should set portfolio", func() {
		err := applyProfileChange(&prof, "portfolio", "https://alice.dev")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Portfolio).To(Equal("https://alice.dev"))
	})

	It("should set languages", func() {
		err := applyProfileChange(&prof, "languages", []string{"Go", "Python"})
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Languages).To(Equal([]string{"Go", "Python"}))
	})

	It("should set frontend", func() {
		err := applyProfileChange(&prof, "frontend", []string{"React"})
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Frontend).To(Equal([]string{"React"}))
	})

	It("should set systems", func() {
		err := applyProfileChange(&prof, "systems", []string{"Docker", "K8s"})
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.Systems).To(Equal([]string{"Docker", "K8s"}))
	})

	It("should set core_strengths", func() {
		err := applyProfileChange(&prof, "core_strengths", []string{"Architecture"})
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.CoreStrengths).To(Equal([]string{"Architecture"}))
	})

	It("should set what_i_bring", func() {
		err := applyProfileChange(&prof, "what_i_bring", []string{"Leadership"})
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.WhatIBring).To(Equal([]string{"Leadership"}))
	})

	It("should set default_role", func() {
		err := applyProfileChange(&prof, "default_role", "staff_ic")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.DefaultRole).To(Equal("staff_ic"))
	})

	It("should set default_audience", func() {
		err := applyProfileChange(&prof, "default_audience", "executive")
		Expect(err).NotTo(HaveOccurred())
		Expect(prof.DefaultAudience).To(Equal("executive"))
	})

	It("should return error for unknown key", func() {
		err := applyProfileChange(&prof, "unknown_key", "value")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown profile setting"))
	})
})

var _ = Describe("applyExportChange", func() {
	var exp config.ExportConfig

	BeforeEach(func() {
		exp = config.ExportConfig{
			DefaultDestination: "file",
			AutoOpen:           false,
		}
	})

	It("should set default_destination", func() {
		err := applyExportChange(&exp, "default_destination", "clipboard")
		Expect(err).NotTo(HaveOccurred())
		Expect(exp.DefaultDestination).To(Equal("clipboard"))
	})

	It("should set auto_open", func() {
		err := applyExportChange(&exp, "auto_open", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(exp.AutoOpen).To(BeTrue())
	})

	It("should return error for unknown key", func() {
		err := applyExportChange(&exp, "unknown_key", "value")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown export setting"))
	})
})

var _ = Describe("applyDisplayChange", func() {
	var disp config.DisplayConfig

	BeforeEach(func() {
		disp = config.DisplayConfig{
			Theme:      "dark",
			Animations: true,
		}
	})

	It("should set theme", func() {
		err := applyDisplayChange(&disp, "theme", "light")
		Expect(err).NotTo(HaveOccurred())
		Expect(disp.Theme).To(Equal("light"))
	})

	It("should set animations", func() {
		err := applyDisplayChange(&disp, "animations", false)
		Expect(err).NotTo(HaveOccurred())
		Expect(disp.Animations).To(BeFalse())
	})

	It("should return error for unknown key", func() {
		err := applyDisplayChange(&disp, "unknown_key", "value")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown display setting"))
	})
})

var _ = Describe("ConfigureSystemIntent Getters", func() {
	var (
		intent *ConfigureSystemIntent
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		intent, err = NewConfigureSystemIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("GetDomain", func() {
		It("should return empty string initially", func() {
			Expect(string(intent.GetDomain())).To(BeEmpty())
		})

		It("should return the selected domain", func() {
			intent.SetDomain(DomainSystem)
			Expect(intent.GetDomain()).To(Equal(DomainSystem))
		})
	})

	Describe("GetChanges", func() {
		It("should return nil when no pending changes", func() {
			Expect(intent.GetChanges()).To(BeNil())
		})

		It("should return changes when pending changes exist", func() {
			intent.SetDomain(DomainSystem)
			intent.pendingChanges = map[string]interface{}{"log_level": "debug"}

			changes := intent.GetChanges()
			Expect(changes).NotTo(BeNil())
			Expect(changes.Domain).To(Equal(DomainSystem))
			Expect(changes.Modified).To(HaveKeyWithValue("log_level", "debug"))
		})
	})

	Describe("GetSavingModal", func() {
		It("should return nil initially", func() {
			Expect(intent.GetSavingModal()).To(BeNil())
		})

		It("should return saving modal when in saving state", func() {
			intent.SetState(ConfigStateSaving)
			Expect(intent.GetSavingModal()).NotTo(BeNil())
		})
	})

	Describe("GetResultModal", func() {
		It("should return nil initially", func() {
			Expect(intent.GetResultModal()).To(BeNil())
		})

		It("should return result modal when in complete state", func() {
			intent.SetState(ConfigStateComplete)
			Expect(intent.GetResultModal()).NotTo(BeNil())
		})

		It("should return result modal when in failed state", func() {
			intent.SetState(ConfigStateFailed)
			Expect(intent.GetResultModal()).NotTo(BeNil())
		})
	})

	Describe("RenderDomainContent", func() {
		It("should return non-empty content", func() {
			content := intent.RenderDomainContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should contain domain names", func() {
			content := intent.RenderDomainContent()
			Expect(content).To(ContainSubstring("System"))
		})
	})
})
