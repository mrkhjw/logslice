package merge

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, msg string) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Message:   msg,
		Level:     "INFO",
		Fields:    map[string]string{},
	}
}

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestApply_EmptySources(t *testing.T) {
	result := Apply(DefaultOptions())
	if result != nil {
		t.Errorf("expected nil for empty sources, got %v", result)
	}
}

func TestApply_SingleSource(t *testing.T) {
	src := []parser.Entry{
		makeEntry(base.Add(2*time.Second), "b"),
		makeEntry(base.Add(1*time.Second), "a"),
	}
	out := Apply(DefaultOptions(), src)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Message != "a" || out[1].Message != "b" {
		t.Errorf("unexpected order: %s, %s", out[0].Message, out[1].Message)
	}
}

func TestApply_MultipleSources_Interleaved(t *testing.T) {
	src1 := []parser.Entry{
		makeEntry(base.Add(1*time.Second), "s1-1"),
		makeEntry(base.Add(3*time.Second), "s1-3"),
	}
	src2 := []parser.Entry{
		makeEntry(base.Add(2*time.Second), "s2-2"),
		makeEntry(base.Add(4*time.Second), "s2-4"),
	}
	out := Apply(DefaultOptions(), src1, src2)
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	want := []string{"s1-1", "s2-2", "s1-3", "s2-4"}
	for i, w := range want {
		if out[i].Message != w {
			t.Errorf("index %d: want %q, got %q", i, w, out[i].Message)
		}
	}
}

func TestApply_StableOrder_SameTimestamp(t *testing.T) {
	ts := base
	src1 := []parser.Entry{makeEntry(ts, "first")}
	src2 := []parser.Entry{makeEntry(ts, "second")}
	out := Apply(Options{Stable: true}, src1, src2)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Message != "first" {
		t.Errorf("stable: expected 'first' first, got %q", out[0].Message)
	}
}

func TestApply_EmptySliceAmongSources(t *testing.T) {
	src1 := []parser.Entry{makeEntry(base, "only")}
	src2 := []parser.Entry{}
	out := Apply(DefaultOptions(), src1, src2)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}
