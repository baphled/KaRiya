package models

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ViewEventWithFacts Integration Tests", func() {
	var (
		repo         *careerrepo.MemoryRepository
		factRepo     *careerrepo.MemoryFactRepository
		service      *careerservice.Service
		ctx          context.Context
		testEvent    *career.CareerEvent
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		factRepo = careerrepo.NewMemoryFactRepository()
		service = careerservice.NewService(repo)
		service.SetFactRepository(factRepo)
		ctx = context.Background()

		// Create test event
		testEvent = &career.CareerEvent{
			ID:      "test-event-1",
			Text:    "Led cross-functional team to deliver critical microservices architecture",
			Date:    time.Now().Add(-1 * time.Hour),
			Company: "TechCorp Inc.",
			Project: "Platform Migration",
			Tags:    []string{"leadership", "technical"},
			Categories: []string{"technical", "leadership"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Persist event
		err := repo.Create(ctx, testEvent)
		Expect(err).ToNot(HaveOccurred())

		// Extract and persist facts
		facts, err := service.ExtractFactsFromEvent(ctx, testEvent)
		Expect(err).ToNot(HaveOccurred())

		// Save extracted facts (convert to pointers)
		for i := range facts {
			fact := &facts[i]
			err := service.SaveFact(ctx, fact)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Describe("Task 11.0: Integrate Fact Display with Event Details", func() {
		It("11.1: should add facts section to event detail view", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Initialize to load facts
			cmd := model.Init()
			Expect(cmd).ToNot(BeNil())

			// Execute load facts command
			msg := cmd()
			Expect(msg).To(BeAssignableToTypeOf(FactsLoadedMsg{}))

			// Update model with loaded facts
			updatedModel, _ := model.Update(msg)
			model = updatedModel.(*ViewEventWithFactsModel)

			// Verify facts section is present in view
			view := model.View()
			Expect(view).To(ContainSubstring("Extracted Facts"))
			Expect(model.factsLoaded).To(BeTrue())
		})

		It("11.2: should display extracted facts for selected event", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Verify facts are loaded
			Expect(model.eventFacts).ToNot(BeEmpty())

			// Check that facts are displayed
			view := model.View()
			for _, fact := range model.eventFacts {
				// Check for fact text (truncated or full)
				if len(fact.Text) > 100 {
					Expect(view).To(ContainSubstring(fact.Text[:50]))
				} else {
					Expect(view).To(ContainSubstring(fact.Text))
				}
			}
		})

		It("11.3: should show facts grouped by source (event vs burst)", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Verify grouping
			Expect(model.eventFacts).ToNot(BeNil())
			Expect(model.burstFacts).ToNot(BeNil())

			// View should show counts
			view := model.View()
			if len(model.eventFacts) > 0 || len(model.burstFacts) > 0 {
				Expect(view).To(MatchRegexp(`\d+ from event`))
			}
		})

		It("11.4: should allow user to review and confirm facts", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Navigate to "View Facts" action
			model.selectedAction = 2 // View Facts

			// Trigger action
			cmd = model.performAction()
			Expect(model.showFactsList).To(BeTrue())
			Expect(model.factListModel).ToNot(BeNil())
		})

		It("11.5: should allow user to edit individual facts through FactEditorModel", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Open facts list
			model.selectedAction = 2
			model.performAction()

			// Simulate fact selection in list
			model.factListModel.submitted = true

			// Update to trigger fact editor
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = updatedModel.(*ViewEventWithFactsModel)

			// Verify fact editor is shown
			Expect(model.showFactEditor).To(BeTrue())
			Expect(model.factEditorModel).ToNot(BeNil())
		})

		It("11.6: should allow user to reject facts", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Verify initial fact count
			initialCount := len(model.eventFacts)
			Expect(initialCount).To(BeNumerically(">", 0))

			// Simulate reject fact message
			if len(model.eventFacts) > 0 {
				factToReject := model.eventFacts[0]
				rejectMsg := RejectFactMsg{FactID: factToReject.ID}

				// Process reject message
				_, cmd := model.Update(rejectMsg)
				Expect(cmd).ToNot(BeNil())

				// Execute delete command
				resultMsg := cmd()
				model.Update(resultMsg)

				// Verify fact was processed
				Expect(resultMsg).To(BeAssignableToTypeOf(FactRejectedMsg{}))
			}
		})

		It("11.7: should implement keyboard navigation for fact review", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Verify shortcut keys
			initialAction := model.selectedAction

			// Navigate up
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			// Can't go below 0
			Expect(model.selectedAction).To(Equal(initialAction))

			// Navigate down
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.selectedAction).To(Equal(initialAction + 1))

			// Quick view facts shortcut
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(model.showFactsList).To(BeTrue())
		})

		It("11.8: should persist fact confirmations to database", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Get a fact to update
			if len(model.eventFacts) > 0 {
				factToUpdate := model.eventFacts[0]

				// Modify fact
				factToUpdate.Text = "Updated fact text for testing"

				// Save fact
				saveMsg := SaveFactMsg{Fact: factToUpdate}
				_, cmd := model.Update(saveMsg)
				Expect(cmd).ToNot(BeNil())

				// Execute save
				resultMsg := cmd()
				Expect(resultMsg).ToNot(BeNil())

				// Verify fact was saved
				savedFact, err := factRepo.GetByID(ctx, factToUpdate.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(savedFact.Text).To(Equal("Updated fact text for testing"))
			}
		})

		It("11.9: should handle error when fact repository is not configured", func() {
			// Create service without fact repository
			serviceWithoutFacts := careerservice.NewService(repo)
			model := NewViewEventWithFactsModel(serviceWithoutFacts, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Load facts (should return empty)
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)

			// Verify no error and empty facts
			Expect(model.err).To(BeNil())
			Expect(model.eventFacts).To(BeEmpty())
		})
	})

	Describe("Integration with FactListModel and FactEditorModel", func() {
		It("should support full workflow: view details → view facts → edit fact → save", func() {
			model := NewViewEventWithFactsModel(service, ctx, testEvent)
			model.width = 100
			model.height = 40

			// Step 1: Load facts
			cmd := model.Init()
			msg := cmd()
			model.Update(msg)
			Expect(model.factsLoaded).To(BeTrue())

			// Step 2: View facts list
			model.selectedAction = 2 // View Facts action
			model.performAction()
			Expect(model.showFactsList).To(BeTrue())

			// Step 3: Select a fact
			if len(model.eventFacts) > 0 {
				model.factListModel.submitted = true

				// Trigger edit
				updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				model = updatedModel.(*ViewEventWithFactsModel)
				Expect(model.showFactEditor).To(BeTrue())

				// Step 4: Edit and save fact
				fact := model.factEditorModel.GetFact()
				fact.Text = "Edited fact text via workflow"

				// Submit editor
				// Move focus to Save button
				model.factEditorModel.focusIndex = 5  // FactSaveButtonIdx
				model.factEditorModel.fact.CompetencyCategories = []string{"technical"}
				model.factEditorModel.fact.RoleFit = career.RoleFitSeniorIC
				model.factEditorModel.fact.AudienceRelevance = []string{"hiring_manager"}
				updatedModel, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				model = updatedModel.(*ViewEventWithFactsModel)

				// Execute save command
				if cmd != nil {
					msg := cmd()
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(*ViewEventWithFactsModel)
				}

				// Verify workflow completed
				Expect(model.showFactEditor).To(BeFalse())
			}
		})
	})
})

