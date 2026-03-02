package career

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockCloser struct {
	closed bool
	err    error
}

func (m *mockCloser) Close() error {
	m.closed = true
	return m.err
}

var _ = Describe("Repositories", func() {
	var repos *Repositories

	BeforeEach(func() {
		repos = NewRepositories()
	})

	Describe("SetCloser", func() {
		It("stores the closer for use by Close", func() {
			closer := &mockCloser{}
			repos.SetCloser(closer)

			err := repos.Close()
			Expect(err).NotTo(HaveOccurred())
			Expect(closer.closed).To(BeTrue())
		})
	})

	Describe("Close", func() {
		Context("when no closer has been set", func() {
			It("returns nil", func() {
				err := repos.Close()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when a closer has been set", func() {
			It("delegates to the closer", func() {
				closer := &mockCloser{}
				repos.SetCloser(closer)

				err := repos.Close()
				Expect(err).NotTo(HaveOccurred())
				Expect(closer.closed).To(BeTrue())
			})
		})

		Context("when the closer returns an error", func() {
			It("propagates the error", func() {
				expectedErr := errors.New("close failed")
				closer := &mockCloser{err: expectedErr}
				repos.SetCloser(closer)

				err := repos.Close()
				Expect(err).To(MatchError(expectedErr))
			})
		})
	})
})
