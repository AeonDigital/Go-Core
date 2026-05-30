//go:build windows

package xfs

import (
	"os"
	"strings"
	"testing"
)

func TestHasSpaceAvailable_Windows_Specialist(t *testing.T) {
	tests := []struct {
		name             string
		inputPath        string
		inputBytes       uint64
		want             bool
		wantErr          bool
		errContains      string
		mockWindowsSpace func(dir *uint16, free, total, freeTotal *uint64) error
	}{
		{
			name:        "Windows Syscall Error: Should return error if GetDiskFreeSpaceEx fails",
			inputPath:   "C:\\mock\\dir",
			inputBytes:  1024,
			want:        false,
			wantErr:     true,
			errContains: "permission denied",
			mockWindowsSpace: func(dir *uint16, free, total, freeTotal *uint64) error {
				return os.ErrPermission // Simula falha de permissão nativa do Windows
			},
		},
		{
			name:       "Windows Math Logic: Should return false if free bytes available is lower than required",
			inputPath:  "C:\\mock\\dir",
			inputBytes: 5000, // Requer 5000 bytes
			want:       false,
			wantErr:    false,
			mockWindowsSpace: func(dir *uint16, free, total, freeTotal *uint64) error {
				*free = 1000 // Injeta apenas 1000 bytes livres informados pelo Windows
				return nil
			},
		},
		{
			name:       "Windows Math Logic: Should return true if free bytes available is plenty",
			inputPath:  "C:\\mock\\dir",
			inputBytes: 1000, // Requer 1000 bytes
			want:       true,
			wantErr:    false,
			mockWindowsSpace: func(dir *uint16, free, total, freeTotal *uint64) error {
				*free = 10000 // Injeta 10000 bytes livres informados pelo Windows
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Isola o ambiente usando os nossos padrões mockables comuns
			mockFunction(t, &fnMockable_RetrieveFullPath, func(p string) (string, error) { return p, nil })
			mockFunction(t, &fnMockable_Exists, func(p string) bool { return true })

			// Devolvemos para o fnMockable_GetVolumeFreeSpace a função REAL do Windows que queremos testar
			mockFunction(t, &fnMockable_GetVolumeFreeSpace, getVolumeFreeSpace)

			// Injeta o mock da Syscall do Windows (aponta para fnMockable_Windows_GetDiskFreeSpaceEx do arquivo _windows.go)
			if tt.mockWindowsSpace != nil {
				mockFunction(t, &fnMockable_Windows_GetDiskFreeSpaceEx, tt.mockWindowsSpace)
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
