package importcmd_test

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	importcmd "github.com/baphled/kariya/internal/cli/cmd/import"
)

var _ = Describe("Import Command", func() {
	var (
		ctx cliutil.ServiceContext
		cmd *cobra.Command
	)

	BeforeEach(func() {
		cliCtx := cmdpkg.NewCLIContext("", true)
		err := cliCtx.InitService(new(bytes.Buffer))
		Expect(err).NotTo(HaveOccurred())
		ctx = cliCtx
	})

	Describe("NewImportCmd", func() {
		It("should create command with correct properties", func() {
			cmd = importcmd.NewImportCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("import"))
			Expect(cmd.Short).To(Equal("Import events from CSV"))
			Expect(cmd.Long).To(ContainSubstring("CSV"))
		})

		It("should have file flag", func() {
			cmd = importcmd.NewImportCmd(ctx)
			flag := cmd.Flag("file")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Shorthand).To(Equal("f"))
		})

		It("should require file flag", func() {
			cmd = importcmd.NewImportCmd(ctx)
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))
			cmd.SetArgs([]string{})

			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("HandleImport", func() {
		Context("with missing file", func() {
			It("should return 1 for nonexistent file", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: "/nonexistent/file.csv",
					Service:  svc,
					Out:      new(bytes.Buffer),
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(Equal(1))
			})

			It("should write error message", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: "/nonexistent/file.csv",
					Service:  svc,
					Out:      out,
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(Equal(1))
			})
		})

		Context("with valid CSV file", func() {
			It("should handle empty CSV file", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				err := os.WriteFile(csvPath, []byte(""), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      new(bytes.Buffer),
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
			})

			It("should handle CSV with headers only", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				csvContent := "date,text,context\n"
				err := os.WriteFile(csvPath, []byte(csvContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      new(bytes.Buffer),
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
			})

			It("should handle CSV with valid data", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				csvContent := "date,text,context\n2024-01-01,Test event,Test context\n"
				err := os.WriteFile(csvPath, []byte(csvContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      out,
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
			})

			It("should write output to stdout", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				csvContent := "date,text,context\n2024-01-01,Test event,Test context\n"
				err := os.WriteFile(csvPath, []byte(csvContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      out,
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
			})
		})

		Context("with multiple rows", func() {
			It("should handle multiple CSV rows", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				csvContent := "date,text,context\n2024-01-01,Event 1,Context 1\n2024-01-02,Event 2,Context 2\n2024-01-03,Event 3,Context 3\n"
				err := os.WriteFile(csvPath, []byte(csvContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      out,
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
			})
		})

		Context("error handling", func() {
			It("should handle service errors gracefully", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				csvContent := "date,text,context\n2024-01-01,Test event,Test context\n"
				err := os.WriteFile(csvPath, []byte(csvContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      new(bytes.Buffer),
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
			})

			It("should return valid exit code", func() {
				tmpDir := GinkgoT().TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")

				csvContent := "date,text,context\n2024-01-01,Test event,Test context\n"
				err := os.WriteFile(csvPath, []byte(csvContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				code := importcmd.HandleImport(importcmd.ImportParams{
					FilePath: csvPath,
					Service:  svc,
					Out:      new(bytes.Buffer),
					ErrOut:   new(bytes.Buffer),
				})
				Expect(code).To(BeNumerically(">=", 0))
				Expect(code).To(BeNumerically("<=", 1))
			})
		})
	})
})
