package main

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Initialization", func() {
	Context("Version Flag", func() {
		It("should print version when --version flag is provided", func() {
			// Save original args and restore after test
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set test arguments
			os.Args = []string{"kariya", "--version"}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run main function
			go func() {
				defer w.Close()
				main()
			}()

			// Read captured output
			w.Close()
			output, _ := ioutil.ReadAll(r)
			os.Stdout = oldStdout

			// Assert
			Expect(string(output)).To(ContainSubstring("KaRiya CLI v"))
		})
	})

	Context("Help Flag", func() {
		It("should print help information when --help flag is provided", func() {
			// Save original args and restore after test
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set test arguments
			os.Args = []string{"kariya", "--help"}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run main function
			go func() {
				defer w.Close()
				main()
			}()

			// Read captured output
			w.Close()
			output, _ := ioutil.ReadAll(r)
			os.Stdout = oldStdout

			// Assert
			Expect(string(output)).To(ContainSubstring("KaRiya CLI - Career Journaling Tool"))
			Expect(string(output)).To(ContainSubstring("Usage: kariya [options]"))
			Expect(string(output)).To(ContainSubstring("Commands:"))
		})
	})
})
