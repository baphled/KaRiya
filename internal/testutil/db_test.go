package testutil

import (
	"database/sql"
	"testing"

	"github.com/baphled/kariya/internal/repository/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testutil Suite")
}

var _ = Describe("Test Database Utilities", func() {
	Describe("SetupTestDB", func() {
		It("should create a database with migrations applied", func() {
			// Use a standard Go test helper for Ginkgo
			db, cleanup := setupTestDBHelper()
			defer cleanup()

			// Verify database connection is valid
			Expect(db.Ping()).To(Succeed())

			// Verify tables exist
			tables := []string{"career_events", "bursts", "facts", "goose_db_version"}
			for _, table := range tables {
				var count int
				err := db.QueryRow(`
					SELECT COUNT(*) FROM sqlite_master 
					WHERE type='table' AND name=?
				`, table).Scan(&count)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1), "table %s should exist", table)
			}

			// Verify categories column exists
			rows, err := db.Query("PRAGMA table_info(career_events)")
			Expect(err).NotTo(HaveOccurred())
			defer rows.Close()

			columnExists := false
			for rows.Next() {
				var cid int
				var name, typ string
				var notnull, pk int
				var dflt interface{}
				err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk)
				Expect(err).NotTo(HaveOccurred())
				if name == "categories" {
					columnExists = true
					break
				}
			}
			Expect(columnExists).To(BeTrue(), "categories column should exist")
		})

		It("should provide a working cleanup function", func() {
			db, cleanup := setupTestDBHelper()

			// Database should be open
			Expect(db.Ping()).To(Succeed())

			// Call cleanup
			cleanup()

			// Database should be closed (Ping will fail)
			err := db.Ping()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SetupTestDBWithPath", func() {
		It("should return both path and database connection", func() {
			dbPath, db, cleanup := setupTestDBWithPathHelper()
			defer cleanup()

			// Verify path is not empty
			Expect(dbPath).NotTo(BeEmpty())

			// Verify database connection is valid
			Expect(db.Ping()).To(Succeed())

			// Verify we can query the database
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM career_events").Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0), "should start with no events")
		})

		It("should create the same database at the returned path", func() {
			dbPath, db1, cleanup1 := setupTestDBWithPathHelper()
			defer cleanup1()

			// Insert a test row
			_, err := db1.Exec(`
				INSERT INTO career_events (id, text, date, tags, company, project, created_at, updated_at)
				VALUES ('test-id', 'Test', '2024-01-01', '', '', '', '2024-01-01', '2024-01-01')
			`)
			Expect(err).NotTo(HaveOccurred())

			// Close first connection
			db1.Close()

			// Open the same database using the path
			db2, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db2.Close()

			// Verify data persisted
			var count int
			err = db2.QueryRow("SELECT COUNT(*) FROM career_events WHERE id = 'test-id'").Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
})

// Helper functions that work with Ginkgo's TempDir.
func setupTestDBHelper() (*sql.DB, func()) {
	tmpDir := GinkgoT().TempDir()
	dbPath := tmpDir + "/test.db"

	db, err := sql.Open("sqlite", dbPath)
	Expect(err).NotTo(HaveOccurred())

	err = career.RunMigrations(db)
	Expect(err).NotTo(HaveOccurred())

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func setupTestDBWithPathHelper() (string, *sql.DB, func()) {
	tmpDir := GinkgoT().TempDir()
	dbPath := tmpDir + "/test.db"

	db, err := sql.Open("sqlite", dbPath)
	Expect(err).NotTo(HaveOccurred())

	err = career.RunMigrations(db)
	Expect(err).NotTo(HaveOccurred())

	cleanup := func() {
		db.Close()
	}

	return dbPath, db, cleanup
}
