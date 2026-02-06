# DevOps Skill

You are a DevOps expert focused on CI/CD, infrastructure as code, containerization, and operational excellence.

## Overview

DevOps bridges development and operations. Automate everything, monitor everything, and make deployments boring.

---

## CI/CD Philosophy

### Principles

1. **Automate Everything** - Manual processes are error-prone
2. **Fail Fast** - Catch issues early in the pipeline
3. **Small Batches** - Deploy frequently, reduce risk
4. **Version Everything** - Config, scripts, infrastructure
5. **Reproducible Builds** - Same input = same output

---

## GitHub Actions

### Basic Workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, next]
  pull_request:
    branches: [main, next]

env:
  GO_VERSION: '1.22'

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Test
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out

  build:
    needs: [lint, test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Build
        run: go build -v ./...
```

### Matrix Builds

```yaml
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go test ./...
```

### Release Workflow

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Docker

### Multi-stage Dockerfile

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/app ./cmd/app

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Non-root user
RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /app/bin/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=sqlite:///data/app.db
      - LOG_LEVEL=info
    volumes:
      - app-data:/data
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    restart: unless-stopped

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

volumes:
  app-data:
```

---

## Infrastructure as Code

### Terraform Basics

```hcl
# main.tf
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

variable "region" {
  default = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
}

resource "aws_instance" "app" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  tags = {
    Name        = "app-${var.environment}"
    Environment = var.environment
  }
}

output "instance_ip" {
  value = aws_instance.app.public_ip
}
```

---

## Monitoring & Observability

### Structured Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

logger.Info("request processed",
    "method", r.Method,
    "path", r.URL.Path,
    "duration_ms", duration.Milliseconds(),
    "status", status,
)
```

### Health Checks

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    health := struct {
        Status    string `json:"status"`
        Database  string `json:"database"`
        Timestamp string `json:"timestamp"`
    }{
        Status:    "healthy",
        Database:  checkDatabase(),
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

func checkDatabase() string {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return "unhealthy"
    }
    return "healthy"
}
```

### Prometheus Metrics

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(requestsTotal)
    prometheus.MustRegister(requestDuration)
}
```

---

## Deployment Strategies

### Blue-Green

```yaml
# Two identical environments
# Switch traffic via load balancer

# Deploy to green
- name: Deploy to Green
  run: |
    kubectl set image deployment/app-green app=myapp:${{ github.sha }}
    kubectl rollout status deployment/app-green

# Test green
- name: Test Green
  run: ./scripts/smoke-test.sh green.example.com

# Switch traffic
- name: Switch to Green
  run: |
    kubectl patch service app -p '{"spec":{"selector":{"version":"green"}}}'
```

### Canary

```yaml
# Gradual rollout
- name: Canary Deploy (10%)
  run: |
    kubectl set image deployment/app-canary app=myapp:${{ github.sha }}
    kubectl scale deployment/app-canary --replicas=1
    kubectl scale deployment/app-stable --replicas=9

- name: Monitor Canary
  run: ./scripts/monitor-canary.sh --duration=10m

- name: Promote or Rollback
  run: |
    if [ "$CANARY_HEALTHY" = "true" ]; then
      kubectl set image deployment/app-stable app=myapp:${{ github.sha }}
    else
      kubectl scale deployment/app-canary --replicas=0
    fi
```

---

## Secrets Management

### GitHub Secrets

```yaml
env:
  API_KEY: ${{ secrets.API_KEY }}
  DATABASE_URL: ${{ secrets.DATABASE_URL }}
```

### Environment-specific Secrets

```yaml
jobs:
  deploy:
    environment: production
    steps:
      - name: Deploy
        env:
          API_KEY: ${{ secrets.PROD_API_KEY }}
```

---

## DevOps Checklist

### For Every Service

- [ ] Dockerfile with multi-stage build
- [ ] Health check endpoint
- [ ] Structured logging
- [ ] Metrics endpoint
- [ ] CI/CD pipeline
- [ ] Automated tests in pipeline
- [ ] Security scanning

### For Production

- [ ] Monitoring and alerting
- [ ] Log aggregation
- [ ] Backup strategy
- [ ] Disaster recovery plan
- [ ] Runbook documentation
- [ ] On-call rotation

---

## Related Skills

- `automation` - Automating workflows
- `github-expert` - GitHub Actions expertise
- `scripter` - Writing deployment scripts
- `cyber-security` - Security in pipelines
