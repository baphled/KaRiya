# Scripter Skill

You are a scripting expert proficient in Bash, Python, and other scripting languages for automation and tooling.

## Overview

Scripts are the glue of software development. Write scripts that are readable, maintainable, and robust.

---

## Bash Scripting

### Script Header Template

```bash
#!/usr/bin/env bash
#
# Script: script-name.sh
# Description: Brief description of what this script does
# Usage: ./script-name.sh [options] <arguments>
#

set -euo pipefail  # Exit on error, undefined vars, pipe failures
IFS=$'\n\t'        # Safer word splitting

# Constants
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_NAME="$(basename "$0")"

# Colors (optional)
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m'  # No Color
```

### Error Handling

```bash
# Error handler
error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
    exit "${2:-1}"
}

# Warning
warn() {
    echo -e "${YELLOW}WARNING: $1${NC}" >&2
}

# Success
success() {
    echo -e "${GREEN}$1${NC}"
}

# Usage with trap
cleanup() {
    # Cleanup temporary files
    rm -f "$TEMP_FILE" 2>/dev/null || true
}
trap cleanup EXIT
```

### Argument Parsing

```bash
# Simple positional args
if [[ $# -lt 1 ]]; then
    error "Usage: $SCRIPT_NAME <argument>"
fi

ARG1="$1"
ARG2="${2:-default}"  # Default value

# Getopts for options
usage() {
    cat << EOF
Usage: $SCRIPT_NAME [OPTIONS] <input>

Options:
    -h, --help      Show this help
    -v, --verbose   Verbose output
    -o, --output    Output file
    -f, --force     Force overwrite

Examples:
    $SCRIPT_NAME -v input.txt
    $SCRIPT_NAME -o output.txt input.txt
EOF
}

VERBOSE=false
FORCE=false
OUTPUT=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -o|--output)
            OUTPUT="$2"
            shift 2
            ;;
        -*)
            error "Unknown option: $1"
            ;;
        *)
            INPUT="$1"
            shift
            ;;
    esac
done
```

### Common Patterns

```bash
# Check command exists
command_exists() {
    command -v "$1" &> /dev/null
}

if ! command_exists go; then
    error "Go is not installed"
fi

# Check file exists
if [[ ! -f "$FILE" ]]; then
    error "File not found: $FILE"
fi

# Check directory
if [[ ! -d "$DIR" ]]; then
    mkdir -p "$DIR"
fi

# Read file line by line
while IFS= read -r line; do
    echo "Processing: $line"
done < "$FILE"

# Loop through files
for file in *.go; do
    [[ -e "$file" ]] || continue  # Handle no matches
    echo "Processing: $file"
done

# Conditional execution
[[ "$VERBOSE" == true ]] && echo "Verbose mode enabled"

# Array operations
declare -a FILES=()
FILES+=("file1.txt")
FILES+=("file2.txt")

for f in "${FILES[@]}"; do
    echo "$f"
done
```

### Logging

```bash
# Log levels
LOG_LEVEL="${LOG_LEVEL:-INFO}"

log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp
    timestamp="$(date '+%Y-%m-%d %H:%M:%S')"
    
    case "$level" in
        DEBUG)
            [[ "$LOG_LEVEL" == "DEBUG" ]] && echo "[$timestamp] DEBUG: $message"
            ;;
        INFO)
            [[ "$LOG_LEVEL" =~ ^(DEBUG|INFO)$ ]] && echo "[$timestamp] INFO: $message"
            ;;
        WARN)
            echo "[$timestamp] WARN: $message" >&2
            ;;
        ERROR)
            echo "[$timestamp] ERROR: $message" >&2
            ;;
    esac
}

log INFO "Starting script"
log DEBUG "Debug information"
log ERROR "Something went wrong"
```

---

## Python Scripting

### Script Template

```python
#!/usr/bin/env python3
"""
Script: script_name.py
Description: Brief description
Usage: python script_name.py [options] <arguments>
"""

import argparse
import logging
import sys
from pathlib import Path

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


def parse_args():
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(
        description='Script description',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument('input', help='Input file')
    parser.add_argument('-o', '--output', help='Output file')
    parser.add_argument('-v', '--verbose', action='store_true', help='Verbose output')
    parser.add_argument('-f', '--force', action='store_true', help='Force overwrite')
    
    return parser.parse_args()


def main():
    """Main entry point."""
    args = parse_args()
    
    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)
    
    logger.info(f"Processing: {args.input}")
    
    # Main logic here
    
    return 0


if __name__ == '__main__':
    sys.exit(main())
```

