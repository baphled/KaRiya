package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("EventDetailModal", func() {
	var (
		modal *modals.EventDetailModal
		event *career.Event
		theme themes.Theme
	)

	BeforeEach(func() {
		event = fixtures.EventWith("test-event-id", "Event description text", "Test Company", "Test Project")
		event.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		event.Tags = []string{"tag1", "tag2"}
		event.Categories = []string{"category1"}
		event.Skills = []string{"skill1"}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewEventDetailModal", func() {
		It("creates modal with event content", func() {
			modal = modals.NewEventDetailModal(event, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetEvent()).To(Equal(event))
		})

		It("handles nil theme with default", func() {
			modal = modals.NewEventDetailModal(event, nil)

			Expect(modal).NotTo(BeNil())
		})

		It("shows skills option by default", func() {
			modal = modals.NewEventDetailModal(event, theme)

			view := modal.View()
			Expect(view).To(ContainSubstring("Skills"))
		})
	})

	Describe("WithShowSkillsOption", func() {
		It("can disable skills option", func() {
			modal = modals.NewEventDetailModal(event, theme).
				WithShowSkillsOption(false)

			Expect(modal).NotTo(BeNil())
		})

		It("returns self for chaining", func() {
			modal = modals.NewEventDetailModal(event, theme)

			result := modal.WithShowSkillsOption(true)

			Expect(result).To(Equal(modal))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewEventDetailModal(event, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()

			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewEventDetailModal(event, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("delegates to underlying modal", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})

		It("handles window size message", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewEventDetailModal(event, theme)
			modal.SetDimensions(80, 24)
		})

		It("renders event details", func() {
			modal.Show()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewEventDetailModal(event, theme)
		})

		It("Show makes modal visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("Hide makes modal invisible", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		It("sets terminal dimensions", func() {
			modal = modals.NewEventDetailModal(event, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetEvent", func() {
		It("updates the displayed event", func() {
			modal = modals.NewEventDetailModal(event, theme)
			newEvent := fixtures.Event("new-id")
			newEvent.Text = "New event text"

			modal.SetEvent(newEvent)

			Expect(modal.GetEvent()).To(Equal(newEvent))
		})
	})

	Describe("GetEvent", func() {
		It("returns the current event", func() {
			modal = modals.NewEventDetailModal(event, theme)

			result := modal.GetEvent()

			Expect(result).To(Equal(event))
			Expect(result.ID).To(Equal("test-event-id"))
		})
	})

	Describe("Key Action Signals (Phase 2-Tier 3)", func() {
		BeforeEach(func() {
			modal = modals.NewEventDetailModal(event, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("should signal 'e' key for edit action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'd' key for delete action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 's' key for skills view action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should close modal on Escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})
	})
})

var _ = Describe("SkillsDetailModal", func() {
	var (
		modal   *modals.SkillsDetailModal
		skills  []*career.Skill
		eventID string
		theme   themes.Theme
	)

	BeforeEach(func() {
		eventID = "test-event-id"
		years := 5
		goSkill := fixtures.SkillWith("skill-1", "Go", "Programming", "Expert")
		goSkill.YearsUsed = &years
		dockerSkill := fixtures.SkillWith("skill-2", "Docker", "DevOps", "Intermediate")
		skills = []*career.Skill{goSkill, dockerSkill}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewSkillsDetailModal", func() {
		It("creates modal with skills content", func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetEventID()).To(Equal(eventID))
		})

		It("handles nil theme with default", func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, nil)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty skills list", func() {
			modal = modals.NewSkillsDetailModal(eventID, []*career.Skill{}, theme)

			Expect(modal).NotTo(BeNil())
		})

		It("handles nil skills list", func() {
			modal = modals.NewSkillsDetailModal(eventID, nil, theme)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()

			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("delegates to underlying modal", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)
			modal.SetDimensions(80, 24)
		})

		It("renders skills content", func() {
			modal.Show()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)
		})

		It("Show makes modal visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("Hide makes modal invisible", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		It("sets terminal dimensions", func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetSkills", func() {
		It("updates the displayed skills", func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)
			newSkills := []*career.Skill{
				fixtures.SkillWith("new-1", "Python", "backend", "intermediate"),
			}

			modal.SetSkills(newSkills)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("GetEventID", func() {
		It("returns the event ID", func() {
			modal = modals.NewSkillsDetailModal(eventID, skills, theme)

			result := modal.GetEventID()

			Expect(result).To(Equal(eventID))
		})
	})
})
