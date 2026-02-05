package generatecv

import (
	"context"
	"errors"

	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
)

// ExtractedTechnology re-exports technology.ExtractedTechnology so that the
// generatecv package can reference extracted skill data without importing the
// technology service package directly.
type ExtractedTechnology = technology.ExtractedTechnology

// FocusAreaSuggestion re-exports technology.FocusAreaSuggestion so that the
// generatecv package can consume AI-generated focus area recommendations
// without a direct dependency on the technology service package.
type FocusAreaSuggestion = technology.FocusAreaSuggestion

// CVProfile is a type alias for the shared types.CVProfile.
type CVProfile = types.CVProfile

// CVExporter defines the interface for CV export operations.
// This interface enables testing by allowing mock implementations.
//
// Side effects:
//   - None at interface definition.
type CVExporter interface {
	ExportToText(
		ctx context.Context, cvView *career.CVView,
		sections []*career.CVSection, bullets map[string][]*career.CVBullet,
	) (string, error)
	ExportToMarkdown(
		ctx context.Context, cvView *career.CVView,
		sections []*career.CVSection, bullets map[string][]*career.CVBullet,
	) (string, error)
	ExportToYAML(
		ctx context.Context, cvView *career.CVView,
		sections []*career.CVSection, bullets map[string][]*career.CVBullet,
	) (string, error)
	SaveToFile(
		ctx context.Context, cvName string, format cv.ExportFormat, content string,
	) (string, error)
	CopyToClipboard(ctx context.Context, content string) error
}

// IntentContext is the input context passed to the GenerateCV intent.
//
// Side effects:
//   - None.
type IntentContext struct {
	AvailableProfiles     []*CVProfile
	Events                []*career.Event
	Facts                 []*career.Fact
	DefaultProfile        *CVProfile
	CVGenerationService   cv.CVGenerationService
	DataProcessingService cv.DataProcessingService
	BulletGenerator       cv.BulletGenerator
	ExportService         CVExporter
	SkillRepository       careerRepo.SkillRepository
	EventRepository       careerRepo.EventRepository
	ProfileConfig         *config.ProfileConfig
	AppContext            context.Context
}

// Validate ensures the context carries the minimum data needed to start the
// CV generation workflow.
//
// Expected:
//   - ctx must be non-nil.
//   - ctx.AvailableProfiles must contain at least one profile.
//   - ctx.Events must contain at least one event.
//
// Returns:
//   - Nil when all preconditions are met.
//   - A descriptive error when any required field is missing or empty.
//
// Side effects:
//   - None.
func (ctx *IntentContext) Validate() error {
	if ctx == nil {
		return errors.New("IntentContext cannot be nil")
	}

	if len(ctx.AvailableProfiles) == 0 {
		return errors.New("IntentContext must have at least one available profile")
	}

	if len(ctx.Events) == 0 {
		return errors.New("IntentContext must have at least one event")
	}

	return nil
}
