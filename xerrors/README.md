xerrors
================================

> Structured error classification and logging utilities for corporate Go applications.

&nbsp;

This package provides a rich error abstraction layer for Go.
It supports domain-specific error codes, component tracing, structured debug output,
and integration with Go's `slog` logging framework.


________________________________________________________________________________

## Purpose

`xerrors` is designed to standardize error handling and observability across the repository.
It provides:

* `xerrors.NewError` — create a rich `DetailedError` with domain code, message, data, and component tracing.
* `xerrors.NewErr` — factory for formatted or corporate-style errors using the `XERR_` token family.
* `xerrors.Log` — publish structured error logs via `slog`.
* `xerrors.WithLogCallerSkip` — adjust caller frame tracing for wrapped logging helpers.

The package improves troubleshooting by keeping error metadata structured,
traceable, and consistent across logging channels.


________________________________________________________________________________

## Installation

Use `go get` pointing to the package module:

```shell
  go get github.com/AeonDigital/Go-Core/xerrors@latest
```

In your code, import the package:

```go
import (
    "context"
    "errors"
    "log/slog"

    "github.com/AeonDigital/Go-Core/xerrors"
)
```


________________________________________________________________________________

## Basic usage

Create a detailed error instance and log it with structured attributes:

```go
err := errors.New("database connection failed")
richErr := xerrors.NewError(
    xerrors.ErrUnknown,
    err,
    "database initialization error",
    `{"retries":3}`,
)

xerrors.Log(
    context.Background(),
    richErr,
    []slog.Attr{
        slog.String("service", "auth"),
    },
)
```

The `NewErr` helper also supports formatted corporate tokens:

```go
err := xerrors.NewErr(
    xerrors.XERR_NOT_FOUND,
    "DB",
    "resource not found",
    "database",
    "users",
    errors.New("missing row"),
)
```

To adjust caller frame detection when `Log` is wrapped by a helper:

```go
xerrors.Log(ctx, err, nil, xerrors.WithLogCallerSkip(1))
```


________________________________________________________________________________

## Supported APIs

The package exposes:

* `xerrors.DetailedError` — interface for rich error payloads.
* `xerrors.NewError` — build structured errors with component and debug payloads.
* `xerrors.NewErr` — format plain or token-driven runtime errors.
* `xerrors.Log` — emit structured `slog` records.
* `xerrors.WithLogCallerSkip` — configure call-site tracing.

It also defines a set of error classification constants such as:

* `xerrors.ErrUnknown`
* `xerrors.XERR_FIELD_REQUIRED`
* `xerrors.XERR_NOT_FOUND`
* `xerrors.XERR_INVALID_FORMAT`
* `xerrors.XERR_MUTUAL_EXCLUSIVITY_VIOLATION`


________________________________________________________________________________

## External dependencies

`xerrors` depends only on the Go standard library:

* `context`
* `errors`
* `fmt`
* `log/slog`
* `runtime`
* `strings`

No third-party packages are required by this package.
