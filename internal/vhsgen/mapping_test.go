package vhsgen

import (
	"testing"
)

func testMenuIntentCommands(t *testing.T, intent string, wantDowns int) {
	stepText := `I select "` + intent + `" from the menu`
	cmds, translatable, reason := TranslateStep(stepText, "When")

	if !translatable {
		t.Fatalf("expected translatable, got reason: %s", reason)
	}

	downCount := 0
	hasEnter := false

	for _, cmd := range cmds {
		switch cmd.Type {
		case Down:
			downCount++
		case Enter:
			hasEnter = true
		}
	}

	if downCount != wantDowns {
		t.Errorf("want %d Down commands, got %d", wantDowns, downCount)
	}

	if !hasEnter {
		t.Error("expected Enter command at end")
	}

	expectedLen := wantDowns + 1
	if len(cmds) != expectedLen {
		t.Errorf("want %d total commands, got %d", expectedLen, len(cmds))
	}
}

func TestMenuAllSevenIntents(t *testing.T) {
	tests := []struct {
		name      string
		intent    string
		wantDowns int
	}{
		{"capture_event is first (0 downs)", "capture_event", 0},
		{"browse_timeline is second (1 down)", "browse_timeline", 1},
		{"manage_skills is third (2 downs)", "manage_skills", 2},
		{"generate_cv is fourth (3 downs)", "generate_cv", 3},
		{"configure_system is fifth (4 downs)", "configure_system", 4},
		{"burst_management is sixth (5 downs)", "burst_management", 5},
		{"fact_management is seventh (6 downs)", "fact_management", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMenuIntentCommands(t, tt.intent, tt.wantDowns)
		})
	}
}

func TestMenuCaptureEventEnterOnly(t *testing.T) {
	cmds, translatable, _ := TranslateStep(`I select "capture_event" from the menu`, "When")

	if !translatable {
		t.Fatal("expected translatable")
	}

	if len(cmds) != 1 {
		t.Fatalf("capture_event should produce 1 command (Enter only), got %d", len(cmds))
	}

	if cmds[0].Type != Enter {
		t.Errorf("expected Enter, got %s", cmds[0].Type)
	}
}

func TestMenuManageSkillsTwoDowns(t *testing.T) {
	cmds, translatable, _ := TranslateStep(`I select "manage_skills" from the menu`, "When")

	if !translatable {
		t.Fatal("expected translatable")
	}

	if len(cmds) != 3 {
		t.Fatalf("manage_skills should produce 3 commands (2 Down + Enter), got %d", len(cmds))
	}

	if cmds[0].Type != Down || cmds[1].Type != Down {
		t.Error("first two commands should be Down")
	}

	if cmds[2].Type != Enter {
		t.Error("last command should be Enter")
	}
}

func TestMenuFactManagementSixDowns(t *testing.T) {
	cmds, translatable, _ := TranslateStep(`I select "fact_management" from the menu`, "When")

	if !translatable {
		t.Fatal("expected translatable")
	}

	if len(cmds) != 7 {
		t.Fatalf("fact_management should produce 7 commands (6 Down + Enter), got %d", len(cmds))
	}

	for i := range 6 {
		if cmds[i].Type != Down {
			t.Errorf("command[%d] should be Down, got %s", i, cmds[i].Type)
		}
	}

	if cmds[6].Type != Enter {
		t.Error("last command should be Enter")
	}
}

func TestMenuUnknownIntentReturnsNilCommands(t *testing.T) {
	cmds, translatable, _ := TranslateStep(`I select "nonexistent" from the menu`, "When")

	if !translatable {
		t.Fatal("pattern should match (translatable) even with unknown intent")
	}

	if cmds != nil {
		t.Errorf("unknown intent should produce nil commands, got %v", cmds)
	}
}

func TestUntranslatableFormBypass(t *testing.T) {
	formBypassSteps := []string{
		"I submit the event",
		"I submit the skill form",
		"I confirm filter",
		"I confirm sort",
		"I accept the suggested burst",
		"I accept all inferred skills",
		"I save the burst edit",
		"I save metadata changes",
		"I confirm the review",
	}

	for _, step := range formBypassSteps {
		t.Run(step, func(t *testing.T) {
			_, translatable, reason := TranslateStep(step, "When")

			if translatable {
				t.Fatalf("expected untranslatable for %q", step)
			}

			expected := "form-bypass: use keyboard navigation instead"
			if reason != expected {
				t.Errorf("want reason %q, got %q", expected, reason)
			}
		})
	}
}

func TestUntranslatableUnknownStep(t *testing.T) {
	_, translatable, reason := TranslateStep("I do something completely unknown", "When")

	if translatable {
		t.Fatal("expected untranslatable for unknown step")
	}

	expected := "unknown step: no matching pattern"
	if reason != expected {
		t.Errorf("want reason %q, got %q", expected, reason)
	}
}

