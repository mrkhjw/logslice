package batch_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/batch"
	"github.com/yourorg/logslice/internal/parser"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeEntries(n int, step time.Duration) []parser.Entry {
	out := make([]parser.Entry, n)
	for i := range out {
		out[i] = parser.Entry{
			Timestamp: base.Add(time.Duration(i) * step),
			Level:     "INFO",
			Message:   "msg",
		}
	}
	return out
}

func TestApply_EmptyInput(t *testing.T) {
	result := batch.Apply(nil, batch.DefaultOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestApply_NoLimits_ReturnsSingleBatch(t *testing.T) {
	entries := makeEntries(10, time.Second)
	result := batch.Apply(entries, batch.Options{})
	if len(result) != 1 || len(result[0]) != 10 {
		t.Fatalf("expected 1 batch of 10, got %d batches", len(result))
	}
}

func TestApply_BySize_EvenSplit(t *testing.T) {
	entries := makeEntries(9, time.Second)
	result := batch.Apply(entries, batch.Options{Size: 3})
	if len(result) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(result))
	}
	for i, b := range result {
		if len(b) != 3 {
			t.Errorf("batch %d: expected 3 entries, got %d", i, len(b))
		}
	}
}

func TestApply_BySize_UnevenSplit(t *testing.T) {
	entries := makeEntries(7, time.Second)
	result := batch.Apply(entries, batch.Options{Size: 3})
	if len(result) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(result))
	}
	if len(result[2]) != 1 {
		t.Errorf("last batch: expected 1 entry, got %d", len(result[2]))
	}
}

func TestApply_ByDuration(t *testing.T) {
	// 6 entries, 10s apart → 3 batches of 2 with a 30s window
	entries := makeEntries(6, 10*time.Second)
	result := batch.Apply(entries, batch.Options{Duration: 25 * time.Second})
	if len(result) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(result))
	}
}

func TestApply_DefaultOptions(t *testing.T) {
	opts := batch.DefaultOptions()
	if opts.Size != 100 {
		t.Errorf("expected default size 100, got %d", opts.Size)
	}
	if opts.Duration != 0 {
		t.Errorf("expected default duration 0, got %v", opts.Duration)
	}
}

func TestApply_SizeAndDuration_SizeTriggersFirst(t *testing.T) {
	// entries 1s apart, size=2, duration=10s → size should trigger first
	entries := makeEntries(4, time.Second)
	result := batch.Apply(entries, batch.Options{Size: 2, Duration: 10 * time.Second})
	if len(result) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(result))
	}
}
