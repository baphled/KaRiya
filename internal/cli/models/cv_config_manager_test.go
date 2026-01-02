package models

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)


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
	})

	Describe("View", func() {
		It("should show loading message when loading", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Loading configuration templates"))
		})

		It("should show empty message when no configs", func() {
			msg := ConfigsLoadedMsg{configs: []*career.CVConfig{}}
			newModel, _ := model.Update(msg)

			view := newModel.View()

			Expect(view).To(ContainSubstring("No CV configuration templates found"))
		})

		It("should render config list with details", func() {
			msg := ConfigsLoadedMsg{configs: testConfigs}
			newModel, _ := model.Update(msg)

			view := newModel.View()

			Expect(view).To(ContainSubstring("Senior IC CV"))
			Expect(view).To(ContainSubstring("senior_ic"))
			Expect(view).To(ContainSubstring("Staff Engineer CV"))
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
