//go:build solaris || illumos

package xfs

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

func TestHasSpaceAvailable_Solaris_Illumos_Specialist(t *testing.T) {
	tests := []struct {
		name               string
		inputPath          string
		inputBytes         uint64
		want               bool
		wantErr            bool
		errContains        string
		mockSolarisStatvfs func(path string, buf *unix.Statvfs_t) error
	}{
		{
			name:        "Solaris Syscall Error: Should return error if statvfs fails",
			inputPath:   "/mock/dir",
			inputBytes:  1024,
			want:        false,
			wantErr:     true,
			errContains: "permission denied",
			mockSolarisStatvfs: func(path string, buf *unix.Statvfs_t) error {
				return os.ErrPermission
			},
		},
		{
			name:       "Solaris Math Logic: Should return false if free space is lower than required",
			inputPath:  "/mock/dir",
			inputBytes: 5000,
			want:       false,
			wantErr:    false,
			mockSolarisStatvfs: func(path string, buf *unix.Statvfs_t) error {
				// O Solaris usa a estrutura Statvfs_t com Bavail e Bsize normais
				buf.Bavail = 10
				buf.Bsize = 100
				return nil
			},
		},
		{
			name:       "Solaris Math Logic: Should return true if free space is plenty",
			inputPath:  "/mock/dir",
			inputBytes: 1000,
			want:       true,
			wantErr:    false,
			mockSolarisStatvfs: func(path string, buf *unix.Statvfs_t) error {
				buf.Bavail = 50
				buf.Bsize = 200
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFunction(t, &fnMockable_RetrieveFullPath, func(p string) (string, error) { return p, nil })
			mockFunction(t, &fnMockable_Exists, func(p string) bool { return true })
			mockFunction(t, &fnMockable_GetVolumeFreeSpace, getVolumeFreeSpace)

			if tt.mockSolarisStatvfs != nil {
				mockFunction(t, &fnMockable_Solaris_Illumos_Statvfs, tt.mockSolarisStatvfs)
			}

			got, err := HasSpaceAvailable(tt.inputPath, tt.inputBytes)

			if (err != nil) != tt.wantErr {
				t.Fatalf("HasSpaceAvailable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if got != tt.want {
				t.Errorf("HasSpaceAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
