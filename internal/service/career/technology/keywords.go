package technology

import "strings"

// TechnologyKeyword maps a search keyword to its canonical skill.
// Keywords are lowercase for case-insensitive matching.
// Skills use canonical names for display (e.g., "Go", "PostgreSQL").
type TechnologyKeyword struct {
	Keyword  string // Lowercase search term for matching (e.g., "golang")
	Skill    string // Canonical skill name for display (e.g., "Go")
	Category string // Skill category: backend, frontend, database, devops, cloud, mobile, tooling, testing, ml, data, monitoring, other
}

// technologyKeywords is the master dictionary organized by category.
// This comprehensive list covers ~224 keywords across 14 categories.
var technologyKeywords = []TechnologyKeyword{
	// ========================================
	// BACKEND LANGUAGES & FRAMEWORKS (28)
	// ========================================
	{"go", "Go", "backend"},
	{"golang", "Go", "backend"},
	{"python", "Python", "backend"},
	{"ruby", "Ruby", "backend"},
	{"rails", "Ruby on Rails", "backend"},
	{"ruby on rails", "Ruby on Rails", "backend"},
	{"java", "Java", "backend"},
	{"spring", "Spring", "backend"},
	{"spring boot", "Spring Boot", "backend"},
	{"scala", "Scala", "backend"},
	{"rust", "Rust", "backend"},
	{"c#", "C#", "backend"},
	{"csharp", "C#", "backend"},
	{".net", ".NET", "backend"},
	{"dotnet", ".NET", "backend"},
	{"node", "Node.js", "backend"},
	{"nodejs", "Node.js", "backend"},
	{"node.js", "Node.js", "backend"},
	{"express", "Express.js", "backend"},
	{"expressjs", "Express.js", "backend"},
	{"php", "PHP", "backend"},
	{"laravel", "Laravel", "backend"},
	{"symfony", "Symfony", "backend"},
	{"elixir", "Elixir", "backend"},
	{"phoenix", "Phoenix", "backend"},

	// ========================================
	// FRONTEND FRAMEWORKS & LIBRARIES (35)
	// ========================================
	{"react", "React", "frontend"},
	{"reactjs", "React", "frontend"},
	{"react.js", "React", "frontend"},
	{"vue", "Vue.js", "frontend"},
	{"vuejs", "Vue.js", "frontend"},
	{"vue.js", "Vue.js", "frontend"},
	{"angular", "Angular", "frontend"},
	{"angularjs", "Angular", "frontend"},
	{"svelte", "Svelte", "frontend"},
	{"typescript", "TypeScript", "frontend"},
	{"javascript", "JavaScript", "frontend"},
	{"nextjs", "Next.js", "frontend"},
	{"next.js", "Next.js", "frontend"},
	{"nuxt", "Nuxt.js", "frontend"},
	{"gatsby", "Gatsby", "frontend"},
	{"ember", "Ember.js", "frontend"},
	{"backbone", "Backbone.js", "frontend"},
	{"jquery", "jQuery", "frontend"},
	{"html", "HTML", "frontend"},
	{"css", "CSS", "frontend"},
	{"tailwind", "Tailwind CSS", "frontend"},
	{"tailwindcss", "Tailwind CSS", "frontend"},
	{"bootstrap", "Bootstrap", "frontend"},
	{"sass", "Sass", "frontend"},
	{"scss", "Sass", "frontend"},
	{"less", "Less", "frontend"},
	{"webpack", "Webpack", "frontend"},
	{"vite", "Vite", "frontend"},
	{"rollup", "Rollup", "frontend"},
	{"parcel", "Parcel", "frontend"},

	// ========================================
	// DATABASES (19)
	// ========================================
	{"postgresql", "PostgreSQL", "database"},
	{"postgres", "PostgreSQL", "database"},
	{"mysql", "MySQL", "database"},
	{"mariadb", "MariaDB", "database"},
	{"mongodb", "MongoDB", "database"},
	{"mongo", "MongoDB", "database"},
	{"redis", "Redis", "database"},
	{"elasticsearch", "Elasticsearch", "database"},
	{"elastic", "Elasticsearch", "database"},
	{"dynamodb", "DynamoDB", "database"},
	{"cassandra", "Cassandra", "database"},
	{"couchdb", "CouchDB", "database"},
	{"sqlite", "SQLite", "database"},
	{"oracle", "Oracle DB", "database"},
	{"sql server", "SQL Server", "database"},
	{"mssql", "SQL Server", "database"},
	{"neo4j", "Neo4j", "database"},
	{"influxdb", "InfluxDB", "database"},
	{"timescaledb", "TimescaleDB", "database"},

	// ========================================
	// DEVOPS & INFRASTRUCTURE (24)
	// ========================================
	{"kubernetes", "Kubernetes", "devops"},
	{"k8s", "Kubernetes", "devops"},
	{"docker", "Docker", "devops"},
	{"terraform", "Terraform", "devops"},
	{"ansible", "Ansible", "devops"},
	{"puppet", "Puppet", "devops"},
	{"chef", "Chef", "devops"},
	{"jenkins", "Jenkins", "devops"},
	{"circleci", "CircleCI", "devops"},
	{"circle ci", "CircleCI", "devops"},
	{"github actions", "GitHub Actions", "devops"},
	{"gitlab ci", "GitLab CI", "devops"},
	{"travis", "Travis CI", "devops"},
	{"helm", "Helm", "devops"},
	{"prometheus", "Prometheus", "devops"},
	{"grafana", "Grafana", "devops"},
	{"datadog", "Datadog", "devops"},
	{"new relic", "New Relic", "devops"},
	{"nginx", "Nginx", "devops"},
	{"apache", "Apache", "devops"},
	{"istio", "Istio", "devops"},
	{"envoy", "Envoy", "devops"},
	{"consul", "Consul", "devops"},
	{"vault", "Vault", "devops"},

	// ========================================
	// CLOUD PLATFORMS & SERVICES (23)
	// ========================================
	{"aws", "AWS", "cloud"},
	{"amazon web services", "AWS", "cloud"},
	{"ec2", "AWS EC2", "cloud"},
	{"s3", "AWS S3", "cloud"},
	{"lambda", "AWS Lambda", "cloud"},
	{"rds", "AWS RDS", "cloud"},
	{"ecs", "AWS ECS", "cloud"},
	{"eks", "AWS EKS", "cloud"},
	{"sqs", "AWS SQS", "cloud"},
	{"sns", "AWS SNS", "cloud"},
	{"cloudformation", "CloudFormation", "cloud"},
	{"gcp", "Google Cloud", "cloud"},
	{"google cloud", "Google Cloud", "cloud"},
	{"bigquery", "BigQuery", "cloud"},
	{"cloud run", "Cloud Run", "cloud"},
	{"gke", "Google Kubernetes Engine", "cloud"},
	{"azure", "Azure", "cloud"},
	{"heroku", "Heroku", "cloud"},
	{"vercel", "Vercel", "cloud"},
	{"netlify", "Netlify", "cloud"},
	{"cloudflare", "Cloudflare", "cloud"},
	{"digitalocean", "DigitalOcean", "cloud"},
	{"linode", "Linode", "cloud"},

	// ========================================
	// MOBILE DEVELOPMENT (10)
	// ========================================
	{"ios", "iOS", "mobile"},
	{"android", "Android", "mobile"},
	{"swift", "Swift", "mobile"},
	{"objective-c", "Objective-C", "mobile"},
	{"objective c", "Objective-C", "mobile"},
	{"kotlin", "Kotlin", "mobile"},
	{"react native", "React Native", "mobile"},
	{"flutter", "Flutter", "mobile"},
	{"xamarin", "Xamarin", "mobile"},
	{"ionic", "Ionic", "mobile"},

	// ========================================
	// TOOLING & PROTOCOLS (19)
	// ========================================
	{"git", "Git", "tooling"},
	{"github", "GitHub", "tooling"},
	{"gitlab", "GitLab", "tooling"},
	{"bitbucket", "Bitbucket", "tooling"},
	{"jira", "Jira", "tooling"},
	{"confluence", "Confluence", "tooling"},
	{"slack", "Slack", "tooling"},
	{"figma", "Figma", "tooling"},
	{"sketch", "Sketch", "tooling"},
	{"graphql", "GraphQL", "tooling"},
	{"rest", "REST API", "tooling"},
	{"rest api", "REST API", "tooling"},
	{"grpc", "gRPC", "tooling"},
	{"rabbitmq", "RabbitMQ", "tooling"},
	{"kafka", "Kafka", "tooling"},
	{"oauth", "OAuth", "tooling"},
	{"jwt", "JWT", "tooling"},
	{"websocket", "WebSocket", "tooling"},
	{"protobuf", "Protocol Buffers", "tooling"},

	// ========================================
	// TESTING FRAMEWORKS (15)
	// ========================================
	{"jest", "Jest", "testing"},
	{"mocha", "Mocha", "testing"},
	{"chai", "Chai", "testing"},
	{"jasmine", "Jasmine", "testing"},
	{"pytest", "Pytest", "testing"},
	{"junit", "JUnit", "testing"},
	{"testng", "TestNG", "testing"},
	{"rspec", "RSpec", "testing"},
	{"cypress", "Cypress", "testing"},
	{"selenium", "Selenium", "testing"},
	{"playwright", "Playwright", "testing"},
	{"webdriverio", "WebDriverIO", "testing"},
	{"cucumber", "Cucumber", "testing"},
	{"postman", "Postman", "testing"},
	{"insomnia", "Insomnia", "testing"},

	// ========================================
	// BUILD TOOLS & PACKAGE MANAGERS (12)
	// ========================================
	{"maven", "Maven", "tooling"},
	{"gradle", "Gradle", "tooling"},
	{"make", "Make", "tooling"},
	{"bazel", "Bazel", "tooling"},
	{"npm", "npm", "tooling"},
	{"yarn", "Yarn", "tooling"},
	{"pnpm", "pnpm", "tooling"},
	{"pip", "pip", "tooling"},
	{"poetry", "Poetry", "tooling"},
	{"bundler", "Bundler", "tooling"},
	{"cargo", "Cargo", "tooling"},
	{"cmake", "CMake", "tooling"},

	// ========================================
	// MACHINE LEARNING & DATA SCIENCE (10)
	// ========================================
	{"tensorflow", "TensorFlow", "ml"},
	{"pytorch", "PyTorch", "ml"},
	{"scikit-learn", "scikit-learn", "ml"},
	{"sklearn", "scikit-learn", "ml"},
	{"pandas", "Pandas", "ml"},
	{"numpy", "NumPy", "ml"},
	{"keras", "Keras", "ml"},
	{"jupyter", "Jupyter", "ml"},
	{"openai", "OpenAI", "ml"},
	{"langchain", "LangChain", "ml"},

	// ========================================
	// DATA ENGINEERING (9)
	// ========================================
	{"apache spark", "Apache Spark", "data"},
	{"spark", "Apache Spark", "data"},
	{"airflow", "Apache Airflow", "data"},
	{"databricks", "Databricks", "data"},
	{"snowflake", "Snowflake", "data"},
	{"dbt", "dbt", "data"},
	{"hadoop", "Hadoop", "data"},
	{"hive", "Hive", "data"},
	{"presto", "Presto", "data"},

	// ========================================
	// MONITORING & OBSERVABILITY (7)
	// ========================================
	{"splunk", "Splunk", "monitoring"},
	{"elk stack", "ELK Stack", "monitoring"},
	{"kibana", "Kibana", "monitoring"},
	{"logstash", "Logstash", "monitoring"},
	{"sentry", "Sentry", "monitoring"},
	{"pagerduty", "PagerDuty", "monitoring"},
	{"honeycomb", "Honeycomb", "monitoring"},

	// ========================================
	// DOCUMENTATION TOOLS (6)
	// ========================================
	{"swagger", "Swagger", "tooling"},
	{"openapi", "OpenAPI", "tooling"},
	{"redoc", "Redoc", "tooling"},
	{"docusaurus", "Docusaurus", "tooling"},
	{"mkdocs", "MkDocs", "tooling"},
	{"sphinx", "Sphinx", "tooling"},

	// ========================================
	// OPERATING SYSTEMS (7)
	// ========================================
	{"linux", "Linux", "devops"},
	{"unix", "Unix", "devops"},
	{"ubuntu", "Ubuntu", "devops"},
	{"debian", "Debian", "devops"},
	{"centos", "CentOS", "devops"},
	{"macos", "macOS", "devops"},
	{"windows server", "Windows Server", "devops"},

	// ========================================
	// MODERN RUNTIMES & FRAMEWORKS (8)
	// ========================================
	{"deno", "Deno", "backend"},
	{"bun", "Bun", "backend"},
	{"astro", "Astro", "frontend"},
	{"remix", "Remix", "frontend"},
	{"solidjs", "SolidJS", "frontend"},
	{"qwik", "Qwik", "frontend"},
	{"fresh", "Fresh", "frontend"},
	{"hono", "Hono", "backend"},
}

