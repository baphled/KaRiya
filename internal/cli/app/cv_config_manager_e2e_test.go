package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVConfigManager E2E - User Gets Stuck on Loading Screen", func() {
	Describe("Bug: Loading configuration templates screen never completes", func() {
		It("should not get stuck on 'Loading configuration templates...' when configs exist", func() {
			// Setup: Create a config manager with configs already available
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			// Pre-populate with some configs
			configs := []*career.CVConfig{
				{
					Name:           "Principal Engineer",
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				},
				{
					Name:           "Staff Engineer",
					TargetRole:     "staff",
					TargetAudience: []string{"recruiter"},
				},
			}

			for _, config := range configs {
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			// Create model
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Check initial loading state
			view := cvModel.View()
			Expect(view).To(ContainSubstring("Loading configuration templates"), "Should show loading message initially")

			// Execute loading
			cmd := cvModel.Init()
			Expect(cmd).NotTo(BeNil(), "Init() must return a command to load configs")

			// Command should complete without hanging
			done := make(chan tea.Msg, 1)
			go func() {
				msg := cmd()
				done <- msg
			}()

			var loadMsg tea.Msg
			select {
			case loadMsg = <-done:
				Expect(loadMsg).NotTo(BeNil())
			case <-time.After(5 * time.Second):
				Fail("Command execution timed out - this is the bug!")
			}

			// Process message
			updatedModel, _ := cvModel.Update(loadMsg)
			Expect(updatedModel).NotTo(BeNil())

			// After processing, should not be loading anymore
			finalView := updatedModel.View()
			Expect(finalView).NotTo(ContainSubstring("Loading configuration templates"),
				"BUG: Still showing loading message - configs never loaded!")

			// Should show actual configs
			Expect(finalView).To(ContainSubstring("Principal Engineer"),
				"BUG: Configs not displayed - user stuck on loading screen!")
			Expect(finalView).To(ContainSubstring("Staff Engineer"),
				"Configs should be visible in the view")

			// User should be able to interact
			keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
			_, escapeCmd := updatedModel.Update(keyMsg)
			Expect(escapeCmd).NotTo(BeNil(), "User should be able to escape")
		})

		It("should render something useful even if configs fail to load", func() {
			// Empty config manager (simulating no configs available)
			configManager := cv.NewMemoryConfigManager()

			// Create model
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Execute the loading flow
			cmd := cvModel.Init()
			Expect(cmd).NotTo(BeNil())

			loadMsg := cmd()
			Expect(loadMsg).NotTo(BeNil())

			updatedModel, _ := cvModel.Update(loadMsg)

			// Should not be stuck on loading
			view := updatedModel.View()
			Expect(view).NotTo(ContainSubstring("Loading configuration templates..."),
				"BUG: Still stuck on loading screen!")

			// Should show something useful
			Expect(view).NotTo(BeEmpty(), "View should render something")

			// User can interact
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, cmd = updatedModel.Update(keyMsg)
			Expect(cmd).NotTo(BeNil(), "Should respond to user input")
		})

		It("should handle timeout gracefully without UI freeze", func() {
			// Setup
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Execute loading with timeout monitoring
			cmd := cvModel.Init()

			startTime := time.Now()
			done := make(chan tea.Msg, 1)
			go func() {
				msg := cmd()
				done <- msg
			}()

			// Command should complete in reasonable time
			select {
			case msg := <-done:
				duration := time.Since(startTime)
				Expect(duration < 5*time.Second).To(BeTrue(),
					"BUG: Loading took too long - UI would freeze!")
				Expect(msg).NotTo(BeNil())

				// Should be able to process the message
				updatedModel, _ := cvModel.Update(msg)
				Expect(updatedModel).NotTo(BeNil())

				// View should be responsive
				view := updatedModel.View()
				Expect(view).NotTo(BeEmpty())

			case <-time.After(10 * time.Second):
				Fail("BUG: Command execution timed out - UI is frozen!")
			}
		})

		It("should allow escape key to work during loading", func() {
			// Setup
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Model is in loading state
			view := cvModel.View()
			Expect(view).To(ContainSubstring("Loading configuration templates"))

			// User presses escape while loading
			keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
			_, cmd := cvModel.Update(keyMsg)

			// BUG: If this fails, user is stuck
			Expect(cmd).NotTo(BeNil(),
				"BUG: Escape key doesn't work during loading - user is stuck!")

			// Command should do something (navigate away)
			msg := cmd()
			Expect(msg).NotTo(BeNil())
		})

		It("should complete full flow without getting stuck", func() {
			// Setup
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			// Add test configs
			for i := 1; i <= 2; i++ {
				config := &career.CVConfig{
					Name:           "test-config-" + string(rune(48+i)),
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				}
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			// Create model
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Check initial loading state
			view := cvModel.View()
			Expect(view).To(ContainSubstring("Loading configuration templates"))

			// Execute loading
			cmd := cvModel.Init()
			msg := cmd()

			// Process message
			updatedModel, _ := cvModel.Update(msg)

			// Check final state
			finalView := updatedModel.View()

			// Should not be stuck on loading
			Expect(finalView).NotTo(ContainSubstring("Loading configuration templates"),
				"BUG: Still stuck on loading screen!")

			// Should show configs
			Expect(finalView).To(ContainSubstring("test-config-1"),
				"BUG: Config 1 not visible!")
			Expect(finalView).To(ContainSubstring("test-config-2"),
				"BUG: Config 2 not visible!")

			// User can navigate
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			navModel, _ := updatedModel.Update(keyMsg)
			Expect(navModel).NotTo(BeNil())

			// User can escape
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			_, escapeCmd := navModel.Update(keyMsg)
			Expect(escapeCmd).NotTo(BeNil())
		})
	})
})

