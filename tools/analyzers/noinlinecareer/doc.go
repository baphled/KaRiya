// Package noinlinecareer provides a static analysis tool that detects
// improper inline usage of career domain types outside the domain and
// model layers.
//
// The analyzer enforces that career.Event, career.Fact, career.Burst,
// and career.Skill types are only referenced directly in
// internal/domain/career, internal/model/career, and
// internal/repository/career packages. Other packages must use display
// types or interfaces to avoid tight coupling to the domain layer.
package noinlinecareer
