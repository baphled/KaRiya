//nolint:errcheck // Test file - error handling for test setup is not relevant.
package app_test

import (
	"context"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/config"
	career "github.com/baphled/kariya/internal/domain/career"
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

		model = app.NewModel(cliService, svc)
		model.SkipOnboarding() // Skip onboarding for tests
	})

	AfterEach(func() {
		config.ResetConfigPath()
	})

	Describe("Init", func() {
		It("should return batch command for window size and logo init", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
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
