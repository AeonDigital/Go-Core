package xcli_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/AeonDigital/Go-Core/xcli"
)

// TestPrint verifies that messages are correctly written to standard output.
func TestPrint(t *testing.T) {
	// 1. Capture os.Stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stdout: %v", err)
	}
	os.Stdout = w

	// 2. Execute the function with a single message (behaves like Println)
	xcli.Print("hello world")

	// 3. Execute the function with arguments (behaves like Printf + \n)
	xcli.Print("user %s has id %d", "john", 42)

	// Close the writer so we can read from the pipe
	w.Close()
	os.Stdout = oldStdout

	// 4. Read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read from stdout pipe: %v", err)
	}
	output := buf.String()

	// 5. Assert the results
	expected1 := "hello world\n"
	expected2 := "user john has id 42\n"

	if !strings.Contains(output, expected1) {
		t.Errorf("Print layout 1 failed.\nExpected to contain: %q\nFull output: %q", expected1, output)
	}
	if !strings.Contains(output, expected2) {
		t.Errorf("Print layout 2 failed.\nExpected to contain: %q\nFull output: %q", expected2, output)
	}
}
