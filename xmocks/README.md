# xmocks

`xmocks` is a powerful, generic CLI generator designed to parse Go interfaces via AST and automatically emit type-safe mock implementations. It is built to seamlessly support structural dependency decoupling throughout the `Go-Core` ecosystem.

## Features

- **Dynamic Package Resolution**: No longer hardcoded. The generator dynamically identifies the source package name and the destination folder name to build clean package headers (e.g., `package mocks`).
- **Automatic Package Prefixing**: Types and structures declared locally in the source file automatically receive their parent package identifier (e.g., `xbridge.IFileBridge`), preventing compilation errors.
- **Variadic Parameter Handling**: Full native support for variadic parameters (`...string`). The generated code correctly propagates slice expansion (`args...`) into underlying mock function fields.
- **Collision-Resistant Receivers**: All mock methods utilize `oMock` as their pointer receiver name to avoid nomenclature conflicts with method parameters (such as `m *runtime.MemStats`).
- **Strict Testing Contract**: Implements `TestCase<Alias>`, `OnCall`, and `SetReturn` mechanics alongside runtime call protection (`panicIfNotConfigured`).

---

## Architectural Conventions (The Bridge Pattern)

To achieve absolute cross-platform isolation and guarantee 100% testability across core infrastructure packages (like `os`, `io`, `runtime`, and `path/filepath`), the `Go-Core` project strictly enforces the **Bridge Pattern** guidelines below:

### 1. Unified Naming Design
- **Interfaces**: Must be public and prefixed with a capital `I` (e.g., `IOSBridge`, `IIOBridge`, `IFileBridge`).
- **Concrete Structs**: Must be unexported/private to ensure strict encapsulation and prevent accidental manual instantiations. They use a lowercase `s` prefix (e.g., `sOSBridge`, `sIOBridge`, `sFileBridge`).
- **Global Pointers**: A single public global variable named after the package must point to the interface type, acting as the standard production gateway (e.g., `var OSBridge IOSBridge = sOSBridge{}`).

### 2. Visibilities and Boundaries
- **Production Boundary**: Code outside the `xbridge` package only communicates with the global variables (`xbridge.OSBridge`). Struct implementations are fully hidden.
- **Compile-Time Safeguards**: Every bridge implementation file must include a compile-time check in English right below its private struct declaration:
  ```go
  // Ensure at compile time that the private struct implements the public interface.
  var _ IOSBridge = sOSBridge{}
  ```

### 3. Dynamic Life Cycles (No Globals for Products)
- Heavy singletons like the operating system layer (`OSBridge`) or platform diagnostic hooks (`RuntimeBridge`) remain global.
- Volatile, short-lived resources like opened files (`IFileBridge`) or metadata snapshots (`IFileInfoBridge`) **do not possess global entry points**. They are dynamically manufactured by the main services and mocked via return-injection inside tests.

---

## Usage

```bash
go run ./xmocks --file=<path/to/source.go> --interface=<InterfaceName> [--alias=<Alias>] [--output=<dir|path>]
```

### Required flags

- `--file` — Path to the Go source file containing the target interface.
- `--interface` — Name of the interface to mock (e.g., `IOSBridge`).

### Optional flags

- `--alias` — Alias used in generated mock type names; defaults to the interface name.
- `--output` — Output directory or direct file path. The generated file will dynamically match the destination package directory name.

---

## Example (Batch Generation)

To generate the standard `xbridge` mocks architecture inside your project, execute the following commands in sequence:

```bash
# 1. Operating System & File descriptors layers
go run ./xmocks --file=xbridge/os.go --interface=IOSBridge --alias=OS --output=xbridge/mocks/
go run ./xmocks --file=xbridge/file.go --interface=IFileBridge --alias=File --output=xbridge/mocks/
go run ./xmocks --file=xbridge/fileinfo.go --interface=IFileInfoBridge --alias=FileInfo --output=xbridge/mocks/

# 2. Streams & Paths utilities
go run ./xmocks --file=xbridge/ioe.go --interface=IIOBridge --alias=IO --output=xbridge/mocks/
go run ./xmocks --file=xbridge/filepath.go --interface=IFilepathBridge --alias=Filepath --output=xbridge/mocks/

# 3. Platform Diagnostics
go run ./xmocks --file=xbridge/runtime.go --interface=IRuntimeBridge --alias=Runtime --output=xbridge/mocks/
```

This creates a clean, independent subpackage under `xbridge/mocks/` with `package mocks` dynamically stamped at the top of each file, ready to be safely imported by external applications.

---

## Testing

The CLI tool preserves a rigid testing coverage standard. To execute the internal generator test battery and review statements coverage:

```bash
go test -cover -v ./xmocks
```
