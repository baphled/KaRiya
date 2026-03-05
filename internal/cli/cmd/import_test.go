package cmd_test

import (
	"bytes"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
)

type dummyFileInfo struct {
	os.FileInfo
}

func (d dummyFileInfo) Name() string { return "test.csv" }
func (d dummyFileInfo) Size() int64  { return 100 }

var _ = Describe("Import Command", func() {
	var (
		ctx    *cmdpkg.CLIContext
		out    *bytes.Buffer
		errOut *bytes.Buffer
		opener cliutil.MockFileOpener
		runner cliutil.MockProgressRunner
	)

	BeforeEach(func() {
		ctx = &cmdpkg.CLIContext{InMemory: true}
		initErr := ctx.InitService()
		Expect(initErr).NotTo(HaveOccurred())

		out = new(bytes.Buffer)
		errOut = new(bytes.Buffer)
		opener = cliutil.MockFileOpener{
			StatFn: func(_ string) (os.FileInfo, error) {
				return dummyFileInfo{}, nil
			},
			OpenFn: func(_ string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader([]byte("Date,Company,Project,Text,Tags\n2023-01-01,Company,Project,Some Achievement,tag1"))), nil
			},
		}
		runner = cliutil.MockProgressRunner{}
	})

	Describe("HandleImport", func() {
		It("returns 0 on successful import", func() {
			code := cmdpkg.HandleImport("test.csv", false, false, ctx.Service, out, errOut, opener, runner)
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Import successful"))
		})

		It("returns 1 when file cannot be accessed", func() {
			opener.StatFn = func(_ string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			}
			code := cmdpkg.HandleImport("nonexistent.csv", false, false, ctx.Service, out, errOut, opener, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Cannot access import file"))
		})
	})
})
