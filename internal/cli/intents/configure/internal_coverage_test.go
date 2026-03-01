package configure

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	configscreens "github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/config"
)

var _ = Describe("Configure intent internals", func() {
	var (
		cfg    *config.Config
		intent *Intent
	)

	BeforeEach(func() {
		cfg = config.DefaultConfig()
		config.SetConfigPathForTesting(filepath.Join(GinkgoT().TempDir(), "config.yaml"))
		ctx := &IntentContext{Cfg: cfg, Settings: settingsFromConfig(cfg)}
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

	Describe("screen result dispatch", func() {
		It("handles nil result", func() {
			cmd := intent.handleScreenResult(nil)
			Expect(cmd).To(BeNil())
		})

		It("handles non-ScreenResult", func() {
			cmd := intent.handleScreenResult("not a screen result")
			Expect(cmd).To(BeNil())
		})

		It("dispatches cancel result", func() {
			result := &screens.CancelResult{}
			cmd := intent.handleScreenResult(result)
			Expect(cmd).To(BeNil())
			Expect(intent.active).To(BeFalse())
		})

		It("dispatches navigate result", func() {
			result := &screens.NavigateResult{ResultData: DomainSystem}
			cmd := intent.handleScreenResult(result)
			Expect(cmd).NotTo(BeNil())
			Expect(intent.state).To(Equal(ConfigStateEditSettings))
		})
	})

	Describe("domain screen update", func() {
		It("updates domain screen and dispatches cancel", func() {
			cmd := intent.updateDomainScreen(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
			Expect(intent.active).To(BeFalse())
		})

		It("updates domain screen with passthrough", func() {
			cmd := intent.updateDomainScreen(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("modal update flows", func() {
		BeforeEach(func() {
			intent.settings = settingsFromConfig(cfg)
			intent.selectedDomain = DomainSystem
		})

		It("edit modal passthrough", func() {
			intent.openEditModal()
			cmd := intent.updateEditModal(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.editModal).NotTo(BeNil())
		})

		It("edit modal submit triggers review", func() {
			intent.openEditModal()
			cmd := intent.updateEditModal(tea.KeyMsg{Type: tea.KeyCtrlS})
			_ = cmd
		})

		It("edit modal cancel returns to select domain", func() {
			intent.openEditModal()
			intent.updateEditModal(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state).To(Equal(ConfigStateSelectDomain))
			Expect(intent.editModal).To(BeNil())
		})

		It("review modal accept triggers confirm", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			intent.updateReviewModal(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state).To(Equal(ConfigStateConfirm))
		})

		It("review modal cancel returns to edit", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			intent.updateReviewModal(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state).To(Equal(ConfigStateEditSettings))
		})

		It("review modal passthrough", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			cmd := intent.updateReviewModal(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.reviewModal).NotTo(BeNil())
		})

		It("confirm modal accept triggers saving", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openConfirmModal()
			intent.updateConfirmModal(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state).To(Equal(ConfigStateSaving))
		})

		It("confirm modal cancel returns to review", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openConfirmModal()
			intent.updateConfirmModal(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state).To(Equal(ConfigStateReviewChanges))
		})

		It("confirm modal passthrough", func() {
			intent.openConfirmModal()
			cmd := intent.updateConfirmModal(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.confirmModal).NotTo(BeNil())
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
		It("renders domain content", func() {
			intent.domainScreen = configscreens.NewDomainSelectScreen([]ConfigurationDomain{DomainSystem})
			Expect(intent.RenderDomainContent()).NotTo(BeEmpty())
		})

		It("reports all state names", func() {
			intent.clearAllModals()
			Expect(intent.getStateName()).To(Equal("Select Domain"))

			intent.openEditModal()
			Expect(intent.getStateName()).To(Equal("Edit Settings"))

			intent.editModal = nil
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			Expect(intent.getStateName()).To(Equal("Review Changes"))

			intent.reviewModal = nil
			intent.openConfirmModal()
			Expect(intent.getStateName()).To(Equal("Confirm"))

			intent.confirmModal = nil
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

			intent.openEditModal()
			Expect(intent.getContextHelp()).To(Equal(""))

			intent.editModal = nil
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			Expect(intent.getContextHelp()).To(Equal(""))

			intent.reviewModal = nil
			intent.openConfirmModal()
			Expect(intent.getContextHelp()).To(Equal(""))

			intent.confirmModal = nil
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
		It("renders edit modal adapter", func() {
			intent.selectedDomain = DomainSystem
			intent.openEditModal()
			adapter := editModalAdapter{intent.editModal, 80, 24}
			Expect(adapter.Render(0, 0)).NotTo(BeEmpty())
		})

		It("renders review modal adapter", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			adapter := reviewModalAdapter{intent.reviewModal, 80, 24}
			Expect(adapter.Render(0, 0)).NotTo(BeEmpty())
		})

		It("renders confirm modal adapter", func() {
			intent.openConfirmModal()
			adapter := confirmModalAdapter{intent.confirmModal, 80, 24}
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

		It("routes to confirm modal", func() {
			intent.openConfirmModal()
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyEnter})
			_ = cmd
		})

		It("routes to review modal", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyEnter})
			_ = cmd
			Expect(intent.state).To(Equal(ConfigStateConfirm))
		})

		It("routes to edit modal", func() {
			intent.openEditModal()
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
		})

		It("routes to domain screen", func() {
			cmd := intent.routeToActiveComponent(tea.KeyMsg{Type: tea.KeyDown})
			_ = cmd
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("window size propagation", func() {
		It("propagates to all modals", func() {
			intent.selectedDomain = DomainSystem
			intent.openEditModal()
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			intent.openConfirmModal()
			intent.handleWindowSizeMsg(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(intent.editModal).NotTo(BeNil())
			Expect(intent.reviewModal).NotTo(BeNil())
			Expect(intent.confirmModal).NotTo(BeNil())
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

		It("handles q key", func() {
			result := intent.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(result).To(BeTrue())
			Expect(intent.active).To(BeFalse())
		})

		It("suppresses q when form modals open", func() {
			intent.openEditModal()
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

		It("processes async completion", func() {
			result := &SystemResult{Success: true}
			cmd := intent.Update(ConfigCompleteMsg{Result: result})
			_ = cmd
			Expect(intent.state).To(Equal(ConfigStateComplete))
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

		It("renders with edit modal overlay", func() {
			intent.selectedDomain = DomainSystem
			intent.openEditModal()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with review modal overlay", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with confirm modal overlay", func() {
			intent.openConfirmModal()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with saving modal overlay", func() {
			intent.savingModal = feedback.NewLoadingModal("Saving...", false)
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with result modal overlay", func() {
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
		It("sets edit with domain", func() {
			intent.selectedDomain = DomainSystem
			intent.SetState(ConfigStateEditSettings)
			Expect(intent.editModal).NotTo(BeNil())
		})

		It("sets edit without domain", func() {
			intent.selectedDomain = ""
			intent.SetState(ConfigStateEditSettings)
			Expect(intent.editModal).To(BeNil())
		})

		It("sets review with pending changes", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.SetState(ConfigStateReviewChanges)
			Expect(intent.reviewModal).NotTo(BeNil())
		})

		It("sets review without pending changes", func() {
			intent.pendingChanges = nil
			intent.SetState(ConfigStateReviewChanges)
			Expect(intent.reviewModal).To(BeNil())
		})

		It("sets confirm", func() {
			intent.SetState(ConfigStateConfirm)
			Expect(intent.confirmModal).NotTo(BeNil())
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

		It("sets select domain", func() {
			intent.SetState(ConfigStateSelectDomain)
			Expect(intent.state).To(Equal(ConfigStateSelectDomain))
		})
	})

	Describe("result modal dismiss", func() {
		It("dismisses success on space", func() {
			intent.configResult = &SystemResult{Success: true}
			intent.resultModal = feedback.NewSuccessModal("Done")
			intent.updateResultModal(tea.KeyMsg{Type: tea.KeySpace})
			Expect(intent.active).To(BeFalse())
		})

		It("dismisses failed and returns to edit", func() {
			intent.configResult = &SystemResult{Success: false}
			intent.resultModal = feedback.NewErrorModal("Failed", "Error")
			intent.selectedDomain = DomainSystem
			intent.settings = settingsFromConfig(cfg)
			cmd := intent.updateResultModal(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.state).To(Equal(ConfigStateEditSettings))
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

		It("dismisses nil result as failed", func() {
			intent.configResult = nil
			intent.resultModal = feedback.NewErrorModal("Failed", "Error")
			intent.selectedDomain = DomainSystem
			intent.settings = settingsFromConfig(cfg)
			cmd := intent.updateResultModal(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.state).To(Equal(ConfigStateEditSettings))
		})
	})

	Describe("async completion", func() {
		It("handles complete", func() {
			result := &SystemResult{Success: true}
			cmd := intent.handleAsyncCompletion(ConfigCompleteMsg{Result: result})
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(ConfigStateComplete))
		})

		It("handles error", func() {
			errMsg := ConfigErrorMsg{Error: &intents.IntentError{Code: "err", Message: "oops"}}
			cmd := intent.handleAsyncCompletion(errMsg)
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(ConfigStateFailed))
		})

		It("ignores unrelated", func() {
			cmd := intent.handleAsyncCompletion(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("navigate handler", func() {
		It("handles invalid domain type", func() {
			cmd := intent.HandleNavigate(&screens.NavigateResult{ResultData: 42})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("hasFormModals", func() {
		It("returns false when no modals", func() {
			Expect(intent.hasFormModals()).To(BeFalse())
		})

		It("returns true with edit modal", func() {
			intent.openEditModal()
			Expect(intent.hasFormModals()).To(BeTrue())
		})

		It("returns true with review modal", func() {
			intent.pendingChanges = map[string]interface{}{"log_level": "warn"}
			intent.openReviewModal()
			Expect(intent.hasFormModals()).To(BeTrue())
		})

		It("returns true with confirm modal", func() {
			intent.openConfirmModal()
			Expect(intent.hasFormModals()).To(BeTrue())
		})
	})
})
