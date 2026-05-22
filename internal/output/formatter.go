package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/user/logslice/internal/parser"
)

// Format represents the output format type.
type Format string

const (
	FormatJSON  Format = "json"
	FormatText  Format = "text"
	FormatTable Format = "table"
)

// Formatter writes log entries to an output stream.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter with the given format and writer.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{format: format, w: w}
}

// Write outputs a slice of log entries in the configured format.
func (f *Formatter) Write(entries []parser.Entry) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(entries)
	case FormatTable:
		return f.writeTable(entries)
	default:
		return f.writeText(entries)
	}
}

func (f *Formatter) writeJSON(entries []parser.Entry) error {
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("json encode: %w", err)
		}
	}
	return nil
}

func (f *Formatter) writeText(entries []parser.Entry) error {
	for _, e := range entries {
		fields := make([]string, 0, len(e.Fields))
		for k, v := range e.Fields {
			fields = append(fields, fmt.Sprintf("%s=%s", k, v))
		}
		line := fmt.Sprintf("%s [%s] %s", e.Timestamp.Format("2006-01-02T15:04:05Z07:00"), e.Level, e.Message)
		if len(fields) > 0 {
			line += " " + strings.Join(fields, " ")
		}
		fmt.Fprintln(f.w, line)
	}
	return nil
}

func (f *Formatter) writeTable(entries []parser.Entry) error {
	tw := tabwriter.NewWriter(f.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tLEVEL\tMESSAGE")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			e.Level,
			e.Message,
		)
	}
	return tw.Flush()
}
