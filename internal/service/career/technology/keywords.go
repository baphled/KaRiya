// Package technology provides skill extraction, focus area analysis, and
// keyword-based categorization for career technologies.
package technology

import (
	"context"
	"sort"
	"strings"
	"unicode"

	careerRepo "github.com/baphled/kariya/internal/repository/career"
)

// Entry maps a lowercase search keyword to a display-friendly skill name
// and its canonical category. The keyword is used for case-insensitive
// matching against user-entered skill names.
type Entry struct {
	Keyword  string
	Skill    string
	Category string
}

// TechnologyKeyword is an alias for Entry for backward compatibility.
// Deprecated: Use Entry instead.
//
//nolint:revive // Intentional naming for backward compatibility
type TechnologyKeyword = Entry

// Keywords is the canonical keyword dictionary mapping technology names to
// their skill categories. Each entry has a lowercase keyword for matching,
// a display-friendly skill name, and a category from
// constants.AllSkillCategories().
var Keywords = []Entry{

	// --- backend ---
	{Keyword: "go", Skill: "Go", Category: "backend"},
	{Keyword: "golang", Skill: "Go", Category: "backend"},
	{Keyword: "ruby", Skill: "Ruby", Category: "backend"},
	{Keyword: "python", Skill: "Python", Category: "backend"},
	{Keyword: "java", Skill: "Java", Category: "backend"},
	{Keyword: "c#", Skill: "C#", Category: "backend"},
	{Keyword: "c", Skill: "C", Category: "backend"},
	{Keyword: "c++", Skill: "C++", Category: "backend"},
	{Keyword: "rust", Skill: "Rust", Category: "backend"},
	{Keyword: "scala", Skill: "Scala", Category: "backend"},
	{Keyword: "kotlin", Skill: "Kotlin", Category: "backend"},
	{Keyword: "elixir", Skill: "Elixir", Category: "backend"},
	{Keyword: "erlang", Skill: "Erlang", Category: "backend"},
	{Keyword: "haskell", Skill: "Haskell", Category: "backend"},
	{Keyword: "clojure", Skill: "Clojure", Category: "backend"},
	{Keyword: "php", Skill: "PHP", Category: "backend"},
	{Keyword: "perl", Skill: "Perl", Category: "backend"},
	{Keyword: "node.js", Skill: "Node.js", Category: "backend"},
	{Keyword: "nodejs", Skill: "Node.js", Category: "backend"},
	{Keyword: "express", Skill: "Express", Category: "backend"},
	{Keyword: "django", Skill: "Django", Category: "backend"},
	{Keyword: "flask", Skill: "Flask", Category: "backend"},
	{Keyword: "rails", Skill: "Rails", Category: "backend"},
	{Keyword: "ruby on rails", Skill: "Ruby on Rails", Category: "backend"},
	{Keyword: "spring", Skill: "Spring", Category: "backend"},
	{Keyword: "spring boot", Skill: "Spring Boot", Category: "backend"},
	{Keyword: "grpc", Skill: "gRPC", Category: "backend"},
	{Keyword: "graphql", Skill: "GraphQL", Category: "backend"},
	{Keyword: "rest", Skill: "REST", Category: "backend"},
	{Keyword: "api design", Skill: "API Design", Category: "backend"},
	{Keyword: "coffeescript", Skill: "CoffeeScript", Category: "backend"},
	{Keyword: "zend", Skill: "Zend", Category: "backend"},
	{Keyword: "wordpress", Skill: "WordPress", Category: "backend"},
	{Keyword: "sidekiq", Skill: "Sidekiq", Category: "backend"},
	{Keyword: "ajax", Skill: "AJAX", Category: "backend"},
	{Keyword: "soap", Skill: "SOAP", Category: "backend"},
	{Keyword: "api", Skill: "API", Category: "backend"},
	{Keyword: "websockets", Skill: "WebSockets", Category: "backend"},
	{Keyword: "sinatra", Skill: "Sinatra", Category: "backend"},
	{Keyword: "laravel", Skill: "Laravel", Category: "backend"},
	{Keyword: "fastapi", Skill: "FastAPI", Category: "backend"},
	{Keyword: ".net", Skill: ".NET", Category: "backend"},
	{Keyword: "dotnet", Skill: ".NET", Category: "backend"},
	{Keyword: "objective-c", Skill: "Objective-C", Category: "backend"},
	{Keyword: "groovy", Skill: "Groovy", Category: "backend"},
	{Keyword: "lua", Skill: "Lua", Category: "backend"},
	{Keyword: "r", Skill: "R", Category: "backend"},
	{Keyword: "dart", Skill: "Dart", Category: "backend"},
	{Keyword: "f#", Skill: "F#", Category: "backend"},
	{Keyword: "ocaml", Skill: "OCaml", Category: "backend"},
	{Keyword: "zig", Skill: "Zig", Category: "backend"},
	{Keyword: "delphi", Skill: "Delphi", Category: "backend"},
	{Keyword: "asp.net", Skill: "ASP.NET", Category: "backend"},
	{Keyword: "nestjs", Skill: "NestJS", Category: "backend"},
	{Keyword: "api development", Skill: "API Development", Category: "backend"},
	{Keyword: "api versioning", Skill: "API Versioning", Category: "backend"},
	{Keyword: "backend development", Skill: "Backend Development", Category: "backend"},
	{Keyword: "backend engineering", Skill: "Backend Engineering", Category: "backend"},
	{Keyword: "cms", Skill: "CMS", Category: "backend"},
	{Keyword: "crud", Skill: "CRUD", Category: "backend"},
	{Keyword: "error handling", Skill: "Error Handling", Category: "backend"},
	{Keyword: "full-stack development", Skill: "Full-stack Development", Category: "backend"},
	{Keyword: "full stack development", Skill: "Full-stack Development", Category: "backend"},
	{Keyword: "integration", Skill: "Integration", Category: "backend"},
	{Keyword: "pagination", Skill: "Pagination", Category: "backend"},
	{Keyword: "serialization", Skill: "Serialization", Category: "backend"},
	{Keyword: "validation", Skill: "Validation", Category: "backend"},
	{Keyword: "web development", Skill: "Web Development", Category: "backend"},
	{Keyword: "web", Skill: "Web", Category: "backend"},
	{Keyword: "websocket", Skill: "WebSocket", Category: "backend"},
	{Keyword: "generics", Skill: "Generics", Category: "backend"},
	{Keyword: "go generics", Skill: "Go Generics", Category: "backend"},

	// --- frontend ---
	{Keyword: "javascript", Skill: "JavaScript", Category: "frontend"},
	{Keyword: "typescript", Skill: "TypeScript", Category: "frontend"},
	{Keyword: "react", Skill: "React", Category: "frontend"},
	{Keyword: "reactjs", Skill: "React", Category: "frontend"},
	{Keyword: "react.js", Skill: "React", Category: "frontend"},
	{Keyword: "angular", Skill: "Angular", Category: "frontend"},
	{Keyword: "vue", Skill: "Vue", Category: "frontend"},
	{Keyword: "vue.js", Skill: "Vue.js", Category: "frontend"},
	{Keyword: "svelte", Skill: "Svelte", Category: "frontend"},
	{Keyword: "next.js", Skill: "Next.js", Category: "frontend"},
	{Keyword: "nextjs", Skill: "Next.js", Category: "frontend"},
	{Keyword: "html", Skill: "HTML", Category: "frontend"},
	{Keyword: "css", Skill: "CSS", Category: "frontend"},
	{Keyword: "sass", Skill: "Sass", Category: "frontend"},
	{Keyword: "tailwind", Skill: "Tailwind CSS", Category: "frontend"},
	{Keyword: "webpack", Skill: "Webpack", Category: "frontend"},
	{Keyword: "jquery", Skill: "jQuery", Category: "frontend"},
	{Keyword: "redux", Skill: "Redux", Category: "frontend"},
	{Keyword: "backbonejs", Skill: "Backbone.js", Category: "frontend"},
	{Keyword: "backbone.js", Skill: "Backbone.js", Category: "frontend"},
	{Keyword: "responsive design", Skill: "Responsive Design", Category: "frontend"},
	{Keyword: "accessibility", Skill: "Accessibility", Category: "frontend"},
	{Keyword: "storybook", Skill: "Storybook", Category: "frontend"},
	{Keyword: "seo", Skill: "SEO", Category: "frontend"},
	{Keyword: "design systems", Skill: "Design Systems", Category: "frontend"},
	{Keyword: "visual design", Skill: "Visual Design", Category: "frontend"},
	{Keyword: "ember", Skill: "Ember.js", Category: "frontend"},
	{Keyword: "ember.js", Skill: "Ember.js", Category: "frontend"},
	{Keyword: "bootstrap", Skill: "Bootstrap", Category: "frontend"},
	{Keyword: "material ui", Skill: "Material UI", Category: "frontend"},
	{Keyword: "d3.js", Skill: "D3.js", Category: "frontend"},
	{Keyword: "gatsby", Skill: "Gatsby", Category: "frontend"},
	{Keyword: "nuxt", Skill: "Nuxt.js", Category: "frontend"},
	{Keyword: "nuxt.js", Skill: "Nuxt.js", Category: "frontend"},
	{Keyword: "less", Skill: "Less", Category: "frontend"},
	{Keyword: "styled-components", Skill: "Styled Components", Category: "frontend"},
	{Keyword: "bubble tea", Skill: "Bubble Tea", Category: "frontend"},
	{Keyword: "bubbletea", Skill: "Bubble Tea", Category: "frontend"},
	{Keyword: "lipgloss", Skill: "Lipgloss", Category: "frontend"},
	{Keyword: "vite", Skill: "Vite", Category: "frontend"},
	{Keyword: "component library", Skill: "Component Library", Category: "frontend"},
	{Keyword: "components", Skill: "Components", Category: "frontend"},
	{Keyword: "form components", Skill: "Form Components", Category: "frontend"},
	{Keyword: "filter design", Skill: "Filter Design", Category: "frontend"},
	{Keyword: "frontend development", Skill: "Frontend Development", Category: "frontend"},
	{Keyword: "frontend", Skill: "Frontend", Category: "frontend"},
	{Keyword: "interface design", Skill: "Interface Design", Category: "frontend"},
	{Keyword: "keyboard navigation", Skill: "Keyboard Navigation", Category: "frontend"},
	{Keyword: "layout design", Skill: "Layout Design", Category: "frontend"},
	{Keyword: "modal design", Skill: "Modal Design", Category: "frontend"},
	{Keyword: "navigation", Skill: "Navigation", Category: "frontend"},
	{Keyword: "theme system", Skill: "Theme System", Category: "frontend"},
	{Keyword: "ux", Skill: "UX", Category: "frontend"},
	{Keyword: "ux delivery", Skill: "UX Delivery", Category: "frontend"},
	{Keyword: "user experience", Skill: "User Experience", Category: "frontend"},
	{Keyword: "wizard design", Skill: "Wizard Design", Category: "frontend"},
	{Keyword: "reusable components", Skill: "Reusable Components", Category: "frontend"},

	// --- devops ---
	{Keyword: "docker", Skill: "Docker", Category: "devops"},
	{Keyword: "kubernetes", Skill: "Kubernetes", Category: "devops"},
	{Keyword: "k8s", Skill: "Kubernetes", Category: "devops"},
	{Keyword: "terraform", Skill: "Terraform", Category: "devops"},
	{Keyword: "ansible", Skill: "Ansible", Category: "devops"},
	{Keyword: "puppet", Skill: "Puppet", Category: "devops"},
	{Keyword: "chef", Skill: "Chef", Category: "devops"},
	{Keyword: "jenkins", Skill: "Jenkins", Category: "devops"},
	{Keyword: "github actions", Skill: "GitHub Actions", Category: "devops"},
	{Keyword: "gitlab ci", Skill: "GitLab CI", Category: "devops"},
	{Keyword: "circleci", Skill: "CircleCI", Category: "devops"},
	{Keyword: "ci/cd", Skill: "CI/CD", Category: "devops"},
	{Keyword: "nginx", Skill: "Nginx", Category: "devops"},
	{Keyword: "apache", Skill: "Apache", Category: "devops"},
	{Keyword: "vagrant", Skill: "Vagrant", Category: "devops"},
	{Keyword: "deployment", Skill: "Deployment", Category: "devops"},
	{Keyword: "infrastructure", Skill: "Infrastructure", Category: "devops"},
	{Keyword: "server management", Skill: "Server Management", Category: "devops"},
	{Keyword: "system administration", Skill: "System Administration", Category: "devops"},
	{Keyword: "shell", Skill: "Shell Scripting", Category: "devops"},
	{Keyword: "bash", Skill: "Bash", Category: "devops"},
	{Keyword: "cron", Skill: "Cron", Category: "devops"},
	{Keyword: "git flow", Skill: "Git Flow", Category: "devops"},
	{Keyword: "git hooks", Skill: "Git Hooks", Category: "devops"},
	{Keyword: "operations", Skill: "Operations", Category: "devops"},
	{Keyword: "networking", Skill: "Networking", Category: "devops"},
	{Keyword: "helm", Skill: "Helm", Category: "devops"},
	{Keyword: "argocd", Skill: "ArgoCD", Category: "devops"},
	{Keyword: "packer", Skill: "Packer", Category: "devops"},
	{Keyword: "consul", Skill: "Consul", Category: "devops"},
	{Keyword: "istio", Skill: "Istio", Category: "devops"},
	{Keyword: "travis ci", Skill: "Travis CI", Category: "devops"},
	{Keyword: "capistrano", Skill: "Capistrano", Category: "devops"},
	{Keyword: "devops", Skill: "DevOps", Category: "devops"},
	{Keyword: "embedded systems", Skill: "Embedded Systems", Category: "devops"},
	{Keyword: "firmware development", Skill: "Firmware Development", Category: "devops"},
	{Keyword: "incident response", Skill: "Incident Response", Category: "devops"},
	{Keyword: "iot", Skill: "IoT", Category: "devops"},
	{Keyword: "migration", Skill: "Migration", Category: "devops"},
	{Keyword: "migrations", Skill: "Migrations", Category: "devops"},
	{Keyword: "operational support", Skill: "Operational Support", Category: "devops"},
	{Keyword: "operational systems", Skill: "Operational Systems", Category: "devops"},
	{Keyword: "production support", Skill: "Production Support", Category: "devops"},
	{Keyword: "reliability engineering", Skill: "Reliability Engineering", Category: "devops"},
	{Keyword: "reliability", Skill: "Reliability", Category: "devops"},
	{Keyword: "system stability", Skill: "System Stability", Category: "devops"},
	{Keyword: "system integration", Skill: "System Integration", Category: "devops"},

	// --- database ---
	{Keyword: "postgresql", Skill: "PostgreSQL", Category: "database"},
	{Keyword: "postgres", Skill: "PostgreSQL", Category: "database"},
	{Keyword: "mysql", Skill: "MySQL", Category: "database"},
	{Keyword: "mongodb", Skill: "MongoDB", Category: "database"},
	{Keyword: "redis", Skill: "Redis", Category: "database"},
	{Keyword: "elasticsearch", Skill: "Elasticsearch", Category: "database"},
	{Keyword: "sqlite", Skill: "SQLite", Category: "database"},
	{Keyword: "sql", Skill: "SQL", Category: "database"},
	{Keyword: "dynamodb", Skill: "DynamoDB", Category: "database"},
	{Keyword: "cassandra", Skill: "Cassandra", Category: "database"},
	{Keyword: "couchdb", Skill: "CouchDB", Category: "database"},
	{Keyword: "oracle", Skill: "Oracle DB", Category: "database"},
	{Keyword: "nosql", Skill: "NoSQL", Category: "database"},
	{Keyword: "database design", Skill: "Database Design", Category: "database"},
	{Keyword: "database migrations", Skill: "Database Migrations", Category: "database"},
	{Keyword: "database optimization", Skill: "Database Optimization", Category: "database"},
	{Keyword: "database performance", Skill: "Database Performance", Category: "database"},
	{Keyword: "data modeling", Skill: "Data Modeling", Category: "database"},
	{Keyword: "data integrity", Skill: "Data Integrity", Category: "database"},
	{Keyword: "mariadb", Skill: "MariaDB", Category: "database"},
	{Keyword: "neo4j", Skill: "Neo4j", Category: "database"},
	{Keyword: "memcached", Skill: "Memcached", Category: "database"},
	{Keyword: "influxdb", Skill: "InfluxDB", Category: "database"},
	{Keyword: "timescaledb", Skill: "TimescaleDB", Category: "database"},

	// --- cloud ---
	{Keyword: "aws", Skill: "AWS", Category: "cloud"},
	{Keyword: "gcp", Skill: "GCP", Category: "cloud"},
	{Keyword: "google cloud", Skill: "Google Cloud", Category: "cloud"},
	{Keyword: "azure", Skill: "Azure", Category: "cloud"},
	{Keyword: "heroku", Skill: "Heroku", Category: "cloud"},
	{Keyword: "digitalocean", Skill: "DigitalOcean", Category: "cloud"},
	{Keyword: "serverless", Skill: "Serverless", Category: "cloud"},
	{Keyword: "cloudflare", Skill: "Cloudflare", Category: "cloud"},

	// --- mobile ---
	{Keyword: "react native", Skill: "React Native", Category: "mobile"},
	{Keyword: "flutter", Skill: "Flutter", Category: "mobile"},
	{Keyword: "ios", Skill: "iOS", Category: "mobile"},
	{Keyword: "android", Skill: "Android", Category: "mobile"},
	{Keyword: "xamarin", Skill: "Xamarin", Category: "mobile"},
	{Keyword: "ionic", Skill: "Ionic", Category: "mobile"},
	{Keyword: "swift", Skill: "Swift", Category: "mobile"},

	// --- tooling ---
	{Keyword: "git", Skill: "Git", Category: "tooling"},
	{Keyword: "github", Skill: "GitHub", Category: "tooling"},
	{Keyword: "gitlab", Skill: "GitLab", Category: "tooling"},
	{Keyword: "bitbucket", Skill: "Bitbucket", Category: "tooling"},
	{Keyword: "jira", Skill: "Jira", Category: "tooling"},
	{Keyword: "confluence", Skill: "Confluence", Category: "tooling"},
	{Keyword: "slack", Skill: "Slack", Category: "tooling"},
	{Keyword: "vim", Skill: "Vim", Category: "tooling"},
	{Keyword: "neovim", Skill: "Neovim", Category: "tooling"},
	{Keyword: "vscode", Skill: "VS Code", Category: "tooling"},
	{Keyword: "intellij", Skill: "IntelliJ", Category: "tooling"},
	{Keyword: "postman", Skill: "Postman", Category: "tooling"},
	{Keyword: "swagger", Skill: "Swagger", Category: "tooling"},
	{Keyword: "make", Skill: "Make", Category: "tooling"},
	{Keyword: "cmake", Skill: "CMake", Category: "tooling"},
	{Keyword: "gradle", Skill: "Gradle", Category: "tooling"},
	{Keyword: "maven", Skill: "Maven", Category: "tooling"},
	{Keyword: "npm", Skill: "npm", Category: "tooling"},
	{Keyword: "yarn", Skill: "Yarn", Category: "tooling"},
	{Keyword: "pip", Skill: "pip", Category: "tooling"},
	{Keyword: "homebrew", Skill: "Homebrew", Category: "tooling"},
	{Keyword: "linux", Skill: "Linux", Category: "tooling"},
	{Keyword: "macos", Skill: "macOS", Category: "tooling"},
	{Keyword: "windows", Skill: "Windows", Category: "tooling"},
	{Keyword: "documentation", Skill: "Documentation", Category: "tooling"},
	{Keyword: "technical documentation", Skill: "Technical Documentation", Category: "tooling"},
	{Keyword: "markdown", Skill: "Markdown", Category: "tooling"},
	{Keyword: "xml", Skill: "XML", Category: "tooling"},
	{Keyword: "json", Skill: "JSON", Category: "tooling"},
	{Keyword: "csv", Skill: "CSV", Category: "tooling"},
	{Keyword: "regex", Skill: "Regex", Category: "tooling"},
	{Keyword: "figma", Skill: "Figma", Category: "tooling"},
	{Keyword: "trello", Skill: "Trello", Category: "tooling"},
	{Keyword: "notion", Skill: "Notion", Category: "tooling"},
	{Keyword: "asana", Skill: "Asana", Category: "tooling"},
	{Keyword: "docker compose", Skill: "Docker Compose", Category: "tooling"},
	{Keyword: "openapi", Skill: "OpenAPI", Category: "tooling"},
	{Keyword: "mermaid", Skill: "Mermaid", Category: "tooling"},
	{Keyword: "graphviz", Skill: "Graphviz", Category: "tooling"},
	{Keyword: "yaml", Skill: "YAML", Category: "tooling"},
	{Keyword: "toml", Skill: "TOML", Category: "tooling"},
	{Keyword: "protobuf", Skill: "Protocol Buffers", Category: "tooling"},
	{Keyword: "diagramming", Skill: "Diagramming", Category: "tooling"},
	{Keyword: "technical writing", Skill: "Technical Writing", Category: "tooling"},
	{Keyword: "automation", Skill: "Automation", Category: "tooling"},
	{Keyword: "code generation", Skill: "Code Generation", Category: "tooling"},

	// --- testing ---
	{Keyword: "testing", Skill: "Testing", Category: "testing"},
	{Keyword: "unit testing", Skill: "Unit Testing", Category: "testing"},
	{Keyword: "integration testing", Skill: "Integration Testing", Category: "testing"},
	{Keyword: "e2e testing", Skill: "E2E Testing", Category: "testing"},
	{Keyword: "automated testing", Skill: "Automated Testing", Category: "testing"},
	{Keyword: "jest", Skill: "Jest", Category: "testing"},
	{Keyword: "mocha", Skill: "Mocha", Category: "testing"},
	{Keyword: "rspec", Skill: "RSpec", Category: "testing"},
	{Keyword: "pytest", Skill: "pytest", Category: "testing"},
	{Keyword: "cypress", Skill: "Cypress", Category: "testing"},
	{Keyword: "selenium", Skill: "Selenium", Category: "testing"},
	{Keyword: "ginkgo", Skill: "Ginkgo", Category: "testing"},
	{Keyword: "gomega", Skill: "Gomega", Category: "testing"},
	{Keyword: "code coverage", Skill: "Code Coverage", Category: "testing"},
	{Keyword: "quality assurance", Skill: "Quality Assurance", Category: "testing"},
	{Keyword: "load testing", Skill: "Load Testing", Category: "testing"},
	{Keyword: "performance testing", Skill: "Performance Testing", Category: "testing"},
	{Keyword: "test automation", Skill: "Test Automation", Category: "testing"},
	{Keyword: "cucumber", Skill: "Cucumber", Category: "testing"},
	{Keyword: "playwright", Skill: "Playwright", Category: "testing"},
	{Keyword: "benchmarking", Skill: "Benchmarking", Category: "testing"},

	// --- monitoring ---
	{Keyword: "prometheus", Skill: "Prometheus", Category: "monitoring"},
	{Keyword: "grafana", Skill: "Grafana", Category: "monitoring"},
	{Keyword: "datadog", Skill: "Datadog", Category: "monitoring"},
	{Keyword: "new relic", Skill: "New Relic", Category: "monitoring"},
	{Keyword: "sentry", Skill: "Sentry", Category: "monitoring"},
	{Keyword: "elk", Skill: "ELK Stack", Category: "monitoring"},
	{Keyword: "logging", Skill: "Logging", Category: "monitoring"},
	{Keyword: "structured logging", Skill: "Structured Logging", Category: "monitoring"},
	{Keyword: "alerting", Skill: "Alerting", Category: "monitoring"},
	{Keyword: "metrics", Skill: "Metrics", Category: "monitoring"},
	{Keyword: "observability", Skill: "Observability", Category: "monitoring"},
	{Keyword: "monitoring", Skill: "Monitoring", Category: "monitoring"},
	{Keyword: "splunk", Skill: "Splunk", Category: "monitoring"},
	{Keyword: "pagerduty", Skill: "PagerDuty", Category: "monitoring"},
	{Keyword: "kibana", Skill: "Kibana", Category: "monitoring"},
	{Keyword: "logstash", Skill: "Logstash", Category: "monitoring"},

	// --- data ---
	{Keyword: "etl", Skill: "ETL", Category: "data"},
	{Keyword: "data processing", Skill: "Data Processing", Category: "data"},
	{Keyword: "data export", Skill: "Data Export", Category: "data"},
	{Keyword: "data platforms", Skill: "Data Platforms", Category: "data"},
	{Keyword: "analytics", Skill: "Analytics", Category: "data"},
	{Keyword: "business intelligence", Skill: "Business Intelligence", Category: "data"},
	{Keyword: "reporting", Skill: "Reporting", Category: "data"},
	{Keyword: "apache spark", Skill: "Apache Spark", Category: "data"},
	{Keyword: "kafka", Skill: "Kafka", Category: "data"},
	{Keyword: "rabbitmq", Skill: "RabbitMQ", Category: "data"},
	{Keyword: "streaming", Skill: "Streaming", Category: "data"},
	{Keyword: "rss", Skill: "RSS", Category: "data"},

	// --- ml ---
	{Keyword: "machine learning", Skill: "Machine Learning", Category: "ml"},
	{Keyword: "deep learning", Skill: "Deep Learning", Category: "ml"},
	{Keyword: "tensorflow", Skill: "TensorFlow", Category: "ml"},
	{Keyword: "pytorch", Skill: "PyTorch", Category: "ml"},
	{Keyword: "natural language processing", Skill: "NLP", Category: "ml"},
	{Keyword: "nlp", Skill: "NLP", Category: "ml"},
	{Keyword: "ai", Skill: "AI", Category: "ml"},
	{Keyword: "generative ai", Skill: "Generative AI", Category: "ml"},
	{Keyword: "llms", Skill: "LLMs", Category: "ml"},
	{Keyword: "prompt engineering", Skill: "Prompt Engineering", Category: "ml"},
	{Keyword: "computer vision", Skill: "Computer Vision", Category: "ml"},
	{Keyword: "scikit-learn", Skill: "scikit-learn", Category: "ml"},

	// --- architecture ---
	{Keyword: "microservices", Skill: "Microservices", Category: "architecture"},
	{Keyword: "monolith", Skill: "Monolithic Architecture", Category: "architecture"},
	{Keyword: "soa", Skill: "SOA", Category: "architecture"},
	{Keyword: "service-oriented architecture", Skill: "SOA", Category: "architecture"},
	{Keyword: "domain-driven design", Skill: "Domain-Driven Design", Category: "architecture"},
	{Keyword: "ddd", Skill: "Domain-Driven Design", Category: "architecture"},
	{Keyword: "event-driven", Skill: "Event-Driven Architecture", Category: "architecture"},
	{Keyword: "event sourcing", Skill: "Event Sourcing", Category: "architecture"},
	{Keyword: "cqrs", Skill: "CQRS", Category: "architecture"},
	{Keyword: "distributed systems", Skill: "Distributed Systems", Category: "architecture"},
	{Keyword: "system design", Skill: "System Design", Category: "architecture"},
	{Keyword: "scalability", Skill: "Scalability", Category: "architecture"},
	{Keyword: "high availability", Skill: "High Availability", Category: "architecture"},
	{Keyword: "fault tolerance", Skill: "Fault Tolerance", Category: "architecture"},
	{Keyword: "load balancing", Skill: "Load Balancing", Category: "architecture"},
	{Keyword: "caching", Skill: "Caching", Category: "architecture"},
	{Keyword: "message queues", Skill: "Message Queues", Category: "architecture"},
	{Keyword: "pub/sub", Skill: "Pub/Sub", Category: "architecture"},
	{Keyword: "dependency injection", Skill: "Dependency Injection", Category: "architecture"},
	{Keyword: "design patterns", Skill: "Design Patterns", Category: "architecture"},
	{Keyword: "solid", Skill: "SOLID Principles", Category: "architecture"},
	{Keyword: "clean architecture", Skill: "Clean Architecture", Category: "architecture"},
	{Keyword: "hexagonal architecture", Skill: "Hexagonal Architecture", Category: "architecture"},
	{Keyword: "api gateway", Skill: "API Gateway", Category: "architecture"},
	{Keyword: "state machine", Skill: "State Machine", Category: "architecture"},
	{Keyword: "concurrency", Skill: "Concurrency", Category: "architecture"},
	{Keyword: "multithreading", Skill: "Multithreading", Category: "architecture"},
	{Keyword: "data structures", Skill: "Data Structures", Category: "architecture"},
	{Keyword: "algorithms", Skill: "Algorithms", Category: "architecture"},
	{Keyword: "software architecture", Skill: "Software Architecture", Category: "architecture"},
	{Keyword: "twelve-factor", Skill: "Twelve-Factor App", Category: "architecture"},
	{Keyword: "saga pattern", Skill: "Saga Pattern", Category: "architecture"},
	{Keyword: "circuit breaker", Skill: "Circuit Breaker", Category: "architecture"},
	{Keyword: "rate limiting", Skill: "Rate Limiting", Category: "architecture"},
	{Keyword: "idempotency", Skill: "Idempotency", Category: "architecture"},
	{Keyword: "adr", Skill: "ADR", Category: "architecture"},
	{Keyword: "architecture", Skill: "Architecture", Category: "architecture"},
	{Keyword: "asynchronous processing", Skill: "Asynchronous Processing", Category: "architecture"},
	{Keyword: "backend architecture", Skill: "Backend Architecture", Category: "architecture"},
	{Keyword: "backend systems", Skill: "Backend Systems", Category: "architecture"},
	{Keyword: "background processing", Skill: "Background Processing", Category: "architecture"},
	{Keyword: "builder pattern", Skill: "Builder Pattern", Category: "architecture"},
	{Keyword: "component architecture", Skill: "Component Architecture", Category: "architecture"},
	{Keyword: "high-throughput systems", Skill: "High-throughput Systems", Category: "architecture"},
	{Keyword: "legacy modernization", Skill: "Legacy Modernization", Category: "architecture"},
	{Keyword: "multi-language systems", Skill: "Multi-language Systems", Category: "architecture"},
	{Keyword: "payment systems", Skill: "Payment Systems", Category: "architecture"},
	{Keyword: "payments", Skill: "Payments", Category: "architecture"},
	{Keyword: "order systems", Skill: "Order Systems", Category: "architecture"},
	{Keyword: "real-time systems", Skill: "Real-time Systems", Category: "architecture"},
	{Keyword: "repository pattern", Skill: "Repository Pattern", Category: "architecture"},
	{Keyword: "resilience", Skill: "Resilience", Category: "architecture"},
	{Keyword: "resilient systems", Skill: "Resilient Systems", Category: "architecture"},
	{Keyword: "scalable systems", Skill: "Scalable Systems", Category: "architecture"},
	{Keyword: "scaling systems", Skill: "Scaling Systems", Category: "architecture"},
	{Keyword: "service architecture", Skill: "Service Architecture", Category: "architecture"},
	{Keyword: "service integration", Skill: "Service Integration", Category: "architecture"},
	{Keyword: "service layer", Skill: "Service Layer", Category: "architecture"},
	{Keyword: "type safety", Skill: "Type Safety", Category: "architecture"},
	{Keyword: "domain", Skill: "Domain", Category: "architecture"},

	// --- security (15) ---
	{Keyword: "security", Skill: "Security", Category: "security"},
	{Keyword: "authentication", Skill: "Authentication", Category: "security"},
	{Keyword: "authorization", Skill: "Authorization", Category: "security"},
	{Keyword: "oauth", Skill: "OAuth", Category: "security"},
	{Keyword: "oauth2", Skill: "OAuth 2.0", Category: "security"},
	{Keyword: "jwt", Skill: "JWT", Category: "security"},
	{Keyword: "tls", Skill: "TLS", Category: "security"},
	{Keyword: "ssl", Skill: "SSL", Category: "security"},
	{Keyword: "pki", Skill: "PKI", Category: "security"},
	{Keyword: "pci", Skill: "PCI Compliance", Category: "security"},
	{Keyword: "saml", Skill: "SAML", Category: "security"},
	{Keyword: "active directory", Skill: "Active Directory", Category: "security"},
	{Keyword: "compliance", Skill: "Compliance", Category: "security"},
	{Keyword: "encryption", Skill: "Encryption", Category: "security"},
	{Keyword: "penetration testing", Skill: "Penetration Testing", Category: "security"},
	{Keyword: "owasp", Skill: "OWASP", Category: "security"},
	{Keyword: "sso", Skill: "Single Sign-On", Category: "security"},
	{Keyword: "rbac", Skill: "RBAC", Category: "security"},
	{Keyword: "gdpr", Skill: "GDPR", Category: "security"},
	{Keyword: "soc2", Skill: "SOC 2", Category: "security"},

	// --- practices ---
	{Keyword: "agile", Skill: "Agile", Category: "practices"},
	{Keyword: "scrum", Skill: "Scrum", Category: "practices"},
	{Keyword: "kanban", Skill: "Kanban", Category: "practices"},
	{Keyword: "tdd", Skill: "TDD", Category: "practices"},
	{Keyword: "test-driven development", Skill: "TDD", Category: "practices"},
	{Keyword: "bdd", Skill: "BDD", Category: "practices"},
	{Keyword: "behavior-driven development", Skill: "BDD", Category: "practices"},
	{Keyword: "code review", Skill: "Code Review", Category: "practices"},
	{Keyword: "pair programming", Skill: "Pair Programming", Category: "practices"},
	{Keyword: "mob programming", Skill: "Mob Programming", Category: "practices"},
	{Keyword: "refactoring", Skill: "Refactoring", Category: "practices"},
	{Keyword: "code quality", Skill: "Code Quality", Category: "practices"},
	{Keyword: "debugging", Skill: "Debugging", Category: "practices"},
	{Keyword: "performance optimization", Skill: "Performance Optimization", Category: "practices"},
	{Keyword: "feature flags", Skill: "Feature Flags", Category: "practices"},
	{Keyword: "trunk-based development", Skill: "Trunk-Based Development", Category: "practices"},
	{Keyword: "continuous integration", Skill: "Continuous Integration", Category: "practices"},
	{Keyword: "continuous delivery", Skill: "Continuous Delivery", Category: "practices"},
	{Keyword: "devops practices", Skill: "DevOps Practices", Category: "practices"},
	{Keyword: "incident management", Skill: "Incident Management", Category: "practices"},
	{Keyword: "on-call", Skill: "On-Call", Category: "practices"},
	{Keyword: "retrospectives", Skill: "Retrospectives", Category: "practices"},
	{Keyword: "sprint planning", Skill: "Sprint Planning", Category: "practices"},
	{Keyword: "estimation", Skill: "Estimation", Category: "practices"},
	{Keyword: "technical debt", Skill: "Technical Debt Management", Category: "practices"},
	{Keyword: "documentation practices", Skill: "Documentation Practices", Category: "practices"},
	{Keyword: "release management", Skill: "Release Management", Category: "practices"},
	{Keyword: "version control", Skill: "Version Control", Category: "practices"},
	{Keyword: "lean", Skill: "Lean", Category: "practices"},
	{Keyword: "extreme programming", Skill: "Extreme Programming", Category: "practices"},
	{Keyword: "post-mortems", Skill: "Post-Mortems", Category: "practices"},
	{Keyword: "capacity planning", Skill: "Capacity Planning", Category: "practices"},
	{Keyword: "chaos engineering", Skill: "Chaos Engineering", Category: "practices"},
	{Keyword: "site reliability", Skill: "Site Reliability Engineering", Category: "practices"},
	{Keyword: "sre", Skill: "SRE", Category: "practices"},
	{Keyword: "runbooks", Skill: "Runbooks", Category: "practices"},
	{Keyword: "knowledge sharing", Skill: "Knowledge Sharing", Category: "practices"},
	{Keyword: "mentoring", Skill: "Mentoring", Category: "practices"},
	{Keyword: "stakeholder management", Skill: "Stakeholder Management", Category: "practices"},
	{Keyword: "agile delivery", Skill: "Agile Delivery", Category: "practices"},
	{Keyword: "application maintenance", Skill: "Application Maintenance", Category: "practices"},
	{Keyword: "client collaboration", Skill: "Client Collaboration", Category: "practices"},
	{Keyword: "client work", Skill: "Client Work", Category: "practices"},
	{Keyword: "code standards", Skill: "Code Standards", Category: "practices"},
	{Keyword: "collaboration", Skill: "Collaboration", Category: "practices"},
	{Keyword: "consulting", Skill: "Consulting", Category: "practices"},
	{Keyword: "consulting delivery", Skill: "Consulting Delivery", Category: "practices"},
	{Keyword: "cost optimization", Skill: "Cost Optimization", Category: "practices"},
	{Keyword: "cross-functional collaboration", Skill: "Cross-functional Collaboration", Category: "practices"},
	{Keyword: "cross-team collaboration", Skill: "Cross-team Collaboration", Category: "practices"},
	{Keyword: "delivery", Skill: "Delivery", Category: "practices"},
	{Keyword: "delivery management", Skill: "Delivery Management", Category: "practices"},
	{Keyword: "developer experience", Skill: "Developer Experience", Category: "practices"},
	{Keyword: "end-to-end delivery", Skill: "End-to-End Delivery", Category: "practices"},
	{Keyword: "engineering judgment", Skill: "Engineering Judgment", Category: "practices"},
	{Keyword: "engineering practices", Skill: "Engineering Practices", Category: "practices"},
	{Keyword: "feature delivery", Skill: "Feature Delivery", Category: "practices"},
	{Keyword: "full lifecycle delivery", Skill: "Full Lifecycle Delivery", Category: "practices"},
	{Keyword: "incremental delivery", Skill: "Incremental Delivery", Category: "practices"},
	{Keyword: "innovation", Skill: "Innovation", Category: "practices"},
	{Keyword: "knowledge transfer", Skill: "Knowledge Transfer", Category: "practices"},
	{Keyword: "leadership", Skill: "Leadership", Category: "practices"},
	{Keyword: "learning", Skill: "Learning", Category: "practices"},
	{Keyword: "maintenance", Skill: "Maintenance", Category: "practices"},
	{Keyword: "onboarding", Skill: "Onboarding", Category: "practices"},
	{Keyword: "operational alignment", Skill: "Operational Alignment", Category: "practices"},
	{Keyword: "optimization", Skill: "Optimization", Category: "practices"},
	{Keyword: "performance engineering", Skill: "Performance Engineering", Category: "practices"},
	{Keyword: "presentation", Skill: "Presentation", Category: "practices"},
	{Keyword: "product development", Skill: "Product Development", Category: "practices"},
	{Keyword: "product engineering", Skill: "Product Engineering", Category: "practices"},
	{Keyword: "product management", Skill: "Product Management", Category: "practices"},
	{Keyword: "product support", Skill: "Product Support", Category: "practices"},
	{Keyword: "project delivery", Skill: "Project Delivery", Category: "practices"},
	{Keyword: "project handover", Skill: "Project Handover", Category: "practices"},
	{Keyword: "project management", Skill: "Project Management", Category: "practices"},
	{Keyword: "remote collaboration", Skill: "Remote Collaboration", Category: "practices"},
	{Keyword: "software craft", Skill: "Software Craft", Category: "practices"},
	{Keyword: "software engineering", Skill: "Software Engineering", Category: "practices"},
	{Keyword: "stakeholder collaboration", Skill: "Stakeholder Collaboration", Category: "practices"},
	{Keyword: "strategy", Skill: "Strategy", Category: "practices"},
	{Keyword: "technical advisory", Skill: "Technical Advisory", Category: "practices"},
	{Keyword: "technical debt management", Skill: "Technical Debt Management", Category: "practices"},
	{Keyword: "technical leadership", Skill: "Technical Leadership", Category: "practices"},
	{Keyword: "time management", Skill: "Time Management", Category: "practices"},
	{Keyword: "training", Skill: "Training", Category: "practices"},
	{Keyword: "workflow design", Skill: "Workflow Design", Category: "practices"},
}

