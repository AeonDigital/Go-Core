xjson
================================

> JSON helper utilities for Go, intended as a small aggregator for CLI and debugging.

&nbsp;

This package aims to collect reusable JSON manipulation helpers in one place.
At the moment it only exposes a debug-oriented `Dump` function, but it is
intended to grow into a broader collection of JSON utilities.


________________________________________________________________________________

## Purpose

`xjson` is designed to centralize JSON helper functions for command-line and
utility code. The package currently focuses on debug serialization and
readable JSON output, with the future goal of aggregating additional JSON
manipulation utilities.


________________________________________________________________________________

## Installation

Use `go get` to add the package:

```shell
  go get github.com/AeonDigital/Go-Core/xjson@latest
```

Import it in your code:

```go
import "github.com/AeonDigital/Go-Core/xjson"
```


________________________________________________________________________________

## Basic usage

Serialize a value to an indented JSON string for debugging:

```go
output := xjson.Dump(myValue, "")
fmt.Println(output)
```

If the indentation string is empty, `Dump` defaults to two spaces.

```go
output := xjson.Dump(myValue, "    ")
```

If marshaling fails, `Dump` returns a string containing the error message,
which makes it safe for debug prints.


________________________________________________________________________________

## Supported APIs

Currently the package exposes:

* `xjson.Dump(v any, i string) string`

This function is intended to provide a single consistent JSON dump helper
while the package expands into a broader JSON utility collection.


________________________________________________________________________________

## External dependencies

`xjson` depends only on the Go standard library:

* `encoding/json`

No third-party packages are required.
