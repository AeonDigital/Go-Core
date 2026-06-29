package xclifn_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliconstt"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// TestValidateStringLimitsSuccessAndErrors validates text constraints tracking
// rune length boundaries and regular expression pattern validations.
func TestValidateStringLimitsSuccessAndErrors(t *testing.T) {
	minVal := 5
	maxVal := 10

	// 1. Success cases
	specSuccess := xclistruc.Flag{
		MinLength: &minVal,
		MaxLength: &maxVal,
		Regex:     "^[a-z]+$",
	}
	err := xclifn.ValidateStringLimits("--name", "abcdef", specSuccess)
	if err != nil {
		t.Errorf("expected no error for valid string boundaries, got: %v", err)
	}

	// 2. Error case: String too short
	specShort := xclistruc.Flag{
		MinLength: &minVal,
	}
	errShort := xclifn.ValidateStringLimits("-n", "abc", specShort)
	if errShort == nil {
		t.Fatalf("expected error for string below minimum limit, got nil")
	}
	msgShort := errShort.Error()
	if !strings.Contains(msgShort, "Rule: min length: 5") || !strings.Contains(msgShort, "Given: abc") {
		t.Errorf("unexpected visual layout structure for short string: %q", msgShort)
	}

	// 3. Error case: String too long
	specLong := xclistruc.Flag{
		MaxLength: &maxVal,
	}
	errLong := xclifn.ValidateStringLimits("--name", "abcdefghijklmnop", specLong)
	if errLong == nil {
		t.Fatalf("expected error for string exceeding maximum limit, got nil")
	}
	msgLong := errLong.Error()
	if !strings.Contains(msgLong, "Rule: max length: 10") {
		t.Errorf("unexpected visual layout structure for long string: %q", msgLong)
	}

	// 4. Error case: Pattern mismatch
	specRegex := xclistruc.Flag{
		Regex: "^[0-9]+$",
	}
	errRegex := xclifn.ValidateStringLimits("--code", "123-ABC", specRegex)
	if errRegex == nil {
		t.Fatalf("expected error for regex pattern mismatch, got nil")
	}

	// 5. Technical Edge Case: Broken regex compiling syntax pattern
	specBroken := xclistruc.Flag{
		Regex: "[a-z",
	}
	errBrokenRegex := xclifn.ValidateStringLimits("--code", "value", specBroken)
	if errBrokenRegex == nil {
		t.Fatalf("expected error for malformed regex syntax compiling, got nil")
	}
}

// TestValidateNumberLimitsAllTypes triggers bounds tracking for all quantifiable numbers
// including integers, float decimals, durations, and timestamps using strongly-typed generic constraints.
func TestValidateNumberLimitsAllTypes(t *testing.T) {
	// 1. Test Int limits
	specIntValid := xclistruc.Flag{Min: 10, Max: 20}
	if err := xclifn.ValidateNumberLimits[int]("--count", 15, specIntValid); err != nil {
		t.Errorf("valid int bounds failed: %v", err)
	}

	specIntLow := xclistruc.Flag{Min: 10}
	if err := xclifn.ValidateNumberLimits[int]("--count", 5, specIntLow); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("expected out of range for low int block")
	}

	specIntHigh := xclistruc.Flag{Max: 20}
	if err := xclifn.ValidateNumberLimits[int]("--count", 25, specIntHigh); err == nil {
		t.Errorf("expected out of range for high int block")
	}

	// 2. Test Float64 limits
	specFloatValid := xclistruc.Flag{Min: 1.0, Max: 5.0}
	if err := xclifn.ValidateNumberLimits[float64]("--ratio", 3.14, specFloatValid); err != nil {
		t.Errorf("valid float bounds failed: %v", err)
	}

	specFloatLow := xclistruc.Flag{Min: 1.0}
	if err := xclifn.ValidateNumberLimits[float64]("--ratio", 0.5, specFloatLow); err == nil {
		t.Errorf("expected out of range for low float block")
	}

	specFloatHigh := xclistruc.Flag{Max: 5.0}
	if err := xclifn.ValidateNumberLimits[float64]("--ratio", 6.0, specFloatHigh); err == nil {
		t.Errorf("expected out of range for high float block")
	}

	// 3. Test Duration limits
	specDurValid := xclistruc.Flag{Min: 1 * time.Minute, Max: 10 * time.Minute}
	if err := xclifn.ValidateNumberLimits[time.Duration]("--timeout", 5*time.Minute, specDurValid); err != nil {
		t.Errorf("valid duration bounds failed: %v", err)
	}

	specDurLow := xclistruc.Flag{Min: 1 * time.Minute}
	if err := xclifn.ValidateNumberLimits[time.Duration]("--timeout", 30*time.Second, specDurLow); err == nil {
		t.Errorf("expected out of range for low duration block")
	}

	specDurHigh := xclistruc.Flag{Max: 10 * time.Minute}
	if err := xclifn.ValidateNumberLimits[time.Duration]("--timeout", 15*time.Minute, specDurHigh); err == nil {
		t.Errorf("expected out of range for high duration block")
	}

	// 4. Test Time limits
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	specTimeValid := xclistruc.Flag{Min: past, Max: future}
	if err := xclifn.ValidateNumberLimits[time.Time]("--date", now, specTimeValid); err != nil {
		t.Errorf("valid time bounds failed: %v", err)
	}

	specTimeLow := xclistruc.Flag{Min: past}
	if err := xclifn.ValidateNumberLimits[time.Time]("--date", past.Add(-1*time.Hour), specTimeLow); err == nil {
		t.Errorf("expected out of range for low time block")
	}

	specTimeHigh := xclistruc.Flag{Max: future}
	if err := xclifn.ValidateNumberLimits[time.Time]("--date", future.Add(1*time.Hour), specTimeHigh); err == nil {
		t.Errorf("expected out of range for high time block")
	}
}

