package importcmd_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		Context("with nil reader", func() {
			It("should return 1", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  nil,
					Service: svc,
					Out:     new(bytes.Buffer),
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(1))
			})
		})

		Context("with empty CSV", func() {
			It("should return 1 for completely empty reader", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				reader := bytes.NewBufferString("")
				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     new(bytes.Buffer),
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(1))
			})

			It("should return 1 for CSV with only headers", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				csvData := "Date,Text,Company\n"
				reader := bytes.NewBufferString(csvData)
				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     new(bytes.Buffer),
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(1))
			})
		})

		Context("with valid CSV data", func() {
			It("should process single valid row", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				csvData := `Date,Text,Company,Skills
2024-01-15,Implemented authentication system,TechCorp,Go;Security
`
				reader := bytes.NewBufferString(csvData)
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     out,
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("Import Complete"))
			})

			It("should process multiple valid rows", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				csvData := `Date,Text,Company,Skills
2024-01-15,Implemented authentication,TechCorp,Go;Security
2024-01-16,Built REST API,TechCorp,Go;API
2024-01-17,Deployed to production,TechCorp,DevOps;Kubernetes
`
				reader := bytes.NewBufferString(csvData)
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     out,
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("Import Complete"))
				Expect(out.String()).To(ContainSubstring("Successfully imported:"))
			})
		})

		Context("with malformed CSV", func() {
			It("should handle invalid CSV structure", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				csvData := "Invalid,CSV,Data\nNo proper structure"
				reader := bytes.NewBufferString(csvData)

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     new(bytes.Buffer),
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(1))
			})
		})

		Context("output formatting", func() {
			It("should display import statistics", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				csvData := `Date,Text,Company,Skills
2024-01-15,Event one,TestCo,Go
2024-01-16,Event two,TestCo,Python
`
				reader := bytes.NewBufferString(csvData)
				out := new(bytes.Buffer)

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     out,
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(0))

				output := out.String()
				Expect(output).To(ContainSubstring("Total rows processed:"))
				Expect(output).To(ContainSubstring("Successfully imported:"))
			})
		})

		Context("edge cases", func() {
			It("should handle CSV with special characters", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				csvData := `Date,Text,Company,Skills
2024-01-15,"Event with, comma",TestCo,Go
`
				reader := bytes.NewBufferString(csvData)

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     new(bytes.Buffer),
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle CSV with long text", func() {
				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())

				longText := strings.Repeat("A very long event description ", 50)
				csvData := fmt.Sprintf(`Date,Text,Company,Skills
2024-01-15,"%s",TestCo,Go
`, longText)
				reader := bytes.NewBufferString(csvData)

				code := importcmd.HandleImport(importcmd.ImportParams{
					Reader:  reader,
					Service: svc,
					Out:     new(bytes.Buffer),
					ErrOut:  new(bytes.Buffer),
				}, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

		})
	})

	Describe("ValidateFilePath", func() {
		It("should return nil for existing file", func() {
			tmpFile, err := os.CreateTemp("", "test-*.csv")
			Expect(err).NotTo(HaveOccurred())
			tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			err = importcmd.ValidateFilePath(tmpFile.Name())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error for non-existent file", func() {
			err := importcmd.ValidateFilePath("/tmp/nonexistent_file_12345.csv")
			Expect(err).To(HaveOccurred())
		})

		It("should return nil for existing directory", func() {
			err := importcmd.ValidateFilePath("/tmp")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

