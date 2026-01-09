# CV Generation Guide

## Overview

The CV Generation feature in KaRiya allows you to transform your raw career events and facts into professional, audience-specific CV views. This guide explains how to use the feature, understand its rules, and export your generated CVs.

### Key Concepts

- **CV Configuration**: A YAML file that defines how to generate a CV (target role, audience, filters)
- **CV View**: An in-memory CV generated from a configuration (NOT stored in the database)
- **Bullets**: Individual achievement statements in your CV, traced to source events and facts
- **Traceability**: Full tracking of which events and facts contributed to each bullet
- **Ephemeral Generation**: CVs are generated fresh each time, always reflecting your latest data

## Getting Started

### 1. Access CV Management

From the main menu, select **"Manage CV Configs"** to view and manage your CV configurations.

### 2. Create a New CV Configuration

1. Select **"New Config"** or press `n` from the config list
2. Fill in the configuration form:
   - **CV Name**: A unique identifier for this CV (e.g., "Staff Engineer CV")
   - **Target Role**: The role level you're applying for (Principal, Staff, EM, Senior IC)
   - **Target Audience**: Who will read this CV (Hiring Manager, Recruiter, Peer)
   - **Date Range** (optional): Filter events to a specific timeframe
   - **Companies** (optional): Include events only from specific companies
   - **Tags** (optional): Include only events with specific tags
   - **Competencies** (optional): Highlight specific skill areas

3. Press `Enter` to save the configuration

### 3. Generate a CV

1. From the config list, select your configuration
2. Press `Enter` or select **"Generate CV"**
3. KaRiya will generate your CV and display it in preview mode

### 4. Preview Your CV

The CV preview shows:
- **Metadata**: Role, audience, and generation date
- **Sections**: Experience, Core Competencies, Professional Summary
- **Bullets**: Achievement statements with source indicators

### 5. Explore Sources

Press `Enter` on any bullet to see which events and facts contributed to it. This helps you understand:
- Which career moments shaped this achievement statement
- How confident KaRiya is in the bullet (confidence score)
- Why the bullet was included (inclusion reason)

### 6. Select CV Structure

Before generation, you can choose how your CV is structured:

- **Standard**: Traditional CV format with Experience, Projects, Skills, and Summary sections. Best for most job applications.
- **Narrative**: Language-agnostic format emphasizing pragmatic expertise with Core Strengths, Technologies, and "What I Bring" sections. Best for senior engineers with cross-domain experience.

Use arrow keys or `j`/`k` to select, then press `Enter` to generate.

### 7. Export Your CV

After generation, you can export your CV in multiple formats:
- **Text Export**: Plain text format, suitable for pasting
- **Markdown Export**: Markdown format with formatting
- **YAML Export**: Machine-readable data format
- **Copy to Clipboard**: Quick copy for pasting elsewhere

Exports are saved to `$HOME/.kariya/cv_exports/` by default.

**Note**: YAML export always uses Standard structure since it's a data format, not a presentation format.

## CV Configuration Format

CV configurations are stored as YAML files in `$HOME/.kariya/cv_configs/`. You can also edit these files directly in a text editor.

### Example Configuration

```yaml
name: "Staff Engineer CV"
targetRole: "Staff"
targetAudience:
  - "HiringManager"
  - "Peer"
eventFilters:
  startDate: "2023-01-01"
  endDate: "2024-12-31"
  companies:
    - "Acme Corp"
    - "Tech Startup"
  tags:
    - "technical"
    - "leadership"
  categories:
    - "technical"
    - "product"
createdAt: "2024-12-01T10:30:00Z"
updatedAt: "2024-12-01T10:30:00Z"
```

