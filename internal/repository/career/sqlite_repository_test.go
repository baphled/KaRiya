package career

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("SQLite Repository", func() {
	var (
		repo    *SQLiteRepository
		ctx     context.Context
		tempDir string
		dbPath  string
	)

	BeforeEach(func() {
		var err error
		// Create a temporary directory for test database
		tempDir, err = os.MkdirTemp("", "kariya-sqlite-test-")
		Expect(err).NotTo(HaveOccurred())

		// Create database path
		dbPath = filepath.Join(tempDir, "test_events.db")

		// Create SQLite repository
		repo, err = NewSQLiteRepository(dbPath)
		Expect(err).NotTo(HaveOccurred())

		ctx = context.Background()
	})

	AfterEach(func() {
		// Clean up: close repository and remove temporary files
		if repo != nil {
			repo.Close()
		}
		os.RemoveAll(tempDir)
	})

	// Helper function to generate test career event
	createSQLiteTestEvent := func() *career.CareerEvent {
		return &career.CareerEvent{
			Text:    "Test Career Event",
			Date:    time.Now(),
			Tags:    []string{"project", "technical"},
			Company: "Test Company",
		}
	}

	Describe("Create", func() {
		It("should create an event and retrieve it successfully", func() {
			event := createSQLiteTestEvent()

			// Create event
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve event
			retrievedEvent, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Validate retrieved event
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Date.Unix()).To(Equal(event.Date.Unix()))
			Expect(retrievedEvent.Tags).To(Equal(event.Tags))
			Expect(retrievedEvent.Company).To(Equal(event.Company))
		})

		It("should prevent duplicate event creation", func() {
			event := createSQLiteTestEvent()
			event.ID = "fixed-id"

			// First creation should succeed
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Second creation with same ID should fail
			err = repo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrDuplicateEvent)).To(BeTrue())
		})
	})

	Describe("Update", func() {
		It("should update event successfully", func() {
			event := createSQLiteTestEvent()

			// Create initial event
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Update event
			event.Text = "Updated Test Career Event"
			event.Tags = []string{"leadership", "project"}
			err = repo.Update(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve updated event
			retrievedEvent, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Validate updates
			Expect(retrievedEvent.Text).To(Equal("Updated Test Career Event"))
			Expect(retrievedEvent.Tags).To(Equal([]string{"leadership", "project"}))
		})

		It("should return error when updating non-existent event", func() {
			event := createSQLiteTestEvent()
			event.ID = "non-existent-id"

			// Update non-existent event should fail
			err := repo.Update(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrEventNotFound)).To(BeTrue())
		})
	})

	Describe("Delete", func() {
		It("should delete event successfully", func() {
			event := createSQLiteTestEvent()

			// Create event
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Delete event
			err = repo.Delete(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve should fail
			_, err = repo.GetByID(ctx, event.ID)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrEventNotFound)).To(BeTrue())
		})

		It("should return error when deleting non-existent event", func() {
			// Delete non-existent event should fail
			err := repo.Delete(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrEventNotFound)).To(BeTrue())
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create multiple test events
			events := []*career.CareerEvent{
				{
					Text:    "Event 1",
					Date:    time.Now().AddDate(0, 0, -30),
					Tags:    []string{"project", "technical"},
					Company: "Company A",
				},
				{
					Text:    "Event 2",
					Date:    time.Now(),
					Tags:    []string{"leadership", "project"},
					Company: "Company B",
				},
			}

			for _, event := range events {
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should filter events by tag", func() {
			// List events with filters
			listedEvents, err := repo.List(ctx, ListFilters{
				Tags: []string{"project"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate number of returned events
			Expect(listedEvents).To(HaveLen(2),
				"Should return both events with project tag")
		})

		It("should filter events by date range", func() {
			startDate := time.Now().AddDate(0, 0, -15)

			// List events with filters
			listedEvents, err := repo.List(ctx, ListFilters{
				StartDate: &startDate,
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate number of returned events
			Expect(listedEvents).To(HaveLen(1),
				"Should return only recent event within date range")
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			// Create multiple test events
			events := []*career.CareerEvent{
				{
					Text:    "Event 1",
					Date:    time.Now().AddDate(0, 0, -30),
					Tags:    []string{"project", "technical"},
					Company: "Company A",
				},
				{
					Text:    "Event 2",
					Date:    time.Now(),
					Tags:    []string{"leadership", "project"},
					Company: "Company B",
				},
			}

			for _, event := range events {
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should count all events", func() {
			// Count events with filters
			count, err := repo.Count(ctx, ListFilters{})
			Expect(err).NotTo(HaveOccurred())

			// Validate count of events
			Expect(count).To(Equal(2),
				"Should count all events")
		})

		It("should count events by tag", func() {
			// Count events with filters
			count, err := repo.Count(ctx, ListFilters{
				Tags: []string{"project"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate count of events
			Expect(count).To(Equal(2),
				"Should count both events with project tag")
		})

		It("should count events by date range", func() {
			startDate := time.Now().AddDate(0, 0, -15)

			// Count events with filters
			count, err := repo.Count(ctx, ListFilters{
				StartDate: &startDate,
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate count of events
			Expect(count).To(Equal(1),
				"Should count only event within recent date range")
		})
	})

	Describe("NewSQLiteRepositoryWithDB", func() {
		It("should create repository with existing database connection", func() {
			// Create a new repository using the existing connection from repo
			repoWithDB := NewSQLiteRepositoryWithDB(repo.GetDB())
			Expect(repoWithDB).NotTo(BeNil())

			// Test that it can perform basic operations
			event := createSQLiteTestEvent()
			err := repoWithDB.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
			Expect(event.ID).NotTo(BeEmpty())
		})

		It("should share database connection across operations", func() {
			// Use the existing database connection
			repoWithDB := NewSQLiteRepositoryWithDB(repo.GetDB())

			// Create multiple events using the same connection
			for i := 0; i < 3; i++ {
				event := createSQLiteTestEvent()
				err := repoWithDB.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			// Verify events were created
			count, err := repoWithDB.Count(ctx, ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeNumerically(">=", 3))
		})
	})

	Describe("GetDB", func() {
		It("should return the underlying database connection", func() {
			db := repo.GetDB()
			Expect(db).NotTo(BeNil())

			// Verify it's a functional connection
			err := db.Ping()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Skill Associations", func() {
		var skillRepo *SQLiteSkillRepository

		BeforeEach(func() {
			skillRepo = NewSQLiteSkillRepositoryWithDB(repo.GetDB())
		})

		Describe("Create with skills", func() {
			It("should save skill associations when creating event", func() {
				// Create skills first
				skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
				skill2 := &career.Skill{Name: "PostgreSQL", Category: "database"}

				err := skillRepo.Create(ctx, skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(ctx, skill2)
				Expect(err).NotTo(HaveOccurred())

				// Create event with skills
				event := createSQLiteTestEvent()
				event.Skills = []string{skill1.ID, skill2.ID}

				err = repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				// Retrieve event and verify skills are loaded
				retrieved, err := repo.GetByID(ctx, event.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Skills).To(HaveLen(2))
				Expect(retrieved.Skills).To(ContainElement(skill1.ID))
				Expect(retrieved.Skills).To(ContainElement(skill2.ID))
			})

			It("should create event without skills", func() {
				event := createSQLiteTestEvent()
				event.Skills = []string{}

				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				retrieved, err := repo.GetByID(ctx, event.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Skills).To(BeEmpty())
			})
		})

		Describe("GetByID with skills", func() {
			It("should load skill IDs when retrieving event", func() {
				// Create skills
				skill := &career.Skill{Name: "Go", Category: "backend"}
				err := skillRepo.Create(ctx, skill)
				Expect(err).NotTo(HaveOccurred())

				// Create event with skill
				event := createSQLiteTestEvent()
				event.Skills = []string{skill.ID}
				err = repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				// Retrieve and verify
				retrieved, err := repo.GetByID(ctx, event.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Skills).To(Equal([]string{skill.ID}))
			})
		})

		Describe("Update with skills", func() {
			It("should update skill associations", func() {
				// Create skills
				skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
				skill2 := &career.Skill{Name: "React", Category: "frontend"}
				skill3 := &career.Skill{Name: "Docker", Category: "devops"}

				err := skillRepo.Create(ctx, skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(ctx, skill2)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(ctx, skill3)
				Expect(err).NotTo(HaveOccurred())

				// Create event with initial skills
				event := createSQLiteTestEvent()
				event.Skills = []string{skill1.ID, skill2.ID}
				err = repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				// Update event skills
				event.Skills = []string{skill2.ID, skill3.ID} // Remove skill1, add skill3
				err = repo.Update(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				// Verify updated skills
				retrieved, err := repo.GetByID(ctx, event.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Skills).To(HaveLen(2))
				Expect(retrieved.Skills).To(ContainElement(skill2.ID))
				Expect(retrieved.Skills).To(ContainElement(skill3.ID))
				Expect(retrieved.Skills).NotTo(ContainElement(skill1.ID))
			})

			It("should allow removing all skills", func() {
				// Create skill and event
				skill := &career.Skill{Name: "Ruby", Category: "backend"}
				err := skillRepo.Create(ctx, skill)
				Expect(err).NotTo(HaveOccurred())

				event := createSQLiteTestEvent()
				event.Skills = []string{skill.ID}
				err = repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				// Remove all skills
				event.Skills = []string{}
				err = repo.Update(ctx, event)
				Expect(err).NotTo(HaveOccurred())

				// Verify skills removed
				retrieved, err := repo.GetByID(ctx, event.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Skills).To(BeEmpty())
			})
		})

		Describe("List with skills", func() {
			It("should load skills for all events in list", func() {
				// Create skills
				skill1 := &career.Skill{Name: "Ruby", Category: "backend"}
				skill2 := &career.Skill{Name: "React", Category: "frontend"}
				err := skillRepo.Create(ctx, skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(ctx, skill2)
				Expect(err).NotTo(HaveOccurred())

				// Create events with different skills
				event1 := createSQLiteTestEvent()
				event1.Skills = []string{skill1.ID}
				err = repo.Create(ctx, event1)
				Expect(err).NotTo(HaveOccurred())

				event2 := createSQLiteTestEvent()
				event2.Skills = []string{skill2.ID}
				err = repo.Create(ctx, event2)
				Expect(err).NotTo(HaveOccurred())

				// List and verify skills loaded
				filters := ListFilters{}
				events, err := repo.List(ctx, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(events).To(HaveLen(2))

				// Find our events and verify skills
				for _, event := range events {
					if event.ID == event1.ID {
						Expect(event.Skills).To(Equal([]string{skill1.ID}))
					} else if event.ID == event2.ID {
						Expect(event.Skills).To(Equal([]string{skill2.ID}))
					}
				}
			})
		})
	})
})
