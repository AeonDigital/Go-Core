//go:build windows

package xfs

import "golang.org/x/sys/windows"

var fnMockable_Windows_GetDiskFreeSpaceEx = windows.GetDiskFreeSpaceEx

func init() {
	fnMockable_GetVolumeFreeSpace = getVolumeFreeSpace
}

func getVolumeFreeSpace(currentPath string, requiredBytes uint64) (bool, error) {
	pathPtr, err := windows.UTF16PtrFromString(currentPath)
	if err != nil {
		return false, err
	}
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	if err := fnMockable_Windows_GetDiskFreeSpaceEx(pathPtr, &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes); err != nil {
		return false, err
	}
	return freeBytesAvailable >= requiredBytes, nil
}
