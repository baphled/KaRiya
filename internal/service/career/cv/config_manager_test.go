//nolint:errcheck // Test file - error handling for test setup is not relevant.
package cv

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// expectWindowsPermissionError handles the platform-specific behaviour of permission errors.
// Windows does not enforce Unix-style permissions, so mkdir/write may succeed or fail
// depending on ACLs. This helper verifies the expected behaviour on each platform.
func expectWindowsPermissionError(err error, errorSubstring string, dirPath string) {
	if runtime.GOOS == "windows" {
		if err == nil {
			_, statErr := os.Stat(dirPath)
			Expect(statErr).NotTo(HaveOccurred())
		} else {
			Expect(err.Error()).To(ContainSubstring(errorSubstring))
		}
	} else {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(errorSubstring))
	}
}

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

		It("should handle context cancellation", func() {
			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			err := manager.SaveConfig(cancelCtx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should set CreatedAt when zero", func() {
			config := fixtures.CVConfigWith("zero-created", "principal", "hiring_manager")
			config.CreatedAt = time.Time{}

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "zero-created")
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.CreatedAt).NotTo(BeZero())
		})

		It("should preserve existing CreatedAt on update", func() {
			config := fixtures.CVConfigWith("preserve-ts", "principal", "hiring_manager")
			originalCreatedAt := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
			config.CreatedAt = originalCreatedAt

			err := manager.SaveConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := manager.LoadConfig(ctx, "preserve-ts")
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.CreatedAt).To(BeTemporally("~", originalCreatedAt, time.Second))
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

		It("should return error for invalid YAML content", func() {
			invalidPath := filepath.Join(tempDir, "bad-yaml.yaml")
			Expect(os.WriteFile(invalidPath, []byte("{{{invalid yaml: [[["), 0o600)).To(Succeed())

			_, err := manager.LoadConfig(ctx, "bad-yaml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unmarshal"))
		})

		It("should return error when config fails validation", func() {
			validYAML := "name: invalid-role-config\ntarget_role: not-a-valid-role\ntarget_audience: hiring_manager\n"
			configPath := filepath.Join(tempDir, "invalid-role-config.yaml")
			Expect(os.WriteFile(configPath, []byte(validYAML), 0o600)).To(Succeed())

			_, err := manager.LoadConfig(ctx, "invalid-role-config")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("validation failed"))
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

		It("should handle context cancellation", func() {
			config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")
			Expect(manager.SaveConfig(context.Background(), config)).To(Succeed())

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			err := manager.DeleteConfig(cancelCtx, "test-config")
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

		It("should return empty list for non-existent directory", func() {
			manager.configDir = filepath.Join(tempDir, "does-not-exist")

			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded).To(BeEmpty())
		})

		It("should skip subdirectories", func() {
			config := fixtures.CVConfigWith("config1", "principal", "hiring_manager")
			Expect(manager.SaveConfig(ctx, config)).To(Succeed())

			subDir := filepath.Join(tempDir, "subdir")
			Expect(os.Mkdir(subDir, 0o750)).To(Succeed())

			loaded, err := manager.ListConfigs(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded).To(HaveLen(1))
			Expect(loaded[0].Name).To(Equal("config1"))
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

	Describe("GetConfigDirectory", func() {
		It("should return the configured directory", func() {
			Expect(manager.GetConfigDirectory()).To(Equal(tempDir))
		})

		It("should return the exact configDir field value", func() {
			customDir := filepath.Join(tempDir, "custom", "path")
			manager.configDir = customDir
			Expect(manager.GetConfigDirectory()).To(Equal(customDir))
		})
	})

	Describe("VerifyDirectory", func() {
		It("should succeed for existing writable directory", func() {
			err := manager.VerifyDirectory()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create directory when it does not exist", func() {
			newDir := filepath.Join(tempDir, "new-subdir", "deep")
			manager.configDir = newDir

			err := manager.VerifyDirectory()
			Expect(err).NotTo(HaveOccurred())

			info, statErr := os.Stat(newDir)
			Expect(statErr).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())
		})

		It("should return error when path is a file not directory", func() {
			filePath := filepath.Join(tempDir, "not-a-dir")
			Expect(os.WriteFile(filePath, []byte("content"), 0o600)).To(Succeed())

			manager.configDir = filePath

			err := manager.VerifyDirectory()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a directory"))
		})

		It("should return error when stat fails with non-NotExist error", func() {
			if runtime.GOOS != "linux" {
				return
			}

			manager.configDir = "/dev/null/impossible"

			err := manager.VerifyDirectory()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to stat"))
		})

		It("should return error when directory creation fails", func() {
			if os.Getuid() == 0 {
				return
			}

			parent := filepath.Join(tempDir, "readonly-parent")
			Expect(os.Mkdir(parent, 0o555)).To(Succeed())
			defer os.Chmod(parent, 0o755)

			manager.configDir = filepath.Join(parent, "child")

			err := manager.VerifyDirectory()
			expectWindowsPermissionError(err, "failed to create", manager.configDir)
		})

		It("should return error when directory is not writable", func() {
			if os.Getuid() == 0 {
				return
			}

			readOnlyDir := filepath.Join(tempDir, "not-writable")
			Expect(os.Mkdir(readOnlyDir, 0o555)).To(Succeed())
			defer os.Chmod(readOnlyDir, 0o755)

			manager.configDir = readOnlyDir

			err := manager.VerifyDirectory()
			expectWindowsPermissionError(err, "not writable", manager.configDir)
		})

		It("should work with nil logger on existing directory", func() {
			manager.logger = nil

			err := manager.VerifyDirectory()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create directory with nil logger when missing", func() {
			newDir := filepath.Join(tempDir, "nil-logger-dir")
			manager.configDir = newDir
			manager.logger = nil

			err := manager.VerifyDirectory()
			Expect(err).NotTo(HaveOccurred())

			_, statErr := os.Stat(newDir)
			Expect(statErr).NotTo(HaveOccurred())
		})

		It("should return error with nil logger when path is a file", func() {
			filePath := filepath.Join(tempDir, "file-not-dir")
			Expect(os.WriteFile(filePath, []byte("data"), 0o600)).To(Succeed())

			manager.configDir = filePath
			manager.logger = nil

			err := manager.VerifyDirectory()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a directory"))
		})

		It("should return stat error with nil logger", func() {
			if runtime.GOOS != "linux" {
				return
			}

			manager.configDir = "/dev/null/impossible"
			manager.logger = nil

			err := manager.VerifyDirectory()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to stat"))
		})

		It("should verify write permissions", func() {
			err := manager.VerifyDirectory()
			Expect(err).NotTo(HaveOccurred())

			// Verify the write test file was cleaned up
			testFile := filepath.Join(tempDir, ".write-test")
			_, statErr := os.Stat(testFile)
			Expect(os.IsNotExist(statErr)).To(BeTrue())
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

	It("should have configDir ending with .kariya/cv_configs", func() {
		log := logger.DefaultLogger()
		manager, err := NewYAMLConfigManager(log)
		Expect(err).NotTo(HaveOccurred())
		Expect(manager.configDir).To(HaveSuffix(filepath.Join(".kariya", "cv_configs")))
	})

	It("should set the logger field", func() {
		log := logger.DefaultLogger()
		manager, err := NewYAMLConfigManager(log)
		Expect(err).NotTo(HaveOccurred())
		Expect(manager.logger).To(Equal(log))
	})
})

var _ = Describe("isPathWithinDirectory", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-path-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	It("should return true for path within directory", func() {
		filePath := filepath.Join(tempDir, "config.yaml")
		Expect(isPathWithinDirectory(filePath, tempDir)).To(BeTrue())
	})

	It("should return true for nested path within directory", func() {
		filePath := filepath.Join(tempDir, "sub", "dir", "config.yaml")
		Expect(isPathWithinDirectory(filePath, tempDir)).To(BeTrue())
	})

	It("should return false for path outside directory", func() {
		Expect(isPathWithinDirectory("/etc/passwd", tempDir)).To(BeFalse())
	})

	It("should return false for sibling directory with same prefix", func() {
		siblingDir := tempDir + "-evil"
		filePath := filepath.Join(siblingDir, "config.yaml")
		Expect(isPathWithinDirectory(filePath, tempDir)).To(BeFalse())
	})

	It("should return false for parent directory traversal", func() {
		filePath := filepath.Join(tempDir, "..", "outside.yaml")
		Expect(isPathWithinDirectory(filePath, tempDir)).To(BeFalse())
	})

	It("should handle directory with trailing separator", func() {
		filePath := filepath.Join(tempDir, "config.yaml")
		dirWithSep := tempDir + string(filepath.Separator)
		Expect(isPathWithinDirectory(filePath, dirWithSep)).To(BeTrue())
	})
})

var _ = Describe("MemoryConfigManager", func() {
	var (
		manager *MemoryConfigManager
		ctx     context.Context
	)

	BeforeEach(func() {
		manager = NewMemoryConfigManager()
		ctx = context.Background()
	})

	Describe("LoadConfig", func() {
		Context("when the config exists", func() {
			It("returns the config", func() {
				config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")
				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				loaded, err := manager.LoadConfig(ctx, "test-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(loaded.Name).To(Equal("test-config"))
				Expect(loaded.TargetRole).To(Equal("principal"))
				Expect(loaded.TargetAudience).To(Equal("hiring_manager"))
			})

			It("returns a copy to prevent external modification", func() {
				config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")
				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				loaded, err := manager.LoadConfig(ctx, "test-config")
				Expect(err).NotTo(HaveOccurred())

				loaded.TargetRole = "em"

				reloaded, err := manager.LoadConfig(ctx, "test-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(reloaded.TargetRole).To(Equal("principal"))
			})
		})

		Context("when the config does not exist", func() {
			It("returns ErrConfigNotFound", func() {
				_, err := manager.LoadConfig(ctx, "nonexistent")
				Expect(err).To(MatchError(ErrConfigNotFound))
			})
		})
	})

	Describe("SaveConfig", func() {
		Context("with a valid config", func() {
			It("stores the config successfully", func() {
				config := fixtures.CVConfigWith("test-config", "principal", "hiring_manager")

				err := manager.SaveConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())

				loaded, err := manager.LoadConfig(ctx, "test-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(loaded.Name).To(Equal("test-config"))
			})

			It("sets CreatedAt when zero", func() {
				config := fixtures.CVConfigWith("save-ts", "staff", "recruiter")
				config.CreatedAt = time.Time{}

				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				loaded, err := manager.LoadConfig(ctx, "save-ts")
				Expect(err).NotTo(HaveOccurred())
				Expect(loaded.CreatedAt).NotTo(BeZero())
			})

			It("preserves existing CreatedAt on update", func() {
				config := fixtures.CVConfigWith("update-ts", "staff", "recruiter")
				originalCreatedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				config.CreatedAt = originalCreatedAt

				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				loaded, err := manager.LoadConfig(ctx, "update-ts")
				Expect(err).NotTo(HaveOccurred())
				Expect(loaded.CreatedAt).To(Equal(originalCreatedAt))
			})

			It("sets UpdatedAt on every save", func() {
				config := fixtures.CVConfigWith("upd-at", "em", "peer")

				before := time.Now()
				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				loaded, err := manager.LoadConfig(ctx, "upd-at")
				Expect(err).NotTo(HaveOccurred())
				Expect(loaded.UpdatedAt).To(BeTemporally(">=", before))
			})
		})

		Context("with nil config", func() {
			It("returns ErrInvalidConfigName", func() {
				err := manager.SaveConfig(ctx, nil)
				Expect(err).To(MatchError(ErrInvalidConfigName))
			})
		})

		Context("with empty name", func() {
			It("returns ErrInvalidConfigName", func() {
				config := fixtures.CVConfig("")
				err := manager.SaveConfig(ctx, config)
				Expect(err).To(MatchError(ErrInvalidConfigName))
			})
		})

		Context("with an invalid config", func() {
			It("returns a validation error", func() {
				config := fixtures.CVConfig("invalid-config")
				config.TargetRole = "not-a-valid-role"
				err := manager.SaveConfig(ctx, config)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("DeleteConfig", func() {
		Context("when the config exists", func() {
			It("removes the config", func() {
				config := fixtures.CVConfigWith("to-delete", "principal", "hiring_manager")
				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				err := manager.DeleteConfig(ctx, "to-delete")
				Expect(err).NotTo(HaveOccurred())

				_, err = manager.LoadConfig(ctx, "to-delete")
				Expect(err).To(MatchError(ErrConfigNotFound))
			})
		})

		Context("when the config does not exist", func() {
			It("returns ErrConfigNotFound", func() {
				err := manager.DeleteConfig(ctx, "nonexistent")
				Expect(err).To(MatchError(ErrConfigNotFound))
			})
		})
	})

	Describe("ListConfigs", func() {
		Context("when configs exist", func() {
			It("returns all stored configs", func() {
				configs := []*career.CVConfig{
					fixtures.CVConfigWith("config-a", "principal", "hiring_manager"),
					fixtures.CVConfigWith("config-b", "staff", "recruiter"),
					fixtures.CVConfigWith("config-c", "em", "peer"),
				}
				for _, c := range configs {
					Expect(manager.SaveConfig(ctx, c)).To(Succeed())
				}

				listed, err := manager.ListConfigs(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(listed).To(HaveLen(3))
			})

			It("returns copies to prevent external modification", func() {
				config := fixtures.CVConfigWith("list-copy", "principal", "hiring_manager")
				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				listed, err := manager.ListConfigs(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(listed).To(HaveLen(1))

				listed[0].TargetRole = "em"

				reloaded, err := manager.LoadConfig(ctx, "list-copy")
				Expect(err).NotTo(HaveOccurred())
				Expect(reloaded.TargetRole).To(Equal("principal"))
			})
		})

		Context("when no configs exist", func() {
			It("returns an empty slice", func() {
				listed, err := manager.ListConfigs(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(listed).To(BeEmpty())
			})
		})
	})

	Describe("GetConfigPath", func() {
		It("returns a memory:// URI for the given name", func() {
			path := manager.GetConfigPath("my-config")
			Expect(path).To(Equal("memory://my-config"))
		})
	})

	Describe("ConfigExists", func() {
		Context("when the config exists", func() {
			It("returns true", func() {
				config := fixtures.CVConfigWith("exists-check", "staff", "recruiter")
				Expect(manager.SaveConfig(ctx, config)).To(Succeed())

				exists, err := manager.ConfigExists(ctx, "exists-check")
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
			})
		})

		Context("when the config does not exist", func() {
			It("returns false", func() {
				exists, err := manager.ConfigExists(ctx, "nonexistent")
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())
			})
		})
	})
})
