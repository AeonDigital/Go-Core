xerrors
================================

> Structured error classification and thread-safe logging utilities for corporate Go applications.

&nbsp;

This package provides a rich, dual-layer error abstraction designed to isolate user-side validation failures from unexpected infrastructure breakdowns. It enforces visual predictability across logs, microservices, and tracking channels, boasting a completely lock-free concurrent reading engine with zero third-party dependencies.

&nbsp;
&nbsp;

________________________________________________________________________________

## Purpose

`xerrors` standardizes error tracking, observability, and debugging ergonomics across repositories by introducing two specialized error families powered by a unified internal layout engine:

*   **Error 400 Family (`IError400`)** — Lightweight, high-performance validation exceptions caused by invalid user inputs. Bypasses reflection and stack inspection to maintain high throughput in API entry points.
*   **Error 500 Family (`IError500`)** — Rich, technical failure instances capturing root causes, thread-safe lazy-loaded runtime component/stack tracing, and unstructured debugging data payloads.

Furthermore, both families feature an advanced, immutable **Fluent API (`WithArgs`)** that maps context variables sequentially into metadata registries, automatically safeguarding missing arguments with the mathematical empty set marker (`ø`) via high-performance string builders.

&nbsp;
&nbsp;

________________________________________________________________________________

## Installation

Use `go get` pointing to the package module:

```shell
  go get github.com/AeonDigital/Go-Core/xerrors@latest
```

In your code, import the package:

```go
import (
    "errors"

    "github.com/AeonDigital/Go-Core/xerrors"
)
```

&nbsp;
&nbsp;

________________________________________________________________________________

## Global Configuration (Thread-Safe Runtime Tuning)

The package includes atomic-driven global configuration modifiers to control error rendering verbosity safely at runtime across multiple goroutines without performance bottlenecks.

```go
// Enable technical layout extensions (Component tracking and root-cause dumps)
xerrors.EnableDebugMode()

// Fallback to sanitized, user-friendly messages for end-users
xerrors.DisableDebugMode()

// Check or alternate current state atomically
isTrue := xerrors.GetDebugMode()
xerrors.ToggleDebugMode()
```

&nbsp;
&nbsp;

________________________________________________________________________________

## Basic Usage

### 1. User Validation Failures (Error 400)
The `NewError400` factory uses a flexible, variadic signature (`args ...any`). It automatically determines if you are triggering a core framework token, an extended domain-specific token, or a standard plain formatted text. Contextual telemetry fields are appended cleanly using the fluent `.WithArgs()` API.

```go
// Scenario A: Emitting a core framework error token with dynamic validation arguments
err := xerrors.NewError400(xerrors.XERR_FIELD_REQUIRED).WithArgs("email")

// Scenario B: Emitting a plain formatted error without a registered token (defaults to XERR_NONE)
errPlain := xerrors.NewError400("invalid temporary session token: %s", tokenID)
```

### 2. Unexpected System Failures (Error 500)
Use `NewError500` to instantiate a rich operational error when an infrastructure layer or processing routine breaks down. It natively supports Go's standard library unwrapping semantics (`errors.Is` / `errors.As`).

```go
dbErr := errors.New("connection timeout downstream")

richErr := xerrors.NewError500(
    xerrors.XERR_PKGCTX,            // Context tracking boundary (errCTX)
    xerrors.XERR_UNKNOWN,           // Internal ErrorCode constant
    dbErr,                          // Raw low-level root cause error
    "database repository failure",  // Human-readable message summary
    `{"retry_count": 3}`,           // Raw debugging context chunk or payload
).WithArgs("user_id_123")           // Injects dynamic tracking tags concurrently
```

To adjust caller frame detection when wrapping the error within helper blocks or custom factories, utilize the fluent `WithCallerSkip` API to securely clone and shift the stack inspection depth:

```go
return err.WithCallerSkip(1)
```

### 3. Quick Debug Printing
For quick debugging or direct CLI application tracking, the package provides a zero-allocation utility to stream error summaries directly to `os.Stderr`. It safely ignores `nil` targets.

```go
xerrors.Print(err)
```

&nbsp;
&nbsp;

________________________________________________________________________________

## Domain Error Expansion (Extending xerrors)

To prevent collision anomalies across independent application packages (e.g., two domains using `E1001`), `xerrors` introduces a **namespaced plugin registration architecture** via `RegisterDomainErrors` powered by an ultra-fast `sync.Map` engine.

### Rules for Domain Expansion:
1.  **Code Syntax Contract**: Custom codes must follow the `E[Family][Sequence]` convention (e.g., `E1001`).
2.  **Isolated Namespace**: Every domain package must declare a distinct package context constant (`XERR_PKGCTX`).
3.  **Bootstrap Hook**: Registrations must be executed during the package `init()` phase.

#### Step 1: Define and Register your Custom Domain Codes
Inside your domain package (e.g., `pkg/orders`), declare your context, codes, and metadata:

```go
package orders

import "github.com/AeonDigital/Go-Core/xerrors"

const (
    // Unique package isolation boundary namespace
    XERR_PKGCTX xerrors.ErrorCode = "ERR_ORDER"

    // Custom domain code: Family 1 (State validation), Sequence 1
    XERR_ORDER_EXPIRED xerrors.ErrorCode = "E1001"
)

func init() {
    // Inject configurations into the centralized core registry safely.
    // Thread-safe lock-free reads are guaranteed from this point onward.
    xerrors.RegisterDomainErrors(
        XERR_PKGCTX,
        map[xerrors.ErrorCode]xerrors.MetaMessage{
            XERR_ORDER_EXPIRED: xerrors.NewMetaMessage(
                "checkout session execution window expired",
                "", 
                []string{"FIELD", "VALUE"}, // Triggers the sequential layout tagging engine
            ),
        },
    )
}
```

#### Step 2: Triggering the Custom Namespaced Error
Thanks to the polymorphic signature of `NewError400` and the fluent API, developers pass the package context alongside the custom token, chaining the target metrics cleanly without manual slice indexing:

```go
package orders

import "github.com/AeonDigital/Go-Core/xerrors"

func ValidateSession(sessionID string, expiredAt int64) error {
    if isExpired {
        // Automatically maps to namespace "ERR_ORDER:E1001" behind the scenes.
        // Arguments sequence aligned with extraTags: FIELD (sessionID) then VALUE (expiredAt)
        return xerrors.NewError400(XERR_PKGCTX, XERR_ORDER_EXPIRED).WithArgs(sessionID, expiredAt)
    }
    return nil
}
```

&nbsp;
&nbsp;

________________________________________________________________________________

## Supported APIs

The package exposes the following thread-safe core contracts:

*   `xerrors.IError400` — Interface contract for client-side and validation failures.
*   `xerrors.IError500` — Interface contract for system-side and infrastructure breakdowns.
*   `xerrors.NewError400(args ...any)` — Polymorphic lightweight validation factory.
*   `xerrors.NewError500(ctx, code, err, msg, info)` — Component-tracing operational error factory.
*   `xerrors.RegisterDomainErrors(pkgCtx, registry)` — Safe concurrent domain registration hook.
*   `xerrors.GetDebugMode() / EnableDebugMode() / DisableDebugMode() / ToggleDebugMode()` — Atomic runtime tuning flags.
*   `xerrors.Print(err)` — Swift zero-allocation writer to `os.Stderr`.

&nbsp;
&nbsp;

________________________________________________________________________________

## External Dependencies

`xerrors` strictly depends only on the Go standard library:

*   `fmt`
*   `os`
*   `runtime`
*   `strings`
*   `sync`
*   `sync/atomic`

No third-party modules or external frameworks are required.
