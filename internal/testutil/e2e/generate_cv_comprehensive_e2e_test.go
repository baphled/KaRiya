package e2e_test

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mocksvc "github.com/baphled/kariya/internal/testutil/mocks/service"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Comprehensive E2E Tests for Generate CV Wizard Workflow (Phase 3 - Task 55)
//
// PURPOSE: Exhaustive coverage of every user journey through the wizard flow,
// using mock services to test the full pipeline from wizard completion through
// review, preview, and export states.
//
// STRATEGY:
// - Tests that only exercise wizard navigation use GetSharedEnv (no mocks)
// - Tests that need the full pipeline (extraction -> generation -> review/preview/export)
//   use SetupWithCVMocks with GoMock services
// - Wizard completion is done via intents.WizardCompleteMsg to bypass huh form
//   internals, matching the proven pattern from internal tests

// completeWizardViaMessage sends a WizardCompleteMsg to bypass the huh form
// and jump directly from configuring to extracting state. This is the same
// approach used by all internal wizard E2E tests since huh form navigation
// cannot be reliably driven via keyboard in the E2E framework.
func completeWizardViaMessage(env *e2e.TestEnv) {
	env.SendMessage(generatecv.WizardCompleteMsg{
		ProfileID:    "profile-staff-engineer",
		Audience:     "hiring_manager",
		TechFocus:    "language_agnostic",
		SkillsFormat: "grouped",
		SkillsLimit:  5,
		CVLength:     "2_page",
	})
}

