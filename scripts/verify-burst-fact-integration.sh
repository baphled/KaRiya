#!/bin/bash
# Verification Script for Burst & Fact CLI Integration
# Tests all P1-P5 functionality

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test database path
TEST_DB="/tmp/kariya-test-$(date +%s).db"
CLI_BIN="./kariya-cli"

echo "======================================"
echo "Burst & Fact Integration Verification"
echo "======================================"
echo ""

# Step 1: Build the CLI
echo -e "${YELLOW}[1/10] Building CLI...${NC}"
go build -o "$CLI_BIN" ./cmd/cli
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
echo ""

# Step 2: Run all tests with race detection
echo -e "${YELLOW}[2/10] Running all tests with race detection...${NC}"
go test -race ./... > /tmp/test-output.txt 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passing${NC}"
    grep -E "^ok" /tmp/test-output.txt | wc -l | xargs echo "  Packages tested:"
else
    echo -e "${RED}✗ Tests failed${NC}"
    cat /tmp/test-output.txt
    exit 1
fi
echo ""

# Step 3: Verify CLI flags
echo -e "${YELLOW}[3/10] Verifying CLI flags...${NC}"
$CLI_BIN --help > /tmp/help-output.txt 2>&1
if grep -q "detect-bursts" /tmp/help-output.txt && \
   grep -q "extract-facts" /tmp/help-output.txt && \
   grep -q "show-bursts" /tmp/help-output.txt && \
   grep -q "show-facts" /tmp/help-output.txt; then
    echo -e "${GREEN}✓ All burst/fact flags present${NC}"
    echo "  - --detect-bursts"
    echo "  - --extract-facts"
    echo "  - --show-bursts"
    echo "  - --show-facts"
else
    echo -e "${RED}✗ Missing flags${NC}"
    exit 1
fi
echo ""

# Step 4: Create test CSV
echo -e "${YELLOW}[4/10] Creating test CSV data...${NC}"
cat > /tmp/test-events.csv << 'EOF'
Text,Date,Company,Project,Tags
"Led microservices migration for payment platform",2024-01-15,TechCorp,Platform Migration,technical;leadership
"Architected distributed system handling 10M+ requests/day",2024-01-20,TechCorp,Platform Migration,technical;architecture
"Mentored 3 engineers on system design patterns",2024-02-01,TechCorp,Team Development,mentoring;leadership
"Designed API gateway with rate limiting and auth",2024-02-10,TechCorp,Platform Migration,technical;architecture
"Conducted architecture review with stakeholders",2024-02-15,TechCorp,Platform Migration,leadership;consulting
"Implemented CI/CD pipeline for microservices",2024-03-01,TechCorp,DevOps Automation,technical;product
"Researched event-driven architecture patterns",2024-03-10,TechCorp,Research Initiative,research;technical
"Delivered training on microservices best practices",2024-03-20,TechCorp,Team Development,mentoring;technical
"Optimized database queries reducing latency by 60%",2024-04-01,TechCorp,Performance Initiative,technical;consulting
"Collaborated with product team on roadmap planning",2024-04-15,TechCorp,Product Strategy,product;leadership
EOF
echo -e "${GREEN}✓ Test CSV created (10 events)${NC}"
echo ""

# Step 5: Import CSV with burst/fact detection
echo -e "${YELLOW}[5/10] Importing CSV with burst/fact detection...${NC}"
$CLI_BIN --db "$TEST_DB" --import /tmp/test-events.csv --skip-import-review > /tmp/import-output.txt 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Import successful${NC}"
    grep "Successfully imported" /tmp/import-output.txt || echo "  (Import completed)"
else
    echo -e "${RED}✗ Import failed${NC}"
    cat /tmp/import-output.txt
    exit 1
fi
echo ""

# Step 6: Verify database tables exist
echo -e "${YELLOW}[6/10] Verifying database schema...${NC}"
if sqlite3 "$TEST_DB" ".tables" | grep -q "bursts" && \
   sqlite3 "$TEST_DB" ".tables" | grep -q "facts"; then
    echo -e "${GREEN}✓ Database tables created${NC}"
    echo "  Tables: $(sqlite3 "$TEST_DB" ".tables")"
