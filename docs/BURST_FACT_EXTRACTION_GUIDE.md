# Burst and Fact Extraction Guide

**Purpose**: Understand how KaRiya automatically detects related events (bursts) and extracts inferred competencies and achievements (facts).

**Target Audience**: Career professionals, consultants, engineers using KaRiya to organize and enrich their career events.

---

## Table of Contents

1. [Overview](#overview)
2. [Understanding Bursts](#understanding-bursts)
3. [Understanding Facts](#understanding-facts)
4. [Burst Detection Algorithm](#burst-detection-algorithm)
5. [Fact Extraction Process](#fact-extraction-process)
6. [Role Fit Classification](#role-fit-classification)
7. [Audience Relevance](#audience-relevance)
8. [Keyboard Shortcuts](#keyboard-shortcuts)
9. [Workflow Examples](#workflow-examples)
10. [Best Practices](#best-practices)
11. [Troubleshooting](#troubleshooting)
12. [Competency Reference](#competency-reference)

---

## Overview

### What are Bursts?

**Bursts** are automatically detected groupings of 2 or more related career events. They represent concentrated periods of work on related projects, technologies, or initiatives.

**Examples**:
- Multiple events about building a microservices platform
- Several events about mentoring team members
- Related events about a major system migration

### What are Facts?

**Facts** are automatically extracted inferences about your competencies, achievements, and professional strengths. They are grounded statements without aspirational language, representing things you've actually done.

**Examples**:
- ✅ "Led cross-functional team through critical product launch"
- ✅ "Architected distributed system supporting 10M+ requests/day"
- ✅ "Mentored 5 junior engineers in system design"
- ❌ "Would like to lead teams" (aspirational - rejected)
- ❌ "Could potentially architect systems" (speculative - rejected)

---

## Understanding Bursts

### Why Bursts Matter

Bursts help you:
1. **Identify Themes**: Recognize patterns in your career progression
2. **Group Related Work**: Organize events about the same project or initiative
3. **Tell Cohesive Stories**: Prepare narrative-based CVs with grouped achievements
4. **Understand Depth**: See where you've invested significant effort
5. **Generate Case Studies**: Create portfolio pieces from burst groupings

### Burst Components

Each burst contains:
- **Name**: Human-readable title (e.g., "Platform Modernization Initiative")
- **Description**: Optional details about the burst (e.g., "Migrated monolith to microservices")
- **Events**: 2 or more related CareerEvents
- **Competency Focus**: Primary competency area (e.g., Technical, Leadership)
- **Metadata**: Creation date, last update date

### Example Burst

```
Burst: Platform Modernization Initiative
Description: Led team through 18-month migration from monolith to microservices
Events:
  1. "Designed microservices architecture for monolith decomposition"
  2. "Led team through iterative migration strategy and rollout"
  3. "Architected service mesh for inter-service communication"
  4. "Implemented observability platform for distributed tracing"
  5. "Mentored team on microservices patterns and best practices"
Competency Focus: Technical
Created: 2023-01-15
```

---

## Understanding Facts

### Fact Properties

Each fact has:
- **Text**: The factual statement (1-2000 characters, no aspirational language)
- **Competencies**: 1+ categories (Leadership, Technical, Mentoring, Product, Consulting, Research)
- **Role Fit**: Your professional tier (Principal, EM, Staff Engineer, Senior IC)
- **Audience Relevance**: Who cares about this (Hiring Manager, Recruiter, Peer)
- **Strength Signal**: Key achievement indicator extracted from event
- **Source**: Inferred from event or burst

### Fact Validation Rules

**Text Rules**:
- ✅ "Led 5-person team through Q3 launch"
- ✅ "Reduced API latency by 40% through optimization"
- ✅ "Designed system supporting 10M daily active users"
- ❌ "Will lead teams" (aspirational: "will")
- ❌ "Should improve system performance" (speculative: "should")
- ❌ "Could potentially architect solutions" (aspirational: "could")

**Aspirational Keywords Rejected**:
will, should, could, might, may, want, wish, hope, plan, intend, attempt, try, would

### Example Facts

```
Fact 1:
Text: "Architected microservices platform supporting 50M+ requests daily"
Competencies: [Technical, Architecture]
Role Fit: Principal
Audience: [Hiring Manager, Peer]
Strength Signal: "Architected system handling massive scale"

Fact 2:
Text: "Led cross-functional team of 12 engineers through critical product launch"
Competencies: [Leadership, Project Management]
Role Fit: EM
Audience: [Hiring Manager, Recruiter, Peer]
Strength Signal: "Led team to successful launch"

Fact 3:
Text: "Mentored 8 junior engineers in system design and code review practices"
Competencies: [Mentoring, Technical]
Role Fit: Staff Engineer
Audience: [Hiring Manager, Peer]
Strength Signal: "Mentored multiple junior engineers"
```

---

## Burst Detection Algorithm

### How Burst Detection Works

KaRiya uses a **three-step algorithm** to suggest related events:

#### Step 1: Similarity Scoring
Calculates how similar two events are using:

**Text Similarity** (word overlap):
- Keyword matching on event text
- Compound similarity (if A~B and B~C, then A~C suggested)
- Score: 0.0 to 1.0

**Metadata Matching**:
- Same company (high weight: +0.3)
- Same project (high weight: +0.3)
- Shared tags (medium weight: +0.1 per tag)

**Combined Score**:
```
Overall Similarity = (Text Similarity × 0.4) +
                     (Company Match × 0.3) +
                     (Project Match × 0.3)
```

#### Step 2: Temporal Grouping
Events within **6 months** are considered related.

**Timeline Example**:
```
Jan 2023: "Started platform redesign"        ←─┐
Feb 2023: "Designed new architecture"          │ Related (within 6 months)
Mar 2023: "Led migration planning"             │
Apr 2023: "Completed migration"                │
May 2023: "Implemented monitoring"           ←─┘

Aug 2023: "Started new project"              ← Outside 6-month window
```

#### Step 3: Suggestion Generation
Combines similarity scores with temporal grouping:

**High Confidence** (score > 0.7):
- Automatically suggested with green indicator
- Usually clear relationship between events

**Medium Confidence** (score 0.4-0.7):
- Suggested with yellow indicator
- May need user review to confirm relationship

**Low Confidence** (score < 0.4):
- Not suggested
- User can manually create bursts if desired

---

## Fact Extraction Process

### How Facts are Extracted

KaRiya uses an **inference engine** to extract facts from events and bursts:

#### From Single Events

**Process**:
1. Analyze event text for keywords and competencies
2. Infer role fit based on keywords (lead, architect, mentor, etc.)
3. Determine audience relevance
4. Extract strength signal from verbs and achievements
5. Validate against aspirational language
6. Present for user confirmation

**Example**:
```
Event: "Led cross-functional team of 8 engineers through Q3 product launch"

Extracted Fact:
Text: "Led cross-functional team of 8 engineers through Q3 product launch"
Competencies: [Leadership, Project Management]
Role Fit: EM
Audience: [Hiring Manager, Recruiter, Peer]
Strength Signal: "Led team to successful launch"
```

#### From Bursts

**Process**:
1. Analyze all events in burst together
2. Identify common themes across events
3. Generate higher-level facts representing the burst
4. Infer role fit from overall contribution
5. Present for user confirmation

**Example**:
```
Burst: Platform Modernization (5 events)

Extracted Facts:
1. "Architected microservices platform replacing monolithic system"
   Competencies: [Technical, Architecture]
   Role Fit: Principal

2. "Led team through 18-month migration with zero service interruption"
   Competencies: [Leadership, Project Management]
   Role Fit: EM

3. "Implemented comprehensive monitoring and observability system"
   Competencies: [Technical, Operations]
   Role Fit: Staff Engineer
```

---

## Role Fit Classification

### What is Role Fit?

**Role Fit** indicates which career level or specialization a fact demonstrates:

### Role Fit Types

#### 🔴 Principal / Director Level
**Indicators**: Vision, strategy, company-wide impact, founding role, architecture

**Keywords**: principal, architect, vision, strategy, roadmap, company-wide, enterprise, organization, technical direction, founding, founder

**Example Facts**:
- "Architected company-wide infrastructure strategy supporting 10+ product lines"
- "Founded and led engineering organization from 0 to 50 engineers"
- "Set technical vision and roadmap for 500+ engineer organization"

**Use When**: You drove strategic initiatives with company-wide impact

#### 🟠 Engineering Manager (EM)
**Indicators**: Team leadership, hiring, people management, cross-team coordination

**Keywords**: manager, director, head, vp, vice president, management, people management, hiring, team

**Example Facts**:
- "Managed team of 12 engineers and contractors through critical launch"
- "Led hiring and onboarding process for engineering organization"
- "Directed cross-functional team coordination across 3 locations"

**Use When**: You led people, managed teams, or coordinated across groups

#### 💜 Staff Engineer
**Indicators**: Deep technical expertise, complex systems, mentoring, influence without authority

**Keywords**: staff engineer, principal engineer, deep expertise, complex, difficult, systems, architecture design

**Example Facts**:
- "Designed and implemented distributed database system for 100B+ record queries"
- "Led technical design review process for critical infrastructure changes"
- "Mentored 5 senior engineers in advanced system design techniques"

**Use When**: You drove technical excellence and complex system work

#### 💙 Senior IC / Individual Contributor
**Indicators**: Strong technical execution, feature delivery, code quality

**Default Role**: Assigned when no other role fit indicators present

**Example Facts**:
- "Implemented caching layer reducing API latency by 60%"
- "Delivered critical feature supporting new product line"
- "Refactored authentication system improving security posture"

**Use When**: You delivered strong individual technical work

### Role Fit Best Practices

1. **One primary role per fact**: Choose the strongest indicator
2. **Consider your career stage**: Staff+ engineers often have multiple roles
3. **Be accurate**: Hiring managers check role fit alignment
4. **Support with examples**: Back up role fit with specific achievements

---

## Audience Relevance

### Understanding Audience Relevance

**Audience Relevance** indicates who cares about a fact:

#### 👔 Hiring Manager
Cares about: Leadership, team impact, technical decision-making, mentoring

**Fact Types**:
- Leadership facts (led teams, managed people)
- Significant technical contributions (architecture, design)
- Mentoring and development facts
- Cross-team coordination

**Example**:
- "Led team through critical product launch" → Hiring Manager ✅
- "Improved code review process efficiency" → Hiring Manager ❌

#### 💼 Recruiter
Cares about: Management experience, team size, company growth, role progression

**Fact Types**:
- Leadership and management facts
- Team building and hiring
- Organization growth

**Example**:
- "Built engineering team from 5 to 50 engineers" → Recruiter ✅
- "Implemented new caching mechanism" → Recruiter ❌

#### 👥 Peer
Cares about: Technical depth, solutions, approaches, practices

**Fact Types**:
- All facts (every fact is relevant to peers)
- Technical details, systems, approaches
- Mentoring and knowledge sharing

**Example**:
- Any fact is relevant to peers ✅

### Default Audience Rules

**Peer**: Always included (all facts relevant to peers)

**Hiring Manager**: Included for:
- Leadership facts (contain keywords: lead, manage, team, mentor)
- Technical facts with significant impact

**Recruiter**: Included for:
- Management/leadership facts
- Organization building facts

---

## Keyboard Shortcuts

### Burst Suggestion Screen

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate between burst suggestions |
| `Enter` | Confirm burst suggestion and create burst |
| `y` | Confirm current burst (alternative) |
| `n` | Reject current burst suggestion |
| `e` | Edit burst name/description before confirming |
| `Escape` | Cancel burst suggestions and return to previous screen |
| `?` | Show help |

### Burst List Screen

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate through bursts |
| `←` / `→` | Navigate between sections (burst list ↔ events in burst) |
| `Space` | Select/deselect burst |
| `Enter` | View burst details |
| `e` | Edit burst (name, description, events) |
| `d` | Delete burst |
| `f` | Filter bursts by competency |
| `s` | Sort bursts (by date, event count, name) |
| `Escape` | Return to previous screen |
| `?` | Show help |

### Fact Editor Screen

| Key | Action |
|-----|--------|
| `Tab` | Move to next field |
| `Shift+Tab` | Move to previous field |
| `↑` / `↓` | Navigate selections in current field |
| `Space` | Toggle selection (competencies, audience) |
| `Enter` | Confirm and save fact |
| `Escape` | Cancel and discard changes |
| `?` | Show help |

---

## Workflow Examples

### Example 1: Single Event Fact Extraction

**Scenario**: You captured an event about leading a product launch.

**Steps**:
1. Capture event: "Led cross-functional team through critical Q3 product launch"
2. Navigate to metadata review (press `u`)
3. Review and confirm metadata (company, project, tags)
4. Return to home
5. View event details
6. Facts are automatically extracted:
   - "Led cross-functional team through Q3 product launch"
   - Competencies: [Leadership, Project Management]
   - Role Fit: EM
   - Audience: [Hiring Manager, Recruiter, Peer]
7. Confirm or edit fact as needed

### Example 2: Burst Detection and Confirmation

**Scenario**: You have 4 related events about platform migration.

**Events**:
1. "Designed microservices architecture"
2. "Led migration planning and strategy"
3. "Implemented new deployment pipeline"
4. "Mentored team on microservices patterns"

**Steps**:
1. Complete metadata clarification for all events
2. System suggests burst: "Platform Modernization Initiative"
   - Similarity score: 0.82 (high confidence)
   - Related events: 4
3. Review suggestion (press `u` for burst suggestions)
4. Confirm burst (press `y`)
5. System creates burst grouping all 4 events
6. Facts extracted from burst:
   - "Architected microservices platform"
   - "Led 18-month migration initiative"
   - "Implemented deployment automation"
   - "Mentored team in microservices patterns"
7. Review and confirm facts as needed

### Example 3: Bulk Fact Enrichment

**Scenario**: You have multiple related events and want to enrich all at once.

**Steps**:
1. View metadata review screen
2. Use bulk operations to select events in a burst
3. Edit metadata for all selected events simultaneously
4. System extracts facts for each event
5. Review extracted facts in batch mode

### Example 4: CSV Import with Automatic Burst/Fact Detection

**Scenario**: You're importing a batch of career events from a CSV file.

**Steps**:
1. Prepare CSV file with events (10+ events recommended for burst detection)
2. Import using CLI flag:
   ```bash
   ./kariya-cli --import events.csv
   ```
3. System imports all events
4. **Automatic burst detection runs**:
   - Analyzes event relationships
   - Groups similar events into bursts
   - Calculates confidence scores
   - Displays burst suggestions in import summary
5. **Automatic fact extraction runs**:
   - Extracts facts from each event
   - Infers competencies, role fit, audience
   - Persists facts to database
   - Displays extraction summary
6. Review import summary showing:
   - Events imported: X
   - Bursts detected: Y (with confidence scores)
   - Facts extracted: Z
7. Optionally review bursts and facts interactively:
   ```bash
   ./kariya-cli --show-bursts
   ./kariya-cli --show-facts
   ```

### Example 5: Re-running Burst Detection

**Scenario**: You've added new events and want to re-detect bursts.

**Steps**:
1. Run burst detection on all events:
   ```bash
   ./kariya-cli --detect-bursts
   ```
2. System analyzes all events in database
3. Displays new burst suggestions
4. Confirm or reject each suggestion
5. Bursts persisted to database

### Example 6: Re-running Fact Extraction

**Scenario**: You've updated event metadata and want to re-extract facts.

**Steps**:
1. Run fact extraction on all events:
   ```bash
   ./kariya-cli --extract-facts
   ```
2. System extracts facts from all events
3. Displays extraction summary by competency
4. Facts persisted to database
5. View facts with:
   ```bash
   ./kariya-cli --show-facts
   ```
6. Confirm/reject facts as group
7. Export enriched events with facts for CV generation

---

## Best Practices

### For Capturing Events

1. **Be Specific**: Include context and impact
   - ✅ "Led team of 8 through critical migration with zero downtime"
   - ❌ "Did some technical work"

2. **Use Action Verbs**: Start with verbs that show agency
   - ✅ "Architected", "Led", "Implemented", "Designed"
   - ❌ "Was responsible for", "Involved in"

3. **Include Metrics**: Quantify impact when possible
   - ✅ "Reduced latency by 40%, serving 10M+ requests daily"
   - ❌ "Improved performance"

4. **Tag for Grouping**: Use tags to hint at related events
   - Project names as tags
   - Technology names
   - Initiative names

### For Reviewing Burst Suggestions

1. **High Confidence (>0.7)**: Usually safe to accept
2. **Medium Confidence (0.4-0.7)**: Review carefully before accepting
3. **Edit Before Confirming**: Customize burst name/description
4. **Keep Related Events**: Don't remove clearly related events
5. **Reject Only If Unrelated**: Use rejection sparingly

### For Fact Extraction

1. **Avoid Aspirational Language**:
   - ❌ "Will lead teams"
   - ✅ "Led teams"

2. **Ground Metrics**: Use real, measured numbers
   - ❌ "Potentially improved performance"
   - ✅ "Improved API latency by 40%"

3. **Be Honest**: Facts are for you first, credibility second
   - ❌ Exaggerate achievements
   - ✅ Accurately describe what you did

4. **Provide Context**: More detail helps with inference
   - ❌ "Architected system"
   - ✅ "Architected distributed system supporting 50M+ daily requests"

5. **Update Role Fit**: Verify role fit matches your actual career level
   - Consider your current career stage
   - Ensure consistency across facts
   - Update if career progresses

---

## Troubleshooting

### Bursts Not Being Detected

**Problem**: You have related events but no burst suggestion appears.

**Causes & Solutions**:
1. **Events too far apart**: Bursts require events within 6 months
   - Solution: Check event dates, move dates if incorrect
2. **Low similarity score**: Events aren't similar enough
   - Solution: Add shared tags or company/project to events
3. **Insufficient events**: Bursts require ≥2 events
   - Solution: Capture more related events

**Manual Override**: You can manually create bursts if automatic detection doesn't work.

### Facts Not Extracted

**Problem**: Events exist but no facts are extracted.

**Causes & Solutions**:
1. **No competency indicators**: Event text lacks indicator keywords
   - Solution: Reword event to include action verbs (led, architected, designed)
2. **Aspirational language**: Event contains prohibited keywords
   - Solution: Rewrite to remove aspirational language (will, could, should)
3. **No event text**: Event with empty text can't generate facts
   - Solution: Add descriptive event text

### Role Fit Seems Wrong

**Problem**: Extracted fact has incorrect role fit.

**Causes & Solutions**:
1. **Keyword mismatch**: Event lacks role-specific keywords
   - Solution: Edit event text to include relevant keywords
2. **Career stage mismatch**: Role fit doesn't match your level
   - Solution: Edit fact to correct role fit (Staff → Principal, etc.)
3. **Default role assigned**: No keywords matched, defaulted to Senior IC
   - Solution: Add role-specific keywords to event text

### Audience Relevance Missing

**Problem**: Fact is missing expected audience.

**Causes & Solutions**:
1. **Leadership keywords missing**: Hiring Manager not included
   - Solution: Add "led", "team", "managed" keywords to event
2. **Incompatible competencies**: Some audiences filtered out
   - Solution: Add relevant competency categories
3. **Senior IC without leadership**: No manager relevance
   - Solution: Change role fit to EM if fact demonstrates leadership

---

## Competency Reference

### Six Competency Categories

#### 🔧 Technical
**Focus**: Deep technical skills, architecture, implementation

**Examples**:
- System design and architecture
- Programming and coding practices
- Infrastructure and DevOps
- Database and data systems
- Security implementation

**Keywords**: develop, engineer, code, implement, architect, backend, frontend, system, algorithm, optimize, design, technical

#### 👔 Leadership
**Focus**: Team leadership, strategic direction, decision-making

**Examples**:
- Team management and hiring
- Strategic planning and vision
- Cross-team coordination
- Roadmap and priority setting
- Organizational growth

**Keywords**: lead, manage, strategy, guide, mentor, direct, coordinate, transform, vision, roadmap, team, organization

#### 🎯 Product
**Focus**: Product thinking, user focus, feature delivery

**Examples**:
- Product strategy and roadmap
- Feature design and prioritization
- User experience and customer focus
- MVP and prototyping
- Market innovation

**Keywords**: product, feature, roadmap, design, user experience, customer, MVP, prototype, innovation, focus

#### 💼 Consulting
**Focus**: Advisory, optimization, transformation

**Examples**:
- Strategic consulting
- Process optimization
- Technology transformation
- Client solutions and recommendations
- Business impact improvement

**Keywords**: consult, advise, strategic, transform, client, solution, recommend, optimize, business, impact

#### 🔬 Research
**Focus**: Investigation, analysis, experimentation

**Examples**:
- Research and experimentation
- Data analysis and interpretation
- Proof of concept and prototyping
- New technology evaluation
- Methodology development

**Keywords**: research, analyze, investigate, discover, study, prototype, experiment, innovation, methodology, data, evaluation

#### 👨‍🏫 Mentoring
**Focus**: Teaching, coaching, skill development

**Examples**:
- Mentoring and coaching
- Training and skill development
- Knowledge sharing and documentation
- Junior engineer support
- Professional growth facilitation

**Keywords**: mentor, train, coach, develop, guide, support, teach, onboard, grow, skill development, knowledge

---

## FAQ

### Can I manually edit bursts?

Yes! You can edit burst names, descriptions, and event memberships. Use the burst editor to customize detected bursts.

### Can I delete facts?

Yes! If a fact doesn't apply, you can delete it. Deleted facts won't reappear during extraction.

### How often are facts extracted?

Facts are extracted:
- Automatically after you confirm metadata (if enabled)
- When viewing event details
- When reviewing burst details
- You can manually request extraction anytime

### Can I have facts without bursts?

Yes! Facts are extracted from individual events and from bursts. Most facts come from events.

### How do I use facts in CV generation?

Facts are used to populate CV content:
1. Select target audience (hiring manager, recruiter, peer)
2. Select role fit level
3. System groups facts by competency
4. Generate role-specific CV sections

### What's the difference between facts and events?

- **Events**: Raw career activities (what you captured)
- **Facts**: Inferred achievements (what you extracted from events)
- **Bursts**: Grouped events (related activities clustered)

### How do I improve fact quality?

1. **Add better event text**: More detail = better facts
2. **Use consistent language**: Similar events get similar facts
3. **Edit extracted facts**: Customize role fit and audience
4. **Add competency tags**: Tag events with relevant competencies

---

## Related Documentation

- [CLI_GUIDE.md](./CLI_GUIDE.md) - Complete keyboard shortcuts reference
- [METADATA_REVIEW_GUIDE.md](./METADATA_REVIEW_GUIDE.md) - Event metadata clarification
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Common issues and solutions
- [TUI_STANDARDS.md](./TUI_STANDARDS.md) - UI/UX standards and navigation

---

**Document Version**: 1.0
**Last Updated**: 2025-12-31
**Status**: Complete and Production-Ready

