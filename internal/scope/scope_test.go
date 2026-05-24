package scope

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Level:     "INFO",
		Message:   msg,
		Fields:    map[string]string{},
	}
}

var base = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestApply_EmptyInput(t *testing.T) {
	result := Apply(nil, DefaultOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestApply_SingleBucket(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(base, "a"),
		makeEntry(base.Add(30*time.Second), "b"),
	}
	buckets := Apply(entries, DefaultOptions())
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if len(buckets[0].Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(buckets[0].Entries))
	}
}

func TestApply_MultipleBuckets(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(base, "a"),
		makeEntry(base.Add(2*time.Minute), "b"),
		makeEntry(base.Add(2*time.Minute+10*time.Second), "c"),
	}
	buckets := Apply(entries, DefaultOptions())
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if len(buckets[0].Entries) != 1 {
		t.Errorf("bucket 0: expected 1 entry, got %d", len(buckets[0].Entries))
	}
	if len(buckets[1].Entries) != 2 {
		t.Errorf("bucket 1: expected 2 entries, got %d", len(buckets[1].Entries))
	}
}

func TestApply_UnscopedZeroTimestamp(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(time.Time{}, "no-ts"),
		makeEntry(base, "with-ts"),
	}
	buckets := Apply(entries, DefaultOptions())
	names := map[string]int{}
	for _, b := range buckets {
		names[b.Name] = len(b.Entries)
	}
	if names["unscoped"] != 1 {
		t.Errorf("expected 1 unscoped entry, got %d", names["unscoped"])
	}
}

func TestApply_CustomDuration(t *testing.T) {
	opts := Options{Duration: 5 * time.Minute, Label: "2006-01-02T15:04"}
	entries := []parser.Entry{
		makeEntry(base, "a"),
		makeEntry(base.Add(4*time.Minute), "b"),
		makeEntry(base.Add(6*time.Minute), "c"),
	}
	buckets := Apply(entries, opts)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
}

func TestApply_BucketBoundaries(t *testing.T) {
	opts := DefaultOptions()
	entries := []parser.Entry{makeEntry(base, "x")}
	buckets := Apply(entries, opts)
	b := buckets[0]
	if !b.From.Equal(base.Truncate(time.Minute)) {
		t.Errorf("unexpected From: %v", b.From)
	}
	if !b.To.Equal(b.From.Add(time.Minute)) {
		t.Errorf("unexpected To: %v", b.To)
	}
}