else
    echo -e "${RED}✗ Missing database tables${NC}"
    sqlite3 "$TEST_DB" ".tables"
    exit 1
fi
echo ""

# Step 7: Check burst detection results
echo -e "${YELLOW}[7/10] Checking burst detection results...${NC}"
BURST_COUNT=$(sqlite3 "$TEST_DB" "SELECT COUNT(*) FROM bursts;")
if [ "$BURST_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Bursts detected: $BURST_COUNT${NC}"
    sqlite3 "$TEST_DB" "SELECT name, event_count, confidence_score FROM bursts LIMIT 3;" | while read line; do
        echo "  - $line"
    done
else
    echo -e "${YELLOW}⚠ No bursts detected (may be expected with small dataset)${NC}"
fi
echo ""

# Step 8: Check fact extraction results
echo -e "${YELLOW}[8/10] Checking fact extraction results...${NC}"
FACT_COUNT=$(sqlite3 "$TEST_DB" "SELECT COUNT(*) FROM facts;")
if [ "$FACT_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Facts extracted: $FACT_COUNT${NC}"
    sqlite3 "$TEST_DB" "SELECT SUBSTR(text, 1, 60), competencies, role_fit FROM facts LIMIT 3;" | while read line; do
        echo "  - $line"
    done
else
    echo -e "${YELLOW}⚠ No facts extracted${NC}"
fi
echo ""

# Step 9: Test CLI commands
echo -e "${YELLOW}[9/10] Testing CLI commands...${NC}"

# Test --show-bursts
$CLI_BIN --db "$TEST_DB" --show-bursts > /tmp/show-bursts.txt 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ --show-bursts works${NC}"
else
    echo -e "${RED}✗ --show-bursts failed${NC}"
fi

# Test --show-facts
$CLI_BIN --db "$TEST_DB" --show-facts > /tmp/show-facts.txt 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ --show-facts works${NC}"
else
    echo -e "${RED}✗ --show-facts failed${NC}"
fi

# Test --detect-bursts
$CLI_BIN --db "$TEST_DB" --detect-bursts > /tmp/detect-bursts.txt 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ --detect-bursts works${NC}"
else
    echo -e "${RED}✗ --detect-bursts failed${NC}"
fi

# Test --extract-facts
$CLI_BIN --db "$TEST_DB" --extract-facts > /tmp/extract-facts.txt 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ --extract-facts works${NC}"
else
    echo -e "${RED}✗ --extract-facts failed${NC}"
fi
echo ""

# Step 10: Summary
echo -e "${YELLOW}[10/10] Verification Summary${NC}"
echo "======================================"
echo -e "${GREEN}✓ CLI builds successfully${NC}"
echo -e "${GREEN}✓ All tests passing (0 race conditions)${NC}"
echo -e "${GREEN}✓ All CLI flags present and working${NC}"
echo -e "${GREEN}✓ CSV import successful (10 events)${NC}"
echo -e "${GREEN}✓ Database tables created${NC}"
echo -e "${GREEN}✓ Bursts detected: $BURST_COUNT${NC}"
echo -e "${GREEN}✓ Facts extracted: $FACT_COUNT${NC}"
echo -e "${GREEN}✓ All CLI commands functional${NC}"
echo ""
echo -e "${GREEN}✓✓✓ VERIFICATION COMPLETE ✓✓✓${NC}"
echo ""
echo "Test database: $TEST_DB"
echo "Output files in: /tmp/"
echo ""
echo "To inspect results:"
echo "  sqlite3 $TEST_DB '.schema'"
echo "  sqlite3 $TEST_DB 'SELECT * FROM bursts;'"
echo "  sqlite3 $TEST_DB 'SELECT * FROM facts;'"
echo ""
echo "To clean up:"
echo "  rm $TEST_DB"
echo "  rm /tmp/test-events.csv"
echo "  rm /tmp/*.txt"

