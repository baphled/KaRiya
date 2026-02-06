---
name: configuration-management
description: Manage configuration properly - environment variables, config files, feature flags, secrets
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide proper configuration management - separating config from code, handling secrets safely, managing feature flags, and supporting different environments.

## When to use me

- Adding configurable behaviour
- Setting up environment-specific settings
- Managing secrets and sensitive data
- Implementing feature flags
- Reviewing configuration practices

## Core Principles

1. **Config separate from code** - Never hardcode environment-specific values
2. **Secrets are special** - Never commit secrets, use proper secret management
3. **Fail fast** - Validate config at startup, not at use time
4. **Document defaults** - Make expected configuration clear

## Configuration Hierarchy

```
Priority (highest to lowest):
1. Command-line flags
2. Environment variables
3. Config file (environment-specific)
4. Config file (default)
5. Hardcoded defaults
```

## Environment Variables

### Naming Convention

```bash
# Format: APPNAME_SECTION_KEY
KARIYA_DATABASE_URL=postgres://...
KARIYA_DATABASE_MAX_CONNECTIONS=10
KARIYA_LOG_LEVEL=debug
KARIYA_FEATURE_NEW_UI=true
```

### Loading in Go

```go
type Config struct {
    Database DatabaseConfig
    Log      LogConfig
    Features FeatureFlags
}

type DatabaseConfig struct {
    URL            string `env:"KARIYA_DATABASE_URL" required:"true"`
    MaxConnections int    `env:"KARIYA_DATABASE_MAX_CONNECTIONS" default:"10"`
}

type LogConfig struct {
    Level  string `env:"KARIYA_LOG_LEVEL" default:"info"`
    Format string `env:"KARIYA_LOG_FORMAT" default:"json"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    if err := envconfig.Process("kariya", cfg); err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }
    
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("validate config: %w", err)
    }
    
    return cfg, nil
}
```

### Validation at Startup

```go
func (c *Config) Validate() error {
    var errs []error
    
    if c.Database.URL == "" {
        errs = append(errs, errors.New("KARIYA_DATABASE_URL is required"))
    }
    
    if c.Database.MaxConnections < 1 {
        errs = append(errs, errors.New("KARIYA_DATABASE_MAX_CONNECTIONS must be >= 1"))
    }
    
    validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLevels[c.Log.Level] {
        errs = append(errs, fmt.Errorf("invalid log level: %s", c.Log.Level))
    }
    
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    
    return nil
}
```

## Config Files

### YAML Structure

```yaml
# config/default.yaml
database:
  max_connections: 10
  timeout: 30s

log:
  level: info
  format: json

features:
  new_ui: false
  beta_export: false

# config/development.yaml (overrides default)
database:
  max_connections: 5

log:
  level: debug
  format: text

features:
  new_ui: true
  beta_export: true
