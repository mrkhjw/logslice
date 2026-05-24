package throttle_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/throttle"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func makeEntry(level, msg string, ts time.Time) parser.Entry {
	return parser.Entry{Level: level, Message: msg, Timestamp: ts}
}

func TestApply_NoCooldown_ReturnsAll(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", "hello", base),
		makeEntry("INFO", "hello", base.Add(time.Second)),
	}
	opts := throttle.DefaultOptions()
	opts.Cooldown = 0
	got := throttle.Apply(entries, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_DropsDuplicateWithinCooldown(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", "hello", base),
		makeEntry("INFO", "hello", base.Add(2*time.Second)),
	}
	opts := throttle.Options{Cooldown: 5 * time.Second}
	got := throttle.Apply(entries, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Timestamp != base {
		t.Errorf("unexpected entry kept")
	}
}

func TestApply_AllowsAfterCooldown(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", "hello", base),
		makeEntry("INFO", "hello", base.Add(10*time.Second)),
	}
	opts := throttle.Options{Cooldown: 5 * time.Second}
	got := throttle.Apply(entries, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_DifferentLevels_BothKept(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", "hello", base),
		makeEntry("ERROR", "hello", base.Add(time.Second)),
	}
	opts := throttle.Options{Cooldown: 5 * time.Second}
	got := throttle.Apply(entries, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_CustomKeyFunc(t *testing.T) {
	// Key only on level — any INFO within cooldown is suppressed regardless of msg.
	entries := []parser.Entry{
		makeEntry("INFO", "msg-a", base),
		makeEntry("INFO", "msg-b", base.Add(time.Second)),
		makeEntry("ERROR", "msg-a", base.Add(time.Second)),
	}
	opts := throttle.Options{
		Cooldown: 5 * time.Second,
		KeyFunc: func(e parser.Entry) string { return e.Level },
	}
	got := throttle.Apply(entries, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 (one INFO + one ERROR), got %d", len(got))
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := throttle.Apply(nil, throttle.DefaultOptions())
	if len(got) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := throttle.DefaultOptions()
	if opts.Cooldown != 5*time.Second {
		t.Errorf("expected 5s cooldown, got %v", opts.Cooldown)
	}
}
