package xjson

import (
	"encoding/json"
)

// Dump is a debugging utility that marshals the given value into an indented JSON string.
// It accepts any value (v) to be serialized and an indentation string (i), which defaults to two spaces ("  ") if empty.
// Since it is intended for debugging, it swallows marshal errors and returns the error message string instead.
func Dump(v any, i string) string {
	if i == "" {
		i = "  "
	}

	bytes, err := json.MarshalIndent(v, "", i)
	if err != nil {
		return "ERROR:" + err.Error()
	}
	return string(bytes)
}
