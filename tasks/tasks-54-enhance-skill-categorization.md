# Task 54: Enhance Skill Categorization Coverage

## Problem

84% of imported skills (257 of 305) are categorized as "other". The keyword
dictionary in `keywords.go` only covers ~224 exact technology tool/framework
names. Imported skills include practices, architectural patterns, security
concepts, and engineering methodologies that have no keyword entries.

### Original Distribution (before any work)

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

### Post-Phase 1-4 Distribution (after first recategorization run)

| Category     | Count | % of Total |
|--------------|------:|----------:|
| other        |   167 |     54.8% |
| backend      |    20 |      6.6% |
| devops       |    15 |      4.9% |
| frontend     |    15 |      4.9% |
| practices    |    14 |      4.6% |
| database     |    13 |      4.3% |
| architecture |    10 |      3.3% |
| testing      |    10 |      3.3% |
| monitoring   |     9 |      3.0% |
| security     |     9 |      3.0% |
| tooling      |     9 |      3.0% |
| data         |     8 |      2.6% |
| ml           |     4 |      1.3% |
| cloud        |     2 |      0.7% |

**167 skills remain "other".** Analysis reveals two root causes:

1. **Exact-match-only algorithm** — `GetCategoryForSkillName` requires the
   entire skill name to match a keyword exactly after lowercasing. "ELK Stack"
   does not match keyword `"elk"`, "Code Reviews" does not match `"code review"`.
2. **Missing keywords for compound phrases** — Skills like "Backend Development",
   "Service Architecture", "Component Library" are multi-word phrases absent
   from the dictionary.

### Breakdown of 167 Remaining "Other" Skills

| Bucket | Count | Fix Required |
|--------|------:|--------------|
| Exact match exists, missed due to bug | 6 | Add missing keywords (`timescaledb`, `kibana`, `logstash`) |
| Substring match would work | 9 | Improve matching algorithm (contains fallback) |
| Need new keywords | 124 | Add ~124 keyword entries for compound phrases |
| Genuinely uncategorizable | 28 | Correct — remain as "other" |

## Goal

Reduce "other" from 84% to ~9% by:

1. Adding 3 new `SkillCategory` constants: `architecture`, `security`,
   `practices`.
2. Expanding the keyword dictionary with ~350 entries (initial) + ~130 more.
3. Adding a `--recategorize-skills` CLI flag to update existing DB skills.
4. Improving the matching algorithm with substring/contains fallback.
5. Adding ~124 new compound-phrase keyword entries.

## Phases

### Phase 1: New SkillCategory Constants

**Status:** DONE (commit `49bc64e4`)

**Files:**

- `internal/constants/constants.go`
- `internal/constants/constants_test.go`

**Changes:**

- Added `SkillCategoryArchitecture`, `SkillCategorySecurity`, `SkillCategoryPractices`.
- Added `AllSkillCategories()`, `IsValidSkillCategory()`, `SkillCategoryStrings()`.
- Total categories: 15 (was 8).

### Phase 2: Initial Keyword Dictionary

**Status:** DONE (commit `05fa5227`)

**Files:**

- `internal/service/career/technology/keywords.go`
- `internal/service/career/technology/keywords_test.go`

**Changes:**

- Created `Entry` struct, `Keywords` slice (350 entries), `GetCategoryForSkillName()`.
- OAuth moved from `tooling` to `security`.
- 44 tests covering dictionary structure, coverage per category, lookup behavior.

### Phase 3: Recategorization CLI Flag

**Status:** DONE (commit `6a16ca04`)

**Files:**

- `cmd/cli/main.go`
- `cmd/cli/main_test.go`
- `internal/service/career/technology/keywords.go` (added `RecategorizeSkills`)

**Changes:**

- Added `RecategorizeResult` struct and `RecategorizeSkills(ctx, repo)` function.
- Added `--recategorize-skills` CLI flag with summary output.
- 6 unit tests + 3 CLI tests.

### Phase 4: Update CV Profile Inference

**Status:** DONE (commit `eabaaf94`)

**Files:**

- `internal/service/career/cv/profile_inference.go`
- `internal/service/career/cv/profile_inference_test.go`

**Changes:**

- Routed `architecture` and `security` to Systems bucket in `InferTechnologies()`.
- Added expertise checks for architecture and security in `InferCoreStrengths()`.
- 4 new tests.

### Phase 5: Improve Matching Algorithm

