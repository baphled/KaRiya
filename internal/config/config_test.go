package config_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	var (
		tempDir    string
		configPath string
	)

	BeforeEach(func() {
		// Create temporary directory for test configs
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-config-test-*")
		Expect(err).NotTo(HaveOccurred())

		configPath = filepath.Join(tempDir, "config.yaml")
	})

	AfterEach(func() {
		// Clean up temp directory
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("DefaultConfig", func() {
		It("should return a config with default values", func() {
			cfg := config.DefaultConfig()

			Expect(cfg).NotTo(BeNil())
			Expect(cfg.System.LogLevel).To(Equal("info"))
			Expect(cfg.System.AutoBackup).To(BeTrue())
			Expect(cfg.System.BackupCount).To(Equal(5))
			Expect(cfg.Profile.DefaultRole).To(Equal("senior_ic"))
			Expect(cfg.Profile.DefaultAudience).To(Equal("technical"))
		})

		It("should set DataDir to user home directory", func() {
			cfg := config.DefaultConfig()

			Expect(cfg.System.DataDir).NotTo(BeEmpty())
			Expect(cfg.System.DataDir).To(ContainSubstring(".kariya"))
		})
	})

	Describe("SaveConfig", func() {
		It("should save config to YAML file", func() {
			cfg := config.DefaultConfig()
			cfg.System.LogLevel = "debug"
			cfg.Profile.Name = "Test User"

			err := config.SaveConfigToPath(cfg, configPath)
			Expect(err).NotTo(HaveOccurred())

			// Verify file exists
			_, err = os.Stat(configPath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create directory if it doesn't exist", func() {
			nestedPath := filepath.Join(tempDir, "nested", "dir", "config.yaml")
			cfg := config.DefaultConfig()

			err := config.SaveConfigToPath(cfg, nestedPath)
			Expect(err).NotTo(HaveOccurred())

			// Verify file exists
			_, err = os.Stat(nestedPath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle invalid path", func() {
			cfg := config.DefaultConfig()
			err := config.SaveConfigToPath(cfg, "/invalid/nonexistent/path/config.yaml")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("LoadConfig", func() {
		It("should load config from YAML file", func() {
			// Save a config first
			original := config.DefaultConfig()
			original.System.LogLevel = "debug"
			original.Profile.Name = "Test User"
			original.Profile.Email = "test@example.com"

			err := config.SaveConfigToPath(original, configPath)
			Expect(err).NotTo(HaveOccurred())

			// Load it back
			loaded, err := config.LoadConfigFromPath(configPath)
			Expect(err).NotTo(HaveOccurred())

			Expect(loaded.System.LogLevel).To(Equal("debug"))
			Expect(loaded.Profile.Name).To(Equal("Test User"))
			Expect(loaded.Profile.Email).To(Equal("test@example.com"))
		})

		It("should return default config if file doesn't exist", func() {
			nonexistentPath := filepath.Join(tempDir, "nonexistent.yaml")
			cfg, err := config.LoadConfigFromPath(nonexistentPath)

			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.System.LogLevel).To(Equal("info")) // Default value
		})

		It("should return error for invalid YAML", func() {
			// Write invalid YAML
			err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
			Expect(err).NotTo(HaveOccurred())

			_, err = config.LoadConfigFromPath(configPath)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to parse config"))
		})
	})

	Describe("GetConfigPath", func() {
		It("should return path in user home directory", func() {
			path, err := config.GetConfigPath()

			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(ContainSubstring(".kariya"))
			Expect(path).To(HaveSuffix("config.yaml"))
		})
	})

	Describe("Config Struct", func() {
		It("should marshal to YAML correctly", func() {
			cfg := config.DefaultConfig()
			cfg.System.LogLevel = "debug"
			cfg.Profile.Name = "Test"

			err := config.SaveConfigToPath(cfg, configPath)
			Expect(err).NotTo(HaveOccurred())

			// Read raw YAML
			data, err := os.ReadFile(configPath)
			Expect(err).NotTo(HaveOccurred())

			yamlStr := string(data)
			Expect(yamlStr).To(ContainSubstring("log_level: debug"))
			Expect(yamlStr).To(ContainSubstring("name: Test"))
		})

		It("should preserve all fields during save/load cycle", func() {
			original := config.DefaultConfig()
			original.System.LogLevel = "warn"
			original.System.DataDir = "/custom/path"
			original.System.AutoBackup = false
			original.System.BackupCount = 10
			original.Profile.Name = "John Doe"
			original.Profile.Email = "john@example.com"
			original.Profile.DefaultRole = "staff_ic"
			original.Profile.DefaultAudience = "executive"
			original.CV.DefaultFormat = "text"
			original.CV.MaxBullets = 75
			original.Export.DefaultDestination = "clipboard"
			original.Export.AutoOpen = true
			original.Display.Theme = "light"
			original.Display.Animations = false

			err := config.SaveConfigToPath(original, configPath)
			Expect(err).NotTo(HaveOccurred())

			loaded, err := config.LoadConfigFromPath(configPath)
			Expect(err).NotTo(HaveOccurred())

			// Verify all fields preserved
			Expect(loaded.System).To(Equal(original.System))
			Expect(loaded.Profile).To(Equal(original.Profile))
			Expect(loaded.CV).To(Equal(original.CV))
			Expect(loaded.Export).To(Equal(original.Export))
			Expect(loaded.Display).To(Equal(original.Display))
		})
	})
})
