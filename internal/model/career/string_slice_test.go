package career

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringSlice", func() {
	Describe("Scan", func() {
		var s StringSlice

		BeforeEach(func() {
			s = nil
		})

		Context("when value is nil", func() {
			It("sets the slice to nil", func() {
				err := s.Scan(nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		Context("when value is a string", func() {
			It("splits a comma-separated string", func() {
				err := s.Scan("a,b,c")

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"a", "b", "c"}))
			})

			It("handles a single element", func() {
				err := s.Scan("solo")

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"solo"}))
			})

			It("sets to nil for an empty string", func() {
				err := s.Scan("")

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		Context("when value is []byte", func() {
			It("splits a comma-separated byte slice", func() {
				err := s.Scan([]byte("x,y,z"))

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"x", "y", "z"}))
			})

			It("handles a single element", func() {
				err := s.Scan([]byte("only"))

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"only"}))
			})

			It("sets to nil for empty bytes", func() {
				err := s.Scan([]byte(""))

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		Context("when value is an unsupported type", func() {
			It("returns an error for int", func() {
				err := s.Scan(42)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported type"))
			})

			It("returns an error for bool", func() {
				err := s.Scan(true)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported type"))
			})
		})
	})

	Describe("Value", func() {
		Context("when slice is nil", func() {
			It("returns an empty string", func() {
				var s StringSlice
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(""))
			})
		})

		Context("when slice is empty", func() {
			It("returns an empty string", func() {
				s := StringSlice{}
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(""))
			})
		})

		Context("when slice has one element", func() {
			It("returns the element without commas", func() {
				s := StringSlice{"single"}
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("single"))
			})
		})

		Context("when slice has multiple elements", func() {
			It("returns a comma-separated string", func() {
				s := StringSlice{"a", "b", "c"}
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("a,b,c"))
			})
		})
	})

	Describe("Scan-Value round trip", func() {
		It("preserves data through serialisation and deserialisation", func() {
			original := StringSlice{"technical", "leadership", "mentoring"}
			val, err := original.Value()
			Expect(err).NotTo(HaveOccurred())

			var restored StringSlice
			err = restored.Scan(val)
			Expect(err).NotTo(HaveOccurred())
			Expect(restored).To(Equal(original))
		})
	})
})
