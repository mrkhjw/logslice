package rotate_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/rotate"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTmp: %v", err)
	}
	return p
}

func TestNew_ValidFile(t *testing.T) {
	p := writeTmp(t, "hello\n")
	d, err := rotate.New(p, rotate.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := rotate.New("/nonexistent/path/file.log", rotate.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRotated_NoChange(t *testing.T) {
	p := writeTmp(t, "line1\n")
	d, _ := rotate.New(p, rotate.DefaultOptions())
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Error("expected file not rotated after no change")
	}
}

func TestRotated_FileTruncated(t *testing.T) {
	p := writeTmp(t, "line1\nline2\nline3\n")
	d, _ := rotate.New(p, rotate.DefaultOptions())

	// Truncate the file
	if err := os.WriteFile(p, []byte(""), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Error("expected rotation detected after truncation")
	}
}

func TestRotated_FileDeleted(t *testing.T) {
	p := writeTmp(t, "data\n")
	d, _ := rotate.New(p, rotate.DefaultOptions())

	if err := os.Remove(p); err != nil {
		t.Fatalf("remove: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Error("expected rotation detected after file deletion")
	}
}

func TestReset_UpdatesBaseline(t *testing.T) {
	p := writeTmp(t, "line1\n")
	d, _ := rotate.New(p, rotate.DefaultOptions())

	// Truncate then reset
	if err := os.WriteFile(p, []byte(""), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	if err := d.Reset(); err != nil {
		t.Fatalf("reset error: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Error("expected no rotation after reset")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := rotate.DefaultOptions()
	if opts.PollInterval != 500*time.Millisecond {
		t.Errorf("unexpected poll interval: %v", opts.PollInterval)
	}
	if opts.MaxReopen != 5 {
		t.Errorf("unexpected max reopen: %d", opts.MaxReopen)
	}
}
