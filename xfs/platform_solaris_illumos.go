//go:build solaris || illumos

package xfs

import "golang.org/x/sys/unix"

var fnMockable_Solaris_Illumos_Statvfs = unix.Statvfs

func init() {
	fnMockable_GetVolumeFreeSpace = getVolumeFreeSpace
}

func getVolumeFreeSpace(currentPath string, requiredBytes uint64) (bool, error) {
	var stat unix.Statvfs_t
	if err := fnMockable_Solaris_Illumos_Statvfs(currentPath, &stat); err != nil {
		return false, err
	}
	freeSpace := uint64(stat.Bavail) * uint64(stat.Bsize)
	return freeSpace >= requiredBytes, nil
}
