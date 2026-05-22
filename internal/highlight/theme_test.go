package highlight

import (
	"testing"
)

func TestSetTheme_Custom(t *testing.T) {
	custom := Theme{
		Error: "[ERR]",
		Warn:  "[WRN]",
		Info:  "[INF]",
		Debug: "[DBG]",
		Fatal: "[FTL]",
		Bold:  "[B]",
		Reset: "[R]",
	}
	SetTheme(custom)
	defer EnableColor()

	if ActiveTheme.Error != "[ERR]" {
		t.Errorf("expected Error=[ERR], got %q", ActiveTheme.Error)
	}
	if ActiveTheme.Info != "[INF]" {
		t.Errorf("expected Info=[INF], got %q", ActiveTheme.Info)
	}
}

func TestDisableColor(t *testing.T) {
	DisableColor()
	defer EnableColor()

	if ActiveTheme.Error != "" {
		t.Errorf("expected empty Error after DisableColor, got %q", ActiveTheme.Error)
	}
	if ActiveTheme.Bold != "" {
		t.Errorf("expected empty Bold after DisableColor, got %q", ActiveTheme.Bold)
	}
}

func TestEnableColor(t *testing.T) {
	DisableColor()
	EnableColor()

	if ActiveTheme.Error == "" {
		t.Error("expected non-empty Error after EnableColor")
	}
	if ActiveTheme.Reset == "" {
		t.Error("expected non-empty Reset after EnableColor")
	}
}

func TestNoColorTheme_AllEmpty(t *testing.T) {
	fields := []struct {
		name  string
		value string
	}{
		{"Error", NoColorTheme.Error},
		{"Warn", NoColorTheme.Warn},
		{"Info", NoColorTheme.Info},
		{"Debug", NoColorTheme.Debug},
		{"Fatal", NoColorTheme.Fatal},
		{"Bold", NoColorTheme.Bold},
		{"Reset", NoColorTheme.Reset},
	}
	for _, f := range fields {
		if f.value != "" {
			t.Errorf("NoColorTheme.%s should be empty, got %q", f.name, f.value)
		}
	}
}

func TestDefaultTheme_NonEmpty(t *testing.T) {
	if DefaultTheme.Error == "" {
		t.Error("DefaultTheme.Error should not be empty")
	}
	if DefaultTheme.Reset == "" {
		t.Error("DefaultTheme.Reset should not be empty")
	}
}
