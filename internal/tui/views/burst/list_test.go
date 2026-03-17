package burst_test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		view   *burst.List
		bursts []display.Burst
	)

	BeforeEach(func() {
		b1 := fixtures.Burst("burst-1", "e1", "e2")
		b1.Name = "API Development Sprint"
		b1.Description = "Built REST APIs for user management"

		b2 := fixtures.BurstConfirmed("burst-2", "e3")
		b2.Name = "DevOps Migration"
		b2.Description = "Migrated infrastructure to Kubernetes"

		b3 := fixtures.Burst("burst-3", "e4", "e5", "e6")
		b3.Name = "Frontend Redesign"

		bursts = display.BurstsFromDomain([]*career.Burst{b1, b2, b3})
	})

	Describe("Construction", func() {
		It("should create a burst list view", func() {
			view = burst.NewList(bursts)
			Expect(view).NotTo(BeNil())
		})

		It("should store bursts", func() {
			view = burst.NewList(bursts)
			Expect(view.GetBursts()).To(Equal(bursts))
		})

		It("should handle empty burst list", func() {
			view = burst.NewList([]display.Burst{})
			Expect(view).NotTo(BeNil())
			Expect(view.GetBursts()).To(BeEmpty())
		})

		It("should handle nil burst list", func() {
			view = burst.NewList(nil)
			Expect(view).NotTo(BeNil())
		})

		It("should start with first item selected", func() {
			view = burst.NewList(bursts)
			Expect(view.GetSelectedIndex()).To(Equal(0))
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			view = burst.NewList(bursts)
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		BeforeEach(func() {
			view = burst.NewList(bursts)
		})

		It("should update terminal dimensions", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			cmd, result := view.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Update with KeyMsg", func() {
		BeforeEach(func() {
			view = burst.NewList(bursts)
			view.SetTerminalInfo(120, 40)
		})

		Context("navigation keys", func() {
			It("should move down with down arrow", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(view.GetSelectedIndex()).To(Equal(1))
			})

			It("should move up with up arrow", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyDown})
				view.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(view.GetSelectedIndex()).To(Equal(0))
			})

			It("should move down with j", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(view.GetSelectedIndex()).To(Equal(1))
			})

			It("should move up with k", func() {
				view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(view.GetSelectedIndex()).To(Equal(0))
			})
		})

		Context("selection", func() {
			It("should return NavigateViewResult with burst on Enter", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				nav, ok := result.Data().(burst.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(burst.ActionView))
				Expect(nav.Burst).NotTo(BeNil())
				Expect(nav.Burst.ID).To(Equal("burst-1"))
			})

			It("should return nil result on Enter with empty list", func() {
				view = burst.NewList([]display.Burst{})
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(result).To(BeNil())
			})
		})

		Context("cancel", func() {
			It("should return CancelViewResult on Escape", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultCancel))
			})
		})

		Context("action keys", func() {
			It("should return navigate result with add action on 'a'", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				nav, ok := result.Data().(burst.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(burst.ActionAdd))
			})

			It("should return navigate result with edit action on 'e'", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(result).NotTo(BeNil())
				nav, ok := result.Data().(burst.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(burst.ActionEdit))
				Expect(nav.Burst).NotTo(BeNil())
			})

			It("should return nil for edit with empty list", func() {
				view = burst.NewList([]display.Burst{})
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(result).To(BeNil())
			})

			It("should return navigate result with delete action on 'd'", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				Expect(result).NotTo(BeNil())
				nav, ok := result.Data().(burst.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(burst.ActionDelete))
				Expect(nav.Burst).NotTo(BeNil())
			})

			It("should return nil for delete with empty list", func() {
				view = burst.NewList([]display.Burst{})
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				Expect(result).To(BeNil())
			})

			It("should return navigate result with suggest action on 's'", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				Expect(result).NotTo(BeNil())
				Expect(result.Type()).To(Equal(widgets.ResultNavigate))
				nav, ok := result.Data().(burst.Nav)
				Expect(ok).To(BeTrue())
				Expect(nav.Action).To(Equal(burst.ActionSuggest))
			})
		})

		Context("unhandled keys", func() {
			It("should return nil for unhandled keys", func() {
				_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
				Expect(result).To(BeNil())
			})
		})

		Context("non-key messages", func() {
			It("should return nil for non-key messages", func() {
				cmd, result := view.Update(tea.MouseMsg{})
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RenderContent", func() {
		It("should return non-empty string", func() {
			view = burst.NewList(bursts)
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should render empty state when no bursts", func() {
			view = burst.NewList([]display.Burst{})
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("HelpText", func() {
		It("should return non-empty string", func() {
			view = burst.NewList(bursts)
			help := view.HelpText()
			Expect(help).NotTo(BeEmpty())
		})

		It("should use theme when set", func() {
			view = burst.NewList(bursts)
			view.SetTheme(nil)
			help := view.HelpText()
			Expect(help).NotTo(BeEmpty())
		})
	})

	Describe("SetTerminalInfo", func() {
		It("should store terminal dimensions", func() {
			view = burst.NewList(bursts)
			view.SetTerminalInfo(100, 50)
			Expect(view.GetTerminalWidth()).To(Equal(100))
			Expect(view.GetTerminalHeight()).To(Equal(50))
		})
	})

	Describe("SetTheme", func() {
		It("should store theme", func() {
			view = burst.NewList(bursts)
			view.SetTheme("test-theme")
			Expect(view.GetTheme()).To(Equal("test-theme"))
		})
	})

	Describe("SetLogo", func() {
		It("should store logo and spacing", func() {
			view = burst.NewList(bursts)
			view.SetLogo("test-logo", 2)
			Expect(view.GetLogo()).To(Equal("test-logo"))
			Expect(view.GetLogoSpacing()).To(Equal(2))
		})
	})

	Describe("burstRowFormatter", func() {
		It("should truncate long names", func() {
			b := fixtures.Burst("burst-long", "e1", "e2")
			b.Name = "This is a very long burst name that exceeds the column width limit"
			view = burst.NewList(display.BurstsFromDomain([]*career.Burst{b}))
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should truncate long descriptions", func() {
			b := fixtures.Burst("burst-desc", "e1", "e2")
			b.Description = "This is a very long description that should be truncated for display"
			view = burst.NewList(display.BurstsFromDomain([]*career.Burst{b}))
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should show dash for empty description", func() {
			b := fixtures.Burst("burst-empty", "e1", "e2")
			b.Description = ""
			view = burst.NewList(display.BurstsFromDomain([]*career.Burst{b}))
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should show confirmed status", func() {
			b := fixtures.BurstConfirmed("burst-conf", "e1")
			view = burst.NewList(display.BurstsFromDomain([]*career.Burst{b}))
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should replace newlines in description", func() {
			b := fixtures.Burst("burst-nl", "e1", "e2")
			b.Description = "Line one\nLine two\rLine three"
			view = burst.NewList(display.BurstsFromDomain([]*career.Burst{b}))
			content := view.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	var _ = Describe("Boundary: burst view import rules", func() {
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
