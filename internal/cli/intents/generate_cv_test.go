package intents

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCVIntent", func() {
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
			{
				ID:         "event2",
				Text:       "Led team on project Y",
				Date:       time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				Company:    "Company B",
				CreatedAt:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				Tags:       []string{"leadership"},
				Categories: []string{"management"},
			},
		}

		// Create test profiles.
		profiles := []*CVProfile{
			{
				ID:             "profile1",
				Name:           "Senior IC",
				TargetRole:     "senior_ic",
				TargetAudience: []string{"hiring_manager", "recruiter"},
				Description:    "Profile for senior individual contributor roles",
			},
			{
				ID:             "profile2",
				Name:           "Staff Engineer",
				TargetRole:     "staff",
				TargetAudience: []string{"peer"},
				Description:    "Profile for staff engineer roles",
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

	Describe("Creation", func() {
		It("should create a valid intent", func() {
			Expect(intent).NotTo(BeNil())
			Expect(intent.active).To(BeTrue())
		})

		It("should initialize with profile selection state", func() {
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectProfile))
		})

		It("should have default profile selected", func() {
			Expect(intent.state.selectedProfile).NotTo(BeNil())
			Expect(intent.state.selectedProfile.ID).To(Equal("profile1"))
		})
	})

	Describe("Init", func() {
		It("should initialize without error", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("should keep default profile selected", func() {
			intent.Init()
			Expect(intent.state.selectedProfile.ID).To(Equal("profile1"))
		})
	})

	Describe("Update - Profile Selection", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should move selection down with down arrow", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.selectedIndex).To(Equal(1))
			Expect(intent.state.selectedProfile.ID).To(Equal("profile2"))
		})

		It("should move selection down with j key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.state.selectedIndex).To(Equal(1))
			Expect(intent.state.selectedProfile.ID).To(Equal("profile2"))
		})

		It("should move selection up with up arrow", func() {
			intent.state.selectedIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(intent.state.selectedIndex).To(Equal(0))
			Expect(intent.state.selectedProfile.ID).To(Equal("profile1"))
		})

		It("should move selection up with k key", func() {
			intent.state.selectedIndex = 1
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.state.selectedIndex).To(Equal(0))
			Expect(intent.state.selectedProfile.ID).To(Equal("profile1"))
		})

		It("should transition to audience selection on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
		})

		It("should set selected audiences from profile", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.selectedAudiences).To(Equal([]string{"hiring_manager", "recruiter"}))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should cancel on ctrl+c", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Update - Audience Selection", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateSelectAudience
		})

		It("should transition to generating on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateGenerating))
		})

		It("should generate CV when async generation completes", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Simulate async generation completion
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{
					ID:               "test-cv",
					Name:             intent.state.selectedProfile.Name,
					TargetRole:       intent.state.selectedProfile.TargetRole,
					TargetAudience:   intent.state.selectedAudiences,
					GeneratedAt:      time.Now(),
					SourceEventCount: 0,
					SourceFactCount:  0,
				},
				Error: nil,
			})
			Expect(intent.state.currentState).To(Equal(GenerateCVStatePreview))
			Expect(intent.state.generatedCV).NotTo(BeNil())
			Expect(intent.state.generatedCV.Name).To(Equal(intent.state.selectedProfile.Name))
		})

		It("should go back to profile selection on esc", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectProfile))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Update - Preview", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
		})

		It("should transition to review on e key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateReview))
		})

		It("should transition to confirm on c key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateConfirm))
		})

		It("should go back to audience selection on esc", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Update - Review", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateReview
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
		})

		It("should transition to confirm on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateConfirm))
		})

		It("should go back to preview on esc", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStatePreview))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Update - Confirm", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
		})

		It("should complete on y key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(intent.result.Status).To(Equal(Completed))
		})

		It("should complete on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.result.Status).To(Equal(Completed))
		})

		It("should go back to review on n key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateReview))
		})

		It("should go back to review on esc", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateReview))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should render profile selection view", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Select CV Profile"))
		})

		It("should show profiles in list", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Senior IC"))
			Expect(view).To(ContainSubstring("Staff Engineer"))
		})

		It("should show selected profile with marker", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring(">"))
		})

		It("should render audience selection view", func() {
			intent.state.currentState = GenerateCVStateSelectAudience
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Select Target Audience"))
		})

		It("should render preview view", func() {
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("CV Preview"))
		})

		It("should render review view", func() {
			intent.state.currentState = GenerateCVStateReview
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Review & Edit CV"))
		})

		It("should render confirm view", func() {
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Confirm CV Generation"))
		})
	})

	Describe("Result", func() {
		It("should return nil when no result set", func() {
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return completed result after completion", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).NotTo(BeNil())
		})

		It("should include generated CV in result data", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			resultData := result.Data.(*GenerateCVResult)
			Expect(resultData.GeneratedCV).NotTo(BeNil())
			Expect(resultData.GeneratedCV.ID).To(Equal("cv_123"))
		})

		It("should include selected profile in result data", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			resultData := result.Data.(*GenerateCVResult)
			Expect(resultData.SelectedProfile).NotTo(BeNil())
			Expect(resultData.SelectedProfile.ID).To(Equal("profile1"))
		})

		It("should include metadata in result", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStateConfirm
			intent.state.generatedCV = &career.CVView{
				ID:               "cv_123",
				Name:             "Senior IC",
				TargetRole:       "senior_ic",
				TargetAudience:   []string{"hiring_manager"},
				GeneratedAt:      time.Now(),
				SourceEventCount: 2,
				SourceFactCount:  0,
			}
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result.Metadata).NotTo(BeNil())
			Expect(result.Metadata["profile"]).To(Equal("profile1"))
			Expect(result.Metadata["event_count"]).To(Equal(2))
		})
	})
})