var _ = Describe("CVConfigManager Keyboard Shortcuts E2E", func() {
	Describe("'n' key - Create new configuration", func() {
		It("should trigger NavigateToScreenMsg when 'n' is pressed", func() {
			// Setup
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'n' to create new config
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, resultCmd := updatedModel.Update(keyMsg)

			// Should return a command
			Expect(resultCmd).NotTo(BeNil(), "BUG: 'n' key doesn't trigger navigation command")

			// Execute the command
			resultMsg := resultCmd()
			Expect(resultMsg).NotTo(BeNil())

			// Should be NavigateToScreenMsg
			navMsg, ok := resultMsg.(models.NavigateToScreenMsg)
			Expect(ok).To(BeTrue(), "BUG: 'n' key should trigger NavigateToScreenMsg")
			Expect(navMsg.ScreenID).To(Equal("cv_config_editor"), "Should navigate to CV config editor")
		})

		It("should work even with empty config list", func() {
			// Setup with no configs
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'n' on empty list
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, resultCmd := updatedModel.Update(keyMsg)

			Expect(resultCmd).NotTo(BeNil(), "BUG: 'n' should work on empty list")

			resultMsg := resultCmd()
			navMsg, ok := resultMsg.(models.NavigateToScreenMsg)
			Expect(ok).To(BeTrue())
			Expect(navMsg.ScreenID).To(Equal("cv_config_editor"))
		})
	})

	Describe("'e' key - Edit selected configuration", func() {
		It("should trigger EditCVConfigMsg when 'e' is pressed with selection", func() {
			// Setup with configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			testConfig := &career.CVConfig{
				Name:           "Test Config",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			err := configManager.SaveConfig(ctx, testConfig)
			Expect(err).NotTo(HaveOccurred())

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'e' to edit selected config
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, resultCmd := updatedModel.Update(keyMsg)

			Expect(resultCmd).NotTo(BeNil(), "BUG: 'e' key doesn't trigger edit command")

			resultMsg := resultCmd()
			Expect(resultMsg).NotTo(BeNil())

			// Should be EditCVConfigMsg
			editMsg, ok := resultMsg.(models.EditCVConfigMsg)
			Expect(ok).To(BeTrue(), "BUG: 'e' key should trigger EditCVConfigMsg")
			Expect(editMsg.Config.Name).To(Equal("Test Config"))
		})

		It("should not trigger EditCVConfigMsg when 'e' is pressed with empty list", func() {
			// Setup with no configs
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'e' on empty list
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, resultCmd := updatedModel.Update(keyMsg)

			// Should return nil (no action)
			Expect(resultCmd).To(BeNil(), "Should not trigger action on empty list")
		})

		It("should edit the correct config when multiple exist", func() {
			// Setup with multiple configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			configs := []*career.CVConfig{
				{
					Name:           "Config 1",
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				},
				{
					Name:           "Config 2",
					TargetRole:     "staff",
					TargetAudience: []string{"recruiter"},
				},
				{
					Name:           "Config 3",
					TargetRole:     "em",
					TargetAudience: []string{"hiring_manager"},
				},
			}

			for _, config := range configs {
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Get the initially selected config
			configModel := updatedModel.(*models.CVConfigManagerModel)
			initialConfig := configModel.GetSelectedConfig()
			Expect(initialConfig).NotTo(BeNil())

			// Press 'e' to edit the selected config
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, resultCmd := updatedModel.Update(keyMsg)

			Expect(resultCmd).NotTo(BeNil(), "Edit command should be returned")

			resultMsg := resultCmd()
			editMsg, ok := resultMsg.(models.EditCVConfigMsg)
			Expect(ok).To(BeTrue(), "'e' key should trigger EditCVConfigMsg")
			// Should be editing the initially selected config
			Expect(editMsg.Config.Name).To(Equal(initialConfig.Name), "Should edit the selected config")
		})
	})

	Describe("'enter' key - Generate CV from configuration", func() {
		It("should trigger GenerateCVFromConfigMsg when 'enter' is pressed", func() {
			// Setup with configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			testConfig := &career.CVConfig{
				Name:           "Test Config",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			err := configManager.SaveConfig(ctx, testConfig)
			Expect(err).NotTo(HaveOccurred())

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'enter' to generate CV
			keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
			_, resultCmd := updatedModel.Update(keyMsg)

			Expect(resultCmd).NotTo(BeNil(), "BUG: 'enter' key doesn't trigger generate command")

			resultMsg := resultCmd()
			Expect(resultMsg).NotTo(BeNil())

			// Should be GenerateCVFromConfigMsg
			genMsg, ok := resultMsg.(models.GenerateCVFromConfigMsg)
			Expect(ok).To(BeTrue(), "BUG: 'enter' should trigger GenerateCVFromConfigMsg")
			Expect(genMsg.Config.Name).To(Equal("Test Config"))
		})

		It("should not trigger GenerateCVFromConfigMsg when 'enter' is pressed with empty list", func() {
			// Setup with no configs
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'enter' on empty list
			keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
			_, resultCmd := updatedModel.Update(keyMsg)

			// Should return nil (no action)
			Expect(resultCmd).To(BeNil(), "Should not trigger action on empty list")
		})
	})

	Describe("'d' key - Delete configuration", func() {
		It("should show confirmation dialog when 'd' is pressed", func() {
			// Setup with configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			testConfig := &career.CVConfig{
				Name:           "Test Config",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			err := configManager.SaveConfig(ctx, testConfig)
			Expect(err).NotTo(HaveOccurred())

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'd' to delete
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			resultModel, _ := updatedModel.Update(keyMsg)

			// Cast to check deletion state
			configModel, ok := resultModel.(*models.CVConfigManagerModel)
			Expect(ok).To(BeTrue())
			Expect(configModel.DeletionState.IsConfirming()).To(BeTrue(), "BUG: 'd' should show confirmation dialog")
		})

		It("should not delete on empty list", func() {
			// Setup with no configs
			configManager := cv.NewMemoryConfigManager()
			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Press 'd' on empty list
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			resultModel, _ := updatedModel.Update(keyMsg)

			configModel, _ := resultModel.(*models.CVConfigManagerModel)
			Expect(configModel.DeletionState.IsConfirming()).To(BeFalse(), "Should not show confirmation on empty list")
		})
	})

	Describe("Navigation keys", func() {
		It("should navigate up with 'k' key", func() {
			// Setup with multiple configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			for i := 1; i <= 3; i++ {
				config := &career.CVConfig{
					Name:           "Config " + string(rune(48+i)),
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				}
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Navigate down first
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			model1, _ := updatedModel.Update(keyMsg)
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			model2, _ := model1.Update(keyMsg)

			// Then navigate up
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			model3, _ := model2.Update(keyMsg)

			Expect(model3).NotTo(BeNil())
		})

		It("should go to home with 'g' key", func() {
			// Setup with multiple configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			for i := 1; i <= 5; i++ {
				config := &career.CVConfig{
					Name:           "Config " + string(rune(48+i)),
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				}
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Navigate down multiple times
			for i := 0; i < 3; i++ {
				keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				updatedModel, _ = updatedModel.Update(keyMsg)
			}

			// Then go to home
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
			model, _ := updatedModel.Update(keyMsg)

			Expect(model).NotTo(BeNil())
		})

		It("should go to end with 'G' key", func() {
			// Setup with multiple configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			for i := 1; i <= 5; i++ {
				config := &career.CVConfig{
					Name:           "Config " + string(rune(48+i)),
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				}
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize and load
			cmd := cvModel.Init()
			msg := cmd()
			updatedModel, _ := cvModel.Update(msg)

			// Go to end
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
			model, _ := updatedModel.Update(keyMsg)

			Expect(model).NotTo(BeNil())
		})
	})

	Describe("Full user workflow", func() {
		It("should support complete user interaction flow", func() {
			// Setup with configs
			configManager := cv.NewMemoryConfigManager()
			ctx := context.Background()

			configs := []*career.CVConfig{
				{
					Name:           "Principal Config",
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				},
				{
					Name:           "Staff Config",
					TargetRole:     "staff",
					TargetAudience: []string{"recruiter"},
				},
			}

			for _, config := range configs {
				err := configManager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			baseModel := models.NewBaseStandardModel()
			cvModel := models.NewCVConfigManagerModel(baseModel, configManager)

			// Initialize
			cmd := cvModel.Init()
			msg := cmd()
			model, _ := cvModel.Update(msg)

			// User sees loaded configs
			view := model.View()
			Expect(view).To(ContainSubstring("Principal Config"))
			Expect(view).To(ContainSubstring("Staff Config"))

			// User can navigate
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			model, _ = model.Update(keyMsg)

			// User can press 'e' to edit (command should be returned)
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, resultCmd := model.Update(keyMsg)
			Expect(resultCmd).NotTo(BeNil(), "Edit command should be returned")

			// User can also press 'n' to create new (command should be returned)
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, newCmd := model.Update(keyMsg)
			Expect(newCmd).NotTo(BeNil(), "New config command should be returned")

			// User can escape
			keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
			_, escapeCmd := model.Update(keyMsg)
			Expect(escapeCmd).NotTo(BeNil(), "Escape should return a command")
		})
	})
})
