package main

import (
	"regexp"
	"strings"
)

/*
  ARCHITECTURE & SCOPE LIMITATION:
  03_functions.go groups stateless, decoupled utility behaviors and pure computational
  routines required to back up the central proposal of the package.

  Design Constraints:
  - Every function placed here should ideally be deterministic (same input produces same output).
  - No global state mutation or complex side-effects are allowed within these routines.
  - If routines grow complex or introduce stateful context, split them into dedicated files.
*/

// Insert standalone functions or mathematical algorithms below.

func toSnakeUpperCase(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	var result strings.Builder
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		if i > 0 && (runes[i] >= 'A' && runes[i] <= 'Z') && !(runes[i-1] >= 'A' && runes[i-1] <= 'Z') && runes[i-1] != '_' {
			result.WriteRune('_')
		}
		result.WriteRune(runes[i])
	}
	return strings.ToUpper(result.String())
}

func validateToken(token string) bool {
	if len(token) == 0 || len(token) > 32 {
		return false
	}
	matched, _ := regexp.MatchString(`^[A-Z][A-Z0-9_]*$`, token)
	return matched
}
