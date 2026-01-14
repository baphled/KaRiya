-- +goose Up
-- Create skills table for user-defined skills and technologies
CREATE TABLE IF NOT EXISTS skills (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL,
    level TEXT,
    years_used INTEGER,
    last_used DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Indexes for querying skills by category and name
CREATE INDEX IF NOT EXISTS idx_skills_category ON skills(category);
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);

-- +goose Down
-- Rollback: Remove indexes and table
DROP INDEX IF EXISTS idx_skills_name;
DROP INDEX IF EXISTS idx_skills_category;
DROP TABLE IF EXISTS skills;
