// Package highlight provides ANSI color highlighting for log output.
package highlight

import "fmt"

// Color represents an ANSI terminal color code.
type Color int

const (
	Reset  Color = 0
	Red    Color = 31
	Yellow Color = 33
	Cyan   Color = 36
	White  Color = 37
	Bold   Color = 1
)

// Colorize wraps text with the given ANSI color code.
func Colorize(text string, c Color) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", c, text)
}

// LevelColor returns the appropriate color for a log level string.
func LevelColor(level string) Color {
	switch level {
	case "ERROR", "FATAL":
		return Red
	case "WARN", "WARNING":
		return Yellow
	case "INFO":
		return Cyan
	default:
		return White
	}
}

// ForLevel returns the level string wrapped in its associated color.
func ForLevel(level string) string {
	return Colorize(level, LevelColor(level))
}

// Bold wraps text in bold ANSI formatting.
func BoldText(text string) string {
	return Colorize(text, Bold)
}
