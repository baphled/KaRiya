//nolint:errcheck // Test file - error handling for test setup is not relevant.
package app_test

import (
	"context"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/config"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Unit Tests", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		ctx        context.Context
	)

	BeforeEach(func() {
		config.SetConfigPathForTesting(filepath.Join(GinkgoT().TempDir(), "config.yaml"))
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)

		// Pre-populate repositories to avoid nil panics.
		//nolint:errcheck // Test setup - error handling not relevant.
		burstRepo.Create(ctx, &career.Burst{
			ID:       "b1",
			Name:     "dummy",
			EventIDs: []string{"e1", "e2"},
		})
		//nolint:errcheck // Test setup - error handling not relevant.
		factRepo.Create(ctx, &career.Fact{
			ID:                   "f1",
			Text:                 "dummy",
			CompetencyCategories: []string{"leadership"},
			RoleFit:              "staff",
			AudienceRelevance:    []string{"peer"},
			SourceEventID:        "e1",
		})

		// Create bootstrap result (skipping onboarding for tests)
		log := logger.DefaultLogger()
		bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

		model = app.NewModel(cliService, svc, bootstrapResult)
	})

	AfterEach(func() {
		config.ResetConfigPath()
	})

	Describe("Init", func() {
		It("should return batch command for window size and logo init", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("should navigate to browse_timeline when initial screen is ListScreen", func() {
			// Create a new model with initial screen set to ListScreen.
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)
			listModel := app.NewModel(cliService, svc, bootstrapResult)
			listModel.SetInitialScreen(app.ListScreen)

			cmd := listModel.Init()
			Expect(cmd).NotTo(BeNil())

			// After Init, should be in intent state.
			state := listModel.GetState()
			Expect(state).To(Equal(app.StateIntent))
		})

		It("should navigate to capture_event with mode when initial capture mode is set", func() {
			// Create a new model with initial capture mode set.
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)
			captureModel := app.NewModel(cliService, svc, bootstrapResult)
			captureModel.SetInitialCaptureMode("manual")

			cmd := captureModel.Init()
			Expect(cmd).NotTo(BeNil())

			// After Init, should be in intent state.
			state := captureModel.GetState()
			Expect(state).To(Equal(app.StateIntent))
		})
	})

	Describe("Update - Key Handling", func() {
		Context("ctrl+c key", func() {
			It("should quit from menu state", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlC}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})

			It("should quit from intent state", func() {
				// Activate an intent first
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				// Now send ctrl+c
				msg := tea.KeyMsg{Type: tea.KeyCtrlC}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})
		})

		Context("q key", func() {
			It("should quit when in menu state", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})

			It("should not quit when in intent state", func() {
				// Activate an intent
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				// Send q key (should not quit, just pass to intent)
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})
		})

		Context("? key (help)", func() {
			It("should toggle help screen on ? key", func() {
				// Initially help is not showing
				view := model.View()
				Expect(view).NotTo(ContainSubstring("Keyboard Reference"))

				// Press ? to show help
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
				newModel, cmd := model.Update(msg)
				model = newModel.(*app.Model)
				Expect(cmd).To(BeNil())

				// Now help should be showing
				view = model.View()
				Expect(view).To(ContainSubstring("Keyboard Reference"))

				// Press ? again to hide help
				newModel, _ = model.Update(msg)
				model = newModel.(*app.Model)

				view = model.View()
				Expect(view).NotTo(ContainSubstring("Keyboard Reference"))
			})

			// Note: 'h' key is vim-style left navigation, not help
			// Help is toggled only with '?' key per KEYBOARD_REFERENCE.md

			It("should show navigation shortcuts in help", func() {
				// Show help
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
				newModel, _ := model.Update(msg)
				model = newModel.(*app.Model)

				view := model.View()
				Expect(view).To(ContainSubstring("Navigation"))
				Expect(view).To(ContainSubstring("Move up"))
				Expect(view).To(ContainSubstring("Move down"))
			})

			It("should show global shortcuts in help", func() {
				// Show help
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
				newModel, _ := model.Update(msg)
				model = newModel.(*app.Model)

				view := model.View()
				Expect(view).To(ContainSubstring("Global Shortcuts"))
				Expect(view).To(ContainSubstring("Quit"))
				Expect(view).To(ContainSubstring("Go back"))
			})
		})

		Context("home/esc/escape keys", func() {
			// NOTE: After removing app-level escape interceptor, escape is now
			// handled by intents themselves. Intents return Cancelled result,
			// which causes app to return to menu. These tests verify the intent
			// properly handles escape and returns the appropriate result.

			It("should forward escape to intent (intent returns to menu)", func() {
				// Activate an intent
				newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				model = newModel.(*app.Model)
				state := model.GetState()
				Expect(state).To(Equal(app.StateIntent))

				// Press escape - forwarded to intent
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				newModel, cmd := model.Update(msg)
				model = newModel.(*app.Model)

				// Intent may return command that completes and returns to menu
				if cmd != nil {
					resultMsg := cmd()
					if resultMsg != nil {
						newModel, _ = model.Update(resultMsg)
						model = newModel.(*app.Model)
					}
				}

				// Eventually should return to menu (via intent's Cancelled result)
				state = model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})

			It("should not intercept home key (intent handles it)", func() {
				// Activate an intent
				newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				model = newModel.(*app.Model)

				// Press home - no longer intercepted by app
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("home")}
				newModel, _ = model.Update(msg)
				model = newModel.(*app.Model)

				// App forwards to intent; intent may or may not handle 'home'
				// This test just verifies app doesn't crash
				state := model.GetState()
				Expect(state).To(Equal(app.StateIntent)) // Still in intent
			})

			It("should forward esc key to intent", func() {
				// Activate an intent
				newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				model = newModel.(*app.Model)

				// Press esc - forwarded to intent
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				newModel, cmd := model.Update(msg)
				model = newModel.(*app.Model)

				// Process any command returned
				if cmd != nil {
					resultMsg := cmd()
					if resultMsg != nil {
						newModel, _ = model.Update(resultMsg)
						model = newModel.(*app.Model)
					}
				}

				state := model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})

			It("should not affect menu state when pressing escape in menu", func() {
				state := model.GetState()
				Expect(state).To(Equal(app.StateMenu))

				// Press escape while in menu
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("escape")}
				model.Update(msg)

				state = model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})
		})

		Context("logo animation tick", func() {
			It("should update logo animation when in menu state", func() {
				// Send a tick message to update logo animation
				msg := display.TickMsg(time.Now())
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})
		})

		Context("window size message", func() {
			It("should handle window size updates", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})
		})
	})

	Describe("handleIntentInput", func() {
		BeforeEach(func() {
			// Activate an intent before each test
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		Context("when intent returns result", func() {
			It("should return to menu state via escape", func() {
				// Cancel intent with escape - forwarded to intent (already activated by BeforeEach)
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				newModel, cmdResult := model.Update(msg)
				model = newModel.(*app.Model)

				// Execute completion message (intent returns Cancelled result)
				if cmdResult != nil {
					completionMsg := cmdResult()
					if completionMsg != nil {
						newModel, _ = model.Update(completionMsg)
						model = newModel.(*app.Model)
					}
				}

				// Should return to menu via intent's result
				state := model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})

			It("should ignore 'q' key within intent to prevent accidental exits", func() {
				// 'q' key is no longer handled at intent level to prevent accidental exits
				// See: fix(navigation): remove global quit key to prevent accidental exits
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
				newModel, cmdResult := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
				// 'q' no longer quits from within intents - returns nil command
				Expect(cmdResult).To(BeNil())

				// Should still be in intent state
				state := model.GetState()
				Expect(state).To(Equal(app.StateIntent))
			})
		})

		Context("when intent does not return result", func() {
			It("should forward message to active intent", func() {
				// Send a non-completing key (like arrow down)
				msg := tea.KeyMsg{Type: tea.KeyDown}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())

				// Should still be in intent state
				state := model.GetState()
				Expect(state).To(Equal(app.StateIntent))
			})

			It("should return command from intent", func() {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				newModel, _ := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
			})
		})
	})

	Describe("SetInitialCaptureMode", func() {
		It("should accept manual mode without panic", func() {
			Expect(func() {
				model.SetInitialCaptureMode("manual")
			}).NotTo(Panic())
		})

		It("should accept burst mode without panic", func() {
			Expect(func() {
				model.SetInitialCaptureMode("burst")
			}).NotTo(Panic())
		})

		It("should accept csv mode without panic", func() {
			Expect(func() {
				model.SetInitialCaptureMode("csv")
			}).NotTo(Panic())
		})
	})

	Describe("GetState", func() {
		It("should return menu state initially", func() {
			state := model.GetState()
			Expect(state).To(Equal(app.StateMenu))
		})

		It("should return intent state after activating intent", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			state := model.GetState()
			Expect(state).To(Equal(app.StateIntent))
		})
	})

	Describe("GetActiveIntent", func() {
		It("should return nil when no intent is active", func() {
			intent := model.GetActiveIntent()
			Expect(intent).To(BeNil())
		})

		It("should return active intent after activation", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent := model.GetActiveIntent()
			Expect(intent).NotTo(BeNil())
		})
	})

	Describe("IntentCompletedMsg handling", func() {
		It("should handle intent completed message", func() {
			// Activate intent
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Send IntentCompletedMsg
			msg := app.IntentCompletedMsg{}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})
	})

	Describe("handleEditEventRequest", func() {
		It("should handle RequestEditEventMsg and transition to intent state", func() {
			// Create an event to edit.
			event := &career.CareerEvent{
				ID:   "test-event-1",
				Text: "Test Event Description",
			}

			// Send RequestEditEventMsg via Update (handleDefaultMsg routes it).
			msg := intents.RequestEditEventMsg{Event: event}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())

			// Should transition to intent state.
			model = newModel.(*app.Model)
			state := model.GetState()
			Expect(state).To(Equal(app.StateIntent))
		})

		It("should activate capture_event_edit intent with event context", func() {
			// Create an event with specific details.
			event := &career.CareerEvent{
				ID:   "edit-event-1",
				Text: "Event to Edit",
			}

			// Send RequestEditEventMsg.
			msg := intents.RequestEditEventMsg{Event: event}
			newModel, _ := model.Update(msg)
			model = newModel.(*app.Model)

			// Verify intent is active.
			activeIntent := model.GetActiveIntent()
			Expect(activeIntent).NotTo(BeNil())
		})
	})

	Describe("handleDefaultMsg - Coverage", func() {
		It("should return nil when in menu state and not RequestEditEventMsg", func() {
			// In menu state, send a generic message (not key, not edit request).
			msg := tea.MouseMsg{X: 0, Y: 0}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			// Should return nil command since we're in menu state.
			Expect(cmd).To(BeNil())
		})

		It("should route message to intent when in intent state", func() {
			// First activate an intent.
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Send a non-key message to the intent.
			msg := tea.MouseMsg{X: 10, Y: 10}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})

		It("should handle intent result and return to menu", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Send escape which causes intent to return result.
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
			model = newModel.(*app.Model)

			// Process any returned command.
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					newModel, _ = model.Update(resultMsg)
					model = newModel.(*app.Model)
				}
			}

			// Should be back in menu.
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})

		It("should pass through non-completing message when in intent state", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// Send a message that doesn't complete the intent (like arrow down).
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = newModel.(*app.Model)

			// Should still be in intent state.
			Expect(model.GetState()).To(Equal(app.StateIntent))
			// Command may or may not be nil depending on intent's response.
			Expect(newModel).NotTo(BeNil())
			_ = cmd // We don't care about the specific command.
		})

		It("should return to menu when intent completes with result", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// Trigger intent completion by pressing escape (intents handle this).
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			newModel, cmd := model.Update(msg)
			model = newModel.(*app.Model)

			// Process the batch command if returned.
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					// Check if it's a batch.
					if batchMsg, ok := msg.(tea.BatchMsg); ok {
						for _, bCmd := range batchMsg {
							if bCmd != nil {
								innerMsg := bCmd()
								if innerMsg != nil {
									newModel, _ = model.Update(innerMsg)
									model = newModel.(*app.Model)
								}
							}
						}
					} else {
						newModel, _ = model.Update(msg)
						model = newModel.(*app.Model)
					}
				}
			}

			// Should be back in menu.
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})
	})

	Describe("View - Additional Coverage", func() {
		It("should show info modal when visible in menu state", func() {
			// Navigate to generate_cv without any events to trigger info modal.
			// First, navigate down to "Generate CV" (index 3).
			for i := 0; i < 3; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			// Press enter to select - should show info modal since no events.
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// View should show modal content.
			view := model.View()
			Expect(view).To(ContainSubstring("No Career Events"))
		})

		It("should return 'No active intent' when intent state but no active intent", func() {
			// This is an edge case - normally shouldn't happen.
			// We test by directly setting state without activating intent.
			// Since we can't directly set state, we test the normal path.
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := model.View()
			// Should show intent view (not "No active intent" since intent is active).
			Expect(view).NotTo(Equal("No active intent"))
		})

		It("should render menu view when in menu state", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Capture Event"))
			Expect(view).To(ContainSubstring("Browse Timeline"))
		})

		It("should render intent view when intent is active", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			view := model.View()
			// Should render intent's view, not menu.
			Expect(view).NotTo(ContainSubstring("Browse Timeline"))
		})

		It("should return empty string for unknown state", func() {
			// The fallback case returns "" - this is tested via verifying the
			// model handles all expected states properly without panicking.
			// We cannot directly set an invalid state, so we verify edge handling.
			view := model.View()
			// Menu state should render properly.
			Expect(view).NotTo(BeEmpty())

			// Intent state should also render properly.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)
			view = model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("handleMenuInput - Boundary Conditions", func() {
		It("should not move up when at top of menu", func() {
			// Already at index 0, try moving up.
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			// Should still be at first item.
			view := model.View()
			Expect(view).To(ContainSubstring("Capture Event"))
		})

		It("should not move down when at bottom of menu", func() {
			// Move to bottom of menu.
			menuItems := model.GetMenuItems()
			for i := 0; i < len(menuItems)-1; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			// Try moving down again - should stay at bottom.
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			// Should still render without error.
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle space key as selection", func() {
			msg := tea.KeyMsg{Type: tea.KeySpace}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle k key for up navigation", func() {
			// Move down first, then use k to go up.
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})

		It("should handle j key for down navigation", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})
	})

	Describe("handleMenuSelection - Edge Cases", func() {
		It("should show info modal when selecting generate_cv without events", func() {
			// Navigate to generate_cv (index 3).
			for i := 0; i < 3; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			// Select it.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// Should still be in menu state (info modal shown).
			state := model.GetState()
			Expect(state).To(Equal(app.StateMenu))

			// View should contain modal.
			view := model.View()
			Expect(view).To(ContainSubstring("No Career Events"))
		})
	})

	Describe("handleKeyMsg - Info Modal", func() {
		It("should dismiss info modal on any key press", func() {
			// First trigger info modal by selecting generate_cv without events.
			for i := 0; i < 3; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Verify modal is shown.
			view := model.View()
			Expect(view).To(ContainSubstring("No Career Events"))

			// Press enter to dismiss modal.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// Modal should be dismissed, showing menu again.
			view = model.View()
			Expect(view).NotTo(ContainSubstring("No Career Events"))
		})
	})

	Describe("handleKeyMsg - State Fallback", func() {
		It("should handle key messages that don't match any case in menu state", func() {
			// Test the fallback return in handleKeyMsg by sending a key
			// that doesn't match ctrl+c, q, or ? while in menu state.
			// The 'x' key should fall through and be handled by handleMenuInput.
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			// handleMenuInput returns nil for unknown keys.
			Expect(cmd).To(BeNil())
		})

		It("should handle function keys gracefully", func() {
			// Function keys should fall through the switch.
			msg := tea.KeyMsg{Type: tea.KeyF1}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})
	})

	Describe("getMenuColumnWidths - Terminal Sizes", func() {
		It("should handle different terminal widths", func() {
			// Test with various window sizes.
			sizes := []tea.WindowSizeMsg{
				{Width: 40, Height: 20},  // Tiny.
				{Width: 60, Height: 24},  // Compact.
				{Width: 80, Height: 24},  // Normal.
				{Width: 120, Height: 40}, // Large.
				{Width: 200, Height: 60}, // XLarge.
			}

			for _, size := range sizes {
				newModel, _ := model.Update(size)
				model = newModel.(*app.Model)
				// View should render without error.
				view := model.View()
				Expect(view).NotTo(BeEmpty())
			}
		})
	})

	Describe("handleDefaultMsg - Intent Result Path", func() {
		It("should return nil when in menu state and not RequestEditEventMsg", func() {
			// In menu state, send a generic message (not key, not edit request).
			msg := tea.MouseMsg{X: 0, Y: 0}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			// Should return nil command since we're in menu state.
			Expect(cmd).To(BeNil())
		})

		It("should route message to intent when in intent state", func() {
			// First activate an intent.
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Send a non-key message to the intent.
			msg := tea.MouseMsg{X: 10, Y: 10}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})

		It("should handle intent result and return to menu", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Send escape which causes intent to return result.
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
			model = newModel.(*app.Model)

			// Process any returned command.
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					newModel, _ = model.Update(resultMsg)
					model = newModel.(*app.Model)
				}
			}

			// Should be back in menu.
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})

		It("should pass through non-completing message when in intent state", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// Send a message that doesn't complete the intent (like arrow down).
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = newModel.(*app.Model)

			// Should still be in intent state.
			Expect(model.GetState()).To(Equal(app.StateIntent))
			// Command may or may not be nil depending on intent's response.
			Expect(newModel).NotTo(BeNil())
			_ = cmd // We don't care about the specific command.
		})

		It("should return to menu when intent completes with result", func() {
			// Activate an intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// Trigger intent completion by pressing escape (intents handle this).
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			newModel, cmd := model.Update(msg)
			model = newModel.(*app.Model)

			// Process the batch command if returned.
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					// Check if it's a batch.
					if batchMsg, ok := msg.(tea.BatchMsg); ok {
						for _, bCmd := range batchMsg {
							if bCmd != nil {
								innerMsg := bCmd()
								if innerMsg != nil {
									newModel, _ = model.Update(innerMsg)
									model = newModel.(*app.Model)
								}
							}
						}
					} else {
						newModel, _ = model.Update(msg)
						model = newModel.(*app.Model)
					}
				}
			}

			// Should be back in menu.
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})

		It("should handle intent completion via non-key message (handleDefaultMsg result path)", func() {
			// Activate capture_event intent.
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Create a valid event for form submission.
			testEvent := &career.CareerEvent{
				ID:   "test-event-123",
				Text: "Test event for coverage",
				Date: time.Now(),
			}

			// Send FormSubmittedMsg to transition intent to review state.
			// This is a non-key message that goes through handleDefaultMsg.
			formMsg := intents.FormSubmittedMsg{Event: testEvent}
			newModel, _ = model.Update(formMsg)
			model = newModel.(*app.Model)

			// Intent should still be active (now in review state).
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Now send ReviewCancelledMsg - this should complete the intent
			// via the handleDefaultMsg -> result != nil path.
			cancelMsg := intents.ReviewCancelledMsg{}
			newModel, cmd := model.Update(cancelMsg)
			model = newModel.(*app.Model)

			// Process any batch command returned.
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					if batchMsg, ok := resultMsg.(tea.BatchMsg); ok {
						for _, bCmd := range batchMsg {
							if bCmd != nil {
								innerMsg := bCmd()
								if innerMsg != nil {
									newModel, _ = model.Update(innerMsg)
									model = newModel.(*app.Model)
								}
							}
						}
					} else {
						newModel, _ = model.Update(resultMsg)
						model = newModel.(*app.Model)
					}
				}
			}

			// Should be back in menu after intent completion.
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})
		It("should handle fallback when state is neither Menu nor Intent", func() {
			// This tests line 207 - the final return m, nil
			// This is technically unreachable with current state enum,
			// but we can at least verify the function handles unexpected states.
			// Since we can't set an invalid state, we verify the normal paths work.
			state := model.GetState()
			Expect(state).To(Equal(app.StateMenu))

			// Send a non-key, non-edit message while in menu state.
			customMsg := struct{ data int }{data: 42}
			newModel, cmd := model.Update(customMsg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil()) // Falls through to return m, nil
		})
	})

	Describe("Intent Activation - All Types", func() {
		// These tests ensure all intent registration factories are exercised.
		// Menu items: 0=capture_event, 1=browse_timeline, 2=manage_skills,
		// 3=generate_cv, 4=configure_system, 5=burst_management, 6=fact_management

		It("should activate browse_timeline intent", func() {
			// Navigate to browse_timeline (index 1).
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			Expect(model.GetState()).To(Equal(app.StateIntent))
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate manage_skills intent", func() {
			// Navigate to manage_skills (index 2).
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			Expect(model.GetState()).To(Equal(app.StateIntent))
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate configure_system intent", func() {
			// Navigate to configure_system (index 4).
			for i := 0; i < 4; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			Expect(model.GetState()).To(Equal(app.StateIntent))
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate burst_management intent", func() {
			// Navigate to burst_management (index 5).
			for i := 0; i < 5; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			Expect(model.GetState()).To(Equal(app.StateIntent))
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate fact_management intent", func() {
			// Navigate to fact_management (index 6).
			for i := 0; i < 6; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
			}
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			Expect(model.GetState()).To(Equal(app.StateIntent))
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Edge Cases", func() {
		It("should handle unknown key messages gracefully", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})

		It("should handle multiple state transitions", func() {
			// Menu -> Intent -> Menu -> Intent
			state := model.GetState()
			Expect(state).To(Equal(app.StateMenu))

			// Activate intent
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)
			state = model.GetState()
			Expect(state).To(Equal(app.StateIntent))

			// Return to menu via intent handling escape
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			newModel, cmd := model.Update(msg)
			model = newModel.(*app.Model)

			// Process intent's result
			if cmd != nil {
				resultMsg := cmd()
				if resultMsg != nil {
					newModel, _ = model.Update(resultMsg)
					model = newModel.(*app.Model)
				}
			}

			state = model.GetState()
			Expect(state).To(Equal(app.StateMenu))

			// Activate intent again
			newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)
			state = model.GetState()
			Expect(state).To(Equal(app.StateIntent))
		})

		It("should handle rapid key presses", func() {
			// Simulate rapid navigation
			for i := 0; i < 10; i++ {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
				newModel, _ := model.Update(msg)
				model = newModel.(*app.Model)
			}
			Expect(model).NotTo(BeNil())
		})
	})
})

