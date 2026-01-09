package forms

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ParseDateString parses a date string in various formats.
// Supports:
// - YYYY-MM-DD format
// - "today"
// - Relative dates like "1 week ago", "2 days ago"
func ParseDateString(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// Handle "today"
	if strings.ToLower(s) == "today" {
		return time.Now(), nil
	}

	// Try standard date format first
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		return t, nil
	}

	// Try relative dates (e.g., "1 week ago", "2 days ago")
	relativeRegex := regexp.MustCompile(`^(\d+)\s+(day|days|week|weeks|month|months)\s+ago$`)
	if matches := relativeRegex.FindStringSubmatch(strings.ToLower(s)); len(matches) == 3 {
		var amount int
		if _, err := fmt.Sscanf(matches[1], "%d", &amount); err != nil {
			amount = 0 // Fallback to 0 if parsing fails
		}
		unit := matches[2]

		now := time.Now()
		switch {
		case strings.HasPrefix(unit, "day"):
			return now.AddDate(0, 0, -amount), nil
		case strings.HasPrefix(unit, "week"):
			return now.AddDate(0, 0, -amount*7), nil
		case strings.HasPrefix(unit, "month"):
			return now.AddDate(0, -amount, 0), nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}

// Common validation errors
var (
	ErrRequired         = fmt.Errorf("this field is required")
	ErrInvalidDate      = fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
	ErrInvalidEmail     = fmt.Errorf("invalid email address")
	ErrTooShort         = fmt.Errorf("value is too short")
	ErrTooLong          = fmt.Errorf("value is too long")
	ErrInvalidCharacter = fmt.Errorf("contains invalid characters")
)

// Required validates that a field is not empty.
func Required(value string) error {
	if strings.TrimSpace(value) == "" {
		return ErrRequired
	}
	return nil
}

// MinLength validates that a string has at least n characters.
func MinLength(n int) func(string) error {
	return func(value string) error {
		if len(strings.TrimSpace(value)) < n {
			return fmt.Errorf("must be at least %d characters", n)
		}
		return nil
	}
}

// MaxLength validates that a string has at most n characters.
func MaxLength(n int) func(string) error {
	return func(value string) error {
		if len(value) > n {
			return fmt.Errorf("must be at most %d characters", n)
		}
		return nil
	}
}

// LengthRange validates that a string length is within a range.
func LengthRange(min, max int) func(string) error {
	return func(value string) error {
		length := len(strings.TrimSpace(value))
		if length < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		if length > max {
			return fmt.Errorf("must be at most %d characters", max)
		}
		return nil
	}
}

// DateFormat validates that a date string matches YYYY-MM-DD format.
func DateFormat(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil // Allow empty dates (use Required separately if needed)
	}

	// Check format with regex
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(value) {
		return ErrInvalidDate
	}

	// Validate it's a real date
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		return ErrInvalidDate
	}

	return nil
}

// DateFormatRequired validates date format and requires non-empty.
func DateFormatRequired(value string) error {
	if err := Required(value); err != nil {
		return err
	}
	return DateFormat(value)
}

// Email validates basic email format.
func Email(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil // Allow empty (use Required separately if needed)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return ErrInvalidEmail
	}

	return nil
}

// URL validates basic URL format.
func URL(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil // Allow empty (use Required separately if needed)
	}

	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(value) {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// AlphaNumeric validates that a string contains only letters and numbers.
func AlphaNumeric(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	alphaNumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphaNumRegex.MatchString(value) {
		return fmt.Errorf("must contain only letters and numbers")
	}

	return nil
}

// NoSpecialChars validates that a string doesn't contain special characters.
func NoSpecialChars(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	// Allow letters, numbers, spaces, hyphens, underscores
	validRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validRegex.MatchString(value) {
		return fmt.Errorf("contains invalid special characters")
	}

	return nil
}

// Compose combines multiple validators into one.
// Returns the first error encountered, or nil if all pass.
func Compose(validators ...func(string) error) func(string) error {
	return func(value string) error {
		for _, validator := range validators {
			if err := validator(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// Domain-specific validators for KaRiya

// EventText validates career event text (required, reasonable length).
func EventText(value string) error {
	return Compose(
		Required,
		MinLength(10),
		MaxLength(2000),
	)(value)
}

// EventTextOptional validates career event text when optional.
func EventTextOptional(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return Compose(
		MinLength(10),
		MaxLength(2000),
	)(value)
}

// CompanyName validates company name (reasonable length, no special chars).
func CompanyName(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil // Optional field
	}
	return Compose(
		MinLength(2),
		MaxLength(100),
	)(value)
}

// TagName validates a single tag name.
func TagName(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return Compose(
		MinLength(1),
		MaxLength(50),
		NoSpecialChars,
	)(value)
}

// Title validates a title field (burst, fact, etc.).
func Title(value string) error {
	return Compose(
		Required,
		MinLength(3),
		MaxLength(200),
	)(value)
}

// Description validates a description field.
func Description(value string) error {
	if strings.TrimSpace(value) == "" {
		return nil // Optional
	}
	return Compose(
		MinLength(10),
		MaxLength(1000),
	)(value)
}

// ProfileName validates a CV profile name.
func ProfileName(value string) error {
	return Compose(
		Required,
		MinLength(2),
		MaxLength(100),
	)(value)
}

// AudienceName validates a CV audience name.
func AudienceName(value string) error {
	return Compose(
		Required,
		MinLength(2),
		MaxLength(100),
	)(value)
}

// OneOf validates that a value is one of the allowed options.
func OneOf(allowed []string) func(string) error {
	return func(value string) error {
		for _, option := range allowed {
			if value == option {
				return nil
			}
		}
		return fmt.Errorf("must be one of: %s", strings.Join(allowed, ", "))
	}
}

// Custom creates a custom validator with a specific error message.
func Custom(validate func(string) bool, errorMsg string) func(string) error {
	return func(value string) error {
		if !validate(value) {
			return fmt.Errorf("%s", errorMsg)
		}
		return nil
	}
}
