package intents

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	careerdom "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

var _ = Describe("BurstManagement Intent", func() {
	var (
		intent    *BurstManagementIntent
		ctx       context.Context
		mockRepo  *MockBurstRepository
		testBurst *careerdom.Burst
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockBurstRepository()

		// Create test burst
		testBurst = &careerdom.Burst{
			ID:              "burst-1",
			Name:            "Test Burst",
			Description:     "A test burst",
			EventIDs:        []string{"event-1", "event-2"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.bursts = []*careerdom.Burst{testBurst}

		// Create intent
		data := NewBurstManagementContext(nil, mockRepo, ctx)
		var err error
		intent, err = NewBurstManagementIntent(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("Initialization", func() {
		It("should initialize with list state", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
			Expect(intent.state.currentState).To(Equal(BurstStateList))
		})

		It("should have bursts in filtered list", func() {
			intent.Init()
			Expect(len(intent.state.filteredBursts)).To(Equal(1))
		})
	})

	Describe("Update - List View", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should move selection down", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.selectedIndex).To(Equal(0))
		})

		It("should transition to detail view on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Result Handling", func() {
		It("should return completed result on success", func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
			intent.setCompleted()

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).To(HaveField("Burst", testBurst))
		})

		It("should return cancelled result on cancellation", func() {
			intent.Init()
			intent.setCancelled()

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})
	})

	Describe("View Rendering", func() {
		It("should render list view", func() {
			intent.Init()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show empty state when no bursts", func() {
			emptyRepo := NewMockBurstRepository()
			emptyCtx := NewBurstManagementContext(nil, emptyRepo, ctx)
			emptyIntent, err := NewBurstManagementIntent(emptyCtx)
			Expect(err).NotTo(HaveOccurred())
			emptyIntent.Init()
			view := emptyIntent.View()
			Expect(view).To(ContainSubstring("No bursts"))
		})

		It("should render detail view for selected burst", func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			view := intent.View()
			Expect(view).To(ContainSubstring("Burst Details"))
			Expect(view).To(ContainSubstring(testBurst.Name))
		})
	})

	Describe("Enhanced Detail View - Events", func() {
		var (
			testEvent1 *careerdom.CareerEvent
			testEvent2 *careerdom.CareerEvent
		)

		BeforeEach(func() {
			// Create test events
			testEvent1 = &careerdom.CareerEvent{
				ID:   "event-1",
				Text: "First event in burst",
				Date: time.Now().AddDate(0, -1, 0),
				Tags: []string{"tag1", "tag2"},
			}
			testEvent2 = &careerdom.CareerEvent{
				ID:   "event-2",
				Text: "Second event in burst",
				Date: time.Now(),
				Tags: []string{"tag3"},
			}

			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to events view on 'e' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.state.currentState).To(Equal(BurstStateDetailEvents))
			Expect(intent.state.loadingEvents).To(BeTrue())
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle events loaded message", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.loadingEvents = true

			msg := BurstEventsLoadedMsg{
				Events: []*careerdom.CareerEvent{testEvent1, testEvent2},
			}

			intent.Update(msg)
			Expect(intent.state.loadingEvents).To(BeFalse())
			Expect(len(intent.state.burstEvents)).To(Equal(2))
			Expect(intent.state.burstEvents[0].ID).To(Equal("event-1"))
		})

		It("should render loading state while loading events", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.loadingEvents = true

			view := intent.View()
			Expect(view).To(ContainSubstring("Loading events"))
		})

		It("should render events view with event details", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.burstEvents = []*careerdom.CareerEvent{testEvent1, testEvent2}
			intent.state.loadingEvents = false

			view := intent.View()
			Expect(view).To(ContainSubstring("Events in Burst"))
			Expect(view).To(ContainSubstring(testBurst.Name))
			Expect(view).To(ContainSubstring(testEvent1.Text))
			Expect(view).To(ContainSubstring(testEvent2.Text))
			Expect(view).To(ContainSubstring("tag1"))
		})

		It("should show empty state when no events found", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.burstEvents = []*careerdom.CareerEvent{}
			intent.state.loadingEvents = false

			view := intent.View()
			Expect(view).To(ContainSubstring("No events found"))
		})

		It("should return to detail view on escape key", func() {
			intent.state.currentState = BurstStateDetailEvents

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
		})

		It("should cancel on q key from events view", func() {
			intent.state.currentState = BurstStateDetailEvents

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should handle no burst selected error when loading events", func() {
			intent.state.selectedBurst = nil

			cmd := intent.loadEventsForBurst()
			msg := cmd().(BurstEventsLoadedMsg)

			Expect(msg.Error).To(HaveOccurred())
			Expect(msg.Error.Error()).To(ContainSubstring("no burst selected"))
		})

		It("should handle error in events loaded message", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.loadingEvents = true

			msg := BurstEventsLoadedMsg{
				Error: errors.New("failed to load events"),
			}

			intent.Update(msg)
			Expect(intent.state.loadingEvents).To(BeFalse())
			Expect(len(intent.state.burstEvents)).To(Equal(0))
		})
	})

	Describe("Enhanced Detail View - Facts", func() {
		var (
			testFact1 *careerdom.Fact
			testFact2 *careerdom.Fact
		)

		BeforeEach(func() {
			// Create test facts
			testFact1 = &careerdom.Fact{
				ID:                   "fact-1",
				Text:                 "Led team of 5 engineers",
				CompetencyCategories: []string{"leadership", "technical"},
				StrengthSignal:       "strong",
				SourceBurstID:        "burst-1",
			}
			testFact2 = &careerdom.Fact{
				ID:                   "fact-2",
				Text:                 "Improved system performance by 50%",
				CompetencyCategories: []string{"technical"},
				StrengthSignal:       "moderate",
				SourceBurstID:        "burst-1",
			}

			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to facts view on 'f' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(intent.state.currentState).To(Equal(BurstStateDetailFacts))
			Expect(intent.state.loadingFacts).To(BeTrue())
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle facts loaded message", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.loadingFacts = true

			msg := BurstFactsLoadedMsg{
				Facts: []*careerdom.Fact{testFact1, testFact2},
			}

			intent.Update(msg)
			Expect(intent.state.loadingFacts).To(BeFalse())
			Expect(len(intent.state.burstFacts)).To(Equal(2))
			Expect(intent.state.burstFacts[0].ID).To(Equal("fact-1"))
		})

		It("should render loading state while loading facts", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.loadingFacts = true

			view := intent.View()
			Expect(view).To(ContainSubstring("Loading facts"))
		})

		It("should render facts view with fact details", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.burstFacts = []*careerdom.Fact{testFact1, testFact2}
			intent.state.loadingFacts = false

			view := intent.View()
			Expect(view).To(ContainSubstring("Facts from Burst"))
			Expect(view).To(ContainSubstring(testBurst.Name))
			Expect(view).To(ContainSubstring(testFact1.Text))
			Expect(view).To(ContainSubstring(testFact2.Text))
			Expect(view).To(ContainSubstring("leadership"))
			Expect(view).To(ContainSubstring("strong"))
		})

		It("should show helpful message when no facts exist", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.burstFacts = []*careerdom.Fact{}
			intent.state.loadingFacts = false

			view := intent.View()
			Expect(view).To(ContainSubstring("No facts extracted yet"))
			Expect(view).To(ContainSubstring("Confirm the burst"))
		})

		It("should return to detail view on escape key", func() {
			intent.state.currentState = BurstStateDetailFacts

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
		})

		It("should cancel on q key from facts view", func() {
			intent.state.currentState = BurstStateDetailFacts

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should handle no burst selected error when loading facts", func() {
			intent.state.selectedBurst = nil

			cmd := intent.loadFactsForBurst()
			msg := cmd().(BurstFactsLoadedMsg)

			Expect(msg.Error).To(HaveOccurred())
			Expect(msg.Error.Error()).To(ContainSubstring("no burst selected"))
		})

		It("should handle error in facts loaded message", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.loadingFacts = true

			msg := BurstFactsLoadedMsg{
				Error: errors.New("failed to load facts"),
			}

			intent.Update(msg)
			Expect(intent.state.loadingFacts).To(BeFalse())
			Expect(len(intent.state.burstFacts)).To(Equal(0))
		})
	})

	Describe("Detail View Enhanced Footer", func() {
		It("should show enhanced keyboard shortcuts in detail view", func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst

			view := intent.View()
			Expect(view).To(ContainSubstring("e=events"))
			Expect(view).To(ContainSubstring("f=facts"))
			Expect(view).To(ContainSubstring("x=edit"))
			Expect(view).To(ContainSubstring("d=delete"))
		})
	})

	Describe("Edit Operations", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to edit state on 'x' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(BurstStateEdit))
			Expect(cmd).To(BeNil()) // initBurstEditor returns nil for now
		})

		It("should render edit view with burst details", func() {
			intent.state.currentState = BurstStateEdit

			view := intent.View()
			Expect(view).To(ContainSubstring("Edit Burst"))
			Expect(view).To(ContainSubstring(testBurst.Name))
			Expect(view).To(ContainSubstring(testBurst.Description))
			Expect(view).To(ContainSubstring("Ctrl+S=save"))
			Expect(view).To(ContainSubstring("Esc=cancel"))
		})

		It("should show placeholder message in edit view", func() {
			intent.state.currentState = BurstStateEdit

			view := intent.View()
			Expect(view).To(ContainSubstring("Full edit functionality coming soon"))
		})

		It("should cancel edit on 'Esc' key and return to detail", func() {
			intent.state.currentState = BurstStateEdit

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
			Expect(intent.state.editError).To(BeNil())
		})

		It("should save burst on 'Ctrl+S' key", func() {
			intent.state.currentState = BurstStateEdit
			intent.state.selectedBurst.Name = "Updated Burst Name"

			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))

			// Verify burst was updated in repository
			saved, err := mockRepo.Read(ctx, testBurst.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(saved.Name).To(Equal("Updated Burst Name"))
		})

		It("should handle save error gracefully", func() {
			// Create a burst that doesn't exist in repository
			nonExistentBurst := &careerdom.Burst{
				ID:   "non-existent",
				Name: "Non-existent Burst",
			}
			intent.state.selectedBurst = nonExistentBurst
			intent.state.currentState = BurstStateEdit

			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

			// Should stay in edit state with error
			Expect(intent.state.currentState).To(Equal(BurstStateEdit))
			Expect(intent.state.editError).NotTo(BeNil())
		})

		It("should display edit error in view", func() {
			intent.state.currentState = BurstStateEdit
			intent.state.editError = errors.New("failed to save burst")

			view := intent.View()
			Expect(view).To(ContainSubstring("Error: failed to save burst"))
		})

		It("should clear edit error when cancelling", func() {
			intent.state.currentState = BurstStateEdit
			intent.state.editError = errors.New("some error")

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.editError).To(BeNil())
		})

		It("should handle nil selectedBurst gracefully", func() {
			intent.state.selectedBurst = nil
			intent.state.currentState = BurstStateEdit

			view := intent.View()
			Expect(view).To(ContainSubstring("No burst selected"))
		})

		It("should reload bursts list after successful save", func() {
			intent.state.currentState = BurstStateEdit
			originalCount := len(intent.state.filteredBursts)

			// Add a new burst to the repository
			newBurst := &careerdom.Burst{
				ID:              "burst-2",
				Name:            "New Burst",
				Description:     "New burst description",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			mockRepo.bursts = append(mockRepo.bursts, newBurst)

			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

			// Verify bursts were reloaded
			Expect(len(intent.state.filteredBursts)).To(Equal(originalCount + 1))
		})
	})

	Describe("Delete Operations", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to delete confirm state on 'd' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.state.currentState).To(Equal(BurstStateDeleteConfirm))
			Expect(cmd).To(BeNil())
		})

		It("should render delete confirmation view with burst details", func() {
			intent.state.currentState = BurstStateDeleteConfirm

			view := intent.View()
			Expect(view).To(ContainSubstring("DELETE BURST"))
			Expect(view).To(ContainSubstring("Are you sure"))
			Expect(view).To(ContainSubstring(testBurst.Name))
			Expect(view).To(ContainSubstring("This will remove the burst grouping but NOT delete the events"))
			Expect(view).To(ContainSubstring("This action cannot be undone"))
			Expect(view).To(ContainSubstring("y=confirm delete"))
			Expect(view).To(ContainSubstring("n/Esc=cancel"))
		})

		It("should cancel delete on 'n' key", func() {
			intent.state.currentState = BurstStateDeleteConfirm

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
			Expect(intent.state.deleteError).To(BeNil())
		})

		It("should cancel delete on 'Esc' key", func() {
			intent.state.currentState = BurstStateDeleteConfirm

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
			Expect(intent.state.deleteError).To(BeNil())
		})

		It("should delete burst on 'y' key and return to list", func() {
			intent.state.currentState = BurstStateDeleteConfirm
			originalCount := len(mockRepo.bursts)

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Should transition to list view
			Expect(intent.state.currentState).To(Equal(BurstStateList))

			// Verify burst was deleted from repository
			Expect(len(mockRepo.bursts)).To(Equal(originalCount - 1))
			_, err := mockRepo.Read(ctx, testBurst.ID)
			Expect(err).To(HaveOccurred())
		})

		It("should reset selectedBurst and selectedIndex after delete", func() {
			intent.state.currentState = BurstStateDeleteConfirm
			intent.state.selectedIndex = 5

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(intent.state.selectedBurst).To(BeNil())
			Expect(intent.state.selectedIndex).To(Equal(0))
		})

		It("should reload bursts list after delete", func() {
			// Add another burst so we can verify the count after delete
			anotherBurst := &careerdom.Burst{
				ID:          "burst-2",
				Name:        "Another Burst",
				Description: "Another test burst",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			mockRepo.bursts = append(mockRepo.bursts, anotherBurst)
			intent.context.LoadBursts()
			intent.state.filteredBursts = intent.context.Bursts

			intent.state.currentState = BurstStateDeleteConfirm

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Should have one burst left
			Expect(len(intent.state.filteredBursts)).To(Equal(1))
			Expect(intent.state.filteredBursts[0].ID).To(Equal("burst-2"))
		})

		It("should handle delete error gracefully", func() {
			// Try to delete a burst that doesn't exist
			nonExistentBurst := &careerdom.Burst{
				ID:   "non-existent",
				Name: "Non-existent Burst",
			}
			intent.state.selectedBurst = nonExistentBurst
			intent.state.currentState = BurstStateDeleteConfirm

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Should stay in delete confirm state with error
			Expect(intent.state.currentState).To(Equal(BurstStateDeleteConfirm))
			Expect(intent.state.deleteError).NotTo(BeNil())
		})

		It("should display delete error in view", func() {
			intent.state.currentState = BurstStateDeleteConfirm
			intent.state.deleteError = errors.New("failed to delete burst")

			view := intent.View()
			Expect(view).To(ContainSubstring("Error: failed to delete burst"))
		})

		It("should clear delete error when cancelling", func() {
			intent.state.currentState = BurstStateDeleteConfirm
			intent.state.deleteError = errors.New("some error")

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.deleteError).To(BeNil())
		})

		It("should handle nil selectedBurst gracefully", func() {
			intent.state.selectedBurst = nil
			intent.state.currentState = BurstStateDeleteConfirm

			view := intent.View()
			Expect(view).To(ContainSubstring("No burst selected"))
		})

		It("should cancel on 'q' key from delete confirm", func() {
			intent.state.currentState = BurstStateDeleteConfirm

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should show warning styling in delete confirm view", func() {
			intent.state.currentState = BurstStateDeleteConfirm

			view := intent.View()
			// The view should contain warning emoji
			Expect(view).To(ContainSubstring("⚠️"))
		})
	})

	Describe("Confirmation & Fact Extraction", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to confirm state on 'c' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			Expect(intent.state.currentState).To(Equal(BurstStateConfirm))
			Expect(cmd).NotTo(BeNil()) // checkForExistingFacts returns a command
		})

		It("should show 'c=confirm burst' in detail view footer", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("c=confirm burst"))
		})

		Context("when no facts exist", func() {
			It("should transition to extracting facts state", func() {
				intent.state.currentState = BurstStateConfirm

				// Simulate no facts loaded
				msg := BurstFactsLoadedMsg{Facts: []*careerdom.Fact{}}
				intent.Update(msg)

				Expect(intent.state.currentState).To(Equal(BurstStateExtractingFacts))
				Expect(intent.state.extractingFacts).To(BeTrue())
			})

			It("should render extracting facts view", func() {
				intent.state.currentState = BurstStateExtractingFacts

				view := intent.View()
				Expect(view).To(ContainSubstring("Extracting Facts"))
				Expect(view).To(ContainSubstring(testBurst.Name))
				Expect(view).To(ContainSubstring("⏳ Extracting and saving facts"))
				Expect(view).To(ContainSubstring("This may take a few moments"))
			})

			It("should handle extraction complete message", func() {
				intent.state.currentState = BurstStateExtractingFacts
				intent.state.extractingFacts = true

				// Simulate successful extraction
				extractedFacts := []*careerdom.Fact{
					{ID: "fact-1", Text: "Test fact 1"},
					{ID: "fact-2", Text: "Test fact 2"},
				}
				msg := FactExtractionCompleteMsg{Facts: extractedFacts}

				cmd := intent.Update(msg)

				Expect(intent.state.extractingFacts).To(BeFalse())
				Expect(intent.state.extractedFactsCount).To(Equal(2))
				Expect(cmd).NotTo(BeNil()) // confirmBurstOnly command
			})

			It("should mark burst as confirmed after extraction", func() {
				intent.state.currentState = BurstStateExtractingFacts

				// Simulate extraction complete
				msg := FactExtractionCompleteMsg{
					Facts: []*careerdom.Fact{
						{ID: "fact-1", Text: "Test fact"},
					},
				}

				cmd := intent.Update(msg)
				Expect(cmd).NotTo(BeNil())

				// Execute the confirm command
				result := cmd().(BurstConfirmedMsg)
				Expect(result.Error).To(BeNil())
				Expect(result.Burst.Confirmed).To(BeTrue())
				Expect(result.Burst.ConfirmedAt).NotTo(BeNil())
			})
		})

		Context("BurstConfirmedMsg handling", func() {
			It("should handle BurstConfirmedMsg successfully", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.extractionComplete = false

				msg := BurstConfirmedMsg{
					Burst: &careerdom.Burst{
						ID:          "burst-1",
						Name:        "Test Burst",
						Confirmed:   true,
						ConfirmedAt: &[]time.Time{time.Now()}[0],
					},
					Error: nil,
				}

				cmd := intent.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(intent.state.confirmError).To(BeNil())
			})

			It("should handle BurstConfirmedMsg with error", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.confirmError = nil

				expectedErr := errors.New("confirmation failed")
				msg := BurstConfirmedMsg{
					Burst: nil,
					Error: expectedErr,
				}

				cmd := intent.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(intent.state.confirmError).To(Equal(expectedErr))
			})
		})

		Context("when facts already exist", func() {
			var existingFacts []*careerdom.Fact

			BeforeEach(func() {
				existingFacts = []*careerdom.Fact{
					{ID: "fact-1", Text: "Existing fact 1"},
					{ID: "fact-2", Text: "Existing fact 2"},
					{ID: "fact-3", Text: "Existing fact 3"},
				}
			})

			It("should show re-extract prompt", func() {
				intent.state.currentState = BurstStateConfirm

				// Simulate facts already loaded
				msg := BurstFactsLoadedMsg{Facts: existingFacts}
				intent.Update(msg)

				Expect(intent.state.showReextractPrompt).To(BeTrue())
				Expect(intent.state.existingFactsCount).To(Equal(3))
			})

			It("should render re-extract prompt view", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.showReextractPrompt = true
				intent.state.existingFactsCount = 3

				view := intent.View()
				Expect(view).To(ContainSubstring("Confirm Burst"))
				Expect(view).To(ContainSubstring("already has 3 facts extracted"))
				Expect(view).To(ContainSubstring("Do you want to extract more facts"))
				Expect(view).To(ContainSubstring("y=re-extract facts"))
				Expect(view).To(ContainSubstring("n=skip re-extraction"))
			})

			It("should re-extract facts on 'y' key", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.showReextractPrompt = true

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				Expect(intent.state.showReextractPrompt).To(BeFalse())
				Expect(intent.state.currentState).To(Equal(BurstStateExtractingFacts))
				Expect(intent.state.extractingFacts).To(BeTrue())
				Expect(cmd).NotTo(BeNil())
			})

			It("should skip re-extraction on 'n' key", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.showReextractPrompt = true

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				Expect(cmd).NotTo(BeNil()) // confirmBurstOnly command

				// Execute the confirm command
				result := cmd().(BurstConfirmedMsg)
				Expect(result.Error).To(BeNil())
				Expect(result.Burst.Confirmed).To(BeTrue())
			})

			It("should skip re-extraction on 'Esc' key", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.showReextractPrompt = true

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(cmd).NotTo(BeNil()) // confirmBurstOnly command
			})
		})

		Context("error handling", func() {
			It("should handle fact loading error", func() {
				intent.state.currentState = BurstStateConfirm

				msg := BurstFactsLoadedMsg{Error: errors.New("failed to load facts")}
				intent.Update(msg)

				Expect(intent.state.confirmError).NotTo(BeNil())
			})

			It("should handle extraction error", func() {
				intent.state.currentState = BurstStateExtractingFacts

				msg := FactExtractionCompleteMsg{Error: errors.New("extraction failed")}
				intent.Update(msg)

				Expect(intent.state.extractingFacts).To(BeFalse())
				Expect(intent.state.confirmError).NotTo(BeNil())
				Expect(intent.state.currentState).To(Equal(BurstStateConfirm))
			})

			It("should display confirmation error in view", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.confirmError = errors.New("confirmation failed")

				view := intent.View()
				Expect(view).To(ContainSubstring("Error: confirmation failed"))
			})

			It("should handle nil selectedBurst gracefully", func() {
				intent.state.selectedBurst = nil
				intent.state.currentState = BurstStateConfirm

				view := intent.View()
				Expect(view).To(ContainSubstring("No burst selected"))
			})
		})

		Context("completion and success", func() {
			It("should show success message after extraction", func() {
				intent.state.currentState = BurstStateConfirm
				intent.state.extractionComplete = true
				intent.state.extractedFactsCount = 5

				view := intent.View()
				Expect(view).To(ContainSubstring("✓ Successfully extracted and saved 5 facts"))
				Expect(view).To(ContainSubstring("Burst has been confirmed"))
			})

			It("should reload bursts list after confirmation", func() {
				originalCount := len(intent.state.filteredBursts)

				// Add a new burst
				newBurst := &careerdom.Burst{
					ID:          "burst-2",
					Name:        "New Burst",
					Description: "New description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.bursts = append(mockRepo.bursts, newBurst)

				// Execute confirm
				cmd := intent.confirmBurstOnly()
				cmd()

				// Verify bursts were reloaded
				Expect(len(intent.state.filteredBursts)).To(Equal(originalCount + 1))
			})

			It("should transition to confirm view and set extractionComplete", func() {
				intent.state.currentState = BurstStateExtractingFacts

				// Execute confirm
				cmd := intent.confirmBurstOnly()
				result := cmd().(BurstConfirmedMsg)

				Expect(result.Error).To(BeNil())
				Expect(intent.state.currentState).To(Equal(BurstStateConfirm))
				Expect(intent.state.extractionComplete).To(BeTrue())
			})
		})

		It("should cancel from confirm state on 'Esc' when not prompting", func() {
			intent.state.currentState = BurstStateConfirm
			intent.state.showReextractPrompt = false

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
			Expect(intent.state.confirmError).To(BeNil())
		})

		It("should cancel from extracting state on 'q' key", func() {
			intent.state.currentState = BurstStateExtractingFacts

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should return to detail view on any key after extraction completes", func() {
			intent.state.currentState = BurstStateConfirm
			intent.state.extractionComplete = true
			intent.state.extractedFactsCount = 3

			// Press any key (e.g., space)
			intent.Update(tea.KeyMsg{Type: tea.KeySpace})

			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
			Expect(intent.state.extractionComplete).To(BeFalse())
			Expect(intent.state.extractedFactsCount).To(Equal(0))
		})
	})

	Describe("UX Polish & Integration", func() {
		BeforeEach(func() {
			intent.Init()
		})

		Context("confirmation status display", func() {
			It("should show confirmation status in detail view for confirmed bursts", func() {
				now := time.Now()
				testBurst.Confirmed = true
				testBurst.ConfirmedAt = &now

				intent.state.currentState = BurstStateDetail
				intent.state.selectedBurst = testBurst

				view := intent.View()
				Expect(view).To(ContainSubstring("✓ Confirmed"))
				Expect(view).To(ContainSubstring("Confirmed:"))
			})

			It("should not show confirmation status for unconfirmed bursts", func() {
				testBurst.Confirmed = false
				testBurst.ConfirmedAt = nil

				intent.state.currentState = BurstStateDetail
				intent.state.selectedBurst = testBurst

				view := intent.View()
				Expect(view).NotTo(ContainSubstring("✓ Confirmed"))
			})
		})

		Context("keyboard shortcuts consistency", func() {
			It("should have consistent help text across all views", func() {
				// List view
				intent.state.currentState = BurstStateList
				listView := intent.View()
				Expect(listView).NotTo(BeEmpty())

				// Detail view
				intent.state.currentState = BurstStateDetail
				intent.state.selectedBurst = testBurst
				detailView := intent.View()
				Expect(detailView).To(ContainSubstring("Esc=back"))
				Expect(detailView).To(ContainSubstring("q=cancel"))
			})

			It("should support escape key from all states", func() {
				states := []string{
					BurstStateDetail,
					BurstStateDetailEvents,
					BurstStateDetailFacts,
					BurstStateEdit,
					BurstStateDeleteConfirm,
					BurstStateConfirm,
				}

				for _, state := range states {
					intent.state.currentState = state
					intent.state.selectedBurst = testBurst
					intent.result = nil // Reset result

					intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

					// Should either go back or cancel
					backStates := []string{BurstStateList, BurstStateDetail}
					shouldGoBack := false
					for _, backState := range backStates {
						if intent.state.currentState == backState {
							shouldGoBack = true
							break
						}
					}

					if !shouldGoBack && intent.result == nil {
						// If not going back and no result, something's wrong
						Expect(intent.state.currentState).To(Or(
							Equal(BurstStateList),
							Equal(BurstStateDetail),
						), fmt.Sprintf("Expected escape from %s to handle correctly", state))
					}
				}
			})
		})

		Context("full workflow integration", func() {
			It("should complete full confirmation workflow", func() {
				// Start in detail view
				intent.state.currentState = BurstStateDetail
				intent.state.selectedBurst = testBurst

				// Press 'c' to confirm
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
				Expect(intent.state.currentState).To(Equal(BurstStateConfirm))
				Expect(cmd).NotTo(BeNil())

				// Simulate no facts exist
				msg := BurstFactsLoadedMsg{Facts: []*careerdom.Fact{}}
				intent.Update(msg)
				Expect(intent.state.currentState).To(Equal(BurstStateExtractingFacts))

				// Simulate extraction complete
				extractMsg := FactExtractionCompleteMsg{
					Facts: []*careerdom.Fact{
						{ID: "fact-1", Text: "Test fact"},
					},
				}
				confirmCmd := intent.Update(extractMsg)
				Expect(confirmCmd).NotTo(BeNil())

				// Execute confirm command
				result := confirmCmd().(BurstConfirmedMsg)
				Expect(result.Error).To(BeNil())
				Expect(result.Burst.Confirmed).To(BeTrue())
			})

			It("should handle burst edit and reload", func() {
				// Start in detail view
				intent.state.currentState = BurstStateDetail
				intent.state.selectedBurst = testBurst
				originalName := testBurst.Name

				// Press 'x' to edit
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(intent.state.currentState).To(Equal(BurstStateEdit))

				// Change burst name
				testBurst.Name = "Updated Burst Name"

				// Press Ctrl+S to save
				intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
				Expect(intent.state.currentState).To(Equal(BurstStateDetail))

				// Verify burst was updated
				saved, err := mockRepo.Read(ctx, testBurst.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(saved.Name).To(Equal("Updated Burst Name"))

				// Reset for other tests
				testBurst.Name = originalName
			})

			It("should handle burst deletion and list refresh", func() {
				// Add another burst
				anotherBurst := &careerdom.Burst{
					ID:          "burst-2",
					Name:        "Another Burst",
					Description: "Test",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.bursts = append(mockRepo.bursts, anotherBurst)
				intent.context.LoadBursts()
				intent.state.filteredBursts = intent.context.Bursts

				// Start in detail view with first burst
				intent.state.currentState = BurstStateDetail
				intent.state.selectedBurst = testBurst

				// Press 'd' to delete
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				Expect(intent.state.currentState).To(Equal(BurstStateDeleteConfirm))

				// Press 'y' to confirm
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(intent.state.currentState).To(Equal(BurstStateList))

				// Verify burst was deleted
				_, err := mockRepo.Read(ctx, testBurst.ID)
				Expect(err).To(HaveOccurred())

				// Verify list was refreshed
				Expect(len(intent.state.filteredBursts)).To(Equal(1))
			})
		})
	})

	Describe("Pagination", func() {
		var (
			manyBurstsIntent *BurstManagementIntent
			manyBurstsRepo   *MockBurstRepository
		)

		BeforeEach(func() {
			// Create 35 bursts to span multiple pages (pageSize = 15)
			manyBurstsRepo = NewMockBurstRepository()
			for i := 0; i < 35; i++ {
				burst := &careerdom.Burst{
					ID:              "burst-" + fmt.Sprintf("%02d", i),
					Name:            fmt.Sprintf("Burst %02d", i+1),
					Description:     fmt.Sprintf("Burst description %d", i),
					EventIDs:        []string{"event-1"},
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				manyBurstsRepo.bursts = append(manyBurstsRepo.bursts, burst)
			}

			manyBurstsCtx := NewBurstManagementContext(nil, manyBurstsRepo, ctx)
			var err error
			manyBurstsIntent, err = NewBurstManagementIntent(manyBurstsCtx)
			Expect(err).NotTo(HaveOccurred())
			manyBurstsIntent.Init()
		})

		It("should display correct bursts on first page", func() {
			// Verify we're on page 1
			view := manyBurstsIntent.View()
			Expect(view).To(ContainSubstring("Page 1 of 3"))

			// Verify table shows first 15 bursts
			rows := manyBurstsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15))
		})

		It("should update table rows when navigating to next page", func() {
			// Navigate to page 2 using ctrl+d key (pgdn)
			manyBurstsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

			// Verify we're on page 2
			view := manyBurstsIntent.View()
			Expect(view).To(ContainSubstring("Page 2 of 3"))

			// Verify the selected index is now in the second page range
			Expect(manyBurstsIntent.state.selectedIndex).To(Equal(15))

			// FAILING TEST: Verify table shows the correct bursts for page 2
			// Currently the table shows ALL bursts (all 35 rows) instead of just the current page (15 rows)
			rows := manyBurstsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15), "Table should show only 15 bursts for page 2, but shows %d bursts", len(rows))
		})

		It("should update table rows when navigating to last page", func() {
			// Navigate to page 3 using ctrl+d key twice
			manyBurstsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			manyBurstsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

			// Verify we're on page 3
			view := manyBurstsIntent.View()
			Expect(view).To(ContainSubstring("Page 3 of 3"))

			// Verify the selected index is now in the third page range
			Expect(manyBurstsIntent.state.selectedIndex).To(Equal(30))

			// FAILING TEST: Verify table shows the correct bursts for page 3
			// Currently the table shows ALL bursts (all 35 rows) instead of just the current page (5 rows)
			rows := manyBurstsIntent.table.Rows()
			Expect(len(rows)).To(Equal(5), "Table should show only 5 bursts for page 3, but shows %d bursts", len(rows))
		})
	})

	Describe("Helper Functions for Enhanced Table", func() {
		Describe("formatConfirmedStatus", func() {
			It("should return green ✓ Yes for confirmed bursts", func() {
				result := intent.formatConfirmedStatus(true)
				Expect(result).To(ContainSubstring("✓"))
				Expect(result).To(ContainSubstring("Yes"))
			})

			It("should return gray ✗ No for unconfirmed bursts", func() {
				result := intent.formatConfirmedStatus(false)
				Expect(result).To(ContainSubstring("✗"))
				Expect(result).To(ContainSubstring("No"))
			})
		})

		Describe("formatDescription", func() {
			It("should truncate long descriptions to 25 chars", func() {
				longDesc := "This is a very long description that should definitely be truncated at 22 characters"
				result := intent.formatDescription(longDesc)
				// Should contain "..." for truncation
				Expect(result).To(ContainSubstring("..."))
			})

			It("should render short descriptions fully", func() {
				shortDesc := "Short desc"
				result := intent.formatDescription(shortDesc)
				Expect(result).To(ContainSubstring("Short desc"))
				Expect(result).NotTo(ContainSubstring("..."))
			})

			It("should return - for empty descriptions", func() {
				result := intent.formatDescription("")
				Expect(result).To(ContainSubstring("-"))
			})

			It("should remove newlines from descriptions", func() {
				descWithNewlines := "Line 1\nLine 2\nLine 3"
				result := intent.formatDescription(descWithNewlines)
				Expect(result).NotTo(ContainSubstring("\n"))
				Expect(result).To(ContainSubstring("Line 1 Line 2"))
			})

			It("should handle whitespace-only descriptions", func() {
				result := intent.formatDescription("   \n\t  ")
				Expect(result).To(ContainSubstring("-"))
			})
		})

		Describe("formatCreatedDate", func() {
			It("should format date as YYYY-MM-DD", func() {
				testDate := time.Date(2024, 12, 15, 10, 30, 0, 0, time.UTC)
				result := intent.formatCreatedDate(testDate)
				Expect(result).To(ContainSubstring("2024-12-15"))
			})

			It("should format different dates correctly", func() {
				testDate := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
				result := intent.formatCreatedDate(testDate)
				Expect(result).To(ContainSubstring("2023-01-05"))
			})
		})
	})

	Describe("Enhanced Table Row Generation", func() {
		var (
			testBurst1 *careerdom.Burst
			testBurst2 *careerdom.Burst
		)

		BeforeEach(func() {
			testBurst1 = &careerdom.Burst{
				ID:          "burst-1",
				Name:        "Backend API Migration",
				Description: "Migrated legacy REST API to GraphQL with performance improvements",
				EventIDs:    []string{"event-1", "event-2", "event-3"},
				Confirmed:   true,
				CreatedAt:   time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Now(),
			}

			testBurst2 = &careerdom.Burst{
				ID:          "burst-2",
				Name:        "Team Leadership",
				Description: "", // Empty description
				EventIDs:    []string{"event-4", "event-5"},
				Confirmed:   false,
				CreatedAt:   time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Now(),
			}

			mockRepo.bursts = []*careerdom.Burst{testBurst1, testBurst2}
			intent.Init()
		})

		It("should generate 5 columns for each row", func() {
			rows := intent.table.Rows()
			Expect(len(rows)).To(BeNumerically(">", 0))
			for _, row := range rows {
				Expect(len(row)).To(Equal(5), "Each row should have 5 columns")
			}
		})

		It("should include confirmed status in column 3", func() {
			rows := intent.table.Rows()
			Expect(len(rows)).To(BeNumerically(">", 0))
			// First row (confirmed=true) should contain "Yes"
			Expect(rows[0][2]).To(ContainSubstring("Yes"))
		})

		It("should include created date in column 5", func() {
			rows := intent.table.Rows()
			Expect(len(rows)).To(BeNumerically(">", 0))
			// First row should have 2024-12-15
			Expect(rows[0][4]).To(ContainSubstring("2024-12-15"))
		})

		It("should handle empty description gracefully", func() {
			rows := intent.table.Rows()
			Expect(len(rows)).To(BeNumerically(">=", 2))
			// Second row (empty description) should show "-"
			Expect(rows[1][1]).To(ContainSubstring("-"))
		})

		It("should truncate long names to fit column width", func() {
			longNameBurst := &careerdom.Burst{
				ID:              "burst-3",
				Name:            "This is an extremely long burst name that should definitely be truncated",
				Description:     "Description",
				EventIDs:        []string{"event-6"},
				Confirmed:       true,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			mockRepo.bursts = []*careerdom.Burst{longNameBurst}
			intent.Init()

			rows := intent.table.Rows()
			Expect(len(rows)).To(Equal(1))
			// Name should be truncated with "..."
			Expect(rows[0][0]).To(ContainSubstring("..."))
		})
	})
})

