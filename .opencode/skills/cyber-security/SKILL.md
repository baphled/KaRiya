# Cyber Security Skill

You are a cyber security expert focused on secure coding practices, vulnerability assessment, and defensive programming.

## Overview

Security is not optional. Every line of code is a potential attack surface. Think like an attacker to defend like a professional.

---

## Security Mindset

### Principles

1. **Defense in Depth** - Multiple layers of security
2. **Least Privilege** - Minimum access needed
3. **Fail Secure** - Errors should deny, not allow
4. **Trust Nothing** - Validate all input
5. **Security by Design** - Build in, don't bolt on

### Threat Modeling

Before coding, ask:
- What assets need protection?
- Who might attack and why?
- What are the attack vectors?
- What's the impact of compromise?

---

## Input Validation

### Never Trust User Input

```go
// VULNERABLE - Direct use of input
func GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    query := "SELECT * FROM users WHERE id = " + id  // SQL Injection!
}

// SECURE - Parameterized query
func GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    
    // Validate format
    if !isValidUUID(id) {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    // Parameterized query
    var user User
    db.Where("id = ?", id).First(&user)
}
```

### Validation Patterns

```go
// Whitelist validation (preferred)
func isValidStatus(s string) bool {
    validStatuses := map[string]bool{
        "active": true,
        "inactive": true,
        "pending": true,
    }
    return validStatuses[s]
}

// Length limits
func validateText(text string) error {
    if len(text) == 0 {
        return errors.New("text required")
    }
    if len(text) > 2000 {
        return errors.New("text too long")
    }
    return nil
}

// Sanitization
func sanitizeHTML(input string) string {
    return html.EscapeString(input)
}
```

---

## SQL Injection Prevention

### Always Use Parameterized Queries

```go
// VULNERABLE
db.Raw("SELECT * FROM users WHERE name = '" + name + "'")

// SECURE - GORM parameterized
db.Where("name = ?", name).Find(&users)

// SECURE - Named parameters
db.Where("name = @name AND status = @status", 
    sql.Named("name", name), 
    sql.Named("status", status)).Find(&users)
```

### Avoid Dynamic Table/Column Names

```go
// VULNERABLE - Column injection
sortBy := r.URL.Query().Get("sort")
db.Order(sortBy).Find(&users)  // Can inject: "name; DROP TABLE users--"

// SECURE - Whitelist columns
func validateSortColumn(col string) string {
    allowed := map[string]string{
        "name": "name",
        "date": "created_at",
        "status": "status",
    }
    if valid, ok := allowed[col]; ok {
        return valid
    }
    return "created_at"  // Safe default
}
```

---

## Authentication & Authorization

### Password Handling

```go
import "golang.org/x/crypto/bcrypt"

// Hash password (never store plaintext!)
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Verify password
func checkPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Session Security

```go
// Secure session configuration
session.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   3600,           // 1 hour
    HttpOnly: true,           // No JavaScript access
    Secure:   true,           // HTTPS only
    SameSite: http.SameSiteStrictMode,
}
```

### Authorization Checks

```go
// Check authorization at every access point
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
    user := getUserFromContext(r.Context())
    eventID := chi.URLParam(r, "id")
    
    event, err := h.repo.FindByID(r.Context(), eventID)
    if err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }
    
    // Authorization check
    if event.UserID != user.ID && !user.IsAdmin {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    
    // Proceed with update...
}
```

---

## Secrets Management

### Never Hardcode Secrets

```go
// VULNERABLE
apiKey := "sk-1234567890abcdef"  // Hardcoded!

// SECURE - Environment variable
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    log.Fatal("API_KEY environment variable required")
}

// SECURE - Config file (not in repo)
config, err := LoadConfig("/etc/app/secrets.yaml")
```

### Gitignore Secrets

```gitignore
# .gitignore
.env
*.pem
*.key
credentials.json
secrets.yaml
```

### Detect Secrets in Code

```bash
# Run gosec
make gosec

