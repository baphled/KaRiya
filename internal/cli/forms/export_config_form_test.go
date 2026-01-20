package forms_test

import (
	"github.com/baphled/kariya/internal/cli/forms"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportConfigForm", func() {
	Describe("NewExportConfigFormData", func() {
		It("creates form data with defaults", func() {
			data := forms.NewExportConfigFormData()
			Expect(data).NotTo(BeNil())
			Expect(data.ArtifactType).To(Equal("events"))
			Expect(data.Format).To(Equal("json"))
			Expect(data.Destination).To(Equal("file"))
		})
	})

	Describe("NewExportConfigForm", func() {
		var data *forms.ExportConfigFormData

		BeforeEach(func() {
			data = forms.NewExportConfigFormData()
		})

		It("creates a form", func() {
			form := forms.NewExportConfigForm(data, 80, 0)
			Expect(form).NotTo(BeNil())
		})

		It("creates a form with dimensions", func() {
			form := forms.NewExportConfigForm(data, 100, 30)
			Expect(form).NotTo(BeNil())
		})

		It("renders the form", func() {
			form := forms.NewExportConfigForm(data, 80, 0)
			view := form.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("shows step 1 group title", func() {
			form := forms.NewExportConfigForm(data, 80, 0)
			view := form.View()
			Expect(view).To(ContainSubstring("Step 1: What to Export"))
		})

		It("shows step 1 title", func() {
			form := forms.NewExportConfigForm(data, 80, 0)
			view := form.View()
			Expect(view).To(ContainSubstring("What to Export"))
		})
	})

	Describe("GetArtifactTypeLabel", func() {
		It("returns correct label for events", func() {
			Expect(forms.GetArtifactTypeLabel("events")).To(Equal("Career Events"))
		})

		It("returns correct label for facts", func() {
			Expect(forms.GetArtifactTypeLabel("facts")).To(Equal("Facts"))
		})

		It("returns correct label for bursts", func() {
			Expect(forms.GetArtifactTypeLabel("bursts")).To(Equal("Bursts"))
		})

		It("returns correct label for cv", func() {
			Expect(forms.GetArtifactTypeLabel("cv")).To(Equal("CV/Resume"))
		})

		It("returns correct label for profile", func() {
			Expect(forms.GetArtifactTypeLabel("profile")).To(Equal("Profile"))
		})

		It("returns raw value for unknown types", func() {
			Expect(forms.GetArtifactTypeLabel("unknown")).To(Equal("unknown"))
		})
	})

	Describe("GetFormatLabel", func() {
		It("returns correct label for json", func() {
			Expect(forms.GetFormatLabel("json")).To(Equal("JSON"))
		})

		It("returns correct label for yaml", func() {
			Expect(forms.GetFormatLabel("yaml")).To(Equal("YAML"))
		})

		It("returns correct label for csv", func() {
			Expect(forms.GetFormatLabel("csv")).To(Equal("CSV"))
		})

		It("returns correct label for markdown", func() {
			Expect(forms.GetFormatLabel("markdown")).To(Equal("Markdown"))
		})

		It("returns correct label for txt", func() {
			Expect(forms.GetFormatLabel("txt")).To(Equal("Plain Text"))
		})

		It("returns raw value for unknown formats", func() {
			Expect(forms.GetFormatLabel("unknown")).To(Equal("unknown"))
		})
	})

	Describe("GetDestinationLabel", func() {
		It("returns correct label for file", func() {
			Expect(forms.GetDestinationLabel("file")).To(Equal("File"))
		})

		It("returns correct label for clipboard", func() {
			Expect(forms.GetDestinationLabel("clipboard")).To(Equal("Clipboard"))
		})

		It("returns raw value for unknown destinations", func() {
			Expect(forms.GetDestinationLabel("unknown")).To(Equal("unknown"))
		})
	})
})