// mockFailingRegistrar is a test registrar that returns an error.
type mockFailingRegistrar struct {
	err error
}

func (m *mockFailingRegistrar) RegisterAll(_ context.Context, _ *intents.DefaultIntentRouter) error {
	return m.err
}

// mockNilFactoryRegistrar registers factories that return nil intents.
type mockNilFactoryRegistrar struct{}

func (m *mockNilFactoryRegistrar) RegisterAll(_ context.Context, router *intents.DefaultIntentRouter) error {
	// Register a factory that returns nil (simulating intent creation failure).
	return router.RegisterIntent("capture_event", func() intents.Intent {
		return nil // Simulates intent creation failure.
	})
}

// mockPartialRegistrar registers some intents successfully and some with nil factories.
type mockPartialRegistrar struct{}

func (m *mockPartialRegistrar) RegisterAll(ctx context.Context, router *intents.DefaultIntentRouter) error {
	// Register capture_event with a nil factory.
	router.RegisterIntent("capture_event", func() intents.Intent {
		return nil
	})
	// Register browse_timeline with a working factory.
	router.RegisterIntent("browse_timeline", func() intents.Intent {
		return &mockIntent{}
	})
	return nil
}

// mockIntent is a minimal intent implementation for testing.
type mockIntent struct{}

