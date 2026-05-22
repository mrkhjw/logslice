package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

// Config holds the parsed CLI arguments.
type Config struct {
	FilePath  string
	Level     string
	Since     string
	Until     string
	Field     string
	Format    string
}

// Run parses arguments and executes the main pipeline.
func Run(args []string) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	r, err := openInput(cfg.FilePath)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	if rc, ok := r.(io.Closer); ok {
		defer rc.Close()
	}

	p := parser.New(r)
	entries, err := p.Parse()
	if err != nil {
		return fmt.Errorf("parsing: %w", err)
	}

	opts := buildFilterOpts(cfg)
	filtered := filter.Filter(entries, opts)

	fmt, err := output.ParseFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	f := output.New(os.Stdout, fmt)
	return f.Write(filtered)
}

func parseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	cfg := &Config{}

	fs.StringVar(&cfg.Level, "level", "", "filter by log level (e.g. INFO, ERROR)")
	fs.StringVar(&cfg.Since, "since", "", "include entries at or after this time (RFC3339)")
	fs.StringVar(&cfg.Until, "until", "", "include entries at or before this time (RFC3339)")
	fs.StringVar(&cfg.Field, "field", "", "filter by field key=value")
	fs.StringVar(&cfg.Format, "format", "text", fmt.Sprintf("output format: %v", output.ValidFormats()))

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() > 0 {
		cfg.FilePath = fs.Arg(0)
	}
	return cfg, nil
}

func openInput(path string) (io.Reader, error) {
	if path == "" || path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func buildFilterOpts(cfg *Config) filter.Options {
	opts := filter.Options{
		Level: cfg.Level,
		Field: cfg.Field,
	}
	if cfg.Since != "" {
		if t, err := time.Parse(time.RFC3339, cfg.Since); err == nil {
			opts.Since = &t
		}
	}
	if cfg.Until != "" {
		if t, err := time.Parse(time.RFC3339, cfg.Until); err == nil {
			opts.Until = &t
		}
	}
	return opts
}
