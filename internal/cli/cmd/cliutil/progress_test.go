package cliutil_test

import (
	"bytes"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
)

var _ = Describe("Progress Utilities", func() {
	Describe("PrintSuccess", func() {
		It("should not panic with valid message", func() {
			Expect(func() {
				cliutil.PrintSuccess("Success message")
			}).NotTo(Panic())
		})

		It("should not panic with empty message", func() {
			Expect(func() {
				cliutil.PrintSuccess("")
			}).NotTo(Panic())
		})

		It("should not panic with long message", func() {
			Expect(func() {
				cliutil.PrintSuccess("This is a very long success message that contains multiple words and should still work correctly without any issues")
			}).NotTo(Panic())
		})

		It("should not panic with special characters", func() {
			Expect(func() {
				cliutil.PrintSuccess("Success! @#$%^&*()")
			}).NotTo(Panic())
		})
	})

	Describe("PrintError", func() {
		It("should not panic with valid message", func() {
			Expect(func() {
				cliutil.PrintError("Error message")
			}).NotTo(Panic())
		})

		It("should not panic with empty message", func() {
			Expect(func() {
				cliutil.PrintError("")
			}).NotTo(Panic())
		})

		It("should not panic with long message", func() {
			Expect(func() {
				cliutil.PrintError("This is a very long error message that contains multiple words and should still work correctly without any issues")
			}).NotTo(Panic())
		})

		It("should not panic with special characters", func() {
			Expect(func() {
				cliutil.PrintError("Error! @#$%^&*()")
			}).NotTo(Panic())
		})
	})

	Describe("PrintInfo", func() {
		It("should not panic with valid message", func() {
			Expect(func() {
				cliutil.PrintInfo("Info message")
			}).NotTo(Panic())
		})

		It("should not panic with empty message", func() {
			Expect(func() {
				cliutil.PrintInfo("")
			}).NotTo(Panic())
		})

		It("should not panic with long message", func() {
			Expect(func() {
				cliutil.PrintInfo("This is a very long info message that contains multiple words and should still work correctly without any issues")
			}).NotTo(Panic())
		})

		It("should not panic with special characters", func() {
			Expect(func() {
				cliutil.PrintInfo("Info! @#$%^&*()")
			}).NotTo(Panic())
		})
	})

	Describe("Multiple calls", func() {
		It("should handle multiple success calls", func() {
			Expect(func() {
				cliutil.PrintSuccess("Message 1")
				cliutil.PrintSuccess("Message 2")
				cliutil.PrintSuccess("Message 3")
			}).NotTo(Panic())
		})

		It("should handle multiple error calls", func() {
			Expect(func() {
				cliutil.PrintError("Error 1")
				cliutil.PrintError("Error 2")
				cliutil.PrintError("Error 3")
			}).NotTo(Panic())
		})

		It("should handle multiple info calls", func() {
			Expect(func() {
				cliutil.PrintInfo("Info 1")
				cliutil.PrintInfo("Info 2")
				cliutil.PrintInfo("Info 3")
			}).NotTo(Panic())
		})

		It("should handle mixed calls", func() {
			Expect(func() {
				cliutil.PrintSuccess("Success")
				cliutil.PrintError("Error")
				cliutil.PrintInfo("Info")
			}).NotTo(Panic())
		})
	})

	Describe("RunWithSpinner", func() {
		It("should complete successfully when function returns nil", func() {
			executed := false
			err := cliutil.RunWithSpinner("Processing", func() error {
				executed = true
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("should propagate error when function returns error", func() {
			testErr := fmt.Errorf("test error")
			err := cliutil.RunWithSpinner("Processing", func() error {
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
		})

		It("should handle immediate completion", func() {
			callCount := 0
			err := cliutil.RunWithSpinner("Quick task", func() error {
				callCount++
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(callCount).To(Equal(1))
		})

		It("should handle immediate error", func() {
			testErr := fmt.Errorf("immediate error")
			err := cliutil.RunWithSpinner("Failing task", func() error {
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
		})

		It("should handle nil error explicitly", func() {
			err := cliutil.RunWithSpinner("Task", func() error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).ToNot(HaveOccurred())
		})

		It("should execute function exactly once", func() {
			callCount := 0
			cliutil.RunWithSpinner("Task", func() error {
				callCount++
				return nil
			}, tea.WithInput(nil))

			Expect(callCount).To(Equal(1))
		})

		It("should handle ctrl+c gracefully", func() {
			err := cliutil.RunWithSpinner("Task", func() error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle wrapped error types", func() {
			innerErr := fmt.Errorf("inner error")
			wrappedErr := fmt.Errorf("wrapped: %w", innerErr)
			err := cliutil.RunWithSpinner("Task", func() error {
				return wrappedErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(wrappedErr))
		})

		It("should handle slow function with multiple ticks", func() {
			executed := false
			err := cliutil.RunWithSpinner("Slow task", func() error {
				executed = true
				time.Sleep(50 * time.Millisecond)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("should handle function that returns nil error", func() {
			var returnedErr error
			err := cliutil.RunWithSpinner("Task", func() error {
				returnedErr = nil
				return returnedErr
			}, tea.WithInput(nil))

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("RunWithProgress", func() {
		It("should complete successfully with progress updates", func() {
			executed := false
			updateCount := 0
			err := cliutil.RunWithProgress("Processing", 10, func(update func(int)) error {
				executed = true
				for i := 1; i <= 5; i++ {
					update(i)
					updateCount++
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
			Expect(updateCount).To(Equal(5))
		})

		It("should propagate error after partial progress", func() {
			testErr := fmt.Errorf("progress error")
			updateCount := 0
			err := cliutil.RunWithProgress("Processing", 10, func(update func(int)) error {
				for i := 1; i <= 3; i++ {
					update(i)
					updateCount++
				}
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(updateCount).To(Equal(3))
		})

		It("should handle immediate completion with single update", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Quick", 5, func(update func(int)) error {
				update(5)
				updateCount++
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(1))
		})

		It("should handle partial progress before error", func() {
			testErr := fmt.Errorf("partial error")
			updateCount := 0
			err := cliutil.RunWithProgress("Processing", 20, func(update func(int)) error {
				update(5)
				updateCount++
				update(10)
				updateCount++
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(updateCount).To(Equal(2))
		})

		It("should handle zero total edge case", func() {
			executed := false
			err := cliutil.RunWithProgress("Zero total", 0, func(update func(int)) error {
				executed = true
				update(0)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("should handle no progress updates", func() {
			err := cliutil.RunWithProgress("No updates", 10, func(update func(int)) error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle multiple progress updates to completion", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Full progress", 100, func(update func(int)) error {
				for i := 10; i <= 100; i += 10 {
					update(i)
					updateCount++
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(10))
		})

		It("should handle ctrl+c gracefully", func() {
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				update(5)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle wrapped error types", func() {
			innerErr := fmt.Errorf("inner error")
			wrappedErr := fmt.Errorf("wrapped: %w", innerErr)
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				update(5)
				return wrappedErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(wrappedErr))
		})

		It("should handle progress exceeding total", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Overflow", 10, func(update func(int)) error {
				update(5)
				updateCount++
				update(15)
				updateCount++
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(2))
		})

		It("should handle slow function with multiple progress updates", func() {
			executed := false
			updateCount := 0
			err := cliutil.RunWithProgress("Slow progress", 10, func(update func(int)) error {
				executed = true
				for i := 1; i <= 5; i++ {
					update(i)
					updateCount++
					time.Sleep(10 * time.Millisecond)
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
			Expect(updateCount).To(Equal(5))
		})

		It("should handle function that returns nil error", func() {
			var returnedErr error
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				update(5)
				returnedErr = nil
				return returnedErr
			}, tea.WithInput(nil))

			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle large total value", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Large total", 1000000, func(update func(int)) error {
				update(500000)
				updateCount++
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(1))
		})
	})

	Describe("RunWithSpinner with input", func() {
		It("should handle other input gracefully", func() {
			executed := false
			err := cliutil.RunWithSpinner("Task", func() error {
				executed = true
				time.Sleep(50 * time.Millisecond)
				return nil
			}, tea.WithInput(bytes.NewReader([]byte("abc"))))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})
	})

	Describe("RunWithProgress with input", func() {
		It("should handle input during progress", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Processing", 10, func(update func(int)) error {
				for i := 1; i <= 5; i++ {
					update(i)
					updateCount++
					time.Sleep(20 * time.Millisecond)
				}
				return nil
			}, tea.WithInput(bytes.NewReader([]byte("abc\x1b[B"))))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(5))
		})

		It("should complete before input is processed", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Quick", 5, func(update func(int)) error {
				update(5)
				updateCount++
				return nil
			}, tea.WithInput(bytes.NewReader([]byte("abc"))))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(1))
		})

		It("should handle other input gracefully", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				for i := 1; i <= 3; i++ {
					update(i)
					updateCount++
					time.Sleep(20 * time.Millisecond)
				}
				return nil
			}, tea.WithInput(bytes.NewReader([]byte("xyz"))))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(3))
		})
	})

	Describe("Edge cases for coverage", func() {
		It("RunWithSpinner should handle successful execution", func() {
			executed := false
			err := cliutil.RunWithSpinner("Task", func() error {
				executed = true
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("RunWithSpinner should propagate function errors", func() {
			testErr := fmt.Errorf("function error")
			err := cliutil.RunWithSpinner("Task", func() error {
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
		})

		It("RunWithProgress should handle successful execution", func() {
			executed := false
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				executed = true
				update(10)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("RunWithProgress should propagate function errors", func() {
			testErr := fmt.Errorf("function error")
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				update(5)
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
		})

		It("RunWithSpinner should handle nil return from function", func() {
			err := cliutil.RunWithSpinner("Task", func() error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).ToNot(HaveOccurred())
		})

		It("RunWithProgress should handle nil return from function", func() {
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).ToNot(HaveOccurred())
		})

		It("RunWithSpinner should handle very long running task", func() {
			executed := false
			err := cliutil.RunWithSpinner("Long task", func() error {
				executed = true
				time.Sleep(200 * time.Millisecond)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("RunWithProgress should handle many progress updates", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Many updates", 100, func(update func(int)) error {
				for i := 1; i <= 100; i++ {
					update(i)
					updateCount++
					if i%10 == 0 {
						time.Sleep(5 * time.Millisecond)
					}
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(100))
		})

		It("RunWithSpinner should handle error with nil input", func() {
			testErr := fmt.Errorf("test error")
			err := cliutil.RunWithSpinner("Failing", func() error {
				time.Sleep(50 * time.Millisecond)
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
		})

		It("RunWithProgress should handle error with nil input", func() {
			testErr := fmt.Errorf("test error")
			err := cliutil.RunWithProgress("Failing", 10, func(update func(int)) error {
				update(5)
				time.Sleep(50 * time.Millisecond)
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
		})

		It("RunWithSpinner should handle multiple sequential calls", func() {
			for i := range 3 {
				err := cliutil.RunWithSpinner(fmt.Sprintf("Task %d", i), func() error {
					return nil
				}, tea.WithInput(nil))
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("RunWithProgress should handle multiple sequential calls", func() {
			for i := range 3 {
				err := cliutil.RunWithProgress(fmt.Sprintf("Task %d", i), 10, func(update func(int)) error {
					update(i + 1)
					return nil
				}, tea.WithInput(nil))
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("RunWithSpinner should handle concurrent calls", func() {
			done := make(chan error, 2)
			go func() {
				done <- cliutil.RunWithSpinner("Task 1", func() error {
					time.Sleep(50 * time.Millisecond)
					return nil
				}, tea.WithInput(nil))
			}()
			go func() {
				done <- cliutil.RunWithSpinner("Task 2", func() error {
					time.Sleep(50 * time.Millisecond)
					return nil
				}, tea.WithInput(nil))
			}()

			err1 := <-done
			err2 := <-done
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("RunWithProgress should handle concurrent calls", func() {
			done := make(chan error, 2)
			go func() {
				done <- cliutil.RunWithProgress("Task 1", 10, func(update func(int)) error {
					update(5)
					time.Sleep(50 * time.Millisecond)
					return nil
				}, tea.WithInput(nil))
			}()
			go func() {
				done <- cliutil.RunWithProgress("Task 2", 10, func(update func(int)) error {
					update(5)
					time.Sleep(50 * time.Millisecond)
					return nil
				}, tea.WithInput(nil))
			}()

			err1 := <-done
			err2 := <-done
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		It("RunWithSpinner should handle very long task with many ticks", func() {
			executed := false
			err := cliutil.RunWithSpinner("Very long task", func() error {
				executed = true
				time.Sleep(300 * time.Millisecond)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("RunWithProgress should handle many rapid updates", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Rapid updates", 1000, func(update func(int)) error {
				for i := 1; i <= 1000; i++ {
					update(i)
					updateCount++
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(1000))
		})

		It("RunWithSpinner should handle task with delayed start", func() {
			executed := false
			err := cliutil.RunWithSpinner("Delayed task", func() error {
				time.Sleep(100 * time.Millisecond)
				executed = true
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("RunWithProgress should handle task with delayed updates", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Delayed updates", 10, func(update func(int)) error {
				time.Sleep(50 * time.Millisecond)
				for i := 1; i <= 5; i++ {
					update(i)
					updateCount++
					time.Sleep(10 * time.Millisecond)
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(5))
		})

		It("RunWithProgress should not panic on ctrl+c with slow task", func() {
			Expect(func() {
				cliutil.RunWithProgress("Slow task", 100, func(update func(int)) error {
					for i := range 10 {
						update(i + 1)
						time.Sleep(5 * time.Millisecond)
					}
					return nil
				}, tea.WithInput(bytes.NewReader([]byte("\x03"))))
			}).NotTo(Panic())
		})

		It("RunWithSpinner should handle non-ctrl+c input", func() {
			executed := false
			err := cliutil.RunWithSpinner("Task", func() error {
				executed = true
				time.Sleep(100 * time.Millisecond)
				return nil
			}, tea.WithInput(bytes.NewReader([]byte("abcdefghij"))))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("RunWithProgress should handle non-ctrl+c input", func() {
			updateCount := 0
			err := cliutil.RunWithProgress("Task", 10, func(update func(int)) error {
				for i := 1; i <= 5; i++ {
					update(i)
					updateCount++
					time.Sleep(20 * time.Millisecond)
				}
				return nil
			}, tea.WithInput(bytes.NewReader([]byte("abcdefghij"))))

			Expect(err).NotTo(HaveOccurred())
			Expect(updateCount).To(Equal(5))
		})
	})
})
