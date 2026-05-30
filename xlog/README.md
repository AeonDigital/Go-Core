xlog
================================

> CLI and registry logging utilities with structured output and optional colorization.

&nbsp;

This package provides a `slog.Handler` implementation aimed at command-line tools.
It supports CLI rendering with optional ANSI colors and an optional file-based
registry log destination for persisted application logs.


________________________________________________________________________________

## Purpose

`xlog` centralizes logging behavior for CLI applications.
It is built around a reusable `LogHandler` that can:

* emit formatted CLI logs,
* apply ANSI colors when the terminal supports them,
* write log records to a file registry,
* manage time format and timezone configuration,
* keep testable I/O boundaries via a small internal bridge.

The package is designed to provide predictable logger handling for command-line
and service workflows.


________________________________________________________________________________

## Installation

Use `go get` to add the package:

```shell
  go get github.com/AeonDigital/Go-Core/xlog@latest
```

Import it in your code:

```go
import (
    "log/slog"

    "github.com/AeonDigital/Go-Core/xlog"
)
```


________________________________________________________________________________

## Basic usage

Create a handler and configure it for CLI and/or file registry logging:

```go
handler := &xlog.LogHandler{
    LogCLI:       true,
    LogCLILevel:  xlog.LevelInfo,
    LogCLIColors: true,
    LogRegistry:  true,
}

err := handler.CheckConfiguration("myapp", "app.log")
if err != nil {
    panic(err)
}

logger := slog.New(handler)
logger.Info("starting application")
```

The handler supports:

* CLI output with time, level and message formatting,
* optional ANSI colors when terminal support is detected,
* optional log file registry writing,
* automatic fallback defaults for time format and file configuration.


________________________________________________________________________________

## Supported APIs

The package exposes the main logging interface and helpers:

* `xlog.LogHandler`
* `(*xlog.LogHandler).CheckConfiguration(appName, logFileName string) error`
* `(*xlog.LogHandler).Enabled(context.Context, slog.Level) bool`
* `(*xlog.LogHandler).Handle(context.Context, slog.Record) error`
* `(*xlog.LogHandler).WithAttrs([]slog.Attr) slog.Handler`
* `(*xlog.LogHandler).WithGroup(name string) slog.Handler`

It also defines severity levels and a default CLI palette:

* `xlog.LevelAll`
* `xlog.LevelNone`
* `xlog.LevelInfo`
* `xlog.LevelWarn`
* `xlog.LevelError`


________________________________________________________________________________

## External dependencies

`xlog` uses the Go standard library and one external Go module:

* `golang.org/x/term`

It also depends on internal repository packages:

* `github.com/AeonDigital/Go-Core/xfs`
* `github.com/AeonDigital/Go-Core/xunits`
* `github.com/AeonDigital/Go-Core/tools`

No additional third-party logging libraries are required.
