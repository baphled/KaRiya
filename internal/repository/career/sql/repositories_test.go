package sql

import (
	stdsql "database/sql"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "modernc.org/sqlite"
)

var _ = Describe("Repositories", func() {
	Describe("NewGormDB", func() {
		It("creates a GORM DB from a sql.DB", func() {
			sqlDB, err := stdsql.Open("sqlite", ":memory:")
			Expect(err).NotTo(HaveOccurred())
			defer sqlDB.Close()

			gormDB, err := NewGormDB(sqlDB)

			Expect(err).NotTo(HaveOccurred())
			Expect(gormDB).NotTo(BeNil())
		})
	})

	Describe("NewRepositoriesFromDB", func() {
		It("creates all repositories from a GORM DB", func() {
			sqlDB, err := stdsql.Open("sqlite", ":memory:")
			Expect(err).NotTo(HaveOccurred())
			defer sqlDB.Close()

			gormDB, err := NewGormDB(sqlDB)
			Expect(err).NotTo(HaveOccurred())

			repos := NewRepositoriesFromDB(gormDB)

			Expect(repos).NotTo(BeNil())
			Expect(repos.Event).NotTo(BeNil())
			Expect(repos.Skill).NotTo(BeNil())
			Expect(repos.Fact).NotTo(BeNil())
			Expect(repos.Burst).NotTo(BeNil())
		})
	})

	Describe("NewRepositories", func() {
		It("creates repositories from a sql.DB", func() {
			sqlDB, err := stdsql.Open("sqlite", ":memory:")
			Expect(err).NotTo(HaveOccurred())
			defer sqlDB.Close()

			repos, err := NewRepositories(sqlDB)

			Expect(err).NotTo(HaveOccurred())
			Expect(repos).NotTo(BeNil())
			Expect(repos.Event).NotTo(BeNil())
			Expect(repos.Skill).NotTo(BeNil())
			Expect(repos.Fact).NotTo(BeNil())
			Expect(repos.Burst).NotTo(BeNil())
		})
	})

	Describe("OpenDB", func() {
		It("opens a SQLite database at the given path", func() {
			tmpDir, err := os.MkdirTemp("", "kariya-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "test.db")
			db, err := OpenDB(dbPath)

			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())
			defer db.Close()

			err = db.Ping()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("NewRepositoriesFromPath", func() {
		It("creates repositories from a database file path", func() {
			tmpDir, err := os.MkdirTemp("", "kariya-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "test.db")
			repos, err := NewRepositoriesFromPath(dbPath)

			Expect(err).NotTo(HaveOccurred())
			Expect(repos).NotTo(BeNil())
			Expect(repos.Event).NotTo(BeNil())
			Expect(repos.Skill).NotTo(BeNil())
			Expect(repos.Fact).NotTo(BeNil())
			Expect(repos.Burst).NotTo(BeNil())
			defer repos.Close()
		})

		It("returns error for invalid path", func() {
			_, err := NewRepositoriesFromPath("/dev/null/impossible/path.db")

			Expect(err).To(HaveOccurred())
		})
	})
})
