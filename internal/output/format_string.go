package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format constant.
// Returns an error if the format is not recognized.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "json":
		return FormatJSON, nil
	case "text", "":
		return FormatText, nil
	case "table":
		return FormatTable, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of json, text, table", s)
	}
}

// String implements the Stringer interface for Format.
func (f Format) String() string {
	return string(f)
}

// ValidFormats returns all supported format names.
func ValidFormats() []string {
	return []string{
		string(FormatJSON),
		string(FormatText),
		string(FormatTable),
	}
}
