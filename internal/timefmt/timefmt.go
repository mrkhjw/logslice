// Package timefmt provides utilities for formatting and parsing timestamps
// in log entries using configurable time zone and layout options.
package timefmt

import (
	"fmt"
	"time"
)

// Layout represents a named time format layout.
type Layout struct {
	Name   string
	Format string
}

// Predefined layouts for common log timestamp formats.
var (
	LayoutRFC3339     = Layout{Name: "rfc3339", Format: time.RFC3339}
	LayoutRFC3339Nano = Layout{Name: "rfc3339nano", Format: time.RFC3339Nano}
	LayoutDateTime    = Layout{Name: "datetime", Format: "2006-01-02 15:04:05"}
	LayoutDateOnly    = Layout{Name: "date", Format: "2006-01-02"}
	LayoutUnixDate    = Layout{Name: "unix", Format: time.UnixDate}
)

// KnownLayouts lists all built-in layouts.
var KnownLayouts = []Layout{
	LayoutRFC3339,
	LayoutRFC3339Nano,
	LayoutDateTime,
	LayoutDateOnly,
	LayoutUnixDate,
}

// Options controls formatting behaviour.
type Options struct {
	Layout   Layout
	Location *time.Location
}

// DefaultOptions returns Options with RFC3339 layout and UTC timezone.
func DefaultOptions() Options {
	return Options{
		Layout:   LayoutRFC3339,
		Location: time.UTC,
	}
}

// Format formats t using the options.
func Format(t time.Time, opts Options) string {
	loc := opts.Location
	if loc == nil {
		loc = time.UTC
	}
	return t.In(loc).Format(opts.Layout.Format)
}

// Parse parses s trying each of the provided layouts in order.
// It returns the first successful parse result.
func Parse(s string, layouts []Layout, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.UTC
	}
	for _, l := range layouts {
		t, err := time.ParseInLocation(l.Format, s, loc)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timefmt: cannot parse %q with any known layout", s)
}

// ParseWithDefaults parses s using all KnownLayouts and UTC.
func ParseWithDefaults(s string) (time.Time, error) {
	return Parse(s, KnownLayouts, time.UTC)
}

// LookupLayout returns the named layout, or an error if not found.
func LookupLayout(name string) (Layout, error) {
	for _, l := range KnownLayouts {
		if l.Name == name {
			return l, nil
		}
	}
	return Layout{}, fmt.Errorf("timefmt: unknown layout %q", name)
}
