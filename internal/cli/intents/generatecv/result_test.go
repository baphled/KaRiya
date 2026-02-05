package generatecv_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Result", func() {
	It("should store generated CV", func() {
		cv := fixtures.CVView("cv-1")
		result := &generatecv.Result{GeneratedCV: cv}

		Expect(result.GeneratedCV).To(Equal(cv))
	})

	It("should store selected profile", func() {
		profile := &generatecv.CVProfile{ID: "p1", Name: "Staff Engineer"}
		result := &generatecv.Result{SelectedProfile: profile}

		Expect(result.SelectedProfile.ID).To(Equal("p1"))
		Expect(result.SelectedProfile.Name).To(Equal("Staff Engineer"))
	})

	It("should store accepted fields", func() {
		result := &generatecv.Result{
			AcceptedFields: map[string]bool{"summary": true, "skills": false},
		}

		Expect(result.AcceptedFields).To(HaveKeyWithValue("summary", true))
		Expect(result.AcceptedFields).To(HaveKeyWithValue("skills", false))
	})

	It("should store export metadata", func() {
		now := time.Now()
		result := &generatecv.Result{
			ExportPath:     "/tmp/cv.md",
			CVExportFormat: "markdown",
			ExportedAt:     &now,
		}

		Expect(result.ExportPath).To(Equal("/tmp/cv.md"))
		Expect(result.CVExportFormat).To(Equal("markdown"))
		Expect(result.ExportedAt).To(Equal(&now))
	})

	It("should default to zero values", func() {
		result := &generatecv.Result{}

		Expect(result.GeneratedCV).To(BeNil())
		Expect(result.SelectedProfile).To(BeNil())
		Expect(result.AcceptedFields).To(BeNil())
		Expect(result.ExportPath).To(BeEmpty())
		Expect(result.CVExportFormat).To(BeEmpty())
		Expect(result.ExportedAt).To(BeNil())
	})
})
