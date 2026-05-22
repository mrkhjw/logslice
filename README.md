# logslice

Fast log file parser and filter tool with structured output and time-range queries.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
# Filter logs between two timestamps
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log

# Filter by log level and output as JSON
logslice --level error --format json app.log

# Read from stdin
cat app.log | logslice --from "2024-01-15T08:00:00" --level warn

# Write filtered output to a file
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log -o out.log
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--level` | Filter by log level (info, warn, error) |
| `--format` | Output format: `text` (default) or `json` |
| `-o` | Write output to file instead of stdout |

---

## Features

- Parses common log formats automatically
- Time-range filtering with RFC3339 timestamps
- Structured JSON output for downstream processing
- Streams large files without loading them fully into memory

---

## License

MIT © 2024 yourusername