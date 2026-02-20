package vhsgen

import (
	"testing"
)

func TestSourceTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      SourceType
		expected SourceType
	}{
		{"SourceBusiness constant", SourceBusiness, "business"},
		{"SourceVHSOnly constant", SourceVHSOnly, "vhs-only"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got %q, want %q", tt.got, tt.expected)
			}
		})
	}
}

func TestVHSCommandTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      VHSCommandType
		expected VHSCommandType
	}{
		{"Type command", Type, "Type"},
		{"Down command", Down, "Down"},
		{"Up command", Up, "Up"},
		{"Enter command", Enter, "Enter"},
		{"Escape command", Escape, "Escape"},
		{"Tab command", Tab, "Tab"},
		{"Sleep command", Sleep, "Sleep"},
		{"Hide command", Hide, "Hide"},
		{"Show command", Show, "Show"},
		{"Screenshot command", Screenshot, "Screenshot"},
		{"Source command", Source, "Source"},
		{"Output command", Output, "Output"},
		{"CtrlC command", CtrlC, "Ctrl+C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got %q, want %q", tt.got, tt.expected)
			}
		})
	}
}

func TestVHSCommandConstruction(t *testing.T) {
	tests := []struct {
		name     string
		cmd      VHSCommand
		wantType VHSCommandType
		wantArgs []string
	}{
		{
			name:     "Type command with text",
			cmd:      VHSCommand{Type: Type, Args: []string{"hello"}},
			wantType: Type,
			wantArgs: []string{"hello"},
		},
		{
			name:     "Sleep command with duration",
			cmd:      VHSCommand{Type: Sleep, Args: []string{"500ms"}},
			wantType: Sleep,
			wantArgs: []string{"500ms"},
		},
		{
			name:     "Screenshot command with path",
			cmd:      VHSCommand{Type: Screenshot, Args: []string{"output.png"}},
			wantType: Screenshot,
			wantArgs: []string{"output.png"},
		},
		{
			name:     "Command with no args",
			cmd:      VHSCommand{Type: Enter},
			wantType: Enter,
			wantArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.Type != tt.wantType {
				t.Errorf("Type: got %q, want %q", tt.cmd.Type, tt.wantType)
			}
			if len(tt.cmd.Args) != len(tt.wantArgs) {
				t.Errorf("Args length: got %d, want %d", len(tt.cmd.Args), len(tt.wantArgs))
			}
			for i, arg := range tt.cmd.Args {
				if arg != tt.wantArgs[i] {
					t.Errorf("Args[%d]: got %q, want %q", i, arg, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestStepIRConstruction(t *testing.T) {
	tests := []struct {
		name     string
		step     StepIR
		wantText string
		wantType string
	}{
		{
			name: "Given step",
			step: StepIR{
				Text:         "a user is logged in",
				StepType:     "Given",
				Translatable: true,
			},
			wantText: "a user is logged in",
			wantType: "Given",
		},
		{
			name: "When step",
			step: StepIR{
				Text:         "the user clicks the button",
				StepType:     "When",
				Translatable: true,
			},
			wantText: "the user clicks the button",
			wantType: "When",
		},
		{
			name: "Then step",
			step: StepIR{
				Text:         "the page should display success",
				StepType:     "Then",
				Translatable: true,
			},
			wantText: "the page should display success",
			wantType: "Then",
		},
		{
			name: "Untranslatable step",
			step: StepIR{
				Text:                 "some complex step",
				StepType:             "When",
				Translatable:         false,
				UntranslatableReason: "no matching pattern",
			},
			wantText: "some complex step",
			wantType: "When",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.step.Text != tt.wantText {
				t.Errorf("Text: got %q, want %q", tt.step.Text, tt.wantText)
			}
			if tt.step.StepType != tt.wantType {
				t.Errorf("StepType: got %q, want %q", tt.step.StepType, tt.wantType)
			}
		})
	}
}

func TestStepIRWithCommands(t *testing.T) {
	step := StepIR{
		Text:     "the user types their name",
		StepType: "When",
		Commands: []VHSCommand{
			{Type: Type, Args: []string{"John Doe"}},
			{Type: Tab},
		},
		Translatable: true,
	}

	if len(step.Commands) != 2 {
		t.Errorf("Commands length: got %d, want 2", len(step.Commands))
	}

	if step.Commands[0].Type != Type {
		t.Errorf("First command type: got %q, want %q", step.Commands[0].Type, Type)
	}

	if step.Commands[1].Type != Tab {
		t.Errorf("Second command type: got %q, want %q", step.Commands[1].Type, Tab)
	}
}

func TestScenarioIRConstruction(t *testing.T) {
	scenario := ScenarioIR{
		Name:    "User login",
		Feature: "Authentication",
		Tags:    []string{"@critical", "@smoke"},
		Source:  SourceBusiness,
		SetupSteps: []StepIR{
			{Text: "setup step", StepType: "Given", Translatable: true},
		},
		DemoSteps: []StepIR{
			{Text: "demo step", StepType: "When", Translatable: true},
		},
		Translatable: true,
	}

	if scenario.Name != "User login" {
		t.Errorf("Name: got %q, want %q", scenario.Name, "User login")
	}

	if scenario.Feature != "Authentication" {
		t.Errorf("Feature: got %q, want %q", scenario.Feature, "Authentication")
	}

	if scenario.Source != SourceBusiness {
		t.Errorf("Source: got %q, want %q", scenario.Source, SourceBusiness)
	}

	if len(scenario.Tags) != 2 {
		t.Errorf("Tags length: got %d, want 2", len(scenario.Tags))
	}

	if len(scenario.SetupSteps) != 1 {
		t.Errorf("SetupSteps length: got %d, want 1", len(scenario.SetupSteps))
	}

	if len(scenario.DemoSteps) != 1 {
		t.Errorf("DemoSteps length: got %d, want 1", len(scenario.DemoSteps))
	}
}

func TestScenarioIRSourceTypes(t *testing.T) {
	tests := []struct {
		name   string
		source SourceType
	}{
		{"Business source", SourceBusiness},
		{"VHS-only source", SourceVHSOnly},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scenario := ScenarioIR{
				Name:   "Test scenario",
				Source: tt.source,
			}

			if scenario.Source != tt.source {
				t.Errorf("Source: got %q, want %q", scenario.Source, tt.source)
			}
		})
	}
}

func TestGeneratorConfigConstruction(t *testing.T) {
	config := GeneratorConfig{
		OutputDir:        "/tmp/output",
		TemplatePath:     "/path/to/template.tape",
		ConfigSourcePath: "demos/vhs/config.tape",
		SleepDuration:    "500ms",
		ScenariosDir:     "features/",
	}

	if config.OutputDir != "/tmp/output" {
		t.Errorf("OutputDir: got %q, want %q", config.OutputDir, "/tmp/output")
	}

	if config.TemplatePath != "/path/to/template.tape" {
		t.Errorf("TemplatePath: got %q, want %q", config.TemplatePath, "/path/to/template.tape")
	}

	if config.ConfigSourcePath != "demos/vhs/config.tape" {
		t.Errorf("ConfigSourcePath: got %q, want %q", config.ConfigSourcePath, "demos/vhs/config.tape")
	}

	if config.SleepDuration != "500ms" {
		t.Errorf("SleepDuration: got %q, want %q", config.SleepDuration, "500ms")
	}

	if config.ScenariosDir != "features/" {
		t.Errorf("ScenariosDir: got %q, want %q", config.ScenariosDir, "features/")
	}
}

func TestAnalysisResultConstruction(t *testing.T) {
	result := AnalysisResult{
		ScenarioName: "Login flow",
		Feature:      "Authentication",
		Translatable: true,
		Source:       SourceBusiness,
		Warnings:     []string{"slow step detected"},
		Errors:       []string{},
	}

	if result.ScenarioName != "Login flow" {
		t.Errorf("ScenarioName: got %q, want %q", result.ScenarioName, "Login flow")
	}

	if result.Feature != "Authentication" {
		t.Errorf("Feature: got %q, want %q", result.Feature, "Authentication")
	}

	if !result.Translatable {
		t.Errorf("Translatable: got false, want true")
	}

	if result.Source != SourceBusiness {
		t.Errorf("Source: got %q, want %q", result.Source, SourceBusiness)
	}

	if len(result.Warnings) != 1 {
		t.Errorf("Warnings length: got %d, want 1", len(result.Warnings))
	}

	if len(result.Errors) != 0 {
		t.Errorf("Errors length: got %d, want 0", len(result.Errors))
	}
}

func TestAnalysisResultWithUntranslatableSteps(t *testing.T) {
	untranslatableStep := StepIR{
		Text:                 "complex step",
		StepType:             "When",
		Translatable:         false,
		UntranslatableReason: "no pattern match",
	}

	result := AnalysisResult{
		ScenarioName:        "Complex scenario",
		Feature:             "Advanced",
		Translatable:        false,
		UntranslatableSteps: []StepIR{untranslatableStep},
		Source:              SourceVHSOnly,
		Errors:              []string{"cannot translate scenario"},
	}

	if result.Translatable {
		t.Errorf("Translatable: got true, want false")
	}

	if len(result.UntranslatableSteps) != 1 {
		t.Errorf("UntranslatableSteps length: got %d, want 1", len(result.UntranslatableSteps))
	}

	if result.UntranslatableSteps[0].Text != "complex step" {
		t.Errorf("Step text: got %q, want %q", result.UntranslatableSteps[0].Text, "complex step")
	}
}

func TestParamConstraintConstruction(t *testing.T) {
	constraint := ParamConstraint{
		Type:   "enum",
		Values: []string{"value1", "value2", "value3"},
	}

	if constraint.Type != "enum" {
		t.Errorf("Type: got %q, want %q", constraint.Type, "enum")
	}

	if len(constraint.Values) != 3 {
		t.Errorf("Values length: got %d, want 3", len(constraint.Values))
	}

	if constraint.Values[0] != "value1" {
		t.Errorf("Values[0]: got %q, want %q", constraint.Values[0], "value1")
	}
}

func TestStepPatternConstruction(t *testing.T) {
	pattern := StepPattern{
		Pattern:  `the user types "([^"]+)"`,
		Type:     "When",
		Category: "input",
		Params: map[string]ParamConstraint{
			"text": {
				Type:   "string",
				Values: nil,
			},
		},
		Example: `the user types "hello"`,
	}

	if pattern.Pattern != `the user types "([^"]+)"` {
		t.Errorf("Pattern: got %q, want %q", pattern.Pattern, `the user types "([^"]+)"`)
	}

	if pattern.Type != "When" {
		t.Errorf("Type: got %q, want %q", pattern.Type, "When")
	}

	if pattern.Category != "input" {
		t.Errorf("Category: got %q, want %q", pattern.Category, "input")
	}

	if len(pattern.Params) != 1 {
		t.Errorf("Params length: got %d, want 1", len(pattern.Params))
	}

	if pattern.Example != `the user types "hello"` {
		t.Errorf("Example: got %q, want %q", pattern.Example, `the user types "hello"`)
	}
}

func testZeroVHSCommand(t *testing.T) {
	var cmd VHSCommand
	if cmd.Type != "" {
		t.Errorf("zero VHSCommand.Type should be empty string, got %q", cmd.Type)
	}
	if cmd.Args != nil {
		t.Errorf("zero VHSCommand.Args should be nil, got %v", cmd.Args)
	}
}

func testZeroStepIR(t *testing.T) {
	var step StepIR
	if step.Text != "" {
		t.Errorf("zero StepIR.Text should be empty string, got %q", step.Text)
	}
	if step.Translatable {
		t.Errorf("zero StepIR.Translatable should be false, got true")
	}
}

func testZeroScenarioIR(t *testing.T) {
	var scenario ScenarioIR
	if scenario.Name != "" {
		t.Errorf("zero ScenarioIR.Name should be empty string, got %q", scenario.Name)
	}
	if scenario.Source != "" {
		t.Errorf("zero ScenarioIR.Source should be empty string, got %q", scenario.Source)
	}
}

func testZeroGeneratorConfig(t *testing.T) {
	var config GeneratorConfig
	if config.OutputDir != "" {
		t.Errorf("zero GeneratorConfig.OutputDir should be empty string, got %q", config.OutputDir)
	}
}

func testZeroAnalysisResult(t *testing.T) {
	var result AnalysisResult
	if result.ScenarioName != "" {
		t.Errorf("zero AnalysisResult.ScenarioName should be empty string, got %q", result.ScenarioName)
	}
	if result.Translatable {
		t.Errorf("zero AnalysisResult.Translatable should be false, got true")
	}
}

func testZeroParamConstraint(t *testing.T) {
	var constraint ParamConstraint
	if constraint.Type != "" {
		t.Errorf("zero ParamConstraint.Type should be empty string, got %q", constraint.Type)
	}
}

func testZeroStepPattern(t *testing.T) {
	var pattern StepPattern
	if pattern.Pattern != "" {
		t.Errorf("zero StepPattern.Pattern should be empty string, got %q", pattern.Pattern)
	}
	if pattern.Params != nil {
		t.Errorf("zero StepPattern.Params should be nil, got %v", pattern.Params)
	}
}

func TestZeroValues(t *testing.T) {
	t.Run("zero VHSCommand", testZeroVHSCommand)
	t.Run("zero StepIR", testZeroStepIR)
	t.Run("zero ScenarioIR", testZeroScenarioIR)
	t.Run("zero GeneratorConfig", testZeroGeneratorConfig)
	t.Run("zero AnalysisResult", testZeroAnalysisResult)
	t.Run("zero ParamConstraint", testZeroParamConstraint)
	t.Run("zero StepPattern", testZeroStepPattern)
}

func TestEmptySlices(t *testing.T) {
	t.Run("ScenarioIR with empty slices", func(t *testing.T) {
		scenario := ScenarioIR{
			Name:       "Empty scenario",
			Tags:       []string{},
			SetupSteps: []StepIR{},
			DemoSteps:  []StepIR{},
		}

		if len(scenario.Tags) != 0 {
			t.Errorf("Tags length: got %d, want 0", len(scenario.Tags))
		}
		if len(scenario.SetupSteps) != 0 {
			t.Errorf("SetupSteps length: got %d, want 0", len(scenario.SetupSteps))
		}
		if len(scenario.DemoSteps) != 0 {
			t.Errorf("DemoSteps length: got %d, want 0", len(scenario.DemoSteps))
		}
	})

	t.Run("AnalysisResult with empty slices", func(t *testing.T) {
		result := AnalysisResult{
			ScenarioName:        "Empty result",
			UntranslatableSteps: []StepIR{},
			Warnings:            []string{},
			Errors:              []string{},
		}

		if len(result.UntranslatableSteps) != 0 {
			t.Errorf("UntranslatableSteps length: got %d, want 0", len(result.UntranslatableSteps))
		}
		if len(result.Warnings) != 0 {
			t.Errorf("Warnings length: got %d, want 0", len(result.Warnings))
		}
		if len(result.Errors) != 0 {
			t.Errorf("Errors length: got %d, want 0", len(result.Errors))
		}
	})
}
