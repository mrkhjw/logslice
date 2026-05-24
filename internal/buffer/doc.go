// Package buffer implements a fixed-capacity, thread-safe ring buffer
// for log entries.
//
// It is useful for maintaining a sliding window of the most recently
// seen log entries — for example, to provide context lines around a
// matched entry (similar to grep -B / -A), or to cap memory usage
// when tailing a high-volume log stream.
//
// Usage:
//
//	b := buffer.New(100)
//	b.Push(entry)
//	entries := b.Entries() // ordered oldest-to-newest
//
// When the buffer is full, the oldest entry is silently overwritten.
// All operations are safe for concurrent use.
package buffer
