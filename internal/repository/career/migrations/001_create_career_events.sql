-- +goose Up
-- Create career events table with indexes for improved query performance
CREATE TABLE IF NOT EXISTS career_events (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    date DATETIME NOT NULL,
    tags TEXT,
    company TEXT,
    project TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_career_events_date ON career_events(date);
CREATE INDEX IF NOT EXISTS idx_career_events_company ON career_events(company);

-- +goose Down
-- Rollback: Remove indexes and table
DROP INDEX IF EXISTS idx_career_events_company;
DROP INDEX IF EXISTS idx_career_events_date;
DROP TABLE IF EXISTS career_events;
