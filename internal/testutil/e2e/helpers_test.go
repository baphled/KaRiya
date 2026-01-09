package e2e_test

import (
	"testing"

	"github.com/baphled/kariya/internal/testutil/e2e"
)

func TestSetup(t *testing.T) {
	env := e2e.Setup(t)
	defer env.Cleanup()

	if env.Model == nil {
		t.Fatal("Model should not be nil")
	}
	if env.DB == nil {
		t.Fatal("DB should not be nil for SQLite setup")
	}
	if env.Service == nil {
		t.Fatal("Service should not be nil")
	}
}

func TestSetupWithMemory(t *testing.T) {
	env := e2e.SetupWithMemory(t)
	defer env.Cleanup()

	if env.Model == nil {
		t.Fatal("Model should not be nil")
	}
	if env.DB != nil {
		t.Fatal("DB should be nil for memory setup")
	}
	if env.Service == nil {
		t.Fatal("Service should not be nil")
	}
}

func TestNavigationHelpers(t *testing.T) {
	env := e2e.SetupWithMemory(t)
	defer env.Cleanup()

	// Test that we start at the menu
	if !env.IsInMenuState() {
		t.Error("Should start in menu state")
	}

	// Test navigation
	env.NavigateDown().NavigateDown()

	// Should still be in menu (just different selection)
	if !env.IsInMenuState() {
		t.Error("Should still be in menu after navigation")
	}
}

func TestSelectIntent(t *testing.T) {
	env := e2e.SetupWithMemory(t)
	defer env.Cleanup()

	// Select CaptureEvent (index 0)
	env.SelectIntent(0)

	// Should no longer be in menu state
	// The view should change to show intent content
	view := env.GetView()
	if view == "" {
		t.Error("View should not be empty")
	}
}

func TestSelectIntentByName(t *testing.T) {
	env := e2e.SetupWithMemory(t)
	defer env.Cleanup()

	// Select by name
	env.SelectIntentByName("browse_timeline")

	// View should show timeline content
	view := env.GetView()
	if view == "" {
		t.Error("View should not be empty")
	}
}

func TestViewAssertions(t *testing.T) {
	env := e2e.SetupWithMemory(t)
	defer env.Cleanup()

	// Menu should contain certain text
	env.AssertViewContains("Career Event Management System")
	env.AssertViewContains("Capture Event")
}

func TestDataPopulation(t *testing.T) {
	env := e2e.Setup(t)
	defer env.Cleanup()

	// Start with empty database
	env.AssertEventCount(0)
	env.AssertBurstCount(0)
	env.AssertFactCount(0)

	// Add test data
	env.PopulateTestData(5, 2, 3)

	// Verify counts
	env.AssertEventCount(5)
	env.AssertBurstCount(2)
	env.AssertFactCount(3)
}

func TestSimulateRestart(t *testing.T) {
	env := e2e.Setup(t)
	defer env.Cleanup()

	// Add some data
	event := e2e.CreateMinimalEvent("test_event_1")
	env.AddEvent(event)
	env.AssertEventCount(1)

	// Simulate restart
	env.SimulateRestart()

	// Data should still be there
	env.AssertEventCount(1)

	// Should be back at menu
	if !env.IsInMenuState() {
		t.Error("Should be in menu state after restart")
	}
}

func TestCreateSampleEvents(t *testing.T) {
	events := e2e.CreateSampleEvents(10)

	if len(events) != 10 {
		t.Errorf("Expected 10 events, got %d", len(events))
	}

	for i, event := range events {
		if event.ID == "" {
			t.Errorf("Event %d has empty ID", i)
		}
		if event.Text == "" {
			t.Errorf("Event %d has empty Text", i)
		}
		if event.Date.IsZero() {
			t.Errorf("Event %d has zero Date", i)
		}
	}
}

func TestCreateSampleBursts(t *testing.T) {
	events := e2e.CreateSampleEvents(5)
	bursts := e2e.CreateSampleBursts(3, events)

	if len(bursts) != 3 {
		t.Errorf("Expected 3 bursts, got %d", len(bursts))
	}

	for i, burst := range bursts {
		if burst.ID == "" {
			t.Errorf("Burst %d has empty ID", i)
		}
		if burst.Name == "" {
			t.Errorf("Burst %d has empty Name", i)
		}
		if len(burst.EventIDs) == 0 {
			t.Errorf("Burst %d has no event IDs", i)
		}
	}
}

func TestCreateSampleFacts(t *testing.T) {
	events := e2e.CreateSampleEvents(5)
	facts := e2e.CreateSampleFacts(5, events)

	if len(facts) != 5 {
		t.Errorf("Expected 5 facts, got %d", len(facts))
	}

	for i, fact := range facts {
		if fact.ID == "" {
			t.Errorf("Fact %d has empty ID", i)
		}
		if fact.Text == "" {
			t.Errorf("Fact %d has empty Text", i)
		}
		if len(fact.CompetencyCategories) == 0 {
			t.Errorf("Fact %d has no competency categories", i)
		}
	}
}

func TestCreateSampleProfiles(t *testing.T) {
	profiles := e2e.CreateSampleProfiles()

	if len(profiles) < 3 {
		t.Errorf("Expected at least 3 profiles, got %d", len(profiles))
	}

	for i, profile := range profiles {
		if profile.ID == "" {
			t.Errorf("Profile %d has empty ID", i)
		}
		if profile.Name == "" {
			t.Errorf("Profile %d has empty Name", i)
		}
		if profile.TargetRole == "" {
			t.Errorf("Profile %d has empty TargetRole", i)
		}
	}
}
