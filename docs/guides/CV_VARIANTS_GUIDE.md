# CV Variants Guide

## Overview

KaRiya's CV Variant system provides 16 pre-configured CV templates that combine two dimensions: **Role Emphasis** (what aspect of your experience to highlight) and **Length Format** (how much detail to include). This two-step selection process makes it easy to generate targeted CVs for different job applications.

### Key Concepts

- **Role Emphasis**: Determines which categories of experience are prioritized (e.g., technical depth vs leadership)
- **Length Format**: Controls CV density through year limits, company limits, and confidence thresholds
- **CV Structure**: The output format determined automatically by your variant selection
- **Variant**: A specific combination of role emphasis + length format with pre-configured settings

## The Two-Step Selection Process

When generating a CV, you select:

1. **Role Emphasis** - What aspect of your career to emphasize
2. **Length Format** - How detailed/long the CV should be

KaRiya then automatically determines the appropriate CV structure based on your selections.

```
Profile → Audience → Role Emphasis → Length Format → Generate CV
```

## Role Emphasis Options

### Senior Backend Developer

**Best for**: Backend-focused engineering roles, API development, database work

| Property | Value |
|----------|-------|
| Primary Categories | technical, architecture |
| Secondary Categories | product, delivery |
| Bullet Focus | Technical depth, product impact |
| Output Structure | Standard (full/standard/short) or Highlights (ultra-short) |

**When to use**:
- Applying for backend engineering positions
- Roles emphasizing system design and implementation
- Positions requiring deep technical expertise

### Staff/Principal Engineer

**Best for**: Senior IC roles, technical leadership, architecture positions

| Property | Value |
|----------|-------|
| Primary Categories | leadership, strategy, architecture |
| Secondary Categories | technical, mentoring |
| Bullet Focus | Architecture decisions, mentorship, cross-team impact |
| Output Structure | Standard (full/standard/short) or Highlights (ultra-short) |

**When to use**:
- Staff or Principal Engineer applications
- Roles requiring technical leadership without management
- Positions emphasizing architectural influence

### Consulting Engineer

**Best for**: Consulting roles, client-facing work, rapid assessment projects

| Property | Value |
|----------|-------|
| Primary Categories | strategy, delivery, consulting |
| Secondary Categories | technical, leadership |
| Bullet Focus | Client engagements, rapid assessment, delivery |
| Output Structure | Consulting (full/standard/short) or Highlights (ultra-short) |

**When to use**:
- Consulting firm applications
- Client-facing technical roles
- Positions requiring diverse project experience

### Language-Agnostic Engineer

**Best for**: Polyglot developers, cross-platform work, adaptability emphasis

| Property | Value |
|----------|-------|
| Primary Categories | technical, architecture |
| Secondary Categories | all categories |
| Bullet Focus | Multi-language evidence, adaptability, pragmatic problem-solving |
| Output Structure | Narrative (full/standard/short) or Highlights (ultra-short) |

**When to use**:
- Roles not tied to specific technologies
- Positions valuing adaptability and learning
- Companies with diverse tech stacks

## Length Format Options

### Full (3+ pages)

**Best for**: Comprehensive career documentation, internal promotions

| Property | Value |
|----------|-------|
| Max Years History | Unlimited |
| Max Companies | Unlimited |
| Max Bullets per Job | 8 |
| Min Confidence | 0.50 |
| Target Pages | 3+ |

**Includes**:
- Complete work history
- All qualifying bullets
- Full skills section
- Education (if configured)

### Standard (2-3 pages)

**Best for**: Most job applications, external positions

| Property | Value |
|----------|-------|
| Max Years History | 10 years |
| Max Companies | Unlimited |
| Max Bullets per Job | 6 |
| Min Confidence | 0.65 |
| Target Pages | 2-3 |

**Includes**:
- Last 10 years of experience
- Higher-confidence bullets only
- Summary section
- Skills section

### Short (1-2 pages)

**Best for**: Quick reviews, networking, initial screenings

| Property | Value |
|----------|-------|
| Max Years History | 5 years |
| Max Companies | 5 |
| Max Bullets per Job | 4 |
| Min Confidence | 0.75 |
| Target Pages | 1-2 |

**Includes**:
- Last 5 years of experience
- Top 5 companies only
- High-confidence bullets
- Condensed skills

### Ultra-Short (1 page)

**Best for**: One-pagers, executive summaries, quick introductions

| Property | Value |
|----------|-------|
| Max Years History | 3 years |
| Max Companies | 3 |
| Max Bullets per Job | 3 |
| Min Confidence | 0.85 |
| Target Pages | 1 |

**Includes**:
- Last 3 years only
- Top 3 companies
- Only highest-confidence bullets
- Key capabilities highlight

## The 16 Built-In Variants

### Variant Matrix

