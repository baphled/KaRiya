# CV Generation Examples

This document provides examples of CVs generated for different roles and audiences, showing how KaRiya transforms the same events into tailored CV bullets.

## Sample Career Events

To illustrate the examples, we'll use these sample career events:

### Event 1: Mobile App Performance
```
Text: "Led performance optimization of mobile app, reduced load time from 8s to 2s. Identified bottlenecks in image processing pipeline and refactored the caching strategy. This improved user retention by 12%."
Date: 2024-06-15
Company: TechCorp
Tags: technical, achievement, leadership
Categories: technical, product
```

### Event 2: API Design
```
Text: "Designed and implemented new REST API for customer data platform. Created comprehensive OpenAPI documentation and led design review with 6 engineers. API is now used by 3 internal teams and 2 external partners."
Date: 2024-03-20
Company: TechCorp
Tags: technical, achievement
Categories: technical, product
```

### Event 3: Team Mentoring
```
Text: "Mentored 4 junior engineers over 6 months. Helped them grow from junior to mid-level roles. Two of them were promoted and now lead their own teams."
Date: 2024-01-10
Company: TechCorp
Tags: leadership, achievement
Categories: mentoring
```

### Event 4: Incident Response
```
Text: "Responded to production incident affecting 50k users. Diagnosed root cause within 2 hours and implemented temporary fix. Worked with team to implement permanent solution within 24 hours."
Date: 2023-11-05
Company: TechCorp
Tags: technical
Categories: technical
```

### Event 5: Cross-team Collaboration
```
Text: "Collaborated with product team to define requirements for new reporting feature. Participated in 8 design meetings and provided technical perspective on feasibility and timeline."
Date: 2023-08-22
Company: TechCorp
Tags: leadership, achievement
Categories: product
```

## Example 1: Staff Engineer for Hiring Manager

**Configuration**:
- Target Role: Staff
- Target Audience: HiringManager
- Date Range: Last 2 years
- Companies: TechCorp
- Tags: None (all tags included)
- Categories: None (all categories included)

**Generated CV**:

```
STAFF ENGINEER CV - FOR HIRING MANAGER
Generated: 2024-12-01

PROFESSIONAL SUMMARY
Led high-impact technical initiatives that improved product performance and user experience,
with measurable business outcomes including 12% improvement in user retention.

EXPERIENCE
- Reduced mobile app load time from 8s to 2s through systematic performance optimization,
  improving user retention by 12%
  [Confidence: 0.95 | Sources: 1 event]

- Designed and implemented REST API for customer data platform serving 3 internal teams
  and 2 external partners
  [Confidence: 0.88 | Sources: 1 event]

- Mentored 4 junior engineers resulting in 2 promotions to mid-level roles
  [Confidence: 0.82 | Sources: 1 event]

- Diagnosed and resolved production incident affecting 50k users within 24 hours
  [Confidence: 0.75 | Sources: 1 event]

CORE COMPETENCIES
- Performance Optimization: Identified bottlenecks, refactored caching strategy,
  improved load time 75%
- API Design: REST API design, OpenAPI documentation, cross-team adoption
- Technical Leadership: Mentorship, code review leadership, incident response
```

**Key Observations**:
- Hiring manager sees business impact (user retention, user count affected)
- Bullets emphasize outcomes and measurable results
- Summary highlights business value
- Older event (incident) is lower in ranking due to temporal decay

## Example 2: Staff Engineer for Recruiter

**Configuration**:
- Target Role: Staff
- Target Audience: Recruiter
- Date Range: Last 2 years
- Companies: TechCorp
- Tags: None
- Categories: None

**Generated CV**:

```
STAFF ENGINEER CV - FOR RECRUITER
Generated: 2024-12-01

PROFESSIONAL SUMMARY
Experienced Staff Engineer with expertise in system design, performance optimization,
and technical mentorship. Demonstrated ability to lead complex projects and grow engineering talent.

EXPERIENCE
- Designed and implemented REST API for customer data platform, demonstrating
  strong system design skills and cross-team collaboration
  [Confidence: 0.88 | Sources: 1 event]

- Led performance optimization of mobile app, reducing load time by 75%,
  showcasing technical depth and impact
  [Confidence: 0.85 | Sources: 1 event]

- Mentored 4 junior engineers to mid-level roles, demonstrating technical leadership
  and people development skills
  [Confidence: 0.82 | Sources: 1 event]

- Responded to and resolved production incident, demonstrating debugging skills
  and operational excellence
  [Confidence: 0.75 | Sources: 1 event]

CORE COMPETENCIES
- System Design: REST API design, cross-team adoption, scalability
- Performance Engineering: Optimization, bottleneck identification, caching strategies
- Technical Leadership: Mentorship, code review, incident response, technical decision-making
```

