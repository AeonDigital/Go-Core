//go:build netbsd || openbsd

package xfs

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

func TestHasSpaceAvailable_NetBSD_OpenBSD_Specialist(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		inputBytes     uint64
		want           bool
		wantErr        bool
		errContains    string
		mockUnixStatfs func(path string, buf *unix.Statfs_t) error
	}{
		{
			name:        "BSD Syscall Error: Should return error if statfs fails",
			inputPath:   "/mock/dir",
			inputBytes:  1024,
			want:        false,
			wantErr:     true,
			errContains: "permission denied",
			mockUnixStatfs: func(path string, buf *unix.Statfs_t) error {
				return os.ErrPermission
			},
		},
		{
			name:       "BSD Math Logic: Should return false if free space is lower than required",
			inputPath:  "/mock/dir",
			inputBytes: 5000,
			want:       false,
			wantErr:    false,
			mockUnixStatfs: func(path string, buf *unix.Statfs_t) error {
				// Usa o formato nativo com F_ exigido pelos BSDs tradicionais
				buf.F_bavail = 10
				buf.F_bsize = 100
				return nil
			},
		},
		{
			name:       "BSD Math Logic: Should return true if free space is plenty",
			inputPath:  "/mock/dir",
			inputBytes: 1000,
			want:       true,
			wantErr:    false,
			mockUnixStatfs: func(path string, buf *unix.Statfs_t) error {
				buf.F_bavail = 50
				buf.F_bsize = 200
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFunction(t, &fnMockable_RetrieveFullPath, func(p string) (string, error) { return p, nil })
			mockFunction(t, &fnMockable_Exists, func(p string) bool { return true })
			mockFunction(t, &fnMockable_GetVolumeFreeSpace, getVolumeFreeSpace)

			if tt.mockUnixStatfs != nil {
				mockFunction(t, &fnMockable_NetBSD_OpenBSD_Statfs, tt.mockUnixStatfs)
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
