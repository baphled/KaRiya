package cv

import (
	"context"
	"testing"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigInitializer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ConfigInitializer Suite")
}

var _ = Describe("ConfigInitializer", func() {
	var (
		initializer *ConfigInitializer
		manager     *MemoryConfigManager
		log         *logger.Logger
		ctx         context.Context
	)

	BeforeEach(func() {
		manager = NewMemoryConfigManager()
		log = logger.DefaultLogger()
		ctx = context.Background()
		initializer = NewConfigInitializer(manager, log)
	})

	Describe("Initialize", func() {
		It("should create default configs when none exist", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(configs)).To(Equal(4))

			// Verify default configs were created
			names := make(map[string]bool)
			for _, cfg := range configs {
				names[cfg.Name] = true
			}

			Expect(names["Principal Engineer"]).To(BeTrue())
			Expect(names["Staff Engineer"]).To(BeTrue())
			Expect(names["Engineering Manager"]).To(BeTrue())
			Expect(names["Senior IC"]).To(BeTrue())
		})

		It("should not create defaults if configs already exist", func() {
			// Create one config
			config := &career.CVConfig{
				Name:           "existing-config",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Initialize
			err = initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			// Should have only 1 config (the existing one)
			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(configs)).To(Equal(1))
			Expect(configs[0].Name).To(Equal("existing-config"))
		})

		It("should handle context cancellation", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			err := initializer.Initialize(cancelCtx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(context.Canceled))
		})

		It("should set timestamps on default configs", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())

			for _, cfg := range configs {
				Expect(cfg.CreatedAt.IsZero()).To(BeFalse())
				Expect(cfg.UpdatedAt.IsZero()).To(BeFalse())
			}
		})

		It("should create valid default configs", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())

			for _, cfg := range configs {
				err := cfg.Validate()
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Describe("ValidateSetup", func() {
		It("should succeed if configuration system is accessible", func() {
			err := initializer.ValidateSetup(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle context cancellation", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			err := initializer.ValidateSetup(cancelCtx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(context.Canceled))
		})
	})

	Describe("EnsureConfigExists", func() {
		It("should create defaults if no configs exist", func() {
			err := initializer.EnsureConfigExists(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(configs)).To(Equal(4))
		})

		It("should not create defaults if configs already exist", func() {
			// Create one config
			config := &career.CVConfig{
				Name:           "existing-config",
				TargetRole:     "staff",
				TargetAudience: []string{"recruiter"},
			}
			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Ensure config exists
			err = initializer.EnsureConfigExists(ctx)
			Expect(err).NotTo(HaveOccurred())

			// Should still have only 1 config
			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(configs)).To(Equal(1))
		})

		It("should handle context cancellation", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			err := initializer.EnsureConfigExists(cancelCtx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(context.Canceled))
		})
	})

	Describe("CheckAndCreateConfig", func() {
		It("should return existing config if it exists", func() {
			// Create a config first
			existingConfig := &career.CVConfig{
				Name:           "test-config",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			err := manager.SaveConfig(ctx, existingConfig)
			Expect(err).NotTo(HaveOccurred())

			// Check and create should return the existing config
			config, err := initializer.CheckAndCreateConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.Name).To(Equal("test-config"))
			Expect(config.TargetRole).To(Equal("principal"))
		})

		It("should create a new config if it doesn't exist", func() {
			// Check and create a non-existent config
			config, err := initializer.CheckAndCreateConfig(ctx, "new-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())
			Expect(config.Name).To(Equal("new-config"))

			// Verify it was saved
			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(configs)).To(Equal(1))
			Expect(configs[0].Name).To(Equal("new-config"))
		})

		It("should reject empty config names", func() {
			config, err := initializer.CheckAndCreateConfig(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		})

		It("should create config with sensible defaults", func() {
			config, err := initializer.CheckAndCreateConfig(ctx, "new-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(config).NotTo(BeNil())

			// Check defaults
			Expect(config.TargetRole).To(Equal("staff"))
			Expect(config.TargetAudience).To(ContainElement("hiring_manager"))
			Expect(config.EventFilters).NotTo(BeNil())
		})

		It("should handle context cancellation", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()

			config, err := initializer.CheckAndCreateConfig(cancelCtx, "test-config")
			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		})

		It("should validate created config", func() {
			config, err := initializer.CheckAndCreateConfig(ctx, "valid-config")
			Expect(err).NotTo(HaveOccurred())

			// Validate the returned config
			err = config.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Default Configuration Content", func() {
		It("should create Principal Engineer config with correct role", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())

			var principalConfig *career.CVConfig
			for _, cfg := range configs {
				if cfg.Name == "Principal Engineer" {
					principalConfig = cfg
					break
				}
			}

			Expect(principalConfig).NotTo(BeNil())
			Expect(principalConfig.TargetRole).To(Equal("principal"))
			Expect(principalConfig.TargetAudience).To(ContainElement("hiring_manager"))
		})

		It("should create Staff Engineer config with correct audiences", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())

			var staffConfig *career.CVConfig
			for _, cfg := range configs {
				if cfg.Name == "Staff Engineer" {
					staffConfig = cfg
					break
				}
			}

			Expect(staffConfig).NotTo(BeNil())
			Expect(staffConfig.TargetRole).To(Equal("staff"))
			Expect(len(staffConfig.TargetAudience)).To(Equal(2))
			Expect(staffConfig.TargetAudience).To(ContainElement("hiring_manager"))
			Expect(staffConfig.TargetAudience).To(ContainElement("peer"))
		})

		It("should create Engineering Manager config with leadership focus", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())

			var emConfig *career.CVConfig
			for _, cfg := range configs {
				if cfg.Name == "Engineering Manager" {
					emConfig = cfg
					break
				}
			}

			Expect(emConfig).NotTo(BeNil())
			Expect(emConfig.TargetRole).To(Equal("em"))
			Expect(emConfig.TargetAudience).To(ContainElement("hiring_manager"))
		})

		It("should create Senior IC config with recruiter audience", func() {
			err := initializer.Initialize(ctx)
			Expect(err).NotTo(HaveOccurred())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())

			var seniorConfig *career.CVConfig
			for _, cfg := range configs {
				if cfg.Name == "Senior IC" {
					seniorConfig = cfg
					break
				}
			}

			Expect(seniorConfig).NotTo(BeNil())
			Expect(seniorConfig.TargetRole).To(Equal("senior_ic"))
			Expect(seniorConfig.TargetAudience).To(ContainElement("recruiter"))
		})
	})
})
