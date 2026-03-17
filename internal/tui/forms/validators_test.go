package forms_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/tui/forms"
)

var _ = Describe("Validators", func() {
	Describe("Required", func() {
		It("should pass for non-empty strings", func() {
			err := forms.Required("hello")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for empty strings", func() {
			err := forms.Required("")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for whitespace-only strings", func() {
			err := forms.Required("   ")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MinLength", func() {
		It("should pass when string meets minimum", func() {
			validator := forms.MinLength(5)
			err := validator("hello")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when string is too short", func() {
			validator := forms.MinLength(10)
			err := validator("hello")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least 10 characters"))
		})
	})

	Describe("MaxLength", func() {
		It("should pass when string is under maximum", func() {
			validator := forms.MaxLength(10)
			err := validator("hello")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when string is too long", func() {
			validator := forms.MaxLength(3)
			err := validator("hello")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at most 3 characters"))
		})
	})

	Describe("LengthRange", func() {
		It("should pass when string is in range", func() {
			validator := forms.LengthRange(3, 10)
			err := validator("hello")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when string is too short", func() {
			validator := forms.LengthRange(10, 20)
			err := validator("hi")
			Expect(err).To(HaveOccurred())
		})

		It("should fail when string is too long", func() {
			validator := forms.LengthRange(1, 3)
			err := validator("hello")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DateFormat", func() {
		It("should pass for valid dates", func() {
			err := forms.DateFormat("2024-01-07")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for empty strings", func() {
			err := forms.DateFormat("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for invalid format", func() {
			err := forms.DateFormat("01/07/2024")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for invalid dates", func() {
			err := forms.DateFormat("2024-13-99")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DateFormatRequired", func() {
		It("should pass for valid dates", func() {
			err := forms.DateFormatRequired("2024-01-07")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for empty strings", func() {
			err := forms.DateFormatRequired("")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Email", func() {
		It("should pass for valid emails", func() {
			err := forms.Email("test@example.com")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for empty strings", func() {
			err := forms.Email("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for invalid emails", func() {
			err := forms.Email("not-an-email")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("URL", func() {
		It("should pass for valid URLs", func() {
			err := forms.URL("https://example.com")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for HTTP URLs", func() {
			err := forms.URL("http://example.com/path")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for empty strings", func() {
			err := forms.URL("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for invalid URLs", func() {
			err := forms.URL("not a url")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("AlphaNumeric", func() {
		It("should pass for alphanumeric strings", func() {
			err := forms.AlphaNumeric("Hello123")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for empty strings", func() {
			err := forms.AlphaNumeric("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for strings with special characters", func() {
			err := forms.AlphaNumeric("hello-world")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("NoSpecialChars", func() {
		It("should pass for strings without special chars", func() {
			err := forms.NoSpecialChars("Hello World 123")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for hyphens and underscores", func() {
			err := forms.NoSpecialChars("hello-world_123")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for special characters", func() {
			err := forms.NoSpecialChars("hello@world")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Compose", func() {
		It("should pass when all validators pass", func() {
			validator := forms.Compose(
				forms.Required,
				forms.MinLength(5),
				forms.MaxLength(10),
			)
			err := validator("hello")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail on first error", func() {
			validator := forms.Compose(
				forms.Required,
				forms.MinLength(10),
			)
			err := validator("hi")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least 10 characters"))
		})

		It("should fail on empty string when Required is first", func() {
			validator := forms.Compose(
				forms.Required,
				forms.MinLength(5),
			)
			err := validator("")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(forms.ErrRequired))
		})
	})

	Describe("Domain-specific validators", func() {
		Describe("EventText", func() {
			It("should pass for valid event text", func() {
				err := forms.EventText("This is a valid career event description")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail for empty strings", func() {
				err := forms.EventText("")
				Expect(err).To(HaveOccurred())
			})

			It("should fail for too short text", func() {
				err := forms.EventText("Too short")
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("EventTextOptional", func() {
			It("should pass for valid event text", func() {
				err := forms.EventTextOptional("This is a valid career event description")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass for empty strings", func() {
				err := forms.EventTextOptional("")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail for too short non-empty text", func() {
				err := forms.EventTextOptional("Too short")
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("CompanyName", func() {
			It("should pass for valid company names", func() {
				err := forms.CompanyName("Acme Corp")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass for empty strings", func() {
				err := forms.CompanyName("")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail for too long names", func() {
				longName := string(make([]byte, 101))
				err := forms.CompanyName(longName)
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("Title", func() {
			It("should pass for valid titles", func() {
				err := forms.Title("Valid Title")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail for empty titles", func() {
				err := forms.Title("")
				Expect(err).To(HaveOccurred())
			})

			It("should fail for too short titles", func() {
				err := forms.Title("AB")
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("ProfileName", func() {
			It("should pass for valid profile names", func() {
				err := forms.ProfileName("Software Engineer")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should fail for empty names", func() {
				err := forms.ProfileName("")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("OneOf", func() {
		It("should pass when value is in allowed list", func() {
			validator := forms.OneOf([]string{"quick", "manual", "csv"})
			err := validator("quick")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when value is not in allowed list", func() {
			validator := forms.OneOf([]string{"quick", "manual", "csv"})
			err := validator("invalid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be one of"))
		})
	})

	Describe("Custom", func() {
		It("should pass when custom logic succeeds", func() {
			validator := forms.Custom(
				func(val string) bool { return len(val) == 6 },
				"must be exactly 6 characters",
			)
			err := validator("hello!")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when custom logic fails", func() {
			validator := forms.Custom(
				func(val string) bool { return len(val) == 6 },
				"must be exactly 6 characters",
			)
			err := validator("hi")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("must be exactly 6 characters"))
		})
	})

	Describe("GitHubUsername", func() {
		It("should pass for valid usernames", func() {
			err := forms.GitHubUsername("baphled")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for usernames with hyphens", func() {
			err := forms.GitHubUsername("john-doe")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for usernames with numbers", func() {
			err := forms.GitHubUsername("user123")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for single character usernames", func() {
			err := forms.GitHubUsername("a")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for empty strings (optional field)", func() {
			err := forms.GitHubUsername("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for full GitHub URLs", func() {
			err := forms.GitHubUsername("https://github.com/baphled")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("username only"))
		})

		It("should fail for URLs with github.com", func() {
			err := forms.GitHubUsername("github.com/baphled")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("username only"))
		})

		It("should fail for usernames starting with hyphen", func() {
			err := forms.GitHubUsername("-baphled")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot start or end with a hyphen"))
		})

		It("should fail for usernames ending with hyphen", func() {
			err := forms.GitHubUsername("baphled-")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot start or end with a hyphen"))
		})

		It("should fail for usernames with consecutive hyphens", func() {
			err := forms.GitHubUsername("baph--led")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("consecutive hyphens"))
		})

		It("should fail for usernames with special characters", func() {
			err := forms.GitHubUsername("baph_led")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("alphanumeric"))
		})

		It("should fail for usernames with spaces", func() {
			err := forms.GitHubUsername("baph led")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("alphanumeric"))
		})

		It("should fail for usernames longer than 39 characters", func() {
			longUsername := "abcdefghijklmnopqrstuvwxyz1234567890abcd" // 40 chars
			err := forms.GitHubUsername(longUsername)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("39 characters"))
		})

		It("should pass for exactly 39 character usernames", func() {
			username39 := "abcdefghijklmnopqrstuvwxyz1234567890abc" // 39 chars
			err := forms.GitHubUsername(username39)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GitHubURL", func() {
		// CONTRACT: The GitHub field in config stores USERNAME ONLY (e.g., "baphled"),
		// not the full URL. GitHubURL() decorates it for display in CV exports.
		// This is validated by GitHubUsername() which rejects full URLs.

		It("should format username as full URL", func() {
			url := forms.GitHubURL("baphled")
			Expect(url).To(Equal("https://github.com/baphled"))
		})

		It("should handle usernames with hyphens", func() {
			url := forms.GitHubURL("john-doe")
			Expect(url).To(Equal("https://github.com/john-doe"))
		})

		It("should return empty string for empty username", func() {
			url := forms.GitHubURL("")
			Expect(url).To(Equal(""))
		})

		It("should trim whitespace from username", func() {
			url := forms.GitHubURL("  baphled  ")
			Expect(url).To(Equal("https://github.com/baphled"))
		})
	})

	Describe("ParseDateString", func() {
		It("should parse YYYY-MM-DD format", func() {
			t, err := forms.ParseDateString("2024-01-15")
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Year()).To(Equal(2024))
			Expect(int(t.Month())).To(Equal(1))
			Expect(t.Day()).To(Equal(15))
		})

		It("should parse 'today' case-insensitively", func() {
			t, err := forms.ParseDateString("today")
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Year()).To(Equal(time.Now().Year()))
		})

		It("should parse 'TODAY' uppercase", func() {
			t, err := forms.ParseDateString("TODAY")
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Year()).To(Equal(time.Now().Year()))
		})

		It("should parse short relative format -Nd", func() {
			t, err := forms.ParseDateString("-7d")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -7)
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(t.Month()).To(Equal(expected.Month()))
			Expect(t.Day()).To(Equal(expected.Day()))
		})

		It("should parse short relative format -Nw", func() {
			t, err := forms.ParseDateString("-2w")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -14)
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(t.Month()).To(Equal(expected.Month()))
			Expect(t.Day()).To(Equal(expected.Day()))
		})

		It("should parse short relative format -Nm", func() {
			t, err := forms.ParseDateString("-1m")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, -1, 0)
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(t.Month()).To(Equal(expected.Month()))
		})

		It("should parse long relative format 'N days ago'", func() {
			t, err := forms.ParseDateString("1 day ago")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -1)
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(t.Month()).To(Equal(expected.Month()))
			Expect(t.Day()).To(Equal(expected.Day()))
		})

		It("should parse long relative format 'N weeks ago'", func() {
			t, err := forms.ParseDateString("2 weeks ago")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, 0, -14)
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(t.Month()).To(Equal(expected.Month()))
			Expect(t.Day()).To(Equal(expected.Day()))
		})

		It("should parse long relative format 'N months ago'", func() {
			t, err := forms.ParseDateString("3 months ago")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now().AddDate(0, -3, 0)
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(t.Month()).To(Equal(expected.Month()))
		})

		It("should fail for invalid format", func() {
			_, err := forms.ParseDateString("invalid")
			Expect(err).To(HaveOccurred())
		})

		It("should trim whitespace", func() {
			t, err := forms.ParseDateString("  2024-01-15  ")
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Year()).To(Equal(2024))
		})

		It("should parse 'today' with mixed case", func() {
			t, err := forms.ParseDateString("ToDay")
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Year()).To(Equal(time.Now().Year()))
		})

		It("should parse -0d (today)", func() {
			t, err := forms.ParseDateString("-0d")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now()
			Expect(t.Year()).To(Equal(expected.Year()))
			Expect(int(t.Month())).To(Equal(int(expected.Month())))
			Expect(t.Day()).To(Equal(expected.Day()))
		})

		It("should parse 0 days ago (today)", func() {
			t, err := forms.ParseDateString("0 days ago")
			Expect(err).NotTo(HaveOccurred())
			expected := time.Now()
			Expect(t.Year()).To(Equal(expected.Year()))
		})

		It("should fail for invalid relative format", func() {
			_, err := forms.ParseDateString("-abc")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for invalid date values", func() {
			_, err := forms.ParseDateString("2024-13-01")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for empty string", func() {
			_, err := forms.ParseDateString("")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("TagName", func() {
		It("should pass for valid tag names", func() {
			err := forms.TagName("technical")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for empty strings", func() {
			err := forms.TagName("")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for whitespace-only strings", func() {
			err := forms.TagName("   ")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for tags with special characters", func() {
			err := forms.TagName("tag@name")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for tags exceeding max length", func() {
			longTag := "a" + strings.Repeat("b", 50)
			err := forms.TagName(longTag)
			Expect(err).To(HaveOccurred())
		})

		It("should pass for tags with hyphens and underscores", func() {
			err := forms.TagName("tag-name_123")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("AudienceName", func() {
		It("should pass for valid audience names", func() {
			err := forms.AudienceName("Hiring Manager")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for empty strings", func() {
			err := forms.AudienceName("")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for whitespace-only strings", func() {
			err := forms.AudienceName("   ")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for names shorter than 2 characters", func() {
			err := forms.AudienceName("a")
			Expect(err).To(HaveOccurred())
		})

		It("should fail for names exceeding max length", func() {
			longName := strings.Repeat("a", 101)
			err := forms.AudienceName(longName)
			Expect(err).To(HaveOccurred())
		})

		It("should pass for exactly 2 character names", func() {
			err := forms.AudienceName("HR")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass for exactly 100 character names", func() {
			name100 := strings.Repeat("a", 100)
			err := forms.AudienceName(name100)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
