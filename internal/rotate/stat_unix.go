//go:build !windows

package rotate

import (
	"fmt"
	"syscall"
)

func stat(path string) (fileInfo, error) {
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return fileInfo{}, fmt.Errorf("rotate: stat %s: %w", path, err)
	}
	return fileInfo{
		inode: s.Ino,
		size:  s.Size,
	}, nil
}
