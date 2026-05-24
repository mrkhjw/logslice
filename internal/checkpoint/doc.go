// Package checkpoint provides persistent read-position (offset) tracking
// for log files processed by logslice.
//
// A Store is backed by a single JSON file on disk. Each tracked log file
// is keyed by its absolute path. On every Set or Delete call the store is
// atomically written via a rename so that a crash cannot corrupt the file.
//
// Typical usage:
//
//	store, err := checkpoint.New("/var/lib/logslice/checkpoint.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	offset := store.Get("/var/log/app.log")
//	// ... parse from offset ...
//	_ = store.Set("/var/log/app.log", newOffset)
package checkpoint
