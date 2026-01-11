#!/usr/bin/env python3
"""
E2E Test Splitting Script

Splits MIXED workflow test files into:
1. E2E tests (SQLite) - stays in internal/testutil/e2e/
2. Navigation tests (Memory) - moves to internal/cli/intents/

Usage: python3 scripts/split-e2e-tests.py <workflow_test_file>
Example: python3 scripts/split-e2e-tests.py browse_workflow_test.go
"""

import sys
import os
import re
from pathlib import Path


def extract_describe_blocks(content, setup_pattern):
    """Extract Describe blocks that contain the specified setup pattern."""
    blocks = []
    lines = content.split("\n")

    i = 0
    while i < len(lines):
        line = lines[i]

        # Find Describe blocks
        if re.match(r"^\s*Describe\(", line):
            # Extract the block
            block_lines = [line]
            brace_count = line.count("{") - line.count("}")
            i += 1

            while i < len(lines) and brace_count > 0:
                block_lines.append(lines[i])
                brace_count += lines[i].count("{") - lines[i].count("}")
                i += 1

            block_text = "\n".join(block_lines)

            # Check if this block contains the setup pattern
            if setup_pattern in block_text:
                blocks.append(block_text)
            continue

        i += 1

    return blocks


def to_pascal_case(snake_str):
    """Convert snake_case to PascalCase."""
    components = snake_str.split("_")
    return "".join(x.title() for x in components)


def main():
    if len(sys.argv) < 2:
        print("❌ ERROR: Workflow test file required\n")
        print("Usage: python3 scripts/split-e2e-tests.py <workflow_test_file>")
        print("\nExample: python3 scripts/split-e2e-tests.py browse_workflow_test.go")
        print("\nRemaining files to split:")
        e2e_dir = Path("internal/testutil/e2e")
        for f in sorted(e2e_dir.glob("*_workflow_test.go")):
            print(f"  {f.name}")
        sys.exit(1)

    workflow_file = sys.argv[1]

    # Paths
    e2e_dir = Path("internal/testutil/e2e")
    intents_dir = Path("internal/cli/intents")
    source_file = e2e_dir / workflow_file

    # Validate source file exists
    if not source_file.exists():
        print(f"❌ ERROR: File not found: {source_file}")
        sys.exit(1)

    # Extract base name (e.g., "browse" from "browse_workflow_test.go")
    base_name = workflow_file.replace("_workflow_test.go", "")

    # Output files
    e2e_output = e2e_dir / f"{base_name}_e2e_test.go"
    nav_output = intents_dir / f"{base_name}_navigation_test.go"

    print(f"\n🔍 Splitting: {workflow_file}")
    print(f"  E2E tests   → {e2e_output.name}")
    print(f"  Navigation  → {nav_output.name}\n")

    # Read source file
    content = source_file.read_text()

    # Extract blocks
    print("📝 Extracting E2E tests...")
    e2e_blocks = extract_describe_blocks(content, "e2e.Setup(GinkgoT())")

    print("📝 Extracting Navigation tests...")
    nav_blocks = extract_describe_blocks(content, "e2e.SetupWithMemory(GinkgoT())")

    # Get intent name for navigation tests
    intent_name = to_pascal_case(
        base_name.replace("_timeline", "Timeline")
        .replace("_management", "Management")
        .replace("_artifact", "Artifact")
        .replace("_system", "System")
        .replace("_wizard", "Wizard")
        .replace("_editor", "Editor")
        .replace("_cv", "CV")
    )

    # Write E2E file
    e2e_content = (
        """package e2e_test

import (
\t"github.com/baphled/kariya/internal/testutil/e2e"
\ttea "github.com/charmbracelet/bubbletea"
\t. "github.com/onsi/ginkgo/v2"
\t. "github.com/onsi/gomega"
)

var _ = Describe("E2E """
        + intent_name
        + """ Workflow", func() {
\tvar env *e2e.TestEnv

"""
    )

    for block in e2e_blocks:
        # Add proper indentation
        indented_block = "\n".join(
            "\t" + line if line.strip() else line for line in block.split("\n")
        )
        e2e_content += indented_block + "\n\n"

    e2e_content += "})\n"

    # Write Navigation file
    nav_content = (
        '''package intents_test

import (
\t"github.com/baphled/kariya/internal/testutil/e2e"
\t. "github.com/onsi/ginkgo/v2"
\t. "github.com/onsi/gomega"
)

var _ = Describe("'''
        + intent_name
        + """ Navigation", func() {
\tvar env *e2e.TestEnv

"""
    )

    for block in nav_blocks:
        # Add proper indentation
        indented_block = "\n".join(
            "\t" + line if line.strip() else line for line in block.split("\n")
        )
        nav_content += indented_block + "\n\n"

    nav_content += "})\n"

    # Write files
    e2e_output.write_text(e2e_content)
    nav_output.write_text(nav_content)

    print(f"✅ E2E tests extracted: {e2e_output} ({len(e2e_blocks)} blocks)")
    print(f"✅ Navigation tests extracted: {nav_output} ({len(nav_blocks)} blocks)")

    # Summary
    print("\n✅ Split complete!\n")
    print("Next steps:")
    print("  1. Review the generated files")
    print(f"  2. Run tests: ginkgo {e2e_output}")
    print(f"               ginkgo {nav_output}")
    print(f"  3. Delete original: rm {source_file}")
    print(
        f'  4. Commit: git add -A && make ai-commit MSG="refactor(tests): split {base_name}_workflow into e2e and navigation tests"'
    )
    print()


if __name__ == "__main__":
    main()
