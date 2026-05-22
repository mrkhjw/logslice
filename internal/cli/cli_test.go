package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return f.Name()
}

const sampleLog = `2024-01-15T10:00:00Z INFO  starting server port=8080
2024-01-15T10:01:00Z ERROR connection refused host=db
2024-01-15T10:02:00Z DEBUG heartbeat ok
`

func TestRun_TextFormat(t *testing.T) {
	path := writeTempLog(t, sampleLog)
	if err := Run([]string{"-format", "text", path}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_JSONFormat(t *testing.T) {
	path := writeTempLog(t, sampleLog)
	if err := Run([]string{"-format", "json", path}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_LevelFilter(t *testing.T) {
	path := writeTempLog(t, sampleLog)
	if err := Run([]string{"-level", "ERROR", path}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	path := writeTempLog(t, sampleLog)
	if err := Run([]string{"-format", "xml", path}); err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestRun_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.log")
	if err := Run([]string{path}); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseArgs_Defaults(t *testing.T) {
	cfg, err := parseArgs([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.Format)
	}
	if cfg.Level != "" {
		t.Errorf("expected empty level, got %q", cfg.Level)
	}
}

func TestParseArgs_AllFlags(t *testing.T) {
	cfg, err := parseArgs([]string{
		"-level", "INFO",
		"-since", "2024-01-01T00:00:00Z",
		"-until", "2024-12-31T23:59:59Z",
		"-format", "json",
		"myfile.log",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Level != "INFO" {
		t.Errorf("expected level INFO, got %q", cfg.Level)
	}
	if cfg.FilePath != "myfile.log" {
		t.Errorf("expected filepath myfile.log, got %q", cfg.FilePath)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format json, got %q", cfg.Format)
	}
}
