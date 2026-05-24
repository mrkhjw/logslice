// Package buffer provides a fixed-size ring buffer for log entries,
// useful for keeping a sliding window of the most recent N entries.
package buffer

import (
	"sync"

	"github.com/yourorg/logslice/internal/parser"
)

// Buffer is a thread-safe fixed-capacity ring buffer of log entries.
type Buffer struct {
	mu       sync.Mutex
	entries  []parser.Entry
	cap      int
	head     int
	count    int
}

// New creates a new ring buffer with the given capacity.
// If capacity is less than 1, it defaults to 1.
func New(capacity int) *Buffer {
	if capacity < 1 {
		capacity = 1
	}
	return &Buffer{
		entries: make([]parser.Entry, capacity),
		cap:     capacity,
	}
}

// Push adds an entry to the buffer, overwriting the oldest entry
// when the buffer is full.
func (b *Buffer) Push(e parser.Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	index := (b.head + b.count) % b.cap
	b.entries[index] = e

	if b.count < b.cap {
		b.count++
	} else {
		// Overwrite oldest: advance head
		b.head = (b.head + 1) % b.cap
	}
}

// Entries returns a snapshot of all buffered entries in insertion order.
func (b *Buffer) Entries() []parser.Entry {
	b.mu.Lock()
	defer b.mu.Unlock()

	result := make([]parser.Entry, b.count)
	for i := 0; i < b.count; i++ {
		result[i] = b.entries[(b.head+i)%b.cap]
	}
	return result
}

// Len returns the number of entries currently in the buffer.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// Reset clears all entries from the buffer.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.head = 0
	b.count = 0
}
