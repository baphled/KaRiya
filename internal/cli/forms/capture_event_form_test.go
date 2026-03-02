package forms_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("CaptureEventForm", func() {
	Describe("NewCaptureEventFormData", func() {
		It("returns a new form data with empty defaults", func() {
			data := forms.NewCaptureEventFormData()

			Expect(data).NotTo(BeNil())
			Expect(data.Text).To(BeEmpty())
			Expect(data.Date).To(BeEmpty())
			Expect(data.Company).To(BeEmpty())
			Expect(data.Project).To(BeEmpty())
			Expect(data.Tags).To(BeEmpty())
			Expect(data.Categories).To(BeEmpty())
			Expect(data.SubmitConfirmed).To(BeFalse())
		})

		It("initialises slices as empty rather than nil", func() {
			data := forms.NewCaptureEventFormData()

			Expect(data.Tags).NotTo(BeNil())
			Expect(data.Categories).NotTo(BeNil())
		})
	})

	Describe("ConfirmSubmit", func() {
		It("sets SubmitConfirmed to true", func() {
			data := forms.NewCaptureEventFormData()
			Expect(data.SubmitConfirmed).To(BeFalse())

			data.ConfirmSubmit()

			Expect(data.SubmitConfirmed).To(BeTrue())
		})
	})

	Describe("NewCaptureEventForm", func() {
		var data *forms.CaptureEventFormData

		BeforeEach(func() {
			data = forms.NewCaptureEventFormData()
		})

		Context("with quick strategy", func() {
			It("creates a non-nil form", func() {
				form := forms.NewCaptureEventForm(data, "quick", 80, 24)

				Expect(form).NotTo(BeNil())
			})
		})

		Context("with manual strategy", func() {
			It("creates a non-nil form", func() {
				form := forms.NewCaptureEventForm(data, "manual", 80, 24)

				Expect(form).NotTo(BeNil())
			})
		})

		Context("with zero dimensions", func() {
			It("creates a form without error", func() {
				form := forms.NewCaptureEventForm(data, "quick", 0, 0)

				Expect(form).NotTo(BeNil())
			})
		})
	})

	Describe("NewCaptureEventFormForModal", func() {
		var data *forms.CaptureEventFormData

		BeforeEach(func() {
			data = forms.NewCaptureEventFormData()
		})

		Context("with quick strategy", func() {
			It("creates a non-nil form", func() {
				form := forms.NewCaptureEventFormForModal(data, "quick", 80, 24)

				Expect(form).NotTo(BeNil())
			})
		})

		Context("with manual strategy", func() {
			It("creates a non-nil form with all fields", func() {
				form := forms.NewCaptureEventFormForModal(data, "manual", 80, 24)

				Expect(form).NotTo(BeNil())
			})
		})

		Context("with zero dimensions", func() {
			It("creates a form without error", func() {
				form := forms.NewCaptureEventFormForModal(data, "quick", 0, 0)

				Expect(form).NotTo(BeNil())
			})
		})
	})

	Describe("GetCaptureEventFormData", func() {
		It("extracts form data from an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.Text = "Delivered a major feature"
			event.Date = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
			event.Company = "Acme Corp"
			event.Project = "Phoenix"
			event.Tags = []string{"delivery", "feature"}
			event.Categories = []string{"engineering"}

			data := forms.GetCaptureEventFormData(event)

			Expect(data.Text).To(Equal("Delivered a major feature"))
			Expect(data.Date).To(Equal("2024-06-15"))
			Expect(data.Company).To(Equal("Acme Corp"))
			Expect(data.Project).To(Equal("Phoenix"))
			Expect(data.Tags).To(Equal([]string{"delivery", "feature"}))
			Expect(data.Categories).To(Equal([]string{"engineering"}))
		})

		It("handles nil tags by returning empty slice", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.Text = "Test event with nil tags"
			event.Tags = nil

			data := forms.GetCaptureEventFormData(event)

			Expect(data.Tags).NotTo(BeNil())
			Expect(data.Tags).To(BeEmpty())
		})

		It("handles nil categories by returning empty slice", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.Text = "Test event with nil categories"
			event.Categories = nil

			data := forms.GetCaptureEventFormData(event)

			Expect(data.Categories).NotTo(BeNil())
			Expect(data.Categories).To(BeEmpty())
		})

		It("handles both nil tags and categories", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.Text = "Test event with nil collections"
			event.Tags = nil
			event.Categories = nil

			data := forms.GetCaptureEventFormData(event)

			Expect(data.Tags).To(Equal([]string{}))
			Expect(data.Categories).To(Equal([]string{}))
		})
	})
})
