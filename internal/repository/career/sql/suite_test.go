package sql

import (
	stdsql "database/sql"
	"os"
	"path/filepath"
	"testing"

	careerrepo "github.com/baphled/kariya/internal/repository/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var sharedGormDB *gorm.DB
var sharedSQLDB *stdsql.DB
var sharedTmpDir string

var _ = BeforeSuite(func() {
	var err error
	sharedTmpDir, err = os.MkdirTemp("", "sql_repo_suite_*")
	Expect(err).NotTo(HaveOccurred())

	dbPath := filepath.Join(sharedTmpDir, "sql_repo_suite.db")

	sharedSQLDB, err = OpenDB(dbPath)
	Expect(err).NotTo(HaveOccurred())

	err = careerrepo.RunMigrationsForTests(sharedSQLDB)
	Expect(err).NotTo(HaveOccurred())

	sharedGormDB, err = NewGormDB(sharedSQLDB)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if sharedSQLDB != nil {
		_ = sharedSQLDB.Close()
	}
	if sharedTmpDir != "" {
		_ = os.RemoveAll(sharedTmpDir)
	}
	sharedGormDB = nil
	sharedSQLDB = nil
})

func TestSQLRepositorySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQL Repository Suite")
}
