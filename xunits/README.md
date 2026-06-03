xunits
================================

> Custom types for managing units of magnitude: bytes, durations, and other common measurements.

&nbsp;

`xunits` aggregates reusable types that either create new magnitude representations
(like `Bytes` for human-readable data sizes) or extend existing Go types
(like `TimeDuration` which wraps `time.Duration` with additional features).

The package is designed to provide convenient, JSON-friendly types for CLI tools,
configuration management, and service applications.


________________________________________________________________________________

## Purpose

`xunits` centralizes custom types for common units of measurement, with the goal of:

* creating new types that don't exist in the standard library (`Bytes`),
* extending existing types with additional capabilities (`TimeDuration`),
* providing JSON marshaling/unmarshaling support for these types,
* making it easy to work with human-readable representations.

This is particularly useful when accepting configuration, command-line flags, or
API inputs that need to support friendly formats like `"1.5GB"` or `"7d"`.


________________________________________________________________________________

## Installation

Add the package to your module:

```shell
  go get github.com/AeonDigital/Go-Core/xunits@latest
```

Import it in your code:

```go
import "github.com/AeonDigital/Go-Core/xunits"
```


________________________________________________________________________________

## Basic usage

### Bytes

Work with data sizes using human-readable notation:

```go
var size xunits.Bytes
json.Unmarshal([]byte(`"1.5GB"`), &size)
fmt.Println(size)  // Output: "1.50GB"
fmt.Println(size.String())  // Output: "1.50GB"

// Access constants
allocSize := 10 * xunits.MB
```

The `Bytes` type supports JSON input like `"10mb"`, `"500K"`, `"1.5GB"`, and formats
output with two decimal places for readability.

### TimeDuration

Parse and work with durations including day notation:

```go
var duration xunits.TimeDuration
json.Unmarshal([]byte(`"7d"`), &duration)
fmt.Println(duration.String())  // Output: "7d" (since it's exactly 7 days)

var mixed xunits.TimeDuration
json.Unmarshal([]byte(`"1h30m"`), &mixed)
fmt.Println(mixed.String())  // Output: "1h30m0s" (not a day multiple)
```

`TimeDuration` extends the standard `time.Duration` with:

* JSON marshaling support,
* parsing of the custom `"d"` suffix for days,
* automatic formatting back to `"Xd"` notation for day multiples.


________________________________________________________________________________

## Supported APIs

The package exposes two main types:

### Bytes

* `xunits.Bytes` — unsigned 64-bit integer representing data size.
* `xunits.Bytes.UnmarshalJSON([]byte) error` — parse JSON strings like `"10mb"`.
* `xunits.Bytes.String() string` — format to human-readable representation.
* Constants: `xunits.B`, `xunits.KB`, `xunits.MB`, `xunits.GB`, `xunits.TB`.

### TimeDuration

* `xunits.TimeDuration` — wraps `time.Duration` with custom marshaling.
* `xunits.TimeDuration.MarshalJSON() ([]byte, error)` — serialize to JSON string.
* `xunits.TimeDuration.UnmarshalJSON([]byte) error` — parse JSON including `"Xd"` syntax.
* `xunits.TimeDuration.String() string` — format with day notation when applicable.


________________________________________________________________________________

## External dependencies

`xunits` depends only on the Go standard library:

* `encoding/json`
* `fmt`
* `strconv`
* `strings`
* `time`
* `unicode`

It also depends on the internal repository package:

* `github.com/AeonDigital/Go-Core/xerrors`

No third-party measurement or time libraries are required.
