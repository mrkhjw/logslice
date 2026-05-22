package highlight

// Theme defines a color theme for log output.
type Theme struct {
	Error string
	Warn  string
	Info  string
	Debug string
	Fatal string
	Bold  string
	Reset string
}

// DefaultTheme is the default ANSI color theme.
var DefaultTheme = Theme{
	Error: "\033[31m", // Red
	Warn:  "\033[33m", // Yellow
	Info:  "\033[32m", // Green
	Debug: "\033[36m", // Cyan
	Fatal: "\033[35m", // Magenta
	Bold:  "\033[1m",
	Reset: "\033[0m",
}

// NoColorTheme is a theme with no ANSI escape codes (plain text).
var NoColorTheme = Theme{
	Error: "",
	Warn:  "",
	Info:  "",
	Debug: "",
	Fatal: "",
	Bold:  "",
	Reset: "",
}

// ActiveTheme is the theme used by highlight functions.
var ActiveTheme = DefaultTheme

// SetTheme sets the active theme used for colorization.
func SetTheme(t Theme) {
	ActiveTheme = t
}

// DisableColor switches to the no-color theme.
func DisableColor() {
	ActiveTheme = NoColorTheme
}

// EnableColor restores the default color theme.
func EnableColor() {
	ActiveTheme = DefaultTheme
}
