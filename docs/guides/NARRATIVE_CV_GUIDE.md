# Narrative CV Guide

## Overview

The **Narrative CV structure** is a language-agnostic format designed for senior engineers who want to emphasize pragmatic expertise over technology identity. This guide explains when and how to use narrative CVs effectively.

> **Note**: The narrative structure is automatically selected when you choose **Language-Agnostic** role emphasis with Full, Standard, or Short length formats. Ultra-Short always uses the Highlights structure. See the [CV Variants Guide](CV_VARIANTS_GUIDE.md) for the complete variant system.

## When to Use Narrative Structure

### Best For

- **Senior engineers with cross-domain experience** - You've worked across backend, frontend, infrastructure, and want to show breadth
- **Language-agnostic professionals** - You treat languages as tools, not identity ("I solve problems, languages are just tools")
- **Pragmatic, outcome-focused engineers** - You care more about results than technology evangelism
- **Roles requiring broad technical expertise** - Staff/Principal roles needing systems thinking

### Less Suitable For

- **Specialized roles** - When deep expertise in one language/framework is required
- **Junior to mid-level positions** - Where demonstrating growth trajectory matters more
- **Traditional company cultures** - Where standard CV formats are expected
- **Automated ATS screening** - Some systems expect standard CV sections

## Narrative CV Structure

### Section Overview

| Section | Purpose | Source |
|---------|---------|--------|
| **Profile Header** | Contact info and professional identity | Profile configuration |
| **Summary** | 2-3 sentence professional summary | CV sections or defaults |
| **Core Strengths** | 6 key competencies | Profile configuration |
| **Languages & Technologies** | Technical skills grouped by type | Profile configuration |
| **Selected Experience** | High-confidence achievements | Generated bullets (>= 0.75) |
| **What I Bring** | 4 value propositions | Profile configuration |

### Example Output

```markdown
# Jane Smith

**Senior Software Engineer / Technical Consultant**  
Remote (UK)  
Email: jane@example.com  
GitHub: https://github.com/janesmith  
Portfolio: https://janesmith.dev

---

## Summary

Senior, language-agnostic software engineer with 15+ years of experience delivering 
production systems across diverse stacks and domains. Strong systems thinker with a 
proven ability to select the right tools for complex problems.

---

## Core Strengths

- Language-agnostic backend and systems engineering  
- System design and architectural ownership  
- Pragmatic problem decomposition  
- Legacy stabilisation and modernisation  
- Product-focused delivery  
- Linux-first operational mindset  

---

## Languages & Technologies

**Languages:** Go, Python, Ruby, JavaScript, Shell  
**Frontend:** React, Vue.js  
**Systems:** Kubernetes, AWS, Terraform, CI/CD  

---

## Selected Experience

### TechCorp Inc.
*Jan 2020 - Present*

- Designed and implemented distributed event processing system handling 1M+ events/day
- Led migration from monolith to microservices, reducing deployment time from hours to minutes
- Mentored team of 5 engineers on system design best practices

### StartupXYZ
*Mar 2017 - Dec 2019*

- Owned end-to-end development of core API platform serving 500K users
- Architected multi-region deployment strategy achieving 99.99% uptime

---

## What I Bring

- Languages as tools, not identity  
- Calm handling of complexity  
- Clear thinking under constraints  
- Long-term maintainability focus  

---

**References available on request.**
```

## Configuring Your Profile

### Access Profile Settings

1. From main menu, select **Configure System**
2. Select **Profile** domain
3. Edit fields and save

### Profile Fields

| Field | Description | Example |
|-------|-------------|---------|
| **Name** | Your full name | Jane Smith |
| **Email** | Contact email | jane@example.com |
| **Title** | Professional title | Senior Software Engineer |
| **Location** | Where you're based | Remote (UK) |
| **GitHub** | GitHub profile URL | https://github.com/janesmith |
| **Portfolio** | Personal site URL | https://janesmith.dev |
| **Languages** | Programming languages (comma-separated) | Go, Python, Ruby, JavaScript |
| **Frontend** | Frontend tech (comma-separated) | React, Vue.js |
| **Systems** | Infrastructure/DevOps (comma-separated) | Kubernetes, AWS, Terraform |

### Default Values

If any field is empty, the narrative CV uses sensible defaults. This ensures the CV always looks complete, even if you haven't configured all fields.

## Generation Workflow

### Step-by-Step

