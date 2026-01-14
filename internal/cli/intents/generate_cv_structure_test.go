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

	Describe("State Transitions (Variant-Based Flow)", func() {
		BeforeEach(func() {
			intent.Init()
			// Navigate to audience selection first
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // select profile
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
		})

		It("should transition from select_audience to select_role_emphasis on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectRoleEmphasis))
		})

		It("should transition from select_role_emphasis to select_length_format on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectRoleEmphasis))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectLengthFormat))
		})

		It("should transition from select_length_format to generating on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format -> generating
			Expect(intent.state.currentState).To(Equal(GenerateCVStateGenerating))
		})

		It("should go back to select_audience on esc from select_role_emphasis", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectRoleEmphasis))
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
		})

		It("should go back to select_role_emphasis on esc from select_length_format", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectLengthFormat))
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectRoleEmphasis))
		})

		It("should cancel on q from select_role_emphasis", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectRoleEmphasis))
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel on m (main menu) from select_length_format", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectLengthFormat))
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Role Emphasis Selection UI", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
		})

		It("should show all 4 role emphasis options", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Senior Backend"))
			Expect(view).To(ContainSubstring("Staff/Principal"))
			Expect(view).To(ContainSubstring("Consulting"))
			Expect(view).To(ContainSubstring("Language-Agnostic"))
		})

		It("should show descriptions for each role emphasis", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("backend engineering"))
			Expect(view).To(ContainSubstring("technical leadership"))
		})

		It("should highlight currently selected role emphasis with marker", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("▶"))
		})

		It("should default to first role emphasis (Senior Backend)", func() {
			Expect(intent.state.roleEmphasisIndex).To(Equal(0))
		})
	})

	Describe("Length Format Selection UI", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
		})

		It("should show all 4 length format options", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Full"))
			Expect(view).To(ContainSubstring("Standard"))
			Expect(view).To(ContainSubstring("Short"))
			Expect(view).To(ContainSubstring("1-Page"))
		})

		It("should show page counts for each length format", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("3+"))
			Expect(view).To(ContainSubstring("2-3"))
			Expect(view).To(ContainSubstring("1-2"))
		})

		It("should highlight currently selected length format with marker", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("▶"))
		})

		It("should default to Standard length format (index 1)", func() {
			Expect(intent.state.lengthFormatIndex).To(Equal(1))
		})
	})

	Describe("Role Emphasis Navigation", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
		})

		It("should navigate down with j key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.state.roleEmphasisIndex).To(Equal(1))
		})

		It("should navigate down with down arrow", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.roleEmphasisIndex).To(Equal(1))
		})

		It("should navigate up with k key", func() {
			intent.state.roleEmphasisIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.state.roleEmphasisIndex).To(Equal(0))
		})

		It("should navigate up with up arrow", func() {
			intent.state.roleEmphasisIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.roleEmphasisIndex).To(Equal(0))
		})

		It("should not go below 0", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.roleEmphasisIndex).To(Equal(0))
		})

		It("should not go above max index (3)", func() {
			for i := 0; i < 5; i++ {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			Expect(intent.state.roleEmphasisIndex).To(Equal(3))
		})
	})

	Describe("Length Format Navigation", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
		})

		It("should navigate down with j key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.state.lengthFormatIndex).To(Equal(2)) // default is 1
		})

		It("should navigate up with k key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.state.lengthFormatIndex).To(Equal(0)) // default is 1, up goes to 0
		})

		It("should not go above max index (3)", func() {
			for i := 0; i < 5; i++ {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			Expect(intent.state.lengthFormatIndex).To(Equal(3))
		})
	})

	Describe("Variant Selection and Storage", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
		})

		It("should store selected role emphasis", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // select Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.selectedRoleEmphasis).To(Equal(RoleEmphasisStaffPrincipal))
		})

		It("should store selected length format and look up variant", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // select Senior Backend
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectLengthFormat))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // select Standard (default at index 1)
			Expect(intent.state.selectedLengthFormat).To(Equal(LengthStandard))
			Expect(intent.state.selectedVariant).NotTo(BeNil())
			Expect(intent.state.selectedVariant.ID).To(Equal("senior_backend_standard"))
		})

		It("should set base structure from variant", func() {
			// Select Consulting role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Consulting
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Select Full length format
			intent.Update(tea.KeyMsg{Type: tea.KeyUp}) // go to Full (index 0)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Consulting Full uses CVStructureConsulting
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructure("consulting")))
		})
	})

	Describe("Breadcrumbs", func() {
		It("should include 'Select Role Emphasis' in breadcrumbs", func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			crumbs := intent.getBreadcrumbs()
			Expect(crumbs).To(ContainElement("Select Role Emphasis"))
		})

		It("should include 'Select Length' in breadcrumbs", func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			crumbs := intent.getBreadcrumbs()
			Expect(crumbs).To(ContainElement("Select Length"))
		})
	})

	Describe("Result Metadata with Variant", func() {
		BeforeEach(func() {
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
			// Select Language-Agnostic (index 3)
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // 1
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // 2
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // 3
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format (Standard) -> generating
		})

		It("should include selected variant in result when completed", func() {
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
			Expect(intent.result.Data.SelectedVariant).NotTo(BeNil())
			Expect(intent.result.Data.SelectedVariant.ID).To(Equal("language_agnostic_standard"))
		})

		It("should include role_emphasis and length_format in metadata", func() {
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
			Expect(intent.result.Metadata).To(HaveKey("role_emphasis"))
			Expect(intent.result.Metadata).To(HaveKey("length_format"))
			Expect(intent.result.Metadata).To(HaveKey("variant_id"))
			Expect(intent.result.Metadata["variant_id"]).To(Equal("language_agnostic_standard"))
		})

		It("should set base structure from variant (Narrative for Language-Agnostic)", func() {
			// Language-Agnostic uses CVStructureNarrative as base
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructure("narrative")))
		})
	})
})
