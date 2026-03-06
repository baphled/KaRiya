package importcmd_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	importcmd "github.com/baphled/kariya/internal/cli/cmd/import"
	"github.com/baphled/kariya/internal/cli/importer"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Formatter Functions", func() {
	Describe("FormatFailedRows", func() {
		Context("when there are no failures", func() {
			It("should return empty string", func() {
				result := &importer.ImportResult{
					FailedCount: 0,
					FailedRows:  nil,
				}

				output := importcmd.FormatFailedRows(result)
				Expect(output).To(BeEmpty())
			})
		})

		Context("when there are failures", func() {
			It("should show count and first 5 rows", func() {
				result := &importer.ImportResult{
					FailedCount: 3,
					FailedRows: []*importer.ParsedRow{
						{Event: fixtures.EventWith("e1", "Failed row 1", "", "")},
						{Event: fixtures.EventWith("e2", "Failed row 2", "", "")},
						{Event: fixtures.EventWith("e3", "Failed row 3", "", "")},
					},
				}

				output := importcmd.FormatFailedRows(result)
				Expect(output).To(ContainSubstring("Failed: 3"))
				Expect(output).To(ContainSubstring("Failed row 1"))
				Expect(output).To(ContainSubstring("Failed row 2"))
				Expect(output).To(ContainSubstring("Failed row 3"))
			})

			It("should limit to 5 rows when more exist", func() {
				rows := make([]*importer.ParsedRow, 10)
				for i := range rows {
					rows[i] = &importer.ParsedRow{
						Event: fixtures.EventWith("e"+string(rune(i)), "Failed row", "", ""),
					}
				}
				result := &importer.ImportResult{
					FailedCount: 10,
					FailedRows:  rows,
				}

				output := importcmd.FormatFailedRows(result)
				Expect(output).To(ContainSubstring("Failed: 10"))
				Expect(output).To(ContainSubstring("first 5"))
			})
		})

		Context("when FailedCount > 0 but FailedRows is empty", func() {
			It("should show count only", func() {
				result := &importer.ImportResult{
					FailedCount: 5,
					FailedRows:  []*importer.ParsedRow{},
				}

				output := importcmd.FormatFailedRows(result)
				Expect(output).To(Equal("Failed: 5\n"))
				Expect(output).NotTo(ContainSubstring("Failed rows"))
			})
		})
	})

	Describe("FormatBurstSuggestions", func() {
		Context("when there are no suggestions", func() {
			It("should return empty string", func() {
				result := &importer.ImportResult{
					BurstSuggestions: []burstfact.BurstSuggestion{},
				}

				output := importcmd.FormatBurstSuggestions(result)
				Expect(output).To(BeEmpty())
			})
		})

		Context("when there are suggestions", func() {
			It("should format all suggestions", func() {
				result := &importer.ImportResult{
					BurstSuggestions: []burstfact.BurstSuggestion{
						{
							Name:            "Q1 Initiative",
							EventIDs:        []string{"e1", "e2"},
							ConfidenceScore: 0.85,
						},
						{
							Name:            "Q2 Migration",
							EventIDs:        []string{"e3", "e4", "e5"},
							ConfidenceScore: 0.92,
						},
					},
				}

				output := importcmd.FormatBurstSuggestions(result)
				Expect(output).To(ContainSubstring("=== Burst Suggestions ==="))
				Expect(output).To(ContainSubstring("Detected 2 potential bursts"))
				Expect(output).To(ContainSubstring("Q1 Initiative"))
				Expect(output).To(ContainSubstring("Events: 2 | Confidence: 85.0%"))
				Expect(output).To(ContainSubstring("Q2 Migration"))
				Expect(output).To(ContainSubstring("Events: 3 | Confidence: 92.0%"))
			})

			It("should handle unnamed bursts", func() {
				result := &importer.ImportResult{
					BurstSuggestions: []burstfact.BurstSuggestion{
						{
							Name:            "",
							EventIDs:        []string{"e1", "e2"},
							ConfidenceScore: 0.75,
						},
					},
				}

				output := importcmd.FormatBurstSuggestions(result)
				Expect(output).To(ContainSubstring("1. Burst 1"))
				Expect(output).To(ContainSubstring("Events: 2 | Confidence: 75.0%"))
			})
		})

		Context("with more than 5 suggestions", func() {
			It("should show all suggestions", func() {
				suggestions := make([]burstfact.BurstSuggestion, 7)
				for i := range suggestions {
					suggestions[i] = burstfact.BurstSuggestion{
						Name:            "Burst " + string(rune('A'+i)),
						EventIDs:        []string{"e1", "e2"},
						ConfidenceScore: 0.8,
					}
				}
				result := &importer.ImportResult{
					BurstSuggestions: suggestions,
				}

				output := importcmd.FormatBurstSuggestions(result)
				Expect(output).To(ContainSubstring("Detected 7 potential bursts"))
			})
		})
	})

	Describe("FormatFactExtraction", func() {
		Context("when there are no facts", func() {
			It("should return empty string", func() {
				result := &importer.ImportResult{
					ExtractedFactsCount: 0,
				}

				output := importcmd.FormatFactExtraction(result)
				Expect(output).To(BeEmpty())
			})
		})

		Context("when there are facts without competencies", func() {
			It("should show count only", func() {
				result := &importer.ImportResult{
					ExtractedFactsCount: 10,
					CreatedEvents:       nil,
					FactsByCompetency:   map[string]int{},
				}

				output := importcmd.FormatFactExtraction(result)
				Expect(output).To(ContainSubstring("=== Fact Extraction ==="))
				Expect(output).To(ContainSubstring("Extracted 10 facts from 0 events"))
				Expect(output).NotTo(ContainSubstring("Competency breakdown"))
			})
		})

		Context("when there are facts with competencies", func() {
			It("should show sorted competency breakdown", func() {
				result := &importer.ImportResult{
					ExtractedFactsCount: 15,
					CreatedEvents:       nil,
					FactsByCompetency: map[string]int{
						"Technical":    8,
						"Leadership":   5,
						"Architecture": 2,
					},
				}

				output := importcmd.FormatFactExtraction(result)
				Expect(output).To(ContainSubstring("=== Fact Extraction ==="))
				Expect(output).To(ContainSubstring("Extracted 15 facts"))
				Expect(output).To(ContainSubstring("Competency breakdown:"))
				Expect(output).To(ContainSubstring("- Architecture: 2 facts"))
				Expect(output).To(ContainSubstring("- Leadership: 5 facts"))
				Expect(output).To(ContainSubstring("- Technical: 8 facts"))

				// Verify alphabetical order
				architectureIdx := bytes.Index([]byte(output), []byte("Architecture"))
				leadershipIdx := bytes.Index([]byte(output), []byte("Leadership"))
				technicalIdx := bytes.Index([]byte(output), []byte("Technical"))
				Expect(architectureIdx).To(BeNumerically("<", leadershipIdx))
				Expect(leadershipIdx).To(BeNumerically("<", technicalIdx))
			})
		})

		Context("with single competency", func() {
			It("should format correctly", func() {
				result := &importer.ImportResult{
					ExtractedFactsCount: 5,
					FactsByCompetency: map[string]int{
						"Technical": 5,
					},
				}

				output := importcmd.FormatFactExtraction(result)
				Expect(output).To(ContainSubstring("- Technical: 5 facts"))
			})
		})
	})

	Describe("CalculateExitCode", func() {
		Context("when events were successfully imported", func() {
			It("should return 0", func() {
				result := &importer.ImportResult{
					SuccessCount: 5,
				}

				exitCode := importcmd.CalculateExitCode(result)
				Expect(exitCode).To(Equal(0))
			})
		})

		Context("when events were skipped", func() {
			It("should return 0", func() {
				result := &importer.ImportResult{
					SuccessCount: 0,
					SkippedCount: 3,
				}

				exitCode := importcmd.CalculateExitCode(result)
				Expect(exitCode).To(Equal(0))
			})
		})

		Context("when no events were processed", func() {
			It("should return 1", func() {
				result := &importer.ImportResult{
					SuccessCount: 0,
					SkippedCount: 0,
				}

				exitCode := importcmd.CalculateExitCode(result)
				Expect(exitCode).To(Equal(1))
			})
		})

		Context("when there were only failures", func() {
			It("should return 1", func() {
				result := &importer.ImportResult{
					SuccessCount: 0,
					SkippedCount: 0,
					FailedCount:  5,
				}

				exitCode := importcmd.CalculateExitCode(result)
				Expect(exitCode).To(Equal(1))
			})
		})
	})

	Describe("DisplayBurstSuggestions", func() {
		It("should write formatted output to writer", func() {
			result := &importer.ImportResult{
				BurstSuggestions: []burstfact.BurstSuggestion{
					{
						Name:            "Test Burst",
						EventIDs:        []string{"e1", "e2"},
						ConfidenceScore: 0.9,
					},
				},
			}

			var buf bytes.Buffer
			importcmd.DisplayBurstSuggestions(&buf, result)

			output := buf.String()
			Expect(output).To(ContainSubstring("Test Burst"))
			Expect(output).To(ContainSubstring("Events: 2 | Confidence: 90.0%"))
		})
	})

	Describe("DisplayFactExtraction", func() {
		It("should write formatted output to writer", func() {
			result := &importer.ImportResult{
				ExtractedFactsCount: 5,
				FactsByCompetency: map[string]int{
					"Technical": 5,
				},
			}

			var buf bytes.Buffer
			importcmd.DisplayFactExtraction(&buf, result)

			output := buf.String()
			Expect(output).To(ContainSubstring("=== Fact Extraction ==="))
			Expect(output).To(ContainSubstring("Extracted 5 facts"))
			Expect(output).To(ContainSubstring("- Technical: 5 facts"))
		})
	})

	Describe("DisplayFailedRows", func() {
		It("should write formatted output to writer", func() {
			result := &importer.ImportResult{
				FailedCount: 2,
				FailedRows: []*importer.ParsedRow{
					{Event: fixtures.EventWith("e1", "Failed 1", "", "")},
					{Event: fixtures.EventWith("e2", "Failed 2", "", "")},
				},
			}

			var buf bytes.Buffer
			importcmd.DisplayFailedRows(&buf, result)

			output := buf.String()
			Expect(output).To(ContainSubstring("Failed: 2"))
			Expect(output).To(ContainSubstring("Failed 1"))
		})
	})
})
