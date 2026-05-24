package replay_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/replay"
)

func makeEntries(timestamps []time.Time) []parser.Entry {
	entries := make([]parser.Entry, len(timestamps))
	for i, ts := range timestamps {
		entries[i] = parser.Entry{
			Timestamp: ts,
			Level:     "INFO",
			Message:   "test message",
		}
	}
	return entries
}

func TestApply_EmptyInput(t *testing.T) {
	ch := replay.Apply(nil, replay.DefaultOptions())
	var got []parser.Entry
	for e := range ch {
		got = append(got, e)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestApply_AllEntriesDelivered(t *testing.T) {
	now := time.Now()
	timestamps := []time.Time{now, now.Add(10 * time.Millisecond), now.Add(20 * time.Millisecond)}
	entries := makeEntries(timestamps)

	opts := replay.Options{Speed: 100.0, MaxDelay: time.Second}
	ch := replay.Apply(entries, opts)

	var got []parser.Entry
	for e := range ch {
		got = append(got, e)
	}
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestApply_OrderPreserved(t *testing.T) {
	now := time.Now()
	timestamps := []time.Time{now, now.Add(5 * time.Millisecond), now.Add(10 * time.Millisecond)}
	entries := makeEntries(timestamps)

	opts := replay.Options{Speed: 1000.0, MaxDelay: time.Second}
	ch := replay.Apply(entries, opts)

	i := 0
	for e := range ch {
		if !e.Timestamp.Equal(entries[i].Timestamp) {
			t.Errorf("entry %d: expected %v, got %v", i, entries[i].Timestamp, e.Timestamp)
		}
		i++
	}
}

func TestApply_ZeroSpeedDefaultsToOne(t *testing.T) {
	now := time.Now()
	entries := makeEntries([]time.Time{now, now.Add(time.Millisecond)})

	// Should not panic or hang; MaxDelay keeps it fast.
	opts := replay.Options{Speed: 0, MaxDelay: 5 * time.Millisecond}
	ch := replay.Apply(entries, opts)
	var got []parser.Entry
	for e := range ch {
		got = append(got, e)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := replay.DefaultOptions()
	if opts.Speed != 1.0 {
		t.Errorf("expected Speed 1.0, got %f", opts.Speed)
	}
	if opts.MaxDelay != 5*time.Second {
		t.Errorf("expected MaxDelay 5s, got %v", opts.MaxDelay)
	}
}
