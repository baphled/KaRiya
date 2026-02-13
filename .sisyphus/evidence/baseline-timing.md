# Baseline Test Performance Measurements

**Date:** 2026-02-13  
**Commit:** Current HEAD  
**Environment:** Linux, Go 1.x, Ginkgo 2.28.1

## Unit Tests

- **Total Time:** 0.133 seconds
- **Status:** ✅ PASS
- **Specs Run:** 56 of 56
- **Passed:** 56
- **Failed:** 0
- **Pending:** 0
- **Skipped:** 0

### Summary
Unit tests are fast and stable. No performance concerns at this tier.

---

## BDD Tests (Acceptance Tests)

- **Total Time:** 163.801 seconds (~2 minutes 44 seconds)
- **Status:** ❌ FAIL (1 failure)
- **Scenarios Run:** 70 total
- **Passed:** 69
- **Failed:** 1
- **Pending:** 0

### Failure Details
- **Scenario:** `Event_with_all_available_tags`
- **Error:** Expected 1 events, got 0 (test setup issue, not performance)

### Performance Analysis

**Current State:**
- BDD tests consume **99.8%** of total test time
- Unit tests consume **0.2%** of total test time
- **Target:** Reduce BDD from 163.8s to <120s (26% reduction needed)

---

## Top 10 Slowest BDD Scenarios

| Rank | Scenario Name | Time | Category |
|------|---------------|------|----------|
| 1 | Accept_inferred_skill_during_review | 25.66s | Skill inference + review |
| 2 | Quick_capture_with_minimal_input | 24.66s | Event capture + burst detection |
| 3 | Edit_suggested_burst_before_accepting | 17.16s | Burst editing + review |
| 4 | Accept_suggested_burst_during_review | 17.14s | Burst acceptance + review |
| 5 | Reject_all_suggestions_and_submit_raw_event | 16.14s | Event submission |
| 6 | Edit_event_metadata_during_review | 13.14s | Event editing + review |
| 7 | Add_event_from_timeline | 11.68s | Event creation |
| 8 | Edit_event_metadata_and_save | 11.17s | Event editing |
| 9 | Edit_burst_description_and_save | 11.13s | Burst editing |
| 10 | Edit_burst_name_and_save | 6.62s | Burst editing |

**Cumulative Time (Top 10):** 154.70 seconds (94.4% of total BDD time)

### Observations

1. **Skill inference scenarios dominate** - Top 2 scenarios (50.32s) involve skill inference
2. **Event capture is expensive** - Scenarios 2, 5, 6, 7 involve event capture/review (68.12s)
3. **Database setup overhead** - Each scenario creates full SQLite DB (6 migrations per scenario)
4. **Remaining 60 scenarios** - Only 9.1 seconds combined (fast, mostly UI navigation)

---

## Database Setup Overhead

Each BDD scenario runs:
```
001_create_career_events.sql    ~0.3-3.5ms
002_add_categories_column.sql   ~0.2-0.6ms
003_create_bursts.sql           ~0.2-0.7ms
004_create_facts.sql            ~0.3-1.1ms
005_create_skills.sql           ~0.3-30ms (variable)
006_create_event_skills.sql     ~0.3-6ms
```

**Total per scenario:** ~2-40ms (mostly negligible vs. scenario time)

---

## Recommendations for Optimization

### High-Impact Opportunities (Target: 26% reduction = 43.8s savings)

1. **Parallel BDD Execution** (Est. 40-50% reduction)
   - Run independent scenarios concurrently
   - Potential savings: 65-82 seconds

2. **Shared Database Setup** (Est. 5-10% reduction)
   - Reuse database across scenarios instead of recreating
   - Potential savings: 8-16 seconds

3. **Skill Inference Caching** (Est. 10-15% reduction)
   - Cache inference results for repeated patterns
   - Potential savings: 16-24 seconds

4. **Event Capture Optimization** (Est. 5-10% reduction)
   - Reduce form interaction delays
   - Potential savings: 8-16 seconds

### Quick Wins (Low effort, measurable impact)

- Profile skill inference service (likely bottleneck)
- Check for unnecessary database queries in review flows
- Verify form rendering isn't blocking on I/O

---

## Next Steps

1. Profile top 3 slowest scenarios to identify exact bottlenecks
2. Implement parallel execution for independent scenarios
3. Measure impact of each optimization
4. Re-run baseline after each change to track progress
