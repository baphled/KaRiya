package vhsgen

import (
	"strings"
	"testing"
)

func TestRenderTape(t *testing.T) {
	data := TapeData{
		FeatureName:      "User Registration",
		ScenarioName:     "Successful registration",
		GIFPath:          "demos/vhs/features/user-registration/happy-path.gif",
		ASCIIPath:        "demos/vhs/features/user-registration/happy-path.ascii",
		ConfigSourcePath: "demos/vhs/config.tape",
		SetupCommands: `Type "mkdir -p /tmp/demo"
Enter
Sleep 300ms`,
		DemoCommands: `Type "./kariya --config /tmp/demo/config.yaml"
Enter
Sleep 2s`,
	}

	result, err := RenderTape(data)
	if err != nil {
		t.Fatalf("RenderTape failed: %v", err)
	}

	// Verify the result contains expected content
	tests := []struct {
		name     string
		expected string
	}{
		{"Feature comment", "# Feature: User Registration"},
		{"Scenario comment", "# Scenario: Successful registration"},
		{"Source directive", "Source demos/vhs/config.tape"},
		{"GIF output", "Output demos/vhs/features/user-registration/happy-path.gif"},
		{"ASCII output", "Output demos/vhs/features/user-registration/happy-path.ascii"},
		{"Hide block", "Hide"},
		{"Show block", "Show"},
		{"Setup commands", "mkdir -p /tmp/demo"},
		{"Demo commands", "./kariya --config /tmp/demo/config.yaml"},
		{"Exit command", "Ctrl+C"},
	}

	for _, tt := range tests {
		if !strings.Contains(result, tt.expected) {
			t.Errorf("Expected %q in rendered output, but not found", tt.expected)
		}
	}

	// Verify no cleanup commands are present
	forbiddenCommands := []string{"rm -rf", "DELETE", "DROP"}
	for _, cmd := range forbiddenCommands {
		if strings.Contains(result, cmd) {
			t.Errorf("Rendered output contains forbidden command: %q", cmd)
		}
	}

	// Verify both Output directives are present
	gifCount := strings.Count(result, "Output demos/vhs/features/user-registration/happy-path.gif")
	asciiCount := strings.Count(result, "Output demos/vhs/features/user-registration/happy-path.ascii")
	if gifCount != 1 {
		t.Errorf("Expected 1 GIF Output directive, found %d", gifCount)
	}
	if asciiCount != 1 {
		t.Errorf("Expected 1 ASCII Output directive, found %d", asciiCount)
	}
}

func TestRenderTapeMinimalData(t *testing.T) {
	data := TapeData{
		FeatureName:      "Minimal",
		ScenarioName:     "Test",
		GIFPath:          "out.gif",
		ASCIIPath:        "out.ascii",
		ConfigSourcePath: "config.tape",
		SetupCommands:    "",
		DemoCommands:     "",
	}

	result, err := RenderTape(data)
	if err != nil {
		t.Fatalf("RenderTape with minimal data failed: %v", err)
	}

	if !strings.Contains(result, "# Feature: Minimal") {
		t.Error("Expected feature name in output")
	}
	if !strings.Contains(result, "# Scenario: Test") {
		t.Error("Expected scenario name in output")
	}
}

func TestRenderTapeWithSpecialCharacters(t *testing.T) {
	data := TapeData{
		FeatureName:      "Feature with \"quotes\" and 'apostrophes'",
		ScenarioName:     "Scenario with special chars: <>&",
		GIFPath:          "path/with spaces/output.gif",
		ASCIIPath:        "path/with spaces/output.ascii",
		ConfigSourcePath: "config.tape",
		SetupCommands:    `Type "echo 'hello world'"`,
		DemoCommands:     `Type "curl http://example.com?foo=bar&baz=qux"`,
	}

	result, err := RenderTape(data)
	if err != nil {
		t.Fatalf("RenderTape with special characters failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}
