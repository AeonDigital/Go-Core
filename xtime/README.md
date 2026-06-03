xtime
================================

> Date/time helpers and small utilities for formatting and layout translation.

&nbsp;

`xtime` aggregates convenience functions to simplify common date/time operations
and formatting tasks across CLI tools and services. The package intends to
collect a small set of helpers that reduce boilerplate when working with Go's
`time` package and custom layout conventions.


________________________________________________________________________________

## Purpose

The goal of `xtime` is to provide lightweight helpers for converting and
managing datetime layout strings, parsing convenience, and other small tools
that are commonly useful in CLI utilities and application scaffolding.

At the moment the package exposes a layout translation helper that maps a
universal format string (e.g. `YYYY-MM-DD HH:mm:ss.SSS a`) into Go's native
reference time layout.


________________________________________________________________________________

## Installation

Use `go get` to add the package to your module:

```shell
  go get github.com/AeonDigital/Go-Core/xtime@latest
```

Import it in your code:

```go
import "github.com/AeonDigital/Go-Core/xtime"
```


________________________________________________________________________________

## Basic usage

Translate a universal datetime layout into Go's layout string:

```go
universal := "YYYY-MM-DD HH:mm:ss.SSS a"
goLayout := xtime.FormatGenericDateTimeToGolangLayout(universal)
// goLayout == "2006-01-02 15:04:05.000 pm"

now := time.Now().In(location)
formatted := now.Format(goLayout)
```

This helper is useful when accepting user-facing formatting patterns or
supporting configuration that uses common, readable datetime tokens.


________________________________________________________________________________

## Supported APIs

Currently the package exposes:

* `xtime.FormatGenericDateTimeToGolangLayout(universalLayout string) string`

This function maps common tokens such as `YYYY`, `MM`, `DD`, `HH`, `mm`, `ss`,
`SSS`, `a/A`, and timezone tokens (`Z`, `ZZ`) into Go's reference time layout
tokens.


________________________________________________________________________________

## External dependencies

`xtime` depends only on the Go standard library (`strings`). No third-party
packages are required.
