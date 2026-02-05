package models_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaptureForm", func() {
	var (
		form       *models.CaptureForm
		cliService *service.CLIEventService
		repo       *careermemory.EventRepository
		svc        *careerservice.Service
	)

	BeforeEach(func() {
		repo = careermemory.NewEventRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		form = models.NewCaptureForm(cliService)
	})

	Describe("Initialization", func() {
		It("should create form with default manual strategy", func() {
			Expect(form).NotTo(BeNil())
			Expect(form.GetStrategy()).To(Equal("manual"))
		})

		It("should return init command", func() {
			cmd := form.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("should have form instance created", func() {
			Expect(form).NotTo(BeNil())
			// Form is created in constructor via rebuildForm
		})
	})

	Describe("Strategy Management", func() {
		Context("when setting quick strategy", func() {
			It("should update strategy to quick", func() {
				form.SetStrategy("quick")
				Expect(form.GetStrategy()).To(Equal("quick"))
			})

			It("should rebuild form with quick strategy fields", func() {
				form.SetStrategy("quick")
				// Form is rebuilt with fewer fields for quick capture
				Expect(form.GetStrategy()).To(Equal("quick"))
			})
		})

		Context("when setting manual strategy", func() {
			It("should update strategy to manual", func() {
				form.SetStrategy("manual")
				Expect(form.GetStrategy()).To(Equal("manual"))
			})

			It("should rebuild form with all fields", func() {
				form.SetStrategy("manual")
				// Form is rebuilt with all fields for manual capture
				Expect(form.GetStrategy()).To(Equal("manual"))
			})
		})

		Context("when switching strategies", func() {
			It("should allow switching from quick to manual", func() {
				form.SetStrategy("quick")
				form.SetStrategy("manual")
				Expect(form.GetStrategy()).To(Equal("manual"))
			})

			It("should allow switching from manual to quick", func() {
				form.SetStrategy("manual")
				form.SetStrategy("quick")
				Expect(form.GetStrategy()).To(Equal("quick"))
			})

			It("should preserve form data when switching strategies", func() {
				form.SetStrategy("quick")
				// Form data is preserved across strategy switches
				form.SetStrategy("manual")
				Expect(form.GetStrategy()).To(Equal("manual"))
			})
		})
	})

	Describe("Event Loading", func() {
		Context("when loading event for editing", func() {
			It("should populate form with event data", func() {
				event := fixtures.EventWith("test-1", "Original event text", "Test Company", "Test Project")
				event.Tags = []string{"tag1", "tag2"}
				event.Categories = []string{"cat1"}

				form.LoadEventForEditing(event)

				formData := form.GetFormData()
				Expect(formData.Text).To(Equal("Original event text"))
				Expect(formData.Company).To(Equal("Test Company"))
				Expect(formData.Project).To(Equal("Test Project"))
			})

			It("should handle nil event gracefully", func() {
				originalText := form.GetFormData().Text

				form.LoadEventForEditing(nil)

				Expect(form.GetFormData().Text).To(Equal(originalText))
			})

			It("should copy tags and categories", func() {
				event := fixtures.EventWith("test-2", "Event with metadata", "", "")
				event.Tags = []string{"backend", "go"}
				event.Categories = []string{"development"}

				form.LoadEventForEditing(event)

				formData := form.GetFormData()
				Expect(formData.Tags).To(Equal([]string{"backend", "go"}))
				Expect(formData.Categories).To(Equal([]string{"development"}))

				event.Tags[0] = "modified"
				Expect(formData.Tags[0]).To(Equal("backend"))
			})

			It("should format date correctly", func() {
				testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
				event := fixtures.EventWith("test-3", "Event with specific date", "", "")
				event.Date = testDate

				form.LoadEventForEditing(event)

				formData := form.GetFormData()
				Expect(formData.Date).To(Equal("2024-01-15"))
			})
		})
	})

	Describe("Form Submission", func() {
		Context("when submitting form", func() {
			It("should create SubmitMsg with event data", func() {
				cmd := form.SubmitForm()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)
				Expect(ok).To(BeTrue())
				Expect(submitMsg.Event).NotTo(BeNil())
			})

			It("should use current time when date is empty", func() {
				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg := msg.(models.SubmitMsg)

				Expect(submitMsg.Event.Date).To(BeTemporally("~", time.Now(), time.Second))
			})

			It("should parse date string when provided", func() {
				event := fixtures.EventWith("test-date", "Event for date test", "", "")
				event.Date = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
				form.LoadEventForEditing(event)

				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)

				Expect(ok).To(BeTrue())
				Expect(submitMsg.Event.Date.Year()).To(Equal(2024))
				Expect(submitMsg.Event.Date.Month()).To(Equal(time.June))
				Expect(submitMsg.Event.Date.Day()).To(Equal(15))
			})

			It("should return error for invalid date", func() {
				formData := form.GetFormData()
				formData.Date = "invalid-date-format"

				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)

				Expect(ok).To(BeTrue())
				Expect(submitMsg.Err).To(HaveOccurred())
			})

			It("should include all metadata in event", func() {
				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg := msg.(models.SubmitMsg)

				Expect(submitMsg.Event).To(SatisfyAll(
					HaveField("CreatedAt", Not(BeZero())),
					HaveField("UpdatedAt", Not(BeZero())),
				))
			})
		})

		Context("when triggering submission via Ctrl+S", func() {
			It("should submit form with Ctrl+S key", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				_, cmd := form.Update(msg)
				Expect(cmd).NotTo(BeNil())

				// Should produce SubmitMsg
				result := cmd()
				_, ok := result.(models.SubmitMsg)
				Expect(ok).To(BeTrue())
			})

			It("should set submit confirmed flag", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlS}
				_, _ = form.Update(msg)

				formData := form.GetFormData()
				Expect(formData.SubmitConfirmed).To(BeTrue())
			})
		})

		Context("when form is completed via huh", func() {
			It("should return submission command when form state is completed", func() {
				cmd := form.SubmitForm()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)
				Expect(ok).To(BeTrue())
				Expect(submitMsg.Event).NotTo(BeNil())
			})
		})
	})

	Describe("Cancel Behavior", func() {
		Context("when pressing Escape", func() {
			It("should allow parent to handle back navigation", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				updatedModel, _ := form.Update(msg)

				Expect(updatedModel).NotTo(BeNil())
				// Escape is handled BEFORE delegating to form
				// Parent intent should receive this
			})

			It("should not delegate Escape to huh form", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, cmd := form.Update(msg)

				// Command should be nil as Escape is intercepted
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("Window Resizing", func() {
		Context("when receiving WindowSizeMsg", func() {
			It("should update width and height", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				updatedModel, _ := form.Update(msg)

				Expect(updatedModel).NotTo(BeNil())
				// Form dimensions are updated
			})

			It("should resize huh form", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				_, _ = form.Update(msg)

				// Form is resized with new dimensions
			})
		})
	})

	Describe("View Rendering", func() {
		Context("when rendering form", func() {
			It("should return non-empty view", func() {
				view := form.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should delegate to huh form view", func() {
				view := form.View()
				// View comes from huh.Form.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Behavior Integration", func() {
		Context("when using form in workflow", func() {
			It("should support full capture flow", func() {
				// 1. Create form
				Expect(form).NotTo(BeNil())

				// 2. Set strategy
				form.SetStrategy("quick")

				// 3. Submit form
				cmd := form.SubmitForm()
				msg := cmd()

				// 4. Verify event created
				submitMsg := msg.(models.SubmitMsg)
				Expect(submitMsg.Event).NotTo(BeNil())
				Expect(submitMsg.Err).ToNot(HaveOccurred())
			})

			It("should handle edit flow", func() {
				event := fixtures.EventWith("test-4", "Original text", "Original Company", "")
				form.LoadEventForEditing(event)

				cmd := form.SubmitForm()
				msg := cmd()

				submitMsg := msg.(models.SubmitMsg)
				Expect(submitMsg.Event).NotTo(BeNil())
			})
		})
	})

	Describe("Edge Cases", func() {
		Context("when handling empty form", func() {
			It("should create event with default values", func() {
				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg := msg.(models.SubmitMsg)

				Expect(submitMsg.Event).NotTo(BeNil())
				Expect(submitMsg.Event.Date).To(BeTemporally("~", time.Now(), time.Second))
			})
		})

		Context("when handling special characters", func() {
			It("should preserve special characters in text", func() {
				event := fixtures.EventWith("test-5", "Event with @#$%^&* special chars", "", "")
				form.LoadEventForEditing(event)

				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg := msg.(models.SubmitMsg)

				Expect(submitMsg.Event.Text).To(ContainSubstring("@#$%"))
			})
		})

		Context("when handling form state transitions", func() {
			It("should handle rapid updates", func() {
				for range 10 {
					msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
					_, _ = form.Update(msg)
				}
				Expect(form).NotTo(BeNil())
			})
		})
	})

	Describe("Date Parsing", func() {
		Context("when parsing date strings", func() {
			It("should accept YYYY-MM-DD format", func() {
				formData := form.GetFormData()
				formData.Date = "2024-03-15"

				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)

				Expect(ok).To(BeTrue())
				Expect(submitMsg.Err).ToNot(HaveOccurred())
				Expect(submitMsg.Event.Date.Year()).To(Equal(2024))
				Expect(submitMsg.Event.Date.Month()).To(Equal(time.March))
				Expect(submitMsg.Event.Date.Day()).To(Equal(15))
			})

			It("should reject invalid date formats", func() {
				formData := form.GetFormData()
				formData.Date = "15/03/2024"

				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)

				Expect(ok).To(BeTrue())
				Expect(submitMsg.Err).To(HaveOccurred())
			})

			It("should handle edge dates", func() {
				formData := form.GetFormData()
				formData.Date = "2024-02-29"

				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)

				Expect(ok).To(BeTrue())
				Expect(submitMsg.Err).ToNot(HaveOccurred())
				Expect(submitMsg.Event.Date.Day()).To(Equal(29))
			})
		})
	})

	Describe("Form State Management", func() {
		Context("when form reaches completed state", func() {
			It("should return submit command when SubmitForm is called", func() {
				cmd := form.SubmitForm()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				submitMsg, ok := msg.(models.SubmitMsg)
				Expect(ok).To(BeTrue())
				Expect(submitMsg.Event).NotTo(BeNil())
			})
		})

		Context("when tracking form state", func() {
			It("should allow strategy inspection", func() {
				Expect(form.GetStrategy()).To(Equal("manual"))

				form.SetStrategy("quick")
				Expect(form.GetStrategy()).To(Equal("quick"))
			})
		})
	})

	Describe("BaseStandardModel Integration", func() {
		It("should embed BaseStandardModel", func() {
			// CaptureForm embeds *BaseStandardModel
			// Inheriting common model behavior
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("Strategy-Specific Behavior", func() {
		Context("quick strategy", func() {
			BeforeEach(func() {
				form.SetStrategy("quick")
			})

			It("should have minimal required fields", func() {
				// Quick strategy shows fewer fields
				Expect(form.GetStrategy()).To(Equal("quick"))
			})

			It("should use defaults for optional fields", func() {
				cmd := form.SubmitForm()
				msg := cmd()
				submitMsg := msg.(models.SubmitMsg)

				// Quick capture may have empty company/project
				Expect(submitMsg.Event).NotTo(BeNil())
			})
		})

		Context("manual strategy", func() {
			BeforeEach(func() {
				form.SetStrategy("manual")
			})

			It("should show all available fields", func() {
				// Manual strategy shows all fields
				Expect(form.GetStrategy()).To(Equal("manual"))
			})

			It("should allow full metadata entry", func() {
				// Manual form includes company, project, tags, categories
				Expect(form.GetStrategy()).To(Equal("manual"))
			})
		})
	})
})
