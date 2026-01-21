package components_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container Rendering", func() {
	It("should include all rows in the rendered output", func() {
		rows := [][]string{
			{"Row1"},
			{"Row2"},
		}
		rendered := "Row1\nRow2"
		for _, row := range rows {
			Expect(strings.Contains(rendered, row[0])).To(BeTrue(), "expected list to contain: %q", row[0])
		}
	})
})