// GetTechnologyKeywords returns the complete list of technology keywords.
// This is the master dictionary used for skill inference.
func GetTechnologyKeywords() []TechnologyKeyword {
	return technologyKeywords
}

// GetKeywordMap returns a map for O(1) keyword lookup.
// This is used by the inference service for fast detection.
func GetKeywordMap() map[string]TechnologyKeyword {
	keywordMap := make(map[string]TechnologyKeyword, len(technologyKeywords))
	for _, kw := range technologyKeywords {
		keywordMap[kw.Keyword] = kw
	}
	return keywordMap
}

// GetCategoryForSkillName returns the category for a skill name by matching
// against both canonical skill names and keywords (case-insensitive).
// Returns empty string if no match found.
func GetCategoryForSkillName(name string) string {
	if name == "" {
		return ""
	}

	lower := strings.ToLower(name)

	for _, kw := range technologyKeywords {
		if strings.ToLower(kw.Skill) == lower || kw.Keyword == lower {
			return kw.Category
		}
	}

	return ""
}

// GetAllSkillNames returns unique canonical skill names.
// Useful for deduplication and reporting.
func GetAllSkillNames() []string {
	seen := make(map[string]bool)
	names := []string{}

	for _, kw := range technologyKeywords {
		if !seen[kw.Skill] {
			seen[kw.Skill] = true
			names = append(names, kw.Skill)
		}
	}

	return names
}