### Configuration Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique identifier for this CV config |
| `targetRole` | string | Yes | One of: `Principal`, `Staff`, `EM`, `SeniorIC` |
| `targetAudience` | array | Yes | One or more of: `HiringManager`, `Recruiter`, `Peer` |
| `eventFilters` | object | No | Filters to apply to events |
| `eventFilters.startDate` | string | No | ISO 8601 date (e.g., "2023-01-01") |
| `eventFilters.endDate` | string | No | ISO 8601 date (e.g., "2024-12-31") |
| `eventFilters.companies` | array | No | List of company names to include |
| `eventFilters.tags` | array | No | List of tags to include |
| `eventFilters.categories` | array | No | List of categories to include |
| `createdAt` | timestamp | Auto | Creation timestamp |
| `updatedAt` | timestamp | Auto | Last update timestamp |

### Valid Values

**Target Roles**:
- `Principal` - Principal Engineer level
- `Staff` - Staff Engineer level
- `EM` - Engineering Manager level
- `SeniorIC` - Senior Individual Contributor level

**Target Audiences**:
- `HiringManager` - Hiring manager or decision maker
- `Recruiter` - Recruiter or talent acquisition
- `Peer` - Peer or technical colleague

**Event Tags**:
- `project` - Project work
- `achievement` - Notable achievements
- `leadership` - Leadership activities
- `technical` - Technical work
- `consulting` - Consulting or advisory work
- `research` - Research or exploration
- `product` - Product-related work
- `mentoring` - Mentoring or coaching

**Event Categories**:
- `technical` - Technical contributions
- `leadership` - Leadership contributions
- `product` - Product contributions
- `consulting` - Consulting contributions
- `research` - Research contributions
- `mentoring` - Mentoring contributions

## Bullet Generation Rules

### Inclusion Criteria

A bullet is included in your CV only if it meets ALL of these criteria:

1. **Traces to at least one event**: Every bullet must have a source event
2. **Single claim**: The bullet makes exactly one claim (not multiple unrelated achievements)
3. **No aspirational language**: Doesn't use future tense or aspirational language (e.g., "will", "aims to", "plans to")
4. **No inferred metrics**: Doesn't claim metrics that weren't explicitly stated in source events
5. **No role inflation**: Doesn't claim responsibilities above your actual role level
6. **Repeated signals**: Preferred if multiple events support the same bullet

### Exclusion Criteria

Bullets are excluded if they:

- **Use aspirational keywords**: "will", "aims to", "plans to", "seeks to", "strives to", "attempts to", "hopes to", "want to", "like to"
- **Infer metrics**: Claim numbers/percentages that weren't in the source events
- **Inflate role**: Claim ownership of work you only contributed to, or claim strategic decisions you only advised on
- **Lack clear ownership**: Are vague about who did what
- **Multiple unrelated claims**: Try to pack too many achievements into one bullet

### Ranking Priority

Bullets are ranked by priority, then by signal strength and recency:

1. **Ownership** (Highest Priority)
   - "Designed and built X from scratch"
   - "Owned X end-to-end"
   - "Led X initiative"

2. **Contribution**
   - "Contributed to X"
   - "Helped build X"
   - "Supported X"

3. **Strategy**
   - "Defined strategy for X"
   - "Planned X"
   - "Architected X"

4. **Execution**
   - "Implemented X"
   - "Executed X"
   - "Delivered X"

5. **Outcome**
   - "Improved X by Y%"
   - "Reduced X by Y%"
   - "Increased X by Y%"

6. **Activity** (Lowest Priority)
   - "Worked on X"
   - "Participated in X"
   - "Attended X"

### Scoring Algorithm

Each bullet receives a confidence score (0.0 to 1.0) based on:

- **Priority Level**: Ownership bullets score higher than Activity bullets
- **Signal Strength**: Multiple sources boost score
- **Recency**: Older events score slightly lower (temporal decay)
- **Audience Fit**: Bullets matching target audience preferences score higher

Example scoring:
- Ownership bullet, 2 sources, recent: 0.95
- Contribution bullet, 1 source, recent: 0.75
- Activity bullet, 1 source, older: 0.45

## Role-Specific Generation

Each target role has specific bullet caps and preferences:

### Principal Role
- **Bullet Cap**: 3-4 bullets maximum
- **Preference**: Strategic impact, organizational influence, system-wide improvements
- **Focus**: Long-term vision, cross-functional impact, mentorship of senior engineers
- **Compression**: Removes lower-priority bullets to stay within cap

