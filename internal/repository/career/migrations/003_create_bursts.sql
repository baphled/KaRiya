-- +goose Up
-- Create bursts table for grouping related career events
CREATE TABLE IF NOT EXISTS bursts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    event_ids TEXT NOT NULL,
    confirmed INTEGER NOT NULL DEFAULT 0,
    confirmed_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Index for filtering confirmed/unconfirmed bursts
CREATE INDEX IF NOT EXISTS idx_bursts_confirmed ON bursts(confirmed);

-- +goose Down
-- Rollback: Remove index and table
DROP INDEX IF EXISTS idx_bursts_confirmed;
DROP TABLE IF EXISTS bursts;
