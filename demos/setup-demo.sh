#!/bin/bash
set -e

# Set up fake home
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
  first_name: Yomi
  last_name: Colledge
  email: yomi@boodah.net
  phone: ""
  linkedin: ""
  country: ""
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

# Initialize DB tables manually to ensure clean state
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

INSERT INTO career_events (id, text, date, created_at, updated_at) VALUES
('evt_01', 'Used BubbleTea for UI development', '2023-01-01 10:00:00', '2023-01-01 10:00:00', '2023-01-01 10:00:00');

INSERT INTO skills (id, name, category, created_at, updated_at) VALUES
('skl_01', 'BubbleTea', 'Framework', '2023-01-01 10:00:00', '2023-01-01 10:00:00');
EOF

echo "Demo environment setup complete at $FAKE_HOME"
