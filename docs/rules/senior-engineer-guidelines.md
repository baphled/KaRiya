---
id: senior-engineer-context-guidelines
aliases:
  - Senior Engineer Context
  - Prompt Guidelines for Engineers
  - Engineering Workflow Checklist
created: 2025-07-28T23:39
lead: Senior Engineer Context & Prompt Guidelines
modified: 2025-12-16T20:37
template-type: Note
---

# Senior‑Engineer Context & Prompt‑Guidelines

A concise, language‑agnostic checklist to guide LLM‑driven “senior engineer”
workflows. Embed these headings and keywords in your prompts to ensure SOLID
design, TDD safety, refactoring rigor, automation, observability, and more.

---

## 1. Principles & Patterns
- **SOLID**
  - **S**ingle Responsibility
  - **O**pen/Closed
  - **L**iskov Substitution
  - **I**nterface Segregation
  - **D**ependency Inversion
- **Clean Code**: KISS, DRY, YAGNI
- **Design Patterns**: Strategy, Factory, Observer, etc.
- **Architectural Decisions**: ADR → context, options, decision, consequences

## 2. Workflow: Red–Green–Refactor
1. **Small Change** → scope one behavior/bug
2. **Write/Update Test** (Red) → failing test first
3. **Implement Minimal Code** (Green) → satisfy test only
4. **Run Tests & CI** → lint, static analysis, coverage, mutation
5. **Refactor** → remove smells (Extract Method, Rename, etc.)
6. **Atomic Commit** → descriptive message; push triggers CI
   - See [Atomic Commit Guidelines](./atomic-commits.md) for detailed guidance
   - Use `make review-commit` before finalizing commits

## 3. Documentation & Communication
- **Docstrings & Comments**: intent‑first, public interfaces
- **README / Runbook**: setup, CI badges, troubleshooting
- **ADR Updates**: maintain decision log

## 4. Error‑Handling & Resilience
- **Fail‑Fast** validations; descriptive exceptions
- **Graceful Degradation** for optional features
- **Retry/Circuit‑Breaker** around unstable calls

## 5. Observability & Monitoring
- **Logging & Tracing**: structured logs, trace spans
- **Metrics & Alerts**: rate, latency, error counts
- **Health Endpoints**: `/healthz`, `/readyz`

## 6. Performance & Scalability
- **Benchmarking & Profiling** → track regressions
- **Caching**: TTL, invalidation patterns
- **Concurrency Controls**: bulkheads, rate limiters

## 7. Security & Compliance
- **Input Validation/Sanitization**
- **Secret Management**: vault/env‑vars only
- **Dependency Scanning**: SAST/DAST, CVE checks

## 8. Release & Dependency Management
- **Semantic Versioning** (MAJOR.MINOR.PATCH)
- **Feature Toggles & Canary Releases**
- **Changelog Automation** from commits

## 9. Developer Experience
- **Branching Model**: GitHub Flow/Git Flow, `feature/…`, `fix/…`
- **Onboarding Scripts**: `make setup`, `npm run setup`
- **Pre‑commit Hooks**: lint, format, tests, static checks
- **Watch Mode**: auto‑rerun tests on save

## 10. Continuous Learning & Improvement
- **Code Reviews**: clarity, patterns, coverage, SOLID
- **Refactoring Backlog**: schedule cleanup sprints
- **Post‑Mortem & RCA**: blameless incident analysis
- **Tech Radar & Spikes**: log experiments, update radar

---

## Prompt‑Engineering Keywords

“Generate minimal change”
“Use Red‑Green‑Refactor cycle”
“Follow SOLID principles”
“Include/update failing test first”
“Ensure CI triggers on file save”
“Provide atomic commit message”
“Refactor code smells”
“Maintain 100% branch coverage”
“Update ADR”
“Instrument logs and metrics”
“Use semantic commit”
“Implement health checks”
“Add secret vault integration”

Embed this checklist into your LLM prompts to drive robust, production‑grade outputs at every step.

