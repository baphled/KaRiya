# CV Generation Troubleshooting Guide

This guide helps you diagnose and resolve common issues with KaRiya's CV generation feature.

## Common Issues and Solutions

### Issue 1: "Not Enough Bullets Generated"

**Symptoms**: Your CV has fewer bullets than expected, or looks too sparse

**Root Causes**:
1. Events don't meet inclusion criteria
2. Events use aspirational language
3. Events claim inferred metrics
4. Events lack clear ownership
5. Filters are too restrictive

**Diagnosis Steps**:
1. Check how many events match your filters
   - Go to List Events
   - Apply the same filters as your CV config
   - Count the events

2. Check event quality
   - Review the events that should contribute bullets
   - Look for aspirational language ("will", "aims to", "plans to")
   - Verify metrics are explicitly stated in events

3. Check role-specific bullet caps
   - Principal: 3-4 bullets max
   - Staff: 4-5 bullets max
   - EM: 3-4 bullets max
   - SeniorIC: 4-5 bullets max

**Solutions**:

**Solution A: Improve Event Quality**
- Remove aspirational language from events
- Only include metrics that were explicitly stated
- Make ownership clear (use "Led", "Designed", "Owned")
- Add more context and detail to events
- Example improvement:
  ```
  Before: "Will improve performance in the future"
  After: "Improved API response time from 500ms to 100ms"
  ```

**Solution B: Adjust Filters**
- Remove date range filters if too restrictive
- Remove company filters to include all companies
- Remove tag filters to include more events
- Remove category filters to expand scope
- Try regenerating with no filters first

**Solution C: Check Inclusion Criteria**
- Ensure events have clear ownership (not just "participated")
- Verify events make single claims (not multiple unrelated achievements)
- Check that events don't claim metrics not in the source
- Ensure events aren't aspirational

**Solution D: Verify Role and Audience**
- Try different role selections (staff has higher cap than principal)
- Try different audience combinations
- Multi-audience CVs often include more bullets

### Issue 2: "Unexpected Bullets in CV"

**Symptoms**: CV contains bullets that don't accurately represent your work

**Root Causes**:
1. Source events have vague descriptions
2. Events have multiple unrelated claims
3. Role inflation in event descriptions
4. Misinterpretation of event content

**Diagnosis Steps**:
1. Click on the unexpected bullet (press Enter in preview)
2. View the source events contributing to it
3. Check if the source event accurately describes the work
4. Verify ownership claims are accurate

**Solutions**:

**Solution A: Improve Source Events**
- Go back to List Events
- Find and edit the source event
- Make the description more specific
- Clarify your actual role and contribution
- Remove vague language
- Example improvement:
  ```
  Before: "Helped with project"
  After: "Contributed performance optimization to core service, reducing load time 30%"
  ```

**Solution B: Split Multi-Claim Events**
- If an event makes multiple unrelated claims, split it
- Create separate events for each achievement
- Each bullet should represent one claim
- Example:
  ```
  Before: "Led team and improved performance and designed API"
  After:
    Event 1: "Led team of 4 engineers on project"
    Event 2: "Improved API performance by 40%"
    Event 3: "Designed new REST API for data platform"
  ```

**Solution C: Clarify Role and Contributions**
- Use clear ownership language: "Led", "Designed", "Owned", "Architected"
- Distinguish between ownership and contribution
- Don't claim strategic decisions you only advised on
- Example improvement:
  ```
  Before: "Worked on performance optimization"
  After: "Led performance optimization effort, reducing load time from 8s to 2s"
  ```

### Issue 3: "Generated CV is Too Short"

**Symptoms**: CV is much shorter than expected, even though you have many events

**Root Causes**:
1. Bullet cap for your role is low
2. Many events don't meet inclusion criteria
3. Filters are too restrictive
4. Compression removed too many bullets

**Diagnosis Steps**:
1. Check your role's bullet cap
   - Principal: 3-4 max
   - Staff: 4-5 max
   - EM: 3-4 max
   - SeniorIC: 4-5 max

