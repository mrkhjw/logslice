package pivot_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/pivot"
)

func makeEntry(level, msg string, fields map[string]string) parser.Entry {
	if fields == nil {
		fields = map[string]string{}
	}
	return parser.Entry{Timestamp: time.Now(), Level: level, Message: msg, Fields: fields}
}

func TestApply_GroupByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", nil),
		makeEntry("error", "c", nil),
	}
	res := pivot.Apply(entries, pivot.DefaultOptions())
	if len(res.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(res.Rows))
	}
	if res.Rows[0].Key != "info" || res.Rows[0].Count != 2 {
		t.Errorf("expected info=2, got %s=%d", res.Rows[0].Key, res.Rows[0].Count)
	}
}

func TestApply_GroupByCustomField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", map[string]string{"host": "web1"}),
		makeEntry("info", "b", map[string]string{"host": "web1"}),
		makeEntry("info", "c", map[string]string{"host": "web2"}),
	}
	opts := pivot.Options{KeyField: "host"}
	res := pivot.Apply(entries, opts)
	if len(res.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(res.Rows))
	}
	if res.Rows[0].Key != "web1" {
		t.Errorf("expected web1 first, got %s", res.Rows[0].Key)
	}
}

func TestApply_TopN(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("warn", "b", nil),
		makeEntry("error", "c", nil),
	}
	opts := pivot.Options{KeyField: "level", TopN: 2}
	res := pivot.Apply(entries, opts)
	if len(res.Rows) != 2 {
		t.Fatalf("expected 2 rows after TopN, got %d", len(res.Rows))
	}
}

func TestApply_ValueField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", map[string]string{"host": "web1"}),
		makeEntry("info", "b", map[string]string{"host": "web2"}),
	}
	opts := pivot.Options{KeyField: "level", ValueField: "host"}
	res := pivot.Apply(entries, opts)
	if len(res.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(res.Rows))
	}
	if len(res.Rows[0].Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(res.Rows[0].Values))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	res := pivot.Apply(nil, pivot.DefaultOptions())
	if len(res.Rows) != 0 {
		t.Errorf("expected empty rows for nil input")
	}
}

func TestApply_EmptyKeyField_FallsBackToLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("debug", "x", nil),
	}
	opts := pivot.Options{KeyField: ""}
	res := pivot.Apply(entries, opts)
	if res.KeyField != "level" {
		t.Errorf("expected KeyField=level, got %s", res.KeyField)
	}
	if len(res.Rows) != 1 || res.Rows[0].Key != "debug" {
		t.Errorf("unexpected rows: %+v", res.Rows)
	}
}
