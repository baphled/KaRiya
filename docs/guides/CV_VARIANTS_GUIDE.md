---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# CV Variants Guide

## Overview

KaRiya uses a **technology-focused variant system** to generate CVs tailored to specific roles and audiences. Instead of static variant templates, KaRiya dynamically generates CV variants based on your technology selections, focus area, and presentation preferences.

This guide explains how CV variants work, how to select the right variant for your needs, and how the variant system adapts to your career profile.

## Technology-Focused Variants

### The Three Technology Focus Types

KaRiya offers three presentation styles that determine how your technical expertise is showcased:

#### 1. Language Agnostic

**When to use**: You want to emphasize **adaptability and breadth** across technologies rather than specific tech stacks.

**Best for**:
- Leadership roles (EM, Staff+, Principal)
- Consulting or advisory positions
- Roles requiring technology agnostic decision-making
- Transitioning between tech stacks

**Characteristics**:
- **Narrative structure**: Emphasizes outcomes, leadership, and impact over specific tools
- **All events included**: No technology-based filtering
- **Skills section optional**: Technical skills listed but not emphasized
- **Bullet style**: Focuses on "what you achieved" not "what you used"

**Example bullet**:
> Led platform modernization initiative, reducing deployment time by 60% and improving system reliability to 99.99% uptime

#### 2. Generalist (2-5 Technologies)

**When to use**: You want to showcase **depth in multiple technologies** that you regularly use together.

**Best for**:
- Full-stack engineers
- Polyglot developers
- Roles requiring expertise in complementary tech stacks
- Modern development roles

**Characteristics**:
- **Standard structure**: Traditional CV format with skills-first approach
- **Technology-boosted bullets**: Events using selected technologies ranked higher
- **Skills section prominent**: Selected technologies highlighted at top
- **Bullet style**: Balances outcomes with technical implementation

**Example selection**: Ruby, PostgreSQL, Docker, React, Redis (5 technologies)

**Example bullet**:
> Built microservices platform using Ruby and PostgreSQL, processing 1M+ daily transactions with 99.9% uptime

#### 3. Specialist (1 Technology)

**When to use**: You want to demonstrate **deep expertise in a specific technology** or ecosystem.

**Best for**:
- Deep technical IC roles
- Technology-specific positions (e.g., "Senior Go Engineer")
- Roles requiring certified expertise
- Niche technology markets

**Characteristics**:
- **Standard structure**: Traditional CV format
- **Highly filtered bullets**: Heavy preference for events using the selected technology
- **Skills section specialized**: Selected technology featured prominently
- **Bullet style**: Technical depth and specifics emphasized

**Example selection**: Ruby (1 technology)

**Example bullet**:
> Optimized Ruby on Rails application performance, reducing response times from 800ms to 120ms through database query optimization and caching strategies

## Focus Areas

In addition to technology focus, you select a **focus area** that categorizes your primary domain:

### Backend

**Specializations**: Server-side applications, APIs, databases, system architecture

**Common skills**: Ruby, Go, Python, Java, PostgreSQL, Redis, Kafka, gRPC

**Suitable for**: Backend Engineers, Platform Engineers, API Developers, Database Specialists

### Frontend

**Specializations**: User interfaces, web applications, client-side frameworks, UX engineering

**Common skills**: React, Vue.js, TypeScript, JavaScript, CSS, HTML, Webpack, Next.js

**Suitable for**: Frontend Engineers, UI Engineers, JavaScript Developers, UX Engineers

### Fullstack

**Specializations**: End-to-end application development, both client and server

**Common skills**: Mixed backend + frontend technologies (e.g., Ruby + React, Node.js + Vue)

**Suitable for**: Full-stack Engineers, Product Engineers, Startup Engineers

### DevOps

**Specializations**: Infrastructure, deployment, CI/CD, monitoring, cloud platforms

**Common skills**: Docker, Kubernetes, Terraform, AWS, GCP, Jenkins, Prometheus, Ansible

