// Package enrich provides log entry enrichment by attaching derived or
// external metadata fields to parsed log entries before output.
//
// Enrichers can inject fields such as hostname, environment tags, or
// computed values based on existing entry content.
package enrich

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Enricher applies additional fields to a log entry.
type Enricher interface {
	Enrich(e parser.Entry) parser.Entry
}

// EnricherFunc is a function adapter implementing Enricher.
type EnricherFunc func(e parser.Entry) parser.Entry

// Enrich implements Enricher.
func (f EnricherFunc) Enrich(e parser.Entry) parser.Entry {
	return f(e)
}

// Apply runs all provided enrichers against each entry in sequence.
func Apply(entries []parser.Entry, enrichers ...Enricher) []parser.Entry {
	if len(enrichers) == 0 {
		return entries
	}
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		for _, enc := range enrichers {
			e = enc.Enrich(e)
		}
		out[i] = e
	}
	return out
}

// WithHostname returns an Enricher that attaches the current machine's
// hostname to every entry under the key "host".
func WithHostname() Enricher {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return EnricherFunc(func(e parser.Entry) parser.Entry {
		e = copyEntry(e)
		e.Fields["host"] = host
		return e
	})
}

// WithStaticField returns an Enricher that sets a fixed key/value pair on
// every entry. Existing fields with the same key are overwritten.
func WithStaticField(key, value string) Enricher {
	return EnricherFunc(func(e parser.Entry) parser.Entry {
		e = copyEntry(e)
		e.Fields[key] = value
		return e
	})
}

// WithElapsed returns an Enricher that computes the elapsed duration since
// origin and stores it as a human-readable string under the key "elapsed".
func WithElapsed(origin time.Time) Enricher {
	return EnricherFunc(func(e parser.Entry) parser.Entry {
		e = copyEntry(e)
		d := e.Timestamp.Sub(origin)
		e.Fields["elapsed"] = formatDuration(d)
		return e
	})
}

// WithPrefix returns an Enricher that prepends a string to the entry's
// message field.
func WithPrefix(prefix string) Enricher {
	return EnricherFunc(func(e parser.Entry) parser.Entry {
		e.Message = fmt.Sprintf("%s%s", prefix, e.Message)
		return e
	})
}

// copyEntry returns a shallow copy of e with a duplicated Fields map so
// enrichers do not mutate the original entry.
func copyEntry(e parser.Entry) parser.Entry {
	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	e.Fields = fields
	return e
}

// formatDuration renders a duration as a compact string, e.g. "1h2m3s".
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "-" + formatDuration(-d)
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	var b strings.Builder
	if h > 0 {
		fmt.Fprintf(&b, "%dh", h)
	}
	if m > 0 || h > 0 {
		fmt.Fprintf(&b, "%dm", m)
	}
	fmt.Fprintf(&b, "%ds", s)
	return b.String()
}
