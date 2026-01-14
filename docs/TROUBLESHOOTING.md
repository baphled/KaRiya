---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya CLI Troubleshooting Guide

## Common Issues and Solutions

### Database & Persistence

#### Events disappear after closing the app

**Cause**: Using the default in-memory database

**Solution**:
```bash
# Specify a database file with --db flag
./kariya-cli --db ./events.db

# Or use absolute path
./kariya-cli --db ~/.local/share/kariya/events.db
```

**Explanation**: By default, KaRiya uses in-memory storage which is lost when the application closes. Use the `--db` flag to enable persistent SQLite storage.

#### "Error initializing database" message

**Cause**: Invalid database path or permission issues

**Solution**:
1. Check the path is writable:
```bash
# Test write permission
touch /path/to/database.db
```

2. Ensure directory exists:
```bash
mkdir -p ~/.local/share/kariya
./kariya-cli --db ~/.local/share/kariya/events.db
```

3. Check disk space:
```bash
df -h /path/to/database
```

#### Database file too large

**Cause**: Many events accumulated

**Solution**:
- Archive old events by creating a new database file
- Use filters to list and archive events before cleanup
- For large deployments, consider PostgreSQL (future enhancement)

---

### Form & Input Issues

#### Form field navigation not working

**Cause**: Terminal too small or keyboard layout issues

**Solution**:
1. Ensure terminal width is at least 80 columns
2. Ensure terminal height is at least 24 rows
3. Try different key combinations:
   - **Tab** or **Shift+Tab**: Navigate between fields
   - **Arrow Keys**: Select options and lists
   - **Enter**: Submit form
   - **Escape**: Cancel operation

#### Text input showing strange characters

**Cause**: Terminal doesn't support UTF-8

**Solution**:
1. Set terminal encoding to UTF-8:
```bash
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8
```

2. Verify encoding:
```bash
locale
```

3. Try different terminal emulator (e.g., iTerm2, Kitty, Alacritty)

#### Special characters cutting off text

**Cause**: Terminal rendering issue with wide characters

**Solution**:
- Use ASCII characters instead of unicode
- Or use a terminal that better supports unicode (Kitty, WezTerm)
- Report on GitHub with terminal version info

#### Date field not accepting input

**Cause**: Incorrect date format

**Accepted formats**:
- ISO format: `2025-12-24`
- Relative: `today`, `1 week ago`, `2 months ago`, `1 year ago`
- Short form: `-1d` (not yet supported, use `1 day ago`)

**Solution**: Try different date format:
```
✓ 2025-12-24
✓ today
✓ 1 day ago
✓ 2 weeks ago
✓ 3 months ago
✗ 12/24/2025 (not supported)
✗ Dec 24, 2025 (not supported)
```

#### Cannot select more than 8 tags

**Cause**: System limit for tag count per event

**Solution**: This is intentional to keep events focused. If you need more tags:
1. Create separate events for different aspects
2. Combine related tags (e.g., "backend-optimization" instead of separate tags)

#### Event text truncated at 2000 characters

**Cause**: Character limit enforced by domain model

**Solution**:
- Create separate events for different aspects of the work
- Use concise, impactful text for each event
- Focus on key accomplishments rather than detailed descriptions

---

### Display & Appearance

#### Colors not showing correctly

**Cause**: Terminal doesn't support 256 colors

**Solution**:
1. Check terminal color support:
```bash
echo $TERM
```

2. Set to a supported value:
```bash
export TERM=xterm-256color  # Standard
export TERM=tmux-256color   # If using tmux
```

3. Try different terminal emulator

#### Text overlapping or misaligned

**Cause**: Terminal font issue or small window

**Solution**:
1. Increase terminal window size
2. Use monospace font (Courier, Monospace, Inconsolata)
3. Check font size (recommend 12-14pt)
4. Ensure terminal width >= 80 chars, height >= 24 lines

#### Help screen not fully visible

**Cause**: Terminal height too small

**Solution**:
- Maximize terminal window vertically
- Minimum height: 24 lines
- Recommended: 30+ lines for comfortable reading

---

### Performance

#### Application startup is slow

**Cause**: Usually disk I/O if using SQLite on slow drive

**Solution**:
1. Use fast storage (SSD instead of HDD)
2. Check disk I/O:
```bash
iostat -x 1 5
```

3. Profile startup:
```bash
time ./kariya-cli --help
```

#### List or search is slow with large database

**Cause**: Database performance with many events

**Solution**:
1. Database file corruption:
```bash
sqlite3 events.db "PRAGMA integrity_check;"
```

2. Optimize database:
```bash
sqlite3 events.db "VACUUM;"
```

3. Create indexes (future feature):
- We'll add indexing on frequently filtered columns

#### Memory usage increasing over time

**Cause**: Possible memory leak or large result set

**Solution**:
1. Restart the application periodically
2. Avoid listing all events at once with large databases
3. Use filters to limit results (by date range, tags, etc.)
4. Report memory profiles on GitHub if issue persists

