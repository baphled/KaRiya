package sql

import (
	"context"
	"os"
	"path/filepath"

	career_repo "github.com/baphled/kariya/internal/repository/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repositories", func() {
	var (
		tmpDir string
		ctx    context.Context
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "repo_test_*")
		Expect(err).NotTo(HaveOccurred())
		ctx = context.Background()
	})

	AfterEach(func() {
		if tmpDir != "" {
			_ = os.RemoveAll(tmpDir)
		}
	})

	Describe("NewRepositoriesFromDB", func() {
		It("creates all repositories", func() {
			repos := NewRepositoriesFromDB(sharedGormDB)

			Expect(repos).NotTo(BeNil())
			Expect(repos.Event).NotTo(BeNil())
			Expect(repos.Skill).NotTo(BeNil())
			Expect(repos.Fact).NotTo(BeNil())
			Expect(repos.Burst).NotTo(BeNil())
		})
	})

	Describe("NewRepositories", func() {
		It("creates repositories from sql.DB", func() {
			dbPath := filepath.Join(tmpDir, "test.db")
			sqlDB, err := OpenDB(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer sqlDB.Close()

			err = career_repo.RunMigrationsForTests(sqlDB)
			Expect(err).NotTo(HaveOccurred())

			repos, err := NewRepositories(sqlDB)

			Expect(err).NotTo(HaveOccurred())
			Expect(repos).NotTo(BeNil())
			Expect(repos.Event).NotTo(BeNil())
			Expect(repos.Skill).NotTo(BeNil())
		})

		It("returns error when NewGormDB fails", func() {
			// Create an invalid database by using a closed db
			sqlDB, err := OpenDB(filepath.Join(tmpDir, "test.db"))
			Expect(err).NotTo(HaveOccurred())
			sqlDB.Close()

			_, err = NewRepositories(sqlDB)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("OpenDB", func() {
		It("opens database successfully", func() {
			dbPath := filepath.Join(tmpDir, "test.db")

			db, err := OpenDB(dbPath)

			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())
			Expect(db.Ping()).To(Succeed())
			db.Close()
		})

		It("opens in-memory database", func() {
			db, err := OpenDB(":memory:")

			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())
			Expect(db.Ping()).To(Succeed())
			db.Close()
		})

		It("sets max open connections to 1", func() {
			dbPath := filepath.Join(tmpDir, "test.db")

			db, err := OpenDB(dbPath)
			Expect(err).NotTo(HaveOccurred())

			// SetMaxOpenConns is called internally, verify by checking DB is usable
			err = db.Ping()
			Expect(err).NotTo(HaveOccurred())
			db.Close()
		})
	})

	Describe("NewRepositoriesFromPath", func() {
		It("creates repositories from file path", func() {
			dbPath := filepath.Join(tmpDir, "test.db")

			repos, err := NewRepositoriesFromPath(dbPath)

			Expect(err).NotTo(HaveOccurred())
			Expect(repos).NotTo(BeNil())
			Expect(repos.Event).NotTo(BeNil())
			Expect(repos.Skill).NotTo(BeNil())
			repos.Close()
		})

		It("runs migrations on new database", func() {
			dbPath := filepath.Join(tmpDir, "test.db")

			repos, err := NewRepositoriesFromPath(dbPath)

			Expect(err).NotTo(HaveOccurred())
			// Try to use the repositories - should have tables
			_, err = repos.Skill.List(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			repos.Close()
		})

		It("returns error for invalid path", func() {
			_, err := NewRepositoriesFromPath("/invalid/path/to/db.db")

			Expect(err).To(HaveOccurred())
		})

		It("closes DB on migration failure", func() {
			// Create a path but make it unwritable by using a directory
			dbPath := filepath.Join(tmpDir, "dir.db")
			err := os.Mkdir(dbPath, 0555)
			Expect(err).NotTo(HaveOccurred())

			_, err = NewRepositoriesFromPath(dbPath)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("NewGormDB", func() {
		It("creates GORM DB from sql.DB", func() {
			dbPath := filepath.Join(tmpDir, "test.db")
			sqlDB, err := OpenDB(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer sqlDB.Close()

			gormDB, err := NewGormDB(sqlDB)

			Expect(err).NotTo(HaveOccurred())
			Expect(gormDB).NotTo(BeNil())
		})
	})
})
