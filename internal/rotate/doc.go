// Package rotate provides log file rotation detection for logslice.
//
// It allows callers to detect when a log file has been rotated (replaced or
// truncated) so that readers can reopen and continue processing from the new
// file. Rotation is detected by comparing the file's inode and size against
// a stored baseline.
//
// Basic usage:
//
//	detector, err := rotate.New("/var/log/app.log", rotate.DefaultOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if rotated, _ := detector.Rotated(); rotated {
//	    // reopen the file
//	    detector.Reset()
//	}
//
// On Windows, inode-based detection is unavailable; only size-based truncation
// detection is supported.
package rotate
