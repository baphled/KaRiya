package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
)

var _ = Describe("CLIContext", func() {
	Describe("NewCLIContext", func() {
		It("creates context with provided parameters", func() {
			ctx := cmdpkg.NewCLIContext("/path/to/db", false)
			Expect(ctx).NotTo(BeNil())
			Expect(ctx.DBPath).To(Equal("/path/to/db"))
			Expect(ctx.InMemory).To(BeFalse())
			Expect(ctx.Service).To(BeNil())
		})

		It("creates context for in-memory mode", func() {
			ctx := cmdpkg.NewCLIContext("", true)
			Expect(ctx).NotTo(BeNil())
			Expect(ctx.InMemory).To(BeTrue())
			Expect(ctx.DBPath).To(Equal(""))
		})
	})

	Describe("InitService", func() {
		It("initializes service with in-memory repositories", func() {
			ctx := cmdpkg.NewCLIContext("", true)
			errOut := &bytes.Buffer{}
			err := ctx.InitService(errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("initializes service with database repositories", func() {
			tmpDir, err := os.MkdirTemp("", "kariya-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "test.db")
			ctx := cmdpkg.NewCLIContext(dbPath, false)
			errOut := &bytes.Buffer{}
			err = ctx.InitService(errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("uses provided DBPath when specified", func() {
			tmpDir, err := os.MkdirTemp("", "kariya-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "test.db")
			ctx := cmdpkg.NewCLIContext(dbPath, false)
			errOut := &bytes.Buffer{}
			err = ctx.InitService(errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("uses default path when DBPath is empty", func() {
			ctx := cmdpkg.NewCLIContext("", false)
			errOut := &bytes.Buffer{}
			err := ctx.InitService(errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("returns error when database path is invalid", func() {
			ctx := cmdpkg.NewCLIContext("/invalid/path/that/does/not/exist/test.db", false)
			errOut := &bytes.Buffer{}
			err := ctx.InitService(errOut)
			Expect(err).To(HaveOccurred())
		})

		It("writes error to errOut when initialization fails", func() {
			ctx := cmdpkg.NewCLIContext("/invalid/path/that/does/not/exist/test.db", false)
			errOut := &bytes.Buffer{}
			err := ctx.InitService(errOut)
			Expect(err).To(HaveOccurred())
			Expect(errOut.String()).NotTo(BeEmpty())
		})

		It("configures all repositories on service", func() {
			ctx := cmdpkg.NewCLIContext("", true)
			errOut := &bytes.Buffer{}
			err := ctx.InitService(errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("creates kariya directory when using default path", func() {
			ctx := cmdpkg.NewCLIContext("", false)
			errOut := &bytes.Buffer{}
			err := ctx.InitService(errOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})
	})
})
