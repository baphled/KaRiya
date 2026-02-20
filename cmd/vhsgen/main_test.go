package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("vhsgen CLI", func() {
	Context("no subcommand", func() {
		It("prints usage and returns 0", func() {
			var out, errOut bytes.Buffer
			code := run([]string{}, &out, &errOut)

			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("vhsgen"))
			Expect(out.String()).To(ContainSubstring("list"))
			Expect(out.String()).To(ContainSubstring("generate"))
		})
	})

	Context("--help flag", func() {
		It("prints usage and returns 0", func() {
			var out, errOut bytes.Buffer
			code := run([]string{"--help"}, &out, &errOut)

			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("vhsgen"))
		})
	})

	Context("unknown subcommand", func() {
		It("returns exit code 1 with error message", func() {
			var out, errOut bytes.Buffer
			code := run([]string{"foobar"}, &out, &errOut)

			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("unknown subcommand"))
		})
	})

	Describe("list subcommand", func() {
		Context("with a valid features directory", func() {
			It("prints a table with Scenario, Feature, Source, Translatable, Reason columns", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--features", "../../features/", "--scenarios-dir", "../../demos/vhs/scenarios/"}, &out, &errOut)

				Expect(code).To(Equal(0))
				output := out.String()
				Expect(output).To(ContainSubstring("Parsing..."))
				Expect(output).To(ContainSubstring("Scenario"))
				Expect(output).To(ContainSubstring("Feature"))
				Expect(output).To(ContainSubstring("Source"))
				Expect(output).To(ContainSubstring("Translatable"))
			})
		})

		Context("--json flag", func() {
			It("produces valid JSON with a scenarios array containing source field", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--features", "../../features/", "--scenarios-dir", "../../demos/vhs/scenarios/", "--json"}, &out, &errOut)

				Expect(code).To(Equal(0))

				var payload map[string]interface{}
				outputBytes := out.Bytes()
				jsonStart := bytes.Index(outputBytes, []byte("{"))
				Expect(jsonStart).To(BeNumerically(">=", 0), "no JSON found in output")

				err := json.Unmarshal(outputBytes[jsonStart:], &payload)
				Expect(err).NotTo(HaveOccurred())
				Expect(payload).To(HaveKey("scenarios"))

				scenarios, ok := payload["scenarios"].([]interface{})
				Expect(ok).To(BeTrue())
				Expect(scenarios).NotTo(BeEmpty())

				first := scenarios[0].(map[string]interface{})
				Expect(first).To(HaveKey("source"))
				Expect(first).To(HaveKey("scenario_name"))
				Expect(first).To(HaveKey("feature"))
				Expect(first).To(HaveKey("translatable"))
			})
		})

		Context("--count flag", func() {
			It("shows counts broken down by source", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--features", "../../features/", "--scenarios-dir", "../../demos/vhs/scenarios/", "--count"}, &out, &errOut)

				Expect(code).To(Equal(0))
				Expect(out.String()).To(MatchRegexp(`Business: \d+/\d+ translatable \| VHS-only: \d+/\d+ translatable`))
			})
		})

		Context("--steps flag", func() {
			It("outputs a readable table of translatable step patterns", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--steps"}, &out, &errOut)

				Expect(code).To(Equal(0))
				output := out.String()
				Expect(output).To(ContainSubstring("Pattern"))
				Expect(output).To(ContainSubstring("Type"))
				Expect(output).To(ContainSubstring("Category"))
				Expect(output).To(ContainSubstring("Example"))
			})

			It("includes navigation patterns", func() {
				var out, errOut bytes.Buffer
				run([]string{"list", "--steps"}, &out, &errOut)

				Expect(out.String()).To(ContainSubstring("navigation"))
			})

			It("includes input patterns", func() {
				var out, errOut bytes.Buffer
				run([]string{"list", "--steps"}, &out, &errOut)

				Expect(out.String()).To(ContainSubstring("input"))
			})
		})

		Context("--steps --json flag", func() {
			It("outputs JSON array with pattern, type, category, params, example fields", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--steps", "--json"}, &out, &errOut)

				Expect(code).To(Equal(0))

				var patterns []map[string]interface{}
				err := json.Unmarshal(out.Bytes(), &patterns)
				Expect(err).NotTo(HaveOccurred())
				Expect(patterns).NotTo(BeEmpty())

				first := patterns[0]
				Expect(first).To(HaveKey("pattern"))
				Expect(first).To(HaveKey("type"))
				Expect(first).To(HaveKey("category"))
				Expect(first).To(HaveKey("example"))
			})
		})

		Context("with non-existent scenarios-dir", func() {
			It("handles missing scenarios dir gracefully (no error)", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--features", "../../features/", "--scenarios-dir", "/nonexistent/path/"}, &out, &errOut)

				Expect(code).To(Equal(0))
				Expect(errOut.String()).To(BeEmpty())
			})
		})

		Context("with non-existent features dir", func() {
			It("returns exit code 1 with error message", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"list", "--features", "/nonexistent/features/", "--scenarios-dir", "/nonexistent/scenarios/"}, &out, &errOut)

				Expect(code).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("Error parsing features dir"))
			})
		})
	})

	Describe("generate subcommand", func() {
		var tmpDir string

		BeforeEach(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "vhsgen-test-*")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
		})

		Context("--output missing", func() {
			It("returns exit code 1 with error message", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"generate", "--all", "--features", "../../features/"}, &out, &errOut)

				Expect(code).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("--output is required"))
			})
		})

		Context("no filter flags", func() {
			It("returns exit code 1 requiring --all, --feature, or --scenario", func() {
				var out, errOut bytes.Buffer
				code := run([]string{"generate", "--output", tmpDir}, &out, &errOut)

				Expect(code).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("--all"))
			})
		})

		Context("--all flag", func() {
			It("generates tape files and shows summary", func() {
				var out, errOut bytes.Buffer
				code := run([]string{
					"generate",
					"--all",
					"--features", "../../features/",
					"--scenarios-dir", "../../demos/vhs/scenarios/",
					"--output", tmpDir,
				}, &out, &errOut)

				Expect(code).To(Equal(0))
				output := out.String()
				Expect(output).To(ContainSubstring("Parsing..."))
				Expect(output).To(ContainSubstring("Generating..."))
				Expect(output).To(MatchRegexp(`Generated \d+ tapes`))
			})

			It("runs without error even when all scenarios are untranslatable", func() {
				var out, errOut bytes.Buffer
				code := run([]string{
					"generate",
					"--all",
					"--features", "../../features/",
					"--scenarios-dir", "../../demos/vhs/scenarios/",
					"--output", tmpDir,
				}, &out, &errOut)

				Expect(code).To(Equal(0))
				Expect(errOut.String()).To(BeEmpty())
			})

			It("reports 'Written: <path>' for each tape file written", func() {
				var out, errOut bytes.Buffer
				run([]string{
					"generate",
					"--all",
					"--features", "../../features/",
					"--scenarios-dir", "../../demos/vhs/scenarios/",
					"--output", tmpDir,
				}, &out, &errOut)

				output := out.String()
				lines := strings.Split(output, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Written:") {
						path := strings.TrimSpace(strings.TrimPrefix(line, "Written:"))
						Expect(path).To(HaveSuffix(".tape"))
					}
				}
			})
		})

		Context("--feature flag", func() {
			It("filters by feature name (case-insensitive) and shows summary", func() {
				var out, errOut bytes.Buffer
				code := run([]string{
					"generate",
					"--feature", "Onboarding",
					"--features", "../../features/",
					"--scenarios-dir", "../../demos/vhs/scenarios/",
					"--output", tmpDir,
				}, &out, &errOut)

				Expect(code).To(Equal(0))
				Expect(out.String()).To(MatchRegexp(`Generated \d+ tapes`))
			})
		})

		Context("--scenario flag", func() {
			It("filters by scenario name", func() {
				var out, errOut bytes.Buffer
				code := run([]string{
					"generate",
					"--scenario", "First launch with no existing data",
					"--features", "../../features/",
					"--scenarios-dir", "../../demos/vhs/scenarios/",
					"--output", tmpDir,
				}, &out, &errOut)

				Expect(code).To(Equal(0))
				Expect(out.String()).To(MatchRegexp(`Generated \d+ tapes`))
			})
		})

		Context("output summary format", func() {
			It("shows summary with from features, from scenarios, and warnings counts", func() {
				var out, errOut bytes.Buffer
				run([]string{
					"generate",
					"--all",
					"--features", "../../features/",
					"--scenarios-dir", "../../demos/vhs/scenarios/",
					"--output", tmpDir,
				}, &out, &errOut)

				Expect(out.String()).To(MatchRegexp(
					`Generated \d+ tapes \(\d+ from features, \d+ from scenarios, \d+ warnings\)`,
				))
			})
		})
	})

	Describe("slugify helper", func() {
		DescribeTable("converts names to URL-safe slugs",
			func(input, expected string) {
				Expect(slugify(input)).To(Equal(expected))
			},
			Entry("lowercase words", "hello world", "hello-world"),
			Entry("mixed case", "Capture Event", "capture-event"),
			Entry("underscores", "manage_skills", "manage-skills"),
			Entry("multiple spaces", "foo  bar", "foo-bar"),
			Entry("special chars stripped", "foo!@bar", "foobar"),
			Entry("leading/trailing hyphens", "-foo-bar-", "foo-bar"),
		)
	})

	Describe("truncate helper", func() {
		It("returns short strings unchanged", func() {
			Expect(truncate("hello", 10)).To(Equal("hello"))
		})

		It("truncates long strings with ellipsis", func() {
			result := truncate("hello world this is long", 10)
			Expect(result).To(HaveLen(10))
			Expect(result).To(HaveSuffix("..."))
		})

		It("handles max <= 3 without panic", func() {
			result := truncate("hello", 2)
			Expect(len(result)).To(BeNumerically("<=", 2))
		})
	})

	Describe("parseAllScenarios", func() {
		It("returns error for missing features directory", func() {
			var errOut bytes.Buffer
			_, err := parseAllScenarios("/nonexistent/features", "../../demos/vhs/scenarios/", &errOut)
			Expect(err).To(HaveOccurred())
		})

		It("parses scenarios successfully", func() {
			var errOut bytes.Buffer
			scenarios, err := parseAllScenarios("../../features/", "../../demos/vhs/scenarios/", &errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(scenarios).NotTo(BeEmpty())
		})
	})

	Describe("parseGenerateFlags", func() {
		It("returns error when output is missing", func() {
			var errOut bytes.Buffer
			_, err := parseGenerateFlags([]string{"--all"}, &errOut)
			Expect(err).To(HaveOccurred())
		})

		It("returns error when no filter is specified", func() {
			var errOut bytes.Buffer
			_, err := parseGenerateFlags([]string{"--output", "/tmp"}, &errOut)
			Expect(err).To(HaveOccurred())
		})

		It("parses valid flags successfully", func() {
			var errOut bytes.Buffer
			opts, err := parseGenerateFlags([]string{"--output", "/tmp", "--all"}, &errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(*opts.outputDir).To(Equal("/tmp"))
			Expect(*opts.generateAll).To(BeTrue())
		})
	})
})
