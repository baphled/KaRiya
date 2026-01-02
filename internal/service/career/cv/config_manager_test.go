package cv

import (
	"context"
	"os"
	"path/filepath"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)


var _ = Describe("YAMLConfigManager", func() {
	var (
		manager   *YAMLConfigManager
		tempDir   string
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-cv-test-*")
		Expect(err).NotTo(HaveOccurred())

		log = logger.DefaultLogger()
		ctx = context.Background()

		// Create manager with temp directory
		manager = &YAMLConfigManager{
			configDir: tempDir,
			logger:    log,
		}
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("SaveConfig", func() {
		It("should save a valid configuration", func() {
			config := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager"},
				EventFilters: map[string]interface{}{
					"tags": []string{"technical"},
				},
			}

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Verify file exists
			filePath := manager.GetConfigPath("test-config")
			_, err = os.Stat(filePath)
			Expect(err).NotTo(HaveOccurred())

			// Verify timestamps were set
			Expect(config.CreatedAt).NotTo(BeZero())
			Expect(config.UpdatedAt).NotTo(BeZero())
		})

		It("should overwrite existing configuration", func() {
			config1 := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			err := manager.SaveConfig(ctx, config1)
			Expect(err).NotTo(HaveOccurred())

			createdAtBefore := config1.CreatedAt
			Expect(createdAtBefore).NotTo(BeZero())

			// Wait a bit and load the config
			time.Sleep(10 * time.Millisecond)

			// Load existing config to preserve CreatedAt
			loaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())

			// Update with new data
			loaded.TargetRole = "principal"
			loaded.TargetAudience = []string{"recruiter"}

			err = manager.SaveConfig(ctx, loaded)
			Expect(err).NotTo(HaveOccurred())

			// Reload and verify updates
			reloaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(reloaded.TargetRole).To(Equal("principal"))
			Expect(reloaded.TargetAudience).To(Equal([]string{"recruiter"}))
			// CreatedAt should be approximately the same (within 1 second)
			Expect(reloaded.CreatedAt.Sub(createdAtBefore).Abs()).To(BeNumerically("<", time.Second))
		})

		It("should reject nil config", func() {
			err := manager.SaveConfig(ctx, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should reject invalid config", func() {
			config := &career.CVConfig{
				Name: "test-config",
				// Missing required fields
			}

			err := manager.SaveConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should create directory if needed", func() {
			subDir := filepath.Join(tempDir, "subdir", "test")
			manager.configDir = subDir

			config := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Verify directory was created
			_, err = os.Stat(subDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should use atomic write", func() {
			config := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			// Count temp files before
			entries, _ := os.ReadDir(tempDir)
			tempBefore := len(entries)

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Count temp files after - should not have leftover temp files
			entries, _ = os.ReadDir(tempDir)
			tempAfter := len(entries)

			// Should have exactly one file (.yaml) and no temp files
			Expect(tempAfter).To(Equal(tempBefore + 1))
		})
	})

	Describe("LoadConfig", func() {
		It("should load an existing configuration", func() {
			config := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager", "recruiter"},
				EventFilters: map[string]interface{}{
					"tags": []string{"technical", "leadership"},
				},
			}

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.Name).To(Equal("test-config"))
			Expect(loaded.TargetRole).To(Equal("principal"))
			Expect(loaded.TargetAudience).To(ContainElement("hiring_manager"))
			Expect(loaded.TargetAudience).To(ContainElement("recruiter"))
			Expect(loaded.EventFilters).NotTo(BeNil())
		})

		It("should return ErrConfigNotFound for missing config", func() {
			_, err := manager.LoadConfig(ctx, "nonexistent")
			Expect(err).To(Equal(ErrConfigNotFound))
		})

		It("should reject empty name", func() {
			_, err := manager.LoadConfig(ctx, "")
			Expect(err).To(HaveOccurred())
		})

		It("should reject whitespace name", func() {
			_, err := manager.LoadConfig(ctx, "   ")
			Expect(err).To(HaveOccurred())
		})

		It("should handle context cancellation", func() {
			config := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			manager.SaveConfig(context.Background(), config)

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := manager.LoadConfig(cancelCtx, "test-config")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteConfig", func() {
		It("should delete an existing configuration", func() {
			config := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			manager.SaveConfig(ctx, config)
			filePath := manager.GetConfigPath("test-config")

			// Verify file exists
			_, err := os.Stat(filePath)
			Expect(err).NotTo(HaveOccurred())

			// Delete
			err = manager.DeleteConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())

			// Verify file is deleted
			_, err = os.Stat(filePath)
			Expect(err).To(HaveOccurred())
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("should return ErrConfigNotFound for missing config", func() {
			err := manager.DeleteConfig(ctx, "nonexistent")
			Expect(err).To(Equal(ErrConfigNotFound))
		})

		It("should reject empty name", func() {
			err := manager.DeleteConfig(ctx, "")
			Expect(err).To(HaveOccurred())
		})

		It("should reject whitespace name", func() {
			err := manager.DeleteConfig(ctx, "   ")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ListConfigs", func() {
		It("should list all configurations", func() {
			configs := []*career.CVConfig{
				{
					Name:           "config1",
					TargetRole:     "principal",
					TargetAudience: []string{"hiring_manager"},
				},
				{
					Name:           "config2",
					TargetRole:     "principal",
					TargetAudience: []string{"recruiter"},
				},
				{
					Name:           "config3",
					TargetRole:     "em",
					TargetAudience: []string{"peer"},
				},
			}

			for _, config := range configs {
				err := manager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(loaded)).To(Equal(3))

			names := make(map[string]bool)
			for _, config := range loaded {
				names[config.Name] = true
			}

			Expect(names["config1"]).To(BeTrue())
			Expect(names["config2"]).To(BeTrue())
			Expect(names["config3"]).To(BeTrue())
		})

		It("should return empty list for empty directory", func() {
			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(loaded)).To(Equal(0))
		})

		It("should skip non-YAML files", func() {
			// Save a config
			config := &career.CVConfig{
				Name:           "config1",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Create a non-YAML file
			nonYamlPath := filepath.Join(tempDir, "readme.txt")
			err = os.WriteFile(nonYamlPath, []byte("test"), 0644)
			Expect(err).NotTo(HaveOccurred())

			// List should only return the YAML config
			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(loaded)).To(Equal(1))
			Expect(loaded[0].Name).To(Equal("config1"))
		})

		It("should skip invalid configs during list", func() {
			// Save valid config
			config1 := &career.CVConfig{
				Name:           "config1",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}
			manager.SaveConfig(ctx, config1)

			// Manually create an invalid YAML file
			invalidPath := filepath.Join(tempDir, "invalid.yaml")
			os.WriteFile(invalidPath, []byte("invalid: yaml: [[["), 0644)

			// List should return valid configs and skip invalid ones
			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(loaded)).To(Equal(1))
			Expect(loaded[0].Name).To(Equal("config1"))
		})

		It("should handle context cancellation", func() {
			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := manager.ListConfigs(cancelCtx)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetConfigPath", func() {
		It("should return correct path for config name", func() {
			path := manager.GetConfigPath("my-config")
			Expect(path).To(ContainSubstring("my-config.yaml"))
			Expect(path).To(ContainSubstring(tempDir))
		})

		It("should sanitize dangerous characters", func() {
			path := manager.GetConfigPath("config/with\bad:chars*?")
			// Check that the filename part (after tempDir) has dangerous chars replaced
			filename := path[len(tempDir)+1:] // Skip the directory path
			Expect(filename).To(ContainSubstring("_"))
			Expect(filename).To(ContainSubstring(".yaml"))
			Expect(filename).NotTo(ContainSubstring(":"))
			Expect(filename).NotTo(ContainSubstring("*"))
			Expect(filename).NotTo(ContainSubstring("?"))
		})

		It("should handle whitespace", func() {
			path := manager.GetConfigPath("  my-config  ")
			Expect(path).To(ContainSubstring("my-config.yaml"))
		})
	})

	Describe("ConfigExists", func() {
		It("should return true for existing config", func() {
			config := &career.CVConfig{
				Name:           "test-config",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			manager.SaveConfig(ctx, config)

			exists, err := manager.ConfigExists(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should return false for non-existing config", func() {
			exists, err := manager.ConfigExists(ctx, "nonexistent")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should reject empty name", func() {
			_, err := manager.ConfigExists(ctx, "")
			Expect(err).To(HaveOccurred())
		})

		It("should handle context cancellation", func() {
			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := manager.ConfigExists(cancelCtx, "test")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("YAML Serialization", func() {
		It("should preserve all config fields", func() {
			originalConfig := &career.CVConfig{
				Name:         "test-config",
				TargetRole:   "principal",
				TargetAudience: []string{"hiring_manager", "recruiter"},
				EventFilters: map[string]interface{}{
					"tags":       []string{"technical", "leadership"},
					"companies":  []string{"Google", "Meta"},
					"minDate":    "2020-01-01",
					"maxDate":    "2023-12-31",
				},
			}

			err := manager.SaveConfig(ctx, originalConfig)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())

			Expect(loaded.Name).To(Equal(originalConfig.Name))
			Expect(loaded.TargetRole).To(Equal(originalConfig.TargetRole))
			Expect(loaded.TargetAudience).To(Equal(originalConfig.TargetAudience))
			// Note: YAML unmarshaling converts []string to []interface{}, so just verify keys exist
			Expect(loaded.EventFilters).To(HaveKey("tags"))
			Expect(loaded.EventFilters).To(HaveKey("companies"))
			Expect(loaded.EventFilters).To(HaveKey("minDate"))
			Expect(loaded.EventFilters).To(HaveKey("maxDate"))
		})

		It("should handle complex nested filters", func() {
			config := &career.CVConfig{
				Name:         "complex-config",
				TargetRole:   "em",
				TargetAudience: []string{"peer"},
				EventFilters: map[string]interface{}{
					"tags":         []string{"leadership", "mentoring"},
					"categories":   []string{"leadership", "product"},
					"companies":    []string{"Company A", "Company B"},
					"roles":        []string{"staff", "principal"},
					"dateRange":    map[string]string{"start": "2021-01-01", "end": "2024-12-31"},
				},
			}

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "complex-config")
			Expect(err).NotTo(HaveOccurred())

			Expect(loaded.EventFilters["companies"]).To(Equal([]interface{}{"Company A", "Company B"}))
		})
	})

	Describe("Security checks", func() {
		It("should prevent path traversal in GetConfigPath", func() {
			path := manager.GetConfigPath("../../../etc/passwd")
			Expect(path).To(ContainSubstring(tempDir))
			// The .. characters are replaced with underscores during sanitization
			Expect(path).To(HavePrefix(tempDir))
		})

		It("should prevent path traversal in LoadConfig", func() {
			_, err := manager.LoadConfig(ctx, "../../../etc/passwd")
			Expect(err).To(HaveOccurred())
		})

		It("should prevent path traversal in SaveConfig", func() {
			config := &career.CVConfig{
				Name:           "../../../etc/passwd",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager"},
			}

			err := manager.SaveConfig(ctx, config)
			// Should either fail or sanitize the name
			if err == nil {
				// If it succeeds, verify the file is in the correct directory
				filePath := manager.GetConfigPath(config.Name)
				Expect(isPathWithinDirectory(filePath, tempDir)).To(BeTrue())
			}
		})

		It("should prevent path traversal in DeleteConfig", func() {
			err := manager.DeleteConfig(ctx, "../../../etc/passwd")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("NewYAMLConfigManager", func() {
	It("should create config directory if it doesn't exist", func() {
		tempDir, err := os.MkdirTemp("", "kariya-cv-test-*")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(tempDir)

		Expect(os.RemoveAll(filepath.Join(tempDir, "kariya"))).NotTo(HaveOccurred())

		// This would require modifying the manager to use a custom home dir
		// For now, we test that it succeeds with the real home dir
		log := logger.DefaultLogger()
		manager, err := NewYAMLConfigManager(log)
		Expect(err).NotTo(HaveOccurred())
		Expect(manager).NotTo(BeNil())

		// Verify directory exists
		_, err = os.Stat(manager.configDir)
		Expect(err).NotTo(HaveOccurred())
	})
})

