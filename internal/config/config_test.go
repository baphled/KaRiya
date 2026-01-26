//nolint:errcheck // Test file - error handling for test setup is not relevant.
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
			// Use a path with null byte which is invalid on all platforms
			err := config.SaveConfigToPath(cfg, "/path/with\x00null/config.yaml")
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
			err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0600)
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

	Describe("SetConfigPathForTesting", func() {
		AfterEach(func() {
			// Always reset after each test to avoid polluting other tests
			config.ResetConfigPath()
		})

		It("should override GetConfigPath to return the test path", func() {
			testPath := filepath.Join(tempDir, "test-config.yaml")

			config.SetConfigPathForTesting(testPath)

			path, err := config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(testPath))
		})

		It("should make SaveConfig use the overridden path", func() {
			testPath := filepath.Join(tempDir, "test-config.yaml")
			config.SetConfigPathForTesting(testPath)

			cfg := config.DefaultConfig()
			cfg.Profile.Name = "Test Override"

			err := config.SaveConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			// Verify file was saved to test path
			_, err = os.Stat(testPath)
			Expect(err).NotTo(HaveOccurred())

			// Verify content
			loaded, err := config.LoadConfigFromPath(testPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.Profile.Name).To(Equal("Test Override"))
		})

		It("should make LoadConfig use the overridden path", func() {
			testPath := filepath.Join(tempDir, "test-config.yaml")

			// Save a config to test path first
			cfg := config.DefaultConfig()
			cfg.Profile.Name = "Override Test"
			err := config.SaveConfigToPath(cfg, testPath)
			Expect(err).NotTo(HaveOccurred())

			// Set override and load
			config.SetConfigPathForTesting(testPath)

			loaded, err := config.LoadConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(loaded.Profile.Name).To(Equal("Override Test"))
		})
	})

	Describe("ResetConfigPath", func() {
		It("should restore GetConfigPath to return the default path", func() {
			testPath := filepath.Join(tempDir, "test-config.yaml")

			// Get original path
			originalPath, err := config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())

			// Override
			config.SetConfigPathForTesting(testPath)
			overriddenPath, _ := config.GetConfigPath()
			Expect(overriddenPath).To(Equal(testPath))

			// Reset
			config.ResetConfigPath()

			// Should return original path
			restoredPath, err := config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(restoredPath).To(Equal(originalPath))
		})

		It("should be safe to call multiple times", func() {
			// Reset without setting should not panic
			Expect(func() {
				config.ResetConfigPath()
				config.ResetConfigPath()
			}).NotTo(Panic())

			// GetConfigPath should still work
			path, err := config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(ContainSubstring(".kariya"))
		})
	})

	Describe("SwapConfigPathForTesting", func() {
		It("should return the previous path", func() {
			firstPath := filepath.Join(tempDir, "first.yaml")
			secondPath := filepath.Join(tempDir, "second.yaml")

			// Set first path
			config.SetConfigPathForTesting(firstPath)

			// Swap to second path - should return first path
			prev := config.SwapConfigPathForTesting(secondPath)
			Expect(prev).To(Equal(firstPath))

			// Current path should be second
			current, _ := config.GetConfigPath()
			Expect(current).To(Equal(secondPath))
		})

		It("should return empty string when no previous override exists", func() {
			config.ResetConfigPath() // Ensure clean state

			testPath := filepath.Join(tempDir, "test.yaml")
			prev := config.SwapConfigPathForTesting(testPath)
			Expect(prev).To(BeEmpty())

			// Clean up
			config.ResetConfigPath()
		})

		It("should support nested isolation pattern", func() {
			// Simulate BeforeSuite setting a path
			suiteConfigPath := filepath.Join(tempDir, "suite.yaml")
			config.SetConfigPathForTesting(suiteConfigPath)

			// Simulate e2e test swapping its own path
			testConfigPath := filepath.Join(tempDir, "test.yaml")
			prevPath := config.SwapConfigPathForTesting(testConfigPath)
			Expect(prevPath).To(Equal(suiteConfigPath))

			// Current should be test path
			current, _ := config.GetConfigPath()
			Expect(current).To(Equal(testConfigPath))

			// Simulate e2e cleanup restoring previous path
			config.SetConfigPathForTesting(prevPath)

			// Should be back to suite path
			restored, _ := config.GetConfigPath()
			Expect(restored).To(Equal(suiteConfigPath))
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

	Describe("Config Defaults", func() {
		It("should apply all defaults when loading minimal config", func() {
			// Write minimal config with just profile name
			yamlContent := `
profile:
  name: Test User
`
			err := os.WriteFile(configPath, []byte(yamlContent), 0600)
			Expect(err).NotTo(HaveOccurred())

			cfg, err := config.LoadConfigFromPath(configPath)
			Expect(err).NotTo(HaveOccurred())

			// Custom value preserved
			Expect(cfg.Profile.Name).To(Equal("Test User"))

			// System defaults
			Expect(cfg.System.DataDir).To(ContainSubstring(".kariya"))
			Expect(cfg.System.LogLevel).To(Equal("info"))
			Expect(cfg.System.BackupCount).To(Equal(5))

			// Profile defaults
			Expect(cfg.Profile.DefaultRole).To(Equal("senior_ic"))
			Expect(cfg.Profile.DefaultAudience).To(Equal("technical"))

			// CV defaults
			Expect(cfg.CV.DefaultFormat).To(Equal("markdown"))
			Expect(cfg.CV.MaxBullets).To(Equal(50))

			// AudienceBullets defaults
			Expect(cfg.CV.AudienceBullets.Recruiter).To(Equal(4))
			Expect(cfg.CV.AudienceBullets.HiringManager).To(Equal(6))
			Expect(cfg.CV.AudienceBullets.Peer).To(Equal(8))
			Expect(cfg.CV.AudienceBullets.Default).To(Equal(5))

			// Export defaults
			Expect(cfg.Export.DefaultDestination).To(Equal("file"))

			// Display defaults
			Expect(cfg.Display.Theme).To(Equal("dark"))
		})

		It("should apply defaults for missing audience bullet values", func() {
			// Write config with partial audience_bullets (only recruiter set)
			yamlContent := `
cv:
  audience_bullets:
    recruiter: 10
`
			err := os.WriteFile(configPath, []byte(yamlContent), 0600)
			Expect(err).NotTo(HaveOccurred())

			cfg, err := config.LoadConfigFromPath(configPath)
			Expect(err).NotTo(HaveOccurred())

			// Recruiter should be the custom value
			Expect(cfg.CV.AudienceBullets.Recruiter).To(Equal(10))
			// Others should have defaults
			Expect(cfg.CV.AudienceBullets.HiringManager).To(Equal(6))
			Expect(cfg.CV.AudienceBullets.Peer).To(Equal(8))
			Expect(cfg.CV.AudienceBullets.Default).To(Equal(5))
		})

		It("should preserve all custom values when fully specified", func() {
			// Write config with all values set
			yamlContent := `
system:
  log_level: debug
  backup_count: 10
profile:
  default_role: staff
  default_audience: executive
cv:
  default_format: text
  max_bullets: 100
  audience_bullets:
    recruiter: 3
    hiring_manager: 5
    peer: 10
    default: 4
export:
  default_destination: clipboard
display:
  theme: light
`
			err := os.WriteFile(configPath, []byte(yamlContent), 0600)
			Expect(err).NotTo(HaveOccurred())

			cfg, err := config.LoadConfigFromPath(configPath)
			Expect(err).NotTo(HaveOccurred())

			// All custom values preserved
			Expect(cfg.System.LogLevel).To(Equal("debug"))
			Expect(cfg.System.BackupCount).To(Equal(10))
			Expect(cfg.Profile.DefaultRole).To(Equal("staff"))
			Expect(cfg.Profile.DefaultAudience).To(Equal("executive"))
			Expect(cfg.CV.DefaultFormat).To(Equal("text"))
			Expect(cfg.CV.MaxBullets).To(Equal(100))
			Expect(cfg.CV.AudienceBullets.Recruiter).To(Equal(3))
			Expect(cfg.CV.AudienceBullets.HiringManager).To(Equal(5))
			Expect(cfg.CV.AudienceBullets.Peer).To(Equal(10))
			Expect(cfg.CV.AudienceBullets.Default).To(Equal(4))
			Expect(cfg.Export.DefaultDestination).To(Equal("clipboard"))
			Expect(cfg.Display.Theme).To(Equal("light"))
		})
	})

	Describe("ScoringConfig", func() {
		Describe("DefaultConfig", func() {
			It("should include default scoring config with valid weights", func() {
				cfg := config.DefaultConfig()

				Expect(cfg.Scoring).NotTo(BeNil())

				// Weights should sum to 1.0
				weights := cfg.Scoring.Weights
				sum := weights.RoleScore + weights.AudienceScore + weights.MetricScore +
					weights.ImpactScore + weights.Confidence
				Expect(sum).To(BeNumerically("~", 1.0, 0.001))
			})

			It("should have default weight values", func() {
				cfg := config.DefaultConfig()

				Expect(cfg.Scoring.Weights.RoleScore).To(Equal(0.25))
				Expect(cfg.Scoring.Weights.AudienceScore).To(Equal(0.20))
				Expect(cfg.Scoring.Weights.MetricScore).To(Equal(0.20))
				Expect(cfg.Scoring.Weights.ImpactScore).To(Equal(0.20))
				Expect(cfg.Scoring.Weights.Confidence).To(Equal(0.15))
			})

			It("should have default confidence thresholds", func() {
				cfg := config.DefaultConfig()

				Expect(cfg.Scoring.Thresholds.FactDefaultConfidence).To(Equal(0.85))
				Expect(cfg.Scoring.Thresholds.EventDefaultConfidence).To(Equal(0.80))
				Expect(cfg.Scoring.Thresholds.HighConfidence).To(Equal(0.80))
				Expect(cfg.Scoring.Thresholds.HighImpactConfidence).To(Equal(0.85))
			})

			It("should have default role settings", func() {
				cfg := config.DefaultConfig()

				Expect(cfg.Scoring.RoleSettings).To(HaveKey("principal"))
				Expect(cfg.Scoring.RoleSettings).To(HaveKey("staff"))
				Expect(cfg.Scoring.RoleSettings).To(HaveKey("em"))
				Expect(cfg.Scoring.RoleSettings).To(HaveKey("senior_ic"))

				principal := cfg.Scoring.RoleSettings["principal"]
				Expect(principal.MinConfidence).To(Equal(0.80))
				Expect(principal.MaxBulletsPerCompany).To(Equal(4))

				staff := cfg.Scoring.RoleSettings["staff"]
				Expect(staff.MinConfidence).To(Equal(0.75))
				Expect(staff.MaxBulletsPerCompany).To(Equal(5))
			})
		})

		Describe("Loading config with scoring section", func() {
			It("should load custom scoring weights", func() {
				yamlContent := `
scoring:
  weights:
    role_score: 0.30
    audience_score: 0.25
    metric_score: 0.20
    impact_score: 0.15
    confidence: 0.10
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(cfg.Scoring.Weights.RoleScore).To(Equal(0.30))
				Expect(cfg.Scoring.Weights.AudienceScore).To(Equal(0.25))
				Expect(cfg.Scoring.Weights.MetricScore).To(Equal(0.20))
				Expect(cfg.Scoring.Weights.ImpactScore).To(Equal(0.15))
				Expect(cfg.Scoring.Weights.Confidence).To(Equal(0.10))
			})

			It("should load custom thresholds", func() {
				yamlContent := `
scoring:
  thresholds:
    fact_default_confidence: 0.90
    event_default_confidence: 0.85
    high_confidence: 0.85
    high_impact_confidence: 0.90
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(cfg.Scoring.Thresholds.FactDefaultConfidence).To(Equal(0.90))
				Expect(cfg.Scoring.Thresholds.EventDefaultConfidence).To(Equal(0.85))
				Expect(cfg.Scoring.Thresholds.HighConfidence).To(Equal(0.85))
				Expect(cfg.Scoring.Thresholds.HighImpactConfidence).To(Equal(0.90))
			})

			It("should load custom role settings", func() {
				yamlContent := `
scoring:
  role_settings:
    principal:
      min_confidence: 0.85
      max_bullets_per_company: 3
    staff:
      min_confidence: 0.80
      max_bullets_per_company: 4
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(cfg.Scoring.RoleSettings["principal"].MinConfidence).To(Equal(0.85))
				Expect(cfg.Scoring.RoleSettings["principal"].MaxBulletsPerCompany).To(Equal(3))
				Expect(cfg.Scoring.RoleSettings["staff"].MinConfidence).To(Equal(0.80))
				Expect(cfg.Scoring.RoleSettings["staff"].MaxBulletsPerCompany).To(Equal(4))
			})
		})

		Describe("Auto-migration of missing scoring section", func() {
			It("should apply default scoring when section is missing", func() {
				yamlContent := `
profile:
  name: Test User
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				// Defaults should be applied
				Expect(cfg.Scoring.Weights.RoleScore).To(Equal(0.25))
				Expect(cfg.Scoring.Thresholds.FactDefaultConfidence).To(Equal(0.85))
				Expect(cfg.Scoring.RoleSettings).To(HaveLen(4))
			})

			It("should preserve existing config values when adding scoring defaults", func() {
				yamlContent := `
profile:
  name: Existing User
  email: existing@example.com
cv:
  max_bullets: 100
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				// Existing values preserved
				Expect(cfg.Profile.Name).To(Equal("Existing User"))
				Expect(cfg.Profile.Email).To(Equal("existing@example.com"))
				Expect(cfg.CV.MaxBullets).To(Equal(100))

				// Scoring defaults applied
				Expect(cfg.Scoring.Weights.RoleScore).To(Equal(0.25))
			})
		})

		Describe("ValidateWeights", func() {
			It("should return nil for valid weights summing to 1.0", func() {
				cfg := config.DefaultConfig()
				err := cfg.Scoring.ValidateWeights()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return error for weights not summing to 1.0", func() {
				cfg := config.DefaultConfig()
				cfg.Scoring.Weights.RoleScore = 0.50 // This breaks the sum

				err := cfg.Scoring.ValidateWeights()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("must sum to 1.0"))
			})

			It("should accept weights within tolerance", func() {
				cfg := config.DefaultConfig()
				// Slightly adjust to test tolerance
				cfg.Scoring.Weights.RoleScore = 0.2501

				err := cfg.Scoring.ValidateWeights()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