| Role Emphasis | Full | Standard | Short | Ultra-Short |
|---------------|------|----------|-------|-------------|
| **Senior Backend** | Standard | Standard | Standard | Highlights |
| **Staff/Principal** | Standard | Standard | Standard | Highlights |
| **Consulting** | Consulting | Consulting | Consulting | Highlights |
| **Language-Agnostic** | Narrative | Narrative | Narrative | Highlights |

### Variant Details

#### Senior Backend Variants

| Variant ID | Min Confidence | Max Years | Max Companies | Structure |
|------------|---------------|-----------|---------------|-----------|
| `senior_backend_full` | 0.60 | Unlimited | Unlimited | Standard |
| `senior_backend_standard` | 0.65 | 10 | Unlimited | Standard |
| `senior_backend_short` | 0.75 | 5 | Unlimited | Standard |
| `senior_backend_ultra_short` | 0.85 | Unlimited | 3 | Highlights |

#### Staff/Principal Variants

| Variant ID | Min Confidence | Max Years | Max Companies | Structure |
|------------|---------------|-----------|---------------|-----------|
| `staff_principal_full` | 0.70 | Unlimited | Unlimited | Standard |
| `staff_principal_standard` | 0.75 | 10 | Unlimited | Standard |
| `staff_principal_short` | 0.80 | 5 | Unlimited | Standard |
| `staff_principal_ultra_short` | 0.90 | Unlimited | 3 | Highlights |

#### Consulting Variants

| Variant ID | Min Confidence | Max Years | Max Companies | Structure |
|------------|---------------|-----------|---------------|-----------|
| `consulting_full` | 0.50 | Unlimited | Unlimited | Consulting |
| `consulting_standard` | 0.60 | 10 | Unlimited | Consulting |
| `consulting_short` | 0.70 | 5 | Unlimited | Consulting |
| `consulting_ultra_short` | 0.80 | Unlimited | 3 | Highlights |

#### Language-Agnostic Variants

| Variant ID | Min Confidence | Max Years | Max Companies | Structure |
|------------|---------------|-----------|---------------|-----------|
| `language_agnostic_full` | 0.65 | Unlimited | Unlimited | Narrative |
| `language_agnostic_standard` | 0.70 | 10 | Unlimited | Narrative |
| `language_agnostic_short` | 0.78 | 5 | Unlimited | Narrative |
| `language_agnostic_ultra_short` | 0.88 | Unlimited | 3 | Highlights |

## CV Structures

The variant system uses four CV structures. The structure is determined automatically based on your role emphasis and length selection.

### Standard Structure

Traditional CV format used by Senior Backend and Staff/Principal variants.

**Sections**:
1. **Summary** - Professional summary statement
2. **Experience** - Work history with achievements
3. **Projects** - Notable projects (if applicable)
4. **Skills** - Technical skills and competencies

### Narrative Structure

Language-agnostic format emphasizing pragmatic expertise.

**Sections**:
1. **Profile Header** - Name, title, location, contact
2. **Positioning Statement** - Career positioning (optional)
3. **Summary** - Professional summary
4. **Core Strengths** - 6 key competencies
5. **Technologies** - Languages, frontend, systems
6. **Selected Experience** - High-confidence achievements
7. **What I Bring** - Value propositions

### Consulting Structure

Client-focused format for consulting professionals.

**Sections**:
1. **Profile Header** - Name, title, location, contact
2. **Summary** - Professional summary
3. **Client Engagements** - Experience grouped by company
4. **Technical Capabilities** - Skills (optional)
5. **What I Bring** - Value propositions

### Highlights Structure

One-page executive summary used by all ultra-short variants.

**Sections**:
1. **Condensed Header** - Single-line profile
2. **Summary** - Brief professional summary
3. **Key Capabilities** - 4-6 core competencies
4. **Selected Highlights** - Top 5 bullets by confidence
5. **Technologies** - Languages and systems

## How Variants Affect Your CV

### Bullet Filtering

Variants filter your bullets in multiple ways:

1. **Confidence Threshold**: Only bullets above the variant's min confidence are included
2. **Category Scoring**: Primary categories score 1.0, secondary 0.6, others 0.3
3. **Date Filtering**: Events older than max years are excluded
4. **Company Limiting**: Only top N companies by recency (for ultra-short)

### Example: Same Career, Different Variants

**Senior Backend Standard** might produce:
```
EXPERIENCE

Acme Corp | Senior Backend Engineer | 2022-Present
- Designed and implemented distributed caching layer (confidence: 0.92)
- Led migration from monolith to microservices (confidence: 0.88)
- Optimized database queries reducing latency by 40% (confidence: 0.78)

Previous Corp | Backend Developer | 2019-2022
- Built real-time notification system (confidence: 0.85)
- Implemented CI/CD pipeline (confidence: 0.72)
```

**Senior Backend Ultra-Short** (same career) might produce:
```
KEY CAPABILITIES
- Distributed systems architecture
- Performance optimization
- Technical leadership

SELECTED HIGHLIGHTS
- Designed distributed caching layer reducing response times by 60%
- Led microservices migration for 50-engineer organization
- Built real-time notification system serving 1M+ users
```

