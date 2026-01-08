package career

import (
	"database/sql"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "modernc.org/sqlite"
)

var _ = Describe("Migrator", func() {
	var (
		db      *sql.DB
		tempDir string
		dbPath  string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-migrator-test-")
		Expect(err).NotTo(HaveOccurred())
		dbPath = filepath.Join(tempDir, "test.db")
	})

	AfterEach(func() {
		if db != nil {
			db.Close()
		}
		os.RemoveAll(tempDir)
	})

	Describe("RunMigrations", func() {
		Context("with a fresh database", func() {
			It("should apply all migrations successfully", func() {
				var err error
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				// Run migrations
				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Verify all tables exist
				tables := []string{"career_events", "bursts", "facts", "goose_db_version"}
				for _, table := range tables {
					Expect(hasTable(db, table)).To(BeTrue(), "table %s should exist", table)
				}

				// Verify categories column exists
				Expect(hasColumn(db, "career_events", "categories")).To(BeTrue())

				// Verify migration version
				version, err := MigrationStatus(db)
				Expect(err).NotTo(HaveOccurred())
				Expect(version).To(Equal(int64(4)))
			})

			It("should create indexes", func() {
				var err error
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Verify indexes exist
				indexes := []string{
					"idx_career_events_date",
					"idx_career_events_company",
					"idx_bursts_confirmed",
					"idx_facts_source_event_id",
					"idx_facts_source_burst_id",
				}

				for _, index := range indexes {
					var count int
					err := db.QueryRow(`
						SELECT COUNT(*) FROM sqlite_master 
						WHERE type='index' AND name=?
					`, index).Scan(&count)
					Expect(err).NotTo(HaveOccurred())
					Expect(count).To(Equal(1), "index %s should exist", index)
				}
			})
		})

		Context("with an existing database (pre-goose)", func() {
			It("should detect baseline and not re-create tables", func() {
				var err error
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				// Manually create career_events table (simulating old database)
				_, err = db.Exec(`
					CREATE TABLE career_events (
						id TEXT PRIMARY KEY,
						text TEXT NOT NULL,
						date DATETIME NOT NULL,
						tags TEXT,
						company TEXT,
						project TEXT,
						created_at DATETIME NOT NULL,
						updated_at DATETIME NOT NULL
					)
				`)
				Expect(err).NotTo(HaveOccurred())

				// Insert a test row
				_, err = db.Exec(`
					INSERT INTO career_events (id, text, date, tags, company, project, created_at, updated_at)
					VALUES ('test-id', 'Test event', '2024-01-01', 'tag1,tag2', 'Test Co', 'Test Project', '2024-01-01', '2024-01-01')
				`)
				Expect(err).NotTo(HaveOccurred())

				// Run migrations (should detect baseline)
				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Verify test row still exists (data not lost)
				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM career_events WHERE id = 'test-id'").Scan(&count)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))

				// Verify categories column was added
				Expect(hasColumn(db, "career_events", "categories")).To(BeTrue())

				// Verify all tables exist
				Expect(hasTable(db, "bursts")).To(BeTrue())
				Expect(hasTable(db, "facts")).To(BeTrue())
			})

			It("should detect categories column and set correct baseline", func() {
				var err error
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				// Create career_events with categories (simulating database at version 2)
				_, err = db.Exec(`
					CREATE TABLE career_events (
						id TEXT PRIMARY KEY,
						text TEXT NOT NULL,
						date DATETIME NOT NULL,
						tags TEXT,
						categories TEXT,
						company TEXT,
						project TEXT,
						created_at DATETIME NOT NULL,
						updated_at DATETIME NOT NULL
					)
				`)
				Expect(err).NotTo(HaveOccurred())

				// Run migrations
				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Verify migration version (should skip migration 002)
				version, err := MigrationStatus(db)
				Expect(err).NotTo(HaveOccurred())
				Expect(version).To(Equal(int64(4)))

				// Verify all subsequent tables were created
				Expect(hasTable(db, "bursts")).To(BeTrue())
				Expect(hasTable(db, "facts")).To(BeTrue())
			})

			It("should detect all existing tables and mark baseline as version 4", func() {
				var err error
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				// Create all tables (simulating fully migrated database before goose)
				_, err = db.Exec(`
					CREATE TABLE career_events (
						id TEXT PRIMARY KEY,
						text TEXT NOT NULL,
						date DATETIME NOT NULL,
						tags TEXT,
						categories TEXT,
						company TEXT,
						project TEXT,
						created_at DATETIME NOT NULL,
						updated_at DATETIME NOT NULL
					)
				`)
				Expect(err).NotTo(HaveOccurred())

				_, err = db.Exec(`
					CREATE TABLE bursts (
						id TEXT PRIMARY KEY,
						name TEXT NOT NULL,
						description TEXT,
						event_ids TEXT NOT NULL,
						confirmed INTEGER NOT NULL DEFAULT 0,
						confirmed_at DATETIME,
						created_at DATETIME NOT NULL,
						updated_at DATETIME NOT NULL
					)
				`)
				Expect(err).NotTo(HaveOccurred())

				_, err = db.Exec(`
					CREATE TABLE facts (
						id TEXT PRIMARY KEY,
						text TEXT NOT NULL,
						competencies TEXT NOT NULL,
						role_fit TEXT NOT NULL,
						audience_relevance TEXT NOT NULL,
						strength_signal TEXT,
						source_event_id TEXT,
						source_burst_id TEXT,
						created_at DATETIME NOT NULL,
						updated_at DATETIME NOT NULL
					)
				`)
				Expect(err).NotTo(HaveOccurred())

				// Run migrations
				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Verify baseline version is 4 (all migrations already applied)
				version, err := MigrationStatus(db)
				Expect(err).NotTo(HaveOccurred())
				Expect(version).To(Equal(int64(4)))
			})
		})

		Context("with goose already initialized", func() {
			It("should not interfere with existing goose tracking", func() {
				var err error
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				// Run migrations first time
				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Close and reopen database
				db.Close()
				db, err = sql.Open("sqlite", dbPath)
				Expect(err).NotTo(HaveOccurred())

				// Run migrations again (should be no-op)
				err = RunMigrations(db)
				Expect(err).NotTo(HaveOccurred())

				// Verify version is still 4
				version, err := MigrationStatus(db)
				Expect(err).NotTo(HaveOccurred())
				Expect(version).To(Equal(int64(4)))
			})
		})
	})

	Describe("MigrationStatus", func() {
		It("should return 0 for a fresh database before migrations", func() {
			var err error
			db, err = sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Don't run migrations, just check status
			version, err := MigrationStatus(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal(int64(0)))
		})

		It("should return correct version after migrations", func() {
			var err error
			db, err = sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())

			err = RunMigrations(db)
			Expect(err).NotTo(HaveOccurred())

			version, err := MigrationStatus(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal(int64(4)))
		})
	})

	Describe("Helper Functions", func() {
		BeforeEach(func() {
			var err error
			db, err = sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())

			// Create a test table
			_, err = db.Exec(`
				CREATE TABLE test_table (
					id TEXT PRIMARY KEY,
					name TEXT,
					value INTEGER
				)
			`)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("hasTable", func() {
			It("should return true for existing table", func() {
				Expect(hasTable(db, "test_table")).To(BeTrue())
			})

			It("should return false for non-existing table", func() {
				Expect(hasTable(db, "non_existing_table")).To(BeFalse())
			})
		})

		Describe("hasColumn", func() {
			It("should return true for existing column", func() {
				Expect(hasColumn(db, "test_table", "name")).To(BeTrue())
				Expect(hasColumn(db, "test_table", "value")).To(BeTrue())
			})

			It("should return false for non-existing column", func() {
				Expect(hasColumn(db, "test_table", "non_existing_column")).To(BeFalse())
			})

			It("should return false for non-existing table", func() {
				Expect(hasColumn(db, "non_existing_table", "name")).To(BeFalse())
			})
		})
	})
})
