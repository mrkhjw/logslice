package headtail_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/headtail"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	base := time.Now()
	for i := 0; i < n; i++ {
		entries[i] = parser.Entry{
			Timestamp: base.Add(time.Duration(i) * time.Second),
			Level:     "INFO",
			Message:   "msg",
			Fields:    map[string]string{"idx": string(rune('0' + i))},
		}
	}
	return entries
}

func TestHead_ZeroN_ReturnsAll(t *testing.T) {
	entries := makeEntries(5)
	got := headtail.Head(entries, headtail.Options{N: 0})
	if len(got) != 5 {
		t.Fatalf("expected 5, got %d", len(got))
	}
}

func TestHead_KeepsFirstN(t *testing.T) {
	entries := makeEntries(6)
	got := headtail.Head(entries, headtail.Options{N: 3})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	for i, e := range got {
		if e.Fields["idx"] != entries[i].Fields["idx"] {
			t.Errorf("entry %d mismatch", i)
		}
	}
}

func TestHead_NGreaterThanLen_ReturnsAll(t *testing.T) {
	entries := makeEntries(3)
	got := headtail.Head(entries, headtail.Options{N: 10})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestHead_EmptyInput(t *testing.T) {
	got := headtail.Head(nil, headtail.Options{N: 5})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestTail_ZeroN_ReturnsAll(t *testing.T) {
	entries := makeEntries(5)
	got := headtail.Tail(entries, headtail.Options{N: 0})
	if len(got) != 5 {
		t.Fatalf("expected 5, got %d", len(got))
	}
}

func TestTail_KeepsLastN(t *testing.T) {
	entries := makeEntries(6)
	got := headtail.Tail(entries, headtail.Options{N: 2})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	if got[0].Fields["idx"] != entries[4].Fields["idx"] {
		t.Errorf("first tail entry mismatch")
	}
	if got[1].Fields["idx"] != entries[5].Fields["idx"] {
		t.Errorf("second tail entry mismatch")
	}
}

func TestTail_NGreaterThanLen_ReturnsAll(t *testing.T) {
	entries := makeEntries(3)
	got := headtail.Tail(entries, headtail.Options{N: 10})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestTail_EmptyInput(t *testing.T) {
	got := headtail.Tail(nil, headtail.Options{N: 5})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := headtail.DefaultOptions()
	if opts.N != 0 {
		t.Errorf("expected N=0, got %d", opts.N)
	}
}
