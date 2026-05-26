package offset

import (
	"bytes"
	"os"
	"testing"
)

func TestNew_DefaultsToZero(t *testing.T) {
	tr := New("/tmp/test.log")
	if tr.Get() != 0 {
		t.Fatalf("expected 0, got %d", tr.Get())
	}
}

func TestSet_StoresOffset(t *testing.T) {
	tr := New("/tmp/test.log")
	if err := tr.Set(128); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Get() != 128 {
		t.Fatalf("expected 128, got %d", tr.Get())
	}
}

func TestSet_NegativeOffset_ReturnsError(t *testing.T) {
	tr := New("/tmp/test.log")
	if err := tr.Set(-1); err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestReset_SetsOffsetToZero(t *testing.T) {
	tr := New("/tmp/test.log")
	_ = tr.Set(99)
	tr.Reset()
	if tr.Get() != 0 {
		t.Fatalf("expected 0 after reset, got %d", tr.Get())
	}
}

func TestSeek_PositionsReader(t *testing.T) {
	data := []byte("hello world")
	rs := bytes.NewReader(data)

	tr := New("")
	_ = tr.Set(6)

	if err := tr.Seek(rs); err != nil {
		t.Fatalf("seek error: %v", err)
	}

	buf := make([]byte, 5)
	if _, err := rs.Read(buf); err != nil {
		t.Fatalf("read error: %v", err)
	}
	if string(buf) != "world" {
		t.Fatalf("expected 'world', got %q", string(buf))
	}
}

func TestSync_SetsOffsetToFileSize(t *testing.T) {
	f, err := os.CreateTemp("", "offset_test_*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer os.Remove(f.Name())

	content := []byte("some log content")
	if _, err := f.Write(content); err != nil {
		t.Fatalf("write: %v", err)
	}
	f.Close()

	tr := New(f.Name())
	if err := tr.Sync(); err != nil {
		t.Fatalf("sync error: %v", err)
	}
	if tr.Get() != int64(len(content)) {
		t.Fatalf("expected %d, got %d", len(content), tr.Get())
	}
}

func TestSync_MissingFile_ReturnsError(t *testing.T) {
	tr := New("/nonexistent/path/file.log")
	if err := tr.Sync(); err == nil {
		t.Fatal("expected error for missing file")
	}
}