// keywordIndex is a pre-built map from lowercase keyword to category for
// constant-time lookups. It is built once at package init.
var keywordIndex map[string]string

// substringKeywords holds keywords sorted by length descending for substring
// matching. Longer keywords are checked first so that more specific matches
// take priority (e.g. "event-driven" before "event").
var substringKeywords []string

func init() {
	keywordIndex = make(map[string]string, len(Keywords))
	for _, kw := range Keywords {
		keywordIndex[kw.Keyword] = kw.Category
	}

	substringKeywords = make([]string, 0, len(keywordIndex))
	for kw := range keywordIndex {
		substringKeywords = append(substringKeywords, kw)
	}
	sort.Slice(substringKeywords, func(i, j int) bool {
		return len(substringKeywords[i]) > len(substringKeywords[j])
	})
}

// GetCategoryForSkillName performs a case-insensitive lookup of the given
// skill name against the keyword dictionary. It first tries an exact match,
// then falls back to substring matching with word-boundary checks.
func GetCategoryForSkillName(name string) string {
	lower := strings.ToLower(name)
	if cat, ok := keywordIndex[lower]; ok {
		return cat
	}

	for _, kw := range substringKeywords {
		if containsKeywordAtBoundary(lower, kw) {
			return keywordIndex[kw]
		}
	}

	return ""
}

