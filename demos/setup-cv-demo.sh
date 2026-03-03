#!/bin/bash
set -e

# Set up fake home for CV review screen demo
FAKE_HOME="$(pwd)/demos/temp_demo_env"
rm -rf "$FAKE_HOME"
mkdir -p "$FAKE_HOME/.kariya"

# Create config with a profile (required for CV generation)
cat <<EOF > "$FAKE_HOME/.kariya/config.yaml"
system:
  data_dir: "$FAKE_HOME/.kariya"
  log_level: info
  auto_backup: true
  backup_count: 5
profile:
  name: Yomi Colledge
  first_name: Yomi
  last_name: Colledge
  email: yomi@boodah.net
  phone: ""
  linkedin: ""
  country: ""
  default_role: ""
  default_audience: ""
  title: "Staff Software Engineer"
  location: Remote
  github: baphled
  portfolio: http://boodah.net
  core_strengths:
    - Technical Leadership
    - System Architecture
    - Go Development
  languages:
    - Go
    - TypeScript
    - Ruby
  frontend: []
  systems:
    - PostgreSQL
    - Redis
  what_i_bring:
    - Deep expertise in TUI development
    - Strong mentoring skills
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

# Initialize DB with events and skills for CV generation
DB_PATH="$FAKE_HOME/.kariya/events.db"
sqlite3 "$DB_PATH" <<EOF
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
CREATE TABLE IF NOT EXISTS event_skills (
    event_id TEXT NOT NULL,
    skill_id TEXT NOT NULL,
    PRIMARY KEY (event_id, skill_id),
    FOREIGN KEY (event_id) REFERENCES career_events(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS facts (
    id TEXT PRIMARY KEY,
    event_id TEXT NOT NULL,
    text TEXT NOT NULL,
    confidence REAL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES career_events(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS bursts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    theme TEXT,
    start_date DATETIME,
    end_date DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS cv_profiles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    target_role TEXT,
    target_audience TEXT,
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Insert career events with company context
INSERT INTO career_events (id, text, date, company, project, created_at, updated_at) VALUES
('evt_01', 'Built REST API in Go using Gin framework and PostgreSQL', '2024-06-15 10:00:00', 'TechCorp', 'Platform API', '2024-06-15 10:00:00', '2024-06-15 10:00:00'),
('evt_02', 'Led migration from monolith to microservices architecture', '2024-05-20 14:00:00', 'TechCorp', 'Platform API', '2024-05-20 14:00:00', '2024-05-20 14:00:00'),
('evt_03', 'Implemented Bubble Tea TUI for internal tooling', '2024-04-10 09:00:00', 'TechCorp', 'DevTools', '2024-04-10 09:00:00', '2024-04-10 09:00:00'),
('evt_04', 'Reduced API latency by 40% through Redis caching', '2024-03-25 11:00:00', 'TechCorp', 'Platform API', '2024-03-25 11:00:00', '2024-03-25 11:00:00'),
('evt_05', 'Mentored 3 junior developers on Go best practices', '2024-02-15 15:00:00', 'TechCorp', 'Team Lead', '2024-02-15 15:00:00', '2024-02-15 15:00:00'),
('evt_06', 'Designed event-driven architecture using NATS', '2024-01-10 10:00:00', 'TechCorp', 'Platform API', '2024-01-10 10:00:00', '2024-01-10 10:00:00');

-- Insert skills
INSERT INTO skills (id, name, category, level, years_used, created_at, updated_at) VALUES
('skl_01', 'Go', 'Language', 'expert', 5, '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
('skl_02', 'PostgreSQL', 'Database', 'advanced', 6, '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
('skl_03', 'Redis', 'Database', 'intermediate', 3, '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
('skl_04', 'Bubble Tea', 'Framework', 'advanced', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
('skl_05', 'Microservices', 'Architecture', 'expert', 4, '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
('skl_06', 'NATS', 'Infrastructure', 'intermediate', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00');

-- Link events to skills
INSERT INTO event_skills (event_id, skill_id) VALUES
('evt_01', 'skl_01'),
('evt_01', 'skl_02'),
('evt_02', 'skl_05'),
('evt_03', 'skl_01'),
('evt_03', 'skl_04'),
('evt_04', 'skl_03'),
('evt_06', 'skl_06');

-- Insert facts derived from events
INSERT INTO facts (id, event_id, text, confidence, created_at, updated_at) VALUES
('fact_01', 'evt_04', 'Reduced API latency by 40%', 0.95, '2024-03-25 11:00:00', '2024-03-25 11:00:00'),
('fact_02', 'evt_05', 'Mentored 3 junior developers', 0.90, '2024-02-15 15:00:00', '2024-02-15 15:00:00');

-- Insert CV profile
INSERT INTO cv_profiles (id, name, target_role, target_audience, description, created_at, updated_at) VALUES
('profile_01', 'Staff Engineer', 'staff', 'hiring_manager', 'Profile for staff engineer positions', '2024-01-01 10:00:00', '2024-01-01 10:00:00');

EOF

echo "CV demo environment setup complete at $FAKE_HOME"
