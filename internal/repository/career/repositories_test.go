package career

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "modernc.org/sqlite"
)

type mockCloser struct {
	closed    bool
	shouldErr bool
}

func (m *mockCloser) Close() error {
	m.closed = true
	if m.shouldErr {
		return errors.New("close error")
	}
	return nil
}

func newEmptyRepositories() *Repositories {
	return new(Repositories)
}

var _ = Describe("Repositories", func() {
	var repos *Repositories

	BeforeEach(func() {
		repos = newEmptyRepositories()
	})

	Describe("SetCloser", func() {
		It("sets the closer on the repositories struct", func() {
			closer := &mockCloser{}

			repos.SetCloser(closer)

			err := repos.Close()
			Expect(err).NotTo(HaveOccurred())
			Expect(closer.closed).To(BeTrue())
		})
	})

	Describe("Close", func() {
		Context("when closer is nil", func() {
			It("returns nil without error", func() {
				err := repos.Close()

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when closer is set", func() {
			It("delegates to the closer", func() {
				closer := &mockCloser{}
				repos.SetCloser(closer)

				err := repos.Close()

				Expect(err).NotTo(HaveOccurred())
				Expect(closer.closed).To(BeTrue())
			})

			It("returns error from closer", func() {
				closer := &mockCloser{shouldErr: true}
				repos.SetCloser(closer)

				err := repos.Close()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("close error"))
			})
		})
	})
})

var _ = Describe("RunMigrationsForTests", func() {
	var (
		db      *sql.DB
		tempDir string
		dbPath  string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-test-migrations-")
		Expect(err).NotTo(HaveOccurred())
		dbPath = filepath.Join(tempDir, "test.db")

		db, err = sql.Open("sqlite", dbPath)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if db != nil {
			db.Close()
		}
		os.RemoveAll(tempDir)
	})

	It("applies all migrations to a fresh database", func() {
		err := RunMigrationsForTests(db)
		Expect(err).NotTo(HaveOccurred())

		Expect(hasTable(db, "career_events")).To(BeTrue())
		Expect(hasTable(db, "bursts")).To(BeTrue())
		Expect(hasTable(db, "facts")).To(BeTrue())
		Expect(hasTable(db, "skills")).To(BeTrue())
	})

	It("is idempotent when run twice", func() {
		err := RunMigrationsForTests(db)
		Expect(err).NotTo(HaveOccurred())

		err = RunMigrationsForTests(db)
		Expect(err).NotTo(HaveOccurred())
	})
})