**Suitable for**: DevOps Engineers, SRE, Platform Engineers, Infrastructure Engineers

## Skills Section Formatting

KaRiya offers two formats for displaying technical skills in your CV:

### Flat Format

**Structure**: One skill per bullet point

**Best for**:
- Short skills lists (< 10 skills)
- Emphasizing each skill individually
- Clean, minimal CVs
- Recruiter-friendly formats (easy scanning)

**Example**:
```
Technical Skills
- Ruby
- Go
- PostgreSQL
- Docker
- Redis
```

### Grouped Format

**Structure**: Skills grouped by category with category headers

**Best for**:
- Long skills lists (10+ skills)
- Showcasing expertise breadth across categories
- Organizing related technologies together
- Technical audience (engineers, architects)

**Example**:
```
Technical Skills

Backend
  - Ruby
  - Go
  - Python

Database
  - PostgreSQL
  - Redis
  - MongoDB

DevOps
  - Docker
  - Kubernetes
  - Terraform
```

**Skill Limits**:
- **Flat format**: Limit applies to total skills shown
- **Grouped format**: Limit applies per category
- **Range**: 0 (no limit) to 50 skills
- **Recommendation**: 10-15 for flat, 5-8 per category for grouped

## Variant Naming Convention

Variants are named using this structure:

```
{focus}_{area}_{length}
```

### Examples

**Language Agnostic Variants**:
- `agnostic_backend_standard` - Language agnostic, backend focus, standard length
- `agnostic_fullstack_short` - Language agnostic, fullstack focus, short length
- `agnostic_devops_ultra_short` - Language agnostic, DevOps focus, ultra short

**Generalist Variants**:
- `generalist_backend_full` - Generalist (2-5 techs), backend focus, full length
- `generalist_frontend_standard` - Generalist (2-5 techs), frontend focus, standard
- `generalist_fullstack_short` - Generalist (2-5 techs), fullstack focus, short

**Specialist Variants**:
- `specialist_ruby_backend_full` - Ruby specialist, backend focus, full length
- `specialist_react_frontend_standard` - React specialist, frontend, standard
- `specialist_kubernetes_devops_short` - Kubernetes specialist, DevOps, short

## Length Formats

CV length determines how many bullets and how much detail to include:

### Ultra Short (1 page)

- **Total bullets**: ~8-12
- **Use case**: LinkedIn summary, brief profiles, initial screening
- **Structure**: Highlights format (top achievements only)

### Short (1-2 pages)

- **Total bullets**: ~15-25
- **Use case**: Standard job applications, recruiter submissions
- **Structure**: Concise standard format

### Standard (2-3 pages)

- **Total bullets**: ~30-40
- **Use case**: Technical roles, detailed applications
- **Structure**: Full standard format

### Full (3-4 pages)

- **Total bullets**: ~50-60
- **Use case**: Senior/leadership roles, comprehensive portfolios
- **Structure**: Detailed narrative with extensive evidence

## Choosing the Right Variant

### Decision Tree

```
1. What's your primary goal?
   └─ Emphasize adaptability → Language Agnostic
   └─ Showcase specific tech stack → Generalist (2-5 techs)
   └─ Demonstrate deep expertise → Specialist (1 tech)

2. What's your domain?
   └─ Backend, Frontend, Fullstack, or DevOps

3. How long should your CV be?
   └─ Ultra Short (screening), Short (standard), Standard (technical), Full (senior/leadership)

4. How should skills be displayed?
   └─ Flat (< 10 skills, recruiter-friendly)
   └─ Grouped (10+ skills, technical audience)
```

### Example Scenarios

**Scenario 1**: "I'm a Senior Backend Engineer applying to a Ruby-specific role"
- **Technology Focus**: Specialist (Ruby)
- **Focus Area**: Backend
- **Length**: Standard
- **Skills Format**: Grouped
- **Variant**: `specialist_ruby_backend_standard`

