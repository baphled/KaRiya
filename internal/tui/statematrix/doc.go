// Package statematrix provides state transition documentation generation.
//
// # Overview
//
// The statematrix package analyzes intent source code to extract state
// transition information and generate documentation. It parses Go source
// files to build a complete picture of state machines.
//
// # Features
//
//   - State extraction from intent constants
//   - Transition detection from Update methods
//   - Mermaid diagram generation
//   - Documentation output in Markdown
//
// # Usage
//
// Generate state matrix for an intent:
//
//	make generate-state-matrix
//
// Or programmatically:
//
//	matrix := statematrix.New("intents/myintent")
//	matrix.Parse()
//	matrix.GenerateDocs("output.md")
package statematrix
