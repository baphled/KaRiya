package vhsgen_test

import (
	"testing"

	"github.com/baphled/kariya/internal/vhsgen"
)

func TestAnalyser_BusinessScenarioAllTranslatable(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:    "Navigate menu",
			Feature: "Navigation",
			Source:  vhsgen.SourceBusiness,
			SetupSteps: []vhsgen.StepIR{
				{Text: "I am on the main menu", StepType: "Given", Translatable: true},
			},
			DemoSteps: []vhsgen.StepIR{
				{Text: `I select "manage_skills" from the menu`, StepType: "When", Translatable: true},
				{Text: "I press enter", StepType: "When", Translatable: true},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if !r.Translatable {
		t.Errorf("expected translatable, got false")
	}
	if len(r.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", r.Warnings)
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
	if len(r.UntranslatableSteps) != 0 {
		t.Errorf("expected no untranslatable steps, got %d", len(r.UntranslatableSteps))
	}
}

func TestAnalyser_BusinessScenarioFormBypass(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:    "Submit event",
			Feature: "Capture Event",
			Source:  vhsgen.SourceBusiness,
			DemoSteps: []vhsgen.StepIR{
				{Text: `I enter event description "Built API"`, StepType: "When", Translatable: true},
				{Text: "I submit the event", StepType: "When", Translatable: false, UntranslatableReason: "form-bypass: use keyboard navigation instead"},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Translatable {
		t.Errorf("expected untranslatable, got true")
	}
	if len(r.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(r.Warnings), r.Warnings)
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
	if len(r.UntranslatableSteps) != 1 {
		t.Fatalf("expected 1 untranslatable step, got %d", len(r.UntranslatableSteps))
	}
	if r.UntranslatableSteps[0].Text != "I submit the event" {
		t.Errorf("expected untranslatable step text 'I submit the event', got %q", r.UntranslatableSteps[0].Text)
	}
}

func TestAnalyser_VHSOnlyUntranslatableStep(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:    "Demo custom flow",
			Feature: "VHS Demo",
			Source:  vhsgen.SourceVHSOnly,
			DemoSteps: []vhsgen.StepIR{
				{Text: "I press enter", StepType: "When", Translatable: true},
				{Text: "I do something custom", StepType: "When", Translatable: false, UntranslatableReason: "unknown step: no matching pattern"},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Translatable {
		t.Errorf("expected untranslatable, got true")
	}
	if len(r.Warnings) != 0 {
		t.Errorf("expected no warnings for VHS-only, got %v", r.Warnings)
	}
	if len(r.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(r.Errors), r.Errors)
	}

	wantErr := "Step 'I do something custom' not found in mapping. Run `vhsgen list --steps` to see available steps."
	if r.Errors[0] != wantErr {
		t.Errorf("expected error %q, got %q", wantErr, r.Errors[0])
	}
}

func TestAnalyser_GivenThenStepsDontAffectTranslatability(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:    "Setup does not block",
			Feature: "Resilience",
			Source:  vhsgen.SourceBusiness,
			SetupSteps: []vhsgen.StepIR{
				{Text: "some untranslatable setup", StepType: "Given", Translatable: false, UntranslatableReason: "unknown"},
				{Text: "some untranslatable assertion", StepType: "Then", Translatable: false, UntranslatableReason: "unknown"},
			},
			DemoSteps: []vhsgen.StepIR{
				{Text: "I press enter", StepType: "When", Translatable: true},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if !r.Translatable {
		t.Errorf("expected translatable (Given/Then don't affect), got false")
	}
	if len(r.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", r.Warnings)
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
}

func TestAnalyser_SourceFieldPropagated(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:   "Business scenario",
			Source: vhsgen.SourceBusiness,
			DemoSteps: []vhsgen.StepIR{
				{Text: "I press enter", StepType: "When", Translatable: true},
			},
		},
		{
			Name:   "VHS-only scenario",
			Source: vhsgen.SourceVHSOnly,
			DemoSteps: []vhsgen.StepIR{
				{Text: "I press enter", StepType: "When", Translatable: true},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Source != vhsgen.SourceBusiness {
		t.Errorf("expected SourceBusiness, got %q", results[0].Source)
	}
	if results[1].Source != vhsgen.SourceVHSOnly {
		t.Errorf("expected SourceVHSOnly, got %q", results[1].Source)
	}
}

func TestAnalyser_ScenarioNameAndFeaturePropagated(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:    "My Scenario",
			Feature: "My Feature",
			Source:  vhsgen.SourceBusiness,
			DemoSteps: []vhsgen.StepIR{
				{Text: "I press enter", StepType: "When", Translatable: true},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if results[0].ScenarioName != "My Scenario" {
		t.Errorf("expected scenario name 'My Scenario', got %q", results[0].ScenarioName)
	}
	if results[0].Feature != "My Feature" {
		t.Errorf("expected feature 'My Feature', got %q", results[0].Feature)
	}
}

func TestAnalyser_MultipleUntranslatableWhenSteps(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:   "Multiple blockers",
			Source: vhsgen.SourceVHSOnly,
			DemoSteps: []vhsgen.StepIR{
				{Text: "I do thing A", StepType: "When", Translatable: false, UntranslatableReason: "unknown"},
				{Text: "I press enter", StepType: "When", Translatable: true},
				{Text: "I do thing B", StepType: "When", Translatable: false, UntranslatableReason: "unknown"},
			},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	r := results[0]
	if r.Translatable {
		t.Errorf("expected untranslatable")
	}
	if len(r.UntranslatableSteps) != 2 {
		t.Fatalf("expected 2 untranslatable steps, got %d", len(r.UntranslatableSteps))
	}
	if len(r.Errors) != 2 {
		t.Fatalf("expected 2 errors for VHS-only, got %d", len(r.Errors))
	}
}

func TestAnalyser_EmptyDemoStepsIsTranslatable(t *testing.T) {
	scenarios := []vhsgen.ScenarioIR{
		{
			Name:   "Setup only",
			Source: vhsgen.SourceBusiness,
			SetupSteps: []vhsgen.StepIR{
				{Text: "I am on the main menu", StepType: "Given", Translatable: true},
			},
			DemoSteps: []vhsgen.StepIR{},
		},
	}

	results := vhsgen.AnalyseScenarios(scenarios)

	if !results[0].Translatable {
		t.Errorf("expected translatable with empty demo steps")
	}
}
