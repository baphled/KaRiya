// Package main provides the vhsgen CLI for listing and generating VHS tapes.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/baphled/kariya/internal/vhsgen"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	if len(args) == 0 {
		printUsageTo(out)
		return 0
	}

	subcommand := args[0]
	rest := args[1:]

	switch subcommand {
	case "list":
		return runList(rest, out, errOut)
	case "generate":
		return runGenerate(rest, out, errOut)
	case "--help", "-h", "help":
		printUsageTo(out)
		return 0
	default:
		fmt.Fprintf(errOut, "Error: unknown subcommand %q\n\n", subcommand)
		printUsageTo(errOut)
		return 1
	}
}

// runList implements the `list` subcommand.
func runList(args []string, out io.Writer, errOut io.Writer) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	featuresDir := fs.String("features", "features/", "Directory containing .feature files")
	scenariosDir := fs.String("scenarios-dir", "demos/vhs/scenarios/", "Directory containing VHS-only .feature files")
	asJSON := fs.Bool("json", false, "Output as JSON")
	showCount := fs.Bool("count", false, "Show counts broken down by source")
	showSteps := fs.Bool("steps", false, "Show translatable step patterns")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(errOut, "Error parsing flags: %v\n", err)
		return 1
	}

	if *showSteps {
		return runListSteps(*asJSON, out)
	}

	fmt.Fprintf(out, "Parsing...\n")

	if _, err := os.Stat(*featuresDir); err != nil {
		fmt.Fprintf(errOut, "Error parsing features dir %q: %v\n", *featuresDir, err)
		return 1
	}

	businessScenarios, err := vhsgen.ParseFeatureDir(*featuresDir, vhsgen.SourceBusiness)
	if err != nil {
		fmt.Fprintf(errOut, "Error parsing features dir %q: %v\n", *featuresDir, err)
		return 1
	}

	vhsOnlyScenarios, err := vhsgen.ParseFeatureDir(*scenariosDir, vhsgen.SourceVHSOnly)
	if err != nil {
		fmt.Fprintf(errOut, "Error parsing scenarios dir %q: %v\n", *scenariosDir, err)
		return 1
	}

	allScenarios := make([]vhsgen.ScenarioIR, 0, len(businessScenarios)+len(vhsOnlyScenarios))
	allScenarios = append(allScenarios, businessScenarios...)
	allScenarios = append(allScenarios, vhsOnlyScenarios...)
	results := vhsgen.AnalyseScenarios(allScenarios)

	if *showCount {
		return runListCount(results, out)
	}

	if *asJSON {
		return runListJSON(results, out, errOut)
	}

	return runListTable(results, out)
}

// runListSteps outputs the translatable step patterns.
func runListSteps(asJSON bool, out io.Writer) int {
	patterns := vhsgen.ListTranslatablePatterns()

	if asJSON {
		type jsonPattern struct {
			Pattern  string                            `json:"pattern"`
			Type     string                            `json:"type"`
			Category string                            `json:"category"`
			Params   map[string]vhsgen.ParamConstraint `json:"params,omitempty"`
			Example  string                            `json:"example"`
		}

		output := make([]jsonPattern, 0, len(patterns))
		for _, p := range patterns {
			output = append(output, jsonPattern{
				Pattern:  p.Pattern,
				Type:     p.Type,
				Category: p.Category,
				Params:   p.Params,
				Example:  p.Example,
			})
		}

		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(output); err != nil {
			return 1
		}
		return 0
	}

	colPattern := 50
	colType := 10
	colCategory := 14

	header := fmt.Sprintf("%-*s  %-*s  %-*s  %s",
		colPattern, "Pattern",
		colType, "Type",
		colCategory, "Category",
		"Example",
	)
	separator := strings.Repeat("-", len(header)+10)

	fmt.Fprintln(out, header)
	fmt.Fprintln(out, separator)

	for _, p := range patterns {
		fmt.Fprintf(out, "%-*s  %-*s  %-*s  %s\n",
			colPattern, truncate(p.Pattern, colPattern),
			colType, truncate(p.Type, colType),
			colCategory, truncate(p.Category, colCategory),
			p.Example,
		)
	}

	return 0
}

