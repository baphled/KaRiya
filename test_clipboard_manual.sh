#!/bin/bash
#
# Manual test script to reproduce the clipboard issue
#
# This script will:
# 1. Create a test database with one event
# 2. Run kariya and attempt to export to clipboard
# 3. Show debug output
#

set -e

echo "=== Setting up test environment ==="

# Create temp directory
TEST_DIR=$(mktemp -d)
export HOME=$TEST_DIR
echo "Test directory: $TEST_DIR"

# Create .kariya directory
mkdir -p "$TEST_DIR/.kariya"

# Create a test database with one event
cat > "$TEST_DIR/create_test_event.sh" << 'EOF'
#!/bin/bash
sqlite3 ~/.kariya/events.db << 'SQL'
CREATE TABLE IF NOT EXISTS career_events (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    date TEXT NOT NULL,
    company TEXT,
    project TEXT,
    tags TEXT,
    categories TEXT,
    created_at TEXT,
    updated_at TEXT
);

INSERT INTO career_events (id, text, date, company, tags, created_at, updated_at)
VALUES ('test-1', 'Test event for clipboard', '2024-01-01', 'Test Corp', '["technical"]', datetime('now'), datetime('now'));
SQL
EOF

chmod +x "$TEST_DIR/create_test_event.sh"
"$TEST_DIR/create_test_event.sh"

echo "=== Test event created ==="
sqlite3 "$TEST_DIR/.kariya/events.db" "SELECT * FROM career_events;"

echo ""
echo "=== Now run ./kariya and try to export to clipboard ==="
echo "=== Watch for DEBUG messages on stderr ==="
echo "=== Navigate: Main Menu > Export Artifact > Events > JSON > Clipboard > Confirm > Execute ==="
echo ""
echo "Press Enter to start kariya..."
read

# Run kariya and redirect stderr to a file so we can see debug output
./kariya 2> "$TEST_DIR/debug.log"

echo ""
echo "=== Debug output: ==="
cat "$TEST_DIR/debug.log"

echo ""
echo "=== Cleanup ==="
rm -rf "$TEST_DIR"
