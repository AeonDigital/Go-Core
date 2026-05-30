//go:build linux || freebsd || dragonfly

package xfs

import "golang.org/x/sys/unix"

var fnMockable_Linux_FreeBSD_Dragonfly_Statfs = unix.Statfs

func init() {
	fnMockable_GetVolumeFreeSpace = getVolumeFreeSpace
}

func getVolumeFreeSpace(currentPath string, requiredBytes uint64) (bool, error) {
	var stat unix.Statfs_t
	if err := fnMockable_Linux_FreeBSD_Dragonfly_Statfs(currentPath, &stat); err != nil {
		return false, err
	}
	freeSpace := uint64(stat.Bavail) * uint64(stat.Bsize)
	return freeSpace >= requiredBytes, nil
}
