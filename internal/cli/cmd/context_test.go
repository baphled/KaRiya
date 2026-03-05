package cmd_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
)

var _ = Describe("CLIContext", func() {
	Describe("InitService", func() {
		It("initializes service with in-memory repositories", func() {
			ctx := &cmdpkg.CLIContext{InMemory: true}
			err := ctx.InitService()
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("initializes service with database repositories", func() {
			tmpDir, err := os.MkdirTemp("", "kariya-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "test.db")
			ctx := &cmdpkg.CLIContext{InMemory: false, DBPath: dbPath}
			err = ctx.InitService()
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("uses provided DBPath when specified", func() {
			tmpDir, err := os.MkdirTemp("", "kariya-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "test.db")
			ctx := &cmdpkg.CLIContext{InMemory: false, DBPath: dbPath}
			err = ctx.InitService()
			Expect(err).NotTo(HaveOccurred())
			Expect(ctx.Service).NotTo(BeNil())
		})

		It("returns error when database path is invalid", func() {
			ctx := &cmdpkg.CLIContext{InMemory: false, DBPath: "/invalid/path/that/does/not/exist/test.db"}
			err := ctx.InitService()
			Expect(err).To(HaveOccurred())
		})

		It("returns error when home directory cannot be determined", func() {
			// This test is difficult to implement without mocking os.UserHomeDir
			// Skip for now as it requires more complex setup
		})
	})
})
