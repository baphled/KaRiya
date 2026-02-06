---
name: api-design
description: Design clean, consistent APIs - RESTful conventions, versioning, backwards compatibility, error handling
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide the design of clean, consistent, and maintainable APIs. This covers internal Go interfaces, service boundaries, and external HTTP APIs.

## When to use me

- Designing new service interfaces
- Creating HTTP/REST endpoints
- Reviewing API contracts
- Planning API versioning
- Ensuring backwards compatibility

## Core Principles

1. **Consistency** - Same patterns everywhere
2. **Predictability** - Behave as users expect
3. **Simplicity** - Easy to use correctly, hard to misuse
4. **Evolvability** - Can change without breaking clients

## Go Interface Design

### Keep Interfaces Small

```go
// GOOD - Small, focused interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Compose when needed
type ReadWriter interface {
    Reader
    Writer
}

// BAD - Kitchen sink interface
type Storage interface {
    Read(id string) (*Entity, error)
    Write(entity *Entity) error
    Delete(id string) error
    List() ([]*Entity, error)
    Search(query string) ([]*Entity, error)
    Count() (int, error)
    Exists(id string) (bool, error)
    // ... 20 more methods
}
```

### Accept Interfaces, Return Structs

```go
// GOOD - Accept interface
func ProcessData(r io.Reader) error {
    // Can accept any Reader
}

// GOOD - Return concrete type
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

// BAD - Return interface (usually)
func NewService(repo Repository) ServiceInterface {
    return &Service{repo: repo}
}
```

### Define Interfaces Where Used

```go
// GOOD - Interface defined by consumer
// In: internal/cli/intents/timeline/intent.go
type eventRepository interface {
    FindByDateRange(start, end time.Time) ([]*Event, error)
}

// BAD - Interface defined by provider
// In: internal/repository/event_repository.go
type EventRepositoryInterface interface {
    // All methods whether needed or not
}
```

## Service API Design

### Constructor Pattern

```go
// Required dependencies via constructor
func NewEventService(repo EventRepository, logger Logger) *EventService {
    return &EventService{
        repo:   repo,
        logger: logger,
    }
}

// Optional dependencies via setters
func (s *EventService) WithCache(cache Cache) *EventService {
    s.cache = cache
    return s
}

// Or functional options
func NewEventService(repo EventRepository, opts ...Option) *EventService {
    s := &EventService{repo: repo}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### Method Naming

```go
// CRUD operations
Create(entity *Entity) (*Entity, error)
Get(id string) (*Entity, error)      // Single item
List(filter Filter) ([]*Entity, error) // Multiple items
Update(entity *Entity) (*Entity, error)
Delete(id string) error

// Queries
FindByName(name string) ([]*Entity, error)
FindByDateRange(start, end time.Time) ([]*Entity, error)
ExistsByID(id string) (bool, error)
CountByStatus(status Status) (int, error)

// Actions
Start(id string) error
Complete(id string) error
Cancel(id string) error
```

### Error Design

```go
// Sentinel errors for expected conditions
var (
    ErrNotFound      = errors.New("entity not found")
    ErrAlreadyExists = errors.New("entity already exists")
    ErrInvalidInput  = errors.New("invalid input")
)

// Wrap with context
func (s *Service) Get(id string) (*Entity, error) {
    entity, err := s.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("get entity %s: %w", id, err)
    }
    return entity, nil
}

// Check specific errors
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

## HTTP API Design

### RESTful Resource Naming

```
# Resources are nouns, plural
GET    /events           # List events
POST   /events           # Create event
GET    /events/{id}      # Get single event
PUT    /events/{id}      # Update event (full)
PATCH  /events/{id}      # Update event (partial)
DELETE /events/{id}      # Delete event

# Nested resources
GET    /events/{id}/skills    # Skills for event
POST   /events/{id}/skills    # Add skill to event

# Actions (when CRUD doesn't fit)
POST   /events/{id}/complete  # Complete event
POST   /events/{id}/cancel    # Cancel event
```

### Query Parameters

```
# Filtering
GET /events?status=active&type=meeting

# Pagination
GET /events?page=2&per_page=20
GET /events?offset=20&limit=20
GET /events?cursor=abc123  # Cursor-based

# Sorting
GET /events?sort=date&order=desc
GET /events?sort=-date,+name  # Prefix notation

# Field selection
GET /events?fields=id,name,date

# Search
GET /events?q=meeting
```

### Response Structure

```json
// Single resource
{
  "data": {
    "id": "123",
    "type": "event",
    "attributes": {
      "name": "Team Meeting",
      "date": "2024-03-15T10:00:00Z"
    }
  }
}

// Collection
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 2,
    "per_page": 20
  },
  "links": {
    "self": "/events?page=2",
    "next": "/events?page=3",
    "prev": "/events?page=1"
  }
}

// Error
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      {"field": "email", "message": "must be valid email"}
    ]
  }
}
```

### HTTP Status Codes

| Code | Meaning | Use For |
|------|---------|---------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST creating resource |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid input, validation error |
| 401 | Unauthorised | Missing/invalid authentication |
| 403 | Forbidden | Authenticated but not allowed |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource state conflict |
| 422 | Unprocessable | Valid syntax but semantic error |
| 500 | Server Error | Unexpected server failure |

## Versioning

### URL Versioning

```
/api/v1/events
/api/v2/events
```

### Header Versioning

```
Accept: application/vnd.kariya.v1+json
```

### Backwards Compatibility Rules

**Safe changes (don't require version bump):**
- Adding new endpoints
- Adding optional fields to requests
- Adding fields to responses
- Adding new enum values (if clients ignore unknown)

**Breaking changes (require version bump):**
- Removing endpoints
- Removing or renaming fields
- Changing field types
- Changing required/optional status
- Changing validation rules
- Changing error codes

### Deprecation Strategy

```go
// Document deprecation
// Deprecated: Use GetEventByID instead. Will be removed in v3.
func (s *Service) GetEvent(id string) (*Event, error)

// HTTP header
Deprecation: true
Sunset: Sat, 1 Jan 2025 00:00:00 GMT
Link: </api/v2/events>; rel="successor-version"
```

## Documentation

### Go Doc for Interfaces

```go
// EventRepository provides access to event storage.
//
// Implementations must be safe for concurrent use.
type EventRepository interface {
    // FindByID retrieves a single event.
    //
    // Expected: id - non-empty event identifier
    // Returns: event if found, ErrNotFound if not exists
    // Side effects: None
    FindByID(id string) (*Event, error)
    
    // Save persists an event, creating or updating as needed.
    //
    // Expected: event - valid event with non-empty ID
    // Returns: saved event with updated timestamps
    // Side effects: Writes to database
    Save(event *Event) (*Event, error)
}
```

### OpenAPI for HTTP APIs

```yaml
paths:
  /events/{id}:
    get:
      summary: Get event by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Event found
        '404':
          description: Event not found
```

## Related Skills

- `domain-modeling` - Domain entities and validation
- `service-layer` - Service implementation
- `error-handling` - Error design
- `documentation-writing` - API documentation