// runListCount outputs counts by source.
func runListCount(results []vhsgen.AnalysisResult, out io.Writer) int {
	var (
		businessTotal        int
		businessTranslatable int
		vhsOnlyTotal         int
		vhsOnlyTranslatable  int
	)

	for i := range results {
		r := &results[i]
		switch r.Source {
		case vhsgen.SourceBusiness:
			businessTotal++
			if r.Translatable {
				businessTranslatable++
			}
		case vhsgen.SourceVHSOnly:
			vhsOnlyTotal++
			if r.Translatable {
				vhsOnlyTranslatable++
			}
		}
	}

	fmt.Fprintf(out, "Business: %d/%d translatable | VHS-only: %d/%d translatable\n",
		businessTranslatable, businessTotal,
		vhsOnlyTranslatable, vhsOnlyTotal,
	)

	return 0
}

// runListJSON outputs the analysis results as JSON.
func runListJSON(results []vhsgen.AnalysisResult, out io.Writer, errOut io.Writer) int {
	type jsonScenario struct {
		ScenarioName string `json:"scenario_name"`
		Feature      string `json:"feature"`
		Source       string `json:"source"`
		Translatable bool   `json:"translatable"`
		Reason       string `json:"reason,omitempty"`
	}

	scenarios := make([]jsonScenario, 0, len(results))
	for i := range results {
		r := &results[i]
		var reason string
		if !r.Translatable && len(r.UntranslatableSteps) > 0 {
			reasons := make([]string, 0, len(r.UntranslatableSteps))
			for _, s := range r.UntranslatableSteps {
				reasons = append(reasons, s.UntranslatableReason)
			}
			reason = strings.Join(reasons, "; ")
		}

		scenarios = append(scenarios, jsonScenario{
			ScenarioName: r.ScenarioName,
			Feature:      r.Feature,
			Source:       string(r.Source),
			Translatable: r.Translatable,
			Reason:       reason,
		})
	}

	payload := map[string]interface{}{
		"scenarios": scenarios,
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintf(errOut, "Error encoding JSON: %v\n", err)
		return 1
	}

	return 0
}

// runListTable outputs the analysis results as a formatted table.
func runListTable(results []vhsgen.AnalysisResult, out io.Writer) int {
	colScenario := 40
	colFeature := 25
	colSource := 10
	colTranslatable := 12
	colReason := 40

	header := fmt.Sprintf("%-*s  %-*s  %-*s  %-*s  %s",
		colScenario, "Scenario",
		colFeature, "Feature",
		colSource, "Source",
		colTranslatable, "Translatable",
		"Reason",
	)
	separator := strings.Repeat("-", colScenario+colFeature+colSource+colTranslatable+colReason+8)

	fmt.Fprintln(out, header)
	fmt.Fprintln(out, separator)

	for i := range results {
		r := &results[i]
		translatable := "yes"
		var reason string
		if !r.Translatable {
			translatable = "no"
			if len(r.UntranslatableSteps) > 0 {
				reasons := make([]string, 0, len(r.UntranslatableSteps))
				for _, s := range r.UntranslatableSteps {
					reasons = append(reasons, s.UntranslatableReason)
				}
				reason = strings.Join(reasons, "; ")
			}
		}

		fmt.Fprintf(out, "%-*s  %-*s  %-*s  %-*s  %s\n",
			colScenario, truncate(r.ScenarioName, colScenario),
			colFeature, truncate(r.Feature, colFeature),
			colSource, truncate(string(r.Source), colSource),
			colTranslatable, translatable,
			truncate(reason, colReason),
		)
	}

	return 0
}

// runGenerate implements the `generate` subcommand.
func runGenerate(args []string, out io.Writer, errOut io.Writer) int {
	opts, err := parseGenerateFlags(args, errOut)
	if err != nil {
		return 1
	}

	fmt.Fprintf(out, "Parsing...\n")

	allScenarios, err := parseAllScenarios(*opts.featuresDir, *opts.scenariosDir, errOut)
	if err != nil {
		return 1
	}

	results := vhsgen.AnalyseScenarios(allScenarios)
	filtered := filterResults(results, allScenarios, *opts.generateAll, *opts.featureFilter, *opts.scenarioFilter)

	fmt.Fprintf(out, "Generating...\n")

	cfg := generateConfig{
		outputDir:    *opts.outputDir,
		configSource: *opts.configSource,
		verbose:      *opts.verbose,
		out:          out,
		errOut:       errOut,
	}
	stats := generateTapes(filtered, cfg)

	fmt.Fprintf(out, "Generated %d tapes (%d from features, %d from scenarios, %d warnings)\n",
		stats.total, stats.fromBusiness, stats.fromVHSOnly, stats.warnings)

	return 0
}