**Status:** IN PROGRESS

**Files:**

- `internal/service/career/technology/keywords.go`
- `internal/service/career/technology/keywords_test.go`

**Problem:** `GetCategoryForSkillName` uses exact match only. "ELK Stack"
lowercases to `"elk stack"` but the keyword is `"elk"`.

**Fix:** Add a substring/contains fallback after exact match fails. Must
guard against short keywords (`"c"`, `"r"`, `"go"`, `"ai"`) matching inside
unrelated words — only apply substring matching for keywords >= 3 characters,
or use word-boundary matching.

**Also add missing exact keywords:** `timescaledb`, `kibana`, `logstash`,
`bubble tea`, `lipgloss`, `cms`, `websocket` (singular).

### Phase 6: Expand Keyword Dictionary for Compound Phrases

**Status:** NOT STARTED

**Files:**

- `internal/service/career/technology/keywords.go`
- `internal/service/career/technology/keywords_test.go`

**~124 new keywords needed for compound phrases like:**

| Category     | Count | Examples |
|--------------|------:|---------|
| architecture | ~25   | backend architecture, component architecture, service architecture, asynchronous processing, real-time systems, scalable systems |
| practices    | ~55   | agile delivery, code standards, collaboration, consulting, delivery management, engineering practices, project management, software engineering |
| backend      | ~15   | backend development, backend engineering, crud, error handling, full-stack development, web development, serialization |
| frontend     | ~15   | component library, form components, keyboard navigation, layout design, modal design, ux, user experience |
| devops       | ~10   | firmware development, incident response, operational support, production support, reliability engineering |
| tooling      | ~2    | diagramming, technical writing |
| data         | ~1    | streaming |

### Phase 7: Final Validation

**Status:** NOT STARTED

- `go test ./... -count=1` passes.
- `staticcheck ./...` no new warnings.
- Run `kariya --recategorize-skills` against actual DB.
- Verify distribution: `sqlite3 ~/.kariya/events.db "SELECT category, COUNT(*) FROM skills GROUP BY category ORDER BY count DESC;"`.
- Target: ~28 skills remain "other" (genuinely uncategorizable).

## Expected Final Outcome

| Category     | Before | After Phase 4 | After Phase 7 |
|--------------|-------:|--------------:|--------------:|
| other        |    257 |           167 |           ~28 |
| practices    |      0 |            14 |           ~69 |
| architecture |      0 |            10 |           ~35 |
| backend      |      7 |            20 |           ~35 |
| frontend     |      9 |            15 |           ~30 |
| devops       |      7 |            15 |           ~25 |
| database     |      8 |            13 |           ~13 |
| testing      |      4 |            10 |           ~10 |
| monitoring   |      3 |             9 |           ~12 |
| security     |      0 |             9 |            ~9 |
| tooling      |      8 |             9 |           ~11 |
| data         |      0 |             8 |            ~9 |
| ml           |      0 |             4 |            ~4 |
| cloud        |      2 |             2 |            ~2 |

## Commit Strategy

7 atomic commits:

1. `feat(domain): add architecture, security, practices skill categories` — DONE
2. `feat(service): add keyword dictionary with 350 entries` — DONE
3. `feat(cli): add --recategorize-skills flag` — DONE
4. `feat(cv): update profile inference for new categories` — DONE
5. `fix(service): add substring matching fallback to keyword lookup`
6. `feat(service): expand keyword dictionary with compound phrase entries`
7. `docs(tasks): update task 54 status after final validation`

## Risks & Notes

- **OAuth move:** Moving OAuth from `tooling` to `security` affects 1 existing
  skill in the DB. The recategorization pass handles this.
- **No schema migration needed:** The `category` column is a free-form string,
  not an enum. New values work immediately.
- **CV rendering:** `buildGroupedSkills` dynamically groups by category string.
  New categories appear automatically with no renderer changes.
- **Short keyword safety:** Substring matching must not allow 1-2 char keywords
  (`"c"`, `"r"`, `"go"`, `"ai"`) to match inside unrelated words. Use a minimum
  keyword length threshold or word-boundary matching.
- **~28 skills remain "other":** These are genuinely uncategorizable (soft
  skills, domain knowledge, niche tech like IRC/WAP/ShoutCast). Acceptable.
- **Dependency:** Builds on Phase 14 of Task 47 (skill category normalization)
  which unified the 12 canonical categories.
