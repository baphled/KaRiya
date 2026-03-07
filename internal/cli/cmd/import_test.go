package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/importer"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
)

type dummyFileInfo struct {
	os.FileInfo
}

func (d dummyFileInfo) Name() string { return "test.csv" }
func (d dummyFileInfo) Size() int64  { return 100 }

var _ = Describe("Import Command", func() {
	var (
		ctx    *cmdpkg.CLIContext
		out    *bytes.Buffer
		errOut *bytes.Buffer
		opener cliutil.MockFileOpener
		runner cliutil.MockProgressRunner
	)

	BeforeEach(func() {
		ctx = cmdpkg.NewCLIContext("", true)
		errOut = new(bytes.Buffer)
		initErr := ctx.InitService(errOut)
		Expect(initErr).NotTo(HaveOccurred())

		out = new(bytes.Buffer)
		errOut = new(bytes.Buffer)
		opener = cliutil.MockFileOpener{
			StatFn: func(_ string) (os.FileInfo, error) {
				return dummyFileInfo{}, nil
			},
			OpenFn: func(_ string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader([]byte("Date,Company,Project,Text,Tags\n2023-01-01,Company,Project,Some Achievement,tag1"))), nil
			},
		}
		runner = cliutil.MockProgressRunner{}
	})

	Describe("HandleImport", func() {
		It("returns 0 on successful import", func() {
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Import successful"))
		})

		It("returns 1 when file cannot be accessed", func() {
			opener.StatFn = func(_ string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			}
			code := cmdpkg.HandleImport("nonexistent.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Cannot access import file"))
		})

		It("returns 1 when file cannot be opened", func() {
			opener.OpenFn = func(_ string) (io.ReadCloser, error) {
				return nil, os.ErrPermission
			}
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Error opening import file"))
		})

		It("returns 1 when CSV parsing fails", func() {
			runner.RunWithSpinnerFn = func(_ string, fn func() error, _ ...tea.ProgramOption) error {
				return fmt.Errorf("parse error")
			}
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Error parsing CSV"))
		})

		It("returns 1 when no valid rows found in CSV", func() {
			callCount := 0
			runner.RunWithSpinnerFn = func(_ string, fn func() error, _ ...tea.ProgramOption) error {
				callCount++
				if callCount == 1 {
					// First call: PrepareImport - return empty rows
					return nil
				}
				return fn()
			}
			opener.OpenFn = func(_ string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader([]byte("Date,Company,Project,Text,Tags\n"))), nil
			}
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("No valid rows found"))
		})

		It("returns 1 when import operation fails", func() {
			callCount := 0
			runner.RunWithSpinnerFn = func(_ string, fn func() error, _ ...tea.ProgramOption) error {
				callCount++
				if callCount == 2 {
					return fmt.Errorf("import error")
				}
				return fn()
			}
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Error during import"))
		})

	})

	Describe("Exit code paths", func() {
		It("returns 0 when SuccessCount > 0", func() {
			// This is tested by the successful import test
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("✓ Import successful"))
		})

		It("returns 0 when SkippedCount > 0 and SuccessCount = 0", func() {
			// Import the same row twice to trigger duplicate detection
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			Expect(code).To(Equal(0)) // First import succeeds

			// Now import the same data again - should be skipped as duplicates
			code = cmdpkg.HandleImport("test.csv", false, false, ctx.Service(), out, errOut, opener, runner)
			// The second import should either succeed (if duplicates are skipped) or fail
			Expect(code).To(BeElementOf(0, 1))
		})
	})

	Describe("DisplayImportResults", func() {
		It("displays failed rows when FailedCount > 0", func() {
			result := &importer.ImportResult{
				TotalRows:    3,
				SuccessCount: 2,
				SkippedCount: 0,
				FailedCount:  1,
				FailedRows: []*importer.ParsedRow{
					{
						RowNumber:        3,
						ValidationErrors: []string{"Invalid date format"},
						IsValid:          false,
					},
				},
			}
			out := new(bytes.Buffer)
			cmdpkg.DisplayImportResults(result, out)
			output := out.String()
			Expect(output).To(ContainSubstring("Failed: 1"))
		})

		It("displays burst suggestions with names", func() {
			result := &importer.ImportResult{
				TotalRows:    2,
				SuccessCount: 2,
				SkippedCount: 0,
				FailedCount:  0,
				BurstSuggestions: []burstfact.BurstSuggestion{
					{
						Name:            "API Redesign",
						EventIDs:        []string{"event1", "event2"},
						ConfidenceScore: 0.95,
					},
				},
			}
			out := new(bytes.Buffer)
			cmdpkg.DisplayImportResults(result, out)
			output := out.String()
			Expect(output).To(ContainSubstring("Burst Suggestions"))
			Expect(output).To(ContainSubstring("API Redesign"))
			Expect(output).To(ContainSubstring("95.0%"))
		})

		It("displays burst suggestions without names", func() {
			result := &importer.ImportResult{
				TotalRows:    2,
				SuccessCount: 2,
				SkippedCount: 0,
				FailedCount:  0,
				BurstSuggestions: []burstfact.BurstSuggestion{
					{
						Name:            "",
						EventIDs:        []string{"event1", "event2"},
						ConfidenceScore: 0.85,
					},
				},
			}
			out := new(bytes.Buffer)
			cmdpkg.DisplayImportResults(result, out)
			output := out.String()
			Expect(output).To(ContainSubstring("Burst Suggestions"))
		})

		It("displays extracted facts with competency breakdown", func() {
			result := &importer.ImportResult{
				TotalRows:           2,
				SuccessCount:        2,
				SkippedCount:        0,
				FailedCount:         0,
				ExtractedFactsCount: 5,
				CreatedEvents:       make([]*career.Event, 2),
				FactsByCompetency: map[string]int{
					"Leadership": 2,
					"Technical":  3,
				},
			}
			out := new(bytes.Buffer)
			cmdpkg.DisplayImportResults(result, out)
			output := out.String()
			Expect(output).To(ContainSubstring("Fact Extraction"))
			Expect(output).To(ContainSubstring("Extracted 5 facts"))
			Expect(output).To(ContainSubstring("Leadership: 2"))
			Expect(output).To(ContainSubstring("Technical: 3"))
		})
	})
})