**Key Observations**:
- Recruiter sees skills and competencies emphasized
- Focus on career growth and technical leadership
- Less emphasis on business metrics, more on technical approach
- Bullets highlight career progression and seniority

## Example 3: Principal Engineer for Hiring Manager

**Configuration**:
- Target Role: Principal
- Target Audience: HiringManager
- Date Range: Last 2 years
- Companies: TechCorp
- Tags: None
- Categories: None

**Generated CV** (Note: Principal role has 3-4 bullet cap):

```
PRINCIPAL ENGINEER CV - FOR HIRING MANAGER
Generated: 2024-12-01

PROFESSIONAL SUMMARY
Strategic technical leader who drives organizational impact through system design,
performance optimization, and talent development.

EXPERIENCE
- Reduced mobile app load time from 8s to 2s through systematic optimization,
  improving user retention by 12% and directly impacting business metrics
  [Confidence: 0.95 | Sources: 1 event]

- Designed REST API platform serving 3 internal teams and 2 external partners,
  establishing architectural standards across the organization
  [Confidence: 0.92 | Sources: 1 event]

- Mentored 4 junior engineers resulting in 2 promotions, demonstrating ability
  to develop senior technical talent and build organizational capability
  [Confidence: 0.88 | Sources: 1 event]

CORE COMPETENCIES
- Strategic Technical Leadership: System design, organizational standards, cross-team influence
- Performance & Scale: Optimization, architectural decisions, measurable business impact
- Talent Development: Mentorship, career growth, leadership development
```

**Key Observations**:
- Principal role has stricter bullet cap (3 bullets)
- Older incident response bullet was compressed out
- Remaining bullets emphasize strategic impact and organizational influence
- Language focuses on scale, standards, and organizational capability

## Example 4: Engineering Manager for Hiring Manager

**Configuration**:
- Target Role: EM
- Target Audience: HiringManager
- Date Range: Last 2 years
- Companies: TechCorp
- Tags: None
- Categories: None

**Generated CV** (Note: EM role has 3-4 bullet cap):

```
ENGINEERING MANAGER CV - FOR HIRING MANAGER
Generated: 2024-12-01

PROFESSIONAL SUMMARY
Engineering leader focused on team growth, delivery excellence, and cross-functional collaboration
to drive product outcomes.

EXPERIENCE
- Mentored 4 junior engineers resulting in 2 promotions to mid-level roles,
  demonstrating ability to grow and develop engineering talent
  [Confidence: 0.88 | Sources: 1 event]

- Led performance optimization initiative improving mobile app load time by 75%,
  directly impacting user retention by 12%
  [Confidence: 0.85 | Sources: 1 event]

- Collaborated with product team on new reporting feature, providing technical perspective
  and ensuring feasibility and timeline alignment across functions
  [Confidence: 0.78 | Sources: 1 event]

CORE COMPETENCIES
- Team Leadership: Mentorship, talent development, career growth
- Technical Delivery: Project leadership, cross-functional collaboration, execution
- Product Impact: User-focused outcomes, business metrics, feature delivery
```

**Key Observations**:
- EM role emphasizes mentorship and team growth
- Bullets focus on people development and cross-functional collaboration
- Technical depth is secondary to leadership and delivery
- Incident response bullet compressed out (activity-level for manager role)

## Example 5: Same Events, Different Audiences (Staff Role)

### Audience: Hiring Manager

```
EXPERIENCE
- Reduced mobile app load time from 8s to 2s, improving user retention by 12%
- Designed REST API for customer data platform serving 5+ teams/partners
- Mentored 4 junior engineers resulting in 2 promotions
```

**Focus**: Business impact, measurable outcomes, customer/user value

### Audience: Recruiter

```
EXPERIENCE
- Designed and implemented REST API, demonstrating system design expertise
- Led performance optimization initiative, showcasing technical depth
- Mentored junior engineers, demonstrating technical leadership capability
```

**Focus**: Skills, competencies, career progression, technical excellence

### Audience: Peer

```
EXPERIENCE
- Optimized mobile app performance through systematic bottleneck analysis and caching
  strategy refactoring
- Designed REST API with comprehensive OpenAPI documentation and cross-team adoption
- Mentored junior engineers on technical best practices and career development
```

**Focus**: Technical approach, design decisions, knowledge sharing, collaboration

## Example 6: Multi-Audience CV (Staff for Hiring Manager + Recruiter)

**Configuration**:
- Target Role: Staff
- Target Audience: HiringManager, Recruiter
- Date Range: Last 2 years
- Companies: TechCorp

