// Package tail implements real-time log file following (akin to `tail -f`).
//
// It polls a file for newly appended bytes, splits them into lines, and
// parses each line through the configured parser.Parser, forwarding
// successfully parsed entries to a caller-supplied channel.
//
// Usage:
//
//	p := parser.New()
//	out := make(chan parser.Entry, 64)
//	stop := make(chan struct{})
//
//	go tail.Tail("/var/log/app.log", tail.DefaultOptions(), p, out, stop)
//
//	for entry := range out {
//		fmt.Println(entry.Message)
//	}
package tail