```

### Loading with Overrides

```go
func LoadConfigFile(env string) (*Config, error) {
    cfg := &Config{}
    
    // Load defaults
    if err := loadYAML("config/default.yaml", cfg); err != nil {
        return nil, err
    }
    
    // Override with environment-specific
    envFile := fmt.Sprintf("config/%s.yaml", env)
    if exists(envFile) {
        if err := loadYAML(envFile, cfg); err != nil {
            return nil, err
        }
    }
    
    // Environment variables override files
    if err := envconfig.Process("kariya", cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}
```

## Secrets Management

### Never Commit Secrets

```bash
# .gitignore
.env
.env.*
*.pem
*.key
secrets/
```

### Use Secret Managers

```go
// Load from secret manager in production
func LoadSecrets(ctx context.Context) (*Secrets, error) {
    if os.Getenv("ENV") == "production" {
        return loadFromVault(ctx)
    }
    
    // Development: load from .env file
    return loadFromEnvFile(".env")
}
```

### Secret Rotation

```go
// Design for rotation - don't cache secrets forever
type SecretProvider interface {
    GetDatabasePassword(ctx context.Context) (string, error)
    GetAPIKey(ctx context.Context) (string, error)
}

// Refresh periodically
func (p *VaultProvider) GetDatabasePassword(ctx context.Context) (string, error) {
    p.mu.RLock()
    if time.Since(p.lastFetch) < 5*time.Minute && p.cachedPassword != "" {
        defer p.mu.RUnlock()
        return p.cachedPassword, nil
    }
    p.mu.RUnlock()
    
    // Fetch fresh
    return p.fetchAndCache(ctx, "database/password")
}
```

## Feature Flags

### Simple Boolean Flags

```go
type FeatureFlags struct {
    NewUI      bool `env:"KARIYA_FEATURE_NEW_UI" default:"false"`
    BetaExport bool `env:"KARIYA_FEATURE_BETA_EXPORT" default:"false"`
}

func (i *Intent) View() string {
    if i.features.NewUI {
        return i.renderNewUI()
    }
    return i.renderClassicUI()
}
```

### Percentage Rollouts

```go
type FeatureConfig struct {
    Name       string
    Enabled    bool
    Percentage int  // 0-100
}

func (f *FeatureFlags) IsEnabled(feature string, userID string) bool {
    cfg := f.configs[feature]
    if !cfg.Enabled {
        return false
    }
    
    if cfg.Percentage >= 100 {
        return true
    }
    
    // Consistent hash for user
    hash := fnv.New32a()
    hash.Write([]byte(userID + feature))
    return int(hash.Sum32()%100) < cfg.Percentage
}
```

### Flag Lifecycle

```
1. DEVELOPMENT: Flag off by default
2. TESTING: Flag on in test/staging
3. ROLLOUT: Gradual percentage increase
4. STABLE: Flag on by default
5. CLEANUP: Remove flag, code is permanent
```

## Environment-Specific Configuration

### Example .env Files

```bash
# .env.example (committed - template)
KARIYA_DATABASE_URL=postgres://user:pass@localhost/kariya_dev
KARIYA_LOG_LEVEL=debug
KARIYA_FEATURE_NEW_UI=true

# .env.development (not committed)
KARIYA_DATABASE_URL=postgres://dev:dev@localhost/kariya_dev

# .env.test (not committed)
KARIYA_DATABASE_URL=postgres://test:test@localhost/kariya_test
KARIYA_LOG_LEVEL=error

# .env.production (never on disk - from secret manager)
```

### Detecting Environment

```go
func GetEnvironment() string {
    env := os.Getenv("KARIYA_ENV")
    if env == "" {
        env = os.Getenv("ENV")
    }
    if env == "" {
        env = "development"
    }
    return env
}

func IsProduction() bool {
    return GetEnvironment() == "production"
}

func IsDevelopment() bool {
    env := GetEnvironment()
    return env == "development" || env == "dev"
}
```

## Configuration Anti-Patterns

### DON'T: Hardcode Values

```go
// BAD
const maxConnections = 10
const apiEndpoint = "https://api.example.com"

// GOOD
cfg.Database.MaxConnections
cfg.API.Endpoint
```

### DON'T: Use Config at Import Time

```go
// BAD - Fails if env not set at import
var dbURL = os.Getenv("DATABASE_URL")  // Package-level

// GOOD - Load at runtime
func NewDB(cfg *Config) (*DB, error) {
    return sql.Open("postgres", cfg.Database.URL)
}
```

### DON'T: Scatter Config Access

```go
// BAD - Config accessed everywhere
func (s *Service) DoThing() {
    if os.Getenv("FEATURE_X") == "true" {  // Direct env access
        // ...
    }
}

// GOOD - Inject config
func NewService(cfg *Config) *Service {
    return &Service{
        featureX: cfg.Features.X,
    }
}
```

## Testing with Config

```go
func TestService(t *testing.T) {
    // Use test configuration
    cfg := &Config{
        Database: DatabaseConfig{
            URL: "postgres://test@localhost/test",
        },
        Features: FeatureFlags{
            NewUI: true,  // Test with feature enabled
        },
    }
    
    svc := NewService(cfg)
    // ... test
}
```

## Related Skills

- `security` - Secure handling of secrets
- `devops` - Deployment configuration
- `error-handling` - Config validation errors
