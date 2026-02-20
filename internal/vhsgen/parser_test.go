package vhsgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func featuresDir(t *testing.T) string {
	t.Helper()

	dir := filepath.Join("..", "..", "features")
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("features directory not found: %v", err)
	}

	return dir
}

func TestParserParsesRealFeatures(t *testing.T) {
	results, err := ParseFeatureDir(featuresDir(t), SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least one scenario, got zero")
	}

	featureNames := make(map[string]bool)
	for _, ir := range results {
		featureNames[ir.Feature] = true
	}

	if len(featureNames) < 2 {
		t.Errorf("expected scenarios from multiple features, got %d feature(s)", len(featureNames))
	}
}

func TestParserSourceField(t *testing.T) {
	dir := featuresDir(t)

	tests := []struct {
		name   string
		source SourceType
	}{
		{"business source", SourceBusiness},
		{"vhs-only source", SourceVHSOnly},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ParseFeatureDir(dir, tt.source)
			if err != nil {
				t.Fatalf("ParseFeatureDir() error: %v", err)
			}

			for _, ir := range results {
				if ir.Source != tt.source {
					t.Errorf("scenario %q: Source = %v, want %v", ir.Name, ir.Source, tt.source)
				}
			}
		})
	}
}

func TestParserBackgroundSteps(t *testing.T) {
	dir := featuresDir(t)

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	var skillsScenarios []ScenarioIR
	for _, ir := range results {
		if ir.Feature == "Manage Skills" {
			skillsScenarios = append(skillsScenarios, ir)
		}
	}

	if len(skillsScenarios) == 0 {
		t.Fatal("expected scenarios from Manage Skills feature")
	}

	for _, ir := range skillsScenarios {
		if len(ir.SetupSteps) == 0 {
			t.Errorf("scenario %q: expected SetupSteps from Background, got none", ir.Name)
			continue
		}

		foundBackground := false
		for _, step := range ir.SetupSteps {
			if step.Text == "I am on the main menu" && step.StepType == "Given" {
				foundBackground = true
				break
			}
		}

		if !foundBackground {
			t.Errorf("scenario %q: SetupSteps missing Background step 'I am on the main menu'", ir.Name)
		}
	}
}

func TestParserBackgroundIsFirstSetupStep(t *testing.T) {
	dir := featuresDir(t)

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	for _, ir := range results {
		if ir.Feature != "Manage Skills" {
			continue
		}

		if len(ir.SetupSteps) == 0 {
			continue
		}

		if ir.SetupSteps[0].Text != "I am on the main menu" {
			t.Errorf("scenario %q: first SetupStep = %q, want 'I am on the main menu'",
				ir.Name, ir.SetupSteps[0].Text)
		}
	}
}

func TestParserEmptyDir(t *testing.T) {
	dir := t.TempDir()

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected zero results, got %d", len(results))
	}
}

func TestParserNonExistentDir(t *testing.T) {
	results, err := ParseFeatureDir("/tmp/nonexistent-vhsgen-dir-"+t.Name(), SourceBusiness)
	if err != nil {
		t.Fatalf("expected no error for non-existent dir, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected zero results, got %d", len(results))
	}
}

func TestParserStepClassification(t *testing.T) {
	dir := featuresDir(t)

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	for _, ir := range results {
		for _, step := range ir.SetupSteps {
			if step.StepType != "Given" {
				t.Errorf("scenario %q: SetupStep %q has StepType %q, want 'Given'",
					ir.Name, step.Text, step.StepType)
			}
		}

		for _, step := range ir.DemoSteps {
			if step.StepType != "When" && step.StepType != "Then" {
				t.Errorf("scenario %q: DemoStep %q has StepType %q, want 'When' or 'Then'",
					ir.Name, step.Text, step.StepType)
			}
		}
	}
}

func TestParserScenarioOutline(t *testing.T) {
	dir := t.TempDir()

	content := `Feature: Test Outline

  Scenario Outline: Greet user
    Given I am on the main menu
    When I enter "<input>"
    Then I should see "<output>"

    Examples:
      | input | output      |
      | Alice | Hello Alice |
      | Bob   | Hello Bob   |
`

	err := os.WriteFile(filepath.Join(dir, "outline.feature"), []byte(content), 0o600)
	if err != nil {
		t.Fatalf("writing test feature: %v", err)
	}

	results, err := ParseFeatureDir(dir, SourceVHSOnly)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 scenario from outline, got %d", len(results))
	}

	ir := results[0]

	if ir.Name != "Greet user" {
		t.Errorf("Name = %q, want 'Greet user'", ir.Name)
	}

	if ir.Source != SourceVHSOnly {
		t.Errorf("Source = %v, want %v", ir.Source, SourceVHSOnly)
	}

	if ir.Feature != "Test Outline" {
		t.Errorf("Feature = %q, want 'Test Outline'", ir.Feature)
	}

	foundSubstitutedInput := false
	foundSubstitutedOutput := false

	for _, step := range ir.DemoSteps {
		if strings.Contains(step.Text, "Alice") {
			foundSubstitutedInput = true
		}

		if strings.Contains(step.Text, "Hello Alice") {
			foundSubstitutedOutput = true
		}
	}

	if !foundSubstitutedInput {
		t.Error("expected DemoStep with 'Alice' after first-row substitution")
	}

	if !foundSubstitutedOutput {
		t.Error("expected DemoStep with 'Hello Alice' after first-row substitution")
	}

	allSteps := append(ir.SetupSteps, ir.DemoSteps...) //nolint:gocritic
	for _, step := range allSteps {
		if strings.Contains(step.Text, "<input>") || strings.Contains(step.Text, "<output>") {
			t.Errorf("step %q still contains unsubstituted placeholder", step.Text)
		}
	}
}

