#!/bin/bash
set -e

# Set up fake home for skill acceptance demo
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

echo "Skill demo environment setup complete at $FAKE_HOME"
