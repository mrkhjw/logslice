// Package replay replays a slice of log entries at a configurable speed,
// simulating the original timing between events.
//
// # Basic usage
//
//	opts := replay.Options{
//		Speed:    2.0,           // twice real-time
//		MaxDelay: 3 * time.Second,
//	}
//	ch := replay.Apply(entries, opts)
//	for entry := range ch {
//		fmt.Println(entry.Message)
//	}
//
// # Speed multiplier
//
// A Speed of 1.0 replays at real time. Values greater than 1.0 speed up
// playback; values between 0 and 1.0 slow it down. A value of 0 or less
// is treated as 1.0.
//
// # MaxDelay
//
// Long silences in a log file can stall a replay. MaxDelay caps the
// inter-entry pause so the replay remains responsive even when the
// original log contains large time gaps.
package replay