func TestListTranslatablePatternsNotEmpty(t *testing.T) {
	patterns := ListTranslatablePatterns()

	if len(patterns) == 0 {
		t.Fatal("ListTranslatablePatterns should return >0 patterns")
	}

	for i, p := range patterns {
		if p.Pattern == "" {
			t.Errorf("pattern[%d] has empty Pattern", i)
		}

		if p.Type == "" {
			t.Errorf("pattern[%d] (%s) has empty Type", i, p.Pattern)
		}

		if p.Category == "" {
			t.Errorf("pattern[%d] (%s) has empty Category", i, p.Pattern)
		}

		if p.Example == "" {
			t.Errorf("pattern[%d] (%s) has empty Example", i, p.Pattern)
		}
	}
}

func TestListTranslatablePatternsMenuEnum(t *testing.T) {
	patterns := ListTranslatablePatterns()

	var found bool

	for _, p := range patterns {
		if p.Pattern != `^I select "([^"]*)" from the menu$` {
			continue
		}

		found = true

		intentParam, ok := p.Params["intent"]
		if !ok {
			t.Fatal("menu pattern should have 'intent' param")
		}

		if intentParam.Type != "enum" {
			t.Errorf("intent param type should be 'enum', got %q", intentParam.Type)
		}

		if len(intentParam.Values) != 7 {
			t.Errorf("want 7 valid intents, got %d", len(intentParam.Values))
		}

		break
	}

	if !found {
		t.Fatal("menu selection pattern not found in translatable patterns")
	}
}

func TestListTranslatablePatternsCategories(t *testing.T) {
	patterns := ListTranslatablePatterns()

	categories := make(map[string]bool)
	for _, p := range patterns {
		categories[p.Category] = true
	}

	expectedCats := []string{"navigation", "input", "setup"}
	for _, cat := range expectedCats {
		if !categories[cat] {
			t.Errorf("expected category %q in patterns", cat)
		}
	}
}

func TestListTranslatablePatternsExcludesFormBypass(t *testing.T) {
	patterns := ListTranslatablePatterns()

	for _, p := range patterns {
		if p.Category == "form-bypass" {
			t.Errorf("form-bypass should not appear in translatable: %s", p.Pattern)
		}
	}
}

func TestTranslateNavigationPrimitives(t *testing.T) {
	tests := []struct {
		step     string
		wantType VHSCommandType
	}{
		{"I press enter", Enter},
		{"I press enter to view event details", Enter},
		{"I press escape", Escape},
		{"I close the modal", Escape},
		{"I cancel", Escape},
		{"I navigate down", Down},
		{`I press "j" to navigate down`, Down},
		{"I navigate up", Up},
		{`I press "k" to navigate up`, Up},
		{"I press tab", Tab},
	}

	for _, tt := range tests {
		t.Run(tt.step, func(t *testing.T) {
			cmds, translatable, reason := TranslateStep(tt.step, "When")

			if !translatable {
				t.Fatalf("expected translatable, got: %s", reason)
			}

			if len(cmds) != 1 {
				t.Fatalf("want 1 command, got %d", len(cmds))
			}

			if cmds[0].Type != tt.wantType {
				t.Errorf("want %s, got %s", tt.wantType, cmds[0].Type)
			}
		})
	}
}

func TestTranslateKeyDiscrepancies(t *testing.T) {
	t.Run("press s to view events sends CtrlE", func(t *testing.T) {
		cmds, translatable, _ := TranslateStep(`I press "s" to view events`, "When")
		if !translatable {
			t.Fatal("expected translatable")
		}

		if len(cmds) != 1 || cmds[0].Type != CtrlE {
			t.Errorf("want Ctrl+E, got %v", cmds)
		}
	})

	t.Run("press m to open metadata sends e", func(t *testing.T) {
		cmds, translatable, _ := TranslateStep(`I press 'm' to open metadata editor`, "When")
		if !translatable {
			t.Fatal("expected translatable")
		}

		if len(cmds) != 1 || cmds[0].Type != Type || cmds[0].Args[0] != "e" {
			t.Errorf("want Type 'e', got %v", cmds)
		}
	})
}

func TestTranslateTextInput(t *testing.T) {
	cmds, translatable, _ := TranslateStep(`I enter event description "Built a REST API"`, "When")

	if !translatable {
		t.Fatal("expected translatable")
	}

	if len(cmds) != 1 {
		t.Fatalf("want 1 command, got %d", len(cmds))
	}

	if cmds[0].Type != Type {
		t.Errorf("want Type command, got %s", cmds[0].Type)
	}

	if len(cmds[0].Args) != 2 || cmds[0].Args[0] != "100ms" || cmds[0].Args[1] != "Built a REST API" {
		t.Errorf("want args [100ms, Built a REST API], got %v", cmds[0].Args)
	}
}

func TestTranslateSetupSteps(t *testing.T) {
	setupSteps := []string{
		"the database is empty",
		"I am on the main menu",
		"I have 3 skills in my profile",
		`I have a skill "Python"`,
		`I have a skill "Go" with category "backend"`,
		`I have an event "Built API" at company "Acme"`,
		`I have 2 events that use skill "Go"`,
	}

	for _, step := range setupSteps {
		t.Run(step, func(t *testing.T) {
			cmds, translatable, reason := TranslateStep(step, "Given")

			if !translatable {
				t.Fatalf("setup step should be translatable, got: %s", reason)
			}

			if cmds != nil {
				t.Errorf("setup step should produce nil commands, got %v", cmds)
			}
		})
	}
}