### Profile Override

Variants can include profile overrides to customize:

- **Professional Title**: Role-specific title (e.g., "Senior Consulting Engineer")
- **Core Strengths**: Variant-specific competencies
- **Career Differentiators**: Unique value propositions

These overrides are applied during export, keeping your base profile unchanged.

## Using Variants in the TUI

### Step-by-Step Workflow

1. **Select Profile**: Choose your CV configuration
2. **Select Audience**: Choose target audience (hiring_manager, recruiter, peer)
3. **Select Role Emphasis**: Choose from 4 role emphases
   - Navigate with `j`/`k` or arrow keys
   - Press `Enter` to select
   - Press `Esc` to go back
4. **Select Length Format**: Choose from 4 length options
   - Navigate with `j`/`k` or arrow keys
   - Press `Enter` to generate
   - Press `Esc` to go back
5. **Review Generated CV**: Preview with full traceability
6. **Export**: Choose format (text, markdown, YAML)

### Keyboard Shortcuts

| State | Key | Action |
|-------|-----|--------|
| Role Emphasis | `j`/`↓` | Move down |
| Role Emphasis | `k`/`↑` | Move up |
| Role Emphasis | `Enter` | Select and proceed |
| Role Emphasis | `Esc` | Go back to audience |
| Length Format | `j`/`↓` | Move down |
| Length Format | `k`/`↑` | Move up |
| Length Format | `Enter` | Generate CV |
| Length Format | `Esc` | Go back to role emphasis |

## Common Workflows

### Workflow 1: Technical Job Application

1. Select **Senior Backend** role emphasis
2. Select **Standard** length (2-3 pages)
3. Review generated CV
4. Export as Markdown for further editing

### Workflow 2: Executive Summary

1. Select appropriate role emphasis for your background
2. Select **Ultra-Short** length (1 page)
3. Review the highlights-format CV
4. Export as Text for quick sharing

### Workflow 3: Consulting Proposal

1. Select **Consulting** role emphasis
2. Select **Full** length to show all client engagements
3. Review client engagement groupings
4. Export as Markdown for proposal documents

### Workflow 4: Tech-Agnostic Application

1. Select **Language-Agnostic** role emphasis
2. Select **Standard** length
3. Review narrative-format CV emphasizing adaptability
4. Export as Markdown

## Tips & Best Practices

### 1. Match Role Emphasis to Job Description

- Backend-heavy JD? Use **Senior Backend**
- Leadership mentioned? Use **Staff/Principal**
- Consulting role? Use **Consulting**
- "Any language" or "polyglot"? Use **Language-Agnostic**

### 2. Length Format Guidelines

| Situation | Recommended Length |
|-----------|-------------------|
| Initial application | Standard |
| Recruiter screening | Short |
| Detailed review | Full |
| Networking/quick share | Ultra-Short |
| Internal promotion | Full |

### 3. Review Confidence Thresholds

Higher role emphasis levels (Staff/Principal) have higher confidence thresholds because:
- Senior roles expect higher-quality achievements
- Leadership bullets need stronger evidence
- Fewer, stronger bullets beat many weak ones

### 4. Iterate on Source Data

If your CV is missing important work:
1. Check the confidence scores of excluded bullets
2. Improve source events to boost confidence
3. Add more facts supporting the achievement
4. Regenerate with the same variant

### 5. Use Ultra-Short for Initial Contact

The Highlights structure is ideal for:
- LinkedIn messages
- Email introductions
- Conference networking
- Quick capability overviews

## Troubleshooting

### Problem: "Not enough bullets in CV"

**Causes**:
- Confidence threshold too high for your data
- Date filter excluding relevant events
- Category mismatch with role emphasis

**Solutions**:
1. Try a lower length format (Full instead of Standard)
2. Use a role emphasis that matches your experience categories
3. Improve source events to increase confidence scores

### Problem: "Wrong sections appearing"

**Causes**:
- Structure mismatch with expectations
- Optional sections appearing/missing

**Solutions**:
1. Check which structure your variant uses (see matrix above)
2. Remember: ultra-short always uses Highlights structure
3. Optional sections only appear if data exists

### Problem: "CV too long for target"

**Causes**:
- Selected length format too permissive
- Too many high-confidence bullets

**Solutions**:
1. Use a shorter length format
2. Use Ultra-Short for strict 1-page limit
3. Export and manually trim if needed

## Related Documentation

- [CV Generation Guide](CV_GENERATION_GUIDE.md) - Complete CV generation documentation
- [Narrative CV Guide](NARRATIVE_CV_GUIDE.md) - Details on narrative structure
- [CV Troubleshooting](CV_TROUBLESHOOTING.md) - Common issues and solutions

---

**Document Version**: 1.0
**Last Updated**: 2026-01-09
