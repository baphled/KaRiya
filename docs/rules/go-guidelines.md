# Go Guidelines

## Go

* **Language:** Go (compiled, statically typed language). Commonly used for backend services and CLI tools. Keep Go code simple and idiomatic, focusing on clarity, efficient concurrency, and explicit error handling.

* **Description:**
  Write Go code that is clear and idiomatic. Follow the official [Effective Go](https://go.dev/doc/effective_go) guidelines. Use Go modules (`go.mod` and `go.sum`) for dependency management to ensure reproducible builds.

## Frameworks

### API

#### Fiber

* **Description:** A Go web framework for building fast and scalable applications. Inspired by Express.js, Fiber offers a familiar syntax and extensive middleware support.
* **Example:**

  ```go
  package main

  import "github.com/gofiber/fiber/v2"

  func main() {
      app := fiber.New()

      app.Get("/", func(c *fiber.Ctx) error {
          return c.SendString("Hello, World!")
      })

      app.Listen(":3000")
  }
  ```

#### Beego

* **Description:** A high-performance web framework for Go with an MVC architecture. Beego includes features like a built-in ORM and a task scheduler.
* **Example:**

  ```go
  package main

  import (
      "github.com/beego/beego/v2/server/web"
  )

  func main() {
      web.Router("/", &MainController{})
      web.Run()
  }

  type MainController struct {
      web.Controller
  }

  func (c *MainController) Get() {
      c.Ctx.WriteString("Hello, Beego!")
  }
  ```

### CLI

#### Cobra

* **Description:** A powerful CLI framework for Go, allowing developers to create command-line applications with ease. It provides features like command parsing, flags, and subcommands.
* **Example:**

  ```go
  package main

  import (
      "fmt"
      "github.com/spf13/cobra"
  )

  func main() {
      var rootCmd = &cobra.Command{
          Use:   "app",
          Short: "A simple CLI application",
          Run: func(cmd *cobra.Command, args []string) {
              fmt.Println("Hello from Cobra!")
          },
      }

      rootCmd.Execute()
  }
  ```

## Framework Details

Go projects favor modular design and simplicity. Frameworks like Fiber and Beego provide a robust foundation for web services, enforcing structured design (MVC patterns) and offering libraries for rapid development. They integrate well with Go’s `net/http` package and encourage clear separation of concerns.

## Project Management Practices

1. **Dependencies:** Use Go Modules for dependency management. Define a `go.mod` file at the project root to pin versions.
2. **Task Automation:** Use Makefiles or simple shell scripts for common tasks (building, testing, linting).
3. **Configuration Files:** Use environment variables (with a `.env` file for local development) for configuration. Libraries like `godotenv` can load these for development.
4. **Version Control:** Use Git for version control. Follow a branching strategy (e.g., Git Flow/GitHub Flow) for features and releases.
5. **Containerisation:** Docker or Nix are used for consistent development and CI environments **(Added)**.
6. **Continuous Integration/Deployment (CI/CD):** Use tools like GitHub Actions, GitLab CI, or CircleCI to automate testing, code linting, and deployment.

* **Conventions:** **(Expanded)**

  1. Handle errors explicitly and immediately. Always check the returned `error` from functions; avoid silent failures. Use `if err != nil { ... }` to handle or return errors. Use `fmt.Errorf` (with `%w`) or the `errors` package to wrap errors with contextual information.
  2. Use generics judiciously. Introduced in Go 1.18, generics are powerful but should be used only when they simplify code or eliminate significant duplication. Do not over-complicate code with generics—clear, simple code is preferred over overly abstract solutions.
  3. Design for testability. Use interfaces to decouple components and facilitate mocking in tests. Dependency injection (passing interfaces or function parameters) helps in writing unit tests for logic without heavy setup.
  4. Use `context.Context` for cancellation and timeouts in long-running operations (e.g., when making external calls or handling requests). Pass context through call chains to enable cancellations and deadlines.
  5. Employ concurrency thoughtfully. Prefer communication over shared memory — use channels to orchestrate goroutines and avoid race conditions. When shared state is unavoidable, protect it with mutexes or other sync primitives. Always strive to write race-free concurrent code (use `go run -race` to detect issues during development).
  6. Prefer structured logging. Use a logging library (e.g., `logrus` or Uber’s `zap`) to include fields (key-value pairs) in logs for better searchability and analysis, instead of simple print statements.
  7. Profile and optimize. Use Go’s built-in profiling tools (the `pprof` package and `go test -bench`) to identify performance bottlenecks and memory issues. Optimize only after identifying real bottlenecks, and prefer readability unless performance measurements justify more complex code.

* **Code Structure:**
  Organize code into packages, each in its own directory:

  * Maintain a clear project layout. For example, use a `cmd/` directory for application entry points (each subdirectory in `cmd/` produces a separate binary), and a `pkg/` (or `internal/`) directory for reusable library code.
  * Keep package names short and lowercase (avoid underscores or mixed-case names).
  * Within packages, follow idiomatic naming for files and identifiers. Exported types and functions should be named with capital letters (and documented with comments), while unexported ones use concise lowercase names.
  * Co-locate test files with the code (e.g., `foo_test.go` in the same package as `foo.go`). Use the standard `testing` package and, if needed, complementary libraries (like `testify`) for assertions.

* **Toolchain:**

  1. **Editor:** Neovim (with Go plugins and LSP support via `gopls`).
  2. **Build System:** Go’s built-in toolchain (`go build`, `go run`, `go test`) for compiling and testing.
  3. **Version Control:** Git (optionally with tools like `pre-commit` hooks to run linters/tests on commits).
  4. **Package Management:** Go Modules (`go mod init`, `go mod tidy`) to manage dependencies and versions.
  5. **Task Automation:** Makefile (common targets for build, test, fmt, vet) or simple shell scripts.
  6. **Documentation:** GoDoc (comments in code) for documentation generation (e.g., hosted on pkg.go.dev).

* **Dependencies:**

  * **Production Dependencies:**

    * `github.com/gofiber/fiber/v2` – Web framework (Fiber).
    * `github.com/beego/beego/v2` – Web framework (Beego MVC).
    * `github.com/sirupsen/logrus` – Structured logging library.
  * **Development Dependencies:**

    * `golangci-lint` – Meta-linter aggregator for code quality.
    * `gofmt` – Code formatter (also via `go fmt` command).
    * `go vet` – Static analyzer for common issues.

## LSP

Neovim (and other editors) integrates with **gopls**, the official Go language server:

* **Setup:** `gopls` provides autocompletion, formatting, refactoring, and diagnostics for Go code. Ensure `gopls` is installed and configured in the editor (e.g., via the LSP client in Neovim).
* **Features:**

  * Intelligent auto-completion of symbols and imports.
  * On-the-fly diagnostics and error highlighting as you type.
  * Jump to definition, find references, and symbol renaming.
  * Code actions for quick fixes and organizing imports.
* **Plugins:** Complementary plugins (for example, an autocompletion plugin like **nvim-cmp** with LSP source) can enhance the editing experience by providing snippet insertion, signature help, and documentation on hover.

## Testing Tools

Go provides a built-in `testing` package for unit tests and benchmarks. For more expressive behavior-driven testing, the Ginkgo framework (often used with the Gomega assertion library) is adopted.

### Ginkgo

* **Description:** A BDD-style testing framework for Go. Ginkgo allows expressive, nested test descriptions, and works in conjunction with Gomega for rich assertions.
* **Example:**

  ```go
  package main

  import (
      "testing"
      . "github.com/onsi/ginkgo/v2"
      . "github.com/onsi/gomega"
  )

  func TestMain(m *testing.M) {
      RunSpecs(m, "Main Suite")
  }

  var _ = Describe("Main", func() {
      It("adds numbers correctly", func() {
          Expect(1 + 1).To(Equal(2))
      })
  })
  ```

### Test Naming Conventions

KaRiya uses specific naming conventions to distinguish between different types of tests:

* **Unit Tests:** `foo_test.go` - Co-located with source files, test individual functions/methods
* **Integration Tests:** `integration_test.go` - Test multiple components working together
* **E2E Tests:** `*_e2e_test.go` - Test end-to-end workflows with SQLite persistence
  * Location: `internal/testutil/e2e/`
  * Use `e2e.Setup(GinkgoT())` for full database setup
  * Verify data persistence with `SimulateRestart()`
  * Example: `browse_e2e_test.go`, `capture_e2e_test.go`
* **Navigation Tests:** `*_navigation_test.go` - Test UI/navigation with in-memory storage
  * Location: `internal/cli/intents/`
  * Use `e2e.SetupWithMemory(GinkgoT())` for fast setup
  * Focus on key presses and screen transitions
  * Example: `browse_navigation_test.go`, `capture_navigation_test.go`

**When to use each type:**
- Use **E2E tests** when you need to verify data persistence across app restarts
- Use **Navigation tests** when you need to verify UI behavior without persistence concerns
- Use **Unit tests** for testing individual functions in isolation
- Use **Integration tests** for testing interactions between services/repositories

See `AGENTS.md` section "E2E Test Structure" for detailed examples.

### Test Fixtures (Factory Pattern) - MANDATORY

* **Description:** Use the `fixtures` package for generating test data. This package uses `factory-go` for the factory pattern and `gofakeit` for realistic fake data generation. **Always use fixtures instead of manually creating test objects.**

* **Location:** `internal/testutil/fixtures/`

* **Why Use Fixtures:**
  1. **Consistency** - All test data follows domain validation rules
  2. **Realism** - Factory-generated data uses realistic values (company names, dates, etc.)
  3. **Maintainability** - Single source of truth for test object creation
  4. **Reproducibility** - Use `SetSeed()` for deterministic tests
  5. **Less Boilerplate** - Quick helpers reduce test setup code

* **Available Factories:**

  | Factory | Domain Object | Quick Helper |
  |---------|--------------|--------------|
  | `EventFactory` | `career.CareerEvent` | `Event(id)`, `EventWith(...)`, `Events(n)` |
  | `BurstFactory` | `career.Burst` | `Burst(id, eventIDs...)`, `BurstConfirmed(...)`, `Bursts(n, events)` |
  | `FactFactory` | `career.Fact` | `Fact(id, sourceEventID)`, `FactFromBurst(...)`, `Facts(n, events)` |

* **Usage Examples:**

  ```go
  import "github.com/baphled/kariya/internal/testutil/fixtures"

  // Quick helpers for minimal valid objects
  event := fixtures.Event("evt-1")
  burst := fixtures.Burst("burst-1", "evt-1", "evt-2")
  fact := fixtures.Fact("fact-1", "evt-1")

  // Factory with realistic random data
  event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)

  // Factory with field overrides
  event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
      "Company": "TechCorp",
      "Project": "Platform",
  }).(*career.CareerEvent)

  // Batch creation
  events := fixtures.Events(10)
  bursts := fixtures.Bursts(5, events)
  facts := fixtures.Facts(20, events)

  // Reproducible tests with seed
  fixtures.SetSeed(42)
  event1 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
  ```

* **Rules:**

  1. **ALWAYS use fixtures** - Never manually construct domain objects in tests
  2. **Use quick helpers for simple tests** - `Event("id")` is cleaner than factory calls
  3. **Use factories for realistic data** - When tests need varied, realistic data
  4. **Use `SetSeed()` for determinism** - When test assertions depend on specific values
  5. **Link related objects** - Use `Bursts(n, events)` to properly link bursts to events

* **Anti-patterns (DO NOT DO):**

  ```go
  // BAD: Manual object construction
  event := &career.CareerEvent{
      ID:        "test-1",
      Text:      "Some text",
      Date:      time.Now(),
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
  }

  // GOOD: Use fixtures
  event := fixtures.Event("test-1")
  ```

## Code Style

1. Format code with the standard **gofmt** tool (or run `go fmt`). This automatically formats code (indentation, spacing, etc.) according to Go conventions. Additionally, run **`go vet`** to catch common issues (unused variables, misuse of `unsafe`, etc.). These tools should be part of your development/CI process so that code is always formatted and vetted.
2. Adopt idiomatic naming conventions. Package names are short and lowercase. Exported functions, types, and variables have descriptive names starting with a capital letter (with a comment above them for documentation). Avoid underscores in names; use MixedCaps (CamelCase) for multi-word identifiers instead. Prefer short, concise names for variables in small scopes (e.g., loop indices, receiver names).
3. Keep the code base idiomatic and simple. Go favors composition over complex inheritance and clear code over clever code. Use standard library functions and idioms in preference to reinventing solutions. When in doubt, refer to “Effective Go” or official style comments for guidance.

## Code Comments

Use Go’s standard documentation comments for packages, types, and functions. Begin each comment with the name of the item it describes (for example, `// Add adds two numbers.` above an `Add` function). These comments are processed by Go’s documentation tools (like `godoc` or displayed on pkg.go.dev) to produce user-facing docs. In general, comment any exported function or type with a brief description. Non-exported code can have comments if the intent or logic is complex. Avoid writing comments that restate the code – instead, explain the rationale or important details that aren’t obvious from the code itself.


> **Important**: For detailed, mandatory documentation rules including package comments, exported identifiers, and inline comment prohibitions, see **[GO_DOCUMENTATION_RULES.md](../conventions/GO_DOCUMENTATION_RULES.md)**. These rules are enforced by `golangci-lint` (stylecheck, revive, godot).
