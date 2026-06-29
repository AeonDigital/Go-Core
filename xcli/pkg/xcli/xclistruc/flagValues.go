package xclistruc

import "time"

// Context holds the completely parsed, converted, and validated flag values
// for the currently executing command.
type FlagValues struct {
	// values stores the final typed data, indexed by the flag's LongName.
	values map[string]any
}

// NewFlagValues initializes an empty container for validated flag values.
func NewFlagValues() *FlagValues {
	return &FlagValues{
		values: make(map[string]any),
	}
}

// SetInternalValue stores a raw or processed value inside the internal map.
//
// Arguments:
//   - name: The long name identifier of the target flag.
//   - val: The computed or raw value payload to hold.
func (c *FlagValues) SetInternalValue(name string, val any) {
	if c.values == nil {
		c.values = make(map[string]any)
	}
	c.values[name] = val
}

// Has checks if a specific flag was explicitly provided by the user in the terminal.
//
// Useful to differentiate between a flag omitted (using default) and one explicitly sent.
func (c *FlagValues) Has(name string) bool {
	_, exists := c.values[name]
	return exists
}

// ============================================================================
// 1. PRIMITIVE TYPES GETTERS
// ============================================================================

// GetString returns the value of a string flag.
//
// Returns empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetString(name string) string {
	if val, exists := c.values[name]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// GetInt returns the value of an integer flag.
//
// Returns 0 if the flag is absent or is not an integer type.
func (c *FlagValues) GetInt(name string) int {
	if val, exists := c.values[name]; exists {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return 0
}

// GetFloat returns the value of a float flag.
//
// Returns 0.0 if the flag is absent or is not a float64 type.
func (c *FlagValues) GetFloat(name string) float64 {
	if val, exists := c.values[name]; exists {
		if floatVal, ok := val.(float64); ok {
			return floatVal
		}
	}
	return 0.0
}

// GetBool returns the value of a boolean flag.
//
// Returns false if the flag is absent or is not a boolean type.
func (c *FlagValues) GetBool(name string) bool {
	if val, exists := c.values[name]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

// ============================================================================
// 2. STRUCTURED DATA GETTERS
// ============================================================================

// GetJSON returns the validated JSON raw payload text.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetJSON(name string) string {
	return c.GetString(name)
}

// GetDuration returns the value of a duration flag as a time.Duration object.
//
// Returns 0 if the flag is absent or is not a time.Duration type.
func (c *FlagValues) GetDuration(name string) time.Duration {
	if val, exists := c.values[name]; exists {
		if durationVal, ok := val.(time.Duration); ok {
			return durationVal
		}
	}
	return 0
}

// GetDate returns the value of a date flag as a time.Time object.
//
// Returns a zero time.Time object if the flag is absent or is not a time.Time type.
func (c *FlagValues) GetDate(name string) time.Time {
	if val, exists := c.values[name]; exists {
		if timeVal, ok := val.(time.Time); ok {
			return timeVal
		}
	}
	return time.Time{}
}

// GetTime returns the value of a time flag as a time.Time object.
//
// Returns a zero time.Time object if the flag is absent or is not a time.Time type.
func (c *FlagValues) GetTime(name string) time.Time {
	if val, exists := c.values[name]; exists {
		if timeVal, ok := val.(time.Time); ok {
			return timeVal
		}
	}
	return time.Time{}
}

// GetDateTime returns the value of a datetime flag as a time.Time object.
//
// Returns a zero time.Time object if the flag is absent or is not a time.Time type.
func (c *FlagValues) GetDateTime(name string) time.Time {
	if val, exists := c.values[name]; exists {
		if timeVal, ok := val.(time.Time); ok {
			return timeVal
		}
	}
	return time.Time{}
}

// GetEmail returns the value of an email flag.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetEmail(name string) string {
	return c.GetString(name)
}

// ============================================================================
// 3. SYSTEM AND NETWORK VALIDATIONS GETTERS
// ============================================================================

