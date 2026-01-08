-- +goose Up
-- Create facts table for extracted career insights
CREATE TABLE IF NOT EXISTS facts (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    competencies TEXT NOT NULL,
    role_fit TEXT NOT NULL,
    audience_relevance TEXT NOT NULL,
    strength_signal TEXT,
    source_event_id TEXT,
    source_burst_id TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Indexes for joining facts with their source events/bursts
CREATE INDEX IF NOT EXISTS idx_facts_source_event_id ON facts(source_event_id);
CREATE INDEX IF NOT EXISTS idx_facts_source_burst_id ON facts(source_burst_id);

-- +goose Down
-- Rollback: Remove indexes and table
DROP INDEX IF EXISTS idx_facts_source_burst_id;
DROP INDEX IF EXISTS idx_facts_source_event_id;
DROP TABLE IF EXISTS facts;
