package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/logslice/internal/output"
)

// args holds all parsed command-line arguments.
type args struct {
	// Input source
	filePath string

	// Filtering options
	level  string
	since  string
	until  string
	fields []fieldArg

	// Output options
	format  string
	noColor bool
	outFile string

	// Stats
	showStats bool

	// Misc
	verbose bool
}

// fieldArg represents a key=value filter supplied via --field flag.
type fieldArg struct {
	Key   string
	Value string
}

// multiFlag allows a flag to be specified multiple times.
type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ", ")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// newFlagSet builds and returns a configured *flag.FlagSet along with pointers
// to the underlying values. The returned FlagSet writes usage to w.
func newFlagSet(name string, w io.Writer) (*flag.FlagSet, *args, *multiFlag) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(w)

	a := &args{}
	var rawFields multiFlag

	fs.StringVar(&a.level, "level", "", "filter by log level (e.g. info, warn, error)")
	fs.StringVar(&a.since, "since", "", "include entries at or after this timestamp (RFC3339 or common log formats)")
	fs.StringVar(&a.until, "until", "", "include entries at or before this timestamp")
	fs.Var(&rawFields, "field", "filter by field key=value (repeatable)")

	fs.StringVar(&a.format, "format", "text", fmt.Sprintf("output format: %s", strings.Join(output.ValidFormats(), ", ")))
	fs.BoolVar(&a.noColor, "no-color", false, "disable ANSI color output")
	fs.StringVar(&a.outFile, "out", "", "write output to file instead of stdout")

	fs.BoolVar(&a.showStats, "stats", false, "print summary statistics after output")
	fs.BoolVar(&a.verbose, "verbose", false, "enable verbose/debug logging")

	fs.Usage = func() {
		fmt.Fprintf(w, "Usage: %s [flags] [logfile]\n\n", name)
		fmt.Fprintf(w, "logslice parses, filters, and formats structured log files.\n\n")
		fmt.Fprintf(w, "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(w, "\nIf no logfile is given, logslice reads from stdin.\n")
	}

	return fs, a, &rawFields
}

// parseFieldArgs converts raw "key=value" strings into typed fieldArg values.
// Returns an error if any entry is malformed.
func parseFieldArgs(raw []string) ([]fieldArg, error) {
	out := make([]fieldArg, 0, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid --field value %q: expected key=value", r)
		}
		out = append(out, fieldArg{Key: parts[0], Value: parts[1]})
	}
	return out, nil
}