// GetPath returns the cleaned semantic path component layout string.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetPath(name string) string {
	return c.GetString(name)
}

// GetFilename returns the syntactically valid standalone filename string.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetFilename(name string) string {
	return c.GetString(name)
}

// GetFilepath returns the normalized physical file path string.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetFilepath(name string) string {
	return c.GetString(name)
}

// GetDirname returns the syntactically valid directory name layout string.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetDirname(name string) string {
	return c.GetString(name)
}

// GetDirpath returns the normalized physical directory path string.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetDirpath(name string) string {
	return c.GetString(name)
}

// GetURL returns the validated basic URL target string layout.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetURL(name string) string {
	return c.GetString(name)
}

// GetFullURL returns the validated full URL domain locator string.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetFullURL(name string) string {
	return c.GetString(name)
}

// GetRelativeURL returns the validated relative target path string layout.
//
// It returns an empty string if the flag is absent or is not a string type.
func (c *FlagValues) GetRelativeURL(name string) string {
	return c.GetString(name)
}

// ============================================================================
// 4. ARRAY / SLICE FORMATS GETTERS
// ============================================================================

// --- Slice of Primitive Types ---

// GetStringSlice returns the value of a string array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetStringSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetIntSlice returns the value of an integer array flag.
//
// Returns nil if the flag is absent or is not a []int type.
func (c *FlagValues) GetIntSlice(name string) []int {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]int); ok {
			return sliceVal
		}
	}
	return nil
}

// GetFloatSlice returns the value of a float array flag.
//
// Returns nil if the flag is absent or is not a []float64 type.
func (c *FlagValues) GetFloatSlice(name string) []float64 {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]float64); ok {
			return sliceVal
		}
	}
	return nil
}

// GetBoolSlice returns the value of a boolean array flag.
//
// Returns nil if the flag is absent or is not a []bool type.
func (c *FlagValues) GetBoolSlice(name string) []bool {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]bool); ok {
			return sliceVal
		}
	}
	return nil
}

// --- Slice of Structured Data ---

// GetDurationSlice returns the value of a duration array flag.
//
// Returns nil if the flag is absent or is not a []time.Duration type.
func (c *FlagValues) GetDurationSlice(name string) []time.Duration {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]time.Duration); ok {
			return sliceVal
		}
	}
	return nil
}

// GetDateSlice returns the value of a date array flag.
//
// Returns nil if the flag is absent or is not a []time.Time type.
func (c *FlagValues) GetDateSlice(name string) []time.Time {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]time.Time); ok {
			return sliceVal
		}
	}
	return nil
}

// GetTimeSlice returns the value of a time array flag.
//
// Returns nil if the flag is absent or is not a []time.Time type.
func (c *FlagValues) GetTimeSlice(name string) []time.Time {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]time.Time); ok {
			return sliceVal
		}
	}
	return nil
}

// GetDateTimeSlice returns the value of a datetime array flag.
//
// Returns nil if the flag is absent or is not a []time.Time type.
func (c *FlagValues) GetDateTimeSlice(name string) []time.Time {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]time.Time); ok {
			return sliceVal
		}
	}
	return nil
}

// GetEmailSlice returns the value of an email array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetEmailSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// --- Slice of System and Network Validations ---

// GetPathSlice returns the value of a path array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetPathSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetFilenameSlice returns the value of a filename array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetFilenameSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetFilepathSlice returns the value of a filepath array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetFilepathSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetDirnameSlice returns the value of a dirname array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetDirnameSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetDirpathSlice returns the value of a dirpath array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetDirpathSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetURLSlice returns the value of a URL array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetURLSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetFullURLSlice returns the value of a full URL array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetFullURLSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}

// GetRelativeURLSlice returns the value of a relative URL array flag.
//
// Returns nil if the flag is absent or is not a []string type.
func (c *FlagValues) GetRelativeURLSlice(name string) []string {
	if val, exists := c.values[name]; exists {
		if sliceVal, ok := val.([]string); ok {
			return sliceVal
		}
	}
	return nil
}
