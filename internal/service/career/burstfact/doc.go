// Package burstfact provides classifiers and detectors for identifying skill
// bursts and extracting facts from career events.
//
// # Overview
//
// The burstfact package implements:
// - Burst detection: Identifying periods of concentrated skill usage
// - Fact extraction: Pulling quantifiable achievements from events
// - Similarity scoring: Matching events to detect skill patterns
// - Temporal grouping: Organizing events by time proximity
//
// # Usage Example
//
//	detector := burstfact.NewBurstDetector()
//	bursts, err := detector.DetectBursts(ctx, events)
package burstfact
