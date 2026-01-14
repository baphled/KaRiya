---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Skills Management Guide

## Overview

KaRiya's Skills Management feature allows you to define, organize, and track your technical and professional skills. Skills can be associated with career events and will be used in CV generation to highlight your expertise.

## Accessing Skills Management

From the main menu:
1. Press `s` to open **Manage Skills**
2. Or select "Manage Skills" from the menu

## Managing Skills

### Viewing Your Skills

The **Skills List** shows all your skills organized by category:

```
┌─────────────────────────────────────────────┐
│ Manage Skills                               │
├─────────────────────────────────────────────┤
│                                             │
│ Backend (5 skills)                          │
│ ▶ Go              Advanced    5 yrs  [12]  │
│   Ruby            Expert      8 yrs  [25]  │
│   Python          Intermediate 3 yrs [8]   │
│   PostgreSQL      Advanced    6 yrs  [18]  │
│   Redis           Intermediate 4 yrs [10]  │
│                                             │
│ Frontend (3 skills)                         │
│   React           Advanced    4 yrs  [15]  │
│   TypeScript      Intermediate 3 yrs [12]  │
│   CSS             Advanced    6 yrs  [20]  │
│                                             │
│ DevOps (4 skills)                           │
│   Kubernetes      Advanced    3 yrs  [10]  │
│   Docker          Expert      5 yrs  [22]  │
│   Jenkins         Intermediate 2 yrs [7]   │
│   Terraform       Beginner    1 yr   [3]   │
└─────────────────────────────────────────────┘

j/k or ↑/↓ - Navigate | Enter - View detail | n - Add skill
e - Edit | d - Delete | f - Filter | s - Sort | x - Clear filters | Esc - Back
```

**Columns**:
- **Name**: Skill name (e.g., "Go", "Kubernetes")
- **Level**: Your proficiency level (Beginner, Intermediate, Advanced, Expert)
- **Years Used**: How many years you've used this skill
- **[Count]**: Number of events using this skill

### Filtering and Sorting Skills

KaRiya provides powerful filtering and sorting capabilities to help you focus on specific skills.

#### Filtering Skills

Press `f` from the Skills List to open the **Filter Menu**:

```
┌─────────────────────────────────────────────┐
│ Filter Skills                               │
├─────────────────────────────────────────────┤
│                                             │
│ Filter by Category:                         │
│ ▶ All Categories                            │
│   Backend                                   │
│   Frontend                                  │
│   DevOps                                    │
│   Database                                  │
│   Cloud                                     │
│                                             │
│ Filter by Level:                            │
│   All Levels                                │
│   Beginner                                  │
│   Intermediate                              │
│   Advanced                                  │
│   Expert                                    │
│                                             │
│ Other Filters:                              │
│   Used skills only (has events)             │
└─────────────────────────────────────────────┘

j/k or ↑/↓ - Navigate | Enter - Apply filter
u - Used skills only (quick) | Esc - Cancel
```

**Filter Options**:

1. **By Category**: Show only skills in a specific category
   - Example: Select "Backend" to see only Go, Ruby, Python, etc.
   - Useful for focusing on one technical area

2. **By Level**: Show only skills at a specific proficiency level
   - Example: Select "Expert" to see your strongest skills
   - Useful for CV preparation targeting senior roles

3. **Used Skills Only**: Show only skills associated with events
   - Press `u` for quick access to this filter
   - Hides skills you've defined but not yet used in events
   - Useful for finding which skills need event associations

**Applying Filters**:
- Navigate with `j/k` or arrow keys
- Press `Enter` to apply the selected filter
- Press `u` for quick "Used skills only" filter
- Press `Esc` to cancel without filtering

#### Sorting Skills

Press `s` from the Skills List to open the **Sort Menu**:

```
┌─────────────────────────────────────────────┐
│ Sort Skills                                 │
├─────────────────────────────────────────────┤
│                                             │
│ ▶ Name (A-Z)                                │
│   Name (Z-A)                                │
│   Most Used (event count)                   │
│   Least Used (event count)                  │
│   Recently Used (last used date)            │
│   Oldest Used (last used date)              │
│   Category                                  │
└─────────────────────────────────────────────┘

j/k or ↑/↓ - Navigate | Enter - Apply sort
e - Most used (quick) | Esc - Cancel
```

**Sort Options**:

1. **By Name (A-Z / Z-A)**: Alphabetical sorting (default)
   - A-Z: Standard alphabetical order
   - Z-A: Reverse alphabetical order

