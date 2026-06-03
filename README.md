Go-Core
================================

> [Aeon Digital](http://www.aeondigital.com.br)  
> rianna@aeondigital.com.br

&nbsp;

A modular repository of general-purpose packages and utilities for Go applications. 
Managed as a monorepo, each directory represents an independent module with its own versioning, allowing you to import only what your application needs.



&nbsp;
&nbsp;


________________________________________________________________________________

## Repository Structure

```text
Go-Core/
├── internal/             # Internal testing utilities
├── tools/                # Core debugging and formatting utilities
├── xconfig/              # Configuration management and parsing
├── xerrors/              # Structured error handling and logging
├── xfs/                  # Cross-platform filesystem utilities
├── xio/                  # CLI input/output helpers
├── xjson/                # JSON manipulation utilities
├── xlog/                 # Logging with CLI and registry support
├── xmocks/               # Mock generation tools
├── xreflect/             # Reflection helpers
├── xtime/                # DateTime formatting and utilities
├── xunits/               # Units of magnitude (bytes, durations)
└── README.md             # This file
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Modules

Each package is independent and can be imported separately. Visit each package directory for detailed documentation.

### xconfig

Configuration management with support for multiple formats and sources.

* **Features:** Multi-parser composition, environment variables, `.env`, JSON, YAML support.
* **Import path:** `github.com/AeonDigital/Go-Core/xconfig`

### xerrors

Structured error classification and logging with corporate formatting.

* **Features:** Domain-specific error codes, component tracing, structured debug output, slog integration.
* **Import path:** `github.com/AeonDigital/Go-Core/xerrors`

### xfs

Cross-platform filesystem utilities with testable boundary isolation.

* **Features:** Path resolution, file/directory operations, permissions, standard user directories.
* **Import path:** `github.com/AeonDigital/Go-Core/xfs`

### xio

CLI input/output helper utilities.

* **Features:** Standard output formatting with predictable line breaks.
* **Import path:** `github.com/AeonDigital/Go-Core/xio`

### xjson

JSON manipulation and debugging helpers.

* **Features:** Debug JSON serialization (`Dump`), human-readable output.
* **Import path:** `github.com/AeonDigital/Go-Core/xjson`

### xlog

Logging framework with CLI and file registry support.

* **Features:** Structured slog handler, ANSI colors, log file registry, timezone support.
* **Import path:** `github.com/AeonDigital/Go-Core/xlog`

### xmocks

Mock generation utilities for testing.

* **Features:** Test data generation, mock scaffolding tools.
* **Import path:** `github.com/AeonDigital/Go-Core/xmocks`

### xreflect

Reflection helpers for runtime type operations.

* **Features:** Generic instance allocation, dynamic type initialization.
* **Import path:** `github.com/AeonDigital/Go-Core/xreflect`

### xtime

DateTime formatting and utilities.

* **Features:** Universal layout translation to Go reference format, common datetime helpers.
* **Import path:** `github.com/AeonDigital/Go-Core/xtime`

### xunits

Custom types for units of magnitude.

* **Features:** Human-readable bytes (`Bytes`), extended durations (`TimeDuration`), JSON support.
* **Import path:** `github.com/AeonDigital/Go-Core/xunits`

### tools

Core debugging and formatting utilities.

* **Features:** JSON dump, generic datetime layout translation, error formatting.
* **Import path:** `github.com/AeonDigital/Go-Core/tools`


&nbsp;
&nbsp;


________________________________________________________________________________

## Install

Install any individual package using `go get`:

```shell
  go get github.com/AeonDigital/Go-Core/xconfig@latest
  go get github.com/AeonDigital/Go-Core/xerrors@latest
  go get github.com/AeonDigital/Go-Core/xfs@latest
```

Or install the entire monorepo:

```shell
  go get github.com/AeonDigital/Go-Core/...@latest
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Test Suite

Each package contains its own isolated test suite. To run tests across all modules locally:

```shell
go test -v -cover ./...
```

Or test individual packages:

```shell
cd xconfig && go test -v -cover ./...
cd ../xfs && go test -v -cover ./...
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Licence

This project is offered under the [MIT license](LICENCE.md).