#!/bin/bash

# Migration script to add 'project' column to existing KaRiya databases
# Usage: ./migrate_add_project_column.sh <database_path>

if [ -z "$1" ]; then
    echo "Usage: $0 <database_path>"
    echo "Example: $0 ./kariya.db"
    exit 1
fi

DB_PATH="$1"

if [ ! -f "$DB_PATH" ]; then
    echo "Error: Database file not found: $DB_PATH"
    exit 1
fi

echo "Migrating database: $DB_PATH"
echo "Adding 'project' column to career_events table..."

# Check if column already exists
COLUMN_EXISTS=$(sqlite3 "$DB_PATH" "PRAGMA table_info(career_events);" | grep -c "project")

if [ "$COLUMN_EXISTS" -gt 0 ]; then
    echo "✅ Column 'project' already exists. No migration needed."
    exit 0
fi

# Add the project column
sqlite3 "$DB_PATH" "ALTER TABLE career_events ADD COLUMN project TEXT;"

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
    echo "   - Added 'project' column to career_events table"
    
    # Verify the migration
    echo ""
    echo "Verifying migration..."
    sqlite3 "$DB_PATH" "PRAGMA table_info(career_events);" | grep "project"
    
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