2. **By Event Count**: Sort by how often you've used each skill
   - Most Used: Skills with highest event count first
   - Least Used: Skills with lowest event count first
   - Press `e` for quick "Most used" sort
   - Useful for identifying your most frequently used skills

3. **By Last Used Date**: Sort by when you last used each skill
   - Recently Used: Most recent usage first
   - Oldest Used: Oldest usage first
   - Useful for CV preparation (highlight recent experience)

4. **By Category**: Group skills by category, then sort by name
   - Shows all Backend skills together, then Frontend, etc.
   - Default view in the Skills List

**Applying Sorting**:
- Navigate with `j/k` or arrow keys
- Press `Enter` to apply the selected sort
- Press `e` for quick "Most used" sort
- Press `Esc` to cancel without sorting

#### Clearing Filters and Sorting

When filters or custom sorting are active, the footer shows:

```
j/k or ↑/↓ - Navigate | Enter - View detail | n - Add skill
e - Edit | d - Delete | f - Filter | s - Sort | x - Clear filters | Esc - Back
```

Press `x` to **clear all filters and sorting**, returning to the default view (grouped by category, sorted by name).

**Active Filter Indicator**:

When filters are active, you'll see an indicator at the top of the list:

```
┌─────────────────────────────────────────────┐
│ Manage Skills                               │
│ 🔍 Filtered: Backend, Advanced level       │
├─────────────────────────────────────────────┤
│                                             │
│ Backend (3 skills)                          │
│ ▶ Go              Advanced    5 yrs  [12]  │
│   PostgreSQL      Advanced    6 yrs  [18]  │
│   Redis           Advanced    4 yrs  [10]  │
└─────────────────────────────────────────────┘
```

#### Common Filtering/Sorting Workflows

**Workflow 1: Find Skills to Add to Events**
1. Press `f` to open filter menu
2. Select "Used skills only"
3. Review skills without events
4. Press `x` to show all skills
5. Press `Esc` to return to main menu

**Workflow 2: Prioritize Skills for CV**
1. Press `s` to open sort menu
2. Select "Most Used (event count)" (or press `e` for quick access)
3. Review top skills with highest event counts
4. These are your strongest skills to highlight

**Workflow 3: Focus on Recent Experience**
1. Press `s` to open sort menu
2. Select "Recently Used (last used date)"
3. Review skills you've used recently
4. Useful for roles requiring current experience

**Workflow 4: Audit a Specific Category**
1. Press `f` to open filter menu
2. Select specific category (e.g., "DevOps")
3. Review all skills in that category
4. Identify gaps or skills to develop

### Adding a New Skill

1. **From Skills List**: Press `n`
2. **Fill out the form**:

```
┌─────────────────────────────────────────────┐
│ Add New Skill                               │
├─────────────────────────────────────────────┤
│                                             │
│ Name: _____________________                 │
│ (e.g., "Go", "Kubernetes", "React")        │
│                                             │
│ Category: _____________________             │
│ Suggestions: backend, frontend, devops,     │
│              database, cloud, tooling       │
│                                             │
│ Level (optional): _____________________     │
│ Options: beginner, intermediate, advanced,  │
│          expert                             │
│                                             │
│ Years Used (optional): _____               │
│ (0-50)                                      │
│                                             │
│ [Submit]  [Cancel]                          │
└─────────────────────────────────────────────┘

Tab/Shift+Tab - Navigate fields
Ctrl+S - Save | Esc - Cancel
```

**Field Details**:
- **Name** (required): 1-100 characters, must be unique
- **Category** (required): 1-50 characters, suggested categories:
  - `backend` - Server-side languages (Go, Ruby, Python, Java, etc.)
  - `frontend` - UI technologies (React, Vue, Angular, CSS, etc.)
  - `devops` - Infrastructure (Kubernetes, Docker, Jenkins, etc.)
  - `database` - Database systems (PostgreSQL, MySQL, MongoDB, etc.)
  - `cloud` - Cloud platforms (AWS, GCP, Azure)
  - `mobile` - Mobile development (Swift, Kotlin, React Native, etc.)
  - `tooling` - Development tools (Git, VS Code, etc.)
  - `other` - Uncategorized skills
- **Level** (optional): beginner, intermediate, advanced, or expert
- **Years Used** (optional): 0-50 years

3. **Submit**: Press `Ctrl+S` or navigate to Submit and press `Enter`

