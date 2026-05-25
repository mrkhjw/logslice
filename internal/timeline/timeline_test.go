package timeline

import (
	"bytes"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{Timestamp: ts, Message: msg, Level: "info"}
}

var base = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestApply_EmptyInput(t *testing.T) {
	buckets := Apply(nil, DefaultOptions())
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(buckets))
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
		t.Errorf("expected 2 entries in bucket, got %d", len(buckets[0].Entries))
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

func TestApply_SkipsZeroTimestamp(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(base, "valid"),
		{Message: "no timestamp"},
	}
	buckets := Apply(entries, DefaultOptions())
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if len(buckets[0].Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(buckets[0].Entries))
	}
}

func TestApply_CustomResolution(t *testing.T) {
	opts := Options{Resolution: time.Hour}
	entries := []parser.Entry{
		makeEntry(base, "a"),
		makeEntry(base.Add(30*time.Minute), "b"),
		makeEntry(base.Add(61*time.Minute), "c"),
	}
	buckets := Apply(entries, opts)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
}

func TestPrintSummary_ContainsHeader(t *testing.T) {
	buckets := Apply([]parser.Entry{makeEntry(base, "x")}, DefaultOptions())
	var buf bytes.Buffer
	PrintSummary(&buf, buckets, "")
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}
	for _, want := range []string{"bucket", "count", "bar"} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestPrintSummary_EmptyBuckets(t *testing.T) {
	var buf bytes.Buffer
	PrintSummary(&buf, nil, "")
	if !bytes.Contains(buf.Bytes(), []byte("no entries")) {
		t.Errorf("expected 'no entries' message, got: %s", buf.String())
	}
}
