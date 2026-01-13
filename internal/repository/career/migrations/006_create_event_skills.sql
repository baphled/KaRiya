-- +goose Up
-- Create event_skills junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS event_skills (
    event_id TEXT NOT NULL,
    skill_id TEXT NOT NULL,
    PRIMARY KEY (event_id, skill_id),
    FOREIGN KEY (event_id) REFERENCES career_events(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

-- Indexes for efficient querying in both directions
CREATE INDEX IF NOT EXISTS idx_event_skills_event ON event_skills(event_id);
CREATE INDEX IF NOT EXISTS idx_event_skills_skill ON event_skills(skill_id);

-- +goose Down
-- Rollback: Remove indexes and table
DROP INDEX IF EXISTS idx_event_skills_skill;
DROP INDEX IF EXISTS idx_event_skills_event;
DROP TABLE IF EXISTS event_skills;