### Staff Role
- **Bullet Cap**: 4-5 bullets maximum
- **Preference**: Technical depth, complex problem-solving, technical leadership
- **Focus**: Technical excellence, architectural decisions, mentorship
- **Compression**: Removes lower-priority bullets to stay within cap

### EM (Engineering Manager) Role
- **Bullet Cap**: 3-4 bullets maximum
- **Preference**: Team building, process improvement, business alignment
- **Focus**: Team growth, delivery, cross-functional collaboration
- **Compression**: Removes lower-priority bullets to stay within cap

### Senior IC Role
- **Bullet Cap**: 4-5 bullets maximum
- **Preference**: Technical contributions, impact, growth trajectory
- **Focus**: Technical leadership without management, mentorship, innovation
- **Compression**: Removes lower-priority bullets to stay within cap

### Compression Logic

When you have more bullets than the role allows:

1. Bullets are sorted by confidence score (highest first)
2. Lower-confidence bullets are removed until within the cap
3. Older role transitions are compressed first
4. At least one bullet per section is preserved

Example: Staff role with 8 bullets and cap of 5
- Top 5 bullets by score are kept
- Bottom 3 are removed
- All sections remain represented

## Audience-Specific Filtering

Each target audience has different preferences for what makes a strong bullet:

### Hiring Manager Audience
- **Focus**: Business impact, outcomes, measurable results
- **Prefers**: Bullets showing ROI, revenue impact, cost savings
- **Values**: Delivery speed, quality, customer satisfaction
- **Includes**: Metrics, timelines, business context
- **Excludes**: Deep technical details, internal processes

### Recruiter Audience
- **Focus**: Skills, competencies, career growth, achievements
- **Prefers**: Bullets showing breadth of skills, seniority level, progression
- **Values**: High-level achievements, role fit, growth trajectory
- **Includes**: Skills demonstrated, scope of work, impact
- **Excludes**: Overly technical details, internal jargon

### Peer Audience
- **Focus**: Technical depth, problem-solving, collaboration
- **Prefers**: Bullets showing technical excellence, architectural thinking, mentorship
- **Values**: Technical leadership, innovation, knowledge sharing
- **Includes**: Technical approach, design decisions, mentorship
- **Excludes**: Business metrics, high-level strategy

### Multi-Audience Generation

When you select multiple audiences:

- A bullet is included if it's relevant to **ANY** of the selected audiences
- Bullets are prioritized by how many audiences they fit
- Bullets fitting all audiences score higher than those fitting only one

## Traceability System

Every bullet in your CV is fully traceable to its source events and facts.

### Viewing Sources

1. In the CV preview, navigate to a bullet
2. Press `Enter` to view sources
3. You'll see:
   - **Source Events**: Which career events contributed to this bullet
   - **Source Facts**: Which extracted facts (from burst analysis) contributed
   - **Confidence Score**: How confident KaRiya is in this bullet
   - **Inclusion Reason**: Why this bullet was included

### Source Event Details

For each source event, you can see:
- **Event Text**: The original event description
- **Date**: When the event occurred
- **Company**: Associated company
- **Tags**: Tags applied to the event
- **Categories**: Categories assigned to the event

### Confidence Score Interpretation

- **0.9-1.0**: Very confident - strong ownership, multiple sources, recent
- **0.7-0.9**: Confident - clear ownership or strong contribution, recent
- **0.5-0.7**: Moderate - contribution or activity, possibly older
- **0.3-0.5**: Lower confidence - activity-level work or older events
- **Below 0.3**: Rarely included in CVs

## Compression Logic Details

When KaRiya needs to compress bullets to meet role caps:

### Step 1: Evaluate Candidates
- All bullets below the cap threshold are candidates for removal
- Bullets are sorted by confidence score (lowest first)

### Step 2: Remove Lower-Priority Bullets
- Removes bullets with Activity or Execution priority first
- Then removes Outcome bullets
- Then removes Strategy bullets
- Never removes Ownership or Contribution bullets (unless absolutely necessary)

