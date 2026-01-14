# Consulting CV Guide

## Overview

The **Consulting CV structure** is designed for consulting professionals who need to showcase client engagements, delivery capability, and adaptability across diverse technical environments. This guide explains when and how to use consulting CVs effectively.

> **Note**: The consulting structure is automatically selected when you choose **Consulting** role emphasis with Full, Standard, or Short length formats. Ultra-Short always uses the Highlights structure. See the [CV Variants Guide](CV_VARIANTS_GUIDE.md) for the complete variant system.

## When to Use Consulting Structure

### Best For

- **Consulting firm applications** - When applying to consulting companies (Big 4, boutique firms, etc.)
- **Client-facing technical roles** - Roles requiring stakeholder management and communication
- **Freelance/contract work** - Showcasing diverse project experience
- **Technical advisory positions** - Where rapid assessment and recommendations matter
- **Roles requiring adaptability** - Demonstrating success across different environments

### Less Suitable For

- **Permanent product company roles** - Where deep product ownership matters more
- **Highly specialized technical roles** - Where depth in one area trumps breadth
- **Junior positions** - Where career trajectory matters more than client diversity
- **Roles with strict ATS screening** - Some systems expect traditional CV formats

## Consulting CV Structure

### Section Overview

| Section | Purpose | Source |
|---------|---------|--------|
| **Profile Header** | Contact info and professional identity | Profile configuration |
| **Summary** | 2-3 sentence professional summary | CV sections or defaults |
| **Client Engagements** | Work history grouped by company/client | Generated from events |
| **Technical Capabilities** | Skills and tools (optional) | Profile configuration |
| **What I Bring** | Value propositions | Profile configuration |

### Example Output

