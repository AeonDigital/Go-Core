xio
================================

> CLI helper utilities for IO operations in command-line tools.

&nbsp;

This package is intended as a small aggregator of helper functions for CLI use.
At the moment it only provides basic standard output formatting helpers, but its
purpose is to centralize reusable IO behavior for command-line applications.


________________________________________________________________________________

## Purpose

`xio` is designed to collect simple CLI-oriented IO utilities in one place.
It currently focuses on printing consistent output to standard output with a
minimal API, while leaving space for future CLI helpers.

The package is meant for command-line tools that need a small, predictable
wrapper around common output operations.


________________________________________________________________________________

## Installation

Use `go get` to install the package:

```shell
  go get github.com/AeonDigital/Go-Core/xio@latest
```

Import it in your code:

```go
import "github.com/AeonDigital/Go-Core/xio"
```


________________________________________________________________________________

## Basic usage

Print a simple line to standard output:

```go
xio.Print("hello world")
```

Use formatted output with arguments:

```go
xio.Print("user %s has id %d", "john", 42)
```

The function behaves like `fmt.Println` when no format arguments are provided,
and like `fmt.Printf` followed by a newline when arguments are passed.


________________________________________________________________________________

## Supported APIs

Currently the package exposes:

* `xio.Print(message string, args ...any)`

This function is intended to provide predictable CLI output behavior while
keeping the API surface small.


________________________________________________________________________________

## External dependencies

`xio` depends only on the Go standard library:

* `fmt`
* `io`
* `os`
* `strings`

No third-party external packages are required.
