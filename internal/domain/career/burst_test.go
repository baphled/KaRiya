package career

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst", func() {
	Context("when creating a burst", func() {
		It("should have all required fields", func() {
			burst := &Burst{
				ID:              "burst-123",
				Name:            "Platform Migration",
				Description:     "Led team to migrate legacy platform to microservices",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "technical",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			Expect(burst.ID).To(Equal("burst-123"))
			Expect(burst.Name).To(Equal("Platform Migration"))
			Expect(burst.Description).To(Equal("Led team to migrate legacy platform to microservices"))
			Expect(burst.EventIDs).To(HaveLen(2))
			Expect(burst.CompetencyFocus).To(Equal("technical"))
			Expect(burst.CreatedAt).NotTo(BeZero())
			Expect(burst.UpdatedAt).NotTo(BeZero())
		})
	})

	Context("when validating a burst", func() {
		It("should validate successfully with minimum required fields", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject burst with empty ID", func() {
			burst := &Burst{
				Name:     "Project Name",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ID cannot be empty"))
		})

		It("should reject burst with empty name", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
		})

		It("should reject burst with whitespace-only name", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "   ",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
		})

		It("should reject burst with less than 2 event IDs", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: []string{"event-1"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least 2 events required"))
		})

		It("should reject burst with empty event IDs list", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: []string{},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least 2 events required"))
		})

		It("should reject burst with nil event IDs", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: nil,
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least 2 events required"))
		})

		It("should reject burst with duplicate event IDs", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: []string{"event-1", "event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate event IDs"))
		})

		It("should reject burst with empty event ID in list", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: []string{"event-1", "", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
		})

		It("should accept burst with valid competency focus", func() {
			burst := &Burst{
				ID:              "burst-123",
				Name:            "Project Name",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "technical",
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject burst with invalid competency focus", func() {
			burst := &Burst{
				ID:              "burst-123",
				Name:            "Project Name",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "invalid-competency",
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid competency focus"))
		})

		It("should accept burst with empty competency focus (optional)", func() {
			burst := &Burst{
				ID:              "burst-123",
				Name:            "Project Name",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "",
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject burst with name exceeding 200 characters", func() {
			longName := ""
			for i := 0; i < 201; i++ {
				longName += "a"
			}

			burst := &Burst{
				ID:       "burst-123",
				Name:     longName,
				EventIDs: []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name cannot exceed 200 characters"))
		})

		It("should accept burst with name at 200 character limit", func() {
			name := ""
			for i := 0; i < 200; i++ {
				name += "a"
			}

			burst := &Burst{
				ID:       "burst-123",
				Name:     name,
				EventIDs: []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject burst with description exceeding 1000 characters", func() {
			longDesc := ""
			for i := 0; i < 1001; i++ {
				longDesc += "a"
			}

			burst := &Burst{
				ID:          "burst-123",
				Name:        "Project Name",
				Description: longDesc,
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("description cannot exceed 1000 characters"))
		})

		It("should accept burst with description at 1000 character limit", func() {
			desc := ""
			for i := 0; i < 1000; i++ {
				desc += "a"
			}

			burst := &Burst{
				ID:          "burst-123",
				Name:        "Project Name",
				Description: desc,
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept burst with empty description (optional)", func() {
			burst := &Burst{
				ID:          "burst-123",
				Name:        "Project Name",
				Description: "",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept burst with many event IDs", func() {
			eventIDs := []string{}
			for i := 0; i < 50; i++ {
				eventIDs = append(eventIDs, "event-"+string(rune(i)))
			}

			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: eventIDs,
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when working with confirmed status", func() {
		It("should default to not confirmed", func() {
			burst := &Burst{
				ID:       "burst-123",
				Name:     "Project Name",
				EventIDs: []string{"event-1", "event-2"},
			}

			Expect(burst.Confirmed).To(BeFalse())
			Expect(burst.ConfirmedAt).To(BeNil())
		})

		It("should allow setting confirmed status", func() {
			now := time.Now()
			burst := &Burst{
				ID:          "burst-123",
				Name:        "Project Name",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   true,
				ConfirmedAt: &now,
			}

			Expect(burst.Confirmed).To(BeTrue())
			Expect(burst.ConfirmedAt).NotTo(BeNil())
			Expect(*burst.ConfirmedAt).To(BeTemporally("~", now, time.Second))
		})

		It("should validate successfully when confirmed", func() {
			now := time.Now()
			burst := &Burst{
				ID:          "burst-123",
				Name:        "Project Name",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   true,
				ConfirmedAt: &now,
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should validate successfully when confirmed without timestamp", func() {
			burst := &Burst{
				ID:        "burst-123",
				Name:      "Project Name",
				EventIDs:  []string{"event-1", "event-2"},
				Confirmed: true,
			}

			err := burst.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
