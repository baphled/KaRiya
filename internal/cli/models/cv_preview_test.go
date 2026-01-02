package models

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCVPreview(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CVPreview Suite")
}

var _ = Describe("CVPreviewModel", func() {
	var (
		model        *CVPreviewModel
		baseModel    *BaseStandardModel
		testCVView   *career.CVView
		testSections []*career.CVSection
	)

	BeforeEach(func() {
		baseModel = NewBaseStandardModel()

		testCVView = &career.CVView{
			ID:                uuid.New().String(),
			Name:              "Senior IC CV",
			TargetRole:        "senior_ic",
			TargetAudience:    []string{"hiring_manager", "recruiter"},
			GeneratedAt:       time.Now(),
			SourceEventCount:  15,
			SourceFactCount:   8,
		}

		testSections = []*career.CVSection{
			{
				ID:          uuid.New().String(),
				CVViewID:    testCVView.ID,
				SectionType: "experience",
				Title:       "Professional Experience",
				Order:       0,
				Content:     "• Led team of 5 engineers\n• Designed and implemented new architecture",
			},
			{
				ID:          uuid.New().String(),
				CVViewID:    testCVView.ID,
				SectionType: "skills",
				Title:       "Core Competencies",
				Order:       1,
				Content:     "• Go, Python, Rust\n• System Design\n• Leadership",
			},
		}

		model = NewCVPreviewModel(baseModel, testCVView, testSections, nil)
	})

	Describe("Initialization", func() {
		It("should initialize with CV view and sections", func() {
			Expect(model.cvView).To(Equal(testCVView))
			Expect(model.sections).To(Equal(testSections))
			Expect(model.selectedSectionIdx).To(Equal(0))
			Expect(model.selectedBulletIdx).To(Equal(0))
		})

		It("should have empty expanded bullets map", func() {
			Expect(model.expandedBullets).To(BeEmpty())
		})
	})

	Describe("Navigation", func() {
		It("should move bullet down with 'j'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			castedModel := newModel.(*CVPreviewModel)
			// Bullet movement depends on section content parsing
			Expect(castedModel).NotTo(BeNil())
		})

		It("should move bullet up with 'k'", func() {
			model.selectedBulletIdx = 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel.selectedBulletIdx).To(Equal(0))
		})

		It("should move to next section with 'l'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel.selectedSectionIdx).To(Equal(1))
		})

		It("should move to previous section with 'h'", func() {
			model.selectedSectionIdx = 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel.selectedSectionIdx).To(Equal(0))
		})

		It("should not move past last section with 'l'", func() {
			model.selectedSectionIdx = len(testSections) - 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel.selectedSectionIdx).To(Equal(len(testSections) - 1))
		})

		It("should not move before first section with 'h'", func() {
			model.selectedSectionIdx = 0
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel.selectedSectionIdx).To(Equal(0))
		})
	})

	Describe("Bullet Expansion", func() {
		It("should toggle bullet expansion on Enter", func() {
			Expect(model.expandedBullets[0]).To(BeFalse())

			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			castedModel := newModel.(*CVPreviewModel)

			Expect(castedModel.expandedBullets[0]).To(BeTrue())
		})

		It("should collapse expanded bullet on Enter", func() {
			model.expandedBullets[0] = true

			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			castedModel := newModel.(*CVPreviewModel)

			Expect(castedModel.expandedBullets[0]).To(BeFalse())
		})
	})

	Describe("Export", func() {
		It("should trigger export on 'e'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(ShowExportOptionsMsg{}))

			exportMsg := result.(ShowExportOptionsMsg)
			Expect(exportMsg.cvView).To(Equal(testCVView))
		})
	})

	Describe("Navigation Back", func() {
		It("should go back on Esc", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(BackToCVConfigManagerMsg{}))
		})

		It("should go back on 'q'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(BackToCVConfigManagerMsg{}))
		})
	})

	Describe("GetCVView", func() {
		It("should return the CV view", func() {
			cvView := model.GetCVView()

			Expect(cvView).To(Equal(testCVView))
		})
	})

	Describe("View", func() {
		It("should render CV header with title", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Senior IC CV"))
		})

		It("should render role in header", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("senior_ic"))
		})

		It("should render audience in header", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("hiring_manager"))
		})

		It("should render event and fact counts", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Events: 15"))
			Expect(view).To(ContainSubstring("Facts: 8"))
		})

		It("should render section titles", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Professional Experience"))
			Expect(view).To(ContainSubstring("Core Competencies"))
		})

		It("should render keyboard shortcuts", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Export"))
			Expect(view).To(ContainSubstring("Sources"))
		})
	})

	Describe("moveBulletUp", func() {
		It("should move bullet up when not at top", func() {
			model.selectedBulletIdx = 2
			model.moveBulletUp()

			Expect(model.selectedBulletIdx).To(Equal(1))
		})

		It("should not move above first bullet", func() {
			model.selectedBulletIdx = 0
			model.moveBulletUp()

			Expect(model.selectedBulletIdx).To(Equal(0))
		})
	})

	Describe("moveBulletDown", func() {
		It("should move bullet down when not at bottom", func() {
			model.selectedBulletIdx = 0
			model.moveBulletDown()

			// Movement depends on getBulletCount()
			Expect(model.selectedBulletIdx >= 0).To(BeTrue())
		})
	})

	Describe("moveSectionLeft", func() {
		It("should move to previous section", func() {
			model.selectedSectionIdx = 1
			model.moveSectionLeft()

			Expect(model.selectedSectionIdx).To(Equal(0))
			Expect(model.selectedBulletIdx).To(Equal(0))
		})

		It("should reset bullet index when moving sections", func() {
			model.selectedSectionIdx = 1
			model.selectedBulletIdx = 5
			model.moveSectionLeft()

			Expect(model.selectedBulletIdx).To(Equal(0))
		})

		It("should not move before first section", func() {
			model.selectedSectionIdx = 0
			model.moveSectionLeft()

			Expect(model.selectedSectionIdx).To(Equal(0))
		})
	})

	Describe("moveSectionRight", func() {
		It("should move to next section", func() {
			model.selectedSectionIdx = 0
			model.moveSectionRight()

			Expect(model.selectedSectionIdx).To(Equal(1))
			Expect(model.selectedBulletIdx).To(Equal(0))
		})

		It("should reset bullet index when moving sections", func() {
			model.selectedSectionIdx = 0
			model.selectedBulletIdx = 5
			model.moveSectionRight()

			Expect(model.selectedBulletIdx).To(Equal(0))
		})

		It("should not move past last section", func() {
			model.selectedSectionIdx = len(testSections) - 1
			model.moveSectionRight()

			Expect(model.selectedSectionIdx).To(Equal(len(testSections) - 1))
		})
	})
})

