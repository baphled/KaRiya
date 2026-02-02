// Package importer provides data import functionality for external sources.
//
// # Overview
//
// The importer package handles importing career data from external sources
// such as LinkedIn exports, JSON files, and other formats. It provides
// validation, transformation, and persistence of imported data.
//
// # Supported Formats
//
//   - LinkedIn profile exports
//   - JSON career data
//   - CSV job history
//   - Custom formats via adapters
//
// # Usage
//
// Import data:
//
//	importer := importer.New(service)
//	events, err := importer.ImportLinkedIn(data)
//
// For more details on import workflows, see the import intent documentation.
package importer
