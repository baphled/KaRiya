package configure

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/terminal"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

var _ = Describe("Configure intent internals", func() {
	var (
		cfg    *config.Config
		intent *Intent
	)

	BeforeEach(func() {
		cfg = config.DefaultConfig()
		config.SetConfigPathForTesting(filepath.Join(GinkgoT().TempDir(), "config.yaml"))
		ctx := &IntentValidator{Cfg: cfg, Settings: settingsFromConfig(cfg)}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("applySystemChange covers all branches", func() {
		It("applies data_dir", func() {
			Expect(applySystemChange(&cfg.System, "data_dir", "/tmp/data")).To(Succeed())
			Expect(cfg.System.DataDir).To(Equal("/tmp/data"))
		})

		It("applies log_level", func() {
			Expect(applySystemChange(&cfg.System, "log_level", "warn")).To(Succeed())
			Expect(cfg.System.LogLevel).To(Equal("warn"))
		})

		It("applies auto_backup", func() {
			Expect(applySystemChange(&cfg.System, "auto_backup", true)).To(Succeed())
			Expect(cfg.System.AutoBackup).To(BeTrue())
		})

		It("applies backup_count", func() {
			Expect(applySystemChange(&cfg.System, "backup_count", 10)).To(Succeed())
			Expect(cfg.System.BackupCount).To(Equal(10))
		})

		It("rejects unknown key", func() {
			Expect(applySystemChange(&cfg.System, "unknown", "value")).To(HaveOccurred())
		})
	})

	Describe("applyConfigChange covers all domains", func() {
		It("applies to system", func() {
			Expect(applyConfigChange(cfg, DomainSystem, "log_level", "debug")).To(Succeed())
		})

		It("applies to profile", func() {
			Expect(applyConfigChange(cfg, DomainProfile, "name", "Ada")).To(Succeed())
		})

		It("applies to export", func() {
			Expect(applyConfigChange(cfg, DomainExport, "default_destination", "clipboard")).To(Succeed())
		})

		It("applies to display", func() {
			Expect(applyConfigChange(cfg, DomainUI, "theme", "light")).To(Succeed())
		})

		It("rejects unknown domain", func() {
			Expect(applyConfigChange(cfg, "unknown", "key", "value")).To(HaveOccurred())
		})
	})

	Describe("applyProfileChange", func() {
		It("applies string fields", func() {
			Expect(applyProfileChange(&cfg.Profile, "email", "test@test.com")).To(Succeed())
			Expect(cfg.Profile.Email).To(Equal("test@test.com"))

			Expect(applyProfileChange(&cfg.Profile, "title", "Engineer")).To(Succeed())
			Expect(cfg.Profile.Title).To(Equal("Engineer"))

			Expect(applyProfileChange(&cfg.Profile, "location", "London")).To(Succeed())
			Expect(cfg.Profile.Location).To(Equal("London"))

			Expect(applyProfileChange(&cfg.Profile, "github", "https://github.com/test")).To(Succeed())
			Expect(cfg.Profile.GitHub).To(Equal("https://github.com/test"))

			Expect(applyProfileChange(&cfg.Profile, "portfolio", "https://test.com")).To(Succeed())
			Expect(cfg.Profile.Portfolio).To(Equal("https://test.com"))

			Expect(applyProfileChange(&cfg.Profile, "default_role", "staff_ic")).To(Succeed())
			Expect(cfg.Profile.DefaultRole).To(Equal("staff_ic"))

			Expect(applyProfileChange(&cfg.Profile, "default_audience", "executive")).To(Succeed())
			Expect(cfg.Profile.DefaultAudience).To(Equal("executive"))
		})

		It("applies slice fields", func() {
			Expect(applyProfileChange(&cfg.Profile, "languages", []string{"Go"})).To(Succeed())
			Expect(cfg.Profile.Languages).To(ContainElement("Go"))

			Expect(applyProfileChange(&cfg.Profile, "frontend", []string{"React"})).To(Succeed())
			Expect(cfg.Profile.Frontend).To(ContainElement("React"))

			Expect(applyProfileChange(&cfg.Profile, "systems", []string{"Linux"})).To(Succeed())
			Expect(cfg.Profile.Systems).To(ContainElement("Linux"))

			Expect(applyProfileChange(&cfg.Profile, "core_strengths", []string{"Go"})).To(Succeed())
			Expect(cfg.Profile.CoreStrengths).To(ContainElement("Go"))

			Expect(applyProfileChange(&cfg.Profile, "what_i_bring", []string{"TDD"})).To(Succeed())
			Expect(cfg.Profile.WhatIBring).To(ContainElement("TDD"))
		})

		It("rejects unknown key", func() {
			Expect(applyProfileChange(&cfg.Profile, "unknown", "value")).To(HaveOccurred())
		})
	})

	Describe("applyExportChange", func() {
		It("applies default_destination", func() {
			Expect(applyExportChange(&cfg.Export, "default_destination", "clipboard")).To(Succeed())
			Expect(cfg.Export.DefaultDestination).To(Equal("clipboard"))
		})

		It("applies auto_open", func() {
			Expect(applyExportChange(&cfg.Export, "auto_open", true)).To(Succeed())
			Expect(cfg.Export.AutoOpen).To(BeTrue())
		})

		It("rejects unknown key", func() {
			Expect(applyExportChange(&cfg.Export, "unknown", "value")).To(HaveOccurred())
		})
	})

	Describe("applyDisplayChange", func() {
		It("applies theme", func() {
			Expect(applyDisplayChange(&cfg.Display, "theme", "light")).To(Succeed())
			Expect(cfg.Display.Theme).To(Equal("light"))
		})

		It("applies animations", func() {
			Expect(applyDisplayChange(&cfg.Display, "animations", false)).To(Succeed())
			Expect(cfg.Display.Animations).To(BeFalse())
		})

		It("rejects unknown key", func() {
			Expect(applyDisplayChange(&cfg.Display, "unknown", "value")).To(HaveOccurred())
		})
	})

	Describe("ApplyChanges", func() {
		It("applies changes across domains", func() {
			changes := map[string]interface{}{
				"log_level": "warn",
				"theme":     "light",
			}
			settings := settingsFromConfig(cfg)
			Expect(ApplyChanges(cfg, settings, changes)).To(Succeed())
			Expect(cfg.System.LogLevel).To(Equal("warn"))
			Expect(cfg.Display.Theme).To(Equal("light"))
		})

		It("ignores keys not present in settings", func() {
			changes := map[string]interface{}{
				"nonexistent_key": "value",
			}
			settings := settingsFromConfig(cfg)
			Expect(ApplyChanges(cfg, settings, changes)).To(Succeed())
		})
	})

	Describe("settings modal flow", func() {
		It("openSettingsModal creates modal", func() {
			intent.clearAllModals()
			intent.openSettingsModal()
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("updateSettingsModal passes through messages", func() {
			intent.openSettingsModal()
			cmd := intent.updateSettingsModal(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("handles cancel from settings modal", func() {
			intent.openSettingsModal()
			intent.updateSettingsModal(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.settingsModal).To(BeNil())
			Expect(intent.active).To(BeFalse())
		})

		It("handles completion from settings modal", func() {
			intent.openSettingsModal()
			intent.selectedDomain = DomainSystem
			cmd := intent.updateSettingsModal(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.settingsModal).To(BeNil())
			Expect(intent.state).To(Equal(ConfigStateSaving))
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("saving modal flow", func() {
		BeforeEach(func() {
			intent.settings = settingsFromConfig(cfg)
			intent.selectedDomain = DomainSystem
		})

		It("saving modal complete", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			result := &SystemResult{Success: true}
			cmd := intent.updateSavingModal(ConfigCompleteMsg{Result: result})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.state).To(Equal(ConfigStateComplete))
		})

		It("saving modal error", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			errMsg := ConfigErrorMsg{Error: &intents.IntentError{Code: "save", Message: "fail"}}
			cmd := intent.updateSavingModal(errMsg)
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(ConfigStateFailed))
		})

		It("saving modal spinner tick", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			cmd := intent.updateSavingModal(feedback.ModalSpinnerTickMsg{})
			_ = cmd
		})

		It("saving modal ignores unrelated", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			cmd := intent.updateSavingModal(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
		})

		It("startSaving creates batch command", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "info"}
			cmd := intent.startSaving()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			Expect(msg).NotTo(BeNil())
		})
	})

	Describe("render helpers", func() {
		It("reports all state names", func() {
			intent.clearAllModals()
			Expect(intent.getStateName()).To(Equal("Select Domain"))

			intent.openSettingsModal()
			Expect(intent.getStateName()).To(Equal("Configure"))

			intent.settingsModal = nil
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			Expect(intent.getStateName()).To(Equal("Saving"))

			intent.savingModal = nil
			intent.configResult = &SystemResult{Success: true}
			intent.resultModal = feedback.NewSuccessModal("Done")
			Expect(intent.getStateName()).To(Equal("Complete"))

			intent.configResult = nil
			intent.resultModal = feedback.NewErrorModal("Failed", "Error")
			Expect(intent.getStateName()).To(Equal("Failed"))
		})

		It("reports context help for all states", func() {
			intent.clearAllModals()
			Expect(intent.getContextHelp()).NotTo(BeEmpty())

			intent.openSettingsModal()
			Expect(intent.getContextHelp()).To(Equal(""))

			intent.settingsModal = nil
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			Expect(intent.getContextHelp()).NotTo(BeEmpty())

			intent.savingModal = nil
			intent.resultModal = feedback.NewSuccessModal("Done")
			Expect(intent.getContextHelp()).NotTo(BeEmpty())
		})

		It("uses terminal info when available", func() {
			intent.UpdateTerminalInfo(&terminal.Info{Width: 200, Height: 100, IsValid: true})
			width, height := intent.getTerminalDimensions()
			Expect(width).To(Equal(200))
			Expect(height).To(Equal(100))
		})

		It("uses defaults when no terminal info", func() {
			width, height := intent.getTerminalDimensions()
			Expect(width).To(BeNumerically(">", 0))
			Expect(height).To(BeNumerically(">", 0))
		})
	})

	Describe("modal adapters", func() {
		It("renders settings modal adapter", func() {
			intent.openSettingsModal()
			adapter := settingsModalAdapter{intent.settingsModal, 80, 24}
			Expect(adapter.Render(0, 0)).NotTo(BeEmpty())
		})
	})

	Describe("routing", func() {
		It("routes to result modal", func() {
			intent.configResult = &SystemResult{Success: true}
			intent.resultModal = feedback.NewSuccessModal("Done")
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(intent.active).To(BeFalse())
		})

		It("routes to saving modal", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			cmd := intent.routeToActiveComponent(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
		})

		It("routes to settings modal", func() {
			intent.openSettingsModal()
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("returns nil when no active component", func() {
			intent.clearAllModals()
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("window size propagation", func() {
		It("propagates to settings modal", func() {
			intent.openSettingsModal()
			intent.handleWindowSizeMsg(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("handles non-window messages", func() {
			intent.handleWindowSizeMsg(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("key handling", func() {
		It("handles help toggle", func() {
			result := intent.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(result).To(BeTrue())
		})

		It("handles q key when no settings modal", func() {
			intent.settingsModal = nil
			result := intent.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(result).To(BeTrue())
			Expect(intent.active).To(BeFalse())
		})

		It("suppresses q when settings modal is open", func() {
			intent.openSettingsModal()
			result := intent.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(result).To(BeFalse())
			Expect(intent.active).To(BeTrue())
		})

		It("handles non-key messages", func() {
			result := intent.handleKeyMsg(tea.MouseMsg{})
			Expect(result).To(BeFalse())
		})

		It("handles unrecognised keys", func() {
			result := intent.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(result).To(BeFalse())
		})
	})

	Describe("Update", func() {
		It("returns nil when inactive", func() {
			intent.active = false
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
		})

		It("processes window size", func() {
			cmd := intent.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			_ = cmd
		})

		It("processes key messages", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(cmd).To(BeNil())
		})

		It("routes to active component", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
		})
	})

	Describe("View", func() {
		It("returns empty when inactive", func() {
			intent.active = false
			Expect(intent.View()).To(BeEmpty())
		})

		It("renders with settings modal overlay", func() {
			intent.openSettingsModal()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with saving modal overlay", func() {
			intent.clearAllModals()
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with result modal overlay", func() {
			intent.clearAllModals()
			intent.resultModal = feedback.NewSuccessModal("Done")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders without overlay", func() {
			intent.clearAllModals()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Result", func() {
		It("returns nil when active", func() {
			Expect(intent.Result()).To(BeNil())
		})

		It("returns result when set", func() {
			intent.active = false
			intent.result = &intents.IntentResult[interface{}]{Status: intents.Completed, Data: "ok"}
			Expect(intent.Result()).NotTo(BeNil())
			Expect(intent.Result().Status).To(Equal(intents.Completed))
		})

		It("returns completed from configResult", func() {
			intent.active = false
			intent.configResult = &SystemResult{Success: true}
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Completed))
		})

		It("returns cancelled when no result", func() {
			intent.active = false
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("SetState", func() {
		It("sets select domain with settings modal", func() {
			intent.SetState(ConfigStateSelectDomain)
			Expect(intent.state).To(Equal(ConfigStateSelectDomain))
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("sets edit settings with settings modal", func() {
			intent.SetState(ConfigStateEditSettings)
			Expect(intent.state).To(Equal(ConfigStateEditSettings))
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("sets review changes with settings modal", func() {
			intent.SetState(ConfigStateReviewChanges)
			Expect(intent.state).To(Equal(ConfigStateReviewChanges))
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("sets confirm with settings modal", func() {
			intent.SetState(ConfigStateConfirm)
			Expect(intent.state).To(Equal(ConfigStateConfirm))
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("sets saving", func() {
			intent.SetState(ConfigStateSaving)
			Expect(intent.savingModal).NotTo(BeNil())
		})

		It("sets complete", func() {
			intent.SetState(ConfigStateComplete)
			Expect(intent.resultModal).NotTo(BeNil())
			Expect(intent.configResult).NotTo(BeNil())
		})

		It("sets failed", func() {
			intent.SetState(ConfigStateFailed)
			Expect(intent.resultModal).NotTo(BeNil())
		})
	})

	Describe("result modal dismiss", func() {
		It("dismisses success on space", func() {
			intent.configResult = &SystemResult{Success: true}
			intent.resultModal = feedback.NewSuccessModal("Done")
			intent.updateResultModal(tea.KeyMsg{Type: tea.KeySpace})
			Expect(intent.active).To(BeFalse())
		})

		It("dismisses failed and reopens settings modal", func() {
			intent.configResult = &SystemResult{Success: false}
			intent.resultModal = feedback.NewErrorModal("Failed", "Error")
			intent.selectedDomain = DomainSystem
			intent.settings = settingsFromConfig(cfg)
			cmd := intent.updateResultModal(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.resultModal).To(BeNil())
			Expect(intent.settingsModal).NotTo(BeNil())
		})

		It("handles countdown for success", func() {
			intent.configResult = &SystemResult{Success: true}
			intent.resultModal = feedback.NewSuccessModal("Done")
			intent.updateResultModal(feedback.ModalCountdownTickMsg{})
		})

		It("ignores countdown for error", func() {
			intent.resultModal = feedback.NewErrorModal("Failed", "Error")
			cmd := intent.updateResultModal(feedback.ModalCountdownTickMsg{})
			Expect(cmd).To(BeNil())
		})

		It("handles auto-dismiss", func() {
			intent.resultModal = feedback.NewSuccessModal("Done")
			intent.updateResultModal(feedback.ModalAutoDismissMsg{})
			Expect(intent.active).To(BeFalse())
		})

		It("dismisses nil result and reopens settings modal", func() {
			intent.configResult = nil
			intent.resultModal = feedback.NewErrorModal("Failed", "Error")
			intent.selectedDomain = DomainSystem
			intent.settings = settingsFromConfig(cfg)
			cmd := intent.updateResultModal(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.resultModal).To(BeNil())
			Expect(intent.settingsModal).NotTo(BeNil())
		})
	})

	Describe("uncovered functions", func() {
		It("SetContext sets the context", func() {
			newCtx := &IntentValidator{Cfg: cfg, Settings: settingsFromConfig(cfg)}
			intent.SetContext(newCtx)
			Expect(intent.GetContext()).To(Equal(newCtx))
		})

		It("SetActive sets active state", func() {
			intent.SetActive(false)
			Expect(intent.IsActive()).To(BeFalse())
			intent.SetActive(true)
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("GetDomain returns selected domain", func() {
			intent.SetDomain(DomainProfile)
			Expect(intent.GetDomain()).To(Equal(DomainProfile))
		})

		It("GetSavingModal returns saving modal", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			Expect(intent.GetSavingModal()).NotTo(BeNil())
		})

		It("HandleCancel deactivates intent", func() {
			cmd := intent.HandleCancel(&widgets.CancelViewResult{})
			Expect(cmd).To(BeNil())
			Expect(intent.active).To(BeFalse())
		})

		It("HandleSubmit returns nil", func() {
			cmd := intent.HandleSubmit(&widgets.SubmitViewResult{})
			Expect(cmd).To(BeNil())
		})

		It("HandleError returns nil", func() {
			cmd := intent.HandleError(&widgets.ErrorViewResult{})
			Expect(cmd).To(BeNil())
		})

		It("startSaving returns batch command", func() {
			intent.selectedDomain = DomainSystem
			intent.pendingChanges = map[string]interface{}{"system.log_level": "debug"}
			cmd := intent.startSaving()
			Expect(cmd).NotTo(BeNil())
		})

		It("ApplyChanges applies config changes", func() {
			changes := map[string]interface{}{"system.log_level": "debug"}
			err := ApplyChanges(cfg, settingsFromConfig(cfg), changes)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Init creates settings modal and returns command", func() {
			intent.settingsModal = nil
			cmd := intent.Init()
			Expect(intent.settingsModal).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("Update with inactive intent returns nil", func() {
			intent.active = false
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("navigate handler", func() {
		It("handles invalid domain type", func() {
			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: 42})
			Expect(cmd).To(BeNil())
		})

		It("sets selected domain from navigation", func() {
			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: DomainSystem})
			Expect(cmd).To(BeNil())
			Expect(intent.selectedDomain).To(Equal(DomainSystem))
		})
	})

	Describe("clearAllModals", func() {
		It("clears all modal references", func() {
			intent.openSettingsModal()
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			intent.resultModal = feedback.NewSuccessModal("Done")
			intent.clearAllModals()
			Expect(intent.settingsModal).To(BeNil())
			Expect(intent.savingModal).To(BeNil())
			Expect(intent.resultModal).To(BeNil())
		})
	})

	Describe("setCancelled", func() {
		It("deactivates intent and clears config result", func() {
			intent.configResult = &SystemResult{Success: true}
			intent.setCancelled()
			Expect(intent.active).To(BeFalse())
			Expect(intent.configResult).To(BeNil())
		})
	})
})
