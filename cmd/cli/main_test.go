package main

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Initialization", func() {
	Context("Version Flag", func() {
		It("should print version when --version flag is provided", func() {
			var buf bytes.Buffer
			exitCode := run([]string{"--version"}, &buf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("KaRiya CLI v"))
		})
	})

	Context("Help Flag", func() {
		It("should print help information when --help flag is provided", func() {
			var buf bytes.Buffer
			exitCode := run([]string{"--help"}, &buf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("KaRiya CLI - Career Journaling Tool"))
			Expect(buf.String()).To(ContainSubstring("Usage: kariya [options]"))
			Expect(buf.String()).To(ContainSubstring("Commands:"))
		})
	})
})
