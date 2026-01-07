.PHONY: test coverage test-suite individual-test review-commit pre-commit build fmt vet check-compliance install-git-hooks check-ai-attribution audit-ai-commits list-ai-commits ci-local ci-install-tools gosec

# Run all tests in verbose mode
test:
	ginkgo -v --race ./...

# Run a specific test suite
test-suite:
	@if [ -z "$(SUITE)" ]; then \
		echo "Please specify a test suite using SUITE=path/to/suite"; \
		exit 1; \
	fi
	ginkgo -v --race $(SUITE)

# Run a specific test
individual-test:
	@if [ -z "$(TEST)" ]; then \
		echo "Please specify a test using TEST=path/to/test/file/TestName"; \
		exit 1; \
	fi
	ginkgo -v -focus="$(TEST)" ./...

coverage:
	@bash scripts/test-coverage.sh

clean-coverage:
	@rm -rf coverage

# Build the application
build:
	@echo "Building KaRiya..."
	@go build -o kariya ./cmd/cli

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run staticcheck (advanced static analysis)
staticcheck:
	@echo "Running staticcheck..."
	@command -v staticcheck >/dev/null 2>&1 || { echo "Installing staticcheck..."; go install honnef.co/go/tools/cmd/staticcheck@latest; }
	@staticcheck ./...

# Run gosec security scanner
gosec:
	@echo "Running gosec security scanner..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	@gosec -no-fail -fmt text ./...

# Pre-commit checks (quick)
pre-commit:
	@echo "Running pre-commit checks..."
	@go fmt ./...
	@go vet ./...
	@command -v staticcheck >/dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest
	@staticcheck ./...
	@go build ./...
	@go test ./...
	@echo "✅ Pre-commit checks passed"

# Review staged commit (comprehensive)
review-commit:
	@bash scripts/review-commit.sh

# Check full project compliance (all rules)
check-compliance: staticcheck
	@bash scripts/check-compliance.sh

# Install all CI tools locally
ci-install-tools:
	@echo "Installing all CI tools..."
	@command -v ginkgo >/dev/null 2>&1 || { echo "Installing ginkgo..."; go install github.com/onsi/ginkgo/v2/ginkgo@latest; }
	@command -v staticcheck >/dev/null 2>&1 || { echo "Installing staticcheck..."; go install honnef.co/go/tools/cmd/staticcheck@latest; }
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	@[ -d node_modules ] || { echo "Installing npm dependencies..."; npm ci; }
	@echo "✅ All CI tools installed"

# Run ALL CI checks locally (mirrors GitHub Actions)
ci-local:
	@bash scripts/ci-local.sh

# Install git hooks for AI attribution
install-git-hooks:
	@bash scripts/install-git-hooks.sh

# Check AI attribution in latest commit
check-ai-attribution:
	@echo "Checking latest commit for AI attribution..."
	@git log -1 --pretty=%B | grep "AI-Generated-By:" || \
		echo "⚠️  No AI attribution found in latest commit"

# Audit all AI-generated commits
audit-ai-commits:
	@if [ -f scripts/audit-ai-commits.sh ]; then \
		bash scripts/audit-ai-commits.sh; \
	else \
		echo "Total AI commits: $$(git log --all --grep='AI-Generated-By:' --oneline | wc -l)"; \
		echo ""; \
		echo "By Assistant:"; \
		git log --all --grep="AI-Generated-By:" --pretty=%B | \
			grep "AI-Generated-By:" | sort | uniq -c; \
	fi

# List all AI-generated commits
list-ai-commits:
	@echo "AI-Generated Commits:"
	@git log --all --grep="AI-Generated-By:" --oneline

# Show token efficiency reminder
token-check:
	@echo "================================================"
	@echo "💬 TOKEN EFFICIENCY REMINDERS"
	@echo "================================================"
	@echo ""
	@echo "Token Thresholds:"
	@echo "  ✅ < 20k: Healthy"
	@echo "  ⚠️  20-50k: Be more concise"
	@echo "  🔶 50-100k: Consider fresh start"
	@echo "  🔴 >100k: Start fresh NOW"
	@echo ""
	@echo "Best Practices:"
	@echo "  • Use tools (view, grep, ls) over text"
	@echo "  • Be concise and specific"
	@echo "  • Batch multiple operations"
	@echo "  • Reference context, don't repeat"
	@echo "  • Focus on deltas, not full state"
	@echo ""

# Show task workflow
task-workflow:
	@cat docs/rules/TASK_QUICK_REF.md

# Show help for all available targets
help:
	@echo "================================================"
	@echo "📋 KARIYA PROJECT - AVAILABLE COMMANDS"
	@echo "================================================"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test              - Run all tests"
	@echo "  make test-suite        - Run specific suite (SUITE=path)"
	@echo "  make individual-test   - Run specific test (TEST=name)"
	@echo "  make coverage          - Generate coverage report"
	@echo ""
	@echo "🔍 Quality Checks:"
	@echo "  make check-compliance  - Full rules compliance check"
	@echo "  make review-commit     - Review staged commit"
	@echo "  make pre-commit        - Quick pre-commit checks"
	@echo "  make ci-local          - Run ALL CI checks locally (mirrors GitHub Actions)"
	@echo "  make ci-install-tools  - Install all required CI tools"
	@echo "  make fmt               - Format code"
	@echo "  make vet               - Run static analysis"
	@echo "  make staticcheck       - Run staticcheck (advanced analysis)"
	@echo "  make gosec             - Run security scanner"
	@echo ""
	@echo "🤖 AI Attribution:"
	@echo "  make install-git-hooks    - Install AI attribution hooks"
	@echo "  make check-ai-attribution - Check latest commit"
	@echo "  make audit-ai-commits     - Audit all AI commits"
	@echo "  make list-ai-commits      - List AI-generated commits"
	@echo ""
	@echo "💬 Workflow Guides:"
	@echo "  make task-workflow     - Show task execution workflow"
	@echo "  make token-check       - Show token efficiency tips"
	@echo "  make help              - Show this help message"
	@echo ""
	@echo "🏗️  Build:"
	@echo "  make build             - Build the application"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  docs/rules/master-task-prompt.md     - Full task guide"
	@echo "  docs/rules/TASK_QUICK_REF.md         - Quick reference"
	@echo "  docs/rules/AI_COMMIT_ATTRIBUTION.md  - AI attribution rules"
	@echo "  MASTER_TASK_PROMPT_SETUP.md          - Setup guide"
	@echo ""

