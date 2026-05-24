package window_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/window"
)

func makeEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{Timestamp: ts, Message: msg, Level: "info"}
}

func TestApply_EmptyInput(t *testing.T) {
	result := window.Apply(nil, window.DefaultOptions())
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestApply_ZeroSize_ReturnsAll(t *testing.T) {
	now := time.Now()
	entries := []parser.Entry{
		makeEntry(now.Add(-10*time.Minute), "old"),
		makeEntry(now, "new"),
	}
	result := window.Apply(entries, window.Options{Size: 0})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestApply_Rolling_DropsOldEntries(t *testing.T) {
	now := time.Now()
	entries := []parser.Entry{
		makeEntry(now.Add(-10*time.Minute), "too old"),
		makeEntry(now.Add(-3*time.Minute), "in window"),
		makeEntry(now, "latest"),
	}
	opts := window.Options{Size: 5 * time.Minute}
	result := window.Apply(entries, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Message != "in window" {
		t.Errorf("unexpected first entry: %s", result[0].Message)
	}
}

func TestApply_Anchored_KeepsFromFirstTimestamp(t *testing.T) {
	now := time.Now()
	entries := []parser.Entry{
		makeEntry(now, "first"),
		makeEntry(now.Add(3*time.Minute), "within"),
		makeEntry(now.Add(7*time.Minute), "outside"),
	}
	opts := window.Options{Size: 5 * time.Minute, Anchor: true}
	result := window.Apply(entries, opts)
	// anchor=now, cutoff=now-5m; all entries are >= now so all kept
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_ZeroTimestamp_AlwaysKept(t *testing.T) {
	now := time.Now()
	entries := []parser.Entry{
		{Message: "no timestamp"},
		makeEntry(now.Add(-10*time.Minute), "old"),
		makeEntry(now, "latest"),
	}
	opts := window.Options{Size: 2 * time.Minute}
	result := window.Apply(entries, opts)
	// zero-ts entry + latest kept; old dropped
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Message != "no timestamp" {
		t.Errorf("expected zero-timestamp entry first, got %q", result[0].Message)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := window.DefaultOptions()
	if opts.Size != 5*time.Minute {
		t.Errorf("expected 5m, got %v", opts.Size)
	}
	if opts.Anchor {
		t.Error("expected Anchor=false by default")
	}
}
