package models

import (
	"context"
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCVConfigManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CVConfigManager Suite")
}

var _ = Describe("CVConfigManagerModel", func() {
	var (
		model         *CVConfigManagerModel
		baseModel     *BaseStandardModel
		configManager cv.ConfigManager
		ctx           context.Context
		testConfigs   []*career.CVConfig
	)

	BeforeEach(func() {
		ctx = context.Background()
		baseModel = NewBaseStandardModel()
		configManager = cv.NewMemoryConfigManager()
		model = NewCVConfigManagerModel(baseModel, configManager)

		// Create test configurations
		testConfigs = []*career.CVConfig{
			{
				Name:            "Senior IC CV",
				TargetRole:      "senior_ic",
				TargetAudience:  []string{"hiring_manager", "recruiter"},
				EventFilters:    map[string]interface{}{},
				CreatedAt:       time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				Name:            "Staff Engineer CV",
				TargetRole:      "staff",
				TargetAudience:  []string{"hiring_manager"},
				EventFilters:    map[string]interface{}{},
				CreatedAt:       time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-20 * 24 * time.Hour),
			},
			{
				Name:            "Engineering Manager CV",
				TargetRole:      "em",
				TargetAudience:  []string{"recruiter", "peer"},
				EventFilters:    map[string]interface{}{},
				CreatedAt:       time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-10 * 24 * time.Hour),
			},
		}

		// Save test configurations
		for _, config := range testConfigs {
			err := configManager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("Initialization", func() {
		It("should initialize with empty configs and loading state", func() {
			newModel := NewCVConfigManagerModel(baseModel, configManager)

			Expect(newModel.configs).To(BeEmpty())
			Expect(newModel.loading).To(BeTrue())
			Expect(newModel.selectedIdx).To(Equal(0))
		})

		It("should load configs on Init", func() {
			cmd := model.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("ConfigsLoadedMsg", func() {
		It("should update configs and set loading to false", func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			newModel, _ := model.Update(msg)

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.loading).To(BeFalse())
			Expect(castedModel.configs).To(HaveLen(3))
		})

		It("should select first config after loading", func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			newModel, _ := model.Update(msg)

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(0))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			updatedModel, _ := model.Update(msg)
			model = updatedModel.(*CVConfigManagerModel)
		})

		It("should move down through configs with 'j'", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(1))
		})

		It("should move up through configs with 'k'", func() {
			model.selectedIdx = 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(0))
		})

		It("should move down with arrow key", func() {
			downKey := tea.KeyMsg{Type: tea.KeyDown}
			newModel, _ := model.Update(downKey)

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(1))
		})

		It("should move up with arrow key", func() {
			model.selectedIdx = 1
			upKey := tea.KeyMsg{Type: tea.KeyUp}
			newModel, _ := model.Update(upKey)

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(0))
		})

		It("should not move past last config", func() {
			model.selectedIdx = len(testConfigs) - 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(len(testConfigs) - 1))
		})

		It("should not move before first config", func() {
			model.selectedIdx = 0
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.selectedIdx).To(Equal(0))
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			updatedModel, _ := model.Update(msg)
			model = updatedModel.(*CVConfigManagerModel)
		})

		It("should trigger GenerateCVFromConfigMsg on Enter", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(cmd).NotTo(BeNil())
			// Execute the command to get the message
			resultMsg := cmd()
			Expect(resultMsg).To(BeAssignableToTypeOf(GenerateCVFromConfigMsg{}))

			genMsg := resultMsg.(GenerateCVFromConfigMsg)
			Expect(genMsg.config).To(Equal(testConfigs[0]))
		})

		It("should trigger NavigateToScreenMsg on 'n'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			Expect(cmd).NotTo(BeNil())
			resultMsg := cmd()
			Expect(resultMsg).To(BeAssignableToTypeOf(NavigateToScreenMsg{}))

			navMsg := resultMsg.(NavigateToScreenMsg)
			Expect(navMsg.screenID).To(Equal("cv_config_editor"))
		})

		It("should trigger EditCVConfigMsg on 'e'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			Expect(cmd).NotTo(BeNil())
			resultMsg := cmd()
			Expect(resultMsg).To(BeAssignableToTypeOf(EditCVConfigMsg{}))

			editMsg := resultMsg.(EditCVConfigMsg)
			Expect(editMsg.config).To(Equal(testConfigs[0]))
		})

		It("should trigger ConfirmDeleteCVConfigMsg on 'd'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			Expect(cmd).NotTo(BeNil())
			resultMsg := cmd()
			Expect(resultMsg).To(BeAssignableToTypeOf(ConfirmDeleteCVConfigMsg{}))

			delMsg := resultMsg.(ConfirmDeleteCVConfigMsg)
			Expect(delMsg.config).To(Equal(testConfigs[0]))
		})

		It("should trigger BackToMainMenuMsg on Esc", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(cmd).NotTo(BeNil())
			resultMsg := cmd()
			Expect(resultMsg).To(BeAssignableToTypeOf(BackToMainMenuMsg{}))
		})

		It("should trigger BackToMainMenuMsg on 'q'", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(cmd).NotTo(BeNil())
			resultMsg := cmd()
			Expect(resultMsg).To(BeAssignableToTypeOf(BackToMainMenuMsg{}))
		})
	})

	Describe("GetSelectedConfig", func() {
		It("should return selected config", func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			updatedModel, _ := model.Update(msg)
			castedModel := updatedModel.(*CVConfigManagerModel)

			selected := castedModel.GetSelectedConfig()

			Expect(selected).NotTo(BeNil())
			Expect(selected).To(Equal(testConfigs[0]))
		})

		It("should return nil when no configs loaded", func() {
			selected := model.GetSelectedConfig()

			Expect(selected).To(BeNil())
		})
	})

	Describe("View", func() {
		It("should show loading message when loading", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Loading configurations"))
		})

		It("should show empty message when no configs", func() {
			msg := ConfigsLoadedMsg{configs: []*career.CVConfig{}}
			newModel, _ := model.Update(msg)

			view := newModel.View()

			Expect(view).To(ContainSubstring("No CV configurations found"))
		})

		It("should render config list with details", func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			newModel, _ := model.Update(msg)

			view := newModel.View()

			Expect(view).To(ContainSubstring("Senior IC CV"))
			Expect(view).To(ContainSubstring("senior_ic"))
			Expect(view).To(ContainSubstring("Staff Engineer CV"))
		})

		It("should show keyboard shortcuts in footer", func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			newModel, _ := model.Update(msg)

			view := newModel.View()

			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("New"))
		})
	})

	Describe("Error Handling", func() {
		It("should display error when config loading fails", func() {
			errMsg := ConfigLoadError{err: fmt.Errorf("failed to load configs")}
			newModel, _ := model.Update(errMsg)

			castedModel := newModel.(*CVConfigManagerModel)
			Expect(castedModel.GetLastError()).NotTo(BeNil())
			Expect(castedModel.loading).To(BeFalse())
		})
	})
})