```markdown
# Jane Smith

**Senior Consulting Engineer**  
Remote (UK)  
Email: jane@example.com  
GitHub: https://github.com/janesmith  
Portfolio: https://janesmith.dev

---

## Summary

Consulting engineer with 12+ years delivering technical solutions across diverse client 
environments. Specializes in rapid assessment, architecture modernization, and team enablement. 
Proven track record of translating complex technical challenges into actionable roadmaps.

---

## Client Engagements

### Acme Corporation
*Jan 2022 - Present*

- Led technical assessment resulting in 18-month modernization roadmap
- Designed and implemented event-driven architecture handling 500K events/day
- Mentored internal team on distributed systems best practices
- Reduced deployment cycle from 2 weeks to 2 hours through CI/CD implementation

### TechStartup Inc
*Jun 2021 - Dec 2021*

- Rapid 2-week assessment of legacy monolith codebase
- Delivered microservices migration strategy with clear milestones
- Implemented automated testing framework increasing coverage from 20% to 85%
- Enabled team to self-serve deployments through infrastructure automation

### Financial Services Corp
*Jan 2020 - May 2021*

- Architected real-time fraud detection system processing 10M transactions/day
- Led security audit and remediation project achieving SOC 2 compliance
- Trained 25 developers on secure coding practices

---

## Technical Capabilities

**Languages:** Go, Python, Ruby, JavaScript, Shell  
**Cloud:** AWS, GCP, Azure  
**Infrastructure:** Kubernetes, Terraform, Docker  
**Practices:** CI/CD, TDD, Agile, DevOps

---

## What I Bring

- Rapid technical assessment and roadmapping  
- Clear communication with technical and non-technical stakeholders  
- Pragmatic solutions within budget and timeline constraints  
- Knowledge transfer and team enablement  

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
| **Title** | Professional title | Senior Consulting Engineer |
| **Location** | Where you're based | Remote (UK) |
| **GitHub** | GitHub profile URL | https://github.com/janesmith |
| **Portfolio** | Personal site URL | https://janesmith.dev |
| **Languages** | Programming languages | Go, Python, Ruby |
| **Systems** | Infrastructure expertise | Kubernetes, AWS, Terraform |
| **What I Bring** | Value propositions (4 items) | Rapid assessment, Clear communication, ... |

### Default Values

If any field is empty, the consulting CV uses sensible defaults for consulting professionals.

## Generation Workflow

### Step-by-Step

1. **Generate CV** from main menu
2. **Select Profile** (e.g., Staff Engineer, Senior IC)
3. **Select Audience** (e.g., Hiring Manager, Recruiter)
4. **Select Role Emphasis** - Choose **Consulting**
5. **Select Length Format** - Choose Full, Standard, or Short
6. **Review Preview** - See how it looks
7. **Export** - Choose Text or Markdown format

### Variant Selection

The consulting structure is used by these variants:

| Variant | Length | Min Confidence | Best For |
|---------|--------|----------------|----------|
| `consulting_full` | 3+ pages | 0.50 | Complete engagement history |
| `consulting_standard` | 2-3 pages | 0.60 | Most consulting applications |
| `consulting_short` | 1-2 pages | 0.70 | Quick reviews, networking |

**Note**: `consulting_ultra_short` uses the Highlights structure for a one-page format.

### Export Formats

| Format | Consulting Support | Use Case |
|--------|-------------------|----------|
| **Text** | Full support | Email, pasting |
| **Markdown** | Full support | Proposals, GitHub |
| **YAML** | Uses Standard | Data interchange |

**Note**: YAML always uses Standard structure because it's a data format.

## Customization Tips

### Tailoring Client Engagements

The consulting structure groups your experience by company/client. To maximize impact:

1. **Use clear company names** - Ensure your events have accurate company fields
2. **Focus on outcomes** - Each bullet should show impact, not just activity
3. **Quantify where possible** - Numbers make consulting achievements concrete
4. **Show client diversity** - Variety demonstrates adaptability

### Tailoring What I Bring

For consulting roles, emphasize:

- **Rapid assessment capability** - How quickly you can understand and act
- **Communication skills** - Working with technical and non-technical stakeholders
- **Pragmatism** - Delivering within real-world constraints
- **Knowledge transfer** - Enabling client teams, not creating dependency

### Confidence Thresholds

Consulting variants have lower confidence thresholds than other role emphases:

| Variant | Min Confidence |
|---------|----------------|
| `consulting_full` | 0.50 |
| `consulting_standard` | 0.60 |
| `consulting_short` | 0.70 |
| `consulting_ultra_short` | 0.80 |

This reflects that consulting work often involves breadth over depth, so more diverse experiences are included.

## Best Practices

### 1. Group Related Work

If you worked with the same client multiple times, ensure events share the same company name so they're grouped together in the Client Engagements section.

### 2. Show Engagement Scope

Include context about engagement type:
- Assessment/Discovery
- Implementation
- Advisory/Strategy
- Training/Enablement

### 3. Demonstrate Adaptability

Show you've succeeded across:
- Different industries
- Different tech stacks
- Different team sizes
- Different engagement lengths

### 4. Highlight Soft Skills

Consulting success depends on:
- Stakeholder management
- Clear communication
- Expectation setting
- Conflict resolution

### 5. Include Measurable Outcomes

Strong consulting bullets include:
- Timeline improvements
- Cost reductions
- Performance gains
- Risk mitigations

## Troubleshooting

### Client Engagements Not Grouping

**Cause**: Events have different company names (e.g., "Acme Corp" vs "Acme Corporation")

**Solution**: Standardize company names in your events

### Too Few Engagements Shown

**Cause**: Events don't meet confidence threshold

**Solution**: 
1. Try a longer length format (Full has lower threshold)
2. Improve event quality to boost confidence scores
3. Add supporting facts to events

### Technical Capabilities Missing

**Cause**: Profile not configured or empty

**Solution**:
1. Configure profile in **Configure System → Profile**
2. Add languages, systems, and tools
3. Regenerate CV

For more troubleshooting, see [CV Troubleshooting Guide](CV_TROUBLESHOOTING.md).

## Comparison: Consulting vs Other Structures

| Aspect | Standard | Narrative | Consulting |
|--------|----------|-----------|------------|
| **Focus** | Chronological history | Pragmatic expertise | Client engagements |
| **Experience Section** | "Experience" | "Selected Experience" | "Client Engagements" |
| **Best For** | Most roles | Language-agnostic | Consulting roles |
| **Skills Presentation** | Skills section | Technologies by type | Technical Capabilities |
| **Value Proposition** | Implicit | "What I Bring" | "What I Bring" |

## Related Documentation

- [CV Variants Guide](CV_VARIANTS_GUIDE.md) - Complete guide to all 16 variants
- [CV Generation Guide](CV_GENERATION_GUIDE.md) - Complete CV generation workflow
- [Narrative CV Guide](NARRATIVE_CV_GUIDE.md) - For language-agnostic emphasis
- [CV Troubleshooting](CV_TROUBLESHOOTING.md) - Common issues and solutions

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-09