**Scenario 2**: "I'm a Staff Engineer applying to a leadership role, tech stack doesn't matter"
- **Technology Focus**: Language Agnostic
- **Focus Area**: Backend (or your primary domain)
- **Length**: Full
- **Skills Format**: Flat or Grouped
- **Variant**: `agnostic_backend_full`

**Scenario 3**: "I'm a Full-stack Engineer who works with React, Node.js, and PostgreSQL"
- **Technology Focus**: Generalist (React, Node.js, PostgreSQL)
- **Focus Area**: Fullstack
- **Length**: Standard
- **Skills Format**: Grouped
- **Variant**: `generalist_fullstack_standard`

## How Variants Affect CV Generation

### Bullet Filtering and Ranking

**Language Agnostic**:
- All events included (no filtering)
- Bullets ranked by quality, recency, and impact
- Technology mentions minimized in bullet text

**Generalist**:
- Events with selected technologies get **+0.15 score bonus**
- Events without skills still included (not penalized)
- High-quality events without selected techs can still rank well
- Technology mentions included where relevant

**Specialist**:
- Heavy preference for events with selected technology
- Other high-quality events included as supporting evidence
- Technology specifics emphasized in bullets

### Skills Section Population

1. **Skills extracted from events**: All unique skills from career events
2. **Prioritization**: Selected technologies appear first
3. **Sorting**: Selected first, then alphabetically
4. **Grouping** (if grouped format): By skill category (backend, frontend, database, devops, etc.)
5. **Limits applied**: Per section (flat) or per category (grouped)

### CV Structure

**Language Agnostic → Narrative Structure**:
- Summary section first (prose)
- Experience section (outcome-focused bullets)
- Projects section (if applicable)
- Skills section (minimal, supporting)

**Generalist/Specialist → Standard Structure**:
- Summary section (brief)
- Skills section (prominent)
- Experience section (tech + outcome bullets)
- Projects section (if applicable)

**Ultra Short → Highlights Structure** (regardless of focus):
- Top 8-12 bullets only
- No detailed sections
- Compact format

## Target Role Impact

The **target role** (from your profile) affects bullet filtering independently of technology focus:

- **Principal/Staff**: Emphasis on leadership, strategy, outcomes
- **EM**: Emphasis on team leadership, process, people
- **Senior IC**: Emphasis on technical depth, ownership, execution

This is **orthogonal to technology focus** - you can be a "Language Agnostic Principal Engineer" or a "Ruby Specialist Staff Engineer."

## Best Practices

### DO:

✅ Choose **Language Agnostic** for leadership roles  
✅ Choose **Generalist** when you regularly use 2-5 complementary technologies  
✅ Choose **Specialist** when the role explicitly requires deep expertise in one technology  
✅ Use **Grouped format** for 10+ skills  
✅ Limit skills to 10-15 (flat) or 5-8 per category (grouped)  
✅ Select focus area matching the role's primary domain  
✅ Use **Standard** or **Full** length for technical roles  

### DON'T:

❌ Don't select all your skills as "selected technologies" (defeats the purpose)  
❌ Don't use **Specialist** for leadership/management roles  
❌ Don't use **Flat format** with 20+ skills (hard to scan)  
❌ Don't use **Ultra Short** for senior technical roles (insufficient detail)  
❌ Don't select technologies you haven't used in 3+ events  

## Examples

See [CV_EXAMPLES.md](CV_EXAMPLES.md) for complete CV examples demonstrating each variant type.

## Related Documentation

- [CV Generation Guide](CV_GENERATION_GUIDE.md) - Complete workflow for generating CVs
- [CV Examples](CV_EXAMPLES.md) - Sample CVs for each variant
- [CV Troubleshooting](CV_TROUBLESHOOTING.md) - Common issues and solutions
- [CV Generation Workflow](../workflows/CV_GENERATION_WORKFLOW.md) - Detailed workflow documentation

---

**Last Updated**: 2026-01-14  
**Version**: 2.0 (Technology-Focused System)
