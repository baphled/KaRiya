#!/bin/bash
set -e

# Set up fake home for burst management acceptance demo
# Needs ≥5 existing events with similar attributes so SuggestBursts can run
FAKE_HOME="$(pwd)/demos/temp_demo_env"
rm -rf "$FAKE_HOME"
mkdir -p "$FAKE_HOME/.kariya"

# Create config
cat <<EOF > "$FAKE_HOME/.kariya/config.yaml"
system:
  data_dir: "$FAKE_HOME/.kariya"
  log_level: info
  auto_backup: true
  backup_count: 5
profile:
  name: Yomi Colledge
  email: yomi@boodah.net
  default_role: ""
  default_audience: ""
  title: ""
  location: Remote
  github: baphled
  portfolio: http://boodah.net
  core_strengths: []
  languages: []
  frontend: []
  systems: []
  what_i_bring: []
cv:
  default_format: markdown
  max_bullets: 50
  audience_bullets:
    recruiter: 4
    hiring_manager: 6
    peer: 8
    default: 5
export:
  default_destination: file
  auto_open: false
display:
  theme: dark
  animations: true
scoring:
  weights:
    role_score: 0.25
    audience_score: 0.2
    metric_score: 0.2
    impact_score: 0.2
    confidence: 0.15
  thresholds:
    fact_default_confidence: 0.85
    event_default_confidence: 0.8
    high_confidence: 0.8
    high_impact_confidence: 0.85
  role_settings:
    em:
      min_confidence: 0.75
      max_bullets_per_company: 4
    principal:
      min_confidence: 0.8
      max_bullets_per_company: 4
    senior_ic:
      min_confidence: 0.75
      max_bullets_per_company: 5
    staff:
      min_confidence: 0.75
      max_bullets_per_company: 5
EOF

# Initialise DB schema manually (matches all 6 goose migrations)
DB_PATH="$FAKE_HOME/.kariya/events.db"
sqlite3 "$DB_PATH" <<'SQLEOF'
-- Migration 001: career_events table + indexes
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
CREATE INDEX IF NOT EXISTS idx_career_events_date ON career_events(date);
CREATE INDEX IF NOT EXISTS idx_career_events_company ON career_events(company);

-- Migration 002: categories column
ALTER TABLE career_events ADD COLUMN categories TEXT;

-- Migration 003: bursts table
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
CREATE INDEX IF NOT EXISTS idx_bursts_confirmed ON bursts(confirmed);

-- Migration 004: facts table
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
CREATE INDEX IF NOT EXISTS idx_facts_source_event_id ON facts(source_event_id);
CREATE INDEX IF NOT EXISTS idx_facts_source_burst_id ON facts(source_burst_id);

-- Migration 005: skills table
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
CREATE INDEX IF NOT EXISTS idx_skills_category ON skills(category);
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);

-- Migration 006: event_skills junction table
CREATE TABLE IF NOT EXISTS event_skills (
    event_id TEXT NOT NULL,
    skill_id TEXT NOT NULL,
    PRIMARY KEY (event_id, skill_id),
    FOREIGN KEY (event_id) REFERENCES career_events(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_event_skills_event ON event_skills(event_id);
CREATE INDEX IF NOT EXISTS idx_event_skills_skill ON event_skills(skill_id);

-- goose version table so app startup skips re-running migrations
CREATE TABLE IF NOT EXISTS goose_db_version (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version_id INTEGER NOT NULL,
    is_applied INTEGER NOT NULL,
    tstamp TIMESTAMP DEFAULT (datetime('now'))
);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (0, 1);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (1, 1);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (2, 1);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (3, 1);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (4, 1);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (5, 1);
INSERT INTO goose_db_version (version_id, is_applied) VALUES (6, 1);

-- Seed 5 events within the last 6 months so SuggestBursts temporal window includes them.
-- All events share the tokens "microservices" and "reducing", giving text similarity of ~0.22.
-- Combined with company match (TechCo) and empty-keyword/project scores, pairs reach
-- combined similarity of ~0.61 — above the 0.60 MinConfidence threshold —
-- so DetectBursts forms a cluster and returns a burst suggestion.
-- This mirrors the WORKING pattern from setup-burst-demo.sh.
INSERT INTO career_events (id, text, date, company, created_at, updated_at) VALUES
  ('evt_01', 'Built microservices in Go, reducing deployment time by 30%', date('now', '-150 days'), 'TechCo', datetime('now', '-150 days'), datetime('now', '-150 days')),
  ('evt_02', 'Deployed microservices to production, reducing query latency by 40%', date('now', '-120 days'), 'TechCo', datetime('now', '-120 days'), datetime('now', '-120 days')),
  ('evt_03', 'Optimised microservices architecture, reducing memory usage by 25%', date('now', '-90 days'), 'TechCo', datetime('now', '-90 days'), datetime('now', '-90 days')),
  ('evt_04', 'Refactored microservices communication, reducing API errors by 50%', date('now', '-60 days'), 'TechCo', datetime('now', '-60 days'), datetime('now', '-60 days')),
  ('evt_05', 'Scaled microservices infrastructure, reducing response time by 35%', date('now', '-30 days'), 'TechCo', datetime('now', '-30 days'), datetime('now', '-30 days'));
SQLEOF

echo "Burst management demo environment setup complete at $FAKE_HOME"
echo "Seeded 5 events for burst inference (company: TechCo, shared tokens: microservices, reducing)"
