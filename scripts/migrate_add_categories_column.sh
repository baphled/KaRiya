#!/bin/bash

# Migration script to add 'categories' column to existing KaRiya databases
# and populate it with fallback to tags for backward compatibility
# Usage: ./migrate_add_categories_column.sh <database_path>

if [ -z "$1" ]; then
    echo "Usage: $0 <database_path>"
    echo "Example: $0 ~/.kariya/events.db"
    exit 1
fi

DB_PATH="$1"

if [ ! -f "$DB_PATH" ]; then
    echo "Error: Database file not found: $DB_PATH"
    exit 1
fi

echo "Migrating database: $DB_PATH"
echo "Adding 'categories' column to career_events table..."

# Check if column already exists
COLUMN_EXISTS=$(sqlite3 "$DB_PATH" "PRAGMA table_info(career_events);" | grep -c "categories")

if [ "$COLUMN_EXISTS" -gt 0 ]; then
    echo "✅ Column 'categories' already exists. No migration needed."
    exit 0
fi

# Add the categories column
sqlite3 "$DB_PATH" "ALTER TABLE career_events ADD COLUMN categories TEXT;"

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
    echo "   - Added 'categories' column to career_events table"
    echo ""
    echo "ℹ️  Note: Existing events will have empty categories."
    echo "   The CV generation will fall back to using tags for filtering."
    echo "   You can manually populate categories or let the application"
    echo "   handle the fallback automatically."
    echo ""

    # Verify the migration
    echo "Verifying migration..."
    sqlite3 "$DB_PATH" "PRAGMA table_info(career_events);" | grep "categories"

    if [ $? -eq 0 ]; then
        echo "✅ Verification successful!"
    else
        echo "❌ Verification failed. Please check the database manually."
        exit 1
    fi
else
    echo "❌ Migration failed. Please check the error message above."
    exit 1
fi

