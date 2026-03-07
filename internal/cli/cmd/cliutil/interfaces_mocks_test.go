package cliutil_test

import (
	"bytes"
	"errors"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
)

var _ = Describe("OSFileOpener", func() {
	var opener cliutil.OSFileOpener

	BeforeEach(func() {
		opener = cliutil.OSFileOpener{}
	})

	Describe("Open", func() {
		It("should open an existing file", func() {
			tmpDir := GinkgoT().TempDir()
			tmpFile := tmpDir + "/test.txt"
			err := os.WriteFile(tmpFile, []byte("test content"), 0600)
			Expect(err).NotTo(HaveOccurred())

			file, err := opener.Open(tmpFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(file).NotTo(BeNil())

			content, err := io.ReadAll(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("test content"))

			err = file.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error for non-existent file", func() {
			file, err := opener.Open("/nonexistent/path/to/file.txt")
			Expect(err).To(HaveOccurred())
			Expect(file).To(BeNil())
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})
	})

	Describe("Stat", func() {
		It("should return file info for existing file", func() {
			tmpDir := GinkgoT().TempDir()
			tmpFile := tmpDir + "/test.txt"
			err := os.WriteFile(tmpFile, []byte("test"), 0600)
			Expect(err).NotTo(HaveOccurred())

			info, err := opener.Stat(tmpFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(info).NotTo(BeNil())
			Expect(info.Name()).To(Equal("test.txt"))
			Expect(info.Size()).To(Equal(int64(4)))
		})

		It("should return error for non-existent file", func() {
			info, err := opener.Stat("/nonexistent/path/to/file.txt")
			Expect(err).To(HaveOccurred())
			Expect(info).To(BeNil())
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})
	})
})

var _ = Describe("DefaultProgressRunner", func() {
	var runner cliutil.DefaultProgressRunner

	BeforeEach(func() {
		runner = cliutil.DefaultProgressRunner{}
	})

	Describe("RunWithSpinner", func() {
		It("should execute function successfully", func() {
			executed := false
			err := runner.RunWithSpinner(
				"Processing",
				func() error {
					executed = true
					return nil
				},
				tea.WithOutput(io.Discard),
				tea.WithInput(nil),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("should return error from function", func() {
			testErr := errors.New("test error")
			err := runner.RunWithSpinner(
				"Processing",
				func() error {
					return testErr
				},
				tea.WithOutput(io.Discard),
				tea.WithInput(nil),
			)
			Expect(err).To(Equal(testErr))
		})

		It("should handle nil error return", func() {
			err := runner.RunWithSpinner(
				"Processing",
				func() error {
					return nil
				},
				tea.WithOutput(io.Discard),
				tea.WithInput(nil),
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("RunWithProgress", func() {
		It("should execute function with progress callback", func() {
			executed := false
			updateCalled := false

			err := runner.RunWithProgress(
				"Processing",
				10,
				func(update func(current int)) error {
					executed = true
					update(5)
					updateCalled = true
					return nil
				},
				tea.WithOutput(io.Discard),
				tea.WithInput(nil),
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
			Expect(updateCalled).To(BeTrue())
		})

		It("should return error from function", func() {
			testErr := errors.New("progress error")
			err := runner.RunWithProgress(
				"Processing",
				10,
				func(update func(current int)) error {
					return testErr
				},
				tea.WithOutput(io.Discard),
				tea.WithInput(nil),
			)
			Expect(err).To(Equal(testErr))
		})

		It("should handle nil error return", func() {
			err := runner.RunWithProgress(
				"Processing",
				10,
				func(update func(current int)) error {
					return nil
				},
				tea.WithOutput(io.Discard),
				tea.WithInput(nil),
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("MockFileOpener", func() {
	Describe("Open", func() {
		It("should call OpenFn when set", func() {
			called := false
			mock := cliutil.MockFileOpener{
				OpenFn: func(path string) (io.ReadCloser, error) {
					called = true
					return io.NopCloser(bytes.NewReader([]byte("mock content"))), nil
				},
			}

			file, err := mock.Open("/some/path")
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(BeTrue())
			Expect(file).NotTo(BeNil())

			content, err := io.ReadAll(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("mock content"))
		})

		It("should return ErrNotExist when OpenFn is nil", func() {
			mock := cliutil.MockFileOpener{}

			file, err := mock.Open("/some/path")
			Expect(err).To(Equal(os.ErrNotExist))
			Expect(file).To(BeNil())
		})

		It("should return error from OpenFn", func() {
			testErr := errors.New("open failed")
			mock := cliutil.MockFileOpener{
				OpenFn: func(path string) (io.ReadCloser, error) {
					return nil, testErr
				},
			}

			file, err := mock.Open("/some/path")
			Expect(err).To(Equal(testErr))
			Expect(file).To(BeNil())
		})
	})

	Describe("Stat", func() {
		It("should call StatFn when set", func() {
			called := false
			mock := cliutil.MockFileOpener{
				StatFn: func(path string) (os.FileInfo, error) {
					called = true
					return os.Stat(GinkgoT().TempDir())
				},
			}

			info, err := mock.Stat("/some/path")
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(BeTrue())
			Expect(info).NotTo(BeNil())
		})

		It("should return ErrNotExist when StatFn is nil", func() {
			mock := cliutil.MockFileOpener{}

			info, err := mock.Stat("/some/path")
			Expect(err).To(Equal(os.ErrNotExist))
			Expect(info).To(BeNil())
		})

		It("should return error from StatFn", func() {
			testErr := errors.New("stat failed")
			mock := cliutil.MockFileOpener{
				StatFn: func(path string) (os.FileInfo, error) {
					return nil, testErr
				},
			}

			info, err := mock.Stat("/some/path")
			Expect(err).To(Equal(testErr))
			Expect(info).To(BeNil())
		})
	})
})

var _ = Describe("MockProgressRunner", func() {
	Describe("RunWithSpinner", func() {
		It("should call RunWithSpinnerFn when set", func() {
			called := false
			mock := cliutil.MockProgressRunner{
				RunWithSpinnerFn: func(message string, fn func() error, opts ...tea.ProgramOption) error {
					called = true
					return fn()
				},
			}

			err := mock.RunWithSpinner("test", func() error {
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(BeTrue())
		})

		It("should execute fn directly when RunWithSpinnerFn is nil", func() {
			executed := false
			mock := cliutil.MockProgressRunner{}

			err := mock.RunWithSpinner("test", func() error {
				executed = true
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("should return error from RunWithSpinnerFn", func() {
			testErr := errors.New("spinner error")
			mock := cliutil.MockProgressRunner{
				RunWithSpinnerFn: func(message string, fn func() error, opts ...tea.ProgramOption) error {
					return testErr
				},
			}

			err := mock.RunWithSpinner("test", func() error {
				return nil
			})
			Expect(err).To(Equal(testErr))
		})

		It("should return error from fn when RunWithSpinnerFn is nil", func() {
			testErr := errors.New("fn error")
			mock := cliutil.MockProgressRunner{}

			err := mock.RunWithSpinner("test", func() error {
				return testErr
			})
			Expect(err).To(Equal(testErr))
		})
	})

	Describe("RunWithProgress", func() {
		It("should call RunWithProgressFn when set", func() {
			called := false
			mock := cliutil.MockProgressRunner{
				RunWithProgressFn: func(message string, total int, fn func(update func(current int)) error, opts ...tea.ProgramOption) error {
					called = true
					return fn(func(_ int) {})
				},
			}

			err := mock.RunWithProgress("test", 10, func(update func(current int)) error {
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(BeTrue())
		})

		It("should execute fn with no-op update when RunWithProgressFn is nil", func() {
			executed := false
			updateCalled := false
			mock := cliutil.MockProgressRunner{}

			err := mock.RunWithProgress("test", 10, func(update func(current int)) error {
				executed = true
				update(5)
				updateCalled = true
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
			Expect(updateCalled).To(BeTrue())
		})

		It("should return error from RunWithProgressFn", func() {
			testErr := errors.New("progress error")
			mock := cliutil.MockProgressRunner{
				RunWithProgressFn: func(message string, total int, fn func(update func(current int)) error, opts ...tea.ProgramOption) error {
					return testErr
				},
			}

			err := mock.RunWithProgress("test", 10, func(update func(current int)) error {
				return nil
			})
			Expect(err).To(Equal(testErr))
		})

		It("should return error from fn when RunWithProgressFn is nil", func() {
			testErr := errors.New("fn error")
			mock := cliutil.MockProgressRunner{}

			err := mock.RunWithProgress("test", 10, func(update func(current int)) error {
				return testErr
			})
			Expect(err).To(Equal(testErr))
		})
	})
})
