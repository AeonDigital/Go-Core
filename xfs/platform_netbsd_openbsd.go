//go:build netbsd || openbsd

package xfs

import "golang.org/x/sys/unix"

var fnMockable_NetBSD_OpenBSD_Statfs = unix.Statfs

func init() {
	fnMockable_GetVolumeFreeSpace = getVolumeFreeSpace
}

func getVolumeFreeSpace(currentPath string, requiredBytes uint64) (bool, error) {
	var stat unix.Statfs_t
	if err := fnMockable_NetBSD_OpenBSD_Statfs(currentPath, &stat); err != nil {
		return false, err
	}
	freeSpace := uint64(stat.F_bavail) * uint64(stat.F_bsize)
	return freeSpace >= requiredBytes, nil
}
