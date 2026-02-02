// Package career provides repository interfaces and types for career domain entities.
package career

import "io"

// Repositories holds all repository implementations.
// Fields are interface types so either SQL or Memory implementations can be used.
type Repositories struct {
	Event EventRepository
	Skill SkillRepository
	Fact  FactRepository
	Burst BurstRepository

	closer io.Closer
}

// SetCloser sets the closer used to release underlying resources (e.g. database connections).
//
// Expected:
//   - closer must be valid.
//
// Side effects:
//   - None.
func (r *Repositories) SetCloser(c io.Closer) {
	r.closer = c
}

// Close releases underlying resources.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *Repositories) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}
