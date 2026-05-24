package buffer_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/buffer"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestNew_DefaultsCapacityToOne(t *testing.T) {
	b := buffer.New(0)
	if b.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", b.Len())
	}
	b.Push(makeEntry("a"))
	b.Push(makeEntry("b"))
	entries := b.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Message != "b" {
		t.Errorf("expected last entry 'b', got %q", entries[0].Message)
	}
}

func TestPush_WithinCapacity(t *testing.T) {
	b := buffer.New(3)
	b.Push(makeEntry("x"))
	b.Push(makeEntry("y"))
	if b.Len() != 2 {
		t.Fatalf("expected 2, got %d", b.Len())
	}
	entries := b.Entries()
	if entries[0].Message != "x" || entries[1].Message != "y" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestPush_OverwritesOldest(t *testing.T) {
	b := buffer.New(3)
	for _, msg := range []string{"a", "b", "c", "d"} {
		b.Push(makeEntry(msg))
	}
	entries := b.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	expected := []string{"b", "c", "d"}
	for i, e := range entries {
		if e.Message != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], e.Message)
		}
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := buffer.New(4)
	b.Push(makeEntry("one"))
	b.Push(makeEntry("two"))
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Len())
	}
	if len(b.Entries()) != 0 {
		t.Error("expected empty entries after reset")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	b := buffer.New(2)
	b.Push(makeEntry("alpha"))
	snap := b.Entries()
	snap[0].Message = "mutated"
	original := b.Entries()
	if original[0].Message == "mutated" {
		t.Error("Entries() should return a copy, not a reference")
	}
}
