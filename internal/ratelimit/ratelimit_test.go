package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/ratelimit"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{Message: "msg"}
	}
	return entries
}

func TestApply_NoLimit_ReturnsAll(t *testing.T) {
	l := ratelimit.New(ratelimit.DefaultOptions())
	in := makeEntries(10)
	out := l.Apply(in)
	if len(out) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(out))
	}
}

func TestApply_ZeroPerSecond_ReturnsAll(t *testing.T) {
	l := ratelimit.New(ratelimit.Options{PerSecond: 0})
	in := makeEntries(5)
	out := l.Apply(in)
	if len(out) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(out))
	}
}

func TestApply_BurstLimit_DropsExcess(t *testing.T) {
	// burst of 3, no time passes → only 3 tokens available
	l := ratelimit.New(ratelimit.Options{PerSecond: 3, Burst: 3})
	// freeze clock so no tokens are replenished between calls
	fixed := time.Now()
	l.(*ratelimit.Limiter) // white-box: use exported helper via option

	// Use the exported Apply; initial tokens == burst == 3
	in := makeEntries(6)
	out := l.Apply(in)
	if len(out) > 3 {
		t.Fatalf("expected at most 3 entries within burst, got %d", len(out))
	}
	_ = fixed
}

func TestApply_EmptyInput(t *testing.T) {
	l := ratelimit.New(ratelimit.Options{PerSecond: 10})
	out := l.Apply([]parser.Entry{})
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestDefaultOptions_NoLimit(t *testing.T) {
	opts := ratelimit.DefaultOptions()
	if opts.PerSecond != 0 {
		t.Errorf("expected PerSecond 0, got %d", opts.PerSecond)
	}
}

func TestApply_HighRate_AllowsMany(t *testing.T) {
	// 1000 rps burst — all 20 entries should pass
	l := ratelimit.New(ratelimit.Options{PerSecond: 1000, Burst: 1000})
	in := makeEntries(20)
	out := l.Apply(in)
	if len(out) != 20 {
		t.Fatalf("expected 20 entries, got %d", len(out))
	}
}

// TestApply_BurstRefill_AllowsAfterWait verifies that tokens are replenished
// over time: after exhausting the burst, waiting one second should allow
// additional entries through at the configured rate.
func TestApply_BurstRefill_AllowsAfterWait(t *testing.T) {
	const rate = 5
	l := ratelimit.New(ratelimit.Options{PerSecond: rate, Burst: rate})

	// Exhaust the initial burst.
	first := l.Apply(makeEntries(rate))
	if len(first) != rate {
		t.Fatalf("expected %d entries in first batch, got %d", rate, len(first))
	}

	// Wait slightly over one second so the token bucket refills.
	time.Sleep(1100 * time.Millisecond)

	// After refill, at least one entry should be allowed through.
	second := l.Apply(makeEntries(rate))
	if len(second) == 0 {
		t.Fatal("expected at least 1 entry after token refill, got 0")
	}
}