### Viewing Skill Details

Press `Enter` on a skill to view full details:

```
┌─────────────────────────────────────────────┐
│ Skill Detail: Go                            │
├─────────────────────────────────────────────┤
│                                             │
│ Name:          Go                           │
│ Category:      backend                      │
│ Level:         Advanced                     │
│ Years Used:    5                            │
│                                             │
│ Event Count:   12 events                    │
│ Last Used:     2024-01-07                   │
│                                             │
│ Created:       2024-01-01 10:30:00          │
│ Updated:       2024-01-07 15:45:00          │
│                                             │
└─────────────────────────────────────────────┘

Enter - View events | e - Edit | d - Delete | Esc - Back
```

**Fields**:
- **Event Count**: Number of events where this skill is used
- **Last Used**: Most recent event date with this skill
- **Created/Updated**: Timestamps for tracking

### Viewing Events Using a Skill

From the **Detail View**, press `Enter` to see all events using this skill:

```
┌─────────────────────────────────────────────┐
│ Events Using Skill: Go                      │
├─────────────────────────────────────────────┤
│                                             │
│ Date       | Event                 | Company│
│────────────┼──────────────────────┼────────│
│ 2024-01-07 │ Architected microse… │ TechCo │
│ 2024-01-05 │ Implemented REST AP… │ TechCo │
│ 2023-12-20 │ Built CLI tool for…  │ TechCo │
│ 2023-12-15 │ Optimized API perfo… │ TechCo │
│ 2023-11-30 │ Led Go migration pr… │ OldCo  │
│                                             │
│ Showing 5 of 12 events                      │
│                                             │
└─────────────────────────────────────────────┘

j/k or ↑/↓ - Navigate events | Esc - Back to detail
```

**Sorted by date** (most recent first)

### Editing a Skill

**From List View**:
1. Navigate to skill
2. Press `e` to edit

**From Detail View**:
1. Press `e` to edit

The edit form is identical to the add form, but pre-filled with existing values.

### Deleting a Skill

**From List View**:
1. Navigate to skill
2. Press `d` to delete

**From Detail View**:
1. Press `d` to delete

**Confirmation Dialog**:

```
┌─────────────────────────────────────────────┐
│ Delete Skill?                               │
├─────────────────────────────────────────────┤
│                                             │
│ Skill: Go                                   │
│ Category: backend                           │
│                                             │
│ This skill is used by 12 events.            │
│                                             │
│ Deleting will remove this skill from all    │
│ events. This action cannot be undone.       │
│                                             │
│ Are you sure you want to delete this skill? │
│                                             │
└─────────────────────────────────────────────┘

y - Confirm deletion | n/Esc - Cancel
```

**Warning**: Deleting a skill removes it from all associated events!

## Associating Skills with Events

### During Event Capture

When capturing a new event, you can associate skills:

**Quick Capture Mode**:
```
┌─────────────────────────────────────────────┐
│ Capture Event - Quick                       │
├─────────────────────────────────────────────┤
│                                             │
│ Text: Built REST API with Go                │
│                                             │
│ Skills (optional): [          ]             │
│ Available: Go, Ruby, PostgreSQL, Redis      │
│                                             │
└─────────────────────────────────────────────┘
```

**Manual Capture Mode**:
```
┌─────────────────────────────────────────────┐
│ Capture Event - Manual                      │
├─────────────────────────────────────────────┤
│                                             │
│ Text: Built REST API with Go                │
│ Date: 2024-01-07                            │
│ Company: TechCorp                           │
│ Project: Platform API                       │
│                                             │
│ Skills (optional): [          ]             │
│ Available: Go, Ruby, PostgreSQL, Redis      │
│                                             │
│ Tags: technical;api                         │
│ Categories: Technical                       │
│                                             │
└─────────────────────────────────────────────┘
```

**Multi-select**: Use arrow keys and space to select multiple skills.

### Via Metadata Editor

Edit skills for existing events:

1. Navigate to event in timeline
2. Press `m` to edit metadata
3. Update Skills field
4. Save changes

### Via CSV Import

Import events with skills from CSV files:

**CSV Format**:
```csv
Text,Date,Categories,Tags,Project,Company,Skills
"Built REST API with Go",2024-01,Technical,technical;api,Platform,TechCorp,"Go;PostgreSQL;Redis"
"Architected microservices",2024-02,Technical,technical;architecture,Platform,TechCorp,"Go;Kubernetes;Docker"
```

