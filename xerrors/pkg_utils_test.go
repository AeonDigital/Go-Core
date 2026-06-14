package xerrors_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xerrors"
)

func TestPrint_CoverageSuite(t *testing.T) {
	// Scenario 1: err is nil, function should return immediately doing nothing
	xerrors.Print(nil)

	// Scenario 2: err is present, must capture os.Stderr buffer output stream
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	testErr := errors.New("severe technical system breakdown")
	xerrors.Print(testErr)

	// Close writer stream and restore original system Stderr block
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	outputStr := buf.String()

	if !strings.Contains(outputStr, "severe technical system breakdown") {
		t.Errorf("Print() output = %q, expected to contain %q", outputStr, testErr.Error())
	}
}
