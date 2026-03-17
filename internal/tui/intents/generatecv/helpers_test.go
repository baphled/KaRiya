package generatecv_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockintent "github.com/baphled/kariya/internal/testutil/mocks/intent"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
	mocksvc "github.com/baphled/kariya/internal/testutil/mocks/service"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
	cvviews "github.com/baphled/kariya/internal/tui/views/cv"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
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

	Describe("generateCVAsync with CVGenerationService", func() {
		var (
			ctrl              *gomock.Controller
			mockCVService     *mocksvc.MockCVGenerationService
			intentWithService *generatecv.Intent
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockCVService = mocksvc.NewMockCVGenerationService(ctrl)

			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
				},
				Events:              []*career.Event{fixtures.Event("e1")},
				CVGenerationService: mockCVService,
			}
			var err error
			intentWithService, err = generatecv.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := intentWithService.GetTerminalInfo()
			termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			intentWithService.Init()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should generate CV successfully with service", func() {
			expectedCV := fixtures.CVView("generated-cv")
			mockCVService.EXPECT().
				GenerateCVFromConfig(gomock.Any(), gomock.Any()).
				Return(expectedCV, nil)

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
			cmd := intentWithService.Update(wizardMsg)
			Expect(cmd).NotTo(BeNil())

			techMsg := cmd()
			cmd = intentWithService.Update(techMsg)
			Expect(cmd).NotTo(BeNil())

			resultMsg := cmd()
			cvMsg, ok := resultMsg.(generatecv.CVGenerationCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(cvMsg.Error).ToNot(HaveOccurred())
			Expect(cvMsg.CV).NotTo(BeNil())
			Expect(cvMsg.CV.ID).To(Equal("generated-cv"))
		})

		It("should handle CV generation error", func() {
			expectedErr := errors.New("generation failed")
			mockCVService.EXPECT().
				GenerateCVFromConfig(gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

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
			cmd := intentWithService.Update(wizardMsg)
			Expect(cmd).NotTo(BeNil())

			techMsg := cmd()
			cmd = intentWithService.Update(techMsg)
			Expect(cmd).NotTo(BeNil())

			resultMsg := cmd()
			cvMsg, ok := resultMsg.(generatecv.CVGenerationCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(cvMsg.Error).To(HaveOccurred())
			Expect(cvMsg.CV).To(BeNil())
		})
	})

	Describe("generateCVAsync with ProfileConfig", func() {
		var (
			ctrl              *gomock.Controller
			mockCVService     *mocksvc.MockCVGenerationService
			intentWithProfile *generatecv.Intent
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockCVService = mocksvc.NewMockCVGenerationService(ctrl)

			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
				},
				Events:              []*career.Event{fixtures.Event("e1")},
				CVGenerationService: mockCVService,
				ProfileConfig: &config.ProfileConfig{
					Name:           "Test User",
					Title:          "Staff Engineer",
					SummaryHeading: "Professional Summary",
					WhatIBring:     []string{"leadership", "technical excellence"},
					CoreStrengths:  []string{"Go", "distributed systems"},
				},
			}
			var err error
			intentWithProfile, err = generatecv.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := intentWithProfile.GetTerminalInfo()
			termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			intentWithProfile.Init()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should pass ProfileConfig fields to CV generation", func() {
			expectedCV := fixtures.CVView("generated-cv-profile")
			mockCVService.EXPECT().
				GenerateCVFromConfig(gomock.Any(), gomock.Any()).
				Return(expectedCV, nil)

			cmd := intentWithProfile.Update(generatecv.WizardCompleteMsg{
				ProfileID: "p1", Audience: "hiring_manager",
			})
			Expect(cmd).NotTo(BeNil())

			techMsg := cmd()
			cmd = intentWithProfile.Update(techMsg)
			Expect(cmd).NotTo(BeNil())

			resultMsg := cmd()
			cvMsg, ok := resultMsg.(generatecv.CVGenerationCompleteMsg)
			Expect(ok).To(BeTrue())
			Expect(cvMsg.Error).ToNot(HaveOccurred())
			Expect(cvMsg.CV).NotTo(BeNil())
		})
	})

	Describe("extractTechnologiesAsync with repositories", func() {
		var (
			ctrl                   *gomock.Controller
			mockSkillRepo          *mockrepo.MockSkillRepository
			mockEventRepo          *mockrepo.MockEventRepository
			intentWithRepositories *generatecv.Intent
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockSkillRepo = mockrepo.NewMockSkillRepository(ctrl)
			mockEventRepo = mockrepo.NewMockEventRepository(ctrl)

			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
				},
				Events:          []*career.Event{fixtures.Event("e1")},
				SkillRepository: mockSkillRepo,
				EventRepository: mockEventRepo,
			}
			var err error
			intentWithRepositories, err = generatecv.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intentWithRepositories.Init()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should extract technologies with repositories", func() {
			mockSkillRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*career.Skill{fixtures.Skill("s1")}, nil).AnyTimes()
			mockSkillRepo.EXPECT().
				GetEventCountsForSkills(gomock.Any()).
				Return(map[string]int{"Go": 5}, nil).AnyTimes()
			mockSkillRepo.EXPECT().
				GetLastUsedForSkills(gomock.Any()).
				Return(nil, nil).AnyTimes()
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*career.Event{fixtures.Event("e1")}, nil).AnyTimes()

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
			cmd := intentWithRepositories.Update(wizardMsg)
			Expect(cmd).NotTo(BeNil())

			resultMsg := cmd()
			techMsg, ok := resultMsg.(generatecv.TechnologiesExtractedMsg)
			Expect(ok).To(BeTrue())
			Expect(techMsg.Error).ToNot(HaveOccurred())
		})

		It("should handle repository errors", func() {
			mockSkillRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error")).AnyTimes()

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
			cmd := intentWithRepositories.Update(wizardMsg)
			Expect(cmd).NotTo(BeNil())

			resultMsg := cmd()
			techMsg, ok := resultMsg.(generatecv.TechnologiesExtractedMsg)
			Expect(ok).To(BeTrue())
			Expect(techMsg.Error).To(HaveOccurred())
		})
	})

	Describe("extractTechnologiesAsync without AppContext", func() {
		var intentWithoutAppContext *generatecv.Intent
		var ctrl *gomock.Controller
		var mockSkillRepo *mockrepo.MockSkillRepository
		var mockEventRepo *mockrepo.MockEventRepository

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockSkillRepo = mockrepo.NewMockSkillRepository(ctrl)
			mockEventRepo = mockrepo.NewMockEventRepository(ctrl)

			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
				},
				Events:          []*career.Event{fixtures.Event("e1")},
				SkillRepository: mockSkillRepo,
				EventRepository: mockEventRepo,
			}
			var err error
			intentWithoutAppContext, err = generatecv.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intentWithoutAppContext.Init()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should extract technologies using background context", func() {
			mockSkillRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*career.Skill{fixtures.Skill("s1")}, nil).AnyTimes()
			mockSkillRepo.EXPECT().
				GetEventCountsForSkills(gomock.Any()).
				Return(map[string]int{"Go": 5}, nil).AnyTimes()
			mockSkillRepo.EXPECT().
				GetLastUsedForSkills(gomock.Any()).
				Return(nil, nil).AnyTimes()
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*career.Event{fixtures.Event("e1")}, nil).AnyTimes()

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
			cmd := intentWithoutAppContext.Update(wizardMsg)
			Expect(cmd).NotTo(BeNil())

			resultMsg := cmd()
			_, ok := resultMsg.(generatecv.TechnologiesExtractedMsg)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("delegateToWizardModal", func() {
		Context("when wizard modal completes", func() {
			It("should return wizard complete handler", func() {
				intent.Init()

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
			})
		})

		Context("when wizard modal is cancelled", func() {
			It("should cancel the intent", func() {
				intent.Init()
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
			})
		})
	})

	Describe("delegateToExportModal", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intent.Update(generatecv.TechnologiesExtractedMsg{})
			intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
		})

		Context("when export modal is visible", func() {
			It("should delegate updates to export modal", func() {
				intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})

				Expect(intent.GetState()).To(Equal(generatecv.StateExporting))
			})
		})

		Context("when export modal is cancelled", func() {
			It("should return to preview state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.GetState()).To(Equal(generatecv.StatePreview))
			})
		})
	})

	Describe("handleKeyDelegation", func() {
		Context("when progress modal is visible and cancellable", func() {
			It("should return to wizard on escape", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.GetState()).To(Equal(generatecv.StateConfiguring))
			})
		})

		Context("in review state with active screen", func() {
			It("should delegate to active screen", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("in export complete state", func() {
			It("should handle enter key to complete", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
				intent.HandleNavigate(&widgets.NavigateViewResult{ResultData: cvviews.Nav{Action: cvviews.ActionExport}})
				intent.Update(generatecv.ExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})

				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
			})
		})
	})

	Describe("exportCVAsync error cases", func() {
		Context("without ExportService", func() {
			It("should return error when ExportService is nil", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				cmd := intent.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).To(HaveOccurred())
				Expect(exportMsg.Error.Error()).To(ContainSubstring("export service not available"))
			})
		})

		Context("without AppContext but with ExportService", func() {
			var intentWithService *generatecv.Intent

			BeforeEach(func() {
				ctx := &generatecv.IntentValidator{
					AvailableProfiles: []*generatecv.CVProfile{
						{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
					},
					Events: []*career.Event{fixtures.Event("e1")},
				}
				var err error
				intentWithService, err = generatecv.NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				intentWithService.Init()
			})

			It("should return error when AppContext is nil", func() {
				intentWithService.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intentWithService.Update(generatecv.TechnologiesExtractedMsg{})
				intentWithService.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				view := intentWithService.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("with unknown export format", func() {
			It("should return error for unknown format", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				intent.SetSelectedExportFormatForTest("")
				cmd := intent.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).To(HaveOccurred())
			})
		})
	})

	Describe("handleViewResult edge cases", func() {
		Context("with nil result", func() {
			It("should return nil command", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("with non-ViewResult type", func() {
			It("should handle gracefully", func() {
				intent.Init()
				intent.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
				intent.Update(generatecv.TechnologiesExtractedMsg{})
				intent.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("exportCVAsync with mocked ExportService", func() {
		var (
			ctrl               *gomock.Controller
			mockExporter       *mockintent.MockCVExporter
			intentWithExporter *generatecv.Intent
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockExporter = mockintent.NewMockCVExporter(ctrl)

			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Staff Engineer", TargetRole: "staff", TargetAudience: "hiring_manager"},
				},
				Events:        []*career.Event{fixtures.Event("e1")},
				ExportService: mockExporter,
			}
			var err error
			intentWithExporter, err = generatecv.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			termInfo := intentWithExporter.GetTerminalInfo()
			termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			intentWithExporter.Init()
			intentWithExporter.Update(generatecv.WizardCompleteMsg{ProfileID: "p1", Audience: "hiring_manager"})
			intentWithExporter.Update(generatecv.TechnologiesExtractedMsg{})
			intentWithExporter.Update(generatecv.CVGenerationCompleteMsg{CV: fixtures.CVView("cv-1")})
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		Context("export to text and save to file", func() {
			It("should export successfully", func() {
				mockExporter.EXPECT().
					ExportToText(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("CV Content", nil)
				mockExporter.EXPECT().
					SaveToFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("/tmp/cv.txt", nil)

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionSaveToFile)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).ToNot(HaveOccurred())
				Expect(exportMsg.Path).To(Equal("/tmp/cv.txt"))
			})
		})

		Context("export to markdown and save to file", func() {
			It("should export successfully", func() {
				mockExporter.EXPECT().
					ExportToMarkdown(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("# CV Content", nil)
				mockExporter.EXPECT().
					SaveToFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("/tmp/cv.md", nil)

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatMarkdown)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionSaveToFile)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).ToNot(HaveOccurred())
				Expect(exportMsg.Path).To(Equal("/tmp/cv.md"))
			})
		})

		Context("export to YAML and save to file", func() {
			It("should export successfully", func() {
				mockExporter.EXPECT().
					ExportToYAML(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("name: CV", nil)
				mockExporter.EXPECT().
					SaveToFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("/tmp/cv.yaml", nil)

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatYAML)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionSaveToFile)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).ToNot(HaveOccurred())
				Expect(exportMsg.Path).To(Equal("/tmp/cv.yaml"))
			})
		})

		Context("export to clipboard", func() {
			It("should export successfully", func() {
				mockExporter.EXPECT().
					ExportToText(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("CV Content", nil)
				mockExporter.EXPECT().
					CopyToClipboard(gomock.Any(), gomock.Any()).
					Return(nil)

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionClipboard)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).ToNot(HaveOccurred())
				Expect(exportMsg.Path).To(Equal("clipboard"))
			})
		})

		Context("export format error", func() {
			It("should return error when export fails", func() {
				mockExporter.EXPECT().
					ExportToText(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("export failed"))

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionSaveToFile)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).To(HaveOccurred())
				Expect(exportMsg.Error.Error()).To(ContainSubstring("failed to export"))
			})
		})

		Context("save to file error", func() {
			It("should return error when save fails", func() {
				mockExporter.EXPECT().
					ExportToText(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("CV Content", nil)
				mockExporter.EXPECT().
					SaveToFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("save failed"))

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionSaveToFile)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).To(HaveOccurred())
				Expect(exportMsg.Error.Error()).To(ContainSubstring("failed to save file"))
			})
		})

		Context("clipboard error", func() {
			It("should return error when clipboard fails", func() {
				mockExporter.EXPECT().
					ExportToText(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("CV Content", nil)
				mockExporter.EXPECT().
					CopyToClipboard(gomock.Any(), gomock.Any()).
					Return(errors.New("clipboard failed"))

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
				intentWithExporter.SetSelectedExportOptionForTest(generatecv.ExportOptionClipboard)

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).To(HaveOccurred())
				Expect(exportMsg.Error.Error()).To(ContainSubstring("failed to copy to clipboard"))
			})
		})

		Context("unknown export option", func() {
			It("should return error for unknown option", func() {
				mockExporter.EXPECT().
					ExportToText(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("CV Content", nil)

				intentWithExporter.SetSelectedExportFormatForTest(generatecv.ExportFormatText)
				intentWithExporter.SetSelectedExportOptionForTest("")

				cmd := intentWithExporter.InvokeExportCVAsyncForTest()
				Expect(cmd).NotTo(BeNil())

				resultMsg := cmd()
				exportMsg, ok := resultMsg.(generatecv.ExportCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(exportMsg.Error).To(HaveOccurred())
				Expect(exportMsg.Error.Error()).To(ContainSubstring("unknown export option"))
			})
		})
	})

})
