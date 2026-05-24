// Package alert implements threshold-based alerting for log entry streams.
//
// Rules are defined by a log level, an entry count threshold, and a sliding
// time window. When the number of entries matching a rule's level exceeds the
// threshold within the window, a handler function is invoked with an Event
// describing the breach.
//
// Basic usage:
//
//	rules := alert.NewBuilder().
//		OnError(10, time.Minute).
//		OnFatal(1, time.Minute).
//		Build()
//
//	a := alert.New(rules, func(e alert.Event) {
//		fmt.Println(e.Message)
//	})
//
//	a.Evaluate(entries)
//
// The Alerter is stateful: repeated calls to Evaluate accumulate timestamps
// across invocations, enabling streaming use-cases where entries arrive in
// batches over time.
package alert
