//nolint:errcheck // Test file - error handling for test setup is not relevant.
package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
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

		It("should use fallback DataDir when home directory is unavailable", func() {
			if runtime.GOOS == "windows" {
				return
			}

			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", "")
			defer os.Setenv("HOME", originalHome)

			cfg := config.DefaultConfig()
			Expect(cfg.System.DataDir).To(Equal(".kariya"))
		})

		It("should fallback to current directory when home dir resolver fails", func() {
			restore := config.SwapHomeDirForTesting(func() (string, error) {
				return "", errors.New("no home directory")
			})
			defer restore()

			cfg := config.DefaultConfig()
			Expect(cfg.System.DataDir).To(Equal(filepath.Join(".", ".kariya")))
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
			err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0o600)
			Expect(err).NotTo(HaveOccurred())

			_, err = config.LoadConfigFromPath(configPath)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to parse config"))
		})
		It("should automatically migrate legacy Name field on load", func() {
			// Create a config with legacy Name field
			cfgData := `
system:
  data_dir: ~/.kariya
  log_level: info
  auto_backup: true
  backup_count: 5
profile:
  name: "Dr. Jane Smith"
  email: test@example.com
  default_role: Software Engineer
  default_audience: Technical
scoring:
  min_score: 0.0
  max_score: 10.0
`
			err := os.WriteFile(configPath, []byte(cfgData), 0o600)
			Expect(err).NotTo(HaveOccurred())

			// Load config - migration should happen automatically
			loaded, err := config.LoadConfigFromPath(configPath)
			Expect(err).NotTo(HaveOccurred())

			// Verify migration occurred
			Expect(loaded.Profile.Prefix).To(Equal("Dr."))
			Expect(loaded.Profile.FirstName).To(Equal("Jane"))
			Expect(loaded.Profile.LastName).To(Equal("Smith"))
			// Legacy Name field should still be present
			Expect(loaded.Profile.Name).To(Equal("Dr. Jane Smith"))
		})

	})

	Describe("GetConfigPath", func() {
		It("should return path in user home directory", func() {
			path, err := config.GetConfigPath()

			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(ContainSubstring(".kariya"))
			Expect(path).To(HaveSuffix("config.yaml"))
		})

		It("should return error when home directory cannot be determined", func() {
			if runtime.GOOS == "windows" {
				return
			}

			config.ResetConfigPath()

			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", "")
			defer os.Setenv("HOME", originalHome)

			_, err := config.GetConfigPath()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get home directory"))
		})

		It("should return error when home dir resolver fails and no override set", func() {
			restore := config.SwapHomeDirForTesting(func() (string, error) {
				return "", errors.New("no home directory")
			})
			defer restore()

			config.ResetConfigPath()

			_, err := config.GetConfigPath()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get home directory"))
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
			// Profile will have auto-migration applied, so check fields individually
			Expect(loaded.Profile.Name).To(Equal(original.Profile.Name))
			Expect(loaded.Profile.Email).To(Equal(original.Profile.Email))
			Expect(loaded.Profile.DefaultRole).To(Equal(original.Profile.DefaultRole))
			Expect(loaded.Profile.DefaultAudience).To(Equal(original.Profile.DefaultAudience))
			// Auto-migration will populate FirstName and LastName from Name
			Expect(loaded.Profile.FirstName).To(Equal("John"))
			Expect(loaded.Profile.LastName).To(Equal("Doe"))
			Expect(loaded.CV).To(Equal(original.CV))
			Expect(loaded.Export).To(Equal(original.Export))
			Expect(loaded.Display).To(Equal(original.Display))
		})

		Describe("Config Defaults", func() {
			It("should apply all defaults when loading minimal config", func() {
				// Write minimal config with just profile name
				yamlContent := `
profile:
  name: Test User
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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

		Describe("Profile fields - new and legacy compatibility", func() {
			It("should load new profile fields from YAML", func() {
				yamlContent := "profile:\n  first_name: Alice\n  last_name: Smith\n  prefix: Dr.\n  phone: \"+123456789\"\n  linkedin: alice-smith\n  country: Wonderland\n  email: alice@example.com\n"
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(cfg.Profile.FirstName).To(Equal("Alice"))
				Expect(cfg.Profile.LastName).To(Equal("Smith"))
				Expect(cfg.Profile.Prefix).To(Equal("Dr."))
				Expect(cfg.Profile.Phone).To(Equal("+123456789"))
				Expect(cfg.Profile.LinkedIn).To(Equal("alice-smith"))
				Expect(cfg.Profile.Country).To(Equal("Wonderland"))
				Expect(cfg.Profile.Email).To(Equal("alice@example.com"))
			})

			It("should auto-migrate legacy `name` field when new fields absent", func() {
				yamlContent := "profile:\n  name: Legacy User\n  email: legacy@example.com\n"
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				// Legacy name still present
				Expect(cfg.Profile.Name).To(Equal("Legacy User"))
				Expect(cfg.Profile.Email).To(Equal("legacy@example.com"))

				// Auto-migration should populate FirstName and LastName
				Expect(cfg.Profile.FirstName).To(Equal("Legacy"))
				Expect(cfg.Profile.LastName).To(Equal("User"))
				Expect(cfg.Profile.Prefix).To(BeEmpty())
				Expect(cfg.Profile.Phone).To(BeEmpty())
				Expect(cfg.Profile.LinkedIn).To(BeEmpty())
				Expect(cfg.Profile.Country).To(BeEmpty())
			})

			It("should round-trip profile fields through save/load", func() {
				cfg := config.DefaultConfig()
				cfg.Profile.FirstName = "Ada"
				cfg.Profile.LastName = "Lovelace"
				cfg.Profile.Prefix = "Dr."
				cfg.Profile.Phone = "+4412345678"
				cfg.Profile.LinkedIn = "ada-lovelace"
				cfg.Profile.Country = "UK"

				err := config.SaveConfigToPath(cfg, configPath)
				Expect(err).NotTo(HaveOccurred())

				loaded, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(loaded.Profile.FirstName).To(Equal("Ada"))
				Expect(loaded.Profile.LastName).To(Equal("Lovelace"))
				Expect(loaded.Profile.Prefix).To(Equal("Dr."))
				Expect(loaded.Profile.Phone).To(Equal("+4412345678"))
				Expect(loaded.Profile.LinkedIn).To(Equal("ada-lovelace"))
				Expect(loaded.Profile.Country).To(Equal("UK"))
			})

			It("should load SummaryHeading from YAML", func() {
				yamlContent := "profile:\n  summary_heading: \"**{{.Title}} | {{.Years}}+ Years**\"\n"
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(cfg.Profile.SummaryHeading).To(Equal("**{{.Title}} | {{.Years}}+ Years**"))
			})

			It("should unmarshal with empty string when summary_heading absent", func() {
				yamlContent := "profile:\n  first_name: Test\n"
				err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
				Expect(err).NotTo(HaveOccurred())

				cfg, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(cfg.Profile.SummaryHeading).To(BeEmpty())
			})

			It("should round-trip SummaryHeading through save/load", func() {
				cfg := config.DefaultConfig()
				cfg.Profile.SummaryHeading = "**{{.Title}} | {{.Specialty}}**"

				err := config.SaveConfigToPath(cfg, configPath)
				Expect(err).NotTo(HaveOccurred())

				loaded, err := config.LoadConfigFromPath(configPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(loaded.Profile.SummaryHeading).To(Equal("**{{.Title}} | {{.Specialty}}**"))
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
					err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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
					err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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
					err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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
					err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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
					err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
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

			Describe("MigrateProfileConfig", func() {
				DescribeTable("migrates legacy profile name into new fields",
					func(name string, wantMigrated bool, wantPrefix, wantFirst, wantLast string) {
						cfg := config.DefaultConfig()
						cfg.Profile.Name = name
						cfg.Profile.FirstName = ""
						cfg.Profile.LastName = ""

						result := config.MigrateProfileConfig(cfg)

						Expect(result).To(Equal(wantMigrated))
						Expect(cfg.Profile.Prefix).To(Equal(wantPrefix))
						Expect(cfg.Profile.FirstName).To(Equal(wantFirst))
						Expect(cfg.Profile.LastName).To(Equal(wantLast))
					},
					Entry("standard name", "Yomi Colledge", true, "", "Yomi", "Colledge"),
					Entry("prefixed name", "Dr. Jane Smith", true, "Dr.", "Jane", "Smith"),
					Entry("multi-word last name", "John Paul Jones", true, "", "John", "Paul Jones"),
					Entry("single name", "Madonna", true, "", "Madonna", ""),
					Entry("empty name", "", false, "", "", ""),
					Entry("hyphenated name", "Mary-Jane Watson-Parker", true, "", "Mary-Jane", "Watson-Parker"),
				)

				It("should migrate when name is populated and new fields are empty", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Ada Lovelace"
					cfg.Profile.FirstName = ""
					cfg.Profile.LastName = ""

					result := config.MigrateProfileConfig(cfg)

					Expect(result).To(BeTrue())
					Expect(cfg.Profile.FirstName).To(Equal("Ada"))
					Expect(cfg.Profile.LastName).To(Equal("Lovelace"))
				})

				It("should return false when FirstName already set", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Yomi Colledge"
					cfg.Profile.FirstName = "Yomi"
					cfg.Profile.LastName = "Colledge"

					result := config.MigrateProfileConfig(cfg)

					Expect(result).To(BeFalse())
				})

				It("should be idempotent when called multiple times", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Yomi Colledge"
					cfg.Profile.FirstName = ""
					cfg.Profile.LastName = ""

					first := config.MigrateProfileConfig(cfg)
					second := config.MigrateProfileConfig(cfg)

					Expect(first).To(BeTrue())
					Expect(second).To(BeFalse())
					Expect(cfg.Profile.FirstName).To(Equal("Yomi"))
					Expect(cfg.Profile.LastName).To(Equal("Colledge"))
				})

				It("should detect Prof. prefix", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Prof. John Doe"
					cfg.Profile.FirstName = ""
					cfg.Profile.LastName = ""

					result := config.MigrateProfileConfig(cfg)

					Expect(result).To(BeTrue())
					Expect(cfg.Profile.Prefix).To(Equal("Prof."))
					Expect(cfg.Profile.FirstName).To(Equal("John"))
					Expect(cfg.Profile.LastName).To(Equal("Doe"))
				})

				It("should detect Mr. prefix", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Mr. John Doe"
					cfg.Profile.FirstName = ""
					cfg.Profile.LastName = ""

					result := config.MigrateProfileConfig(cfg)

					Expect(result).To(BeTrue())
					Expect(cfg.Profile.Prefix).To(Equal("Mr."))
					Expect(cfg.Profile.FirstName).To(Equal("John"))
					Expect(cfg.Profile.LastName).To(Equal("Doe"))
				})

				It("should detect Mrs. prefix", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Mrs. Jane Doe"
					cfg.Profile.FirstName = ""
					cfg.Profile.LastName = ""

					result := config.MigrateProfileConfig(cfg)

					Expect(result).To(BeTrue())
					Expect(cfg.Profile.Prefix).To(Equal("Mrs."))
					Expect(cfg.Profile.FirstName).To(Equal("Jane"))
					Expect(cfg.Profile.LastName).To(Equal("Doe"))
				})

				It("should detect Ms. prefix", func() {
					cfg := config.DefaultConfig()
					cfg.Profile.Name = "Ms. Jane Doe"
					cfg.Profile.FirstName = ""
					cfg.Profile.LastName = ""

					result := config.MigrateProfileConfig(cfg)

					Expect(result).To(BeTrue())
					Expect(cfg.Profile.Prefix).To(Equal("Ms."))
					Expect(cfg.Profile.FirstName).To(Equal("Jane"))
					Expect(cfg.Profile.LastName).To(Equal("Doe"))
				})
			})
		})
	})

	Describe("requireTestIsolation", func() {
		It("should panic when LoadConfig is called without config path override", func() {
			config.ResetConfigPath()

			Expect(func() {
				config.LoadConfig()
			}).To(Panic())
		})

		It("should panic when SaveConfig is called without config path override", func() {
			config.ResetConfigPath()

			Expect(func() {
				config.SaveConfig(config.DefaultConfig())
			}).To(Panic())
		})
	})

	Describe("LoadConfigFromPath error handling", func() {
		It("should return error when file exists but is not readable", func() {
			if runtime.GOOS == "windows" {
				return
			}

			err := os.WriteFile(configPath, []byte("system:\n  log_level: info\n"), 0o600)
			Expect(err).NotTo(HaveOccurred())

			err = os.Chmod(configPath, 0o000)
			Expect(err).NotTo(HaveOccurred())

			_, err = config.LoadConfigFromPath(configPath)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to read config"))
		})
	})

	Describe("SaveConfigToPath error handling", func() {
		It("should return error when file path is a directory", func() {
			dirAsFile := filepath.Join(tempDir, "is-a-dir")
			err := os.Mkdir(dirAsFile, 0o750)
			Expect(err).NotTo(HaveOccurred())

			cfg := config.DefaultConfig()
			err = config.SaveConfigToPath(cfg, dirAsFile)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to write config"))
		})

		It("should return error when YAML marshalling fails", func() {
			restore := config.SwapYamlMarshalForTesting(func(_ interface{}) ([]byte, error) {
				return nil, errors.New("marshal failure")
			})
			defer restore()

			cfg := config.DefaultConfig()
			err := config.SaveConfigToPath(cfg, filepath.Join(tempDir, "marshal-fail.yaml"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to marshal config"))
		})
	})

	Describe("MigrateProfileConfig with whitespace-only remaining name", func() {
		It("should handle name that is only a prefix followed by spaces", func() {
			cfg := config.DefaultConfig()
			cfg.Profile.Name = "Dr.  "
			cfg.Profile.FirstName = ""
			cfg.Profile.LastName = ""

			result := config.MigrateProfileConfig(cfg)

			Expect(result).To(BeTrue())
			Expect(cfg.Profile.Prefix).To(Equal("Dr."))
			Expect(cfg.Profile.FirstName).To(BeEmpty())
			Expect(cfg.Profile.LastName).To(BeEmpty())
		})

		It("should handle name that is only whitespace", func() {
			cfg := config.DefaultConfig()
			cfg.Profile.Name = "   "
			cfg.Profile.FirstName = ""
			cfg.Profile.LastName = ""

			result := config.MigrateProfileConfig(cfg)

			Expect(result).To(BeTrue())
			Expect(cfg.Profile.FirstName).To(BeEmpty())
		})
	})
})