### Common Patterns

```python
# File operations
from pathlib import Path

path = Path('data/file.txt')

# Check exists
if not path.exists():
    raise FileNotFoundError(f"File not found: {path}")

# Read file
content = path.read_text()

# Write file
path.write_text("content")

# Iterate directory
for file in Path('.').glob('**/*.go'):
    print(file)

# JSON handling
import json

data = json.loads(path.read_text())
path.write_text(json.dumps(data, indent=2))

# YAML handling
import yaml

with open('config.yaml') as f:
    config = yaml.safe_load(f)

# Subprocess
import subprocess

result = subprocess.run(
    ['go', 'test', './...'],
    capture_output=True,
    text=True,
    check=True
)
print(result.stdout)

# HTTP requests
import requests

response = requests.get('https://api.example.com/data')
response.raise_for_status()
data = response.json()
```

---

## Go Scripts

### Quick Scripts with go run

```go
//go:build ignore
// +build ignore

// Script: check-files.go
// Usage: go run check-files.go

package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "Usage: go run check-files.go <directory>")
        os.Exit(1)
    }
    
    dir := os.Args[1]
    
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".go" {
            fmt.Println(path)
        }
        return nil
    })
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

---

## Script Best Practices

### 1. Always Use Strict Mode (Bash)

```bash
set -euo pipefail
```

### 2. Quote Variables

```bash
# WRONG
if [ $FILE = "test" ]; then

# RIGHT
if [[ "$FILE" = "test" ]]; then
```

### 3. Use Functions

```bash
# WRONG - Inline everything
echo "Step 1"
# 50 lines of code
echo "Step 2"
# 50 more lines

# RIGHT - Functions
step_one() {
    echo "Step 1"
    # Implementation
}

step_two() {
    echo "Step 2"
    # Implementation
}

main() {
    step_one
    step_two
}

main "$@"
```

### 4. Handle Errors Gracefully

```bash
# WRONG
rm important-file.txt

# RIGHT
rm important-file.txt || error "Failed to remove file"
```

### 5. Make Scripts Idempotent

```bash
# WRONG - Fails if already exists
mkdir data

# RIGHT - Idempotent
mkdir -p data

# WRONG - Appends duplicate
echo "export PATH=..." >> ~/.bashrc

# RIGHT - Check first
if ! grep -q "export PATH=..." ~/.bashrc; then
    echo "export PATH=..." >> ~/.bashrc
fi
```

---

## Common Script Tasks

### Check Go Code Quality

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "Running Go quality checks..."

echo "==> go fmt"
UNFORMATTED=$(gofmt -l .)
if [[ -n "$UNFORMATTED" ]]; then
    echo "Unformatted files:"
    echo "$UNFORMATTED"
    exit 1
fi

echo "==> go vet"
go vet ./...

echo "==> staticcheck"
staticcheck ./...

echo "==> go test"
go test -race ./...

echo "All checks passed!"
```

### Database Operations

```bash
#!/usr/bin/env bash
set -euo pipefail

DB_PATH="${DB_PATH:-./data/app.db}"
BACKUP_DIR="${BACKUP_DIR:-./backups}"

backup() {
    mkdir -p "$BACKUP_DIR"
    local backup_file="$BACKUP_DIR/app-$(date +%Y%m%d-%H%M%S).db"
    cp "$DB_PATH" "$backup_file"
    echo "Backup created: $backup_file"
}

restore() {
    local backup_file="$1"
    if [[ ! -f "$backup_file" ]]; then
        error "Backup file not found: $backup_file"
    fi
    cp "$backup_file" "$DB_PATH"
    echo "Restored from: $backup_file"
}

case "${1:-}" in
    backup)
        backup
        ;;
    restore)
        restore "${2:-}"
        ;;
    *)
        echo "Usage: $0 {backup|restore <file>}"
        exit 1
        ;;
esac
```

---

## Related Skills

- `automation` - Automating workflows
- `devops` - CI/CD scripting
- `go-expert` - Go tooling
