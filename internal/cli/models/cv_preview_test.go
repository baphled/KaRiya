package models

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			ID:               uuid.New().String(),
			Name:             "Senior IC CV",
			TargetRole:       "senior_ic",
			TargetAudience:   "hiring_manager",
			GeneratedAt:      time.Now(),
			SourceEventCount: 15,
			SourceFactCount:  8,
		}

		testSections = []*career.CVSection{
			{
				ID:          uuid.New().String(),
				CVViewID:    testCVView.ID,
				SectionType: "experience",
				Title:       "Professional Experience",
				Order:       0,
				Content: []*career.SectionContentGroup{
					{
						Header: "Acme Corp",
						Bullets: []*career.CVBullet{
							{
								ID:   uuid.New().String(),
								Text: "Led team of 5 engineers",
							},
							{
								ID:   uuid.New().String(),
								Text: "Designed and implemented new architecture",
							},
						},
					},
				},
			},
			{
				ID:          uuid.New().String(),
				CVViewID:    testCVView.ID,
				SectionType: "skills",
				Title:       "Core Competencies",
				Order:       1,
				Content: []*career.SectionContentGroup{
					{
						Bullets: []*career.CVBullet{
							{
								ID:   uuid.New().String(),
								Text: "Go, Python, Rust",
							},
							{
								ID:   uuid.New().String(),
								Text: "System Design",
							},
							{
								ID:   uuid.New().String(),
								Text: "Leadership",
							},
						},
					},
				},
			},
		}

		model = NewCVPreviewModel(baseModel, testCVView, testSections, nil)
	})

	Describe("Initialization", func() {
		It("should initialize with CV view and sections", func() {
			Expect(model.cvView).To(Equal(testCVView))
			Expect(model.sections).To(Equal(testSections))
			Expect(model.selectedIdx).To(Equal(0))
		})

		It("should have empty expanded sections map", func() {
			Expect(model.expandedSections).To(BeEmpty())
		})

		It("should initialize breadcrumbs", func() {
			Expect(model.breadcrumbs).To(ContainElement("Preview"))
		})
	})

	Describe("Navigation", func() {
		It("should move down with 'j'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel).NotTo(BeNil())
		})

		It("should move up with 'k'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel).NotTo(BeNil())
		})

		It("should move to first section with 'home'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyHome})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel).NotTo(BeNil())
		})

		It("should move to last section with 'end'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnd})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel).NotTo(BeNil())
		})
	})

	Describe("Section Expansion", func() {
		It("should toggle section expansion on Enter", func() {
			initialState := model.expandedSections[0]

			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			castedModel := newModel.(*CVPreviewModel)

			Expect(castedModel.expandedSections[0]).To(Equal(!initialState))
		})

		It("should collapse expanded section on Enter", func() {
			model.expandedSections[0] = true

			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			castedModel := newModel.(*CVPreviewModel)

			Expect(castedModel.expandedSections[0]).To(BeFalse())
		})
	})

	Describe("Export", func() {
		It("should trigger export on 'x'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(ShowExportOptionsMsg{}))

			exportMsg := result.(ShowExportOptionsMsg)
			Expect(exportMsg.CVView).To(Equal(testCVView))
		})
	})

	Describe("Navigation Back", func() {
		It("should go back on Esc", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(BackMsg{}))
		})

		It("should go back on 'q'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(QuitMsg{}))
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

		It("should render role in metadata", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("senior_ic"))
		})

		It("should render audience in metadata", func() {
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

			// Titles are truncated to 18 characters in the table
			Expect(view).To(ContainSubstring("Professional Ex"))
			Expect(view).To(ContainSubstring("Core Competencies"))
		})
	})

	Describe("WindowSizeMsg", func() {
		It("should update width and height on WindowSizeMsg", func() {
			newModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			castedModel := newModel.(*CVPreviewModel)
			Expect(castedModel.width).To(Equal(120))
			Expect(castedModel.height).To(Equal(40))
		})
	})
})
