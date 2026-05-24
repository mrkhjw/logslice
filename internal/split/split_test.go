package split_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/split"
)

func makeEntries(timestamps ...time.Time) []parser.Entry {
	entries := make([]parser.Entry, len(timestamps))
	for i, ts := range timestamps {
		entries[i] = parser.Entry{
			Timestamp: ts,
			Message:   "msg",
			Level:     "INFO",
		}
	}
	return entries
}

func TestByCount_EvenSplit(t *testing.T) {
	entries := makeEntries(make([]time.Time, 6)...)
	chunks := split.ByCount(entries, 2)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	for _, c := range chunks {
		if len(c) != 2 {
			t.Errorf("expected chunk size 2, got %d", len(c))
		}
	}
}

func TestByCount_UnevenSplit(t *testing.T) {
	entries := makeEntries(make([]time.Time, 5)...)
	chunks := split.ByCount(entries, 2)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if len(chunks[2]) != 1 {
		t.Errorf("expected last chunk size 1, got %d", len(chunks[2]))
	}
}

func TestByCount_ZeroN_ReturnsSingleChunk(t *testing.T) {
	entries := makeEntries(make([]time.Time, 4)...)
	chunks := split.ByCount(entries, 0)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
}

func TestByCount_EmptyInput(t *testing.T) {
	chunks := split.ByCount(nil, 3)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk for empty input, got %d", len(chunks))
	}
}

func TestByDuration_SplitsCorrectly(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	timestamps := []time.Time{
		base,
		base.Add(30 * time.Second),
		base.Add(61 * time.Second),
		base.Add(90 * time.Second),
	}
	entries := makeEntries(timestamps...)
	chunks := split.ByDuration(entries, time.Minute)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 {
		t.Errorf("expected first chunk size 2, got %d", len(chunks[0]))
	}
	if len(chunks[1]) != 2 {
		t.Errorf("expected second chunk size 2, got %d", len(chunks[1]))
	}
}

func TestByDuration_ZeroDuration_ReturnsSingleChunk(t *testing.T) {
	base := time.Now()
	entries := makeEntries(base, base.Add(time.Hour))
	chunks := split.ByDuration(entries, 0)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
}
