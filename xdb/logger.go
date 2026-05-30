package xdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/AeonDigital/Go-Core/xerrors"
)

// retrieveDbType inspects the underlying driver of an active sql.DB instance safely to identify the target dialect.
func retrieveDbType(db *sql.DB) string {
	resource := "DB"

	if db != nil {
		if drv := db.Driver(); drv != nil {
			drvType := strings.ToLower(reflect.TypeOf(drv).String())
			switch {
			case strings.Contains(drvType, "sqlite"):
				resource = "sqlite"
			case strings.Contains(drvType, "mysql"):
				resource = "mysql"
			case strings.Contains(drvType, "postgres") || strings.Contains(drvType, "pq"):
				resource = "postgres"
			}
		}
	}

	return resource
}

// logRepoError normalizes engine metadata and delegates structured repository failure telemetry to the global xlog framework.
func logRepoError(
	ctx context.Context,
	db *sql.DB,
	errorCode xerrors.ErrorCode,
	originalErr error,
	query string,
	args []any,
) {
	resource := retrieveDbType(db)

	attrs := []slog.Attr{
		slog.String("resource", resource),
	}

	if query != "" {
		attrs = append(attrs, slog.String("sql_query", query))
	}

	if len(args) > 0 {
		attrs = append(attrs, slog.Any("args", args))
	}

	attrs = append(attrs, slog.String("error_cod", string(errorCode)))
	attrs = append(attrs, slog.String("error_str", fmt.Sprintf("%s", errorCode)))

	// If the original error implements IError500, adjust its caller skip
	// so the reported component points to the real origin.
	de, ok := originalErr.(xerrors.IError500)
	if ok {
		de = de.WithCallerSkip(1)
		attrs = append(attrs, slog.String("component", de.Component()))
	}

	if originalErr != nil {
		attrs = append(attrs, slog.String("error", originalErr.Error()))
	}

	// Emit structured error via slog so that consumer projects can configure
	// `slog` with an `xlog.LogHandler` to control formatting/outputs.
	slog.LogAttrs(ctx, slog.LevelError, "operation failure detected", attrs...)
}
