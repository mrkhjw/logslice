package sampler_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/sampler"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Message:   "test message",
		}
	}
	return entries
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	s := sampler.New(sampler.Options{}, 42)
	entries := makeEntries(10)
	result := s.Apply(entries)
	if len(result) != 10 {
		t.Errorf("expected 10 entries, got %d", len(result))
	}
}

func TestApply_EveryN_KeepsCorrectEntries(t *testing.T) {
	s := sampler.New(sampler.Options{EveryN: 3}, 42)
	entries := makeEntries(9)
	result := s.Apply(entries)
	// entries 3, 6, 9 should be kept
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_EveryN_One_KeepsAll(t *testing.T) {
	s := sampler.New(sampler.Options{EveryN: 1}, 42)
	entries := makeEntries(5)
	result := s.Apply(entries)
	if len(result) != 5 {
		t.Errorf("expected 5 entries, got %d", len(result))
	}
}

func TestApply_Percent100_KeepsAll(t *testing.T) {
	s := sampler.New(sampler.Options{Percent: 100}, 42)
	entries := makeEntries(20)
	result := s.Apply(entries)
	if len(result) != 20 {
		t.Errorf("expected 20 entries, got %d", len(result))
	}
}

func TestApply_Percent0_NoFiltering(t *testing.T) {
	s := sampler.New(sampler.Options{Percent: 0}, 42)
	entries := makeEntries(10)
	result := s.Apply(entries)
	if len(result) != 10 {
		t.Errorf("expected 10 entries, got %d", len(result))
	}
}

func TestApply_Percent_ReducesEntries(t *testing.T) {
	s := sampler.New(sampler.Options{Percent: 50}, 99)
	entries := makeEntries(1000)
	result := s.Apply(entries)
	// With 50% we expect roughly half; allow wide tolerance.
	if len(result) < 300 || len(result) > 700 {
		t.Errorf("expected ~500 entries with 50%% sample, got %d", len(result))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	s := sampler.New(sampler.Options{EveryN: 2}, 42)
	result := s.Apply([]parser.Entry{})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}
