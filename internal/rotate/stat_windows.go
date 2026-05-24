//go:build windows

package rotate

import (
	"fmt"
	"os"
)

// On Windows inode is not available via syscall in a portable way;
// we fall back to using file size only for rotation detection.
func stat(path string) (fileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return fileInfo{}, fmt.Errorf("rotate: stat %s: %w", path, err)
	}
	return fileInfo{
		inode: 0,
		size:  fi.Size(),
	}, nil
}