// containsKeywordAtBoundary checks whether keyword appears in text with
// word boundaries on both sides. A word boundary is the start/end of the
// string, any non-alphanumeric character, or a trailing plural "s".
func containsKeywordAtBoundary(text, keyword string) bool {
	offset := 0
	for {
		idx := strings.Index(text[offset:], keyword)
		if idx < 0 {
			return false
		}
		idx += offset

		leftOK := idx == 0 || !isAlphanumeric(rune(text[idx-1]))
		endPos := idx + len(keyword)
		rightOK := endPos == len(text) ||
			!isAlphanumeric(rune(text[endPos])) ||
			(text[endPos] == 's' && (endPos+1 == len(text) || !isAlphanumeric(rune(text[endPos+1]))))

		if leftOK && rightOK {
			return true
		}
		offset = idx + 1
	}
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// RecategorizeResult summarises the outcome of a bulk skill recategorization.
type RecategorizeResult struct {
	Total      int
	Updated    int
	Skipped    int
	NoMatch    int
	ByCategory map[string]int
}

// RecategorizeSkills iterates over all skills in the repository and updates
// their category when the keyword dictionary provides a better match than
// the current value. Skills that already have the correct category or have
// no keyword match are left unchanged.
func RecategorizeSkills(ctx context.Context, repo careerRepo.SkillRepository) (*RecategorizeResult, error) {
	skills, err := repo.List(ctx, &careerRepo.SkillListFilters{})
	if err != nil {
		return nil, err
	}

	result := &RecategorizeResult{
		Total:      len(skills),
		ByCategory: make(map[string]int),
	}

	for _, skill := range skills {
		newCategory := GetCategoryForSkillName(skill.Name)
		if newCategory == "" {
			result.NoMatch++
			continue
		}

		if skill.Category == newCategory {
			result.Skipped++
			continue
		}

		skill.Category = newCategory
		if err := repo.Update(ctx, skill); err != nil {
			return nil, err
		}

		result.Updated++
		result.ByCategory[newCategory]++
	}

	return result, nil
}

// GetKeywordMap returns a map from keyword to TechnologyKeyword for backward compatibility.
// Deprecated: Use the keywordIndex map or GetCategoryForSkillName instead.
func GetKeywordMap() map[string]TechnologyKeyword {
	result := make(map[string]TechnologyKeyword, len(Keywords))
	for _, entry := range Keywords {
		result[entry.Keyword] = entry
	}
	return result
}
