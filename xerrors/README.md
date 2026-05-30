xerrors
================================

> Structured error classification and extensible logging utilities for corporate Go applications.

&nbsp;

This package provides a rich, dual-layer error abstraction designed to isolate user-side validation failures from unexpected infrastructure breakdowns. It enforces visual predictability across logs, microservices, and tracking channels without third-party dependencies.


&nbsp;
&nbsp;


________________________________________________________________________________

## Purpose

`xerrors` standardizes error tracking, observability, and debugging ergonomics across repositories by introducing two specialized error families:

*   **Error 400 Family (`IError400`)** — Lightweight, high-performance validation exceptions caused by invalid user inputs. Bypasses reflection and stack inspection to maintain high throughput.
*   **Error 500 Family (`IError500`)** — Rich, technical failure instances capturing root causes, lazy-loaded runtime component/stack tracing, and debugging data payloads.

Furthermore, the package features an advanced polymorphic parsing engine (`buildMask`) that dynamically normalizes visual log layout blocks, automatically safe-guarding missing arguments with the mathematical empty set marker (`ø`).


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
    "context"
    "errors"
    "log/slog"

    "github.com/AeonDigital/Go-Core/xerrors"
)
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Basic Usage

### 1. User Validation Failures (Error 400)
The `NewError400` factory uses a flexible, variadic signature (`args ...any`). It automatically determines if you are triggering a core framework token, an extended domain-specific token, or a standard string format.

```go
// Scenario A: Emitting a core framework error token
err := xerrors.NewError400(xerrors.XERR_FIELD_REQUIRED, "USER-FLOW", "", "email")

// Scenario B: Emitting a plain formatted error without a registered token (defaults to XERR_NONE)
errPlain := xerrors.NewError400("invalid temporary session token: %s", tokenID)
```

### 2. Unexpected System Failures (Error 500)
Use `NewError500` to instantiate a rich error when an infrastructure or unexpected processing routine breaks down. It natively supports Go's modern unwrapping semantics.

```go
dbErr := errors.New("connection timeout downstream")

richErr := xerrors.NewError500(
    xerrors.XERR_PKGCTX,            // Context tracking boundary (errCTX)
    xerrors.XERR_UNKNOWN,           // Internal ErrorCode constant
    dbErr,                          // Raw low-level root cause error
    "database repository failure",  // Human-readable message summary
    `{"retry_count": 3}`,           // Raw debugging context chunk or payload
)
```

To adjust caller frame detection when wrapping the error within helper blocks or custom factories, utilize the fluent `WithCallerSkip` API to securely clone and shift the stack inspection depth:

```go
return err.WithCallerSkip(1)
```


&nbsp;
&nbsp;


________________________________________________________________________________

## Domain Error Expansion (Extending xerrors)

To prevent collision anomalies across independent application packages (e.g., two domains using `E1001`), `xerrors` introduces a **namespaced plugin registration architecture** via `RegisterDomainErrors`. 

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
    // Inject configurations into the centralized core registry safely
    xerrors.RegisterDomainErrors(
        XERR_PKGCTX,
        map[xerrors.ErrorCode]xerrors.MetaMessage{
            XERR_ORDER_EXPIRED: {
                message:   "checkout session execution window expired",
                extraTags: []string{"FIELD", "VALUE"}, // Triggers the visual layout engine
            },
        },
    )
}
```

#### Step 2: Triggering the Custom Namespaced Error
Thanks to the polymorphic signature of `NewError400`, developers can seamlessly pass the package context alongside the custom token without manual string typecasting:

```go
package orders

import "github.com/AeonDigital/Go-Core/xerrors"

func ValidateSession(ctx string, expiredAt int64) error {
    if isExpired {
        // Automatically maps to namespace "ERR_ORDER:E1001" behind the scenes.
        // Arguments sequence aligned with extraTags: FIELD ("session_id") then VALUE (expiredAt)
        return xerrors.NewError400(
            XERR_PKGCTX, 
            XERR_ORDER_EXPIRED, 
            ctx, 
            "", 
            "session_id", 
            expiredAt,
        )
    }
    return nil
}

```


&nbsp;
&nbsp;


________________________________________________________________________________

## Supported APIs

The package exposes the following core contracts:

*   `xerrors.IError400` — Interface contract for client-side and validation failures.
*   `xerrors.IError500` — Interface contract for system-side and infrastructure breakdowns.
*   `xerrors.NewError400(args ...any)` — Polymorphic lightweight validation factory.
*   `xerrors.NewError500(ctx, code, err, msg, data)` — Component-tracing structural error factory.
*   `xerrors.RegisterDomainErrors(pkgCtx, registry)` — Domain registration hook for extension mapping.
*   `xerrors.Log(ctx, err, attrs, options...)` — Structured `slog` publishing channel.


&nbsp;
&nbsp;


________________________________________________________________________________

## External Dependencies

`xerrors` strictly depends only on the Go standard library:

*   `context`
*   `errors`
*   `fmt`
*   `log/slog`
*   `runtime`
*   `strings`

No third-party modules or frameworks are required.
