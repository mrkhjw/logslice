package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"time"
)

// Common log timestamp formats to attempt parsing.
var timestampFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// logLineRegex matches: TIMESTAMP LEVEL message [optional key=value pairs]
var logLineRegex = regexp.MustCompile(
	`^(\S+(?:\s\S+)?)\s+(DEBUG|INFO|WARN|ERROR|FATAL)\s+(.+)$`,
)

// Parser reads log entries from an io.Reader.
type Parser struct {
	reader io.Reader
}

// New creates a new Parser for the given reader.
func New(r io.Reader) *Parser {
	return &Parser{reader: r}
}

// Parse reads all log entries from the reader.
func (p *Parser) Parse() ([]*Entry, error) {
	var entries []*Entry
	scanner := bufio.NewScanner(p.reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		entry := parseLine(line)
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func parseLine(line string) *Entry {
	entry := &Entry{Raw: line, Fields: make(map[string]string)}
	matches := logLineRegex.FindStringSubmatch(line)
	if matches == nil {
		entry.Level = LevelUnknown
		entry.Message = line
		return entry
	}
	entry.Timestamp = parseTimestamp(matches[1])
	entry.Level = ParseLevel(matches[2])
	entry.Message, entry.Fields = extractFields(matches[3])
	return entry
}

func parseTimestamp(s string) time.Time {
	for _, format := range timestampFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func extractFields(s string) (string, map[string]string) {
	fields := make(map[string]string)
	kvRegex := regexp.MustCompile(`(\w+)=(\S+)`)
	msg := kvRegex.ReplaceAllStringFunc(s, func(m string) string {
		parts := strings.SplitN(m, "=", 2)
		if len(parts) == 2 {
			fields[parts[0]] = parts[1]
		}
		return ""
	})
	return strings.TrimSpace(msg), fields
}