// TestValidateArrayLimitsSuccessAndErrors verifies slice capacity boundaries.
func TestValidateArrayLimitsSuccessAndErrors(t *testing.T) {
	minItems := 2
	maxItems := 4

	// 1. Success cases
	specValid := xclistruc.Flag{MinItems: &minItems, MaxItems: &maxItems}
	// Passed the length directly (3 elements) matching the new function signature
	err := xclifn.ValidateArrayLimits("--ids", 3, specValid)
	if err != nil {
		t.Errorf("valid array evaluation failed: %v", err)
	}

	// 2. Empty slice case
	specEmpty := xclistruc.Flag{}
	errEmpty := xclifn.ValidateArrayLimits("--ids", 0, specEmpty)
	if errEmpty != nil {
		t.Errorf("blank array payload fallback track failed: %v", errEmpty)
	}

	// 3. Error case: Too few elements
	specFew := xclistruc.Flag{MinItems: &minItems}
	errFew := xclifn.ValidateArrayLimits("--ids", 1, specFew)
	if errFew == nil {
		t.Fatalf("expected error for array elements under min capacity limits, got nil")
	}
	if !strings.Contains(errFew.Error(), "min items: 2") {
		t.Errorf("unexpected elements count error output: %q", errFew.Error())
	}

	// 4. Error case: Too many elements
	specMany := xclistruc.Flag{MaxItems: &maxItems}
	errMany := xclifn.ValidateArrayLimits("--ids", 5, specMany)
	if errMany == nil {
		t.Fatalf("expected error for array elements exceeding max capacity limits, got nil")
	}
}

// TestNewValidationErrorExpectedLayout forces a full structural track inside the private
// generator via exported bridge to guarantee coverage across the Expected field rendering block.
func TestNewValidationErrorExpectedLayout(t *testing.T) {
	// Chamada através da função exposta pelo arquivo de gateway
	err := xclifn.ExportNewValidationError(
		"--output",
		"format failure",
		"raw-text",
		"valid-json-format", // Preenche o campo para ativar a linha do IF
		"type boundary",
	)

	if err == nil {
		t.Fatalf("expected structured error instance, got nil")
	}

	msg := err.Error()

	// Garante a cobertura da linha interna e a validação do formato visual
	if !strings.Contains(msg, "Expected: valid-json-format") {
		t.Errorf("missing Expected visual layout segment in validation error block: %q", msg)
	}

	if !strings.Contains(msg, "Rule: type boundary") || !strings.Contains(msg, "Given: raw-text") {
		t.Errorf("unexpected visual layout structure in generated error: %q", msg)
	}
}

// TestValidateDiskResourcesSuccessAndErrors executes physical filesystem evaluation bounds
// ensuring telemetry lookups, structural targets type and permissions match the flag specs.
func TestValidateDiskResourcesSuccessAndErrors(t *testing.T) {
	// Setup stable local OS environment traces using standard testing capabilities
	tmpDir := t.TempDir()

	// Create a real mock file inside the isolated sandbox directory for lookups tracking
	mockFilePath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(mockFilePath, []byte("{}"), 0666); err != nil {
		t.Fatalf("failed to setup mock test files infrastructure layout: %v", err)
	}

	// 1. Success cases: Filepath validation tracking
	specFileValid := xclistruc.Flag{
		Type:      xcliconstt.FlagType("filepath"),
		MustExist: true,
		Access:    xcliconstt.AccessMode("readwrite"),
	}
	if err := xclifn.ValidateDiskResources("--cfg", mockFilePath, specFileValid); err != nil {
		t.Errorf("expected no validation errors for valid existing files, got: %v", err)
	}

	// 2. Success cases: Dirpath validation tracking
	specDirValid := xclistruc.Flag{
		Type:      xcliconstt.FlagType("dirpath"),
		MustExist: true,
		Access:    xcliconstt.AccessMode("readwrite"),
	}
	if err := xclifn.ValidateDiskResources("--src", tmpDir, specDirValid); err != nil {
		t.Errorf("expected no validation errors for valid existing directory, got: %v", err)
	}

	// 3. Error case: Missing mandatory resource existence boundary rule
	specMustExist := xclistruc.Flag{
		Type:      xcliconstt.FlagType("filepath"),
		MustExist: true,
	}
	ghostPath := filepath.Join(tmpDir, "ghost_file.txt")
	errExist := xclifn.ValidateDiskResources("--input", ghostPath, specMustExist)
	if errExist == nil {
		t.Fatalf("expected error for non-existent path resource required by spec, got nil")
	}
	if !strings.Contains(errExist.Error(), "path resource does not exist") {
		t.Errorf("unexpected error payload structure for non-existent bounds: %q", errExist.Error())
	}

	// 4. Error case: Structural mismatch constraint (directory provided, but file required)
	specFileConstraint := xclistruc.Flag{
		Type: xcliconstt.FlagType("filepath"),
	}
	errFileConstraint := xclifn.ValidateDiskResources("--file", tmpDir, specFileConstraint)
	if errFileConstraint == nil {
		t.Fatalf("expected error when directory layout targets a strict filepath rule, got nil")
	}
	if !strings.Contains(errFileConstraint.Error(), "resource location is a directory, file required") {
		t.Errorf("unexpected error payload structure for file constraint type: %q", errFileConstraint.Error())
	}

	// 5. Error case: Structural mismatch constraint (file provided, but directory required)
	specDirConstraint := xclistruc.Flag{
		Type: xcliconstt.FlagType("dirpath"),
	}
	errDirConstraint := xclifn.ValidateDiskResources("--dir", mockFilePath, specDirConstraint)
	if errDirConstraint == nil {
		t.Fatalf("expected error when file layout targets a strict dirpath rule, got nil")
	}
	if !strings.Contains(errDirConstraint.Error(), "resource location is a file, directory required") {
		t.Errorf("unexpected error payload structure for directory constraint type: %q", errDirConstraint.Error())
	}
}

