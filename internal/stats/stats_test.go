package stats_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/stats"
)

func makeEntry(level parser.Level, ts time.Time, fields map[string]string) parser.Entry {
	return parser.Entry{
		Level:     level,
		Timestamp: ts,
		Message:   "test message",
		Fields:    fields,
	}
}

func TestCompute_Total(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, time.Now(), nil),
		makeEntry(parser.LevelError, time.Now(), nil),
	}
	s := stats.Compute(entries)
	if s.Total != 2 {
		t.Errorf("expected Total=2, got %d", s.Total)
	}
}

func TestCompute_ByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, time.Now(), nil),
		makeEntry(parser.LevelInfo, time.Now(), nil),
		makeEntry(parser.LevelError, time.Now(), nil),
	}
	s := stats.Compute(entries)
	if s.ByLevel["INFO"] != 2 {
		t.Errorf("expected INFO=2, got %d", s.ByLevel["INFO"])
	}
	if s.ByLevel["ERROR"] != 1 {
		t.Errorf("expected ERROR=1, got %d", s.ByLevel["ERROR"])
	}
}

func TestCompute_TimeRange(t *testing.T) {
	early := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	late := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, late, nil),
		makeEntry(parser.LevelInfo, early, nil),
	}
	s := stats.Compute(entries)
	if s.FirstTime == nil || !s.FirstTime.Equal(early) {
		t.Errorf("expected FirstTime=%v, got %v", early, s.FirstTime)
	}
	if s.LastTime == nil || !s.LastTime.Equal(late) {
		t.Errorf("expected LastTime=%v, got %v", late, s.LastTime)
	}
}

func TestCompute_Duration(t *testing.T) {
	early := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	late := early.Add(2 * time.Hour)
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, early, nil),
		makeEntry(parser.LevelInfo, late, nil),
	}
	s := stats.Compute(entries)
	if s.Duration() != 2*time.Hour {
		t.Errorf("expected 2h duration, got %v", s.Duration())
	}
}

func TestCompute_Fields(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, time.Now(), map[string]string{"service": "api"}),
		makeEntry(parser.LevelInfo, time.Now(), map[string]string{"service": "api"}),
		makeEntry(parser.LevelInfo, time.Now(), map[string]string{"service": "worker"}),
	}
	s := stats.Compute(entries)
	if s.Fields["service"]["api"] != 2 {
		t.Errorf("expected service.api=2, got %d", s.Fields["service"]["api"])
	}
	if s.Fields["service"]["worker"] != 1 {
		t.Errorf("expected service.worker=1, got %d", s.Fields["service"]["worker"])
	}
}

func TestCompute_Empty(t *testing.T) {
	s := stats.Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if s.Duration() != 0 {
		t.Errorf("expected zero duration, got %v", s.Duration())
	}
}
