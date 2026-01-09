package intents

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCVIntent - Structure Selection", func() {
	var (
		intent *GenerateCVIntent
		ctx    *GenerateCVContext
	)

	BeforeEach(func() {
		// Create test events.
		events := []*career.CareerEvent{
			{
				ID:         "event1",
				Text:       "Implemented feature X",
				Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Company:    "Company A",
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:       []string{"go", "backend"},
				Categories: []string{"technical"},
			},
		}

		// Create test profiles.
		profiles := []*CVProfile{
			{
				ID:             "profile1",
				Name:           "Senior IC",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
				Description:    "Profile for senior individual contributor roles",
			},
		}

		ctx = &GenerateCVContext{
			AvailableProfiles: profiles,
			Events:            events,
			Facts:             make([]*career.Fact, 0),
			DefaultProfile:    profiles[0],
		}

		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("CVStructure Type", func() {
		It("should define Standard structure constant", func() {
			Expect(CVStructureStandard).To(Equal(CVStructure("standard")))
		})

		It("should define Narrative structure constant", func() {
			Expect(CVStructureNarrative).To(Equal(CVStructure("narrative")))
		})
	})

	Describe("State Transitions", func() {
		BeforeEach(func() {
			intent.Init()
			// Navigate to audience selection first
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // select profile
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
		})

		It("should transition from select_audience to select_structure on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectStructure))
		})

		It("should transition from select_structure to generating on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectStructure))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // structure -> generating
			Expect(intent.state.currentState).To(Equal(GenerateCVStateGenerating))
		})

		It("should go back to select_audience on esc from select_structure", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectStructure))
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
		})

		It("should cancel on q from select_structure", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectStructure))
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel on m (main menu) from select_structure", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectStructure))
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Structure Selection UI", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
		})

		It("should show Standard and Narrative structure options", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Standard"))
			Expect(view).To(ContainSubstring("Narrative"))
		})

		It("should show descriptions for each structure", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Traditional CV"))
			Expect(view).To(ContainSubstring("professional format"))
		})

		It("should highlight currently selected structure with marker", func() {
			view := intent.View()
			// First structure (Standard) should be selected by default
			Expect(view).To(ContainSubstring("▶"))
		})

		It("should default to Standard structure", func() {
			Expect(intent.state.structureIndex).To(Equal(0))
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructureStandard))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
		})

		It("should navigate down with j key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.state.structureIndex).To(Equal(1))
		})

		It("should navigate down with down arrow", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.structureIndex).To(Equal(1))
		})

		It("should navigate up with k key", func() {
			intent.state.structureIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.state.structureIndex).To(Equal(0))
		})

		It("should navigate up with up arrow", func() {
			intent.state.structureIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.structureIndex).To(Equal(0))
		})

		It("should not go below 0", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.structureIndex).To(Equal(0))
		})

		It("should not go above max index (1)", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.structureIndex).To(Equal(1))
		})
	})

	Describe("Structure Storage", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
		})

		It("should store Standard structure when selected at index 0", func() {
			Expect(intent.state.structureIndex).To(Equal(0))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructureStandard))
		})

		It("should store Narrative structure when selected at index 1", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // select Narrative
			Expect(intent.state.structureIndex).To(Equal(1))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructureNarrative))
		})
	})

	Describe("Breadcrumbs", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
		})

		It("should include 'Select Structure' in breadcrumbs", func() {
			crumbs := intent.getBreadcrumbs()
			Expect(crumbs).To(ContainElement("Select Structure"))
		})
	})

	Describe("Result Metadata", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> structure
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // select Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // structure -> generating
		})

		It("should include selected structure in result when completed", func() {
			// Simulate CV generation completion
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{
					ID:               "test-cv",
					Name:             "Test CV",
					TargetRole:       "senior_ic",
					TargetAudience:   "hiring_manager",
					GeneratedAt:      time.Now(),
					SourceEventCount: 1,
					SourceFactCount:  0,
				},
				Error: nil,
			})

			// Navigate to confirm and complete
			intent.state.currentState = GenerateCVStateConfirm
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Data).NotTo(BeNil())
			Expect(intent.result.Data.SelectedStructure).To(Equal(CVStructureNarrative))
		})

		It("should include structure in metadata", func() {
			// Simulate CV generation completion
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{
					ID:               "test-cv",
					Name:             "Test CV",
					TargetRole:       "senior_ic",
					TargetAudience:   "hiring_manager",
					GeneratedAt:      time.Now(),
					SourceEventCount: 1,
					SourceFactCount:  0,
				},
				Error: nil,
			})

			// Navigate to confirm and complete
			intent.state.currentState = GenerateCVStateConfirm
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Metadata).To(HaveKey("structure"))
			Expect(intent.result.Metadata["structure"]).To(Equal(CVStructureNarrative))
		})
	})
})