func TestParserScenarioOutlinePreservesSetup(t *testing.T) {
	dir := t.TempDir()

	content := `Feature: Outline With Setup

  Scenario Outline: Test with setup
    Given I am on the main menu
    When I enter "<text>"

    Examples:
      | text  |
      | hello |
      | world |
`

	err := os.WriteFile(filepath.Join(dir, "setup_outline.feature"), []byte(content), 0o600)
	if err != nil {
		t.Fatalf("writing test feature: %v", err)
	}

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(results))
	}

	ir := results[0]

	if len(ir.SetupSteps) != 1 {
		t.Fatalf("expected 1 SetupStep, got %d", len(ir.SetupSteps))
	}

	if ir.SetupSteps[0].Text != "I am on the main menu" {
		t.Errorf("SetupStep = %q, want 'I am on the main menu'", ir.SetupSteps[0].Text)
	}

	if len(ir.DemoSteps) != 1 {
		t.Fatalf("expected 1 DemoStep, got %d", len(ir.DemoSteps))
	}

	if !strings.Contains(ir.DemoSteps[0].Text, "hello") {
		t.Errorf("DemoStep = %q, expected 'hello' from first example row", ir.DemoSteps[0].Text)
	}
}

func TestParserTranslateStepIntegration(t *testing.T) {
	dir := featuresDir(t)

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	translatableFound := false
	untranslatableFound := false

	for _, ir := range results {
		for _, step := range append(ir.SetupSteps, ir.DemoSteps...) {
			if step.Translatable {
				translatableFound = true
			} else {
				untranslatableFound = true
			}
		}
	}

	if !translatableFound {
		t.Error("expected at least one translatable step across all features")
	}

	if !untranslatableFound {
		t.Error("expected at least one untranslatable step across all features")
	}
}

func TestParserFeatureAndScenarioTags(t *testing.T) {
	dir := featuresDir(t)

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	taggedCount := 0
	for _, ir := range results {
		if len(ir.Tags) > 0 {
			taggedCount++
		}
	}

	if taggedCount == 0 {
		t.Error("expected at least some scenarios to have tags")
	}
}

func TestParserDirWithNoFeatureFiles(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not a feature"), 0o600)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected zero results, got %d", len(results))
	}
}

func TestParserFeatureFileWithNoFeature(t *testing.T) {
	dir := t.TempDir()

	content := `# This is just a comment, no feature`
	err := os.WriteFile(filepath.Join(dir, "empty.feature"), []byte(content), 0o600)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected zero results for file with no feature, got %d", len(results))
	}
}

func TestParserMultipleFeatureFiles(t *testing.T) {
	dir := t.TempDir()

	content1 := `Feature: Feature One
  Scenario: Scenario One
    Given setup
    When action
    Then result`

	content2 := `Feature: Feature Two
  Scenario: Scenario Two
    Given setup
    When action
    Then result`

	err := os.WriteFile(filepath.Join(dir, "feature1.feature"), []byte(content1), 0o600)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	err = os.WriteFile(filepath.Join(dir, "feature2.feature"), []byte(content2), 0o600)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 scenarios, got %d", len(results))
	}
}

func TestParserSubstituteExampleValues(t *testing.T) {
	dir := t.TempDir()

	content := `Feature: Substitution Test
  Scenario Outline: Test with examples
    Given I have <count> items
    When I add <item>
    Then I have <total> items

    Examples:
      | count | item | total |
      | 5     | one  | 6     |
      | 10    | two  | 12    |`

	err := os.WriteFile(filepath.Join(dir, "outline.feature"), []byte(content), 0o600)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	results, err := ParseFeatureDir(dir, SourceBusiness)
	if err != nil {
		t.Fatalf("ParseFeatureDir() error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 scenario from outline (uses first row), got %d", len(results))
	}

	if !strings.Contains(results[0].SetupSteps[0].Text, "5") {
		t.Errorf("expected scenario to have substituted value '5', got %q", results[0].SetupSteps[0].Text)
	}

	if !strings.Contains(results[0].DemoSteps[0].Text, "one") {
		t.Errorf("expected scenario to have substituted value 'one', got %q", results[0].DemoSteps[0].Text)
	}
}
