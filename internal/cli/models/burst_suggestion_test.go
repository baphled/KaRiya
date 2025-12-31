package models

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstSuggestionModel", func() {
	var (
		repo           *careerrepo.MemoryRepository
		svc            *careerservice.Service
		ctx            context.Context
		suggestions    []burstfact.BurstSuggestion
		event1, event2 *career.CareerEvent
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		ctx = context.Background()

		// Create test events
		event1 = &career.CareerEvent{
			ID:   "event1",
			Text: "Led backend team on microservices migration project",
			Date: time.Now().Add(-30 * 24 * time.Hour),
			Tags: []string{"leadership", "technical"},
		}

		event2 = &career.CareerEvent{
			ID:   "event2",
			Text: "Architected service mesh for improved scalability",
			Date: time.Now().Add(-25 * 24 * time.Hour),
			Tags: []string{"technical", "architecture"},
		}

		// Save events to repository
		repo.Create(ctx, event1)
		repo.Create(ctx, event2)

		// Create test suggestions
		suggestions = []burstfact.BurstSuggestion{
			{
				EventIDs:        []string{"event1", "event2"},
				ConfidenceScore: 0.85,
				Name:            "Microservices Architecture Initiative",
				Description:     "Led migration to microservices with service mesh implementation",
			},
			{
				EventIDs:        []string{"event1"},
				ConfidenceScore: 0.6,
				Name:            "",
				Description:     "",
			},
		}
	})

	Describe("NewBurstSuggestionModel", func() {
		It("creates model with empty suggestions", func() {
			model := NewBurstSuggestionModel(svc, []burstfact.BurstSuggestion{}, ctx)
			Expect(model).NotTo(BeNil())
			Expect(model.currentIdx).To(Equal(0))
			Expect(model.editing).To(BeFalse())
			Expect(len(model.confirmed)).To(Equal(0))
			Expect(len(model.rejected)).To(Equal(0))
		})

		It("creates model with suggestions", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			Expect(model).NotTo(BeNil())
			Expect(len(model.suggestions)).To(Equal(2))
			Expect(model.currentIdx).To(Equal(0))
		})

		It("initializes input fields for editing", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			Expect(len(model.inputs)).To(Equal(2))
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("handles window size message", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			updatedModel, _ := model.Update(msg)

			Expect(updatedModel).NotTo(BeNil())
			Expect(model.width).To(Equal(120))
			Expect(model.height).To(Equal(40))
		})

		It("handles key message for navigation down", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			initialIdx := model.currentIdx

			msg := tea.KeyMsg{Type: tea.KeyDown}
			model.Update(msg)

			Expect(model.currentIdx).To(Equal(initialIdx + 1))
		})

		It("handles key message for navigation up", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.currentIdx = 1

			msg := tea.KeyMsg{Type: tea.KeyUp}
			model.Update(msg)

			Expect(model.currentIdx).To(Equal(0))
		})

		It("handles Escape key", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			updatedModel, _ := model.Update(msg)

			Expect(updatedModel).NotTo(BeNil())
		})
	})

	Describe("Confirmation Workflow", func() {
		It("confirms current suggestion with 'y' key", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			model.Update(msg)

			Expect(len(model.confirmed)).To(Equal(1))
			Expect(model.confirmed[0].EventIDs).To(Equal(suggestions[0].EventIDs))
		})

		It("rejects current suggestion with 'n' key", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			model.Update(msg)

			Expect(len(model.rejected)).To(Equal(1))
			Expect(model.rejected[0].EventIDs).To(Equal(suggestions[0].EventIDs))
		})

		It("moves to next suggestion after confirmation", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			model.Update(msg)

			Expect(model.currentIdx).To(Equal(1))
		})

		It("moves to next suggestion after rejection", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			model.Update(msg)

			Expect(model.currentIdx).To(Equal(1))
		})

		It("processes multiple confirmations and rejections", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)

			// Confirm first
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(len(model.confirmed)).To(Equal(1))

			// Reject second
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(len(model.rejected)).To(Equal(1))

			// Both should be marked as done
			Expect(model.IsDone()).To(BeTrue())
		})
	})

	Describe("Editing Workflow", func() {
		It("starts editing with 'e' key", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			model.Update(msg)

			Expect(model.editing).To(BeTrue())
		})

		It("populates edit inputs with current suggestion values", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			Expect(model.inputs[0].Value()).To(Equal(suggestions[0].Name))
			Expect(model.inputs[1].Value()).To(Equal(suggestions[0].Description))
		})

		It("saves edited values on Enter", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			model.inputs[0].SetValue("Updated Name")
			model.inputs[1].SetValue("Updated Description")
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(model.editing).To(BeFalse())
			Expect(model.editedNames[0]).To(Equal("Updated Name"))
			Expect(model.editedDescs[0]).To(Equal("Updated Description"))
		})

		It("exits editing on Escape", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			Expect(model.editing).To(BeTrue())
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.editing).To(BeFalse())
		})

		It("uses edited values when confirming", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			model.inputs[0].SetValue("Custom Name")
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(model.confirmed[0].Name).To(Equal("Custom Name"))
		})
	})

	Describe("Edge Cases", func() {
		It("handles empty suggestions list gracefully", func() {
			model := NewBurstSuggestionModel(svc, []burstfact.BurstSuggestion{}, ctx)
			view := model.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("handles navigation with single suggestion", func() {
			singleSuggestion := []burstfact.BurstSuggestion{suggestions[0]}
			model := NewBurstSuggestionModel(svc, singleSuggestion, ctx)

			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.currentIdx).To(Equal(0))
		})

		It("displays related events correctly", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			view := model.View()

			Expect(view).To(ContainSubstring("Related Events"))
		})

		It("maintains related events cache", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)

			// Cache map should be initialized and empty
			Expect(model.relatedEvents).NotTo(BeNil())
			Expect(len(model.relatedEvents)).To(Equal(0))

			// Populate cache
			events := []*career.CareerEvent{event1}
			model.relatedEvents[0] = events

			// Cache entry should exist
			Expect(model.relatedEvents[0]).NotTo(BeNil())
			Expect(len(model.relatedEvents[0])).To(Equal(1))
		})
	})

	Describe("View Rendering", func() {
		It("renders review view by default", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			view := model.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Burst Suggestion"))
		})

		It("includes confidence score in view", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			view := model.View()

			Expect(view).To(ContainSubstring("Confidence"))
			Expect(view).To(ContainSubstring("85%"))
		})

		It("includes burst name in view", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			view := model.View()

			Expect(view).To(ContainSubstring(suggestions[0].Name))
		})

		It("includes navigation help text", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			view := model.View()

			Expect(view).To(ContainSubstring("↑/↓"))
			Expect(view).To(ContainSubstring("y confirm"))
		})

		It("renders edit view when editing", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			view := model.View()
			Expect(view).To(ContainSubstring("Edit Burst"))
		})

		It("shows progress in burst suggestion view", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			view := model.View()

			Expect(view).To(ContainSubstring("1 of 2"))
		})

		It("updates progress after confirmation", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			view := model.View()
			Expect(view).To(ContainSubstring("2 of 2"))
		})
	})

	Describe("GetConfirmed", func() {
		It("returns empty list initially", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			Expect(model.GetConfirmed()).To(BeEmpty())
		})

		It("returns confirmed suggestions", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			confirmed := model.GetConfirmed()
			Expect(len(confirmed)).To(Equal(1))
			Expect(confirmed[0].EventIDs).To(Equal(suggestions[0].EventIDs))
		})
	})

	Describe("GetRejected", func() {
		It("returns empty list initially", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			Expect(model.GetRejected()).To(BeEmpty())
		})

		It("returns rejected suggestions", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			rejected := model.GetRejected()
			Expect(len(rejected)).To(Equal(1))
			Expect(rejected[0].EventIDs).To(Equal(suggestions[0].EventIDs))
		})
	})

	Describe("IsDone", func() {
		It("returns false when suggestions remain", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)
			Expect(model.IsDone()).To(BeFalse())
		})

		It("returns true when all suggestions processed", func() {
			model := NewBurstSuggestionModel(svc, suggestions, ctx)

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(model.IsDone()).To(BeTrue())
		})
	})
})