### Step 3: Preserve Section Representation
- At least one bullet per section is always kept
- If a section would have zero bullets, the highest-scoring bullet is preserved

### Step 4: Temporal Compression
- When multiple bullets have similar scores, older ones are compressed first
- Ensures CV represents your most recent work

### Example Compression

**Before Compression (7 bullets)**:
1. Ownership bullet, score 0.95 ✓ Keep
2. Contribution bullet, score 0.85 ✓ Keep
3. Strategy bullet, score 0.80 ✓ Keep
4. Execution bullet, score 0.75 ✓ Keep
5. Outcome bullet, score 0.70 → Remove
6. Activity bullet, score 0.60 → Remove
7. Activity bullet, score 0.55 → Remove

**After Compression (4 bullets - Staff cap)**:
- Keeps top 4 bullets by score
- Removes lower-priority bullets first
- Result: 1 Ownership + 1 Contribution + 1 Strategy + 1 Execution

## Export Formats

### Plain Text Export

Format:
```
CV NAME - TARGET ROLE FOR TARGET AUDIENCE
Generated: 2024-12-01

SECTION NAME
- Bullet 1
- Bullet 2
- Bullet 3

ANOTHER SECTION
- Bullet A
- Bullet B
```

**Use Cases**:
- Pasting into email or messaging
- Importing into word processors
- Quick sharing with colleagues
- Simple, clean format

**File Location**: `$HOME/.kariya/cv_exports/cv_name_YYYY-MM-DD_HH-MM-SS.txt`

### Markdown Export

Format:
```markdown
# CV NAME - TARGET ROLE FOR TARGET AUDIENCE

**Generated**: 2024-12-01

## SECTION NAME

- Bullet 1
- Bullet 2
- Bullet 3

## ANOTHER SECTION

- Bullet A
- Bullet B
```

**Use Cases**:
- Publishing to GitHub or documentation sites
- Using in markdown-based resume builders
- Preserving formatting in markdown tools
- Version control in git

**File Location**: `$HOME/.kariya/cv_exports/cv_name_YYYY-MM-DD_HH-MM-SS.md`

### YAML Export

Format:
```yaml
name: "Staff Engineer CV"
targetRole: "staff"
targetAudience: "hiring_manager"
generatedAt: "2024-12-01T10:30:00Z"
sections:
  - type: "experience"
    title: "Experience"
    bullets:
      - text: "Led team of 5 engineers..."
        confidence: 0.95
```

**Use Cases**:
- Machine processing and integration
- Data interchange between tools
- Programmatic CV manipulation

**File Location**: `$HOME/.kariya/cv_exports/cv_name_YYYY-MM-DD_HH-MM-SS.yaml`

**Note**: YAML always uses Standard structure regardless of selection.

### Copy to Clipboard

- Copies the generated CV text to your system clipboard
- Same format as plain text export
- No file saved
- Quick sharing or pasting

## CV Structures

KaRiya supports two CV structures that determine how your content is organized.

### Standard Structure

The traditional CV format with these sections:

| Section | Content |
|---------|---------|
| **Summary** | Professional summary statement |
| **Experience** | Work history with company, dates, and achievements |
| **Projects** | Notable projects and contributions |
| **Skills** | Technical skills and competencies |

**Best for**:
- Most job applications
- Traditional company cultures
- When specific role experience matters

### Narrative Structure

A language-agnostic format for professionals emphasizing pragmatic expertise:

| Section | Content |
|---------|---------|
| **Profile Header** | Name, title, location, contact info |
| **Summary** | Professional summary with language-agnostic emphasis |
| **Core Strengths** | Key competencies (6 bullet points) |
| **Languages & Technologies** | Languages, Frontend, Systems |
| **Selected Experience** | High-confidence achievements (>= 0.75) |
| **What I Bring** | Value propositions (4 bullet points) |

**Best for**:
- Senior engineers with cross-domain experience
- Emphasizing "languages as tools, not identity"
- Pragmatic, outcome-focused professionals
- Roles requiring broad technical expertise

