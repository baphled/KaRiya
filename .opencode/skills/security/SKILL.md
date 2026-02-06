---
name: security
description: Security best practices for KaRiya including input validation, SQL injection prevention, and secure coding
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Enforce security best practices in the KaRiya codebase.

## When to use me

Use this skill when:
- Handling user input
- Writing database queries
- Working with file paths
- Implementing authentication/authorization
- Reviewing code for security issues

## Security Principles

1. **Never trust user input** - Validate everything
2. **Principle of least privilege** - Minimal permissions
3. **Defense in depth** - Multiple security layers
4. **Fail securely** - Errors should deny access
5. **Keep it simple** - Complex code has more vulnerabilities

## Input Validation

### Always Validate

```go
// BAD - No validation
func ProcessEvent(title string) error {
    return repo.Create(&Event{Title: title})
}

// GOOD - Validate input
func ProcessEvent(title string) error {
    if title == "" {
        return errors.New("title is required")
    }
    if len(title) > 500 {
        return errors.New("title too long")
    }
    if !isValidTitle(title) {
        return errors.New("title contains invalid characters")
    }
    return repo.Create(&Event{Title: title})
}
```

### Validation Patterns

```go
// Use validation in Context structs
type CreateEventContext struct {
    Title       string
    Description string
    Date        time.Time
}

func (c *CreateEventContext) Validate() error {
    if c.Title == "" {
        return errors.New("title is required")
    }
    if c.Date.IsZero() {
        return errors.New("date is required")
    }
    if c.Date.Before(time.Now().AddDate(-50, 0, 0)) {
        return errors.New("date too far in the past")
    }
    return nil
}
```

## SQL Injection Prevention

### GORM is Safe by Default

```go
// GOOD - Parameterized query (GORM default)
db.Where("id = ?", userInput).First(&event)

// GOOD - Struct-based queries
db.Where(&Event{Title: userInput}).First(&event)

// BAD - String concatenation (NEVER DO THIS)
db.Raw("SELECT * FROM events WHERE id = '" + userInput + "'")

// BAD - fmt.Sprintf in queries (VULNERABLE)
db.Raw(fmt.Sprintf("SELECT * FROM events WHERE id = '%s'", userInput))
```

### Safe Raw Queries

```go
// If you must use raw SQL, use parameters
db.Raw("SELECT * FROM events WHERE title LIKE ?", "%"+searchTerm+"%")
```

## File Path Security

### Path Traversal Prevention

```go
// BAD - Allows path traversal
func ReadFile(filename string) ([]byte, error) {
    return os.ReadFile("/data/" + filename)
    // Attacker: filename = "../../../etc/passwd"
}

// GOOD - Validate and sanitize paths
func ReadFile(filename string) ([]byte, error) {
    // Remove any path components
    cleanName := filepath.Base(filename)
    
    // Validate against allowed pattern
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+\.json$`).MatchString(cleanName) {
        return nil, errors.New("invalid filename")
    }
    
    fullPath := filepath.Join("/data", cleanName)
    
    // Verify it's still under the base directory
    if !strings.HasPrefix(fullPath, "/data/") {
        return nil, errors.New("invalid path")
    }
    
    return os.ReadFile(fullPath)
}
```

## Error Handling Security

### Don't Leak Internal Details

```go
// BAD - Exposes internal structure
if err != nil {
    return fmt.Errorf("database error: %v", err)
    // Reveals: "database error: pq: relation 'users' does not exist"
}

// GOOD - Generic external message, log details
if err != nil {
    log.Printf("database error: %v", err) // Internal log
    return errors.New("an error occurred processing your request")
}
```

### Secure Error Types

```go
// Use sentinel errors for expected cases
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
)

// Check specific errors
if errors.Is(err, ErrNotFound) {
    return nil, ErrNotFound  // Safe to return
}
// Unknown error - don't expose details
return nil, errors.New("internal error")
```

## Secrets Management

### Never Hardcode Secrets

```go
// BAD
apiKey := "sk_live_abc123..."

// GOOD - Environment variables
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    return errors.New("API_KEY not configured")
}
```

### Files to Never Commit

| File Pattern | Risk |
|--------------|------|
| `.env` | Environment secrets |
| `*.pem`, `*.key` | Private keys |
| `credentials.json` | API credentials |
| `*_secret*` | Any secrets |
| `config.local.*` | Local overrides |

Check `.gitignore` includes these patterns.

## Secure Defaults

### Timeouts

```go
// Always set timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// HTTP clients need timeouts
client := &http.Client{
    Timeout: 30 * time.Second,
}
```

### File Permissions

```go
// Restrictive permissions for sensitive files
os.WriteFile(path, data, 0600)  // Owner read/write only

// Directories
os.MkdirAll(path, 0700)  // Owner only
```

## Security Scanning

### Run gosec

```bash
make gosec
```

This checks for:
- Hardcoded credentials
- SQL injection risks
- Insecure random number generation
- Path traversal vulnerabilities
- And more...

### Common gosec Warnings

| Warning | Fix |
|---------|-----|
| G101: Hardcoded credentials | Use environment variables |
| G201: SQL string concatenation | Use parameterized queries |
| G304: File path from variable | Validate and sanitize path |
| G401: Use of weak crypto | Use crypto/rand, not math/rand |

## Security Checklist

Before merging:

- [ ] All user input is validated
- [ ] No SQL string concatenation
- [ ] No hardcoded secrets
- [ ] File paths are sanitized
- [ ] Errors don't leak internal details
- [ ] Timeouts are set for external calls
- [ ] `make gosec` passes

## Related skills

- `code-reviewer` - Security review
- `db-operations` - Secure database access
- `clean-code` - Security through simplicity
