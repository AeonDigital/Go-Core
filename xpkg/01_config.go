package main

import "embed"

/*
	ARCHITECTURE & SCOPE LIMITATION:
	00_config.go isolates global runtime configurations, operational switches, and
	thread-safe internal states that dictate the behavioral mechanics of this package.

	Design Constraints:
	- No core domain logic execution or heavy computing should exist here.
	- Any mutable global state MUST utilize atomic operations or synchronization primitives.
*/

// Insert configuration structures, operational constants, or atomic state flags below.

//go:embed templates/*
var TemplateFS embed.FS