1. **Generate CV** from main menu
2. **Select Profile** (e.g., Staff Engineer)
3. **Select Audience** (e.g., Hiring Manager)
4. **Select Role Emphasis** - Choose **Language-Agnostic**
5. **Select Length Format** - Choose Full, Standard, or Short (Ultra-Short uses Highlights structure)
6. **Review Preview** - See how it looks
7. **Export** - Choose Text or Markdown format

### Variant Selection

The narrative structure is used by these variants:

| Variant | Length | Best For |
|---------|--------|----------|
| `language_agnostic_full` | 3+ pages | Complete history, internal promotion |
| `language_agnostic_standard` | 2-3 pages | Most job applications |
| `language_agnostic_short` | 1-2 pages | Quick reviews, networking |

**Note**: `language_agnostic_ultra_short` uses the Highlights structure for a one-page format.

### Export Formats

| Format | Narrative Support | Use Case |
|--------|-------------------|----------|
| **Text** | ✅ Full support | Email, pasting |
| **Markdown** | ✅ Full support | GitHub, docs |
| **YAML** | ❌ Uses Standard | Data interchange |

**Note**: YAML always uses Standard structure because it's a data format, not a presentation format.

## Customization Tips

### Tailoring Core Strengths

Edit your profile to emphasize different strengths for different roles:

**For Backend-Heavy Roles**:
- Language-agnostic backend and systems engineering
- API design and microservices architecture
- Database optimization and data modeling

**For Leadership Roles**:
- Technical team leadership and mentorship
- Cross-functional collaboration
- Strategic technical planning

### Tailoring Technologies

Group your skills strategically:

**Languages**: List your strongest/most recent first
**Frontend**: Include frameworks you're comfortable with
**Systems**: Emphasize infrastructure you've owned

### Tailoring What I Bring

These should be your unique value propositions:

- What makes you different from other candidates?
- What do you consistently bring to teams?
- What's your working philosophy?

## Best Practices

### 1. Keep It Concise

- Summary: 2-3 sentences max
- Core Strengths: 6 bullets
- What I Bring: 4 bullets
- Experience: Only high-confidence achievements

### 2. Be Specific

Instead of:
> "Good at backend development"

Write:
> "Language-agnostic backend engineering with 15+ years across Ruby, Go, and Python"

### 3. Show, Don't Tell

Let your Selected Experience demonstrate your claims:
- If you claim "system design expertise", show system design achievements
- If you claim "pragmatic problem-solving", show practical solutions

### 4. Update Regularly

As your career evolves:
- Update profile fields quarterly
- Re-export CVs after significant achievements
- Keep technologies current

## Troubleshooting

### Profile Not Appearing in Export

1. Verify profile is configured in **Configure System → Profile**
2. Regenerate CV after profile changes
3. Use Text or Markdown format (not YAML)

### Too Few Experience Bullets

The narrative structure only shows bullets with confidence >= 0.75. To get more bullets:
- Improve event quality (clearer ownership, more detail)
- Add supporting facts to events
- Ensure events use action verbs ("Led", "Designed", "Owned")

### Technologies Look Wrong

1. Edit profile in **Configure System → Profile**
2. Use comma-separated format: "Go, Python, Ruby"
3. Save and re-export

For more troubleshooting, see [CV Troubleshooting Guide](CV_TROUBLESHOOTING.md).

## Comparison: Standard vs Narrative

| Aspect | Standard | Narrative |
|--------|----------|-----------|
| **Sections** | Experience, Projects, Skills | Core Strengths, Technologies, What I Bring |
| **Focus** | Chronological history | Pragmatic expertise |
| **Profile Header** | CV name only | Full contact info |
| **Technologies** | Skills list | Grouped by type |
| **Experience Filter** | All bullets | High-confidence only (>= 0.75) |
| **Value Proposition** | Implicit | Explicit ("What I Bring") |

## Related Documentation

- [CV Variants Guide](CV_VARIANTS_GUIDE.md) - Complete guide to all 16 variants
- [CV Generation Guide](CV_GENERATION_GUIDE.md) - Complete CV generation workflow
- [CV Troubleshooting](CV_TROUBLESHOOTING.md) - Common issues and solutions
- [Custom CV Format Proposal](../CUSTOM_CV_FORMAT_PROPOSAL.md) - Implementation details

---

**Document Version**: 1.1  
**Last Updated**: 2026-01-09