// MockBurstRepository is a mock implementation of BurstRepository for testing
type MockBurstRepository struct {
	bursts []*careerdom.Burst
}

func NewMockBurstRepository() *MockBurstRepository {
	return &MockBurstRepository{
		bursts: make([]*careerdom.Burst, 0),
	}
}

func (m *MockBurstRepository) Create(ctx context.Context, burst *careerdom.Burst) error {
	burst.ID = "burst-" + time.Now().Format("20060102150405")
	burst.CreatedAt = time.Now()
	burst.UpdatedAt = time.Now()
	m.bursts = append(m.bursts, burst)
	return nil
}

func (m *MockBurstRepository) Read(ctx context.Context, id string) (*careerdom.Burst, error) {
	for _, b := range m.bursts {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, errors.New("burst not found")
}

func (m *MockBurstRepository) Update(ctx context.Context, burst *careerdom.Burst) error {
	for i, b := range m.bursts {
		if b.ID == burst.ID {
			burst.UpdatedAt = time.Now()
			m.bursts[i] = burst
			return nil
		}
	}
	return errors.New("burst not found")
}

func (m *MockBurstRepository) Delete(ctx context.Context, id string) error {
	for i, b := range m.bursts {
		if b.ID == id {
			m.bursts = append(m.bursts[:i], m.bursts[i+1:]...)
			return nil
		}
	}
	return errors.New("burst not found")
}

func (m *MockBurstRepository) List(ctx context.Context, filters careerrepo.BurstListFilters) ([]*careerdom.Burst, error) {
	return m.bursts, nil
}

func (m *MockBurstRepository) Count(ctx context.Context, filters careerrepo.BurstListFilters) (int, error) {
	return len(m.bursts), nil
}

func (m *MockBurstRepository) GetByID(ctx context.Context, id string) (*careerdom.Burst, error) {
	return m.Read(ctx, id)
}
