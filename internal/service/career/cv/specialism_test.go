package cv

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeriveSpecialism", func() {
	DescribeTable("should derive correct specialism",
		func(technology, position, sector, expected string) {
			actual := DeriveSpecialism(technology, position, sector)
			Expect(actual).To(Equal(expected))
		},
		Entry("Go backend startup", "Go", "backend", "startup", "Platform Engineering"),
		Entry("Ruby backend enterprise", "Ruby", "backend", "enterprise", "Enterprise Engineering"),
		Entry("JavaScript frontend public-sector", "JavaScript", "frontend", "public-sector", "Public Sector Engineering"),
		Entry("DevOps ai/ml startup", "DevOps", "ai/ml", "startup", "AI/ML Engineering"),
		Entry("Fullstack backend startup", "Fullstack", "backend", "startup", "Platform Engineering"),
		Entry("Go frontend enterprise", "Go", "frontend", "enterprise", "Enterprise Engineering"),
		Entry("Ruby devops public-sector", "Ruby", "devops", "public-sector", "Public Sector Engineering"),
		Entry("Unknown mapping", "Cobol", "mainframe", "bank", "General Engineering"),
	)
})
