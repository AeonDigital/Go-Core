package xclistruc_test

import (
	"testing"
	"time"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclistruc"
)

// TestFlagValuesLifecycle enforces full verification across the entire reading API
// processing valid entries, incorrect type assertions, and missing key triggers.
func TestFlagValuesLifecycle(t *testing.T) {
	ctx := xclistruc.NewFlagValues()
	now := time.Now()

	// ------------------------------------------------------------------------
	// SETUP DATA INJECTION
	// ------------------------------------------------------------------------

	// Primitive Types Data
	ctx.SetInternalValue("v_string", "GoCLI")
	ctx.SetInternalValue("v_int", 42)
	ctx.SetInternalValue("v_float", 3.14)
	ctx.SetInternalValue("v_bool", true)

	// Structured Data
	ctx.SetInternalValue("v_json", `{"status":"ok"}`)
	ctx.SetInternalValue("v_duration", 5*time.Minute)
	ctx.SetInternalValue("v_date", now)
	ctx.SetInternalValue("v_time", now)
	ctx.SetInternalValue("v_datetime", now)
	ctx.SetInternalValue("v_email", "test@aeon.digital")

	// System and Network Validations
	ctx.SetInternalValue("v_path", "/usr/local")
	ctx.SetInternalValue("v_filename", "main.go")
	ctx.SetInternalValue("v_filepath", "/usr/local/main.go")
	ctx.SetInternalValue("v_dirname", "bin")
	ctx.SetInternalValue("v_dirpath", "/usr/local/bin")
	ctx.SetInternalValue("v_url", "https://aeon.digital")
	ctx.SetInternalValue("v_fullurl", "https://aeon.digital")
	ctx.SetInternalValue("v_relativeurl", "/api/v1")

	// Array / Slice Formats
	ctx.SetInternalValue("v_string_arr", []string{"A", "B"})
	ctx.SetInternalValue("v_int_arr", []int{1, 2})
	ctx.SetInternalValue("v_float_arr", []float64{1.1, 2.2})
	ctx.SetInternalValue("v_bool_arr", []bool{true, false})
	ctx.SetInternalValue("v_duration_arr", []time.Duration{1 * time.Second})
	ctx.SetInternalValue("v_date_arr", []time.Time{now})
	ctx.SetInternalValue("v_time_arr", []time.Time{now})
	ctx.SetInternalValue("v_datetime_arr", []time.Time{now})
	ctx.SetInternalValue("v_email_arr", []string{"a@a.com"})
	ctx.SetInternalValue("v_path_arr", []string{"/path"})
	ctx.SetInternalValue("v_filename_arr", []string{"file.txt"})
	ctx.SetInternalValue("v_filepath_arr", []string{"/path/file.txt"})
	ctx.SetInternalValue("v_dirname_arr", []string{"dir"})
	ctx.SetInternalValue("v_dirpath_arr", []string{"/path/dir"})
	ctx.SetInternalValue("v_url_arr", []string{"http://a.com"})
	ctx.SetInternalValue("v_fullurl_arr", []string{"http://a.com"})
	ctx.SetInternalValue("v_relativeurl_arr", []string{"/b"})

	// Defensive tracking validation paths
	if !ctx.Has("v_string") {
		t.Errorf("expected key 'v_string' to exist inside the context map")
	}
	if ctx.Has("missing_key") {
		t.Errorf("unexpected track: 'missing_key' should not exist")
	}

	// ============================================================================
	// 1. PRIMITIVE TYPES VERIFICATION TRACKS
	// ============================================================================

	// GetString
	if ctx.GetString("v_string") != "GoCLI" {
		t.Errorf("expected 'GoCLI', got '%s'", ctx.GetString("v_string"))
	}
	if ctx.GetString("missing_key") != "" {
		t.Errorf("absent string key must return empty string")
	}
	if ctx.GetString("v_int") != "" {
		t.Errorf("invalid type assertion for string must fallback to empty string")
	}

	// GetInt
	if ctx.GetInt("v_int") != 42 {
		t.Errorf("expected 42, got %d", ctx.GetInt("v_int"))
	}
	if ctx.GetInt("missing_key") != 0 {
		t.Errorf("absent integer key must return 0")
	}
	if ctx.GetInt("v_string") != 0 {
		t.Errorf("invalid type assertion for int must fallback to 0")
	}

	// GetFloat
	if ctx.GetFloat("v_float") != 3.14 {
		t.Errorf("expected 3.14, got %f", ctx.GetFloat("v_float"))
	}
	if ctx.GetFloat("missing_key") != 0.0 {
		t.Errorf("absent float key must return 0.0")
	}
	if ctx.GetFloat("v_string") != 0.0 {
		t.Errorf("invalid type assertion for float must fallback to 0.0")
	}

	// GetBool
	if ctx.GetBool("v_bool") != true {
		t.Errorf("expected true, got %t", ctx.GetBool("v_bool"))
	}
	if ctx.GetBool("missing_key") != false {
		t.Errorf("absent boolean key must return false")
	}
	if ctx.GetBool("v_string") != false {
		t.Errorf("invalid type assertion for bool must fallback to false")
	}

	// ============================================================================
	// 2. STRUCTURED DATA VERIFICATION TRACKS
	// ============================================================================

	// GetJSON
	if ctx.GetJSON("v_json") != `{"status":"ok"}` {
		t.Errorf("expected JSON string, got '%s'", ctx.GetJSON("v_json"))
	}
	if ctx.GetJSON("missing_key") != "" {
		t.Errorf("absent json key must return empty string")
	}
	if ctx.GetJSON("v_int") != "" {
		t.Errorf("invalid type assertion for json must fallback to empty string")
	}

	// GetDuration
	if ctx.GetDuration("v_duration") != 5*time.Minute {
		t.Errorf("expected 5m, got %v", ctx.GetDuration("v_duration"))
	}
	if ctx.GetDuration("missing_key") != 0 {
		t.Errorf("absent duration key must return 0")
	}
	if ctx.GetDuration("v_string") != 0 {
		t.Errorf("invalid type assertion for duration must fallback to 0")
	}

	// GetDate
	if ctx.GetDate("v_date") != now {
		t.Errorf("expected time %v, got %v", now, ctx.GetDate("v_date"))
	}
	if !ctx.GetDate("missing_key").IsZero() {
		t.Errorf("absent date key must return a zero time object")
	}
	if !ctx.GetDate("v_string").IsZero() {
		t.Errorf("invalid type assertion for date must fallback to a zero time object")
	}

	// GetTime
	if ctx.GetTime("v_time") != now {
		t.Errorf("expected time %v, got %v", now, ctx.GetTime("v_time"))
	}
	if !ctx.GetTime("missing_key").IsZero() {
		t.Errorf("absent time key must return a zero time object")
	}
	if !ctx.GetTime("v_string").IsZero() {
		t.Errorf("invalid type assertion for time must fallback to a zero time object")
	}

	// GetDateTime
	if ctx.GetDateTime("v_datetime") != now {
		t.Errorf("expected time %v, got %v", now, ctx.GetDateTime("v_datetime"))
	}
	if !ctx.GetDateTime("missing_key").IsZero() {
		t.Errorf("absent datetime key must return a zero time object")
	}
	if !ctx.GetDateTime("v_string").IsZero() {
		t.Errorf("invalid type assertion for datetime must fallback to a zero time object")
	}

	// GetEmail
	if ctx.GetEmail("v_email") != "test@aeon.digital" {
		t.Errorf("expected email, got '%s'", ctx.GetEmail("v_email"))
	}
	if ctx.GetEmail("missing_key") != "" {
		t.Errorf("absent email key must return empty string")
	}
	if ctx.GetEmail("v_int") != "" {
		t.Errorf("invalid type assertion for email must fallback to empty string")
	}

	// ============================================================================
	// 3. SYSTEM AND NETWORK VALIDATIONS VERIFICATION TRACKS
	// ============================================================================

	// GetPath
	if ctx.GetPath("v_path") != "/usr/local" {
		t.Errorf("error on GetPath")
	}
	if ctx.GetPath("missing_key") != "" {
		t.Errorf("error on GetPath empty")
	}
	if ctx.GetPath("v_int") != "" {
		t.Errorf("error on GetPath fallback")
	}

	// GetFilename
	if ctx.GetFilename("v_filename") != "main.go" {
		t.Errorf("error on GetFilename")
	}
	if ctx.GetFilename("missing_key") != "" {
		t.Errorf("error on GetFilename empty")
	}
	if ctx.GetFilename("v_int") != "" {
		t.Errorf("error on GetFilename fallback")
	}

	// GetFilepath
	if ctx.GetFilepath("v_filepath") != "/usr/local/main.go" {
		t.Errorf("error on GetFilepath")
	}
	if ctx.GetFilepath("missing_key") != "" {
		t.Errorf("error on GetFilepath empty")
	}
	if ctx.GetFilepath("v_int") != "" {
		t.Errorf("error on GetFilepath fallback")
	}

	// GetDirname
	if ctx.GetDirname("v_dirname") != "bin" {
		t.Errorf("error on GetDirname")
	}
	if ctx.GetDirname("missing_key") != "" {
		t.Errorf("error on GetDirname empty")
	}
	if ctx.GetDirname("v_int") != "" {
		t.Errorf("error on GetDirname fallback")
	}

	// GetDirpath
	if ctx.GetDirpath("v_dirpath") != "/usr/local/bin" {
		t.Errorf("error on GetDirpath")
	}
	if ctx.GetDirpath("missing_key") != "" {
		t.Errorf("error on GetDirpath empty")
	}
	if ctx.GetDirpath("v_int") != "" {
		t.Errorf("error on GetDirpath fallback")
	}

	// GetURL
	if ctx.GetURL("v_url") != "https://aeon.digital" {
		t.Errorf("error on GetURL")
	}
	if ctx.GetURL("missing_key") != "" {
		t.Errorf("error on GetURL empty")
	}
	if ctx.GetURL("v_int") != "" {
		t.Errorf("error on GetURL fallback")
	}

	// GetFullURL
	if ctx.GetFullURL("v_fullurl") != "https://aeon.digital" {
		t.Errorf("error on GetFullURL")
	}
	if ctx.GetFullURL("missing_key") != "" {
		t.Errorf("error on GetFullURL empty")
	}
	if ctx.GetFullURL("v_int") != "" {
		t.Errorf("error on GetFullURL fallback")
	}

	// GetRelativeURL
	if ctx.GetRelativeURL("v_relativeurl") != "/api/v1" {
		t.Errorf("error on GetRelativeURL")
	}
	if ctx.GetRelativeURL("missing_key") != "" {
		t.Errorf("error on GetRelativeURL empty")
	}
	if ctx.GetRelativeURL("v_int") != "" {
		t.Errorf("error on GetRelativeURL fallback")
	}

	// ============================================================================
	// 4. ARRAY / SLICE FORMATS VERIFICATION TRACKS
	// ============================================================================

	// --- Slice of Primitive Types ---

	// GetStringSlice
	if len(ctx.GetStringSlice("v_string_arr")) != 2 {
		t.Errorf("error on GetStringSlice")
	}
	if ctx.GetStringSlice("missing_key") != nil {
		t.Errorf("error on GetStringSlice empty")
	}
	if ctx.GetStringSlice("v_int") != nil {
		t.Errorf("error on GetStringSlice fallback")
	}

	// GetIntSlice
	if len(ctx.GetIntSlice("v_int_arr")) != 2 {
		t.Errorf("error on GetIntSlice")
	}
	if ctx.GetIntSlice("missing_key") != nil {
		t.Errorf("error on GetIntSlice empty")
	}
	if ctx.GetIntSlice("v_string") != nil {
		t.Errorf("error on GetIntSlice fallback")
	}

	// GetFloatSlice
	if len(ctx.GetFloatSlice("v_float_arr")) != 2 {
		t.Errorf("error on GetFloatSlice")
	}
	if ctx.GetFloatSlice("missing_key") != nil {
		t.Errorf("error on GetFloatSlice empty")
	}
	if ctx.GetFloatSlice("v_string") != nil {
		t.Errorf("error on GetFloatSlice fallback")
	}

	// GetBoolSlice
	if len(ctx.GetBoolSlice("v_bool_arr")) != 2 {
		t.Errorf("error on GetBoolSlice")
	}
	if ctx.GetBoolSlice("missing_key") != nil {
		t.Errorf("error on GetBoolSlice empty")
	}
	if ctx.GetBoolSlice("v_string") != nil {
		t.Errorf("error on GetBoolSlice fallback")
	}

	// --- Slice of Structured Data ---

	// GetDurationSlice
	if len(ctx.GetDurationSlice("v_duration_arr")) != 1 {
		t.Errorf("error on GetDurationSlice")
	}
	if ctx.GetDurationSlice("missing_key") != nil {
		t.Errorf("error on GetDurationSlice empty")
	}
	if ctx.GetDurationSlice("v_string") != nil {
		t.Errorf("error on GetDurationSlice fallback")
	}

	// GetDateSlice
	if len(ctx.GetDateSlice("v_date_arr")) != 1 {
		t.Errorf("error on GetDateSlice")
	}
	if ctx.GetDateSlice("missing_key") != nil {
		t.Errorf("error on GetDateSlice empty")
	}
	if ctx.GetDateSlice("v_string") != nil {
		t.Errorf("error on GetDateSlice fallback")
	}

	// GetTimeSlice
	if len(ctx.GetTimeSlice("v_time_arr")) != 1 {
		t.Errorf("error on GetTimeSlice")
	}
	if ctx.GetTimeSlice("missing_key") != nil {
		t.Errorf("error on GetTimeSlice empty")
	}
	if ctx.GetTimeSlice("v_string") != nil {
		t.Errorf("error on GetTimeSlice fallback")
	}

	// GetDateTimeSlice
	if len(ctx.GetDateTimeSlice("v_datetime_arr")) != 1 {
		t.Errorf("error on GetDateTimeSlice")
	}
	if ctx.GetDateTimeSlice("missing_key") != nil {
		t.Errorf("error on GetDateTimeSlice empty")
	}
	if ctx.GetDateTimeSlice("v_string") != nil {
		t.Errorf("error on GetDateTimeSlice fallback")
	}

	// GetEmailSlice
	if len(ctx.GetEmailSlice("v_email_arr")) != 1 {
		t.Errorf("error on GetEmailSlice")
	}
	if ctx.GetEmailSlice("missing_key") != nil {
		t.Errorf("error on GetEmailSlice empty")
	}
	if ctx.GetEmailSlice("v_int") != nil {
		t.Errorf("error on GetEmailSlice fallback")
	}

	// --- Slice of System and Network Validations ---

	// GetPathSlice
	if len(ctx.GetPathSlice("v_path_arr")) != 1 {
		t.Errorf("error on GetPathSlice")
	}
	if ctx.GetPathSlice("missing_key") != nil {
		t.Errorf("error on GetPathSlice empty")
	}
	if ctx.GetPathSlice("v_int") != nil {
		t.Errorf("error on GetPathSlice fallback")
	}

	// GetFilenameSlice
	if len(ctx.GetFilenameSlice("v_filename_arr")) != 1 {
		t.Errorf("error on GetFilenameSlice")
	}
	if ctx.GetFilenameSlice("missing_key") != nil {
		t.Errorf("error on GetFilenameSlice empty")
	}
	if ctx.GetFilenameSlice("v_int") != nil {
		t.Errorf("error on GetFilenameSlice fallback")
	}

	// GetFilepathSlice
	if len(ctx.GetFilepathSlice("v_filepath_arr")) != 1 {
		t.Errorf("error on GetFilepathSlice")
	}
	if ctx.GetFilepathSlice("missing_key") != nil {
		t.Errorf("error on GetFilepathSlice empty")
	}
	if ctx.GetFilepathSlice("v_int") != nil {
		t.Errorf("error on GetFilepathSlice fallback")
	}

	// GetDirnameSlice
	if len(ctx.GetDirnameSlice("v_dirname_arr")) != 1 {
		t.Errorf("error on GetDirnameSlice")
	}
	if ctx.GetDirnameSlice("missing_key") != nil {
		t.Errorf("error on GetDirnameSlice empty")
	}
	if ctx.GetDirnameSlice("v_int") != nil {
		t.Errorf("error on GetDirnameSlice fallback")
	}

	// GetDirpathSlice
	if len(ctx.GetDirpathSlice("v_dirpath_arr")) != 1 {
		t.Errorf("error on GetDirpathSlice")
	}
	if ctx.GetDirpathSlice("missing_key") != nil {
		t.Errorf("error on GetDirpathSlice empty")
	}
	if ctx.GetDirpathSlice("v_int") != nil {
		t.Errorf("error on GetDirpathSlice fallback")
	}

	// GetURLSlice
	if len(ctx.GetURLSlice("v_url_arr")) != 1 {
		t.Errorf("error on GetURLSlice")
	}
	if ctx.GetURLSlice("missing_key") != nil {
		t.Errorf("error on GetURLSlice empty")
	}
	if ctx.GetURLSlice("v_int") != nil {
		t.Errorf("error on GetURLSlice fallback")
	}

	// GetFullURLSlice
	if len(ctx.GetFullURLSlice("v_fullurl_arr")) != 1 {
		t.Errorf("error on GetFullURLSlice")
	}
	if ctx.GetFullURLSlice("missing_key") != nil {
		t.Errorf("error on GetFullURLSlice empty")
	}
	if ctx.GetFullURLSlice("v_int") != nil {
		t.Errorf("error on GetFullURLSlice fallback")
	}

	// GetRelativeURLSlice
	if len(ctx.GetRelativeURLSlice("v_relativeurl_arr")) != 1 {
		t.Errorf("error on GetRelativeURLSlice")
	}
	if ctx.GetRelativeURLSlice("missing_key") != nil {
		t.Errorf("error on GetRelativeURLSlice empty")
	}
	if ctx.GetRelativeURLSlice("v_int") != nil {
		t.Errorf("error on GetRelativeURLSlice fallback")
	}
}

// TestSetInternalValueWithNilMap ensures that if the internal values map
// is uninitialized (nil), the engine safely instantiates it before assignment.
func TestSetInternalValueWithNilMap(t *testing.T) {
	ctx := &xclistruc.FlagValues{}
	ctx.SetInternalValue("lazy_init", "success")

	if !ctx.Has("lazy_init") {
		t.Errorf("expected map to be instantiated and hold the value")
	}
	if ctx.GetString("lazy_init") != "success" {
		t.Errorf("expected to retrieve 'success', got '%s'", ctx.GetString("lazy_init"))
	}
}