2. Check which events are filtered out
   - Edit your CV config
   - Temporarily remove filters
   - Regenerate to see if more bullets appear

3. Check event quality scores
   - Go to metadata review
   - Look at data quality scores
   - Lower-quality events may not generate bullets

**Solutions**:

**Solution A: Increase Bullet Cap**
- Change your target role to one with higher cap
- Staff and SeniorIC have 4-5 bullets (vs Principal's 3-4)
- Generate multiple CVs with different roles

**Solution B: Improve Event Quality**
- Go to Metadata Review
- Enrich events with better descriptions
- Add missing company, project, tags, or categories
- Improve data quality scores
- Regenerate CV after improvements

**Solution C: Adjust Filters**
- Remove date range to include older work
- Remove company filters to include all companies
- Remove tag/category filters to expand scope
- Try with no filters first, then add back incrementally

**Solution D: Check Compression**
- If you have many events, compression might remove bullets
- Try different audience combinations
- Multi-audience CVs often include more bullets
- Review which bullets were kept vs removed

### Issue 4: "Export File Not Found"

**Symptoms**: Export appears to complete but file isn't where expected

**Root Causes**:
1. Directory doesn't exist or isn't writable
2. File permissions issue
3. Wrong export location
4. Terminal output error not noticed

**Diagnosis Steps**:
1. Check if export directory exists
   ```bash
   ls -la ~/.kariya/cv_exports/
   ```

2. Check directory permissions
   ```bash
   ls -ld ~/.kariya/cv_exports/
   ```

3. Try creating the directory manually
   ```bash
   mkdir -p ~/.kariya/cv_exports/
   ```

4. Check if you have write permissions
   ```bash
   touch ~/.kariya/cv_exports/test.txt
   ```

**Solutions**:

**Solution A: Create Export Directory**
```bash
mkdir -p ~/.kariya/cv_exports/
chmod 755 ~/.kariya/cv_exports/
```

**Solution B: Fix Directory Permissions**
```bash
chmod 755 ~/.kariya/cv_exports/
chmod 644 ~/.kariya/cv_exports/*.txt
chmod 644 ~/.kariya/cv_exports/*.md
```

**Solution C: Use Clipboard Instead**
- If file export fails, use "Copy to Clipboard"
- Paste the CV into your text editor
- Save manually to your preferred location

**Solution D: Check Terminal Output**
- Look for error messages in the terminal
- Check if the app showed a success message
- Verify the export location shown in the message

### Issue 5: "Configuration Won't Save"

**Symptoms**: CV configuration doesn't save or disappears after restart

**Root Causes**:
1. Config directory doesn't exist
2. File permissions issue
3. Invalid YAML in config
4. Insufficient disk space

**Diagnosis Steps**:
1. Check if config directory exists
   ```bash
   ls -la ~/.kariya/cv_configs/
   ```

2. Check directory permissions
   ```bash
   ls -ld ~/.kariya/cv_configs/
   ```

3. Try to list configs
   - Go to "Manage CV Configs"
   - Check if any configs appear

4. Check YAML syntax if editing manually
   ```bash
   # YAML files should be valid
   # Use 2-space indentation
   # No tabs
   ```

**Solutions**:

**Solution A: Create Config Directory**
```bash
mkdir -p ~/.kariya/cv_configs/
chmod 755 ~/.kariya/cv_configs/
```

**Solution B: Check File Permissions**
```bash
chmod 755 ~/.kariya/cv_configs/
chmod 644 ~/.kariya/cv_configs/*.yaml
```

**Solution C: Verify YAML Syntax**
If editing configs manually:
- Use 2-space indentation (not tabs)
- Use proper YAML syntax
- Validate with online YAML validator
- Example:
  ```yaml
  name: "My CV"
  targetRole: "Staff"
  targetAudience:
    - "HiringManager"
  ```

**Solution D: Check Disk Space**
```bash
df -h ~
```

### Issue 6: "Configuration Shows Wrong Values"

**Symptoms**: CV config displays different values than what you entered

**Root Causes**:
1. Config file was overwritten
2. Manual YAML edit had syntax errors
3. Cached values not refreshed
4. Partial save due to error

**Diagnosis Steps**:
1. Check the YAML file directly
   ```bash
   cat ~/.kariya/cv_configs/my_config.yaml
   ```

2. Verify what values are displayed in the UI
   - Go to "Manage CV Configs"
   - Select the config and press 'e' to edit
   - Compare displayed values with file contents

3. Check if there are multiple config files
   ```bash
   ls ~/.kariya/cv_configs/
   ```

**Solutions**:

**Solution A: Re-edit Configuration**
1. Go to "Manage CV Configs"
2. Select the config
3. Press 'e' to edit
4. Correct the values
5. Press Enter to save

**Solution B: Delete and Recreate**
1. Go to "Manage CV Configs"
2. Select the config
3. Press 'd' to delete (with confirmation)
4. Press 'n' to create new
5. Re-enter all values carefully

**Solution C: Edit YAML Directly**
```bash
nano ~/.kariya/cv_configs/my_config.yaml
```
- Make sure syntax is correct
- Use proper indentation
- Save and close

### Issue 7: "Filtering Not Working"

**Symptoms**: CV still includes events that should be filtered out

**Root Causes**:
1. Filter values not saved correctly
2. Misunderstanding of how filters work
3. Events don't have the expected metadata
4. Filter values don't match event values exactly

**Diagnosis Steps**:
1. Check your CV config
   - Go to "Manage CV Configs"
   - Select config and press 'e'
   - Verify filter values are correct

2. Check event metadata
   - Go to List Events
   - View events that should be filtered
   - Verify they have the expected company/tags/categories

3. Try without filters
   - Edit config and clear all filters
   - Regenerate to see baseline

**Solutions**:

**Solution A: Verify Filter Values**
- Company names must match exactly (case-sensitive)
- Tags must be from allowed list
- Categories must be from allowed list
- Dates must be in ISO format (YYYY-MM-DD)

**Solution B: Update Event Metadata**
- Go to Metadata Review
- Add missing company, tags, or categories
- Make sure values match your filter criteria
- Regenerate CV after updates

**Solution C: Use Correct Filter Format**
- In YAML config:
  ```yaml
  eventFilters:
    companies:
      - "Exact Company Name"
    tags:
      - "technical"
      - "leadership"
    categories:
      - "technical"
  ```

**Solution D: Test Filters Incrementally**
- Start with no filters
- Add one filter at a time
- Regenerate after each change
- Identify which filter is problematic

### Issue 8: "Compression Removes Important Bullets"

**Symptoms**: Important bullets are removed when exceeding role cap

**Root Causes**:
1. Important bullets have lower confidence scores
2. Older work is deprioritized
3. Bullet priority doesn't match your expectations
4. Events lack supporting facts/signals

**Diagnosis Steps**:
1. Generate CV and check which bullets remain
2. Click on removed bullets to see their confidence scores
3. Understand why they scored lower
4. Check source events for quality

**Solutions**:

**Solution A: Improve Event Quality**
- Bullets with higher confidence are kept
- Improve lower-scoring event descriptions
- Add more detail and context
- Make ownership clearer
- Regenerate to see improvement

**Solution B: Choose Different Role**
- Different roles have different bullet caps
- Staff (4-5) has higher cap than Principal (3-4)
- Try generating for different roles
- Keep the CV with more bullets

**Solution C: Use Targeted Filters**
- Create focused CVs with filters
- "Recent Technical Work" (recent date range, technical tags)
- "Leadership Focus" (leadership tags)
- Fewer events = fewer compression needed

**Solution D: Add Facts to Improve Scores**
- Run burst detection on events
- Run fact extraction on bursts
- Facts improve bullet confidence scores
- Regenerate CV after facts are added

### Issue 9: "Traceability Shows Wrong Sources"

**Symptoms**: Bullet sources don't match the bullet content

**Root Causes**:
1. Events were edited after CV generation
2. Multiple events contributed to one bullet
3. Misunderstanding of how traceability works
4. Display bug (rare)

**Diagnosis Steps**:
1. Click on the bullet to view sources
2. Read the source event text carefully
3. Check if multiple events contributed
4. Verify the inclusion reason

**Solutions**:

**Solution A: Understand Multi-Source Bullets**
- A bullet can come from multiple events
- Click to see all sources
- Read all source events to understand the bullet
- This is intentional and shows signal strength

**Solution B: Check Event History**
- If events were edited, regenerate CV
- CV is always generated from current event state
- Older CV generations won't update automatically
- Generate fresh CV to get current sources

**Solution C: Review Inclusion Reason**
- View the source tracer (press Enter on bullet)
- Read the "Inclusion Reason" field
- This explains why the bullet was included
- Helps understand the ranking logic

### Issue 10: "Terminal Display Issues"

**Symptoms**: CV screens don't display correctly, text is cut off, or formatting is wrong

**Root Causes**:
1. Terminal window too small
2. Terminal doesn't support UTF-8
3. Terminal font issues
4. Terminal color support limited

**Diagnosis Steps**:
1. Check terminal size
   ```bash
   echo $COLUMNS x $LINES
   ```

2. Check UTF-8 support
   ```bash
   locale
   # Should show UTF-8
   ```

3. Try maximizing terminal window
4. Try different terminal emulator

**Solutions**:

**Solution A: Increase Terminal Size**
- Minimum width: 80 columns
- Minimum height: 24 lines
- Recommended: 120+ columns, 30+ lines
- Maximize your terminal window

**Solution B: Enable UTF-8**
```bash
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
```

**Solution C: Try Different Terminal**
- Try: iTerm2 (macOS), GNOME Terminal (Linux), Windows Terminal
- Avoid: Older terminals, limited color support
- Ensure terminal supports 256 colors

**Solution D: Disable Colors (Last Resort)**
- Some terminals don't support colors well
- Try redirecting output to text file instead
- Export CV to text/markdown format
- View exported file in text editor

## Performance Issues

### Issue: "CV Generation is Slow"

**Symptoms**: CV generation takes >2 seconds

**Root Causes**:
1. Large number of events (1000+)
2. Slow disk I/O
3. Complex filters requiring many comparisons
4. System resources constrained

**Solutions**:

**Solution A: Optimize Filters**
- Use date range filters to reduce event set
- Filter by company to focus scope
- Use tag filters to target specific work
- Fewer events = faster generation

**Solution B: Check System Resources**
```bash
# Check available memory
free -h

# Check disk space
df -h

# Check CPU usage
top
```

**Solution C: Generate Incrementally**
- Create focused CVs with narrow filters
- Rather than one big CV with all events
- Faster generation, more relevant results

**Solution D: Improve Event Quality**
- Better event descriptions = better ranking
- Less time spent on low-quality events
- Faster overall generation

## Getting More Help

### If Issue Persists

1. **Check Logs**: Look for error messages in terminal
2. **Verify Setup**: Run through setup steps again
3. **Try Defaults**: Reset to default configuration
4. **Restart App**: Close and reopen KaRiya
5. **Review Documentation**: Check CV_GENERATION_GUIDE.md

### Provide Feedback

If you encounter an issue not listed here:
1. Note the exact steps to reproduce
2. Check the terminal for error messages
3. Document what you expected vs what happened
4. Share your CV config (YAML file)
5. Include sample events that cause the issue

## Quick Reference

### Directory Locations
- CV Configs: `~/.kariya/cv_configs/`
- CV Exports: `~/.kariya/cv_exports/`
- Main Database: `~/.kariya/events.db`

### File Formats
- Configs: YAML format
- Exports: Text or Markdown format
- Timestamps: ISO 8601 format (YYYY-MM-DD)

### Role Bullet Caps
- Principal: 3-4 bullets
- Staff: 4-5 bullets
- EM: 3-4 bullets
- SeniorIC: 4-5 bullets

### Keyboard Shortcuts
- `j`/`k`: Navigate
- `Enter`: Select/confirm
- `e`: Edit config
- `d`: Delete config
- `n`: New config
- `Esc`: Back/cancel

---

**Document Version**: 1.0
**Last Updated**: 2026-01-02