type generateOptions struct {
	generateAll    *bool
	featureFilter  *string
	scenarioFilter *string
	featuresDir    *string
	scenariosDir   *string
	outputDir      *string
	configSource   *string
	verbose        *bool
}

func parseGenerateFlags(args []string, errOut io.Writer) (*generateOptions, error) {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	opts := &generateOptions{
		generateAll:    fs.Bool("all", false, "Generate for all translatable scenarios"),
		featureFilter:  fs.String("feature", "", "Filter by feature name"),
		scenarioFilter: fs.String("scenario", "", "Filter by scenario name"),
		featuresDir:    fs.String("features", "features/", "Directory containing .feature files"),
		scenariosDir:   fs.String("scenarios-dir", "demos/vhs/scenarios/", "Directory containing VHS-only .feature files"),
		outputDir:      fs.String("output", "", "Output directory (required)"),
		configSource:   fs.String("config-source", "demos/vhs/config.tape", "Path to config tape file"),
		verbose:        fs.Bool("verbose", false, "Verbose output"),
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(errOut, "Error parsing flags: %v\n", err)
		return nil, err
	}

	if *opts.outputDir == "" {
		fmt.Fprintf(errOut, "Error: --output is required\n")
		return nil, errors.New("output required")
	}

	if !*opts.generateAll && *opts.featureFilter == "" && *opts.scenarioFilter == "" {
		fmt.Fprintf(errOut, "Error: one of --all, --feature, or --scenario is required\n")
		return nil, errors.New("filter required")
	}

	return opts, nil
}

func parseAllScenarios(featuresDir, scenariosDir string, errOut io.Writer) ([]vhsgen.ScenarioIR, error) {
	if _, err := os.Stat(featuresDir); err != nil {
		fmt.Fprintf(errOut, "Error parsing features dir %q: %v\n", featuresDir, err)
		return nil, err
	}

	businessScenarios, err := vhsgen.ParseFeatureDir(featuresDir, vhsgen.SourceBusiness)
	if err != nil {
		fmt.Fprintf(errOut, "Error parsing features dir %q: %v\n", featuresDir, err)
		return nil, err
	}

	vhsOnlyScenarios, err := vhsgen.ParseFeatureDir(scenariosDir, vhsgen.SourceVHSOnly)
	if err != nil {
		fmt.Fprintf(errOut, "Error parsing scenarios dir %q: %v\n", scenariosDir, err)
		return nil, err
	}

	allScenarios := make([]vhsgen.ScenarioIR, 0, len(businessScenarios)+len(vhsOnlyScenarios))
	allScenarios = append(allScenarios, businessScenarios...)
	allScenarios = append(allScenarios, vhsOnlyScenarios...)
	return allScenarios, nil
}

type scenarioWithResult struct {
	scenario vhsgen.ScenarioIR
	result   vhsgen.AnalysisResult
}

type generateStats struct {
	total        int
	fromBusiness int
	fromVHSOnly  int
	warnings     int
}

type generateConfig struct {
	outputDir    string
	configSource string
	verbose      bool
	out          io.Writer
	errOut       io.Writer
}

// generateTapes processes filtered scenarios and writes tape files.
func generateTapes(filtered []scenarioWithResult, cfg generateConfig) generateStats {
	var stats generateStats

	for i := range filtered {
		entry := &filtered[i]
		scenario := entry.scenario
		result := entry.result

		if !result.Translatable {
			if cfg.verbose {
				fmt.Fprintf(cfg.out, "Skipping %q (not translatable)\n", scenario.Name)
			}
			stats.warnings++
			continue
		}

		outPath, tapeErr := writeScenarioTape(scenario, cfg.outputDir, cfg.configSource)
		if tapeErr != nil {
			fmt.Fprintf(cfg.errOut, "Error generating tape for %q: %v\n", scenario.Name, tapeErr)
			continue
		}

		fmt.Fprintf(cfg.out, "Written: %s\n", outPath)

		switch scenario.Source {
		case vhsgen.SourceBusiness:
			stats.fromBusiness++
		case vhsgen.SourceVHSOnly:
			stats.fromVHSOnly++
		}
	}

	stats.total = stats.fromBusiness + stats.fromVHSOnly
	return stats
}

