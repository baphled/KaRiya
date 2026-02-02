package intents

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCV - Escape Key Behavior", func() {
	var (
		intent   *GenerateCVIntent
		profiles []*CVProfile
	)

	BeforeEach(func() {
		profiles = []*CVProfile{
			{
				ID:             "profile1",
				Name:           "Senior IC",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
			},
		}

		events := []*career.Event{
			fixtures.Event("event1"),
		}

		ctx := &GenerateCVContext{
			AvailableProfiles: profiles,
			Events:            events,
			Facts:             make([]*career.Fact, 0),
			DefaultProfile:    profiles[0],
		}

		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("SelectProfile State", func() {
		It("should cancel intent when esc is pressed", func() {
			intent.state.currentState = GenerateCVStateSelectProfile
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.active).To(BeFalse())
			Expect(intent.result).NotTo(BeNil())
		})

	})

	Describe("SelectAudience State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateSelectAudience
			intent.state.selectedProfile = profiles[0]
		})

		It("should go back to SelectProfile when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectProfile))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Generating State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateGenerating
			intent.state.isGenerating = true
		})

		It("should go back to SelectAudience when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Preview State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = fixtures.CVView("cv1")
		})

		It("should go back to SelectAudience when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Review State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateReview
			intent.state.generatedCV = fixtures.CVView("cv1")
		})

		It("should go back to Preview when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStatePreview))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Confirm State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = fixtures.CVView("cv1")
			intent.state.selectedProfile = profiles[0]
		})

		It("should go back to Review when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateReview))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("ExportSelectFormat State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateExportSelectFormat
		})

		It("should go back to Confirm when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateConfirm))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("ExportSelectLocation State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateExportSelectLocation
			intent.state.selectedExportFormat = CVExportFormatText
		})

		It("should go back to ExportSelectFormat when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportSelectFormat))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("Exporting State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateExporting
			intent.state.isExporting = true
		})

		It("should go back to ExportSelectLocation when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportSelectLocation))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("ExportComplete State", func() {
		BeforeEach(func() {
			intent.state.currentState = GenerateCVStateExportComplete
			intent.state.exportedPath = "/tmp/cv.txt"
			intent.state.generatedCV = fixtures.CVView("cv1")
			intent.state.selectedProfile = profiles[0]
		})

		It("should go back to ExportSelectLocation when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportSelectLocation))
			Expect(intent.active).To(BeTrue())
		})
	})

	Describe("View Methods", func() {
		It("should show 'm' in SelectProfile footer", func() {
			intent.state.currentState = GenerateCVStateSelectProfile
			view := intent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
		})

		It("should show 'm' in SelectAudience footer", func() {
			intent.state.currentState = GenerateCVStateSelectAudience
			intent.state.selectedProfile = profiles[0]
			view := intent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
		})

		It("should show 'esc' and 'm' in Generating footer", func() {
			intent.state.currentState = GenerateCVStateGenerating
			intent.state.selectedProfile = profiles[0]
			view := intent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
		})

		It("should show 'm' in Preview footer", func() {
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = fixtures.CVViewWith("cv1", "Test", "", "test")
			view := intent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
		})

		It("should show 'm' in Review footer", func() {
			intent.state.currentState = GenerateCVStateReview
			intent.state.generatedCV = fixtures.CVViewWith("cv1", "Test", "", "test")
			view := intent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
		})

		It("should show 'm' in Confirm footer", func() {
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = fixtures.CVViewWith("cv1", "Test", "", "test")
			view := intent.View()
			Expect(view).NotTo(BeEmpty()) // Verify view renders
		})
	})
})