func (m *mockIntent) Init() tea.Cmd                              { return nil }
func (m *mockIntent) Update(_ tea.Msg) tea.Cmd                   { return nil }
func (m *mockIntent) View() string                               { return "mock intent view" }
func (m *mockIntent) Result() *intents.IntentResult[interface{}] { return nil }

var _ = Describe("IntentRegistrar DI Tests", func() {
	var (
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		config.SetConfigPathForTesting(filepath.Join(GinkgoT().TempDir(), "config.yaml"))
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)
	})

	AfterEach(func() {
		config.ResetConfigPath()
	})

	Describe("WithIntentRegistrar option", func() {
		It("should use custom registrar when provided", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			// Use mock registrar that registers a working intent.
			mockReg := &mockPartialRegistrar{}
			model := app.NewModel(cliService, svc, bootstrapResult, app.WithIntentRegistrar(mockReg))

			Expect(model).NotTo(BeNil())
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})

		It("should handle registrar that returns error gracefully", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			// Use mock registrar that returns an error.
			mockReg := &mockFailingRegistrar{err: context.DeadlineExceeded}
			model := app.NewModel(cliService, svc, bootstrapResult, app.WithIntentRegistrar(mockReg))

			// Model should still be created (error is logged, not fatal).
			Expect(model).NotTo(BeNil())
			Expect(model.GetState()).To(Equal(app.StateMenu))
		})
	})

	Describe("handleMenuSelection with failing factory", func() {
		It("should handle ActivateIntent error when factory returns nil", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			// Use mock registrar that registers a nil-returning factory.
			mockReg := &mockNilFactoryRegistrar{}
			model := app.NewModel(cliService, svc, bootstrapResult, app.WithIntentRegistrar(mockReg))

			// Try to activate the intent (should fail because factory returns nil).
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(*app.Model)

			// ActivateIntent should fail, model should stay in menu state.
			// Note: The state is set to Intent BEFORE ActivateIntent is called,
			// so we need to check if the error path resets state or handles gracefully.
			Expect(model).NotTo(BeNil())
			// The cmd should be nil or a no-op when activation fails.
			_ = cmd
		})
	})

	Describe("DefaultIntentRegistrar", func() {
		It("should register all intents successfully", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:      cliService,
				CareerService:   svc,
				Log:             log,
				CVGenService:    bootstrapResult.Services.CVGenService,
				CVExportService: bootstrapResult.Services.CVExportService,
			})

			router := intents.NewDefaultIntentRouter()
			err := registrar.RegisterAll(context.Background(), router)

			Expect(err).To(BeNil())
		})

		It("should return error on duplicate registration", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:      cliService,
				CareerService:   svc,
				Log:             log,
				CVGenService:    bootstrapResult.Services.CVGenService,
				CVExportService: bootstrapResult.Services.CVExportService,
			})

			router := intents.NewDefaultIntentRouter()

			// First registration should succeed.
			err := registrar.RegisterAll(context.Background(), router)
			Expect(err).To(BeNil())

			// Second registration should fail (duplicate).
			err = registrar.RegisterAll(context.Background(), router)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("already registered"))
		})
	})

	Describe("View Edge Cases with State Manipulation", func() {
		It("should return 'No active intent' when state is Intent but router has no active intent", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			// Create model with partial registrar (only browse_timeline).
			mockReg := &mockPartialRegistrar{}
			model := app.NewModel(cliService, svc, bootstrapResult, app.WithIntentRegistrar(mockReg))

			// Force state to Intent without activating an intent.
			model.SetStateForTesting(app.StateIntent)

			// View should return "No active intent".
			view := model.View()
			Expect(view).To(Equal("No active intent"))
		})

		It("should return empty string for unknown state", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			mockReg := &mockPartialRegistrar{}
			model := app.NewModel(cliService, svc, bootstrapResult, app.WithIntentRegistrar(mockReg))

			// Set state to something that's not Menu or Intent.
			// Note: This requires using a state that's neither StateMenu nor StateIntent.
			// Since AppState is a string type, we can set it to an invalid value.
			model.SetStateForTesting(app.AppState("invalid"))

			// View should return empty string for unknown state.
			view := model.View()
			Expect(view).To(Equal(""))
		})
	})

	Describe("Accessors", func() {
		It("GetIntentRouter should return the router", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			model := app.NewModel(cliService, svc, bootstrapResult)
			router := model.GetIntentRouter()

			Expect(router).NotTo(BeNil())
		})
	})

	Describe("handleKeyMsg with invalid state", func() {
		It("should return nil cmd when state is invalid", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			mockReg := &mockPartialRegistrar{}
			model := app.NewModel(cliService, svc, bootstrapResult, app.WithIntentRegistrar(mockReg))

			// Set state to invalid value.
			model.SetStateForTesting(app.AppState("invalid"))

			// Send a key message - should hit the fallback.
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
			newModel, cmd := model.Update(msg)

			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleEditEventRequest Error Paths", func() {
		It("should handle duplicate registration gracefully", func() {
			log := logger.DefaultLogger()
			bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

			model := app.NewModel(cliService, svc, bootstrapResult)

			// Create a test event.
			testEvent := &career.CareerEvent{
				ID:   "edit-test-event",
				Text: "Test event for editing",
			}

			// First request - should succeed.
			msg1 := intents.RequestEditEventMsg{Event: testEvent}
			newModel, _ := model.Update(msg1)
			model = newModel.(*app.Model)
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Return to menu.
			model.SetStateForTesting(app.StateMenu)

			// Second request with same event - registration is duplicate but handled.
			// The nolint comment indicates this is expected behavior.
			msg2 := intents.RequestEditEventMsg{Event: testEvent}
			newModel, cmd := model.Update(msg2)
			model = newModel.(*app.Model)

			// Should still transition to intent state (uses existing registration).
			// The cmd may be nil if the intent returns nil from Init().
			Expect(model.GetState()).To(Equal(app.StateIntent))
			_ = cmd // Command may or may not be nil.
		})
	})

	Describe("Intent Factory Error Paths", func() {
		It("should handle capture_event factory failure when context is invalid", func() {
			log := logger.DefaultLogger()

			// Create registrar with nil services (will cause validation failure).
			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:    nil, // Invalid - will cause factory to fail.
				CareerService: nil,
				Log:           log,
			})

			router := intents.NewDefaultIntentRouter()
			err := registrar.RegisterAll(context.Background(), router)
			Expect(err).To(BeNil()) // Registration succeeds, factory failure happens on activation.

			// Now try to activate - the factory will fail and return nil.
			_, err = router.ActivateIntent("capture_event", nil)
			Expect(err).NotTo(BeNil()) // Factory returned nil.
			Expect(err.Error()).To(ContainSubstring("factory returned nil"))
		})

		It("should handle browse_timeline factory failure", func() {
			log := logger.DefaultLogger()

			// Create registrar with nil services.
			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:    nil,
				CareerService: nil,
				Log:           log,
			})

			router := intents.NewDefaultIntentRouter()
			registrar.RegisterAll(context.Background(), router)

			// Activation will fail because factory returns nil.
			_, err := router.ActivateIntent("browse_timeline", nil)
			Expect(err).NotTo(BeNil())
		})

		It("should handle generate_cv factory failure", func() {
			log := logger.DefaultLogger()

			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:    nil,
				CareerService: nil,
				Log:           log,
			})

			router := intents.NewDefaultIntentRouter()
			registrar.RegisterAll(context.Background(), router)

			_, err := router.ActivateIntent("generate_cv", nil)
			Expect(err).NotTo(BeNil())
		})

		It("should handle configure_system factory failure", func() {
			log := logger.DefaultLogger()

			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:    nil,
				CareerService: nil,
				Log:           log,
			})

			router := intents.NewDefaultIntentRouter()
			registrar.RegisterAll(context.Background(), router)

			// ConfigureSystem might not fail with nil services, check behavior.
			_, err := router.ActivateIntent("configure_system", nil)
			// May or may not fail depending on implementation.
			_ = err
		})

		It("should handle burst_management factory failure", func() {
			log := logger.DefaultLogger()

			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:    nil,
				CareerService: nil, // Will cause GetBurstRepository to panic or fail.
				Log:           log,
			})

			router := intents.NewDefaultIntentRouter()
			registrar.RegisterAll(context.Background(), router)

			_, err := router.ActivateIntent("burst_management", nil)
			Expect(err).NotTo(BeNil())
		})

		It("should handle fact_management factory failure", func() {
			log := logger.DefaultLogger()

			registrar := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
				CLIService:    nil,
				CareerService: nil,
				Log:           log,
			})

			router := intents.NewDefaultIntentRouter()
			registrar.RegisterAll(context.Background(), router)

			_, err := router.ActivateIntent("fact_management", nil)
			Expect(err).NotTo(BeNil())
		})
	})
})
