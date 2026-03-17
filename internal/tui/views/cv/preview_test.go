package cv_test

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/cv"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

var _ = Describe("Preview", func() {
	var (
		sampleCV      display.CVView
		profileConfig *config.ProfileConfig
	)

	BeforeEach(func() {
		bullet1 := fixtures.CVBulletWith("b1", "s1", "Built scalable systems")
		bullet1.Confidence = 0.9
		bullet2 := fixtures.CVBulletWith("b2", "s1", "Improved performance by 40%")
		bullet2.Confidence = 0.85
		group1 := fixtures.ContentGroupWithBullets("Company A", []*career.CVBullet{bullet1, bullet2})

		section1 := fixtures.CVSectionWithContent("s1", "cv-prev", []*career.SectionContentGroup{group1})
		section1.Title = "Experience"
		section1.SectionType = "experience"

		summarySection := fixtures.CVSectionWithSummary("s-summary", "cv-prev", "Professional summary text for preview testing")
		summarySection.Title = "Summary"

		sampleCVDomain := fixtures.CVViewWithSections("cv-prev", []*career.CVSection{summarySection, section1})
		sampleCVDomain.Name = "Test CV"
		sampleCVDomain.TargetRole = "Engineer"
		sampleCVDomain.TargetAudience = "Hiring Manager"
		sampleCV = display.CVViewFromDomain(sampleCVDomain)

		profileConfig = &config.ProfileConfig{
			Name:     "Test User",
			Email:    "test@example.com",
			Location: "Remote",
		}
	})

	Describe("Construction", func() {
		It("should create Preview via NewPreview", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			Expect(p).ToNot(BeNil())
			Expect(p.GetCV()).To(Equal(sampleCV))
		})

		It("should handle nil CV without panic", func() {
			Expect(func() {
				p := cv.NewPreview(display.CVView{}, profileConfig)
				Expect(p).ToNot(BeNil())
			}).NotTo(Panic())
		})
	})

	Describe("Init", func() {
		It("should return nil", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			Expect(p.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		var view *cv.Preview

		BeforeEach(func() {
			view = cv.NewPreview(sampleCV, profileConfig)
		})

		It("should handle WindowSizeMsg", func() {
			cmd, result := view.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(view.GetTerminalWidth()).To(Equal(120))
			Expect(view.GetTerminalHeight()).To(Equal(40))
		})

		It("should return CancelViewResult on esc", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(widgets.ResultCancel))
		})

		It("should return CancelViewResult type on esc", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_, ok := result.(*widgets.CancelViewResult)
			Expect(ok).To(BeTrue())
		})

		It("should return NavigateViewResult for confirm on enter", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionConfirm))
			Expect(payload.CV).To(Equal(sampleCV))
		})

		It("should return NavigateViewResult for confirm on y key", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionConfirm))
		})

		It("should return NavigateViewResult for edit on e key", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionEdit))
		})

		It("should return NavigateViewResult for export on x key", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionExport))
		})

		It("should go to top on g key", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should go to bottom on G key", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return nil for scroll keys when viewport not ready", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should handle scroll keys when viewport is ready", func() {
			view.RenderContent()
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(result).To(BeNil())
		})

		It("should handle up key when viewport is ready", func() {
			view.RenderContent()
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(result).To(BeNil())
		})

		It("should handle down key when viewport is ready", func() {
			view.RenderContent()
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(result).To(BeNil())
		})

		DescribeTable("should handle scroll keys when viewport is ready", func(msg tea.KeyMsg) {
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view.RenderContent()
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		},
			Entry("up", tea.KeyMsg{Type: tea.KeyUp}),
			Entry("down", tea.KeyMsg{Type: tea.KeyDown}),
			Entry("k", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}),
			Entry("pgup", tea.KeyMsg{Type: tea.KeyPgUp}),
			Entry("pgdown", tea.KeyMsg{Type: tea.KeyPgDown}),
			Entry("ctrl+u", tea.KeyMsg{Type: tea.KeyCtrlU}),
			Entry("ctrl+d", tea.KeyMsg{Type: tea.KeyCtrlD}),
		)

		It("should return nil for unhandled messages", func() {
			cmd, result := view.Update(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("GetCV", func() {
		It("should return the CV data", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			Expect(p.GetCV()).To(Equal(sampleCV))
		})

		It("should return nil when CV is nil", func() {
			p := cv.NewPreview(display.CVView{}, profileConfig)
			Expect(p.GetCV()).To(Equal(display.CVView{}))
		})
	})

	Describe("RenderContent", func() {
		It("should render CV Preview title", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("CV Preview"))
		})

		It("should render summary section text", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Professional summary text"))
		})

		It("should render with nil CV", func() {
			p := cv.NewPreview(display.CVView{}, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("No CV data available"))
		})

		It("should render with empty sections", func() {
			emptyCVDomain := fixtures.CVViewWithSections("cv-empty", []*career.CVSection{})
			emptyCVDomain.Name = "Empty CV"
			emptyCVDomain.TargetRole = "Engineer"
			emptyCVDomain.TargetAudience = "Manager"
			emptyCV := display.CVViewFromDomain(emptyCVDomain)

			p := cv.NewPreview(emptyCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("No sections generated"))
		})

		It("should render content groups with same start and end date", func() {
			bullet := fixtures.CVBulletWith("b-date", "s-date", "Worked on project")
			group := fixtures.ContentGroupFull("Company B", "2023-01", "2023-01", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s-date", "cv-date", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			dateCVDomain := fixtures.CVViewWithSections("cv-date", []*career.CVSection{section})
			dateCVDomain.Name = "Date CV"
			dateCVDomain.TargetRole = "Engineer"
			dateCVDomain.TargetAudience = "Manager"
			dateCV := display.CVViewFromDomain(dateCVDomain)

			p := cv.NewPreview(dateCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Company B (2023-01)"))
		})

		It("should render content groups with different start and end date", func() {
			bullet := fixtures.CVBulletWith("b-range", "s-range", "Long tenure project")
			group := fixtures.ContentGroupFull("Company C", "2020-01", "2023-06", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s-range", "cv-range", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			rangeCVDomain := fixtures.CVViewWithSections("cv-range", []*career.CVSection{section})
			rangeCVDomain.Name = "Range CV"
			rangeCVDomain.TargetRole = "Engineer"
			rangeCVDomain.TargetAudience = "Manager"
			rangeCV := display.CVViewFromDomain(rangeCVDomain)

			p := cv.NewPreview(rangeCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Company C (2020-01 - 2023-06)"))
		})

		It("should score bullets using audience relevance in highlights", func() {
			bullet := fixtures.CVBulletWith("b-aud", "s1", "Audience-relevant bullet")
			bullet.Confidence = 0.8
			bullet.AudienceRelevance = map[string]float64{"hiring manager": 0.95}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-aud", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			audienceCVDomain := fixtures.CVViewWithSections("cv-aud", []*career.CVSection{section})
			audienceCVDomain.Name = "Audience CV"
			audienceCVDomain.TargetRole = "Engineer"
			audienceCVDomain.TargetAudience = "Hiring Manager"
			audienceCV := display.CVViewFromDomain(audienceCVDomain)

			p := cv.NewPreview(audienceCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Audience-relevant bullet"))
		})

		It("should fall back to confidence when audience key not found", func() {
			bullet := fixtures.CVBulletWith("b-nokey", "s1", "No key bullet")
			bullet.Confidence = 0.75
			bullet.AudienceRelevance = map[string]float64{"recruiter": 0.5}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-nokey", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			nokeyCVDomain := fixtures.CVViewWithSections("cv-nokey", []*career.CVSection{section})
			nokeyCVDomain.Name = "No Key CV"
			nokeyCVDomain.TargetRole = "Engineer"
			nokeyCVDomain.TargetAudience = "Hiring Manager"
			nokeyCV := display.CVViewFromDomain(nokeyCVDomain)

			p := cv.NewPreview(nokeyCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("No key bullet"))
		})

		It("should use confidence when audience is master", func() {
			bullet := fixtures.CVBulletWith("b-master", "s1", "Master bullet")
			bullet.Confidence = 0.9
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-master", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			masterCVDomain := fixtures.CVViewWithSections("cv-master", []*career.CVSection{section})
			masterCVDomain.Name = "Master CV"
			masterCVDomain.TargetRole = "Engineer"
			masterCVDomain.TargetAudience = "master"
			masterCV := display.CVViewFromDomain(masterCVDomain)

			p := cv.NewPreview(masterCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Master bullet"))
		})

		It("should use confidence when audience is empty", func() {
			bullet := fixtures.CVBulletWith("b-empty-aud", "s1", "Empty audience bullet")
			bullet.Confidence = 0.8
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-emptyaud", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			emptyAudCVDomain := fixtures.CVViewWithSections("cv-emptyaud", []*career.CVSection{section})
			emptyAudCVDomain.Name = "Empty Audience CV"
			emptyAudCVDomain.TargetRole = "Engineer"
			emptyAudCVDomain.TargetAudience = ""
			emptyAudCV := display.CVViewFromDomain(emptyAudCVDomain)

			p := cv.NewPreview(emptyAudCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Empty audience bullet"))
		})

		It("should use enhanced text when text is empty in highlights", func() {
			bullet := fixtures.CVBulletWith("b-enh", "s1", "")
			bullet.EnhancedText = "Enhanced highlight text"
			bullet.Confidence = 0.9
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-enh", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			enhCVDomain := fixtures.CVViewWithSections("cv-enh", []*career.CVSection{section})
			enhCVDomain.Name = "Enhanced CV"
			enhCVDomain.TargetRole = "Engineer"
			enhCVDomain.TargetAudience = "Manager"
			enhCV := display.CVViewFromDomain(enhCVDomain)

			p := cv.NewPreview(enhCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Enhanced highlight text"))
		})

		It("should respect MaxHighlights from profile config", func() {
			highlightProfile := &config.ProfileConfig{
				Name:          "Test",
				Email:         "test@example.com",
				Location:      "Remote",
				MaxHighlights: 1,
			}

			p := cv.NewPreview(sampleCV, highlightProfile)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("Key Highlights"))
		})

		It("should render with nil sections for highlights", func() {
			noSectionsCVDomain := fixtures.CVViewWithSections("cv-nosec", []*career.CVSection{})
			noSectionsCVDomain.Name = "No Sections CV"
			noSectionsCVDomain.TargetRole = "Engineer"
			noSectionsCVDomain.TargetAudience = "Manager"
			noSectionsCV := display.CVViewFromDomain(noSectionsCVDomain)

			p := cv.NewPreview(noSectionsCV, profileConfig)
			output := p.RenderContent()
			Expect(output).NotTo(ContainSubstring("Key Highlights"))
		})

		It("should render with set theme", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			p.SetTheme(themes.NewDefaultTheme())
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("CV Preview"))
		})

		It("should render with small viewport dimensions", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			p.Update(tea.WindowSizeMsg{Width: 20, Height: 3})
			output := p.RenderContent()
			Expect(output).ToNot(BeEmpty())
		})

		It("should render line count for tall content", func() {
			var sections []*career.CVSection
			for i := range 5 {
				var bullets []*career.CVBullet
				for j := range 10 {
					b := fixtures.CVBulletWith(
						"b-tall-"+string(rune('a'+i))+string(rune('0'+j)),
						"s-tall",
						"This is a long bullet point that creates tall content for testing scroll behavior in the viewport",
					)
					b.Confidence = 0.8
					bullets = append(bullets, b)
				}
				group := fixtures.ContentGroupWithBullets("Company "+string(rune('A'+i)), bullets)
				sec := fixtures.CVSectionWithContent("s-tall-"+string(rune('a'+i)), "cv-tall", []*career.SectionContentGroup{group})
				sec.Title = "Section " + string(rune('A'+i))
				sec.SectionType = "experience"
				sections = append(sections, sec)
			}

			tallCVDomain := fixtures.CVViewWithSections("cv-tall", sections)
			tallCVDomain.Name = "Tall CV"
			tallCVDomain.TargetRole = "Engineer"
			tallCVDomain.TargetAudience = "Manager"
			tallCV := display.CVViewFromDomain(tallCVDomain)

			p := cv.NewPreview(tallCV, profileConfig)
			p.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
			output := p.RenderContent()
			Expect(output).ToNot(BeEmpty())
			Expect(output).To(ContainSubstring("CV Preview"))
		})

		It("should include line count when content exceeds viewport", func() {
			var bullets []*career.CVBullet
			for i := range 12 {
				bullet := fixtures.CVBulletWith(
					"b-lines-"+string(rune('a'+i)),
					"s-lines",
					"Long bullet line to create content that exceeds the viewport height",
				)
				bullet.Confidence = 0.8
				bullets = append(bullets, bullet)
			}
			group := fixtures.ContentGroupWithBullets("Company Lines", bullets)
			section := fixtures.CVSectionWithContent("s-lines", "cv-lines", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			lineCVDomain := fixtures.CVViewWithSections("cv-lines", []*career.CVSection{section})
			lineCVDomain.Name = "Line Count CV"
			lineCVDomain.TargetRole = "Engineer"
			lineCVDomain.TargetAudience = "Manager"
			lineCV := display.CVViewFromDomain(lineCVDomain)

			p := cv.NewPreview(lineCV, profileConfig)
			p.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
			p.RenderContent()
			p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("lines)"))
		})

		It("should render with nil sections", func() {
			nilSectionsCVDomain := fixtures.CVViewWithSections("cv-nil", nil)
			nilSectionsCVDomain.Name = "Nil Sections CV"
			nilSectionsCVDomain.TargetRole = "Engineer"
			nilSectionsCVDomain.TargetAudience = "Manager"
			nilSectionsCV := display.CVViewFromDomain(nilSectionsCVDomain)

			p := cv.NewPreview(nilSectionsCV, profileConfig)
			output := p.RenderContent()
			Expect(output).To(ContainSubstring("No sections generated yet"))
			Expect(output).NotTo(ContainSubstring("Key Highlights"))
		})
	})

	Describe("HelpText", func() {
		It("should return footer with help text", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			footer := p.HelpText()
			Expect(footer).To(ContainSubstring("scroll"))
		})

		It("should render basic help when viewport is not ready", func() {
			p := cv.NewPreview(sampleCV, profileConfig)
			footer := p.HelpText()
			Expect(footer).NotTo(ContainSubstring("%"))
			Expect(footer).To(ContainSubstring("enter/y: confirm"))
		})

		It("should include scroll percentage when viewport is ready with tall content", func() {
			var bullets []*career.CVBullet
			for i := range 50 {
				b := fixtures.CVBulletWith(
					"b-help-"+string(rune('a'+i%26)),
					"s-help",
					"Bullet for help text testing with enough content to require scrolling",
				)
				b.Confidence = 0.8
				bullets = append(bullets, b)
			}
			group := fixtures.ContentGroupWithBullets("Company", bullets)
			sec := fixtures.CVSectionWithContent("s-help", "cv-help", []*career.SectionContentGroup{group})
			sec.Title = "Experience"
			sec.SectionType = "experience"

			tallCVDomain := fixtures.CVViewWithSections("cv-help", []*career.CVSection{sec})
			tallCVDomain.Name = "Help CV"
			tallCVDomain.TargetRole = "Engineer"
			tallCVDomain.TargetAudience = "Manager"
			tallCV := display.CVViewFromDomain(tallCVDomain)

			p := cv.NewPreview(tallCV, profileConfig)
			p.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
			p.RenderContent()
			footer := p.HelpText()
			Expect(footer).To(ContainSubstring("%"))
		})
	})

	var _ = Describe("Boundary: cv view import rules", func() {
		forbidden := []string{
			"github.com/baphled/kariya/internal/domain/career",
			"github.com/baphled/kariya/internal/service/career",
			"github.com/baphled/kariya/internal/tui/intents",
		}

		viewDir := "./"
		files, err := filepath.Glob(filepath.Join(viewDir, "*.go"))
		It("should not error when globbing", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") || strings.HasSuffix(file, "suite_test.go") {
				continue
			}
			It("should not import forbidden packages in "+file, func() {
				content, err := os.ReadFile(file)
				Expect(err).ToNot(HaveOccurred())
				for _, f := range forbidden {
					Expect(string(content)).NotTo(ContainSubstring(f), "Forbidden import: %s in %s", f, file)
				}
			})
		}
	})
})
