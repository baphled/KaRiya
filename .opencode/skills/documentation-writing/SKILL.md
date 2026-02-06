---
name: documentation-writing
description: Write clear technical documentation - READMEs, ADRs, runbooks, API docs, inline documentation
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide writing clear, useful technical documentation that helps readers accomplish their goals. Good documentation saves time and reduces confusion.

## When to use me

- Writing READMEs
- Creating architecture decision records (ADRs)
- Writing runbooks and guides
- Documenting APIs
- Adding code documentation

## Core Principles

1. **Reader-focused** - Write for the audience, not yourself
2. **Task-oriented** - Help readers accomplish goals
3. **Accurate** - Wrong docs are worse than no docs
4. **Maintainable** - Keep docs close to code, update together
5. **Scannable** - Use structure for quick navigation

## Document Types

### README.md

```markdown
# Project Name

Brief description of what this project does.

## Quick Start

```bash
# Minimal steps to get running
git clone ...
make run
```

## Features

- Feature 1
- Feature 2

## Installation

Detailed installation instructions...

## Usage

Common use cases with examples...

## Configuration

Environment variables, config files...

## Development

How to set up development environment...

## Testing

How to run tests...

## Contributing

How to contribute...

## License

License information...
```

### Architecture Decision Records (ADRs)

```markdown
# ADR-001: Use SQLite for Local Storage

## Status

Accepted

## Context

We need persistent storage for career events. Options considered:
- SQLite
- PostgreSQL
- File-based JSON

## Decision

Use SQLite because:
- Zero configuration
- Single file database
- Sufficient for expected data volume
- Good Go support (mattn/go-sqlite3)

## Consequences

**Positive:**
- Simple deployment
- No external dependencies
- Easy backup (copy file)

**Negative:**
- Limited concurrent writes
- No built-in replication
- May need migration if we scale

## References

- SQLite docs: https://sqlite.org
- Go driver: https://github.com/mattn/go-sqlite3
```

### Runbooks

```markdown
# Runbook: Database Migration

## When to Use

When deploying changes that require database schema updates.

## Prerequisites

- [ ] Database backup taken
- [ ] Deployment window scheduled
- [ ] Rollback plan ready

## Procedure

### 1. Pre-checks

```bash
# Verify current schema version
make db-version
```

### 2. Run Migration

```bash
# Dry run first
make migrate-dry-run

# If successful, apply
make migrate
```

### 3. Verify

```bash
# Check schema version
make db-version

# Run smoke tests
make smoke-test
```

## Rollback

If issues occur:

```bash
# Rollback last migration
make migrate-down

# Verify
make db-version
```

## Troubleshooting

### Migration fails with "table exists"

1. Check current schema state
2. Manually verify table structure
3. If safe, use `--force` flag

### Application errors after migration

1. Check logs for specific errors
2. Verify all columns exist
3. Rollback if necessary
```

## Writing Guidelines

### Use Active Voice

```markdown
# GOOD - Active voice
"Run the test suite before committing."
"The service validates input before processing."

# BAD - Passive voice
"The test suite should be run before committing."
"Input is validated by the service before being processed."
```

### Be Direct

```markdown
# GOOD - Direct
"Set `LOG_LEVEL=debug` to enable debug logging."

# BAD - Indirect
"If you want to see debug logs, you might want to consider 
setting the LOG_LEVEL environment variable to debug."
```

### Use Consistent Terminology

```markdown
# Pick one and stick to it
"event" not sometimes "event" and sometimes "entry"
"repository" not sometimes "repository" and sometimes "repo"

# Create a glossary if needed
## Glossary
- **Event**: A career event (meeting, task, achievement)
- **Burst**: A cluster of related events
```

### Structure for Scanning

```markdown
# Use headings hierarchically
## Main Section
### Subsection
#### Detail

# Use lists for sequences
1. First step
2. Second step
3. Third step

# Use tables for comparisons
| Option | Pros | Cons |
|--------|------|------|
| A      | Fast | Complex |
| B      | Simple | Slow |

# Use code blocks for commands
```bash
make test
```
```

## Code Documentation

### Go Doc Comments

```go
// Package events provides career event management.
//
// # Overview
//
// The events package handles creation, storage, and retrieval
// of career events. Events represent activities like meetings,
// tasks, and achievements.
//
// # Usage
//
//	svc := events.NewService(repo)
//	event, err := svc.Create(ctx, &events.Event{
//	    Title: "Team Meeting",
//	    Type:  events.TypeMeeting,
//	})
package events

// Event represents a career event.
//
// Events are the core unit of career tracking. Each event
// has a type, timestamp, and associated metadata.
type Event struct {
    ID        string
    Title     string
    Type      EventType
    Timestamp time.Time
}

// Create persists a new event.
//
// Expected: event with non-empty Title and valid Type
// Returns: created event with generated ID, or error
// Side effects: Writes to database, emits EventCreated metric
func (s *Service) Create(ctx context.Context, event *Event) (*Event, error)
```

### Inline Comments (When Necessary)

```go
// GOOD - Explains WHY
// We retry 3 times because the external API has occasional transient failures
for i := 0; i < 3; i++ {

// BAD - Explains WHAT (obvious from code)
// Loop 3 times
for i := 0; i < 3; i++ {
```

## API Documentation

### Endpoint Documentation

```markdown
## Create Event

Creates a new career event.

### Request

`POST /api/v1/events`

```json
{
  "title": "Team Meeting",
  "type": "meeting",
  "timestamp": "2024-03-15T10:00:00Z",
  "description": "Weekly sync"
}
```

### Response

**Success (201 Created)**

```json
{
  "data": {
    "id": "evt_123",
    "title": "Team Meeting",
    "type": "meeting",
    "timestamp": "2024-03-15T10:00:00Z",
    "created_at": "2024-03-14T08:30:00Z"
  }
}
```

**Error (400 Bad Request)**

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      {"field": "title", "message": "required"}
    ]
  }
}
```

### Example

```bash
curl -X POST https://api.example.com/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"title": "Team Meeting", "type": "meeting"}'
```
```

## Documentation Maintenance

### Keep Docs Near Code

```
internal/service/
├── event_service.go
├── event_service_test.go
└── README.md           # Package-level docs

docs/
├── architecture/       # High-level architecture
├── adr/               # Architecture decisions
└── runbooks/          # Operational guides
```

### Update Together

```bash
# Commit code and docs together
git add internal/service/event_service.go
git add internal/service/README.md
git commit -m "feat(events): add event validation

Updated README with validation rules."
```

### Review Docs in PRs

```markdown
## PR Checklist
- [ ] Code changes
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] README reflects changes
```

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Outdated docs | Update docs with code |
| Too much detail | Focus on what reader needs |
| Missing examples | Add concrete examples |
| Jargon without explanation | Define terms or link to glossary |
| Assuming knowledge | State prerequisites |
| Wall of text | Use structure, lists, headings |

## Related Skills

- `british-english` - Language standards
- `api-design` - API documentation
- `clean-code` - Self-documenting code