**Auto-Creation**: Skills that don't exist are automatically created with category `"other"`.

**Refine After Import**:
1. Press `s` to open Manage Skills
2. Find auto-created skills (category: "other")
3. Edit each skill to set proper category
4. Skills are now categorized correctly

See [CSV_IMPORT_GUIDE.md](./CSV_IMPORT_GUIDE.md) for details.

## Keyboard Shortcuts

### Skills List

| Key | Action |
|-----|--------|
| `j` / `k` or `↑` / `↓` | Navigate skills |
| `Enter` | View skill detail |
| `n` | Add new skill |
| `e` | Edit selected skill (quick) |
| `d` | Delete selected skill (quick) |
| `f` | Open filter menu |
| `s` | Open sort menu |
| `x` | Clear all filters and sorting |
| `Esc` | Back to main menu |

### Filter Menu

| Key | Action |
|-----|--------|
| `j` / `k` or `↑` / `↓` | Navigate filter options |
| `Enter` | Apply selected filter |
| `u` | Quick: Used skills only |
| `Esc` | Cancel without filtering |

### Sort Menu

| Key | Action |
|-----|--------|
| `j` / `k` or `↑` / `↓` | Navigate sort options |
| `Enter` | Apply selected sort |
| `e` | Quick: Most used (event count desc) |
| `Esc` | Cancel without sorting |

### Skill Detail View

| Key | Action |
|-----|--------|
| `Enter` | View events using this skill |
| `e` | Edit skill |
| `d` | Delete skill |
| `Esc` | Back to skills list |

### Events View

| Key | Action |
|-----|--------|
| `j` / `k` or `↑` / `↓` | Navigate events |
| `Esc` | Back to skill detail |

### Add/Edit Form

