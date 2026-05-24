package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempStore(t *testing.T) (*checkpoint.Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "checkpoint.json")
	s, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s, path
}

func TestGet_DefaultsToZero(t *testing.T) {
	s, _ := tempStore(t)
	if got := s.Get("/var/log/app.log"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestSet_PersistsOffset(t *testing.T) {
	s, path := tempStore(t)
	if err := s.Set("/var/log/app.log", 1024); err != nil {
		t.Fatalf("Set: %v", err)
	}
	// Re-open from same file to verify persistence.
	s2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("New (reload): %v", err)
	}
	if got := s2.Get("/var/log/app.log"); got != 1024 {
		t.Errorf("expected 1024, got %d", got)
	}
}

func TestSet_MultipleFiles(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Set("/a.log", 100)
	_ = s.Set("/b.log", 200)
	if got := s.Get("/a.log"); got != 100 {
		t.Errorf("a.log: expected 100, got %d", got)
	}
	if got := s.Get("/b.log"); got != 200 {
		t.Errorf("b.log: expected 200, got %d", got)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s, path := tempStore(t)
	_ = s.Set("/var/log/app.log", 512)
	if err := s.Delete("/var/log/app.log"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got := s.Get("/var/log/app.log"); got != 0 {
		t.Errorf("expected 0 after delete, got %d", got)
	}
	// Reload and confirm deletion persisted.
	s2, _ := checkpoint.New(path)
	if got := s2.Get("/var/log/app.log"); got != 0 {
		t.Errorf("reload: expected 0 after delete, got %d", got)
	}
}

func TestNew_MissingFile_ReturnsEmptyStore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does_not_exist.json")
	s, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if got := s.Get("/any.log"); got != 0 {
		t.Errorf("expected 0 for empty store, got %d", got)
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)
	_, err := checkpoint.New(path)
	if err == nil {
		t.Error("expected error for corrupt checkpoint file")
	}
}
