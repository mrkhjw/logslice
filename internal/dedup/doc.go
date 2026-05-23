// Package dedup implements log entry deduplication for logslice.
//
// It detects consecutive log entries that share the same severity level and
// message text and collapses them into a single annotated entry. The collapsed
// entry carries a "dedup_count" field recording how many times the message
// appeared, and its message is suffixed with "[repeated N times]".
//
// Deduplication is controlled by an Options struct:
//
//	Window   — only entries within this duration of each other are candidates.
//	MaxCount — flush a group after this many duplicates (0 = unlimited).
//
// Usage:
//
//	result := dedup.Dedup(entries, dedup.DefaultOptions())
package dedup
