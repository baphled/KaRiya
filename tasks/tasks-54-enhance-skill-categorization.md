# Task 54: Enhance Skill Categorization Coverage

## Problem

84% of imported skills (257 of 305) are categorized as "other". The keyword
dictionary in `keywords.go` only covers ~224 exact technology tool/framework
names. Imported skills include practices, architectural patterns, security
concepts, and engineering methodologies that have no keyword entries.

### Current Distribution

| Category   | Count | % of Total |
|------------|------:|----------:|
| other      |   257 |     84.3% |
| frontend   |     9 |      3.0% |
| database   |     8 |      2.6% |
| tooling    |     8 |      2.6% |
| backend    |     7 |      2.3% |
| devops     |     7 |      2.3% |
| testing    |     4 |      1.3% |
| monitoring |     3 |      1.0% |
| cloud      |     2 |      0.7% |
| mobile     |     0 |      0.0% |
| data       |     0 |      0.0% |
| ml         |     0 |      0.0% |

## Goal

Reduce "other" from 84% to ~28% by:

1. Adding 3 new `SkillCategory` constants: `architecture`, `security`,
   `practices`.
2. Expanding the keyword dictionary with ~130 new entries.
3. Adding a `--recategorize-skills` CLI flag to update existing DB skills.

## Phases

### Phase 1: New SkillCategory Constants

**Status:** NOT STARTED

**Files:**

- `internal/constants/constants.go`
- `internal/constants/constants_test.go`

**Changes:**

- Add `SkillCategoryArchitecture SkillCategory = "architecture"`.
- Add `SkillCategorySecurity SkillCategory = "security"`.
- Add `SkillCategoryPractices SkillCategory = "practices"`.
- Update `AllSkillCategories()` to return 15 (currently 12).
- `IsValidSkillCategory()` works automatically via `AllSkillCategories`.

**No other files need changes.** Domain validation in `skill.go` already calls
`constants.IsValidSkillCategory()`.

### Phase 2: Expand Keyword Dictionary

**Status:** NOT STARTED

**Files:**

- `internal/service/career/technology/keywords.go`
- `internal/service/career/technology/keywords_test.go`

**New keyword sections (~130 entries):**

| Category       | Count | Examples                                                     |
|----------------|------:|--------------------------------------------------------------|
| architecture   |   ~30 | microservices, soa, domain-driven design, distributed systems, system design, scalability, dependency injection, state machine |
| security       |   ~12 | security, authentication, tls, pki, pci, saml, active directory, compliance, encryption |
| practices      |   ~30 | agile, tdd, bdd, code review, pair programming, refactoring, code quality, debugging, performance optimization, feature flags |
| backend        |   ~15 | c, c++, coffeescript, zend, wordpress, sidekiq, ajax, soap, api, api design, websockets |
| frontend       |    ~8 | backbonejs, responsive design, accessibility, storybook, seo, design systems, visual design |
| devops         |   ~10 | deployment, infrastructure, server management, system administration, shell, cron, git flow, git hooks, operations, networking |
| database       |    ~6 | nosql, database migrations, database optimization, database performance, data modeling, data integrity |
| testing        |    ~8 | ginkgo, gomega, e2e testing, automated testing, testing, code coverage, quality assurance |
| monitoring     |    ~7 | logging, structured logging, alerting, metrics, observability, monitoring, elk |
| data           |    ~7 | etl, data processing, data export, data platforms, analytics, business intelligence, reporting |
| ml             |    ~4 | ai, generative ai, llms, prompt engineering |
| tooling        |    ~6 | documentation, technical documentation, markdown, xml, json, csv |

**Decision:** Move OAuth from `tooling` to `security` (authentication protocol).

### Phase 3: Recategorization CLI Flag

**Status:** NOT STARTED

**Files:**

- `cmd/cli/main.go`
- `cmd/cli/main_test.go`
- `internal/service/career/technology/keywords.go` (add `RecategorizeSkills`)

**CLI usage:**

```
kariya --recategorize-skills                    # Default DB
kariya --recategorize-skills --db ./events.db   # Custom DB
```

**Logic for `RecategorizeSkills(ctx, repo)`:**

1. Load all skills via `repo.List(ctx, nil)`.
2. For each skill call `GetCategoryForSkillName(skill.Name)`.
3. If result is non-empty AND differs from current category, update via
   `repo.Update(ctx, skill)`.
4. Print summary: `Recategorized N skills (X backend, Y architecture, ...)`.
5. Skip skills where `GetCategoryForSkillName` returns empty (no match).

Exits after recategorization (does not launch TUI).

### Phase 4: Update CV Profile Inference

**Status:** NOT STARTED

**Files:**

- `internal/service/career/cv/profile_inference.go`

Add awareness of new categories in the skill categorization switch:

- `architecture` counts toward technical/backend weight.
- `security` counts toward technical weight.
- `practices` counts toward general technical weight.

### Phase 5: Validation & Cleanup

**Status:** NOT STARTED

- `go test ./... -count=1` passes.
- `staticcheck ./...` no new warnings.
- `make check-compliance` passes (minus pre-existing violations).
- Run `kariya --recategorize-skills` against actual DB.
- Verify: `sqlite3 ~/.kariya/events.db "SELECT category, COUNT(*) FROM skills GROUP BY category ORDER BY count DESC;"`.

## Expected Outcome

| Category        | Before | After |
|-----------------|-------:|------:|
| other           |    257 |   ~85 |
| practices (new) |      0 |   ~34 |
| architecture (new) |   0 |   ~30 |
| backend         |      7 |   ~22 |
| frontend        |      9 |   ~17 |
| devops          |      7 |   ~17 |
| database        |      8 |   ~14 |
| tooling         |      8 |   ~13 |
| testing         |      4 |   ~12 |
| security (new)  |      0 |   ~10 |
| monitoring      |      3 |   ~10 |
| data            |      0 |    ~7 |
| ml              |      0 |    ~4 |
| cloud           |      2 |    ~3 |

## Commit Strategy

4-5 atomic commits:

1. `feat(constants): add architecture, security, practices skill categories`
2. `feat(technology): expand keyword dictionary with ~130 entries`
3. `feat(cli): add --recategorize-skills flag for existing skill migration`
4. `feat(cv): update profile inference for new skill categories`
5. `test: validate full recategorization flow`

## Risks & Notes

- **OAuth move:** Moving OAuth from `tooling` to `security` affects 1 existing
  skill in the DB. The recategorization pass handles this.
- **No schema migration needed:** The `category` column is a free-form string,
  not an enum. New values work immediately.
- **CV rendering:** `buildGroupedSkills` dynamically groups by category string.
  New categories appear automatically with no renderer changes.
- **~85 skills remain "other":** These are genuinely uncategorizable (soft
  skills, domain knowledge, niche tech like IRC/WAP/ShoutCast). Acceptable.
- **Dependency:** Builds on Phase 14 of Task 47 (skill category normalization)
  which unified the 12 canonical categories.
