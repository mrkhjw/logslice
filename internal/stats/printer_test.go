package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/stats"
)

func TestPrint_ContainsTotals(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, time.Now(), nil),
		makeEntry(parser.LevelError, time.Now(), nil),
	}
	s := stats.Compute(entries)
	var buf bytes.Buffer
	stats.Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "Total entries:") {
		t.Error("expected 'Total entries:' in output")
	}
	if !strings.Contains(out, "2") {
		t.Error("expected count '2' in output")
	}
}

func TestPrint_ContainsLevels(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelWarn, time.Now(), nil),
		makeEntry(parser.LevelWarn, time.Now(), nil),
		makeEntry(parser.LevelError, time.Now(), nil),
	}
	s := stats.Compute(entries)
	var buf bytes.Buffer
	stats.Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strings.Contains(out, "ERROR") {
		t.Error("expected ERROR in output")
	}
	if !strings.Contains(out, "66.7%") {
		t.Errorf("expected 66.7%% in output, got:\n%s", out)
	}
}

func TestPrint_EmptySummary(t *testing.T) {
	s := stats.Compute(nil)
	var buf bytes.Buffer
	stats.Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "0") {
		t.Error("expected '0' in empty summary output")
	}
}

func TestPrint_TimeRange(t *testing.T) {
	early := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	late := early.Add(30 * time.Minute)
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, early, nil),
		makeEntry(parser.LevelInfo, late, nil),
	}
	s := stats.Compute(entries)
	var buf bytes.Buffer
	stats.Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "2024-03-15") {
		t.Errorf("expected date in output, got:\n%s", out)
	}
	if !strings.Contains(out, "30m") {
		t.Errorf("expected duration '30m' in output, got:\n%s", out)
	}
}
