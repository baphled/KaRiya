package models

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVGeneratorModel", func() {
	var (
		model      *CVGeneratorModel
		baseModel  *BaseStandardModel
		cvService  cv.CVGenerationService
		testConfig *career.CVConfig
	)

	BeforeEach(func() {
		baseModel = NewBaseStandardModel()

		// Create a mock CV service
		cvService = NewMockCVGenerationService()

		testConfig = &career.CVConfig{
			Name:           "Test CV",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			EventFilters:   map[string]interface{}{},
		}

		model = NewCVGeneratorModel(baseModel, cvService, testConfig)
	})

	Describe("Initialization", func() {
		It("should initialize with generating state true", func() {
			Expect(model.generating).To(BeTrue())
			Expect(model.generatedCV).To(BeNil())
		})

		It("should have config set", func() {
			Expect(model.config).To(Equal(testConfig))
		})
	})

	Describe("CVGeneratedMsg", func() {
		It("should update model with generated CV", func() {
			cvView := &career.CVView{
				ID:               uuid.New().String(),
				Name:             "Test CV",
				TargetRole:       "senior_ic",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 5,
				SourceFactCount:  3,
			}

			msg := CVGeneratedMsg{cvView: cvView}
			newModel, cmd := model.Update(msg)

			castedModel := newModel.(*CVGeneratorModel)
			Expect(castedModel.generating).To(BeFalse())
			Expect(castedModel.generatedCV).To(Equal(cvView))
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("CVGenerationErrorMsg", func() {
		It("should set error and stop generating", func() {
			errMsg := CVGenerationErrorMsg{err: ErrNoSourceEvents}
			newModel, _ := model.Update(errMsg)

			castedModel := newModel.(*CVGeneratorModel)
			Expect(castedModel.generating).To(BeFalse())
			Expect(castedModel.GetLastError()).NotTo(BeNil())
		})
	})

	Describe("Navigation", func() {
		It("should go back on Esc when not generating", func() {
			model.generating = false
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(BackMsg{}))
		})

		It("should go back on 'q'", func() {
			model.generating = false
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			Expect(result).To(BeAssignableToTypeOf(QuitMsg{}))
		})

		It("should ignore input when generating", func() {
			model.generating = true
			newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})

			castedModel := newModel.(*CVGeneratorModel)
			Expect(castedModel.generating).To(BeTrue())
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		It("should show generating message", func() {
			view := model.View()

			Expect(view).To(ContainSubstring("Generating CV"))
		})

		It("should show error when generation fails", func() {
			model.generating = false
			model.SetError(ErrNoSourceEvents)

			view := model.View()

			Expect(view).To(ContainSubstring("Failed"))
		})
	})
})

// MockCVGenerationService is a mock implementation for testing
type MockCVGenerationService struct{}

func NewMockCVGenerationService() cv.CVGenerationService {
	return &MockCVGenerationService{}
}

func (m *MockCVGenerationService) GenerateCV(ctx context.Context, configName string) (*career.CVView, error) {
	return &career.CVView{
		ID:          uuid.New().String(),
		Name:        configName,
		GeneratedAt: time.Now(),
	}, nil
}

func (m *MockCVGenerationService) GenerateCVFromConfig(ctx context.Context, config *career.CVConfig) (*career.CVView, error) {
	return &career.CVView{
		ID:               uuid.New().String(),
		Name:             config.Name,
		TargetRole:       config.TargetRole,
		TargetAudience:   config.TargetAudience,
		GeneratedAt:      time.Now(),
		SourceEventCount: 10,
		SourceFactCount:  5,
	}, nil
}

// ErrNoSourceEvents is a test error
var ErrNoSourceEvents = fmt.Errorf("no source events found")