---

### Navigation & Workflow

#### Can't navigate between screens

**Cause**: Keyboard shortcut not working

**Shortcuts**:
- `h` or `Home`: Go to home screen
- `c` or `C`: Capture event
- `l` or `L`: List events
- `q` or `Q` or `Ctrl+C`: Quit
- `Backspace`: Go back to previous screen

**Solution**:
1. Check if keys are being captured by terminal
2. Disable terminal key bindings conflicting with KaRiya
3. Try numeric keypad or different keyboard layout

#### Can't exit the application

**Cause**: Quit key not responding

**Solution**:
1. Press `q` followed by `Enter`
2. Or use `Ctrl+C` (force quit)
3. Check terminal isn't capturing `q` key

#### Form submission appears to hang

**Cause**: Slow database write or network issue

**Solution**:
1. Wait a few seconds (writing to disk takes time)
2. Check database is accessible:
```bash
ls -la ./events.db
```

3. Check disk space:
```bash
df -h .
```

---

### Capture Strategy Issues

#### Date Handling

Both capture strategies accept any date (past or present). There are no date restrictions.

**Capture Strategies Explained**:
- **Quick Capture**: Rapid event logging with minimal fields visible
  - Event text and date are always visible
  - Press 't' to toggle optional fields (Company, Project, Tags, Categories)
- **Manual Capture**: Detailed entry with all fields visible by default
  - All fields visible: Event, Date, Company, Project, Tags, Categories
  - Press 't' to hide optional fields if desired

#### Event date in the future is rejected

**Cause**: System validation prevents future-dated events

**Solution**:
- Use correct date (today or earlier)
- Check system clock if needed:
```bash
date
```

---

### Advanced Troubleshooting

#### Enable debug logging

**How**:
- Application logs to stdout during operation
- Watch for ERROR, WARN, and INFO messages
- These help identify issues

**Interpretation**:
```
ERROR: text cannot be empty
  → Event text field is required and empty

WARN: Timeline event outside 30-day window
  → Event date is too old for Timeline mode

INFO: Event captured successfully
  → Event was saved to database
```

#### Check database integrity

```bash
# Verify SQLite database
sqlite3 events.db ".tables"

# Check event count
sqlite3 events.db "SELECT COUNT(*) FROM career_events;"

# Export events for backup
sqlite3 events.db ".mode json" "SELECT * FROM career_events;" > backup.json
```

#### File locations

**Default locations**:
- Linux: `~/.local/share/kariya/events.db`
- macOS: `~/Library/Application Support/kariya/events.db`
- Windows: `%APPDATA%\kariya\events.db`

**Current directory**:
```bash
./events.db  # Default when using --db ./events.db
```

---

### Getting Help

#### Help within the CLI

Press `h` in the application to open interactive help with:
- Overview of KaRiya
- Explanation of capture modes
- Tagging and organization tips
- Keyboard shortcuts reference
- Best practices guide

#### Documentation

- **CLI Guide**: `docs/CLI_GUIDE.md` - Comprehensive usage guide
- **Main README**: `README.md` - Project overview
- **AGENTS.md**: Architecture and developer documentation

#### Report Issues

When reporting an issue on GitHub:

1. Include version:
```bash
./kariya-cli --version
```

2. Include environment:
```bash
go version
echo $TERM
uname -a
```

3. Include exact error message and steps to reproduce
4. Attach screenshot if display issue
5. Provide relevant log output

---

### Configuration Tips

#### Set default database

Create a shell alias:
```bash
# Add to ~/.bashrc or ~/.zshrc
alias kariya='kariya-cli --db ~/.local/share/kariya/events.db'
```

#### Start in specific mode

```bash
# Always start in Timeline mode
alias kariya-timeline='kariya-cli --mode timeline'

# Always start with events list
alias kariya-list='kariya-cli --list'
```

#### Script multiple events

```bash
# (Future feature - would use API or batch mode)
```

---

## FAQ

**Q: Where are my events stored?**
A: Use `--db` flag to specify location. Default is in-memory (lost on close).

**Q: Can I backup my events?**
A: Export SQLite database or use:
```bash
cp events.db events.db.backup
```

**Q: Can I share events with others?**
A: Future feature - currently events are local to your database.

**Q: Is my data encrypted?**
A: No, use encrypted filesystem if needed (e.g., LUKS, BitLocker).

**Q: How many events can the system handle?**
A: Tested with 10,000+ events, should handle 100k+ with optimization.

**Q: Can I edit existing events?**
A: Not yet - future feature in Phase 3+.

**Q: What terminal is best for KaRiya?**
A: Any terminal with UTF-8 support. Recommended: iTerm2 (Mac), Windows Terminal (Windows), Kitty (Linux).

---

**Last Updated**: 2025-12-24  
**Version**: 1.0  
**Status**: Complete for Phase 3

For the latest information, visit: https://github.com/baphled/kariya
