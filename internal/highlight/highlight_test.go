package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestColorize_ContainsText(t *testing.T) {
	result := highlight.Colorize("hello", highlight.Red)
	if !strings.Contains(result, "hello") {
		t.Errorf("expected result to contain original text, got: %q", result)
	}
}

func TestColorize_ContainsEscapeCode(t *testing.T) {
	result := highlight.Colorize("hello", highlight.Red)
	if !strings.Contains(result, "\033[") {
		t.Errorf("expected ANSI escape code in result, got: %q", result)
	}
}

func TestLevelColor_Error(t *testing.T) {
	if highlight.LevelColor("ERROR") != highlight.Red {
		t.Error("expected ERROR to map to Red")
	}
}

func TestLevelColor_Fatal(t *testing.T) {
	if highlight.LevelColor("FATAL") != highlight.Red {
		t.Error("expected FATAL to map to Red")
	}
}

func TestLevelColor_Warn(t *testing.T) {
	if highlight.LevelColor("WARN") != highlight.Yellow {
		t.Error("expected WARN to map to Yellow")
	}
}

func TestLevelColor_Info(t *testing.T) {
	if highlight.LevelColor("INFO") != highlight.Cyan {
		t.Error("expected INFO to map to Cyan")
	}
}

func TestLevelColor_Unknown(t *testing.T) {
	if highlight.LevelColor("DEBUG") != highlight.White {
		t.Error("expected unknown level to map to White")
	}
}

func TestForLevel_ContainsLevel(t *testing.T) {
	levels := []string{"INFO", "WARN", "ERROR", "FATAL", "DEBUG"}
	for _, lvl := range levels {
		result := highlight.ForLevel(lvl)
		if !strings.Contains(result, lvl) {
			t.Errorf("ForLevel(%q) = %q, want to contain level string", lvl, result)
		}
	}
}

func TestBoldText_ContainsText(t *testing.T) {
	result := highlight.BoldText("timestamp")
	if !strings.Contains(result, "timestamp") {
		t.Errorf("BoldText result missing original text: %q", result)
	}
}
