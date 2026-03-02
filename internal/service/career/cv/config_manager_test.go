//nolint:errcheck // Test file - error handling for test setup is not relevant.
package cv

import (
	"context"
	"os"
	"path/filepath"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("YAMLConfigManager", func() {
	var (
		manager *YAMLConfigManager
		tempDir string
		log     *logger.Logger
		ctx     context.Context
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
			config := fixtures.CVConfigWithFilters("test-config", map[string]interface{}{
				"tags": []string{"technical"},
			})
			config.TargetRole = "principal"

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
			config1 := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

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
			loaded.TargetAudience = "recruiter"

			err = manager.SaveConfig(ctx, loaded)
			Expect(err).NotTo(HaveOccurred())

			// Reload and verify updates
			reloaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(reloaded.TargetRole).To(Equal("principal"))
			Expect(reloaded.TargetAudience).To(Equal("recruiter"))
			// CreatedAt should be approximately the same (within 1 second)
			Expect(reloaded.CreatedAt.Sub(createdAtBefore).Abs()).To(BeNumerically("<", time.Second))
		})

		It("should reject nil config", func() {
			err := manager.SaveConfig(ctx, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should reject invalid config", func() {
			config := fixtures.CVConfig("test-config")
			config.TargetRole = ""
			config.TargetAudience = ""

			err := manager.SaveConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should create directory if needed", func() {
			subDir := filepath.Join(tempDir, "subdir", "test")
			manager.configDir = subDir

			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Verify directory was created
			_, err = os.Stat(subDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should use atomic write", func() {
			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

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
			config := fixtures.CVConfigWithFilters("test-config", map[string]interface{}{
				"tags": []string{"technical", "leadership"},
			})
			config.TargetRole = "principal"

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.Name).To(Equal("test-config"))
			Expect(loaded.TargetRole).To(Equal("principal"))
			Expect(loaded.TargetAudience).To(Equal("hiring_manager"))
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
			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

			Expect(manager.SaveConfig(context.Background(), config)).To(Succeed())

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := manager.LoadConfig(cancelCtx, "test-config")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteConfig", func() {
		It("should delete an existing configuration", func() {
			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

			Expect(manager.SaveConfig(ctx, config)).To(Succeed())
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
				fixtures.CVConfigWith("config1", "principal", "hiring_manager"),
				fixtures.CVConfigWith("config2", "principal", "recruiter"),
				fixtures.CVConfigWith("config3", "em", "peer"),
			}

			for _, config := range configs {
				err := manager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
			}

			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded).To(HaveLen(3))

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
			Expect(loaded).To(BeEmpty())
		})

		It("should skip non-YAML files", func() {
			// Save a config
			config := fixtures.CVConfigWith("config1", "principal", "hiring_manager")
			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// Create a non-YAML file
			nonYamlPath := filepath.Join(tempDir, "readme.txt")
			err = os.WriteFile(nonYamlPath, []byte("test"), 0o600)
			Expect(err).NotTo(HaveOccurred())

			// List should only return the YAML config
			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded).To(HaveLen(1))
			Expect(loaded[0].Name).To(Equal("config1"))
		})

		It("should skip invalid configs during list", func() {
			// Save valid config
			config1 := fixtures.CVConfigWith("config1", "principal", "hiring_manager")
			Expect(manager.SaveConfig(ctx, config1)).To(Succeed())

			// Manually create an invalid YAML file
			invalidPath := filepath.Join(tempDir, "invalid.yaml")
			Expect(os.WriteFile(invalidPath, []byte("invalid: yaml: [[["), 0o600)).To(Succeed())

			// List should return valid configs and skip invalid ones
			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded).To(HaveLen(1))
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
			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

			Expect(manager.SaveConfig(ctx, config)).To(Succeed())

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
			originalConfig := fixtures.CVConfigWithFilters("test-config", map[string]interface{}{
				"tags":      []string{"technical", "leadership"},
				"companies": []string{"Google", "Meta"},
				"minDate":   "2020-01-01",
				"maxDate":   "2023-12-31",
			})
			originalConfig.TargetRole = "principal"

			err := manager.SaveConfig(ctx, originalConfig)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "test-config")
			Expect(err).NotTo(HaveOccurred())

			Expect(loaded.Name).To(Equal(originalConfig.Name))
			Expect(loaded.TargetRole).To(Equal(originalConfig.TargetRole))
			Expect(loaded.TargetAudience).To(Equal(originalConfig.TargetAudience))

			Expect(loaded.EventFilters).To(HaveKey("tags"))
			Expect(loaded.EventFilters).To(HaveKey("companies"))
			Expect(loaded.EventFilters).To(HaveKey("minDate"))
			Expect(loaded.EventFilters).To(HaveKey("maxDate"))
		})

		It("should handle complex nested filters", func() {
			config := fixtures.CVConfigWithFilters("complex-config", map[string]interface{}{
				"tags":       []string{"leadership", "mentoring"},
				"categories": []string{"leadership", "product"},
				"companies":  []string{"Company A", "Company B"},
				"roles":      []string{"staff", "principal"},
				"dateRange":  map[string]string{"start": "2021-01-01", "end": "2024-12-31"},
			})
			config.TargetRole = "em"
			config.TargetAudience = "peer"

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
			config := fixtures.CVConfigWith("../../../etc/passwd", "principal", "hiring_manager")

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

var _ = Describe("MemoryConfigManager", func() {
	var (
		mgr *MemoryConfigManager
		ctx context.Context
	)

	BeforeEach(func() {
		mgr = NewMemoryConfigManager()
		ctx = context.Background()
	})

	Describe("LoadConfig", func() {
		It("should return ErrConfigNotFound for missing config", func() {
			_, err := mgr.LoadConfig(ctx, "nonexistent")
			Expect(err).To(Equal(ErrConfigNotFound))
		})

		It("should load a previously saved config", func() {
			config := fixtures.CVConfigWith("test-memory", "principal", "hiring_manager")
			Expect(mgr.SaveConfig(ctx, config)).To(Succeed())

			loaded, err := mgr.LoadConfig(ctx, "test-memory")
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.Name).To(Equal("test-memory"))
			Expect(loaded.TargetRole).To(Equal("principal"))
		})

		It("should return a copy to prevent external modification", func() {
			config := fixtures.CVConfigWith("copy-test", "principal", "hiring_manager")
			Expect(mgr.SaveConfig(ctx, config)).To(Succeed())

			loaded, err := mgr.LoadConfig(ctx, "copy-test")
			Expect(err).NotTo(HaveOccurred())
			loaded.TargetRole = "modified"

			reloaded, err := mgr.LoadConfig(ctx, "copy-test")
			Expect(err).NotTo(HaveOccurred())
			Expect(reloaded.TargetRole).To(Equal("principal"))
		})
	})

	Describe("SaveConfig", func() {
		It("should reject nil config", func() {
			err := mgr.SaveConfig(ctx, nil)
			Expect(err).To(Equal(ErrInvalidConfigName))
		})

		It("should reject config with empty name", func() {
			config := fixtures.CVConfigWith("", "principal", "hiring_manager")
			err := mgr.SaveConfig(ctx, config)
			Expect(err).To(Equal(ErrInvalidConfigName))
		})

		It("should reject invalid config", func() {
			config := fixtures.CVConfig("test")
			config.TargetRole = ""
			config.TargetAudience = ""
			err := mgr.SaveConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should set timestamps on save", func() {
			config := fixtures.CVConfigWith("ts-test", "principal", "hiring_manager")
			Expect(mgr.SaveConfig(ctx, config)).To(Succeed())

			Expect(config.CreatedAt).NotTo(BeZero())
			Expect(config.UpdatedAt).NotTo(BeZero())
		})

		It("should preserve CreatedAt on update", func() {
			config := fixtures.CVConfigWith("update-test", "principal", "hiring_manager")
			Expect(mgr.SaveConfig(ctx, config)).To(Succeed())
			origCreated := config.CreatedAt

			time.Sleep(5 * time.Millisecond)
			config.TargetAudience = "recruiter"
			Expect(mgr.SaveConfig(ctx, config)).To(Succeed())

			Expect(config.CreatedAt).To(Equal(origCreated))
			Expect(config.UpdatedAt).To(BeTemporally(">", origCreated))
		})
	})

	Describe("DeleteConfig", func() {
		It("should delete an existing config", func() {
			config := fixtures.CVConfigWith("del-test", "principal", "hiring_manager")
			Expect(mgr.SaveConfig(ctx, config)).To(Succeed())

			Expect(mgr.DeleteConfig(ctx, "del-test")).To(Succeed())

			_, err := mgr.LoadConfig(ctx, "del-test")
			Expect(err).To(Equal(ErrConfigNotFound))
		})

		It("should return ErrConfigNotFound for missing config", func() {
			err := mgr.DeleteConfig(ctx, "nonexistent")
			Expect(err).To(Equal(ErrConfigNotFound))
		})
	})

	Describe("ListConfigs", func() {
		It("should return empty list when no configs exist", func() {
			configs, err := mgr.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).To(BeEmpty())
		})

		It("should list all saved configs", func() {
			Expect(mgr.SaveConfig(ctx, fixtures.CVConfigWith("a", "principal", "hiring_manager"))).To(Succeed())
			Expect(mgr.SaveConfig(ctx, fixtures.CVConfigWith("b", "staff", "recruiter"))).To(Succeed())

			configs, err := mgr.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).To(HaveLen(2))
		})

		It("should return copies of configs", func() {
			Expect(mgr.SaveConfig(ctx, fixtures.CVConfigWith("c", "principal", "hiring_manager"))).To(Succeed())

			configs, err := mgr.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			configs[0].TargetRole = "modified"

			original, err := mgr.LoadConfig(ctx, "c")
			Expect(err).NotTo(HaveOccurred())
			Expect(original.TargetRole).To(Equal("principal"))
		})
	})

	Describe("GetConfigPath", func() {
		It("should return memory:// prefixed path", func() {
			path := mgr.GetConfigPath("test-config")
			Expect(path).To(Equal("memory://test-config"))
		})
	})

	Describe("ConfigExists", func() {
		It("should return false for non-existing config", func() {
			exists, err := mgr.ConfigExists(ctx, "nonexistent")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should return true for existing config", func() {
			Expect(mgr.SaveConfig(ctx, fixtures.CVConfigWith("exists-test", "principal", "hiring_manager"))).To(Succeed())

			exists, err := mgr.ConfigExists(ctx, "exists-test")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})
})

var _ = Describe("YAMLConfigManager - GetConfigDirectory and VerifyDirectory", func() {
	var (
		manager *YAMLConfigManager
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-cv-verify-*")
		Expect(err).NotTo(HaveOccurred())

		manager = &YAMLConfigManager{
			configDir: tempDir,
			logger:    logger.DefaultLogger(),
		}
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("GetConfigDirectory", func() {
		It("should return the config directory path", func() {
			Expect(manager.GetConfigDirectory()).To(Equal(tempDir))
		})
	})

	Describe("VerifyDirectory", func() {
		It("should succeed for a valid writable directory", func() {
			Expect(manager.VerifyDirectory()).To(Succeed())
		})

		It("should create directory if it does not exist", func() {
			newDir := filepath.Join(tempDir, "new-subdir")
			manager.configDir = newDir

			Expect(manager.VerifyDirectory()).To(Succeed())

			info, err := os.Stat(newDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())
		})

		It("should return error if path is a file not a directory", func() {
			filePath := filepath.Join(tempDir, "not-a-dir")
			Expect(os.WriteFile(filePath, []byte("data"), 0o600)).To(Succeed())
			manager.configDir = filePath

			err := manager.VerifyDirectory()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a directory"))
		})

		It("should return error for unwritable directory", func() {
			readOnlyDir := filepath.Join(tempDir, "readonly")
			Expect(os.MkdirAll(readOnlyDir, 0o500)).To(Succeed())
			manager.configDir = readOnlyDir

			err := manager.VerifyDirectory()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not writable"))
		})
	})
})

var _ = Describe("YAMLConfigManager additional branches", func() {
	var (
		manager *YAMLConfigManager
		tempDir string
		ctx     context.Context
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-yaml-branch-test-*")
		Expect(err).NotTo(HaveOccurred())

		log := logger.DefaultLogger()
		manager = &YAMLConfigManager{
			configDir: tempDir,
			logger:    log,
		}
		ctx = context.Background()
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("SaveConfig", func() {
		It("should return error for nil config", func() {
			err := manager.SaveConfig(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("nil"))
		})

		It("should return error for cancelled context", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			config := fixtures.CVConfigWith("test", "principal", "hiring_manager")
			err := manager.SaveConfig(cancelCtx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should save and load config successfully", func() {
			config := fixtures.CVConfigWith("save-test", "staff", "recruiter")
			Expect(manager.SaveConfig(ctx, config)).To(Succeed())

			loaded, err := manager.LoadConfig(ctx, "save-test")
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.Name).To(Equal("save-test"))
		})
	})

	Describe("DeleteConfig", func() {
		It("should return error for cancelled context", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			err := manager.DeleteConfig(cancelCtx, "test")
			Expect(err).To(HaveOccurred())
		})

		It("should return error for empty name", func() {
			err := manager.DeleteConfig(ctx, "")
			Expect(err).To(HaveOccurred())
		})

		It("should return ErrConfigNotFound for missing config", func() {
			err := manager.DeleteConfig(ctx, "nonexistent")
			Expect(err).To(Equal(ErrConfigNotFound))
		})

		It("should delete existing config", func() {
			config := fixtures.CVConfigWith("delete-me", "principal", "hiring_manager")
			Expect(manager.SaveConfig(ctx, config)).To(Succeed())

			Expect(manager.DeleteConfig(ctx, "delete-me")).To(Succeed())

			_, err := manager.LoadConfig(ctx, "delete-me")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ListConfigs", func() {
		It("should return error for cancelled context", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			_, err := manager.ListConfigs(cancelCtx)
			Expect(err).To(HaveOccurred())
		})

		It("should return empty list for empty directory", func() {
			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).To(BeEmpty())
		})

		It("should list saved configs", func() {
			config1 := fixtures.CVConfigWith("list-1", "principal", "hiring_manager")
			config2 := fixtures.CVConfigWith("list-2", "staff", "recruiter")
			Expect(manager.SaveConfig(ctx, config1)).To(Succeed())
			Expect(manager.SaveConfig(ctx, config2)).To(Succeed())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).To(HaveLen(2))
		})

		It("should skip non-yaml files", func() {
			config := fixtures.CVConfigWith("yaml-file", "principal", "hiring_manager")
			Expect(manager.SaveConfig(ctx, config)).To(Succeed())

			Expect(os.WriteFile(filepath.Join(tempDir, "not-yaml.txt"), []byte("data"), 0o600)).To(Succeed())

			configs, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).To(HaveLen(1))
		})
	})

	Describe("LoadConfig", func() {
		It("should return error for cancelled context", func() {
			cancelCtx, cancel := context.WithCancel(ctx)
			cancel()
			_, err := manager.LoadConfig(cancelCtx, "test")
			Expect(err).To(HaveOccurred())
		})

		It("should return error for empty name", func() {
			_, err := manager.LoadConfig(ctx, "")
			Expect(err).To(HaveOccurred())
		})

		It("should return ErrConfigNotFound for missing file", func() {
			_, err := manager.LoadConfig(ctx, "nonexistent")
			Expect(err).To(HaveOccurred())
		})

		It("should return error for invalid YAML", func() {
			invalidPath := filepath.Join(tempDir, "invalid.yaml")
			Expect(os.WriteFile(invalidPath, []byte("{{invalid yaml"), 0o600)).To(Succeed())

			_, err := manager.LoadConfig(ctx, "invalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("isPathWithinDirectory", func() {
		It("should return true for path within directory", func() {
			Expect(isPathWithinDirectory("/home/user/configs/test.yaml", "/home/user/configs")).To(BeTrue())
		})

		It("should return false for path outside directory", func() {
			Expect(isPathWithinDirectory("/etc/passwd", "/home/user/configs")).To(BeFalse())
		})

		It("should return false for path traversal", func() {
			Expect(isPathWithinDirectory("/home/user/configs/../../../etc/passwd", "/home/user/configs")).To(BeFalse())
		})
	})
})