**Generated CV** (More bullets included to satisfy both audiences):

```
STAFF ENGINEER CV - FOR HIRING MANAGER & RECRUITER
Generated: 2024-12-01

PROFESSIONAL SUMMARY
Staff Engineer with expertise in system design, performance optimization, and technical mentorship.
Demonstrated ability to drive measurable business impact while developing engineering talent.

EXPERIENCE
- Reduced mobile app load time from 8s to 2s through systematic performance optimization,
  improving user retention by 12%
  [Confidence: 0.95 | Sources: 1 event]

- Designed and implemented REST API for customer data platform serving 3 internal teams
  and 2 external partners, demonstrating strong system design and cross-team collaboration
  [Confidence: 0.88 | Sources: 1 event]

- Mentored 4 junior engineers resulting in 2 promotions to mid-level roles,
  demonstrating technical leadership and people development skills
  [Confidence: 0.82 | Sources: 1 event]

- Diagnosed and resolved production incident affecting 50k users within 24 hours,
  demonstrating debugging and operational excellence
  [Confidence: 0.75 | Sources: 1 event]

CORE COMPETENCIES
- System Design: REST API design, scalability, cross-team adoption
- Performance Engineering: Optimization techniques, bottleneck identification, measurable impact
- Technical Leadership: Mentorship, code review, incident response, team development
- Operational Excellence: Incident response, production systems, reliability
```

**Key Observations**:
- Multi-audience CVs include more bullets (satisfy both audiences)
- Bullets are phrased to appeal to both audiences
- Business impact AND technical depth both highlighted
- More comprehensive representation of capabilities

## Example 7: Filtered CV (Recent Technical Work Only)

**Configuration**:
- Target Role: Staff
- Target Audience: Peer
- Date Range: Last 6 months
- Companies: TechCorp
- Tags: technical
- Categories: technical

**Generated CV** (Only recent technical work):

```
STAFF ENGINEER CV - RECENT TECHNICAL WORK
Generated: 2024-12-01

PROFESSIONAL SUMMARY
Recent technical contributions demonstrating expertise in performance optimization
and API design for modern systems.

EXPERIENCE
- Reduced mobile app load time from 8s to 2s through systematic bottleneck analysis
  and caching strategy refactoring
  [Confidence: 0.95 | Sources: 1 event]

- Designed and implemented REST API with comprehensive documentation, establishing
  reusable patterns and cross-team adoption
  [Confidence: 0.88 | Sources: 1 event]

CORE COMPETENCIES
- Performance Optimization: Bottleneck analysis, caching strategies, measurement
- API Design: REST design, OpenAPI documentation, scalability patterns
- Technical Excellence: Code quality, design patterns, best practices
```

**Key Observations**:
- Filters remove mentoring and incident response events
- Only technical work (tags: technical) included
- More focused, specialized CV
- Useful for peer review or technical interviews

## Example 8: Compression in Action

**Before Compression** (8 bullets generated):
1. Mobile app optimization (0.95) - Ownership
2. REST API design (0.88) - Ownership
3. Mentoring (0.82) - Contribution
4. Incident response (0.75) - Execution
5. Design review participation (0.65) - Contribution
6. Code review leadership (0.55) - Activity
7. Meeting participation (0.45) - Activity
8. Debugging assistance (0.35) - Activity

**After Compression for Staff Role** (4-5 bullet cap):
1. Mobile app optimization (0.95) ✓ Keep
2. REST API design (0.88) ✓ Keep
3. Mentoring (0.82) ✓ Keep
4. Incident response (0.75) ✓ Keep
5. Design review participation (0.65) → Removed (lower priority)
6. Code review leadership (0.55) → Removed
7. Meeting participation (0.45) → Removed
8. Debugging assistance (0.35) → Removed

**Result**: 4 strongest bullets kept, lower-priority activities removed

## Key Takeaways from Examples

1. **Same events, different roles**: Principal role emphasizes strategic impact, while EM emphasizes team growth

2. **Same events, different audiences**: Hiring managers focus on business outcomes, recruiters on skills, peers on technical depth

3. **Filtering matters**: Date range and tag filters create focused CVs for specific purposes

4. **Compression is intelligent**: Lower-priority bullets are removed while maintaining representation across sections

5. **Confidence scores guide inclusion**: Higher-confidence bullets (with multiple sources or clear ownership) are always kept

6. **Traceability is preserved**: Every bullet traces back to source events, allowing verification

7. **Multi-audience CVs are richer**: Including multiple audiences results in more comprehensive CVs

8. **Quality events produce quality bullets**: Events with clear ownership, metrics, and context produce higher-scoring bullets

---

**Document Version**: 1.0
**Last Updated**: 2026-01-02

