package models_test

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
	careerdom "github.com/baphled/kariya/internal/domain/career"
)


var _ = ginkgo.Describe("SourceEventTracer Model", func() {
	var (
		model  *models.SourceEventTracerModel
		bullet *careerdom.CVBullet
	)

	ginkgo.BeforeEach(func() {
		baseModel := models.NewBaseStandardModel()
		bullet = &careerdom.CVBullet{
			ID:              "bullet-1",
			Text:            "Led team through major architectural redesign",
			SourceEventIDs:  []string{"event-1", "event-2"},
			SourceFactIDs:   []string{"fact-1"},
			Rank:            0.85,
			InclusionReason: "ownership",
			Confidence:      0.92,
		}

		// Create model with nil repositories for testing
		model = models.NewSourceEventTracerModel(
			baseModel,
			bullet,
			nil,
			nil,
			nil,
		)
	})

	ginkgo.Describe("initialization", func() {
		ginkgo.It("should create a new source event tracer model", func() {
			gomega.Expect(model).NotTo(gomega.BeNil())
		})

		ginkgo.It("should have empty error initially", func() {
			gomega.Expect(model.GetLastError()).To(gomega.BeNil())
		})

		ginkgo.It("should store bullet information", func() {
			gomega.Expect(model.GetBullet()).To(gomega.Equal(bullet))
		})
	})

	ginkgo.Describe("rendering", func() {
		ginkgo.BeforeEach(func() {
			// Disable loading for rendering tests
			model.SetLoading(false)
		})

		ginkgo.It("should render without error", func() {
			view := model.View()
			gomega.Expect(view).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("should include header", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Bullet Source Traceability"))
		})

		ginkgo.It("should display bullet text", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Led team through major architectural redesign"))
		})

		ginkgo.It("should display bullet metadata", func() {
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Rank"))
			gomega.Expect(view).To(gomega.ContainSubstring("Confidence"))
			gomega.Expect(view).To(gomega.ContainSubstring("ownership"))
		})

		ginkgo.It("should display loading state initially", func() {
			// Create a fresh model without SetLoading(false)
			baseModel := models.NewBaseStandardModel()
			freshModel := models.NewSourceEventTracerModel(
				baseModel,
				bullet,
				nil,
				nil,
				nil,
			)
			// Model starts in loading state
			view := freshModel.View()
			gomega.Expect(view).To(gomega.ContainSubstring("Loading"))
		})

		ginkgo.It("should handle window size updates", func() {
			_, cmd := model.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
			gomega.Expect(cmd).To(gomega.BeNil())
			gomega.Expect(model.GetWidth()).To(gomega.Equal(120))
			gomega.Expect(model.GetHeight()).To(gomega.Equal(30))
		})
	})

	ginkgo.Describe("model interface", func() {
		ginkgo.It("should implement Model interface", func() {
			var _ tea.Model = model
		})

		ginkgo.It("should implement StandardModel interface", func() {
			var _ models.StandardModel = model
		})

		ginkgo.It("should have Init method", func() {
			cmd := model.Init()
			gomega.Expect(cmd).NotTo(gomega.BeNil())
		})

		ginkgo.It("should return BaseStandardModel", func() {
			gomega.Expect(model.BaseStandardModel).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("escape handling", func() {
		ginkgo.It("should handle escape key", func() {
			var backPressed bool
			model.SetOnBack(func() {
				backPressed = true
			})

			model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			gomega.Expect(backPressed).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("message handling", func() {
		ginkgo.It("should handle SourcesLoadedMsg", func() {
			events := []*careerdom.CareerEvent{
				{
					ID:         "event-1",
					Text:       "Redesigned system architecture",
					Date:       time.Now(),
					Company:    "TechCorp",
					Tags:       []string{"architecture", "technical"},
					Categories: []string{"technical"},
				},
			}
			facts := []*careerdom.Fact{
				{
					ID:             "fact-1",
					Text:           "Technical leadership on complex architecture",
					StrengthSignal: "ownership",
				},
			}

			msg := models.SourcesLoadedMsg{Events: events, Facts: facts}
			newModel, _ := model.Update(msg)

			tracerModel := newModel.(*models.SourceEventTracerModel)
			gomega.Expect(tracerModel.GetSourceEvents()).To(gomega.HaveLen(1))
			gomega.Expect(tracerModel.GetSourceFacts()).To(gomega.HaveLen(1))
		})

		ginkgo.It("should handle SourcesLoadErrorMsg", func() {
			msg := models.SourcesLoadErrorMsg{Err: "Failed to load events"}
			newModel, _ := model.Update(msg)

			tracerModel := newModel.(*models.SourceEventTracerModel)
			gomega.Expect(tracerModel.GetLastError()).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("scroll navigation", func() {
		ginkgo.It("should initialize with zero scroll offset", func() {
			gomega.Expect(model.GetScrollOffset()).To(gomega.Equal(0))
		})

		ginkgo.It("should handle up/k navigation in loading state", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			gomega.Expect(model.GetScrollOffset()).To(gomega.Equal(0))
		})
	})
})