// filterResults selects scenarios to generate based on flags.
func filterResults(
	results []vhsgen.AnalysisResult,
	scenarios []vhsgen.ScenarioIR,
	all bool,
	featureFilter, scenarioFilter string,
) []scenarioWithResult {
	var out []scenarioWithResult

	resultByName := make(map[string]vhsgen.AnalysisResult, len(results))
	for i := range results {
		r := &results[i]
		resultByName[r.ScenarioName] = *r
	}

	for i := range scenarios {
		s := &scenarios[i]
		result, ok := resultByName[s.Name]
		if !ok {
			continue
		}

		if !all {
			if featureFilter != "" && !strings.EqualFold(s.Feature, featureFilter) {
				continue
			}
			if scenarioFilter != "" && !strings.EqualFold(s.Name, scenarioFilter) {
				continue
			}
		}

		out = append(out, scenarioWithResult{scenario: *s, result: result})
	}

	return out
}

// writeScenarioTape generates and writes a tape file with source-aware routing:
// Business tapes → {output}/{feature-slug}/{scenario-slug}.tape.
// VHS-only tapes → {output}/scenarios/{subdirectory}/{scenario-slug}.tape.
func writeScenarioTape(scenario vhsgen.ScenarioIR, outputDir, configSourcePath string) (string, error) {
	featureSlug := slugify(scenario.Feature)
	scenarioSlug := slugify(scenario.Name)

	var tapeDir string
	switch scenario.Source {
	case vhsgen.SourceVHSOnly:
		tapeDir = filepath.Join(outputDir, "scenarios", featureSlug)
	default:
		tapeDir = filepath.Join(outputDir, featureSlug)
	}

	config := vhsgen.GeneratorConfig{
		OutputDir:        outputDir,
		ConfigSourcePath: configSourcePath,
	}

	content, err := vhsgen.GenerateTape(scenario, config)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(tapeDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory %q: %w", tapeDir, err)
	}

	outPath := filepath.Join(tapeDir, scenarioSlug+".tape")
	if err := os.WriteFile(outPath, []byte(content), 0o600); err != nil {
		return "", fmt.Errorf("writing tape file %q: %w", outPath, err)
	}

	return outPath, nil
}

var (
	slugStripRe    = regexp.MustCompile(`[^a-z0-9-]`)
	slugCollapseRe = regexp.MustCompile(`-{2,}`)
)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = slugStripRe.ReplaceAllString(s, "")
	s = slugCollapseRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func printUsageTo(out io.Writer) {
	fmt.Fprintln(out, "vhsgen — VHS tape generator for KaRiya")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, "  vhsgen list [flags]      List scenarios and their translatability")
	fmt.Fprintln(out, "  vhsgen generate [flags]  Generate VHS tape files from scenarios")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "list flags:")
	fmt.Fprintln(out, "  --features DIR       Directory with .feature files (default: features/)")
	fmt.Fprintln(out, "  --scenarios-dir DIR  Directory with VHS-only .feature files (default: demos/vhs/scenarios/)")
	fmt.Fprintln(out, "  --json               Output as JSON")
	fmt.Fprintln(out, "  --count              Show counts broken down by source")
	fmt.Fprintln(out, "  --steps              Show translatable step patterns")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "generate flags:")
	fmt.Fprintln(out, "  --all                Generate for all translatable scenarios")
	fmt.Fprintln(out, "  --feature NAME       Filter by feature name")
	fmt.Fprintln(out, "  --scenario NAME      Filter by scenario name")
	fmt.Fprintln(out, "  --features DIR       Directory with .feature files (default: features/)")
	fmt.Fprintln(out, "  --scenarios-dir DIR  Directory with VHS-only .feature files (default: demos/vhs/scenarios/)")
	fmt.Fprintln(out, "  --output DIR         Output directory (required)")
	fmt.Fprintln(out, "  --config-source PATH Path to config tape file (default: demos/vhs/config.tape)")
	fmt.Fprintln(out, "  --verbose            Verbose output")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintln(out, "  vhsgen list --features features/ --scenarios-dir demos/vhs/scenarios/")
	fmt.Fprintln(out, "  vhsgen list --json")
	fmt.Fprintln(out, "  vhsgen list --count")
	fmt.Fprintln(out, "  vhsgen list --steps")
	fmt.Fprintln(out, "  vhsgen list --steps --json")
	fmt.Fprintln(out, "  vhsgen generate --all --features features/ --scenarios-dir demos/vhs/scenarios/ --output /tmp/tapes/")
	fmt.Fprintln(out, "  vhsgen generate --feature onboarding --output /tmp/test/")
}
