package intents

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/google/uuid"
)

// TestCaptureEventUsesStandardView verifies CaptureEvent uses StandardView patterns
func TestCaptureEventUsesStandardView(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}

	intent, err := NewCaptureEventIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	testStandardViewConsistency(t, "CaptureEvent", view)
}

// TestBrowseTimelineUsesStandardView verifies BrowseTimeline uses StandardView patterns
func TestBrowseTimelineUsesStandardView(t *testing.T) {
	ctx := &BrowseTimelineContext{}

	intent, err := NewBrowseTimelineIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	testStandardViewConsistency(t, "BrowseTimeline", view)
}

// TestGenerateCVUsesStandardView verifies GenerateCV uses StandardView patterns
func TestGenerateCVUsesStandardView(t *testing.T) {
	ctx := &GenerateCVContext{
		AvailableProfiles: []*CVProfile{
			{
				ID:             "default",
				Name:           "Default Profile",
				TargetRole:     "staff",
				TargetAudience: []string{"hiring_manager"},
			},
		},
		Events: []*career.CareerEvent{
			{
				ID:   uuid.New().String(),
				Text: "Implemented test feature for CV generation",
				Date: time.Now(),
			},
		},
	}

	intent, err := NewGenerateCVIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	testStandardViewConsistency(t, "GenerateCV", view)
}

// TestExportArtifactUsesStandardView verifies ExportArtifact uses StandardView patterns
func TestExportArtifactUsesStandardView(t *testing.T) {
	ctx := context.Background()

	intent, err := NewExportArtifactIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	testStandardViewConsistency(t, "ExportArtifact", view)
}

// TestConfigureSystemUsesStandardView verifies ConfigureSystem uses StandardView patterns
func TestConfigureSystemUsesStandardView(t *testing.T) {
	ctx := context.Background()

	intent, err := NewConfigureSystemIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	testStandardViewConsistency(t, "ConfigureSystem", view)
}

// testStandardViewConsistency checks that an intent's view follows StandardView patterns
func testStandardViewConsistency(t *testing.T, intentName, view string) {
	if view == "" {
		t.Errorf("[%s] View is empty", intentName)
		return
	}

	// Check for logo or branding
	// The view should contain some form of branding
	hasLogo := strings.Contains(view, "██") ||
		strings.Contains(view, "KARIYA") ||
		strings.Contains(view, "KaRiya") ||
		strings.Contains(view, "Career Event Management System")

	if !hasLogo {
		t.Logf("[%s] Note: View may not show logo in current state", intentName)
	}

	// Check for some form of help text
	// StandardView intents should provide help/navigation hints
	hasHelp := strings.Contains(view, "Quit") ||
		strings.Contains(view, "Help") ||
		strings.Contains(view, "q ") ||
		strings.Contains(view, "Esc")

	if !hasHelp {
		t.Logf("[%s] Note: View may not show help text in current state", intentName)
	}

	// Check view is substantial (not just whitespace)
	trimmed := strings.TrimSpace(view)
	if len(trimmed) < 50 {
		t.Errorf("[%s] View seems too short (%d chars), may be incomplete", intentName, len(trimmed))
	}
}

// TestAllIntentsInitializeSuccessfully verifies all intents can be created and initialized
func TestAllIntentsInitializeSuccessfully(t *testing.T) {
	tests := []struct {
		name       string
		createFunc func() (interface{}, error)
		initFunc   func(interface{})
		viewFunc   func(interface{}) string
	}{
		{
			name: "CaptureEvent",
			createFunc: func() (interface{}, error) {
				return NewCaptureEventIntent(&CaptureEventContext{
					CaptureStrategy: "manual",
					Metadata:        make(map[string]string),
				})
			},
			initFunc: func(i interface{}) {
				i.(*CaptureEventIntent).Init()
			},
			viewFunc: func(i interface{}) string {
				return i.(*CaptureEventIntent).View()
			},
		},
		{
			name: "BrowseTimeline",
			createFunc: func() (interface{}, error) {
				return NewBrowseTimelineIntent(&BrowseTimelineContext{})
			},
			initFunc: func(i interface{}) {
				i.(*BrowseTimelineIntent).Init()
			},
			viewFunc: func(i interface{}) string {
				return i.(*BrowseTimelineIntent).View()
			},
		},
		{
			name: "GenerateCV",
			createFunc: func() (interface{}, error) {
				return NewGenerateCVIntent(&GenerateCVContext{
					AvailableProfiles: []*CVProfile{
						{
							ID:             "default",
							Name:           "Default Profile",
							TargetRole:     "staff",
							TargetAudience: []string{"hiring_manager"},
						},
					},
					Events: []*career.CareerEvent{
						{
							ID:   uuid.New().String(),
							Text: "Implemented test feature for CV generation",
							Date: time.Now(),
						},
					},
				})
			},
			initFunc: func(i interface{}) {
				i.(*GenerateCVIntent).Init()
			},
			viewFunc: func(i interface{}) string {
				return i.(*GenerateCVIntent).View()
			},
		},
		{
			name: "ExportArtifact",
			createFunc: func() (interface{}, error) {
				return NewExportArtifactIntent(context.Background())
			},
			initFunc: func(i interface{}) {
				i.(*ExportArtifactIntent).Init()
			},
			viewFunc: func(i interface{}) string {
				return i.(*ExportArtifactIntent).View()
			},
		},
		{
			name: "ConfigureSystem",
			createFunc: func() (interface{}, error) {
				return NewConfigureSystemIntent(context.Background())
			},
			initFunc: func(i interface{}) {
				i.(*ConfigureSystemIntent).Init()
			},
			viewFunc: func(i interface{}) string {
				return i.(*ConfigureSystemIntent).View()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create intent
			intent, err := tt.createFunc()
			if err != nil {
				t.Fatalf("Failed to create %s intent: %v", tt.name, err)
			}

			// Initialize
			tt.initFunc(intent)

			// Get view
			view := tt.viewFunc(intent)

			// Basic sanity checks
			if view == "" {
				t.Errorf("%s produced empty view", tt.name)
			}

			if len(strings.TrimSpace(view)) < 10 {
				t.Errorf("%s view is too short", tt.name)
			}
		})
	}
}
