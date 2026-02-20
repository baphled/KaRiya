package vhsgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateTapeContainsExpectedDirectives(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Happy path registration",
		Feature: "User Registration",
		SetupSteps: []StepIR{
			{
				Text:         "the database is empty",
				StepType:     "Given",
				Translatable: true,
				Commands:     nil,
			},
		},
		DemoSteps: []StepIR{
			{
				Text:         `I select "capture_event" from the menu`,
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
			{
				Text:         `I enter event description "Built API"`,
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Type, Args: []string{"100ms", "Built API"}}},
			},
		},
		Translatable: true,
	}

	config := GeneratorConfig{
		OutputDir:        "demos/vhs/generated",
		ConfigSourcePath: "demos/vhs/config.tape",
		SleepDuration:    "2s",
	}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	checks := []struct {
		name     string
		expected string
	}{
		{"Source directive", "Source demos/vhs/config.tape"},
		{"GIF output", "Output demos/vhs/generated/user-registration/happy-path-registration.gif"},
		{"ASCII output", "Output demos/vhs/generated/user-registration/happy-path-registration.ascii"},
		{"Hide block", "Hide"},
		{"Show block", "Show"},
		{"Ctrl+C ending", "Ctrl+C"},
		{"Enter command", "Enter"},
		{"Type command", `Type@100ms "Built API"`},
		{"Sleep between steps", "Sleep 2s"},
		{"Feature comment", "# Feature: User Registration"},
		{"Scenario comment", "# Scenario: Happy path registration"},
	}

	for _, tt := range checks {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected %q in output, not found.\nOutput:\n%s", tt.expected, result)
			}
		})
	}
}