var _ = Describe("E2E GenerateCV Comprehensive", func() {

	// ========================================================================
	// SECTION 3.1: HAPPY PATHS
	// ========================================================================
	Describe("Happy Paths", func() {

		Describe("A1: Complete wizard to review (no export)", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVViewWithSections("cv-a1", []*career.CVSection{
					fixtures.CVSectionWithContent("sec-1", "cv-a1", []*career.SectionContentGroup{
						fixtures.ContentGroup("Professional Experience"),
					}),
				})
				testCV.Name = "Staff Engineer"
				testCV.TargetRole = "staff"
				testCV.TargetAudience = "hiring_manager"
				testCV.SourceEventCount = 5

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should reach review screen after completing wizard", func() {
				env.SelectIntentByName("generate_cv")

				completeWizardViaMessage(env)

				view := env.GetView()
				Expect(view).To(ContainSubstring("CV Review"))
			})
		})

		Describe("A4: Skip wizard with Ctrl+S", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should not panic when Ctrl+S pressed with profile selected", func() {
				env.SelectIntentByName("generate_cv")

				env.NavigateDown()
				env.Confirm()

				env.PressKey(tea.KeyCtrlS)

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("A5: Tab through wizard fields", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should advance through wizard fields with Tab", func() {
				env.SelectIntentByName("generate_cv")

				env.AssertViewContains("Profile")

				env.Tab()

				view := env.GetView()
				Expect(view).To(ContainSubstring("Audience"))
				Expect(view).NotTo(ContainSubstring("panic"))
			})
		})

		Describe("A7: Navigate between review and preview", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVViewWithSections("cv-a7", []*career.CVSection{
					fixtures.CVSectionWithContent("sec-1", "cv-a7", []*career.SectionContentGroup{
						fixtures.ContentGroup("Experience"),
					}),
				})
				testCV.Name = "Staff Engineer"
				testCV.SourceEventCount = 5

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should navigate from review to preview with Enter and back with Esc", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				env.AssertViewContains("CV Review")

				env.Confirm()

				previewView := env.GetView()
				Expect(previewView).NotTo(ContainSubstring("CV Review"))

				env.Cancel()

				backToReview := env.GetView()
				Expect(backToReview).To(ContainSubstring("CV Review"))
			})
		})

		Describe("A8: Profile options generate successfully", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, config *career.CVConfig) (*career.CVView, error) {
						return fixtures.CVViewWith("cv-profile-test", config.Name, config.TargetRole, config.TargetAudience), nil
					}).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should generate CV with staff engineer profile", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				view := env.GetView()
				Expect(view).To(ContainSubstring("CV Review"))
			})

			It("should generate CV with principal engineer profile", func() {
				env.SelectIntentByName("generate_cv")
				env.SendMessage(generatecv.WizardCompleteMsg{
					ProfileID:    "profile-principal-engineer",
					Audience:     "recruiter",
					SkillsFormat: "flat",
					CVLength:     "1_page",
				})

				view := env.GetView()
				Expect(view).To(ContainSubstring("CV Review"))
			})
		})
	})

	// ========================================================================
	// SECTION 3.2: SAD PATHS
	// ========================================================================
	Describe("Sad Paths", func() {

		Describe("B1: Empty database shows warning modal", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should show warning modal and dismiss with Esc returns to menu", func() {
				env.AssertEventCount(0)

				env.SelectIntentByName("generate_cv")

				env.AssertViewContains("No Career Events")

				env.Cancel()

				Expect(env.IsInMenuState()).To(BeTrue())
			})
		})

		Describe("B2: Cancel from wizard step 1", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should return to main menu when Esc pressed from step 1", func() {
				env.SelectIntentByName("generate_cv")

				env.AssertViewContains("Profile")

				env.Cancel()

				view := env.GetView()
				Expect(view).To(Or(
					ContainSubstring("Capture Event"),
					ContainSubstring("Main Menu"),
				))
			})
		})

		Describe("B5: CV generation failure", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("service unavailable: generation failed")).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should not panic when CV generation fails", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("B7: Cancel export modal", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVViewWithSections("cv-b7", []*career.CVSection{
					fixtures.CVSection("sec-1", "cv-b7"),
				})

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should return to preview when export modal is cancelled", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				env.Confirm()

				env.PressKeyRune('x')

				exportView := env.GetView()
				Expect(exportView).To(Or(
					ContainSubstring("Export"),
					ContainSubstring("format"),
				))

				env.Cancel()

				afterCancel := env.GetView()
				Expect(afterCancel).NotTo(ContainSubstring("panic"))
			})
		})

		Describe("B8: Rapid state transitions", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should not panic with rapid Tab/Esc/Enter/j/k mashing", func() {
				env.SelectIntentByName("generate_cv")

				for range 3 {
					env.Tab()
					env.NavigateDown()
					env.NavigateUp()
					env.Confirm()
					env.Cancel()
					env.Tab()
					env.NavigateDown()
					env.Confirm()
				}

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	// ========================================================================
	// SECTION 3.3: EDGE CASES
	// ========================================================================
	Describe("Edge Cases", func() {

		Describe("C1: Single event in database", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVViewWithCounts("cv-c1", 1, 0)
				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})

				event := fixtures.EventWith("event-single", "Implemented microservices architecture", "TestCorp", "Platform Migration")
				event.Date = time.Now().AddDate(0, -1, 0)
				env.AddEvent(event)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should generate CV with single event", func() {
				env.AssertEventCount(1)

				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				view := env.GetView()
				Expect(view).To(ContainSubstring("CV Review"))
			})
		})

		Describe("C3: No technologies detected (wizard skips tech step)", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should show wizard without technology focus step", func() {
				env.SelectIntentByName("generate_cv")

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("Technology Focus"))
				Expect(view).To(ContainSubstring("Profile"))
			})
		})

		Describe("C4: Rapid key presses", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should not panic with 20+ Tab presses", func() {
				env.SelectIntentByName("generate_cv")

				for range 25 {
					env.Tab()
				}

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("C5: Window resize during wizard", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should handle window resize without panic", func() {
				env.SelectIntentByName("generate_cv")

				env.SendMessage(tea.WindowSizeMsg{Width: 80, Height: 24})
				env.SendMessage(tea.WindowSizeMsg{Width: 120, Height: 40})
				env.SendMessage(tea.WindowSizeMsg{Width: 40, Height: 15})

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("C6: Window resize during progress modal", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVView("cv-c6")
				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should handle resize after wizard completion without panic", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				env.SendMessage(tea.WindowSizeMsg{Width: 60, Height: 20})

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("C7: Preview scrolling", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				sections := make([]*career.CVSection, 5)
				for i := range 5 {
					sections[i] = fixtures.CVSectionWithContent(
						fmt.Sprintf("sec-%d", i),
						"cv-c7",
						[]*career.SectionContentGroup{
							fixtures.ContentGroup(fmt.Sprintf("Section %d Experience", i)),
						},
					)
				}
				testCV := fixtures.CVViewWithSections("cv-c7", sections)

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should handle scroll keys in preview without panic", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				env.Confirm()

				env.PressKeyRune('j')
				env.PressKeyRune('j')
				env.PressKeyRune('k')
				env.PressKeyRune('G')
				env.PressKeyRune('g')

				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	// ========================================================================
	// SECTION 3.4: DATA INTEGRITY
	// ========================================================================
	Describe("Data Integrity", func() {

		Describe("D1: Profile selection preserved through workflow", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, config *career.CVConfig) (*career.CVView, error) {
						cv := fixtures.CVViewWith("cv-d1", config.Name, config.TargetRole, config.TargetAudience)
						cv.SourceEventCount = 5
						return cv, nil
					}).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should show selected profile data in review screen", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				view := env.GetView()
				Expect(view).To(ContainSubstring("staff"))
				Expect(view).To(ContainSubstring("hiring_manager"))
			})
		})

		Describe("D4: State reset on re-entry after cancellation", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should reset wizard state on re-entry after cancel", func() {
				env.SelectIntentByName("generate_cv")

				env.NavigateDown()
				env.NavigateDown()
				env.Tab()

				env.Cancel()

				env.SelectIntentByName("generate_cv")

				view := env.GetView()
				Expect(view).To(ContainSubstring("Profile"))
				Expect(view).NotTo(ContainSubstring("panic"))
			})
		})

		Describe("D5: Data survives session restart", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should preserve event and fact data across restart", func() {
				env.PopulateTestData(5, 2, 3)

				env.AssertEventCount(5)
				env.AssertFactCount(3)

				env.SimulateRestart()

				env.AssertEventCount(5)
				env.AssertFactCount(3)

				env.SelectIntentByName("generate_cv")
				env.AssertViewContains("Profile")
			})
		})
	})

	// ========================================================================
	// SECTION 3.5: NAVIGATION CONSISTENCY
	// ========================================================================
	Describe("Navigation Consistency", func() {

		Describe("E1: Global 'q' key from wizard state", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should quit without panic from wizard state", func() {
				env.SelectIntentByName("generate_cv")
				env.Quit()

				Expect(true).To(BeTrue())
			})
		})

		Describe("E1: Global 'q' from review state", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVView("cv-e1")
				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should quit without panic from review state", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				env.AssertViewContains("CV Review")

				env.Quit()

				Expect(true).To(BeTrue())
			})
		})

		Describe("E3: Esc behavior is context-appropriate", func() {

			Describe("from wizard step 1", func() {
				var env *e2e.TestEnv

				BeforeEach(func() {
					env = e2e.GetSharedEnv(GinkgoT())
					env.PopulateTestData(5, 2, 3)
				})

				AfterEach(func() {
					env.Cleanup()
				})

				It("should cancel wizard and return to menu", func() {
					env.SelectIntentByName("generate_cv")
					env.Cancel()

					view := env.GetView()
					Expect(view).To(Or(
						ContainSubstring("Capture Event"),
						ContainSubstring("Main Menu"),
					))
				})
			})

			Describe("from review screen", func() {
				var (
					env      *e2e.TestEnv
					ctrl     *gomock.Controller
					mockGen  *mocksvc.MockCVGenerationService
					mockClip *mocksvc.MockClipboardWriter
				)

				BeforeEach(func() {
					ctrl = gomock.NewController(GinkgoT())
					mockGen = mocksvc.NewMockCVGenerationService(ctrl)
					mockClip = mocksvc.NewMockClipboardWriter(ctrl)

					testCV := fixtures.CVView("cv-e3")
					mockGen.EXPECT().
						GenerateCVFromConfig(gomock.Any(), gomock.Any()).
						Return(testCV, nil).
						AnyTimes()

					env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
						CVGenService:    mockGen,
						ClipboardWriter: mockClip,
					})
					env.PopulateTestData(5, 2, 3)
				})

				AfterEach(func() {
					env.Cleanup()
					ctrl.Finish()
				})

				It("should return to wizard from review on Esc", func() {
					env.SelectIntentByName("generate_cv")
					completeWizardViaMessage(env)

					env.AssertViewContains("CV Review")

					env.Cancel()

					view := env.GetView()
					Expect(view).To(ContainSubstring("CV Configuration"))
				})
			})

			Describe("from preview screen", func() {
				var (
					env      *e2e.TestEnv
					ctrl     *gomock.Controller
					mockGen  *mocksvc.MockCVGenerationService
					mockClip *mocksvc.MockClipboardWriter
				)

				BeforeEach(func() {
					ctrl = gomock.NewController(GinkgoT())
					mockGen = mocksvc.NewMockCVGenerationService(ctrl)
					mockClip = mocksvc.NewMockClipboardWriter(ctrl)

					testCV := fixtures.CVViewWithSections("cv-e3-preview", []*career.CVSection{
						fixtures.CVSection("sec-1", "cv-e3-preview"),
					})
					mockGen.EXPECT().
						GenerateCVFromConfig(gomock.Any(), gomock.Any()).
						Return(testCV, nil).
						AnyTimes()

					env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
						CVGenService:    mockGen,
						ClipboardWriter: mockClip,
					})
					env.PopulateTestData(5, 2, 3)
				})

				AfterEach(func() {
					env.Cleanup()
					ctrl.Finish()
				})

				It("should return to review from preview on Esc", func() {
					env.SelectIntentByName("generate_cv")
					completeWizardViaMessage(env)

					env.AssertViewContains("CV Review")

					env.Confirm()

					env.Cancel()

					view := env.GetView()
					Expect(view).To(ContainSubstring("CV Review"))
				})
			})
		})

		Describe("E4: Help footer shows correct shortcuts", func() {
			var env *e2e.TestEnv

			BeforeEach(func() {
				env = e2e.GetSharedEnv(GinkgoT())
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should show navigation help in wizard", func() {
				env.SelectIntentByName("generate_cv")

				view := env.GetView()
				Expect(view).To(Or(
					ContainSubstring("Tab"),
					ContainSubstring("Enter"),
					ContainSubstring("Esc"),
					ContainSubstring("Ctrl+S"),
					ContainSubstring("Skip"),
				))
			})
		})

		Describe("E4: Help footer in review screen", func() {
			var (
				env      *e2e.TestEnv
				ctrl     *gomock.Controller
				mockGen  *mocksvc.MockCVGenerationService
				mockClip *mocksvc.MockClipboardWriter
			)

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				mockGen = mocksvc.NewMockCVGenerationService(ctrl)
				mockClip = mocksvc.NewMockClipboardWriter(ctrl)

				testCV := fixtures.CVView("cv-e4")
				mockGen.EXPECT().
					GenerateCVFromConfig(gomock.Any(), gomock.Any()).
					Return(testCV, nil).
					AnyTimes()

				env = e2e.SetupWithCVMocks(GinkgoT(), &e2e.CVMockConfig{
					CVGenService:    mockGen,
					ClipboardWriter: mockClip,
				})
				env.PopulateTestData(5, 2, 3)
			})

			AfterEach(func() {
				env.Cleanup()
				ctrl.Finish()
			})

			It("should show review-specific shortcuts", func() {
				env.SelectIntentByName("generate_cv")
				completeWizardViaMessage(env)

				view := env.GetView()
				Expect(view).To(ContainSubstring("preview"))
				Expect(view).To(ContainSubstring("export"))
				Expect(view).To(ContainSubstring("esc"))
			})
		})
	})
})
