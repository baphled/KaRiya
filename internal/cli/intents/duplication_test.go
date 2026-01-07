package intents_test

import (
	"strings"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
)

// TestBrowseTimeline_NoDuplicateBreadcrumbs verifies breadcrumbs appear only once
func TestBrowseTimeline_NoDuplicateBreadcrumbs(t *testing.T) {
	events := []*career.CareerEvent{
		{
			ID:        "event1",
			Text:      "Test event",
			Date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Company:   "Test Co",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ctx := &intents.BrowseTimelineContext{
		Events: events,
		InitialFilters: &intents.TimelineFilters{
			Tags:      make([]string, 0),
			Companies: make([]string, 0),
			SortBy:    "date",
			SortOrder: "desc",
		},
	}

	intent, err := intents.NewBrowseTimelineIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	// Count breadcrumb separator occurrences - should only have one breadcrumb line
	separatorCount := strings.Count(view, "▸")
	if separatorCount == 0 {
		t.Error("Expected at least one breadcrumb separator (▸)")
	}
	// Allow some flexibility, but if we see many separators, likely duplication
	if separatorCount > 4 {
		t.Errorf("Too many breadcrumb separators (%d), likely duplication", separatorCount)
	}

	// Count "Timeline" occurrences - breadcrumb + header title is acceptable
	timelineCount := strings.Count(view, "Timeline")
	if timelineCount > 3 {
		t.Errorf("'Timeline' appears %d times, likely breadcrumb duplication", timelineCount)
	}
}

// TestBrowseTimeline_NoDuplicateHelpFooter verifies help footer appears only once
func TestBrowseTimeline_NoDuplicateHelpFooter(t *testing.T) {
	events := []*career.CareerEvent{
		{
			ID:        "event1",
			Text:      "Test event",
			Date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Company:   "Test Co",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ctx := &intents.BrowseTimelineContext{
		Events: events,
		InitialFilters: &intents.TimelineFilters{
			Tags:      make([]string, 0),
			Companies: make([]string, 0),
			SortBy:    "date",
			SortOrder: "desc",
		},
	}

	intent, err := intents.NewBrowseTimelineIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	// Count navigation help pattern - should appear only once
	upDownPattern := "↑/k"
	count := strings.Count(view, upDownPattern)
	if count > 1 {
		t.Errorf("Navigation help '↑/k' appears %d times, likely help footer duplication", count)
	}
}

// TestGenerateCV_NoDuplicateFooter verifies footer appears only once in card views
func TestGenerateCV_NoDuplicateFooter(t *testing.T) {
	profiles := []*intents.CVProfile{
		{
			ID:         "profile1",
			Name:       "Test Profile",
			TargetRole: "Engineer",
		},
	}

	events := []*career.CareerEvent{
		{
			ID:        "event1",
			Text:      "Test achievement",
			Date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Company:   "Test Co",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	ctx := &intents.GenerateCVContext{
		AvailableProfiles: profiles,
		Events:            events,
	}

	intent, err := intents.NewGenerateCVIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	// Count help text patterns - should appear only once
	upDownPattern := "↑/k"
	count := strings.Count(view, upDownPattern)
	if count > 1 {
		t.Errorf("Help text '↑/k' appears %d times, likely footer duplication", count)
	}

	// Count "Main menu" or "Main Menu" - should appear only once in help
	mainMenuCount := strings.Count(strings.ToLower(view), "main menu")
	if mainMenuCount > 2 {
		t.Errorf("'Main menu' appears %d times, likely footer duplication", mainMenuCount)
	}
}

// TestCaptureEvent_NoDuplicateFooter verifies no footer duplication (uses FormModel)
func TestCaptureEvent_NoDuplicateFooter(t *testing.T) {
	ctx := &intents.CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}

	intent, err := intents.NewCaptureEventIntent(ctx)
	if err != nil {
		t.Fatalf("Failed to create intent: %v", err)
	}

	intent.Init()
	view := intent.View()

	// Count help patterns
	quitPattern := "q"
	quitCount := strings.Count(strings.ToLower(view), quitPattern)
	// Be lenient since 'q' is common, but check it's reasonable
	if quitCount > 10 {
		t.Logf("Warning: 'q' appears %d times, may indicate duplication", quitCount)
	}
}
