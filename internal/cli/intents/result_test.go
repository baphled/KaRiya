package intents_test

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/intents"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntentResult", func() {
	Describe("NewCompletedResult", func() {
		It("should create completed result with data", func() {
			result := intents.NewCompletedResult("test data")

			Expect(result.Status).To(Equal(intents.Completed))
			Expect(result.Data).To(Equal("test data"))
			Expect(result.Error).ToNot(HaveOccurred())
			Expect(result.IsSuccessful()).To(BeTrue())
		})

		It("should create completed result with empty string", func() {
			result := intents.NewCompletedResult("")

			Expect(result.Status).To(Equal(intents.Completed))
			Expect(result.Data).To(Equal(""))
			Expect(result.Error).ToNot(HaveOccurred())
			Expect(result.IsSuccessful()).To(BeTrue())
		})
	})

	Describe("NewCancelledResult", func() {
		It("should create cancelled result", func() {
			result := intents.NewCancelledResult[string]()

			Expect(result.Status).To(Equal(intents.Cancelled))
			Expect(result.IsCancelled()).To(BeTrue())
			Expect(result.IsSuccessful()).To(BeFalse())
		})
	})

	Describe("NewFailedResult", func() {
		It("should create failed result with error", func() {
			cause := errors.New("test error")
			result := intents.NewFailedResult[string]("test_code", "test message", cause)

			Expect(result.Status).To(Equal(intents.Failed))
			Expect(result.IsFailed()).To(BeTrue())
			Expect(result.Error).To(HaveOccurred())
			Expect(result.Error.Code).To(Equal("test_code"))
			Expect(result.Error.Cause).To(Equal(cause))
			Expect(result.IsSuccessful()).To(BeFalse())
		})
	})

	Describe("NewPartialResult", func() {
		It("should create partial result with data and error", func() {
			data := "partial data"
			result := intents.NewPartialResult(data, "partial_code", "partial message")

			Expect(result.Status).To(Equal(intents.Partial))
			Expect(result.Data).To(Equal(data))
			Expect(result.Error).To(HaveOccurred())
			Expect(result.IsSuccessful()).To(BeTrue())
		})
	})

	Describe("Metadata", func() {
		Describe("WithMetadata", func() {
			It("should set and get metadata", func() {
				result := intents.NewCompletedResult("test")

				result.WithMetadata("key1", "value1")
				result.WithMetadata("key2", 42)

				val, ok := result.GetMetadata("key1")
				Expect(ok).To(BeTrue())
				Expect(val).To(Equal("value1"))

				val, ok = result.GetMetadata("key2")
				Expect(ok).To(BeTrue())
				Expect(val).To(Equal(42))
			})

			It("should return false for nonexistent key", func() {
				result := intents.NewCompletedResult("test")
				_, ok := result.GetMetadata("nonexistent")
				Expect(ok).To(BeFalse())
			})
		})

		Describe("GetAllMetadata", func() {
			It("should return all metadata", func() {
				result := intents.NewCompletedResult("test")
				result.WithMetadata("key1", "value1")
				result.WithMetadata("key2", 42)

				all := result.GetAllMetadata()

				Expect(all).To(HaveLen(2))
				Expect(all["key1"]).To(Equal("value1"))
				Expect(all["key2"]).To(Equal(42))
			})
		})
	})

	Describe("IsTerminal", func() {
		It("should return true for completed result", func() {
			result := intents.NewCompletedResult("data")
			Expect(result.IsTerminal()).To(BeTrue())
		})

		It("should return true for cancelled result", func() {
			result := intents.NewCancelledResult[string]()
			Expect(result.IsTerminal()).To(BeTrue())
		})

		It("should return true for failed result", func() {
			result := intents.NewFailedResult[string]("code", "message", nil)
			Expect(result.IsTerminal()).To(BeTrue())
		})

		It("should return true for partial result", func() {
			result := intents.NewPartialResult("data", "code", "message")
			Expect(result.IsTerminal()).To(BeTrue())
		})
	})

	Describe("Method Chaining", func() {
		It("should support method chaining for metadata", func() {
			result := intents.NewCompletedResult("data").
				WithMetadata("key1", "value1").
				WithMetadata("key2", 42)

			val, ok := result.GetMetadata("key1")
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("value1"))
		})
	})

	Describe("IsValid", func() {
		It("should return nil for completed result", func() {
			result := intents.NewCompletedResult("data")
			Expect(result.IsValid()).To(Succeed())
		})

		It("should return nil for cancelled result", func() {
			result := intents.NewCancelledResult[string]()
			Expect(result.IsValid()).To(Succeed())
		})

		It("should return error for failed result without error", func() {
			result := &intents.IntentResult[string]{Status: intents.Failed}
			Expect(result.IsValid()).To(HaveOccurred())
		})

		It("should return nil for failed result with error", func() {
			result := intents.NewFailedResult[string]("code", "message", nil)
			Expect(result.IsValid()).To(Succeed())
		})

		It("should return error for partial result without error", func() {
			result := &intents.IntentResult[string]{Status: intents.Partial, Data: "data"}
			Expect(result.IsValid()).To(HaveOccurred())
		})

		It("should return nil for partial result with error", func() {
			result := intents.NewPartialResult("data", "code", "message")
			Expect(result.IsValid()).To(Succeed())
		})
	})
})

var _ = Describe("IntentError", func() {
	Describe("WithCause", func() {
		It("should set the cause", func() {
			err := &intents.IntentError{
				Code:    "code1",
				Message: "message1",
			}

			cause := errors.New("cause error")
			//nolint:errcheck // Test - intentionally ignoring return.
			err.WithCause(cause)

			Expect(err.Cause).To(Equal(cause))
		})
	})

	Describe("WithMessage", func() {
		It("should update the message", func() {
			err := &intents.IntentError{
				Code:    "code1",
				Message: "original message",
			}

			//nolint:errcheck // Test - intentionally ignoring return.
			err.WithMessage("updated message")

			Expect(err.Message).To(Equal("updated message"))
		})
	})

	Describe("Error", func() {
		It("should format error with cause", func() {
			err := &intents.IntentError{
				Code:    "code1",
				Message: "message1",
				Cause:   errors.New("cause error"),
			}

			Expect(err.Error()).To(Equal("code1: message1 (cause: cause error)"))
		})

		It("should format error without cause", func() {
			err := &intents.IntentError{
				Code:    "code2",
				Message: "message2",
				Cause:   nil,
			}

			Expect(err.Error()).To(Equal("code2: message2"))
		})
	})
})