### Configuring Your Profile for Narrative CVs

For narrative CVs, you can customize your profile:

1. Go to **Configure System** from main menu
2. Select **Profile** domain
3. Edit these fields:
   - **Name**: Your full name
   - **Email**: Your email address
   - **Title**: Professional title (e.g., "Senior Software Engineer")
   - **Location**: Your location (e.g., "Remote (UK)")
   - **GitHub**: Your GitHub profile URL
   - **Portfolio**: Your portfolio/website URL
   - **Languages**: Programming languages (comma-separated)
   - **Frontend**: Frontend technologies (comma-separated)
   - **Systems**: Systems/infrastructure expertise (comma-separated)
4. Save changes

These values are used when exporting narrative CVs. Empty fields use sensible defaults.

For more details, see the [Narrative CV Guide](NARRATIVE_CV_GUIDE.md).

## Keyboard Shortcuts

### Config Management Screen
| Key | Action |
|-----|--------|
| `j` / `↓` | Move down in list |
| `k` / `↑` | Move up in list |
| `Enter` | Generate CV from selected config |
| `n` | Create new config |
| `e` | Edit selected config |
| `d` | Delete selected config (with confirmation) |
| `Esc` | Back to main menu |
| `?` | Show help |

### Config Editor Screen
| Key | Action |
|-----|--------|
| `Tab` | Move to next field |
| `Shift+Tab` | Move to previous field |
| `↑` / `↓` | Navigate within multi-select fields |
| `Space` | Toggle multi-select item |
| `Enter` | Save configuration |
| `Esc` | Cancel without saving |
| `?` | Show help |

### CV Preview Screen
| Key | Action |
|-----|--------|
| `j` / `↓` | Move to next bullet |
| `k` / `↑` | Move to previous bullet |
| `→` | Move to next section |
| `←` | Move to previous section |
| `Enter` | View sources for selected bullet |
| `e` | Export CV (shows export menu) |
| `Esc` | Back to config list |
| `?` | Show help |

### Source Event Tracer Screen
| Key | Action |
|-----|--------|
| `j` / `↓` | Scroll down |
| `k` / `↑` | Scroll up |
| `Esc` | Back to CV preview |
| `?` | Show help |

## Common Workflows

### Workflow 1: Create Your First CV

1. Open KaRiya and go to "Manage CV Configs"
2. Press `n` to create new config
3. Enter CV name: "My First CV"
4. Select target role: "Staff"
5. Select target audience: "HiringManager"
6. Leave filters blank to include all events
7. Press `Enter` to save
8. Select your new config and press `Enter` to generate
9. Review the generated CV
10. Press `e` to export as markdown
11. Open the exported file in your text editor

### Workflow 2: Create Role-Specific CVs

1. Create multiple configs with the same filters but different target roles:
   - "Staff Engineer CV" (target: Staff)
   - "Principal Engineer CV" (target: Principal)
   - "Manager CV" (target: EM)

2. For each config, generate and review
3. Export each to markdown
4. Tailor the exported versions for specific job applications

### Workflow 3: Create Audience-Specific CVs

1. Create two configs with same role but different audiences:
   - "Staff for Hiring Manager" (audience: HiringManager)
   - "Staff for Recruiter" (audience: Recruiter)

2. Generate both and compare
3. Notice how the same events produce different bullets for different audiences

### Workflow 4: Use Filters for Targeted CVs

1. Create config "Recent Technical Work"
2. Set date range: Last 2 years
3. Set tags: "technical"
4. Set companies: Current company only
5. Generate to see focused CV of recent technical achievements

## Tips & Best Practices

### 1. Event Quality Matters

CV generation quality depends on event quality:
- **Good events**: Specific, measurable, with clear ownership
- **Poor events**: Vague, activity-level, missing context
- **Tip**: Review and improve your events before CV generation

### 2. Use Tags Strategically

Tags help filter events for specific CVs:
- Tag technical work with "technical"
- Tag leadership work with "leadership"
- Tag product work with "product"
- **Tip**: Consistently tag events for easier filtering

