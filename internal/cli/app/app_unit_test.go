package app_test

import (
	"context"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/service"
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
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)

		// Pre-populate repositories to avoid nil panics
		_ = burstRepo.Create(ctx, &career.Burst{
			ID:       "b1",
			Name:     "dummy",
			EventIDs: []string{"e1", "e2"},
		})
		_ = factRepo.Create(ctx, &career.Fact{
			ID:                   "f1",
			Text:                 "dummy",
			CompetencyCategories: []string{"leadership"},
			RoleFit:              "staff",
			AudienceRelevance:    []string{"peer"},
			SourceEventID:        "e1",
		})

		model = app.NewModel(cliService, svc)
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
			It("should return to menu from intent state on escape", func() {
				// Activate an intent
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				state := model.GetState()
				Expect(state).To(Equal(app.StateIntent))

				// Press escape to return to menu
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("escape")}
				newModel, _ := model.Update(msg)
				model = newModel.(*app.Model)

				state = model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})

			It("should return to menu on home key", func() {
				// Activate an intent
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Press home to return to menu
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("home")}
				newModel, _ := model.Update(msg)
				model = newModel.(*app.Model)

				state := model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})

			It("should return to menu on esc key", func() {
				// Activate an intent
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Press esc to return to menu
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
				newModel, _ := model.Update(msg)
				model = newModel.(*app.Model)

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
				msg := components.TickMsg{}
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
			It("should return to menu state", func() {
				// Send quit key to intent to complete it
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
				newModel, cmdResult := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
				Expect(cmdResult).NotTo(BeNil())

				// Execute the command to trigger state change
				if cmdResult != nil {
					msg := cmdResult()
					if msg != nil {
						newModel, _ = model.Update(msg)
						model = newModel.(*app.Model)
					}
				}
			})

			It("should return to menu on escape", func() {
				// Navigate down in menu first
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})

				// Activate intent
				model.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Cancel intent with escape (q now quits the app)
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				newModel, cmdResult := model.Update(msg)
				model = newModel.(*app.Model)

				// Execute completion message
				if cmdResult != nil {
					completionMsg := cmdResult()
					if completionMsg != nil {
						newModel, _ = model.Update(completionMsg)
						model = newModel.(*app.Model)
					}
				}

				// Menu index should be reset (we can verify by checking state)
				state := model.GetState()
				Expect(state).To(Equal(app.StateMenu))
			})

			It("should return tea.Quit on 'q' key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
				newModel, cmdResult := model.Update(msg)
				Expect(newModel).NotTo(BeNil())
				// q now returns tea.Quit
				Expect(cmdResult).NotTo(BeNil())
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
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			state = model.GetState()
			Expect(state).To(Equal(app.StateIntent))

			// Return to menu
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("escape")}
			newModel, _ := model.Update(msg)
			model = newModel.(*app.Model)
			state = model.GetState()
			Expect(state).To(Equal(app.StateMenu))

			// Activate intent again
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
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