// TestValidateDiskResourcesHardwareAndPermissionBlocks injects simulated OS failures
// via exported bridges to guarantee 100% coverage across complex disk error blocks.
func TestValidateDiskResourcesHardwareAndPermissionBlocks(t *testing.T) {
	// Garante o reset das pontes ao final da execução do teste
	defer xclifn.ResetDiskBridges()

	specFile := xclistruc.Flag{
		Type:      xcliconstt.FlagType("filepath"),
		MustExist: true,
		Access:    xcliconstt.AccessMode("readwrite"),
	}

	specDir := xclistruc.Flag{
		Type:      xcliconstt.FlagType("dirpath"),
		MustExist: true,
		Access:    xcliconstt.AccessMode("readwrite"),
	}

	// 1. Cobrir o erro de telemetria de disco ("disk telemetry access block")
	xclifn.SetDiskBridgesMocks(
		func(name string) (os.FileInfo, error) {
			return nil, errors.New("hardware I/O block failure simulation")
		},
		nil, nil, nil,
	)
	errStat := xclifn.ValidateDiskResources("--cfg", "/any/path", specFile)
	if errStat == nil || !strings.Contains(errStat.Error(), "disk telemetry access block") {
		t.Errorf("expected telemetry access error, got: %v", errStat)
	}

	// Reseta para os testes de permissão conseguirem avançar após o os.Stat legítimo
	xclifn.ResetDiskBridges()

	// 2. Cobrir a falta de permissão de LEITURA ("missing read permission access privileges")
	readTmpFile := filepath.Join(t.TempDir(), "read_test.txt")
	_ = os.WriteFile(readTmpFile, []byte(""), 0666)

	xclifn.SetDiskBridgesMocks(
		nil,
		func(name string) (*os.File, error) {
			return nil, errors.New("permission denied simulation")
		},
		nil, nil,
	)
	errOpen := xclifn.ValidateDiskResources("--cfg", readTmpFile, specFile)
	if errOpen == nil || !strings.Contains(errOpen.Error(), "missing read permission access privileges") {
		t.Errorf("expected read permission error, got: %v", errOpen)
	}

	xclifn.ResetDiskBridges()

	// 3. Cobrir a falta de permissão de ESCRITA EM DIRETÓRIO ("missing write permission access privileges" para dirpath)
	xclifn.SetDiskBridgesMocks(
		nil, nil,
		func(name string) (*os.File, error) {
			return nil, errors.New("write permission denied simulation")
		},
		nil,
	)
	errCreate := xclifn.ValidateDiskResources("--src", t.TempDir(), specDir)
	if errCreate == nil || !strings.Contains(errCreate.Error(), "missing write permission access privileges") {
		t.Errorf("expected directory write permission error, got: %v", errCreate)
	}

	xclifn.ResetDiskBridges()

	// 4. Cobrir a falta de permissão de ESCRITA EM ARQUIVO ("missing write permission access privileges" para filepath)
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	_ = os.WriteFile(tmpFile, []byte(""), 0666)

	xclifn.SetDiskBridgesMocks(
		nil, nil, nil,
		func(name string, flag int, perm os.FileMode) (*os.File, error) {
			return nil, errors.New("file write permission denied simulation")
		},
	)
	errOpenFile := xclifn.ValidateDiskResources("--cfg", tmpFile, specFile)
	if errOpenFile == nil || !strings.Contains(errOpenFile.Error(), "missing write permission access privileges") {
		t.Errorf("expected file write permission error, got: %v", errOpenFile)
	}
}
