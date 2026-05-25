package timefmt_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timefmt"
)

func TestFormat_RFC3339(t *testing.T) {
	ts := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	opts := timefmt.DefaultOptions()
	got := timefmt.Format(ts, opts)
	want := "2024-06-15T12:30:00Z"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormat_CustomLayout(t *testing.T) {
	ts := time.Date(2024, 1, 2, 9, 5, 3, 0, time.UTC)
	opts := timefmt.Options{Layout: timefmt.LayoutDateTime, Location: time.UTC}
	got := timefmt.Format(ts, opts)
	want := "2024-01-02 09:05:03"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormat_NonUTCLocation(t *testing.T) {
	ts := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	loc, _ := time.LoadLocation("America/New_York")
	opts := timefmt.Options{Layout: timefmt.LayoutDateTime, Location: loc}
	got := timefmt.Format(ts, opts)
	want := "2024-03-09 19:00:00"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestParse_RFC3339(t *testing.T) {
	s := "2024-06-15T12:30:00Z"
	got, err := timefmt.ParseWithDefaults(s)
	if err != nil {
		t.Fatalf("ParseWithDefaults() error: %v", err)
	}
	if got.Year() != 2024 || got.Month() != 6 || got.Day() != 15 {
		t.Errorf("unexpected parsed time: %v", got)
	}
}

func TestParse_DateTime(t *testing.T) {
	s := "2024-01-02 09:05:03"
	got, err := timefmt.ParseWithDefaults(s)
	if err != nil {
		t.Fatalf("ParseWithDefaults() error: %v", err)
	}
	if got.Hour() != 9 || got.Minute() != 5 {
		t.Errorf("unexpected parsed time: %v", got)
	}
}

func TestParse_UnknownFormat_ReturnsError(t *testing.T) {
	_, err := timefmt.ParseWithDefaults("not-a-timestamp")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestLookupLayout_Known(t *testing.T) {
	l, err := timefmt.LookupLayout("rfc3339")
	if err != nil {
		t.Fatalf("LookupLayout() error: %v", err)
	}
	if l.Name != "rfc3339" {
		t.Errorf("unexpected layout name: %q", l.Name)
	}
}

func TestLookupLayout_Unknown(t *testing.T) {
	_, err := timefmt.LookupLayout("nosuchlayout")
	if err == nil {
		t.Error("expected error for unknown layout, got nil")
	}
}

func TestDefaultOptions_UTC(t *testing.T) {
	opts := timefmt.DefaultOptions()
	if opts.Location != time.UTC {
		t.Errorf("expected UTC location, got %v", opts.Location)
	}
	if opts.Layout.Name != "rfc3339" {
		t.Errorf("expected rfc3339 layout, got %q", opts.Layout.Name)
	}
}
