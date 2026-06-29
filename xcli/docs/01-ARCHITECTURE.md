ARCHITECTURE.md
================================

## 1. OVERVIEW & PHILOSOPHY

This document establishes the structural, architectural, and documentation guidelines for Go modules within this ecosystem.
Our core development philosophy balances simplicity with rigid code boundaries. We value explicit design over implicit behaviors, ensuring zero black boxes. Every component must have a predictable location, clear constraints, and maximum human readability.



&nbsp;
&nbsp;


________________________________________________________________________________

## 2. REPOSITORY STRUCTURE

Replace the `pkgname/` root directory with your module name, and substitute the `<pkg>` placeholder with the principal shorthand name of your library to enforce context-prefixed public packages.

When a public package mirrors the core purpose of an internal one (e.g., `<pkg>cfg` vs `internal/config`), the public layer must act exclusively as an exposed gateway or contract, handling data primitives or public structures, while delegating heavy orchestration to its private counterpart.

``` 
mainrepo/
└── pkgname/
    ├── docs/                   
    │   ├── 00-MD_RULES.md      
    │   └── 01-ARCHITECTURE.md  
    │
    ├── internal/               # PRIVATE: Hidden core, zero external imports allowed
    │   ├── config/             # Heavy loaders (flags, envs, file parsers)
    │   ├── contt/              # Internal immutable constants and private magic numbers
    │   ├── fn/                 # Core stateless algorithms and computational math
    │   ├── intfc/              # Private decoupled contracts
    │   ├── struc/              # Domain models, strict validations, and stateful structures
    │   └── xerrors/            # Private named sentinel errors
    │
    └── pkg/                    # PUBLIC: Exportable gateway, exposes primitives and contracts
        └── <pkgname>/          
            ├── <pkg>cfg/       # Public configuration structures and primitives
            ├── <pkg>constt/    # Public type definitions, enums, and validation constants
            ├── <pkg>fn/        # Public stateless utility wrappers (Print, Formatter)
            ├── <pkg>intfc/     # Public behavioral interfaces for external extendability
            ├── <pkg>struc/     # Public DTOs and data models needed by external callers
            └── <pkg>errors/    # Publicly interceptable named sentinel errors
```


&nbsp;
&nbsp;


________________________________________________________________________________

## 3. COMPONENT DIRECTORY RULES

### /internal

Contains the private core of the application or library. The Go compiler strictly forbids any external micro-project or module from importing code residing here.

* internal/config/ (package config)
  Entry point for parsing flags, environment variables, or files. No business logic allowed.
* internal/contt/ (package contt)
  Central registry for immutable primitive values. Mutable global states or pointers are strictly banned.
* internal/fn/ (package fn)
  Hosts stateless, deterministic computational routines and mathematical algorithms.
* internal/intfc/ (package intfc)
  Defines decoupled behavioral contracts. Interfaces must be small (1 to 3 methods).
* internal/struc/ (package struc)
  Hosts pure data models and DTOs. Methods are limited to basic validation or getters/setters.
* internal/xerrors/ (package xerrors)
  Central registry for named sentinel errors (e.g., ErrPermissionDenied).


&nbsp;


### /pkg

Contains the public, exportable interface of the library or application. Other micro-projects within the monorepo can freely import packages from this directory.

* Architectural Symmetry Rule: The 6 core functional layers (config, constt, fn, intfc, struc, xerrors) can exist both in /internal (for private orchestration) and /pkg (for public exposure).
* Context-Prefixed Naming Constraint: To prevent package name collisions, any of the 6 core layers exposed in /pkg must be prefixed with the shorthand name of the project (e.g., <pkg>cfg, <pkg>constt, <pkg>fn, <pkg>intfc, <pkg>struc, <pkg>errors).
* Public Constants & Errors Rule: Layers like <pkg>constt and <pkg>errors are explicitly designed to expose types, sentinel errors, and validation primitives that external consumers or package objects need to interact with the library's API.


&nbsp;
&nbsp;


________________________________________________________________________________

## 4. DOCUMENTATION STANDARD (HUMAN READABILITY FIRST)

### 4.1 General Constraints

* Language: 100% technical English for all code documentation, comments, and structure naming.
* Format: Avoid dense, blocky, inline comments. Use generous visual spacing, line breaks, and clear indentation so developers can scan the file effortlessly.


&nbsp;


### 4.2 Code Coverage Requirements

* Functions & Methods: 100% Mandatory documentation coverage. The only exception is trivial, obvious Getters/Setters with zero logic.
* Structs & Fields: Highly Preferred. Omit comments only if the field is standard and completely self-explanatory (e.g., ID string, CreatedAt time.Time). If any business rule applies, document it.


&nbsp;


### 4.3 Function Comment Anatomy

Every function documentation must follow a strict tiered format:

   1. Line 1 (Summary): A brief, single-line summary of what the function accomplishes.
   2. Line 2: A mandatory empty line (visual breathing space).
   3. Detailed View: Technical explanation detailing arguments, return values, special side-effects, and operational constraints.


&nbsp;
&nbsp;


________________________________________________________________________________

## Error & Panic Documentation Criteria:

* Panics: All potential application-halting execution paths (panic) must be explicitly declared.
* Return Errors (Contextual): If a failure is obvious (e.g., simple validation), no further explanation is needed. If a function is complex and can fail due to multiple distinct natures (e.g., network timeout vs. data corruption), those failure natures must be explicitly itemized.


&nbsp;
&nbsp;


________________________________________________________________________________

## 5. CODE EXAMPLE (THE GOLD STANDARD)

``` example.go
package fn

// FormatBytes converts an integer byte count into a string.
//
// Arguments:
//   - bytes: The raw payload size in bytes. Must be a positive integer.
//
// Returns:
//   - string: The formatted value (e.g., "1.5 GB", "42 MB").
//
// Error & Panic Natures:
//   - Panics: Will trigger a panic if the input bytes argument is a negative value.
//   - Complex Errors: Returns an empty string if the internal system math precision 
//     overflows during calculation.
func FormatBytes(bytes int64) string {
    if bytes < 0 {
        panic("fn: bytes argument cannot be negative")
    }
    // Implementation goes here
    return ""
}
```


&nbsp;
&nbsp;


________________________________________________________________________________

## 6. VERSIONING POLICY

This project and all libraries inside the monorepo strictly enforce Semantic Versioning (SemVer) via Git tags (vX.Y.Z) to ensure predictability and prevent breaking downstream applications.
