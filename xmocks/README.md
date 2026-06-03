# xmocks

`xmocks` is a small CLI generator for creating package mocks in the style of `internal/pkgxmock/xfs.go`.

## Purpose

This tool parses a Go source file, extracts a named interface, and emits a mocked implementation that follows the existing `xfs` mock conventions:

- `TestCase<Alias>`
- `NewMock<Alias>`
- `Mock<Alias>`
- `mockOnCall<Alias>`
- `mockSetReturn<Alias>`

The generated file is created under `internal/pkgxmock/<alias>.go` by default, or at a path passed by the developer.

## Usage

```bash
go run ./xmocks --file=<path/to/source.go> --interface=<InterfaceName> [--alias=<Alias>] [--output=<dir|path>]
```

### Required flags

- `--file` — path to the Go source file containing the interface
- `--interface` — name of the interface to mock

### Optional flags

- `--alias` — alias used in generated mock type names; defaults to the interface name
- `--output` — output directory or file path for generated mock; defaults to `internal/pkgxmock/<alias>.go`

## Example

```bash
go run ./xmocks --file=xfs/pkgbridge.go --interface=IPkgBridgeXFS --alias=XFS
```

This generates:

```text
internal/pkgxmock/xfs.go
```

or, if `--output` is provided with a directory:

```bash
go run ./xmocks --file=xfs/pkgbridge.go --interface=IPkgBridgeXFS --alias=XFS --output=tmp/test
```

This generates:

```text
tmp/test/xfs.go
```

## Testing

To run the package tests:

```bash
go test ./xmocks
```
