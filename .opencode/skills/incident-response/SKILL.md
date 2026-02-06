---
name: incident-response
description: Handle production incidents effectively - diagnose, mitigate, resolve, and learn from failures
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide effective incident response - from initial detection through resolution and post-mortem. Minimise impact, restore service, and prevent recurrence.

## When to use me

- Production issues reported
- Alerts firing
- User-reported problems
- System degradation
- Post-incident analysis

## Core Principles

1. **Mitigate first** - Stop the bleeding before finding root cause
2. **Communicate early** - Keep stakeholders informed
3. **Document as you go** - You'll need this for post-mortem
4. **Blameless culture** - Focus on systems, not people
5. **Learn and improve** - Every incident is an opportunity

## Incident Severity Levels

| Level | Impact | Response |
|-------|--------|----------|
| **SEV1** | Complete outage, data loss risk | All hands, immediate |
| **SEV2** | Major feature broken, many users affected | Primary on-call + backup |
| **SEV3** | Minor feature broken, workaround exists | Normal priority |
| **SEV4** | Cosmetic issue, minimal impact | Backlog |

## Incident Response Process

### 1. Detection & Acknowledge

```markdown
## Incident Detected

**Time:** 2024-03-15 14:32 UTC
**Detected by:** Monitoring alert / User report / Manual discovery
**Initial symptoms:** [What we observed]

**Acknowledged by:** [Your name]
**Time to acknowledge:** [X minutes]
```

### 2. Assess & Classify

```markdown
## Assessment

**Severity:** SEV2
**Impact:** 
- Users affected: ~30% (estimate)
- Functionality impaired: Event creation failing
- Business impact: Users cannot log career events

**Initial hypothesis:** Database connection exhaustion
```

### 3. Communicate

```markdown
## Status Update #1 (14:45 UTC)

**Status:** Investigating
**Impact:** Event creation is failing for some users
**Current action:** Investigating database connectivity
**ETA:** Unknown, next update in 15 minutes
```

### 4. Mitigate

**Goal:** Restore service, even if imperfectly

```markdown
## Mitigation Options

1. **Rollback** - Revert last deployment
   - Risk: Low
   - Time: 5 minutes
   
2. **Scale up** - Add database connections
   - Risk: Low
   - Time: 2 minutes

3. **Feature flag** - Disable problematic feature
   - Risk: Low
   - Time: 1 minute

**Action taken:** Increased database connection pool
**Result:** Service restored at 14:52 UTC
```

### 5. Investigate Root Cause

Only after mitigation:

```markdown
## Investigation

### Timeline
- 14:00 - Deployment of v1.2.3
- 14:30 - Connection pool exhausted
- 14:32 - Alerts fired
- 14:45 - Investigation started
- 14:52 - Mitigated

### Root Cause
New feature introduced connection leak in error path.
When validation fails, connection not returned to pool.

### Evidence
- Logs show connections climbing: [link]
- Code review found leak: [file:line]
- Reproducer: [steps]
```

### 6. Resolve

```markdown
## Resolution

**Fix:** Added `defer conn.Close()` in error path
**PR:** #456
**Deployed:** 16:30 UTC
**Verified:** Connection pool stable for 1 hour
**Incident closed:** 17:30 UTC
```

### 7. Post-Mortem

```markdown
# Post-Mortem: Event Creation Outage

**Date:** 2024-03-15
**Duration:** 20 minutes (14:32 - 14:52)
**Severity:** SEV2
**Author:** [Name]

## Summary

Event creation failed for ~30% of users due to database 
connection exhaustion caused by a connection leak.

## Impact

- 20 minutes of degraded service
- ~150 failed event creations
- No data loss

## Timeline

[Detailed timeline from investigation]

## Root Cause

Connection leak in error handling path introduced in v1.2.3.
When validation fails, database connection was not returned
to pool, causing exhaustion under load.

## What Went Well

- Alert fired within 2 minutes
- Mitigation was quick (increased pool)
- Root cause identified same day

## What Went Poorly

- No integration test for error path
- Connection pool sizing had no alerting
- Took 13 minutes to acknowledge alert

## Action Items

| Action | Owner | Due |
|--------|-------|-----|
| Add connection leak test | @dev | 2024-03-20 |
| Add pool exhaustion alert | @ops | 2024-03-18 |
| Review error handling patterns | @team | 2024-03-22 |
| Improve on-call response time | @manager | 2024-03-25 |

## Lessons Learned

1. Error paths need same testing rigour as happy paths
2. Resource exhaustion should have earlier alerting
3. On-call acknowledgment SLA needs work
```

## Debugging in Production

### Safe Information Gathering

```bash
# Check application logs
kubectl logs deployment/kariya -f --tail=100

# Check error rates
curl localhost:9090/metrics | grep error

# Check resource usage
kubectl top pods

# Check recent events
kubectl get events --sort-by='.lastTimestamp'
```

### What NOT to Do

- Don't make changes without documenting
- Don't test fixes in production without mitigation ready
- Don't ignore the incident channel to "just fix it"
- Don't blame individuals

## Rollback Procedures

### Application Rollback

```bash
# Check current version
kubectl get deployment kariya -o jsonpath='{.spec.template.spec.containers[0].image}'

# Rollback to previous
kubectl rollout undo deployment/kariya

# Verify
kubectl rollout status deployment/kariya
```

### Database Rollback

```bash
# Check current migration
make db-version

# Rollback last migration
make migrate-down

# Verify
make db-version
```

### Feature Flag Rollback

```bash
# Disable problematic feature
export KARIYA_FEATURE_NEW_EVENTS=false

# Restart application
kubectl rollout restart deployment/kariya
```

## Communication Templates

### Initial Alert

```
🔴 INCIDENT: [Brief description]
Severity: SEV[X]
Impact: [Who/what affected]
Status: Investigating
Lead: @[name]
Channel: #incident-[date]
```

### Status Update

```
🟡 UPDATE: [Incident name]
Status: [Investigating/Mitigating/Monitoring]
Impact: [Current impact]
Action: [What we're doing]
Next update: [Time]
```

### Resolution

```
🟢 RESOLVED: [Incident name]
Duration: [X minutes/hours]
Impact: [Final impact assessment]
Root cause: [Brief summary]
Post-mortem: [Link] (within 48h)
```

## On-Call Best Practices

### Before On-Call

- [ ] Laptop charged and accessible
- [ ] VPN working
- [ ] Access to all systems verified
- [ ] Runbooks reviewed
- [ ] Escalation contacts known

### During On-Call

- Acknowledge alerts within 5 minutes
- Communicate early and often
- Document everything
- Don't hero - escalate when needed
- Rest when off-call

### After Incident

- Complete post-mortem within 48 hours
- Follow up on action items
- Share learnings with team

## Related Skills

- `logging-observability` - Understanding logs and metrics
- `debug-test` - Debugging techniques
- `devops` - Deployment and operations
