package models

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCVConfigEditor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CVConfigEditor Suite")
}

var _ = Describe("CVConfigEditorModel", func() {
	var (
		model         *CVConfigEditorModel
		baseModel     *BaseStandardModel
		configManager cv.ConfigManager
		testConfig    *career.CVConfig
	)

	BeforeEach(func() {
		baseModel = NewBaseStandardModel()
		configManager = cv.NewMemoryConfigManager()

		testConfig = &career.CVConfig{
			Name:           "Test CV",
			TargetRole:     "senior_ic",
			TargetAudience: []string{"hiring_manager", "recruiter"},
			EventFilters:   map[string]interface{}{},
			CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:      time.Now().Add(-30 * 24 * time.Hour),
		}
	})

	Describe("NewCVConfigEditorModel", func() {
		It("should create editor for new config", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			Expect(model.isNew).To(BeTrue())
			Expect(model.config).NotTo(BeNil())
			Expect(model.config.Name).To(BeEmpty())
			Expect(model.focusIndex).To(Equal(0))
			Expect(model.formFields).To(HaveLen(3))
		})

		It("should create editor for existing config", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, testConfig)

			Expect(model.isNew).To(BeFalse())
			Expect(model.config).To(Equal(testConfig))
			Expect(model.formFields[0].Value).To(Equal("Test CV"))
		})

		It("should initialize form fields", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, testConfig)

			Expect(model.formFields).To(HaveLen(3))
			Expect(model.formFields[0].Name).To(Equal("cv_name"))
			Expect(model.formFields[1].Name).To(Equal("target_role"))
			Expect(model.formFields[2].Name).To(Equal("target_audience"))
		})

		It("should mark required fields", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, testConfig)

			for _, field := range model.formFields {
				Expect(field.Required).To(BeTrue())
			}
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)
		})

		It("should move focus down with Tab", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(1))
		})

		It("should move focus up with Shift+Tab", func() {
			model.focusIndex = 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(0))
		})

		It("should move focus down with arrow key", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(1))
		})

		It("should move focus up with arrow key", func() {
			model.focusIndex = 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(0))
		})

		It("should wrap around when pressing Tab at last field", func() {
			model.focusIndex = len(model.formFields) - 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(0))
		})

		It("should not go below first field with up arrow", func() {
			model.focusIndex = 0
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(0))
		})

		It("should not go above last field with down arrow", func() {
			model.focusIndex = len(model.formFields) - 1
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.focusIndex).To(Equal(len(model.formFields) - 1))
		})
	})

	Describe("Text Input", func() {
		BeforeEach(func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)
		})

		It("should add characters to focused field", func() {
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.formFields[0].Value).To(Equal("T"))
		})

		It("should build up text with multiple keystrokes", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'M'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.formFields[0].Value).To(Equal("My "))
		})

		It("should handle backspace", func() {
			model.formFields[0].Value = "Test"
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.formFields[0].Value).To(Equal("Tes"))
		})

		It("should not go below zero length with backspace", func() {
			model.formFields[0].Value = ""
			newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

			castedModel := newModel.(*CVConfigEditorModel)
			Expect(castedModel.formFields[0].Value).To(Equal(""))
		})
	})

	Describe("Validation", func() {
		It("should reject empty required fields", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			cmd := model.saveConfig()
			result := cmd()

			Expect(result).To(BeAssignableToTypeOf(ConfigValidationErrorMsg{}))
			errMsg := result.(ConfigValidationErrorMsg)
			Expect(errMsg.errors).To(HaveKey("cv_name"))
			Expect(errMsg.errors).To(HaveKey("target_role"))
		})

		It("should validate role is one of allowed values", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)
			model.formFields[0].Value = "My CV"
			model.formFields[1].Value = "invalid_role"
			model.formFields[2].Value = "hiring_manager"

			cmd := model.saveConfig()
			result := cmd()

			Expect(result).To(BeAssignableToTypeOf(ConfigSaveErrorMsg{}))
		})

		It("should accept valid configuration", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)
			model.formFields[0].Value = "Valid CV"
			model.formFields[1].Value = "senior_ic"
			model.formFields[2].Value = "hiring_manager"

			cmd := model.saveConfig()
			result := cmd()

			Expect(result).To(BeAssignableToTypeOf(ConfigSavedMsg{}))
			savedMsg := result.(ConfigSavedMsg)
			Expect(savedMsg.config.Name).To(Equal("Valid CV"))
		})
	})

	Describe("Escape Key", func() {
		It("should navigate back on Esc", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(BackToCVConfigManagerMsg{}))
		})
	})

	Describe("View", func() {
		It("should show new config title for new config", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			view := model.View()

			Expect(view).To(ContainSubstring("New CV Configuration"))
		})

		It("should show edit title for existing config", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, testConfig)

			view := model.View()

			Expect(view).To(ContainSubstring("Edit CV Configuration"))
		})

		It("should show required field indicator", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			view := model.View()

			Expect(view).To(ContainSubstring("*"))
		})

		It("should show save button", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			view := model.View()

			Expect(view).To(ContainSubstring("Save"))
		})

		It("should show help text", func() {
			model = NewCVConfigEditorModel(baseModel, configManager, nil)

			view := model.View()

			Expect(view).To(ContainSubstring("Tab"))
		})
	})

	Describe("AudienceListToString", func() {
		It("should convert empty list to empty string", func() {
			result := audienceListToString([]string{})
			Expect(result).To(Equal(""))
		})

		It("should convert single audience", func() {
			result := audienceListToString([]string{"hiring_manager"})
			Expect(result).To(Equal("hiring_manager"))
		})

		It("should convert multiple audiences with comma separator", func() {
			result := audienceListToString([]string{"hiring_manager", "recruiter", "peer"})
			Expect(result).To(Equal("hiring_manager, recruiter, peer"))
		})
	})

	Describe("AudienceStringToList", func() {
		It("should convert empty string to empty list", func() {
			result := audienceStringToList("")
			Expect(result).To(Equal([]string{}))
		})

		It("should convert single audience", func() {
			result := audienceStringToList("hiring_manager")
			Expect(result).To(Equal([]string{"hiring_manager"}))
		})

		It("should convert comma-separated audiences", func() {
			result := audienceStringToList("hiring_manager, recruiter, peer")
			Expect(result).To(Equal([]string{"hiring_manager", "recruiter", "peer"}))
		})
	})
})