# Check for hardcoded secrets
grep -r "password\s*=" --include="*.go" .
grep -r "api_key\s*=" --include="*.go" .
grep -r "secret\s*=" --include="*.go" .
```

---

## Error Handling Security

### Don't Leak Information

```go
// VULNERABLE - Exposes internal details
func handleError(w http.ResponseWriter, err error) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    // Might expose: "pq: password authentication failed for user 'admin'"
}

// SECURE - Generic message, log details
func handleError(w http.ResponseWriter, err error) {
    log.Error("Internal error", "error", err)  // Log for debugging
    http.Error(w, "An error occurred", http.StatusInternalServerError)
}
```

### Fail Secure

```go
// VULNERABLE - Fails open
func isAuthorized(user *User, resource string) bool {
    perms, err := getPermissions(user.ID)
    if err != nil {
        return true  // WRONG: Error = allow access
    }
    return perms.CanAccess(resource)
}

// SECURE - Fails closed
func isAuthorized(user *User, resource string) bool {
    perms, err := getPermissions(user.ID)
    if err != nil {
        log.Error("Permission check failed", "error", err)
        return false  // Error = deny access
    }
    return perms.CanAccess(resource)
}
```

---

## Cryptography

### Use Standard Libraries

```go
import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
)

// Generate secure random token
func generateToken(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// Hash sensitive data
func hashData(data string) string {
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

### Don't Roll Your Own Crypto

```go
// WRONG - Custom "encryption"
func encrypt(data string) string {
    result := ""
    for _, c := range data {
        result += string(c + 1)  // Caesar cipher is not encryption!
    }
    return result
}

// RIGHT - Use standard library
import "crypto/aes"
// Use established packages for encryption
```

---

## File Operations Security

### Path Traversal Prevention

```go
// VULNERABLE - Path traversal
func serveFile(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    http.ServeFile(w, r, "/uploads/" + filename)
    // Attacker: ?file=../../../etc/passwd
}

// SECURE - Validate and clean path
func serveFile(w http.ResponseWriter, r *http.Request) {
    filename := filepath.Base(r.URL.Query().Get("file"))  // Remove path
    
    fullPath := filepath.Join("/uploads", filename)
    
    // Ensure still within allowed directory
    if !strings.HasPrefix(fullPath, "/uploads/") {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }
    
    http.ServeFile(w, r, fullPath)
}
```

### File Upload Security

```go
func handleUpload(w http.ResponseWriter, r *http.Request) {
    // Limit file size
    r.Body = http.MaxBytesReader(w, r.Body, 10<<20)  // 10MB
    
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "File too large", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Validate content type
    buffer := make([]byte, 512)
    file.Read(buffer)
    contentType := http.DetectContentType(buffer)
    
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }
    
    if !allowedTypes[contentType] {
        http.Error(w, "Invalid file type", http.StatusBadRequest)
        return
    }
    
    // Generate safe filename
    ext := filepath.Ext(header.Filename)
    safeFilename := uuid.New().String() + ext
    
    // Save file...
}
```

---

## Security Tools

### gosec

```bash
# Install
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run
gosec ./...

# Exclude rules
gosec -exclude=G104 ./...
```

### Common gosec Findings

| Rule | Issue | Fix |
|------|-------|-----|
| G101 | Hardcoded credentials | Use environment variables |
| G102 | Bind to all interfaces | Bind to specific interface |
| G104 | Unhandled errors | Check all errors |
| G201 | SQL string formatting | Use parameterized queries |
| G401 | Weak crypto (MD5/SHA1) | Use SHA256+ |
| G501 | Blacklisted imports | Remove unsafe packages |

---

## Security Checklist

Before committing:
- [ ] No hardcoded secrets
- [ ] All inputs validated
- [ ] Parameterized queries used
- [ ] Errors don't leak information
- [ ] Authorization checked at every endpoint
- [ ] gosec passes

Before release:
- [ ] Dependencies scanned for vulnerabilities
- [ ] Security headers configured
- [ ] TLS/HTTPS enforced
- [ ] Rate limiting in place
- [ ] Logging captures security events

---

## Related Skills

- `security` - General security practices
- `go-expert` - Go-specific security patterns
- `code-reviewer` - Security review focus
- `error-handling` - Secure error handling
