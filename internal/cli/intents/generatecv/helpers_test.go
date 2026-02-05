package generatecv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Helpers", func() {
	var (
		intent *generatecv.Intent
	)

	BeforeEach(func() {
		ctx := validContext()
		var err error
		intent, err = generatecv.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("initWizardFlow", func() {
		Context("via Init", func() {
			It("should return a non-nil command", func() {
				freshIntent, err := generatecv.NewIntent(validContext())
				Expect(err).NotTo(HaveOccurred())

				cmd := freshIntent.Init()
				Expect(cmd).NotTo(BeNil())
			})

			It("should render wizard content in view after init", func() {
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("updateWizardFlow", func() {
		Context("window resize handling", func() {
			It("should handle window size messages without error", func() {
				cmd := intent.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				Expect(cmd).To(BeNil())
			})
		})

		Context("quit key handling", func() {
			It("should cancel intent on ctrl+c", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
			})

			It("should cancel intent on q key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
			})
		})

		Context("help toggle", func() {
			It("should handle help key without crashing", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("inactive intent", func() {
			It("should not process updates when inactive", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("technologies extracted message", func() {
			It("should transition to generating state", func() {
				msg := generatecv.TechnologiesExtractedMsg{
					Technologies: []*generatecv.ExtractedTechnology{},
					Suggestion:   nil,
					Error:        nil,
				}

				cmd := intent.Update(msg)
				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("CV generation complete message", func() {
			It("should transition to review state", func() {
				msg := generatecv.CVGenerationCompleteMsg{
					CV:    fixtures.CVView("cv-1"),
					Error: nil,
				}

				cmd := intent.Update(msg)
				Expect(cmd).To(BeNil())
			})
		})

		Context("export complete message", func() {
			It("should transition to export complete state on success", func() {
				msg := generatecv.ExportCompleteMsg{
					Path:  "/tmp/cv.md",
					Error: nil,
				}

				cmd := intent.Update(msg)
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("wizardView", func() {
		Context("in configuring state", func() {
			It("should render non-empty view", func() {
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("when inactive", func() {
			It("should show inactive message", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

				view := intent.View()
				Expect(view).To(ContainSubstring("not active"))
			})
		})
	})

	Describe("generateCVAsync", func() {
		Context("without generation service", func() {
			It("should return fallback CV via wizard complete then tech extracted", func() {
				wizardMsg := generatecv.WizardCompleteMsg{
					ProfileID:    "p1",
					Audience:     "hiring_manager",
					TechFocus:    "all",
					Technologies: []string{},
					FocusArea:    "backend",
					SkillsFormat: "grouped",
					SkillsLimit:  10,
					CVLength:     "standard",
				}

				cmd := intent.Update(wizardMsg)
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				techMsg, ok := resultMsg.(generatecv.TechnologiesExtractedMsg)
				Expect(ok).To(BeTrue())
				Expect(techMsg.Error).ToNot(HaveOccurred())

				cmd = intent.Update(techMsg)
				Expect(cmd).NotTo(BeNil())

				resultMsg = cmd()
				cvMsg, ok := resultMsg.(generatecv.CVGenerationCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(cvMsg.Error).ToNot(HaveOccurred())
				Expect(cvMsg.CV).NotTo(BeNil())
			})
		})
	})

	Describe("extractTechnologiesAsync", func() {
		Context("without repositories", func() {
			It("should return empty technologies via wizard complete", func() {
				wizardMsg := generatecv.WizardCompleteMsg{
					ProfileID:    "p1",
					Audience:     "hiring_manager",
					TechFocus:    "all",
					Technologies: []string{},
					FocusArea:    "backend",
					SkillsFormat: "grouped",
					SkillsLimit:  10,
					CVLength:     "standard",
				}

				cmd := intent.Update(wizardMsg)
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				techMsg, ok := resultMsg.(generatecv.TechnologiesExtractedMsg)
				Expect(ok).To(BeTrue())
				Expect(techMsg.Error).ToNot(HaveOccurred())
				Expect(techMsg.Technologies).To(BeEmpty())
				Expect(techMsg.Suggestion).NotTo(BeNil())
			})
		})
	})

	Describe("export complete handling", func() {
		Context("after successful export via message", func() {
			It("should not complete the intent yet", func() {
				exportMsg := generatecv.ExportCompleteMsg{
					Path:  "/tmp/cv.md",
					Error: nil,
				}

				cmd := intent.Update(exportMsg)
				Expect(cmd).To(BeNil())

				Expect(intent.Result()).To(BeNil())
			})
		})

		Context("after failed export via message", func() {
			It("should not complete the intent", func() {
				exportMsg := generatecv.ExportCompleteMsg{
					Path:  "",
					Error: errors.New("export failed"),
				}

				cmd := intent.Update(exportMsg)
				Expect(cmd).To(BeNil())

				Expect(intent.Result()).To(BeNil())
			})
		})
	})
})
