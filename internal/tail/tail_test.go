package tail

import (
	"os"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	_, err := f.WriteString(line + "\n")
	if err != nil {
		t.Fatalf("writeLine: %v", err)
	}
}

func TestTail_ReceivesNewEntries(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	p := parser.New()
	out := make(chan parser.Entry, 8)
	stop := make(chan struct{})

	opts := Options{
		PollInterval: 20 * time.Millisecond,
		MaxRetries:   50, // ~1 s total
	}

	go func() {
		_ = Tail(f.Name(), opts, p, out, stop)
	}()

	time.Sleep(30 * time.Millisecond)
	writeLine(t, f, `2024-01-15T10:00:00Z INFO  msg="hello tail" service=api`)
	writeLine(t, f, `2024-01-15T10:00:01Z ERROR msg="something broke" code=500`)

	var entries []parser.Entry
	timeout := time.After(2 * time.Second)
collect:
	for len(entries) < 2 {
		select {
		case e := <-out:
			entries = append(entries, e)
		case <-timeout:
			break collect
		}
	}

	close(stop)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Message != "hello tail" {
		t.Errorf("entry[0] message: got %q, want %q", entries[0].Message, "hello tail")
	}
	if entries[1].Level != "ERROR" {
		t.Errorf("entry[1] level: got %q, want ERROR", entries[1].Level)
	}
}

func TestTail_StopsOnMaxRetries(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-empty-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	p := parser.New()
	out := make(chan parser.Entry, 4)
	stop := make(chan struct{})

	opts := Options{
		PollInterval: 10 * time.Millisecond,
		MaxRetries:   3,
	}

	done := make(chan error, 1)
	go func() {
		done <- Tail(f.Name(), opts, p, out, stop)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Tail did not stop after MaxRetries")
	}
	close(stop)
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.PollInterval != 250*time.Millisecond {
		t.Errorf("PollInterval: got %v, want 250ms", opts.PollInterval)
	}
	if opts.MaxRetries != 0 {
		t.Errorf("MaxRetries: got %d, want 0", opts.MaxRetries)
	}
}
