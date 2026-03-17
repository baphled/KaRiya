package career

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringSlice represents a slice of strings that serializes to and from a
// comma-separated string for storage in SQLite text columns.
//
// The encoding preserves backward compatibility with existing data stored in
// the format "a,b,c". StringSlice satisfies both sql.Scanner and
// driver.Valuer so GORM handles the conversion transparently during reads
// and writes.
type StringSlice []string

// Scan deserializes a raw database column value into the StringSlice receiver.
//
// Expected:
//   - interface{} must be valid.
//
// Returns:
//   - An error value.
//
// Side effects:
//   - None.
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("unsupported type for StringSlice: %T", value)
	}
	if str == "" {
		*s = nil
		return nil
	}
	*s = strings.Split(str, ",")
	return nil
}

// Value serializes the StringSlice into a comma-separated string suitable for
// database storage.
//
// Returns:
//   - An empty string and nil error when the slice is nil or empty.
//   - A single comma-joined string and nil error when the slice contains
//     one or more elements.
//
// Side effects:
//   - None.
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "", nil
	}
	return strings.Join(s, ","), nil
}