| Key | Action |
|-----|--------|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Enter` | Submit (when on submit button) |
| `Ctrl+S` | Save (from any field) |
| `Esc` | Cancel |

### Delete Confirmation

| Key | Action |
|-----|--------|
| `y` | Confirm deletion |
| `n` / `Esc` | Cancel |

## Best Practices

### Skill Naming

✅ **Good Examples**:
- "Go" (programming language)
- "Kubernetes" (platform)
- "PostgreSQL" (database)
- "React" (framework)
- "Git" (tool)

❌ **Avoid**:
- "I know Go" (too verbose)
- "go language" (lowercase, inconsistent)
- "Ruby on Rails" (use "Ruby" and "Rails" separately)
- "AWS, GCP, Azure" (create 3 separate skills)

### Categorization

**Backend Skills**: Go, Ruby, Python, Java, Node.js, Elixir
**Frontend Skills**: React, Vue, Angular, TypeScript, CSS, HTML
**DevOps Skills**: Kubernetes, Docker, Jenkins, Terraform, Ansible
**Database Skills**: PostgreSQL, MySQL, MongoDB, Redis, Elasticsearch
**Cloud Skills**: AWS, GCP, Azure, Heroku, DigitalOcean
**Mobile Skills**: Swift, Kotlin, React Native, Flutter
**Tooling Skills**: Git, VS Code, IntelliJ, Vim
**Other**: Agile, Scrum, Code Review, Mentoring

### Skill Levels

- **Beginner**: Learning, basic understanding, needs guidance
- **Intermediate**: Productive, can work independently on typical tasks
- **Advanced**: Deep knowledge, can handle complex scenarios, mentor others
- **Expert**: Industry-recognized expertise, thought leader, architect-level

**Tip**: Be honest! It's better to be accurate than to overstate proficiency.

### Years Used

Track how long you've used each skill:
- **1-2 years**: Recent skills, actively learning
- **3-5 years**: Solid experience, core competencies
- **5-10 years**: Deep expertise, mature understanding
- **10+ years**: Veteran, likely expert level

**Tip**: Count total years, not just current role. If you used Go at 3 companies over 6 years, enter 6 years.

## Skills in CV Generation

Skills are used in CV generation to:

1. **Skills Section**: Show skills grouped by category
2. **Technology Focus**: Filter events by technology
3. **Proficiency Display**: Highlight expertise levels
4. **Relevance Scoring**: Prioritize relevant skills for target role

**Coming Soon** (Task 40): Technology-focused CV generation will allow selecting technologies to emphasize, filtering bullets by skills, and determining focus areas from skill categories.

## Common Workflows

### Workflow 1: Initial Skill Setup

**Goal**: Create your skills inventory from scratch

1. **Brainstorm skills**: List all technologies you've used
2. **Add to KaRiya**: Press `s`, then `n` for each skill
3. **Categorize**: Set proper categories (backend, frontend, devops, etc.)
4. **Set levels**: Assign proficiency levels (beginner, intermediate, advanced, expert)
5. **Associate with events**: Go through events and add skill associations

### Workflow 2: CSV Import with Skills

**Goal**: Import career events with skills from a CSV file

1. **Prepare CSV**: Add Skills column with semicolon-separated skill names
2. **Import**: Press `i` from main menu, select CSV file
3. **Review import**: Verify events imported successfully
4. **Refine skills**: Press `s`, find auto-created skills (category: "other")
5. **Categorize**: Edit each skill to set proper category
6. **Verify**: View skill details to see event counts

### Workflow 3: Enriching Existing Events

**Goal**: Add skills to events that don't have them yet

1. **Open skills management**: Press `s` to create/review your skills
2. **Navigate timeline**: Browse your career events
3. **Edit metadata**: Press `m` on each event
4. **Add skills**: Select relevant skills from multi-select field
5. **Save**: Press `Enter` to save changes
6. **Verify**: Check skill detail views to see event associations

### Workflow 4: Preparing for CV Generation

**Goal**: Ensure skills are ready for technology-focused CVs

1. **Review skills inventory**: Press `s` to view all skills
2. **Verify categories**: Ensure skills are categorized correctly
3. **Update proficiency**: Edit skills to set accurate levels
4. **Check event counts**: View skill details to see usage frequency
5. **Identify gaps**: Look for missing skills in your events
6. **Enrich events**: Add skills to events via metadata editor
7. **Generate CV**: Use skills to create technology-focused CVs

## Troubleshooting

### Issue: Duplicate skill names

**Symptom**: Error when adding a skill: "Skill with this name already exists"

**Cause**: Skill names must be unique (case-insensitive)

**Solution**:
- Check existing skills for similar names
- Use consistent naming (e.g., "JavaScript" not "javascript" or "JS")
- If truly different, use distinguishing names (e.g., "React" vs "React Native")

### Issue: Can't delete a skill

**Symptom**: Delete option is disabled or fails

**Cause**: Skills in use by events cannot be deleted

**Solution**:
1. View skill detail to see event count
2. Press `Enter` to see events using this skill
3. Edit each event to remove the skill association
4. Once event count is 0, you can delete the skill

**Alternative**: Keep the skill but mark it as unused by not associating with new events

### Issue: Auto-created skills from CSV import

**Symptom**: Many skills with category "other" after CSV import

**Cause**: CSV import auto-creates skills that don't exist

**Solution**:
1. Press `s` to open Manage Skills
2. Navigate through skills with category "other"
3. Press `e` to edit each skill
4. Update category to proper value (backend, frontend, devops, etc.)
5. Repeat for all auto-created skills

**Tip**: Do this immediately after CSV import for best results

### Issue: Skill levels inconsistent

**Symptom**: Unsure what level to set for each skill

**Cause**: Subjective proficiency assessment

**Solution**:
- **Beginner**: Can read and understand code, basic tasks with guidance
- **Intermediate**: Can build features independently, common patterns known
- **Advanced**: Deep understanding, can architect solutions, debug complex issues
- **Expert**: Industry authority, speak at conferences, write books/articles

**Alternative**: Leave level empty if uncertain. It's optional!

## Integration with Other Features

### Event Capture

- Skills field available in quick and manual capture modes
- Multi-select from existing skills
- Optional field (events without skills work fine)

### Metadata Editor

- Edit skills for existing events
- Add/remove skill associations
- Changes reflected in skill detail views

### CSV Import

- Skills column in CSV files (optional)
- Auto-creates skills that don't exist (category: "other")
- Matches existing skills by name (case-insensitive)

### Browse Timeline

- Events show associated skills
- Filter events by skill (coming soon)

### CV Generation

- Skills section in CVs (grouped by category)
- Technology-focused CV generation (Task 40, coming soon)

## See Also

- [CSV_IMPORT_GUIDE.md](./CSV_IMPORT_GUIDE.md) - CSV import with skills
- [CSV_FORMAT_GUIDE.md](./CSV_FORMAT_GUIDE.md) - CSV format specification
- [CLI_GUIDE.md](./CLI_GUIDE.md) - General CLI usage
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - General troubleshooting

