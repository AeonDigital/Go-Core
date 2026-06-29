package xclifn_test

import (
	"testing"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
)

// TestParseStringSuccess verifies text trimming standardization tracks.
func TestParseStringSuccess(t *testing.T) {
	res, err := xclifn.ParseString("  aeon digital  ")
	if err != nil {
		t.Errorf("expected string parser to never fail, got: %v", err)
	}
	if res != "aeon digital" {
		t.Errorf("expected trimmed text layout 'aeon digital', got: %q", res)
	}
}

// TestParsePrimitivosSuccessAndErrors validates type coercion and syntax checking
// for boolean, integer, and floating-point primitive parser functions.
func TestParsePrimitivosSuccessAndErrors(t *testing.T) {
	// 1. Validate ParseBool variations
	boolCases := []struct {
		input    string
		expected bool
		hasErr   bool
	}{
		{"true", true, false},
		{"TRUE", true, false},
		{"1", true, false},
		{"false", false, false},
		{"0", false, false},
		{"invalid_bool", false, true},
	}
	for _, tc := range boolCases {
		res, err := xclifn.ParseBool(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseBool(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && (err != nil || res != tc.expected) {
			t.Errorf("ParseBool(%q) expected %t with no error, got %t (err: %v)", tc.input, tc.expected, res, err)
		}
	}

	// 2. Validate ParseInt variations
	intCases := []struct {
		input    string
		expected int
		hasErr   bool
	}{
		{" 42  ", 42, false},
		{"-100", -100, false},
		{"not_a_number", 0, true},
	}
	for _, tc := range intCases {
		res, err := xclifn.ParseInt(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseInt(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && (err != nil || res != tc.expected) {
			t.Errorf("ParseInt(%q) expected %d, got %d", tc.input, tc.expected, res)
		}
	}

	// 3. Validate ParseFloat variations
	floatCases := []struct {
		input    string
		expected float64
		hasErr   bool
	}{
		{" 3.1415  ", 3.1415, false},
		{"-0.005", -0.005, false},
		{"broken_float", 0.0, true},
	}
	for _, tc := range floatCases {
		res, err := xclifn.ParseFloat(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseFloat(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && (err != nil || res != tc.expected) {
			t.Errorf("ParseFloat(%q) expected %f, got %f", tc.input, tc.expected, res)
		}
	}
}

// TestParseEstruturadosSuccessAndErrors validates validation logic for email addresses
// and standardized JSON payload syntax constraints.
func TestParseEstruturadosSuccessAndErrors(t *testing.T) {
	// 1. Validate ParseEmail variations
	emailCases := []struct {
		input    string
		expected string
		hasErr   bool
	}{
		{"  dev@aeondigital.com  ", "dev@aeondigital.com", false},
		{"invalid-email-address", "", true},
		{"John Doe <john@example.com>", "", true}, // Rejects loose naming annotations
	}
	for _, tc := range emailCases {
		res, err := xclifn.ParseEmail(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseEmail(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && (err != nil || res != tc.expected) {
			t.Errorf("ParseEmail(%q) expected %q, got %q (err: %v)", tc.input, tc.expected, res, err)
		}
	}

	// 2. Validate ParseJSON variations
	jsonCases := []struct {
		input    string
		expected string
		hasErr   bool
	}{
		{`{"valid": true, "id": 1}`, `{"valid": true, "id": 1}`, false},
		{`   [1, 2, "three"]   `, `[1, 2, "three"]`, false},
		{`{"broken_json": `, "", true},
	}
	for _, tc := range jsonCases {
		res, err := xclifn.ParseJSON(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseJSON(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && (err != nil || res != tc.expected) {
			t.Errorf("ParseJSON(%q) expected %q, got %q", tc.input, tc.expected, res)
		}
	}
}

// TestParseTemporaisSuccessAndErrors validates date, duration, and flexible clock format
// parsing logic, testing both short and full time syntax options.
func TestParseTemporaisSuccessAndErrors(t *testing.T) {
	// 1. Validate ParseDuration variations
	durationCases := []struct {
		input  string
		hasErr bool
	}{
		{"  5m30s ", false},
		{"2h", false},
		{"invalid_duration", true},
	}
	for _, tc := range durationCases {
		_, err := xclifn.ParseDuration(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseDuration(%q) expected error, got nil", tc.input)
		}
	}

	// 2. Validate ParseDate variations
	dateCases := []struct {
		input  string
		hasErr bool
	}{
		{"2026-06-28", false},
		{"28-06-2026", true}, // Wrong order layout
		{"broken_date", true},
	}
	for _, tc := range dateCases {
		_, err := xclifn.ParseDate(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseDate(%q) expected error, got nil", tc.input)
		}
	}

	// 3. Validate ParseTime variations (Flexible clock syntax checker)
	timeCases := []struct {
		input  string
		hasErr bool
	}{
		{"14:30:15", false}, // Full layout (HH:MM:SS)
		{"09:15", false},    // Short layout (HH:MM)
		{"25:00", true},     // Invalid hour number bounds
		{"broken_time", true},
	}
	for _, tc := range timeCases {
		_, err := xclifn.ParseTime(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseTime(%q) expected error, got nil", tc.input)
		}
	}

	// 4. Validate ParseDateTime variations
	dateTimeCases := []struct {
		input  string
		hasErr bool
	}{
		{"2026-06-28 14:30:00", false},
		{"2026-06-28T14:30:00", true}, // Missing exact space separation token
		{"broken_datetime", true},
	}
	for _, tc := range dateTimeCases {
		_, err := xclifn.ParseDateTime(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseDateTime(%q) expected error, got nil", tc.input)
		}
	}
}

// TestParseFileSystemSemanticSuccessAndErrors executes syntax validation rules
// for dedicated semantic path functions (Filename, Dirname, and Path).
func TestParseFileSystemSemanticSuccessAndErrors(t *testing.T) {
	// 1. Validate ParseFilename
	filenameCases := []struct {
		input  string
		hasErr bool
	}{
		{"document.pdf", false},
		{"nested/folder/file.pdf", true}, // Filename must not contain directory slashes
		{"", true},                       // Empty payload check
	}
	for _, tc := range filenameCases {
		_, err := xclifn.ParseFilename(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseFilename(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseFilename(%q) unexpected error: %v", tc.input, err)
		}
	}

	// 2. Validate ParseDirname
	dirnameCases := []struct {
		input  string
		hasErr bool
	}{
		{"/var/log", false},
		{"subfolder", false},
		{"", true}, // Empty payload check
	}
	for _, tc := range dirnameCases {
		_, err := xclifn.ParseDirname(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseDirname(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseDirname(%q) unexpected error: %v", tc.input, err)
		}
	}

	// 3. Validate ParsePath
	pathCases := []struct {
		input  string
		hasErr bool
	}{
		{"reports/file.txt", false},
		{"clean/path/to/resource", false},
		{"", true}, // Empty payload check
	}
	for _, tc := range pathCases {
		_, err := xclifn.ParsePath(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParsePath(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParsePath(%q) unexpected error: %v", tc.input, err)
		}
	}
}

// TestParseURLVariantsSuccessAndErrors executes syntax evaluation criteria
// across all independent network URL function variations.
func TestParseURLVariantsSuccessAndErrors(t *testing.T) {
	// 1. Validate ParseURLStandard
	urlCases := []struct {
		input  string
		hasErr bool
	}{
		{"https://aeondigital.com", false},
		{"/v1/users?id=1", false},
		{"%%broken-url%%", true}, // Forces url.Parse failure
	}
	for _, tc := range urlCases {
		_, err := xclifn.ParseURLStandard(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseURLStandard(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseURLStandard(%q) unexpected error: %v", tc.input, err)
		}
	}

	// 2. Validate ParseFullURL
	fullURLCases := []struct {
		input  string
		hasErr bool
	}{
		{"https://aeondigital.com", false},
		{"aeondigital.com", true}, // Missing mandatory protocol scheme
		{"%%broken-url%%", true},  // Forces url.Parse failure
	}
	for _, tc := range fullURLCases {
		_, err := xclifn.ParseFullURL(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseFullURL(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseFullURL(%q) unexpected error: %v", tc.input, err)
		}
	}

	// 3. Validate ParseRelativeURL
	relativeURLCases := []struct {
		input  string
		hasErr bool
	}{
		{"/api/v1/health", false},
		{"https://aeondigital.com", true}, // Scheme/Host must be absent
		{"%%broken-url%%", true},          // Forces url.Parse failure
	}
	for _, tc := range relativeURLCases {
		_, err := xclifn.ParseRelativeURL(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseRelativeURL(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseRelativeURL(%q) unexpected error: %v", tc.input, err)
		}
	}
}

// TestParsePhysicalDiskInputsSuccessAndErrors executes baseline normalization
// evaluation bounds for physical path inputs (Filepath and Dirpath).
func TestParsePhysicalDiskInputsSuccessAndErrors(t *testing.T) {
	// 1. Validate ParseFilepath
	filepathCases := []struct {
		input  string
		hasErr bool
	}{
		{"/etc/config.json", false},
		{"relative/file.txt", false},
		{"", true}, // Empty payload check
	}
	for _, tc := range filepathCases {
		_, err := xclifn.ParseFilepath(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseFilepath(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseFilepath(%q) unexpected error: %v", tc.input, err)
		}
	}

	// 2. Validate ParseDirpath
	dirpathCases := []struct {
		input  string
		hasErr bool
	}{
		{"/usr/local/bin", false},
		{"src/pkg", false},
		{"", true}, // Empty payload check
	}
	for _, tc := range dirpathCases {
		_, err := xclifn.ParseDirpath(tc.input)
		if tc.hasErr && err == nil {
			t.Errorf("ParseDirpath(%q) expected error, got nil", tc.input)
		}
		if !tc.hasErr && err != nil {
			t.Errorf("ParseDirpath(%q) unexpected error: %v", tc.input, err)
		}
	}
}
