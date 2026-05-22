package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(level parser.Level, ts time.Time, fields map[string]string) parser.Entry {
	return parser.Entry{
		Level:     level,
		Timestamp: ts,
		Message:   "test message",
		Fields:    fields,
	}
}

var (
	t1 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC)
)

func TestFilter_ByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, t2, nil),
		makeEntry(parser.LevelError, t2, nil),
		makeEntry(parser.LevelWarn, t2, nil),
	}
	got := filter.Filter(entries, filter.Options{Level: "ERROR"})
	if len(got) != 1 || got[0].Level != parser.LevelError {
		t.Errorf("expected 1 ERROR entry, got %d", len(got))
	}
}

func TestFilter_BySince(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, t1, nil),
		makeEntry(parser.LevelInfo, t2, nil),
		makeEntry(parser.LevelInfo, t3, nil),
	}
	got := filter.Filter(entries, filter.Options{Since: t2})
	if len(got) != 2 {
		t.Errorf("expected 2 entries since t2, got %d", len(got))
	}
}

func TestFilter_ByUntil(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, t1, nil),
		makeEntry(parser.LevelInfo, t2, nil),
		makeEntry(parser.LevelInfo, t3, nil),
	}
	got := filter.Filter(entries, filter.Options{Until: t2})
	if len(got) != 2 {
		t.Errorf("expected 2 entries until t2, got %d", len(got))
	}
}

func TestFilter_ByTimeRange(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, t1, nil),
		makeEntry(parser.LevelInfo, t2, nil),
		makeEntry(parser.LevelInfo, t3, nil),
	}
	got := filter.Filter(entries, filter.Options{Since: t2, Until: t2})
	if len(got) != 1 {
		t.Errorf("expected 1 entry in exact range, got %d", len(got))
	}
}

func TestFilter_ByField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, t1, map[string]string{"service": "api"}),
		makeEntry(parser.LevelInfo, t1, map[string]string{"service": "worker"}),
		makeEntry(parser.LevelInfo, t1, nil),
	}
	got := filter.Filter(entries, filter.Options{FieldKey: "service", FieldVal: "api"})
	if len(got) != 1 || got[0].Fields["service"] != "api" {
		t.Errorf("expected 1 entry with service=api, got %d", len(got))
	}
}

func TestFilter_NoOptions(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, t1, nil),
		makeEntry(parser.LevelError, t2, nil),
	}
	got := filter.Filter(entries, filter.Options{})
	if len(got) != len(entries) {
		t.Errorf("expected all %d entries, got %d", len(entries), len(got))
	}
}
