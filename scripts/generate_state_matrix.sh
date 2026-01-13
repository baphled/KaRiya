#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "Generating state matrix documentation..."
go run ./cmd/generate-state-matrix/main.go

echo "✅ State matrix generated successfully"
echo "   - docs/STATE_MATRIX.md"
echo "   - docs/state_matrix.json"
