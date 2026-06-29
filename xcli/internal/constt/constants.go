package constt

/*
	ARCHITECTURE & SCOPE LIMITATION:
	constants.go centralizes immutable, read-only values and global literals
	exclusively required to configure or support the package logic.

	Design Constraints:
	- Only truly constant and stateless values (string, int, time.Duration primitives)
		are allowed here.
	- Never declare mutable global variables or pointers inside this package.
	- If a distinct domain domain context (e.g., custom error codes or CLI flag keys)
		grows large, split those constants into a separate, context-named file here.
*/

// Insert global constants below.

//
//
//
