//go:build darwin

package xfs

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

func TestHasSpaceAvailable_DarwinSpecialist(t *testing.T) {
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
			name:        "Darwin Syscall Error: Should return error if statfs fails",
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
			name:       "Darwin Math Logic: Should return false if free space is lower than required",
			inputPath:  "/mock/dir",
			inputBytes: 5000,
			want:       false,
			wantErr:    false,
			mockUnixStatfs: func(path string, buf *unix.Statfs_t) error {
				buf.Bavail = 10 // 10 blocos
				buf.Bsize = 100 // 100 bytes por bloco = 1000 bytes livres
				return nil
			},
		},
		{
			name:       "Darwin Math Logic: Should return true if free space is plenty",
			inputPath:  "/mock/dir",
			inputBytes: 1000,
			want:       true,
			wantErr:    false,
			mockUnixStatfs: func(path string, buf *unix.Statfs_t) error {
				buf.Bavail = 50
				buf.Bsize = 200 // 50 * 200 = 10000 bytes livres
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Isola o ambiente usando os nossos padrões mockables comuns
			mockFunction(t, &fnMockable_RetrieveFullPath, func(p string) (string, error) { return p, nil })
			mockFunction(t, &fnMockable_Exists, func(p string) bool { return true })

			// Devolvemos para o fnMockable_GetVolumeFreeSpace a função REAL Darwin que queremos testar
			mockFunction(t, &fnMockable_GetVolumeFreeSpace, getVolumeFreeSpace)

			// Injeta o mock da Syscall Darwin (aponta para fnMockable_Darwin_Statfs do arquivo principal _darwin.go)
			if tt.mockUnixStatfs != nil {
				mockFunction(t, &fnMockable_Darwin_Statfs, tt.mockUnixStatfs)
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