func TestGenerateTapeDualOutputDirectives(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Test scenario",
		Feature: "Test Feature",
		DemoSteps: []StepIR{
			{
				Text:         "do something",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	gifCount := strings.Count(result, "Output out/test-feature/test-scenario.gif")
	asciiCount := strings.Count(result, "Output out/test-feature/test-scenario.ascii")

	if gifCount != 1 {
		t.Errorf("Expected exactly 1 GIF Output directive, found %d", gifCount)
	}

	if asciiCount != 1 {
		t.Errorf("Expected exactly 1 ASCII Output directive, found %d", asciiCount)
	}
}

func TestGenerateTapeUntranslatableStepProducesTODO(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Form submission",
		Feature: "Event Capture",
		DemoSteps: []StepIR{
			{
				Text:                 "I submit the event",
				StepType:             "When",
				Translatable:         false,
				UntranslatableReason: "form-bypass: use keyboard navigation instead",
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	expected := "# [Manual step needed] — I submit the event (form-bypass: use keyboard navigation instead)"
	if !strings.Contains(result, expected) {
		t.Errorf("Expected manual step marker in output.\nExpected: %s\nGot:\n%s", expected, result)
	}
}

func TestGenerateTapeNoCleanupCommands(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Test scenario",
		Feature: "Test Feature",
		DemoSteps: []StepIR{
			{
				Text:         "do something",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	forbidden := []string{"rm -rf", "DELETE", "DROP"}
	for _, cmd := range forbidden {
		if strings.Contains(result, cmd) {
			t.Errorf("Output contains forbidden command: %q", cmd)
		}
	}
}

func TestGenerateTapeRejectsForbiddenInSteps(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Dangerous",
		Feature: "Danger",
		DemoSteps: []StepIR{
			{
				Text:         "clean up",
				StepType:     "When",
				Translatable: true,
				Commands: []VHSCommand{
					{Type: Type, Args: []string{"rm -rf /tmp/data"}},
				},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	_, err := GenerateTape(scenario, config)
	if err == nil {
		t.Fatal("Expected error for forbidden pattern, got nil")
	}

	if !strings.Contains(err.Error(), "forbidden pattern") {
		t.Errorf("Expected forbidden pattern error, got: %v", err)
	}
}

func TestGenerateTapeDefaultConfig(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Defaults",
		Feature: "Config",
		DemoSteps: []StepIR{
			{
				Text:         "step one",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Down}},
			},
			{
				Text:         "step two",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	if !strings.Contains(result, "Source demos/vhs/config.tape") {
		t.Error("Expected default ConfigSourcePath")
	}

	if !strings.Contains(result, "Sleep 2s") {
		t.Error("Expected default SleepDuration of 2s")
	}
}

func TestGenerateTapeSleepBetweenStepsOnly(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Sleep test",
		Feature: "Sleep",
		DemoSteps: []StepIR{
			{
				Text:         "step one",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Down}, {Type: Down}, {Type: Enter}},
			},
			{
				Text:         "step two",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{
		OutputDir:     "out",
		SleepDuration: "3s",
	}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	sleepCount := strings.Count(result, "Sleep 3s")
	if sleepCount != 1 {
		t.Errorf("Expected exactly 1 Sleep between 2 steps, found %d.\nOutput:\n%s", sleepCount, result)
	}
}

func TestGenerateTapeSetupStepsNoCommands(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Setup only",
		Feature: "Setup",
		SetupSteps: []StepIR{
			{Text: "the database is empty", StepType: "Given", Translatable: true, Commands: nil},
			{Text: "I am on the main menu", StepType: "Given", Translatable: true, Commands: nil},
		},
		DemoSteps: []StepIR{
			{
				Text:         "do something",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	hideIdx := strings.Index(result, "Hide")
	showIdx := strings.Index(result, "Show")

	if hideIdx == -1 || showIdx == -1 {
		t.Fatal("Expected both Hide and Show blocks")
	}

	setupBlock := result[hideIdx+4 : showIdx]
	trimmed := strings.TrimSpace(setupBlock)

	if trimmed != "" {
		t.Errorf("Expected empty setup block for commandless steps, got: %q", trimmed)
	}
}

func TestGenerateTapeRenderCommandVariants(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Command variants",
		Feature: "Commands",
		DemoSteps: []StepIR{
			{
				Text:         "various commands",
				StepType:     "When",
				Translatable: true,
				Commands: []VHSCommand{
					{Type: Down},
					{Type: Up},
					{Type: Enter},
					{Type: Escape},
					{Type: Tab},
					{Type: Type, Args: []string{"a"}},
					{Type: Type, Args: []string{"100ms", "hello world"}},
					{Type: CtrlC},
					{Type: CtrlE},
					{Type: CtrlS},
				},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	result, err := GenerateTape(scenario, config)
	if err != nil {
		t.Fatalf("GenerateTape failed: %v", err)
	}

	expectations := []struct {
		name     string
		expected string
	}{
		{"Down key", "\nDown\n"},
		{"Up key", "\nUp\n"},
		{"Enter key", "\nEnter\n"},
		{"Escape key", "\nEscape\n"},
		{"Tab key", "\nTab\n"},
		{"Type char", `Type "a"`},
		{"Type with speed", `Type@100ms "hello world"`},
		{"Ctrl+E", "Ctrl+E"},
		{"Ctrl+S", "Ctrl+S"},
	}

	for _, tt := range expectations {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected %q in output, not found.\nOutput:\n%s", tt.expected, result)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple spaces", "Hello World", "hello-world"},
		{"already slug", "hello-world", "hello-world"},
		{"special chars", "Hello, World! 123", "hello-world-123"},
		{"multiple spaces", "hello   world", "hello-world"},
		{"leading trailing spaces", " hello world ", "hello-world"},
		{"empty string", "", ""},
		{"underscores to hyphens", "capture_event", "capture-event"},
		{"mixed separators", "hello_world test", "hello-world-test"},
		{"consecutive special", "a!!b##c", "abc"},
		{"only special chars", "!@#$%", ""},
		{"numbers", "version 2", "version-2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWriteTapeCreatesFile(t *testing.T) {
	tmpDir := t.TempDir()

	scenario := ScenarioIR{
		Name:    "Write test",
		Feature: "File Output",
		DemoSteps: []StepIR{
			{
				Text:         "action",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: tmpDir}

	err := WriteTape(scenario, config)
	if err != nil {
		t.Fatalf("WriteTape failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "file-output", "write-test.tape")

	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read written tape file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Source demos/vhs/config.tape") {
		t.Error("Written file missing Source directive")
	}

	if !strings.Contains(content, "Ctrl+C") {
		t.Error("Written file missing Ctrl+C ending")
	}
}

func TestWriteTapeCreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "deep", "nested")

	scenario := ScenarioIR{
		Name:    "Nested",
		Feature: "Deep",
		DemoSteps: []StepIR{
			{
				Text:         "action",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: nestedDir}

	err := WriteTape(scenario, config)
	if err != nil {
		t.Fatalf("WriteTape failed: %v", err)
	}

	expectedPath := filepath.Join(nestedDir, "deep", "nested.tape")

	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file at %s, but it does not exist", expectedPath)
	}
}

func TestWriteTapeErrorOnInvalidOutputDir(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Test",
		Feature: "Test",
		DemoSteps: []StepIR{
			{
				Text:         "action",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Enter}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "/invalid/nonexistent/path/that/cannot/be/created"}

	err := WriteTape(scenario, config)
	if err == nil {
		t.Error("Expected WriteTape to fail with invalid output directory")
	}
}

func TestRenderCommandWithNoArgs(t *testing.T) {
	tests := []struct {
		name     string
		cmd      VHSCommand
		expected string
	}{
		{"Down key", VHSCommand{Type: Down}, "Down"},
		{"Up key", VHSCommand{Type: Up}, "Up"},
		{"Enter key", VHSCommand{Type: Enter}, "Enter"},
		{"Escape key", VHSCommand{Type: Escape}, "Escape"},
		{"Tab key", VHSCommand{Type: Tab}, "Tab"},
		{"CtrlC", VHSCommand{Type: CtrlC}, "Ctrl+C"},
		{"CtrlE", VHSCommand{Type: CtrlE}, "Ctrl+E"},
		{"CtrlS", VHSCommand{Type: CtrlS}, "Ctrl+S"},
		{"Sleep no args", VHSCommand{Type: Sleep}, "Sleep 1s"},
		{"Type no args", VHSCommand{Type: Type}, "Type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderCommand(tt.cmd)
			if result != tt.expected {
				t.Errorf("renderCommand(%v) = %q, want %q", tt.cmd, result, tt.expected)
			}
		})
	}
}

func TestRenderCommandWithArgs(t *testing.T) {
	tests := []struct {
		name     string
		cmd      VHSCommand
		expected string
	}{
		{"Type with text", VHSCommand{Type: Type, Args: []string{"hello"}}, "Type \"hello\""},
		{"Type with speed and text", VHSCommand{Type: Type, Args: []string{"100ms", "world"}}, "Type@100ms \"world\""},
		{"Sleep with duration", VHSCommand{Type: Sleep, Args: []string{"3s"}}, "Sleep 3s"},
		{"Generic command with arg", VHSCommand{Type: "Custom", Args: []string{"arg1"}}, "Custom arg1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderCommand(tt.cmd)
			if result != tt.expected {
				t.Errorf("renderCommand(%v) = %q, want %q", tt.cmd, result, tt.expected)
			}
		})
	}
}

func TestGenerateTapeWithForbiddenPattern(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "Dangerous",
		Feature: "Security",
		DemoSteps: []StepIR{
			{
				Text:         "delete everything",
				StepType:     "When",
				Translatable: true,
				Commands:     []VHSCommand{{Type: Type, Args: []string{"rm -rf /tmp/data"}}},
			},
		},
	}

	config := GeneratorConfig{OutputDir: "out"}

	_, err := GenerateTape(scenario, config)
	if err == nil {
		t.Error("expected GenerateTape to reject forbidden pattern 'rm -rf'")
	}
}
