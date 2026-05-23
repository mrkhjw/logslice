// Package tail provides functionality for following log files in real-time,
// similar to `tail -f`, emitting new log entries as they are written.
package tail

import (
	"io"
	"os"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options configures the tail behaviour.
type Options struct {
	// PollInterval is how often to check for new data when no inotify is available.
	PollInterval time.Duration
	// MaxRetries is the number of consecutive empty reads before giving up (0 = infinite).
	MaxRetries int
}

// DefaultOptions returns sensible defaults for tail.
func DefaultOptions() Options {
	return Options{
		PollInterval: 250 * time.Millisecond,
		MaxRetries:   0,
	}
}

// Tail opens the named file and streams newly appended log entries into out.
// It blocks until the context is cancelled via the stop channel or MaxRetries
// consecutive empty reads occur (when MaxRetries > 0).
func Tail(path string, opts Options, p *parser.Parser, out chan<- parser.Entry, stop <-chan struct{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only emit new lines.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	if opts.PollInterval <= 0 {
		opts.PollInterval = DefaultOptions().PollInterval
	}

	buf := make([]byte, 0, 4096)
	retries := 0

	for {
		select {
		case <-stop:
			return nil
		default:
		}

		chunk := make([]byte, 4096)
		n, err := f.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			if opts.MaxRetries > 0 {
				retries++
				if retries >= opts.MaxRetries {
					return nil
				}
			}
			time.Sleep(opts.PollInterval)
			continue
		}

		retries = 0
		buf = append(buf, chunk[:n]...)
		buf = flushLines(buf, p, out)
	}
}

// flushLines splits buf on newlines, parses complete lines, and returns any
// remaining incomplete line fragment.
func flushLines(buf []byte, p *parser.Parser, out chan<- parser.Entry) []byte {
	start := 0
	for i, b := range buf {
		if b == '\n' {
			line := string(buf[start:i])
			start = i + 1
			if entry, ok := p.ParseLine(line); ok {
				out <- entry
			}
		}
	}
	return buf[start:]
}
