package mocks_test

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

const aliasLinterName = "import" + "as"

func TestAliasRulesEnforceMockAliases(t *testing.T) {
	config := loadLintConfig(t)

	requiredAliases := map[string]string{
		"github.com/baphled/kariya/internal/testutil/mocks/repository": "mockrepo",
		"github.com/baphled/kariya/internal/testutil/mocks/service":    "mocksvc",
		"github.com/baphled/kariya/internal/testutil/mocks/intent":     "mockintent",
	}

	configuredAliases := extractAliasRules(t, config)

	for pkg, expectedAlias := range requiredAliases {
		t.Run(expectedAlias, func(t *testing.T) {
			actualAlias, exists := configuredAliases[pkg]
			if !exists {
				t.Errorf("missing alias rule for package %s: expected alias %q", pkg, expectedAlias)
				return
			}
			if actualAlias != expectedAlias {
				t.Errorf("wrong alias for package %s: got %q, want %q", pkg, actualAlias, expectedAlias)
			}
		})
	}
}

func TestAliasRulesRejectUnderscoreAliases(t *testing.T) {
	config := loadLintConfig(t)
	configuredAliases := extractAliasRules(t, config)

	forbiddenAliases := []struct {
		pkg          string
		badAlias     string
		correctAlias string
	}{
		{
			pkg:          "github.com/baphled/kariya/internal/testutil/mocks/repository",
			badAlias:     "mock_repo",
			correctAlias: "mockrepo",
		},
		{
			pkg:          "github.com/baphled/kariya/internal/testutil/mocks/repository",
			badAlias:     "mock_repository",
			correctAlias: "mockrepo",
		},
		{
			pkg:          "github.com/baphled/kariya/internal/testutil/mocks/service",
			badAlias:     "mock_svc",
			correctAlias: "mocksvc",
		},
		{
			pkg:          "github.com/baphled/kariya/internal/testutil/mocks/service",
			badAlias:     "mock_service",
			correctAlias: "mocksvc",
		},
		{
			pkg:          "github.com/baphled/kariya/internal/testutil/mocks/intent",
			badAlias:     "mock_intent",
			correctAlias: "mockintent",
		},
	}

	for _, tc := range forbiddenAliases {
		t.Run(tc.badAlias+"_is_rejected", func(t *testing.T) {
			actualAlias, exists := configuredAliases[tc.pkg]
			if !exists {
				t.Errorf("no alias rule for %s: underscore alias %q would not be caught", tc.pkg, tc.badAlias)
				return
			}
			if actualAlias == tc.badAlias {
				t.Errorf("alias rule for %s uses forbidden alias %q: should be %q", tc.pkg, tc.badAlias, tc.correctAlias)
			}
		})
	}
}

func TestAliasLinterIsEnabled(t *testing.T) {
	config := loadLintConfig(t)

	linters, ok := config["linters"].(map[string]interface{})
	if !ok {
		t.Fatal("missing 'linters' section in .golangci.yml")
	}

	enableList, ok := linters["enable"].([]interface{})
	if !ok {
		t.Fatal("missing 'linters.enable' list in .golangci.yml")
	}

	found := false
	for _, v := range enableList {
		if s, ok := v.(string); ok && s == aliasLinterName {
			found = true
			break
		}
	}

	if !found {
		t.Error("alias linter is not enabled in .golangci.yml: underscore import aliases will not be caught")
	}
}

func loadLintConfig(t *testing.T) map[string]interface{} {
	t.Helper()

	projectRoot := findProjectRoot(t)
	configPath := filepath.Join(projectRoot, ".golangci.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read .golangci.yml: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("failed to parse .golangci.yml: %v", err)
	}

	return config
}

func extractAliasRules(t *testing.T, config map[string]interface{}) map[string]string {
	t.Helper()

	result := make(map[string]string)

	linters, ok := config["linters"].(map[string]interface{})
	if !ok {
		t.Fatal("missing 'linters' section in .golangci.yml")
	}

	settings, ok := linters["settings"].(map[string]interface{})
	if !ok {
		t.Fatal("missing 'linters.settings' section in .golangci.yml")
	}

	linterSettings, ok := settings[aliasLinterName].(map[string]interface{})
	if !ok {
		t.Fatalf("missing 'linters.settings.%s' section in .golangci.yml", aliasLinterName)
	}

	aliases, ok := linterSettings["alias"].([]interface{})
	if !ok {
		t.Fatal("missing alias list in linter settings")
	}

	for _, entry := range aliases {
		m, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		pkg, _ := m["pkg"].(string)
		alias, _ := m["alias"].(string)
		if pkg != "" && alias != "" {
			result[pkg] = alias
		}
	}

	return result
}

func findProjectRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}