### 3. Create Multiple Configurations

Don't try to make one CV for all purposes:
- Create "Staff for FAANG" config
- Create "Staff for Startup" config
- Create "Staff for PM collaboration" config
- **Tip**: Different audiences and companies need different CVs

### 4. Review Generated Bullets

Always review the generated CV:
- Check that bullets accurately represent your work
- Verify source events are correct
- Look for any aspirational language
- **Tip**: If a bullet doesn't feel right, improve the source event

### 5. Understand Compression

When you hit the bullet cap:
- Lower-priority bullets are removed first
- Older work is deprioritized
- **Tip**: If important work is compressed out, improve those events to increase their score

### 6. Use Traceability

Take advantage of full source tracking:
- Click on any bullet to see its sources
- Understand why bullets were ranked certain ways
- **Tip**: Use this to improve events that should score higher

### 7. Export for Customization

Exported CVs are starting points, not final products:
- Export to markdown or text
- Customize the exported version for specific roles
- Keep the original configuration for regeneration
- **Tip**: Regenerate and re-export when you add new events

## Troubleshooting

### Problem: "Not enough bullets generated"

**Causes**:
- Events don't meet inclusion criteria
- Events use aspirational language
- Events claim inferred metrics
- Events lack clear ownership

**Solutions**:
1. Review your events for quality
2. Remove aspirational language ("will", "aims to", etc.)
3. Only include metrics that were explicitly in source events
4. Make ownership clear in event descriptions
5. Try different role or audience filters

### Problem: "Unexpected bullets in CV"

**Causes**:
- Events have vague descriptions
- Multiple unrelated claims in one event
- Role inflation in event descriptions

**Solutions**:
1. Check the source events (press `Enter` on bullet)
2. Improve the source event descriptions
3. Split multi-claim events into separate events
4. Clarify your actual role vs. contributions

### Problem: "Generated CV is too short"

**Causes**:
- Bullet cap for your role is low
- Many events don't meet inclusion criteria
- Filters are too restrictive

**Solutions**:
1. Check the bullet cap for your target role
2. Review which events are being filtered out
3. Try removing date range or company filters
4. Improve lower-scoring events

### Problem: "Export file not found"

**Causes**:
- Directory doesn't exist
- File permissions issue
- Wrong export location

**Solutions**:
1. Check `$HOME/.kariya/cv_exports/` exists
2. Verify write permissions
3. Try "Copy to Clipboard" instead
4. Check system error message for details

## Advanced Usage

### Manual Configuration Editing

You can edit YAML config files directly:

```bash
# Edit config in your editor
nano ~/.kariya/cv_configs/my_cv.yaml

# KaRiya will reload the config next time you use it
```

### Configuration Reuse

You can copy a config file and modify it:

```bash
cp ~/.kariya/cv_configs/staff_cv.yaml ~/.kariya/cv_configs/principal_cv.yaml
# Edit the new file to change target role and other settings
```

### Batch Export

To export all your CVs:

1. Open "Manage CV Configs"
2. For each config:
   - Select and generate
   - Press `e` to export
3. All exports are saved to `$HOME/.kariya/cv_exports/`

## Integration with Other KaRiya Features

### With Event Timeline

- Generate CV from selected events
- View events used in each bullet
- Trace back to original event

### With Burst Detection

- Facts from burst detection contribute to bullets
- See fact sources in traceability view
- Use burst-generated facts to improve bullet quality

### With Fact Extraction

- Extracted facts improve bullet confidence scores
- Multiple facts supporting a bullet increase score
- Fact sources are tracked in traceability

## Next Steps

1. **Create Your First Config**: Start with one target role and audience
2. **Generate and Review**: See how your events become bullets
3. **Explore Sources**: Click bullets to understand traceability
4. **Iterate on Events**: Improve events to improve CV quality
5. **Export and Customize**: Export and tailor for specific roles

---

**Document Version**: 1.1
**Last Updated**: 2026-01-09

